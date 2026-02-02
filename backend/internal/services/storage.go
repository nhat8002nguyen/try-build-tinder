package services

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/config"
)

type StorageService struct {
	cfg *config.Config
}

func NewStorageService(cfg *config.Config) *StorageService {
	return &StorageService{cfg: cfg}
}

func (s *StorageService) UploadFile(file *multipart.FileHeader, userID uuid.UUID) (string, error) {
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
