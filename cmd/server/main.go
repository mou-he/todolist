package main

import (
	"github/mouhe/todolist/internal/config"
	"github/mouhe/todolist/internal/database"
	"github/mouhe/todolist/internal/handler"
	"github/mouhe/todolist/internal/middleware"
	"github/mouhe/todolist/internal/pkg/logger"
	"github/mouhe/todolist/internal/repository"
	"github/mouhe/todolist/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	logger.InitLogger(logger.LogConfig{
		Level:      cfg.Log.Level,
		OutputPath: cfg.Log.OutputPath,
		Format:     cfg.Log.Format,
	})
	defer logger.Log.Sync() // 确保日志缓冲区的日志被写出
	db, err := database.InitMySQL(cfg.Database.DSN)
	if err != nil {
		logger.Log.Fatal("Failed to connect to database", zap.Error(err))
	}
	redisClient, err := database.InitRedis(cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	if err != nil {
		logger.Log.Fatal("Failed to connect to redis", zap.Error(err))
	}
	defer redisClient.Close()
	logger.Log.Info("Database connection established")
	taskrepo := repository.NewTaskRepository(db)
	taskservice := service.NewTaskService(taskrepo, redisClient)
	taskHandler := handler.NewTaskHandler(taskservice)
	router := gin.Default()
	v1 := router.Group("/api/v1")
	v1.Use(middleware.Logger())
	{
		v1.GET("/tasks", taskHandler.List)
		v1.POST("/tasks", taskHandler.Create)
		v1.PUT("/tasks/:id", taskHandler.Update)
		v1.DELETE("/tasks/:id", taskHandler.Delete)
		v1.GET("/tasks/:id", taskHandler.GetTaskById)
	}
	router.Run(cfg.Server.Port)

}
