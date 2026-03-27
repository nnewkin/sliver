package bot

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/MonkeyCode/sliver-web/internal/core"
	"github.com/MonkeyCode/sliver-web/internal/db"
	"github.com/MonkeyCode/sliver-web/internal/db/models"
	"github.com/MonkeyCode/sliver-web/internal/db/repository"
	"github.com/MonkeyCode/sliver-web/internal/rpc"
	"github.com/MonkeyCode/sliver-web/internal/rpc/services"
	"github.com/MonkeyCode/sliver-web/internal/ws"
	"github.com/sirupsen/logrus"
)

type BotServer struct {
	config       *core.Config
	database     *db.Database
	sliverClient *rpc.SliverClient
	wsServer     *ws.Server
	userRepo     *repository.UserRepository
	tgRepo       *repository.TelegramWhitelistRepository
	botCfgRepo   *repository.BotConfigRepository
	auditRepo    *repository.AuditRepository

	running bool
	mu      sync.RWMutex
}

func NewBotServer(config *core.Config, database *db.Database, sliverClient *rpc.SliverClient, wsServer *ws.Server) *BotServer {
	return &BotServer{
		config:       config,
		database:     database,
		sliverClient: sliverClient,
		wsServer:     wsServer,
		userRepo:     repository.NewUserRepository(database),
		tgRepo:       repository.NewTelegramWhitelistRepository(database),
		botCfgRepo:   repository.NewBotConfigRepository(database),
		auditRepo:    repository.NewAuditRepository(database),
	}
}

func (b *BotServer) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.running {
		return fmt.Errorf("bot is already running")
	}

	cfg, err := b.botCfgRepo.Get()
	if err != nil {
		return fmt.Errorf("failed to get bot config: %w", err)
	}

	if cfg.Token == "" {
		return fmt.Errorf("bot token not configured")
	}

	if !cfg.Enabled {
		return fmt.Errorf("bot is disabled")
	}

	b.running = true

	go b.runEventLoop()

	logrus.Info("Telegram bot server started")
	return nil
}

func (b *BotServer) Stop() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.running = false
	logrus.Info("Telegram bot server stopped")
}

func (b *BotServer) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

func (b *BotServer) runEventLoop() {
	logrus.Info("Bot event loop started")
}

func (b *BotServer) IsWhitelisted(telegramUserID string) bool {
	entry, err := b.tgRepo.GetByTelegramUserID(telegramUserID)
	if err != nil {
		return false
	}
	return entry.IsActive
}

func (b *BotServer) GetUserPermissions(telegramUserID string) []string {
	entry, err := b.tgRepo.GetByTelegramUserID(telegramUserID)
	if err != nil {
		return nil
	}
	return entry.Permissions
}

func (b *BotServer) HasPermission(telegramUserID, permission string) bool {
	perms := b.GetUserPermissions(telegramUserID)
	for _, p := range perms {
		if p == permission || p == "*" {
			return true
		}
	}
	return false
}

func (b *BotServer) SendAlert(rule *models.AlertRule, message string) error {
	if !b.IsRunning() {
		return fmt.Errorf("bot is not running")
	}

	for _, tgID := range rule.TelegramIDs {
		b.sendMessage(tgID, message)
	}

	return nil
}

func (b *BotServer) sendMessage(telegramUserID, text string) error {
	return nil
}

func (b *BotServer) handleCommand(telegramUserID string, command string) (string, error) {
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return "", nil
	}

	cmd := parts[0]
	args := parts[1:]

	if !b.IsWhitelisted(telegramUserID) {
		return "Access denied. You are not in the whitelist.", nil
	}

	switch cmd {
	case "/start":
		return b.handleStart(telegramUserID)
	case "/help":
		return b.handleHelp(telegramUserID)
	case "/sessions":
		return b.handleSessions(telegramUserID)
	case "/beacons":
		return b.handleBeacons(telegramUserID)
	case "/exec":
		return b.handleExec(telegramUserID, args)
	case "/generate":
		return b.handleGenerate(telegramUserID)
	case "/listener":
		return b.handleListener(telegramUserID)
	case "/alerts":
		return b.handleAlerts(telegramUserID)
	case "/lock":
		return b.handleLock(telegramUserID)
	default:
		return fmt.Sprintf("Unknown command: %s", cmd), nil
	}
}

func (b *BotServer) handleStart(telegramUserID string) (string, error) {
	if !b.IsWhitelisted(telegramUserID) {
		return "Access denied. Please contact an administrator.", nil
	}

	entry, err := b.tgRepo.GetByTelegramUserID(telegramUserID)
	if err != nil {
		return "Error retrieving user information.", nil
	}

	return fmt.Sprintf("Welcome to Sliver Web Bot!\n\nUser: %s\nRole: %s\n\nUse /help to see available commands.", entry.TelegramUsername, b.getUserRole(entry.UserID)), nil
}

func (b *BotServer) handleHelp(telegramUserID string) (string, error) {
	if !b.IsWhitelisted(telegramUserID) {
		return "Access denied.", nil
	}

	helpText := `
Available Commands:
/start - Start the bot
/help - Show this help message
/sessions - List all active sessions
/beacons - List all beacons
/exec <session_id> <command> - Execute a command
/generate - Generate an implant
/listener - Manage listeners
/alerts - View alert history
/lock - Lock the bot (admin only)
`
	return helpText, nil
}

func (b *BotServer) handleSessions(telegramUserID string) (string, error) {
	if !b.IsWhitelisted(telegramUserID) || !b.HasPermission(telegramUserID, "sessions:read") {
		return "You don't have permission to view sessions.", nil
	}

	ctx := context.Background()
	sessionService := services.NewSessionService(b.sliverClient.GetConn())
	sessions, err := sessionService.GetAllSessions(ctx)
	if err != nil {
		return fmt.Sprintf("Error fetching sessions: %v", err), nil
	}

	if len(sessions) == 0 {
		return "No active sessions.", nil
	}

	var sb strings.Builder
	sb.WriteString("Active Sessions:\n\n")
	for _, s := range sessions {
		sb.WriteString(fmt.Sprintf("[%s] %s@%s (%s)\n", s.ID, s.Username, s.Hostname, s.OS))
	}
	return sb.String(), nil
}

func (b *BotServer) handleBeacons(telegramUserID string) (string, error) {
	if !b.IsWhitelisted(telegramUserID) || !b.HasPermission(telegramUserID, "beacons:read") {
		return "You don't have permission to view beacons.", nil
	}

	return "Beacon listing not yet implemented.", nil
}

func (b *BotServer) handleExec(telegramUserID string, args []string) (string, error) {
	if !b.IsWhitelisted(telegramUserID) || !b.HasPermission(telegramUserID, "exec:write") {
		return "You don't have permission to execute commands.", nil
	}

	if len(args) < 2 {
		return "Usage: /exec <session_id> <command>", nil
	}

	sessionID := args[0]
	command := strings.Join(args[1:], " ")

	ctx := context.Background()
	interactiveSvc := services.NewInteractiveService(b.sliverClient.GetConn())
	output, err := interactiveSvc.Execute(ctx, sessionID, command)
	if err != nil {
		return fmt.Sprintf("Error executing command: %v", err), nil
	}

	return fmt.Sprintf("Output:\n%s", output), nil
}

func (b *BotServer) handleGenerate(telegramUserID string) (string, error) {
	if !b.IsWhitelisted(telegramUserID) || !b.HasPermission(telegramUserID, "generate:write") {
		return "You don't have permission to generate implants.", nil
	}

	return "Implant generation should be done via the Web UI.", nil
}

func (b *BotServer) handleListener(telegramUserID string) (string, error) {
	if !b.IsWhitelisted(telegramUserID) || !b.HasPermission(telegramUserID, "listeners:write") {
		return "You don't have permission to manage listeners.", nil
	}

	return "Listener management should be done via the Web UI.", nil
}

func (b *BotServer) handleAlerts(telegramUserID string) (string, error) {
	if !b.IsWhitelisted(telegramUserID) {
		return "Access denied.", nil
	}

	return "Alert history not yet implemented.", nil
}

func (b *BotServer) handleLock(telegramUserID string) (string, error) {
	entry, err := b.tgRepo.GetByTelegramUserID(telegramUserID)
	if err != nil {
		return "Error checking permissions.", nil
	}

	user, err := b.userRepo.GetByID(entry.UserID)
	if err != nil || user.Role != models.RoleSuperAdmin {
		return "Access denied. Super admin only.", nil
	}

	cfg, _ := b.botCfgRepo.Get()
	cfg.Locked = true
	b.botCfgRepo.Update(cfg)

	b.wsServer.BroadcastEvent("system", "bot_locked", nil)

	return "Bot has been locked.", nil
}

func (b *BotServer) getUserRole(userID int64) string {
	user, err := b.userRepo.GetByID(userID)
	if err != nil {
		return "unknown"
	}
	return user.Role
}

func (b *BotServer) LogAudit(userID int64, username, action, resource, resourceID, details, ipAddress, userAgent string) {
	b.auditRepo.Log(userID, username, action, resource, resourceID, details, ipAddress, userAgent)
}
