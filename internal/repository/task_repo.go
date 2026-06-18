package repository

import (
	"errors"
	"github/mouhe/todolist/internal/model"

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
	return t.db.Create(task).Error
}
func (t *taskRepository) GetByID(id uint) (*model.Task, error) {
	var task model.Task
	err := t.db.First(&task, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &task, nil
}
func (t *taskRepository) List(page, pageSize int) ([]model.Task, int64, error) {
	var tasks []model.Task
	var total int64
	if err := t.db.Model(&model.Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询（偏移量 = (页码-1) * 每页条数）
	offset := (page - 1) * pageSize
	err := t.db.Offset(offset).Limit(pageSize).Order("created_at desc").Find(&tasks).Error
	return tasks, total, err
}
func (t *taskRepository) Update(task *model.Task) error {
	return t.db.Save(task).Error
}

func (t *taskRepository) Delete(id uint) error {
	return t.db.Delete(&model.Task{}, id).Error

}
