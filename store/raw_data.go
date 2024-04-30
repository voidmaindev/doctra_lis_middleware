package store

import (
	"errors"
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"gorm.io/gorm"
)

// RawDataStore is the store for the RawData model.
type RawDataStore struct {
	db *gorm.DB
}

// NewRawDataStore creates a new RawDataStore.
func NewRawDataStore(db *gorm.DB) (*RawDataStore, error) {
	store := &RawDataStore{db: db}
	err := store.db.AutoMigrate(&model.RawData{})
	if err != nil {
		return nil, errors.New("failed to migrate RawData model")
	}

	return store, nil
}

// Create creates a new raw data.
func (s *RawDataStore) Create(rawData *model.RawData) error {
	err := s.db.Create(rawData).Error
	if err != nil {
		return errors.New("failed to create raw data")
	}

	return nil
}

// GetByID gets a raw data by ID.
func (s *RawDataStore) GetByID(id uint) (*model.RawData, error) {
	rawData := &model.RawData{}
	err := s.db.First(rawData, id).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get raw data by ID: %v", id)
	}

	return rawData, nil
}

// GetAll gets all raw data.
func (s *RawDataStore) GetAll() ([]*model.RawData, error) {
	rawData := []*model.RawData{}
	err := s.db.Find(&rawData).Error
	if err != nil {
		return nil, errors.New("failed to get all raw data")
	}

	return rawData, nil
}

// Update updates a raw data.
func (s *RawDataStore) Update(rawData *model.RawData) error {
	err := s.db.Save(rawData).Error
	if err != nil {
		return errors.New("failed to update raw data")
	}

	return nil
}

// Delete deletes a raw data by ID.
func (s *RawDataStore) Delete(id uint) error {
	err := s.db.Delete(&model.RawData{}, id).Error
	if err != nil {
		return fmt.Errorf("failed to delete raw data by ID: %v", id)
	}

	return nil
}
