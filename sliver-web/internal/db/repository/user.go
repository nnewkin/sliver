package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/MonkeyCode/sliver-web/internal/core"
	"github.com/MonkeyCode/sliver-web/internal/db"
	"github.com/MonkeyCode/sliver-web/internal/db/models"
)

type UserRepository struct {
	db *db.Database
}

func NewUserRepository(database *db.Database) *UserRepository {
	return &UserRepository{db: database}
}

func (r *UserRepository) Create(username, password, role string) (*models.User, error) {
	hash, err := core.HashPassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	result, err := r.db.Exec(
		"INSERT INTO users (username, password_hash, role) VALUES (?, ?, ?)",
		username, hash, role,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	id, _ := result.LastInsertId()
	return r.GetByID(id)
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	row := r.db.QueryRow(
		"SELECT id, username, password_hash, role, totp_secret, totp_enabled, is_active, created_at, updated_at FROM users WHERE id = ?",
		id,
	)

	var user models.User
	var totpSecret sql.NullString
	var totpEnabled int

	err := row.Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Role,
		&totpSecret, &totpEnabled, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.TOTPSecret = totpSecret.String
	user.TOTPEnabled = totpEnabled == 1
	user.IsActive = user.IsActive

	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	row := r.db.QueryRow(
		"SELECT id, username, password_hash, role, totp_secret, totp_enabled, is_active, created_at, updated_at FROM users WHERE username = ?",
		username,
	)

	var user models.User
	var totpSecret sql.NullString
	var totpEnabled int

	err := row.Scan(
		&user.ID, &user.Username, &user.PasswordHash, &user.Role,
		&totpSecret, &totpEnabled, &user.IsActive,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user.TOTPSecret = totpSecret.String
	user.TOTPEnabled = totpEnabled == 1
	user.IsActive = user.IsActive

	return &user, nil
}

func (r *UserRepository) List() ([]*models.User, error) {
	rows, err := r.db.Query(
		"SELECT id, username, password_hash, role, totp_secret, totp_enabled, is_active, created_at, updated_at FROM users ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		var totpSecret sql.NullString
		var totpEnabled int

		err := rows.Scan(
			&user.ID, &user.Username, &user.PasswordHash, &user.Role,
			&totpSecret, &totpEnabled, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			continue
		}

		user.TOTPSecret = totpSecret.String
		user.TOTPEnabled = totpEnabled == 1
		user.IsActive = user.IsActive
		users = append(users, &user)
	}

	return users, nil
}

func (r *UserRepository) Update(user *models.User) error {
	_, err := r.db.Exec(
		"UPDATE users SET username = ?, role = ?, totp_secret = ?, totp_enabled = ?, is_active = ?, updated_at = ? WHERE id = ?",
		user.Username, user.Role, user.TOTPSecret, user.TOTPEnabled, user.IsActive, time.Now(), user.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

func (r *UserRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	return nil
}

func (r *UserRepository) SetPassword(id int64, password string) error {
	hash, err := core.HashPassword(password)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	_, err = r.db.Exec("UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?", hash, time.Now(), id)
	if err != nil {
		return fmt.Errorf("failed to set password: %w", err)
	}
	return nil
}

func (r *UserRepository) EnableTOTP(id int64, secret string) error {
	_, err := r.db.Exec(
		"UPDATE users SET totp_secret = ?, totp_enabled = 1, updated_at = ? WHERE id = ?",
		secret, time.Now(), id,
	)
	if err != nil {
		return fmt.Errorf("failed to enable TOTP: %w", err)
	}
	return nil
}

func (r *UserRepository) DisableTOTP(id int64) error {
	_, err := r.db.Exec(
		"UPDATE users SET totp_secret = '', totp_enabled = 0, updated_at = ? WHERE id = ?",
		time.Now(), id,
	)
	if err != nil {
		return fmt.Errorf("failed to disable TOTP: %w", err)
	}
	return nil
}

type TelegramWhitelistRepository struct {
	db *db.Database
}

func NewTelegramWhitelistRepository(database *db.Database) *TelegramWhitelistRepository {
	return &TelegramWhitelistRepository{db: database}
}

func (r *TelegramWhitelistRepository) Create(userID int64, telegramUserID, username string, permissions []string) (*models.TelegramWhitelist, error) {
	permsJSON, _ := json.Marshal(permissions)

	result, err := r.db.Exec(
		"INSERT INTO telegram_whitelist (user_id, telegram_user_id, telegram_username, permissions) VALUES (?, ?, ?, ?)",
		userID, telegramUserID, username, string(permsJSON),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create telegram whitelist: %w", err)
	}

	id, _ := result.LastInsertId()
	return r.GetByID(id)
}

func (r *TelegramWhitelistRepository) GetByID(id int64) (*models.TelegramWhitelist, error) {
	row := r.db.QueryRow(
		"SELECT id, user_id, telegram_user_id, telegram_username, permissions, is_active, created_at FROM telegram_whitelist WHERE id = ?",
		id,
	)

	var entry models.TelegramWhitelist
	var permissionsJSON string

	err := row.Scan(
		&entry.ID, &entry.UserID, &entry.TelegramUserID, &entry.TelegramUsername,
		&permissionsJSON, &entry.IsActive, &entry.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get telegram whitelist: %w", err)
	}

	json.Unmarshal([]byte(permissionsJSON), &entry.Permissions)
	return &entry, nil
}

func (r *TelegramWhitelistRepository) GetByTelegramUserID(telegramUserID string) (*models.TelegramWhitelist, error) {
	row := r.db.QueryRow(
		"SELECT id, user_id, telegram_user_id, telegram_username, permissions, is_active, created_at FROM telegram_whitelist WHERE telegram_user_id = ? AND is_active = 1",
		telegramUserID,
	)

	var entry models.TelegramWhitelist
	var permissionsJSON string

	err := row.Scan(
		&entry.ID, &entry.UserID, &entry.TelegramUserID, &entry.TelegramUsername,
		&permissionsJSON, &entry.IsActive, &entry.CreatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get telegram whitelist: %w", err)
	}

	json.Unmarshal([]byte(permissionsJSON), &entry.Permissions)
	return &entry, nil
}

func (r *TelegramWhitelistRepository) List() ([]*models.TelegramWhitelist, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, telegram_user_id, telegram_username, permissions, is_active, created_at FROM telegram_whitelist ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to list telegram whitelist: %w", err)
	}
	defer rows.Close()

	var entries []*models.TelegramWhitelist
	for rows.Next() {
		var entry models.TelegramWhitelist
		var permissionsJSON string

		err := rows.Scan(
			&entry.ID, &entry.UserID, &entry.TelegramUserID, &entry.TelegramUsername,
			&permissionsJSON, &entry.IsActive, &entry.CreatedAt,
		)
		if err != nil {
			continue
		}

		json.Unmarshal([]byte(permissionsJSON), &entry.Permissions)
		entries = append(entries, &entry)
	}

	return entries, nil
}

func (r *TelegramWhitelistRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM telegram_whitelist WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete telegram whitelist: %w", err)
	}
	return nil
}
