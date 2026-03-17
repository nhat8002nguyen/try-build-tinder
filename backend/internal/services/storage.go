package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/config"
)

var ErrNSFWContent = errors.New("image contains disallowed NSFW content")

type NSFWClient struct {
	httpClient *http.Client
	baseURL    string
}

type debugLogEntry struct {
	SessionID    string         `json:"sessionId,omitempty"`
	RunID        string         `json:"runId,omitempty"`
	HypothesisID string         `json:"hypothesisId,omitempty"`
	Location     string         `json:"location,omitempty"`
	Message      string         `json:"message,omitempty"`
	Data         map[string]any `json:"data,omitempty"`
	Timestamp    int64          `json:"timestamp,omitempty"`
}

// #region agent log
func writeNSFWDebugLog(hypothesisID, location, message string, data map[string]any) {
	entry := debugLogEntry{
		SessionID:    "0f9762",
		RunID:        "initial",
		HypothesisID: hypothesisID,
		Location:     location,
		Message:      message,
		Data:         data,
		Timestamp:    time.Now().UnixMilli(),
	}

	payload, err := json.Marshal(entry)
	if err != nil {
		return
	}

	req, err := http.NewRequest(
		http.MethodPost,
		"http://host.docker.internal:7586/ingest/e28d96d1-7f10-4f83-951e-74e65210567a",
		bytes.NewReader(payload),
	)
	if err != nil {
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Debug-Session-Id", "0f9762")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
}

// #endregion

func NewNSFWClient(baseURL string) *NSFWClient {
	if baseURL == "" {
		return nil
	}

	return &NSFWClient{
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
		baseURL: strings.TrimRight(baseURL, "/"),
	}
}

func (c *NSFWClient) Classify(file *multipart.FileHeader) (bool, float64, error) {
	if c == nil {
		return false, 0, nil
	}

	src, err := file.Open()
	if err != nil {
		return false, 0, fmt.Errorf("failed to open uploaded file for nsfw scan: %w", err)
	}
	defer src.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", file.Filename)
	if err != nil {
		return false, 0, fmt.Errorf("failed to create multipart form file: %w", err)
	}

	if _, err := io.Copy(part, src); err != nil {
		return false, 0, fmt.Errorf("failed to copy file to multipart form: %w", err)
	}

	if err := writer.Close(); err != nil {
		return false, 0, fmt.Errorf("failed to finalize multipart form: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, c.baseURL+"/classify", &buf)
	if err != nil {
		writeNSFWDebugLog(
			"A",
			"storage.go:nsfw-classify",
			"failed_to_create_nsfw_request",
			map[string]any{"error": err.Error()},
		)
		return false, 0, fmt.Errorf("failed to create nsfw request: %w", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := c.httpClient.Do(req)
	if err != nil {
		writeNSFWDebugLog(
			"A",
			"storage.go:nsfw-classify",
			"failed_to_call_nsfw_service",
			map[string]any{"error": err.Error()},
		)
		return false, 0, fmt.Errorf("failed to call nsfw service: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		writeNSFWDebugLog(
			"A",
			"storage.go:nsfw-classify",
			"nsfw_service_unavailable",
			map[string]any{"status_code": resp.StatusCode},
		)
		return false, 0, nil
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024))
		writeNSFWDebugLog(
			"A",
			"storage.go:nsfw-classify",
			"nsfw_service_http_error",
			map[string]any{"status_code": resp.StatusCode, "body": strings.TrimSpace(string(body))},
		)
		return false, 0, fmt.Errorf("nsfw service returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var result struct {
		IsNSFW bool    `json:"is_nsfw"`
		Label  string  `json:"label"`
		Score  float64 `json:"score"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		writeNSFWDebugLog(
			"B",
			"storage.go:nsfw-classify",
			"failed_to_decode_nsfw_response",
			map[string]any{"error": err.Error()},
		)
		return false, 0, fmt.Errorf("failed to decode nsfw response: %w", err)
	}

	writeNSFWDebugLog(
		"B",
		"storage.go:nsfw-classify",
		"nsfw_classification_result",
		map[string]any{
			"is_nsfw": result.IsNSFW,
			"label":   result.Label,
			"score":   result.Score,
			"file":    file.Filename,
		},
	)

	return result.IsNSFW, result.Score, nil
}

type StorageService struct {
	cfg        *config.Config
	nsfwClient *NSFWClient
}

func NewStorageService(cfg *config.Config) *StorageService {
	return &StorageService{
		cfg:        cfg,
		nsfwClient: NewNSFWClient(cfg.NSFWServiceURL),
	}
}

func (s *StorageService) UploadFile(file *multipart.FileHeader, userID uuid.UUID) (string, error) {
	if s.nsfwClient != nil {
		isNSFW, _, err := s.nsfwClient.Classify(file)
		if err != nil {
			return "", fmt.Errorf("failed to classify image: %w", err)
		}
		if isNSFW {
			return "", ErrNSFWContent
		}
	}

	if s.cfg.StorageType == "s3" {
		return s.uploadToS3(file, userID)
	}
	return s.uploadToLocal(file, userID)
}

func (s *StorageService) uploadToLocal(file *multipart.FileHeader, userID uuid.UUID) (string, error) {
	uploadDir := filepath.Join(s.cfg.LocalStorageDir, userID.String())
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create upload directory: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	if !isValidImageExtension(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	filename := fmt.Sprintf("%d_%s%s", time.Now().UnixNano(), uuid.New().String()[:8], ext)
	filePath := filepath.Join(uploadDir, filename)

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create destination file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	relativePath := filepath.Join("/uploads", userID.String(), filename)
	return relativePath, nil
}

func (s *StorageService) uploadToS3(file *multipart.FileHeader, userID uuid.UUID) (string, error) {
	ext := filepath.Ext(file.Filename)
	if !isValidImageExtension(ext) {
		return "", fmt.Errorf("invalid file type: %s", ext)
	}

	key := fmt.Sprintf("users/%s/%d_%s%s",
		userID.String(),
		time.Now().UnixNano(),
		uuid.New().String()[:8],
		ext,
	)

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
		s.cfg.S3Bucket,
		s.cfg.S3Region,
		key,
	), nil
}

func (s *StorageService) DeleteFile(fileURL string) error {
	if s.cfg.StorageType == "s3" {
		return s.deleteFromS3(fileURL)
	}
	return s.deleteFromLocal(fileURL)
}

func (s *StorageService) deleteFromLocal(fileURL string) error {
	relativePath := strings.TrimPrefix(fileURL, "/uploads")
	filePath := filepath.Join(s.cfg.LocalStorageDir, relativePath)

	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (s *StorageService) deleteFromS3(fileURL string) error {
	return nil
}

func isValidImageExtension(ext string) bool {
	validExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".webp": true,
	}
	return validExtensions[strings.ToLower(ext)]
}

func (s *StorageService) ValidateImage(file *multipart.FileHeader) error {
	maxSize := int64(10 * 1024 * 1024)
	if file.Size > maxSize {
		return fmt.Errorf("file size exceeds maximum allowed (10MB)")
	}

	ext := filepath.Ext(file.Filename)
	if !isValidImageExtension(ext) {
		return fmt.Errorf("invalid file type: only jpg, jpeg, png, gif, webp allowed")
	}

	return nil
}
