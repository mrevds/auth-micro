package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server    ServerConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	RateLimit RateLimitConfig
}

type ServerConfig struct {
	GRPCPort string
	Host     string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	MaxConns int32
	MinConns int32
}

type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  time.Duration
	RefreshTokenDuration time.Duration
}

type RateLimitConfig struct {
	RequestsPerSecond int
}

// Load загружает конфигурацию из файла и переменных окружения
func Load() (*Config, error) {
	v := viper.New()

	// Настройка Viper
	v.SetConfigName("config")     // имя файла конфигурации (без расширения)
	v.SetConfigType("yaml")       // тип файла
	v.AddConfigPath(".")          // путь поиска в текущей директории
	v.AddConfigPath("./config")   // путь поиска в папке config
	v.AddConfigPath("/etc/auth/") // путь для production

	// Автоматическое чтение переменных окружения
	v.AutomaticEnv()

	// Привязка переменных окружения к ключам конфига
	v.BindEnv("server.grpc_port", "GRPC_PORT")
	v.BindEnv("server.host", "SERVER_HOST")

	v.BindEnv("database.host", "PG_HOST")
	v.BindEnv("database.port", "PG_PORT")
	v.BindEnv("database.user", "PG_USER")
	v.BindEnv("database.password", "PG_PASSWORD")
	v.BindEnv("database.dbname", "PG_DATABASE_NAME")
	v.BindEnv("database.sslmode", "PG_SSL_MODE")

	v.BindEnv("jwt.secret_key", "SECRET_KEY")

	// Значения по умолчанию
	v.SetDefault("server.grpc_port", "50051")
	v.SetDefault("server.host", "localhost")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", "54322")
	v.SetDefault("database.user", "auth_db_user")
	v.SetDefault("database.password", "auth_db_password")
	v.SetDefault("database.dbname", "auth_db")
	v.SetDefault("database.sslmode", "disable")
	v.SetDefault("database.max_conns", 25)
	v.SetDefault("database.min_conns", 5)

	v.SetDefault("jwt.secret_key", "s12dasd1a3s1d6as5d1a3s1d6as5d")
	v.SetDefault("jwt.access_token_duration", "15m")
	v.SetDefault("jwt.refresh_token_duration", "168h") // 7 дней

	v.SetDefault("rate_limit.requests_per_second", 100)

	// Попытка прочитать файл конфигурации (не критично если нет)
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	// Парсинг длительности токенов
	accessDuration, err := time.ParseDuration(v.GetString("jwt.access_token_duration"))
	if err != nil {
		accessDuration = 15 * time.Minute
	}

	refreshDuration, err := time.ParseDuration(v.GetString("jwt.refresh_token_duration"))
	if err != nil {
		refreshDuration = 7 * 24 * time.Hour
	}

	cfg := &Config{
		Server: ServerConfig{
			GRPCPort: v.GetString("server.grpc_port"),
			Host:     v.GetString("server.host"),
		},
		Database: DatabaseConfig{
			Host:     v.GetString("database.host"),
			Port:     v.GetString("database.port"),
			User:     v.GetString("database.user"),
			Password: v.GetString("database.password"),
			DBName:   v.GetString("database.dbname"),
			SSLMode:  v.GetString("database.sslmode"),
			MaxConns: v.GetInt32("database.max_conns"),
			MinConns: v.GetInt32("database.min_conns"),
		},
		JWT: JWTConfig{
			SecretKey:            v.GetString("jwt.secret_key"),
			AccessTokenDuration:  accessDuration,
			RefreshTokenDuration: refreshDuration,
		},
        RateLimit: RateLimitConfig{
            RequestsPerSecond: v.GetInt("rate_limit.requests_per_second"),
        },
	}

	return cfg, nil
}

// GetDSN возвращает строку подключения к БД
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"postgresql://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}
