package model

import "gorm.io/gorm"

// RawData represents a raw data received from the device
type RawData struct {
	gorm.Model
	ConnString string `json:"conn_string" gorm:"not null;index"`
	DeviceID   uint   `json:"device_id" gorm:"not null;index"`
	Data       []byte `json:"data" gorm:"not null"`
}
