package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestMainOutput(t *testing.T) {
	// Build the binary
	buildCmd := exec.Command("go", "build", "-o", "test-bin", ".")
	buildCmd.Dir = "."
	out, err := buildCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("build failed: %v\n%s", err, out)
	}
	defer os.Remove("test-bin")

	// Run the binary
	runCmd := exec.Command("./test-bin")
	runOut, err := runCmd.CombinedOutput()
	if err != nil {
		t.Fatalf("run failed: %v\n%s", err, runOut)
	}

	expected := "openlibing-cli\n"
	got := string(runOut)
	if got != expected {
		t.Fatalf("expected %q, got %q", expected, got)
	}
}
