// Package model provides the models for the application.
package model

import "gorm.io/gorm"

// Role constants.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

// User represents a user.
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null;unique;index"`
	Password string `json:"password,omitempty" gorm:"not null"`
	Role     string `json:"role,omitempty" gorm:"not null;index"`
}
