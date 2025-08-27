# Deploying to Render

This guide will help you deploy your Log Ingestion Server to Render as a web service.

## 🚀 Quick Deploy

### 1. Fork/Clone Repository
```bash
git clone <your-repo-url>
cd Change66-Log-Server
```

### 2. Push to GitHub
```bash
git add .
git commit -m "Production ready for Render deployment"
git push origin main
```

## 📋 Render Setup

### 1. Create New Web Service
- Go to [Render Dashboard](https://dashboard.render.com/)
- Click "New +" → "Web Service"
- Connect your GitHub repository

### 2. Configure Service
- **Name**: `log-ingestion-server` (or your preferred name)
- **Environment**: `Docker`
- **Region**: Choose closest to your users
- **Branch**: `main`
- **Root Directory**: Leave empty (root of repo)

### 3. Environment Variables
Set these in Render dashboard:

```bash
# Required Database Variables
DB_HOST=your-render-postgres-host.render.com
DB_NAME=your_database_name
DB_USER=your_database_user
DB_PASSWORD=your_database_password
DB_SSL_MODE=require

# Required API Keys
API_KEYS=your-production-api-key-1,your-production-api-key-2

# Optional (have defaults)
PORT=8080
LOG_LEVEL=info
ENABLE_METRICS=true
ENABLE_CORS=true
ALLOWED_ORIGINS=*
```

### 4. Build Command
```bash
docker build -t log-ingestion-server .
```

### 5. Start Command
```bash
docker run -p $PORT:8080 --env-file .env log-ingestion-server
```

## 🗄️ Database Setup

### 1. Create PostgreSQL Database
- In Render dashboard: "New +" → "PostgreSQL"
- Choose plan (Free tier available)
- Note down connection details

### 2. Run Migrations
Your database migrations will run automatically on startup, but you can also run them manually:

```bash
# Connect to your Render PostgreSQL
psql "postgresql://user:password@host:port/database"

# Or use the migration files
# The server will handle this automatically
```

## 🔑 API Key Generation

Generate secure API keys for production:

```bash
# Generate a secure random key
openssl rand -hex 32

# Example output: a1b2c3d4e5f6...
# Use this as your API_KEY value
```

## 📊 Health Check

Render will automatically check:
- **Health Endpoint**: `GET /health`
- **Expected Response**: `{"status":"healthy"}`

## 🚀 Deployment Steps

### 1. Automatic Deploy
- Render will automatically build and deploy on every push to main
- Monitor build logs for any issues

### 2. Manual Deploy
- Go to your service in Render dashboard
- Click "Manual Deploy" → "Deploy latest commit"

### 3. Verify Deployment
```bash
# Test health endpoint
curl https://your-service-name.onrender.com/health

# Test API with key
curl -H "X-API-Key: your-api-key" \
  https://your-service-name.onrender.com/api/v1/status
```

## 🔧 Environment Variables Reference

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| `DB_HOST` | ✅ | - | Render PostgreSQL host |
| `DB_NAME` | ✅ | - | Database name |
| `DB_USER` | ✅ | - | Database username |
| `DB_PASSWORD` | ✅ | - | Database password |
| `DB_SSL_MODE` | ❌ | `require` | SSL mode for database |
| `API_KEYS` | ✅ | - | Comma-separated API keys |
| `PORT` | ❌ | `8080` | Server port |
| `LOG_LEVEL` | ❌ | `info` | Logging level |
| `ENABLE_METRICS` | ❌ | `true` | Enable Prometheus metrics |
| `ENABLE_CORS` | ❌ | `true` | Enable CORS |
| `ALLOWED_ORIGINS` | ❌ | `*` | CORS allowed origins |

## 📈 Monitoring & Scaling

### 1. Auto-scaling
- Render automatically scales based on traffic
- Free tier: 750 hours/month
- Paid plans: Unlimited scaling

### 2. Logs
- View logs in Render dashboard
- Real-time log streaming
- Log retention based on plan

### 3. Metrics
- Built-in Render metrics
- Custom Prometheus metrics at `/metrics`
- Health check monitoring

## 🚨 Troubleshooting

### Common Issues

1. **Build Fails**
   - Check Dockerfile syntax
   - Verify Go version compatibility
   - Check build logs in Render

2. **Database Connection Fails**
   - Verify environment variables
   - Check PostgreSQL service status
   - Verify SSL requirements

3. **Health Check Fails**
   - Check application logs
   - Verify port configuration
   - Check health endpoint implementation

### Debug Commands

```bash
# Check service status
curl -v https://your-service.onrender.com/health

# Test database connection
curl -H "X-API-Key: your-key" \
  https://your-service.onrender.com/api/v1/status

# View logs in Render dashboard
# Go to your service → Logs tab
```

## 🔒 Security Considerations

1. **API Keys**: Use strong, randomly generated keys
2. **Database**: Enable SSL connections
3. **CORS**: Restrict origins in production
4. **Rate Limiting**: Configure appropriate limits
5. **Logging**: Set appropriate log levels

## 📚 Next Steps

After successful deployment:

1. **Test API endpoints** with your production API keys
2. **Monitor performance** using Render metrics
3. **Set up alerts** for health check failures
4. **Configure custom domain** if needed
5. **Set up CI/CD** for automatic deployments

## 🎯 Production Checklist

- [ ] Environment variables configured
- [ ] Database migrations run
- [ ] API keys generated and configured
- [ ] Health checks passing
- [ ] SSL certificates working
- [ ] CORS properly configured
- [ ] Rate limiting enabled
- [ ] Monitoring set up
- [ ] Logs accessible
- [ ] Backup strategy in place

Your Log Ingestion Server is now ready for production deployment on Render! 🚀
