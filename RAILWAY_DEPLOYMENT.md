# Railway Deployment Guide - HTTP Support ✅

Deploy your log server to Railway with full HTTP support (no forced redirects!)

---

## 🎯 Why Railway?

- ✅ **HTTP Support**: No forced HTTPS redirects (perfect for your production app)
- ✅ **HTTPS Support**: Also works with HTTPS for future updates
- ✅ **Easy Deployment**: Similar to Render
- ✅ **Free Tier**: $5 credit/month
- ✅ **PostgreSQL**: Built-in database

---

## 📋 Step-by-Step Deployment

### Step 1: Install Railway CLI

Run these commands in your terminal:

```bash
# Fix npm permissions first (if needed)
sudo chown -R $(whoami) ~/.npm

# Install Railway CLI
npm install -g @railway/cli

# Verify installation
railway --version
```

**OR use Homebrew:**

```bash
brew install railway
```

---

### Step 2: Login to Railway

```bash
# This will open your browser to authenticate
railway login
```

You'll need to:
1. Create a Railway account (free): https://railway.app
2. Authorize the CLI

---

### Step 3: Navigate to Your Project

```bash
cd /Users/aneesh/Change66-Log-Server
```

---

### Step 4: Initialize Railway Project

```bash
# Create new Railway project
railway init

# You'll be prompted:
# - Project name: change66-log-server (or your choice)
# - Select "Empty Project"
```

---

### Step 5: Add PostgreSQL Database

```bash
# Add PostgreSQL to your project
railway add --database postgres

# Railway will automatically provision a PostgreSQL database
```

---

### Step 6: Set Environment Variables

Create a `.env.railway` file (Railway-specific):

```bash
# Copy your environment config
cp env.production .env.railway
```

Then set the variables on Railway:

```bash
# Set your API key
railway variables set API_KEYS=habit-tracker-key-dev

# Set server configuration
railway variables set PORT=8080
railway variables set GIN_MODE=release
railway variables set LOG_LEVEL=info
railway variables set ENVIRONMENT=production

# Performance settings
railway variables set RATE_LIMIT_REQUESTS_PER_MINUTE=1000
railway variables set RATE_LIMIT_BURST=100
railway variables set MAX_BATCH_SIZE=1000
railway variables set WORKER_POOL_SIZE=10
railway variables set REQUEST_TIMEOUT_SECONDS=30
railway variables set MAX_REQUEST_SIZE_MB=10

# Features
railway variables set ENABLE_METRICS=true
railway variables set ENABLE_CORS=true
railway variables set ALLOWED_ORIGINS=*

# Database config (Railway auto-provides these from the PostgreSQL service)
railway variables set DB_HOST=${{Postgres.PGHOST}}
railway variables set DB_PORT=${{Postgres.PGPORT}}
railway variables set DB_NAME=${{Postgres.PGDATABASE}}
railway variables set DB_USER=${{Postgres.PGUSER}}
railway variables set DB_PASSWORD=${{Postgres.PGPASSWORD}}
railway variables set DB_SSL_MODE=require
```

**OR set them via Railway Dashboard:**
1. Go to https://railway.app/dashboard
2. Select your project
3. Click on your service → Variables
4. Add each variable manually

---

### Step 7: Deploy to Railway

```bash
# Link your local directory to Railway project
railway link

# Deploy your application
railway up

# Railway will:
# ✅ Build your Docker image
# ✅ Deploy to production
# ✅ Assign a public URL
# ✅ Connect to PostgreSQL
```

---

### Step 8: Get Your Railway URL

```bash
# Get your deployment URL
railway status

# Or view in dashboard
railway open
```

You'll get a URL like:
```
https://change66-log-server-production.up.railway.app
```

**Important: This URL works with BOTH HTTP and HTTPS!**

---

### Step 9: Configure Custom Domain

#### In Railway Dashboard:

1. Go to your project → **Settings** → **Domains**
2. Click **Add Custom Domain**
3. Enter: `logs.biopeak.authify.tech`
4. Railway will show you the CNAME record

#### Update Your DNS:

Add this CNAME record at your DNS provider:

```
Type: CNAME
Name: logs.biopeak.authify
Value: change66-log-server-production.up.railway.app
TTL: 3600
```

**DNS providers:**
- Cloudflare: Dashboard → DNS
- AWS Route 53: Hosted zones → Create record
- Namecheap: Advanced DNS
- GoDaddy: DNS Management

---

### Step 10: Test Your Deployment

#### Test with Railway URL (HTTP):

```bash
# Health check
curl http://change66-log-server-production.up.railway.app/health

# Test with API key
curl -H "X-API-Key: habit-tracker-key-dev" \
  http://change66-log-server-production.up.railway.app/api/v1/status
```

#### Test log ingestion:

```bash
curl -X POST http://change66-log-server-production.up.railway.app/api/v1/ingest \
  -H "Content-Type: application/json" \
  -H "X-API-Key: habit-tracker-key-dev" \
  -d '{
    "event_id": "railway-test-'$(date +%s)'",
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    "event_type": "behavioral",
    "event_name": "railway_deployment_test",
    "properties": {"test": true, "platform": "railway"},
    "user_id": "test-user",
    "session_id": "test-session",
    "app_version": "1.0.0",
    "device_info": {"platform": "test"},
    "sequence_number": 1,
    "priority": "normal"
  }'
```

#### After DNS propagation, test with custom domain (HTTP):

```bash
# This is what your production app will use!
curl http://logs.biopeak.authify.tech/health

# Should return 200 (not 307!)
```

---

## 🔍 Verify HTTP Works (No Redirects!)

```bash
# Check that HTTP doesn't redirect
curl -v http://logs.biopeak.authify.tech/health 2>&1 | grep "HTTP/"

# Should see:
# < HTTP/1.1 200 OK
# NOT HTTP/1.1 307 Temporary Redirect
```

---

## 📊 View Logs and Monitoring

```bash
# View live logs
railway logs

# Follow logs in real-time
railway logs --follow

# Open Railway dashboard
railway open
```

**In the Dashboard you can see:**
- Real-time metrics (CPU, Memory, Network)
- Deployment history
- Environment variables
- Database queries
- Request logs

---

## 🔧 Useful Railway Commands

```bash
# Check deployment status
railway status

# Restart service
railway restart

# View environment variables
railway variables

# Connect to PostgreSQL
railway connect postgres

# Open project in browser
railway open

# View service logs
railway logs --tail 100

# Redeploy
railway up --detach
```

---

## 🔄 Update Your Flutter App Config

Once deployed, your production app with HTTP will work:

```dart
// Your existing production app (works now!)
static const String LOG_SERVER_URL = 'http://logs.biopeak.authify.tech/api/v1';
static const String API_KEY = 'habit-tracker-key-dev';
```

**Both HTTP and HTTPS work!** ✅

---

## 💰 Railway Pricing

### Free Tier:
- $5 credit/month
- ~500 hours of usage
- Perfect for testing

### Paid (if you exceed free tier):
- $5/month base
- Pay-as-you-go for usage
- ~$5-10/month typical for this app

---

## 🚨 Troubleshooting

### Build Fails

```bash
# Check logs
railway logs

# Common issues:
# - Go version mismatch (fixed in Dockerfile)
# - Missing environment variables
# - Database connection issues
```

### Database Connection Issues

```bash
# Connect to database to test
railway connect postgres

# Check if migrations ran
\dt

# Should see 'analytics_logs' table
```

### HTTP Still Redirects

If HTTP still redirects after DNS update:
1. Clear DNS cache: `sudo dscacheutil -flushcache`
2. Wait 5-10 minutes for DNS propagation
3. Test with Railway URL first before custom domain

### Environment Variables Not Working

```bash
# List all variables
railway variables

# Set missing ones
railway variables set KEY=value

# Redeploy
railway up
```

---

## 📈 Migration from Render

If you want to migrate your existing data from Render:

### Option 1: Export/Import Database

```bash
# From Render PostgreSQL
pg_dump -h [render-host] -U [user] -d analytics_logs > logs_backup.sql

# To Railway PostgreSQL
railway connect postgres
\i logs_backup.sql
```

### Option 2: Fresh Start (Recommended)

Just start fresh on Railway - your app will start sending new logs immediately.

---

## ✅ Deployment Checklist

- [ ] Railway CLI installed
- [ ] Logged into Railway
- [ ] Project initialized
- [ ] PostgreSQL database added
- [ ] Environment variables set
- [ ] Application deployed
- [ ] Deployment URL working (HTTP & HTTPS)
- [ ] Custom domain configured
- [ ] DNS CNAME record added
- [ ] HTTP endpoint tested (no redirects!)
- [ ] Logs flowing from app
- [ ] Monitoring set up

---

## 🎉 Success Criteria

You'll know it's working when:

1. ✅ `curl http://logs.biopeak.authify.tech/health` returns 200 (not 307)
2. ✅ Your production app can send logs via HTTP
3. ✅ Logs appear in Railway dashboard
4. ✅ Database shows increasing log count

---

## 🆘 Need Help?

- Railway Docs: https://docs.railway.app
- Railway Discord: https://discord.gg/railway
- Railway Status: https://status.railway.app

---

## 🚀 Quick Start Summary

```bash
# 1. Install
npm install -g @railway/cli

# 2. Login
railway login

# 3. Navigate
cd /Users/aneesh/Change66-Log-Server

# 4. Initialize & Deploy
railway init
railway add --database postgres
railway up

# 5. Set variables
railway variables set API_KEYS=habit-tracker-key-dev
railway variables set PORT=8080
# ... (set all other variables)

# 6. Get URL
railway status

# 7. Configure DNS
# Add CNAME: logs.biopeak.authify.tech → [railway-url]

# 8. Test
curl http://logs.biopeak.authify.tech/health
```

---

**Your log server will now work with HTTP! 🎉**

No more 307 redirects - your production app can send logs immediately!

