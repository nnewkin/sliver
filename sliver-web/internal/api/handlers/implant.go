package handlers

import (
	"context"
	"net/http"

	"github.com/MonkeyCode/sliver-web/internal/db/repository"
	"github.com/MonkeyCode/sliver-web/internal/rpc/services"
	"github.com/gin-gonic/gin"
)

type ImplantHandler struct {
	implantService *services.ImplantService
	auditRepo      *repository.AuditRepository
}

func NewImplantHandler(implantService *services.ImplantService, auditRepo *repository.AuditRepository) *ImplantHandler {
	return &ImplantHandler{
		implantService: implantService,
		auditRepo:      auditRepo,
	}
}

func (h *ImplantHandler) ListProfiles(c *gin.Context) {
	ctx := context.Background()
	profiles, err := h.implantService.ListProfiles(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to list profiles",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": profiles,
	})
}

func (h *ImplantHandler) CreateProfile(c *gin.Context) {
	var req struct {
		Name   string                  `json:"name" binding:"required"`
		Config *services.ImplantConfig `json:"config" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	if err := h.implantService.SaveProfile(ctx, req.Name, req.Config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to create profile",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "create_profile", "profile", req.Name, "", c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "profile created successfully",
	})
}

func (h *ImplantHandler) UpdateProfile(c *gin.Context) {
	profileName := c.Param("id")

	var req struct {
		Config *services.ImplantConfig `json:"config" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	if err := h.implantService.SaveProfile(ctx, profileName, req.Config); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to update profile",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "update_profile", "profile", profileName, "", c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "profile updated successfully",
	})
}

func (h *ImplantHandler) DeleteProfile(c *gin.Context) {
	profileName := c.Param("name")
	ctx := context.Background()

	if err := h.implantService.DeleteProfile(ctx, profileName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to delete profile",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "delete_profile", "profile", profileName, "", c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "profile deleted successfully",
	})
}

func (h *ImplantHandler) GenerateImplant(c *gin.Context) {
	var req *services.ImplantConfig
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	data, err := h.implantService.GenerateImplant(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to generate implant",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "generate_implant", "implant", req.Name, req.Format, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"binary": data,
			"name":   req.Name,
		},
	})
}

func (h *ImplantHandler) DownloadImplant(c *gin.Context) {
	profileName := c.Param("id")
	ctx := context.Background()

	data, err := h.implantService.GenerateFromProfile(ctx, profileName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to generate implant",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "download_implant", "implant", profileName, "", c.ClientIP(), c.GetHeader("User-Agent"))

	c.Header("Content-Disposition", "attachment; filename="+profileName+".exe")
	c.Data(http.StatusOK, "application/octet-stream", data)
}
