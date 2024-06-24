package driver_text_humalyzer_primus_human

import (
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/store"
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
	return ""
}

// RawDataEndString returns the end string of the raw data.
func (d *driver_text_humalyzer_primus_human) RawDataEndString() string {
	return "TRANSFER FINISH"
}

// DataToBeReplaced returns the data to be replaced.
func (d *driver_text_humalyzer_primus_human) DataToBeReplaced() map[string]string {
	return map[string]string{}
}

// SendACK sends an ACK message.
func (d *driver_text_humalyzer_primus_human) SendACK(conn net.Conn) error {
	return nil
}
