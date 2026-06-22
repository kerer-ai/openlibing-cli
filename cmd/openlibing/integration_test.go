package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/openlibing/openlibing-cli/internal/testutil"
)

// buildBinary compiles the CLI binary for integration tests.
// Returns an absolute path to avoid resolution issues when cmd.Dir is set.
func buildBinary(t *testing.T) string {
	t.Helper()
	relPath := "../../bin/openlibing"
	absPath, err := filepath.Abs(relPath)
	if err != nil {
		t.Fatalf("resolve binary path: %v", err)
	}
	cmd := exec.Command("go", "build", "-o", absPath, ".")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	return absPath
}

// runCLI executes the openlibing binary with args, returns stdout+stderr and exit code.
// Uses os.Environ() to inherit HOME set by SetupTestHome.
func runCLI(t *testing.T, bin string, args ...string) (string, int) {
	t.Helper()
	cmd := exec.Command(bin, args...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		}
	}
	return string(out), code
}

// ============================================================================
// Integration Tests — full CLI against mock OpenLibing server
// ============================================================================

func TestIntegration_FullPipeline(t *testing.T) {
	server := testutil.NewMockOpenLibingServer(nil, "")
	defer server.Close()

	_, cleanup := testutil.SetupTestHome(t, server.URL)
	defer cleanup()

	bin := buildBinary(t)

	// ── list ──────────────────────────────────────────────
	t.Run("list_table", func(t *testing.T) {
		out, code := runCLI(t, bin, "list")
		if code != 0 {
			t.Fatalf("exit %d", code)
		}
		for _, name := range []string{"pipeline-list", "pipeline-detail", "pipeline-logs"} {
			if !strings.Contains(out, name) {
				t.Errorf("list output missing SPC %q\n%s", name, out)
			}
		}
	})

	t.Run("list_json", func(t *testing.T) {
		out, code := runCLI(t, bin, "list", "--format", "json")
		if code != 0 {
			t.Fatalf("exit %d", code)
		}
		var rows []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &rows); err != nil {
			t.Fatalf("invalid JSON: %v\n%s", err, out)
		}
		if len(rows) != 3 {
			t.Errorf("expected 3 SPCs, got %d", len(rows))
		}
		for _, r := range rows {
			if r["origin"] != "builtin" {
				t.Errorf("origin = %q, want builtin", r["origin"])
			}
		}
	})

	t.Run("list_category", func(t *testing.T) {
		out, code := runCLI(t, bin, "list", "--category", "pipeline")
		if code != 0 {
			t.Fatalf("exit %d", code)
		}
		if !strings.Contains(out, "pipeline-list") {
			t.Errorf("filtered list missing pipeline-list\n%s", out)
		}
	})

	t.Run("list_empty_category", func(t *testing.T) {
		out, code := runCLI(t, bin, "list", "--category", "nonexistent")
		if code != 0 {
			t.Fatalf("exit %d", code)
		}
		if !strings.Contains(out, "No Super Powers") {
			t.Errorf("expected 'No Super Powers' for empty category\n%s", out)
		}
	})

	// ── inspect ──────────────────────────────────────────
	t.Run("inspect_pipeline_list", func(t *testing.T) {
		out, code := runCLI(t, bin, "inspect", "pipeline-list")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		checks := []string{
			"pipeline-list", "query", "pipeline", "builtin",
			"--project-id", "required",
			"--limit", "10",
			"--status", "running", "failed",
			"GET", "pipeline-run/detail",
			"ID", "Status", "Branch", "Duration", "Created",
		}
		for _, c := range checks {
			if !strings.Contains(out, c) {
				t.Errorf("inspect pipeline-list missing %q\n%s", c, out)
			}
		}
	})

	t.Run("inspect_pipeline_detail", func(t *testing.T) {
		out, code := runCLI(t, bin, "inspect", "pipeline-detail")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		if !strings.Contains(out, "--run-id") {
			t.Errorf("missing --run-id param\n%s", out)
		}
		if !strings.Contains(out, "json") {
			t.Errorf("missing json output format\n%s", out)
		}
	})

	t.Run("inspect_pipeline_logs", func(t *testing.T) {
		out, code := runCLI(t, bin, "inspect", "pipeline-logs")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		if !strings.Contains(out, "POST") {
			t.Errorf("missing POST method\n%s", out)
		}
	})

	t.Run("inspect_not_found", func(t *testing.T) {
		_, code := runCLI(t, bin, "inspect", "nonexistent")
		if code == 0 {
			t.Fatal("expected non-zero exit for missing SPC")
		}
	})

	t.Run("inspect_no_args", func(t *testing.T) {
		_, code := runCLI(t, bin, "inspect")
		if code == 0 {
			t.Fatal("expected non-zero exit for missing args")
		}
	})

	// ── run: pipeline-list ───────────────────────────────
	t.Run("run_pipeline_list_table", func(t *testing.T) {
		out, code := runCLI(t, bin, "run", "pipeline-list", "--project-id", "123", "--limit", "3")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		for _, id := range []string{"run-abc-001", "run-def-002", "run-ghi-003"} {
			if !strings.Contains(out, id) {
				t.Errorf("output missing pipeline %q\n%s", id, out)
			}
		}
		if !strings.Contains(out, "SUCCESS") {
			t.Error("status transform (upper) not applied")
		}
	})

	t.Run("run_pipeline_list_json", func(t *testing.T) {
		out, code := runCLI(t, bin, "run", "pipeline-list", "--project-id", "123", "--limit", "3", "--output", "json")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		var rows []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &rows); err != nil {
			t.Fatalf("invalid JSON: %v\n%s", err, out)
		}
		if len(rows) != 3 {
			t.Errorf("expected 3 runs, got %d", len(rows))
		}
	})

	t.Run("run_pipeline_list_yaml", func(t *testing.T) {
		out, code := runCLI(t, bin, "run", "pipeline-list", "--project-id", "123", "--limit", "3", "--output", "yaml")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		if !strings.Contains(out, "id:") || !strings.Contains(out, "status:") {
			t.Errorf("yaml output missing expected keys\n%s", out)
		}
	})

	t.Run("run_missing_required", func(t *testing.T) {
		_, code := runCLI(t, bin, "run", "pipeline-list")
		if code == 0 {
			t.Fatal("expected non-zero exit for missing project_id")
		}
	})

	t.Run("run_nonexistent_spc", func(t *testing.T) {
		_, code := runCLI(t, bin, "run", "nonexistent-spc")
		if code == 0 {
			t.Fatal("expected non-zero exit for missing SPC")
		}
	})

	// ── run: pipeline-detail ─────────────────────────────
	t.Run("run_pipeline_detail", func(t *testing.T) {
		out, code := runCLI(t, bin, "run", "pipeline-detail", "--run-id", "run-abc-001", "--output", "json")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		var rows []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &rows); err != nil {
			t.Fatalf("invalid JSON: %v\n%s", err, out)
		}
		if len(rows) == 0 {
			t.Fatal("expected at least 1 result")
		}
	})

	// ── run: pipeline-logs (POST) ────────────────────────
	t.Run("run_pipeline_logs_post", func(t *testing.T) {
		out, code := runCLI(t, bin, "run", "pipeline-logs",
			"--project-id", "123",
			"--pipeline-run-id", "run-abc",
			"--job-run-id", "job-1",
			"--step-run-id", "step-a")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		if !strings.Contains(out, "Build SUCCESS") {
			t.Errorf("log output missing expected content\n%s", out)
		}
	})

	// ── chat ─────────────────────────────────────────────
	t.Run("chat_stub", func(t *testing.T) {
		out, code := runCLI(t, bin, "chat")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		if !strings.Contains(out, "Chat mode is coming") {
			t.Errorf("chat stub missing expected message\n%s", out)
		}
	})

	// ── help ─────────────────────────────────────────────
	t.Run("help_all_commands", func(t *testing.T) {
		out, code := runCLI(t, bin, "--help")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		for _, cmd := range []string{"run", "list", "inspect", "chat"} {
			if !strings.Contains(out, cmd) {
				t.Errorf("help missing command %q\n%s", cmd, out)
			}
		}
	})

	t.Run("run_help_flags", func(t *testing.T) {
		out, code := runCLI(t, bin, "run", "--help")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		for _, flag := range []string{"--project-id", "--run-id", "--limit", "--output"} {
			if !strings.Contains(out, flag) {
				t.Errorf("run --help missing flag %q\n%s", flag, out)
			}
		}
	})
}

// ============================================================================
// Flag name normalization: hyphens → underscores
// ============================================================================

func TestFlagNameNormalization(t *testing.T) {
	server := testutil.NewMockOpenLibingServer(nil, "")
	defer server.Close()

	_, cleanup := testutil.SetupTestHome(t, server.URL)
	defer cleanup()

	bin := buildBinary(t)

	out, code := runCLI(t, bin, "run", "pipeline-list", "--project-id", "test-456", "--limit", "1", "--output", "json")
	if code != 0 {
		t.Fatalf("flag normalization failed: exit %d\n%s", code, out)
	}

	var rows []map[string]interface{}
	if err := json.Unmarshal([]byte(out), &rows); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out)
	}
	if len(rows) == 0 {
		t.Fatal("expected data from mock, got empty — flag normalization may have failed")
	}
}

// ============================================================================
// Offline: local commands work without network
// ============================================================================

func TestNoRealNetworkDependency(t *testing.T) {
	home, cleanup := testutil.SetupTestHome(t, "http://127.0.0.1:1")
	defer cleanup()

	bin := buildBinary(t)

	t.Run("list_offline", func(t *testing.T) {
		out, code := runCLI(t, bin, "list")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		if !strings.Contains(out, "pipeline-list") {
			t.Error("list should work without network")
		}
	})

	t.Run("inspect_offline", func(t *testing.T) {
		out, code := runCLI(t, bin, "inspect", "pipeline-list")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		if !strings.Contains(out, "project-id") {
			t.Error("inspect should work without network")
		}
	})

	t.Run("help_offline", func(t *testing.T) {
		_, code := runCLI(t, bin, "--help")
		if code != 0 {
			t.Fatal("--help should work without network")
		}
	})

	t.Run("chat_offline", func(t *testing.T) {
		_, code := runCLI(t, bin, "chat")
		if code != 0 {
			t.Fatal("chat stub should work without network")
		}
	})

	_ = home
}

// ============================================================================
// Custom SPC: verify user-defined SPCs are discovered and executable
// ============================================================================

func TestCustomSPCDiscovery(t *testing.T) {
	server := testutil.NewMockOpenLibingServer(nil, "")
	defer server.Close()

	home, cleanup := testutil.SetupTestHome(t, server.URL)
	defer cleanup()

	// Write custom SPC to user spc dir
	spcDir := filepath.Join(home, ".openlibing", "spc")
	os.MkdirAll(spcDir, 0755)

	customSPC := `name: my-custom-test
version: "1.0"
description: Custom test SPC
type: query
category: test
parameters:
  - name: search
    type: string
    required: true
source:
  method: GET
  endpoint: gateway/openlibing-cicd/project/pipeline/pipeline-run/detail
  query_params:
    projectId: "{{.search}}"
output:
  format: json
  fields:
    - name: id
      header: "ID"
      path: "pipelineRunId"
    - name: status
      header: "Status"
      path: "status"
`
	os.WriteFile(filepath.Join(spcDir, "my-custom-test.spc.yaml"), []byte(customSPC), 0644)

	bin := buildBinary(t)

	// 1. Custom SPC appears in list
	t.Run("custom_in_list", func(t *testing.T) {
		out, code := runCLI(t, bin, "list")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		if !strings.Contains(out, "my-custom-test") {
			t.Errorf("custom SPC not found in list\n%s", out)
		}
	})

	// 2. Custom SPC is inspectable
	t.Run("custom_inspect", func(t *testing.T) {
		out, code := runCLI(t, bin, "inspect", "my-custom-test")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		if !strings.Contains(out, "Custom test SPC") {
			t.Errorf("inspect missing description\n%s", out)
		}
	})

	// 3. Custom SPC executes and returns data
	t.Run("custom_run", func(t *testing.T) {
		out, code := runCLI(t, bin, "run", "my-custom-test", "--project-id", "test", "--param", "search=my-project", "--output", "json")
		if code != 0 {
			t.Fatalf("exit %d\n%s", code, out)
		}
		var rows []map[string]interface{}
		if err := json.Unmarshal([]byte(out), &rows); err != nil {
			t.Fatalf("invalid JSON: %v\n%s", err, out)
		}
		if len(rows) == 0 {
			t.Fatal("expected data, got empty")
		}
	})

	_ = home
}
