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

// VersionAvailabilityRecord is a test fixture for the version-availability SPC.
type VersionAvailabilityRecord struct {
	ProjectID               int     `json:"projectId"`
	PipelineID              string  `json:"pipelineId"`
	PipelineName            string  `json:"pipelineName"`
	PipelineRunCount        int     `json:"pipelineRunCount"`
	PipelineSuccessRate     float64 `json:"pipelineSuccessRate"`
	BuildSuccessRate        float64 `json:"buildSuccessRate"`
	TestSuccessRate         float64 `json:"testSuccessRate"`
	VersionAvailabilityRate float64 `json:"versionAvailabilityRate"`
	RunFrequencyPerDay      float64 `json:"runFrequencyPerDay"`
	ActualDurationAvgMinutes float64 `json:"actualDurationAvgMinutes"`
}

// DefaultVersionAvailabilityRecords returns standard test fixtures.
func DefaultVersionAvailabilityRecords() []VersionAvailabilityRecord {
	return []VersionAvailabilityRecord{
		{PipelineName: "Nightly-CI_MindIE-Motor", PipelineRunCount: 29, PipelineSuccessRate: 51.72, BuildSuccessRate: 92.0, TestSuccessRate: 58.82, VersionAvailabilityRate: 63.64, RunFrequencyPerDay: 1.32, ActualDurationAvgMinutes: 59.59},
		{PipelineName: "Nightly-CI_ATB-Models", PipelineRunCount: 22, PipelineSuccessRate: 100.0, BuildSuccessRate: 100.0, TestSuccessRate: 100.0, VersionAvailabilityRate: 95.45, RunFrequencyPerDay: 1.0, ActualDurationAvgMinutes: 22.18},
		{PipelineName: "Nightly-CI_MindIE-LLM", PipelineRunCount: 22, PipelineSuccessRate: 77.27, BuildSuccessRate: 100.0, TestSuccessRate: 88.64, VersionAvailabilityRate: 81.82, RunFrequencyPerDay: 1.0, ActualDurationAvgMinutes: 113.61},
	}
}

// NewMockOpenLibingServer creates an httptest server that mimics the
// OpenLibing gateway API. It responds to:
//
//	GET  /gateway/openlibing-cicd/project/pipeline/pipeline-run/detail → {"data": [...]}
//	POST /gateway/openlibing-cicd/project/pipeline/exec-log          → raw log text
//	POST /gateway/openlibing-ops/manage/common/detail               → {"data":{"records":[...]}}
func NewMockOpenLibingServer(runs []PipelineRun, execLog string) *httptest.Server {
	if runs == nil {
		runs = DefaultPipelineRuns()
	}
	if execLog == "" {
		execLog = "[2026-06-22 10:30:01] Starting build...\n[2026-06-22 10:32:00] Build SUCCESS"
	}

	versionRecords := DefaultVersionAvailabilityRecords()

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

		case strings.Contains(path, "manage/common/detail"):
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"code":      200,
				"messageCn": "成功",
				"data": map[string]interface{}{
					"records":  versionRecords,
					"total":    len(versionRecords),
					"pageSize": 10,
					"page":     1,
				},
			})

		default:
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(`{"code":404,"msg":"not found"}`))
		}
	}))
}
