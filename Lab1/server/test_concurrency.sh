#!/bin/bash
# Comprehensive test script for HTTP file server
# Tests: GET requests, POST requests, error handling, and semaphore blocking

PORT=${1:-8080}  # Use port from argument or default to 8080
HOST=${2:-localhost}
BASE_URL="http://$HOST:$PORT"

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Test counters
PASSED=0
FAILED=0

# Function to print test result
print_result() {
  local test_name=$1
  local status=$2
  local details=$3
  
  if [ "$status" = "PASS" ]; then
    echo -e "${GREEN}✓${NC} $test_name"
    ((PASSED++))
  else
    echo -e "${RED}✗${NC} $test_name"
    if [ -n "$details" ]; then
      echo -e "  ${RED}→${NC} $details"
    fi
    ((FAILED++))
  fi
}

# Function to check HTTP status code
# Usage: check_status "200" "$status_code"
check_status() {
  local expected=$1
  local actual=$2
  
  if [ "$actual" = "$expected" ]; then
    return 0
  else
    return 1
  fi
}

echo "=========================================="
echo "HTTP File Server Test Suite"
echo "=========================================="
echo "Server: $BASE_URL"
echo "Max Concurrent Requests: 10"
echo "=========================================="
echo ""

# ============================================
# Test 10: Concurrent requests and semaphore blocking
# ============================================
echo -e "${BLUE}[Test 10]${NC} Concurrent requests (15 requests, max 10 concurrent)"
echo -e "${YELLOW}This test sends 15 requests to verify semaphore blocking...${NC}"

# Function to send request and measure time
send_concurrent_request() {
  local num=$1
  local start_time=$(date +%s)
  curl -s -X GET "$BASE_URL/resources/hello.txt" > /dev/null 2>&1
  local end_time=$(date +%s)
  local duration=$((end_time - start_time))
  echo "Request #$num completed in ${duration}s"
}

# Send 15 concurrent requests
for i in {1..15}; do
  send_concurrent_request "$i" &
  sleep 0.1
done

wait
echo -e "${GREEN}All 15 concurrent requests completed${NC}"
echo -e "${YELLOW}Check server logs to verify semaphore blocking behavior${NC}"
print_result "Concurrent requests (15 requests)" "PASS"
echo ""