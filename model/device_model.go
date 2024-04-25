package model

import "gorm.io/gorm"

// DeviceModel represents a device model.
type DeviceModel struct {
	gorm.Model
	Name   string `json:"name" gorm:"not null;unique"`
	Driver string `json:"driver" gorm:"not null"`
}
