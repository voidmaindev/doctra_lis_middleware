package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/voidmaindev/doctra_lis_middleware/log"
)

func TestDeviceQueryService_Query(t *testing.T) {
	tests := []struct {
		name        string
		status      int
		contentType string
		body        string
		wantParams  []string
		wantErr     string
	}{
		{
			name:        "json with indicators",
			status:      http.StatusOK,
			contentType: "application/json; charset=utf-8",
			body:        `{"indicators":[{"indicator":"ALP"},{"indicator":"TSH"}]}`,
			wantParams:  []string{"ALP", "TSH"},
		},
		{
			name:        "json with empty indicators returns no data and no error",
			status:      http.StatusOK,
			contentType: "application/json",
			body:        `{"indicators":[]}`,
			wantParams:  nil,
		},
		{
			name:        "html 200 is rejected as non-JSON",
			status:      http.StatusOK,
			contentType: "text/html",
			body:        "<html>welcome</html>",
			wantErr:     "non-JSON content type",
		},
		{
			name:        "http 500 is an error",
			status:      http.StatusInternalServerError,
			contentType: "application/json",
			body:        `{}`,
			wantErr:     "HTTP error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotRequest DeviceQueryRequestBody
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if err := json.NewDecoder(r.Body).Decode(&gotRequest); err != nil {
					t.Errorf("decode request: %v", err)
				}
				w.Header().Set(contentTypeHeader, tt.contentType)
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.body))
			}))
			defer server.Close()

			service := NewDeviceQueryService(&log.Logger{Disabled: true}, server.URL, "SN-1")
			got, err := service.Query("998877")

			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("Query() error = %v, want containing %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("Query() unexpected error: %v", err)
			}
			if gotRequest.Barcode != "998877" || gotRequest.HardwareSN != "SN-1" {
				t.Errorf("request body = %+v, want barcode 998877 and hardware_sn SN-1", gotRequest)
			}
			if len(got) != len(tt.wantParams) {
				t.Fatalf("Query() returned %d params %+v, want %d", len(got), got, len(tt.wantParams))
			}
			for i, param := range tt.wantParams {
				if got[i].Param != param {
					t.Errorf("param[%d] = %q, want %q", i, got[i].Param, param)
				}
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name      string
		value     string
		maxLength int
		want      string
	}{
		{name: "shorter than limit is unchanged", value: "abc", maxLength: 5, want: "abc"},
		{name: "equal to limit is unchanged", value: "abcde", maxLength: 5, want: "abcde"},
		{name: "longer than limit is cut with ellipsis", value: "abcdefgh", maxLength: 5, want: "abcde..."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := truncate(tt.value, tt.maxLength); got != tt.want {
				t.Errorf("truncate() = %q, want %q", got, tt.want)
			}
		})
	}
}
