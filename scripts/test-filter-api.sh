#!/bin/bash

# Test script for the new filtering API endpoint
# Usage: ./test-filter-api.sh

BASE_URL="http://localhost:8080"
API_KEY="habit-tracker-key-dev"

echo "=== Testing Log Filtering API ==="
echo

# Test 1: Filter by event_type
echo "1. Filter by event_type=behavioral:"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?event_type=behavioral&limit=2" | jq .
echo
echo "---"

# Test 2: Filter by event_name
echo "2. Filter by event_name=habit_fetched:"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?event_name=habit_fetched&limit=2" | jq .
echo
echo "---"

# Test 3: Filter by user_id
echo "3. Filter by user_id:"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?user_id=5bfLjXAYIAQkw1nTr4ScO1xmn5o1&limit=2" | jq .
echo
echo "---"

# Test 4: Filter by priority
echo "4. Filter by priority=normal:"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?priority=normal&limit=2" | jq .
echo
echo "---"

# Test 5: Filter by app_version
echo "5. Filter by app_version:"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?app_version=2.0.0+6&limit=2" | jq .
echo
echo "---"

# Test 6: Filter by provider_name
echo "6. Filter by provider_name=nfc_notifier:"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?provider_name=nfc_notifier&limit=2" | jq .
echo
echo "---"

# Test 6b: Combined filters
echo "6b. Combined filters (event_type + user_id):"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?event_type=behavioral&user_id=5bfLjXAYIAQkw1nTr4ScO1xmn5o1&limit=2" | jq .
echo
echo "---"

# Test 7: Pagination
echo "7. Pagination (page 1, page_size 1):"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?page=1&page_size=1" | jq .
echo
echo "---"

# Test 8: Sorting
echo "8. Sort by event_name ASC:"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?sort_by=event_name&sort_order=ASC&limit=3" | jq .
echo
echo "---"

# Test 9: Time range filter (last 24 hours)
echo "9. Time range filter (last 24 hours):"
START_TIME=$(date -u -v-1d +"%Y-%m-%dT%H:%M:%SZ" 2>/dev/null || date -u -d "1 day ago" +"%Y-%m-%dT%H:%M:%SZ")
END_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?start_time=$START_TIME&end_time=$END_TIME&limit=2" | jq .
echo
echo "---"

# Test 10: Invalid event_type (should return error)
echo "10. Invalid event_type (should return error):"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?event_type=invalid_type" | jq .
echo
echo "---"

# Test 11: Invalid time format (should return error)
echo "11. Invalid time format (should return error):"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?start_time=invalid-time" | jq .
echo
echo "---"

# Test 12: No filters (get all logs with pagination)
echo "12. No filters (all logs, first page):"
curl -s -H "X-API-Key: $API_KEY" \
  "$BASE_URL/api/v1/logs/filter?page=1&page_size=5" | jq .
echo

echo "=== Filter API Testing Complete ==="
