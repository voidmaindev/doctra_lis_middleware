// Package model provides the models for the application.
package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Role constants.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

const defaultHashCost = 10

// User represents a user.
type User struct {
	gorm.Model
	Username string `json:"username" gorm:"not null;unique;index"`
	Password string `json:"-" gorm:"not null"`
	Role     string `json:"role,omitempty" gorm:"not null;index"`
}

// HashPassword hashes the password.
func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), defaultHashCost)
	if err != nil {
		return err
	}

	u.Password = string(hash)

	return nil
}

// CheckPassword checks the password.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// SetDefaultRole sets the default role.
func (u *User) SetDefaultRole() {
	if u.Role == "" {
		u.Role = RoleUser
	}
}

// IsAdmin checks if the user is an admin.
func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

// IsUser checks if the user is a user.
func (u *User) IsUser() bool {
	return u.Role == RoleUser
}
