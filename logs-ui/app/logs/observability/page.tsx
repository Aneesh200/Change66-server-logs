import { LogsPageTemplate } from '@/components/logs-page-template';
import { Eye } from 'lucide-react';

export default function ObservabilityLogsPage() {
  return (
    <LogsPageTemplate
      title="Observability Logs"
      description="System monitoring and observability traces"
      eventType="observability"
      icon={<Eye className="h-10 w-10 text-purple-500" />}
    />
  );
}

