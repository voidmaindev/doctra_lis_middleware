// Package store provides the interface that defines the methods that a store should implement.
package store

import (
	"github.com/voidmaindev/doctra_lis_middleware/config"
	"github.com/voidmaindev/doctra_lis_middleware/log"
	"gorm.io/gorm"
)

// Store is the store for the application.
type Store struct {
	db               *gorm.DB
	UserStore        *UserStore
	DeviceModelStore *DeviceModelStore
	DeviceStore      *DeviceStore
	LabDataStore     *LabDataStore
}

// NewStore creates a new Store.
func NewStore(log *log.Logger) (*Store, error) {
	settings, err := config.ReadDBConfig()
	if err != nil {
		log.Error("failed to read the DB config")
		return nil, err
	}

	dbConfig := &gorm.Config{}
	db, err := NewDB(settings, dbConfig)
	if err != nil {
		log.Err(err, "failed to connect to DB")
		return nil, err
	}

	userStore, err := NewUserStore(db)
	if err != nil {
		log.Err(err, "failed to create UserStore")
		return nil, err
	}

	deviceModelStore, err := NewDeviceModelStore(db)
	if err != nil {
		log.Err(err, "failed to create DeviceModelStore")
		return nil, err
	}

	deviceStore, err := NewDeviceStore(db)
	if err != nil {
		log.Err(err, "failed to create DeviceStore")
		return nil, err
	}

	labDataStore, err := NewLabDataStore(db)
	if err != nil {
		log.Err(err, "failed to create LabDataStore")
		return nil, err
	}

	store := &Store{
		db:               db,
		UserStore:        userStore,
		DeviceModelStore: deviceModelStore,
		DeviceStore:      deviceStore,
		LabDataStore:     labDataStore,
	}

	return store, nil
}
