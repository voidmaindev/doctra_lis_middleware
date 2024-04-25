package store

import (
	"errors"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"gorm.io/gorm"
)

type LabDataStore struct {
	db *gorm.DB
}

func NewLabDataStore(db *gorm.DB) (*LabDataStore, error) {
	store := &LabDataStore{db: db}
	err := store.db.AutoMigrate(&model.LabData{})
	if err != nil {
		return nil, errors.New("failed to migrate LabData model")
	}

	return store, nil
}
