package user

import (
	"github.com/labbs/alfred/pkg/database"
	"github.com/labbs/alfred/pkg/exception"
)

type UserRepositoryDB struct {
	client database.DbConnection
}

func NewUserRepository() UserRepositoryDB {
	client := database.GetDbConnection()
	return UserRepositoryDB{client: client}
}

func (d UserRepositoryDB) FindUserByUsername(username string) (User, *exception.AppError) {
	u := User{}
	r := d.client.DB.Where("username = ?", username).First(&u)
	if r.Error != nil {
		return User{}, exception.NewUnexpectedError("unable to find user", r.Error)
	}
	return u, nil
}

func (d UserRepositoryDB) UpdateUser(user User) *exception.AppError {
	r := d.client.DB.Save(&user)
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to update user", r.Error)
	}
	return nil
}
