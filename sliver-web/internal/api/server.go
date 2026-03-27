package api

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/MonkeyCode/sliver-web/internal/bot"
	"github.com/MonkeyCode/sliver-web/internal/core"
	"github.com/MonkeyCode/sliver-web/internal/db"
	"github.com/MonkeyCode/sliver-web/internal/rpc"
	"github.com/MonkeyCode/sliver-web/internal/ws"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Server struct {
	config       *core.Config
	database     *db.Database
	sliverClient *rpc.SliverClient
	wsServer     *ws.Server
	botServer    *bot.BotServer
	engine       *gin.Engine
	authManager  *core.AuthManager
	totpManager  *core.TOTPManager
	mu           sync.RWMutex
	running      bool
}

func NewServer(config *core.Config, database *db.Database, sliverClient *rpc.SliverClient, wsServer *ws.Server, botServer *bot.BotServer) *Server {
	gin.SetMode(gin.ReleaseMode)

	authManager := core.NewAuthManager(config.Server.JWTSecret, config.GetJWTExpiry())
	totpManager := core.NewTOTPManager("sliver-web")

	server := &Server{
		config:       config,
		database:     database,
		sliverClient: sliverClient,
		wsServer:     wsServer,
		botServer:    botServer,
		authManager:  authManager,
		totpManager:  totpManager,
	}

	server.setupRouter()
	return server
}

func (s *Server) setupRouter() {
	s.engine = gin.New()
	s.engine.Use(gin.Recovery())
	s.engine.Use(s.loggingMiddleware())
	s.engine.Use(s.corsMiddleware())

	api := s.engine.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", s.handleLogin)
			auth.POST("/2fa/verify", s.handle2FAVerify)
			auth.POST("/logout", s.authMiddleware(), s.handleLogout)
			auth.GET("/profile", s.authMiddleware(), s.handleGetProfile)
		}

		sessions := api.Group("/sessions")
		sessions.Use(s.authMiddleware())
		{
			sessions.GET("", s.operatorMiddleware(), s.handleListSessions)
			sessions.GET("/:id", s.operatorMiddleware(), s.handleGetSession)
			sessions.DELETE("/:id", s.operatorMiddleware(), s.handleKillSession)
			sessions.POST("/:id/rename", s.operatorMiddleware(), s.handleRenameSession)
			sessions.GET("/:id/processes", s.operatorMiddleware(), s.handleGetProcesses)
			sessions.POST("/:id/exec", s.operatorMiddleware(), s.handleExecuteCommand)
			sessions.GET("/:id/ls", s.operatorMiddleware(), s.handleListDirectory)
			sessions.GET("/:id/download", s.operatorMiddleware(), s.handleDownloadFile)
			sessions.POST("/:id/upload", s.operatorMiddleware(), s.handleUploadFile)
			sessions.GET("/:id/screenshot", s.operatorMiddleware(), s.handleScreenshot)
			sessions.GET("/:id/privs", s.operatorMiddleware(), s.handleGetPrivileges)
			sessions.POST("/:id/token", s.operatorMiddleware(), s.handleMakeToken)
			sessions.POST("/:id/dump", s.adminMiddleware(), s.handleProcessDump)
		}

		beacons := api.Group("/beacons")
		beacons.Use(s.authMiddleware())
		{
			beacons.GET("", s.operatorMiddleware(), s.handleListBeacons)
			beacons.GET("/:id", s.operatorMiddleware(), s.handleGetBeacon)
			beacons.GET("/:id/tasks", s.operatorMiddleware(), s.handleGetBeaconTasks)
			beacons.POST("/:id/task", s.operatorMiddleware(), s.handleExecuteBeaconTask)
			beacons.DELETE("/:id/tasks/:tid", s.operatorMiddleware(), s.handleCancelBeaconTask)
			beacons.GET("/:id/tasks/:tid/output", s.operatorMiddleware(), s.handleGetTaskOutput)
		}

		listeners := api.Group("/listeners")
		listeners.Use(s.authMiddleware())
		{
			listeners.GET("", s.operatorMiddleware(), s.handleListListeners)
			listeners.POST("/mtls", s.adminMiddleware(), s.handleStartMTLSListener)
			listeners.POST("/dns", s.adminMiddleware(), s.handleStartDNSListener)
			listeners.POST("/http", s.adminMiddleware(), s.handleStartHTTPListener)
			listeners.POST("/https", s.adminMiddleware(), s.handleStartHTTPSListener)
			listeners.POST("/wg", s.adminMiddleware(), s.handleStartWGListener)
			listeners.DELETE("/:id", s.adminMiddleware(), s.handleStopListener)
		}

		profiles := api.Group("/profiles")
		profiles.Use(s.authMiddleware())
		{
			profiles.GET("", s.operatorMiddleware(), s.handleListProfiles)
			profiles.POST("", s.adminMiddleware(), s.handleCreateProfile)
			profiles.PUT("/:id", s.adminMiddleware(), s.handleUpdateProfile)
			profiles.DELETE("/:name", s.adminMiddleware(), s.handleDeleteProfile)
		}

		generate := api.Group("/generate")
		generate.Use(s.authMiddleware())
		{
			generate.POST("", s.operatorMiddleware(), s.handleGenerateImplant)
			generate.GET("/download/:id", s.operatorMiddleware(), s.handleDownloadImplant)
		}

		users := api.Group("/users")
		users.Use(s.authMiddleware())
		{
			users.GET("", s.adminMiddleware(), s.handleListUsers)
			users.POST("", s.adminMiddleware(), s.handleCreateUser)
			users.GET("/:id", s.adminMiddleware(), s.handleGetUser)
			users.PUT("/:id", s.adminMiddleware(), s.handleUpdateUser)
			users.DELETE("/:id", s.adminMiddleware(), s.handleDeleteUser)
			users.POST("/:id/password", s.adminMiddleware(), s.handleChangePassword)
			users.POST("/:id/totp/enable", s.adminMiddleware(), s.handleEnableTOTP)
			users.POST("/:id/totp/disable", s.adminMiddleware(), s.handleDisableTOTP)
		}

		bot := api.Group("/bot")
		bot.Use(s.adminMiddleware())
		{
			bot.GET("/config", s.handleGetBotConfig)
			bot.PUT("/config", s.handleUpdateBotConfig)
			bot.POST("/start", s.handleStartBot)
			bot.POST("/stop", s.handleStopBot)
			bot.GET("/whitelist", s.handleListWhitelist)
			bot.POST("/whitelist", s.handleAddWhitelist)
			bot.DELETE("/whitelist/:id", s.handleRemoveWhitelist)
			bot.GET("/alerts", s.handleListAlertRules)
			bot.POST("/alerts", s.handleCreateAlertRule)
			bot.PUT("/alerts/:id", s.handleUpdateAlertRule)
			bot.DELETE("/alerts/:id", s.handleDeleteAlertRule)
		}

		audit := api.Group("/audit")
		audit.Use(s.authMiddleware())
		{
			audit.GET("", s.auditorMiddleware(), s.handleListAuditLogs)
			audit.GET("/:id", s.auditorMiddleware(), s.handleGetAuditLog)
		}
	}

	s.engine.GET("/ws", s.wsHandler())
	s.engine.GET("/health", s.handleHealth)
}

func (s *Server) Start() error {
	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return fmt.Errorf("server is already running")
	}
	s.running = true
	s.mu.Unlock()

	addr := s.config.GetServerAddr()
	logrus.Infof("Starting API server on %s", addr)

	return s.engine.Run(addr)
}

func (s *Server) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.running = false
	return nil
}

func (s *Server) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}

func (s *Server) loggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		logrus.Infof("%s %s", c.Request.Method, c.Request.URL.Path)
		c.Next()
	}
}

func (s *Server) corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func (s *Server) wsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		claims, err := s.authManager.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		s.wsServer.HandleWebSocket(c.Writer, c.Request, claims.UserID, claims.Username, claims.Role)
	}
}

func (s *Server) handleHealth(c *gin.Context) {
	status := "healthy"
	if !s.sliverClient.IsConnected() {
		status = "degraded"
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        status,
		"sliver_server": s.sliverClient.IsConnected(),
	})
}
