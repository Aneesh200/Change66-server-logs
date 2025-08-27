# Log Filtering API

This document describes the advanced log filtering endpoint that allows you to search and filter logs based on various criteria.

## Endpoint

```
GET /api/v1/logs/filter
```

## Authentication

Requires `X-API-Key` header with a valid API key.

## Query Parameters

### Filtering Parameters

| Parameter | Type | Description | Example |
|-----------|------|-------------|---------|
| `event_type` | string | Filter by event type | `behavioral`, `telemetry`, `observability`, `error`, `performance` |
| `event_name` | string | Filter by specific event name | `habit_fetched`, `user_login`, `page_view` |
| `user_id` | string | Filter by user ID | `5bfLjXAYIAQkw1nTr4ScO1xmn5o1` |
| `session_id` | string | Filter by session ID | `1756288356532_5bfLjXAYIAQkw1nTr4ScO1xmn5o1` |
| `app_version` | string | Filter by application version | `2.0.0+6`, `1.5.2` |
| `priority` | string | Filter by log priority | `normal`, `high` |
| `provider_name` | string | Filter by provider name | `nfc_notifier`, `deep_link_notifier` |
| `start_time` | string | Filter logs after this time (RFC3339 format) | `2025-01-01T00:00:00Z` |
| `end_time` | string | Filter logs before this time (RFC3339 format) | `2025-12-31T23:59:59Z` |

### Pagination Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `page` | integer | `1` | Page number (1-based) |
| `page_size` | integer | `50` | Number of logs per page (max 1000) |

### Sorting Parameters

| Parameter | Type | Default | Description |
|-----------|------|---------|-------------|
| `sort_by` | string | `created_at` | Field to sort by |
| `sort_order` | string | `DESC` | Sort order (`ASC` or `DESC`) |

#### Valid Sort Fields
- `id`, `event_id`, `timestamp`, `event_type`, `event_name`
- `user_id`, `session_id`, `app_version`, `priority`, `created_at`

## Response Format

```json
{
  "success": true,
  "message": "Retrieved X filtered logs (page Y of Z)",
  "data": {
    "logs": [...],
    "total_count": 1234,
    "page": 1,
    "page_size": 50,
    "total_pages": 25
  }
}
```

## Example Requests

### 1. Filter by Event Type

```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/api/v1/logs/filter?event_type=behavioral&limit=10"
```

### 2. Filter by User and Time Range

```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/api/v1/logs/filter?user_id=5bfLjXAYIAQkw1nTr4ScO1xmn5o1&start_time=2025-01-01T00:00:00Z&end_time=2025-01-31T23:59:59Z"
```

### 3. Combined Filters with Pagination

```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/api/v1/logs/filter?event_type=behavioral&priority=high&page=2&page_size=25"
```

### 4. Sort by Event Name (Ascending)

```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/api/v1/logs/filter?sort_by=event_name&sort_order=ASC&page_size=20"
```

### 5. Filter by Multiple Criteria

```bash
curl -H "X-API-Key: your-api-key" \
  "http://localhost:8080/api/v1/logs/filter?event_type=behavioral&event_name=habit_fetched&app_version=2.0.0+6&page_size=10"
```

## Error Responses

### Invalid Event Type
```json
{
  "error": "invalid_event_type",
  "message": "event_type must be one of: behavioral, telemetry, observability, error, performance"
}
```

### Invalid Time Format
```json
{
  "error": "invalid_start_time",
  "message": "start_time must be in RFC3339 format (e.g., 2023-01-01T00:00:00Z)"
}
```

### Invalid Priority
```json
{
  "error": "invalid_priority",
  "message": "priority must be either 'normal' or 'high'"
}
```

## Performance Notes

- Use pagination for large result sets
- Time range filters are highly recommended for better performance
- Combine multiple filters to narrow down results
- The endpoint supports up to 1000 logs per page
- Results are cached for optimal performance

## Use Cases

1. **Debugging**: Find specific error events for a user
2. **Analytics**: Get behavioral events for analysis
3. **Monitoring**: Filter by priority to see high-priority logs
4. **User Journey**: Track all events for a specific session
5. **Performance Analysis**: Get telemetry events within a time range
6. **Version Comparison**: Compare logs across different app versions

## Testing

Use the provided test script to verify the filtering functionality:

```bash
./scripts/test-filter-api.sh
```

This script tests various filter combinations and edge cases.
