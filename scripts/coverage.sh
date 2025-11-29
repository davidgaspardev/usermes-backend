#!/bin/bash

# Coverage threshold
THRESHOLD=90.0

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "=========================================="
echo "Running tests with coverage..."
echo "=========================================="

# Run tests with coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...

if [ $? -ne 0 ]; then
    echo -e "${RED}❌ Tests failed${NC}"
    exit 1
fi

echo ""
echo "=========================================="
echo "Calculating coverage..."
echo "=========================================="

# Calculate total coverage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

echo ""
echo -e "Total coverage: ${YELLOW}${COVERAGE}%${NC}"
echo -e "Threshold:      ${YELLOW}${THRESHOLD}%${NC}"
echo ""

# Check if coverage meets threshold
if (( $(echo "$COVERAGE < $THRESHOLD" | bc -l) )); then
    echo -e "${RED}❌ Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%${NC}"
    echo ""
    echo "Coverage by package:"
    go tool cover -func=coverage.out | grep -v "total" | awk '{if ($3+0 < '${THRESHOLD}') print $1, $3}'
    exit 1
else
    echo -e "${GREEN}✅ Coverage ${COVERAGE}% meets threshold ${THRESHOLD}%${NC}"
fi

echo ""
echo "=========================================="
echo "Generating HTML coverage report..."
echo "=========================================="

# Generate HTML report
go tool cover -html=coverage.out -o coverage.html

echo -e "${GREEN}✓${NC} HTML report generated: coverage.html"
echo ""
echo "To view the report, open: file://$(pwd)/coverage.html"
echo ""

exit 0
