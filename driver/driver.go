// Package driver provides the interface for the driver of the laboratory device.
package driver

import (
	"fmt"
	"strings"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/model"
	"github.com/voidmaindev/doctra_lis_middleware/store"
	"github.com/voidmaindev/doctra_lis_middleware/tcp"
)

// Driver is the interface for the driver of the laboratory device.
type Driver interface {
	ProcessDeviceMessage([]byte, *tcp.ConnData, *model.Device) error
}

func NewDriver(driverName string, logger *log.Logger, store *store.Store) (Driver, error) {
	normalizedDriverName := normalizeDriverName(driverName)

	switch normalizedDriverName {
	case "hl7_231":
		return NewDriver_hl7_231(logger, store), nil
	}

	return nil, fmt.Errorf("unknown driver: %s", driverName)
}

func normalizeDriverName(driverName string) string {
	rv := driverName
	rv = strings.Trim(rv, " ")
	rv = strings.ToLower(rv)
	rv = strings.ReplaceAll(rv, ".", "")
	rv = strings.ReplaceAll(rv, " ", "_")

	return rv
}
