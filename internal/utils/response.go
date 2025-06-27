package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// APIResponse defines the standard structure for API responses.
type APIResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// NewSuccessResponse creates a success API response.
func NewSuccessResponse(c *gin.Context, httpStatus int, message string, data interface{}) {
	c.JSON(httpStatus, APIResponse{
		Status:  httpStatus,
		Message: message,
		Data:    data,
	})
}

// NewErrorResponse creates an error API response.
func NewErrorResponse(c *gin.Context, httpStatus int, message string, err error) {
	// For production, you might want to log the 'err' but not expose it directly in the message
	errorMessage := message
	if err != nil {
		errorMessage = fmt.Sprintf("%s: %v", message, err) // Include specific error detail for development
	}
	c.JSON(httpStatus, APIResponse{
		Status:  httpStatus,
		Message: errorMessage,
		Data:    nil,
	})
}
