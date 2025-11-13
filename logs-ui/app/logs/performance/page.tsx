import { LogsPageTemplate } from '@/components/logs-page-template';
import { Zap } from 'lucide-react';

export default function PerformanceLogsPage() {
  return (
    <LogsPageTemplate
      title="Performance Logs"
      description="Performance metrics and timing data"
      eventType="performance"
      icon={<Zap className="h-10 w-10 text-yellow-500" />}
    />
  );
}


