package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/melnik-dev/go_todo_jwt/internal/auth"
	"github.com/melnik-dev/go_todo_jwt/internal/task"
	"github.com/melnik-dev/go_todo_jwt/internal/user"
	"github.com/melnik-dev/go_todo_jwt/pkg/db"
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

	pgDB, err := db.InitPostgres(db.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		DBName:   cfg.DB.DBName,
		SSLMode:  cfg.DB.SSLMode,
		Password: cfg.DB.Password,
	})
	if err != nil {
		log.Fatalf("Error connecting to database, %s", err)
	}
	route := gin.Default()
	route.GET("/ping", ping)

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
	task.NewHandler(route, task.HandlerDeps{
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
		log.Printf("Server starting on %s", serverAddress)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// канал для сигналов Graceful Shutdown
	quit := make(chan os.Signal, 1)
	// перехват сигналов SIGINT (Ctrl C) и SIGTERM
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // Блокируемся до получения сигнала

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// остановка сервера
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server stopped")
}

func ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
