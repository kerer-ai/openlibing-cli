package engine

import (
	"os"
	"strings"
	"testing"
)

func TestParseSPC_ValidFile(t *testing.T) {
	data, err := os.ReadFile("testdata/valid.spc.yaml")
	if err != nil {
		t.Fatalf("read test fixture: %v", err)
	}

	spc, err := ParseSPC(strings.NewReader(string(data)))
	if err != nil {
		t.Fatalf("ParseSPC failed: %v", err)
	}

	if spc.Name != "pipeline-list" {
		t.Errorf("name = %q", spc.Name)
	}
	if spc.Type != "query" {
		t.Errorf("type = %q", spc.Type)
	}
	if len(spc.Parameters) != 3 {
		t.Errorf("parameters = %d, want 3", len(spc.Parameters))
	}
}

func TestParseSPC_InvalidSyntax(t *testing.T) {
	data, err := os.ReadFile("testdata/invalid_syntax.spc.yaml")
	if err != nil {
		t.Fatalf("read test fixture: %v", err)
	}

	_, err = ParseSPC(strings.NewReader(string(data)))
	if err == nil {
		t.Fatal("expected error for invalid YAML, got nil")
	}
}

func TestParseSPCFile(t *testing.T) {
	spc, err := ParseSPCFile("testdata/valid.spc.yaml")
	if err != nil {
		t.Fatalf("ParseSPCFile failed: %v", err)
	}
	if spc.FilePath != "testdata/valid.spc.yaml" {
		t.Errorf("FilePath = %q", spc.FilePath)
	}
}

func TestParseSPCFile_NotFound(t *testing.T) {
	_, err := ParseSPCFile("testdata/does_not_exist.spc.yaml")
	if err == nil {
		t.Fatal("expected error for nonexistent file, got nil")
	}
}
