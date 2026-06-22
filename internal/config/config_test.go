package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	if cfg.Endpoint != "https://www.openlibing.com" {
		t.Errorf("endpoint = %q", cfg.Endpoint)
	}
	if cfg.Defaults.Limit != 10 {
		t.Errorf("defaults.limit = %d, want 10", cfg.Defaults.Limit)
	}
	if cfg.Output.Format != "table" {
		t.Errorf("output.format = %q, want table", cfg.Output.Format)
	}
	if !cfg.Output.Color {
		t.Error("output.color should default to true")
	}
}

func TestLoadConfig_DefaultsWhenNoFile(t *testing.T) {
	// Override HOME to a temp dir with no config
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Endpoint != "https://www.openlibing.com" {
		t.Errorf("endpoint = %q, want default", cfg.Endpoint)
	}
}

func TestLoadConfig_ReadsFile(t *testing.T) {
	tmpDir := t.TempDir()
	openlibingDir := filepath.Join(tmpDir, ".openlibing")
	os.MkdirAll(openlibingDir, 0755)

	configYAML := `
endpoint: https://custom.openlibing.com
defaults:
  project_id: my-project
  limit: 25
output:
  format: json
  color: false
`
	os.WriteFile(filepath.Join(openlibingDir, "config.yaml"), []byte(configYAML), 0644)
	t.Setenv("HOME", tmpDir)

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}

	if cfg.Endpoint != "https://custom.openlibing.com" {
		t.Errorf("endpoint = %q", cfg.Endpoint)
	}
	if cfg.Defaults.ProjectID != "my-project" {
		t.Errorf("defaults.project_id = %q", cfg.Defaults.ProjectID)
	}
	if cfg.Defaults.Limit != 25 {
		t.Errorf("defaults.limit = %d, want 25", cfg.Defaults.Limit)
	}
	if cfg.Output.Format != "json" {
		t.Errorf("output.format = %q, want json", cfg.Output.Format)
	}
}

func TestLoadConfig_DirIsCreated(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	dir, err := ConfigDir()
	if err != nil {
		t.Fatalf("ConfigDir failed: %v", err)
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Error("ConfigDir should have created the directory")
	}

	spcDir, err := SPCUserDir()
	if err != nil {
		t.Fatalf("SPCUserDir failed: %v", err)
	}
	if _, err := os.Stat(spcDir); os.IsNotExist(err) {
		t.Error("SPCUserDir should have created the spc directory")
	}
}

func TestLoadAuth_EmptyWhenNoFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("HOME", tmpDir)

	auth, err := LoadAuth()
	if err != nil {
		t.Fatalf("LoadAuth failed: %v", err)
	}
	if auth.OpenLibing.TokenType != "Bearer" {
		t.Errorf("token_type = %q, want Bearer", auth.OpenLibing.TokenType)
	}
}

func TestLoadAuth_ReadsFile(t *testing.T) {
	tmpDir := t.TempDir()
	openlibingDir := filepath.Join(tmpDir, ".openlibing")
	os.MkdirAll(openlibingDir, 0755)

	authYAML := `
openlibing:
  token: test-token-123
  token_type: Bearer
llm:
  provider: anthropic
  api_key: sk-ant-test
`
	os.WriteFile(filepath.Join(openlibingDir, "auth.yaml"), []byte(authYAML), 0600)
	t.Setenv("HOME", tmpDir)

	auth, err := LoadAuth()
	if err != nil {
		t.Fatalf("LoadAuth failed: %v", err)
	}
	if auth.OpenLibing.Token != "test-token-123" {
		t.Errorf("token = %q", auth.OpenLibing.Token)
	}
	if auth.LLM.Provider != "anthropic" {
		t.Errorf("llm.provider = %q", auth.LLM.Provider)
	}
}
