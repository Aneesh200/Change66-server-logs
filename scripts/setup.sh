#!/bin/bash

# Log Ingestion Server Setup Script

set -e

echo "🚀 Setting up Log Ingestion Server..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ Go is not installed. Please install Go 1.21+ first.${NC}"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+')
echo -e "${GREEN}✅ Found Go version: $GO_VERSION${NC}"

# Check if Docker is installed
if command -v docker &> /dev/null; then
    echo -e "${GREEN}✅ Docker found${NC}"
    DOCKER_AVAILABLE=true
else
    echo -e "${YELLOW}⚠️  Docker not found. Some features may not be available.${NC}"
    DOCKER_AVAILABLE=false
fi

# Check if PostgreSQL is available
if command -v psql &> /dev/null; then
    echo -e "${GREEN}✅ PostgreSQL client found${NC}"
    POSTGRES_CLIENT=true
else
    echo -e "${YELLOW}⚠️  PostgreSQL client not found${NC}"
    POSTGRES_CLIENT=false
fi

# Create .env file if it doesn't exist
if [ ! -f .env ]; then
    echo -e "${YELLOW}📝 Creating .env file from template...${NC}"
    cp config.example.env .env
    echo -e "${GREEN}✅ Created .env file. Please edit it with your configuration.${NC}"
else
    echo -e "${GREEN}✅ .env file already exists${NC}"
fi

# Install Go dependencies
echo -e "${YELLOW}📦 Installing Go dependencies...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✅ Dependencies installed${NC}"

# Install development tools
echo -e "${YELLOW}🔧 Installing development tools...${NC}"
if command -v golangci-lint &> /dev/null; then
    echo -e "${GREEN}✅ golangci-lint already installed${NC}"
else
    echo "Installing golangci-lint..."
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
fi

if command -v goimports &> /dev/null; then
    echo -e "${GREEN}✅ goimports already installed${NC}"
else
    echo "Installing goimports..."
    go install golang.org/x/tools/cmd/goimports@latest
fi

if command -v migrate &> /dev/null; then
    echo -e "${GREEN}✅ migrate tool already installed${NC}"
else
    echo "Installing migrate tool..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

echo -e "${GREEN}✅ Development tools installed${NC}"

# Build the application
echo -e "${YELLOW}🏗️  Building application...${NC}"
make build
echo -e "${GREEN}✅ Application built successfully${NC}"

# Setup database (if Docker is available)
if [ "$DOCKER_AVAILABLE" = true ]; then
    echo -e "${YELLOW}🐳 Setting up database with Docker...${NC}"
    
    # Check if docker-compose is available
    if command -v docker-compose &> /dev/null; then
        COMPOSE_CMD="docker-compose"
    elif command -v docker &> /dev/null && docker compose version &> /dev/null; then
        COMPOSE_CMD="docker compose"
    else
        echo -e "${RED}❌ Docker Compose not found${NC}"
        DOCKER_AVAILABLE=false
    fi
    
    if [ "$DOCKER_AVAILABLE" = true ]; then
        # Start only the database
        $COMPOSE_CMD up -d postgres
        echo -e "${GREEN}✅ Database started with Docker${NC}"
        
        # Wait for database to be ready
        echo -e "${YELLOW}⏳ Waiting for database to be ready...${NC}"
        sleep 10
        
        # Run migrations
        echo -e "${YELLOW}🔄 Running database migrations...${NC}"
        if [ "$POSTGRES_CLIENT" = true ]; then
            export DATABASE_URL="postgres://postgres:password123@localhost:5432/analytics_logs?sslmode=disable"
            make migrate-up
            echo -e "${GREEN}✅ Database migrations completed${NC}"
        else
            echo -e "${YELLOW}⚠️  PostgreSQL client not available. Please run migrations manually.${NC}"
        fi
    fi
fi

# Generate API keys
echo -e "${YELLOW}🔑 Generating API keys...${NC}"
API_KEY_1="hta_$(date +%s)_$(openssl rand -hex 16)"
API_KEY_2="hta_$(date +%s)_$(openssl rand -hex 16)"

echo -e "${GREEN}✅ Generated API keys:${NC}"
echo -e "${GREEN}   Development: $API_KEY_1${NC}"
echo -e "${GREEN}   Production:  $API_KEY_2${NC}"

# Update .env file with generated keys
if command -v sed &> /dev/null; then
    sed -i.bak "s/API_KEYS=.*/API_KEYS=$API_KEY_1,$API_KEY_2/" .env
    rm .env.bak
    echo -e "${GREEN}✅ Updated .env file with new API keys${NC}"
fi

echo ""
echo -e "${GREEN}🎉 Setup completed successfully!${NC}"
echo ""
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Review and edit .env file if needed"
echo "2. Start the server: make run"
echo "3. Test the server: curl -f http://localhost:8080/health"
echo "4. View API documentation in README.md"
echo ""
echo -e "${YELLOW}Your API keys:${NC}"
echo -e "${GREEN}Development: $API_KEY_1${NC}"
echo -e "${GREEN}Production:  $API_KEY_2${NC}"
echo ""
echo -e "${YELLOW}Keep these keys secure and update your Flutter app configuration!${NC}"
