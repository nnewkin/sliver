package core

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Sliver   SliverConfig   `mapstructure:"sliver"`
	Database DatabaseConfig `mapstructure:"database"`
	Telegram TelegramConfig `mapstructure:"telegram"`
	Logging  LoggingConfig  `mapstructure:"logging"`
	Security SecurityConfig `mapstructure:"security"`
}

type ServerConfig struct {
	Host           string `mapstructure:"host"`
	Port           int    `mapstructure:"port"`
	Mode           string `mapstructure:"mode"`
	JWTSecret      string `mapstructure:"jwt_secret"`
	JWTExpiryHours int    `mapstructure:"jwt_expiry_hours"`
}

type SliverConfig struct {
	GRPCAddr string `mapstructure:"grpc_addr"`
	CertsDir string `mapstructure:"certs_dir"`
}

type DatabaseConfig struct {
	Path string `mapstructure:"path"`
}

type TelegramConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	BotToken        string `mapstructure:"bot_token"`
	Mode            string `mapstructure:"mode"`
	WebhookURL      string `mapstructure:"webhook_url"`
	PollingInterval int    `mapstructure:"polling_interval"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	Output string `mapstructure:"output"`
}

type SecurityConfig struct {
	TOTPEnabled            bool `mapstructure:"totp_enabled"`
	SessionTimeoutMinutes  int  `mapstructure:"session_timeout_minutes"`
	MaxLoginAttempts       int  `mapstructure:"max_login_attempts"`
	LockoutDurationMinutes int  `mapstructure:"lockout_duration_minutes"`
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("yaml")

	viper.SetDefault("server.host", "0.0.0.0")
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.mode", "release")
	viper.SetDefault("server.jwt_expiry_hours", 24)
	viper.SetDefault("sliver.grpc_addr", "localhost:50050")
	viper.SetDefault("sliver.certs_dir", "./certs")
	viper.SetDefault("database.path", "./data/sliver-web.db")
	viper.SetDefault("telegram.enabled", false)
	viper.SetDefault("telegram.mode", "polling")
	viper.SetDefault("telegram.polling_interval", 1)
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("security.totp_enabled", true)
	viper.SetDefault("security.session_timeout_minutes", 30)
	viper.SetDefault("security.max_login_attempts", 5)
	viper.SetDefault("security.lockout_duration_minutes", 15)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if config.Server.JWTSecret == "" || config.Server.JWTSecret == "CHANGE_ME_IN_PRODUCTION" {
		return nil, fmt.Errorf("JWT secret must be set in configuration")
	}

	return &config, nil
}

func (c *Config) GetServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

func (c *Config) GetJWTExpiry() time.Duration {
	return time.Duration(c.Server.JWTExpiryHours) * time.Hour
}

func (c *Config) GetSessionTimeout() time.Duration {
	return time.Duration(c.Security.SessionTimeoutMinutes) * time.Minute
}

func (c *Config) GetLockoutDuration() time.Duration {
	return time.Duration(c.Security.LockoutDurationMinutes) * time.Minute
}
