package store

import (
	"errors"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"gorm.io/gorm"
)

type UserStore struct {
	db *gorm.DB
}

func NewUserStore(db *gorm.DB) (*UserStore, error) {
	store := &UserStore{db: db}
	err := store.db.AutoMigrate(&model.User{})
	if err != nil {
		return nil, errors.New("failed to migrate User model")
	}

	return store, nil
}
