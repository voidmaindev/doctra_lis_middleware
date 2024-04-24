package model

import (
	"time"

	"gorm.io/gorm"
)

type LabData struct {
	gorm.Model
	DeviceID      uint      `json:"device_id" gorm:"not null;index"`
	Barcode       string    `json:"barcode" gorm:"not null"`
	Param         string    `json:"param" gorm:"not null"`
	Result        float64   `json:"result" gorm:"not null"`
	Unit          string    `json:"unit" gorm:"not null"`
	CompletedDate time.Time `json:"completed_date" gorm:"type:datetime;not null"`
}
