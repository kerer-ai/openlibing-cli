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
