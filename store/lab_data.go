package store

import (
	"errors"
	"fmt"

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

func (s *LabDataStore) Create(labData *model.LabData) error {
	err := s.db.Create(labData).Error
	if err != nil {
		return fmt.Errorf("failed to create lab data for device: %v and barcode: %v", labData.DeviceID, labData.Barcode)
	}

	return nil
}

func (s *LabDataStore) GetByID(id uint) (*model.LabData, error) {
	labData := &model.LabData{}
	err := s.db.First(labData, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by ID: %v", id)
	}

	return labData, nil
}

func (s *LabDataStore) GetByBarcode(barcode string) ([]*model.LabData, error) {
	labData := []*model.LabData{}
	err := s.db.Where("barcode = ?", barcode).Find(&labData).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by barcode: %v", barcode)
	}

	return labData, nil
}

func (s *LabDataStore) GetByDeviceID(deviceID uint) ([]*model.LabData, error) {
	labData := []*model.LabData{}
	err := s.db.Where("device_id = ?", deviceID).Find(&labData).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by device ID: %v", deviceID)
	}

	return labData, nil
}

func (s *LabDataStore) GetByDeviceIDAndBarcode(deviceID uint, barcode string) (*model.LabData, error) {
	labData := &model.LabData{}
	err := s.db.Where("device_id = ? AND barcode = ?", deviceID, barcode).First(labData).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by device ID: %v and barcode: %v", deviceID, barcode)
	}

	return labData, nil
}

func (s *LabDataStore) GetAll() ([]*model.LabData, error) {
	labData := []*model.LabData{}
	err := s.db.Find(&labData).Error
	if err != nil {
		return nil, errors.New("failed to get all lab data")
	}

	return labData, nil
}

func (s *LabDataStore) Update(labData *model.LabData) error {
	err := s.db.Save(labData).Error
	if err != nil {
		return fmt.Errorf("failed to update lab data: %v", labData.ID)
	}

	return nil
}

func (s *LabDataStore) Delete(labData *model.LabData) error {
	err := s.db.Delete(labData).Error
	if err != nil {
		return fmt.Errorf("failed to delete lab data: %v", labData.ID)
	}

	return nil
}
