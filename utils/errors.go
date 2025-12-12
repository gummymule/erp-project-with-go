package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message, details string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Details: details,
	}
}

// Common errors
var (
	ErrNotFound     = NewAppError(http.StatusNotFound, "Resource not found", "")
	ErrBadRequest   = NewAppError(http.StatusBadRequest, "Invalid request", "")
	ErrUnauthorized = NewAppError(http.StatusUnauthorized, "Unauthorized", "")
	ErrInternal     = NewAppError(http.StatusInternalServerError, "Internal server error", "")
)

// Error response helper
func ErrorResponse(c *gin.Context, err error) {
	if appErr, ok := err.(*AppError); ok {
		c.JSON(appErr.Code, gin.H{
			"error":   appErr.Message,
			"details": appErr.Details,
		})
		return
	}

	// Default error
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal server error",
	})
}

// Success response helper
func SuccessResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    data,
	})
}

// Created response helper
func CreatedResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    data,
	})
}
