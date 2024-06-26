package driver_hl7_231

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/voidmaindev/doctra_lis_middleware/model"
)

// hl7Message represents the entire HL7 message with segments stored in a map where keys are segment types.
type hl7Message struct {
	Segments map[string][]map[string]interface{} `json:"segments"`
}

// Unmarshal unmarshals the raw data.
func (d *Driver_hl7_231) Unmarshal(rawData string) (labDatas []*model.LabData, additionalData map[string]interface{}, err error) {
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

	hl7msg, err := parseHL7Message(rawData)
	if err != nil {
		fmt.Println("failed to parse HL7 message")
		return labDatas, nil, err
	}

	ackMsg := hl7msg.Segments["ACK"][0]["ACK"].(string)
	ackMsg = fmt.Sprintf("%c", rawDataStartString) + ackMsg + fmt.Sprintf("%c", rawDataEndString)
	additionalData = map[string]interface{}{"ACK": ackMsg}

	checkObrObx := len(hl7msg.Segments["OBR"]) > 1
	for _, obr := range hl7msg.Segments["OBR"] {
		for _, obx := range hl7msg.Segments["OBX"] {
			if !checkObrObx || obr["Set ID - OBR"] == obx["Set ID - OBX"] {
				barcode, err := getBarcodeForUnmarshalRawData(obr, hl7msg)
				if err != nil {
					fmt.Println("failed to get barcode for unmarshalRawData")
					return labDatas, nil, err
				}

				index, err := getIndexForUnmarshalRawData(obx)
				if err != nil {
					fmt.Println("failed to get index for unmarshalRawData")
					return labDatas, nil, err
				}

				param, err := getParamForUnmarshalRawData(obx)
				if err != nil {
					fmt.Println("failed to get param for unmarshalRawData")
					return labDatas, nil, err
				}

				result, err := getResultForUnmarshalRawData(obx)
				if err != nil {
					fmt.Println("failed to get result for unmarshalRawData")
					return labDatas, nil, err
				}

				unit, err := getUnitForUnmarshalRawData(obx)
				if err != nil {
					fmt.Println("failed to get unit for unmarshalRawData")
					return labDatas, nil, err
				}

				completedDate, err := getCompleteDateForUnmarshalRawData(obr, obx)
				if err != nil {
					fmt.Println("failed to get completed date for unmarshalRawData")
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
		}
	}

	return labDatas, additionalData, nil
}

// getBarcodeForUnmarshalRawData gets the barcode for unmarshalling the raw data.
func getBarcodeForUnmarshalRawData(obr map[string]interface{}, hl7msg *hl7Message) (string, error) {
	spm, ok := hl7msg.Segments["SPM"]
	if ok {
		for _, v := range spm {
			SpecimenID, ok := v["Specimen ID"]
			if ok {
				barcode1, ok := SpecimenID.(string)
				if ok {
					return barcode1, nil
				}
				SpecimenIDData, ok := SpecimenID.(map[string]interface{})
				if ok {
					barcodePre1, ok := SpecimenIDData["Component1"]
					if ok {
						barcode2, ok := barcodePre1.(string)
						if ok && barcode2 != "" {
							return barcode2, nil
						}
					}
					barcodePre2, ok := SpecimenIDData["Component2"]
					if ok {
						barcode2, ok := barcodePre2.(string)
						if ok && barcode2 != "" {
							return barcode2, nil
						}
						barcodePre2Pre, ok := barcodePre2.(map[string]string)
						if ok {
							barcode3, ok := barcodePre2Pre["Component2.Subcomponent1"]
							if ok {
								return barcode3, nil
							}
						}
					}
				}
			}
		}
	}

	barcode, ok := obr["Filler Order Number"].(string)
	if ok {
		return barcode, nil
	}

	return "", errors.New("failed to get barcode")
}

// getIndexForUnmarshalRawData gets the index for unmarshalling the raw data.
func getIndexForUnmarshalRawData(obx map[string]interface{}) (uint, error) {
	index1, err := strconv.Atoi(obx["Set ID - OBX"].(string))
	if err == nil {
		return uint(index1), nil
	}

	return 0, errors.New("failed to get index")
}

// getCompleteDateForUnmarshalRawData gets the completed date for unmarshalling the raw data.
func getCompleteDateForUnmarshalRawData(obr, obx map[string]interface{}) (time.Time, error) {
	completedDateString, ok := obr["Observation Date/Time"].(string)
	if !ok {
		if completedDateString, ok = obx["Date/Time of the Observation"].(string); !ok {
			completedDateString, ok = obx["Date/Time of the Analysis"].(string)
			if !ok {
				return time.Time{}, errors.New("failed to get completed date")
			}
		}
	}

	completedDate, err := time.Parse(completedDateFormat, completedDateString)
	if err != nil {
		return time.Time{}, err
	}

	return completedDate, nil
}

// getParamForUnmarshalRawData gets the param for unmarshalling the raw data.
func getParamForUnmarshalRawData(obx map[string]interface{}) (string, error) {
	param, ok := obx["Observation Identifier"].(string)
	if !ok {
		param1, ok := obx["Observation Identifier"].(map[string]interface{})
		if ok {
			param = param1["Component2"].(string)
		}
	}

	if param == "" {
		return "", errors.New("failed to get param")
	}

	param = strings.TrimPrefix(param, "*")

	return param, nil
}

// getResultForUnmarshalRawData gets the result for unmarshalling the raw data.
func getResultForUnmarshalRawData(obx map[string]interface{}) (string, error) {
	result1, ok := obx["Observation Value"].(string)
	if ok {
		return result1, nil
	}

	result2, ok := obx["Observation Value"].(map[string]interface{})
	if ok {
		return result2["Type"].(string), nil
	}

	return "", errors.New("failed to get result")
}

// getUnitForUnmarshalRawData gets the unit for unmarshalling the raw data.
func getUnitForUnmarshalRawData(obx map[string]interface{}) (string, error) {
	unit1, ok := obx["Units"].(string)
	if ok {
		return unit1, nil
	}

	return "", errors.New("failed to get unit")
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

			if strings.ContainsAny(field, delimiters.component+delimiters.repetition+delimiters.subComponent) && !(segmentName == "MSH" && fieldName == "Encoding Characters") {
				segmentFields[fieldName] = parseComplexField(field, delimiters, componentNames)
			} else {
				segmentFields[fieldName] = field
			}
		}
		message.Segments[segmentName] = append(message.Segments[segmentName], segmentFields)

		// Add ACK segment
		if segmentName == "MSH" {
			addAckSegment(message, fields)
		}
	}

	return message, nil
}

// addAckSegment adds the ACK segment to the HL7 message.
func addAckSegment(message *hl7Message, fields []string) {
	ackSegment := "MSH"
	fieldDefinitions := getFieldDefinitions("MSH")
	msgControlID := ""
	for i, fieldName := range fieldDefinitions {
		if fieldName == "Message Type" {
			ackSegment += "|ACK"
			continue
		}
		if fieldName == "Message Control ID" {
			msgControlID = fields[i+1]
		}
		ackSegment += "|" + fields[i+1]
	}
	ackSegment += "\rMSA|AA|" + msgControlID + "\r"
	message.Segments["ACK"] = []map[string]interface{}{{"ACK": ackSegment}}
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
	case "SPM":
		return []string{"Set ID - SPM", "Specimen ID", "Specimen Parent IDs", "Specimen Type", "Specimen Type Modifier", "Specimen Additives", "Specimen Collection Method", "Specimen Source Site", "Specimen Source Site Modifier", "Specimen Collection Site", "Specimen Role", "Specimen Collection Amount", "Grouped Specimen Count", "Specimen Description", "Specimen Handling Code", "Specimen Risk Code", "Specimen Collection Date/Time", "Specimen Received Date/Time", "Specimen Expiration Date/Time", "Specimen Availability", "Specimen Reject Reason", "Specimen Quality", "Specimen Appropriateness", "Specimen Condition", "Specimen Child Role"}
	case "ORC":
		return []string{"Order Control", "Placer Order Number", "Filler Order Number", "Placer Group Number", "Order Status", "Response Flag", "Quantity/Timing", "Parent", "Date/Time of Transaction", "Entered By", "Verified By", "Ordering Provider", "Enterer's Location", "Call Back Phone Number", "Order Effective Date/Time", "Order Control Code Reason", "Entering Organization", "Entering Device", "Action By", "Advanced Beneficiary Notice Code", "Ordering Facility Name", "Ordering Facility Address", "Ordering Facility Phone Number", "Ordering Provider Address"}
	case "NTE":
		return []string{"Set ID - NTE", "Source of Comment", "Comment", "Comment Type", "Entered By", "Entered Date/Time", "Effective Start Date", "Expiration Date", "Comment Completion Date"}
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
