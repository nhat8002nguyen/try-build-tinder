#!/bin/bash

# Restore Script for Tinder Clone Application
# Restores database, redis data, and uploaded files from a backup

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
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

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
BACKUP_ROOT="/opt/tinder-app/backups"

echo "=========================================="
echo "   Restore Script"
echo "=========================================="
echo ""

if [ -z "$1" ]; then
    print_error "Usage: ./restore.sh <backup_timestamp>"
    print_info "Available backups:"
    ls -1 "$BACKUP_ROOT" 2>/dev/null || echo "  No backups found"
    exit 1
fi

BACKUP_TIMESTAMP="$1"
BACKUP_DIR="$BACKUP_ROOT/$BACKUP_TIMESTAMP"

if [ ! -d "$BACKUP_DIR" ]; then
    print_error "Backup directory not found: $BACKUP_DIR"
    print_info "Available backups:"
    ls -1 "$BACKUP_ROOT"
    exit 1
fi

print_info "Backup directory: $BACKUP_DIR"
print_warning "This will overwrite current data. Are you sure? (yes/no)"
read -p "> " confirm

if [ "$confirm" != "yes" ]; then
    print_info "Restore cancelled"
    exit 0
fi

cd "$DEPLOY_DIR"

print_info "Stopping application services..."
docker compose -f docker-compose.prod.yml stop backend nginx

print_info "Restoring PostgreSQL data..."
if [ -f "$BACKUP_DIR/postgres_data.tar.gz" ]; then
    docker run --rm -v tinder_postgres_data:/data -v "$BACKUP_DIR":/backup alpine sh -c "rm -rf /data/* && tar xzf /backup/postgres_data.tar.gz -C /data"
    print_info "✓ PostgreSQL data restored"
fi

if [ -f "$BACKUP_DIR/postgres_dump.sql.gz" ]; then
    print_info "Restoring PostgreSQL database from SQL dump..."
    gunzip < "$BACKUP_DIR/postgres_dump.sql.gz" | docker compose -f docker-compose.prod.yml exec -T postgres psql -U postgres tinder_clone
    print_info "✓ PostgreSQL database restored"
fi

print_info "Restoring Redis data..."
if [ -f "$BACKUP_DIR/redis_data.tar.gz" ]; then
    docker run --rm -v tinder_redis_data:/data -v "$BACKUP_DIR":/backup alpine sh -c "rm -rf /data/* && tar xzf /backup/redis_data.tar.gz -C /data"
    print_info "✓ Redis data restored"
fi

print_info "Restoring uploaded files..."
if [ -f "$BACKUP_DIR/uploads_data.tar.gz" ]; then
    docker run --rm -v tinder_uploads_data:/data -v "$BACKUP_DIR":/backup alpine sh -c "rm -rf /data/* && tar xzf /backup/uploads_data.tar.gz -C /data"
    print_info "✓ Uploaded files restored"
fi

print_info "Starting services..."
docker compose -f docker-compose.prod.yml start postgres redis
sleep 5
docker compose -f docker-compose.prod.yml start backend nginx

print_info "Waiting for services to be ready..."
sleep 10

echo ""
print_info "Restore completed!"
print_info "Verifying services..."
docker compose -f docker-compose.prod.yml ps
echo ""
