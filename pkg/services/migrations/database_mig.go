package migrations

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/jslee/JiRigo/pkg/infra/config"
)

// 데이터베이스 생성 마이그레이션
type DatabaseMigration struct {
	cfg *config.DatabaseConfig
}

// 새 데이터베이스 마이그레이션 생성
func NewDatabaseMigration(cfg *config.DatabaseConfig) *DatabaseMigration {
	return &DatabaseMigration{
		cfg: cfg,
	}
}

// 마이그레이션 이름 반환
func (m *DatabaseMigration) Name() string {
	// 데이터베이스 이름을 사용하여 마이그레이션 이름 생성
	return fmt.Sprintf("%s 데이터베이스 생성", m.cfg.DBName)
}

// 마이그레이션 버전 반환
func (m *DatabaseMigration) Version() string {
	// 현재 시간을 "YYYYMMDD" 형식으로 포맷팅
	currentTime := time.Now().Format("20060102")
	return fmt.Sprintf("%s_create_database", currentTime)
}

// 마이그레이션 실행
func (m *DatabaseMigration) Run(db *sql.DB) error {
	// 기본 postgres 데이터베이스에 연결
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		m.cfg.Host, m.cfg.Port, m.cfg.User, m.cfg.Password, m.cfg.SSLMode,
	)

	postgresDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("postgres 기본 DB 연결 실패: %w", err)
	}
	defer postgresDB.Close()

	// 데이터베이스 존재 여부 확인
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = postgresDB.QueryRow(query, m.cfg.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("데이터베이스 존재 확인 실패: %w", err)
	}

	// 데이터베이스가 없으면 생성
	if !exists {
		log.Printf("데이터베이스 '%s'가 존재하지 않습니다. 생성합니다...", m.cfg.DBName)

		// CREATE DATABASE는 파라미터화 쿼리를 지원하지 않으므로 문자열 포맷팅 사용
		createQuery := fmt.Sprintf("CREATE DATABASE %s", m.cfg.DBName)
		_, err = postgresDB.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("데이터베이스 생성 실패: %w", err)
		}
		log.Printf("데이터베이스 '%s'가 성공적으로 생성되었습니다.", m.cfg.DBName)
	} else {
		log.Printf("데이터베이스 '%s'가 이미 존재합니다.", m.cfg.DBName)
	}

	return nil
}
