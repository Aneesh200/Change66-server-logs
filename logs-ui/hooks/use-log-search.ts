import { useMemo } from 'react';
import { AnalyticsLog } from '@/lib/types';

/**
 * Client-side search hook for filtering logs by search query
 * Searches through event names, user IDs, session IDs, and metadata/properties
 */
export function useLogSearch(logs: AnalyticsLog[], searchQuery: string) {
  return useMemo(() => {
    if (!searchQuery || searchQuery.trim() === '') {
      return logs;
    }

    const query = searchQuery.toLowerCase().trim();

    return logs.filter((log) => {
      // Search in basic fields
      if (log.event_name?.toLowerCase().includes(query)) return true;
      if (log.event_type?.toLowerCase().includes(query)) return true;
      if (log.user_id?.toLowerCase().includes(query)) return true;
      if (log.session_id?.toLowerCase().includes(query)) return true;
      if (log.app_version?.toLowerCase().includes(query)) return true;
      if (log.event_id?.toLowerCase().includes(query)) return true;
      if (log.priority?.toLowerCase().includes(query)) return true;

      // Search in properties (metadata)
      if (log.properties) {
        const propsString = JSON.stringify(log.properties).toLowerCase();
        if (propsString.includes(query)) return true;
      }

      // Search in device_info
      if (log.device_info) {
        const deviceString = JSON.stringify(log.device_info).toLowerCase();
        if (deviceString.includes(query)) return true;
      }

      return false;
    });
  }, [logs, searchQuery]);
}

