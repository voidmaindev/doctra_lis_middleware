package store

import (
	"errors"
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"gorm.io/gorm"
)

// DeviceStore is the store for the Device model.
type DeviceStore struct {
	db *gorm.DB
}

// NewDeviceStore creates a new DeviceStore.
func NewDeviceStore(db *gorm.DB) (*DeviceStore, error) {
	store := &DeviceStore{db: db}
	err := store.db.AutoMigrate(&model.Device{})
	if err != nil {
		return nil, errors.New("failed to migrate Device model")
	}

	return store, nil
}

// Create creates a new device.
func (s *DeviceStore) Create(device *model.Device) error {
	err := s.db.Create(device).Error
	if err != nil {
		return fmt.Errorf("failed to create device: %v", device.ID)
	}

	return nil
}

// GetByID gets a device by ID.
func (s *DeviceStore) GetByID(id uint) (*model.Device, error) {
	device := &model.Device{}
	err := s.db.Preload("DeviceModel").First(device, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get device by ID: %v", id)
	}

	return device, nil
}

// GetBySerial gets a device by serial.
func (s *DeviceStore) GetBySerial(serial string) (*model.Device, error) {
	device := &model.Device{}
	err := s.db.Preload("DeviceModel").Where("serial = ?", serial).First(device).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get device by serial: %v", serial)
	}

	return device, nil
}

// GetByNetworkAddress gets a device by network address.
func (s *DeviceStore) GetByNetAddress(networkAddress string) (*model.Device, error) {
	device := &model.Device{}
	err := s.db.Preload("DeviceModel").Where("net_address = ?", networkAddress).First(device).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get device by network address: %v", networkAddress)
	}

	return device, nil
}

// GetByModelID gets all devices by model ID.
func (s *DeviceStore) GetByModelID(modelID uint) ([]model.Device, error) {
	devices := []model.Device{}
	err := s.db.Where("model_id = ?", modelID).Preload("DeviceModel").Find(&devices).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get devices by model ID: %v", modelID)
	}

	return devices, nil
}

// GetAll gets all devices.
func (s *DeviceStore) GetAll() ([]model.Device, error) {
	devices := []model.Device{}
	err := s.db.Preload("DeviceModel").Find(&devices).Error
	if err != nil {
		return nil, errors.New("failed to get all devices")
	}

	return devices, nil
}

// Update updates a device.
func (s *DeviceStore) Update(device *model.Device) error {
	err := s.db.Save(device).Error
	if err != nil {
		return fmt.Errorf("failed to update device: %v", device.ID)
	}

	return nil
}

// Delete deletes a device.
func (s *DeviceStore) Delete(device *model.Device) error {
	err := s.db.Delete(device).Error
	if err != nil {
		return fmt.Errorf("failed to delete device: %v", device.ID)
	}

	return nil
}
