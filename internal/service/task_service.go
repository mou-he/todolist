package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github/mouhe/todolist/internal/database"
	"github/mouhe/todolist/internal/model"
	"github/mouhe/todolist/internal/pkg/logger"
	"github/mouhe/todolist/internal/repository"
	"time"

	"go.uber.org/zap"
)

type TaskService struct {
	repo repository.TaskRepository
}

func NewTaskService(repo repository.TaskRepository) *TaskService {
	return &TaskService{repo: repo}
}

// 创建任务（你可以在这里加业务校验，比如标题不能为空）
func (s *TaskService) CreateTask(title, desc string, priority int) (*model.Task, error) {
	logger.Debug("Creating task in service", zap.String("title", title))
	task := &model.Task{
		Title:       title,
		Description: desc,
		Priority:    priority,
		Status:      "pending",
	}

	err := s.repo.Create(task)
	if err != nil {
		logger.Error("Failed to create task in repository", zap.Error(err))
	}
	return task, err
}

// 获取任务列表（封装分页默认值）
// internal/service/task_service.go
func (s *TaskService) GetTasks(page, pageSize int) ([]model.Task, int64, error) {
	ctx := context.Background()

	// 先尝试从缓存获取
	cacheKey := fmt.Sprintf("tasks:page:%d:size:%d", page, pageSize)
	redisClient, err := database.GetRedisClient()
	if err != nil {
		return nil, 0, err
	}

	cachedData, err := redisClient.Get(ctx, cacheKey).Result()

	if err == nil {
		// 缓存命中，直接返回
		var tasks []model.Task
		err = json.Unmarshal([]byte(cachedData), &tasks)
		if err == nil {
			logger.Info("Cache hit for tasks", zap.String("key", cacheKey))
			return tasks, int64(len(tasks)), nil
		}
	}

	// 缓存未命中，从数据库获取
	tasks, total, err := s.repo.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	// 将结果存入缓存（有效期 5 分钟）
	data, _ := json.Marshal(tasks)
	redisClient.Set(ctx, cacheKey, data, 5*time.Minute)

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
	s.CreateTaskCache(task)
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

func (s *TaskService) GetTaskById(id uint) (*model.Task, error) {
	logger.Debug("Getting task by ID in service", zap.Uint("task_id", id))

	// 先尝试从缓存获取
	ctx := context.Background()
	cacheKey := fmt.Sprintf("tasks:%d", id)
	redisClient, err := database.GetRedisClient()
	if err == nil {
		cachedData, cacheErr := redisClient.Get(ctx, cacheKey).Result()
		if cacheErr == nil {
			var task model.Task
			if json.Unmarshal([]byte(cachedData), &task) == nil {
				logger.Info("Cache hit for task", zap.Uint("task_id", id))
				return &task, nil
			}
		}
	}

	// 缓存未命中，从数据库获取
	task, err := s.repo.GetByID(id)
	if err != nil {
		logger.Error("Failed to get task from repository", zap.Uint("task_id", id), zap.Error(err))
		return nil, err
	}

	// 更新缓存
	s.CreateTaskCache(task)
	return task, nil
}

func (s *TaskService) invalidateTaskListCache() {
	ctx := context.Background()
	redisClient, err := database.GetRedisClient()
	if err != nil {
		logger.Error("Failed to get Redis client for cache invalidation", zap.Error(err))
		return
	}
	keys, err := redisClient.Keys(ctx, "tasks:page:*").Result()
	if err != nil {
		logger.Error("Failed to get cache keys for invalidation", zap.Error(err))
		return
	}
	redisClient.Del(ctx, keys...)
	logger.Info("Cache invalidation completed", zap.Int("keys_deleted", len(keys)))
}
func (s *TaskService) invalidateTaskCache(id uint) {
	ctx := context.Background()
	redisClient, err := database.GetRedisClient()
	if err != nil {
		logger.Error("Failed to get Redis client for cache invalidation", zap.Error(err))
		return
	}
	cacheKey := fmt.Sprintf("tasks:%d", id)
	redisClient.Del(ctx, cacheKey)
	logger.Info("Cache invalidation completed for task", zap.Uint("task_id", id))
}
func (s *TaskService) CreateTaskCache(task *model.Task) {
	ctx := context.Background()
	redisClient, err := database.GetRedisClient()
	if err != nil {
		logger.Error("Failed to get Redis client for cache invalidation", zap.Error(err))
		return
	}
	cacheKey := fmt.Sprintf("tasks:%d", task.ID)
	data, _ := json.Marshal(task)
	redisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	logger.Info("Cache created for task", zap.Uint("task_id", task.ID))
}

func (s *TaskService) CreateTaskListCache(tasks []model.Task, page, pageSize int) {
	ctx := context.Background()
	redisClient, err := database.GetRedisClient()
	if err != nil {
		logger.Error("Failed to get Redis client for cache creation", zap.Error(err))
		return
	}
	cacheKey := fmt.Sprintf("tasks:page:%d:size:%d", page, pageSize)
	data, _ := json.Marshal(tasks)
	redisClient.Set(ctx, cacheKey, data, 5*time.Minute)
	logger.Info("Cache created for tasks", zap.Int("page", page), zap.Int("size", pageSize))
}
