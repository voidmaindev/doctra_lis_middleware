package driver_astm

import (
	"fmt"
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString  = 0x05
	rawDataEndString    = 0x04
	completedDateFormat = "20060102150405"
	ack                 = 0x06
)

// Driver_astm is the driver for the "ASTM" laboratory device data format.
type Driver_astm struct {
	log   *log.Logger
	store *store.Store
}

// NewDriver creates a new "ASTM" driver.
func NewDriver(logger *log.Logger, store *store.Store) *Driver_astm {
	return &Driver_astm{
		log:   logger,
		store: store,
	}
}

// Log returns the logger.
func (d *Driver_astm) Log() *log.Logger {
	return d.log
}

// Store returns the store.
func (d *Driver_astm) Store() *store.Store {
	return d.store
}

// RawDataStartString returns the start string of the raw data.
func (d *Driver_astm) RawDataStartString() string {
	return fmt.Sprintf("%c", rawDataStartString)
}

// RawDataEndString returns the end string of the raw data.
func (d *Driver_astm) RawDataEndString() string {
	return fmt.Sprintf("%c", rawDataEndString)
}

// DataToBeReplaced returns the data to be replaced.
func (d *Driver_astm) DataToBeReplaced() map[string]string {
	return map[string]string{}
}

// SendACK sends an ACK message.
func (d *Driver_astm) SendACK(conn net.Conn) error {
	_, err := conn.Write([]byte{ack})
	if err != nil {
		return err
	}

	return nil
}
