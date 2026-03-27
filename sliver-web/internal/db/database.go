package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	conn *sql.DB
	mu   sync.RWMutex
	path string
}

func NewDatabase(path string) (*Database, error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create database directory: %w", err)
	}

	conn, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	conn.SetMaxOpenConns(1)
	conn.SetMaxIdleConns(1)

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	db := &Database{
		conn: conn,
		path: path,
	}

	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return db, nil
}

func (d *Database) Close() error {
	return d.conn.Close()
}

func (d *Database) Begin() (*sql.Tx, error) {
	return d.conn.Begin()
}

func (d *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.conn.Query(query, args...)
}

func (d *Database) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.conn.QueryRow(query, args...)
}

func (d *Database) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.conn.Exec(query, args...)
}

func (d *Database) GetConn() *sql.DB {
	return d.conn
}

func (d *Database) migrate() error {
	migrations := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			role TEXT NOT NULL CHECK(role IN ('super_admin', 'operator', 'auditor')),
			totp_secret TEXT,
			totp_enabled INTEGER DEFAULT 0,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS telegram_whitelist (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER REFERENCES users(id),
			telegram_user_id TEXT UNIQUE NOT NULL,
			telegram_username TEXT,
			permissions TEXT,
			is_active INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS bot_config (
			id INTEGER PRIMARY KEY CHECK (id = 1),
			token TEXT,
			mode TEXT CHECK(mode IN ('polling', 'webhook')) DEFAULT 'polling',
			webhook_url TEXT,
			enabled INTEGER DEFAULT 0,
			locked INTEGER DEFAULT 0,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS audit_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER REFERENCES users(id),
			action TEXT NOT NULL,
			resource TEXT NOT NULL,
			resource_id TEXT,
			details TEXT,
			ip_address TEXT,
			user_agent TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS alert_rules (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			event_type TEXT NOT NULL,
			severity TEXT CHECK(severity IN ('info', 'warning', 'critical')) DEFAULT 'info',
			telegram_ids TEXT,
			email_enabled INTEGER DEFAULT 0,
			email_recipients TEXT,
			enabled INTEGER DEFAULT 1,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_user ON audit_logs(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action)`,
		`CREATE INDEX IF NOT EXISTS idx_audit_logs_created ON audit_logs(created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_telegram_whitelist_tgid ON telegram_whitelist(telegram_user_id)`,
	}

	for _, migration := range migrations {
		if _, err := d.Exec(migration); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
