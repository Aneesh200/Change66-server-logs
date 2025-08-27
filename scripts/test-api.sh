#!/bin/bash

# API Testing Script for Log Ingestion Server

set -e

# Configuration
SERVER_URL="http://localhost:8080"
API_KEY="habit-tracker-key-dev"  # Update this with your actual API key

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🧪 Testing Log Ingestion Server API${NC}"
echo -e "${BLUE}Server: $SERVER_URL${NC}"
echo ""

# Function to make HTTP request and check response
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local description=$5
    
    echo -e "${YELLOW}Testing: $description${NC}"
    echo -e "${BLUE}$method $endpoint${NC}"
    
    if [ -n "$data" ]; then
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "Content-Type: application/json" \
            -H "X-API-Key: $API_KEY" \
            -d "$data" \
            "$SERVER_URL$endpoint")
    else
        response=$(curl -s -w "\n%{http_code}" -X $method \
            -H "X-API-Key: $API_KEY" \
            "$SERVER_URL$endpoint")
    fi
    
    # Extract status code (last line) and body (everything else)
    status_code=$(echo "$response" | tail -n1)
    body=$(echo "$response" | head -n -1)
    
    if [ "$status_code" = "$expected_status" ]; then
        echo -e "${GREEN}✅ Success ($status_code)${NC}"
        echo "$body" | jq . 2>/dev/null || echo "$body"
    else
        echo -e "${RED}❌ Failed (expected $expected_status, got $status_code)${NC}"
        echo "$body"
    fi
    
    echo ""
}

# Test 1: Health Check (no auth required)
test_endpoint "GET" "/health" "" "200" "Health Check"

# Test 2: Readiness Check
test_endpoint "GET" "/readiness" "" "200" "Readiness Check"

# Test 3: Liveness Check
test_endpoint "GET" "/liveness" "" "200" "Liveness Check"

# Test 4: Single Log Ingestion
single_log='{
  "event_id": "test_'$(date +%s)'_1",
  "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
  "event_type": "behavioral",
  "event_name": "habit_completed",
  "properties": {
    "habit_id": "test_habit_123",
    "streak": 5,
    "category": "health"
  },
  "user_id": "test_user_456",
  "session_id": "test_session_789",
  "app_version": "1.0.0",
  "device_info": {
    "platform": "test",
    "model": "Test Device",
    "version": "1.0"
  },
  "sequence_number": 1,
  "priority": "normal"
}'

test_endpoint "POST" "/api/v1/ingest" "$single_log" "201" "Single Log Ingestion"

# Test 5: Batch Log Ingestion
batch_logs='{
  "logs": [
    {
      "event_id": "test_'$(date +%s)'_batch_1",
      "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
      "event_type": "telemetry",
      "event_name": "service_call",
      "properties": {
        "service_name": "habit_service",
        "operation": "create_habit",
        "duration_ms": 150
      },
      "user_id": "test_user_456",
      "session_id": "test_session_789",
      "app_version": "1.0.0",
      "priority": "normal"
    },
    {
      "event_id": "test_'$(date +%s)'_batch_2",
      "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
      "event_type": "performance",
      "event_name": "metric_recorded",
      "properties": {
        "metric_name": "app_startup_time",
        "value": 2.5,
        "unit": "seconds"
      },
      "user_id": "test_user_456",
      "session_id": "test_session_789",
      "app_version": "1.0.0",
      "priority": "normal"
    }
  ]
}'

test_endpoint "POST" "/api/v1/batch-ingest" "$batch_logs" "201" "Batch Log Ingestion"

# Test 6: Service Status
test_endpoint "GET" "/api/v1/status" "" "200" "Service Status"

# Test 7: Analytics Metrics
test_endpoint "GET" "/api/v1/metrics" "" "200" "Analytics Metrics"

# Test 8: Recent Logs
test_endpoint "GET" "/api/v1/logs/recent?limit=5" "" "200" "Recent Logs"

# Test 9: Prometheus Metrics
echo -e "${YELLOW}Testing: Prometheus Metrics${NC}"
echo -e "${BLUE}GET /metrics${NC}"
metrics_response=$(curl -s -w "%{http_code}" "$SERVER_URL/metrics")
if [[ "$metrics_response" == *"200" ]]; then
    echo -e "${GREEN}✅ Success (200)${NC}"
    echo "Metrics endpoint is working (output truncated)"
else
    echo -e "${RED}❌ Failed${NC}"
    echo "$metrics_response"
fi
echo ""

# Test 10: Invalid API Key
echo -e "${YELLOW}Testing: Invalid API Key${NC}"
echo -e "${BLUE}POST /api/v1/ingest (with invalid key)${NC}"
invalid_response=$(curl -s -w "%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "X-API-Key: invalid-key" \
    -d "$single_log" \
    "$SERVER_URL/api/v1/ingest")
if [[ "$invalid_response" == *"401" ]]; then
    echo -e "${GREEN}✅ Correctly rejected invalid API key (401)${NC}"
else
    echo -e "${RED}❌ Should have rejected invalid API key${NC}"
fi
echo ""

# Test 11: Missing API Key
echo -e "${YELLOW}Testing: Missing API Key${NC}"
echo -e "${BLUE}POST /api/v1/ingest (no key)${NC}"
no_key_response=$(curl -s -w "%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -d "$single_log" \
    "$SERVER_URL/api/v1/ingest")
if [[ "$no_key_response" == *"401" ]]; then
    echo -e "${GREEN}✅ Correctly rejected missing API key (401)${NC}"
else
    echo -e "${RED}❌ Should have rejected missing API key${NC}"
fi
echo ""

# Test 12: Invalid JSON
echo -e "${YELLOW}Testing: Invalid JSON${NC}"
echo -e "${BLUE}POST /api/v1/ingest (invalid JSON)${NC}"
invalid_json_response=$(curl -s -w "%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $API_KEY" \
    -d '{"invalid": json}' \
    "$SERVER_URL/api/v1/ingest")
if [[ "$invalid_json_response" == *"400" ]]; then
    echo -e "${GREEN}✅ Correctly rejected invalid JSON (400)${NC}"
else
    echo -e "${RED}❌ Should have rejected invalid JSON${NC}"
fi
echo ""

# Test 13: Large Batch (should be rejected)
echo -e "${YELLOW}Testing: Large Batch (>1000 items)${NC}"
echo -e "${BLUE}POST /api/v1/batch-ingest (large batch)${NC}"
# Create a batch with 1001 items (should be rejected)
large_batch='{"logs":['
for i in {1..1001}; do
    large_batch+='{
        "event_id": "test_large_'$i'",
        "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
        "event_type": "behavioral",
        "event_name": "test_event",
        "properties": {},
        "priority": "normal"
    }'
    if [ $i -lt 1001 ]; then
        large_batch+=','
    fi
done
large_batch+=']}'

large_batch_response=$(curl -s -w "%{http_code}" -X POST \
    -H "Content-Type: application/json" \
    -H "X-API-Key: $API_KEY" \
    -d "$large_batch" \
    "$SERVER_URL/api/v1/batch-ingest")
if [[ "$large_batch_response" == *"400" ]]; then
    echo -e "${GREEN}✅ Correctly rejected large batch (400)${NC}"
else
    echo -e "${RED}❌ Should have rejected large batch${NC}"
fi
echo ""

echo -e "${GREEN}🎉 API testing completed!${NC}"
echo ""
echo -e "${YELLOW}Summary:${NC}"
echo "- Test your specific API keys with your Flutter app"
echo "- Monitor logs for any errors: docker-compose logs -f log-server"
echo "- Check metrics at: $SERVER_URL/metrics"
echo "- View recent logs at: $SERVER_URL/api/v1/logs/recent"
