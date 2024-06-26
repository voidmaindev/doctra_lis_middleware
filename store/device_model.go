package store

import (
	"errors"
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"gorm.io/gorm"
)

// DeviceModelStore is the store for the DeviceModel model.
type DeviceModelStore struct {
	db *gorm.DB
}

// NewDeviceModelStore creates a new DeviceModelStore.
func NewDeviceModelStore(db *gorm.DB) (*DeviceModelStore, error) {
	store := &DeviceModelStore{db: db}
	err := store.db.AutoMigrate(&model.DeviceModel{})
	if err != nil {
		return nil, errors.New("failed to migrate DeviceModel model")
	}

	return store, nil
}

// Create creates a new device model.
func (s *DeviceModelStore) Create(deviceModel *model.DeviceModel) error {
	err := s.db.Create(deviceModel).Error
	if err != nil {
		return fmt.Errorf("failed to create device model: %v", deviceModel.ID)
	}

	return nil
}

// GetByID gets a device model by ID.
func (s *DeviceModelStore) GetByID(id uint) (*model.DeviceModel, error) {
	deviceModel := &model.DeviceModel{}
	err := s.db.First(deviceModel, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get device model by ID: %v", id)
	}

	return deviceModel, nil
}

// GetByModel gets all device models.
func (s *DeviceModelStore) GetAll() ([]model.DeviceModel, error) {
	deviceModels := []model.DeviceModel{}
	err := s.db.Find(&deviceModels).Error
	if err != nil {
		return nil, errors.New("failed to get all device models")
	}

	return deviceModels, nil
}

// Update updates a device model.
func (s *DeviceModelStore) Update(deviceModel *model.DeviceModel) error {
	err := s.db.Save(deviceModel).Error
	if err != nil {
		return fmt.Errorf("failed to update device model: %v", deviceModel.ID)
	}

	return nil
}

// Delete deletes a device model.
func (s *DeviceModelStore) Delete(deviceModel *model.DeviceModel) error {
	err := s.db.Delete(deviceModel).Error
	if err != nil {
		return fmt.Errorf("failed to delete device model: %v", deviceModel.ID)
	}

	return nil
}
