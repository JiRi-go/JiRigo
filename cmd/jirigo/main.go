package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/jslee/JiRigo/pkg/infra/config"
	"github.com/jslee/JiRigo/pkg/infra/db"
	"github.com/jslee/JiRigo/pkg/infra/postgres"
	"github.com/jslee/JiRigo/pkg/services/migrations"
	_ "github.com/lib/pq" // PostgreSQL 드라이버 임포트 추가
)

// loadEnv .env 파일을 로드하는 함수
func loadEnv() {
	// 여러 경로에서 .env 파일 로드 시도
	envFiles := []string{
		".env",
		"../../.env",
		"../../../.env",
	}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err == nil {
			log.Printf("환경 변수 파일 로드됨: %s", file)
			return
		}
	}

	log.Println("경고: .env 파일을 찾을 수 없습니다. 환경 변수를 시스템 환경에서 읽습니다.")
}

// 데이터베이스 초기화 및 마이그레이션 수행 후 DB 인터페이스 반환
func initDatabase(dbConfig *config.DatabaseConfig) (db.DB, error) {
	// 임시 클라이언트로 데이터베이스 존재 여부 확인
	_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tempClient, err := postgres.NewClient(&config.DatabaseConfig{
		Host:     dbConfig.Host,
		Port:     dbConfig.Port,
		User:     dbConfig.User,
		Password: dbConfig.Password,
		DBName:   "postgres", // 기본 데이터베이스 연결
		SSLMode:  dbConfig.SSLMode,
	})
	if err != nil {
		return nil, fmt.Errorf("임시 PostgreSQL 클라이언트 생성 실패: %w", err)
	}
	defer tempClient.Close()

	// 마이그레이션 관리자 생성
	migrationManager := migrations.NewMigrationManager(tempClient.DB(), dbConfig)

	// 데이터베이스 존재 확인 및 생성
	if err := migrationManager.EnsureDatabaseExists(); err != nil {
		return nil, fmt.Errorf("데이터베이스 확인/생성 실패: %w", err)
	}

	// 실제 애플리케이션 데이터베이스로 연결하는 GORM DB 생성
	gormDB, err := db.NewPostgresDB(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("GORM DB 생성 실패: %w", err)
	}

	// 마이그레이션 실행
	// 기존 마이그레이션 시스템을 계속 사용하려면:
	pgClient, err := postgres.NewClient(dbConfig)
	if err != nil {
		gormDB.Close() // GORM DB 연결 닫기
		return nil, fmt.Errorf("PostgreSQL 클라이언트 생성 실패: %w", err)
	}

	migrationManager = migrations.NewMigrationManager(pgClient.DB(), dbConfig)
	if err := migrationManager.RunMigrations(); err != nil {
		pgClient.Close() // PostgreSQL 클라이언트 연결 닫기
		gormDB.Close()   // GORM DB 연결 닫기
		return nil, fmt.Errorf("마이그레이션 실행 실패: %w", err)
	}

	pgClient.Close() // 마이그레이션에만 사용한 클라이언트 연결 닫기

	log.Println("모든 마이그레이션이 성공적으로 실행되었습니다.")

	return gormDB, nil
}

func main() {
	// 환경 변수 로드
	loadEnv()

	// 설정 로드
	dbConfig := config.NewDatabaseConfig()

	// 데이터베이스 초기화 및 DB 인터페이스 획득
	gormDB, err := initDatabase(dbConfig)
	if err != nil {
		log.Fatalf("데이터베이스 초기화 실패: %v", err)
	}
	defer gormDB.Close()

	// 서비스 초기화
	// signinService := signin.NewService(gormDB)

	// API 서버 초기화 및 실행
	// TODO: 서버 시작 로직 추가

	log.Println("서버 시작됨")
}
