package database

import (
	"log"
	"os"
	"path/filepath"

	"github.com/cworld1/aniya-blog/backend/internal/config"
	"github.com/cworld1/aniya-blog/backend/internal/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// Init 初始化数据库
func Init(cfg *config.DatabaseConfig) error {
	// 确保数据库目录存在
	dir := filepath.Dir(cfg.DataSource)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 设置日志级别
	logLevel := logger.Silent
	if os.Getenv("GIN_MODE") == "debug" {
		logLevel = logger.Info
	}

	var err error
	DB, err = gorm.Open(sqlite.Open(cfg.DataSource), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

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
