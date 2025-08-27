#!/bin/bash

# Start the complete analytics stack with Grafana visualization
# Usage: ./start-with-grafana.sh

set -e

echo "🚀 Starting Analytics Log Server with Grafana Visualization"
echo "=========================================================="

# Check if Docker and Docker Compose are available
if ! command -v docker &> /dev/null; then
    echo "❌ Docker is not installed. Please install Docker first."
    exit 1
fi

if ! command -v docker-compose &> /dev/null; then
    echo "❌ Docker Compose is not installed. Please install Docker Compose first."
    exit 1
fi

# Check if .env file exists, if not copy from example
if [ ! -f ".env" ]; then
    if [ -f "config.example.env" ]; then
        echo "📋 Copying example environment configuration..."
        cp config.example.env .env
        echo "✅ Created .env file from config.example.env"
        echo "⚠️  Please review and update .env file with your settings"
    else
        echo "⚠️  No .env file found. Please create one or copy from config.example.env"
    fi
fi

# Create necessary directories
echo "📁 Creating necessary directories..."
mkdir -p grafana/datasources grafana/dashboards grafana/provisioning

# Stop any running containers
echo "🛑 Stopping existing containers..."
docker-compose down

# Pull latest images
echo "📥 Pulling latest Docker images..."
docker-compose pull

# Start the stack
echo "🚀 Starting the complete stack..."
docker-compose up -d

# Wait for services to be ready
echo "⏳ Waiting for services to start..."
sleep 10

# Check service health
echo "🏥 Checking service health..."

# Check PostgreSQL
if docker exec log-ingestion-db pg_isready -U postgres > /dev/null 2>&1; then
    echo "✅ PostgreSQL is ready"
else
    echo "❌ PostgreSQL is not ready"
fi

# Check Redis
if docker exec log-ingestion-redis redis-cli ping > /dev/null 2>&1; then
    echo "✅ Redis is ready"
else
    echo "❌ Redis is not ready"
fi

# Check Log Server
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    echo "✅ Log Server is ready"
else
    echo "⏳ Log Server is starting... (this may take a moment)"
    sleep 5
    if curl -f http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ Log Server is now ready"
    else
        echo "❌ Log Server is not responding"
    fi
fi

# Check Prometheus
if curl -f http://localhost:9090/-/healthy > /dev/null 2>&1; then
    echo "✅ Prometheus is ready"
else
    echo "❌ Prometheus is not ready"
fi

# Check Grafana
if curl -f http://localhost:3000/api/health > /dev/null 2>&1; then
    echo "✅ Grafana is ready"
else
    echo "⏳ Grafana is starting... (this may take a moment)"
    sleep 10
    if curl -f http://localhost:3000/api/health > /dev/null 2>&1; then
        echo "✅ Grafana is now ready"
    else
        echo "❌ Grafana is not responding"
    fi
fi

echo ""
echo "🎉 Analytics Stack is now running!"
echo "=================================="
echo ""
echo "📊 Access your services:"
echo "   • Log Server:     http://localhost:8080"
echo "   • Grafana:        http://localhost:3000 (admin/admin)"
echo "   • Prometheus:     http://localhost:9090"
echo "   • PostgreSQL:     localhost:5432 (postgres/password123)"
echo "   • Redis:          localhost:6379"
echo ""
echo "📈 Pre-configured Grafana Dashboards:"
echo "   • Analytics Overview: http://localhost:3000/d/analytics-overview"
echo "   • Server Metrics:     http://localhost:3000/d/server-metrics"
echo ""
echo "🔧 Useful commands:"
echo "   • View logs:          docker-compose logs -f [service_name]"
echo "   • Stop stack:         docker-compose down"
echo "   • Restart service:    docker-compose restart [service_name]"
echo ""
echo "📚 Documentation:"
echo "   • API Filtering:      cat API_FILTERING.md"
echo "   • Grafana Setup:      cat GRAFANA_SETUP.md"
echo "   • Main README:        cat README.md"
echo ""
echo "🧪 Test the API:"
echo "   • Run filter tests:   ./scripts/test-filter-api.sh"
echo "   • Run basic tests:    ./scripts/test-api.sh"
echo ""

# Show container status
echo "🐳 Container Status:"
docker-compose ps

echo ""
echo "✨ Happy monitoring! Your analytics platform is ready to use."
