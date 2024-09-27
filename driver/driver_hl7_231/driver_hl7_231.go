package driver_hl7_231

import (
	"fmt"
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/services"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString  = 0xb
	rawDataEndString    = 0x1c
	completedDateFormat = "20060102150405"
)

// Driver_hl7_231 is the driver for the "HL7 2.3.1" laboratory device data format.
type Driver_hl7_231 struct {
	log                *log.Logger
	store              *store.Store
	deviceQueryService *services.DeviceQueryService
}

// NewDriver creates a new "HL7 2.3.1" driver.
func NewDriver(logger *log.Logger, store *store.Store, deviceQueryService *services.DeviceQueryService) *Driver_hl7_231 {
	return &Driver_hl7_231{
		log:                logger,
		store:              store,
		deviceQueryService: deviceQueryService,
	}
}

// Log returns the logger.
func (d *Driver_hl7_231) Log() *log.Logger {
	return d.log
}

// Store returns the store.
func (d *Driver_hl7_231) Store() *store.Store {
	return d.store
}

// RawDataStartString returns the start string of the raw data.
func (d *Driver_hl7_231) RawDataStartString() string {
	return fmt.Sprintf("%c", rawDataStartString)
}

// RawDataEndString returns the end string of the raw data.
func (d *Driver_hl7_231) RawDataEndString() string {
	return fmt.Sprintf("%c", rawDataEndString)
}

// DataToBeReplaced returns the data to be replaced.
func (d *Driver_hl7_231) DataToBeReplaced() map[string]string {
	return map[string]string{"\\r": "\n"}
}

// SendSimpleACK sends an ACK message.
func (d *Driver_hl7_231) SendSimpleACK(conn net.Conn) error {
	return nil
}

// ReceivedSimpleACK checks if the message is a simple ACK message.
func (d *Driver_hl7_231) ReceivedSimpleACK(msg string) bool {
	return false
}

// PostUnmarshalActions performs the post-unmarshal actions.
func (d *Driver_hl7_231) PostUnmarshalActions(conn net.Conn, data map[string]interface{}) error {
	// Send ACK
	if ackMsg, ok := data["ACK"]; ok {
		_, err := conn.Write([]byte(ackMsg.(string)))
		if err != nil {
			return err
		}
	}

	return nil
}
