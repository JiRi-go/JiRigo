package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jslee/JiRigo/pkg/infra/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgresDB는 PostgreSQL용 GORM DB 인터페이스를 구현합니다
type PostgresDB struct {
	gormDB *gorm.DB
	sqlDB  *sql.DB
}

// NewPostgresDB는 새 PostgresDB 인스턴스를 생성합니다
func NewPostgresDB(cfg *config.DatabaseConfig) (*PostgresDB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	gormDB, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("GORM DB 인스턴스 생성 실패: %w", err)
	}

	// SQL DB 인스턴스 가져오기 (커넥션 풀 설정을 위해)
	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, fmt.Errorf("SQL DB 인스턴스 가져오기 실패: %w", err)
	}

	// 커넥션 풀 설정
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.ConnMaxLifetime)

	return &PostgresDB{
		gormDB: gormDB,
		sqlDB:  sqlDB,
	}, nil
}

// WithTx는 트랜잭션 내에서 함수를 실행합니다
func (p *PostgresDB) WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error {
	return p.gormDB.WithContext(ctx).Transaction(fn)
}

// DB는 GORM DB 인스턴스를 반환합니다
func (p *PostgresDB) DB() *gorm.DB {
	return p.gormDB
}

// WithContext는 주어진 컨텍스트와 함께 GORM DB 인스턴스를 반환합니다
func (p *PostgresDB) WithContext(ctx context.Context) *gorm.DB {
	return p.gormDB.WithContext(ctx)
}

// Close는 데이터베이스 연결을 닫습니다
func (p *PostgresDB) Close() error {
	return p.sqlDB.Close()
}
