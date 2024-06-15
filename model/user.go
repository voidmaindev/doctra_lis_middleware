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
	Username string `json:"username" gorm:"not null;index"`
	Password []byte `json:"-" gorm:"not null"`
	Role     string `json:"role,omitempty" gorm:"not null;index"`
}

// AuthUser represents an authentication user.
type AuthUser struct {
	Username       string `json:"username"`
	Password       string `json:"password"`
	Role           string `json:"role,omitempty"`
	HashedPassword []byte `json:"-"`
}

// NewUserFromAuthUser creates a new user from the authentication user.
func NewUserFromAuthUser(authUser *AuthUser) (*User, error) {
	err := authUser.HashPassword()
	if err != nil {
		return nil, err
	}

	user := &User{
		Username: authUser.Username,
		Password: authUser.HashedPassword,
		Role:     authUser.Role,
	}

	return user, nil
}

// UpdateFromAuthUser updates the user from the authentication user.
func (u *User) UpdateFromAuthUser(authUser *AuthUser) error {
	err := authUser.HashPassword()
	if err != nil {
		return err
	}

	u.Username = authUser.Username
	u.Password = authUser.HashedPassword
	u.Role = authUser.Role

	if u.Role == "" {
		u.SetDefaultRole()
	}

	return nil
}

// HashPassword hashes the password.
func (u *AuthUser) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), defaultHashCost)
	if err != nil {
		return err
	}

	u.HashedPassword = hash

	return nil
}

// CheckPassword checks the password.
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(u.Password, []byte(password))
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
