import { LogsPageTemplate } from '@/components/logs-page-template';
import { AlertTriangle } from 'lucide-react';

export default function ErrorLogsPage() {
  return (
    <LogsPageTemplate
      title="Error Logs"
      description="Application errors and exceptions"
      eventType="error"
      icon={<AlertTriangle className="h-10 w-10 text-red-500" />}
    />
  );
}


