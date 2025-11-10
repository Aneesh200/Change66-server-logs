# Quick Reference Guide

## 🚀 Quick Start

```bash
cd logs-ui
npm install
npm run dev
```

Visit: `http://localhost:3000`

## 📄 Pages

| Route | Description | Filter |
|-------|-------------|--------|
| `/` | Analytics Dashboard | None |
| `/logs` | All Logs | All types |
| `/logs/behavioral` | Behavioral Logs | Behavioral only |
| `/logs/telemetry` | Telemetry Logs | Telemetry only |
| `/logs/observability` | Observability Logs | Observability only |
| `/logs/errors` | Error Logs | Errors only |
| `/logs/performance` | Performance Logs | Performance only |

## 🎨 Component Reference

### Core Components

```typescript
// Sidebar Navigation
<Sidebar />

// Page Template (reusable)
<LogsPageTemplate 
  title="Page Title"
  description="Page description"
  eventType="behavioral" // optional
  icon={<Icon />}        // optional
/>

// Log Display
<LogsTable logs={logs} />

// Filters
<LogsFilter onFilter={handleFilter} isLoading={false} />

// Pagination
<Pagination 
  currentPage={1}
  totalPages={10}
  pageSize={50}
  totalCount={500}
  onPageChange={handlePageChange}
  onPageSizeChange={handlePageSizeChange}
/>
```

## 🔌 API Methods

```typescript
import { logsAPI } from '@/lib/api';

// Fetch filtered logs
const response = await logsAPI.getFilteredLogs({
  event_type: 'behavioral',
  priority: 'high',
  page: 1,
  page_size: 50
});

// Fetch recent logs
const logs = await logsAPI.getRecentLogs(100);

// Get metrics
const metrics = await logsAPI.getMetrics();

// Get server status
const status = await logsAPI.getStatus();

// Health check
const health = await logsAPI.healthCheck();
```

## 🎨 Color Codes

```typescript
const COLORS = {
  behavioral: '#3b82f6',    // Blue
  telemetry: '#10b981',     // Green
  observability: '#8b5cf6', // Purple
  error: '#ef4444',         // Red
  performance: '#f59e0b',   // Yellow
};
```

## ⚙️ Environment Variables

```env
# .env.local
NEXT_PUBLIC_API_URL=http://logs.biopeak.authify.tech
NEXT_PUBLIC_API_KEY=habit-tracker-key-dev
```

## 📊 Chart Examples

### Pie Chart
```typescript
<PieChart>
  <Pie 
    data={data}
    dataKey="value"
    nameKey="name"
    cx="50%"
    cy="50%"
  >
    {data.map((entry, index) => (
      <Cell key={index} fill={COLORS[entry.name]} />
    ))}
  </Pie>
</PieChart>
```

### Area Chart
```typescript
<AreaChart data={data}>
  <Area 
    type="monotone"
    dataKey="count"
    stroke="#3b82f6"
    fill="url(#colorGradient)"
  />
</AreaChart>
```

## 🔧 Common Tasks

### Add New Page

1. Create page file:
```typescript
// app/logs/custom/page.tsx
import { LogsPageTemplate } from '@/components/logs-page-template';

export default function CustomPage() {
  return (
    <LogsPageTemplate
      title="Custom Logs"
      description="Custom log view"
      eventType="custom"
    />
  );
}
```

2. Add to sidebar:
```typescript
// components/sidebar.tsx
const navItems = [
  // ... existing items
  { name: 'Custom', href: '/logs/custom', icon: Icon },
];
```

### Filter Logs by Multiple Criteria

```typescript
const handleFilter = (newFilter: LogFilter) => {
  const response = await logsAPI.getFilteredLogs({
    event_type: 'error',
    priority: 'high',
    start_time: '2024-01-01T00:00:00Z',
    end_time: '2024-12-31T23:59:59Z',
    page: 1,
    page_size: 100
  });
};
```

### Customize Chart Colors

```typescript
const customColors = {
  primary: '#ff0000',
  secondary: '#00ff00',
};

<Bar dataKey="value" fill={customColors.primary} />
```

## 🐛 Troubleshooting

### 404 Errors
```bash
# Check API URL in browser console
# Verify: http://logs.biopeak.authify.tech/api/v1/logs/filter

# Test API directly:
curl -H "X-API-Key: habit-tracker-key-dev" \
  http://logs.biopeak.authify.tech/api/v1/logs/filter?page=1&page_size=10
```

### No Data Showing
1. Check API key is correct
2. Verify server is running
3. Check browser console for errors
4. Ensure CORS is enabled on server

### Build Errors
```bash
# Clear cache and rebuild
rm -rf .next
npm run build
```

## 📦 Key Dependencies

```json
{
  "next": "16.0.1",
  "react": "19.2.0",
  "recharts": "^2.x",
  "shadcn/ui": "latest",
  "lucide-react": "^0.553.0",
  "date-fns": "^4.1.0",
  "sonner": "^2.0.7"
}
```

## 🎯 Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `Esc` | Close mobile menu |
| `Click outside` | Close dialogs/menus |

## 📱 Responsive Breakpoints

```css
sm: 640px   /* Mobile landscape */
md: 768px   /* Tablet */
lg: 1024px  /* Desktop */
xl: 1280px  /* Large desktop */
```

## 🔒 Security Notes

- Never commit `.env.local`
- API keys should be kept secret
- Use environment variables for all config
- Enable HTTPS in production

## 📈 Performance Tips

1. Use pagination (don't load all logs)
2. Apply filters to reduce data
3. Limit chart data points
4. Use React.memo for expensive components
5. Implement virtual scrolling for large tables

## 🎨 Customization

### Change Theme Colors
Edit `app/globals.css`:
```css
:root {
  --primary: 222.2 47.4% 11.2%;
  --background: 0 0% 100%;
  /* ... more variables */
}
```

### Modify Sidebar Width
Edit `components/sidebar.tsx`:
```typescript
className="w-64" // Change to w-72, w-80, etc.
```

## ✅ Checklist

Before deployment:
- [ ] Update `.env.local` with production values
- [ ] Test all pages
- [ ] Verify API connectivity
- [ ] Check responsive design
- [ ] Test on mobile devices
- [ ] Run `npm run build`
- [ ] Verify all charts render
- [ ] Test error states
- [ ] Check loading states

## 🚀 Deploy

```bash
# Build for production
npm run build

# Start production server
npm start

# Or deploy to Vercel
vercel deploy
```

## 📞 Quick Links

- API Documentation: Check `API_FILTERING.md`
- Server Setup: Check `README.md`
- Detailed Guide: Check `ENHANCED_FEATURES.md`

---

**Pro Tip**: Use the refresh button on each page to get the latest data without reloading the entire page!

