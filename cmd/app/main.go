package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/requestid"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/internal/auth"
	"github.com/melnik-dev/go_todo_jwt/internal/task"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/melnik-dev/go_todo_jwt/pkg/db"
	"github.com/melnik-dev/go_todo_jwt/pkg/logger"
	"github.com/melnik-dev/go_todo_jwt/pkg/middleware"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configs: %s", err)
	}

	mainLogger, err := logger.NewLogger(cfg.Log)
	if err != nil {
		log.Fatalf("Error initializing logger: %s", err)
	}

	mainLogger.Info("Application configuration loaded successfully")

	pgDB, err := db.InitPostgres(cfg, mainLogger)
	if err != nil {
		mainLogger.Fatalf("Error connecting to database: %s", err)
		log.Fatalf("Error connecting to database, %s", err)
	}

	route := gin.Default()
	route.Use(requestid.New())
	route.Use(middleware.Logger())
	route.GET("/ping", ping(mainLogger))

	// Repositories
	userRepo := user.NewRepository(pgDB)
	taskRepo := task.NewRepository(pgDB)
	// Services
	authService := auth.NewService(userRepo)
	taskService := task.NewService(taskRepo)
	// Handler
	auth.NewHandler(route, auth.HandlerDeps{
		AuthService: authService,
		Config:      cfg,
	})
	task.NewHandler(route, &task.HandlerDeps{
		TaskService: taskService,
		Config:      cfg,
	})

	serverAddress := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	srv := &http.Server{
		Addr:           serverAddress,
		Handler:        route,
		ReadTimeout:    cfg.HTTP.ReadTimeout,
		WriteTimeout:   cfg.HTTP.WriteTimeout,
		IdleTimeout:    cfg.HTTP.IdleTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB

	}
	// Запуск сервера
	go func() {
		mainLogger.Infof("Server starting on %s", serverAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			mainLogger.Fatalf("Server failed: %v", err)
		}
	}()

	// канал для сигналов Graceful Shutdown
	quit := make(chan os.Signal, 1)
	// перехват сигналов SIGINT (Ctrl C) и SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Блокируемся до получения сигнала

	mainLogger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// остановка сервера
	if err := srv.Shutdown(ctx); err != nil {
		mainLogger.Fatalf("Server forced to shutdown: %v", err)
	}

	mainLogger.Info("Server stopped")
}

func ping(logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		requestLogger, exists := c.Get("logger")
		if !exists {
			// Если middleware не был установлен, используем основной логгер
			requestLogger = logger
		}
		// Приводим к типу *logrus.Entry, чтобы использовать его методы
		logEntry := requestLogger.(*logrus.Entry)

		logEntry.Info("Ping endpoint")
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	}
}
