// Package driver provides the interface for the driver of the laboratory device.
package driver

import (
	"fmt"
	"net"
	"strings"

	"github.com/voidmaindev/doctra_lis_middleware/driver/driver_astm"
	"github.com/voidmaindev/doctra_lis_middleware/driver/driver_hl7_231"
	"github.com/voidmaindev/doctra_lis_middleware/driver/driver_text_combilyzer_13_Human"
	"github.com/voidmaindev/doctra_lis_middleware/driver/driver_text_humalyzer_primus_human"
	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/model"
	"github.com/voidmaindev/doctra_lis_middleware/store"
	"github.com/voidmaindev/doctra_lis_middleware/tcp"
)

// Driver is the interface for the driver of the laboratory device.
type Driver interface {
	Log() *log.Logger
	Store() *store.Store
	RawDataStartString() string
	RawDataEndString() string
	DataToBeReplaced() map[string]string
	Unmarshal(string) ([]*model.LabData, map[string]interface{}, error)
	SendSimpleACK(net.Conn) error
	PostUnmarshalACtions(net.Conn, map[string]interface{}) error
}

// NewDriver creates a new driver.
func NewDriver(driverName string, logger *log.Logger, store *store.Store) (Driver, error) {
	normalizedDriverName := normalizeDriverName(driverName)

	switch normalizedDriverName {
	case "hl7231":
		return driver_hl7_231.NewDriver(logger, store), nil
	case "astm":
		return driver_astm.NewDriver(logger, store), nil
	case "textcombilyzer13human":
		return driver_text_combilyzer_13_Human.NewDriver(logger, store), nil
	case "texthumalyzerprimushuman":
		return driver_text_humalyzer_primus_human.NewDriver(logger, store), nil
	}

	return nil, fmt.Errorf("unknown driver: %s", driverName)
}

// normalizeDriverName normalizes the driver name.
func normalizeDriverName(driverName string) string {
	rv := driverName
	rv = strings.Trim(rv, " ")
	rv = strings.ToLower(rv)

	deletions := []string{".", " ", "_"}
	for _, deletion := range deletions {
		rv = strings.ReplaceAll(rv, deletion, "")
	}

	return rv
}

func Drivers() []string {
	return []string{
		"hl7_231",
		"astm",
		"text_combilyzer_13_human",
		"text_humalyzer_primus_human",
	}
}

// GetRawDatas gets the raw datas from the message.
func GetRawDatas(d Driver, msg string, prds *tcp.PrevData) []string {
	rawDatas := []string{}

	rawDataStartString := d.RawDataStartString()
	rawDataEndString := d.RawDataEndString()

	for len(msg) > 0 {
		if !prds.Started {
			startIndex := strings.Index(msg, d.RawDataStartString())
			if startIndex == -1 {
				break
			}
			prds.Started = true
			msg = msg[startIndex+len(rawDataStartString):]
		} else {
			endIndex := strings.Index(msg, d.RawDataEndString())
			if endIndex == -1 {
				prds.Data += msg
				break
			}
			rawData := prds.Data + msg[:endIndex]
			if len(rawData) > 0 {
				rawDatas = append(rawDatas, rawData)
			}
			msg = msg[endIndex+len(rawDataEndString):]
			prds.Data = ""
			prds.Started = false
		}
	}

	return rawDatas
}
