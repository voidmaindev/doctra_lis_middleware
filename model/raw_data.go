package model

import "gorm.io/gorm"

// RawData represents a raw data received from the device
type RawData struct {
	gorm.Model
	ConnString string `json:"conn_string" gorm:"not null;index"`
	DeviceID   uint   `json:"device_id" gorm:"not null;index"`
	Data       []byte `json:"data" gorm:"not null"`
	Processed  bool   `json:"processed" gorm:"default:false"`
}

// RawDataApi represents a raw data API.
type RawDataApi struct {
	gorm.Model
	ConnString string `json:"conn_string"`
	DeviceID   uint   `json:"device_id"`
	Data       string `json:"data"`
	Processed  bool   `json:"processed"`
}

// NewRawDataApi creates a new raw data API.
func NewRawDataApi(rawData *RawData) *RawDataApi {
	return &RawDataApi{
		Model:      rawData.Model,
		ConnString: rawData.ConnString,
		DeviceID:   rawData.DeviceID,
		Data:       string(rawData.Data),
		Processed:  rawData.Processed,
	}
}
