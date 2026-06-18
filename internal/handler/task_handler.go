package handler

import (
	"github/mouhe/todolist/internal/pkg/response"
	"github/mouhe/todolist/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(s *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: s}
}
func (h *TaskHandler) Create(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Desc     string `json:"description"`
		Priority int    `json:"priority"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "参数格式错误: "+err.Error())
		return
	}
	task, err := h.taskService.CreateTask(req.Title, req.Desc, req.Priority)
	if err != nil {
		response.Error(c, "创建失败，请稍后重试")
		return
	}
	response.Success(c, gin.H{"task": task})
}
func (h *TaskHandler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.Request.URL.Query().Get("page"))
	if err != nil {
		response.Error(c, "获取页数传参失败")
		return
	}
	size, err := strconv.Atoi(c.Request.URL.Query().Get("size"))
	if err != nil {
		response.Error(c, "获取页码范围传参失败")
		return
	}
	tasks, total, err := h.taskService.GetTasks(page, size)
	if err != nil {
		response.Error(c, "创建失败请稍后重试")
	}
	response.Success(c, gin.H{"task": tasks, "total": total})

}
func (h *TaskHandler) GetTaskById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, "id转换失败")
	}
	task, err := h.taskService.GetTaskById(uint(id))
	if err != nil {
		response.Error(c, "创建失败，请稍后重试")
		return
	}
	response.Success(c, gin.H{"task": task})
}

func (h *TaskHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, "id转换失败")
		return
	}
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Priority    int    `json:"priority"`
		Status      string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, "参数格式错误: "+err.Error())
		return
	}
	task, err := h.taskService.GetTaskById(uint(id))
	if err != nil {
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
	if err := h.taskService.Update(task); err != nil {
		response.Error(c, "更新失败，请稍后重试")
		return
	}
	response.Success(c, gin.H{"task": task})
}

func (h *TaskHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		response.Error(c, "id转换失败")
		return
	}
	if err := h.taskService.DeleteTask(id); err != nil {
		response.Error(c, "删除失败，请稍后重试")
		return
	}
	response.Success(c, gin.H{"message": "删除成功"})
}
