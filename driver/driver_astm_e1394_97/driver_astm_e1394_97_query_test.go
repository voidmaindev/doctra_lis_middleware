package driver_astm_e1394_97

import (
	"bytes"
	"io"
	"net"
	"strings"
	"testing"

	"github.com/voidmaindev/doctra_lis_middleware/services"
)

// Frames captured from a Snibe Maglumi query session (see _samples/snibe240p3.txt), control chars stripped by DataToBeReplaced.
const sampleQueryRawData = stx + "1H|\\^&||PSWD|MaglumiBio|||||Lis||P|E1394-97|20261021" + etx + "92" +
	stx + "2Q|1|^998877||ALL||||||||O" + etx + "6A" +
	stx + "3L|1|N" + etx + "06"

const sampleHeaderPrefix = "H|\\^&||PSWD|MaglumiBio|||||Lis||P|E1394-97|"

func TestGenerateASTMMessagesFromQuery(t *testing.T) {
	queryMessages := parseASTMMessages(sampleQueryRawData)

	tests := []struct {
		name         string
		dataToReturn []services.DeviceQueryDataToReturn
		wantFrames   int
		wantOrder    string
	}{
		{
			name:         "no orders yields empty H/P/L reply",
			dataToReturn: nil,
			wantFrames:   3,
			wantOrder:    "",
		},
		{
			name:         "one order",
			dataToReturn: []services.DeviceQueryDataToReturn{{Param: "ALP"}},
			wantFrames:   4,
			wantOrder:    "O|1|998877||^^^ALP|P|R",
		},
		{
			name:         "several orders are joined with repeat delimiter",
			dataToReturn: []services.DeviceQueryDataToReturn{{Param: "ALP"}, {Param: "TSH"}},
			wantFrames:   4,
			wantOrder:    "O|1|998877||^^^ALP\\^^^TSH|P|R",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frames := generateASTMMessagesFromQuery(queryMessages, tt.dataToReturn)

			if len(frames) != tt.wantFrames {
				t.Fatalf("got %d frames %q, want %d", len(frames), frames, tt.wantFrames)
			}
			if !strings.HasPrefix(frames[0], sampleHeaderPrefix) {
				t.Errorf("header frame = %q, want prefix %q", frames[0], sampleHeaderPrefix)
			}
			if frames[1] != "P|1" {
				t.Errorf("patient frame = %q, want %q", frames[1], "P|1")
			}
			if frames[len(frames)-1] != "L|1|N" {
				t.Errorf("terminator frame = %q, want %q", frames[len(frames)-1], "L|1|N")
			}
			if tt.wantOrder != "" && frames[2] != tt.wantOrder {
				t.Errorf("order frame = %q, want %q", frames[2], tt.wantOrder)
			}
		})
	}
}

func TestFormatASTMFrame(t *testing.T) {
	tests := []struct {
		name        string
		frameNumber int
		record      string
		want        string
	}{
		{
			name:        "header checksum matches device capture",
			frameNumber: 1,
			record:      "H|\\^&||PSWD|MaglumiBio|||||Lis||P|E1394-97|20261021",
			want:        stx + "1H|\\^&||PSWD|MaglumiBio|||||Lis||P|E1394-97|20261021" + cr + etx + "92" + cr + lf,
		},
		{
			name:        "terminator checksum matches device capture",
			frameNumber: 3,
			record:      "L|1|N",
			want:        stx + "3L|1|N" + cr + etx + "06" + cr + lf,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatASTMFrame(tt.frameNumber, tt.record)
			if got != tt.want {
				t.Errorf("formatASTMFrame() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestSendASTMFrames(t *testing.T) {
	tests := []struct {
		name     string
		frames   []string
		wantSent string
	}{
		{
			name:     "empty reply sends ENQ, H, P, L, EOT",
			frames:   []string{"H|\\^&", "P|1", "L|1|N"},
			wantSent: "\x05" + formatASTMFrame(1, "H|\\^&") + formatASTMFrame(2, "P|1") + formatASTMFrame(3, "L|1|N") + "\x04",
		},
		{
			name:     "no frames sends ENQ and EOT only",
			frames:   nil,
			wantSent: "\x05\x04",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hostConn, deviceConn := net.Pipe()
			received := make(chan []byte, 1)
			go runAckingDevice(t, deviceConn, received)

			err := sendASTMFrames(hostConn, tt.frames)
			if err != nil {
				t.Fatalf("sendASTMFrames() error = %v", err)
			}
			hostConn.Close()

			if got := string(<-received); got != tt.wantSent {
				t.Errorf("device received %q, want %q", got, tt.wantSent)
			}
		})
	}
}

// runAckingDevice plays a device that ACKs everything until the host sends EOT.
func runAckingDevice(t *testing.T, conn net.Conn, received chan<- []byte) {
	t.Helper()
	defer conn.Close()

	var all bytes.Buffer
	buf := make([]byte, 1<<10)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			if err != io.EOF {
				t.Errorf("device read error: %v", err)
			}
			break
		}
		all.Write(buf[:n])
		if buf[n-1] == rawDataEndString {
			break
		}
		if _, err := conn.Write([]byte{ack}); err != nil {
			t.Errorf("device ack error: %v", err)
			break
		}
	}
	received <- all.Bytes()
}
