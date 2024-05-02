package model

import (
	"time"

	"gorm.io/gorm"
)

// LabData represents a lab data received from the device
type LabData struct {
	gorm.Model
	RawDataID     uint      `json:"raw_data_id" gorm:"not null;index"`
	DeviceID      uint      `json:"device_id" gorm:"not null;index"`
	Device        Device    `json:"device" gorm:"foreignKey:device_id"`
	Barcode       string    `json:"barcode" gorm:"not null;index"`
	Index         uint      `json:"index" gorm:"not null"`
	Param         string    `json:"param" gorm:"not null"`
	Result        string    `json:"result" gorm:"not null"`
	Unit          string    `json:"unit" gorm:"not null"`
	CompletedDate time.Time `json:"completed_date" gorm:"type:datetime;not null"`
}
