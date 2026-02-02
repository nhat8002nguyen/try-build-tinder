# Quick Deployment Reference

## 🚀 Initial Setup (One-time)

### 1. Launch EC2 Instance
- Ubuntu 22.04 LTS
- t3.medium or larger
- Open ports: 22, 80, 443
- Assign Elastic IP

### 2. Run Setup Script
```bash
scp -i key.pem deploy/scripts/setup-ec2.sh ubuntu@EC2_IP:~/
ssh -i key.pem ubuntu@EC2_IP
sudo bash setup-ec2.sh
```

### 3. Upload Code
```bash
# From local machine
cd /Users/nhatnguyen/Workspaces/try-build-tinder
tar --exclude='node_modules' --exclude='frontend' --exclude='.git' -czf app.tar.gz backend/ deploy/
scp -i key.pem app.tar.gz ubuntu@EC2_IP:/opt/tinder-app/

# On EC2
ssh -i key.pem ubuntu@EC2_IP
cd /opt/tinder-app && tar -xzf app.tar.gz
```

### 4. Configure Environment
```bash
cd /opt/tinder-app/deploy
cp .env.production.example .env.production
nano .env.production
```

**Must change**:
- `POSTGRES_PASSWORD`
- `REDIS_PASSWORD`
- `JWT_SECRET`
- `DOMAIN_NAME` (for SSL)
- `EMAIL_FOR_SSL` (for SSL)

### 5. Deploy
```bash
cd /opt/tinder-app/deploy/scripts
chmod +x *.sh
./deploy.sh
```

### 6. Setup SSL (optional but recommended)
```bash
# After DNS is configured
./setup-ssl.sh
```

## 📋 Common Commands

### Service Management
```bash
cd /opt/tinder-app/deploy

# Start services
docker compose -f docker-compose.prod.yml up -d

# Stop services
docker compose -f docker-compose.prod.yml down

# Restart services
docker compose -f docker-compose.prod.yml restart

# Restart specific service
docker compose -f docker-compose.prod.yml restart backend
```

### Logs
```bash
# All logs (follow)
docker compose -f docker-compose.prod.yml logs -f

# Backend only
docker compose -f docker-compose.prod.yml logs -f backend

# Last 100 lines
docker compose -f docker-compose.prod.yml logs --tail=100

# Nginx access logs
docker exec tinder_nginx tail -f /var/log/nginx/access.log

# Nginx error logs
docker exec tinder_nginx tail -f /var/log/nginx/error.log
```

### Monitoring
```bash
# Status dashboard
cd /opt/tinder-app/deploy/scripts && ./monitor.sh

# Container status
docker compose -f docker-compose.prod.yml ps

# Resource usage
docker stats

# Health check
curl http://localhost/health
curl https://your-domain.com/health
```

### Backup & Restore
```bash
# Create backup
cd /opt/tinder-app/deploy/scripts && ./backup.sh

# List backups
ls /opt/tinder-app/backups/

# Restore
./restore.sh 20260125_120000
```

### Database Access
```bash
# PostgreSQL shell
docker compose -f docker-compose.prod.yml exec postgres psql -U postgres -d tinder_clone

# Redis CLI
docker compose -f docker-compose.prod.yml exec redis redis-cli
```

### Updates
```bash
# Pull latest code
cd /opt/tinder-app && git pull

# Rebuild and deploy
cd deploy/scripts && ./deploy.sh
```

## 🧪 Testing Endpoints

### Health Check
```bash
curl http://your-domain.com/health
```

### Register User
```bash
curl -X POST http://your-domain.com/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123",
    "name": "Test User",
    "gender": "male",
    "birth_date": "1990-01-01",
    "bio": "Test bio"
  }'
```

### Login
```bash
curl -X POST http://your-domain.com/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Authenticated Request
```bash
curl http://your-domain.com/api/auth/me \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 🔧 Troubleshooting

### Services won't start
```bash
docker compose -f docker-compose.prod.yml logs
docker compose -f docker-compose.prod.yml down
docker compose -f docker-compose.prod.yml up -d
```

### 502 Bad Gateway
```bash
# Check if backend is running
docker ps | grep tinder_backend

# Test backend directly
curl http://localhost:8080/health

# Restart
docker compose -f docker-compose.prod.yml restart backend nginx
```

### SSL Issues
```bash
# Check certificates
docker compose -f docker-compose.prod.yml run --rm certbot certificates

# Renew manually
docker compose -f docker-compose.prod.yml run --rm certbot renew
docker compose -f docker-compose.prod.yml restart nginx
```

### Database Issues
```bash
# Check PostgreSQL
docker compose -f docker-compose.prod.yml logs postgres
docker compose -f docker-compose.prod.yml exec postgres pg_isready -U postgres

# Reset if needed
docker compose -f docker-compose.prod.yml down -v
docker compose -f docker-compose.prod.yml up -d
```

## 📊 Postman Testing

### Environment Variables
```
base_url: https://api.yourdomain.com
token: <from-login-response>
```

### Collections
1. Import `Spark_API.postman_collection.json`
2. Update base URL to your domain
3. After login, save token to environment
4. Use `{{token}}` in Authorization headers

## 🛡 Security Checklist

- ✅ Strong passwords in `.env.production`
- ✅ Firewall configured (UFW or Security Group)
- ✅ SSH key-based authentication only
- ✅ SSL/TLS enabled
- ✅ Regular backups scheduled
- ✅ Monitoring in place
- ✅ `.env.production` NOT in git

## 📈 Rate Limits (Default)

- Auth endpoints: 5 req/min per IP
- API endpoints: 30 req/sec per IP
- General: 10 req/sec per IP
- WebSocket: 10 concurrent connections per IP

Edit `deploy/nginx/nginx.conf` to adjust.

## 💾 File Locations

- Application: `/opt/tinder-app/`
- Backups: `/opt/tinder-app/backups/`
- Docker data: `/var/lib/docker/volumes/`
- Nginx logs: `/var/log/nginx/` (in container)
- SSL certs: `/etc/letsencrypt/` (in volume)

## 🔗 Useful Links

- EC2 Dashboard: https://console.aws.amazon.com/ec2/
- Route53 (DNS): https://console.aws.amazon.com/route53/
- CloudWatch: https://console.aws.amazon.com/cloudwatch/
