#!/bin/bash

# Script to set up systemd service for automatic startup

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

if [ "$EUID" -ne 0 ]; then 
    echo "Please run with sudo"
    exit 1
fi

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DEPLOY_DIR="$(dirname "$SCRIPT_DIR")"

print_info "Installing systemd service for automatic startup..."

# Copy service file
cp "$DEPLOY_DIR/systemd/tinder-app.service" /etc/systemd/system/

# Reload systemd
systemctl daemon-reload

# Enable service
systemctl enable tinder-app.service

print_info "Systemd service installed and enabled"
print_info ""
print_info "Commands:"
echo "  Start:   sudo systemctl start tinder-app"
echo "  Stop:    sudo systemctl stop tinder-app"
echo "  Status:  sudo systemctl status tinder-app"
echo "  Disable: sudo systemctl disable tinder-app"
echo ""
print_info "The application will now start automatically on system boot"
