package middleware

import (
	"net/http"
	"strings"

	"github.com/MonkeyCode/sliver-web/internal/core"
	"github.com/MonkeyCode/sliver-web/internal/db/repository"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authManager *core.AuthManager
	userRepo    *repository.UserRepository
}

func NewAuthMiddleware(authManager *core.AuthManager, userRepo *repository.UserRepository) *AuthMiddleware {
	return &AuthMiddleware{
		authManager: authManager,
		userRepo:    userRepo,
	}
}

func (m *AuthMiddleware) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    10001,
				"message": "missing authorization header",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    10001,
				"message": "invalid authorization header format",
			})
			return
		}

		tokenString := parts[1]
		claims, err := m.authManager.ValidateToken(tokenString)
		if err != nil {
			code := 10001
			message := "invalid token"
			if err == core.ErrExpiredToken {
				message = "token has expired"
				code = 10002
			}
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    code,
				"message": message,
			})
			return
		}

		user, err := m.userRepo.GetByID(int64(claims.UserID))
		if err != nil || !user.IsActive {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    10001,
				"message": "user not found or inactive",
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("role", claims.Role)
		c.Set("user", user)

		c.Next()
	}
}

func (m *AuthMiddleware) RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get("user")
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    10001,
				"message": "user not authenticated",
			})
			return
		}

		u, ok := user.(*repository.User)
		if !ok {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"code":    10000,
				"message": "internal error",
			})
			return
		}

		hasRole := false
		for _, role := range roles {
			if u.HasRole(role) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"code":    10004,
				"message": "insufficient permissions",
			})
			return
		}

		c.Next()
	}
}
