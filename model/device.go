package model

import "gorm.io/gorm"

type Device struct {
	gorm.Model
	Name          string `json:"name" gorm:"not null"`
	DeviceModelID uint   `json:"device_model_id" gorm:"not null;index"`
	NetAddress    string `json:"net_address" gorm:"not null;unique"`
}