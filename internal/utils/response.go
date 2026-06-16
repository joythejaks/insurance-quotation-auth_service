package utils

import (
	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func SuccessResponse(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	})
}

func ErrorResponse(c *gin.Context, code int, message string, err interface{}) {
	c.JSON(code, APIResponse{
		Status:  "error",
		Message: message,
		Errors:  err,
	})
}
