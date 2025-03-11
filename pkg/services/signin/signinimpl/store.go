package signinimpl

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jslee/JiRigo/pkg/infra/db"
	"github.com/jslee/JiRigo/pkg/services/signin/signinmodel"
	"gorm.io/gorm"
)

type store interface {
	GetUsers(ctx context.Context) ([]signinmodel.Users, error)
	GetUserByUID(ctx context.Context, userUID string) (*signinmodel.Users, error)
	Create(ctx context.Context, cmd signinmodel.CreateUserCmd) error
	Update(ctx context.Context, userUID string, cmd signinmodel.UpdateUserCmd) error
	Delete(ctx context.Context, userUID string) error
}

type gormStore struct {
	db      db.DB
	deletes []string
}

// TODO_JS : api 반환값에 맞춰 조회하는 데이터 수정 필요 (25-03-11)
// 전체 사용자 목록 조회
func (ss *gormStore) GetUsers(ctx context.Context) ([]signinmodel.Users, error) {
	var users []signinmodel.Users

	result := ss.db.WithContext(ctx).Find(&users)
	if result.Error != nil {
		return nil, fmt.Errorf("사용자 목록 조회 실패: %w", result.Error)
	}

	return users, nil
}

// userUID를 이용하여 단일 목록 조회
func (ss *gormStore) GetUserByUID(
	ctx context.Context,
	userUID string,
) (*signinmodel.Users, error) {
	var user signinmodel.Users

	result := ss.db.WithContext(ctx).First(&user, "uid = ?", userUID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("사용자를 찾을 수 없습니다")
		}
		return nil, fmt.Errorf("사용자 조회 실패: %w", result.Error)
	}

	return &user, nil
}

// 사용자 정보 생성
func (ss *gormStore) Create(
	ctx context.Context,
	cmd signinmodel.CreateUserCmd,
) error {
	return ss.db.WithTx(ctx, func(tx *gorm.DB) error {
		// 이메일 중복 확인
		if err := checkEmailExists(tx, cmd.Email); err != nil {
			return err
		}

		// 고유 UID 생성
		uid, err := ss.generateUniqueUID(tx)
		if err != nil {
			return err
		}

		// 사용자 정보 설정
		now := time.Now()
		cmd.CreatedAt = now
		cmd.UpdatedAt = now
		cmd.UID = uid

		// 사용자 생성
		return tx.Create(cmd).Error
	})
}

// 사용자 정보 수정
func (ss *gormStore) Update(
	ctx context.Context,
	userUID string,
	cmd signinmodel.UpdateUserCmd,
) error {
	cmd.UpdatedAt = time.Now()

	result := ss.db.WithContext(ctx).Model(&signinmodel.Users{}).
		Where("uid = ?", userUID).Updates(cmd)

	if result.Error != nil {
		return fmt.Errorf("비밀번호 업데이트 실패: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("사용자를 찾을 수 없습니다")
	}

	return nil
}

// 사용자 삭제 구현
func (ss *gormStore) Delete(ctx context.Context, userUID string) error {
	result := ss.db.WithContext(ctx).Delete(&signinmodel.Users{}, "uid = ?", userUID)

	if result.Error != nil {
		return fmt.Errorf("사용자 삭제 실패: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return errors.New("삭제할 사용자를 찾을 수 없습니다")
	}

	return nil
}

// ----------------------------

// 이메일이 존재하는지 확인하는 함수
func checkEmailExists(tx *gorm.DB, email string) error {
	var count int64
	if err := tx.Model(&signinmodel.Users{}).
		Where("email = ?", email).Count(&count).Error; err != nil {
		return fmt.Errorf("이메일 확인 실패: %w", err)
	}

	if count > 0 {
		return errors.New("이미 사용 중인 이메일입니다")
	}

	return nil
}

// 데이터베이스에서 고유한 UID를 생성하는 함수
func (ss *gormStore) generateUniqueUID(tx *gorm.DB) (string, error) {
	for {
		uid := ss.generateUID()

		var exists bool
		if err := tx.Model(&signinmodel.Users{}).Select("count(*) > 0").Where("uid = ?", uid).Scan(&exists).Error; err != nil {
			return "", fmt.Errorf("UID 확인 실패: %w", err)
		}

		if !exists {
			return uid, nil // 중복되지 않는 UID 반환
		}
	}
}

// generateUID는 새로운 UID를 생성합니다
func (s *gormStore) generateUID() string {
	return uuid.New().String() // 또는 다른 UID 생성 방식 사용
}
