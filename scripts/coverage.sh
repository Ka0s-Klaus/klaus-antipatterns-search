#!/usr/bin/env bash

# Generate test coverage report and display statistics.
# Usage: ./scripts/coverage.sh [html|report]
#   html   — generate HTML coverage report (coverage.html)
#   report — print coverage statistics to stdout (default)

set -e

MODE="${1:-report}"
COVERAGE_FILE="coverage.out"

echo "🧪 Running tests with coverage..."
go test -v -race -coverprofile="$COVERAGE_FILE" ./...

# Parse coverage percentage
COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' | sed 's/%//')

echo ""
echo "📊 Coverage Report"
echo "=================="
echo "Total coverage: ${COVERAGE}%"
echo ""

if [ "$MODE" = "html" ]; then
	echo "Generating HTML report..."
	go tool cover -html="$COVERAGE_FILE" -o coverage.html
	echo "✅ Report generated: coverage.html"
elif [ "$MODE" = "report" ]; then
	# Print function coverage statistics
	go tool cover -func="$COVERAGE_FILE" | head -20
	echo "..."
	go tool cover -func="$COVERAGE_FILE" | tail -3
else
	echo "Unknown mode: $MODE"
	echo "Usage: $0 [html|report]"
	exit 1
fi

# Warn if coverage < 80%
if (( $(echo "$COVERAGE < 80" | bc -l) )); then
	echo ""
	echo "⚠️  Coverage is below 80% target: ${COVERAGE}%"
	exit 1
fi

echo ""
echo "✅ Coverage check passed (>= 80%)"
