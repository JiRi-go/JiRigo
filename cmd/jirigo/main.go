package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/jslee/JiRigo/pkg/infra/config"
	"github.com/jslee/JiRigo/pkg/infra/postgres"
	"github.com/jslee/JiRigo/pkg/services/migrations"
	_ "github.com/lib/pq" // PostgreSQL 드라이버 임포트 추가
)

// init 함수는 main 함수 실행 전에 자동으로 호출됨
func init() {
	// 여러 가능한 위치 시도
	err1 := godotenv.Load()                              // 현재 작업 디렉토리
	err2 := godotenv.Load("../../.env")                  // main.go 기준 상대 경로
	err3 := godotenv.Load("/home/jslee/dev/JiRigo/.env") // 절대 경로

	if err1 != nil && err2 != nil && err3 != nil {
		log.Printf("경고: .env 파일을 찾을 수 없습니다")
	}
}

func main() {
	// 설정 로드
	dbConfig := config.NewDatabaseConfig()

	// 데이터베이스 존재 여부 확인 및 생성 (임시 클라이언트 사용)
	log.Println("데이터베이스 확인 중...")
	tempClient, err := postgres.NewClient(&config.DatabaseConfig{
		Host:     dbConfig.Host,
		Port:     dbConfig.Port,
		User:     dbConfig.User,
		Password: dbConfig.Password,
		DBName:   "postgres", // 기본 데이터베이스 연결
		SSLMode:  dbConfig.SSLMode,
	})
	if err != nil {
		log.Fatalf("임시 PostgreSQL 클라이언트 생성 실패: %v", err)
	}

	// 마이그레이션 관리자 생성
	migrationManager := migrations.NewMigrationManager(tempClient, dbConfig)

	// 데이터베이스 존재 확인 및 생성
	if err := migrationManager.EnsureDatabaseExists(); err != nil {
		log.Fatalf("데이터베이스 확인/생성 실패: %v", err)
	}

	// 임시 클라이언트 닫기
	tempClient.Close()

	// 실제 애플리케이션 데이터베이스로 연결하는 클라이언트 생성
	pgClient, err := postgres.NewClient(dbConfig)
	if err != nil {
		log.Fatalf("PostgreSQL 클라이언트 생성 실패: %v", err)
	}
	defer pgClient.Close()

	// 마이그레이션 실행
	migrationManager = migrations.NewMigrationManager(pgClient, dbConfig)
	if err := migrationManager.RunMigrations(); err != nil {
		log.Fatalf("마이그레이션 실행 실패: %v", err)
	}

	// 저장소 및 서비스 초기화
	// userRepo := users.NewPostgresRepository(pgClient)
	// userService := users.NewService(userRepo)

	// API 서버 초기화 및 실행
	// ...

	log.Println("서버 시작됨")
}
