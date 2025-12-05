#!/bin/bash

# CI Checks Runner Script
# This script runs all the same checks that CI runs, so you can verify locally before pushing

set -e  # Exit on any error

echo "=========================================="
echo "  Running CI Checks Locally"
echo "=========================================="
echo ""

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Track overall status
OVERALL_STATUS=0

# Function to print status
print_status() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✅ $2${NC}"
    else
        echo -e "${RED}❌ $2${NC}"
        OVERALL_STATUS=1
    fi
}

echo "1. Running go vet..."
echo "-------------------------------------------"
if go vet ./...; then
    print_status 0 "go vet passed"
else
    print_status 1 "go vet failed"
fi
echo ""

echo "2. Running go fmt check..."
echo "-------------------------------------------"
UNFORMATTED=$(gofmt -s -l . | grep -v vendor || true)
if [ -z "$UNFORMATTED" ]; then
    print_status 0 "go fmt passed"
else
    print_status 1 "go fmt failed - The following files need formatting:"
    echo "$UNFORMATTED"
    echo ""
    echo "Run: gofmt -s -w ."
fi
echo ""

echo "3. Running tests with coverage..."
echo "-------------------------------------------"
if go test -coverprofile=coverage.out -covermode=atomic ./internal/modules/...; then
    print_status 0 "All tests passed"
else
    print_status 1 "Tests failed"
fi
echo ""

echo "4. Checking coverage threshold..."
echo "-------------------------------------------"
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
THRESHOLD=37.0

echo "Coverage: $COVERAGE%"
echo "Threshold: $THRESHOLD%"

if (( $(echo "$COVERAGE >= $THRESHOLD" | bc -l) )); then
    print_status 0 "Coverage meets threshold"
    echo ""
    echo -e "${YELLOW}📋 Coverage Improvement Plan:${NC}"
    echo "  Current: $COVERAGE%"
    echo "  Target:  90.0%"
    echo "  Next milestone: 50% (add Resource service + repository tests)"
    echo "  Future milestone: 70% (add HTTP handler tests)"
    echo "  Final milestone: 90% (add integration tests)"
else
    print_status 1 "Coverage below threshold"
fi
echo ""

echo "5. Running golangci-lint..."
echo "-------------------------------------------"
if command -v golangci-lint &> /dev/null; then
    if golangci-lint run --timeout=5m; then
        print_status 0 "golangci-lint passed"
    else
        print_status 1 "golangci-lint failed"
    fi
else
    echo -e "${YELLOW}⚠️  golangci-lint not installed, skipping${NC}"
    echo "Install: brew install golangci-lint (macOS)"
    echo "         go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"
fi
echo ""

echo "6. Building binary..."
echo "-------------------------------------------"
if go build -v -o bin/usermes-backend ./cmd/main.go; then
    print_status 0 "Build successful"
else
    print_status 1 "Build failed"
fi
echo ""

echo "=========================================="
if [ $OVERALL_STATUS -eq 0 ]; then
    echo -e "${GREEN}✅ ALL CI CHECKS PASSED!${NC}"
    echo "=========================================="
    echo ""
    echo "Your code is ready to push! 🚀"
    exit 0
else
    echo -e "${RED}❌ SOME CI CHECKS FAILED${NC}"
    echo "=========================================="
    echo ""
    echo "Please fix the issues above before pushing."
    exit 1
fi
