#!/bin/bash

# Backup Script for Tinder Clone Application
# Creates backups of database, redis data, and uploaded files

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
TIMESTAMP=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="$BACKUP_ROOT/$TIMESTAMP"

echo "=========================================="
echo "   Backup Script"
echo "=========================================="
echo ""

mkdir -p "$BACKUP_DIR"

print_info "Backup directory: $BACKUP_DIR"

cd "$DEPLOY_DIR"

if ! docker ps | grep -q "tinder_"; then
    print_error "Application containers are not running"
    exit 1
fi

print_info "Backing up PostgreSQL database..."
docker compose -f docker-compose.prod.yml exec -T postgres pg_dump -U postgres tinder_clone | gzip > "$BACKUP_DIR/postgres_dump.sql.gz"

print_info "Backing up PostgreSQL data volume..."
docker run --rm -v tinder_postgres_data:/data -v "$BACKUP_DIR":/backup alpine tar czf /backup/postgres_data.tar.gz -C /data .

print_info "Backing up Redis data..."
docker run --rm -v tinder_redis_data:/data -v "$BACKUP_DIR":/backup alpine tar czf /backup/redis_data.tar.gz -C /data .

print_info "Backing up uploaded files..."
docker run --rm -v tinder_uploads_data:/data -v "$BACKUP_DIR":/backup alpine tar czf /backup/uploads_data.tar.gz -C /data .

print_info "Backing up configuration files..."
cp .env.production "$BACKUP_DIR/env.production.backup"

BACKUP_SIZE=$(du -sh "$BACKUP_DIR" | cut -f1)
print_info "Backup completed! Size: $BACKUP_SIZE"
print_info "Location: $BACKUP_DIR"

print_info "Cleaning up old backups (keeping last 7 days)..."
find "$BACKUP_ROOT" -maxdepth 1 -type d -mtime +7 -exec rm -rf {} \; 2>/dev/null || true

echo ""
print_info "Backup contents:"
ls -lh "$BACKUP_DIR"
echo ""
print_info "To restore from this backup, run: ./restore.sh $TIMESTAMP"
echo ""
