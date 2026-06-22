package main

import "testing"

func TestMainOutput(t *testing.T) {
	// Test that main doesn't panic
	// (functional test; binary compilation verified by make build)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("main panicked: %v", r)
		}
	}()
	main()
}
