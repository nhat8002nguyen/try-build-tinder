package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Error:   message,
	})
}

func MessageResponse(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
	})
}

func BadRequest(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusBadRequest, message)
}

func Unauthorized(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusUnauthorized, message)
}

func Forbidden(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusForbidden, message)
}

func NotFound(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusNotFound, message)
}

func InternalError(c *gin.Context, message string) {
	ErrorResponse(c, http.StatusInternalServerError, message)
}
