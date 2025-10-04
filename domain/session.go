package domain

import "time"

type Session struct {
	Id       string
	Username string

	User      User `gorm:"foreignKey:Username;references:Username"`
	UserAgent string
	IpAddress string
	ExpiresAt time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}

func (s *Session) TableName() string {
	return "session"
}

type SessionPers interface{}
