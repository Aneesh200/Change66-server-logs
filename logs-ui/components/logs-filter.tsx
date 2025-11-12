'use client';

import { useState } from 'react';
import { Input } from '@/components/ui/input';
import { Button } from '@/components/ui/button';
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Search, X, Calendar, Filter, SortAsc } from 'lucide-react';
import { LogFilter } from '@/lib/types';

interface LogsFilterProps {
  onFilter: (filter: LogFilter) => void;
  isLoading?: boolean;
  hideEventTypeFilter?: boolean;
}

export function LogsFilter({ onFilter, isLoading, hideEventTypeFilter = false }: LogsFilterProps) {
  const [filter, setFilter] = useState<LogFilter>({
    event_type: 'all',
    event_name: '',
    user_id: '',
    session_id: '',
    priority: 'all',
    app_version: '',
    provider_name: '',
    start_time: '',
    end_time: '',
    sort_by: 'timestamp',
    sort_order: 'desc',
  });

  const [searchQuery, setSearchQuery] = useState('');
  const [showAdvanced, setShowAdvanced] = useState(false);

  const handleFilterChange = (key: keyof LogFilter, value: string) => {
    setFilter((prev) => ({ ...prev, [key]: value }));
  };

  const handleApplyFilter = () => {
    // Remove empty values and "all" values
    const cleanFilter: LogFilter = {};
    Object.entries(filter).forEach(([key, value]) => {
      if (value && value !== '' && value !== 'all') {
        cleanFilter[key as keyof LogFilter] = value;
      }
    });
    
    // Add search query if present
    if (searchQuery.trim()) {
      cleanFilter.search = searchQuery.trim();
    }
    
    onFilter(cleanFilter);
  };

  const handleClearFilter = () => {
    setFilter({
      event_type: 'all',
      event_name: '',
      user_id: '',
      session_id: '',
      priority: 'all',
      app_version: '',
      provider_name: '',
      start_time: '',
      end_time: '',
      sort_by: 'timestamp',
      sort_order: 'desc',
    });
    setSearchQuery('');
    onFilter({});
  };

  // Get date in YYYY-MM-DD format for input
  const getDateInputValue = (dateString: string) => {
    if (!dateString) return '';
    try {
      return new Date(dateString).toISOString().split('T')[0];
    } catch {
      return '';
    }
  };

  // Convert date input to ISO string
  const handleDateChange = (key: 'start_time' | 'end_time', value: string) => {
    if (!value) {
      setFilter((prev) => ({ ...prev, [key]: '' }));
      return;
    }
    
    // Create ISO string for the date
    const date = new Date(value);
    if (key === 'start_time') {
      date.setHours(0, 0, 0, 0);
    } else {
      date.setHours(23, 59, 59, 999);
    }
    setFilter((prev) => ({ ...prev, [key]: date.toISOString() }));
  };

  return (
    <Card>
      <CardHeader>
        <div className="flex items-center justify-between">
          <CardTitle className="text-lg flex items-center gap-2">
            <Filter className="h-5 w-5" />
            Filters & Search
          </CardTitle>
          <Button
            variant="ghost"
            size="sm"
            onClick={() => setShowAdvanced(!showAdvanced)}
          >
            {showAdvanced ? 'Hide' : 'Show'} Advanced
          </Button>
        </div>
      </CardHeader>
      <CardContent>
        {/* Search Bar */}
        <div className="mb-6">
          <label className="text-sm font-medium mb-2 block">Quick Search</label>
          <div className="relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <Input
              placeholder="Search logs, metadata, properties... (e.g., user_123, error message, habit_id)"
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              onKeyDown={(e) => {
                if (e.key === 'Enter') {
                  handleApplyFilter();
                }
              }}
              className="pl-10 h-11 text-base bg-muted/50 border-muted-foreground/20 focus:border-primary"
            />
            {searchQuery && (
              <button
                onClick={() => setSearchQuery('')}
                className="absolute right-3 top-1/2 -translate-y-1/2 text-muted-foreground hover:text-foreground"
              >
                <X className="h-4 w-4" />
              </button>
            )}
          </div>
          <p className="text-xs text-muted-foreground mt-2">
            💡 Searches through event names, user IDs, session IDs, properties, and metadata
          </p>
        </div>

        {/* Basic Filters */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {!hideEventTypeFilter && (
            <div className="space-y-2">
              <label className="text-sm font-medium">Event Type</label>
              <Select
                value={filter.event_type}
                onValueChange={(value) => handleFilterChange('event_type', value)}
              >
                <SelectTrigger>
                  <SelectValue placeholder="All types" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="all">All types</SelectItem>
                  <SelectItem value="behavioral">Behavioral</SelectItem>
                  <SelectItem value="telemetry">Telemetry</SelectItem>
                  <SelectItem value="observability">Observability</SelectItem>
                  <SelectItem value="error">Error</SelectItem>
                  <SelectItem value="performance">Performance</SelectItem>
                </SelectContent>
              </Select>
            </div>
          )}

          <div className="space-y-2">
            <label className="text-sm font-medium">Priority</label>
            <Select
              value={filter.priority}
              onValueChange={(value) => handleFilterChange('priority', value)}
            >
              <SelectTrigger>
                <SelectValue placeholder="All priorities" />
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="all">All priorities</SelectItem>
                <SelectItem value="normal">Normal</SelectItem>
                <SelectItem value="high">High</SelectItem>
              </SelectContent>
            </Select>
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Event Name</label>
            <Input
              placeholder="e.g., user_login, habit_created"
              value={filter.event_name}
              onChange={(e) => handleFilterChange('event_name', e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">User ID</label>
            <Input
              placeholder="Filter by user ID"
              value={filter.user_id}
              onChange={(e) => handleFilterChange('user_id', e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">Session ID</label>
            <Input
              placeholder="Filter by session ID"
              value={filter.session_id}
              onChange={(e) => handleFilterChange('session_id', e.target.value)}
            />
          </div>

          <div className="space-y-2">
            <label className="text-sm font-medium">App Version</label>
            <Input
              placeholder="e.g., 1.0.0, 2.1.3"
              value={filter.app_version}
              onChange={(e) => handleFilterChange('app_version', e.target.value)}
            />
          </div>
        </div>

        {/* Advanced Filters */}
        {showAdvanced && (
          <div className="mt-6 pt-6 border-t space-y-4">
            <h3 className="text-sm font-semibold flex items-center gap-2">
              <Calendar className="h-4 w-4" />
              Date Range & Advanced Options
            </h3>
            
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <div className="space-y-2">
                <label className="text-sm font-medium">Start Date</label>
                <Input
                  type="date"
                  value={getDateInputValue(filter.start_time || '')}
                  onChange={(e) => handleDateChange('start_time', e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">End Date</label>
                <Input
                  type="date"
                  value={getDateInputValue(filter.end_time || '')}
                  onChange={(e) => handleDateChange('end_time', e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Provider Name</label>
                <Input
                  placeholder="e.g., firebase, analytics"
                  value={filter.provider_name}
                  onChange={(e) => handleFilterChange('provider_name', e.target.value)}
                />
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium flex items-center gap-1">
                  <SortAsc className="h-3 w-3" />
                  Sort By
                </label>
                <Select
                  value={filter.sort_by}
                  onValueChange={(value) => handleFilterChange('sort_by', value)}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="timestamp">Timestamp</SelectItem>
                    <SelectItem value="event_type">Event Type</SelectItem>
                    <SelectItem value="priority">Priority</SelectItem>
                    <SelectItem value="event_name">Event Name</SelectItem>
                  </SelectContent>
                </Select>
              </div>

              <div className="space-y-2">
                <label className="text-sm font-medium">Sort Order</label>
                <Select
                  value={filter.sort_order}
                  onValueChange={(value) => handleFilterChange('sort_order', value as 'asc' | 'desc')}
                >
                  <SelectTrigger>
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="desc">Newest First</SelectItem>
                    <SelectItem value="asc">Oldest First</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>

            {/* Quick Date Presets */}
            <div className="flex flex-wrap gap-2">
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  const now = new Date();
                  const today = new Date(now.getFullYear(), now.getMonth(), now.getDate());
                  setFilter((prev) => ({
                    ...prev,
                    start_time: today.toISOString(),
                    end_time: new Date(today.getTime() + 86400000 - 1).toISOString(),
                  }));
                }}
              >
                Today
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  const now = new Date();
                  const yesterday = new Date(now.getTime() - 86400000);
                  const start = new Date(yesterday.getFullYear(), yesterday.getMonth(), yesterday.getDate());
                  setFilter((prev) => ({
                    ...prev,
                    start_time: start.toISOString(),
                    end_time: new Date(start.getTime() + 86400000 - 1).toISOString(),
                  }));
                }}
              >
                Yesterday
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  const now = new Date();
                  const lastWeek = new Date(now.getTime() - 7 * 86400000);
                  setFilter((prev) => ({
                    ...prev,
                    start_time: lastWeek.toISOString(),
                    end_time: now.toISOString(),
                  }));
                }}
              >
                Last 7 Days
              </Button>
              <Button
                variant="outline"
                size="sm"
                onClick={() => {
                  const now = new Date();
                  const lastMonth = new Date(now.getTime() - 30 * 86400000);
                  setFilter((prev) => ({
                    ...prev,
                    start_time: lastMonth.toISOString(),
                    end_time: now.toISOString(),
                  }));
                }}
              >
                Last 30 Days
              </Button>
            </div>
          </div>
        )}

        {/* Action Buttons */}
        <div className="flex gap-2 mt-6">
          <Button 
            onClick={handleApplyFilter} 
            disabled={isLoading}
            className="flex-1"
          >
            <Search className="h-4 w-4 mr-2" />
            Apply Filters
          </Button>
          <Button 
            onClick={handleClearFilter} 
            variant="outline"
            disabled={isLoading}
          >
            <X className="h-4 w-4 mr-2" />
            Clear All
          </Button>
        </div>

        {/* Active Filters Summary */}
        {(searchQuery || Object.entries(filter).some(([key, value]) => 
          value && value !== 'all' && value !== 'timestamp' && value !== 'desc'
        )) && (
          <div className="mt-4 pt-4 border-t">
            <p className="text-xs text-muted-foreground mb-2">Active Filters:</p>
            <div className="flex flex-wrap gap-2">
              {searchQuery && (
                <div className="inline-flex items-center gap-1 bg-blue-500/10 text-blue-400 px-2 py-1 rounded-md text-xs border border-blue-500/20">
                  <Search className="h-3 w-3" />
                  <span className="font-medium">Search:</span>
                  <span className="max-w-[200px] truncate">{searchQuery}</span>
                </div>
              )}
              {Object.entries(filter).map(([key, value]) => {
                if (!value || value === 'all' || value === 'timestamp' || value === 'desc') return null;
                
                let displayValue = value;
                if (key === 'start_time' || key === 'end_time') {
                  try {
                    displayValue = new Date(value as string).toLocaleDateString();
                  } catch {}
                }
                
                return (
                  <div
                    key={key}
                    className="inline-flex items-center gap-1 bg-primary/10 text-primary px-2 py-1 rounded-md text-xs"
                  >
                    <span className="font-medium">{key.replace('_', ' ')}:</span>
                    <span>{displayValue}</span>
                  </div>
                );
              })}
            </div>
          </div>
        )}
      </CardContent>
    </Card>
  );
}

