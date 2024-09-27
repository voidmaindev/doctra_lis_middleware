package driver_astm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/voidmaindev/doctra_lis_middleware/model"
	"github.com/voidmaindev/doctra_lis_middleware/services"
)

const (
	stx = "\x02"
	etb = "\x17"
	etx = "\x03"
	cr  = "\x0D"
	lf  = "\x0A"
)

// Unmarshal unmarshals the raw data to the lab data.
func (d *Driver_astm) Unmarshal(rawData string) (labDatas []*model.LabData, additionalData map[string]interface{}, err error) {
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

	astm_msgs := parseASTMMessages(rawData)

	if isQuery, qMsg := isQueryMessage(astm_msgs); isQuery {
		additionalData = map[string]interface{}{queryName: true,
			queryMessagesName: astm_msgs,
			"sample_id":       qMsg.SampleID,
			"test_id":         qMsg.TestID,
		}
		return labDatas, additionalData, nil
	}

	for _, astm_msg := range astm_msgs {
		barcode, err := getBarcodeForUnmarshalRawData(astm_msg)
		if err != nil {
			fmt.Println("failed to get barcode for unmarshalRawData")
			return labDatas, nil, err
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
				return labDatas, nil, err
			}

			param, err := getParamForUnmarshalRawData(result)
			if err != nil {
				fmt.Println("failed to get param for unmarshalRawData")
				return labDatas, nil, err
			}

			res, err := getResultForUnmarshalRawData(result)
			if err != nil {
				fmt.Println("failed to get result for unmarshalRawData")
				return labDatas, nil, err
			}

			unit, err := getUnitForUnmarshalRawData(result)
			if err != nil {
				fmt.Println("failed to get unit for unmarshalRawData")
				return labDatas, nil, err
			}

			completedDate, err := getCompleteDateForUnmarshalRawData(result)
			if err != nil {
				fmt.Println("failed to get completed date for unmarshalRawData")
				return labDatas, nil, err
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
	}

	return labDatas, additionalData, nil
}

// getBarcodeForUnmarshalRawData gets the barcode for unmarshalling the raw data.
func getBarcodeForUnmarshalRawData(msg Message) (string, error) {
	barcode := strings.Split(msg.Order.PatientName, "^")[0]

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
	Header       Header
	Patient      Patient
	Order        Order
	Results      []Result
	Comments     []Comment
	Term         Termination
	Query        Query
	Notification Notification
}

// Query represents the query/request segment (Q)
type Query struct {
	Type      string
	QueryType string
	SampleID  string
	TestID    string
}

// Notification represents the notification segment (N)
type Notification struct {
	Type             string
	NotificationType string
	Details          string
}

// parseASTMMessages parses the raw data into a slice of structured messages
func parseASTMMessages(data string) []Message {
	var messages []Message
	segments := strings.Split(data, stx)
	var currentMessage Message
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
				if currentMessage.Header.Type != "" {
					currentMessage.Results = results
					currentMessage.Comments = comments
					messages = append(messages, currentMessage)
					currentMessage = Message{}
					results = []Result{}
					comments = []Comment{}
				}
				currentMessage.Header = parseHeader(content)
			case "P":
				currentMessage.Patient = parsePatient(content)
			case "O":
				currentMessage.Order = parseOrder(content)
			case "R":
				results = append(results, parseResult(content))
			case "C":
				comments = append(comments, parseComment(content))
			case "L":
				currentMessage.Term = parseTermination(content)
			case "Q":
				currentMessage.Query = parseQuery(content)
			case "N":
				currentMessage.Notification = parseNotification(content)
			}
		}
	}

	if currentMessage.Header.Type != "" {
		currentMessage.Results = results
		currentMessage.Comments = comments
		messages = append(messages, currentMessage)
	}

	return messages
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
	testId := testIDs[3]
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

// parseQuery parses the query segment (Q)
func parseQuery(content string) Query {
	parts := strings.Split(content, "|")
	return Query{
		Type:      "Q",
		QueryType: parts[1],
		SampleID:  parts[2],
		TestID:    parts[3],
	}
}

// parseNotification parses the notification segment (N)
func parseNotification(content string) Notification {
	parts := strings.Split(content, "|")
	return Notification{
		Type:             "N",
		NotificationType: parts[1],
		Details:          parts[2],
	}
}

// isQueryMessage checks if the message is a query message
func isQueryMessage(messages []Message) (bool, *Query) {
	for _, msg := range messages {
		if msg.Query.Type == "Q" {
			return true, &msg.Query
		}
	}
	return false, nil
}

// QueryAnswerOrder represents the structure of the query answer order
type QueryAnswerOrder struct {
	ID        string
	PatientID string
	Param     string
	Priority  string
	Report    string
}

// formatQueryAnswerOrderMessage formats the structured QueryAnswerOrder message
func formatQueryAnswerOrderMessage(order QueryAnswerOrder) string {
	return stx + fmt.Sprintf("O|%s|%s||^^^%s^\\^^^555|%s||||||%s||||||||||||||O\\Q",
		order.ID,
		order.PatientID,
		order.Param,
		order.Priority,
		order.Report,
	) + cr + string(etx) + cr + lf
}

// generateASTMMessagesFromQuery formats the queryMessages and appends data from dataToReturn
func generateASTMMessagesFromQuery(queryMessages []Message, dataToReturn []services.DeviceQueryDataToReturn) []string {
	var formattedMessages []string

	// Loop over queryMessages and format them
	for _, msg := range queryMessages {
		var formattedMsg string

		// Handle different message types
		switch msg.Header.Type {
		case "H": // Header
			formattedMsg = fmt.Sprintf("H|\\^&|%s|%s|%s|%s|%s",
				msg.Header.Sender,
				msg.Header.Receiver,
				msg.Header.AnalyzerType,
				msg.Header.Version,
				msg.Header.Timestamp,
			)
		case "P": // Patient
			formattedMsg = fmt.Sprintf("P|%s", msg.Patient.ID)
		case "Q": // Query
			for i, data := range dataToReturn {
				order := QueryAnswerOrder{
					ID:        fmt.Sprintf("%d", i+1),
					PatientID: msg.Query.SampleID,
					Param:     data.Param,
					Priority:  "R",
					Report:    "A",
				}
				formattedMessages = append(formattedMessages, formatQueryAnswerOrderMessage(order))
			}

		case "C": // Comments
			for _, comment := range msg.Comments {
				formattedMsg = fmt.Sprintf("C|%s", comment.Comment)
				formattedMessages = append(formattedMessages, stx+formattedMsg+cr+lf)
			}
		case "L": // Termination
			formattedMsg = fmt.Sprintf("L|%s", msg.Term.TerminationCode)
		}

		// Add STX and ETX framing
		if formattedMsg != "" {
			formattedMsg = stx + formattedMsg + cr + string(etx) + cr + lf
			formattedMessages = append(formattedMessages, formattedMsg)
		}
	}

	return formattedMessages
}
