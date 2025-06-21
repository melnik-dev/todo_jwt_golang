package logger

import (
	"fmt"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	globalLogger *logrus.Logger
	initOnce     sync.Once
)

func InitLogger(cfg configs.ConfLog) error {
	var err error
	initOnce.Do(func() {
		newLogger := logrus.New()

		// 1. Уровень логирования
		level, parseErr := logrus.ParseLevel(cfg.Level)
		if parseErr != nil {
			newLogger.SetLevel(logrus.InfoLevel)
			newLogger.Warnf("Invalid log level '%s', defaulting to 'info'. Error: %v", cfg.Level, parseErr)
		} else {
			newLogger.SetLevel(level)
		}

		// 2. Формат логирования
		switch strings.ToLower(cfg.Format) {
		case "json":
			newLogger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: "2006-01-02T15:04:05.000Z07:00",
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:  "timestamp",
					logrus.FieldKeyLevel: "level",
					logrus.FieldKeyMsg:   "message",
					logrus.FieldKeyFunc:  "caller",
					logrus.FieldKeyFile:  "file",
				},
			})
		default: // "text"
			newLogger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
				ForceColors:     true,
				DisableQuote:    true,
			})
		}

		// 3. Вывод логирования
		switch strings.ToLower(cfg.Output) {
		case "file":
			if cfg.FilePath == "" {
				err = fmt.Errorf("log file path is required when output is 'file'")
				return
			}

			lumberjackLogger, err := newLumberjack(cfg)
			if err != nil {
				return
			}

			newLogger.SetOutput(lumberjackLogger)
		case "both":
			if cfg.FilePath != "" {
				lumberjackLogger, err := newLumberjack(cfg)
				if err != nil {
					return
				}

				newLogger.SetOutput(io.MultiWriter(os.Stdout, lumberjackLogger))
			} else {
				newLogger.SetOutput(os.Stdout)
			}
		default:
			newLogger.SetOutput(os.Stdout)
		}

		// 4. Отметка места вызова
		newLogger.SetReportCaller(true)
		newLogger.AddHook(&callerHook{skip: cfg.CallerSkip})

		globalLogger = newLogger
	})
	return err
}

func newLumberjack(cfg configs.ConfLog) (*lumberjack.Logger, error) {
	dir := filepath.Dir(cfg.FilePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		err = fmt.Errorf("failed to create log directory %s: %w", dir, err)
		return nil, err
	}

	if err := checkFileWritable(cfg.FilePath); err != nil {
		err = fmt.Errorf("log file %s is not writable: %w", cfg.FilePath, err)
		return nil, err
	}

	if cfg.MaxSize == 0 {
		cfg.MaxSize = 100
	}
	if cfg.MaxBackups == 0 {
		cfg.MaxBackups = 5
	}
	if cfg.MaxAge == 0 {
		cfg.MaxAge = 30
	}

	return &lumberjack.Logger{
		Filename:   cfg.FilePath,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
	}, nil
}

func checkFileWritable(filename string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return file.Close()
}

func GetLogger() *logrus.Logger {
	if globalLogger == nil {
		panic("Logger not initialized. Call NewLogger first")
	}
	return globalLogger
}

func FromContext(c *gin.Context) *logrus.Entry {
	if l, exists := c.Get("logger"); exists {
		if entry, ok := l.(*logrus.Entry); ok {
			return entry
		}
	}
	globalLogger.Warn("Logger not found in Gin context, using main logger")
	return globalLogger.WithContext(c.Request.Context())
}

// callerHook - Hook для корректировки отчета о вызывающем объекте
type callerHook struct {
	skip int
}

// Fire корректирует поле caller, чтобы оно указывало на реальное место вызова в коде
func (hook *callerHook) Fire(entry *logrus.Entry) error {
	if entry.HasCaller() {
		entry.Caller.File = fmt.Sprintf("%s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line)
	}
	return nil
}

// Levels возвращает все уровни логирования, для которых должен срабатывать хук
func (hook *callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
