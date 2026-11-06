package driver_astm_wondfo_accre_8

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/voidmaindev/doctra_lis_middleware/model"
)

const (
	stx = "\x02"
	etb = "\x17"
	etx = "\x03"
	cr  = "\x0D"
	lf  = "\x0A"

	// The firmware writes these control characters as text instead of single bytes.
	literalSTX   = `\x02`
	literalETX   = `\x03`
	literalCaret = `\x5e`

	fieldSeparator     = "|"
	componentSeparator = "^"
	frameNumberDigits  = "01234567"

	recordTypeOrder  = "O"
	recordTypeResult = "R"

	// Order record: O|seq|<constant>^<sample ID>^tray^pos|<specimen barcode>|...
	orderSampleIDField         = 2
	orderSpecimenIDField       = 3
	orderSampleIDComponent     = 1
	orderSampleIDFallbackComp  = 0
	resultTestField            = 2 // ^^^<test name>
	resultTestNameComponent    = 3
	resultValueField           = 3
	resultUnitField            = 4
	resultCompletedDateField   = 12
	completedDateLayout        = "20060102 15:04:05"
	completedDateLayoutCompact = "20060102150405"
)

var (
	errNoBarcode       = errors.New("result record without a sample ID")
	errUnmarshalFailed = errors.New("failed to unmarshal raw data")
)

// Unmarshal unmarshals the raw data (everything between ENQ and EOT) to the lab data.
func (d *Driver_astm_wondfo_accre_8) Unmarshal(rawData string) (labDatas []*model.LabData, additionalData map[string]interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			d.log.Error("unknown error occurred while unmarshalling raw data: " + fmt.Sprint(r))
			labDatas = []*model.LabData{}
			err = errUnmarshalFailed
		}
	}()

	barcode := ""
	index := uint(0)

	for _, record := range splitRecords(rawData) {
		fields := strings.Split(record, fieldSeparator)

		switch fields[0] {
		case recordTypeOrder:
			barcode = orderBarcode(fields)
			index = 0
		case recordTypeResult:
			labData, err := resultLabData(fields, barcode)
			if err != nil {
				d.log.Err(err, "failed to parse the result record: "+record)
				return labDatas, nil, err
			}
			if labData == nil {
				continue
			}

			index++
			labData.Index = index
			labDatas = append(labDatas, labData)
		}
	}

	return labDatas, additionalData, nil
}

// normalizeControlChars turns the control characters that the firmware sends as text into real ones.
func normalizeControlChars(data string) string {
	replacer := strings.NewReplacer(
		literalSTX, stx,
		literalETX, etx,
		literalCaret, componentSeparator,
	)

	return replacer.Replace(data)
}

// splitRecords splits the raw data into ASTM records (H, P, O, R, L, ...) without framing.
// Intermediate frames (ended by ETB) are joined with the frame that follows them.
// Checksums are ignored: the analyzer sends constants that do not match the frame content.
func splitRecords(rawData string) []string {
	records := []string{}
	pending := ""

	for _, frame := range strings.Split(normalizeControlChars(rawData), stx) {
		content, isIntermediate := frameContent(frame)
		pending += content
		if isIntermediate {
			continue
		}

		if pending != "" {
			records = append(records, pending)
		}
		pending = ""
	}

	return records
}

// frameContent returns the record text of a frame and whether the frame is an intermediate one.
func frameContent(frame string) (string, bool) {
	end := strings.IndexAny(frame, etb+etx)
	isIntermediate := end != -1 && frame[end:end+1] == etb
	if end == -1 {
		end = len(frame)
	}

	content := strings.Trim(frame[:end], cr+lf)
	if content != "" && strings.Contains(frameNumberDigits, content[:1]) {
		content = content[1:]
	}

	return content, isIntermediate
}

// orderBarcode gets the sample identifier from the order record: the scanned specimen barcode
// when present, otherwise the sample ID typed by the operator.
func orderBarcode(fields []string) string {
	specimenID := strings.TrimSpace(fieldAt(fields, orderSpecimenIDField))
	if specimenID != "" {
		return specimenID
	}

	sampleField := fieldAt(fields, orderSampleIDField)
	sampleID := strings.TrimSpace(componentAt(sampleField, orderSampleIDComponent))
	if sampleID != "" {
		return sampleID
	}

	return strings.TrimSpace(componentAt(sampleField, orderSampleIDFallbackComp))
}

// resultLabData converts a result record to a lab data. It returns nil for a record without a test name or a value.
func resultLabData(fields []string, barcode string) (*model.LabData, error) {
	param := strings.TrimSpace(componentAt(fieldAt(fields, resultTestField), resultTestNameComponent))
	value := strings.TrimSpace(fieldAt(fields, resultValueField))
	if param == "" || value == "" {
		return nil, nil
	}

	if barcode == "" {
		return nil, errNoBarcode
	}

	completedDate, err := parseCompletedDate(fieldAt(fields, resultCompletedDateField))
	if err != nil {
		return nil, err
	}

	return &model.LabData{
		Barcode:       barcode,
		Param:         param,
		Result:        value,
		Unit:          strings.TrimSpace(fieldAt(fields, resultUnitField)),
		CompletedDate: completedDate,
	}, nil
}

// parseCompletedDate parses the test completion time ("20261103 17:30:48", or "20261103173048" as in the vendor manual).
func parseCompletedDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)

	completedDate, err := time.Parse(completedDateLayout, value)
	if err == nil {
		return completedDate, nil
	}

	completedDate, err = time.Parse(completedDateLayoutCompact, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("failed to get completed date from %q: %w", value, err)
	}

	return completedDate, nil
}

// fieldAt returns the field at the index, or an empty string when the record is shorter.
func fieldAt(fields []string, index int) string {
	if index < len(fields) {
		return fields[index]
	}

	return ""
}

// componentAt returns the component of a field at the index, or an empty string when the field has fewer components.
func componentAt(field string, index int) string {
	return fieldAt(strings.Split(field, componentSeparator), index)
}
