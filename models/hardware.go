package models

import "gorm.io/gorm"

type Hardware struct {
	gorm.Model
	Name string
}