package persistence

import (
	"github.com/labbs/alfred/domain"
	"gorm.io/gorm"
)

type userPers struct {
	Db *gorm.DB
}

func NewUserPers(db *gorm.DB) *userPers {
	return &userPers{Db: db}
}

func (u *userPers) GetByUsername(username string) (domain.User, error) {
	var user domain.User
	err := u.Db.Debug().Where("username = ?", username).First(&user).Error
	return user, err
}
