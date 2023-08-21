package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email string `gorm:"unique"`
	Passowrd string
}

type UserResponse struct{
	gorm.Model
	Email string `gorm:"unique"`
}