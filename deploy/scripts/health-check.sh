#!/bin/bash

# Health Check Script
# Returns 0 if all services are healthy, 1 otherwise
# Useful for monitoring tools and cron jobs

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"

cd "$DEPLOY_DIR"

EXIT_CODE=0

# Check if containers are running
if ! docker ps | grep -q "tinder_nginx"; then
    echo "ERROR: Nginx container is not running"
    EXIT_CODE=1
fi

if ! docker ps | grep -q "tinder_backend"; then
    echo "ERROR: Backend container is not running"
    EXIT_CODE=1
fi

if ! docker ps | grep -q "tinder_postgres"; then
    echo "ERROR: PostgreSQL container is not running"
    EXIT_CODE=1
fi

if ! docker ps | grep -q "tinder_redis"; then
    echo "ERROR: Redis container is not running"
    EXIT_CODE=1
fi

# Check health endpoint
if ! curl -f -s http://localhost/health > /dev/null 2>&1; then
    echo "ERROR: Health endpoint is not responding"
    EXIT_CODE=1
fi

# Check PostgreSQL
if ! docker exec tinder_postgres pg_isready -U postgres > /dev/null 2>&1; then
    echo "ERROR: PostgreSQL is not ready"
    EXIT_CODE=1
fi

# Check Redis
if ! docker exec tinder_redis redis-cli ping > /dev/null 2>&1; then
    echo "ERROR: Redis is not responding"
    EXIT_CODE=1
fi

if [ $EXIT_CODE -eq 0 ]; then
    echo "OK: All services are healthy"
fi

exit $EXIT_CODE
