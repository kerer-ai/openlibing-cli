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
