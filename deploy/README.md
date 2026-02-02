# Production Deployment Guide - Tinder Clone Backend

This directory contains production-grade deployment configurations and scripts for deploying the Tinder Clone backend to AWS EC2 with Docker, nginx, and SSL/TLS support.

## 📋 Table of Contents

- [Prerequisites](#prerequisites)
- [Architecture](#architecture)
- [Initial Setup](#initial-setup)
- [Deployment Steps](#deployment-steps)
- [SSL/TLS Configuration](#ssltls-configuration)
- [Monitoring & Maintenance](#monitoring--maintenance)
- [Troubleshooting](#troubleshooting)

## 🎯 Prerequisites

### AWS Requirements
- AWS EC2 instance (recommended: t3.medium or larger)
- Ubuntu 22.04 LTS or 20.04 LTS
- At least 20GB storage
- Security group allowing ports: 22 (SSH), 80 (HTTP), 443 (HTTPS)
- Elastic IP (recommended for production)

### Domain & DNS
- Domain name (for SSL/HTTPS)
- DNS A record pointing to your EC2 instance IP

### Local Requirements
- SSH access to EC2 instance
- Git installed
- Basic knowledge of Linux/Unix commands

## 🏗 Architecture

```
Internet
   ↓
[AWS EC2 Instance]
   ↓
[Nginx] (Reverse Proxy + Rate Limiting + SSL/TLS)
   ↓
[Backend API] (Go/Gin) ← WebSocket support
   ↓
[PostgreSQL] + [Redis]
```

### Components

- **Nginx**: Reverse proxy, rate limiting, SSL termination
- **Backend**: Go application (Gin framework)
- **PostgreSQL**: Primary database
- **Redis**: Caching and session storage
- **Certbot**: SSL certificate management (Let's Encrypt)

### Rate Limiting Configuration

- **Authentication endpoints** (`/api/auth/login`, `/api/auth/register`): 5 requests/minute
- **API endpoints** (`/api/*`): 30 requests/second
- **WebSocket** (`/ws`): 5 connections burst, 10 max connections per IP
- **General**: 10 requests/second

## 🚀 Initial Setup

### Step 1: Launch EC2 Instance

1. Launch Ubuntu 22.04 LTS instance
2. Configure security group:
   ```
   Inbound Rules:
   - SSH (22): Your IP
   - HTTP (80): 0.0.0.0/0
   - HTTPS (443): 0.0.0.0/0
   ```
3. Allocate and associate Elastic IP (recommended)

### Step 2: SSH into Instance

```bash
ssh -i your-key.pem ubuntu@your-ec2-ip
```

### Step 3: Run EC2 Setup Script

```bash
# Upload the setup script
scp -i your-key.pem deploy/scripts/setup-ec2.sh ubuntu@your-ec2-ip:~/

# SSH into instance
ssh -i your-key.pem ubuntu@your-ec2-ip

# Run setup
sudo bash setup-ec2.sh

# Logout and login again for docker group to take effect
exit
ssh -i your-key.pem ubuntu@your-ec2-ip
```

This script will:
- Update system packages
- Install Docker and Docker Compose
- Configure firewall (UFW)
- Set up security and performance optimizations
- Create application directory at `/opt/tinder-app`

## 📦 Deployment Steps

### Step 1: Upload Application Code

From your local machine:

```bash
# Navigate to project root
cd /Users/nhatnguyen/Workspaces/try-build-tinder

# Create a tarball (excluding unnecessary files)
# On macOS: COPYFILE_DISABLE=1 avoids adding extended attributes that cause "unknown extended header" warnings on Linux
COPYFILE_DISABLE=1 tar --exclude='node_modules' \
    --exclude='frontend' \
    --exclude='.git' \
    --exclude='backend/uploads' \
    -czf tinder-backend.tar.gz backend/ deploy/

# Upload to EC2
scp -i your-key.pem tinder-backend.tar.gz ubuntu@your-ec2-ip:/opt/tinder-app/

# SSH into instance
ssh -i your-key.pem ubuntu@your-ec2-ip

# Extract files (--warning=no-unknown-keyword suppresses harmless macOS xattr messages on Linux)
cd /opt/tinder-app
tar --warning=no-unknown-keyword -xzf tinder-backend.tar.gz
```

Alternatively, use Git:

```bash
ssh -i your-key.pem ubuntu@your-ec2-ip
cd /opt/tinder-app
git clone <your-repo-url> .
```

### Step 2: Configure Environment Variables

```bash
cd /opt/tinder-app/deploy

# Copy and edit production environment file
cp .env.production.example .env.production
nano .env.production
```

**Important**: Update these values:

```bash
# Strong passwords (use: openssl rand -base64 32)
POSTGRES_PASSWORD=<strong-random-password>
REDIS_PASSWORD=<strong-random-password>
JWT_SECRET=<strong-random-secret-min-32-chars>

# Domain configuration (for SSL)
DOMAIN_NAME=api.yourdomain.com
EMAIL_FOR_SSL=you@example.com

# Optional: OAuth credentials
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
```

### Step 3: Deploy Application

```bash
cd /opt/tinder-app/deploy/scripts

# Make scripts executable
chmod +x *.sh

# Deploy
./deploy.sh
```

The deployment script will:
1. Check Docker installation
2. Create backup of existing deployment
3. Stop old containers
4. Build new Docker images
5. Start all services
6. Verify deployment

### Step 4: Verify Deployment

```bash
# Check container status
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml ps

# Test health endpoint
curl http://localhost/health

# Should return: {"status":"healthy"}
```

## 🔒 SSL/TLS Configuration

### Prerequisites

1. Domain name configured (A record pointing to EC2 IP)
2. Application deployed and running
3. Ports 80 and 443 open in security group

### Step 1: Verify DNS

```bash
# Check if domain resolves to your server
dig +short api.yourdomain.com

# Should return your EC2 public IP
```

### Step 2: Run SSL Setup Script

```bash
cd /opt/tinder-app/deploy/scripts
./setup-ssl.sh
```

The script will:
1. Verify DNS configuration
2. Switch nginx to HTTP-only mode
3. Obtain SSL certificate from Let's Encrypt
4. Update nginx to use SSL
5. Configure automatic renewal

### Step 3: Verify HTTPS

```bash
# Test HTTPS endpoint
curl https://api.yourdomain.com/health

# Check certificate
curl -vI https://api.yourdomain.com 2>&1 | grep -A 10 "Server certificate"
```

## 📊 Monitoring & Maintenance

### View Application Status

```bash
cd /opt/tinder-app/deploy/scripts
./monitor.sh
```

This shows:
- Container status
- Resource usage (CPU, Memory, Network)
- Health checks
- Recent logs

### View Logs

```bash
cd /opt/tinder-app/deploy

# All logs
docker compose -f docker-compose.prod.yml logs -f

# Backend only
docker compose -f docker-compose.prod.yml logs -f backend

# Nginx only
docker compose -f docker-compose.prod.yml logs -f nginx

# Last 100 lines
docker compose -f docker-compose.prod.yml logs --tail=100
```

### Backup

```bash
cd /opt/tinder-app/deploy/scripts
./backup.sh
```

Backups are stored in `/opt/tinder-app/backups/` and include:
- PostgreSQL database dump
- Redis data
- Uploaded files
- Configuration files

Old backups (>7 days) are automatically removed.

### Restore

```bash
cd /opt/tinder-app/deploy/scripts

# List available backups
ls /opt/tinder-app/backups/

# Restore from specific backup
./restore.sh 20260125_120000
```

### Update Application

```bash
cd /opt/tinder-app

# Pull latest code
git pull origin main

# Redeploy
cd deploy/scripts
./deploy.sh
```

### Restart Services

```bash
cd /opt/tinder-app/deploy

# Restart all
docker compose -f docker-compose.prod.yml restart

# Restart specific service
docker compose -f docker-compose.prod.yml restart backend
docker compose -f docker-compose.prod.yml restart nginx
```

## 🧪 Testing with Postman

### Without HTTPS (Initial Testing)

1. In Postman, use: `http://your-ec2-ip/api/auth/register`
2. Headers:
   ```
   Content-Type: application/json
   ```
3. Body:
   ```json
   {
     "email": "test@example.com",
     "password": "password123",
     "name": "Test User"
   }
   ```

### With HTTPS (After SSL Setup)

1. In Postman, use: `https://api.yourdomain.com/api/auth/register`
2. Postman will automatically handle SSL verification
3. Save the token from login response for authenticated requests

### Environment Variables in Postman

Create an environment with:
```
base_url: https://api.yourdomain.com
token: <your-jwt-token>
```

Use `{{base_url}}/api/...` in requests and `{{token}}` in Authorization header.

## 🔧 Troubleshooting

### "tar: Ignoring unknown extended header keyword" when extracting

These messages appear when extracting a tarball created on macOS (it includes extended attributes like `com.apple.provenance` that Linux tar does not use). They are harmless and extraction still succeeds.

- **When creating the tarball on macOS**: Use `COPYFILE_DISABLE=1 tar ...` so the archive does not include macOS extended attributes.
- **When extracting on EC2**: Use `tar --warning=no-unknown-keyword -xzf tinder-backend.tar.gz` to suppress the messages.

### Services Not Starting

```bash
cd /opt/tinder-app/try-build-tinder/deploy   # or your deploy path

# Check logs (use --env-file so variables are set)
docker compose --env-file .env.production -f docker-compose.prod.yml logs -f

# Check container status
docker ps -a

# Restart services
docker compose --env-file .env.production -f docker-compose.prod.yml down
docker compose --env-file .env.production -f docker-compose.prod.yml up -d
```

### SSL Certificate Issues

```bash
# Test certificate manually
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml run --rm certbot certificates

# Renew certificate manually
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml run --rm certbot renew

# Restart nginx
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart nginx
```

### Database Connection Issues

```bash
# From deploy/ directory, always pass env file for variable substitution
cd /opt/tinder-app/try-build-tinder/deploy
docker compose --env-file .env.production -f docker-compose.prod.yml logs postgres

# Connect to database
docker compose --env-file .env.production -f docker-compose.prod.yml exec postgres psql -U postgres -d tinder_clone

# Check connections
docker compose --env-file .env.production -f docker-compose.prod.yml exec postgres psql -U postgres -c "SELECT count(*) FROM pg_stat_activity;"
```

### "password authentication failed for user postgres" / Variables not set

If you see `POSTGRES_PASSWORD variable is not set` or backend fails with "password authentication failed for user postgres":

1. **Use the env file with Compose**  
   Always run docker compose from the `deploy/` directory and pass the env file:
   ```bash
   cd /opt/tinder-app/try-build-tinder/deploy
   docker compose --env-file .env.production -f docker-compose.prod.yml up -d
   ```
   Use the same `--env-file .env.production` for logs, exec, down, etc.

2. **If you changed `POSTGRES_PASSWORD` after the first deploy**  
   Postgres initializes the data volume only on first run. The existing volume still has the old password. Either:
   - **Reset and re-deploy** (erases DB data):
     ```bash
     cd /opt/tinder-app/try-build-tinder/deploy
     docker compose --env-file .env.production -f docker-compose.prod.yml down
     docker volume rm deploy_postgres_data
     ./scripts/deploy.sh
     ```
   - Or change the postgres user password inside the running container to match your new `.env.production` (see Postgres docs).

3. **Passwords with special characters**  
   In `.env.production`, put values in double quotes if they contain `$`, `#`, or spaces, e.g. `POSTGRES_PASSWORD="my$ecure#pass"`.

### Rate Limiting Too Aggressive

Edit `/opt/tinder-app/deploy/nginx/nginx.conf` and adjust rate limits:

```nginx
limit_req_zone $binary_remote_addr zone=api:10m rate=30r/s;  # Increase rate
```

Then restart nginx:
```bash
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart nginx
```

### High Memory Usage

```bash
# Check memory usage
free -h
docker stats --no-stream

# Restart containers to free memory
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart
```

### Nginx 502 Bad Gateway

```bash
# Check if backend is running
docker ps | grep tinder_backend

# Check backend logs
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs backend

# Test backend directly
curl http://localhost:8080/health
```

## 📝 Additional Configuration

### Custom Domain (Subdomain)

If using a subdomain like `api.yourdomain.com`:

1. Create DNS A record: `api.yourdomain.com` → EC2 IP
2. Update `.env.production`:
   ```bash
   DOMAIN_NAME=api.yourdomain.com
   ```
3. Run SSL setup: `./setup-ssl.sh`

### Multiple Domains

To support multiple domains, edit nginx config:

```nginx
server_name api.yourdomain.com api2.yourdomain.com;
```

Then obtain certificates for both:
```bash
docker compose run --rm certbot certonly ... -d api.yourdomain.com -d api2.yourdomain.com
```

### Enable S3 Storage

Update `.env.production`:

```bash
STORAGE_TYPE=s3
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
S3_BUCKET_NAME=your-bucket-name
```

Redeploy:
```bash
./deploy.sh
```

## 🛡 Security Best Practices

1. **Firewall**: Use UFW or AWS Security Groups to restrict access
2. **SSH**: Use key-based authentication, disable password login
3. **Updates**: Regularly update system packages and Docker images
4. **Secrets**: Never commit `.env.production` to version control
5. **Backups**: Run daily backups and test restores regularly
6. **Monitoring**: Set up CloudWatch or external monitoring
7. **Logs**: Regularly review logs for suspicious activity

## 📞 Support

For issues or questions:
1. Check logs: `docker compose logs`
2. Review troubleshooting section above
3. Check nginx error logs: `docker exec tinder_nginx cat /var/log/nginx/error.log`

---

**Last Updated**: January 2026
