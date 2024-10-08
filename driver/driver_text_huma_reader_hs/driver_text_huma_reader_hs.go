package driver_text_huma_reader_hs

import (
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/services"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString = ""
	rawDataEndString   = "TRANSFER FINISH"
)

// driver_text_huma_reader_hs is the driver for the "HumaLyzer Primus Human" laboratory device data format.
type driver_text_huma_reader_hs struct {
	log                *log.Logger
	store              *store.Store
	deviceQueryService *services.DeviceQueryService
}

// NewDriver creates a new "Text HumaLyzer Primus Human" driver.
func NewDriver(logger *log.Logger, store *store.Store, deviceQueryService *services.DeviceQueryService) *driver_text_huma_reader_hs {
	return &driver_text_huma_reader_hs{
		log:                logger,
		store:              store,
		deviceQueryService: deviceQueryService,
	}
}

// Log returns the logger.
func (d *driver_text_huma_reader_hs) Log() *log.Logger {
	return d.log
}

// Store returns the store.
func (d *driver_text_huma_reader_hs) Store() *store.Store {
	return d.store
}

// RawDataStartString returns the start string of the raw data.
func (d *driver_text_huma_reader_hs) RawDataStartString() string {
	return rawDataStartString
}

// RawDataEndString returns the end string of the raw data.
func (d *driver_text_huma_reader_hs) RawDataEndString() string {
	return rawDataEndString
}

// DataToBeReplaced returns the data to be replaced.
func (d *driver_text_huma_reader_hs) DataToBeReplaced() map[string]string {
	return map[string]string{}
}

// SendSimpleACK sends an ACK message.
func (d *driver_text_huma_reader_hs) SendSimpleACK(conn net.Conn) error {
	return nil
}

// ReceivedSimpleACK checks if the message is an ACK message.
func (d *driver_text_huma_reader_hs) ReceivedSimpleACK(msg string) bool {
	return false
}

// PostUnmarshalActions performs post-unmarshal actions.
func (d *driver_text_huma_reader_hs) PostUnmarshalActions(conn net.Conn, data map[string]interface{}) error {
	return nil
}
