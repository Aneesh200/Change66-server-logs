-- Drop indexes
DROP INDEX IF EXISTS idx_server_metrics_timestamp;
DROP INDEX IF EXISTS idx_server_metrics_name_timestamp;
DROP INDEX IF EXISTS idx_api_keys_active;
DROP INDEX IF EXISTS idx_api_keys_hash;
DROP INDEX IF EXISTS idx_analytics_logs_device_info_gin;
DROP INDEX IF EXISTS idx_analytics_logs_properties_gin;
DROP INDEX IF EXISTS idx_analytics_logs_priority;
DROP INDEX IF EXISTS idx_analytics_logs_created_at;
DROP INDEX IF EXISTS idx_analytics_logs_session_id;
DROP INDEX IF EXISTS idx_analytics_logs_user_id;
DROP INDEX IF EXISTS idx_analytics_logs_event_name;
DROP INDEX IF EXISTS idx_analytics_logs_event_type;
DROP INDEX IF EXISTS idx_analytics_logs_timestamp;

-- Drop tables
DROP TABLE IF EXISTS server_metrics;
DROP TABLE IF EXISTS api_keys;
DROP TABLE IF EXISTS analytics_logs;
