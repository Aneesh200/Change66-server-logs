# Log Ingestion Server

A high-performance, multi-threaded Go server for ingesting analytics logs from the Habit Tracker Flutter app. Built with robust authentication, rate limiting, and comprehensive monitoring.

## Features

- **High Performance**: Multi-threaded architecture with configurable worker pools
- **Secure Authentication**: API key-based authentication with automatic key rotation
- **Rate Limiting**: Configurable rate limiting to prevent abuse
- **Batch Processing**: Support for both single and batch log ingestion
- **Advanced Filtering**: Comprehensive log filtering API with pagination and sorting
- **Grafana Visualization**: Pre-built dashboards for analytics and server metrics
- **Comprehensive Monitoring**: Health checks, metrics, and Prometheus integration
- **Database Integration**: PostgreSQL with connection pooling and migrations
- **Graceful Shutdown**: Proper cleanup and graceful server shutdown
- **Docker Support**: Complete Docker and Docker Compose setup
- **Validation**: Comprehensive request validation and error handling

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Docker & Docker Compose (optional)

### Development Setup

1. **Clone and setup**:
```bash
git clone <repository>
cd Change66-Log-Server
make dev-setup
```

2. **Configure environment**:
```bash
# Edit .env file with your configuration
cp config.example.env .env
```

3. **Start PostgreSQL** (if not using Docker):
```bash
# Install and start PostgreSQL
createdb analytics_logs
```

4. **Run migrations**:
```bash
make migrate-up
```

5. **Start the server**:
```bash
make run
```

### Docker Setup with Grafana (Recommended)

```bash
# Start all services including Grafana visualization
./scripts/start-with-grafana.sh

# Or manually with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f log-server

# Stop services
docker-compose down
```

**Services will be available at**:
- **Server**: http://localhost:8080
- **Grafana Dashboards**: http://localhost:3000 (admin/admin)
- **Prometheus**: http://localhost:9090
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

**Pre-built Grafana Dashboards**:
- **Analytics Overview**: http://localhost:3000/d/analytics-overview
- **Server Metrics**: http://localhost:3000/d/server-metrics

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `DB_HOST` | Database host | `localhost` |
| `DB_PORT` | Database port | `5432` |
| `DB_NAME` | Database name | `analytics_logs` |
| `DB_USER` | Database user | `postgres` |
| `DB_PASSWORD` | Database password | **required** |
| `API_KEYS` | Comma-separated API keys | **required** |
| `RATE_LIMIT_REQUESTS_PER_MINUTE` | Rate limit | `1000` |
| `MAX_BATCH_SIZE` | Maximum batch size | `1000` |
| `WORKER_POOL_SIZE` | Worker pool size | `10` |

See `config.example.env` for all available options.

## API Endpoints

### Authentication
All API endpoints require authentication via API key in the header:
```
Authorization: Bearer your-api-key
# OR
X-API-Key: your-api-key
```

### Log Ingestion

#### Single Log Ingestion
```http
POST /api/v1/ingest
Content-Type: application/json
X-API-Key: your-api-key

{
  "event_id": "unique-event-id",
  "timestamp": "2024-01-15T10:30:00Z",
  "event_type": "behavioral",
  "event_name": "habit_completed",
  "properties": {
    "habit_id": "123",
    "streak": 5
  },
  "user_id": "user-123",
  "session_id": "session-456",
  "app_version": "1.0.0",
  "device_info": {
    "platform": "ios",
    "model": "iPhone 14"
  },
  "sequence_number": 1,
  "priority": "normal"
}
```

#### Batch Log Ingestion
```http
POST /api/v1/batch-ingest
Content-Type: application/json
X-API-Key: your-api-key

{
  "logs": [
    { /* log object 1 */ },
    { /* log object 2 */ },
    // ... up to 1000 logs
  ]
}
```

### Monitoring

#### Health Check
```http
GET /health
```

#### Service Status
```http
GET /api/v1/status
X-API-Key: your-api-key
```

#### Analytics Metrics
```http
GET /api/v1/metrics
X-API-Key: your-api-key
```

#### Recent Logs (Debug)
```http
GET /api/v1/logs/recent?limit=50
X-API-Key: your-api-key
```

#### Advanced Log Filtering
```http
GET /api/v1/logs/filter?event_type=behavioral&user_id=USER_ID&page=1&page_size=50
X-API-Key: your-api-key
```

**Filter Parameters**:
- `event_type`: behavioral, telemetry, observability, error, performance
- `event_name`: Specific event name
- `user_id`: Filter by user ID
- `session_id`: Filter by session ID
- `app_version`: Filter by app version
- `priority`: normal, high
- `provider_name`: Filter by provider (e.g., nfc_notifier, deep_link_notifier)
- `start_time`, `end_time`: Time range (RFC3339 format)
- `page`, `page_size`: Pagination
- `sort_by`, `sort_order`: Sorting options

See [API_FILTERING.md](API_FILTERING.md) for detailed documentation.

#### Prometheus Metrics
```http
GET /metrics
```

## Event Types

The server supports the following event types:

- **behavioral**: User interactions, habits completed, etc.
- **telemetry**: Service performance, API calls, etc.
- **observability**: Provider state changes, system events
- **error**: Error tracking with context and stack traces
- **performance**: Performance metrics and measurements

## Database Schema

### analytics_logs
- `id`: Primary key
- `event_id`: Unique event identifier
- `timestamp`: Event timestamp
- `event_type`: Type of event (behavioral, telemetry, etc.)
- `event_name`: Specific event name
- `properties`: JSON properties
- `user_id`: User identifier
- `session_id`: Session identifier
- `app_version`: App version
- `device_info`: Device information (JSON)
- `sequence_number`: Event sequence number
- `priority`: Event priority (normal, high)
- `created_at`: Record creation time
- `processed_at`: Processing timestamp

## Performance & Scaling

### Multi-threading
- Configurable worker pool for handling requests
- Database connection pooling
- Asynchronous API key usage updates

### Rate Limiting
- Per-client rate limiting
- Configurable burst allowance
- Automatic backoff for exceeded limits

### Monitoring
- Prometheus metrics for all operations
- Health checks for dependencies
- Structured logging with request tracing

## Development

### Building
```bash
make build
```

### Testing
```bash
make test
make test-coverage
```

### Linting & Formatting
```bash
make lint
make format
```

### Database Migrations
```bash
# Create new migration
make migrate-create name=add_new_field

# Run migrations
make migrate-up

# Rollback migrations
make migrate-down
```

## Deployment

### Production Build
```bash
make build-prod
```

### Docker Production
```bash
docker build -t log-ingestion-server:latest .
docker run -p 8080:8080 --env-file .env log-ingestion-server:latest
```

### Kubernetes
Kubernetes manifests are available in the `k8s/` directory (to be created).

## Monitoring & Observability

### Prometheus Metrics
- HTTP request metrics (duration, status codes)
- Log ingestion metrics (count by type, batch sizes)
- Database operation metrics
- Rate limiting metrics
- Error tracking

### Grafana Visualization

**Access Grafana**: http://localhost:3000 (admin/admin)

**Pre-built Dashboards**:
- **Analytics Overview** (`/d/analytics-overview`):
  - Total logs and recent activity
  - Active users and sessions
  - Event type distribution
  - Log volume trends
  - Top events and app versions

- **Server Metrics** (`/d/server-metrics`):
  - Request rate and latency
  - Error rate monitoring
  - Log ingestion performance
  - Validation errors tracking
  - Batch size analytics

**Data Sources**:
- **Prometheus**: Server metrics and performance data
- **PostgreSQL**: Raw log data for detailed analytics

See [GRAFANA_SETUP.md](GRAFANA_SETUP.md) for detailed setup and customization.

### Health Checks
- `/health`: Overall service health
- `/readiness`: Kubernetes readiness probe
- `/liveness`: Kubernetes liveness probe

## Security

### API Key Management
- SHA-256 hashed storage
- Automatic key rotation support
- Usage tracking and analytics
- Configurable expiration

### Request Security
- Request size limits
- Rate limiting
- Input validation
- SQL injection protection
- XSS protection headers

## Troubleshooting

### Common Issues

1. **Database Connection Issues**:
   - Check database credentials in `.env`
   - Ensure PostgreSQL is running
   - Verify network connectivity

2. **Authentication Failures**:
   - Verify API key format
   - Check key is properly configured
   - Ensure key hasn't expired

3. **Rate Limiting**:
   - Check rate limit configuration
   - Monitor request patterns
   - Adjust limits if needed

### Logs
```bash
# View application logs
docker-compose logs -f log-server

# View database logs
docker-compose logs -f postgres
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run `make lint format test`
5. Submit a pull request

## License

[License information]

## Support

For issues and questions:
- Create an issue in the repository
- Check the troubleshooting section
- Review application logs
# Change66-server-logs
