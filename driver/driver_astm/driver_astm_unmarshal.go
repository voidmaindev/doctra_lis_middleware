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
		additionalData = map[string]interface{}{queryName: qMsg, queryMessagesName: astm_msgs}
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
	Type               string    // H: Identifies the segment as a Header
	SenderID           string    // Identifier of the sender (usually the device or lab system)
	ReceiverID         string    // Identifier of the receiver (usually the host or middleware)
	SenderName         string    // Full name of the sending entity (device manufacturer or lab name)
	ReceiverName       string    // Full name of the receiving entity
	AnalyzerType       string    // Analyzer model type and version
	Version            string    // ASTM version used for communication
	ControlID          string    // ID for tracking and logging purposes
	AccessPassword     string    // Security feature, if applicable
	SequenceNumber     string    // Sequence number of the message for tracking
	Timestamp          time.Time // Timestamp of when the message was generated (formatted properly)
	MessageType        string    // Type of the message (e.g., Real-Time, Query, Response)
	ProcessingMode     string    // Processing mode (e.g., REAL, SIM for Simulation)
	ProcessingID       string    // ID indicating if the message is a Test, Production, or Simulation
	AcknowledgmentCode string    // Code used to indicate acknowledgment status of previous message
	CharacterSet       string    // Defines character set used in message encoding (e.g., ASCII, UTF-8)
	SecurityCode       string    // Optional field for security authentication or encryption code
	SoftwareVersion    string    // Version of the software running on the device
	ApplicationName    string    // Name of the application (if applicable)
	FacilityID         string    // Identifier of the facility where the test is being conducted
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
	Type                string    // "Q" identifies the segment as a query
	QueryType           string    // Type of query (e.g., order request, test query)
	SampleID            string    // Identifier for the sample being queried
	TestID              string    // Identifier for the specific test being requested
	Priority            string    // Priority level of the query (e.g., Routine, STAT)
	RequestedTimestamp  time.Time // Timestamp of when the query was requested
	RequestingFacility  string    // Name or ID of the facility making the request
	RequestingPhysician string    // Name or ID of the requesting physician
	TestStatus          string    // Status of the test being requested
	ReasonForQuery      string    // Reason for the query (optional for tracking or auditing)
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
			segmentStartIndex := 2
			if segmentType == "|" {
				segmentType = string(segment[0])
				segmentStartIndex = 1
			}
			contentEnd := strings.Index(segment, etb)
			if contentEnd == -1 {
				contentEnd = strings.Index(segment, etx)
			}
			if contentEnd == -1 {
				contentEnd = len(segment)
			}

			content := segment[segmentStartIndex:contentEnd]

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

	// Split for sender and receiver names, if applicable
	senderParts := strings.Split(parts[4], "^")
	receiverParts := strings.Split(parts[9], "^")
	analyzerParts := strings.Split(parts[6], "^")

	messageProcessingParts := strings.Split(parts[10], "^") // TSREQ^REAL

	return Header{
		Type:            "H",
		SenderID:        senderParts[0],
		SenderName:      ifExists(senderParts, 1),
		ReceiverID:      receiverParts[0],
		ReceiverName:    ifExists(receiverParts, 1),
		AnalyzerType:    analyzerParts[0],
		Version:         ifExists(analyzerParts, 1),
		ControlID:       parts[8],           // Assuming part[8] is for control tracking/logging
		AccessPassword:  ifExists(parts, 9), // Assigns password/security, part[9] possibly the missing field
		Timestamp:       parseTimestamp(parts[13]),
		MessageType:     messageProcessingParts[0],                   // TSREQ or other message type
		ProcessingMode:  ifExists(messageProcessingParts, 1, "REAL"), // REAL, SIM, etc.
		ProcessingID:    parts[11],
		SequenceNumber:  ifExists(parts, 12),
		CharacterSet:    ifExists(parts, 15, "ASCII"), // Default to ASCII if not present
		SecurityCode:    ifExists(parts, 16),
		SoftwareVersion: ifExists(parts, 17),
		ApplicationName: ifExists(parts, 18),
		FacilityID:      ifExists(parts, 19),
	}
}

// Helper function to safely extract a value if exists in the slice
func ifExists(parts []string, idx int, defaults ...string) string {
	if idx < len(parts) {
		return parts[idx]
	}
	if len(defaults) > 0 {
		return defaults[0]
	}
	return ""
}

// Helper function to parse timestamp safely
func parseTimestamp(ts string) time.Time {
	parsedTime, err := time.Parse("20060102150405", ts)
	if err != nil {
		return time.Now() // Fallback to current time in case of parsing error
	}
	return parsedTime
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
	patientName := parts[3]
	if len(patientName) == 0 {
		patientName = parts[2]
	}

	return Order{
		Type:        "O",
		ID:          parts[1],
		PatientID:   parts[2],
		PatientName: patientName,
		// Priority:    parts[5],
		// Timestamp:   parts[22],
		// Status:      parts[21],
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

	// Extract timestamps safely and set defaults if parsing fails
	timestamp, err := time.Parse("20060102150405", parts[11])
	if err != nil {
		timestamp = time.Now() // Default to current time if parsing fails
	}

	return Query{
		Type:                "Q",                            // Fixed as this is a query segment
		QueryType:           ifExists(parts, 1),             // Query type (part[1])
		SampleID:            strings.Trim(parts[2], "^"),    // Sample ID (part[2])
		TestID:              ifExists(parts, 3),             // Test ID (part[3])
		Priority:            ifExists(parts, 5, "Routine"),  // Priority (optional, defaults to Routine if not provided)
		RequestedTimestamp:  timestamp,                      // Parsed timestamp
		RequestingFacility:  ifExists(parts, 8),             // Requesting facility (optional)
		RequestingPhysician: ifExists(parts, 9),             // Requesting physician (optional)
		TestStatus:          ifExists(parts, 10, "Pending"), // Status (defaults to "Pending")
		ReasonForQuery:      ifExists(parts, 12),            // Reason for the query (optional)
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

// generateASTMMessagesFromQuery formats the queryMessages and appends data from dataToReturn
func generateASTMMessagesFromQuery(queryMessages []Message, dataToReturn []services.DeviceQueryDataToReturn) []string {
	var formattedMessages []string

	// Loop over queryMessages and format them
	for _, msg := range queryMessages {
		var formattedMsg string

		// Handle different message types
		if msg.Header.Type == "H" {
			if msg.Header.MessageType == "TSREQ" && msg.Header.ProcessingMode == "REAL" {
				msg.Header.MessageType = "TSDWN"
				msg.Header.ProcessingMode = "REPLY"
			}
			formattedMsg = fmt.Sprintf("H|\\^&|||%s|||||%s|%s^%s|%s|%s|%s",
				msg.Header.ReceiverID,                         // ID of the receiver
				msg.Header.SenderID,                           // ID of the sender
				msg.Header.MessageType,                        // TSDWN (Test Shutdown)
				msg.Header.ProcessingMode,                     // REPLY
				msg.Header.ProcessingID,                       // Processing ID (e.g., P for production)
				msg.Header.SequenceNumber,                     // Sequence number of the message
				msg.Header.Timestamp.Format("20060102150405"), // Timestamp in the correct format
			)
			formattedMessages = addFormattedMessage(formattedMessages, formattedMsg)
		}
		if msg.Query.Type == "Q" {
			for i, data := range dataToReturn {
				order := QueryAnswerOrder{
					ID:        fmt.Sprintf("%d", i+1),
					PatientID: msg.Query.SampleID,
					Param:     data.Param,
					Priority:  "R",
					Report:    "A",
				}
				formattedMsg = fmt.Sprintf("O|%s|%s||^^^%s^\\^^^555|%s||||||%s||||||||||||||O\\Q",
					order.ID,
					order.PatientID,
					order.Param,
					order.Priority,
					order.Report,
				)
				formattedMessages = addFormattedMessage(formattedMessages, formattedMsg)
			}
		}

		// if msg.Query.Type == "C" {
		// 	for _, comment := range msg.Comments {
		// 		formattedMsg = fmt.Sprintf("C|%s", comment.Comment)
		// 		formattedMessages = append(formattedMessages, stx+formattedMsg+cr+lf)
		// 	}
		if msg.Term.Type == "L" {
			formattedMessages = addFormattedMessage(formattedMessages, fmt.Sprintf("L|%s|N", msg.Term.TerminationCode))
		}
	}

	return formattedMessages
}

// addFormattedMessage adds the formatted message to the formattedMessages slice
func addFormattedMessage(formattedMessages []string, formattedMsg string) []string {
	checkSum := calculateASTMChecksum(formattedMsg)
	formattedMsg = formattedMsg + cr + string(etx) + fmt.Sprintf("%02X", checkSum) + cr + lf

	return append(formattedMessages, formattedMsg)
}

// calculateASTMChecksum calculates the checksum for ASTM protocol
func calculateASTMChecksum(content string) int {
	checkSum := 0

	for i := 0; i < len(content); i++ {
		checkSum ^= int(content[i])
	}

	return checkSum
}
