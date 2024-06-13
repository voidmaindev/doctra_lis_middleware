package driver_hl7_231

import (
	"fmt"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString  = 0xb
	rawDataEndString    = 0x1c
	completedDateFormat = "20060102150405"
)

// Driver_hl7_231 is the driver for the "HL7 2.3.1" laboratory device data format.
type Driver_hl7_231 struct {
	log   *log.Logger
	store *store.Store
}

// NewDriver creates a new "HL7 2.3.1" driver.
func NewDriver(logger *log.Logger, store *store.Store) *Driver_hl7_231 {
	return &Driver_hl7_231{
		log:   logger,
		store: store,
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
