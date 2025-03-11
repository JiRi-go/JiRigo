package db

import (
	"context"

	"gorm.io/gorm"
)

// GORM 기반 데이터베이스 작업을 위한 인터페이스
type DB interface {
	// 트랜잭션 내에서 함수를 실행
	WithTx(ctx context.Context, fn func(tx *gorm.DB) error) error

	// GORM DB 인스턴스를 반환
	DB() *gorm.DB

	// 주어진 컨텍스트와 함께 GORM DB 인스턴스를 반환
	WithContext(ctx context.Context) *gorm.DB

	// 데이터베이스 연결을 닫음
	Close() error
}
