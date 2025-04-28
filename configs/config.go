package configs

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"os"
	"time"
)

type Config struct {
	HTTP struct {
		Port         string        `mapstructure:"port"`
		ReadTimeout  time.Duration `mapstructure:"readTimeout"`
		WriteTimeout time.Duration `mapstructure:"writeTimeout"`
		IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
	} `mapstructure:"http"`
	DB struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Username string `mapstructure:"username"`
		DBName   string `mapstructure:"dbname"`
		SSLMode  string `mapstructure:"sslmode"`
		Password string `mapstructure:"password"` // не читаем из файла
	} `mapstructure:"db"`
	JWT struct {
		Secret string `mapstructure:"secret"`
	} `mapstructure:"jwt"`
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("failed to read config.yaml file: %w", err)
	}

	//viper.AutomaticEnv() // HTTP_PORT == http.port
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.BindEnv("db.host", "DB_HOST"); err != nil {
		return nil, fmt.Errorf("failed binding DB_HOST: %w", err)
	}
	if err := viper.BindEnv("db.password", "DB_PASSWORD"); err != nil {
		return nil, fmt.Errorf("failed binding DB_PASSWORD: %w", err)
	}
	if err := viper.BindEnv("jwt.secret", "JWT_SECRET"); err != nil {
		return nil, fmt.Errorf("failed binding JWT_SECRET: %w", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if cfg.DB.Host == "" {
		return nil, errors.New("DB_HOST environment variable not set")
	}
	if cfg.DB.Password == "" {
		return nil, errors.New("DB_PASSWORD environment variable not set")
	}
	if cfg.JWT.Secret == "" {
		return nil, errors.New("JWT_SECRET environment variable not set")
	}

	fmt.Println(cfg)
	return &cfg, nil
}
