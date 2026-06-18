package service

import (
	"github/mouhe/todolist/internal/model"
	"github/mouhe/todolist/internal/repository"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// 创建任务（你可以在这里加业务校验，比如标题不能为空）
func (s *TaskService) CreateTask(title, desc string, priority int) (*model.Task, error) {
	task := &model.Task{
		Title:       title,
		Description: desc,
		Priority:    priority,
		Status:      "pending",
	}
	err := s.repo.Create(task)
	return task, err
}

// 获取任务列表（封装分页默认值）
func (s *TaskService) GetTasks(page, pageSize int) ([]model.Task, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}
	return s.repo.List(page, pageSize)
}

func (s *TaskService) Update(task *model.Task) error {
	return s.repo.Update(task)
}
func (s *TaskService) DeleteTask(id int) error {
	return s.repo.Delete(uint(id))
}

func (s *TaskService) GetTaskById(id uint) (*model.Task, error) {
	return s.repo.GetByID(id)
}
