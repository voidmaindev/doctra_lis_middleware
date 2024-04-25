package store

import (
	"errors"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"gorm.io/gorm"
)

type DeviceModelStore struct {
	db *gorm.DB
}

func NewDeviceModelStore(db *gorm.DB) (*DeviceModelStore, error) {
	store := &DeviceModelStore{db: db}
	err := store.db.AutoMigrate(&model.DeviceModel{})
	if err != nil {
		return nil, errors.New("failed to migrate DeviceModel model")
	}

	return store, nil
}
