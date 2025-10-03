package domain

import "time"

type User struct {
	Username string
	Email    string
	Password string

	AvatarUrl string

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) TableName() string {
	return "user"
}

type UserPers interface {
	GetByUsername(username string) (User, error)
}
