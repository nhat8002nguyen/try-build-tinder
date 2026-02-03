#!/bin/bash
# Run this script on the EC2 instance when curl https://your-domain.com hangs.
# A hang usually means port 443 is not reachable (Security Group or OS firewall).

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_ok() { echo -e "${GREEN}[OK]${NC} $1"; }
print_fail() { echo -e "${RED}[FAIL]${NC} $1"; }
print_warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
print_step() { echo -e "${BLUE}[CHECK]${NC} $1"; }

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"
cd "$DEPLOY_DIR"

echo "=========================================="
echo "   HTTPS connectivity troubleshooting"
echo "=========================================="
echo ""

if [ ! -f ".env.production" ]; then
    print_fail ".env.production not found. Run from deploy dir with .env.production."
    exit 1
fi

set -a
source .env.production
set +a

DOMAIN="${DOMAIN_NAME:-your-domain.com}"
COMPOSE_CMD="docker compose --env-file .env.production -f docker-compose.prod.yml"

print_step "1. UFW status (port 443 must be ALLOW)"
if command -v ufw >/dev/null 2>&1; then
    if sudo ufw status 2>/dev/null | grep -q "443.*ALLOW"; then
        print_ok "UFW allows 443"
    else
        print_fail "UFW does not show 443 ALLOW. Run: sudo ufw allow 443/tcp && sudo ufw reload"
    fi
else
    print_warn "UFW not installed; skipping"
fi

print_step "2. Nginx container and port 443"
if $COMPOSE_CMD ps nginx 2>/dev/null | grep -q "Up"; then
    print_ok "Nginx container is running"
    if $COMPOSE_CMD exec -T nginx ss -tlnp 2>/dev/null | grep -q ":443"; then
        print_ok "Nginx is listening on 443"
    else
        print_fail "Nginx is not listening on 443. Check nginx config (default.conf should have listen 443 ssl)"
    fi
else
    print_fail "Nginx container is not running. Run: $COMPOSE_CMD up -d nginx"
fi

print_step "3. Host listening on 443"
if ss -tlnp 2>/dev/null | grep -q ":443 " || netstat -tlnp 2>/dev/null | grep -q ":443 "; then
    print_ok "Something is listening on 443 on the host"
else
    print_fail "Nothing is listening on 443 on the host (Docker may not be publishing 443)"
fi

print_step "4. Local HTTPS (from this server)"
if curl -f -s -o /dev/null --connect-timeout 5 -k "https://localhost/health" 2>/dev/null \
   || curl -f -s -o /dev/null --connect-timeout 5 -k "https://127.0.0.1/health" 2>/dev/null; then
    print_ok "Local https://localhost/health works (nginx and SSL are fine on this host)"
else
    print_fail "Local HTTPS failed; check nginx and SSL config"
fi

echo ""
echo "If local HTTPS works but curl from outside hangs:"
echo "  → Port 443 is blocked before reaching this server."
echo "  → Fix: AWS EC2 → Security Groups → select the SG attached to this instance"
echo "  → Inbound rules: Add rule Type=HTTPS, Port=443, Source=0.0.0.0/0"
echo ""
echo "Verify from your laptop:"
echo "  curl -v --connect-timeout 5 https://$DOMAIN/health"
echo ""
