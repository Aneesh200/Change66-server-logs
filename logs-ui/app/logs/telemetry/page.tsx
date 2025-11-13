import { LogsPageTemplate } from '@/components/logs-page-template';
import { Activity } from 'lucide-react';

export default function TelemetryLogsPage() {
  return (
    <LogsPageTemplate
      title="Telemetry Logs"
      description="System telemetry and diagnostic data"
      eventType="telemetry"
      icon={<Activity className="h-10 w-10 text-green-500" />}
    />
  );
}


