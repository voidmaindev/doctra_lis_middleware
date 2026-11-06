package driver_astm_wondfo_accre_8

import (
	"testing"
	"time"

	"github.com/voidmaindev/doctra_lis_middleware/log"
)

// Blocks captured from a Wondfo Accre 8 (see _samples/jc/wondfo_accre_8.log), ENQ/EOT removed by the framing.
// The control characters are literal text ("\x02"), exactly as the analyzer sends them.
const (
	sampleHeaderFrames = `\x021H|\x5e&||| ACCRE8 |||||Host||1|20090119131415` + "\r" + `\x0353` + "\r\n" +
		`\x022P|1|||||||||||||^` + "\r" + `\x03E9` + "\r\n"
	sampleTerminatorFrame = `\x026L|1|N` + "\r" + `\x030A` + "\r\n"

	sampleDDimer = sampleHeaderFrames +
		`\x023O|1|128123^nadiraze^1^2|||R|20261103 17:30:48|||||||||Serum||||||||||O` + "\r" + `\x0355` + "\r\n" +
		`\x024R|1|^^^D-Dimer|46.33|μg/mL(FEU)|^\x5e|N||F||||20261103 17:30:48` + "\r" + `\x0332` + "\r\n" +
		sampleTerminatorFrame

	sampleWithoutResult = sampleHeaderFrames +
		`\x023O|1|128123^mosize^1^2|||R|20261031 15:01:45|||||||||Serum||||||||||O` + "\r" + `\x0355` + "\r\n"

	sampleScannedBarcode = sampleHeaderFrames +
		`\x023O|1|128123^typed^1^2|998877||R|20261103 16:13:43|||||||||Serum||||||||||O` + "\r" + `\x0355` + "\r\n" +
		`\x024R|1|^^^Vit D|41|ng/mL|^\x5e|N||F||||20261103 16:13:43` + "\r" + `\x0332` + "\r\n" +
		`\x025R|2|^^^TSH|4.096|μIU/mL|^\x5e|N||F||||20261103 16:13:43` + "\r" + `\x0332` + "\r\n" +
		sampleTerminatorFrame

	// The same message as a firmware that sends real control characters and the manual's compact date would send it.
	sampleRealControlChars = "\x021H|\\^&||| ACCRE8 |||||Host||1|20090119131415\r\x0353\r\n" +
		"\x022P|1|||||||||||||^\r\x03E9\r\n" +
		"\x023O|1|128123^cincadze^1^2|||R|20261103161343|||||||||Serum||||||||||O\r\x0355\r\n" +
		"\x024R|1|^^^Vit D|41|ng/mL|^\\^|N||F||||20261103161343\r\x0332\r\n" +
		"\x025L|1|N\r\x030A\r\n"

	sampleBadDate = sampleHeaderFrames +
		`\x023O|1|128123^ruso^1^2|||R||||||||||Serum||||||||||O` + "\r" + `\x0355` + "\r\n" +
		`\x024R|1|^^^Vit D|27.42|ng/mL|^\x5e|N||F||||yesterday` + "\r" + `\x0332` + "\r\n"

	sampleResultWithoutOrder = sampleHeaderFrames +
		`\x024R|1|^^^Vit D|27.42|ng/mL|^\x5e|N||F||||20261103 16:13:43` + "\r" + `\x0332` + "\r\n"
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
	tests := []struct {
		name    string
		rawData string
		want    []wantLabData
		wantErr bool
	}{
		{
			name:    "single result with UTF-8 unit, sample ID typed by the operator",
			rawData: sampleDDimer,
			want: []wantLabData{
				{"nadiraze", 1, "D-Dimer", "46.33", "μg/mL(FEU)", time.Date(2026, 11, 3, 17, 30, 48, 0, time.UTC)},
			},
		},
		{
			name:    "order without a result record yields no lab data and no error",
			rawData: sampleWithoutResult,
			want:    nil,
		},
		{
			name:    "scanned specimen barcode wins over the typed sample ID, results are indexed",
			rawData: sampleScannedBarcode,
			want: []wantLabData{
				{"998877", 1, "Vit D", "41", "ng/mL", time.Date(2026, 11, 3, 16, 13, 43, 0, time.UTC)},
				{"998877", 2, "TSH", "4.096", "μIU/mL", time.Date(2026, 11, 3, 16, 13, 43, 0, time.UTC)},
			},
		},
		{
			name:    "real control characters and compact date",
			rawData: sampleRealControlChars,
			want: []wantLabData{
				{"cincadze", 1, "Vit D", "41", "ng/mL", time.Date(2026, 11, 3, 16, 13, 43, 0, time.UTC)},
			},
		},
		{
			name:    "unparseable completion date is an error",
			rawData: sampleBadDate,
			wantErr: true,
		},
		{
			name:    "result without an order record is an error",
			rawData: sampleResultWithoutOrder,
			wantErr: true,
		},
		{
			name:    "text without frames yields nothing",
			rawData: "Network is connected!",
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

func TestSplitRecords(t *testing.T) {
	tests := []struct {
		name    string
		rawData string
		want    []string
	}{
		{
			name:    "literal control characters, frame numbers and checksums are stripped",
			rawData: sampleHeaderFrames + sampleTerminatorFrame,
			want:    []string{"H|^&||| ACCRE8 |||||Host||1|20090119131415", "P|1|||||||||||||^", "L|1|N"},
		},
		{
			name:    "intermediate frame is joined with the next frame",
			rawData: "\x021R|1|^^^Vit D|4\x17AB\r\n\x0221|ng/mL\r\x03CD\r\n",
			want:    []string{"R|1|^^^Vit D|41|ng/mL"},
		},
		{
			name:    "no frames",
			rawData: "",
			want:    []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitRecords(tt.rawData)

			if len(got) != len(tt.want) {
				t.Fatalf("got %d records %q, want %d", len(got), got, len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("record %d = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
