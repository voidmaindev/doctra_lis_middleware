package model

import "gorm.io/gorm"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null;unique;index"`
	Password string `json:"password,omitempty" gorm:"not null"`
	Role     string `json:"role" gorm:"not null;index"`
}
