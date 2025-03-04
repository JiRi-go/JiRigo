package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

// migrateUsersTable 사용자 테이블 생성 마이그레이션
func (m *MigrationManager) migrateUsersTable(db *sql.DB) error {
	version := "20240303_create_users_table"

	// 이미 적용된 마이그레이션인지 확인
	applied, err := m.checkMigrationApplied(db, version)
	if err != nil {
		return err
	}

	if applied {
		log.Printf("마이그레이션 '%s'는 이미 적용되었습니다.", version)
		return nil
	}

	log.Printf("마이그레이션 '%s' 실행 중...", version)

	// 트랜잭션 시작
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("트랜잭션 시작 실패: %w", err)
	}

	// 롤백 준비
	defer func() {
		if err != nil {
			tx.Rollback()
			log.Printf("마이그레이션 '%s' 실패, 롤백 완료", version)
		}
	}()

	// 사용자 테이블 생성
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS users (
		uid VARCHAR(50) PRIMARY KEY,
		name VARCHAR(122) NOT NULL,
		nick_name VARCHAR(122) NOT NULL,
		email VARCHAR(255) UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role VARCHAR(5) NOT NULL,
		created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
		updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
	)`

	_, err = tx.Exec(createTableQuery)
	if err != nil {
		return fmt.Errorf("사용자 테이블 생성 실패: %w", err)
	}

	// 인덱스 생성
	createIndexQuery := `CREATE INDEX IF NOT EXISTS idx_users_email ON users(email)`

	_, err = tx.Exec(createIndexQuery)
	if err != nil {
		return fmt.Errorf("사용자 이메일 인덱스 생성 실패: %w", err)
	}

	// 마이그레이션 기록
	_, err = tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
	if err != nil {
		return fmt.Errorf("마이그레이션 기록 실패: %w", err)
	}

	// 트랜잭션 커밋
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("트랜잭션 커밋 실패: %w", err)
	}

	log.Printf("마이그레이션 '%s' 성공적으로 적용되었습니다.", version)
	return nil
}
