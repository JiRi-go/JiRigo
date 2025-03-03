package postgres

import (
	"database/sql"
	"fmt"

	"github.com/jslee/JiRigo/pkg/infra/config"
)

// Client PostgreSQL 클라이언트
type Client struct {
	db  *sql.DB
	cfg *config.DatabaseConfig
}

// NewClient 새 PostgreSQL 클라이언트 생성
func NewClient(cfg *config.DatabaseConfig) (*Client, error) {
	connStr := BuildConnString(
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("postgres 연결 실패: %w", err)
	}

	// 연결 풀 설정
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	// 연결 테스트
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("postgres ping 실패: %w", err)
	}

	return &Client{
		db:  db,
		cfg: cfg,
	}, nil
}

// BuildConnString 연결 문자열 생성
func BuildConnString(host, port, user, password, dbname, sslmode string) string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode,
	)
}

// DB 데이터베이스 객체 반환
func (c *Client) DB() *sql.DB {
	return c.db
}

// Close 데이터베이스 연결 종료
func (c *Client) Close() error {
	return c.db.Close()
}
