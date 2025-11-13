# GCP Deployment - Quick Start Guide

Deploy your log server to Google Cloud Platform with **HTTP support** (no forced HTTPS redirects) and **automatic cost optimization**.

---

## 🚀 One-Command Deployment

```bash
cd /Users/aneesh/Change66-Log-Server
./scripts/deploy-to-gcp.sh
```

That's it! The script will handle everything automatically.

---

## 📋 What the Script Does

### 1. **Creates VM with Build Configuration**
   - Machine Type: `e2-standard-2` (2 vCPUs, 8GB RAM)
   - Purpose: Fast Docker image building

### 2. **Installs & Builds**
   - Installs Docker
   - Clones your repository
   - Builds Docker image
   - Creates systemd service
   - Starts your server

### 3. **Downgrades to Runtime Configuration** 🎯
   - Machine Type: `e2-micro` (0.25-1 vCPU, 1GB RAM)
   - **Cost: FREE TIER!** ✅
   - Saves ~$15/month vs keeping build machine

### 4. **Configures Networking**
   - Opens HTTP ports (80, 8080)
   - No HTTPS redirect enforced
   - Your production app works instantly!

---

## 🔧 Prerequisites

### Step 1: Install Google Cloud SDK

**Mac (Homebrew):**
```bash
brew install --cask google-cloud-sdk
```

**Or download from:**
https://cloud.google.com/sdk/docs/install

### Step 2: Login to GCP

```bash
gcloud auth login
```

### Step 3: Create/Select Project

```bash
# Create new project
gcloud projects create my-log-server --name="Log Server"

# Or use existing project
gcloud config set project YOUR_PROJECT_ID
```

### Step 4: Enable Billing

- Go to: https://console.cloud.google.com/billing
- Link your project to a billing account
- **Note**: e2-micro is FREE within free tier limits!

### Step 5: Enable APIs

```bash
gcloud services enable compute.googleapis.com
```

---

## 🎯 Deploy!

```bash
cd /Users/aneesh/Change66-Log-Server
./scripts/deploy-to-gcp.sh
```

The script will:
- ✅ Create VM (build config)
- ✅ Install dependencies
- ✅ Build Docker image (~3-5 minutes)
- ✅ Start service
- ✅ Test HTTP endpoint
- ✅ Downgrade to free tier
- ✅ Verify still working

---

## 🌐 After Deployment

### Update Your DNS

Script will show you the IP address. Add this DNS record:

```
Type: A
Name: logs.biopeak.authify
Value: [IP from script output]
TTL: 3600
```

**DNS Providers:**
- **Cloudflare**: Dashboard → DNS → Add Record
- **AWS Route 53**: Hosted zones → Create record
- **Namecheap**: Advanced DNS → Add record
- **GoDaddy**: DNS Management → Add record

### Test Your Endpoint

```bash
# Replace with your IP
export SERVER_IP="your.gcp.ip"

# Test HTTP (no redirect!)
curl http://$SERVER_IP/health

# Test with domain (after DNS)
curl http://logs.biopeak.authify.tech/health

# Test API
curl -H "X-API-Key: habit-tracker-key-dev" \
  http://logs.biopeak.authify.tech/api/v1/status
```

---

## 💰 Cost Breakdown

### Free Tier (Your Configuration)
```
VM (e2-micro):        FREE (within 744 hrs/month)
Disk (20GB):          $0.40/month
Network egress:       $0.01/GB (1GB free/month)
───────────────────────────────────────────
Total:                ~$0.40 - $2/month
```

### If You Used Build Machine All Time
```
VM (e2-standard-2):   $48/month
Disk (20GB):          $0.40/month
Network egress:       $0.01/GB
───────────────────────────────────────────
Total:                ~$48/month

SAVINGS: ~$46/month! 🎉
```

---

## 🔧 Useful Commands

### View Logs
```bash
gcloud compute ssh log-server --zone=us-central1-a --command="sudo journalctl -u log-server -f"
```

### Restart Service
```bash
gcloud compute ssh log-server --zone=us-central1-a --command="sudo systemctl restart log-server"
```

### SSH into Server
```bash
gcloud compute ssh log-server --zone=us-central1-a
```

### Check Service Status
```bash
gcloud compute ssh log-server --zone=us-central1-a --command="sudo systemctl status log-server"
```

### Stop Instance (to save even more)
```bash
gcloud compute instances stop log-server --zone=us-central1-a
```

### Start Instance
```bash
gcloud compute instances start log-server --zone=us-central1-a
```

### Delete Everything (cleanup)
```bash
# Delete instance
gcloud compute instances delete log-server --zone=us-central1-a

# Delete firewall rule
gcloud compute firewall-rules delete allow-log-server-http
```

---

## 🐛 Troubleshooting

### Build Fails

```bash
# SSH into server
gcloud compute ssh log-server --zone=us-central1-a

# Check Docker status
sudo systemctl status docker

# Check build logs
cd ~/Change66-Log-Server
sudo docker build -t log-server:latest .
```

### Service Won't Start

```bash
# View logs
gcloud compute ssh log-server --zone=us-central1-a --command="sudo journalctl -u log-server -n 100"

# Check if Docker container is running
gcloud compute ssh log-server --zone=us-central1-a --command="sudo docker ps -a"

# Restart service
gcloud compute ssh log-server --zone=us-central1-a --command="sudo systemctl restart log-server"
```

### Can't Connect to Server

```bash
# Check firewall rules
gcloud compute firewall-rules list

# Check if instance is running
gcloud compute instances list

# Get external IP
gcloud compute instances describe log-server --zone=us-central1-a --format='get(networkInterfaces[0].accessConfigs[0].natIP)'
```

### Database Connection Issues

```bash
# SSH into server
gcloud compute ssh log-server --zone=us-central1-a

# Check environment file
cat ~/Change66-Log-Server/.env

# Test database connection
nc -zv dpg-d493bsmmcj7s73e9c5pg-a.oregon-postgres.render.com 5432
```

---

## 🎯 Why GCP?

| Feature | GCP | Render | Railway |
|---------|-----|--------|---------|
| **HTTP Support** | ✅ Yes | ❌ Forces HTTPS | ❌ Forces HTTPS |
| **Free Tier** | ✅ e2-micro | ❌ No | ❌ $5 credit |
| **Full Control** | ✅ Complete | ❌ Limited | ❌ Limited |
| **Cost** | ~$0.40/month | $7/month min | $5-10/month |
| **Flexibility** | ✅ High | ⚠️ Medium | ⚠️ Medium |

---

## ✅ Post-Deployment Checklist

- [ ] Script ran successfully
- [ ] Health check returns 200 (not 307!)
- [ ] External IP noted
- [ ] DNS A record added
- [ ] DNS propagated (10-60 minutes)
- [ ] HTTP endpoint tested with domain
- [ ] Flutter app updated with new URL
- [ ] Logs flowing from production app
- [ ] Instance downgraded to e2-micro

---

## 📱 Update Your Flutter App

After DNS propagates:

```dart
// lib/services/analytics_service.dart
class AnalyticsService {
  // Use HTTP - it works on GCP!
  static const String LOG_SERVER_URL = 'http://logs.biopeak.authify.tech/api/v1';
  static const String API_KEY = 'habit-tracker-key-dev';
  ...
}
```

---

## 🎉 Success!

Your log server is now:
- ✅ **Running on GCP**
- ✅ **HTTP working** (no redirects!)
- ✅ **Cost-optimized** (FREE tier!)
- ✅ **Production-ready**
- ✅ **Your app can connect** immediately!

**Total deployment time**: ~10 minutes
**Monthly cost**: ~$0.40 - $2
**HTTP working**: ✅ YES!

---

## 🆘 Need Help?

**View deployment info:**
```bash
cat deployment-info.txt
```

**GCP Console:**
https://console.cloud.google.com/compute/instances

**Check server status anytime:**
```bash
curl http://logs.biopeak.authify.tech/health
```

---

**Happy logging! 🚀📊**


