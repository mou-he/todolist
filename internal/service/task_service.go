package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github/mouhe/todolist/internal/model"
	"github/mouhe/todolist/internal/pkg/logger"
	"github/mouhe/todolist/internal/repository"
	"time"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

type TaskService struct {
	repo        repository.TaskRepository
	redisClient *redis.Client
}

func NewTaskService(repo repository.TaskRepository, redisClient *redis.Client) *TaskService {
	return &TaskService{repo: repo, redisClient: redisClient}
}

// 创建任务（你可以在这里加业务校验，比如标题不能为空）
func (s *TaskService) CreateTask(title, desc string, priority int, deadline time.Time) (*model.Task, error) {
	logger.Debug("Creating task in service", zap.String("title", title))
	task := &model.Task{

		Title:       title,
		Description: desc,
		Priority:    priority,
		Status:      "pending",
		Deadline:    deadline,
		CreatedAt:   time.Now(),
	}

	err := s.repo.Create(task)
	if err != nil {
		logger.Error("Failed to create task in repository", zap.Error(err))
	}
	return task, err
}

// GetTasks 获取任务列表（封装分页默认值）
func (s *TaskService) GetTasks(page, pageSize int) ([]model.Task, int64, error) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("tasks:page:%d:size:%d", page, pageSize)

	// 先尝试从缓存获取
	cachedData, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var tasks []model.Task
		if err := json.Unmarshal([]byte(cachedData), &tasks); err == nil {
			logger.Info("Cache hit for tasks", zap.String("key", cacheKey))
			return tasks, int64(len(tasks)), nil
		}
		logger.Warn("Failed to unmarshal cached tasks", zap.Error(err))
	}

	// 缓存未命中，从数据库获取
	tasks, total, err := s.repo.List(page, pageSize)
	if err != nil {
		logger.Error("Failed to get tasks from repository", zap.Error(err))
		return nil, 0, err
	}

	// 将结果存入缓存（有效期 5 分钟）
	data, err := json.Marshal(tasks)
	if err != nil {
		logger.Error("Failed to marshal tasks for cache", zap.Error(err))
	} else {
		s.redisClient.Set(ctx, cacheKey, data, 5*time.Minute)
		logger.Info("Cache created for task list", zap.String("key", cacheKey))
	}

	return tasks, total, nil
}

func (s *TaskService) Update(task *model.Task) error {
	logger.Debug("Updating task in service", zap.Uint("task_id", task.ID))
	err := s.repo.Update(task)
	if err != nil {
		logger.Error("Failed to update task in repository", zap.Uint("task_id", task.ID), zap.Error(err))
		return err
	}
	s.invalidateTaskCache(task.ID)
	s.createTaskCache(task)
	s.invalidateTaskListCache()

	return nil
}

func (s *TaskService) DeleteTask(id int) error {
	logger.Debug("Deleting task in service", zap.Int("task_id", id))
	err := s.repo.Delete(uint(id))
	if err != nil {
		logger.Error("Failed to delete task in repository", zap.Int("task_id", id), zap.Error(err))
		return err
	}
	s.invalidateTaskCache(uint(id))
	s.invalidateTaskListCache()

	return nil
}

// GetTaskById 根据ID获取任务
func (s *TaskService) GetTaskById(id uint) (*model.Task, error) {
	logger.Debug("Getting task by ID in service", zap.Uint("task_id", id))

	// 先尝试从缓存获取
	ctx := context.Background()
	cacheKey := fmt.Sprintf("tasks:%d", id)
	cachedData, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var task model.Task
		if err := json.Unmarshal([]byte(cachedData), &task); err == nil {
			logger.Info("Cache hit for task", zap.Uint("task_id", id))
			return &task, nil
		}
		logger.Warn("Failed to unmarshal cached task", zap.Uint("task_id", id), zap.Error(err))
	}

	// 缓存未命中，从数据库获取
	task, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get task from repository", zap.Uint("task_id", id), zap.Error(err))
		return nil, err
	}

	// 更新缓存
	s.createTaskCache(task)
	return task, nil
}

// invalidateTaskListCache 清除任务列表缓存
func (s *TaskService) invalidateTaskListCache() {
	ctx := context.Background()
	keys, err := s.redisClient.Keys(ctx, "tasks:page:*").Result()
	if err != nil {
		logger.Error("Failed to get cache keys for invalidation", zap.Error(err))
		return
	}
	if len(keys) > 0 {
		s.redisClient.Del(ctx, keys...)
		logger.Info("Task list cache invalidated", zap.Int("keys_deleted", len(keys)))
	}
}

// invalidateTaskCache 清除单个任务缓存
func (s *TaskService) invalidateTaskCache(id uint) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("tasks:%d", id)
	if err := s.redisClient.Del(ctx, cacheKey).Err(); err != nil {
		logger.Error("Failed to invalidate task cache", zap.Uint("task_id", id), zap.Error(err))
		return
	}
	logger.Info("Task cache invalidated", zap.Uint("task_id", id))
}

// createTaskCache 创建单个任务缓存
func (s *TaskService) createTaskCache(task *model.Task) {
	ctx := context.Background()
	cacheKey := fmt.Sprintf("tasks:%d", task.ID)
	data, err := json.Marshal(task)
	if err != nil {
		logger.Error("Failed to marshal task for cache", zap.Uint("task_id", task.ID), zap.Error(err))
		return
	}
	if err := s.redisClient.Set(ctx, cacheKey, data, 5*time.Minute).Err(); err != nil {
		logger.Error("Failed to set task cache", zap.Uint("task_id", task.ID), zap.Error(err))
		return
	}
	logger.Info("Task cache created", zap.Uint("task_id", task.ID))
}
