# OpenLibing CLI MVP — 实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an AI-native CLI for querying OpenLibing CI/CD pipeline data using SPC-First architecture.

**Architecture:** SPC-First — all query capabilities defined as `.spc.yaml` files. A thin Go binary (Cobra) loads SPCs via embedded filesystem + user/project directories, executes them through a 7-step engine (Parse→Validate→Resolve→Call→Extract→Transform→Output), and exposes 4 meta-commands (run/list/inspect/chat).

**Tech Stack:** Go 1.22+, Cobra, gopkg.in/yaml.v3, tidwall/gjson, text/template, net/http, go:embed

## Global Constraints

- Go version >= 1.22
- Minimum external dependencies: cobra, yaml.v3, gjson only
- All built-in SPCs embedded via `//go:embed` — zero external file requirement at runtime
- `pkg/spc/types.go` is the single shared type definition — no circular imports
- `api/` package must not import `engine` or `registry`
- SPC type support: `query` only in MVP
- Test coverage: `pkg/spc/`, `internal/engine/`, `internal/registry/`, `internal/api/`
- No hardcoded domain subcommands — all query capabilities from SPC files

---

### Task 1: Project Scaffolding

**Files:**
- Create: `go.mod`
- Create: `Makefile`

**Interfaces:**
- Consumes: nothing
- Produces: Go module `github.com/openlibing/openlibing-cli`, directory structure ready for code

- [ ] **Step 1: Initialize Go module**

```bash
cd /home/wangsike/workspace/openlibing/openlibing-cli
go mod init github.com/openlibing/openlibing-cli
```

Expected: `go.mod` created with module path.

- [ ] **Step 2: Add core dependencies**

```bash
go get github.com/spf13/cobra@latest
go get gopkg.in/yaml.v3@latest
go get github.com/tidwall/gjson@latest
```

Expected: `go.mod` and `go.sum` updated with 3 dependencies.

- [ ] **Step 3: Create directory structure**

```bash
mkdir -p cmd/openlibing
mkdir -p internal/engine
mkdir -p internal/registry
mkdir -p internal/api
mkdir -p internal/cli
mkdir -p internal/chat
mkdir -p internal/config
mkdir -p embed/spc
mkdir -p pkg/spc
mkdir -p docs/superpowers/specs
mkdir -p docs/superpowers/plans
```

Expected: All directories created.

- [ ] **Step 4: Create Makefile**

File: `Makefile`
```makefile
.PHONY: build test run clean lint

BINARY=openlibing

build:
	go build -o bin/$(BINARY) ./cmd/openlibing/

test:
	go test ./... -v -count=1

test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out

run: build
	./bin/$(BINARY)

clean:
	rm -rf bin/

lint:
	go vet ./...
```

- [ ] **Step 5: Verify build compiles (empty main)**

File: `cmd/openlibing/main.go`
```go
package main

import "fmt"

func main() {
	fmt.Println("openlibing-cli")
}
```

Run: `make build`
Expected: `bin/openlibing` binary created, runs and prints "openlibing-cli".

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "chore: scaffold project structure and Makefile"
```

---

### Task 2: Core SPC Types

**Files:**
- Create: `pkg/spc/types.go`
- Create: `pkg/spc/types_test.go`

**Interfaces:**
- Consumes: nothing
- Produces:
  - `type SPCDefinition struct { ... }` — all SPC YAML fields + runtime metadata
  - `type Parameter struct { ... }` — with Name, Type, Required, Default, Enum, Validation
  - `type Validation struct { ... }` — Min *int, Max *int
  - `type Source struct { ... }` — Method, Endpoint, QueryParams, Headers, Body
  - `type Output struct { ... }` — Format, Fields
  - `type Field struct { ... }` — Name, Header, Path, Width, Transform
  - `type AIConfig struct { ... }` — PromptHint, NaturalLanguage, ResultHint
  - `type Example struct { ... }` — Command, Description
  - `type Result struct { ... }` — Format, Rows ([]map[string]interface{}), Raw ([]byte)

- [ ] **Step 1: Write the types file**

File: `pkg/spc/types.go`
```go
package spc

// SPCDefinition is the complete parsed representation of a .spc.yaml file.
type SPCDefinition struct {
	Name        string      `yaml:"name"`
	Version     string      `yaml:"version"`
	Description string      `yaml:"description"`
	Type        string      `yaml:"type"`    // query | action | workflow
	Category   string      `yaml:"category"`
	Tags       []string    `yaml:"tags"`
	Parameters []Parameter `yaml:"parameters"`
	Source     Source      `yaml:"source"`
	Output     Output      `yaml:"output"`
	AI         *AIConfig   `yaml:"ai,omitempty"`
	Examples   []Example   `yaml:"examples,omitempty"`

	// Runtime metadata (not from YAML)
	Origin   string `yaml:"-"` // "builtin" | "user" | "project"
	FilePath string `yaml:"-"` // absolute path to the SPC file
}

// Parameter defines an input parameter for a Super Power.
type Parameter struct {
	Name        string      `yaml:"name"`
	Type        string      `yaml:"type"` // string | integer | boolean
	Required    bool        `yaml:"required"`
	Default     interface{} `yaml:"default,omitempty"`
	Description string      `yaml:"description,omitempty"`
	Enum        []string    `yaml:"enum,omitempty"`
	Validation  *Validation `yaml:"validation,omitempty"`
}

// Validation defines numeric range constraints for a parameter.
type Validation struct {
	Min *int `yaml:"min,omitempty"`
	Max *int `yaml:"max,omitempty"`
}

// Source defines the HTTP data source for a Super Power.
type Source struct {
	Method      string            `yaml:"method"` // GET | POST
	Endpoint    string            `yaml:"endpoint"`
	QueryParams map[string]string `yaml:"query_params,omitempty"`
	Headers     map[string]string `yaml:"headers,omitempty"`
	Body        string            `yaml:"body,omitempty"`
}

// Output defines how to format the response data.
type Output struct {
	Format string  `yaml:"format"` // table | json | yaml | raw
	Fields []Field `yaml:"fields,omitempty"`
}

// Field defines a single column/field in the output.
type Field struct {
	Name      string `yaml:"name"`
	Header    string `yaml:"header"`
	Path      string `yaml:"path"`      // gjson path
	Width     int    `yaml:"width,omitempty"`
	Transform string `yaml:"transform,omitempty"` // upper | lower | duration | truncate:N
}

// AIConfig provides hints for the AI chat mode.
type AIConfig struct {
	PromptHint      string   `yaml:"prompt_hint,omitempty"`
	NaturalLanguage []string `yaml:"natural_language,omitempty"`
	ResultHint      string   `yaml:"result_hint,omitempty"`
}

// Example provides usage examples for --help output.
type Example struct {
	Command     string `yaml:"command"`
	Description string `yaml:"description"`
}

// Result holds the processed output of an SPC execution.
type Result struct {
	Format string                   // output format
	Rows   []map[string]interface{} // extracted rows
	Raw    []byte                   // raw HTTP response body
}
```

- [ ] **Step 2: Write the types test — YAML round-trip**

File: `pkg/spc/types_test.go`
```go
package spc

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSPCDefinition_Unmarshal_FullSchema(t *testing.T) {
	input := `
name: pipeline-list
version: "1.0"
description: List pipeline runs
type: query
category: pipeline
tags: [ci, status]
parameters:
  - name: project_id
    type: string
    required: true
    description: Project ID
  - name: limit
    type: integer
    default: 10
    validation:
      min: 1
      max: 100
source:
  method: GET
  endpoint: /gateway/openlibing-cicd/project/pipeline/pipeline-run/detail
  query_params:
    projectId: "{{.project_id}}"
    pageSize: "{{.limit}}"
  headers:
    Content-Type: application/json
output:
  format: table
  fields:
    - name: id
      header: "ID"
      path: ".pipelineRunId"
      width: 36
    - name: status
      header: "Status"
      path: ".status"
      transform: upper
ai:
  prompt_hint: Use for pipeline queries
  natural_language:
    - "show pipelines"
  result_hint: Focus on failures
examples:
  - command: openlibing run pipeline-list --project-id 123
    description: List pipelines
`

	var spc SPCDefinition
	err := yaml.Unmarshal([]byte(input), &spc)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}

	if spc.Name != "pipeline-list" {
		t.Errorf("name = %q, want %q", spc.Name, "pipeline-list")
	}
	if spc.Version != "1.0" {
		t.Errorf("version = %q, want %q", spc.Version, "1.0")
	}
	if spc.Type != "query" {
		t.Errorf("type = %q, want %q", spc.Type, "query")
	}
	if spc.Category != "pipeline" {
		t.Errorf("category = %q, want %q", spc.Category, "pipeline")
	}
	if len(spc.Tags) != 2 || spc.Tags[0] != "ci" || spc.Tags[1] != "status" {
		t.Errorf("tags = %v, want [ci status]", spc.Tags)
	}

	// Parameters
	if len(spc.Parameters) != 2 {
		t.Fatalf("parameters len = %d, want 2", len(spc.Parameters))
	}
	p0 := spc.Parameters[0]
	if p0.Name != "project_id" || p0.Type != "string" || !p0.Required {
		t.Errorf("param[0] = %+v, want project_id/string/required", p0)
	}
	p1 := spc.Parameters[1]
	if p1.Name != "limit" || p1.Default != 10 {
		t.Errorf("param[1] = %+v, want limit/10", p1)
	}
	if p1.Validation == nil || *p1.Validation.Min != 1 || *p1.Validation.Max != 100 {
		t.Errorf("param[1].Validation = %+v, want min=1 max=100", p1.Validation)
	}

	// Source
	if spc.Source.Method != "GET" {
		t.Errorf("source.method = %q, want GET", spc.Source.Method)
	}
	if spc.Source.QueryParams["projectId"] != "{{.project_id}}" {
		t.Errorf("query_params[projectId] = %q", spc.Source.QueryParams["projectId"])
	}

	// Output
	if spc.Output.Format != "table" {
		t.Errorf("output.format = %q, want table", spc.Output.Format)
	}
	if len(spc.Output.Fields) != 2 {
		t.Errorf("output.fields len = %d, want 2", len(spc.Output.Fields))
	}
	if spc.Output.Fields[1].Transform != "upper" {
		t.Errorf("fields[1].transform = %q, want upper", spc.Output.Fields[1].Transform)
	}

	// AI
	if spc.AI == nil {
		t.Fatal("ai config is nil")
	}
	if len(spc.AI.NaturalLanguage) != 1 || spc.AI.NaturalLanguage[0] != "show pipelines" {
		t.Errorf("ai.natural_language = %v", spc.AI.NaturalLanguage)
	}

	// Examples
	if len(spc.Examples) != 1 {
		t.Errorf("examples len = %d, want 1", len(spc.Examples))
	}
}

func TestSPCDefinition_Unmarshal_MinimalSchema(t *testing.T) {
	input := `
name: minimal
version: "1.0"
description: Minimal SPC
type: query
category: test
source:
  method: GET
  endpoint: /api/test
output:
  format: json
`

	var spc SPCDefinition
	err := yaml.Unmarshal([]byte(input), &spc)
	if err != nil {
		t.Fatalf("unmarshal failed: %v", err)
	}
	if spc.Name != "minimal" {
		t.Errorf("name = %q, want minimal", spc.Name)
	}
	if len(spc.Parameters) != 0 {
		t.Errorf("parameters should be empty, got %d", len(spc.Parameters))
	}
	if spc.AI != nil {
		t.Errorf("ai should be nil when not specified")
	}
	if len(spc.Examples) != 0 {
		t.Errorf("examples should be empty, got %d", len(spc.Examples))
	}
}
```

- [ ] **Step 3: Run tests to verify**

```bash
go test ./pkg/spc/... -v
```

Expected: 2 tests PASS.

- [ ] **Step 4: Commit**

```bash
git add pkg/spc/
git commit -m "feat(spc): add core type definitions with YAML support"
```

---

### Task 3: Configuration Management

**Files:**
- Create: `internal/config/config.go`
- Create: `internal/config/auth.go`
- Create: `internal/config/config_test.go`

**Interfaces:**
- Consumes: nothing (file I/O only)
- Produces:
  - `type Config struct { Endpoint string; Defaults Defaults; Output OutputConfig }`
  - `type Defaults struct { ProjectID string; Limit int }`
  - `type OutputConfig struct { Format string; Color bool; Pager string }`
  - `type Auth struct { OpenLibing OpenLibingAuth; LLM LLMAuth }`
  - `type OpenLibingAuth struct { Token string; TokenType string }`
  - `type LLMAuth struct { Provider string; APIKey string }`
  - `func LoadConfig() (*Config, error)` — reads ~/.openlibing/config.yaml
  - `func LoadAuth() (*Auth, error)` — reads ~/.openlibing/auth.yaml
  - `func DefaultConfig() *Config` — returns hardcoded defaults

- [ ] **Step 1: Write config.go**

File: `internal/config/config.go`
```go
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config holds the user-level CLI configuration.
type Config struct {
	Endpoint string      `yaml:"endpoint"`
	Defaults Defaults    `yaml:"defaults"`
	Output   OutputConfig `yaml:"output"`
}

// Defaults holds default parameter values.
type Defaults struct {
	ProjectID string `yaml:"project_id"`
	Limit     int    `yaml:"limit"`
}

// OutputConfig holds output preferences.
type OutputConfig struct {
	Format string `yaml:"format"` // table | json | yaml
	Color  bool   `yaml:"color"`
	Pager  string `yaml:"pager"` // auto | never | always
}

// DefaultConfig returns the hardcoded default configuration.
func DefaultConfig() *Config {
	return &Config{
		Endpoint: "https://www.openlibing.com",
		Defaults: Defaults{
			Limit: 10,
		},
		Output: OutputConfig{
			Format: "table",
			Color:  true,
			Pager:  "auto",
		},
	}
}

// LoadConfig reads configuration from ~/.openlibing/config.yaml.
// If the file does not exist, returns DefaultConfig().
func LoadConfig() (*Config, error) {
	cfg := DefaultConfig()

	home, err := os.UserHomeDir()
	if err != nil {
		return cfg, nil // can't determine home, use defaults
	}

	path := filepath.Join(home, ".openlibing", "config.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return cfg, nil // file doesn't exist, use defaults
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}

// ConfigDir returns the user's openlibing config directory, creating it if needed.
func ConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dir := filepath.Join(home, ".openlibing")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", err
	}
	return dir, nil
}

// SPCUserDir returns the user's custom SPC directory.
func SPCUserDir() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	spcDir := filepath.Join(dir, "spc")
	if err := os.MkdirAll(spcDir, 0755); err != nil {
		return "", err
	}
	return spcDir, nil
}
```

- [ ] **Step 2: Write auth.go**

File: `internal/config/auth.go`
```go
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Auth holds authentication credentials.
type Auth struct {
	OpenLibing OpenLibingAuth `yaml:"openlibing"`
	LLM        LLMAuth        `yaml:"llm,omitempty"`
}

// OpenLibingAuth holds OpenLibing platform credentials.
type OpenLibingAuth struct {
	Token     string `yaml:"token"`
	TokenType string `yaml:"token_type"` // Bearer
}

// LLMAuth holds LLM API credentials.
type LLMAuth struct {
	Provider string `yaml:"provider"`
	APIKey   string `yaml:"api_key"`
}

// LoadAuth reads authentication from ~/.openlibing/auth.yaml.
// Returns empty Auth if file doesn't exist.
func LoadAuth() (*Auth, error) {
	auth := &Auth{
		OpenLibing: OpenLibingAuth{
			TokenType: "Bearer",
		},
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return auth, nil
	}

	path := filepath.Join(home, ".openlibing", "auth.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return auth, nil
		}
		return nil, err
	}

	if err := yaml.Unmarshal(data, auth); err != nil {
		return nil, err
	}

	return auth, nil
}
```

- [ ] **Step 3: Write config_test.go**

File: `internal/config/config_test.go`
```go
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
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/config/... -v
```

Expected: 6 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/config/
go mod tidy
git add go.mod go.sum
git commit -m "feat(config): add configuration and auth management"
```

---

### Task 4: API Client

**Files:**
- Create: `internal/api/client.go`
- Create: `internal/api/cicd.go`
- Create: `internal/api/client_test.go`

**Interfaces:**
- Consumes: `*config.Auth`, `*config.Config`
- Produces:
  - `type Client struct { ... }` — wraps http.Client with auth, retry, base URL
  - `func NewClient(cfg *config.Config, auth *config.Auth) *Client`
  - `func (c *Client) Do(method, endpoint string, queryParams map[string]string, headers map[string]string, body io.Reader) (*http.Response, error)`
  - `func (c *Client) Get(endpoint string, queryParams map[string]string) (*http.Response, error)`
  - `func (c *Client) Post(endpoint string, body io.Reader) (*http.Response, error)`

- [ ] **Step 1: Write client.go**

File: `internal/api/client.go`
```go
package api

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/openlibing/openlibing-cli/internal/config"
)

// Client wraps http.Client with OpenLibing-specific behavior:
// auth injection, base URL resolution, and retry logic.
type Client struct {
	httpClient *http.Client
	baseURL    string
	auth       *config.OpenLibingAuth
	maxRetries int
}

// NewClient creates a new API client.
func NewClient(cfg *config.Config, auth *config.Auth) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:    strings.TrimRight(cfg.Endpoint, "/"),
		auth:       &auth.OpenLibing,
		maxRetries: 3,
	}
}

// Do executes an HTTP request with auth and retry logic.
func (c *Client) Do(method, endpoint string, queryParams map[string]string, headers map[string]string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + "/" + strings.TrimLeft(endpoint, "/")

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	// Set query parameters
	if len(queryParams) > 0 {
		q := req.URL.Query()
		for k, v := range queryParams {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Set headers from SPC
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// Inject auth
	if c.auth.Token != "" {
		req.Header.Set("Authorization", c.auth.TokenType+" "+c.auth.Token)
	}

	// Retry loop
	var lastErr error
	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s
			time.Sleep(time.Duration(1<<(attempt-1)) * time.Second)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue // network error, retry
		}

		// Only retry on 5xx server errors
		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", c.maxRetries, lastErr)
}

// Get is a convenience method for GET requests with query params.
func (c *Client) Get(endpoint string, queryParams map[string]string) (*http.Response, error) {
	return c.Do("GET", endpoint, queryParams, nil, nil)
}

// Post is a convenience method for POST requests.
func (c *Client) Post(endpoint string, body io.Reader) (*http.Response, error) {
	return c.Do("POST", endpoint, nil, nil, body)
}
```

- [ ] **Step 2: Write cicd.go (MVP: pipeline detail endpoint)**

File: `internal/api/cicd.go`
```go
package api

import (
	"fmt"
	"io"
	"net/http"
)

// GetPipelineDetail fetches pipeline run detail from openlibing-cicd.
func (c *Client) GetPipelineDetail(projectID string, limit int) ([]byte, error) {
	resp, err := c.Do("GET", "gateway/openlibing-cicd/project/pipeline/pipeline-run/detail", map[string]string{
		"projectId": projectID,
		"pageSize":  fmt.Sprintf("%d", limit),
	}, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("get pipeline detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	return io.ReadAll(resp.Body)
}
```

- [ ] **Step 3: Write client_test.go (mock HTTP server)**

File: `internal/api/client_test.go`
```go
package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/openlibing/openlibing-cli/internal/config"
)

func TestClient_Do_InjectsAuth(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer test-token" {
			t.Errorf("Authorization = %q, want %q", auth, "Bearer test-token")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.Endpoint = server.URL
	auth := &config.Auth{
		OpenLibing: config.OpenLibingAuth{
			Token:     "test-token",
			TokenType: "Bearer",
		},
	}

	client := NewClient(cfg, auth)
	resp, err := client.Get("test/endpoint", nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestClient_Do_RetriesOn5xx(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.Endpoint = server.URL
	auth := &config.Auth{}

	client := NewClient(cfg, auth)
	resp, err := client.Get("test/endpoint", nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if attempts != 3 {
		t.Errorf("attempts = %d, want 3", attempts)
	}
	if resp.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want 200", resp.StatusCode)
	}
}

func TestClient_Do_NoRetryOn4xx(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.Endpoint = server.URL
	auth := &config.Auth{}

	client := NewClient(cfg, auth)
	_, err := client.Get("test/endpoint", nil)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if attempts != 1 {
		t.Errorf("attempts = %d, want 1 (no retry on 4xx)", attempts)
	}
}

func TestClient_GetPipelineDetail_BuildsURL(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		projectID := r.URL.Query().Get("projectId")
		pageSize := r.URL.Query().Get("pageSize")
		if projectID != "123" || pageSize != "10" {
			t.Errorf("query params: projectId=%q pageSize=%q", projectID, pageSize)
		}

		resp := map[string]interface{}{
			"data": []map[string]interface{}{
				{"pipelineRunId": "run-1", "status": "SUCCESS"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	cfg := config.DefaultConfig()
	cfg.Endpoint = server.URL
	auth := &config.Auth{}

	client := NewClient(cfg, auth)
	body, err := client.GetPipelineDetail("123", 10)
	if err != nil {
		t.Fatalf("GetPipelineDetail failed: %v", err)
	}
	if len(body) == 0 {
		t.Error("response body is empty")
	}
}
```

- [ ] **Step 4: Run tests**

```bash
go test ./internal/api/... -v
```

Expected: 4 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/api/
git commit -m "feat(api): add HTTP client with auth injection and retry logic"
```

---

### Task 5: SPC Parser and Validator

**Files:**
- Create: `internal/engine/parser.go`
- Create: `internal/engine/validator.go`
- Create: `internal/engine/parser_test.go`
- Create: `internal/engine/validator_test.go`
- Create: `internal/engine/testdata/valid.spc.yaml`
- Create: `internal/engine/testdata/invalid_syntax.spc.yaml`

**Interfaces:**
- Consumes: `*spc.SPCDefinition`, `spc.Parameter`
- Produces:
  - `func ParseSPC(r io.Reader) (*spc.SPCDefinition, error)` — YAML → SPCDefinition
  - `func ParseSPCFile(path string) (*spc.SPCDefinition, error)` — file → SPCDefinition
  - `func Validate(spc *spc.SPCDefinition, params map[string]interface{}) error`

- [ ] **Step 1: Write parser.go**

File: `internal/engine/parser.go`
```go
package engine

import (
	"fmt"
	"io"
	"os"

	"github.com/openlibing/openlibing-cli/pkg/spc"
	"gopkg.in/yaml.v3"
)

// ParseSPC reads a YAML stream and returns a parsed SPCDefinition.
func ParseSPC(r io.Reader) (*spc.SPCDefinition, error) {
	var def spc.SPCDefinition
	decoder := yaml.NewDecoder(r)
	decoder.KnownFields(true)
	if err := decoder.Decode(&def); err != nil {
		return nil, fmt.Errorf("parse SPC: %w", err)
	}
	return &def, nil
}

// ParseSPCFile reads a .spc.yaml file and returns a parsed SPCDefinition.
func ParseSPCFile(path string) (*spc.SPCDefinition, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open SPC file %s: %w", path, err)
	}
	defer f.Close()

	def, err := ParseSPC(f)
	if err != nil {
		return nil, fmt.Errorf("file %s: %w", path, err)
	}
	def.FilePath = path
	return def, nil
}
```

- [ ] **Step 2: Write validator.go**

File: `internal/engine/validator.go`
```go
package engine

import (
	"fmt"
	"strconv"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// Validate checks that the provided parameters satisfy the SPC definition.
// It verifies required fields, type correctness, enum membership, and range constraints.
func Validate(def *spc.SPCDefinition, params map[string]interface{}) error {
	for _, p := range def.Parameters {
		val, exists := params[p.Name]

		// Required check
		if p.Required && (!exists || val == nil || val == "") {
			return &ValidationError{
				Parameter: p.Name,
				Message:   fmt.Sprintf("parameter '%s' is required", p.Name),
			}
		}

		// Skip further checks if no value and not required
		if !exists || val == nil || val == "" {
			// Apply default if available
			if p.Default != nil && exists {
				// User explicitly passed empty → use default
				// This covers the case where value is provided but empty
			}
			continue
		}

		// Type check
		switch p.Type {
		case "integer":
			if err := checkInteger(val); err != nil {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s': %v", p.Name, err),
				}
			}
		case "boolean":
			if err := checkBoolean(val); err != nil {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s': %v", p.Name, err),
				}
			}
		// string type needs no explicit check — everything can be a string
		}

		// Enum check
		if len(p.Enum) > 0 {
			strVal := fmt.Sprintf("%v", val)
			found := false
			for _, e := range p.Enum {
				if e == strVal {
					found = true
					break
				}
			}
			if !found {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s' must be one of %v, got '%s'", p.Name, p.Enum, strVal),
				}
			}
		}

		// Range validation (integer only)
		if p.Validation != nil && p.Type == "integer" {
			n, _ := toInt(val)
			if p.Validation.Min != nil && n < *p.Validation.Min {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s' must be >= %d, got %d", p.Name, *p.Validation.Min, n),
				}
			}
			if p.Validation.Max != nil && n > *p.Validation.Max {
				return &ValidationError{
					Parameter: p.Name,
					Message:   fmt.Sprintf("parameter '%s' must be <= %d, got %d", p.Name, *p.Validation.Max, n),
				}
			}
		}
	}

	return nil
}

// ValidationError represents a parameter validation failure.
type ValidationError struct {
	Parameter string
	Message   string
}

func (e *ValidationError) Error() string {
	return e.Message
}

func checkInteger(val interface{}) error {
	switch v := val.(type) {
	case int, int64, float64:
		return nil
	case string:
		_, err := strconv.Atoi(v)
		return err
	default:
		return fmt.Errorf("expected integer, got %T", val)
	}
}

func checkBoolean(val interface{}) error {
	switch v := val.(type) {
	case bool:
		return nil
	case string:
		if v != "true" && v != "false" {
			return fmt.Errorf("expected boolean, got %q", v)
		}
		return nil
	default:
		return fmt.Errorf("expected boolean, got %T", val)
	}
}

func toInt(val interface{}) (int, error) {
	switch v := val.(type) {
	case int:
		return v, nil
	case int64:
		return int(v), nil
	case float64:
		return int(v), nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, fmt.Errorf("cannot convert %T to int", val)
	}
}
```

- [ ] **Step 3: Create test fixtures**

File: `internal/engine/testdata/valid.spc.yaml`
```yaml
name: pipeline-list
version: "1.0"
description: List pipeline runs
type: query
category: pipeline
tags: [ci]
parameters:
  - name: project_id
    type: string
    required: true
    description: Project ID
  - name: limit
    type: integer
    default: 10
    validation:
      min: 1
      max: 100
  - name: status
    type: string
    enum: [running, success, failed]
source:
  method: GET
  endpoint: /gateway/test
output:
  format: table
  fields:
    - name: id
      header: "ID"
      path: ".id"
```

File: `internal/engine/testdata/invalid_syntax.spc.yaml`
```yaml
name: [ this is broken yaml
  - bad indent
```

- [ ] **Step 4: Write parser_test.go**

File: `internal/engine/parser_test.go`
```go
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
```

- [ ] **Step 5: Write validator_test.go**

File: `internal/engine/validator_test.go`
```go
package engine

import (
	"os"
	"strings"
	"testing"
)

func loadTestSPC(t *testing.T) *SPCDefinition {
	t.Helper()
	data, err := os.ReadFile("testdata/valid.spc.yaml")
	if err != nil {
		t.Fatalf("read fixture: %v", err)
	}
	// Note: this imports from pkg/spc — adjust import in actual file
	return nil // placeholder — actual implementation imports pkg/spc
}
```

Wait — the validator_test.go needs to import pkg/spc. Let me rewrite this properly.

File: `internal/engine/validator_test.go`
```go
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
```

- [ ] **Step 6: Run parser tests**

```bash
go test ./internal/engine/ -run TestParse -v
```

Expected: 4 tests PASS.

- [ ] **Step 7: Run validator tests**

```bash
go test ./internal/engine/ -run TestValidate -v
```

Expected: 7 tests PASS.

- [ ] **Step 8: Commit**

```bash
git add internal/engine/parser.go internal/engine/validator.go internal/engine/testdata/
git add internal/engine/parser_test.go internal/engine/validator_test.go
git commit -m "feat(engine): add SPC parser and parameter validator"
```

---

### Task 6: Template Resolver and HTTP Executor

**Files:**
- Create: `internal/engine/resolver.go`
- Create: `internal/engine/executor.go`
- Create: `internal/engine/resolver_test.go`
- Create: `internal/engine/executor_test.go`

**Interfaces:**
- Consumes: `spc.Source`, `api.Client`
- Produces:
  - `func Resolve(source *spc.Source, params map[string]interface{}) (*ResolvedRequest, error)`
  - `type ResolvedRequest struct { Method, URL string; Headers map[string]string; Body string }`
  - `func ExecuteRequest(client *api.Client, req *ResolvedRequest) (*http.Response, error)`

- [ ] **Step 1: Write resolver.go**

File: `internal/engine/resolver.go`
```go
package engine

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// ResolvedRequest is a fully-materialized HTTP request ready to execute.
type ResolvedRequest struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    string
}

// Resolve renders all templates in a Source definition against the given parameters.
func Resolve(source *spc.Source, params map[string]interface{}, baseURL string) (*ResolvedRequest, error) {
	// Build template context: parameter values keyed by name
	tmplCtx := make(map[string]interface{})
	for k, v := range params {
		tmplCtx[k] = v
	}

	// Resolve endpoint
	endpoint, err := renderTemplate("endpoint", source.Endpoint, tmplCtx)
	if err != nil {
		return nil, err
	}

	rr := &ResolvedRequest{
		Method:  source.Method,
		URL:     strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(endpoint, "/"),
		Headers: make(map[string]string),
		Body:    source.Body,
	}

	// Resolve query params
	if len(source.QueryParams) > 0 || source.Method == "GET" {
		queryParts := []string{}
		for k, v := range source.QueryParams {
			resolvedV, err := renderTemplate("query_param:"+k, v, tmplCtx)
			if err != nil {
				return nil, err
			}
			queryParts = append(queryParts, k+"="+resolvedV)
		}
		if len(queryParts) > 0 {
			rr.URL += "?" + strings.Join(queryParts, "&")
		}
	}

	// Resolve headers
	for k, v := range source.Headers {
		resolvedV, err := renderTemplate("header:"+k, v, tmplCtx)
		if err != nil {
			return nil, err
		}
		rr.Headers[k] = resolvedV
	}

	// Resolve body template (for POST)
	if source.Body != "" {
		resolvedBody, err := renderTemplate("body", source.Body, tmplCtx)
		if err != nil {
			return nil, err
		}
		rr.Body = resolvedBody
	}

	return rr, nil
}

func renderTemplate(name, tmpl string, ctx map[string]interface{}) (string, error) {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("template parse [%s]: %w", name, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("template execute [%s]: %w", name, err)
	}
	return buf.String(), nil
}
```

- [ ] **Step 2: Write executor.go**

File: `internal/engine/executor.go`
```go
package engine

import (
	"io"
	"net/http"
	"strings"

	"github.com/openlibing/openlibing-cli/internal/api"
)

// ExecuteRequest sends a resolved request through the API client and returns the response.
func ExecuteRequest(client *api.Client, req *ResolvedRequest) (*http.Response, []byte, error) {
	// Extract endpoint from URL
	endpoint := strings.TrimPrefix(req.URL, client.BaseURL())
	endpoint = strings.TrimLeft(endpoint, "/")

	// Remove query string from endpoint for API client (it handles query params separately)
	if idx := strings.Index(endpoint, "?"); idx >= 0 {
		endpoint = endpoint[:idx]
	}

	// Map headers
	headers := req.Headers

	// Execute through API client
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	resp, err := client.Do(req.Method, endpoint, nil, headers, bodyReader)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp, nil, err
	}

	return resp, respBody, nil
}
```

- [ ] **Step 8: Commit**

```bash
git add internal/engine/parser.go internal/engine/validator.go internal/engine/testdata/
git add internal/engine/parser_test.go internal/engine/validator_test.go
git commit -m "feat(engine): add SPC parser and parameter validator"
```

---

### Task 6: Template Resolver and HTTP Executor

**Files:**
- Create: `internal/engine/resolver.go`
- Create: `internal/engine/resolver_test.go`

**Interfaces:**
- Consumes: `spc.Source`, `spc.Parameter`
- Produces:
  - `type ResolvedRequest struct { Method, Endpoint string; QueryParams, Headers map[string]string; Body string }`
  - `func Resolve(source *spc.Source, params map[string]interface{}) (*ResolvedRequest, error)` — renders all `{{.var}}` templates in source fields

**Notes:** The resolver does NOT build a full URL — it only resolves template variables. The `api.Client` handles base URL assembly in its `Do()` method. This keeps the resolver pure and avoids circular dependencies.

- [ ] **Step 1: Write resolver.go**

File: `internal/engine/resolver.go`
```go
package engine

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// ResolvedRequest holds all template-resolved fields from a Source definition.
// Templates like {{.project_id}} have been rendered against user-supplied params.
type ResolvedRequest struct {
	Method      string
	Endpoint    string            // e.g. "gateway/openlibing-cicd/..."
	QueryParams map[string]string // resolved key=value pairs
	Headers     map[string]string
	Body        string
}

// Resolve renders all Go templates in a Source definition against the provided params.
func Resolve(source *spc.Source, params map[string]interface{}) (*ResolvedRequest, error) {
	ctx := make(map[string]interface{})
	for k, v := range params {
		ctx[k] = fmt.Sprintf("%v", v)
	}

	rr := &ResolvedRequest{
		Method:      source.Method,
		QueryParams: make(map[string]string),
		Headers:     make(map[string]string),
	}

	// Resolve endpoint
	var err error
	rr.Endpoint, err = render("endpoint", source.Endpoint, ctx)
	if err != nil {
		return nil, err
	}

	// Resolve query params
	for k, v := range source.QueryParams {
		rr.QueryParams[k], err = render("query:"+k, v, ctx)
		if err != nil {
			return nil, err
		}
	}

	// Resolve headers
	for k, v := range source.Headers {
		rr.Headers[k], err = render("header:"+k, v, ctx)
		if err != nil {
			return nil, err
		}
	}

	// Resolve body
	if source.Body != "" {
		rr.Body, err = render("body", source.Body, ctx)
		if err != nil {
			return nil, err
		}
	}

	return rr, nil
}

func render(name, tmpl string, ctx map[string]interface{}) (string, error) {
	t, err := template.New(name).Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("template parse [%s]: %w", name, err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("template execute [%s]: %w", name, err)
	}
	return buf.String(), nil
}
```

- [ ] **Step 2: Write resolver_test.go**

File: `internal/engine/resolver_test.go`
```go
package engine

import (
	"testing"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

func TestResolve_RendersEndpoint(t *testing.T) {
	source := &spc.Source{
		Method:   "GET",
		Endpoint: "gateway/test/{{.project_id}}",
	}
	params := map[string]interface{}{
		"project_id": "my-project-123",
	}

	rr, err := Resolve(source, params)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if rr.Endpoint != "gateway/test/my-project-123" {
		t.Errorf("endpoint = %q", rr.Endpoint)
	}
	if rr.Method != "GET" {
		t.Errorf("method = %q", rr.Method)
	}
}

func TestResolve_RendersQueryParams(t *testing.T) {
	source := &spc.Source{
		Method:   "GET",
		Endpoint: "gateway/test",
		QueryParams: map[string]string{
			"projectId": "{{.project_id}}",
			"pageSize":  "{{.limit}}",
		},
	}
	params := map[string]interface{}{
		"project_id": "456",
		"limit":      10,
	}

	rr, err := Resolve(source, params)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if rr.QueryParams["projectId"] != "456" {
		t.Errorf("queryParams[projectId] = %q", rr.QueryParams["projectId"])
	}
	if rr.QueryParams["pageSize"] != "10" {
		t.Errorf("queryParams[pageSize] = %q", rr.QueryParams["pageSize"])
	}
}

func TestResolve_RendersHeaders(t *testing.T) {
	source := &spc.Source{
		Method:   "POST",
		Endpoint: "gateway/test",
		Headers: map[string]string{
			"Content-Type": "application/json",
			"X-Project":    "{{.project_id}}",
		},
	}
	params := map[string]interface{}{
		"project_id": "789",
	}

	rr, err := Resolve(source, params)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if rr.Headers["X-Project"] != "789" {
		t.Errorf("headers[X-Project] = %q", rr.Headers["X-Project"])
	}
	// Static header should pass through unchanged
	if rr.Headers["Content-Type"] != "application/json" {
		t.Errorf("headers[Content-Type] = %q", rr.Headers["Content-Type"])
	}
}

func TestResolve_RendersBody(t *testing.T) {
	source := &spc.Source{
		Method:   "POST",
		Endpoint: "gateway/test",
		Body:     `{"projectId":"{{.project_id}}"}`,
	}
	params := map[string]interface{}{
		"project_id": "body-test",
	}

	rr, err := Resolve(source, params)
	if err != nil {
		t.Fatalf("Resolve failed: %v", err)
	}
	if rr.Body != `{"projectId":"body-test"}` {
		t.Errorf("body = %q", rr.Body)
	}
}

func TestResolve_BadTemplate_ReturnsError(t *testing.T) {
	source := &spc.Source{
		Method:   "GET",
		Endpoint: "gateway/{{.nonexistent}}",
	}
	params := map[string]interface{}{}

	_, err := Resolve(source, params)
	if err == nil {
		t.Fatal("expected error for missing template variable, got nil")
	}
}
```

- [ ] **Step 3: Run resolver tests**

```bash
go test ./internal/engine/ -run TestResolve -v
```

Expected: 5 tests PASS.

- [ ] **Step 4: Commit**

```bash
git add internal/engine/resolver.go internal/engine/resolver_test.go
git commit -m "feat(engine): add template resolver for SPC source fields"
```<｜end▁of▁thinking｜>Let me continue the plan:

<｜｜DSML｜｜tool_calls>
<｜｜DSML｜｜invoke name="Read">
<｜｜DSML｜｜parameter name="file_path" string="true">/home/wangsike/workspace/openlibing/openlibing-cli/docs/superpowers/plans/2026-06-22-openlibing-cli-mvp.md
---

### Task 7: Output Formatter

**Files:**
- Create: `internal/cli/output.go`
- Create: `internal/cli/output_test.go`

**Interfaces:**
- Consumes: `spc.Result`, `spc.Field`
- Produces:
  - `func FormatTable(rows []map[string]interface{}, fields []spc.Field) string`
  - `func FormatJSON(rows []map[string]interface{}) (string, error)`
  - `func FormatRaw(raw []byte) string`
  - `func FormatResult(result *spc.Result) (string, error)` — dispatches by result.Format
  - `func TransformDuration(val interface{}) string` — ms → "3m 42s"

- [ ] **Step 1: Write output.go**

File: `internal/cli/output.go`
```go
package cli

import (
	"encoding/json"
	"fmt"
	"strings"
	"text/tabwriter"

	"bytes"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// FormatResult formats a Result according to its Format field.
func FormatResult(result *spc.Result) (string, error) {
	switch result.Format {
	case "json":
		return FormatJSON(result.Rows)
	case "table":
		return FormatTable(result.Rows, nil), nil
	case "raw":
		return FormatRaw(result.Raw), nil
	default:
		return FormatTable(result.Rows, nil), nil
	}
}

// FormatTable renders rows as an aligned terminal table.
func FormatTable(rows []map[string]interface{}, fields []spc.Field) string {
	if len(rows) == 0 {
		return "(no results)"
	}

	// If no fields specified, infer from first row keys
	if len(fields) == 0 {
		first := rows[0]
		for k := range first {
			fields = append(fields, spc.Field{Name: k, Header: k, Path: k})
		}
	}

	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 2, ' ', 0)

	// Headers
	headers := make([]string, len(fields))
	for i, f := range fields {
		headers[i] = f.Header
	}
	fmt.Fprintln(w, strings.Join(headers, "\t"))

	// Separator
	seps := make([]string, len(fields))
	for i := range fields {
		seps[i] = strings.Repeat("─", len(headers[i])+4)
	}
	fmt.Fprintln(w, strings.Join(seps, "\t"))

	// Rows
	for _, row := range rows {
		cols := make([]string, len(fields))
		for i, f := range fields {
			val := row[f.Name]
			col := formatValue(val, f.Transform)
			cols[i] = col
		}
		fmt.Fprintln(w, strings.Join(cols, "\t"))
	}

	w.Flush()
	return buf.String()
}

// FormatJSON renders rows as indented JSON.
func FormatJSON(rows []map[string]interface{}) (string, error) {
	if rows == nil {
		rows = []map[string]interface{}{}
	}
	data, err := json.MarshalIndent(rows, "", "  ")
	if err != nil {
		return "", fmt.Errorf("json marshal: %w", err)
	}
	return string(data), nil
}

// FormatRaw returns raw bytes as a string.
func FormatRaw(raw []byte) string {
	return string(raw)
}

func formatValue(val interface{}, transform string) string {
	if val == nil {
		return "-"
	}

	str := fmt.Sprintf("%v", val)

	switch transform {
	case "upper":
		return strings.ToUpper(str)
	case "lower":
		return strings.ToLower(str)
	case "duration":
		return TransformDuration(val)
	case "":
		return str
	default:
		if strings.HasPrefix(transform, "truncate:") {
			var n int
			fmt.Sscanf(transform, "truncate:%d", &n)
			if len(str) > n {
				return str[:n] + "…"
			}
		}
		return str
	}
}

// TransformDuration converts milliseconds to a human-readable duration string.
func TransformDuration(val interface{}) string {
	var ms int64
	switch v := val.(type) {
	case float64:
		ms = int64(v)
	case int64:
		ms = v
	case int:
		ms = int64(v)
	default:
		return fmt.Sprintf("%v", val)
	}

	if ms < 1000 {
		return fmt.Sprintf("%dms", ms)
	}

	seconds := ms / 1000
	if seconds < 60 {
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := seconds / 60
	remainingSeconds := seconds % 60
	if minutes < 60 {
		return fmt.Sprintf("%dm %ds", minutes, remainingSeconds)
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60
	return fmt.Sprintf("%dh %dm", hours, remainingMinutes)
}
```

- [ ] **Step 2: Write output_test.go**

File: `internal/cli/output_test.go`
```go
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
```

- [ ] **Step 3: Run tests**

```bash
go test ./internal/cli/ -run "TestFormat|TestTransform" -v
```

Expected: 8 tests PASS.

- [ ] **Step 4: Commit**

```bash
git add internal/cli/output.go internal/cli/output_test.go
git commit -m "feat(cli): add output formatter (table/json/raw) with transforms"
```

---

### Task 8: Engine Core (Orchestrator)

**Files:**
- Create: `internal/engine/engine.go`
- Create: `internal/engine/engine_test.go`
- Create: `internal/engine/testdata/pipeline_list_response.json`

**Interfaces:**
- Consumes: `api.Client`, `*spc.SPCDefinition`, all engine sub-packages (parse, validate, resolve, extract)
- Produces:
  - `type Engine struct { registry interface{ Get(string) (*spc.SPCDefinition, error) }; client *api.Client }`
  - `func NewEngine(registry interface{ Get(name string) (*spc.SPCDefinition, error) }, client *api.Client) *Engine`
  - `func (e *Engine) Execute(name string, params map[string]interface{}) (*spc.Result, error)` — 7-step pipeline

- [ ] **Step 1: Write engine.go (orchestrator)**

File: `internal/engine/engine.go`
```go
package engine

import (
	"encoding/json"
	"fmt"

	"github.com/openlibing/openlibing-cli/internal/api"
	"github.com/openlibing/openlibing-cli/pkg/spc"
	"github.com/tidwall/gjson"
)

// SPCLookup is the minimal interface Engine needs from Registry.
type SPCLookup interface {
	Get(name string) (*spc.SPCDefinition, error)
}

// Engine orchestrates the 7-step SPC execution pipeline.
type Engine struct {
	registry SPCLookup
	client   *api.Client
}

// NewEngine creates a new Engine.
func NewEngine(registry SPCLookup, client *api.Client) *Engine {
	return &Engine{
		registry: registry,
		client:   client,
	}
}

// Execute runs the full 7-step pipeline: Parse → Validate → Resolve → Call → Extract → Transform → Output.
func (e *Engine) Execute(name string, params map[string]interface{}) (*spc.Result, error) {
	// Step 1: Parse — lookup SPC from registry
	def, err := e.registry.Get(name)
	if err != nil {
		return nil, fmt.Errorf("parse: %w", err)
	}

	// Step 2: Validate — check parameters
	if err := Validate(def, params); err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	// Apply defaults for unset parameters
	params = applyDefaults(def, params)

	// Step 3: Resolve — render templates
	resolved, err := Resolve(&def.Source, params)
	if err != nil {
		return nil, fmt.Errorf("resolve: %w", err)
	}

	// Step 4: Call — execute HTTP request
	respBody, err := e.call(resolved)
	if err != nil {
		return nil, fmt.Errorf("call: %w", err)
	}

	// Step 5: Extract — parse JSON and apply gjson paths
	rows, err := extract(def, respBody)
	if err != nil {
		// Degrade: return raw JSON
		return &spc.Result{
			Format: "raw",
			Raw:    respBody,
		}, nil
	}

	// Steps 6 & 7: Transform (applied during format stage in CLI) and Output
	return &spc.Result{
		Format: def.Output.Format,
		Rows:   rows,
		Raw:    respBody,
	}, nil
}

func (e *Engine) call(resolved *ResolvedRequest) ([]byte, error) {
	resp, err := e.client.Do(
		resolved.Method,
		resolved.Endpoint,
		resolved.QueryParams,
		resolved.Headers,
		nil, // body handled via query params for now
	)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read body
	buf := make([]byte, 0)
	b := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(b)
		buf = append(buf, b[:n]...)
		if err != nil {
			break
		}
	}
	return buf, nil
}

func extract(def *spc.SPCDefinition, rawJSON []byte) ([]map[string]interface{}, error) {
	jsonStr := string(rawJSON)

	// Determine if response wraps data in a "data" array
	dataResult := gjson.Get(jsonStr, "data")
	var arr []gjson.Result
	if dataResult.IsArray() {
		arr = dataResult.Array()
	} else {
		// Try root-level array
		rootResult := gjson.Parse(jsonStr)
		if rootResult.IsArray() {
			arr = rootResult.Array()
		} else {
			// Single object — wrap in array
			arr = []gjson.Result{rootResult}
		}
	}

	var rows []map[string]interface{}
	for _, item := range arr {
		row := make(map[string]interface{})
		for _, field := range def.Output.Fields {
			val := item.Get(field.Path)
			if val.Exists() {
				row[field.Name] = val.Value()
			}
		}
		// If no fields defined, include all keys
		if len(def.Output.Fields) == 0 {
			item.ForEach(func(key, value gjson.Result) bool {
				row[key.String()] = value.Value()
				return true
			})
		}
		rows = append(rows, row)
	}

	return rows, nil
}

func applyDefaults(def *spc.SPCDefinition, params map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for k, v := range params {
		result[k] = v
	}
	for _, p := range def.Parameters {
		if _, exists := result[p.Name]; !exists && p.Default != nil {
			result[p.Name] = p.Default
		}
	}
	return result
}

// Ensure json import is used
var _ = json.Valid
```

- [ ] **Step 2: Create test fixture**

File: `internal/engine/testdata/pipeline_list_response.json`
```json
{
  "data": [
    {
      "pipelineRunId": "run-abc-123",
      "status": "SUCCESS",
      "ref": "main",
      "durationMillis": 125000,
      "createTime": "2026-06-22T10:30:00Z"
    },
    {
      "pipelineRunId": "run-def-456",
      "status": "FAILED",
      "ref": "feat/new-thing",
      "durationMillis": 380000,
      "createTime": "2026-06-22T09:15:00Z"
    }
  ]
}
```

- [ ] **Step 3: Write engine_test.go**

File: `internal/engine/engine_test.go`
```go
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
```

- [ ] **Step 4: Run engine tests**

```bash
go test ./internal/engine/ -run TestEngine -v
```

Expected: 3 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add internal/engine/engine.go internal/engine/engine_test.go internal/engine/testdata/pipeline_list_response.json
git commit -m "feat(engine): add 7-step execution orchestrator"
```

---

### Task 9: Registry (SPC Discovery and Indexing)

**Files:**
- Create: `internal/registry/registry.go`
- Create: `internal/registry/loader.go`
- Create: `internal/registry/resolver.go`
- Create: `internal/registry/registry_test.go`

**Interfaces:**
- Consumes: `spc.SPCDefinition`, `engine.ParseSPC` / `engine.ParseSPCFile`
- Produces:
  - `type RegistryImpl struct { ... }`
  - `func NewRegistry(embedFS embed.FS) *RegistryImpl`
  - `func (r *RegistryImpl) LoadAll() error` — 3-layer load with override
  - `func (r *RegistryImpl) Get(name string) (*spc.SPCDefinition, error)`
  - `func (r *RegistryImpl) ListAll() []*spc.SPCDefinition`
  - `func (r *RegistryImpl) ListByCategory(cat string) []*spc.SPCDefinition`
  - `func (r *RegistryImpl) Search(query string) []*spc.SPCDefinition`
  - `func ResolveSearch(search string, defs []*spc.SPCDefinition) []*spc.SPCDefinition` — keyword matcher

- [ ] **Step 1: Write loader.go**

File: `internal/registry/loader.go`
```go
package registry

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/openlibing/openlibing-cli/internal/config"
	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// LoadFromEmbed reads all .spc.yaml files from an embedded filesystem.
func LoadFromEmbed(embedFS fs.FS, dir string) ([]*spc.SPCDefinition, error) {
	entries, err := fs.ReadDir(embedFS, dir)
	if err != nil {
		return nil, err
	}

	var defs []*spc.SPCDefinition
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".spc.yaml") {
			continue
		}

		f, err := embedFS.Open(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		defer f.Close()

		def, err := engine.ParseSPC(f)
		if err != nil {
			continue // skip malformed SPCs
		}
		def.Origin = "builtin"
		defs = append(defs, def)
	}

	return defs, nil
}

// LoadFromDir scans a directory for .spc.yaml files and parses them.
func LoadFromDir(dir, origin string) ([]*spc.SPCDefinition, error) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil // directory doesn't exist, no SPCs to load
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var defs []*spc.SPCDefinition
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".spc.yaml") {
			continue
		}

		path := filepath.Join(dir, entry.Name())
		def, err := engine.ParseSPCFile(path)
		if err != nil {
			continue // skip malformed SPCs
		}
		def.Origin = origin
		defs = append(defs, def)
	}

	return defs, nil
}

// LoadUserSPCs loads SPCs from ~/.openlibing/spc/
func LoadUserSPCs() ([]*spc.SPCDefinition, error) {
	dir, err := config.SPCUserDir()
	if err != nil {
		return nil, nil
	}
	return LoadFromDir(dir, "user")
}

// LoadProjectSPCs loads SPCs from ./.openlibing/spc/
func LoadProjectSPCs() ([]*spc.SPCDefinition, error) {
	return LoadFromDir(".openlibing/spc", "project")
}
```

- [ ] **Step 2: Write registry.go**

File: `internal/registry/registry.go`
```go
package registry

import (
	"fmt"
	"io/fs"
	"strings"
	"sync"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// RegistryImpl implements the SPC registry with 3-layer loading.
type RegistryImpl struct {
	mu    sync.RWMutex
	index map[string]*spc.SPCDefinition // name → definition
	all   []*spc.SPCDefinition          // all definitions in load order
}

// NewRegistry creates an empty registry.
func NewRegistry() *RegistryImpl {
	return &RegistryImpl{
		index: make(map[string]*spc.SPCDefinition),
	}
}

// LoadAll performs the full 3-layer discovery: builtin → user → project.
func (r *RegistryImpl) LoadAll(embedFS fs.FS, embedDir string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Layer 1: Built-in (embedded)
	builtins, _ := LoadFromEmbed(embedFS, embedDir)
	r.merge(builtins)

	// Layer 2: User custom (~/.openlibing/spc/)
	userSPCs, _ := LoadUserSPCs()
	r.merge(userSPCs)

	// Layer 3: Project local (./.openlibing/spc/)
	projectSPCs, _ := LoadProjectSPCs()
	r.merge(projectSPCs)

	return nil
}

// merge inserts or overrides definitions. Later layers take precedence.
func (r *RegistryImpl) merge(defs []*spc.SPCDefinition) {
	for _, def := range defs {
		// Check for override
		if existing, ok := r.index[def.Name]; ok {
			// Remove from all list
			for i, d := range r.all {
				if d == existing {
					r.all = append(r.all[:i], r.all[i+1:]...)
					break
				}
			}
		}
		r.index[def.Name] = def
		r.all = append(r.all, def)
	}
}

// Get retrieves an SPC by exact name.
func (r *RegistryImpl) Get(name string) (*spc.SPCDefinition, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	def, ok := r.index[name]
	if !ok {
		return nil, fmt.Errorf("SPC '%s' not found — run 'openlibing list' to see available super powers", name)
	}
	return def, nil
}

// ListAll returns all registered SPC definitions.
func (r *RegistryImpl) ListAll() []*spc.SPCDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]*spc.SPCDefinition, len(r.all))
	copy(result, r.all)
	return result
}

// ListByCategory returns SPCs filtered by category.
func (r *RegistryImpl) ListByCategory(cat string) []*spc.SPCDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []*spc.SPCDefinition
	for _, def := range r.all {
		if strings.EqualFold(def.Category, cat) {
			result = append(result, def)
		}
	}
	return result
}

// Search performs keyword-based SPC search. MVP: keyword matching only.
func (r *RegistryImpl) Search(query string) []*spc.SPCDefinition {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return ResolveSearch(query, r.all)
}
```

- [ ] **Step 3: Write resolver.go**

File: `internal/registry/resolver.go`
```go
package registry

import (
	"sort"
	"strings"

	"github.com/openlibing/openlibing-cli/pkg/spc"
)

// searchResult holds a scored SPC match.
type searchResult struct {
	def   *spc.SPCDefinition
	score float64
}

// ResolveSearch performs keyword-based matching of a query against SPC definitions.
// Scoring: name match * 3.0 + description match * 2.0 + tag match * 1.0 + category match * 1.5
func ResolveSearch(query string, defs []*spc.SPCDefinition) []*spc.SPCDefinition {
	query = strings.ToLower(query)
	var results []searchResult

	for _, def := range defs {
		score := 0.0

		// Name keyword match
		if strings.Contains(strings.ToLower(def.Name), query) {
			score += 3.0
		}

		// Description keyword match
		if strings.Contains(strings.ToLower(def.Description), query) {
			score += 2.0
		}

		// Category match
		if strings.Contains(strings.ToLower(def.Category), query) {
			score += 1.5
		}

		// Tag match
		for _, tag := range def.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				score += 1.0
				break
			}
		}

		// Also check individual words in query against name/desc
		words := strings.Fields(query)
		for _, word := range words {
			if len(word) < 3 {
				continue
			}
			for _, tag := range def.Tags {
				if strings.Contains(strings.ToLower(tag), word) {
					score += 0.5
				}
			}
		}

		if score > 0 {
			results = append(results, searchResult{def: def, score: score})
		}
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	result := make([]*spc.SPCDefinition, len(results))
	for i, r := range results {
		result[i] = r.def
	}
	return result
}
```

- [ ] **Step 4: Write registry_test.go**

File: `internal/registry/registry_test.go`
```go
package registry

import (
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"
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
```

- [ ] **Step 5: Run registry tests**

```bash
go test ./internal/registry/... -v
```

Expected: 5 tests PASS.

- [ ] **Step 6: Commit**

```bash
git add internal/registry/
git commit -m "feat(registry): add SPC discovery, 3-layer loading, and keyword search"
```

---

### Task 10: Built-in SPC Files

**Files:**
- Create: `embed/spc/pipeline-list.spc.yaml`
- Create: `embed/spc/pipeline-detail.spc.yaml`
- Create: `embed/spc/pipeline-logs.spc.yaml`
- Create: `embed/embed.go`

**Interfaces:**
- Consumes: SPC YAML format
- Produces: `embed.FS` with `//go:embed spc/*` — 3 built-in Super Powers

- [ ] **Step 1: Write embed.go**

File: `embed/embed.go`
```go
package embed

import "embed"

//go:embed spc/*
var SPCs embed.FS
```

- [ ] **Step 2: Write pipeline-list.spc.yaml**

File: `embed/spc/pipeline-list.spc.yaml`
```yaml
name: pipeline-list
version: "1.0"
description: >
  Query pipeline runs for a GitCode project.
  Returns status, branch, duration, and trigger info.
type: query
category: pipeline
tags: [ci, status, list]

parameters:
  - name: project_id
    type: string
    required: true
    description: GitCode project identifier

  - name: limit
    type: integer
    default: 10
    description: Max number of results to return
    validation:
      min: 1
      max: 100

  - name: status
    type: string
    enum: [running, success, failed, pending]
    description: Filter by pipeline status

source:
  method: GET
  endpoint: gateway/openlibing-cicd/project/pipeline/pipeline-run/detail
  query_params:
    projectId: "{{.project_id}}"
    pageSize: "{{.limit}}"
  headers:
    Content-Type: application/json

output:
  format: table
  fields:
    - name: id
      header: "ID"
      path: "pipelineRunId"
      width: 36
    - name: status
      header: "Status"
      path: "status"
      transform: upper
    - name: branch
      header: "Branch"
      path: "ref"
    - name: duration_ms
      header: "Duration"
      path: "durationMillis"
      transform: duration
    - name: created
      header: "Created"
      path: "createTime"

ai:
  prompt_hint: >
    Use this when users ask about pipeline runs, CI status, or build history
    for a specific project.
  natural_language:
    - "show pipelines for project 123"
    - "list recent builds"
    - "what's the CI status of project X"
  result_hint: >
    Focus on status and duration. If there are failures, highlight them first.

examples:
  - command: |
      openlibing run pipeline-list --project-id 123
    description: List last 10 pipelines for project 123
  - command: |
      openlibing run pipeline-list --project-id 123 --status failed --limit 5
    description: Show 5 most recent failures
```

- [ ] **Step 3: Write pipeline-detail.spc.yaml**

File: `embed/spc/pipeline-detail.spc.yaml`
```yaml
name: pipeline-detail
version: "1.0"
description: >
  Get detailed information about a specific pipeline run,
  including all stages and jobs with their status and duration.
type: query
category: pipeline
tags: [ci, detail, stages]

parameters:
  - name: run_id
    type: string
    required: true
    description: Pipeline run ID (pipelineRunId)

source:
  method: GET
  endpoint: gateway/openlibing-cicd/project/pipeline/pipeline-run/detail
  query_params:
    pipelineRunId: "{{.run_id}}"
  headers:
    Content-Type: application/json

output:
  format: json
  fields:
    - name: id
      header: "Pipeline ID"
      path: "pipelineRunId"
    - name: status
      header: "Status"
      path: "status"
    - name: branch
      header: "Branch"
      path: "ref"
    - name: duration_ms
      header: "Duration"
      path: "durationMillis"
      transform: duration
    - name: stages_count
      header: "Stages"
      path: "stages.#"

ai:
  prompt_hint: >
    Use this when users ask for details about a specific pipeline run,
    want to see stages and jobs, or need to understand why a build failed.
  natural_language:
    - "show details of pipeline run abc-123"
    - "what stages are in the latest build"
    - "why did pipeline X fail"

examples:
  - command: |
      openlibing run pipeline-detail --run-id abc-def-123
    description: Show full detail of a pipeline run
```

- [ ] **Step 4: Write pipeline-logs.spc.yaml**

File: `embed/spc/pipeline-logs.spc.yaml`
```yaml
name: pipeline-logs
version: "1.0"
description: >
  Fetch execution logs for a specific job step within a pipeline run.
  Requires project ID, pipeline run ID, job run ID, and step run ID.
type: query
category: pipeline
tags: [ci, logs, debug]

parameters:
  - name: project_id
    type: string
    required: true
    description: GitCode project identifier
  - name: pipeline_run_id
    type: string
    required: true
    description: Pipeline run ID
  - name: job_run_id
    type: string
    required: true
    description: Job run ID
  - name: step_run_id
    type: string
    required: true
    description: Step run ID

source:
  method: POST
  endpoint: gateway/openlibing-cicd/project/pipeline/exec-log
  query_params:
    projectId: "{{.project_id}}"
    pipelineRunId: "{{.pipeline_run_id}}"
    jobRunId: "{{.job_run_id}}"
    stepRunId: "{{.step_run_id}}"
  headers:
    Content-Type: application/json

output:
  format: raw

ai:
  prompt_hint: >
    Use this when users need to see raw execution logs for debugging
    a failed job or step.
  natural_language:
    - "show me the logs for the failed step"
    - "get execution logs for job X"
    - "what happened in step Y"

examples:
  - command: |
      openlibing run pipeline-logs \
        --project-id 123 \
        --pipeline-run-id run-abc \
        --job-run-id job-1 \
        --step-run-id step-a
    description: Fetch execution logs for a specific step
```

- [ ] **Step 5: Verify embed compiles**

File: `cmd/openlibing/main.go` (temporary)
```go
package main

import (
	"fmt"
	"io/fs"

	embedspc "github.com/openlibing/openlibing-cli/embed"
)

func main() {
	entries, _ := fs.ReadDir(embedspc.SPCs, "spc")
	for _, e := range entries {
		fmt.Println(e.Name())
	}
}
```

Run: `go run ./cmd/openlibing/`
Expected: lists 3 SPC file names.

- [ ] **Step 6: Commit**

```bash
git add embed/
git add cmd/openlibing/main.go
git commit -m "feat(embed): add 3 built-in SPC files for pipeline queries"
```

---

### Task 11: CLI Commands

**Files:**
- Create: `internal/cli/root.go`
- Create: `internal/cli/run.go`
- Create: `internal/cli/list.go`
- Create: `internal/cli/inspect.go`
- Create: `internal/cli/chat.go`

**Interfaces:**
- Consumes: `*engine.Engine`, `*registry.RegistryImpl`
- Produces:
  - `func NewRootCmd(engine *engine.Engine, reg *registry.RegistryImpl) *cobra.Command` — root with all subcommands wired

- [ ] **Step 1: Write root.go**

File: `internal/cli/root.go`
```go
package cli

import (
	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/openlibing/openlibing-cli/internal/registry"
	"github.com/spf13/cobra"
)

var (
	outputFormat string
)

// NewRootCmd creates the root command with all subcommands.
func NewRootCmd(eng *engine.Engine, reg *registry.RegistryImpl) *cobra.Command {
	root := &cobra.Command{
		Use:   "openlibing",
		Short: "OpenLibing CLI — AI-native CI/CD data query tool",
		Long: `openlibing is a CLI tool for querying OpenLibing CI/CD platform data.

All capabilities are defined as SPC (Skill Pipeline Configuration) files.
Use 'openlibing list' to see available Super Powers.
Use 'openlibing run <name>' to execute a Super Power.
Use 'openlibing inspect <name>' to view SPC details.
Use 'openlibing chat' to enter AI conversational mode.`,
		SilenceUsage: true,
	}

	root.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "Output format (table|json|yaml|raw)")

	root.AddCommand(newRunCmd(eng))
	root.AddCommand(newListCmd(reg))
	root.AddCommand(newInspectCmd(reg))
	root.AddCommand(newChatCmd(eng, reg))

	return root
}
```

- [ ] **Step 2: Write run.go**

File: `internal/cli/run.go`
```go
package cli

import (
	"fmt"
	"strings"

	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/openlibing/openlibing-cli/pkg/spc"
	"github.com/spf13/cobra"
)

func newRunCmd(eng *engine.Engine) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run <spc-name>",
		Short: "Execute a Super Power",
		Long:  "Execute a Super Power by name. Parameters are passed as flags (--param-name value).",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			spcName := args[0]

			// Collect all non-standard flags as parameters
			params := make(map[string]interface{})
			flags := cmd.Flags()
			flags.Visit(func(f *cobra.Flag) {
				if f.Name == "output" {
					return // skip global flags
				}
				params[f.Name] = f.Value.String()
			})

			result, err := eng.Execute(spcName, params)
			if err != nil {
				return err
			}

			// Override format if --output flag was set
			if outputFormat != "" {
				result.Format = outputFormat
			}

			output, err := FormatResult(result)
			if err != nil {
				return err
			}

			fmt.Print(output)
			return nil
		},
	}

	// Note: dynamic flag registration happens at execution time
	// The engine validates parameters; we accept arbitrary flags here
	cmd.Flags().String("project-id", "", "Project ID")
	cmd.Flags().String("run-id", "", "Pipeline run ID")
	cmd.Flags().String("pipeline-run-id", "", "Pipeline run ID")
	cmd.Flags().String("job-run-id", "", "Job run ID")
	cmd.Flags().String("step-run-id", "", "Step run ID")
	cmd.Flags().String("status", "", "Filter by status")
	cmd.Flags().Int("limit", 0, "Max results")

	return cmd
}

// Ensure unused imports compile
var _ = strings.TrimSpace
var _ = spc.SPCDefinition{}
```

- [ ] **Step 3: Write list.go**

File: `internal/cli/list.go`
```go
package cli

import (
	"fmt"

	"github.com/openlibing/openlibing-cli/internal/registry"
	"github.com/spf13/cobra"
)

func newListCmd(reg *registry.RegistryImpl) *cobra.Command {
	var category string
	var format string

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List available Super Powers",
		Long:  "List all discoverable SPC (Super Power) definitions.",
		RunE: func(cmd *cobra.Command, args []string) error {
			var defs []*spc.SPCDefinition
			if category != "" {
				defs = reg.ListByCategory(category)
			} else {
				defs = reg.ListAll()
			}

			if len(defs) == 0 {
				fmt.Println("No Super Powers found.")
				return nil
			}

			switch format {
			case "json":
				result := &spc.Result{
					Format: "json",
					Rows:   make([]map[string]interface{}, len(defs)),
				}
				for i, d := range defs {
					result.Rows[i] = map[string]interface{}{
						"name":        d.Name,
						"type":        d.Type,
						"category":    d.Category,
						"description": d.Description,
						"origin":      d.Origin,
						"version":     d.Version,
					}
				}
				out, _ := FormatJSON(result.Rows)
				fmt.Print(out)
			default:
				// Table format
				rows := make([]map[string]interface{}, len(defs))
				for i, d := range defs {
					rows[i] = map[string]interface{}{
						"name":        d.Name,
						"type":        d.Type,
						"category":    d.Category,
						"description": truncateString(d.Description, 60),
					}
				}
				fields := []spc.Field{
					{Name: "name", Header: "NAME"},
					{Name: "type", Header: "TYPE"},
					{Name: "category", Header: "CATEGORY"},
					{Name: "description", Header: "DESCRIPTION"},
				}
				fmt.Print(FormatTable(rows, fields))
			}

			return nil
		},
	}

	cmd.Flags().StringVar(&category, "category", "", "Filter by category (pipeline|codecheck|sca)")
	cmd.Flags().StringVar(&format, "format", "table", "Output format (table|json)")

	return cmd
}

func truncateString(s string, maxLen int) string {
	if len(s) > maxLen {
		return s[:maxLen-3] + "..."
	}
	return s
}
```

- [ ] **Step 4: Write inspect.go**

File: `internal/cli/inspect.go`
```go
package cli

import (
	"fmt"
	"strings"

	"github.com/openlibing/openlibing-cli/internal/registry"
	"github.com/spf13/cobra"
)

func newInspectCmd(reg *registry.RegistryImpl) *cobra.Command {
	return &cobra.Command{
		Use:   "inspect <spc-name>",
		Short: "Show full definition of a Super Power",
		Long:  "Display the complete SPC definition including parameters, output fields, examples, and AI hints.",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			def, err := reg.Get(args[0])
			if err != nil {
				return err
			}

			fmt.Printf("Name:        %s\n", def.Name)
			fmt.Printf("Version:     %s\n", def.Version)
			fmt.Printf("Type:        %s\n", def.Type)
			fmt.Printf("Category:    %s\n", def.Category)
			fmt.Printf("Origin:      %s\n", def.Origin)
			fmt.Printf("Tags:        %s\n", strings.Join(def.Tags, ", "))
			fmt.Printf("\nDescription:\n  %s\n", def.Description)

			if len(def.Parameters) > 0 {
				fmt.Printf("\nParameters:\n")
				for _, p := range def.Parameters {
					required := ""
					if p.Required {
						required = " (required)"
					}
					defaultVal := ""
					if p.Default != nil {
						defaultVal = fmt.Sprintf(" [default: %v]", p.Default)
					}
					fmt.Printf("  --%s  %s%s%s\n", toFlagName(p.Name), p.Type, required, defaultVal)
					if p.Description != "" {
						fmt.Printf("        %s\n", p.Description)
					}
					if len(p.Enum) > 0 {
						fmt.Printf("        Values: %s\n", strings.Join(p.Enum, ", "))
					}
				}
			}

			fmt.Printf("\nSource:\n")
			fmt.Printf("  %s %s\n", def.Source.Method, def.Source.Endpoint)

			fmt.Printf("\nOutput: %s", def.Output.Format)
			if len(def.Output.Fields) > 0 {
				fmt.Printf(" (%d fields)\n", len(def.Output.Fields))
				for _, f := range def.Output.Fields {
					fmt.Printf("  %-15s → %s", f.Name, f.Header)
					if f.Transform != "" {
						fmt.Printf(" [%s]", f.Transform)
					}
					fmt.Println()
				}
			} else {
				fmt.Println()
			}

			if len(def.Examples) > 0 {
				fmt.Printf("\nExamples:\n")
				for _, ex := range def.Examples {
					fmt.Printf("  %s\n", ex.Command)
					fmt.Printf("  # %s\n\n", ex.Description)
				}
			}

			return nil
		},
	}
}

// toFlagName converts parameter name to CLI flag format.
func toFlagName(name string) string {
	return strings.ReplaceAll(name, "_", "-")
}
```

- [ ] **Step 5: Write chat.go (stub)**

File: `internal/cli/chat.go`
```go
package cli

import (
	"fmt"

	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/openlibing/openlibing-cli/internal/registry"
	"github.com/spf13/cobra"
)

func newChatCmd(eng *engine.Engine, reg *registry.RegistryImpl) *cobra.Command {
	return &cobra.Command{
		Use:   "chat",
		Short: "Enter AI conversational mode (coming soon)",
		Long: `chat launches an interactive AI-powered REPL for querying OpenLibing data.

In chat mode, you use natural language to describe what you want to query.
The AI automatically selects the right Super Power, fills in parameters,
and explains the results.

Note: Chat mode requires LLM API credentials in ~/.openlibing/auth.yaml`,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Chat mode is coming in a future release.")
			fmt.Println("For now, use 'openlibing run <spc-name>' for direct queries.")
			fmt.Println("Run 'openlibing list' to see available Super Powers.")
			return nil
		},
	}
}
```

- [ ] **Step 6: Fix main.go to use cobra**

File: `cmd/openlibing/main.go`
```go
package main

import (
	"os"

	embedspc "github.com/openlibing/openlibing-cli/embed"
	"github.com/openlibing/openlibing-cli/internal/api"
	"github.com/openlibing/openlibing-cli/internal/cli"
	"github.com/openlibing/openlibing-cli/internal/config"
	"github.com/openlibing/openlibing-cli/internal/engine"
	"github.com/openlibing/openlibing-cli/internal/registry"
)

func main() {
	// Load configuration
	cfg, _ := config.LoadConfig()
	auth, _ := config.LoadAuth()

	// Initialize API client
	client := api.NewClient(cfg, auth)

	// Initialize Registry and load SPCs (3-layer)
	reg := registry.NewRegistry()
	if err := reg.LoadAll(embedspc.SPCs, "spc"); err != nil {
		// Non-fatal: registry errors are logged but CLI still works
		// (at minimum, run/openlibing list will show empty)
	}

	// Initialize Engine
	eng := engine.NewEngine(reg, client)

	// Build and execute CLI
	rootCmd := cli.NewRootCmd(eng, reg)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
```

- [ ] **Step 7: Build and smoke test**

```bash
go build -o bin/openlibing ./cmd/openlibing/

# Test list command
./bin/openlibing list
# Expected: table with 3 built-in SPCs

# Test inspect command
./bin/openlibing inspect pipeline-list
# Expected: full SPC definition displayed

# Test list with JSON
./bin/openlibing list --format json
# Expected: JSON array with 3 entries

# Test help
./bin/openlibing --help
# Expected: usage with all subcommands
```

- [ ] **Step 8: Commit**

```bash
git add internal/cli/ cmd/openlibing/main.go
git commit -m "feat(cli): add cobra commands (run, list, inspect, chat stub)"
```

---

### Task 12: Integration Test and Final Verification

**Files:**
- No new files — run full test suite and verify the binary

**Interfaces:**
- Consumes: everything
- Produces: green test suite, working binary

- [ ] **Step 1: Run all unit tests**

```bash
go test ./... -v -count=1
```

Expected: ALL tests PASS (~30+ tests across all packages).

- [ ] **Step 2: Run go vet**

```bash
go vet ./...
```

Expected: No issues found.

- [ ] **Step 3: Final build**

```bash
make build
```

Expected: `bin/openlibing` binary produced with no errors.

- [ ] **Step 4: Manual end-to-end test with mock server**

```bash
# Start a mock server that returns pipeline data (use netcat or python)
python3 -c "
import http.server, json
class H(http.server.BaseHTTPRequestHandler):
    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-Type', 'application/json')
        self.end_headers()
        self.wfile.write(json.dumps({'data': [
            {'pipelineRunId': 'abc', 'status': 'SUCCESS', 'ref': 'main', 'durationMillis': 120000, 'createTime': '2026-01-01'}
        ]}).encode())
    def do_POST(self):
        self.do_GET()
H = http.server.HTTPServer
h = H(('', 9999), H)
h.serve_forever()
" &
MOCK_PID=$!
sleep 1

# Run query against mock
./bin/openlibing run pipeline-list --project-id 123 --endpoint http://localhost:9999
# Expected: table with 1 row

kill $MOCK_PID
```

- [ ] **Step 5: Commit and tag**

```bash
git add -A
git commit -m "chore: final integration verification — all tests pass"
```

---

## Implementation Order Summary

```
Task 1  → Scaffold (go mod, dirs, Makefile)
Task 2  → SPC Types (pkg/spc/types.go)
Task 3  → Config (internal/config/)
Task 4  → API Client (internal/api/)
Task 5  → Parser + Validator (internal/engine/)
Task 6  → Template Resolver (internal/engine/)
Task 7  → Output Formatter (internal/cli/)
Task 8  → Engine Core (internal/engine/)
Task 9  → Registry (internal/registry/)
Task 10 → Built-in SPCs (embed/spc/)
Task 11 → CLI Commands (internal/cli/)
Task 12 → Integration & Final Verify
```

Tasks 2-4 can run in parallel (no mutual dependencies).
Tasks 5-7 can run in parallel after Task 2 (share types only).
Tasks 8-11 are sequential (each builds on the previous).

---
🤖 Generated with [Claude Code](https://claude.com/claude-code)
