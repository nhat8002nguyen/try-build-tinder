#!/bin/bash

# Production EC2 Setup Script for Tinder Clone Backend
# This script sets up a fresh Ubuntu EC2 instance with Docker, Docker Compose, and required dependencies

set -e

echo "=========================================="
echo "Tinder Clone Backend - EC2 Setup Script"
echo "=========================================="

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [ "$EUID" -ne 0 ]; then 
    print_error "Please run as root or with sudo"
    exit 1
fi

print_info "Updating system packages..."
apt-get update
apt-get upgrade -y

print_info "Installing required packages..."
apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    gnupg \
    lsb-release \
    software-properties-common \
    ufw \
    git \
    htop \
    vim

print_info "Installing Docker..."
# Add Docker's official GPG key
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /etc/apt/keyrings/docker.gpg
chmod a+r /etc/apt/keyrings/docker.gpg

# Add Docker repository
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list > /dev/null

# Install Docker Engine
apt-get update
apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

print_info "Starting and enabling Docker service..."
systemctl start docker
systemctl enable docker

print_info "Installing Docker Compose standalone (if needed)..."
DOCKER_COMPOSE_VERSION=$(curl -s https://api.github.com/repos/docker/compose/releases/latest | grep 'tag_name' | cut -d\" -f4)
curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
chmod +x /usr/local/bin/docker-compose

print_info "Configuring firewall (UFW)..."
ufw --force enable
ufw default deny incoming
ufw default allow outgoing
ufw allow ssh
ufw allow 80/tcp
ufw allow 443/tcp
ufw status

print_info "Creating application directory..."
mkdir -p /opt/tinder-app
chown -R $SUDO_USER:$SUDO_USER /opt/tinder-app

print_info "Configuring Docker log rotation..."
cat > /etc/docker/daemon.json <<EOF
{
  "log-driver": "json-file",
  "log-opts": {
    "max-size": "10m",
    "max-file": "3"
  }
}
EOF

systemctl restart docker

print_info "Setting up automatic security updates..."
apt-get install -y unattended-upgrades
dpkg-reconfigure -plow unattended-upgrades

print_info "Optimizing system settings for production..."
# Increase file descriptor limits
cat >> /etc/security/limits.conf <<EOF
* soft nofile 65536
* hard nofile 65536
EOF

# Sysctl optimizations
cat >> /etc/sysctl.conf <<EOF
# Network performance tuning
net.core.somaxconn = 65536
net.ipv4.tcp_max_syn_backlog = 8192
net.ipv4.ip_local_port_range = 1024 65535
net.ipv4.tcp_tw_reuse = 1
net.ipv4.tcp_fin_timeout = 30

# Security
net.ipv4.conf.all.rp_filter = 1
net.ipv4.conf.default.rp_filter = 1
net.ipv4.icmp_echo_ignore_broadcasts = 1
net.ipv4.conf.all.accept_source_route = 0
net.ipv4.conf.default.accept_source_route = 0
EOF

sysctl -p

print_info "Creating docker group and adding user..."
groupadd -f docker
usermod -aG docker $SUDO_USER || true

print_info "Setting up log rotation for application logs..."
cat > /etc/logrotate.d/tinder-app <<EOF
/opt/tinder-app/logs/*.log {
    daily
    missingok
    rotate 14
    compress
    delaycompress
    notifempty
    create 0640 $SUDO_USER $SUDO_USER
    sharedscripts
}
EOF

echo ""
echo "=========================================="
print_info "EC2 Setup Complete!"
echo "=========================================="
echo ""
print_info "Next steps:"
echo "1. Logout and login again (or run: newgrp docker)"
echo "2. Upload your application code to /opt/tinder-app"
echo "3. Configure your .env.production file"
echo "4. Run the deployment script: ./deploy.sh"
echo ""
print_warning "Remember to:"
echo "  - Configure your DNS A record to point to this EC2 instance"
echo "  - Update security group to allow ports 80 and 443"
echo "  - Set up SSL certificates using Certbot"
echo ""
