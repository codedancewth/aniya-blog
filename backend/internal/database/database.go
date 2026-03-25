package database

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/cworld1/aniya-blog/backend/internal/config"
	"github.com/cworld1/aniya-blog/backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库
func Init(cfg *config.DatabaseConfig) error {
	// 根据驱动构建 DSN
	var dsn string
	var err error

	switch cfg.Driver {
	case "mysql":
		dsn = cfg.DataSource
		if dsn == "" {
			// 从各个配置项构建 DSN
			dsn = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
				cfg.Username,
				cfg.Password,
				cfg.Host,
				cfg.Port,
				cfg.Database,
			)
		}
		log.Printf("Connecting to MySQL database: %s:%s/%s", cfg.Host, cfg.Port, cfg.Database)

	case "sqlite":
		fallthrough
	default:
		dsn = cfg.DataSource
		if dsn == "" {
			dsn = "data/aniya.db"
		}
		// 确保数据库目录存在
		dir := filepath.Dir(dsn)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
		log.Printf("Connecting to SQLite database: %s", dsn)
	}

	// 设置日志级别
	logLevel := logger.Silent
	if os.Getenv("GIN_MODE") == "debug" {
		logLevel = logger.Info
	}

	// 根据驱动打开数据库连接
	switch cfg.Driver {
	case "mysql":
		DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
	case "sqlite":
		fallthrough
	default:
		DB, err = gorm.Open(sqlite.Open(dsn), &gorm.Config{
			Logger: logger.Default.LogMode(logLevel),
		})
	}

	if err != nil {
		return err
	}

	// 获取底层 SQL DB
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}

	// 设置连接池配置
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLife)

	// 自动迁移
	if cfg.AutoMigrate {
		if err := Migrate(); err != nil {
			return err
		}
	}

	log.Println("Database initialized successfully")
	return nil
}

// Migrate 自动迁移数据库表
func Migrate() error {
	return DB.AutoMigrate(
		&models.User{},
		&models.Post{},
		&models.Category{},
		&models.Tag{},
		&models.Comment{},
		&models.PageView{},
		&models.Link{},
		&models.Config{},
		&models.Like{},
	)
}

// Close 关闭数据库连接
func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	return DB
}
