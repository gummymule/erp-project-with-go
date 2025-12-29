package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// response structure for all API responses
type Response struct {
	ResponseCode string      `json:"responseCode"`
	ResponseDesc string      `json:"responseDesc"`
	ResponseData interface{} `json:"responseData,omitempty"`
}

// success codes
const (
	CodeSuccess   = "00"
	CodeCreated   = "01"
	CodeNoContent = "02"
	CodePartial   = "03"
)

// error codes
const (
	CodeBadRequest   = "10"
	CodeUnauthorized = "11"
	CodeForbidden    = "12"
	CodeNotFound     = "13"
	CodeValidation   = "14"
	CodeDuplicate    = "15"
	CodeInternal     = "99"
)

// success response
func SuccessResponse(c *gin.Context, message string, data interface{}) {
	response := Response{
		ResponseCode: CodeSuccess,
		ResponseDesc: message,
		ResponseData: data,
	}
	c.JSON(http.StatusOK, response)
}

func CreatedResponse(c *gin.Context, message string, data interface{}) {
	response := Response{
		ResponseCode: CodeCreated,
		ResponseDesc: message,
		ResponseData: data,
	}
	c.JSON(http.StatusCreated, response)
}

func NoContentResponse(c *gin.Context, message string) {
	response := Response{
		ResponseCode: CodeNoContent,
		ResponseDesc: message,
	}
	c.JSON(http.StatusNoContent, response)
}

// error response
func ErrorResponse(c *gin.Context, code string, message string, details interface{}) {
	response := Response{
		ResponseCode: code,
		ResponseDesc: message,
		ResponseData: details,
	}

	// map error codes to HTTP status codes
	switch code {
	case CodeBadRequest, CodeValidation, CodeDuplicate:
		c.JSON(http.StatusBadRequest, response)
	case CodeUnauthorized:
		c.JSON(http.StatusUnauthorized, response)
	case CodeForbidden:
		c.JSON(http.StatusForbidden, response)
	case CodeNotFound:
		c.JSON(http.StatusNotFound, response)
	default:
		c.JSON(http.StatusInternalServerError, response)
	}
}

// convienience error functions
func BadRequestResponse(c *gin.Context, message string, details interface{}) {
	ErrorResponse(c, CodeBadRequest, message, details)
}

func ValidationErrorResponse(c *gin.Context, message string, details interface{}) {
	ErrorResponse(c, CodeValidation, message, details)
}

func NotFoundResponse(c *gin.Context, message string) {
	ErrorResponse(c, CodeNotFound, message, nil)
}

func DuplicateErrorResponse(c *gin.Context, message string, details interface{}) {
	ErrorResponse(c, CodeDuplicate, message, details)
}

func InternalErrorResponse(c *gin.Context, message string, details interface{}) {
	ErrorResponse(c, CodeInternal, message, details)
}
