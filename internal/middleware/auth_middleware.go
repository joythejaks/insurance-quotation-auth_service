package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jordisetiawan/insurance-auth-service/internal/utils"
)

func AuthMiddleware(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Authorization header is required", nil)
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid authorization header format", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString, secret)
		if err != nil {
			utils.ErrorResponse(c, http.StatusUnauthorized, "Invalid or expired token", err.Error())
			c.Abort()
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("permissions", claims.Permissions)
		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, _ := c.Get("role")
		for _, role := range roles {
			if role == userRole {
				c.Next()
				return
			}
		}
		utils.ErrorResponse(c, http.StatusForbidden, "Access denied: insufficient role permissions", nil)
		c.Abort()
	}
}

func PermissionMiddleware(permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		permissions, exists := c.Get("permissions")
		if !exists {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied: no permissions assigned", nil)
			c.Abort()
			return
		}

		userPerms := permissions.([]string)
		found := false
		for _, p := range userPerms {
			if p == permission {
				found = true
				break
			}
		}

		if !found {
			utils.ErrorResponse(c, http.StatusForbidden, "Access denied: missing required permission", permission)
			c.Abort()
			return
		}

		c.Next()
	}
}
