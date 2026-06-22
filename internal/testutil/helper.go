package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// SetupTestHome creates a temporary HOME directory with an openlibing config
// that points to the given server URL. Returns the path and a cleanup func.
//
// Usage:
//
//	home, cleanup := testutil.SetupTestHome(t, server.URL)
//	defer cleanup()
func SetupTestHome(t *testing.T, serverURL string) (home string, cleanup func()) {
	t.Helper()

	// Use os.MkdirTemp instead of t.TempDir() to avoid Go module cache
	// cleanup issues on WSL2 (read-only files in GOPATH/pkg/mod).
	home, err := os.MkdirTemp("", "openlibing-test-*")
	if err != nil {
		t.Fatalf("create temp home: %v", err)
	}

	openlibingDir := filepath.Join(home, ".openlibing")
	os.MkdirAll(openlibingDir, 0755)

	configYAML := `endpoint: ` + serverURL + `
defaults:
  limit: 10
output:
  format: table
  color: false
`
	os.WriteFile(filepath.Join(openlibingDir, "config.yaml"), []byte(configYAML), 0644)

	// Save and restore HOME
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", home)

	cleanup = func() {
		os.Setenv("HOME", oldHome)
		os.RemoveAll(home) // best-effort cleanup
	}

	return home, cleanup
}
