package database

import (
	"fmt"
	"github/mouhe/todolist/internal/model"
	"time"

	"gorm.io/gorm/logger" // GORM 内置 logger

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// internal/database/mysql.go
func InitMySQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info), // 开启 GORM 日志
	})
	if err != nil {
		return nil, fmt.Errorf("open database failed: %w", err)
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("get sql.DB failed: %w", err)
	}

	sqlDB.SetMaxOpenConns(45)                  // 最大打开连接数
	sqlDB.SetMaxIdleConns(20)                  // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour)        // 连接最大存活时间
	sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 连接最大空闲时间

	// 健康检查
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("ping database failed: %w", err)
	}

	// 自动迁移
	if err := db.AutoMigrate(&model.Task{}); err != nil {
		return nil, fmt.Errorf("auto migrate failed: %w", err)
	}

	return db, nil
}
