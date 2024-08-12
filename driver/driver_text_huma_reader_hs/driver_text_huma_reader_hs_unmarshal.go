package driver_text_huma_reader_hs

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/voidmaindev/doctra_lis_middleware/model"
)

// Unmarshal unmarshals the raw data.
func (d *driver_text_huma_reader_hs) Unmarshal(rawData string) (labDatas []*model.LabData, additionalData map[string]interface{}, err error) {
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

	lines := strings.Split(rawData, "B,")
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			lines = lines[1:]
		} else {
			break
		}
	}

	completedDate, err := getCompleteDateForUnmarshalRawData()
	if err != nil {
		fmt.Println("failed to get completed date for unmarshalRawData")
		return labDatas, nil, err
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) < 9 {
			continue
		}

		barcode, err := getBarcodeForUnmarshalRawData(parts)
		if err != nil {
			fmt.Println("failed to get barcode for unmarshalRawData")
			return labDatas, nil, err
		}

		index, err := getIndexForUnmarshalRawData()
		if err != nil {
			fmt.Println("failed to get index for unmarshalRawData")
			return labDatas, nil, err
		}

		param, err := getParamForUnmarshalRawData(parts)
		if err != nil {
			fmt.Println("failed to get param for unmarshalRawData")
			return labDatas, nil, err
		}

		result, err := getResultForUnmarshalRawData(parts)
		if err != nil {
			fmt.Println("failed to get result for unmarshalRawData")
			return labDatas, nil, err
		}

		unit, err := getUnitForUnmarshalRawData()
		if err != nil {
			fmt.Println("failed to get unit for unmarshalRawData")
			return labDatas, nil, err
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

	return labDatas, additionalData, nil
}

// getBarcodeForUnmarshalRawData gets the barcode for unmarshalling the raw data.
func getBarcodeForUnmarshalRawData(parts []string) (string, error) {
	// barcode := parts[1] + strconv.Itoa(int(parts[2][0]-'A')) + parts[2][1:]
	barcode := parts[0] + parts[7]

	return barcode, nil
}

// getCompleteDateForUnmarshalRawData gets the completed date for unmarshalling the raw data.
func getCompleteDateForUnmarshalRawData() (time.Time, error) {
	completedDate := time.Now()

	return completedDate, nil
}

// getIndexForUnmarshalRawData gets the index for unmarshalling the raw data.
func getIndexForUnmarshalRawData() (uint, error) {
	index := uint(1)

	return index, nil
}

// getParamForUnmarshalRawData gets the param for unmarshalling the raw data.
func getParamForUnmarshalRawData(parts []string) (string, error) {
	param := parts[3]

	return param, nil
}

// getResultForUnmarshalRawData gets the result for unmarshalling the raw data.
func getResultForUnmarshalRawData(parts []string) (string, error) {
	result := parts[8]

	return result, nil
}

// getUnitForUnmarshalRawData gets the unit for unmarshalling the raw data.
func getUnitForUnmarshalRawData() (string, error) {
	unit := ""

	return unit, nil
}
