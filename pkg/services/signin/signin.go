package signin

import (
	"context"

	"github.com/jslee/JiRigo/pkg/services/signin/signinmodel"
)

type Service interface {
	GetUsers(ctx context.Context) ([]signinmodel.Users, error)
	GetUserByUID(ctx context.Context, userUID string) (*signinmodel.Users, error)
	CreateUser(ctx context.Context, cmd signinmodel.CreateUserCmd) error
	UpdateUser(ctx context.Context, userUID string, cmd signinmodel.UpdateUserCmd) error
	DeleteUser(ctx context.Context, userUID string) error
}
