package config

import (
	"os"
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
	Driver        string
	DataSource    string
	MaxIdleConns  int
	MaxOpenConns  int
	ConnMaxLife   time.Duration
	AutoMigrate   bool
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
			Port:      getEnv("SERVER_PORT", "8080"),
			Mode:      getEnv("SERVER_MODE", "debug"),
			AllowOrigins: []string{
				getEnv("ALLOW_ORIGIN", "http://localhost:4321"),
			},
		},
		Database: DatabaseConfig{
			Driver:        getEnv("DB_DRIVER", "sqlite"),
			DataSource:    getEnv("DB_SOURCE", "data/aniya.db"),
			MaxIdleConns:  10,
			MaxOpenConns:  100,
			ConnMaxLife:   time.Hour,
			AutoMigrate:   true,
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
