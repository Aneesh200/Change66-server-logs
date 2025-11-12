import { LogsPageTemplate } from '@/components/logs-page-template';
import { Users } from 'lucide-react';

export default function BehavioralLogsPage() {
  return (
    <LogsPageTemplate
      title="Behavioral Logs"
      description="User interactions and behavioral analytics"
      eventType="behavioral"
      icon={<Users className="h-10 w-10 text-blue-500" />}
    />
  );
}

