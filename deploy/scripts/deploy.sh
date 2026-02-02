#!/bin/bash

# Production Deployment Script for Tinder Clone Backend
# This script deploys the application to a production environment

set -e

# Color codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
PROJECT_DIR="$(dirname "$DEPLOY_DIR")"

cd "$DEPLOY_DIR"

echo "=========================================="
echo "   Tinder Clone - Production Deployment"
echo "=========================================="
echo ""

if [ ! -f ".env.production" ]; then
    print_error ".env.production file not found!"
    print_info "Please create .env.production from .env.production.example"
    exit 1
fi

COMPOSE_CMD="docker compose --env-file .env.production -f docker-compose.prod.yml"

set -a
source .env.production
set +a

if [ -z "$POSTGRES_PASSWORD" ] || [ -z "$JWT_SECRET" ] || [ -z "$REDIS_PASSWORD" ]; then
    print_error "Required environment variables are not set!"
    print_error "Please check your .env.production file"
    exit 1
fi

print_step "1/8 - Checking Docker installation..."
if ! command -v docker &> /dev/null; then
    print_error "Docker is not installed. Please run setup-ec2.sh first"
    exit 1
fi

if ! command -v docker compose &> /dev/null; then
    print_error "Docker Compose is not installed. Please run setup-ec2.sh first"
    exit 1
fi

print_info "Docker and Docker Compose are installed"

print_step "2/8 - Creating backup of current deployment (if exists)..."
if [ "$(docker ps -q -f name=tinder_)" ]; then
    BACKUP_DIR="/opt/tinder-app/backups/$(date +%Y%m%d_%H%M%S)"
    mkdir -p "$BACKUP_DIR"
    
    print_info "Backing up volumes..."
    docker run --rm -v tinder_postgres_data:/data -v "$BACKUP_DIR":/backup alpine tar czf /backup/postgres_data.tar.gz -C /data . || true
    docker run --rm -v tinder_uploads_data:/data -v "$BACKUP_DIR":/backup alpine tar czf /backup/uploads_data.tar.gz -C /data . || true
    
    print_info "Backup created at: $BACKUP_DIR"
fi

print_step "3/8 - Stopping existing containers..."
$COMPOSE_CMD down || true

print_step "4/8 - Pruning old Docker resources..."
docker system prune -f --volumes=false || true

print_step "5/8 - Building Docker images..."
$COMPOSE_CMD build --no-cache

print_step "6/8 - Starting services..."
$COMPOSE_CMD up -d

print_step "7/8 - Waiting for services to be healthy..."
sleep 10

RETRIES=30
for i in $(seq 1 $RETRIES); do
    if $COMPOSE_CMD ps | grep -q "healthy"; then
        print_info "Services are starting up... ($i/$RETRIES)"
        sleep 2
    else
        break
    fi
done

print_step "8/8 - Verifying deployment..."
echo ""

container_status() {
    docker ps -a --filter "name=$1" --format "{{.Status}}" | head -n1
}

if container_status "tinder_nginx" | grep -q "^Up"; then
    print_info "✓ Nginx is running"
else
    print_error "✗ Nginx is not running"
fi

if container_status "tinder_backend" | grep -q "^Up"; then
    print_info "✓ Backend is running"
else
    print_error "✗ Backend is not running"
fi

if container_status "tinder_postgres" | grep -q "^Up"; then
    print_info "✓ PostgreSQL is running"
else
    print_error "✗ PostgreSQL is not running"
fi

if container_status "tinder_redis" | grep -q "^Up"; then
    print_info "✓ Redis is running"
else
    print_error "✗ Redis is not running"
fi

echo ""
print_info "Testing health endpoint..."
sleep 5

if curl -f http://localhost/health > /dev/null 2>&1; then
    print_info "✓ Health check passed"
else
    print_warning "Health check failed - services may still be starting"
fi

echo ""
echo "=========================================="
print_info "Deployment Complete!"
echo "=========================================="
echo ""
print_info "Service Status:"
$COMPOSE_CMD ps
echo ""
print_info "Useful Commands (run from deploy/ directory):"
echo "  - View logs: $COMPOSE_CMD logs -f"
echo "  - View backend logs: $COMPOSE_CMD logs -f backend"
echo "  - View nginx logs: $COMPOSE_CMD logs -f nginx"
echo "  - Stop services: $COMPOSE_CMD down"
echo "  - Restart services: $COMPOSE_CMD restart"
echo ""

if [ -z "$DOMAIN_NAME" ] || [ "$DOMAIN_NAME" == "your-domain.com" ]; then
    print_warning "SSL is not configured yet. Run ./setup-ssl.sh after configuring your domain"
    print_info "Current access: http://$(curl -s ifconfig.me)"
else
    print_info "Access your application at: https://$DOMAIN_NAME"
fi

echo ""
