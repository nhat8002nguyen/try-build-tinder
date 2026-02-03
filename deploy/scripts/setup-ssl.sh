#!/bin/bash

# SSL Certificate Setup Script using Let's Encrypt
# This script sets up SSL/TLS certificates for your domain

set -e

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

cd "$DEPLOY_DIR"

echo "=========================================="
echo "   SSL Certificate Setup (Let's Encrypt)"
echo "=========================================="
echo ""

if [ ! -f ".env.production" ]; then
    print_error ".env.production file not found!"
    exit 1
fi

set -a
source .env.production
set +a

if [ -z "$DOMAIN_NAME" ] || [ "$DOMAIN_NAME" == "your-domain.com" ]; then
    print_error "DOMAIN_NAME is not set in .env.production"
    print_error "Please set your domain name before running this script"
    exit 1
fi

if [ -z "$EMAIL_FOR_SSL" ] || [ "$EMAIL_FOR_SSL" == "your-email@example.com" ]; then
    print_error "EMAIL_FOR_SSL is not set in .env.production"
    print_error "Please set your email address for SSL notifications"
    exit 1
fi

print_info "Domain: $DOMAIN_NAME"
print_info "Email: $EMAIL_FOR_SSL"
echo ""

print_warning "Before continuing, please ensure:"
echo "  1. Your DNS A record points to this server's IP"
echo "  2. Port 80 and 443 are open in AWS Security Group (Inbound: HTTPS 443 from 0.0.0.0/0)"
echo "  3. Port 443 is allowed in OS firewall (UFW) on this server"
echo "  4. Nginx is running with HTTP configuration"
echo ""
read -p "Have you completed the above steps? (yes/no): " confirm

if [ "$confirm" != "yes" ]; then
    print_warning "Please complete the prerequisites and run this script again"
    exit 0
fi

print_step "1/4 - Testing domain resolution..."
RESOLVED_IP=$(dig +short "$DOMAIN_NAME" | tail -n1)
SERVER_IP=$(curl -s ifconfig.me)

if [ -z "$RESOLVED_IP" ]; then
    print_error "Domain $DOMAIN_NAME does not resolve to any IP"
    print_error "Please configure your DNS settings first"
    exit 1
fi

print_info "Domain resolves to: $RESOLVED_IP"
print_info "Server IP is: $SERVER_IP"

if [ "$RESOLVED_IP" != "$SERVER_IP" ]; then
    print_warning "Domain IP ($RESOLVED_IP) doesn't match server IP ($SERVER_IP)"
    print_warning "SSL certificate may fail. Continue anyway? (yes/no)"
    read -p "> " continue_anyway
    if [ "$continue_anyway" != "yes" ]; then
        exit 0
    fi
fi

COMPOSE_CMD="docker compose --env-file .env.production -f docker-compose.prod.yml"

print_step "2/4 - Obtaining SSL certificate from Let's Encrypt..."
$COMPOSE_CMD run --rm --entrypoint certbot certbot certonly \
    --webroot \
    --webroot-path=/var/www/certbot \
    --email "$EMAIL_FOR_SSL" \
    --agree-tos \
    --no-eff-email \
    -d "$DOMAIN_NAME"

if [ $? -ne 0 ]; then
    print_error "Failed to obtain SSL certificate"
    print_error "Please check the error messages above"
    exit 1
fi

print_step "3/4 - Switching nginx to SSL configuration..."
cp ./nginx/conf.d/default-ssl.conf ./nginx/conf.d/default.conf
sed -i "s/your-domain.com/$DOMAIN_NAME/g" ./nginx/conf.d/default.conf

print_step "4/4 - Reloading nginx with SSL configuration..."
$COMPOSE_CMD restart nginx
sleep 3

echo ""
echo "=========================================="
print_info "SSL Setup Complete!"
echo "=========================================="
echo ""
print_info "Testing HTTPS endpoint..."
sleep 3

if curl -f --connect-timeout 10 https://"$DOMAIN_NAME"/health > /dev/null 2>&1; then
    print_info "✓ HTTPS is working correctly!"
else
    print_warning "HTTPS test failed or timed out."
    echo ""
    print_step "Troubleshooting (if curl hangs = port 443 likely blocked):"
    echo "  1. AWS Security Group: EC2 → Security Groups → your instance's SG → Inbound rules"
    echo "     Add: Type=HTTPS, Port=443, Source=0.0.0.0/0 (or your IP)"
    echo "  2. On this server, check UFW: sudo ufw status | grep 443"
    echo "     If 443 is denied, run: sudo ufw allow 443/tcp && sudo ufw reload"
    echo "  3. Check nginx is listening: $COMPOSE_CMD exec nginx ss -tlnp | grep 443"
    echo "  4. Check nginx logs: $COMPOSE_CMD logs nginx"
    echo ""
    print_info "Run ./scripts/troubleshoot-https.sh on the server for full diagnostics."
fi

echo ""
print_info "Your application is now accessible at: https://$DOMAIN_NAME"
echo ""
print_info "SSL Certificate Details:"
$COMPOSE_CMD run --rm --entrypoint certbot certbot certificates
echo ""
print_info "Certificate auto-renewal is enabled"
print_info "Certificates will be automatically renewed every 12 hours"
echo ""
