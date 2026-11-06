// Package driver_astm_wondfo_accre_8 is the driver for the Wondfo Accre 8 CLIA analyzer.
//
// The analyzer uploads results in an ASTM E1394-like format, but its firmware sends the frame
// control characters as literal text ("\x02", "\x03", "\x5e") and constant checksums, so the
// generic ASTM drivers cannot parse it. See docs/devices-jc.md.
package driver_astm_wondfo_accre_8

import (
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/services"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString = "\x05" // ENQ
	rawDataEndString   = "\x04" // EOT
	ack                = 0x06
)

// Driver_astm_wondfo_accre_8 is the driver for the "Wondfo Accre 8" laboratory device data format.
type Driver_astm_wondfo_accre_8 struct {
	log                *log.Logger
	store              *store.Store
	deviceQueryService *services.DeviceQueryService
}

// NewDriver creates a new "Wondfo Accre 8" driver.
func NewDriver(logger *log.Logger, store *store.Store, deviceQueryService *services.DeviceQueryService) *Driver_astm_wondfo_accre_8 {
	return &Driver_astm_wondfo_accre_8{
		log:                logger,
		store:              store,
		deviceQueryService: deviceQueryService,
	}
}

// Log returns the logger.
func (d *Driver_astm_wondfo_accre_8) Log() *log.Logger {
	return d.log
}

// Store returns the store.
func (d *Driver_astm_wondfo_accre_8) Store() *store.Store {
	return d.store
}

// RawDataStartString returns the start string of the raw data.
func (d *Driver_astm_wondfo_accre_8) RawDataStartString() string {
	return rawDataStartString
}

// RawDataEndString returns the end string of the raw data.
func (d *Driver_astm_wondfo_accre_8) RawDataEndString() string {
	return rawDataEndString
}

// DataToBeReplaced returns the data to be replaced.
// Nothing is replaced: the literal control characters are handled while unmarshalling,
// so the stored raw data stays identical to what the analyzer sent.
func (d *Driver_astm_wondfo_accre_8) DataToBeReplaced() map[string]string {
	return map[string]string{}
}

// SendSimpleACK sends an ACK message.
func (d *Driver_astm_wondfo_accre_8) SendSimpleACK(conn net.Conn) error {
	_, err := conn.Write([]byte{ack})
	if err != nil {
		return err
	}

	return nil
}

// ReceivedSimpleACK checks if the message is an ACK message.
func (d *Driver_astm_wondfo_accre_8) ReceivedSimpleACK(msg string) bool {
	return msg == string(rune(ack))
}

// PostUnmarshalActions performs the post-unmarshal actions.
// The Accre 8 only uploads results (no host query), so there is nothing to do.
func (d *Driver_astm_wondfo_accre_8) PostUnmarshalActions(conn net.Conn, data map[string]interface{}) error {
	return nil
}
