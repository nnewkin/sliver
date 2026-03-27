package repository

import (
	"encoding/json"
	"fmt"

	"github.com/MonkeyCode/sliver-web/internal/db"
	"github.com/MonkeyCode/sliver-web/internal/db/models"
)

type AuditRepository struct {
	db *db.Database
}

func NewAuditRepository(database *db.Database) *AuditRepository {
	return &AuditRepository{db: database}
}

func (r *AuditRepository) Create(log *models.AuditLog) error {
	result, err := r.db.Exec(
		"INSERT INTO audit_logs (user_id, action, resource, resource_id, details, ip_address, user_agent) VALUES (?, ?, ?, ?, ?, ?, ?)",
		log.UserID, log.Action, log.Resource, log.ResourceID, log.Details, log.IPAddress, log.UserAgent,
	)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	id, _ := result.LastInsertId()
	log.ID = id
	return nil
}

func (r *AuditRepository) Log(userID int64, username, action, resource, resourceID, details, ipAddress, userAgent string) error {
	log := &models.AuditLog{
		UserID:     userID,
		Username:   username,
		Action:     action,
		Resource:   resource,
		ResourceID: resourceID,
		Details:    details,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}
	return r.Create(log)
}

func (r *AuditRepository) GetByID(id int64) (*models.AuditLog, error) {
	row := r.db.QueryRow(
		"SELECT id, user_id, action, resource, resource_id, details, ip_address, user_agent, created_at FROM audit_logs WHERE id = ?",
		id,
	)

	var log models.AuditLog
	var details, ipAddress, userAgent, username sql.NullString

	err := row.Scan(
		&log.ID, &log.UserID, &log.Action, &log.Resource, &log.ResourceID,
		&details, &ipAddress, &userAgent, &log.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get audit log: %w", err)
	}

	log.Details = details.String
	log.IPAddress = ipAddress.String
	log.UserAgent = userAgent.String
	return &log, nil
}

func (r *AuditRepository) List(filter *AuditFilter) ([]*models.AuditLog, error) {
	query := "SELECT id, user_id, action, resource, resource_id, details, ip_address, user_agent, created_at FROM audit_logs WHERE 1=1"
	args := []interface{}{}

	if filter.UserID > 0 {
		query += " AND user_id = ?"
		args = append(args, filter.UserID)
	}

	if filter.Action != "" {
		query += " AND action = ?"
		args = append(args, filter.Action)
	}

	if filter.Resource != "" {
		query += " AND resource = ?"
		args = append(args, filter.Resource)
	}

	if !filter.StartDate.IsZero() {
		query += " AND created_at >= ?"
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		query += " AND created_at <= ?"
		args = append(args, filter.EndDate)
	}

	query += " ORDER BY created_at DESC"

	if filter.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", filter.Limit)
		if filter.Offset > 0 {
			query += fmt.Sprintf(" OFFSET %d", filter.Offset)
		}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to list audit logs: %w", err)
	}
	defer rows.Close()

	var logs []*models.AuditLog
	for rows.Next() {
		var log models.AuditLog
		var details, ipAddress, userAgent sql.NullString

		err := rows.Scan(
			&log.ID, &log.UserID, &log.Action, &log.Resource, &log.ResourceID,
			&details, &ipAddress, &userAgent, &log.CreatedAt,
		)
		if err != nil {
			continue
		}

		log.Details = details.String
		log.IPAddress = ipAddress.String
		log.UserAgent = userAgent.String
		logs = append(logs, &log)
	}

	return logs, nil
}

func (r *AuditRepository) Count(filter *AuditFilter) (int, error) {
	query := "SELECT COUNT(*) FROM audit_logs WHERE 1=1"
	args := []interface{}{}

	if filter.UserID > 0 {
		query += " AND user_id = ?"
		args = append(args, filter.UserID)
	}

	if filter.Action != "" {
		query += " AND action = ?"
		args = append(args, filter.Action)
	}

	if filter.Resource != "" {
		query += " AND resource = ?"
		args = append(args, filter.Resource)
	}

	if !filter.StartDate.IsZero() {
		query += " AND created_at >= ?"
		args = append(args, filter.StartDate)
	}

	if !filter.EndDate.IsZero() {
		query += " AND created_at <= ?"
		args = append(args, filter.EndDate)
	}

	var count int
	row := r.db.QueryRow(query, args...)
	if err := row.Scan(&count); err != nil {
		return 0, fmt.Errorf("failed to count audit logs: %w", err)
	}

	return count, nil
}

func (r *AuditRepository) DeleteOlderThan(days int) (int64, error) {
	result, err := r.db.Exec(
		"DELETE FROM audit_logs WHERE created_at < datetime('now', ?)",
		fmt.Sprintf("-%d days", days),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to delete old audit logs: %w", err)
	}

	return result.RowsAffected()
}

type AuditFilter struct {
	UserID    int64
	Action    string
	Resource  string
	StartDate interface{}
	EndDate   interface{}
	Limit     int
	Offset    int
}

type BotConfigRepository struct {
	db *db.Database
}

func NewBotConfigRepository(database *db.Database) *BotConfigRepository {
	return &BotConfigRepository{db: database}
}

func (r *BotConfigRepository) Get() (*models.BotConfig, error) {
	row := r.db.QueryRow(
		"SELECT id, token, mode, webhook_url, enabled, locked, updated_at FROM bot_config WHERE id = 1",
	)

	var config models.BotConfig
	var token, webhookURL sql.NullString

	err := row.Scan(&config.ID, &token, &config.Mode, &webhookURL, &config.Enabled, &config.Locked, &config.UpdatedAt)
	if err != nil {
		if err.Error() == "sql: no rows in result" {
			return &models.BotConfig{ID: 1, Mode: "polling", Enabled: false}, nil
		}
		return nil, fmt.Errorf("failed to get bot config: %w", err)
	}

	config.Token = token.String
	config.WebhookURL = webhookURL.String
	return &config, nil
}

func (r *BotConfigRepository) Update(config *models.BotConfig) error {
	_, err := r.db.Exec(
		"INSERT OR REPLACE INTO bot_config (id, token, mode, webhook_url, enabled, locked, updated_at) VALUES (1, ?, ?, ?, ?, ?, datetime('now'))",
		config.Token, config.Mode, config.WebhookURL, config.Enabled, config.Locked,
	)
	if err != nil {
		return fmt.Errorf("failed to update bot config: %w", err)
	}
	return nil
}

type AlertRuleRepository struct {
	db *db.Database
}

func NewAlertRuleRepository(database *db.Database) *AlertRuleRepository {
	return &AlertRuleRepository{db: database}
}

func (r *AlertRuleRepository) Create(rule *models.AlertRule) error {
	telegramIDsJSON, _ := json.Marshal(rule.TelegramIDs)
	emailRecipientsJSON, _ := json.Marshal(rule.EmailRecipients)

	result, err := r.db.Exec(
		"INSERT INTO alert_rules (name, event_type, severity, telegram_ids, email_enabled, email_recipients, enabled) VALUES (?, ?, ?, ?, ?, ?, ?)",
		rule.Name, rule.EventType, rule.Severity, string(telegramIDsJSON), rule.EmailEnabled, string(emailRecipientsJSON), rule.Enabled,
	)
	if err != nil {
		return fmt.Errorf("failed to create alert rule: %w", err)
	}

	id, _ := result.LastInsertId()
	rule.ID = id
	return nil
}

func (r *AlertRuleRepository) GetByID(id int64) (*models.AlertRule, error) {
	row := r.db.QueryRow(
		"SELECT id, name, event_type, severity, telegram_ids, email_enabled, email_recipients, enabled, created_at FROM alert_rules WHERE id = ?",
		id,
	)

	var rule models.AlertRule
	var telegramIDsJSON, emailRecipientsJSON string

	err := row.Scan(
		&rule.ID, &rule.Name, &rule.EventType, &rule.Severity,
		&telegramIDsJSON, &rule.EmailEnabled, &emailRecipientsJSON, &rule.Enabled, &rule.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get alert rule: %w", err)
	}

	json.Unmarshal([]byte(telegramIDsJSON), &rule.TelegramIDs)
	json.Unmarshal([]byte(emailRecipientsJSON), &rule.EmailRecipients)
	return &rule, nil
}

func (r *AlertRuleRepository) List() ([]*models.AlertRule, error) {
	rows, err := r.db.Query(
		"SELECT id, name, event_type, severity, telegram_ids, email_enabled, email_recipients, enabled, created_at FROM alert_rules ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list alert rules: %w", err)
	}
	defer rows.Close()

	var rules []*models.AlertRule
	for rows.Next() {
		var rule models.AlertRule
		var telegramIDsJSON, emailRecipientsJSON string

		err := rows.Scan(
			&rule.ID, &rule.Name, &rule.EventType, &rule.Severity,
			&telegramIDsJSON, &rule.EmailEnabled, &emailRecipientsJSON, &rule.Enabled, &rule.CreatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(telegramIDsJSON), &rule.TelegramIDs)
		json.Unmarshal([]byte(emailRecipientsJSON), &rule.EmailRecipients)
		rules = append(rules, &rule)
	}

	return rules, nil
}

func (r *AlertRuleRepository) ListByEventType(eventType string) ([]*models.AlertRule, error) {
	rows, err := r.db.Query(
		"SELECT id, name, event_type, severity, telegram_ids, email_enabled, email_recipients, enabled, created_at FROM alert_rules WHERE event_type = ? AND enabled = 1",
		eventType,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list alert rules: %w", err)
	}
	defer rows.Close()

	var rules []*models.AlertRule
	for rows.Next() {
		var rule models.AlertRule
		var telegramIDsJSON, emailRecipientsJSON string

		err := rows.Scan(
			&rule.ID, &rule.Name, &rule.EventType, &rule.Severity,
			&telegramIDsJSON, &rule.EmailEnabled, &emailRecipientsJSON, &rule.Enabled, &rule.CreatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(telegramIDsJSON), &rule.TelegramIDs)
		json.Unmarshal([]byte(emailRecipientsJSON), &rule.EmailRecipients)
		rules = append(rules, &rule)
	}

	return rules, nil
}

func (r *AlertRuleRepository) Update(rule *models.AlertRule) error {
	telegramIDsJSON, _ := json.Marshal(rule.TelegramIDs)
	emailRecipientsJSON, _ := json.Marshal(rule.EmailRecipients)

	_, err := r.db.Exec(
		"UPDATE alert_rules SET name = ?, event_type = ?, severity = ?, telegram_ids = ?, email_enabled = ?, email_recipients = ?, enabled = ? WHERE id = ?",
		rule.Name, rule.EventType, rule.Severity, string(telegramIDsJSON), rule.EmailEnabled, string(emailRecipientsJSON), rule.Enabled, rule.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update alert rule: %w", err)
	}
	return nil
}

func (r *AlertRuleRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM alert_rules WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete alert rule: %w", err)
	}
	return nil
}
