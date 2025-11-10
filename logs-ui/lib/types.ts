// TypeScript types for the logs API

export interface AnalyticsLog {
  id: number;
  event_id: string;
  timestamp: string;
  event_type: 'behavioral' | 'telemetry' | 'observability' | 'error' | 'performance';
  event_name: string;
  properties: Record<string, any>;
  user_id?: string | null;
  session_id?: string | null;
  app_version?: string | null;
  device_info: Record<string, any>;
  sequence_number?: number | null;
  priority: 'normal' | 'high';
  created_at: string;
  processed_at?: string | null;
}

export interface FilteredLogsResponse {
  logs: AnalyticsLog[];
  total_count: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface SuccessResponse<T = any> {
  success: boolean;
  message: string;
  data?: T;
}

export interface ErrorResponse {
  error: string;
  message: string;
  details?: Array<{
    field: string;
    message: string;
  }>;
}

export interface LogFilter {
  event_type?: string;
  event_name?: string;
  user_id?: string;
  session_id?: string;
  app_version?: string;
  priority?: string;
  provider_name?: string;
  start_time?: string;
  end_time?: string;
  page?: number;
  page_size?: number;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

