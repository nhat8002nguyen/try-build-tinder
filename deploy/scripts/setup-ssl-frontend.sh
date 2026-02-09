#!/bin/bash
# Obtain Let's Encrypt certificate for the frontend domain (e.g. spark.vnhatng.com).
# Run from the deploy directory on the server. Nginx must be running (port 80).
# Uses 'docker compose run' so certbot does not need to be a running service.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
cd "$DEPLOY_DIR"

if [ ! -f ".env.production" ]; then
    echo -e "${RED}[ERROR]${NC} .env.production not found. Run from deploy directory."
    exit 1
fi

set -a
source .env.production
set +a

if [ -z "$FRONTEND_DOMAIN" ]; then
    echo -e "${RED}[ERROR]${NC} FRONTEND_DOMAIN not set in .env.production (e.g. spark.vnhatng.com)"
    exit 1
fi

if [ -z "$EMAIL_FOR_SSL" ]; then
    echo -e "${RED}[ERROR]${NC} EMAIL_FOR_SSL not set in .env.production"
    exit 1
fi

COMPOSE_CMD="docker compose --env-file .env.production -f docker-compose.prod.yml"

echo -e "${GREEN}[INFO]${NC} Obtaining SSL certificate for $FRONTEND_DOMAIN..."
echo -e "${YELLOW}[TIP]${NC} Ensure DNS A record for $FRONTEND_DOMAIN points to this server and port 80 is open."
echo ""

$COMPOSE_CMD run --rm --entrypoint certbot certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email "$EMAIL_FOR_SSL" \
    --agree-tos \
    --no-eff-email \
    -d "$FRONTEND_DOMAIN"

echo ""
echo -e "${GREEN}[OK]${NC} Certificate obtained for $FRONTEND_DOMAIN"
echo "Run ./scripts/switch-nginx-to-ssl.sh to enable HTTPS for both API and frontend (if not already)."
