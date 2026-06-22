.PHONY: build test run clean lint test-cover test-unit test-integration test-all check

BINARY=openlibing
BIN_DIR=bin

# Quick build
build:
	go build -o $(BIN_DIR)/$(BINARY) ./cmd/openlibing/

# ── Test targets ──────────────────────────────────────

# All tests (unit + integration)
test:
	go test ./... -v -count=1

# Unit tests only (fast, no binary build needed)
test-unit:
	go test $$(go list ./... | grep -v /cmd/) -v -count=1

# Integration tests (build binary, run against mock server)
test-integration: build
	go test ./cmd/openlibing/ -v -count=1 -run 'TestIntegration|TestFlag|TestCustom|TestNoReal'

# Coverage report
test-cover:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	@echo ""
	@echo "HTML report: go tool cover -html=coverage.out"

# ── Quality gates (CI-ready) ──────────────────────────

# Full check: build + unit + integration + vet
# Use this as the CI gate and pre-commit hook
check: build test-unit test-integration lint
	@echo ""
	@echo "========================================"
	@echo "  All checks passed — ready to commit"
	@echo "========================================"

# ── Utility ──────────────────────────────────────────

run: build
	./$(BIN_DIR)/$(BINARY)

clean:
	rm -rf $(BIN_DIR)/
	rm -f coverage.out

lint:
	go vet ./...
