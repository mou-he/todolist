package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	Log     *zap.Logger
	Sugared *zap.SugaredLogger

	once sync.Once // 确保初始化只执行一次
)

// LogConfig 日志配置结构
type LogConfig struct {
	Level      string `mapstructure:"level"`
	OutputPath string `mapstructure:"output_path"`
	Format     string `mapstructure:"format"`
}

// InitLogger 初始化 zap 日志
func InitLogger(cfg LogConfig) {
	once.Do(func() {
		// 1. 设置日志级别
		level := zap.NewAtomicLevel()
		switch cfg.Level {
		case "debug":
			level.SetLevel(zap.DebugLevel)
		case "info":
			level.SetLevel(zap.InfoLevel)
		case "warn":
			level.SetLevel(zap.WarnLevel)
		case "error":
			level.SetLevel(zap.ErrorLevel)
		default:
			level.SetLevel(zap.InfoLevel)
		}

		// 2. 设置输出路径
		var writer zapcore.WriteSyncer
		if cfg.OutputPath == "" || cfg.OutputPath == "stdout" {
			writer = zapcore.AddSync(os.Stdout)
		} else {
			// 在这里如果不是直接输出到控制台，则会在initFlieWriter中同时输出到文件和控制台，确保日志不会丢失
			// 使用 filepath.Dir 获取目录路径，避免硬编码
			dir := filepath.Dir(cfg.OutputPath)
			if dir != "." {
				err := os.MkdirAll(dir, os.ModePerm)
				if err != nil {
					fmt.Printf("Failed to create log directory %s: %v, falling back to stdout\n", dir, err)
					writer = zapcore.AddSync(os.Stdout)
				} else {
					initFileWriter(cfg.OutputPath, &writer)
				}
			} else {
				initFileWriter(cfg.OutputPath, &writer)
			}
		}

		// 3. 设置编码器
		var encoder zapcore.Encoder
		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 设置时间格式
		encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

		if cfg.Format == "json" {
			encoder = zapcore.NewJSONEncoder(encoderConfig)
		} else {
			encoder = zapcore.NewConsoleEncoder(encoderConfig)
		}

		// 4. 创建核心
		core := zapcore.NewCore(encoder, writer, level)

		// 5. 创建 logger（添加调用者信息）
		Log = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		Sugared = Log.Sugar()

		// 输出初始化日志
		Log.Info("Logger initialized",
			zap.String("level", cfg.Level),
			zap.String("output_path", cfg.OutputPath),
			zap.String("format", cfg.Format),
		)
	})
}

// initFileWriter 初始化文件写入器
func initFileWriter(outputPath string, writer *zapcore.WriteSyncer) {
	file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer func() {
		file.Close()
	}()
	if err != nil {
		fmt.Printf("Failed to open log file %s: %v, falling back to stdout\n", outputPath, err)
		*writer = zapcore.AddSync(os.Stdout)
		return
	}
	*writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(file))
	runtime.SetFinalizer(file, func(f *os.File) {
		f.Close()
	})
}

// Debug 调试级别日志
func Debug(msg string, fields ...zap.Field) {
	Log.Debug(msg, fields...)
}

// Info 信息级别日志
func Info(msg string, fields ...zap.Field) {
	Log.Info(msg, fields...)
}

// Warn 警告级别日志
func Warn(msg string, fields ...zap.Field) {
	Log.Warn(msg, fields...)
}

// Error 错误级别日志
func Error(msg string, fields ...zap.Field) {
	Log.Error(msg, fields...)
}

// Panic 恐慌级别日志（会触发 panic）
func Panic(msg string, fields ...zap.Field) {
	Log.Panic(msg, fields...)
}

// Fatal 致命级别日志（会调用 os.Exit(1)）
func Fatal(msg string, fields ...zap.Field) {
	Log.Fatal(msg, fields...)
}

// Debugf 格式化调试日志
func Debugf(format string, args ...interface{}) {
	Sugared.Debugf(format, args...)
}

// Infof 格式化信息日志
func Infof(format string, args ...interface{}) {
	Sugared.Infof(format, args...)
}

// Warnf 格式化警告日志
func Warnf(format string, args ...interface{}) {
	Sugared.Warnf(format, args...)
}

// Errorf 格式化错误日志
func Errorf(format string, args ...interface{}) {
	Sugared.Errorf(format, args...)
}

// Sync 刷新缓冲区
func Sync() error {
	return Log.Sync()
}
