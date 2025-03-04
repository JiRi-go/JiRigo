package config

import "time"

// DatabaseConfig 데이터베이스 연결 설정
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string

	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// NewDatabaseConfig 새 데이터베이스 설정 생성
func NewDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Host:            GetEnvOrDefault("DB_HOST", "localhost"),
		Port:            GetEnvOrDefault("DB_PORT", "5432"),
		User:            GetEnvOrDefault("DB_USER", "postgres"),
		Password:        GetEnv("DB_PASSWORD"),
		DBName:          GetEnvOrDefault("DB_NAME", "girigo"),
		SSLMode:         GetEnvOrDefault("DB_SSL_MODE", "disable"),
		MaxOpenConns:    10,
		MaxIdleConns:    5,
		ConnMaxLifetime: time.Minute * 5,
	}
}
