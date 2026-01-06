#!/bin/bash
# Check if coverage meets minimum threshold

THRESHOLD=${1:-90}
COVERAGE_FILE=${2:-coverage.out}

if [ ! -f "$COVERAGE_FILE" ]; then
    echo "❌ Coverage file not found: $COVERAGE_FILE"
    exit 1
fi

# Extract total coverage percentage
COVERAGE=$(go tool cover -func="$COVERAGE_FILE" | grep total | awk '{print $3}' | sed 's/%//')

echo "Current coverage: ${COVERAGE}%"
echo "Required coverage: ${THRESHOLD}%"

# Compare coverage with threshold using awk for portability
BELOW_THRESHOLD=$(awk -v cov="$COVERAGE" -v thr="$THRESHOLD" 'BEGIN {print (cov < thr) ? 1 : 0}')

if [ "$BELOW_THRESHOLD" -eq 1 ]; then
    echo "❌ Coverage ${COVERAGE}% is below threshold ${THRESHOLD}%"
    exit 1
fi

echo "✓ Coverage check passed!"
exit 0
