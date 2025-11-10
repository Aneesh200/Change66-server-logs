'use client';

import { AnalyticsLog } from '@/lib/types';
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from '@/components/ui/dialog';
import { formatDistanceToNow } from 'date-fns';
import { Eye } from 'lucide-react';

interface LogsTableProps {
  logs: AnalyticsLog[];
}

export function LogsTable({ logs }: LogsTableProps) {
  const getEventTypeBadge = (eventType: string) => {
    const colors: Record<string, string> = {
      behavioral: 'bg-blue-500',
      telemetry: 'bg-green-500',
      observability: 'bg-purple-500',
      error: 'bg-red-500',
      performance: 'bg-yellow-500',
    };

    return (
      <Badge className={colors[eventType] || 'bg-gray-500'}>
        {eventType}
      </Badge>
    );
  };

  const getPriorityBadge = (priority: string) => {
    return (
      <Badge variant={priority === 'high' ? 'destructive' : 'secondary'}>
        {priority}
      </Badge>
    );
  };

  const formatTimestamp = (timestamp: string) => {
    try {
      return formatDistanceToNow(new Date(timestamp), { addSuffix: true });
    } catch (error) {
      return timestamp;
    }
  };

  return (
    <div className="rounded-md border">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Event Type</TableHead>
            <TableHead>Event Name</TableHead>
            <TableHead>Priority</TableHead>
            <TableHead>User ID</TableHead>
            <TableHead>Session ID</TableHead>
            <TableHead>App Version</TableHead>
            <TableHead>Timestamp</TableHead>
            <TableHead className="text-right">Actions</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {logs.length === 0 ? (
            <TableRow>
              <TableCell colSpan={8} className="text-center text-muted-foreground">
                No logs found
              </TableCell>
            </TableRow>
          ) : (
            logs.map((log) => (
              <TableRow key={log.id}>
                <TableCell>{getEventTypeBadge(log.event_type)}</TableCell>
                <TableCell className="font-medium">{log.event_name}</TableCell>
                <TableCell>{getPriorityBadge(log.priority)}</TableCell>
                <TableCell className="text-sm text-muted-foreground">
                  {log.user_id || '-'}
                </TableCell>
                <TableCell className="text-sm text-muted-foreground truncate max-w-[150px]">
                  {log.session_id || '-'}
                </TableCell>
                <TableCell className="text-sm">{log.app_version || '-'}</TableCell>
                <TableCell className="text-sm text-muted-foreground">
                  {formatTimestamp(log.timestamp)}
                </TableCell>
                <TableCell className="text-right">
                  <Dialog>
                    <DialogTrigger asChild>
                      <Button variant="ghost" size="sm">
                        <Eye className="h-4 w-4" />
                      </Button>
                    </DialogTrigger>
                    <DialogContent className="max-w-3xl max-h-[80vh] overflow-y-auto">
                      <DialogHeader>
                        <DialogTitle>Log Details</DialogTitle>
                        <DialogDescription>
                          Event ID: {log.event_id}
                        </DialogDescription>
                      </DialogHeader>
                      <div className="space-y-4">
                        <div className="grid grid-cols-2 gap-4">
                          <div>
                            <p className="text-sm font-semibold text-muted-foreground">Event Type</p>
                            <div className="mt-1">{getEventTypeBadge(log.event_type)}</div>
                          </div>
                          <div>
                            <p className="text-sm font-semibold text-muted-foreground">Priority</p>
                            <div className="mt-1">{getPriorityBadge(log.priority)}</div>
                          </div>
                          <div>
                            <p className="text-sm font-semibold text-muted-foreground">Event Name</p>
                            <p className="mt-1">{log.event_name}</p>
                          </div>
                          <div>
                            <p className="text-sm font-semibold text-muted-foreground">Timestamp</p>
                            <p className="mt-1 text-sm">{new Date(log.timestamp).toLocaleString()}</p>
                          </div>
                          {log.user_id && (
                            <div>
                              <p className="text-sm font-semibold text-muted-foreground">User ID</p>
                              <p className="mt-1 text-sm font-mono">{log.user_id}</p>
                            </div>
                          )}
                          {log.session_id && (
                            <div>
                              <p className="text-sm font-semibold text-muted-foreground">Session ID</p>
                              <p className="mt-1 text-sm font-mono truncate">{log.session_id}</p>
                            </div>
                          )}
                          {log.app_version && (
                            <div>
                              <p className="text-sm font-semibold text-muted-foreground">App Version</p>
                              <p className="mt-1">{log.app_version}</p>
                            </div>
                          )}
                          {log.sequence_number && (
                            <div>
                              <p className="text-sm font-semibold text-muted-foreground">Sequence Number</p>
                              <p className="mt-1">{log.sequence_number}</p>
                            </div>
                          )}
                        </div>
                        
                        {Object.keys(log.properties || {}).length > 0 && (
                          <div>
                            <p className="text-sm font-semibold text-muted-foreground mb-2">Properties</p>
                            <pre className="bg-muted p-4 rounded-lg text-xs overflow-x-auto">
                              {JSON.stringify(log.properties, null, 2)}
                            </pre>
                          </div>
                        )}
                        
                        {Object.keys(log.device_info || {}).length > 0 && (
                          <div>
                            <p className="text-sm font-semibold text-muted-foreground mb-2">Device Info</p>
                            <pre className="bg-muted p-4 rounded-lg text-xs overflow-x-auto">
                              {JSON.stringify(log.device_info, null, 2)}
                            </pre>
                          </div>
                        )}
                      </div>
                    </DialogContent>
                  </Dialog>
                </TableCell>
              </TableRow>
            ))
          )}
        </TableBody>
      </Table>
    </div>
  );
}

