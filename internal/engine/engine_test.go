package engine

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/openlibing/openlibing-cli/internal/api"
	"github.com/openlibing/openlibing-cli/internal/config"
	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// mockRegistry implements SPCLookup for testing.
type mockRegistry struct {
	defs map[string]*spc.SPCDefinition
}

func (m *mockRegistry) Get(name string) (*spc.SPCDefinition, error) {
	def, ok := m.defs[name]
	if !ok {
		return nil, &ValidationError{Parameter: name, Message: "SPC not found: " + name}
	}
	return def, nil
}

func TestEngine_Execute_PipelineList_Success(t *testing.T) {
	// Setup mock HTTP server with fixture data
	fixture, _ := os.ReadFile("testdata/pipeline_list_response.json")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write(fixture)
	}))
	defer server.Close()

	// Setup engine
	cfg := config.DefaultConfig()
	cfg.Endpoint = server.URL
	client := api.NewClient(cfg, &config.Auth{})

	registry := &mockRegistry{
		defs: map[string]*spc.SPCDefinition{
			"pipeline-list": {
				Name:   "pipeline-list",
				Type:   "query",
				Source: spc.Source{
					Method:   "GET",
					Endpoint: "gateway/test/detail",
					QueryParams: map[string]string{
						"projectId": "{{.project_id}}",
						"pageSize":  "{{.limit}}",
					},
				},
				Output: spc.Output{
					Format: "table",
					Fields: []spc.Field{
						{Name: "id", Header: "ID", Path: "pipelineRunId"},
						{Name: "status", Header: "Status", Path: "status"},
					},
				},
				Parameters: []spc.Parameter{
					{Name: "project_id", Type: "string", Required: true},
					{Name: "limit", Type: "integer", Default: 10},
				},
			},
		},
	}

	engine := NewEngine(registry, client)

	// Execute
	result, err := engine.Execute("pipeline-list", map[string]interface{}{
		"project_id": "123",
		"limit":      5,
	})
	if err != nil {
		t.Fatalf("Execute failed: %v", err)
	}

	if len(result.Rows) != 2 {
		t.Fatalf("rows = %d, want 2", len(result.Rows))
	}
	if result.Rows[0]["id"] != "run-abc-123" {
		t.Errorf("first row id = %q", result.Rows[0]["id"])
	}
	if result.Rows[1]["status"] != "FAILED" {
		t.Errorf("second row status = %q", result.Rows[1]["status"])
	}
}

func TestEngine_Execute_MissingSPC(t *testing.T) {
	client := api.NewClient(config.DefaultConfig(), &config.Auth{})
	registry := &mockRegistry{defs: map[string]*spc.SPCDefinition{}}

	engine := NewEngine(registry, client)
	_, err := engine.Execute("nonexistent", nil)
	if err == nil {
		t.Fatal("expected error for missing SPC, got nil")
	}
}

func TestEngine_Execute_MissingRequiredParam(t *testing.T) {
	client := api.NewClient(config.DefaultConfig(), &config.Auth{})
	registry := &mockRegistry{
		defs: map[string]*spc.SPCDefinition{
			"test": {
				Name: "test",
				Source: spc.Source{
					Method:   "GET",
					Endpoint: "/test",
				},
				Output: spc.Output{Format: "json"},
				Parameters: []spc.Parameter{
					{Name: "required_param", Type: "string", Required: true},
				},
			},
		},
	}

	engine := NewEngine(registry, client)
	_, err := engine.Execute("test", map[string]interface{}{})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}
