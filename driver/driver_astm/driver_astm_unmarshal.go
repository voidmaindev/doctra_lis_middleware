package driver_astm

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
)

// Unmarshal unmarshals the raw data to the lab data.
func (d *Driver_astm) Unmarshal(rawData string) (labDatas []*model.LabData, err error) {
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

	astm_msg := parseASTMMessage(rawData)

	barcode, err := getBarcodeForUnmarshalRawData(astm_msg)
	if err != nil {
		fmt.Println("failed to get barcode for unmarshalRawData")
		return labDatas, err
	}

	index := 0
	for _, result := range astm_msg.Results {
		if result.TestID == "" || result.Value == "" || result.Units == "" || result.Timestamp == "" {
			continue
		}

		index++

		index, err := getIndexForUnmarshalRawData(index)
		if err != nil {
			fmt.Println("failed to get index for unmarshalRawData")
			return labDatas, err
		}

		param, err := getParamForUnmarshalRawData(result)
		if err != nil {
			fmt.Println("failed to get param for unmarshalRawData")
			return labDatas, err
		}

		res, err := getResultForUnmarshalRawData(result)
		if err != nil {
			fmt.Println("failed to get result for unmarshalRawData")
			return labDatas, err
		}

		unit, err := getUnitForUnmarshalRawData(result)
		if err != nil {
			fmt.Println("failed to get unit for unmarshalRawData")
			return labDatas, err
		}

		completedDate, err := getCompleteDateForUnmarshalRawData(result)
		if err != nil {
			fmt.Println("failed to get completed date for unmarshalRawData")
			return labDatas, err
		}

		labData := &model.LabData{
			Barcode:       barcode,
			Index:         index,
			Param:         param,
			Result:        res,
			Unit:          unit,
			CompletedDate: completedDate,
		}

		labDatas = append(labDatas, labData)
	}

	return nil, nil
}

// getBarcodeForUnmarshalRawData gets the barcode for unmarshalling the raw data.
func getBarcodeForUnmarshalRawData(msg Message) (string, error) {
	barcode := msg.Order.ID

	return barcode, nil
}

// getCompleteDateForUnmarshalRawData gets the completed date for unmarshalling the raw data.
func getCompleteDateForUnmarshalRawData(result Result) (time.Time, error) {
	parsedDate, err := time.Parse(completedDateFormat, result.Timestamp)
	if err != nil {
		return time.Time{}, errors.New("failed to get completed date")
	}

	return parsedDate, nil
}

// getIndexForUnmarshalRawData gets the index for unmarshalling the raw data.
func getIndexForUnmarshalRawData(i int) (uint, error) {
	return uint(i), nil
}

// getParamForUnmarshalRawData gets the param for unmarshalling the raw data.
func getParamForUnmarshalRawData(result Result) (string, error) {
	param := result.TestID

	return param, nil
}

// getResultForUnmarshalRawData gets the result for unmarshalling the raw data.
func getResultForUnmarshalRawData(result Result) (string, error) {
	res := result.Value

	return res, nil
}

// getUnitForUnmarshalRawData gets the unit for unmarshalling the raw data.
func getUnitForUnmarshalRawData(result Result) (string, error) {
	unit := result.Units

	return unit, nil
}

// Segment represents a parsed segment of the message
type Segment struct {
	Type    string
	Content string
}

// Header represents the header segment
type Header struct {
	Type           string
	Sender         string
	Receiver       string
	AnalyzerType   string
	Version        string
	SequenceNumber string
	Timestamp      string
	MessageType    string
	ProcessingID   string
}

// Patient represents the patient segment
type Patient struct {
	Type string
	ID   string
}

// Order represents the order segment
type Order struct {
	Type        string
	ID          string
	PatientID   string
	PatientName string
	Priority    string
	Timestamp   string
	Status      string
}

// Result represents the result segment
type Result struct {
	Type           string
	TestID         string
	Value          string
	Units          string
	ReferenceRange string
	Status         string
	Timestamp      string
}

// Comment represents the comment segment
type Comment struct {
	Type    string
	Comment string
}

// Termination represents the termination segment
type Termination struct {
	Type            string
	TerminationCode string
}

// Message represents a parsed message from the communication
type Message struct {
	Header   Header
	Patient  Patient
	Order    Order
	Results  []Result
	Comments []Comment
	Term     Termination
}

// parseMessage parses the raw data into a structured message
func parseASTMMessage(data string) Message {
	var message Message
	segments := strings.Split(data, stx)
	var results []Result
	var comments []Comment

	for _, segment := range segments {
		if len(segment) > 2 {
			segmentType := string(segment[1])
			contentEnd := strings.Index(segment, etb)
			if contentEnd == -1 {
				contentEnd = strings.Index(segment, etx)
			}
			if contentEnd == -1 {
				contentEnd = len(segment)
			}

			content := segment[2:contentEnd]

			switch segmentType {
			case "H":
				message.Header = parseHeader(content)
			case "P":
				message.Patient = parsePatient(content)
			case "O":
				message.Order = parseOrder(content)
			case "R":
				results = append(results, parseResult(content))
			case "C":
				comments = append(comments, parseComment(content))
			case "L":
				message.Term = parseTermination(content)
			}
		}
	}

	message.Results = results
	message.Comments = comments
	return message
}

// parseHeader parses the header segment
func parseHeader(content string) Header {
	parts := strings.Split(content, "|")
	return Header{
		Type:           "H",
		Sender:         parts[4],
		Receiver:       parts[5],
		AnalyzerType:   strings.Split(parts[6], "^")[0],
		Version:        parts[6],
		SequenceNumber: parts[7],
		Timestamp:      parts[13],
		MessageType:    parts[10],
		ProcessingID:   parts[11],
	}
}

// parsePatient parses the patient segment
func parsePatient(content string) Patient {
	parts := strings.Split(content, "|")
	return Patient{
		Type: "P",
		ID:   parts[1],
	}
}

// parseOrder parses the order segment
func parseOrder(content string) Order {
	parts := strings.Split(content, "|")
	return Order{
		Type:        "O",
		ID:          parts[1],
		PatientID:   parts[2],
		PatientName: parts[3],
		Priority:    parts[5],
		Timestamp:   parts[22],
		Status:      parts[21],
	}
}

// parseResult parses the result segment
func parseResult(content string) Result {
	parts := strings.Split(content, "|")
	testIDs := strings.Split(parts[2], "^")
	testId := testIDs[len(testIDs)-1]
	return Result{
		Type:           "R",
		TestID:         testId,
		Value:          parts[3],
		Units:          parts[4],
		ReferenceRange: parts[5],
		Status:         parts[6],
		Timestamp:      parts[12],
	}
}

// parseComment parses the comment segment
func parseComment(content string) Comment {
	parts := strings.Split(content, "|")
	return Comment{
		Type:    "C",
		Comment: parts[3],
	}
}

// parseTermination parses the termination segment
func parseTermination(content string) Termination {
	parts := strings.Split(content, "|")
	return Termination{
		Type:            "L",
		TerminationCode: parts[1],
	}
}
