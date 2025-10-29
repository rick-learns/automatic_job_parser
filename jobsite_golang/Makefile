.PHONY: deps build seed e2e serve clean test vet fmt static lint sec secrets vac lock-info help frontend frontend-serve build-all

# Version info for build
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X 'main.version=$(VERSION)' -X 'main.commit=$(COMMIT)' -X 'main.date=$(DATE)'

# Default target
help:
	@echo "Jobsite - QA/SDET Job Parser"
	@echo ""
	@echo "Targets:"
	@echo "  make deps       - Download Go dependencies"
	@echo "  make build      - Build binary with version info"
	@echo "  make seed       - Build and run with seed data"
	@echo "  make e2e        - Build → seed → validate → report"
	@echo "  make serve      - Serve public/latest on http://localhost:8080"
	@echo "  make frontend   - Build React frontend"
	@echo "  make frontend-serve - Serve React frontend in dev mode"
	@echo "  make build-all  - Build Go binary + React frontend"
	@echo "  make clean      - Remove build artifacts and outputs"
	@echo "  make test       - Run tests with race detector"
	@echo "  make vet        - Run go vet"
	@echo "  make fmt        - Format code"
	@echo "  make static     - Run staticcheck"
	@echo "  make vac        - Run SQLite VACUUM/ANALYZE"
	@echo "  make lock-info  - Check lock file status"
	@echo ""
	@echo "Build Info:"
	@echo "  Version: $(VERSION)"
	@echo "  Commit:  $(COMMIT)"
	@echo "  Date:    $(DATE)"

# Download dependencies
deps:
	go mod download

# Build binary with optimized flags
build:
	go build -trimpath -ldflags "$(LDFLAGS)" -o jobsite ./cmd/jobsite
	@echo "Built: jobsite (version: $(VERSION), commit: $(COMMIT))"

# Build React frontend
frontend:
	cd ../frontend && npm ci && npm run build
	@echo "Frontend built successfully"

# Serve React frontend in dev mode
frontend-serve:
	cd ../frontend && npm run dev

# Build everything: Go binary + React frontend
build-all: build frontend
	@echo "All targets built successfully"

# Build and run with seed data
seed: build
	@echo "Running seed mode..."
	./jobsite seed

# End-to-end test: build → seed → validate
e2e: build
	@echo "=== E2E Test ==="
	./jobsite seed
	@echo "Validating outputs..."
	@test -f public/latest/index.html || (echo "❌ index.html missing"; exit 1)
	@test -s public/latest/jobs.csv || (echo "❌ jobs.csv missing or empty"; exit 1)
	@which jq >/dev/null && jq -e 'length>0' public/latest/jobs.json >/dev/null || echo "⚠️  jq not available, skipping JSON validation"
	@echo "✅ E2E test passed"

# Serve generated site locally
serve:
	@echo "Serving on http://localhost:8080 (Ctrl+C to stop)"
	python3 -m http.server -d public/latest 8080

# Clean build artifacts and outputs
clean:
	rm -f jobsite
	rm -f public/latest/*.html public/latest/*.csv public/latest/*.json
	rm -f *.lock
	@echo "Clean complete"

# Run tests with race detector
test:
	go test ./... -race -count=1

# Run go vet
vet:
	go vet ./...

# Format code
fmt:
	go fmt ./...

# Run staticcheck (optional, install: go install honnef.co/go/tools/cmd/staticcheck@latest)
static:
	@which staticcheck >/dev/null && staticcheck ./... || echo "⚠️  staticcheck not installed"

# Run golangci-lint (optional)
lint:
	@which golangci-lint >/dev/null && golangci-lint run ./... || echo "⚠️  golangci-lint not installed"

# Run gosec security scan (optional)
sec:
	@which gosec >/dev/null && gosec ./... || echo "⚠️  gosec not installed"

# Run gitleaks secrets scan (optional)
secrets:
	@which gitleaks >/dev/null && gitleaks detect --no-banner --redact || echo "⚠️  gitleaks not installed"

# SQLite maintenance
vac:
	@test -f data/jobs.sqlite && sqlite3 data/jobs.sqlite "VACUUM; ANALYZE;" || echo "⚠️  Database not found"

# Check lock file status
lock-info:
	@if [ -f jobsite.lock ]; then \
		echo "Lock file exists:"; \
		cat jobsite.lock; \
	else \
		echo "No lock file found"; \
	fi

