package database

import (
	"context"
	"errors"
	"github/mouhe/todolist/internal/pkg/logger"
	"sync"

	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
)

var (
	redisClient *redis.Client
	redisOnce   sync.Once
)

// InitRedis 初始化 Redis 客户端（线程安全）
func InitRedis(addr, password string, db int) (*redis.Client, error) {
	var errResult error

	redisOnce.Do(func() {
		logger.Info("Connecting to Redis...",
			zap.String("addr", addr),
			zap.Int("db", db),
		)

		client := redis.NewClient(&redis.Options{
			Addr:         addr,
			Password:     password,
			DB:           db,
			PoolSize:     10,
			MinIdleConns: 2,
			IdleTimeout:  5 * 60, // 5分钟
		})

		ctx := context.Background()
		_, err := client.Ping(ctx).Result()
		if err != nil {
			logger.Error("Failed to connect to Redis", zap.Error(err))
			errResult = err
			return
		}

		logger.Info("Redis connection established")
		redisClient = client
		errResult = nil
	})

	return redisClient, errResult
}

// GetRedisClient 获取 Redis 客户端（确保已初始化）
func GetRedisClient() (*redis.Client, error) {
	if redisClient == nil {
		return nil, errors.New("Redis client not initialized")
	}
	return redisClient, nil
}
