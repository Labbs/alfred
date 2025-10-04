package persistence

import "gorm.io/gorm"

type sessionPers struct {
	db *gorm.DB
}

func NewSessionPers(db *gorm.DB) *sessionPers {
	return &sessionPers{db: db}
}
