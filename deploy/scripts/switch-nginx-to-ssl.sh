#!/bin/bash
# Switch nginx to SSL config when certs already exist (e.g. after deploy overwrote default.conf).
# Run from deploy directory on the server.

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
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

if [ -z "$DOMAIN_NAME" ]; then
    echo -e "${RED}[ERROR]${NC} DOMAIN_NAME not set in .env.production"
    exit 1
fi

if [ -z "$FRONTEND_DOMAIN" ]; then
    echo -e "${RED}[ERROR]${NC} FRONTEND_DOMAIN not set in .env.production (e.g. spark.vnhatng.com)"
    exit 1
fi

COMPOSE_CMD="docker compose --env-file .env.production -f docker-compose.prod.yml"

if [ ! -d "./nginx/conf.d" ]; then
    echo -e "${RED}[ERROR]${NC} nginx/conf.d not found"
    exit 1
fi

CERT_PATH="/etc/letsencrypt/live/$DOMAIN_NAME/fullchain.pem"
if ! $COMPOSE_CMD exec -T nginx test -f "$CERT_PATH" 2>/dev/null; then
    echo -e "${RED}[ERROR]${NC} SSL cert not found at $CERT_PATH"
    echo ""
    echo "Existing certs (if any):"
    $COMPOSE_CMD exec -T nginx ls -la /etc/letsencrypt/live/ 2>/dev/null || true
    echo ""
    echo "To obtain a cert for $DOMAIN_NAME, run: ./scripts/setup-ssl.sh"
    echo "  (Requires: DNS A record for $DOMAIN_NAME pointing to this server, ports 80 and 443 open.)"
    exit 1
fi

FRONTEND_CERT_PATH="/etc/letsencrypt/live/$FRONTEND_DOMAIN/fullchain.pem"
if ! $COMPOSE_CMD exec -T nginx test -f "$FRONTEND_CERT_PATH" 2>/dev/null; then
    echo -e "${RED}[ERROR]${NC} SSL cert not found at $FRONTEND_CERT_PATH"
    echo ""
    echo "Obtain a cert for $FRONTEND_DOMAIN (e.g. certbot certonly --webroot -w /var/www/certbot -d $FRONTEND_DOMAIN)"
    exit 1
fi

echo "Switching nginx to SSL config for $DOMAIN_NAME (API) and $FRONTEND_DOMAIN (frontend)..."
cp ./nginx/conf.d/default-ssl.conf ./nginx/conf.d/default.conf
sed -i.bak "s/your-domain.com/$DOMAIN_NAME/g" ./nginx/conf.d/default.conf
sed -i.bak "s/your-frontend-domain.com/$FRONTEND_DOMAIN/g" ./nginx/conf.d/default.conf

$COMPOSE_CMD restart nginx
sleep 2

echo -e "${GREEN}[OK]${NC} Nginx restarted with SSL config. Test: curl -k https://localhost/health"
