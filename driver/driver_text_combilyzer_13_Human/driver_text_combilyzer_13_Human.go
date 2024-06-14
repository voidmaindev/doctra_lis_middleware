package driver_text_combilyzer_13_Human

import (
	"fmt"
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString  = 0x2
	rawDataEndString    = 0x3
	completedDateFormat = "02-01-2006 15:04"
)

// Driver_text_Combilyzer_13_Human is the driver for the "HL7 2.3.1" laboratory device data format.
type Driver_text_Combilyzer_13_Human struct {
	log   *log.Logger
	store *store.Store
}

// NewDriver creates a new "Text Combilyzer 13 Human" driver.
func NewDriver(logger *log.Logger, store *store.Store) *Driver_text_Combilyzer_13_Human {
	return &Driver_text_Combilyzer_13_Human{
		log:   logger,
		store: store,
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

// SendACK sends an ACK message.
func (d *Driver_text_Combilyzer_13_Human) SendACK(conn net.Conn) error {
	return nil
}
