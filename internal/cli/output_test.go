package cli

import (
	"strings"
	"testing"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

func TestFormatJSON(t *testing.T) {
	rows := []map[string]interface{}{
		{"id": "run-1", "status": "SUCCESS"},
	}
	result, err := FormatJSON(rows)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}
	if !strings.Contains(result, "run-1") {
		t.Errorf("JSON output missing expected value: %s", result)
	}
}

func TestFormatJSON_NilRows(t *testing.T) {
	result, err := FormatJSON(nil)
	if err != nil {
		t.Fatalf("FormatJSON failed: %v", err)
	}
	if result != "[]" {
		t.Errorf("empty JSON = %q, want []", result)
	}
}

func TestFormatTable_NoResults(t *testing.T) {
	result := FormatTable([]map[string]interface{}{}, nil)
	if !strings.Contains(result, "no results") {
		t.Errorf("empty table = %q", result)
	}
}

func TestFormatTable_WithRows(t *testing.T) {
	rows := []map[string]interface{}{
		{"id": "abc-123", "status": "SUCCESS"},
		{"id": "def-456", "status": "FAILED"},
	}
	fields := []spc.Field{
		{Name: "id", Header: "ID"},
		{Name: "status", Header: "Status"},
	}
	result := FormatTable(rows, fields)
	if !strings.Contains(result, "abc-123") {
		t.Errorf("table output missing row value: %s", result)
	}
	if !strings.Contains(result, "FAILED") {
		t.Errorf("table output missing row value: %s", result)
	}
}

func TestTransformDuration_Milliseconds(t *testing.T) {
	result := TransformDuration(float64(500))
	if result != "500ms" {
		t.Errorf("500ms = %q", result)
	}
}

func TestTransformDuration_Seconds(t *testing.T) {
	result := TransformDuration(float64(5500))
	if result != "5s" {
		t.Errorf("5500ms = %q", result)
	}
}

func TestTransformDuration_MinutesAndSeconds(t *testing.T) {
	result := TransformDuration(float64(125000))
	if result != "2m 5s" {
		t.Errorf("125000ms = %q", result)
	}
}

func TestTransformDuration_HoursAndMinutes(t *testing.T) {
	result := TransformDuration(float64(7500000))
	if result != "2h 5m" {
		t.Errorf("7500000ms = %q", result)
	}
}

func TestFormatResult_DispatchesByFormat(t *testing.T) {
	tests := []struct {
		name   string
		result *spc.Result
		check  func(string) bool
	}{
		{
			name: "json",
			result: &spc.Result{
				Format: "json",
				Rows:   []map[string]interface{}{{"key": "val"}},
			},
			check: func(s string) bool { return strings.Contains(s, "\"key\"") },
		},
		{
			name: "raw",
			result: &spc.Result{
				Format: "raw",
				Raw:    []byte("raw output"),
			},
			check: func(s string) bool { return s == "raw output" },
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out, err := FormatResult(tt.result)
			if err != nil {
				t.Fatalf("FormatResult failed: %v", err)
			}
			if !tt.check(out) {
				t.Errorf("unexpected output: %s", out)
			}
		})
	}
}
