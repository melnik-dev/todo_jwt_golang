package logger

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Config struct {
	Level      string
	Format     string
	Output     string
	FilePath   string
	CallerSkip int
}

var (
	globalLogger *logrus.Logger
	initOnce     sync.Once
)

func NewLogger(cfg Config) (*logrus.Logger, error) {
	var err error
	initOnce.Do(func() {
		newLogger := logrus.New()

		// 1. Уровень
		level, parseErr := logrus.ParseLevel(cfg.Level)
		if parseErr != nil {
			newLogger.SetLevel(logrus.InfoLevel) // По умолчанию Info, если парсинг не удался
			newLogger.Warnf("Invalid log level '%s', defaulting to 'info'. Error: %v", cfg.Level, parseErr)
		} else {
			newLogger.SetLevel(level)
		}

		// 2. Формат
		switch strings.ToLower(cfg.Format) {
		case "json":
			newLogger.SetFormatter(&logrus.JSONFormatter{
				TimestampFormat: "2006-01-02 15:04:05",
				FieldMap: logrus.FieldMap{
					logrus.FieldKeyTime:  "timestamp",
					logrus.FieldKeyLevel: "level",
					logrus.FieldKeyMsg:   "message",
					logrus.FieldKeyFunc:  "caller", // Добавляем caller в JSON
				},
			})
		default: // "text"
			newLogger.SetFormatter(&logrus.TextFormatter{
				FullTimestamp:   true,
				TimestampFormat: "2006-01-02 15:04:05",
				ForceColors:     true,
			})
		}

		// 3. Вывод
		switch strings.ToLower(cfg.Output) {
		case "file":
			if cfg.FilePath == "" {
				err = fmt.Errorf("log file path is required when output is 'file'")
				return
			}
			dir := filepath.Dir(cfg.FilePath)
			if err = os.MkdirAll(dir, 0755); err != nil {
				err = fmt.Errorf("failed to create log directory %s: %w", dir, err)
				return
			}
			file, openErr := os.OpenFile(cfg.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if openErr != nil {
				err = fmt.Errorf("failed to open log file %s: %w", cfg.FilePath, openErr)
				return
			}
			// файл нужно закрыть при завершении работы приложения
			// defer file.Close() сработает при выходе из функции NewLogger
			// Используем Lumberjack для ротации логов или сигнал
			//lumberjackLogger := &lumberjack.Logger{
			//	Filename:   cfg.FilePath,
			//	MaxSize:    100, // MB
			//	MaxBackups: 5,   // количество старых файлов
			//	MaxAge:     30,   // дней
			//	Compress:   true, // сжимать старые файлы
			//}
			newLogger.SetOutput(io.MultiWriter(os.Stdout, file)) // Логи и в консоль, и в файл
		default:
			newLogger.SetOutput(os.Stdout)
		}

		// 4. Отметка места вызова файл:строка
		// true, чтобы видеть файл и номер строки, откуда был вызван лог
		newLogger.SetReportCaller(true)
		// Пропускаем фреймы вызова, чтобы получить реальное место вызова в коде а не внутренности Logrus
		newLogger.AddHook(&callerHook{skip: cfg.CallerSkip}) // Добавляем хук для корректного отображения caller

		globalLogger = newLogger
	})
	return globalLogger, err
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

// callerHook - Hook, позволяет корректировать отчет о вызывающем объекте
type callerHook struct {
	skip int
}

// Fire корректирует поле caller, чтобы оно указывало на реальное место вызова в коде
func (hook *callerHook) Fire(entry *logrus.Entry) error {
	if entry.HasCaller() {
		// +3 - глубина вызовов Logrus/Hook
		entry.Caller.File = fmt.Sprintf("%s:%d", filepath.Base(entry.Caller.File), entry.Caller.Line)
	}
	return nil
}

// Levels возвращает все уровни логирования, для которых должен срабатывать хук
func (hook *callerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
