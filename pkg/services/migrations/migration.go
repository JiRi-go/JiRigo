package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jslee/JiRigo/pkg/infra/config"
	. "github.com/jslee/JiRigo/pkg/infra/postgres"
)

// MigrationManager 데이터베이스 마이그레이션 관리자
type MigrationManager struct {
	client *Client
	cfg    *config.DatabaseConfig
}

// NewMigrationManager 새로운 마이그레이션 관리자 생성
func NewMigrationManager(client *Client, cfg *config.DatabaseConfig) *MigrationManager {
	return &MigrationManager{
		client: client,
		cfg:    cfg,
	}
}

// EnsureDatabaseExists 데이터베이스가 존재하는지 확인하고 없다면 생성
func (m *MigrationManager) EnsureDatabaseExists() error {
	// postgres 데이터베이스에 연결 (기본 데이터베이스)
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		m.cfg.Host, m.cfg.Port, m.cfg.User, m.cfg.Password, m.cfg.SSLMode,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("postgres 기본 DB 연결 실패: %w", err)
	}
	defer db.Close()

	// 데이터베이스 존재 여부 확인
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = db.QueryRow(query, m.cfg.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("데이터베이스 존재 확인 실패: %w", err)
	}

	// 데이터베이스가 없으면 생성
	if !exists {
		log.Printf("데이터베이스 '%s'가 존재하지 않습니다. 생성합니다...", m.cfg.DBName)

		// CREATE DATABASE는 파라미터화 쿼리를 지원하지 않으므로 문자열 포맷팅 사용
		createQuery := fmt.Sprintf("CREATE DATABASE %s", m.cfg.DBName)
		_, err = db.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("데이터베이스 생성 실패: %w", err)
		}
		log.Printf("데이터베이스 '%s'가 성공적으로 생성되었습니다.", m.cfg.DBName)
	} else {
		log.Printf("데이터베이스 '%s'가 이미 존재합니다.", m.cfg.DBName)
	}

	return nil
}

// RunMigrations 데이터베이스 마이그레이션 실행
func (m *MigrationManager) RunMigrations() error {
	db := m.client.DB()

	// 마이그레이션 이력 테이블 생성
	if err := m.createMigrationsTable(db); err != nil {
		return err
	}

	// 테이블이 존재하는지 확인 먼저 수행
	log.Println("테이블 존재 여부 확인 중...")

	// 사용자 테이블 존재 여부 확인
	usersTableExists, err := m.tableExists(db, "users")
	if err != nil {
		return fmt.Errorf("테이블 존재 여부 확인 실패: %w", err)
	}

	if !usersTableExists {
		log.Println("사용자 테이블이 존재하지 않습니다. 생성합니다...")
		// 사용자 테이블 마이그레이션 실행 - users_mig.go에 있는 구현 호출
		if err := m.migrateUsersTable(db); err != nil {
			return err
		}
	} else {
		log.Println("사용자 테이블이 이미 존재합니다.")
	}

	log.Println("모든 마이그레이션이 성공적으로 실행되었습니다.")
	return nil
}

// tableExists 테이블이 존재하는지 확인
func (m *MigrationManager) tableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	query := `
        SELECT EXISTS (
            SELECT FROM information_schema.tables 
            WHERE table_schema = 'public' 
            AND table_name = $1
        )
    `

	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s 테이블 존재 여부 확인 실패: %w", tableName, err)
	}

	return exists, nil
}

// createMigrationsTable 마이그레이션 이력을 저장할 테이블 생성
func (m *MigrationManager) createMigrationsTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS schema_migrations (
		version VARCHAR(255) PRIMARY KEY,
		applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("마이그레이션 테이블 생성 실패: %w", err)
	}

	return nil
}

// checkMigrationApplied 특정 마이그레이션이 적용되었는지 확인
func (m *MigrationManager) checkMigrationApplied(db *sql.DB, version string) (bool, error) {
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM schema_migrations WHERE version = $1)"

	err := db.QueryRow(query, version).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("마이그레이션 확인 실패: %w", err)
	}

	return exists, nil
}

// recordMigration 마이그레이션 실행 기록
func (m *MigrationManager) recordMigration(db *sql.DB, version string) error {
	query := "INSERT INTO schema_migrations (version) VALUES ($1)"

	_, err := db.Exec(query, version)
	if err != nil {
		return fmt.Errorf("마이그레이션 기록 실패: %w", err)
	}

	return nil
}
