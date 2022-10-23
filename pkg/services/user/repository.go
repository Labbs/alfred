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
	r := d.client.DB.Preload("Tokens").Where("username = ?", username).First(&u)
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

func (d UserRepositoryDB) CreateToken(token Token) *exception.AppError {
	r := d.client.DB.Create(&token)
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to create token", r.Error)
	}
	return nil
}

func (d UserRepositoryDB) FindTokenById(id string) (Token, *exception.AppError) {
	t := Token{}
	r := d.client.DB.Where("id = ?", id).First(&t)
	if r.Error != nil {
		return Token{}, exception.NewUnexpectedError("unable to find token", r.Error)
	}
	return t, nil
}

func (d UserRepositoryDB) DeleteTokenById(id, userId string) *exception.AppError {
	r := d.client.DB.Where("id = ? and user_id = ?", id, userId).Delete(&Token{})
	if r.Error != nil {
		return exception.NewUnexpectedError("unable to delete token", r.Error)
	}
	return nil
}
