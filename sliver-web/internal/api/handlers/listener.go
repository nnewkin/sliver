package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/MonkeyCode/sliver-web/internal/db/repository"
	"github.com/MonkeyCode/sliver-web/internal/rpc/services"
	"github.com/gin-gonic/gin"
)

type ListenerHandler struct {
	listenerService *services.ListenerService
	auditRepo       *repository.AuditRepository
}

func NewListenerHandler(listenerService *services.ListenerService, auditRepo *repository.AuditRepository) *ListenerHandler {
	return &ListenerHandler{
		listenerService: listenerService,
		auditRepo:       auditRepo,
	}
}

func (h *ListenerHandler) ListListeners(c *gin.Context) {
	ctx := context.Background()
	listeners, err := h.listenerService.GetAllListeners(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to get listeners",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": listeners,
	})
}

func (h *ListenerHandler) StartMTLSListener(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		BindAddr string `json:"bind_address" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	listener, err := h.listenerService.StartMTLSListener(ctx, req.Name, req.BindAddr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to start MTLS listener",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "start_listener", "listener", req.Name, req.BindAddr, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": listener,
	})
}

func (h *ListenerHandler) StartDNSListener(c *gin.Context) {
	var req struct {
		Name     string   `json:"name" binding:"required"`
		BindAddr string   `json:"bind_address" binding:"required"`
		Domains  []string `json:"domains"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	listener, err := h.listenerService.StartDNSListener(ctx, req.Name, req.BindAddr, req.Domains)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to start DNS listener",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "start_listener", "listener", req.Name, req.BindAddr, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": listener,
	})
}

func (h *ListenerHandler) StartHTTPListener(c *gin.Context) {
	var req struct {
		Name     string   `json:"name" binding:"required"`
		BindAddr string   `json:"bind_address" binding:"required"`
		Domains  []string `json:"domains"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	listener, err := h.listenerService.StartHTTPListener(ctx, req.Name, req.BindAddr, req.Domains)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to start HTTP listener",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "start_listener", "listener", req.Name, req.BindAddr, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": listener,
	})
}

func (h *ListenerHandler) StartHTTPSListener(c *gin.Context) {
	var req struct {
		Name     string   `json:"name" binding:"required"`
		BindAddr string   `json:"bind_address" binding:"required"`
		Domains  []string `json:"domains"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	listener, err := h.listenerService.StartHTTPSListener(ctx, req.Name, req.BindAddr, req.Domains)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to start HTTPS listener",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "start_listener", "listener", req.Name, req.BindAddr, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": listener,
	})
}

func (h *ListenerHandler) StartWGListener(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		BindAddr string `json:"bind_address" binding:"required"`
		Port     int    `json:"port" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	listener, err := h.listenerService.StartWGListener(ctx, req.Name, req.BindAddr, req.Port)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to start WireGuard listener",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "start_listener", "listener", req.Name, req.BindAddr, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": listener,
	})
}

func (h *ListenerHandler) StopListener(c *gin.Context) {
	listenerIDStr := c.Param("id")
	listenerID, err := strconv.ParseInt(listenerIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid listener id",
		})
		return
	}

	ctx := context.Background()
	if err := h.listenerService.StopListener(ctx, listenerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to stop listener",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "stop_listener", "listener", listenerIDStr, "", c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "listener stopped successfully",
	})
}
