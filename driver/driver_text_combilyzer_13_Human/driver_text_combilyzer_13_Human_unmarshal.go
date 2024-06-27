package driver_text_combilyzer_13_Human

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/voidmaindev/doctra_lis_middleware/model"
)

// Unmarshal unmarshals the raw data.
func (d *Driver_text_Combilyzer_13_Human) Unmarshal(rawData string) (labDatas []*model.LabData, additionalData map[string]interface{}, err error) {
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
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			lines = lines[1:]
		} else {
			break
		}
	}

	barcode, err := getBarcodeForUnmarshalRawData(lines)
	if err != nil {
		fmt.Println("failed to get barcode for unmarshalRawData")
		return labDatas, nil, err
	}

	completedDate, err := getCompleteDateForUnmarshalRawData(lines)
	if err != nil {
		fmt.Println("failed to get completed date for unmarshalRawData")
		return labDatas, nil, err
	}

	for i := 3; i < len(lines); i++ {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) == 2 && strings.HasSuffix(parts[1], "mg/mmol") {
			parts = []string{parts[0], strings.TrimSuffix(parts[1], "mg/mmol"), "mg/mmol"}
		}

		index, err := getIndexForUnmarshalRawData(i)
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

		unit, err := getUnitForUnmarshalRawData(parts)
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
func getBarcodeForUnmarshalRawData(lines []string) (string, error) {
	barcode := strings.TrimSpace(lines[2])
	barcode = strings.TrimPrefix(barcode, "ID:")

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
	param := parts[0]
	param = strings.TrimPrefix(param, "*")

	return param, nil
}

// getResultForUnmarshalRawData gets the result for unmarshalling the raw data.
func getResultForUnmarshalRawData(parts []string) (string, error) {
	result := ""

	switch len(parts) {
	case 1:
		result = ""
	case 2:
		result = parts[1]
	case 3:
		result = parts[1]
	case 4:
		result = parts[2]
	default:
		result = strings.Join(parts[1:len(parts)-1], " ")
	}

	result = strings.TrimPrefix(result, "Normal")

	return result, nil
}

// getUnitForUnmarshalRawData gets the unit for unmarshalling the raw data.
func getUnitForUnmarshalRawData(parts []string) (string, error) {
	unit := ""
	if len(parts) > 2 {
		unit = parts[len(parts)-1]
	}

	return unit, nil
}
