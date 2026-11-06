// Package driver_text_dri_chem_nx600 is the driver for the Fujifilm DRI-CHEM NX600 dry chemistry analyzer.
//
// It handles the analyzer's native Type 1 / Type 2 host protocol: comma separated records framed as
// STX ... ETX followed by a one byte BCC. See docs/devices-jc.md.
package driver_text_dri_chem_nx600

import (
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/services"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString = "\x02" // STX
	rawDataEndString   = "\x03" // ETX
)

// Driver_text_dri_chem_nx600 is the driver for the "Fujifilm DRI-CHEM NX600" laboratory device data format.
type Driver_text_dri_chem_nx600 struct {
	log                *log.Logger
	store              *store.Store
	deviceQueryService *services.DeviceQueryService
}

// NewDriver creates a new "Fujifilm DRI-CHEM NX600" driver.
func NewDriver(logger *log.Logger, store *store.Store, deviceQueryService *services.DeviceQueryService) *Driver_text_dri_chem_nx600 {
	return &Driver_text_dri_chem_nx600{
		log:                logger,
		store:              store,
		deviceQueryService: deviceQueryService,
	}
}

// Log returns the logger.
func (d *Driver_text_dri_chem_nx600) Log() *log.Logger {
	return d.log
}

// Store returns the store.
func (d *Driver_text_dri_chem_nx600) Store() *store.Store {
	return d.store
}

// RawDataStartString returns the start string of the raw data.
func (d *Driver_text_dri_chem_nx600) RawDataStartString() string {
	return rawDataStartString
}

// RawDataEndString returns the end string of the raw data.
func (d *Driver_text_dri_chem_nx600) RawDataEndString() string {
	return rawDataEndString
}

// DataToBeReplaced returns the data to be replaced.
func (d *Driver_text_dri_chem_nx600) DataToBeReplaced() map[string]string {
	return map[string]string{}
}

// SendSimpleACK sends an ACK message.
// The protocol defines no reply to test results, so nothing is sent.
func (d *Driver_text_dri_chem_nx600) SendSimpleACK(conn net.Conn) error {
	return nil
}

// ReceivedSimpleACK checks if the message is an ACK message.
func (d *Driver_text_dri_chem_nx600) ReceivedSimpleACK(msg string) bool {
	return false
}

// PostUnmarshalActions performs the post-unmarshal actions.
func (d *Driver_text_dri_chem_nx600) PostUnmarshalActions(conn net.Conn, data map[string]interface{}) error {
	return nil
}
