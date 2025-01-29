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
	return msg == fmt.Sprintf("%c", ack)
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
		d.log.Err(err, "failed to query the service")
		return err
	}

	queryMessages := data[queryMessagesName].([]Message)

	formattedMessages := generateASTMMessagesFromQuery(queryMessages, dataToReturn)

	err = SendToConn(conn, []byte{rawDataStartString})
	if err != nil {
		d.log.Err(err, "failed to send the start string")
		return err
	}

	err = getAckFromDevice(conn)
	if err != nil {
		d.log.Err(err, "failed to get an ACK from the device")
		return err
	}

	for i, msg := range formattedMessages {
		msgForChecksum := fmt.Sprintf("%d", i+1) + msg + cr + etx
		checkSum := calculateASTMChecksum(msgForChecksum)

		formattedMsg := stx + msgForChecksum + checkSum + cr + lf

		err := SendToConn(conn, []byte(formattedMsg))
		if err != nil {
			d.log.Err(err, fmt.Sprintf("failed to send the message: %v", msg))
			return err
		}

		err = getAckFromDevice(conn)
		if err != nil {
			d.log.Err(err, "failed to get an ACK from the device")
			return err
		}
	}

	err = SendToConn(conn, []byte{rawDataEndString})
	if err != nil {
		return err
	}

	return nil
}

// calculateASTMChecksum calculates the checksum for ASTM protocol
func calculateASTMChecksum(content string) string {
	sum := 0

	for _, char := range content {
		sum += int(char)
	}

	checksum := sum % 0x100

	return fmt.Sprintf("%02X", checksum)
}

// SendToConn sends the message to the connection.
func SendToConn(conn net.Conn, msg []byte) error {
	_, err := conn.Write(msg)
	if err != nil {
		return err
	}

	return nil
}

// getAckFromDevice waits for an ACK message from the device.
func getAckFromDevice(conn net.Conn) error {
	buf := make([]byte, 1<<5)

	n, err := conn.Read(buf)
	if err != nil {
		return err
	}

	if buf[0] != byte(ack) {
		return fmt.Errorf("expected ACK from device, got: \"%v\"", buf[:n])
	}

	return nil
}
