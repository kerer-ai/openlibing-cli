package registry

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

func TestRegistry_LoadAll_FromEmbed(t *testing.T) {
	// Create an in-memory filesystem with 2 SPCs
	embedFS := fstest.MapFS{
		"spc/pipeline-list.spc.yaml": {
			Data: []byte(`name: pipeline-list
version: "1.0"
description: List pipelines
type: query
category: pipeline
source:
  method: GET
  endpoint: /test
output:
  format: table`),
		},
		"spc/pipeline-detail.spc.yaml": {
			Data: []byte(`name: pipeline-detail
version: "1.0"
description: Get pipeline detail
type: query
category: pipeline
source:
  method: GET
  endpoint: /test/detail
output:
  format: table`),
		},
	}

	r := NewRegistry()
	err := r.LoadAll(embedFS, "spc")
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	all := r.ListAll()
	if len(all) != 2 {
		t.Fatalf("ListAll = %d, want 2", len(all))
	}

	// Test Get
	def, err := r.Get("pipeline-list")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if def.Origin != "builtin" {
		t.Errorf("origin = %q, want builtin", def.Origin)
	}

	// Test ListByCategory
	cat := r.ListByCategory("pipeline")
	if len(cat) != 2 {
		t.Errorf("ListByCategory(pipeline) = %d, want 2", len(cat))
	}
}

func TestRegistry_Get_NotFound(t *testing.T) {
	r := NewRegistry()
	_, err := r.Get("nonexistent")
	if err == nil {
		t.Fatal("expected error for missing SPC, got nil")
	}
}

func TestRegistry_UserOverrides(t *testing.T) {
	embedFS := fstest.MapFS{
		"spc/test.spc.yaml": {
			Data: []byte(`name: test
version: "1.0"
description: Built-in version
type: query
category: test
source:
  method: GET
  endpoint: /builtin
output:
  format: json`),
		},
	}

	// Create a temp user SPC directory that overrides the same name
	tmpDir := t.TempDir()
	userSPCDir := filepath.Join(tmpDir, ".openlibing", "spc")
	os.MkdirAll(userSPCDir, 0755)
	os.WriteFile(filepath.Join(userSPCDir, "test.spc.yaml"), []byte(`name: test
version: "1.0"
description: User override version
type: query
category: test
source:
  method: GET
  endpoint: /user-override
output:
  format: json`), 0644)

	// Override HOME for the test
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	r := NewRegistry()
	err := r.LoadAll(embedFS, "spc")
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	def, err := r.Get("test")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if def.Origin != "user" {
		t.Errorf("origin = %q, want user", def.Origin)
	}
	if def.Source.Endpoint != "/user-override" {
		t.Errorf("endpoint = %q, want /user-override (user override)", def.Source.Endpoint)
	}
}

func TestResolveSearch_Scores(t *testing.T) {
	defs := []*spc.SPCDefinition{
		{Name: "pipeline-list", Description: "List pipeline runs", Category: "pipeline", Tags: []string{"ci", "status"}},
		{Name: "codecheck-report", Description: "Code check issues", Category: "codecheck", Tags: []string{"quality"}},
	}

	results := ResolveSearch("pipeline", defs)
	if len(results) != 1 {
		t.Fatalf("search results = %d, want 1", len(results))
	}
	if results[0].Name != "pipeline-list" {
		t.Errorf("top result = %q, want pipeline-list", results[0].Name)
	}
}

func TestResolveSearch_NoMatch(t *testing.T) {
	defs := []*spc.SPCDefinition{
		{Name: "pipeline-list", Description: "List pipelines", Category: "pipeline"},
	}
	results := ResolveSearch("xyz-nonexistent", defs)
	if len(results) != 0 {
		t.Errorf("results = %d, want 0", len(results))
	}
}
