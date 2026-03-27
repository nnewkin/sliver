package models

import (
	"time"
)

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	TOTPSecret   string    `json:"-"`
	TOTPEnabled  bool      `json:"totp_enabled"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type TelegramWhitelist struct {
	ID               int64     `json:"id"`
	UserID           int64     `json:"user_id"`
	TelegramUserID   string    `json:"telegram_user_id"`
	TelegramUsername string    `json:"telegram_username"`
	Permissions      []string  `json:"permissions"`
	IsActive         bool      `json:"is_active"`
	CreatedAt        time.Time `json:"created_at"`
}

type BotConfig struct {
	ID         int64     `json:"id"`
	Token      string    `json:"-"`
	Mode       string    `json:"mode"`
	WebhookURL string    `json:"webhook_url"`
	Enabled    bool      `json:"enabled"`
	Locked     bool      `json:"locked"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	Username   string    `json:"username,omitempty"`
	Action     string    `json:"action"`
	Resource   string    `json:"resource"`
	ResourceID string    `json:"resource_id,omitempty"`
	Details    string    `json:"details,omitempty"`
	IPAddress  string    `json:"ip_address,omitempty"`
	UserAgent  string    `json:"user_agent,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type AlertRule struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	EventType       string    `json:"event_type"`
	Severity        string    `json:"severity"`
	TelegramIDs     []string  `json:"telegram_ids"`
	EmailEnabled    bool      `json:"email_enabled"`
	EmailRecipients []string  `json:"email_recipients"`
	Enabled         bool      `json:"enabled"`
	CreatedAt       time.Time `json:"created_at"`
}

const (
	RoleSuperAdmin = "super_admin"
	RoleOperator   = "operator"
	RoleAuditor    = "auditor"
)

const (
	EventTypeSessionOnline       = "session:online"
	EventTypeSessionOffline      = "session:offline"
	EventTypeBeaconCheckin       = "beacon:checkin"
	EventTypeBeaconTaskDone      = "beacon:task_complete"
	EventTypeAlertSession        = "alert:session"
	EventTypeAlertSecurity       = "alert:security"
	EventTypeListenerChange      = "listener:change"
	EventTypePrivilegeEscalation = "privilege:escalation"
	EventTypeCredentialDump      = "credential:dump"
)

const (
	SeverityInfo     = "info"
	SeverityWarning  = "warning"
	SeverityCritical = "critical"
)

func (u *User) HasRole(required string) bool {
	roleHierarchy := map[string]int{
		RoleSuperAdmin: 3,
		RoleOperator:   2,
		RoleAuditor:    1,
	}

	userLevel := roleHierarchy[u.Role]
	requiredLevel := roleHierarchy[required]

	return userLevel >= requiredLevel
}

func (u *User) CanManageUsers() bool {
	return u.Role == RoleSuperAdmin
}

func (u *User) CanExecuteCommands() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleOperator
}

func (u *User) CanViewAuditLogs() bool {
	return u.Role == RoleSuperAdmin || u.Role == RoleAuditor
}
