# Deploy to Render - Step-by-Step Guide

## 🎯 Goal
Deploy your log server to: `http://logs.biopeak.authify.tech/api/v1`

API Key: `habit-tracker-key-dev`

---

## 📋 Prerequisites

1. GitHub account with your code pushed
2. Render account (sign up at https://render.com - it's free)
3. Domain access to `biopeak.authify.tech` for DNS configuration

---

## 🚀 Step-by-Step Deployment

### Step 1: Push Your Code to GitHub

```bash
# Make sure your code is committed
git add .
git commit -m "Ready for Render deployment"
git push origin main
```

---

### Step 2: Create PostgreSQL Database on Render

1. Go to https://dashboard.render.com/
2. Click **"New +"** → **"PostgreSQL"**
3. Configure:
   - **Name**: `log-server-db` (or any name you prefer)
   - **Database**: `analytics_logs`
   - **User**: (auto-generated, note it down)
   - **Region**: Choose closest to you (e.g., Oregon, Frankfurt)
   - **Plan**: Free (or paid for better performance)
4. Click **"Create Database"**
5. Wait for database to be created (~1-2 minutes)
6. **IMPORTANT**: Copy these values (you'll need them soon):
   - **Internal Database URL** (starts with `postgres://`) postgresql://    analytics_logs_user:Gcai1cMF6763HjkA9aMpTddhm8nH3U1m@dpg-d493bsmmcj7s73e9c5pg-a/analytics_logs
   - **Hostname** (e.g., `dpg-xxxxx.oregon-postgres.render.com`) dpg-d493bsmmcj7s73e9c5pg-a
   - **Port** (usually `5432`)
   - **Database** name analytics_logs
   - **Username** analytics_logs_user
   - **Password** Gcai1cMF6763HjkA9aMpTddhm8nH3U1m

---

### Step 3: Create Web Service on Render

1. In Render Dashboard, click **"New +"** → **"Web Service"**
2. Click **"Connect GitHub"** and authorize Render
3. Select your repository: `Change66-Log-Server`
4. Configure the service:

   **Basic Settings:**
   - **Name**: `log-ingestion-server` (this becomes part of your URL)
   - **Region**: Same as your database (important!)
   - **Branch**: `main`
   - **Root Directory**: Leave empty
   - **Environment**: **Docker**
   - **Plan**: Free (or paid for production)

5. Click **"Create Web Service"** (don't worry about env vars yet)

---

### Step 4: Configure Environment Variables

1. Go to your web service in Render
2. Click **"Environment"** tab
3. Add these environment variables (click "Add Environment Variable" for each):

```bash
# Database Configuration (Use values from Step 2)
DB_HOST=dpg-xxxxx.oregon-postgres.render.com
DB_PORT=5432
DB_NAME=analytics_logs
DB_USER=your_db_username_from_step2
DB_PASSWORD=your_db_password_from_step2
DB_SSL_MODE=require

# API Configuration - YOUR SPECIFIC KEY
API_KEYS=habit-tracker-key-dev

# Server Configuration
PORT=8080
GIN_MODE=release
LOG_LEVEL=info
ENVIRONMENT=production

# Performance Settings
RATE_LIMIT_REQUESTS_PER_MINUTE=1000
RATE_LIMIT_BURST=100
MAX_BATCH_SIZE=1000
WORKER_POOL_SIZE=10

# Features
ENABLE_METRICS=true
ENABLE_CORS=true
ALLOWED_ORIGINS=*
```

4. Click **"Save Changes"**
5. Render will automatically redeploy with new environment variables

---

### Step 5: Wait for Deployment

1. Go to **"Logs"** tab to watch the build process
2. You'll see:
   - Docker image being built
   - Application starting
   - Health checks passing
3. Wait until you see: **"Your service is live 🎉"**
4. This usually takes 2-5 minutes

**Note the auto-generated URL**: `https://log-ingestion-server.onrender.com`

---

### Step 6: Test Your Deployment

```bash
# Test health endpoint
curl https://log-ingestion-server.onrender.com/health

# Expected response:
# {"status":"healthy","timestamp":"..."}

# Test API with your key
curl -H "X-API-Key: habit-tracker-key-dev" \
  https://log-ingestion-server.onrender.com/api/v1/status

# Expected response with database stats
```

---

### Step 7: Configure Custom Domain (logs.biopeak.authify.tech)

#### 7.1 Add Custom Domain in Render

1. In your web service, click **"Settings"** tab
2. Scroll to **"Custom Domains"** section
3. Click **"Add Custom Domain"**
4. Enter: `logs.biopeak.authify.tech`
5. Click **"Save"**
6. Render will show you DNS configuration instructions

#### 7.2 Configure DNS at Your Domain Provider

You need to add **one of these** DNS records (Render will tell you which):

**Option A: CNAME Record (Recommended)**
```
Type: CNAME
Name: logs.biopeak.authify
Value: log-ingestion-server.onrender.com
TTL: 3600 (or Auto)
```

**Option B: A Record**
```
Type: A
Name: logs.biopeak.authify
Value: [IP address from Render]
TTL: 3600 (or Auto)
```

**Where to do this:**
- If using Cloudflare: Dashboard → DNS → Add record
- If using AWS Route 53: Hosted zones → biopeak.authify.tech → Create record
- If using Namecheap: Domain → Advanced DNS → Add record
- If using GoDaddy: DNS Management → Add record

#### 7.3 Wait for DNS Propagation

- DNS changes take 5-60 minutes to propagate
- Check status in Render Dashboard (it will show "Verifying...")
- Once verified, Render automatically provisions SSL certificate (free)

#### 7.4 Verify Custom Domain

```bash
# Test with custom domain (after DNS propagates)
curl https://logs.biopeak.authify.tech/health

# Test API endpoint
curl -H "X-API-Key: habit-tracker-key-dev" \
  https://logs.biopeak.authify.tech/api/v1/status
```

---

## ✅ Your Flutter App Configuration

Once deployed, update your Flutter app:

```dart
class ApiConfig {
  // Production endpoint
  static const String LOG_SERVER_URL = 'https://logs.biopeak.authify.tech/api/v1';
  static const String API_KEY = 'habit-tracker-key-dev';
}
```

**Note**: Use `https://` not `http://` - Render provides free SSL!

---

## 🔧 Post-Deployment Configuration

### Monitor Your Service

1. **View Logs**: Render Dashboard → Your Service → Logs
2. **Metrics**: Dashboard shows CPU, Memory, Response times
3. **Health Checks**: Automatic monitoring at `/health`

### Set Up Alerts (Recommended)

1. Go to **"Settings"** → **"Alerts"**
2. Enable:
   - Health check failures
   - Deploy notifications
   - Resource usage alerts

### Auto-Deploy from GitHub

Already configured! Every push to `main` branch will:
1. Trigger automatic rebuild
2. Deploy new version
3. Run health checks
4. Rollback if health checks fail

---

## 🚨 Troubleshooting

### Issue: Build Fails

**Check:**
```bash
# In Render logs, look for:
- Go version compatibility
- Missing dependencies
- Docker build errors
```

**Fix:**
- Verify Dockerfile is correct
- Check go.mod and go.sum are committed
- Ensure all imports are available

### Issue: Database Connection Failed

**Check:**
```bash
# Verify environment variables:
- DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD
- DB_SSL_MODE should be "require" for Render
```

**Fix:**
1. Go to Database → Connections
2. Copy "Internal Database URL"
3. Paste in environment variables

### Issue: Health Check Failing

**Check:**
```bash
# View logs for errors:
- Database migration failures
- Port binding issues
- Configuration errors
```

**Fix:**
- Ensure PORT=8080 in environment
- Verify database is accessible
- Check API_KEYS is set

### Issue: Custom Domain Not Working

**Check:**
1. DNS propagation: https://dnschecker.org/
2. Verify CNAME/A record is correct
3. Check Render shows "Verified" status

**Fix:**
- Wait 30-60 minutes for DNS
- Verify DNS record at your provider
- Check for conflicting records

---

## 📊 Performance Optimization

### Free Tier Limitations

- Service sleeps after 15 minutes of inactivity
- First request after sleep takes ~30 seconds (cold start)
- 750 hours/month free

### Upgrade to Paid ($7/month) for:

- No sleeping (always on)
- Instant response times
- Better performance
- More CPU/Memory

### Database Optimization

Free PostgreSQL includes:
- 90 days data retention
- 1 GB storage
- Automatic backups (7 days)

Upgrade for:
- More storage
- Longer retention
- Better performance

---

## 🎯 Testing Your API Endpoints

### Full API Test Script

```bash
#!/bin/bash

BASE_URL="https://logs.biopeak.authify.tech/api/v1"
API_KEY="habit-tracker-key-dev"

echo "Testing Log Server..."

# Test 1: Health Check
echo -e "\n1. Health Check:"
curl -s "$BASE_URL/../health" | jq .

# Test 2: Status Check
echo -e "\n2. Status Check:"
curl -s -H "X-API-Key: $API_KEY" "$BASE_URL/status" | jq .

# Test 3: Ingest Log
echo -e "\n3. Ingest Log:"
curl -s -X POST "$BASE_URL/ingest" \
  -H "Content-Type: application/json" \
  -H "X-API-Key: $API_KEY" \
  -d '{
    "event_id": "test-'$(date +%s)'",
    "timestamp": "'$(date -u +%Y-%m-%dT%H:%M:%SZ)'",
    "event_type": "behavioral",
    "event_name": "test_event",
    "properties": {"test": true},
    "user_id": "test-user",
    "session_id": "test-session",
    "app_version": "1.0.0",
    "device_info": {"platform": "test"},
    "sequence_number": 1,
    "priority": "normal"
  }' | jq .

# Test 4: Recent Logs
echo -e "\n4. Recent Logs:"
curl -s -H "X-API-Key: $API_KEY" "$BASE_URL/logs/recent?limit=5" | jq .

echo -e "\n✅ All tests complete!"
```

Save as `test-production.sh` and run:
```bash
chmod +x test-production.sh
./test-production.sh
```

---

## 🔐 Security Best Practices

### 1. API Key Management

```bash
# Generate a stronger production API key
openssl rand -hex 32

# Update in Render:
API_KEYS=habit-tracker-key-dev,habit-tracker-key-prod-abc123xyz789
```

### 2. CORS Configuration

For production, restrict origins:

```bash
# Instead of ALLOWED_ORIGINS=*
ALLOWED_ORIGINS=https://yourhealthapp.com,https://app.yourhealthapp.com
```

### 3. Rate Limiting

Adjust based on your app usage:

```bash
# For high-volume apps:
RATE_LIMIT_REQUESTS_PER_MINUTE=5000
RATE_LIMIT_BURST=500

# For low-volume apps:
RATE_LIMIT_REQUESTS_PER_MINUTE=100
RATE_LIMIT_BURST=20
```

---

## 📈 Scaling Your Service

### Monitor Usage

```bash
# Check metrics endpoint
curl -H "X-API-Key: habit-tracker-key-dev" \
  https://logs.biopeak.authify.tech/api/v1/metrics
```

### When to Scale Up

- Response times > 500ms consistently
- CPU usage > 80%
- Memory usage > 90%
- Frequent cold starts affecting users

### Scaling Options

1. **Vertical**: Upgrade Render plan (more CPU/RAM)
2. **Horizontal**: Add load balancer + multiple instances
3. **Database**: Upgrade PostgreSQL plan
4. **Caching**: Add Redis for rate limiting (already supported)

---

## 🎉 Success Checklist

- [ ] Code pushed to GitHub
- [ ] PostgreSQL database created on Render
- [ ] Web service created and deployed
- [ ] Environment variables configured
- [ ] Health check passing
- [ ] API endpoints responding
- [ ] Custom domain configured
- [ ] DNS records added
- [ ] SSL certificate active
- [ ] Flutter app updated with new URL
- [ ] Test logs successfully ingested
- [ ] Monitoring and alerts set up

---

## 🆘 Need Help?

### Render Support

- Docs: https://render.com/docs
- Community: https://community.render.com
- Status: https://status.render.com

### Common Commands

```bash
# View recent logs
# Go to: Dashboard → Service → Logs

# Manual deploy
# Go to: Dashboard → Service → Manual Deploy

# Restart service
# Go to: Dashboard → Service → Settings → Restart

# View environment variables
# Go to: Dashboard → Service → Environment
```

---

## 💰 Cost Estimation

### Free Tier (Starter)
- Web Service: 750 hours/month (FREE)
- PostgreSQL: 90 days retention, 1GB (FREE)
- **Total**: $0/month (with limitations)

### Paid (Recommended for Production)
- Web Service (Starter): $7/month
- PostgreSQL (Starter): $7/month
- **Total**: $14/month
- Benefits: Always on, better performance, no cold starts

### Production (High Traffic)
- Web Service (Standard): $25/month
- PostgreSQL (Standard): $20/month
- **Total**: $45/month
- Benefits: 2GB RAM, 1 CPU, 20GB storage, high performance

---

## 🎯 Next Steps After Deployment

1. **Test thoroughly** from your Flutter app
2. **Monitor logs** for the first few days
3. **Set up alerts** for failures
4. **Generate a production API key** (stronger than dev key)
5. **Configure CORS** properly for your app domain
6. **Set up backup strategy** for database
7. **Document your API endpoints** for team members
8. **Consider adding authentication** for sensitive endpoints

---

**Your log server is now live at:**
🚀 **https://logs.biopeak.authify.tech/api/v1**

**API Key:** `habit-tracker-key-dev`

Happy logging! 📊

