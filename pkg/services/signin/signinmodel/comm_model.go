package signinmodel

import "time"

type RoleType string

const (
	ROLE_ADMIN RoleType = "ADMIN"
	ROLE_USER  RoleType = "USER"
)

type BaseUser struct {
	Name     string   `gorm:"column:user_name" json:"name"`
	NickName string   `gorm:"column:nick_name" json:"nick_name"`
	Email    string   `gorm:"column:email" json:"email"`
	Password string   `gorm:"column:password" json:"-"`
	RoleType RoleType `gorm:"column:role_type" json:"role_type"`
}

type Users struct {
	BaseUser
	UID       string    `gorm:"column:uid" json:"uid"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type CreateUserCmd struct {
	BaseUser
	UID       string    `gorm:"column:uid" json:"uid"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

type UpdateUserCmd struct {
	BaseUser
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}
