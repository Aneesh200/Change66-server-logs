-- Create logs table for storing analytics events
CREATE TABLE IF NOT EXISTS analytics_logs (
    id BIGSERIAL PRIMARY KEY,
    event_id VARCHAR(255) UNIQUE NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    event_type VARCHAR(50) NOT NULL,
    event_name VARCHAR(100) NOT NULL,
    properties JSONB,
    user_id VARCHAR(255),
    session_id VARCHAR(255),
    app_version VARCHAR(50),
    device_info JSONB,
    sequence_number INTEGER,
    priority VARCHAR(20) DEFAULT 'normal',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    processed_at TIMESTAMPTZ
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_analytics_logs_timestamp ON analytics_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_analytics_logs_event_type ON analytics_logs(event_type);
CREATE INDEX IF NOT EXISTS idx_analytics_logs_event_name ON analytics_logs(event_name);
CREATE INDEX IF NOT EXISTS idx_analytics_logs_user_id ON analytics_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_analytics_logs_session_id ON analytics_logs(session_id);
CREATE INDEX IF NOT EXISTS idx_analytics_logs_created_at ON analytics_logs(created_at);
CREATE INDEX IF NOT EXISTS idx_analytics_logs_priority ON analytics_logs(priority);

-- Create GIN index for JSONB properties for efficient querying
CREATE INDEX IF NOT EXISTS idx_analytics_logs_properties_gin ON analytics_logs USING GIN(properties);
CREATE INDEX IF NOT EXISTS idx_analytics_logs_device_info_gin ON analytics_logs USING GIN(device_info);

-- Create table for API key management
CREATE TABLE IF NOT EXISTS api_keys (
    id BIGSERIAL PRIMARY KEY,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_used_at TIMESTAMPTZ,
    expires_at TIMESTAMPTZ,
    usage_count BIGINT DEFAULT 0
);

-- Create index for API key lookups
CREATE INDEX IF NOT EXISTS idx_api_keys_hash ON api_keys(key_hash);
CREATE INDEX IF NOT EXISTS idx_api_keys_active ON api_keys(is_active);

-- Create table for storing server metrics
CREATE TABLE IF NOT EXISTS server_metrics (
    id BIGSERIAL PRIMARY KEY,
    metric_name VARCHAR(100) NOT NULL,
    metric_value DOUBLE PRECISION NOT NULL,
    labels JSONB,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create index for metrics
CREATE INDEX IF NOT EXISTS idx_server_metrics_name_timestamp ON server_metrics(metric_name, timestamp);
CREATE INDEX IF NOT EXISTS idx_server_metrics_timestamp ON server_metrics(timestamp);
