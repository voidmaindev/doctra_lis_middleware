package driver_text_dri_chem_nx600

import (
	"testing"
	"time"

	"github.com/voidmaindev/doctra_lis_middleware/log"
)

// Records captured from a DRI-CHEM NX600 (see _samples/jc/dri_chem_nx6000.log), STX/ETX/BCC removed by the framing.
const (
	sampleOneTest = "R,NORMAL ,03-11-2026,17:16,8            ,RUSO AL RE   ,             ,00,9,999,01,01," +
		"GLU-PS  ,=,5.3      mmol/l,01,3.9  ,6.1  ,           "

	sampleSixTests = "R,NORMAL ,31-10-2026,14:44,2            ,SUXANSKAIA IR,             ,01,9,999,01,06," +
		"TCHO-PS ,=,7.38     mmol/l,01,3.88 ,5.66 ,H          ," +
		"HDLC-PS ,=,2.12     mmol/l,01,1.03 ,1.84 ,H          ," +
		"TG-PS   ,=,0.77     mmol/l,01,0.56 ,1.68 ,           ," +
		"GLU-PS  ,=,6.0      mmol/l,01,3.9  ,6.1  ,           ," +
		"AST-PS  ,=,20       U/l   ,01,8    ,38   ,           ," +
		"ALT-PS  ,=,15       U/l   ,01,4    ,44   ,           "

	sampleTestStart = "S,NORMAL ,04-11-2026,12:23,1            ,CHKADUA      ,             ,01"
)

// Records built from the format in the NX600 Communication Interface manual.
const (
	sampleControl = "R,CONTROL,03-11-2026,17:16,8            ,             ,QP-L         ,00,9,999,01,01," +
		"GLU-PS  ,=,5.3      mmol/l,01,3.9  ,6.1  ,           "

	sampleAboveRangeISODate = "R,NORMAL ,2026-11-03,10:50,2026110301   ,998877       ,Taro Fuji    ,02,1,003,01,02," +
		"GLU-PS  ,=,75       mg/dl ,01,50.0 ,100.0,           ," +
		"AMYL-PS ,>,1500     U/l   ,01,500  ,1500 ,H#         "

	sampleBlankSampleID = "R,NORMAL ,03-11-2026,17:16,8            ,             ,             ,00,9,999,01,01," +
		"GLU-PS  ,=,5.3      mmol/l,01,3.9  ,6.1  ,           "

	sampleTruncated = "R,NORMAL ,03-11-2026,17:16,8            ,RUSO AL RE   ,             ,00,9,999,01,02," +
		"GLU-PS  ,=,5.3      mmol/l,01,3.9  ,6.1  ,           "

	sampleBadTestCount = "R,NORMAL ,03-11-2026,17:16,8            ,RUSO AL RE   ,             ,00,9,999,01,xx"

	sampleBadDate = "R,NORMAL ,11/03/2026,17:16,8            ,RUSO AL RE   ,             ,00,9,999,01,01," +
		"GLU-PS  ,=,5.3      mmol/l,01,3.9  ,6.1  ,           "
)

type wantLabData struct {
	barcode       string
	index         uint
	param         string
	result        string
	unit          string
	completedDate time.Time
}

func TestUnmarshal(t *testing.T) {
	glucoseTime := time.Date(2026, 11, 3, 17, 16, 0, 0, time.UTC)
	panelTime := time.Date(2026, 10, 31, 14, 44, 0, 0, time.UTC)
	isoTime := time.Date(2026, 11, 3, 10, 50, 0, 0, time.UTC)

	tests := []struct {
		name    string
		rawData string
		want    []wantLabData
		wantErr bool
	}{
		{
			name:    "one test, day-first date",
			rawData: sampleOneTest,
			want:    []wantLabData{{"RUSO AL RE", 1, "GLU-PS", "5.3", "mmol/l", glucoseTime}},
		},
		{
			name:    "six tests",
			rawData: sampleSixTests,
			want: []wantLabData{
				{"SUXANSKAIA IR", 1, "TCHO-PS", "7.38", "mmol/l", panelTime},
				{"SUXANSKAIA IR", 2, "HDLC-PS", "2.12", "mmol/l", panelTime},
				{"SUXANSKAIA IR", 3, "TG-PS", "0.77", "mmol/l", panelTime},
				{"SUXANSKAIA IR", 4, "GLU-PS", "6.0", "mmol/l", panelTime},
				{"SUXANSKAIA IR", 5, "AST-PS", "20", "U/l", panelTime},
				{"SUXANSKAIA IR", 6, "ALT-PS", "15", "U/l", panelTime},
			},
		},
		{
			name:    "leading STX left by a BCC byte equal to STX is dropped",
			rawData: "\x02" + sampleOneTest,
			want:    []wantLabData{{"RUSO AL RE", 1, "GLU-PS", "5.3", "mmol/l", glucoseTime}},
		},
		{
			name:    "test start info is ignored",
			rawData: sampleTestStart,
			want:    nil,
		},
		{
			name:    "control measurement is ignored",
			rawData: sampleControl,
			want:    nil,
		},
		{
			name:    "ISO date, sign other than = is kept in the result",
			rawData: sampleAboveRangeISODate,
			want: []wantLabData{
				{"998877", 1, "GLU-PS", "75", "mg/dl", isoTime},
				{"998877", 2, "AMYL-PS", ">1500", "U/l", isoTime},
			},
		},
		{
			name:    "blank sample ID falls back to the sample No.",
			rawData: sampleBlankSampleID,
			want:    []wantLabData{{"8", 1, "GLU-PS", "5.3", "mmol/l", glucoseTime}},
		},
		{
			name:    "fewer tests than announced is an error",
			rawData: sampleTruncated,
			wantErr: true,
		},
		{
			name:    "invalid number of tests is an error",
			rawData: sampleBadTestCount,
			wantErr: true,
		},
		{
			name:    "unknown date format is an error",
			rawData: sampleBadDate,
			wantErr: true,
		},
		{
			name:    "empty record is ignored",
			rawData: "",
			want:    nil,
		},
	}

	driver := NewDriver(&log.Logger{Disabled: true}, nil, nil)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			labDatas, _, err := driver.Unmarshal(tt.rawData)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if len(labDatas) != len(tt.want) {
				t.Fatalf("got %d lab datas, want %d", len(labDatas), len(tt.want))
			}

			for i, want := range tt.want {
				got := wantLabData{
					labDatas[i].Barcode, labDatas[i].Index, labDatas[i].Param,
					labDatas[i].Result, labDatas[i].Unit, labDatas[i].CompletedDate,
				}
				if got != want {
					t.Errorf("lab data %d = %+v, want %+v", i, got, want)
				}
			}
		})
	}
}

func TestSplitResultAndUnit(t *testing.T) {
	tests := []struct {
		name       string
		field      string
		wantResult string
		wantUnit   string
	}{
		{"result and unit", "5.3      mmol/l", "5.3", "mmol/l"},
		{"padded unit", "20       U/l   ", "20", "U/l"},
		{"full width result", "123456.78mg/dl ", "123456.78", "mg/dl"},
		{"error marker without unit", "****           ", "****", ""},
		{"short field", "75", "75", ""},
		{"empty field", "", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, unit := splitResultAndUnit(tt.field)

			if result != tt.wantResult || unit != tt.wantUnit {
				t.Errorf("splitResultAndUnit(%q) = (%q, %q), want (%q, %q)", tt.field, result, unit, tt.wantResult, tt.wantUnit)
			}
		})
	}
}
