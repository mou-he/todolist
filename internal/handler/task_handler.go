package handler

import (
	"github/mouhe/todolist/internal/pkg/logger"
	"github/mouhe/todolist/internal/pkg/response"
	"github/mouhe/todolist/internal/service"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: s}
}
func (h *TaskHandler) Create(c *gin.Context) {
	var req struct {
		Title     string    `json:"title" binding:"required"`
		Desc      string    `json:"description"`
		Priority  int       `json:"priority"`
		Deadline  time.Time `json:"deadline"`
		CreatedAt time.Time `json:"created_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind request body", zap.Error(err))
		response.ValidationError(c, "参数格式错误: "+err.Error())
		return
	}
	if req.Deadline.IsZero() {
		req.Deadline = time.Now().Add(24 * time.Hour) // 默认截止时间为24小时后
	}
	logger.Info("Creating task", zap.String("title", req.Title), zap.Int("priority", req.Priority))
	task, err := h.taskService.CreateTask(req.Title, req.Desc, req.Priority, req.Deadline)
	if err != nil {
		logger.Error("Failed to create task", zap.Error(err))
		response.Error(c, "创建失败，请稍后重试")
		return
	}
	logger.Info("Task created successfully", zap.Uint("task_id", task.ID))
	response.Success(c, gin.H{"task": task})
}
func (h *TaskHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.Request.URL.Query().Get("page"))
	if err != nil {
		logger.Warn("Invalid page parameter", zap.Error(err))
		response.Error(c, "获取页数传参失败")
		return
	}
	size, err := strconv.Atoi(c.Request.URL.Query().Get("size"))
	if err != nil {
		logger.Warn("Invalid size parameter", zap.Error(err))
		response.Error(c, "获取页码范围传参失败")
		return
	}
	logger.Info("Listing tasks", zap.Int("page", page), zap.Int("size", size))
	tasks, total, err := h.taskService.GetTasks(page, size)
	if err != nil {
		logger.Error("Failed to list tasks", zap.Error(err))
		response.Error(c, "查询失败，请稍后重试")
		return
	}
	logger.Info("Tasks listed successfully", zap.Int64("total", total))
	response.Success(c, gin.H{"task": tasks, "total": total})
}
func (h *TaskHandler) GetTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid task ID", zap.Error(err))
		response.Error(c, "id转换失败")
		return
	}
	logger.Info("Getting task by ID", zap.Int("task_id", id))
	task, err := h.taskService.GetTaskById(uint(id))
	if err != nil {
		logger.Error("Failed to get task", zap.Int("task_id", id), zap.Error(err))
		response.Error(c, "查询失败，请稍后重试")
		return
	}
	response.Success(c, gin.H{"task": task})
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid task ID", zap.Error(err))
		response.Error(c, "id转换失败")
		return
	}
	var req struct {
		Title       string    `json:"title"`
		Description string    `json:"description"`
		Priority    int       `json:"priority"`
		Status      string    `json:"status"`
		Deadline    time.Time `json:"deadline"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warn("Failed to bind request body", zap.Error(err))
		response.ValidationError(c, "参数格式错误: "+err.Error())
		return
	}
	logger.Info("Updating task", zap.Int("task_id", id), zap.String("title", req.Title), zap.String("status", req.Status))
	task, err := h.taskService.GetTaskById(uint(id))
	if err != nil {
		logger.Error("Task not found", zap.Int("task_id", id), zap.Error(err))
		response.Error(c, "任务不存在")
		return
	}
	if req.Title != "" {
		task.Title = req.Title
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Priority != 0 {
		task.Priority = req.Priority
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.Deadline != (time.Time{}) {
		task.Deadline = req.Deadline
	}
	task.UpdatedAt = time.Now()
	if err := h.taskService.Update(task); err != nil {
		logger.Error("Failed to update task", zap.Int("task_id", id), zap.Error(err))
		response.Error(c, "更新失败，请稍后重试")
		return
	}
	logger.Info("Task updated successfully", zap.Int("task_id", id))
	response.Success(c, gin.H{"task": task})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Warn("Invalid task ID", zap.Error(err))
		response.Error(c, "id转换失败")
		return
	}
	logger.Info("Deleting task", zap.Int("task_id", id))
	if err := h.taskService.DeleteTask(id); err != nil {
		logger.Error("Failed to delete task", zap.Int("task_id", id), zap.Error(err))
		response.Error(c, "删除失败，请稍后重试")
		return
	}
	logger.Info("Task deleted successfully", zap.Int("task_id", id))
	response.Success(c, gin.H{"message": "删除成功"})
}
