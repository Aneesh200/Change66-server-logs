import { LogsPageTemplate } from '@/components/logs-page-template';
import { Database } from 'lucide-react';

export default function AllLogsPage() {
  return (
    <LogsPageTemplate
      title="All Logs"
      description="View and filter all application logs"
      icon={<Database className="h-10 w-10 text-primary" />}
    />
  );
}


