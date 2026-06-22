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
