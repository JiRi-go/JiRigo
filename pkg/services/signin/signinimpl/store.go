package signinimpl

import "github.com/jslee/JiRigo/pkg/infra/db"

type store interface{}

type gormStore struct {
	db      db.DB
	deletes []string
}
