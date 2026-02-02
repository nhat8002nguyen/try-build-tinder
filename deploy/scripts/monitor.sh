#!/bin/bash

# Monitoring Script for Tinder Clone Application
# Displays real-time status and metrics

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"

cd "$DEPLOY_DIR"

clear
echo "=========================================="
echo "   Application Monitoring Dashboard"
echo "=========================================="
echo ""

echo -e "${BLUE}=== Container Status ===${NC}"
docker compose -f docker-compose.prod.yml ps
echo ""

echo -e "${BLUE}=== Resource Usage ===${NC}"
docker stats --no-stream --format "table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}"
echo ""

echo -e "${BLUE}=== Disk Usage ===${NC}"
echo "Docker Volumes:"
docker system df -v | grep -A 20 "Local Volumes" | head -n 10
echo ""

echo -e "${BLUE}=== Health Checks ===${NC}"
if curl -f http://localhost/health > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Backend health: OK"
else
    echo -e "${RED}✗${NC} Backend health: FAILED"
fi

if docker exec tinder_postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} PostgreSQL: OK"
else
    echo -e "${RED}✗${NC} PostgreSQL: FAILED"
fi

if docker exec tinder_redis redis-cli ping > /dev/null 2>&1; then
    echo -e "${GREEN}✓${NC} Redis: OK"
else
    echo -e "${RED}✗${NC} Redis: FAILED"
fi
echo ""

echo -e "${BLUE}=== Recent Logs ===${NC}"
echo "Backend (last 10 lines):"
docker compose -f docker-compose.prod.yml logs --tail=10 backend
echo ""

echo -e "${BLUE}=== Nginx Access (last 5 requests) ===${NC}"
docker exec tinder_nginx tail -n 5 /var/log/nginx/access.log 2>/dev/null || echo "No logs available"
echo ""

echo -e "${BLUE}=== Nginx Errors (last 5) ===${NC}"
docker exec tinder_nginx tail -n 5 /var/log/nginx/error.log 2>/dev/null || echo "No errors"
echo ""

echo "=========================================="
echo "Useful commands:"
echo "  Follow logs: docker compose -f docker-compose.prod.yml logs -f"
echo "  Restart app: docker compose -f docker-compose.prod.yml restart"
echo "  Shell access: docker compose -f docker-compose.prod.yml exec backend sh"
echo "=========================================="
