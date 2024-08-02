package driver_text_humalyzer_primus_human

import (
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString = ""
	rawDataEndString   = "TRANSFER FINISH"
)

// driver_text_humalyzer_primus_human is the driver for the "HumaLyzer Primus Human" laboratory device data format.
type driver_text_humalyzer_primus_human struct {
	log   *log.Logger
	store *store.Store
}

// NewDriver creates a new "Text HumaLyzer Primus Human" driver.
func NewDriver(logger *log.Logger, store *store.Store) *driver_text_humalyzer_primus_human {
	return &driver_text_humalyzer_primus_human{
		log:   logger,
		store: store,
	}
}

// Log returns the logger.
func (d *driver_text_humalyzer_primus_human) Log() *log.Logger {
	return d.log
}

// Store returns the store.
func (d *driver_text_humalyzer_primus_human) Store() *store.Store {
	return d.store
}

// RawDataStartString returns the start string of the raw data.
func (d *driver_text_humalyzer_primus_human) RawDataStartString() string {
	return rawDataStartString
}

// RawDataEndString returns the end string of the raw data.
func (d *driver_text_humalyzer_primus_human) RawDataEndString() string {
	return rawDataEndString
}

// DataToBeReplaced returns the data to be replaced.
func (d *driver_text_humalyzer_primus_human) DataToBeReplaced() map[string]string {
	return map[string]string{}
}

// SendSimpleACK sends an ACK message.
func (d *driver_text_humalyzer_primus_human) SendSimpleACK(conn net.Conn) error {
	return nil
}

// PostUnmarshalActions performs post-unmarshal actions.
func (d *driver_text_humalyzer_primus_human) PostUnmarshalActions(conn net.Conn, data map[string]interface{}) error {
	return nil
}
