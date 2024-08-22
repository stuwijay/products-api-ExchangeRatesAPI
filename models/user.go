package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Role     string `json:"role"`
}
