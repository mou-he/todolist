package main

import (
	"github/mouhe/todolist/internal/config"
	"github/mouhe/todolist/internal/database"
	"github/mouhe/todolist/internal/handler"
	"github/mouhe/todolist/internal/repository"
	"github/mouhe/todolist/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	db, err := database.InitMySQL(cfg.Database.DSN)
	if err != nil {
		panic(err)
	}
	taskrepo := repository.NewTaskRepository(db)
	taskservice := service.NewTaskService(taskrepo)
	taskHandler := handler.NewTaskHandler(taskservice)
	router := gin.Default()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/tasks", taskHandler.List)
		v1.POST("/tasks", taskHandler.Create)
		v1.PUT("/tasks/:id", taskHandler.Update)
		v1.DELETE("/tasks/:id", taskHandler.Delete)
		v1.GET("/tasks/:id", taskHandler.GetTaskById)
	}
	router.Run(cfg.Server.Port)

}
