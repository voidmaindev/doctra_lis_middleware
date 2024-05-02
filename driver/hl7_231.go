package driver

import (
	"fmt"
	"strings"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/model"
	"github.com/voidmaindev/doctra_lis_middleware/store"
	"github.com/voidmaindev/doctra_lis_middleware/tcp"
)

const (
	rawDataStartChar = 0xb
	rawDataEndChar   = 0x1c
)

// Driver_hl7_231 is the driver for the "HL7 2.3.1" laboratory device data format.
type Driver_hl7_231 struct {
	log                *log.Logger
	store              *store.Store
	rawDataStartString string
	rawDataEndString   string
	dataToBeReplaced   map[string]string
}

// hl7Message represents the entire HL7 message with segments stored in a map where keys are segment types.
type hl7Message struct {
	Segments map[string][]map[string]interface{} `json:"segments"`
}

// NewDriver_hl7_231 creates a new "HL7 2.3.1" driver.
func NewDriver_hl7_231(logger *log.Logger, store *store.Store) *Driver_hl7_231 {
	d := &Driver_hl7_231{
		log:   logger,
		store: store,
	}

	d.rawDataStartString = fmt.Sprintf("%c", rawDataStartChar)
	d.rawDataEndString = fmt.Sprintf("%c", rawDataEndChar)
	d.dataToBeReplaced = map[string]string{"\\r": "\n"}

	return d
}

// ProcessDeviceMessage processes the device message.
func (d *Driver_hl7_231) ProcessDeviceMessage(deviceMsg []byte, conn *tcp.ConnData, device *model.Device) error {
	msg := string(deviceMsg)
	for k, v := range d.dataToBeReplaced {
		msg = strings.ReplaceAll(msg, k, v)
	}

	rawDatas := d.getRawDatas(msg, &conn.Data)

	for _, rawData := range rawDatas {
		rd := &model.RawData{
			ConnString: conn.ConnString,
			DeviceID:   device.ID,
			Data:       []byte(rawData),
			Processed:  true,
		}

		labDatas, err := d.unmarshalRawData(rawData)
		if err != nil {
			d.log.Error("failed to unmarshal a raw data from " + device.Name)
			rd.Processed = false
		}

		err = d.store.RawDataStore.Create(rd)
		if err != nil {
			d.log.Error("failed to create a raw data from " + device.Name)
			return err
		}

		for _, labData := range labDatas {
			err = d.store.LabDataStore.Create(labData)
			if err != nil {
				d.log.Error(fmt.Sprintf("failed to create a lab data from %s with barcode %s and index %d", device.Name, labData.Barcode, labData.Index))
				continue
			}
		}
	}

	return nil
}

// getRawDatas gets the raw datas from the message.
func (d *Driver_hl7_231) getRawDatas(msg string, prevRawData *string) []string {
	rawDatas := []string{}

	bigChunks := strings.Split(msg, d.rawDataStartString)
	for _, bigChunk := range bigChunks {
		if len(bigChunk) == 0 {
			continue
		}

		chunks := strings.Split(bigChunk, d.rawDataEndString)
		switch len(chunks) {
		case 1:
			*prevRawData += bigChunk
		case 2:
			rawDatas = append(rawDatas, *prevRawData+chunks[0])
			*prevRawData = ""
		}
	}

	return rawDatas
}

// unmarshalRawData unmarshals the raw data.
func (d *Driver_hl7_231) unmarshalRawData(rawData string) ([]*model.LabData, error) {
	labDatas := []*model.LabData{}

	// // hl7msg, err := parseHL7Message(rawData)
	// if err != nil {
	// 	d.log.Error("failed to parse HL7 message")
	// 	return labDatas, err
	// }

	return labDatas, nil
}

// parseHL7Message parses the HL7 message.
func parseHL7Message(rawMessage string) (*hl7Message, error) {
	message := &hl7Message{Segments: make(map[string][]map[string]interface{})}
	normalized := strings.ReplaceAll(rawMessage, "\r\n", "\r")
	normalized = strings.ReplaceAll(normalized, "\n", "\r")
	segments := strings.Split(normalized, "\r")

	delimiters := struct {
		field, component, repetition, escape, subComponent string
	}{
		field:        "|",
		component:    "^",
		repetition:   "~",
		escape:       "\\",
		subComponent: "&",
	}

	for _, segment := range segments {
		if len(segment) < 3 {
			continue
		}
		if segment[:3] == "MSH" {
			delimiters = parseDelimiters(segment)
		}
		fields := strings.Split(segment, delimiters.field)
		if len(fields) == 0 {
			continue
		}
		segmentName := fields[0]
		segmentFields := make(map[string]interface{})
		fieldDefinitions := getFieldDefinitions(segmentName)

		for i, field := range fields[1:] {
			fieldName := "Unknown"
			if i < len(fieldDefinitions) {
				fieldName = fieldDefinitions[i]
			}
			componentNames := getComponentNames(segmentName, fieldName)

			if strings.ContainsAny(field, delimiters.component+delimiters.repetition+delimiters.subComponent) {
				segmentFields[fieldName] = parseComplexField(field, delimiters, componentNames)
			} else {
				segmentFields[fieldName] = field
			}
		}
		message.Segments[segmentName] = append(message.Segments[segmentName], segmentFields)
	}

	return message, nil
}

// parseDelimiters parses the delimiters of MSH segment.
func parseDelimiters(mshSegment string) struct{ field, component, repetition, escape, subComponent string } {
	return struct{ field, component, repetition, escape, subComponent string }{
		field:        "|",
		component:    string(mshSegment[4]),
		repetition:   string(mshSegment[5]),
		escape:       string(mshSegment[6]),
		subComponent: string(mshSegment[7]),
	}
}

// parseComplexField parses the complex field of the HL7 message.
func parseComplexField(field string, delimiters struct{ field, component, repetition, escape, subComponent string }, componentNames []string) map[string]interface{} {
	componentParts := strings.Split(field, delimiters.component)
	result := make(map[string]interface{})
	for i, component := range componentParts {
		componentName := fmt.Sprintf("Component%d", i+1) // Default naming
		if i < len(componentNames) {
			componentName = componentNames[i]
		}
		subComponents := strings.Split(component, delimiters.subComponent)
		if len(subComponents) > 1 {
			subResult := make(map[string]string)
			for j, subComponent := range subComponents {
				subComponentName := fmt.Sprintf("%s.Subcomponent%d", componentName, j+1)
				subResult[subComponentName] = subComponent
			}
			result[componentName] = subResult
		} else {
			result[componentName] = component
		}
	}
	return result
}

// getFieldDefinitions gets the field definitions of the HL7 message.
func getFieldDefinitions(segmentName string) []string {
	switch segmentName {
	case "MSH":
		return []string{"Encoding Characters", "Sending Application", "Sending Facility", "Receiving Application", "Receiving Facility", "Date/Time of Message", "Security", "Message Type", "Message Control ID", "Processing ID", "Version ID"}
	case "PID":
		return []string{"Set ID - PID", "Patient ID", "Patient Identifier List", "Alternate Patient ID - PID", "Patient Name", "Mother’s Maiden Name", "Date/Time of Birth", "Sex", "Patient Alias", "Race", "Patient Address"}
	case "OBX":
		return []string{"Set ID - OBX", "Value Type", "Observation Identifier", "Observation Sub-ID", "Observation Value", "Units", "References Range", "Abnormal Flags", "Probability", "Nature of Abnormal Test", "Observation Result Status", "Effective Date of Reference Range", "User Defined Access Checks", "Date/Time of the Observation", "Producer's ID", "Responsible Observer", "Observation Method", "Equipment Instance Identifier", "Date/Time of the Analysis"}
	case "OBR":
		return []string{"Set ID - OBR", "Placer Order Number", "Filler Order Number", "Universal Service ID", "Priority", "Requested Date/Time", "Observation Date/Time", "Observation End Date/Time", "Collection Volume", "Collector Identifier", "Specimen Action Code", "Danger Code", "Relevant Clinical Info", "Specimen Received Date/Time", "Specimen Source", "Ordering Provider", "Order Callback Phone Number", "Placer Field 1", "Placer Field 2", "Filler Field 1", "Filler Field 2", "Results Rpt/Status Chng - Date/Time", "Charge to Practice", "Diagnostic Serv Sect ID", "Result Status", "Parent Result", "Quantity/Timing", "Result Copies To", "Parent", "Transportation Mode", "Reason for Study", "Principal Result Interpreter", "Assistant Result Interpreter", "Technician", "Transcriptionist", "Scheduled Date/Time", "Number of Sample Containers", "Transport Logistics of Collected Sample", "Collector’s Comment", "Transport Arrangement Responsibility", "Transport Arranged", "Escort Required", "Planned Patient Transport Comment"}
	case "PV1":
		return []string{"Set ID - PV1", "Patient Class", "Assigned Patient Location", "Admission Type", "Preadmit Number", "Prior Patient Location", "Attending Doctor", "Referring Doctor", "Consulting Doctor", "Hospital Service", "Temporary Location", "Preadmit Test Indicator", "Readmission Indicator", "Admit Source", "Ambulatory Status", "VIP Indicator", "Admitting Doctor", "Patient Type", "Visit Number", "Financial Class", "Charge Price Indicator", "Courtesy Code", "Credit Rating", "Contract Code", "Contract Effective Date", "Contract Amount", "Contract Period", "Interest Code", "Transfer to Bad Debt Code", "Transfer to Bad Debt Date", "Bad Debt Agency Code", "Bad Debt Transfer Amount", "Bad Debt Recovery Amount", "Delete Account Indicator", "Delete Account Date", "Discharge Disposition", "Discharged to Location", "Diet Type", "Servicing Facility", "Bed Status", "Account Status", "Pending Location", "Prior Temporary Location", "Admit Date/Time", "Discharge Date/Time", "Current Patient Balance", "Total Charges", "Total Adjustments", "Total Payments", "Alternate Visit ID", "Visit Indicator", "Other Healthcare Provider"}
	default:
		return []string{}
	}
}

// getComponentNames gets the component names of the HL7 message.
func getComponentNames(segmentName, fieldName string) []string {
	if segmentName == "OBX" && fieldName == "Observation Value" {
		return []string{"Type", "Data", "Descriptor", "Unit"}
	}
	return []string{}
}