# Deployment Checklist

Use this checklist to ensure a smooth deployment process.

## Pre-Deployment

### AWS Setup
- [ ] EC2 instance launched (Ubuntu 22.04, t3.medium+)
- [ ] Security group configured (ports 22, 80, 443)
- [ ] Elastic IP assigned
- [ ] SSH key pair available
- [ ] Can SSH into instance

### DNS Configuration
- [ ] Domain name purchased/available
- [ ] DNS A record created pointing to EC2 IP
- [ ] DNS propagation verified (`dig +short your-domain.com`)

### Code Preparation
- [ ] Backend code tested locally
- [ ] All tests passing
- [ ] Environment variables documented
- [ ] Secrets generated (use `openssl rand -base64 32`)

## Initial Setup

### EC2 Setup
- [ ] Uploaded `setup-ec2.sh` to instance
- [ ] Ran `sudo bash setup-ec2.sh`
- [ ] Logged out and back in (for docker group)
- [ ] Verified Docker installed: `docker --version`
- [ ] Verified Docker Compose installed: `docker compose version`

### Code Upload
- [ ] Application code uploaded to `/opt/tinder-app/`
- [ ] Files extracted correctly
- [ ] Proper permissions set: `chown -R ubuntu:ubuntu /opt/tinder-app`

### Environment Configuration
- [ ] Copied `.env.production.example` to `.env.production`
- [ ] Set `POSTGRES_PASSWORD` (strong, random)
- [ ] Set `REDIS_PASSWORD` (strong, random)
- [ ] Set `JWT_SECRET` (32+ characters, random)
- [ ] Set `DOMAIN_NAME` (your actual domain)
- [ ] Set `EMAIL_FOR_SSL` (your email)
- [ ] Configured OAuth credentials (if using)
- [ ] Configured S3 credentials (if using)

### Scripts Setup
- [ ] Made scripts executable: `chmod +x deploy/scripts/*.sh`

## Deployment

### Initial Deployment
- [ ] Ran `./deploy.sh` successfully
- [ ] All containers started
- [ ] Checked container status: `docker compose ps`
- [ ] Verified no errors in logs
- [ ] Health check passes: `curl http://localhost/health`

### Testing
- [ ] Tested via EC2 public IP: `curl http://EC2_IP/health`
- [ ] Registered test user via Postman
- [ ] Logged in successfully
- [ ] Token authentication works
- [ ] WebSocket connection works (if applicable)

## SSL/TLS Setup

### SSL Configuration
- [ ] DNS resolves correctly to server IP
- [ ] Port 80 accessible from internet
- [ ] Port 443 accessible from internet
- [ ] Ran `./setup-ssl.sh` successfully
- [ ] Certificate obtained from Let's Encrypt
- [ ] HTTPS endpoint works: `curl https://your-domain.com/health`
- [ ] HTTP redirects to HTTPS
- [ ] SSL certificate verified in browser

## Post-Deployment

### Monitoring Setup
- [ ] Ran `./monitor.sh` to verify status
- [ ] All health checks passing
- [ ] No critical errors in logs
- [ ] Resource usage acceptable

### Backup Configuration
- [ ] Ran `./backup.sh` to test backup
- [ ] Backup created successfully
- [ ] Tested restore process (optional but recommended)
- [ ] Set up cron job for automatic backups (see below)

### Systemd Service (Optional)
- [ ] Ran `sudo ./setup-systemd.sh`
- [ ] Service enabled: `sudo systemctl status tinder-app`
- [ ] Tested reboot (optional): service starts automatically

### Security
- [ ] Firewall enabled and configured
- [ ] SSH password authentication disabled (key-only)
- [ ] `.env.production` not in version control
- [ ] Strong passwords used for all services
- [ ] Security group rules restrictive

## Testing with Postman

### Postman Setup
- [ ] Imported `Spark_API.postman_collection.json`
- [ ] Created environment with:
  - `base_url`: `https://your-domain.com`
  - `token`: (will be set after login)
- [ ] Updated collection to use `{{base_url}}`

### API Testing
- [ ] Health check works
- [ ] User registration works
- [ ] User login works
- [ ] Token saved to environment
- [ ] Authenticated endpoints work
- [ ] File upload works
- [ ] WebSocket connection works

## Ongoing Maintenance

### Regular Tasks
- [ ] Set up cron job for daily backups
- [ ] Set up monitoring/alerting (CloudWatch, etc.)
- [ ] Review logs weekly
- [ ] Update system packages monthly
- [ ] Review SSL certificate renewal (automatic)
- [ ] Test backup restore quarterly

### Cron Jobs Setup

```bash
# Edit crontab
crontab -e

# Add these lines:

# Daily backup at 2 AM
0 2 * * * /opt/tinder-app/deploy/scripts/backup.sh >> /opt/tinder-app/logs/backup.log 2>&1

# Health check every 5 minutes
*/5 * * * * /opt/tinder-app/deploy/scripts/health-check.sh || echo "Health check failed at $(date)" >> /opt/tinder-app/logs/health.log

# SSL certificate renewal check (Let's Encrypt auto-renews, this is backup)
0 3 * * 0 docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml run --rm certbot renew >> /opt/tinder-app/logs/ssl.log 2>&1
```

## Documentation
- [ ] Updated internal documentation with server details
- [ ] Documented environment variables
- [ ] Shared access credentials securely with team
- [ ] Created runbook for common operations

## Final Verification

### Smoke Tests
- [ ] Can register new user
- [ ] Can login with credentials
- [ ] Can upload profile photo
- [ ] Can discover other users
- [ ] Can send swipes
- [ ] Can match with users
- [ ] Can send messages
- [ ] Notifications work
- [ ] WebSocket updates work

### Performance
- [ ] Response times acceptable (< 200ms for most endpoints)
- [ ] Rate limiting working correctly
- [ ] No memory leaks
- [ ] Database queries optimized
- [ ] Redis caching working

### Security
- [ ] All secrets are secure
- [ ] HTTPS enforced
- [ ] Security headers present
- [ ] Rate limiting effective
- [ ] No sensitive data in logs

---

## Notes

**Date Deployed**: _________________

**Deployed By**: _________________

**Server Details**:
- EC2 Instance ID: _________________
- Elastic IP: _________________
- Domain: _________________

**Issues Encountered**: 

_________________

**Resolved By**:

_________________

---

## Quick Commands Reference

```bash
# View logs
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs -f

# Restart app
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart

# Create backup
/opt/tinder-app/deploy/scripts/backup.sh

# Monitor status
/opt/tinder-app/deploy/scripts/monitor.sh

# Health check
curl https://your-domain.com/health
```
