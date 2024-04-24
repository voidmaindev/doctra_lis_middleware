package model

import "gorm.io/gorm"

type DeviceModel struct {
	gorm.Model
	Name   string `json:"name" gorm:"not null;unique"`
	Driver string `json:"driver" gorm:"not null"`
}
