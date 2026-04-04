package config

import (
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         string
	Mode         string // debug, release, test
	AllowOrigins []string
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver       string
	DataSource   string
	Host         string // MySQL/PostgreSQL 主机地址
	Port         string // MySQL/PostgreSQL 端口
	Database     string // 数据库名
	Username     string // 数据库用户名
	Password     string // 数据库密码
	MaxIdleConns int
	MaxOpenConns int
	ConnMaxLife  time.Duration
	AutoMigrate  bool
}

// JWTConfig JWT 配置
type JWTConfig struct {
	Secret     string
	ExpireTime time.Duration
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string
	FilePath   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// Load 加载配置
func Load() *Config {
	// 尝试加载 .env 文件
	_ = godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Mode: getEnv("SERVER_MODE", "debug"),
			AllowOrigins: parseAllowOrigins(getEnv("ALLOW_ORIGIN", "http://localhost:4321,http://127.0.0.1:4321")),
		},
		Database: DatabaseConfig{
			Driver:       getEnv("DB_DRIVER", "sqlite"),
			Host:         getEnv("DB_HOST", "localhost"),
			Port:         getEnv("DB_PORT", "3306"),
			Database:     getEnv("DB_DATABASE", "aniya_blog"),
			Username:     getEnv("DB_USERNAME", "root"),
			Password:     getEnv("DB_PASSWORD", ""),
			DataSource:   getEnv("DB_SOURCE", ""),
			MaxIdleConns: 10,
			MaxOpenConns: 100,
			ConnMaxLife:  time.Hour,
			AutoMigrate:  true,
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "aniya-blog-secret-key-change-in-production"),
			ExpireTime: 24 * time.Hour,
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			FilePath:   getEnv("LOG_FILE", "logs/server.log"),
			MaxSize:    100,
			MaxBackups: 3,
			MaxAge:     28,
			Compress:   true,
		},
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// parseAllowOrigins 解析 CORS 允许的域名列表（逗号分隔）
func parseAllowOrigins(allowOrigin string) []string {
	if allowOrigin == "" {
		return []string{"http://localhost:4321"}
	}
	origins := strings.Split(allowOrigin, ",")
	for i, origin := range origins {
		origins[i] = strings.TrimSpace(origin)
	}
	return origins
}
