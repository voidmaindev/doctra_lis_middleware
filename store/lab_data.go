package store

import (
	"errors"
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"gorm.io/gorm"
)

// LabDataStore is the store for the LabData model.
type LabDataStore struct {
	db *gorm.DB
}

// NewLabDataStore creates a new LabDataStore.
func NewLabDataStore(db *gorm.DB) (*LabDataStore, error) {
	store := &LabDataStore{db: db}
	err := store.db.AutoMigrate(&model.LabData{})
	if err != nil {
		return nil, errors.New("failed to migrate LabData model")
	}

	return store, nil
}

// Create creates a new lab data.
func (s *LabDataStore) Create(labData *model.LabData) error {
	err := s.db.Create(labData).Error
	if err != nil {
		return fmt.Errorf("failed to create lab data for device: %v and barcode: %v", labData.DeviceID, labData.Barcode)
	}

	return nil
}

// CreateOrUpdate creates or updates a lab data.
func (s *LabDataStore) CreateOrUpdate(labData *model.LabData) error {
	labDataOld, err := s.GetByDeviceIDAndBarcodeAndParam(labData.DeviceID, labData.Barcode, labData.Param)
	if err != nil {
		labDataOld.Result = labData.Result
		labDataOld.Unit = labData.Unit
		labDataOld.CompletedDate = labData.CompletedDate

		errUpd := s.Update(labDataOld)
		if errUpd != nil {
			return fmt.Errorf("failed to update lab data for device: %v and barcode: %v", labData.DeviceID, labData.Barcode)
		}

		return nil
	}

	err = s.db.Save(labData).Error
	if err != nil {
		return fmt.Errorf("failed to create or update lab data for device: %v and barcode: %v", labData.DeviceID, labData.Barcode)
	}

	return nil
}

// GetByID gets a lab data by ID.
func (s *LabDataStore) GetByID(id uint) (*model.LabData, error) {
	labData := &model.LabData{}
	err := s.db.Preload("Device").First(labData, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by ID: %v", id)
	}

	return labData, nil
}

// GetByBarcode gets a lab data by barcode.
func (s *LabDataStore) GetByBarcode(barcode string) ([]*model.LabData, error) {
	labData := []*model.LabData{}
	err := s.db.Preload("Device").Where("barcode = ?", barcode).Find(&labData).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by barcode: %v", barcode)
	}

	return labData, nil
}

// GetByDeviceID gets lab data by device ID.
func (s *LabDataStore) GetByDeviceID(deviceID uint) ([]*model.LabData, error) {
	labData := []*model.LabData{}
	err := s.db.Preload("Device").Where("device_id = ?", deviceID).Find(&labData).Order("ID").Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by device ID: %v", deviceID)
	}

	return labData, nil
}

// GetByDeviceIDAndBarcode gets a lab data by device ID and barcode.
func (s *LabDataStore) GetByDeviceIDAndBarcode(deviceID uint, barcode string) ([]*model.LabData, error) {
	labData := []*model.LabData{}
	err := s.db.Preload("Device").Where("device_id = ? AND barcode = ?", deviceID, barcode).Find(&labData).Order("ID").Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by device ID: %v and barcode: %v", deviceID, barcode)
	}

	return labData, nil
}

// GetByDeviceIDAndBarcodeAndParam gets a lab data by device ID, barcode and param.
func (s *LabDataStore) GetByDeviceIDAndBarcodeAndParam(deviceID uint, barcode, param string) (*model.LabData, error) {
	labData := &model.LabData{}
	err := s.db.Preload("Device").Where("device_id = ? AND barcode = ? AND param = ?", deviceID, barcode, param).First(labData).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by device ID: %v, barcode: %v and param: %v", deviceID, barcode, param)
	}

	return labData, nil
}

// GetBySerial gets a lab data by serial.
func (s *LabDataStore) GetBySerial(serial string) ([]*model.LabData, error) {
	labData := []*model.LabData{}
	err := s.db.Preload("Device").Where("serial = ?", serial).Find(&labData).Order("ID").Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by serial: %v", serial)
	}

	return labData, nil
}

// GetBySerialAndBarcode gets a lab data by serial and barcode.
func (s *LabDataStore) GetBySerialAndBarcode(serial, barcode string) ([]*model.LabData, error) {
	labData := []*model.LabData{}
	err := s.db.Preload("Device").Where("serial = ? AND barcode = ?", serial, barcode).Find(&labData).Order("ID").Error
	if err != nil {
		return nil, fmt.Errorf("failed to get lab data by serial: %v and barcode: %v", serial, barcode)
	}

	return labData, nil
}

// GetAll gets all lab data.
func (s *LabDataStore) GetAll() ([]*model.LabData, error) {
	labData := []*model.LabData{}
	err := s.db.Preload("Device").Find(&labData).Order("ID").Error
	if err != nil {
		return nil, errors.New("failed to get all lab data")
	}

	return labData, nil
}

// Update updates a lab data.
func (s *LabDataStore) Update(labData *model.LabData) error {
	err := s.db.Save(labData).Error
	if err != nil {
		return fmt.Errorf("failed to update lab data: %v", labData.ID)
	}

	return nil
}

// Delete deletes a lab data.
func (s *LabDataStore) Delete(labData *model.LabData) error {
	err := s.db.Delete(labData).Error
	if err != nil {
		return fmt.Errorf("failed to delete lab data: %v", labData.ID)
	}

	return nil
}
