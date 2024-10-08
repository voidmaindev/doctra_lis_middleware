package driver_text_combilyzer_13_Human

import (
	"fmt"
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/services"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString  = 0x2
	rawDataEndString    = 0x3
	completedDateFormat = "02-01-2006 15:04"
)

// Driver_text_Combilyzer_13_Human is the driver for the "Combilyzer 13 Human" laboratory device data format.
type Driver_text_Combilyzer_13_Human struct {
	log                *log.Logger
	store              *store.Store
	deviceQueryService *services.DeviceQueryService
}

// NewDriver creates a new "Text Combilyzer 13 Human" driver.
func NewDriver(logger *log.Logger, store *store.Store, deviceQueryService *services.DeviceQueryService) *Driver_text_Combilyzer_13_Human {
	return &Driver_text_Combilyzer_13_Human{
		log:                logger,
		store:              store,
		deviceQueryService: deviceQueryService,
	}
}

// Log returns the logger.
func (d *Driver_text_Combilyzer_13_Human) Log() *log.Logger {
	return d.log
}

// Store returns the store.
func (d *Driver_text_Combilyzer_13_Human) Store() *store.Store {
	return d.store
}

// RawDataStartString returns the start string of the raw data.
func (d *Driver_text_Combilyzer_13_Human) RawDataStartString() string {
	return fmt.Sprintf("%c", rawDataStartString)
}

// RawDataEndString returns the end string of the raw data.
func (d *Driver_text_Combilyzer_13_Human) RawDataEndString() string {
	return fmt.Sprintf("%c", rawDataEndString)
}

// DataToBeReplaced returns the data to be replaced.
func (d *Driver_text_Combilyzer_13_Human) DataToBeReplaced() map[string]string {
	return map[string]string{}
}

// SendSimpleACK sends an ACK message.
func (d *Driver_text_Combilyzer_13_Human) SendSimpleACK(conn net.Conn) error {
	return nil
}

// ReceivedSimpleACK checks if the message is a simple ACK message.
func (d *Driver_text_Combilyzer_13_Human) ReceivedSimpleACK(msg string) bool {
	return false
}

// PostUnmarshalActions performs the post-unmarshal actions.
func (d *Driver_text_Combilyzer_13_Human) PostUnmarshalActions(conn net.Conn, data map[string]interface{}) error {
	return nil
}
