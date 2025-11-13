# Logs UI Setup Guide

This guide will walk you through setting up and running the Analytics Logs Dashboard alongside your log server.

## Prerequisites

Before you begin, ensure you have:

- ✅ Node.js 18.x or later installed
- ✅ npm or yarn package manager
- ✅ Your Change66 Log Server running and accessible
- ✅ Basic knowledge of Next.js and React

## Quick Start (5 minutes)

### Step 1: Install Dependencies

Navigate to the logs-ui directory and install dependencies:

```bash
cd logs-ui
npm install
```

### Step 2: Configure the API Connection

Create a `.env.local` file in the `logs-ui` directory:

```bash
# For local development
NEXT_PUBLIC_API_URL=http://localhost:8080
```

If your log server is running on a different host or port, update the URL accordingly.

### Step 3: Enable CORS on the Server

Your log server needs to allow requests from the UI. Update your server's configuration file (usually `config.env` or `.env`):

```env
ENABLE_CORS=true
ALLOWED_ORIGINS=http://localhost:3000
```

If you want to allow all origins (not recommended for production):

```env
ENABLE_CORS=true
ALLOWED_ORIGINS=*
```

Restart your log server for the changes to take effect.

### Step 4: Start the Development Server

```bash
npm run dev
```

The UI will be available at [http://localhost:3000](http://localhost:3000)

## Authentication Setup

If your log server requires API key authentication:

1. Get your API key from your server configuration or admin panel

2. Add it to your `.env.local` file:

```env
NEXT_PUBLIC_API_URL=http://localhost:8080
NEXT_PUBLIC_API_KEY=your-api-key-here
```

3. Restart the Next.js development server

## Production Deployment

### Build the Application

```bash
npm run build
```

### Run in Production Mode

```bash
npm start
```

The production server will start on port 3000 by default.

### Environment Variables for Production

Create a `.env.production` file or set environment variables in your hosting platform:

```env
NEXT_PUBLIC_API_URL=https://your-production-server.com
NEXT_PUBLIC_API_KEY=your-production-api-key
```

## Deployment Options

### Option 1: Vercel (Recommended for Next.js)

1. Push your code to GitHub
2. Import the project in [Vercel](https://vercel.com)
3. Set the root directory to `logs-ui`
4. Add environment variables in Vercel dashboard
5. Deploy

### Option 2: Docker

Create a `Dockerfile` in the `logs-ui` directory:

```dockerfile
FROM node:18-alpine AS base

# Install dependencies
FROM base AS deps
WORKDIR /app
COPY package*.json ./
RUN npm ci

# Build the application
FROM base AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

# Production image
FROM base AS runner
WORKDIR /app

ENV NODE_ENV=production

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000

CMD ["node", "server.js"]
```

Build and run:

```bash
docker build -t logs-ui .
docker run -p 3000:3000 -e NEXT_PUBLIC_API_URL=http://your-server:8080 logs-ui
```

### Option 3: Static Export (if you don't need server-side features)

1. Update `next.config.js`:

```javascript
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'export',
}

module.exports = nextConfig
```

2. Build the static site:

```bash
npm run build
```

3. The static files will be in the `out` directory. Deploy them to any static hosting service (Netlify, GitHub Pages, AWS S3, etc.)

## Troubleshooting

### Issue: "Failed to fetch logs"

**Possible causes:**

1. **Server not running**: Ensure your log server is running and accessible
   ```bash
   curl http://localhost:8080/health
   ```

2. **CORS not enabled**: Check server logs for CORS errors and ensure `ENABLE_CORS=true`

3. **Wrong API URL**: Verify the `NEXT_PUBLIC_API_URL` in your `.env.local` matches your server

4. **Authentication required**: Add `NEXT_PUBLIC_API_KEY` if your server requires it

### Issue: "No logs found"

**Possible causes:**

1. **No logs in database**: The server database might be empty. Try ingesting some test logs first

2. **Authentication issue**: Check browser console for 401/403 errors

3. **Wrong API endpoint**: Verify the API client is using `/api/v1/logs/filter`

### Issue: Styles not loading correctly

**Solution:**

1. Clear Next.js cache:
   ```bash
   rm -rf .next
   npm run dev
   ```

2. Verify Tailwind CSS is configured correctly in `tailwind.config.ts`

### Issue: Environment variables not working

**Remember:**

- Environment variables in Next.js must be prefixed with `NEXT_PUBLIC_` to be available in the browser
- Restart the dev server after changing `.env.local`
- For production builds, environment variables must be set before building

## Customization

### Changing the Port

To run the UI on a different port:

```bash
PORT=4000 npm run dev
```

### Updating the Theme

Edit `app/globals.css` to customize colors:

```css
:root {
  --background: 0 0% 100%;
  --foreground: 222.2 84% 4.9%;
  /* ... other color variables ... */
}
```

### Adding New Filter Fields

1. Update the `LogFilter` interface in `lib/types.ts`
2. Add the filter input in `components/logs-filter.tsx`
3. The API client will automatically include it in requests

## Performance Tips

1. **Pagination**: Use reasonable page sizes (50-100) for best performance
2. **Filters**: Apply filters to reduce the amount of data transferred
3. **Caching**: Consider adding React Query for client-side caching
4. **CDN**: In production, serve static assets from a CDN

## Security Best Practices

1. **Never commit `.env.local`**: Add it to `.gitignore`
2. **Use environment variables**: Don't hardcode API keys in code
3. **HTTPS in production**: Always use HTTPS for production deployments
4. **Restrict CORS**: Don't use `ALLOWED_ORIGINS=*` in production
5. **Rate limiting**: Ensure your server has rate limiting enabled

## Getting Help

If you encounter issues:

1. Check the browser console for errors
2. Check the server logs for API errors
3. Verify CORS configuration
4. Test the API endpoints directly with curl or Postman
5. Refer to the main README.md for server-specific issues

## Next Steps

- Explore the filter options to find specific logs
- View detailed log information by clicking the eye icon
- Adjust page sizes based on your needs
- Consider adding more features like export to CSV, real-time updates, or custom dashboards

Happy logging! 🚀


