package models

import (
	"time"
)

type User struct {
	ID int `db:"id" json:"id"`
	UserInfo
}

type UserInfo struct {
	FirstName    string    `db:"first_name" json:"first_name" validate:"required,min=3,max=30"`
	LastName     string    `db:"last_name" json:"last_name" validate:"required,min=3,max=30"`
	Password     string    `db:"password" json:"password" validate:"required,min=3"`
	Email        string    `db:"email" json:"email" validate:"required,email"`
	Phone        string    `db:"phone" json:"phone" validate:"required,min=11,max=11"`
	Token        string    `db:"token" json:"token"`
	UserType     string    `db:"user_type" json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
	RefreshToken string    `db:"refresh_token" json:"refresh_token"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
	UserID       string    `db:"user_id" json:"user_id"`
}
