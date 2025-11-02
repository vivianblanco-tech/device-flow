package config

import (
	"os"
	"strconv"
)

// Config holds all application configuration
type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	Session  SessionConfig
	Google   GoogleOAuthConfig
	SMTP     SMTPConfig
	JIRA     JIRAConfig
	Upload   UploadConfig
	Security SecurityConfig
	Logging  LoggingConfig
}

// AppConfig contains general application settings
type AppConfig struct {
	Environment string
	BaseURL     string
}

// ServerConfig contains HTTP server settings
type ServerConfig struct {
	Host string
	Port string
}

// DatabaseConfig contains database connection settings
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// SessionConfig contains session management settings
type SessionConfig struct {
	Secret string
	MaxAge int
}

// GoogleOAuthConfig contains Google OAuth settings
type GoogleOAuthConfig struct {
	ClientID      string
	ClientSecret  string
	RedirectURL   string
	AllowedDomain string
}

// SMTPConfig contains email server settings
type SMTPConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	From     string
	FromName string
}

// JIRAConfig contains JIRA integration settings
type JIRAConfig struct {
	URL            string
	Username       string
	APIToken       string
	DefaultProject string
	// Legacy OAuth fields (kept for backward compatibility)
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

// UploadConfig contains file upload settings
type UploadConfig struct {
	MaxSize int64
	Path    string
}

// SecurityConfig contains security settings
type SecurityConfig struct {
	CSRFSecret string
}

// LoggingConfig contains logging settings
type LoggingConfig struct {
	Level  string
	Format string
}

// Load reads configuration from environment variables
func Load() *Config {
	return &Config{
		App: AppConfig{
			Environment: getEnv("APP_ENV", "development"),
			BaseURL:     getEnv("APP_BASE_URL", "http://localhost:8080"),
		},
		Server: ServerConfig{
			Host: getEnv("APP_HOST", "localhost"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "laptop_tracking_dev"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Session: SessionConfig{
			Secret: getEnv("SESSION_SECRET", "change-me-in-production"),
			MaxAge: getEnvAsInt("SESSION_MAX_AGE", 86400),
		},
		Google: GoogleOAuthConfig{
			ClientID:      getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret:  getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:   getEnv("GOOGLE_REDIRECT_URL", ""),
			AllowedDomain: getEnv("GOOGLE_ALLOWED_DOMAIN", "bairesdev.com"),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "localhost"),
			Port:     getEnv("SMTP_PORT", "1025"),
			User:     getEnv("SMTP_USER", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", "noreply@laptop-tracking.com"),
			FromName: getEnv("SMTP_FROM_NAME", "Laptop Tracking System"),
		},
		JIRA: JIRAConfig{
			URL:            getEnv("JIRA_URL", ""),
			Username:       getEnv("JIRA_USERNAME", ""),
			APIToken:       getEnv("JIRA_API_TOKEN", ""),
			DefaultProject: getEnv("JIRA_DEFAULT_PROJECT", ""),
			ClientID:       getEnv("JIRA_CLIENT_ID", ""),
			ClientSecret:   getEnv("JIRA_CLIENT_SECRET", ""),
			RedirectURL:    getEnv("JIRA_REDIRECT_URL", ""),
		},
		Upload: UploadConfig{
			MaxSize: getEnvAsInt64("MAX_UPLOAD_SIZE", 10485760), // 10MB default
			Path:    getEnv("UPLOAD_PATH", "./uploads"),
		},
		Security: SecurityConfig{
			CSRFSecret: getEnv("CSRF_SECRET", "change-me-in-production"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
	}
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt retrieves an environment variable as int or returns default
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsInt64 retrieves an environment variable as int64 or returns default
func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}
