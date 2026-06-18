package database

import (
	"github/mouhe/todolist/internal/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 初始化数据库
func InitMySQL(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		// 禁止创建外键约束（强烈建议生产环境关闭）
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, err
	}
	// 自动迁移
	if err := db.AutoMigrate(&model.Task{}); err != nil {
		return nil, err
	}

	return db, nil
}
