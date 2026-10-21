// Package services provides services for the application
package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/voidmaindev/doctra_lis_middleware/log"
)

const (
	contentTypeHeader   = "Content-Type"
	jsonContentType     = "application/json"
	maxLoggedBodyLength = 2048
)

type DeviceQueryService struct {
	client    *resty.Client
	log       *log.Logger
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
func NewDeviceQueryService(logger *log.Logger, queryHost, device string) *DeviceQueryService {
	return &DeviceQueryService{
		client:    resty.New(),
		log:       logger,
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

	err = checkJSONContentType(resp)
	if err != nil {
		return nil, err
	}

	if len(respBody.Indicators) == 0 {
		s.logEmptyResponse(reqBody, resp)
	}

	return respBody.ToReturn(), nil
}

// checkJSONContentType rejects successful responses that are not JSON, since resty
// silently skips unmarshalling them and the result would look like "no orders".
func checkJSONContentType(resp *resty.Response) error {
	contentType := resp.Header().Get(contentTypeHeader)
	if !strings.Contains(strings.ToLower(contentType), jsonContentType) {
		return fmt.Errorf("query host returned non-JSON content type %q with status %s", contentType, resp.Status())
	}

	return nil
}

// logEmptyResponse records what was sent and what came back when the query host returned no indicators.
func (s *DeviceQueryService) logEmptyResponse(reqBody *DeviceQueryRequestBody, resp *resty.Response) {
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		reqJSON = []byte(fmt.Sprintf("%+v", reqBody))
	}

	s.log.Warn(fmt.Sprintf("query host returned no indicators: url=%s request=%s status=%s content_type=%q body=%s",
		s.queryHost,
		reqJSON,
		resp.Status(),
		resp.Header().Get(contentTypeHeader),
		truncate(string(resp.Body()), maxLoggedBodyLength),
	))
}

// truncate shortens a string to at most maxLength characters.
func truncate(value string, maxLength int) string {
	if len(value) <= maxLength {
		return value
	}

	return value[:maxLength] + "..."
}
