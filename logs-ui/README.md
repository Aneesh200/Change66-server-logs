# Analytics Logs Dashboard

A beautiful, modern UI for monitoring and analyzing your application logs in real-time. Built with Next.js and shadcn/ui.

## Features

- 📊 **Real-time Log Monitoring** - View and analyze logs as they come in
- 🔍 **Advanced Filtering** - Filter logs by event type, priority, user ID, session ID, and more
- 📄 **Pagination** - Efficiently browse through large log datasets
- 📱 **Responsive Design** - Works seamlessly on desktop, tablet, and mobile
- 🎨 **Beautiful UI** - Modern design with shadcn/ui components
- ⚡ **Fast Performance** - Optimized with Next.js 15 and React 19
- 📈 **Statistics Dashboard** - View key metrics at a glance

## Getting Started

### Prerequisites

- Node.js 18.x or later
- npm or yarn
- A running instance of the Change66 Log Server

### Installation

1. Navigate to the logs-ui directory:

```bash
cd logs-ui
```

2. Install dependencies:

```bash
npm install
```

3. Create a `.env.local` file with your configuration:

```bash
cp .env.example .env.local
```

4. Update the `.env.local` file with your log server URL:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

If your server requires authentication, add:

```env
NEXT_PUBLIC_API_KEY=your-api-key-here
```

### Running the Development Server

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) in your browser to see the dashboard.

### Building for Production

```bash
npm run build
npm start
```

## Configuration

### Environment Variables

- `NEXT_PUBLIC_API_URL` - The base URL of your log server (default: `http://localhost:8080`)
- `NEXT_PUBLIC_API_KEY` - Optional API key for authentication

## Features Overview

### Log Table

- View all logs in a clean, organized table
- Click on any log to view detailed information
- See event types with color-coded badges
- View timestamps in relative format (e.g., "2 minutes ago")

### Filters

Filter logs by:
- Event Type (behavioral, telemetry, observability, error, performance)
- Priority (normal, high)
- Event Name
- User ID
- Session ID
- App Version

### Statistics

View real-time statistics:
- Total log count
- Behavioral events count
- Error events count
- Performance metrics count

### Pagination

- Navigate through pages with first/previous/next/last buttons
- Customize page size (10, 25, 50, 100 logs per page)
- View current page and total results

## Tech Stack

- **Framework**: [Next.js 15](https://nextjs.org/) with App Router
- **UI Library**: [shadcn/ui](https://ui.shadcn.com/)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **Icons**: [Lucide React](https://lucide.dev/)
- **Date Formatting**: [date-fns](https://date-fns.org/)
- **Notifications**: [Sonner](https://sonner.emilkowal.ski/)

## Project Structure

```
logs-ui/
├── app/
│   ├── layout.tsx          # Root layout with metadata
│   ├── page.tsx            # Main dashboard page
│   └── globals.css         # Global styles
├── components/
│   ├── ui/                 # shadcn/ui components
│   ├── logs-table.tsx      # Log display table
│   ├── logs-filter.tsx     # Filter controls
│   └── pagination.tsx      # Pagination controls
├── lib/
│   ├── api.ts              # API client
│   ├── types.ts            # TypeScript types
│   └── utils.ts            # Utility functions
└── public/                 # Static assets
```

## API Integration

The dashboard connects to your log server's REST API:

- `GET /logs/filter` - Fetch filtered logs with pagination
- `GET /logs/recent` - Fetch recent logs
- `GET /health` - Health check endpoint

## Customization

### Changing the Theme

The dashboard uses shadcn/ui with Tailwind CSS. You can customize the theme by modifying:

- `app/globals.css` - CSS variables for colors
- `tailwind.config.ts` - Tailwind configuration
- `components.json` - shadcn/ui configuration

### Adding New Filters

To add new filter options:

1. Update the `LogFilter` type in `lib/types.ts`
2. Add the filter input in `components/logs-filter.tsx`
3. The API client will automatically include it in requests

## Troubleshooting

### Cannot connect to the server

Make sure:
1. The log server is running
2. The `NEXT_PUBLIC_API_URL` in `.env.local` is correct
3. CORS is enabled on the server (if running on different ports)

### Logs not displaying

Check:
1. The server has logs in the database
2. The API key is correct (if authentication is enabled)
3. Browser console for any error messages

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is part of the Change66 Log Server system.
