package driver

import (
	"testing"

	"github.com/voidmaindev/doctra_lis_middleware/log"
	"github.com/voidmaindev/doctra_lis_middleware/tcp"
)

func TestNewDriverKnowsEveryListedDriver(t *testing.T) {
	logger := &log.Logger{Disabled: true}

	for _, name := range Drivers() {
		t.Run(name, func(t *testing.T) {
			created, err := NewDriver(name, logger, nil, nil)
			if err != nil {
				t.Fatalf("NewDriver(%q) error = %v", name, err)
			}
			if created == nil {
				t.Fatalf("NewDriver(%q) returned no driver", name)
			}
		})
	}
}

func TestGetRawDatas(t *testing.T) {
	const (
		enq = "\x05"
		eot = "\x04"
	)

	tests := []struct {
		name       string
		driverName string
		reads      []string
		want       []string
	}{
		{
			name:       "Accre 8 greeting is ignored and blocks sent in one read are separated",
			driverName: "astm_wondfo_accre_8",
			reads:      []string{"Network is connected!", enq + "first" + eot + enq + "second" + eot},
			want:       []string{"first", "second"},
		},
		{
			name:       "Accre 8 block split over two reads is joined",
			driverName: "astm_wondfo_accre_8",
			reads:      []string{enq + "fir", "st" + eot},
			want:       []string{"first"},
		},
	}

	logger := &log.Logger{Disabled: true}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			created, err := NewDriver(tt.driverName, logger, nil, nil)
			if err != nil {
				t.Fatalf("NewDriver(%q) error = %v", tt.driverName, err)
			}

			got := []string{}
			prevData := &tcp.PrevData{}
			for _, read := range tt.reads {
				got = append(got, GetRawDatas(created, read, prevData)...)
			}

			if len(got) != len(tt.want) {
				t.Fatalf("got %d raw datas %q, want %d", len(got), got, len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Errorf("raw data %d = %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}
