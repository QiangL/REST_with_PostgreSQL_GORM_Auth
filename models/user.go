package models

import "time"

type Ae86user struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	IsAdmin    bool      `json:"is_admin"`
	ExpireDate time.Time `json:"expire_date"`
	RateLimit  int64     `json:"rate_limit"`
}
