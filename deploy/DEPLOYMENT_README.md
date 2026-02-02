# 🚀 AWS EC2 Production Deployment - Complete Package

## Overview

This is a **production-grade deployment package** for deploying the Tinder Clone backend to AWS EC2 with Docker, nginx, SSL/TLS, rate limiting, monitoring, and automated backups.

## ✨ What You Get

🎯 **Complete Production Infrastructure**
- Nginx reverse proxy with rate limiting
- SSL/TLS with Let's Encrypt (automatic renewal)
- Docker + Docker Compose setup
- PostgreSQL 15 database
- Redis 7 caching
- Automated backups
- Health monitoring
- Security hardening

📚 **Comprehensive Documentation** (3,500+ lines)
- 7 detailed guides
- Step-by-step checklists
- Architecture diagrams
- Troubleshooting guides
- Quick reference cards

🔧 **8 Automated Scripts**
- EC2 setup automation
- One-command deployment
- SSL certificate management
- Backup/restore utilities
- Monitoring dashboard
- Health checks

🔒 **Enterprise Security**
- SSL/TLS encryption (HTTPS only)
- Rate limiting per IP
- JWT authentication
- Password hashing (bcrypt)
- Firewall configuration
- Security headers (XSS, CSRF, etc.)

## 📦 What's Included

```
deploy/
├── 📖 Documentation (7 guides)
│   ├── INDEX.md                   ← Documentation index
│   ├── GETTING_STARTED.md         ← Start here!
│   ├── README.md                  ← Complete guide (500+ lines)
│   ├── DEPLOYMENT_CHECKLIST.md    ← Step-by-step checklist
│   ├── QUICK_REFERENCE.md         ← Common commands
│   ├── ARCHITECTURE.md            ← System architecture
│   ├── CI_CD_SETUP.md            ← GitHub Actions guide
│   └── SUMMARY.md                 ← Package overview
│
├── ⚙️ Configuration Files
│   ├── docker-compose.prod.yml    ← Production Docker Compose
│   ├── .env.production.example    ← Environment template
│   ├── Makefile                   ← Command shortcuts
│   └── .gitignore                 ← Security (excludes secrets)
│
├── 🌐 Nginx Configuration
│   ├── nginx.conf                 ← Main nginx config
│   └── conf.d/
│       ├── default.conf          ← HTTPS with rate limiting
│       └── http-only.conf        ← HTTP-only (for testing)
│
├── 🔧 Deployment Scripts (8 scripts)
│   ├── setup-ec2.sh              ← Initial EC2 setup
│   ├── deploy.sh                 ← Main deployment
│   ├── setup-ssl.sh              ← SSL/TLS setup
│   ├── backup.sh                 ← Create backups
│   ├── restore.sh                ← Restore from backup
│   ├── monitor.sh                ← Monitoring dashboard
│   ├── health-check.sh           ← Health checks
│   └── setup-systemd.sh          ← Auto-start config
│
└── 🔄 Auto-start Configuration
    └── systemd/
        └── tinder-app.service    ← Systemd service

.github/workflows/
└── deploy.yml                     ← GitHub Actions CI/CD

backend/
└── .dockerignore                  ← Docker build optimization
```

## 🚀 Quick Start (5 Steps)

### 1. Launch EC2 Instance
```
Instance: t3.medium (2 vCPU, 4GB RAM)
OS: Ubuntu 22.04 LTS
Storage: 30GB
Security Group: Ports 22, 80, 443 open
```

### 2. Run EC2 Setup
```bash
scp -i key.pem deploy/scripts/setup-ec2.sh ubuntu@EC2_IP:~/
ssh -i key.pem ubuntu@EC2_IP
sudo bash setup-ec2.sh
exit && ssh -i key.pem ubuntu@EC2_IP
```

### 3. Upload Code
```bash
tar --exclude='node_modules' --exclude='frontend' --exclude='.git' \
    -czf app.tar.gz backend/ deploy/
scp -i key.pem app.tar.gz ubuntu@EC2_IP:/opt/tinder-app/
ssh -i key.pem ubuntu@EC2_IP
cd /opt/tinder-app && tar -xzf app.tar.gz
```

### 4. Configure Environment
```bash
cd /opt/tinder-app/deploy
cp .env.production.example .env.production
nano .env.production  # Set passwords and secrets
```

### 5. Deploy
```bash
cd scripts && chmod +x *.sh
./deploy.sh              # Deploy application
./setup-ssl.sh           # Setup SSL (after DNS configured)
```

**That's it!** Your production backend is now running at `https://your-domain.com`

## 📚 Documentation Guide

### 🎯 **Start Here**
👉 **[deploy/GETTING_STARTED.md](deploy/GETTING_STARTED.md)** - Overview and quick start (10 min)  
👉 **[deploy/DEPLOYMENT_CHECKLIST.md](deploy/DEPLOYMENT_CHECKLIST.md)** - Step-by-step guide (30 min)

### 📖 **Complete Reference**
- **[deploy/README.md](deploy/README.md)** - Complete deployment guide (30 min read)
- **[deploy/ARCHITECTURE.md](deploy/ARCHITECTURE.md)** - System architecture and design (15 min)

### ⚡ **Quick Reference**
- **[deploy/QUICK_REFERENCE.md](deploy/QUICK_REFERENCE.md)** - Common commands (2 min)
- **[deploy/INDEX.md](deploy/INDEX.md)** - Documentation index

### 🤖 **Automation**
- **[deploy/CI_CD_SETUP.md](deploy/CI_CD_SETUP.md)** - GitHub Actions setup (20 min)

## 🎯 Key Features

### Production-Grade Configuration
✅ **SSL/TLS** - Free Let's Encrypt certificates with auto-renewal  
✅ **Rate Limiting** - Protect against abuse (5 req/min auth, 30 req/sec API)  
✅ **Reverse Proxy** - Nginx with security headers and compression  
✅ **Containerization** - Docker for easy deployment and scaling  
✅ **Monitoring** - Health checks and resource monitoring  
✅ **Backups** - Automated daily backups with 7-day retention  
✅ **Logging** - Structured logs with automatic rotation  
✅ **Security** - Firewall, JWT auth, password hashing  

### Rate Limiting Configuration
| Endpoint Type | Rate Limit | Burst Limit |
|--------------|------------|-------------|
| Auth (login/register) | 5 req/min | 3 requests |
| API endpoints | 30 req/sec | 20 requests |
| WebSocket | 10 connections | 5 connections |
| General | 10 req/sec | 10 requests |

## 🛠 Common Operations

### Using Makefile (Recommended)
```bash
cd /opt/tinder-app/deploy

make help           # Show all commands
make deploy         # Deploy/update application
make logs           # View logs (follow mode)
make monitor        # Monitoring dashboard
make backup         # Create backup
make restore        # Restore from backup
make health         # Health check
make ssl            # Setup SSL/TLS
make restart        # Restart all services
```

### Manual Commands
```bash
# View logs
docker compose -f docker-compose.prod.yml logs -f

# Restart services
docker compose -f docker-compose.prod.yml restart

# Check status
docker compose -f docker-compose.prod.yml ps

# Health check
curl https://your-domain.com/health
```

## 🧪 Testing with Postman

### 1. Import Collection
Use `Spark_API.postman_collection.json` at project root

### 2. Create Environment
```
base_url: https://your-domain.com
token: (will be set after login)
```

### 3. Test Endpoints

**Health Check:**
```
GET {{base_url}}/health
```

**Register:**
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

**Login:**
```
POST {{base_url}}/api/auth/login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "password123"
}
```

**Get Profile (Authenticated):**
```
GET {{base_url}}/api/auth/me
Authorization: Bearer {{token}}
```

## 🔒 Security Best Practices

### Before Deployment
✅ Generate strong passwords: `openssl rand -base64 32`  
✅ Update all values in `.env.production`  
✅ **Never commit** `.env.production` to git  
✅ Configure AWS Security Group properly  
✅ Use Elastic IP for stable DNS  

### After Deployment
✅ Enable SSH key-only authentication  
✅ Disable password SSH login  
✅ Set up CloudWatch monitoring  
✅ Regular security updates: `sudo apt update && sudo apt upgrade`  
✅ Review logs weekly  
✅ Test backup/restore monthly  

## 💰 Cost Estimation

### Monthly AWS Costs
```
EC2 t3.medium (on-demand):    $30/month
EBS 30GB storage:             $3/month
Data transfer (50GB):         $5/month
Elastic IP (attached):        Free
─────────────────────────────────────
Total:                        ~$38/month
```

### Optimizations
- Reserved Instance (1 year): Save 30% → **$26/month**
- Spot Instance (dev/staging): Save 70% → **$9/month**

## 📊 Performance Expectations

### Single t3.medium Instance
- **Concurrent Users:** ~1,000
- **Requests/Second:** ~500
- **Response Time:** <200ms (typical)
- **Database Capacity:** ~10,000 users
- **WebSocket Connections:** ~100

## 🆘 Quick Troubleshooting

### Services won't start
```bash
make logs                    # Check logs
make restart                 # Restart services
```

### 502 Bad Gateway
```bash
make logs-backend           # Check backend logs
make restart-backend        # Restart backend
```

### SSL Issues
```bash
make ssl-check              # Check certificate
make ssl-renew              # Renew certificate
```

### Out of disk space
```bash
make clean                  # Clean Docker resources
docker system prune -a      # Deep clean
```

## 📞 Getting Help

### Check Documentation
1. **[deploy/QUICK_REFERENCE.md](deploy/QUICK_REFERENCE.md)** - Common commands
2. **[deploy/README.md](deploy/README.md)** - Troubleshooting section
3. **[deploy/ARCHITECTURE.md](deploy/ARCHITECTURE.md)** - System design

### Debug Steps
```bash
make logs                   # View all logs
make health                 # Check health
make monitor                # View dashboard
make ps                     # Container status
```

## ✅ Pre-Launch Checklist

- [ ] EC2 instance launched and accessible
- [ ] DNS A record pointing to EC2 IP
- [ ] `.env.production` configured with strong secrets
- [ ] Application deployed successfully
- [ ] Health check passing
- [ ] SSL certificate valid (HTTPS working)
- [ ] Rate limiting tested
- [ ] Can register/login via Postman
- [ ] Backups configured
- [ ] Monitoring in place

## 🎓 Next Steps After Deployment

### Immediate (Today)
1. Test all API endpoints with Postman
2. Verify SSL certificate in browser
3. Check logs for any errors
4. Test rate limiting
5. Create first manual backup

### Short-term (This Week)
1. Set up automated backups (cron)
2. Configure CloudWatch monitoring
3. Set up GitHub Actions CI/CD
4. Load test with 100 concurrent users
5. Document credentials securely

### Long-term (This Month)
1. Set up staging environment
2. Configure CloudWatch alarms
3. Implement centralized logging
4. Performance optimization
5. Security audit
6. Disaster recovery testing

## 🎉 What You've Achieved

You now have a **production-grade backend deployment** with:

✅ **3,500+ lines** of documentation  
✅ **8 automated scripts** for all operations  
✅ **Enterprise security** (SSL, rate limiting, firewall)  
✅ **Automated backups** with restore capability  
✅ **Monitoring and health checks**  
✅ **CI/CD ready** (GitHub Actions)  
✅ **Cost-optimized** (~$38/month)  
✅ **Scalable architecture**  

**Deployment time:** 30-45 minutes  
**Update time:** 5 minutes  

## 📝 Important Notes

- **Security:** Never commit `.env.production` to version control
- **Backups:** Test restore process regularly
- **Updates:** Keep system and packages updated
- **Monitoring:** Check logs weekly for issues
- **Secrets:** Rotate passwords and tokens quarterly

## 🚀 Ready to Deploy?

### First-Time Deployment
👉 Start with **[deploy/GETTING_STARTED.md](deploy/GETTING_STARTED.md)**

### Need Quick Commands
👉 Go to **[deploy/QUICK_REFERENCE.md](deploy/QUICK_REFERENCE.md)**

### Want Complete Guide
👉 Read **[deploy/README.md](deploy/README.md)**

---

## 📞 Support

For questions or issues:
1. Check documentation in `deploy/` directory
2. Review [deploy/QUICK_REFERENCE.md](deploy/QUICK_REFERENCE.md)
3. Check logs: `make logs`
4. Review troubleshooting in [deploy/README.md](deploy/README.md)

---

## 📄 License

This deployment package is part of the Tinder Clone project.

---

**Version:** 1.0  
**Created:** January 2026  
**Platform:** AWS EC2 + Docker + Nginx  
**Application:** Tinder Clone Backend (Go)  

**Good luck with your deployment! 🎉**

Need help? Start with [deploy/GETTING_STARTED.md](deploy/GETTING_STARTED.md)!
