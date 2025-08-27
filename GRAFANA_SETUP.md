# Grafana Visualization Setup

This document explains how to set up and use Grafana for visualizing your analytics logs and server metrics.

## 🚀 Quick Start

### 1. Start the Stack

```bash
docker-compose up -d
```

### 2. Access Grafana

- **URL**: http://localhost:3000
- **Username**: `admin`
- **Password**: `admin`

### 3. View Pre-configured Dashboards

Two dashboards are automatically provisioned:

1. **Analytics Overview Dashboard** - Log analytics and user behavior insights
2. **Server Metrics Dashboard** - Server performance and operational metrics

## 📊 Available Dashboards

### Analytics Overview Dashboard

**URL**: http://localhost:3000/d/analytics-overview

**Features**:
- **Total Logs**: Overall count of analytics logs
- **Logs Last Hour**: Recent activity tracking
- **Active Users**: Unique users in the last hour
- **Active Sessions**: Current session count
- **Event Types Distribution**: Pie chart of event types (behavioral, telemetry, etc.)
- **Priority Distribution**: Normal vs high priority logs
- **Log Volume Timeline**: Hourly log volume by event type
- **Top Event Names**: Most frequent events
- **App Versions**: Version distribution analysis

**Use Cases**:
- Monitor user engagement and activity patterns
- Track feature usage and adoption
- Identify popular app versions
- Analyze user behavior trends

### Server Metrics Dashboard

**URL**: http://localhost:3000/d/server-metrics

**Features**:
- **Request Rate**: Requests per second
- **95th Percentile Latency**: Response time performance
- **Error Rate**: Percentage of 5xx errors
- **Logs Ingested/sec**: Log ingestion rate
- **Request Rate by Endpoint**: Traffic breakdown by API endpoint
- **Response Time Percentiles**: 50th, 95th, and 99th percentile latencies
- **Log Ingestion by Type**: Ingestion rates by event type and priority
- **Validation Errors**: Rate of validation failures
- **Batch Size Distribution**: Log batch size analytics

**Use Cases**:
- Monitor server performance and health
- Identify performance bottlenecks
- Track API usage patterns
- Monitor error rates and validation issues

## 🔧 Configuration

### Data Sources

Two data sources are automatically configured:

1. **Prometheus** (Default)
   - URL: `http://prometheus:9090`
   - Used for: Server metrics, performance data

2. **PostgreSQL**
   - Host: `postgres:5432`
   - Database: `analytics_logs`
   - Used for: Raw log data analysis

### Custom Queries

You can create custom panels using these data sources:

#### PostgreSQL Queries Examples

```sql
-- User activity by hour
SELECT 
  DATE_TRUNC('hour', created_at) as time,
  COUNT(DISTINCT user_id) as active_users
FROM analytics_logs 
WHERE created_at >= NOW() - INTERVAL '24 hours'
  AND user_id IS NOT NULL
GROUP BY DATE_TRUNC('hour', created_at)
ORDER BY time;

-- Event funnel analysis
SELECT 
  event_name,
  COUNT(*) as count,
  COUNT(DISTINCT user_id) as unique_users
FROM analytics_logs 
WHERE event_type = 'behavioral'
  AND created_at >= NOW() - INTERVAL '7 days'
GROUP BY event_name
ORDER BY count DESC;

-- Device platform distribution
SELECT 
  device_info->>'platform' as platform,
  COUNT(*) as count
FROM analytics_logs 
WHERE device_info->>'platform' IS NOT NULL
  AND created_at >= NOW() - INTERVAL '24 hours'
GROUP BY device_info->>'platform'
ORDER BY count DESC;
```

#### Prometheus Queries Examples

```promql
# Request rate by status code
sum(rate(http_requests_total[5m])) by (status)

# Average response time
rate(http_request_duration_seconds_sum[5m]) / rate(http_request_duration_seconds_count[5m])

# Log ingestion throughput
sum(rate(logs_ingested_total[5m])) by (event_type)

# Error rate percentage
rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) * 100
```

## 📈 Creating Custom Dashboards

### 1. Access Dashboard Creation

- Go to Grafana UI (http://localhost:3000)
- Click "+" → "Dashboard"
- Add panels with your custom queries

### 2. Panel Types for Different Data

**For Time Series Data**:
- Use "Time series" panel
- Good for: Log volume over time, user activity trends

**For Current Values**:
- Use "Stat" panel
- Good for: Total counts, current rates, percentages

**For Distributions**:
- Use "Pie chart" or "Bar chart" panels
- Good for: Event type distribution, device platforms

**For Detailed Data**:
- Use "Table" panel
- Good for: Top events, user lists, detailed breakdowns

### 3. Variables and Templating

Create dynamic dashboards with variables:

```
# Time range variable
$__timeFrom and $__timeTo

# Custom variables for filtering
$event_type, $user_id, $app_version
```

## 🎨 Dashboard Customization

### Themes and Appearance

- **Dark Theme**: Default, great for monitoring
- **Light Theme**: Available in preferences
- **Custom Colors**: Configure in panel settings

### Auto-Refresh

- Set refresh intervals: 5s, 10s, 30s, 1m, 5m
- Useful for real-time monitoring

### Time Ranges

- **Default**: Last 6 hours for analytics, last 1 hour for metrics
- **Custom**: Set specific time ranges for analysis
- **Relative**: Last N hours/days/weeks

## 🚨 Alerting

### Setting Up Alerts

1. **Create Alert Rules**:
   - Go to Alerting → Alert Rules
   - Define conditions (e.g., error rate > 5%)

2. **Notification Channels**:
   - Email, Slack, Discord, Webhook
   - Configure in Alerting → Notification channels

3. **Example Alert Rules**:
   ```
   # High error rate alert
   rate(http_requests_total{status=~"5.."}[5m]) / rate(http_requests_total[5m]) * 100 > 5
   
   # Low log ingestion alert
   rate(logs_ingested_total[5m]) < 0.1
   
   # High response time alert
   histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m])) > 2
   ```

## 🔍 Troubleshooting

### Common Issues

**Dashboard not loading data**:
- Check data source connectivity
- Verify database credentials in datasources.yml
- Ensure Prometheus is scraping metrics

**Grafana won't start**:
- Check docker logs: `docker logs log-ingestion-grafana`
- Verify volume mounts and permissions

**Missing metrics**:
- Ensure your application is exposing metrics on `/metrics`
- Check Prometheus targets: http://localhost:9090/targets

### Useful Commands

```bash
# Restart Grafana
docker-compose restart grafana

# View Grafana logs
docker logs -f log-ingestion-grafana

# Access Grafana container
docker exec -it log-ingestion-grafana /bin/bash

# Backup dashboards
docker cp log-ingestion-grafana:/var/lib/grafana ./grafana-backup
```

## 📚 Advanced Features

### 1. Data Source Plugins

Add more data sources:
```yaml
environment:
  - GF_INSTALL_PLUGINS=grafana-piechart-panel,redis-datasource
```

### 2. External Authentication

Configure LDAP, OAuth, or SAML:
```yaml
environment:
  - GF_AUTH_GOOGLE_ENABLED=true
  - GF_AUTH_GOOGLE_CLIENT_ID=your-client-id
```

### 3. High Availability

For production, consider:
- External database for Grafana config
- Load balancing multiple Grafana instances
- Persistent storage for dashboards

## 🎯 Best Practices

1. **Dashboard Organization**:
   - Group related metrics together
   - Use consistent naming conventions
   - Add descriptions to panels

2. **Performance**:
   - Limit time ranges for heavy queries
   - Use appropriate refresh intervals
   - Optimize PostgreSQL queries with indexes

3. **Security**:
   - Change default admin password
   - Set up proper user roles and permissions
   - Use HTTPS in production

4. **Monitoring**:
   - Set up alerts for critical metrics
   - Monitor Grafana itself
   - Regular backup of dashboards

## 📖 Resources

- [Grafana Documentation](https://grafana.com/docs/)
- [Prometheus Query Language](https://prometheus.io/docs/prometheus/latest/querying/)
- [PostgreSQL Functions](https://www.postgresql.org/docs/current/functions.html)
- [Dashboard Examples](https://grafana.com/grafana/dashboards/)

---

Your Grafana setup is now ready to provide powerful insights into your analytics logs and server performance! 🎉
