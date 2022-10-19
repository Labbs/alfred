package user

import (
	"github.com/labbs/alfred/pkg/exception"
)

type User struct {
	Id       string `gorm:"primaryKey" json:"id"`
	Username string `gorm:"index" json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type UserRepository interface {
	FindUserByUsername(username string) (User, *exception.AppError)
	UpdateUser(user User) *exception.AppError
}
