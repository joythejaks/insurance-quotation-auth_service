package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jordisetiawan/insurance-auth-service/internal/utils"
	"go.uber.org/zap"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			utils.Log.Error("Internal Server Error", zap.Any("errors", c.Errors))
			utils.ErrorResponse(c, http.StatusInternalServerError, "Internal Server Error", nil)
		}
	}
}
