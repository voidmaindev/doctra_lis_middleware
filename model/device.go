package model

import "gorm.io/gorm"

// Device represents a device.
type Device struct {
	gorm.Model
	Name          string      `json:"name" gorm:"not null"`
	DeviceModelID uint        `json:"device_model_id" gorm:"not null;index"`
	DeviceModel   DeviceModel `json:"device_model" gorm:"foreignKey:device_model_id"`
	Serial        string      `json:"serial"`
	NetAddress    string      `json:"net_address" gorm:"not null"`
}
