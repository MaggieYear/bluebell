package models

import "time"

type User struct {
	ID         uint64    `db:"id"`
	UserID     uint64    `json:"user_id" db:"user_id"`
	Username   string    `json:"username" db:"username"`
	Password   string    `db:"password"`
	Email      string    `db:"email"`
	Gender     int8      `db:"gender"`
	CreateTime time.Time `db:"create_time"`
	UpdateTime time.Time `db:"update_time"`
	Token      string    `json:"token"`
}
