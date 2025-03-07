package migrations

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/jslee/JiRigo/pkg/infra/config"
	"github.com/jslee/JiRigo/pkg/services/migrations/mig"
	"github.com/jslee/JiRigo/pkg/services/migrations/migrator"
)

// MigrationManager 마이그레이션 관리를 위한 구조체
type MigrationManager struct {
	client *sql.DB
	config *config.DatabaseConfig
}

// NewMigrationManager 새로운 마이그레이션 매니저 생성
func NewMigrationManager(client *sql.DB, cfg *config.DatabaseConfig) *MigrationManager {
	return &MigrationManager{
		client: client,
		config: cfg,
	}
}

// EnsureDatabaseExists 데이터베이스 존재 확인 및 생성
func (m *MigrationManager) EnsureDatabaseExists() error {
	// 기본 postgres 데이터베이스에 연결
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=postgres sslmode=%s",
		m.config.Host, m.config.Port, m.config.User, m.config.Password, m.config.SSLMode,
	)

	postgresDB, err := sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("postgres 기본 DB 연결 실패: %w", err)
	}
	defer postgresDB.Close()

	// 데이터베이스 존재 여부 확인
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM pg_database WHERE datname = $1)"
	err = postgresDB.QueryRow(query, m.config.DBName).Scan(&exists)
	if err != nil {
		return fmt.Errorf("데이터베이스 존재 확인 실패: %w", err)
	}

	// 데이터베이스가 없으면 생성
	if !exists {
		log.Printf("데이터베이스 '%s'가 존재하지 않습니다. 생성합니다...", m.config.DBName)

		// CREATE DATABASE는 파라미터화 쿼리를 지원하지 않으므로 문자열 포맷팅 사용
		createQuery := fmt.Sprintf("CREATE DATABASE %s", m.config.DBName)
		_, err = postgresDB.Exec(createQuery)
		if err != nil {
			return fmt.Errorf("데이터베이스 생성 실패: %w", err)
		}
		log.Printf("데이터베이스 '%s'가 성공적으로 생성되었습니다.", m.config.DBName)
	} else {
		log.Printf("데이터베이스 '%s'가 이미 존재합니다.", m.config.DBName)
	}

	return nil
}

// RunMigrations 모든 마이그레이션 실행
func (m *MigrationManager) RunMigrations() error {
	mg := migrator.NewMigrator(m.client)

	// 각 도메인별 마이그레이션 등록
	mig.AddUserMigrations(mg)
	mig.AddOauthAccountsMigrations(mg)
	mig.AddDiariesMigrations(mg)
	mig.AddDiaryImagesMigrations(mg)
	mig.AddDiaryCommentsMigrations(mg)
	mig.AddDiaryLikesMigrations(mg)

	// 마이그레이션 실행
	return mg.RunMigrations()
}
