package model

import "time"

type Task struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Title       string    `json:"title" gorm:"not null"` // 改这里！json:"title"
	Description string    `json:"description"`           // 建议也加上 json 标签
	Status      string    `json:"status" gorm:"default:'pending'"`
	Priority    int       `json:"priority" gorm:"default:1"`
	Deadline    time.Time `json:"deadline"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
