package handlers

import (
	"context"
	"net/http"
	"strconv"

	"github.com/MonkeyCode/sliver-web/internal/db/repository"
	"github.com/MonkeyCode/sliver-web/internal/rpc/services"
	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService *services.SessionService
	interactiveSvc *services.InteractiveService
	auditRepo      *repository.AuditRepository
}

func NewSessionHandler(sessionService *services.SessionService, interactiveService *services.InteractiveService, auditRepo *repository.AuditRepository) *SessionHandler {
	return &SessionHandler{
		sessionService: sessionService,
		interactiveSvc: interactiveService,
		auditRepo:      auditRepo,
	}
}

func (h *SessionHandler) ListSessions(c *gin.Context) {
	ctx := context.Background()
	sessions, err := h.sessionService.GetAllSessions(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to get sessions",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": sessions,
	})
}

func (h *SessionHandler) GetSession(c *gin.Context) {
	sessionID := c.Param("id")
	ctx := context.Background()

	session, err := h.sessionService.GetSession(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"code":    10005,
			"message": "session not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": session,
	})
}

func (h *SessionHandler) KillSession(c *gin.Context) {
	sessionID := c.Param("id")
	ctx := context.Background()

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")

	if err := h.sessionService.KillSession(ctx, sessionID); err != nil {
		h.auditRepo.Log(userID.(int64), username.(string), "kill_session", "session", sessionID, err.Error(), c.ClientIP(), c.GetHeader("User-Agent"))
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to kill session",
		})
		return
	}

	h.auditRepo.Log(userID.(int64), username.(string), "kill_session", "session", sessionID, "", c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "session killed successfully",
	})
}

func (h *SessionHandler) RenameSession(c *gin.Context) {
	sessionID := c.Param("id")

	var req struct {
		NewName string `json:"new_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	session, err := h.sessionService.RenameSession(ctx, sessionID, req.NewName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to rename session",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "rename_session", "session", sessionID, req.NewName, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": session,
	})
}

func (h *SessionHandler) GetProcesses(c *gin.Context) {
	sessionID := c.Param("id")
	ctx := context.Background()

	processes, err := h.interactiveSvc.GetProcesses(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to get processes",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": processes,
	})
}

func (h *SessionHandler) ExecuteCommand(c *gin.Context) {
	sessionID := c.Param("id")

	var req struct {
		Command string `json:"command" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	output, err := h.interactiveSvc.Execute(ctx, sessionID, req.Command)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to execute command",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "exec", "session", sessionID, req.Command, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"output": output,
		},
	})
}

func (h *SessionHandler) ListDirectory(c *gin.Context) {
	sessionID := c.Param("id")
	path := c.Query("path")

	if path == "" {
		path = "."
	}

	ctx := context.Background()
	files, err := h.interactiveSvc.ListDirectory(ctx, sessionID, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to list directory",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": files,
	})
}

func (h *SessionHandler) DownloadFile(c *gin.Context) {
	sessionID := c.Param("id")
	path := c.Query("path")

	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "path is required",
		})
		return
	}

	ctx := context.Background()
	data, err := h.interactiveSvc.Download(ctx, sessionID, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to download file",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "download", "file", sessionID, path, c.ClientIP(), c.GetHeader("User-Agent"))

	c.Header("Content-Disposition", "attachment; filename="+path)
	c.Data(http.StatusOK, "application/octet-stream", data)
}

func (h *SessionHandler) UploadFile(c *gin.Context) {
	sessionID := c.Param("id")

	var req struct {
		Path string `json:"path" binding:"required"`
		Data string `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	if err := h.interactiveSvc.Upload(ctx, sessionID, req.Path, []byte(req.Data)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to upload file",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "upload", "file", sessionID, req.Path, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "file uploaded successfully",
	})
}

func (h *SessionHandler) Screenshot(c *gin.Context) {
	sessionID := c.Param("id")
	ctx := context.Background()

	data, err := h.interactiveSvc.Screenshot(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to take screenshot",
		})
		return
	}

	c.Data(http.StatusOK, "image/png", data)
}

func (h *SessionHandler) GetPrivileges(c *gin.Context) {
	sessionID := c.Param("id")
	ctx := context.Background()

	privs, err := h.interactiveSvc.GetPrivileges(ctx, sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to get privileges",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": privs,
	})
}

func (h *SessionHandler) MakeToken(c *gin.Context) {
	sessionID := c.Param("id")

	var req struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid request parameters",
		})
		return
	}

	ctx := context.Background()
	if err := h.interactiveSvc.MakeToken(ctx, sessionID, req.Username, req.Password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to make token",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "make_token", "session", sessionID, req.Username, c.ClientIP(), c.GetHeader("User-Agent"))

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "token created successfully",
	})
}

func (h *SessionHandler) ProcessDump(c *gin.Context) {
	sessionID := c.Param("id")

	pidStr := c.Query("pid")
	pid, err := strconv.ParseInt(pidStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    10007,
			"message": "invalid pid",
		})
		return
	}

	ctx := context.Background()
	data, err := h.interactiveSvc.ProcessDump(ctx, sessionID, int32(pid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    10008,
			"message": "failed to dump process",
		})
		return
	}

	userID, _ := c.Get("user_id")
	username, _ := c.Get("username")
	h.auditRepo.Log(userID.(int64), username.(string), "process_dump", "session", sessionID, pidStr, c.ClientIP(), c.GetHeader("User-Agent"))

	c.Header("Content-Disposition", "attachment; filename=process_"+pidStr+".dmp")
	c.Data(http.StatusOK, "application/octet-stream", data)
}
