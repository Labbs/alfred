package user

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/labbs/alfred/pkg/exception"
)

type User struct {
	Id       string `gorm:"primaryKey" json:"id"`
	Username string `gorm:"index" json:"username"`
	Password string `json:"password,omitempty"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar" gorm:"default:''"`
	DarkMode string `json:"dark_mode" gorm:"default:'auto'"`

	Tokens []Token `json:"tokens" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Token struct {
	Id     string `gorm:"primaryKey" json:"id"`
	Name   string `json:"name"`
	UserId string `gorm:"index" json:"-"`

	Scope ScopeStruct `gorm:"type:longtext" json:"scope"`

	CreatedAt int64 `json:"created_at"`
}

type ScopeStruct []Scope

type Scope struct {
	Name string
}

func (sla *ScopeStruct) Scan(value interface{}) error {
	var skills []Scope
	err := json.Unmarshal([]byte(value.([]uint8)), &skills)
	if err != nil {
		return err
	}
	*sla = skills
	return nil
}

func (sla ScopeStruct) Value() (driver.Value, error) {
	val, err := json.Marshal(sla)
	return string(val), err
}

type UserRepository interface {
	FindUserByUsername(username string) (User, *exception.AppError)
	UpdateUser(user User) *exception.AppError
	CreateToken(token Token) *exception.AppError
	FindTokenById(id string) (Token, *exception.AppError)
	DeleteTokenById(id, userId string) *exception.AppError
}
