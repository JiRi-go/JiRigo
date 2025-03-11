package signinimpl

import (
	"context"

	"github.com/jslee/JiRigo/pkg/infra/db"
	"github.com/jslee/JiRigo/pkg/services/signin"
	"github.com/jslee/JiRigo/pkg/services/signin/signinmodel"
)

type Service struct {
	store store
}

func ProvideService(db db.DB) (signin.Service, error) {
	return &Service{
		store: &gormStore{db: db, deletes: []string{}},
	}, nil
}

// 모든 사용자 목록 조회
func (s *Service) GetUsers(ctx context.Context) ([]signinmodel.Users, error) {
	return s.store.GetUsers(ctx)
}

// userUID를 이용하여 단일 사용자 조회
func (s *Service) GetUserByUID(ctx context.Context, userUID string) (*signinmodel.Users, error) {
	return s.store.GetUserByUID(ctx, userUID)
}

// 사용자 생성
func (s *Service) CreateUser(ctx context.Context, cmd signinmodel.CreateUserCmd) error {
	return s.store.Create(ctx, cmd)
}

// 사용자 정보 수정
func (s *Service) UpdateUser(ctx context.Context, userUID string, cmd signinmodel.UpdateUserCmd) error {
	return s.store.Update(ctx, userUID, cmd)
}

// 사용자 삭제
func (s *Service) DeleteUser(ctx context.Context, userUID string) error {
	return s.store.Delete(ctx, userUID)
}
