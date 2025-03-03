package config

import "os"

// GetEnv 환경변수 값 가져오기
func GetEnv(key string) string {
	return os.Getenv(key)
}

// GetEnvOrDefault 기본값이 있는 환경변수
func GetEnvOrDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
