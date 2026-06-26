package engine

import (
	"os"
	"strings"
	"testing"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

func loadSPCFile(t *testing.T, path string) *spc.SPCDefinition {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	def, err := ParseSPC(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("ParseSPC %s: %v", path, err)
	}
	return def
}

func TestPrGateMetrics_ParsesSuccessfully(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml")

	if def.Name != "pr-gate-metrics" {
		t.Errorf("name = %q, want pr-gate-metrics", def.Name)
	}
	if def.Type != "query" {
		t.Errorf("type = %q, want query", def.Type)
	}
	if def.Output.Format != "json" {
		t.Errorf("output.format = %q, want json", def.Output.Format)
	}
	if len(def.Output.Fields) != 0 {
		t.Errorf("fields = %d, want 0 (full output)", len(def.Output.Fields))
	}
	if len(def.Parameters) != 7 {
		t.Errorf("parameters = %d, want 7", len(def.Parameters))
	}
}

func TestPrGateMetricsSummary_ParsesSuccessfully(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics-summary.spc.yaml")

	if def.Name != "pr-gate-metrics-summary" {
		t.Errorf("name = %q, want pr-gate-metrics-summary", def.Name)
	}
	if def.Type != "query" {
		t.Errorf("type = %q, want query", def.Type)
	}
	if def.Output.Format != "table" {
		t.Errorf("output.format = %q, want table", def.Output.Format)
	}
	if len(def.Output.Fields) != 14 {
		t.Fatalf("fields = %d, want 14", len(def.Output.Fields))
	}

	// Verify the 14 field names are as expected
	expectedFields := []string{
		"repo", "branch",
		"e2e_p90", "e2e_avg", "e2e_p50",
		"build_p50", "build_p90", "build_avg", "build_queue_p90",
		"test_p90", "test_p50", "test_avg", "test_queue_p90",
		"check_p90",
	}
	for i, expected := range expectedFields {
		if def.Output.Fields[i].Name != expected {
			t.Errorf("field[%d].name = %q, want %q", i, def.Output.Fields[i].Name, expected)
		}
	}
}

func TestPrGateMetrics_Validate_AllRequired(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml")

	params := map[string]interface{}{
		"project_id": 4,
		"start_date": "2026-06-01",
		"end_date":   "2026-06-26",
	}
	err := Validate(def, params)
	if err != nil {
		t.Errorf("unexpected validation error: %v", err)
	}
}

func TestPrGateMetrics_Validate_MissingRequired(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml")

	params := map[string]interface{}{
		"project_id": 4,
		// start_date and end_date missing
	}
	err := Validate(def, params)
	if err == nil {
		t.Fatal("expected validation error for missing required params")
	}
}

func TestPrGateMetrics_Resolve_BodyTemplate(t *testing.T) {
	def := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml")

	params := map[string]interface{}{
		"project_id": 4,
		"start_date": "2026-06-01",
		"end_date":   "2026-06-26",
		"sort_field": "total",
		"sort_rule":  "desc",
		"page":       1,
		"page_size":  10,
	}

	resolved, err := Resolve(&def.Source, params)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}

	// Integer params must NOT be quoted in the body
	if strings.Contains(resolved.Body, `"projectId":"4"`) {
		t.Error("projectId should not be quoted (integer)")
	}
	if strings.Contains(resolved.Body, `"page":"1"`) {
		t.Error("page should not be quoted (integer)")
	}
	if strings.Contains(resolved.Body, `"pageSize":"10"`) {
		t.Error("pageSize should not be quoted (integer)")
	}

	// String params must be quoted
	if !strings.Contains(resolved.Body, `"startDate":"2026-06-01"`) {
		t.Error("startDate missing or incorrectly formatted")
	}
	if !strings.Contains(resolved.Body, `"category":"repo-pr-pipeline"`) {
		t.Error("category missing or incorrect")
	}
}

func TestPrGateMetricsSummary_SameParamsAsBase(t *testing.T) {
	base := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics.spc.yaml")
	summary := loadSPCFile(t, "../../embedded/spc/pr-gate-metrics-summary.spc.yaml")

	// Both SPCs should have identical parameters
	if len(base.Parameters) != len(summary.Parameters) {
		t.Fatalf("param count mismatch: base=%d summary=%d",
			len(base.Parameters), len(summary.Parameters))
	}
	for i := range base.Parameters {
		b, s := base.Parameters[i], summary.Parameters[i]
		if b.Name != s.Name || b.Type != s.Type || b.Required != s.Required {
			t.Errorf("param[%d] mismatch: base={%s,%s,%v} summary={%s,%s,%v}",
				i, b.Name, b.Type, b.Required, s.Name, s.Type, s.Required)
		}
	}

	// Both SPCs should have identical source (endpoint, method, headers, body)
	if base.Source.Endpoint != summary.Source.Endpoint {
		t.Error("endpoint mismatch")
	}
	if base.Source.Method != summary.Source.Method {
		t.Error("method mismatch")
	}
	if base.Source.Body != summary.Source.Body {
		t.Error("body mismatch")
	}
}
