package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/melnik-dev/go_todo_jwt/configs"
	"github.com/sirupsen/logrus"
)

type Db struct {
	*sqlx.DB
}

func InitPostgres(cfg *configs.Config, logger *logrus.Logger) (*Db, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.DBName, cfg.DB.Password, cfg.DB.SSLMode)
	db, err := sqlx.Open("postgres", connStr) // Connect
	fields := logrus.Fields{
		"db_host": cfg.DB.Host,
		"db_port": cfg.DB.Port,
		"db_name": cfg.DB.DBName,
	}
	if err != nil {
		logger.WithError(err).WithFields(fields).Error("Failed to open database connection")
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.WithError(err).WithFields(fields).Fatal("Failed to ping database")
		return nil, err
	}

	logger.WithFields(fields).Info("Database connected successfully")
	return &Db{db}, nil
}
