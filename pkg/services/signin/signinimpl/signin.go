package signinimpl

import (
	"github.com/jslee/JiRigo/pkg/infra/db"
	"github.com/jslee/JiRigo/pkg/services/signin"
)

type Service struct {
	store store
}

func ProvideService(db db.DB) (signin.Service, error) {
	return &Service{
		store: &gormStore{db: db, deletes: []string{}},
	}, nil
}
