package store

import (
	"github.com/voidmaindev/doctra_lis_middleware/config"
	"github.com/voidmaindev/doctra_lis_middleware/log"
	"gorm.io/gorm"
)

type Store struct {
	db *gorm.DB
}

func NewStore(log log.Logger) (*Store, error) {
	settings, err := config.ReadDBConfig()
	if err != nil {
		log.Error("Failed to read the DB config")
		return nil, err
	}

	config := &gorm.Config{}
	db, err := NewDB(settings, config)
	if err != nil {
		log.Err(err, "Failed to connect to DB")
		return nil, err
	}

	store := &Store{db: db}

	return store, nil
}
