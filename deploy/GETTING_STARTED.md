# AWS EC2 Production Deployment - Complete Package

## 📦 What's Included

This deployment package provides everything you need to deploy your Tinder Clone backend to AWS EC2 with production-grade configuration.

### Directory Structure

```
deploy/
├── README.md                    # Main deployment guide
├── QUICK_REFERENCE.md          # Quick command reference
├── DEPLOYMENT_CHECKLIST.md     # Step-by-step checklist
├── ARCHITECTURE.md             # System architecture details
├── CI_CD_SETUP.md              # GitHub Actions setup
├── .env.production.example     # Environment template
│
├── nginx/                      # Nginx configuration
│   ├── nginx.conf             # Main nginx config
│   └── conf.d/
│       ├── default.conf       # HTTPS configuration
│       └── http-only.conf     # HTTP-only (for testing)
│
├── scripts/                    # Deployment scripts
│   ├── setup-ec2.sh           # Initial EC2 setup
│   ├── deploy.sh              # Main deployment script
│   ├── setup-ssl.sh           # SSL/TLS configuration
│   ├── backup.sh              # Backup utility
│   ├── restore.sh             # Restore utility
│   ├── monitor.sh             # Monitoring dashboard
│   ├── setup-systemd.sh       # Auto-start configuration
│   └── health-check.sh        # Health check utility
│
├── systemd/                    # Systemd service files
│   └── tinder-app.service     # Auto-start service
│
└── docker-compose.prod.yml     # Production Docker Compose

.github/
└── workflows/
    └── deploy.yml              # GitHub Actions CI/CD

backend/
└── .dockerignore               # Docker build optimization
```

## 🚀 Quick Start

### For First-Time Deployment

1. **Launch EC2 Instance** (Ubuntu 22.04, t3.medium+)
2. **Run Setup**:
   ```bash
   scp -i key.pem deploy/scripts/setup-ec2.sh ubuntu@EC2_IP:~/
   ssh -i key.pem ubuntu@EC2_IP
   sudo bash setup-ec2.sh
   ```

3. **Upload Code**:
   ```bash
   cd /Users/nhatnguyen/Workspaces/try-build-tinder
   tar --exclude='node_modules' --exclude='frontend' -czf app.tar.gz backend/ deploy/
   scp -i key.pem app.tar.gz ubuntu@EC2_IP:/opt/tinder-app/
   ```

4. **Configure Environment**:
   ```bash
   ssh -i key.pem ubuntu@EC2_IP
   cd /opt/tinder-app/deploy
   cp .env.production.example .env.production
   nano .env.production  # Edit with your values
   ```

5. **Deploy**:
   ```bash
   cd scripts
   chmod +x *.sh
   ./deploy.sh
   ```

6. **Setup SSL** (after DNS configured):
   ```bash
   ./setup-ssl.sh
   ```

### For Updates

```bash
cd /opt/tinder-app
git pull
cd deploy/scripts
./deploy.sh
```

## 📚 Documentation Guide

### Start Here
1. **README.md** - Complete deployment guide with all steps
2. **DEPLOYMENT_CHECKLIST.md** - Follow this step-by-step

### For Quick Reference
- **QUICK_REFERENCE.md** - Common commands and quick answers
- **ARCHITECTURE.md** - System design and architecture

### For Advanced Setup
- **CI_CD_SETUP.md** - Automated deployments with GitHub Actions

## 🎯 Key Features

### Production-Grade Configuration
✅ SSL/TLS with Let's Encrypt  
✅ Nginx reverse proxy with rate limiting  
✅ Docker containerization  
✅ Automated backups  
✅ Health monitoring  
✅ Log management  
✅ Auto-restart on failure  
✅ Security hardening  

### Rate Limiting
- Authentication: 5 requests/minute per IP
- API calls: 30 requests/second per IP
- WebSocket: 10 concurrent connections per IP
- General: 10 requests/second per IP

### Security Features
- SSL/TLS encryption (HTTPS only)
- Firewall configuration (UFW + Security Groups)
- JWT authentication
- Password hashing (bcrypt)
- Security headers (XSS, CSRF protection)
- Input validation
- SQL injection prevention

### Monitoring & Maintenance
- Automated health checks
- Container monitoring
- Log rotation (keep last 3 files, 10MB each)
- Daily backups (7-day retention)
- SSL auto-renewal
- System metrics dashboard

## 🛠 Scripts Reference

### setup-ec2.sh
**Purpose**: Initial EC2 instance setup  
**Run Once**: Yes  
**Requires**: Root/sudo  
**Does**:
- Installs Docker & Docker Compose
- Configures firewall (UFW)
- System optimizations
- Creates app directory
- Security hardening

### deploy.sh
**Purpose**: Deploy/update application  
**Run**: Every deployment  
**Requires**: User permissions  
**Does**:
- Creates backup
- Stops old containers
- Builds new images
- Starts services
- Verifies health

### setup-ssl.sh
**Purpose**: Configure SSL/TLS  
**Run Once**: After DNS setup  
**Requires**: Domain configured  
**Does**:
- Obtains Let's Encrypt certificate
- Configures nginx for HTTPS
- Sets up auto-renewal

### backup.sh
**Purpose**: Create backup  
**Run**: Daily via cron (or manual)  
**Does**:
- Backs up PostgreSQL
- Backs up Redis
- Backs up uploads
- Backs up configs
- Cleans old backups (>7 days)

### restore.sh
**Purpose**: Restore from backup  
**Run**: When needed  
**Usage**: `./restore.sh TIMESTAMP`

### monitor.sh
**Purpose**: View system status  
**Run**: Anytime  
**Shows**:
- Container status
- Resource usage
- Health checks
- Recent logs

### setup-systemd.sh
**Purpose**: Enable auto-start on boot  
**Run Once**: Optional  
**Requires**: Root/sudo

### health-check.sh
**Purpose**: Check if services are healthy  
**Run**: Via cron or monitoring tools  
**Returns**: Exit 0 if healthy, 1 if not

## 🧪 Testing Guide

### Test with Postman

1. **Import Collection**: Use `Spark_API.postman_collection.json`

2. **Create Environment**:
   ```
   base_url: https://your-domain.com
   token: (empty, will be set after login)
   ```

3. **Test Endpoints**:

**Health Check**:
```
GET {{base_url}}/health
```

**Register**:
```
POST {{base_url}}/api/auth/register
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123",
  "name": "Test User",
  "gender": "male",
  "birth_date": "1990-01-01"
}
```

**Login**:
```
POST {{base_url}}/api/auth/login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

**Save token from response**, then test authenticated endpoints:

**Get Profile**:
```
GET {{base_url}}/api/auth/me
Authorization: Bearer {{token}}
```

### Test Rate Limiting

Run this script to test rate limits:
```bash
# Test auth rate limit (should block after 5 requests)
for i in {1..10}; do
  curl -X POST https://your-domain.com/api/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"wrong"}' \
    -w "\n%{http_code}\n"
  sleep 1
done
```

Expected: First 5 requests succeed, then 429 (Too Many Requests)

## 🔒 Security Best Practices

### Before Deployment
✅ Generate strong passwords (use: `openssl rand -base64 32`)  
✅ Update `.env.production` with secure values  
✅ Never commit `.env.production` to git  
✅ Configure AWS Security Group properly  

### After Deployment
✅ Enable SSH key-only authentication  
✅ Disable password SSH login  
✅ Set up CloudWatch or monitoring  
✅ Regular security updates: `sudo apt update && sudo apt upgrade`  
✅ Review logs weekly  
✅ Test backup/restore process  

### Recommended
✅ Use AWS Secrets Manager for sensitive data  
✅ Enable AWS CloudTrail  
✅ Set up AWS GuardDuty  
✅ Configure CloudWatch alarms  
✅ Use IAM roles instead of access keys  

## 💰 Cost Estimation

### Basic Production Setup
```
EC2 t3.medium (on-demand):  $30/month
EBS 30GB:                   $3/month
Data transfer (50GB):       $5/month
Domain + SSL:               $12/year (domain only, SSL free)
──────────────────────────────────────
Total:                      ~$38-40/month
```

### Cost Optimization
- Use Reserved Instances: Save 30-40%
- Use Spot Instances (dev/staging): Save 70%
- Schedule auto-stop during off-hours
- Use S3 for file storage (cheaper than EBS)

## 📊 Performance Expectations

### Single t3.medium Instance
- **Concurrent Users**: ~1000
- **Requests/Second**: ~500
- **WebSocket Connections**: ~100
- **Database**: ~10K users comfortably
- **Response Time**: <200ms for most endpoints

### When to Scale Up
- CPU consistently >70%
- Memory consistently >80%
- Response times >500ms
- Database connections maxed out
- Frequent 429 errors (rate limiting)

## 🆘 Quick Troubleshooting

### Services Won't Start
```bash
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml down
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml up -d
```

### 502 Bad Gateway
```bash
docker ps | grep tinder_backend
curl http://localhost:8080/health
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart backend
```

### SSL Issues
```bash
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml run --rm certbot certificates
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml run --rm certbot renew
```

### Out of Disk Space
```bash
df -h
docker system prune -a
docker volume prune
```

### Check Logs
```bash
# All logs
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs -f

# Specific service
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs -f backend

# Nginx errors
docker exec tinder_nginx cat /var/log/nginx/error.log
```

## 📞 Support & Resources

### Included Documentation
- Complete deployment guide (README.md)
- Quick reference (QUICK_REFERENCE.md)
- Architecture details (ARCHITECTURE.md)
- Deployment checklist (DEPLOYMENT_CHECKLIST.md)
- CI/CD setup (CI_CD_SETUP.md)

### Useful Commands
```bash
# View all containers
docker ps -a

# Check resource usage
docker stats

# View logs
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs -f

# Restart everything
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart

# Access backend shell
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml exec backend sh

# Access database
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml exec postgres psql -U postgres -d tinder_clone
```

## 🎓 Learning Resources

### AWS Documentation
- [EC2 User Guide](https://docs.aws.amazon.com/ec2/)
- [Security Groups](https://docs.aws.amazon.com/vpc/latest/userguide/VPC_SecurityGroups.html)
- [Route 53 DNS](https://docs.aws.amazon.com/route53/)

### Docker Documentation
- [Docker Compose](https://docs.docker.com/compose/)
- [Docker Best Practices](https://docs.docker.com/develop/dev-best-practices/)

### Nginx Documentation
- [Nginx Admin Guide](https://docs.nginx.com/nginx/admin-guide/)
- [SSL Configuration](https://ssl-config.mozilla.org/)

## ✅ Success Checklist

After deployment, verify:

- [ ] All containers running (`docker ps`)
- [ ] Health endpoint works (`curl https://your-domain.com/health`)
- [ ] Can register user via Postman
- [ ] Can login and get token
- [ ] Authenticated endpoints work
- [ ] HTTPS redirects working
- [ ] SSL certificate valid (green lock in browser)
- [ ] Rate limiting working (test with multiple requests)
- [ ] Logs are being generated
- [ ] Backups configured (cron job)
- [ ] Monitoring dashboard works (`./monitor.sh`)

## 🎉 Next Steps

After successful deployment:

1. **Set up monitoring**: CloudWatch, Datadog, or similar
2. **Configure alerts**: Email/Slack for errors
3. **Set up CI/CD**: Use GitHub Actions (see CI_CD_SETUP.md)
4. **Schedule backups**: Set up cron job for daily backups
5. **Document credentials**: Store securely (1Password, Vault, etc.)
6. **Load testing**: Use k6, JMeter, or similar
7. **Set up staging**: Clone setup for staging environment
8. **Configure logging**: Centralized logging (CloudWatch, ELK)

---

## 📝 Important Notes

- **Don't commit** `.env.production` to version control
- **Backup regularly** - test restore process
- **Update regularly** - system and application
- **Monitor logs** - check for security issues
- **Rotate secrets** - passwords, keys, tokens
- **Test before deploy** - use staging environment

---

**Version**: 1.0  
**Created**: January 2026  
**For**: Tinder Clone Backend Production Deployment  

**Questions or issues?** Review the documentation or check logs first!

**Good luck with your deployment! 🚀**
