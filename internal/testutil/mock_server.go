// Package testutil provides reusable test helpers for openlibing-cli:
// mock HTTP servers, temp config directories, and common fixtures.
package testutil

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
)

// PipelineRun is a simplified pipeline run fixture.
type PipelineRun struct {
	PipelineRunID string `json:"pipelineRunId"`
	Status        string `json:"status"`
	Ref           string `json:"ref"`
	DurationMillis int64 `json:"durationMillis"`
	CreateTime    string `json:"createTime"`
}

// DefaultPipelineRuns returns a standard set of test fixtures.
func DefaultPipelineRuns() []PipelineRun {
	return []PipelineRun{
		{PipelineRunID: "run-abc-001", Status: "SUCCESS", Ref: "main", DurationMillis: 125000, CreateTime: "2026-06-22T10:30:00Z"},
		{PipelineRunID: "run-def-002", Status: "FAILED", Ref: "feat/new-thing", DurationMillis: 380000, CreateTime: "2026-06-22T09:15:00Z"},
		{PipelineRunID: "run-ghi-003", Status: "RUNNING", Ref: "fix/urgent", DurationMillis: 45000, CreateTime: "2026-06-22T11:00:00Z"},
	}
}

// NewMockOpenLibingServer creates an httptest server that mimics the
// OpenLibing gateway API. It responds to:
//
//	GET  /gateway/openlibing-cicd/project/pipeline/pipeline-run/detail → {"data": [...]}
//	POST /gateway/openlibing-cicd/project/pipeline/exec-log          → raw log text
func NewMockOpenLibingServer(runs []PipelineRun, execLog string) *httptest.Server {
	if runs == nil {
		runs = DefaultPipelineRuns()
	}
	if execLog == "" {
		execLog = "[2026-06-22 10:30:01] Starting build...\n[2026-06-22 10:32:00] Build SUCCESS"
	}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		switch {
		case strings.Contains(path, "pipeline-run/detail"):
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"data": runs,
			})

		case strings.Contains(path, "exec-log"):
			w.Header().Set("Content-Type", "text/plain")
			w.Write([]byte(execLog))

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code":404,"msg":"not found"}`))
		}
	}))
}
