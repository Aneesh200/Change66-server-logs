'use client';

import { useEffect, useState } from 'react';
import { LogsTable } from '@/components/logs-table';
import { LogsFilter } from '@/components/logs-filter';
import { Pagination } from '@/components/pagination';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { logsAPI } from '@/lib/api';
import { AnalyticsLog, LogFilter } from '@/lib/types';
import { RefreshCw } from 'lucide-react';
import { toast } from 'sonner';

interface LogsPageTemplateProps {
  title: string;
  description: string;
  eventType?: string;
  icon?: React.ReactNode;
}

export function LogsPageTemplate({ title, description, eventType, icon }: LogsPageTemplateProps) {
  const [logs, setLogs] = useState<AnalyticsLog[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [pageSize, setPageSize] = useState(50);
  const [totalPages, setTotalPages] = useState(0);
  const [totalCount, setTotalCount] = useState(0);
  const [filter, setFilter] = useState<LogFilter>(eventType ? { event_type: eventType } : {});
  const [searchQuery, setSearchQuery] = useState('');

  const fetchLogs = async (newFilter?: LogFilter, page?: number, size?: number) => {
    setIsLoading(true);
    try {
      const finalFilter = {
        ...newFilter,
        ...(eventType && { event_type: eventType }), // Always filter by event type if specified
        page: page || currentPage,
        page_size: size || pageSize,
      };

      const response = await logsAPI.getFilteredLogs(finalFilter);
      setLogs(response.logs);
      setTotalPages(response.total_pages);
      setTotalCount(response.total_count);
    } catch (error) {
      console.error('Failed to fetch logs:', error);
      toast.error('Failed to fetch logs', {
        description: error instanceof Error ? error.message : 'Unknown error',
      });
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchLogs();
  }, []);

  const handleFilterChange = (newFilter: LogFilter) => {
    setFilter(newFilter);
    setSearchQuery(newFilter.search || '');
    setCurrentPage(1);
    fetchLogs(newFilter, 1, pageSize);
  };

  const handlePageChange = (page: number) => {
    setCurrentPage(page);
    fetchLogs(filter, page, pageSize);
  };

  const handlePageSizeChange = (size: number) => {
    setPageSize(size);
    setCurrentPage(1);
    fetchLogs(filter, 1, size);
  };

  const handleRefresh = () => {
    fetchLogs(filter, currentPage, pageSize);
    toast.success('Logs refreshed');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-zinc-950 to-black">
      <div className="container mx-auto p-6 space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            {icon}
            <div>
              <h1 className="text-4xl font-bold tracking-tight">{title}</h1>
              <p className="text-muted-foreground mt-2">{description}</p>
            </div>
          </div>
          <Button onClick={handleRefresh} disabled={isLoading}>
            <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>

        {/* Filters */}
        <LogsFilter 
          onFilter={handleFilterChange} 
          isLoading={isLoading} 
          hideEventTypeFilter={!!eventType}
        />

        {/* Logs Table */}
        <Card>
          <CardHeader>
            <CardTitle>Logs</CardTitle>
            <CardDescription>
              Showing {logs.length} of {totalCount} total logs
            </CardDescription>
          </CardHeader>
          <CardContent>
            {isLoading ? (
              <div className="flex items-center justify-center py-12">
                <RefreshCw className="h-8 w-8 animate-spin text-muted-foreground" />
              </div>
            ) : (
              <>
                <LogsTable logs={logs} searchQuery={searchQuery} />
                {totalCount > 0 && (
                  <Pagination
                    currentPage={currentPage}
                    totalPages={totalPages}
                    pageSize={pageSize}
                    totalCount={totalCount}
                    onPageChange={handlePageChange}
                    onPageSizeChange={handlePageSizeChange}
                  />
                )}
              </>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}

