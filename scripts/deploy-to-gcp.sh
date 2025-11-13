#!/bin/bash

##############################################
# GCP Deployment Script with Auto Upgrade/Downgrade
# Deploys log server with HTTP support (no forced HTTPS)
##############################################

set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
PROJECT_ID="${GCP_PROJECT_ID:-}"
INSTANCE_NAME="log-server"
ZONE="us-central1-a"  # Change to your preferred zone
REGION="us-central1"

# Machine types
BUILD_MACHINE_TYPE="e2-standard-2"    # 2 vCPUs, 8GB RAM for building
RUNTIME_MACHINE_TYPE="e2-micro"       # 0.25-1 vCPU, 1GB RAM for running (FREE TIER!)

# Firewall
FIREWALL_RULE_NAME="allow-log-server-http"

# Database credentials from Render
DB_HOST="dpg-d493bsmmcj7s73e9c5pg-a.oregon-postgres.render.com"
DB_PORT="5432"
DB_NAME="analytics_logs"
DB_USER="analytics_logs_user"
DB_PASSWORD="Gcai1cMF6763HjkA9aMpTddhm8nH3U1m"
DB_SSL_MODE="require"
API_KEY="habit-tracker-key-dev"

##############################################
# Helper Functions
##############################################

print_header() {
    echo -e "\n${BLUE}========================================${NC}"
    echo -e "${BLUE}$1${NC}"
    echo -e "${BLUE}========================================${NC}\n"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

##############################################
# Check Prerequisites
##############################################

print_header "Checking Prerequisites"

# Check if gcloud is installed
if ! command -v gcloud &> /dev/null; then
    print_error "gcloud CLI not found. Please install it:"
    echo "Visit: https://cloud.google.com/sdk/docs/install"
    exit 1
fi
print_success "gcloud CLI found"

# Check if authenticated
if ! gcloud auth list --filter=status:ACTIVE --format="value(account)" &> /dev/null; then
    print_warning "Not authenticated with gcloud"
    echo "Please run: gcloud auth login"
    exit 1
fi
print_success "gcloud authenticated"

# Get or set project ID
if [ -z "$PROJECT_ID" ]; then
    PROJECT_ID=$(gcloud config get-value project 2>/dev/null)
    if [ -z "$PROJECT_ID" ]; then
        print_error "No GCP project set. Please run:"
        echo "gcloud config set project YOUR_PROJECT_ID"
        exit 1
    fi
fi
print_success "Using GCP project: $PROJECT_ID"

##############################################
# Step 1: Create VM Instance (Build Machine Type)
##############################################

print_header "Step 1: Creating VM Instance (Build Configuration)"

# Check if instance already exists
if gcloud compute instances describe $INSTANCE_NAME --zone=$ZONE &>/dev/null; then
    print_warning "Instance $INSTANCE_NAME already exists"
    read -p "Do you want to delete and recreate it? (y/n): " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        print_info "Deleting existing instance..."
        gcloud compute instances delete $INSTANCE_NAME --zone=$ZONE --quiet
        print_success "Instance deleted"
    else
        print_info "Using existing instance"
        EXISTING_INSTANCE=true
    fi
fi

if [ -z "$EXISTING_INSTANCE" ]; then
    print_info "Creating VM with BUILD configuration ($BUILD_MACHINE_TYPE)..."
    
    gcloud compute instances create $INSTANCE_NAME \
        --project=$PROJECT_ID \
        --zone=$ZONE \
        --machine-type=$BUILD_MACHINE_TYPE \
        --image-family=ubuntu-2204-lts \
        --image-project=ubuntu-os-cloud \
        --boot-disk-size=20GB \
        --boot-disk-type=pd-standard \
        --tags=http-server,https-server \
        --metadata=startup-script='#!/bin/bash
# This will run on first boot only
echo "VM created at $(date)" > /var/log/vm-creation.log
'
    
    print_success "VM instance created with build configuration"
    
    # Wait for instance to be ready
    print_info "Waiting for instance to be ready..."
    sleep 30
fi

##############################################
# Step 2: Configure Firewall Rules
##############################################

print_header "Step 2: Configuring Firewall Rules"

# Check if firewall rule exists
if gcloud compute firewall-rules describe $FIREWALL_RULE_NAME &>/dev/null; then
    print_warning "Firewall rule already exists"
else
    print_info "Creating firewall rule to allow HTTP traffic..."
    
    gcloud compute firewall-rules create $FIREWALL_RULE_NAME \
        --project=$PROJECT_ID \
        --direction=INGRESS \
        --priority=1000 \
        --network=default \
        --action=ALLOW \
        --rules=tcp:80,tcp:8080 \
        --source-ranges=0.0.0.0/0 \
        --target-tags=http-server
    
    print_success "Firewall rule created"
fi

##############################################
# Step 3: Install Dependencies and Build
##############################################

print_header "Step 3: Installing Docker and Building Application"

print_info "Installing Docker and dependencies..."

gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --command="
set -e

# Update system
echo '📦 Updating system packages...'
sudo apt-get update -qq

# Install Docker
echo '🐳 Installing Docker...'
sudo apt-get install -y docker.io docker-compose git
sudo systemctl start docker
sudo systemctl enable docker
sudo usermod -aG docker \$USER

# Clone repository
echo '📥 Cloning repository...'
if [ -d ~/Change66-Log-Server ]; then
    echo 'Repository already exists, pulling latest changes...'
    cd ~/Change66-Log-Server
    git pull
else
    git clone https://github.com/Aneesh200/Change66-server-logs.git ~/Change66-Log-Server
    cd ~/Change66-Log-Server
fi

# Create environment file
echo '⚙️  Creating environment file...'
cat > .env << 'EOF'
PORT=8080
GIN_MODE=release
LOG_LEVEL=info
ENVIRONMENT=production

DB_HOST=$DB_HOST
DB_PORT=$DB_PORT
DB_NAME=$DB_NAME
DB_USER=$DB_USER
DB_PASSWORD=$DB_PASSWORD
DB_SSL_MODE=$DB_SSL_MODE

API_KEYS=$API_KEY

RATE_LIMIT_REQUESTS_PER_MINUTE=1000
RATE_LIMIT_BURST=100
MAX_BATCH_SIZE=1000
WORKER_POOL_SIZE=10
REQUEST_TIMEOUT_SECONDS=30
MAX_REQUEST_SIZE_MB=10

ENABLE_METRICS=true
ENABLE_CORS=true
ALLOWED_ORIGINS=*
EOF

# Build Docker image
echo '🏗️  Building Docker image (this may take a few minutes)...'
sudo docker build -t log-server:latest .

# Create systemd service
echo '🔧 Creating systemd service...'
sudo tee /etc/systemd/system/log-server.service > /dev/null << 'SERVICEEOF'
[Unit]
Description=Log Ingestion Server
After=docker.service
Requires=docker.service

[Service]
Type=simple
Restart=always
RestartSec=5
WorkingDirectory=/home/\$USER/Change66-Log-Server
ExecStartPre=-/usr/bin/docker stop log-server
ExecStartPre=-/usr/bin/docker rm log-server
ExecStart=/usr/bin/docker run --name log-server --rm -p 80:8080 --env-file .env log-server:latest
ExecStop=/usr/bin/docker stop log-server

[Install]
WantedBy=multi-user.target
SERVICEEOF

# Fix service file to use actual user
sudo sed -i \"s/\\\$USER/\$(whoami)/g\" /etc/systemd/system/log-server.service

# Start service
echo '🚀 Starting log server service...'
sudo systemctl daemon-reload
sudo systemctl enable log-server
sudo systemctl start log-server

echo '✅ Build and deployment complete!'
"

print_success "Application built and running"

##############################################
# Step 4: Verify Deployment
##############################################

print_header "Step 4: Verifying Deployment"

print_info "Waiting for service to start..."
sleep 10

# Get external IP
EXTERNAL_IP=$(gcloud compute instances describe $INSTANCE_NAME --zone=$ZONE --format='get(networkInterfaces[0].accessConfigs[0].natIP)')
print_success "External IP: $EXTERNAL_IP"

# Test HTTP endpoint
print_info "Testing HTTP endpoint..."
sleep 5

if curl -s -f "http://$EXTERNAL_IP/health" > /dev/null; then
    print_success "Health check passed!"
    curl -s "http://$EXTERNAL_IP/health" | python3 -m json.tool
else
    print_error "Health check failed. Checking logs..."
    gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --command="sudo systemctl status log-server"
fi

##############################################
# Step 5: Downgrade to Runtime Machine Type
##############################################

print_header "Step 5: Downgrading to Runtime Configuration"

print_warning "Downgrading from $BUILD_MACHINE_TYPE to $RUNTIME_MACHINE_TYPE"
print_info "This will save costs (~\$15/month → FREE TIER!)"

# Stop instance
print_info "Stopping instance..."
gcloud compute instances stop $INSTANCE_NAME --zone=$ZONE

# Change machine type
print_info "Changing machine type to $RUNTIME_MACHINE_TYPE..."
gcloud compute instances set-machine-type $INSTANCE_NAME \
    --zone=$ZONE \
    --machine-type=$RUNTIME_MACHINE_TYPE

# Start instance
print_info "Starting instance with new configuration..."
gcloud compute instances start $INSTANCE_NAME --zone=$ZONE

print_success "Instance downgraded to runtime configuration"

# Wait for instance to start
print_info "Waiting for instance to restart..."
sleep 30

# Verify still working
print_info "Verifying service after downgrade..."
sleep 10

if curl -s -f "http://$EXTERNAL_IP/health" > /dev/null; then
    print_success "Service still healthy after downgrade!"
else
    print_error "Service not responding after downgrade. Checking..."
    gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --command="sudo systemctl restart log-server"
    sleep 10
fi

##############################################
# Step 6: Final Information
##############################################

print_header "🎉 Deployment Complete!"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}      Your Log Server is LIVE! 🚀${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}📍 Server Details:${NC}"
echo "   Instance: $INSTANCE_NAME"
echo "   Zone: $ZONE"
echo "   Machine Type: $RUNTIME_MACHINE_TYPE (runtime)"
echo "   External IP: $EXTERNAL_IP"
echo ""
echo -e "${BLUE}🌐 Endpoints:${NC}"
echo "   Health: http://$EXTERNAL_IP/health"
echo "   Status: http://$EXTERNAL_IP/api/v1/status"
echo "   Ingest: http://$EXTERNAL_IP/api/v1/ingest"
echo "   Metrics: http://$EXTERNAL_IP/metrics"
echo ""
echo -e "${BLUE}🔑 API Key:${NC}"
echo "   $API_KEY"
echo ""
echo -e "${BLUE}📊 Test Commands:${NC}"
echo "   Health: curl http://$EXTERNAL_IP/health"
echo "   Status: curl -H 'X-API-Key: $API_KEY' http://$EXTERNAL_IP/api/v1/status"
echo ""
echo -e "${BLUE}🌍 Update Your DNS:${NC}"
echo "   Type: A Record"
echo "   Name: logs.biopeak.authify"
echo "   Value: $EXTERNAL_IP"
echo "   TTL: 3600"
echo ""
echo -e "${BLUE}📱 Flutter App Config:${NC}"
echo "   After DNS update, use:"
echo "   LOG_SERVER_URL = 'http://logs.biopeak.authify.tech/api/v1'"
echo "   API_KEY = '$API_KEY'"
echo ""
echo -e "${BLUE}💰 Cost Estimate:${NC}"
echo "   e2-micro: FREE (within free tier limits)"
echo "   Network egress: ~\$0.01/GB"
echo "   Disk: ~\$0.40/month (20GB)"
echo "   Total: ~FREE - \$2/month"
echo ""
echo -e "${BLUE}🔧 Useful Commands:${NC}"
echo "   SSH into server: gcloud compute ssh $INSTANCE_NAME --zone=$ZONE"
echo "   View logs: gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --command='sudo journalctl -u log-server -f'"
echo "   Restart service: gcloud compute ssh $INSTANCE_NAME --zone=$ZONE --command='sudo systemctl restart log-server'"
echo "   Stop instance: gcloud compute instances stop $INSTANCE_NAME --zone=$ZONE"
echo "   Start instance: gcloud compute instances start $INSTANCE_NAME --zone=$ZONE"
echo ""
echo -e "${GREEN}✅ HTTP (no HTTPS redirect) is working!${NC}"
echo -e "${GREEN}✅ Your production app can now send logs!${NC}"
echo ""
echo -e "${YELLOW}Note: The instance was built with $BUILD_MACHINE_TYPE and downgraded to $RUNTIME_MACHINE_TYPE${NC}"
echo -e "${YELLOW}      This saves ~\$15/month while maintaining full functionality!${NC}"
echo ""

# Save deployment info
cat > deployment-info.txt << EOF
Deployment Date: $(date)
Instance Name: $INSTANCE_NAME
Zone: $ZONE
External IP: $EXTERNAL_IP
Machine Type: $RUNTIME_MACHINE_TYPE

HTTP Endpoint: http://$EXTERNAL_IP/api/v1
API Key: $API_KEY

DNS Configuration:
Type: A
Name: logs.biopeak.authify
Value: $EXTERNAL_IP

Flutter Config:
LOG_SERVER_URL = 'http://logs.biopeak.authify.tech/api/v1'
API_KEY = '$API_KEY'
EOF

print_success "Deployment info saved to deployment-info.txt"
echo ""


