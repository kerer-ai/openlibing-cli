package engine

import (
	"os"
	"strings"
	"testing"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

func loadFixtureSPC(t *testing.T) *spc.SPCDefinition {
	t.Helper()
	data, err := os.ReadFile("testdata/valid.spc.yaml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	def, err := ParseSPC(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("parse fixture: %v", err)
	}
	return def
}

func TestValidate_AllRequiredPresent(t *testing.T) {
	def := loadFixtureSPC(t)
	params := map[string]interface{}{
		"project_id": "my-project",
		"limit":      20,
	}
	err := Validate(def, params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidate_MissingRequired(t *testing.T) {
	def := loadFixtureSPC(t)
	params := map[string]interface{}{
		"limit": 20,
	}
	err := Validate(def, params)
	if err == nil {
		t.Fatal("expected error for missing required param, got nil")
	}
	verr, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("expected *ValidationError, got %T", err)
	}
	if verr.Parameter != "project_id" {
		t.Errorf("error param = %q, want project_id", verr.Parameter)
	}
}

func TestValidate_EnumViolation(t *testing.T) {
	def := loadFixtureSPC(t)
	params := map[string]interface{}{
		"project_id": "p",
		"status":     "unknown",
	}
	err := Validate(def, params)
	if err == nil {
		t.Fatal("expected enum error, got nil")
	}
}

func TestValidate_RangeViolation(t *testing.T) {
	def := loadFixtureSPC(t)
	params := map[string]interface{}{
		"project_id": "p",
		"limit":      200,
	}
	err := Validate(def, params)
	if err == nil {
		t.Fatal("expected range error for limit=200, got nil")
	}
}

func TestValidate_RangeMinimum(t *testing.T) {
	def := loadFixtureSPC(t)
	params := map[string]interface{}{
		"project_id": "p",
		"limit":      -1,
	}
	err := Validate(def, params)
	if err == nil {
		t.Fatal("expected range error for limit=-1, got nil")
	}
}

func TestValidate_NoParams(t *testing.T) {
	// SPC with no parameters
	def := &spc.SPCDefinition{
		Name: "no-params",
		Type: "query",
	}
	err := Validate(def, map[string]interface{}{})
	if err != nil {
		t.Errorf("unexpected error for SPC with no params: %v", err)
	}
}

func TestValidate_OptionalParamSkipped(t *testing.T) {
	def := loadFixtureSPC(t)
	// Only provide required; optional 'status' and 'limit' should be fine
	params := map[string]interface{}{
		"project_id": "p",
	}
	err := Validate(def, params)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
