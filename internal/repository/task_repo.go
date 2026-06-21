package repository

import (
	"errors"
	"github/mouhe/todolist/internal/model"
	"github/mouhe/todolist/internal/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *model.Task) error
	GetByID(id uint) (*model.Task, error)
	List(page, pageSize int) ([]model.Task, int64, error)
	Update(task *model.Task) error
	Delete(id uint) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *taskRepository {
	return &taskRepository{db: db}
}

func (t *taskRepository) Create(task *model.Task) error {
	logger.Debug("Creating task in repository", zap.String("title", task.Title))
	err := t.db.Create(task).Error
	if err != nil {
		logger.Error("Database error creating task", zap.Error(err))
	}
	return err
}
func (t *taskRepository) GetByID(id uint) (*model.Task, error) {
	logger.Debug("Getting task by ID in repository", zap.Uint("id", id))
	var task model.Task
	err := t.db.First(&task, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Warn("Task not found", zap.Uint("id", id))
		return nil, gorm.ErrRecordNotFound
	}
	if err != nil {
		logger.Error("Database error getting task", zap.Uint("id", id), zap.Error(err))
	}
	return &task, nil
}
func (t *taskRepository) List(page, pageSize int) ([]model.Task, int64, error) {
	logger.Debug("Listing tasks in repository", zap.Int("page", page), zap.Int("pageSize", pageSize))
	var tasks []model.Task
	var total int64
	if err := t.db.Model(&model.Task{}).Count(&total).Error; err != nil {
		logger.Error("Database error counting tasks", zap.Error(err))
		return nil, 0, err
	}

	// 分页查询（偏移量 = (页码-1) * 每页条数）
	offset := (page - 1) * pageSize
	err := t.db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&tasks).Error
	if err != nil {
		logger.Error("Database error listing tasks", zap.Error(err))
	}
	return tasks, total, err
}
func (t *taskRepository) Update(task *model.Task) error {
	logger.Debug("Updating task in repository", zap.Uint("id", task.ID))
	err := t.db.Save(task).Error
	if err != nil {
		logger.Error("Database error updating task", zap.Uint("id", task.ID), zap.Error(err))
	}
	return err
}

func (t *taskRepository) Delete(id uint) error {
	logger.Debug("Deleting task in repository", zap.Uint("id", id))
	err := t.db.Delete(&model.Task{}, id).Error
	if err != nil {
		logger.Error("Database error deleting task", zap.Uint("id", id), zap.Error(err))
	}
	return err
}
