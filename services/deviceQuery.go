// Package services provides services for the application
package services

import (
	"errors"

	"github.com/go-resty/resty/v2"
)

type DeviceQueryService struct {
	client    *resty.Client
	queryHost string
	device    string
}

// RequestBody represents the structure of your request
type DeviceQueryRequestBody struct {
	Barcode    string `json:"barcode"`
	HardwareSN string `json:"hardware_sn"`
}

// Indicator represents each indicator in the response
type DeviceQueryIndicator struct {
	Indicator string `json:"indicator"`
	Material  int    `json:"material"`
	Dilution  string `json:"dilution"`
	Barcode   string `json:"barcode"`
}

// ResponseBody represents the structure of your response
type DeviceQueryResponseBody struct {
	Indicators []DeviceQueryIndicator `json:"indicators"`
}

// DeviceQueryResponce represents the structure of your response
type DeviceQueryDataToReturn struct {
	Param string
}

// ToReturn converts the response to the structure you want to return
func (r *DeviceQueryResponseBody) ToReturn() []DeviceQueryDataToReturn {
	var data []DeviceQueryDataToReturn
	for _, indicator := range r.Indicators {
		data = append(data, DeviceQueryDataToReturn{
			Param: indicator.Indicator,
		})
	}

	return data
}

// NewDeviceQueryService creates a new DeviceQueryService
func NewDeviceQueryService(queryHost, device string) *DeviceQueryService {
	return &DeviceQueryService{
		client:    resty.New(),
		queryHost: queryHost,
		device:    device,
	}
}

// Query queries the device
func (s *DeviceQueryService) Query(barcode string) ([]DeviceQueryDataToReturn, error) {
	reqBody := &DeviceQueryRequestBody{
		Barcode:    barcode,
		HardwareSN: s.device,
	}

	respBody := &DeviceQueryResponseBody{}

	resp, err := s.client.R().
		SetBody(reqBody).
		SetResult(respBody).
		Post(s.queryHost)
	if err != nil {
		return nil, errors.New("" + err.Error())
	}

	// Check for HTTP status code
	if resp.IsError() {
		return nil, errors.New("HTTP error: " + resp.Status())
	}

	return respBody.ToReturn(), nil
}
