package driver_text_dri_chem_nx600

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/voidmaindev/doctra_lis_middleware/model"
)

const (
	stx            = "\x02"
	fieldSeparator = ","

	commandResult   = "R"
	conditionNormal = "NORMAL"
	signEqual       = "="

	// Header of a test results record:
	// R,condition,date,time,sample No.,patient ID,patient name,species,sex,age,position,number of tests
	commandField     = 0
	conditionField   = 1
	dateField        = 2
	timeField        = 3
	sampleNoField    = 4
	patientIDField   = 5
	testCountField   = 11
	headerFieldCount = 12

	// Every test adds: name,sign,result+unit,dilution,lower limit,upper limit,warning
	testNameOffset   = 0
	testSignOffset   = 1
	testResultOffset = 2
	testFieldCount   = 7

	// The result field is fixed width: 9 characters of result followed by 6 characters of unit.
	resultWidth = 9

	// The date format follows the analyzer's regional setting.
	dateLayoutDayFirst = "02-01-2006 15:04"
	dateLayoutISO      = "2006-01-02 15:04"
)

var (
	errNoBarcode       = errors.New("test results record without a sample ID and a sample No.")
	errUnmarshalFailed = errors.New("failed to unmarshal raw data")
)

// resultsHeader is the part of a test results record that is shared by all its tests.
type resultsHeader struct {
	barcode       string
	completedDate time.Time
	testCount     int
}

// Unmarshal unmarshals the raw data (everything between STX and ETX) to the lab data.
// Only test results of patient samples produce lab data; test start info, error info,
// worklist requests and control (QC) measurements are accepted and ignored.
func (d *Driver_text_dri_chem_nx600) Unmarshal(rawData string) (labDatas []*model.LabData, additionalData map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			d.log.Error("unknown error occurred while unmarshalling raw data: " + fmt.Sprint(r))
			labDatas = []*model.LabData{}
			err = errUnmarshalFailed
		}
	}()

	fields := splitFields(rawData)
	if !isPatientResult(fields) {
		return nil, nil, nil
	}

	labDatas, err = parseResults(fields)
	if err != nil {
		d.log.Err(err, "failed to parse the test results record")
		return nil, nil, err
	}

	return labDatas, nil, nil
}

// splitFields splits a record into its fields.
// The BCC byte that follows ETX can itself be an STX, which makes the framing hand over
// a record with a leading STX; everything up to the last STX is dropped.
func splitFields(rawData string) []string {
	record := rawData[strings.LastIndex(rawData, stx)+1:]

	return strings.Split(record, fieldSeparator)
}

// isPatientResult checks if the record holds the test results of a patient sample.
func isPatientResult(fields []string) bool {
	return fieldAt(fields, commandField) == commandResult && fieldAt(fields, conditionField) == conditionNormal
}

// parseResults converts a test results record to lab datas, one per test.
func parseResults(fields []string) ([]*model.LabData, error) {
	header, err := parseHeader(fields)
	if err != nil {
		return nil, err
	}

	labDatas := make([]*model.LabData, 0, header.testCount)
	index := uint(0)

	for testNumber := 0; testNumber < header.testCount; testNumber++ {
		offset := headerFieldCount + testNumber*testFieldCount

		labData := testLabData(fields[offset : offset+testFieldCount])
		if labData == nil {
			continue
		}

		index++
		labData.Index = index
		labData.Barcode = header.barcode
		labData.CompletedDate = header.completedDate
		labDatas = append(labDatas, labData)
	}

	return labDatas, nil
}

// parseHeader reads the sample identifier, the test time and the number of tests,
// and checks that the record is long enough to hold that many tests.
func parseHeader(fields []string) (resultsHeader, error) {
	if len(fields) < headerFieldCount {
		return resultsHeader{}, fmt.Errorf("test results header has %d fields, want %d", len(fields), headerFieldCount)
	}

	testCount, err := strconv.Atoi(fieldAt(fields, testCountField))
	if err != nil || testCount < 0 {
		return resultsHeader{}, fmt.Errorf("invalid number of tests %q", fieldAt(fields, testCountField))
	}

	wantFields := headerFieldCount + testCount*testFieldCount
	if len(fields) < wantFields {
		return resultsHeader{}, fmt.Errorf("test results record has %d fields, want %d for %d tests", len(fields), wantFields, testCount)
	}

	barcode := sampleBarcode(fields)
	if barcode == "" {
		return resultsHeader{}, errNoBarcode
	}

	completedDate, err := parseCompletedDate(fieldAt(fields, dateField), fieldAt(fields, timeField))
	if err != nil {
		return resultsHeader{}, err
	}

	return resultsHeader{barcode: barcode, completedDate: completedDate, testCount: testCount}, nil
}

// sampleBarcode gets the sample identifier: the patient ID field (shown as "Sample ID" on the
// analyzer screen, where a barcode scan lands), or the running sample No. when it is blank.
func sampleBarcode(fields []string) string {
	patientID := fieldAt(fields, patientIDField)
	if patientID != "" {
		return patientID
	}

	return fieldAt(fields, sampleNoField)
}

// parseCompletedDate parses the test date and time.
func parseCompletedDate(date, clock string) (time.Time, error) {
	value := date + " " + clock

	completedDate, err := time.Parse(dateLayoutDayFirst, value)
	if err == nil {
		return completedDate, nil
	}

	completedDate, err = time.Parse(dateLayoutISO, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get completed date from %q: %w", value, err)
	}

	return completedDate, nil
}

// testLabData converts the fields of one test to a lab data. It returns nil for a test without a name.
func testLabData(testFields []string) *model.LabData {
	param := fieldAt(testFields, testNameOffset)
	if param == "" {
		return nil
	}

	result, unit := splitResultAndUnit(testFields[testResultOffset])

	sign := fieldAt(testFields, testSignOffset)
	if sign != "" && sign != signEqual {
		result = sign + result
	}

	return &model.LabData{
		Param:  param,
		Result: result,
		Unit:   unit,
	}
}

// splitResultAndUnit splits the fixed width result field into the result and its unit.
func splitResultAndUnit(field string) (string, string) {
	if len(field) <= resultWidth {
		return strings.TrimSpace(field), ""
	}

	return strings.TrimSpace(field[:resultWidth]), strings.TrimSpace(field[resultWidth:])
}

// fieldAt returns the field at the index without its padding, or an empty string when the record is shorter.
func fieldAt(fields []string, index int) string {
	if index < len(fields) {
		return strings.TrimSpace(fields[index])
	}

	return ""
}
