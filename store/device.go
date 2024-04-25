package store

import (
	"errors"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"gorm.io/gorm"
)

type DeviceStore struct {
	db *gorm.DB
}

func NewDeviceStore(db *gorm.DB) (*DeviceStore, error) {
	store := &DeviceStore{db: db}
	err := store.db.AutoMigrate(&model.Device{})
	if err != nil {
		return nil, errors.New("failed to migrate Device model")
	}	

	return store, nil
}
