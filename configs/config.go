package configs

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	App  ConfApp  `mapstructure:"app"`
	HTTP ConfHTTP `mapstructure:"http"`
	DB   ConfDB   `mapstructure:"db"`
	JWT  ConfJWT  `mapstructure:"jwt"`
	Log  ConfLog  `mapstructure:"log"`
}

type ConfApp struct {
	Name    string `mapstructure:"name"`
	Env     string `mapstructure:"env"`
	Version string `mapstructure:"version"`
}

type ConfHTTP struct {
	Port         string        `mapstructure:"port"`
	Host         string        `mapstructure:"host"`
	ReadTimeout  time.Duration `mapstructure:"readTimeout"`
	WriteTimeout time.Duration `mapstructure:"writeTimeout"`
	IdleTimeout  time.Duration `mapstructure:"idleTimeout"`
}

type ConfDB struct {
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	Username string `mapstructure:"username"`
	DBName   string `mapstructure:"dbname"`
	SSLMode  string `mapstructure:"sslmode"`
	Password string `mapstructure:"password"`
}

type ConfJWT struct {
	Secret   string        `mapstructure:"secret"`
	TokenTTL time.Duration `mapstructure:"tokenTTL"`
}

type ConfLog struct {
	Level      string `mapstructure:"level"`
	Format     string `mapstructure:"format"`
	Output     string `mapstructure:"output"`
	FilePath   string `mapstructure:"filepath"`
	CallerSkip int    `mapstructure:"callerskip"`
	// Настройки ротации логов
	MaxSize    int  `mapstructure:"maxsize"`    // MB
	MaxBackups int  `mapstructure:"maxbackups"` // количество файлов
	MaxAge     int  `mapstructure:"maxage"`     // дни
	Compress   bool `mapstructure:"compress"`   // сжатие
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

	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	envBindings := map[string]string{
		"app.name":     "APP_NAME",
		"app.env":      "APP_ENV",
		"app.version":  "APP_VERSION",
		"http.host":    "HTTP_HOST",
		"http.port":    "HTTP_PORT",
		"db.host":      "DB_HOST",
		"db.port":      "DB_PORT",
		"db.username":  "DB_USERNAME",
		"db.password":  "DB_PASSWORD",
		"db.dbname":    "DB_NAME",
		"db.sslmode":   "DB_SSLMODE",
		"jwt.secret":   "JWT_SECRET",
		"log.level":    "LOG_LEVEL",
		"log.format":   "LOG_FORMAT",
		"log.output":   "LOG_OUTPUT",
		"log.filepath": "LOG_FILE_PATH",
	}

	for configKey, envKey := range envBindings {
		if err := viper.BindEnv(configKey, envKey); err != nil {
			return nil, fmt.Errorf("failed binding %s: %w", envKey, err)
		}
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	setDefaults(&cfg)

	return &cfg, nil
}

func validateConfig(cfg *Config) error {
	var errors []string

	if cfg.HTTP.Host == "" {
		errors = append(errors, "HTTP_HOST environment variable not set")
	}
	if cfg.HTTP.Port == "" {
		errors = append(errors, "HTTP_PORT environment variable not set")
	}

	if cfg.DB.Host == "" {
		errors = append(errors, "DB_HOST environment variable not set")
	}
	if cfg.DB.Port == "" {
		errors = append(errors, "DB_PORT environment variable not set")
	}
	if cfg.DB.Username == "" {
		errors = append(errors, "DB_USERNAME environment variable not set")
	}
	if cfg.DB.Password == "" {
		errors = append(errors, "DB_PASSWORD environment variable not set")
	}
	if cfg.DB.DBName == "" {
		errors = append(errors, "DB_NAME environment variable not set")
	}

	if cfg.JWT.Secret == "" {
		errors = append(errors, "JWT_SECRET environment variable not set")
	}

	if len(errors) > 0 {
		return fmt.Errorf("configuration errors: %s", strings.Join(errors, "; "))
	}

	return nil
}

func setDefaults(cfg *Config) {
	if cfg.App.Name == "" {
		cfg.App.Name = "todo_jwt"
	}
	if cfg.App.Env == "" {
		cfg.App.Env = "development"
	}

	if cfg.JWT.TokenTTL == 0 {
		cfg.JWT.TokenTTL = time.Hour
	}

	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Log.Format == "" {
		cfg.Log.Format = "text"
	}
	if cfg.Log.Output == "" {
		cfg.Log.Output = "stdout"
	}
}
