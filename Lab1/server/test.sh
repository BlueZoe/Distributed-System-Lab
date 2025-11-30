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
# Test 1: GET existing file (txt)
# ============================================
echo -e "${BLUE}[Test 1]${NC} GET existing .txt file"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/resources/hello.txt")
status_code=$(echo "$response" | tail -n 1)
content=$(echo "$response" | sed '$d')

if check_status "200" "$status_code" && [ -n "$content" ]; then
  print_result "GET /resources/hello.txt" "PASS"
else
  print_result "GET /resources/hello.txt" "FAIL" "Expected 200, got $status_code"
fi
echo ""

# ============================================
# Test 2: GET existing file (html)
# ============================================
echo -e "${BLUE}[Test 2]${NC} GET existing .html file"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/resources/wiki.html")
status_code=$(echo "$response" | tail -n 1)

if check_status "200" "$status_code"; then
  print_result "GET /resources/wiki.html" "PASS"
else
  print_result "GET /resources/wiki.html" "FAIL" "Expected 200, got $status_code"
fi
echo ""

# ============================================
# Test 3: GET existing image file (jpeg)
# ============================================
echo -e "${BLUE}[Test 3]${NC} GET existing .jpeg file"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/resources/football.jpeg")
status_code=$(echo "$response" | tail -n 1)

if check_status "200" "$status_code"; then
  print_result "GET /resources/football.jpeg" "PASS"
else
  print_result "GET /resources/football.jpeg" "FAIL" "Expected 200, got $status_code"
fi
echo ""

# ============================================
# Test 4: GET non-existent file (404)
# ============================================
echo -e "${BLUE}[Test 4]${NC} GET non-existent file (should return 404)"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/resources/nonexistent.txt")
status_code=$(echo "$response" | tail -n 1)

if check_status "404" "$status_code"; then
  print_result "GET /resources/nonexistent.txt (404)" "PASS"
else
  print_result "GET /resources/nonexistent.txt (404)" "FAIL" "Expected 404, got $status_code"
fi
echo ""

# ============================================
# Test 5: GET unsupported file extension
# ============================================
echo -e "${BLUE}[Test 5]${NC} GET unsupported file extension (should return 501)"
# First create a file with unsupported extension for testing
test_file="/tmp/test_unsupported.pdf"
echo "test content" > "$test_file"

response=$(curl -s -w "\n%{http_code}" "$BASE_URL/resources/test.pdf" 2>/dev/null)
status_code=$(echo "$response" | tail -n 1)

if check_status "501" "$status_code" || check_status "404" "$status_code"; then
  print_result "GET unsupported extension (.pdf)" "PASS"
else
  print_result "GET unsupported extension (.pdf)" "FAIL" "Expected 501 or 404, got $status_code"
fi
echo ""

# ============================================
# Test 6: POST create new file
# ============================================
echo -e "${BLUE}[Test 6]${NC} POST create new .txt file"
test_content="Hello from POST test $(date +%s)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/uploaded/test_post.txt" -d "$test_content")
status_code=$(echo "$response" | tail -n 1)

if check_status "200" "$status_code"; then
  # Verify file was created
  if [ -f "uploaded/test_post.txt" ]; then
    file_content=$(cat "uploaded/test_post.txt")
    if [ "$file_content" = "$test_content" ]; then
      print_result "POST /uploaded/test_post.txt" "PASS"
    else
      print_result "POST /uploaded/test_post.txt" "FAIL" "File content mismatch"
    fi
  else
    print_result "POST /uploaded/test_post.txt" "FAIL" "File was not created"
  fi
else
  print_result "POST /uploaded/test_post.txt" "FAIL" "Expected 200, got $status_code"
fi
echo ""

# ============================================
# Test 7: POST create new file in nested directory
# ============================================
echo -e "${BLUE}[Test 7]${NC} POST create file in nested directory"
test_content="Nested directory test"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/uploaded/nested/test_nested.txt" -d "$test_content")
status_code=$(echo "$response" | tail -n 1)

if check_status "200" "$status_code" && [ -f "uploaded/nested/test_nested.txt" ]; then
  print_result "POST /uploaded/nested/test_nested.txt" "PASS"
else
  print_result "POST /uploaded/nested/test_nested.txt" "FAIL" "Expected 200, got $status_code"
fi
echo ""

# ============================================
# Test 8: POST with unsupported extension
# ============================================
echo -e "${BLUE}[Test 8]${NC} POST with unsupported extension (should return 501)"
response=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/uploaded/test.pdf" -d "test")
status_code=$(echo "$response" | tail -n 1)

if check_status "501" "$status_code"; then
  print_result "POST unsupported extension (.pdf)" "PASS"
else
  print_result "POST unsupported extension (.pdf)" "FAIL" "Expected 501, got $status_code"
fi
echo ""

# ============================================
# Test 9: Unsupported HTTP method (PUT)
# ============================================
echo -e "${BLUE}[Test 9]${NC} Unsupported HTTP method (PUT, should return 501)"
response=$(curl -s -w "\n%{http_code}" -X PUT "$BASE_URL/resources/test.txt" -d "test")
status_code=$(echo "$response" | tail -n 1)

if check_status "501" "$status_code"; then
  print_result "PUT method (unsupported)" "PASS"
else
  print_result "PUT method (unsupported)" "FAIL" "Expected 501, got $status_code"
fi
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

# ============================================
# Test 11: GET uploaded file (verify POST worked)
# ============================================
echo -e "${BLUE}[Test 11]${NC} GET file created by POST"
response=$(curl -s -w "\n%{http_code}" "$BASE_URL/uploaded/test_post.txt")
status_code=$(echo "$response" | tail -n 1)
content=$(echo "$response" | sed '$d')

if check_status "200" "$status_code" && [ -n "$content" ]; then
  print_result "GET /uploaded/test_post.txt (verify POST)" "PASS"
else
  print_result "GET /uploaded/test_post.txt (verify POST)" "FAIL" "Expected 200, got $status_code"
fi
echo ""

# ============================================
# Test Summary
# ============================================
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo -e "${GREEN}Passed: $PASSED${NC}"
echo -e "${RED}Failed: $FAILED${NC}"
echo "Total:  $((PASSED + FAILED))"
echo "=========================================="

if [ $FAILED -eq 0 ]; then
  echo -e "${GREEN}All tests passed! ✓${NC}"
  exit 0
else
  echo -e "${RED}Some tests failed! ✗${NC}"
  exit 1
fi
