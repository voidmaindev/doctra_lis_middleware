package driver_text_combilyzer_13_Human

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/model"
	"github.com/voidmaindev/doctra_lis_middleware/store"
)

const (
	rawDataStartString  = 0x2
	rawDataEndString    = 0x3
	completedDateFormat = "20060102150405"
)

// Driver_text_Combilyzer_13_Human is the driver for the "HL7 2.3.1" laboratory device data format.
type Driver_text_Combilyzer_13_Human struct {
	log   *log.Logger
	store *store.Store
}

// NewDriver_text_Combilyzer_13_Human creates a new "Text Combilyzer 13 Human" driver.
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

// Unmarshal unmarshals the raw data.
func (d *Driver_text_Combilyzer_13_Human) Unmarshal(rawData string) (labDatas []*model.LabData, err error) {
	defer func() {
		if r := recover(); r != nil {
			if err, ok := r.(error); ok {
				d.log.Error("unknown error occurred while unmarshalling raw data: " + err.Error())
			} else {
				d.log.Error("unknown error occurred while unmarshalling raw data: " + fmt.Sprint(r))
			}
			labDatas = []*model.LabData{}
			err = errors.New("failed to unmarshal raw data")
		}
	}()

	lines := strings.Split(rawData, "\n")

	barcode, err := getBarcodeForUnmarshalRawData(lines)
	if err != nil {
		fmt.Println("failed to get barcode for unmarshalRawData")
		return labDatas, err
	}

	completedDate, err := getCompleteDateForUnmarshalRawData(lines)
	if err != nil {
		fmt.Println("failed to get completed date for unmarshalRawData")
		return labDatas, err
	}

	for i := 3; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		parts := strings.Fields(line)

		index, err := getIndexForUnmarshalRawData(i)
		if err != nil {
			fmt.Println("failed to get index for unmarshalRawData")
			return labDatas, err
		}

		param, err := getParamForUnmarshalRawData(parts)
		if err != nil {
			fmt.Println("failed to get param for unmarshalRawData")
			return labDatas, err
		}

		result, err := getResultForUnmarshalRawData(parts)
		if err != nil {
			fmt.Println("failed to get result for unmarshalRawData")
			return labDatas, err
		}

		unit, err := getUnitForUnmarshalRawData(parts)
		if err != nil {
			fmt.Println("failed to get unit for unmarshalRawData")
			return labDatas, err
		}

		labData := &model.LabData{
			Barcode:       barcode,
			Index:         index,
			Param:         param,
			Result:        result,
			Unit:          unit,
			CompletedDate: completedDate,
		}

		labDatas = append(labDatas, labData)
	}

	return labDatas, nil
}

// getBarcodeForUnmarshalRawData gets the barcode for unmarshalling the raw data.
func getBarcodeForUnmarshalRawData(lines []string) (string, error) {
	barcode := strings.TrimSpace(lines[1])

	return barcode, nil
}

// getCompleteDateForUnmarshalRawData gets the completed date for unmarshalling the raw data.
func getCompleteDateForUnmarshalRawData(lines []string) (time.Time, error) {
	dateLine := strings.TrimSpace(lines[0])
	dateString := strings.TrimPrefix(dateLine, "Date:")
	parsedDate, err := time.Parse(completedDateFormat, dateString)
	if err != nil {
		return time.Time{}, errors.New("failed to get completed date")
	}

	return parsedDate, nil
}

// getIndexForUnmarshalRawData gets the index for unmarshalling the raw data.
func getIndexForUnmarshalRawData(i int) (uint, error) {
	return uint(i - 2), nil
}

// getParamForUnmarshalRawData gets the param for unmarshalling the raw data.
func getParamForUnmarshalRawData(parts []string) (string, error) {
	return parts[0], nil
}

// getResultForUnmarshalRawData gets the result for unmarshalling the raw data.
func getResultForUnmarshalRawData(parts []string) (string, error) {
	result := strings.Join(parts[1:len(parts)-1], " ")

	return result, nil
}

// getUnitForUnmarshalRawData gets the unit for unmarshalling the raw data.
func getUnitForUnmarshalRawData(parts []string) (string, error) {
	unit := parts[len(parts)-1]

	return unit, nil
}
