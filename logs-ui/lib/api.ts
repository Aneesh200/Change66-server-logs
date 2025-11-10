// API client for the logs server

import { AnalyticsLog, FilteredLogsResponse, SuccessResponse, LogFilter } from './types';

// Get API URL and API key from environment variables or use defaults
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://logs.biopeak.authify.tech';
const API_KEY = process.env.NEXT_PUBLIC_API_KEY || 'habit-tracker-key-dev';

export class LogsAPIClient {
  private baseUrl: string;
  private apiKey?: string;

  constructor(baseUrl: string = API_BASE_URL, apiKey?: string) {
    this.baseUrl = baseUrl;
    this.apiKey = apiKey;
  }

  private getHeaders(): HeadersInit {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
    };

    if (this.apiKey) {
      headers['X-API-Key'] = this.apiKey;
    }

    return headers;
  }

  /**
   * Fetch filtered logs with pagination
   */
  async getFilteredLogs(filter: LogFilter = {}): Promise<FilteredLogsResponse> {
    const params = new URLSearchParams();

    Object.entries(filter).forEach(([key, value]) => {
      if (value !== undefined && value !== null && value !== '') {
        params.append(key, value.toString());
      }
    });

    const url = `${this.baseUrl}/api/v1/logs/filter?${params.toString()}`;
    
    const response = await fetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch logs: ${response.statusText}`);
    }

    const result: SuccessResponse<FilteredLogsResponse> = await response.json();
    return result.data || { logs: [], total_count: 0, page: 1, page_size: 50, total_pages: 0 };
  }

  /**
   * Fetch recent logs (simple endpoint)
   */
  async getRecentLogs(limit: number = 50): Promise<AnalyticsLog[]> {
    const url = `${this.baseUrl}/api/v1/logs/recent?limit=${limit}`;
    
    const response = await fetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch recent logs: ${response.statusText}`);
    }

    const result: SuccessResponse<AnalyticsLog[]> = await response.json();
    return result.data || [];
  }

  /**
   * Check API health
   */
  async healthCheck(): Promise<{ status: string }> {
    const url = `${this.baseUrl}/health`;
    
    const response = await fetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Health check failed: ${response.statusText}`);
    }

    return response.json();
  }

  /**
   * Get analytics metrics
   */
  async getMetrics(): Promise<any> {
    const url = `${this.baseUrl}/api/v1/metrics`;
    
    const response = await fetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch metrics: ${response.statusText}`);
    }

    const result: SuccessResponse<any> = await response.json();
    return result.data || {};
  }

  /**
   * Get server status
   */
  async getStatus(): Promise<any> {
    const url = `${this.baseUrl}/api/v1/status`;
    
    const response = await fetch(url, {
      method: 'GET',
      headers: this.getHeaders(),
    });

    if (!response.ok) {
      throw new Error(`Failed to fetch status: ${response.statusText}`);
    }

    const result: SuccessResponse<any> = await response.json();
    return result.data || {};
  }
}

// Export a default instance with API key
export const logsAPI = new LogsAPIClient(API_BASE_URL, API_KEY);

