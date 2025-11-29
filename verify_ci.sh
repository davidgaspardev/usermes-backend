#!/bin/bash

echo "=========================================="
echo "CI/CD Verification Script"
echo "=========================================="
echo ""

# Check if CI workflow exists
if [ -f ".github/workflows/ci.yml" ]; then
    echo "✅ CI workflow file exists"
else
    echo "❌ CI workflow file missing"
    exit 1
fi

# Check if tests exist
TEST_FILES=$(find internal/modules/user -name "*_test.go" | wc -l)
if [ "$TEST_FILES" -gt 0 ]; then
    echo "✅ Test files found: $TEST_FILES"
else
    echo "❌ No test files found"
    exit 1
fi

# Run tests and check coverage
echo ""
echo "Running tests with coverage..."
go test -coverprofile=coverage.out ./internal/modules/user/domain/... ./internal/modules/user/application/... > /dev/null 2>&1

if [ $? -eq 0 ]; then
    echo "✅ All tests passed"
else
    echo "❌ Tests failed"
    exit 1
fi

# Check coverage
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
echo "📊 Coverage: ${COVERAGE}%"

if (( $(echo "$COVERAGE >= 90.0" | bc -l) )); then
    echo "✅ Coverage meets 90% threshold"
else
    echo "❌ Coverage below 90% threshold"
    exit 1
fi

echo ""
echo "=========================================="
echo "✅ CI/CD VERIFICATION PASSED"
echo "=========================================="
echo ""
echo "Summary:"
echo "  • CI Workflow: ✅ Ready"
echo "  • Test Suite: ✅ $TEST_FILES test files"
echo "  • Coverage: ✅ ${COVERAGE}% (≥90%)"
echo "  • Status: ✅ PRODUCTION READY"
echo ""
