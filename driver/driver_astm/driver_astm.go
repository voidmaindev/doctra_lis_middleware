package driver_astm

import (
	"fmt"
	"net"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/services"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString  = 0x05
	rawDataEndString    = 0x04
	completedDateFormat = "20060102150405"
	ack                 = 0x06
	queryName           = "QRY"
	queryMessagesName   = "MSGS"
)

// Driver_astm is the driver for the "ASTM" laboratory device data format.
type Driver_astm struct {
	log                *log.Logger
	store              *store.Store
	deviceQueryService *services.DeviceQueryService
}

// NewDriver creates a new "ASTM" driver.
func NewDriver(logger *log.Logger, store *store.Store, deviceQueryService *services.DeviceQueryService) *Driver_astm {
	return &Driver_astm{
		log:                logger,
		store:              store,
		deviceQueryService: deviceQueryService,
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
	return map[string]string{"\r": "", "\n": ""}
}

// SendSimpleACK sends an ACK message.
func (d *Driver_astm) SendSimpleACK(conn net.Conn) error {
	_, err := conn.Write([]byte{ack})
	if err != nil {
		return err
	}

	return nil
}

// ReceivedSimpleACK checks if the message is an ACK message.
func (d *Driver_astm) ReceivedSimpleACK(msg string) bool {
	return msg == string(ack)
}

// PostUnmarshalActions performs the post-unmarshal actions.
func (d *Driver_astm) PostUnmarshalActions(conn net.Conn, data map[string]interface{}) error {
	err := d.doQuery(conn, data)
	if err != nil {
		d.log.Error("failed to do the query action")
		return err
	}

	return nil
}

// doQuery does the query action if the message is a query message.
func (d *Driver_astm) doQuery(conn net.Conn, data map[string]interface{}) error {
	query, ok := data[queryName]
	if !ok {
		return nil
	}

	barcode := query.(*Query).SampleID

	dataToReturn, err := d.deviceQueryService.Query(barcode)
	if err != nil {
		d.log.Error("failed to query the service")
		return err
	}

	queryMessages := data[queryMessagesName].([]Message)

	// Send ENQ (Enquiry)
	if _, err := conn.Write([]byte(fmt.Sprintf("%c", rawDataStartString))); err != nil {
		return err
	}

	// // Wait for ACK from the device
	// buf := make([]byte, 1)
	// if _, err := conn.Read(buf); err != nil {
	// 	return err
	// }
	// if string(buf) != ack {
	// 	return fmt.Errorf("expected ACK from device, got: %v", buf)
	// }

	// Generate messages based on queryMessages and dataToReturn
	formattedMessages := generateASTMMessagesFromQuery(queryMessages, dataToReturn)

	// Send the formatted messages over the connection
	for i, msg := range formattedMessages {
		msg = fmt.Sprintf("%d", i+1) + stx + msg
		// Send the message with STX and ETX framing
		if _, err := conn.Write([]byte(msg)); err != nil {
			return err
		}

		// // Wait for ACK from the device after each message
		// if _, err := conn.Read(buf); err != nil {
		// 	return err
		// }
		// if string(buf) != ack {
		// 	return fmt.Errorf("expected ACK after sending message, got: %v", buf)
		// }
	}

	// Finally, send EOT (End of Transmission)
	if _, err := conn.Write([]byte(fmt.Sprintf("%c", rawDataEndString))); err != nil {
		return err
	}

	return nil
}
