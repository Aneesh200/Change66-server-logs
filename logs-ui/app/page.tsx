'use client';

import { useEffect, useState } from 'react';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { logsAPI } from '@/lib/api';
import { RefreshCw, Activity, TrendingUp, Users, Zap, AlertTriangle, Database, Clock } from 'lucide-react';
import { Button } from '@/components/ui/button';
import { toast } from 'sonner';
import {
  LineChart,
  Line,
  BarChart,
  Bar,
  PieChart,
  Pie,
  Cell,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  Legend,
  ResponsiveContainer,
  AreaChart,
  Area
} from 'recharts';

export default function DashboardPage() {
  const [metrics, setMetrics] = useState<any>(null);
  const [recentLogs, setRecentLogs] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);

  const fetchDashboardData = async () => {
    setIsLoading(true);
    try {
      const [metricsData, logsData] = await Promise.all([
        logsAPI.getMetrics(),
        logsAPI.getRecentLogs(100),
      ]);

      setMetrics(metricsData);
      setRecentLogs(logsData);
    } catch (error) {
      console.error('Failed to fetch dashboard data:', error);
      toast.error('Failed to load dashboard data');
    } finally {
      setIsLoading(false);
    }
  };

  useEffect(() => {
    fetchDashboardData();
    // Refresh every 30 seconds
    const interval = setInterval(fetchDashboardData, 30000);
    return () => clearInterval(interval);
  }, []);

  // Process logs for charts
  const eventTypeData = recentLogs.reduce((acc: any, log) => {
    acc[log.event_type] = (acc[log.event_type] || 0) + 1;
    return acc;
  }, {});

  const pieData = Object.entries(eventTypeData).map(([name, value]) => ({
    name,
    value,
  }));

  // Group logs by hour for timeline
  const timelineData = recentLogs.reduce((acc: any, log) => {
    const hour = new Date(log.timestamp).getHours();
    const key = `${hour}:00`;
    if (!acc[key]) {
      acc[key] = { time: key, count: 0, errors: 0 };
    }
    acc[key].count++;
    if (log.event_type === 'error') {
      acc[key].errors++;
    }
    return acc;
  }, {});

  const timelineArray = Object.values(timelineData);

  // Priority distribution
  const priorityData = [
    { name: 'Normal', value: recentLogs.filter(l => l.priority === 'normal').length },
    { name: 'High', value: recentLogs.filter(l => l.priority === 'high').length },
  ];

  const COLORS = {
    behavioral: '#3b82f6',
    telemetry: '#10b981',
    observability: '#8b5cf6',
    error: '#ef4444',
    performance: '#f59e0b',
  };

  const PRIORITY_COLORS = ['#10b981', '#ef4444'];

  return (
    <div className="min-h-screen bg-gradient-to-br from-black via-zinc-950 to-black">
      <div className="container mx-auto p-6 space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-4xl font-bold tracking-tight">Dashboard</h1>
            <p className="text-muted-foreground mt-2">
              Real-time analytics and insights
            </p>
          </div>
          <Button onClick={fetchDashboardData} disabled={isLoading}>
            <RefreshCw className={`h-4 w-4 mr-2 ${isLoading ? 'animate-spin' : ''}`} />
            Refresh
          </Button>
        </div>

        {/* Stats Cards */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Total Logs</CardTitle>
              <Database className="h-4 w-4 text-muted-foreground" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metrics?.total_logs?.toLocaleString() || '0'}</div>
              <p className="text-xs text-muted-foreground">
                {metrics?.logs_last_hour || 0} in last hour
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Avg Latency</CardTitle>
              <Clock className="h-4 w-4 text-blue-500" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metrics?.average_latency?.toFixed(2) || '0'}ms</div>
              <p className="text-xs text-muted-foreground">Response time</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Error Rate</CardTitle>
              <AlertTriangle className="h-4 w-4 text-red-500" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metrics?.error_rate?.toFixed(2) || '0'}%</div>
              <p className="text-xs text-muted-foreground">
                {recentLogs.filter(l => l.event_type === 'error').length} errors
              </p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
              <CardTitle className="text-sm font-medium">Active Sessions</CardTitle>
              <Users className="h-4 w-4 text-green-500" />
            </CardHeader>
            <CardContent>
              <div className="text-2xl font-bold">{metrics?.active_sessions || '0'}</div>
              <p className="text-xs text-muted-foreground">Unique sessions</p>
            </CardContent>
          </Card>
        </div>

        {/* Charts Row 1 */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Event Distribution */}
          <Card>
            <CardHeader>
              <CardTitle>Event Type Distribution</CardTitle>
              <CardDescription>Breakdown of log events by type</CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <PieChart>
                  <Pie
                    data={pieData}
                    cx="50%"
                    cy="50%"
                    labelLine={false}
                    label={({ name, percent }) => `${name} ${percent ? (percent * 100).toFixed(0) : 0}%`}
                    outerRadius={80}
                    fill="#8884d8"
                    dataKey="value"
                  >
                    {pieData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={COLORS[entry.name as keyof typeof COLORS] || '#999'} />
                    ))}
                  </Pie>
                  <Tooltip />
                </PieChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>

          {/* Priority Distribution */}
          <Card>
            <CardHeader>
              <CardTitle>Priority Levels</CardTitle>
              <CardDescription>Distribution by priority</CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={priorityData}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="name" />
                  <YAxis />
                  <Tooltip />
                  <Bar dataKey="value" fill="#8884d8">
                    {priorityData.map((entry, index) => (
                      <Cell key={`cell-${index}`} fill={PRIORITY_COLORS[index]} />
                    ))}
                  </Bar>
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </div>

        {/* Charts Row 2 */}
        <div className="grid grid-cols-1 gap-6">
          {/* Timeline */}
          <Card>
            <CardHeader>
              <CardTitle>Activity Timeline</CardTitle>
              <CardDescription>Log volume over time (last 24 hours)</CardDescription>
            </CardHeader>
            <CardContent>
              <ResponsiveContainer width="100%" height={300}>
                <AreaChart data={timelineArray}>
                  <defs>
                    <linearGradient id="colorCount" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#3b82f6" stopOpacity={0.8}/>
                      <stop offset="95%" stopColor="#3b82f6" stopOpacity={0}/>
                    </linearGradient>
                    <linearGradient id="colorErrors" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="#ef4444" stopOpacity={0.8}/>
                      <stop offset="95%" stopColor="#ef4444" stopOpacity={0}/>
                    </linearGradient>
                  </defs>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="time" />
                  <YAxis />
                  <Tooltip />
                  <Legend />
                  <Area 
                    type="monotone" 
                    dataKey="count" 
                    stroke="#3b82f6" 
                    fillOpacity={1} 
                    fill="url(#colorCount)"
                    name="Total Logs"
                  />
                  <Area 
                    type="monotone" 
                    dataKey="errors" 
                    stroke="#ef4444" 
                    fillOpacity={1} 
                    fill="url(#colorErrors)"
                    name="Errors"
                  />
                </AreaChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </div>

        {/* Top Event Types */}
        <Card>
          <CardHeader>
            <CardTitle>Top Event Types</CardTitle>
            <CardDescription>Most frequent events</CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              {metrics?.top_event_types?.slice(0, 5).map((item: any, index: number) => (
                <div key={index} className="flex items-center">
                  <div className="flex-1">
                    <div className="flex items-center justify-between mb-1">
                      <span className="text-sm font-medium">{item.event_type}</span>
                      <span className="text-sm text-muted-foreground">{item.count}</span>
                    </div>
                    <div className="w-full bg-muted rounded-full h-2">
                      <div
                        className="h-2 rounded-full"
                        style={{
                          width: `${(item.count / recentLogs.length) * 100}%`,
                          backgroundColor: COLORS[item.event_type as keyof typeof COLORS] || '#999',
                        }}
                      />
                    </div>
                  </div>
                </div>
              )) || <p className="text-muted-foreground">No data available</p>}
            </div>
          </CardContent>
        </Card>

        {/* Quick Stats */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
          <Card>
            <CardHeader>
              <CardTitle className="text-sm">Logs Today</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold">{metrics?.logs_last_day?.toLocaleString() || '0'}</div>
              <p className="text-xs text-muted-foreground mt-2">24-hour period</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-sm">Performance Events</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-yellow-500">
                {recentLogs.filter(l => l.event_type === 'performance').length}
              </div>
              <p className="text-xs text-muted-foreground mt-2">Monitoring metrics</p>
            </CardContent>
          </Card>

          <Card>
            <CardHeader>
              <CardTitle className="text-sm">Observability</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-3xl font-bold text-purple-500">
                {recentLogs.filter(l => l.event_type === 'observability').length}
              </div>
              <p className="text-xs text-muted-foreground mt-2">System traces</p>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
