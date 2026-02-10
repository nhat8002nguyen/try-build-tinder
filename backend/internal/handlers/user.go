package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/tinder-clone/backend/internal/apperrors"
	"github.com/tinder-clone/backend/internal/constants"
	"github.com/tinder-clone/backend/internal/middleware"
	"github.com/tinder-clone/backend/internal/models"
	"github.com/tinder-clone/backend/internal/services"
	"github.com/tinder-clone/backend/internal/utils"
)

type UserHandler struct {
	userService    *services.UserService
	storageService *services.StorageService
}

func NewUserHandler(userService *services.UserService, storageService *services.StorageService) *UserHandler {
	return &UserHandler{
		userService:    userService,
		storageService: storageService,
	}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, err := uuid.Parse(idParam)
	if err != nil {
		utils.BadRequest(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetByIDSimple(userID)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	var input services.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	user, err := h.userService.UpdateProfile(userID, input)
	if err != nil {
		utils.InternalError(c, "Failed to update profile")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, user)
}

func (h *UserHandler) UploadPhoto(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	file, err := c.FormFile("photo")
	if err != nil {
		utils.BadRequest(c, "Photo file required")
		return
	}

	if err := h.storageService.ValidateImage(file); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	photoCount, err := h.userService.GetPhotoCount(userID)
	if err != nil {
		utils.InternalError(c, "Failed to check photo count")
		return
	}

	if photoCount >= int64(constants.MaxPhotosPerUser) {
		utils.BadRequest(c, "Maximum of 6 photos allowed")
		return
	}

	photoURL, err := h.storageService.UploadFile(file, userID)
	if err != nil {
		utils.InternalError(c, "Failed to upload photo")
		return
	}

	photo, err := h.userService.AddPhoto(userID, photoURL, int(photoCount))
	if err != nil {
		if errors.Is(err, apperrors.ErrMaxPhotosReached) {
			utils.BadRequest(c, "Maximum of 6 photos allowed")
			return
		}
		utils.InternalError(c, "Failed to save photo")
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, photo)
}

func (h *UserHandler) DeletePhoto(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	photoIDParam := c.Param("photoId")
	photoID, err := uuid.Parse(photoIDParam)
	if err != nil {
		utils.BadRequest(c, "Invalid photo ID")
		return
	}

	if err := h.userService.DeletePhoto(userID, photoID); err != nil {
		utils.InternalError(c, "Failed to delete photo")
		return
	}

	utils.MessageResponse(c, http.StatusOK, "Photo deleted successfully")
}

type UpdateLocationRequest struct {
	Latitude  float64 `json:"latitude" binding:"required"`
	Longitude float64 `json:"longitude" binding:"required"`
}

func (h *UserHandler) UpdateLocation(c *gin.Context) {
	userID := c.MustGet(middleware.UserIDKey).(uuid.UUID)

	var req UpdateLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	if req.Latitude < -90 || req.Latitude > 90 {
		utils.BadRequest(c, "Invalid latitude")
		return
	}
	if req.Longitude < -180 || req.Longitude > 180 {
		utils.BadRequest(c, "Invalid longitude")
		return
	}

	location := models.Location{
		Latitude:  req.Latitude,
		Longitude: req.Longitude,
	}

	if err := h.userService.UpdateLocation(userID, location); err != nil {
		utils.InternalError(c, "Failed to update location")
		return
	}

	utils.MessageResponse(c, http.StatusOK, "Location updated successfully")
}
