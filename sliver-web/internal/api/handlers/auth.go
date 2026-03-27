package handlers

import (
	"net/http"

	"github.com/MonkeyCode/sliver-web/internal/core"
	"github.com/MonkeyCode/sliver-web/internal/db/models"
	"github.com/MonkeyCode/sliver-web/internal/db/repository"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authManager *core.AuthManager
	totpManager *core.TOTPManager
	userRepo    *repository.UserRepository
	auditRepo   *repository.AuditRepository
}

func NewAuthHandler(authManager *core.AuthManager, totpManager *core.TOTPManager, userRepo *repository.UserRepository, auditRepo *repository.AuditRepository) *AuthHandler {
	return &AuthHandler{
		authManager: authManager,
		totpManager: totpManager,
		userRepo:    userRepo,
		auditRepo:   auditRepo,
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token     string        `json:"token"`
	ExpiresAt string        `json:"expires_at"`
	User      *UserResponse `json:"user"`
	Need2FA   bool          `json:"need_2fa"`
}

type UserResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type Verify2FARequest struct {
	Token string `json:"token" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	user, err := h.userRepo.GetByUsername(req.Username)
	if err != nil {
		h.auditRepo.Log(0, req.Username, "login", "auth", "", "", c.ClientIP(), c.GetHeader("User-Agent"))
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    10001,
			"message": "invalid credentials",
		})
		return
	}

	if !user.IsActive {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    10001,
			"message": "user account is disabled",
		})
		return
	}

	if !core.CheckPassword(req.Password, user.PasswordHash) {
		h.auditRepo.Log(user.ID, user.Username, "login_failed", "auth", "", "invalid password", c.ClientIP(), c.GetHeader("User-Agent"))
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    10001,
			"message": "invalid credentials",
		})
		return
	}

	if user.TOTPEnabled && user.TOTPSecret != "" {
		h.auditRepo.Log(user.ID, user.Username, "login", "auth", "", "2fa required", c.ClientIP(), c.GetHeader("User-Agent"))
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": LoginResponse{
				Need2FA: true,
				User: &UserResponse{
					ID:       user.ID,
					Username: user.Username,
					Role:     user.Role,
				},
			},
		})
		return
	}

	token, err := h.authManager.GenerateToken(int(user.ID), user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10000,
			"message": "failed to generate token",
		})
		return
	}

	h.auditRepo.Log(user.ID, user.Username, "login", "auth", "", "", c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": LoginResponse{
			Token:     token,
			ExpiresAt: "",
			User: &UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Role:     user.Role,
			},
			Need2FA: false,
		},
	})
}

func (h *AuthHandler) Verify2FA(c *gin.Context) {
	var req Verify2FARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	user, err := h.userRepo.GetByUsername(req.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    10001,
			"message": "user not found",
		})
		return
	}

	if !user.TOTPEnabled || user.TOTPSecret == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "2fa not enabled for user",
		})
		return
	}

	if err := h.totpManager.ValidateCode(user.TOTPSecret, req.Code); err != nil {
		h.auditRepo.Log(user.ID, user.Username, "2fa_failed", "auth", "", "", c.ClientIP(), c.GetHeader("User-Agent"))
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    10003,
			"message": "invalid 2fa code",
		})
		return
	}

	token, err := h.authManager.GenerateToken(int(user.ID), user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10000,
			"message": "failed to generate token",
		})
		return
	}

	h.auditRepo.Log(user.ID, user.Username, "login", "auth", "", "2fa success", c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": LoginResponse{
			Token:     token,
			ExpiresAt: "",
			User: &UserResponse{
				ID:       user.ID,
				Username: user.Username,
				Role:     user.Role,
			},
			Need2FA: false,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	h.auditRepo.Log(userID.(int64), username.(string), "logout", "auth", "", "", c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "logged out successfully",
	})
}

func (h *AuthHandler) GetProfile(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    10001,
			"message": "not authenticated",
		})
		return
	}

	u := user.(*models.User)

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"id":           u.ID,
			"username":     u.Username,
			"role":         u.Role,
			"totp_enabled": u.TOTPEnabled,
			"is_active":    u.IsActive,
			"created_at":   u.CreatedAt,
		},
	})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	user, err := h.userRepo.GetByID(userID.(int64))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10000,
			"message": "user not found",
		})
		return
	}

	if !core.CheckPassword(req.OldPassword, user.PasswordHash) {
		h.auditRepo.Log(userID.(int64), username.(string), "change_password", "auth", "", "invalid old password", c.ClientIP(), c.GetHeader("User-Agent"))
		c.JSON(http.StatusUnauthorized, gin.H{
			"code":    10001,
			"message": "invalid old password",
		})
		return
	}

	if err := h.userRepo.SetPassword(userID.(int64), req.NewPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10000,
			"message": "failed to change password",
		})
		return
	}

	h.auditRepo.Log(userID.(int64), username.(string), "change_password", "auth", "", "", c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "password changed successfully",
	})
}
