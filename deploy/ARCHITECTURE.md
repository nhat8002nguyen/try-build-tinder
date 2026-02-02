# Production Architecture Overview

## System Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                          INTERNET                                │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             │ HTTPS (443)
                             │ HTTP (80) → Redirect to HTTPS
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      AWS EC2 Instance                            │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                   Nginx (Port 80/443)                      │  │
│  │  - Reverse Proxy                                           │  │
│  │  - SSL/TLS Termination (Let's Encrypt)                    │  │
│  │  - Rate Limiting                                           │  │
│  │  - Security Headers                                        │  │
│  │  - Gzip Compression                                        │  │
│  └─────────────────────────┬─────────────────────────────────┘  │
│                            │                                     │
│                            │ HTTP (8080)                        │
│                            ▼                                     │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │              Backend API (Go/Gin)                          │  │
│  │  - REST API Endpoints                                      │  │
│  │  - WebSocket Support                                       │  │
│  │  - JWT Authentication                                      │  │
│  │  - Business Logic                                          │  │
│  │  - File Upload Handling                                    │  │
│  └─────────────┬─────────────────────┬───────────────────────┘  │
│                │                     │                           │
│                │                     │                           │
│     ┌──────────▼──────────┐   ┌──────▼──────────┐             │
│     │   PostgreSQL (5432) │   │  Redis (6379)   │             │
│     │  - User Data        │   │  - Sessions     │             │
│     │  - Matches          │   │  - Cache        │             │
│     │  - Messages         │   │  - Rate Limit   │             │
│     │  - Notifications    │   └─────────────────┘             │
│     └─────────────────────┘                                     │
│                                                                  │
│  ┌───────────────────────────────────────────────────────────┐  │
│  │                    Docker Volumes                          │  │
│  │  - postgres_data: Database persistence                     │  │
│  │  - redis_data: Redis persistence                           │  │
│  │  - uploads_data: User uploaded files                       │  │
│  │  - certbot_conf: SSL certificates                          │  │
│  │  - nginx_logs: Access and error logs                       │  │
│  └───────────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────────┘
```

## Network Flow

### 1. HTTPS Request Flow
```
User → DNS → AWS EC2 IP → Nginx (443) → SSL Termination → 
Backend (8080) → Database/Redis → Response → Nginx → User
```

### 2. WebSocket Flow
```
User → Nginx (upgrade connection) → Backend WebSocket Hub → 
Real-time updates to connected clients
```

### 3. File Upload Flow
```
User → Nginx (max 20MB) → Backend → 
Local Storage (/app/uploads) or S3 → Response URL
```

## Security Layers

### Layer 1: Network (AWS Security Group)
- Port 22: SSH (restricted to admin IPs)
- Port 80: HTTP (public, redirects to HTTPS)
- Port 443: HTTPS (public)
- All other ports: Blocked

### Layer 2: Firewall (UFW on EC2)
- SSH: Allowed
- HTTP/HTTPS: Allowed
- Everything else: Denied by default

### Layer 3: Nginx
- Rate limiting per IP
- SSL/TLS encryption
- Security headers (XSS, CSRF, etc.)
- Request size limits
- Connection limits

### Layer 4: Application
- JWT token authentication
- Password hashing (bcrypt)
- Input validation
- SQL injection prevention (GORM)
- CORS configuration

### Layer 5: Database
- Password protected
- Internal network only (not exposed externally)
- Encrypted connections

## Rate Limiting Configuration

```
┌─────────────────────────┬──────────────────┬─────────────────┐
│ Endpoint Type           │ Rate Limit       │ Burst Limit     │
├─────────────────────────┼──────────────────┼─────────────────┤
│ Auth (login/register)   │ 5 req/min        │ 3 requests      │
│ API (general)           │ 30 req/sec       │ 20 requests     │
│ WebSocket               │ 10 connections   │ 5 connections   │
│ General                 │ 10 req/sec       │ 10 requests     │
└─────────────────────────┴──────────────────┴─────────────────┘
```

## High Availability Features

### Health Checks
- Backend: HTTP GET /health every 30s
- PostgreSQL: pg_isready every 10s
- Redis: PING every 10s
- Nginx: Process monitoring

### Automatic Restart
- Docker restart policy: `unless-stopped`
- Systemd service (optional): Auto-start on boot
- Health-based restarts

### Backup Strategy
- Automated backups via cron (daily at 2 AM)
- Backup includes:
  - PostgreSQL database dump
  - Redis data
  - Uploaded files
  - Configuration files
- Retention: 7 days
- Location: `/opt/tinder-app/backups/`

### SSL/TLS Certificate Renewal
- Let's Encrypt certificates
- Auto-renewal every 12 hours
- 90-day validity
- Nginx reload on renewal

## Monitoring Points

### Application Metrics
- Container status (up/down)
- CPU usage per container
- Memory usage per container
- Network I/O
- Disk usage

### Application Logs
- Nginx access logs
- Nginx error logs
- Backend application logs
- PostgreSQL logs
- Redis logs

### Health Endpoints
- `/health` - Overall system health
- Database connection test
- Redis connection test

## Scaling Considerations

### Vertical Scaling (Current Setup)
- Upgrade EC2 instance type
- Increase volume size
- Add more RAM/CPU

### Horizontal Scaling (Future)
```
┌────────────────────────────────────────────────┐
│          Load Balancer (ALB/ELB)               │
└──────┬──────────┬──────────┬─────────┬─────────┘
       │          │          │         │
   ┌───▼───┐  ┌──▼───┐  ┌───▼───┐ ┌──▼───┐
   │ EC2-1 │  │ EC2-2│  │ EC2-3 │ │EC2-N │
   └───┬───┘  └──┬───┘  └───┬───┘ └──┬───┘
       │         │          │         │
       └─────────┴──────────┴─────────┘
                    │
        ┌───────────┴───────────┐
        │                       │
   ┌────▼────┐            ┌─────▼─────┐
   │   RDS   │            │ ElastiCache│
   │PostgreSQL            │   Redis    │
   └─────────┘            └────────────┘
```

### Database Scaling
- RDS PostgreSQL (managed)
- Read replicas for queries
- ElastiCache for Redis
- S3 for file storage

## Resource Requirements

### Minimum (Development/Testing)
- EC2: t3.small (2 vCPU, 2GB RAM)
- Storage: 20GB
- Bandwidth: Low

### Recommended (Production - Small)
- EC2: t3.medium (2 vCPU, 4GB RAM)
- Storage: 30GB
- Bandwidth: Moderate

### Recommended (Production - Medium)
- EC2: t3.large (2 vCPU, 8GB RAM)
- Storage: 50GB
- Bandwidth: High

### Enterprise (High Traffic)
- EC2: Multiple t3.xlarge instances
- RDS: db.t3.large
- ElastiCache: cache.t3.medium
- ALB for load balancing
- S3 for file storage
- CloudWatch for monitoring

## Deployment Pipeline

```
┌──────────────┐
│  Developer   │
│  Git Push    │
└──────┬───────┘
       │
       ▼
┌──────────────────┐
│  GitHub Actions  │
│  - Run Tests     │
│  - Build Image   │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│   Upload to EC2  │
│   via SSH        │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│  Create Backup   │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│  Stop Containers │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│  Load New Image  │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│ Start Containers │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│  Health Check    │
└──────┬───────────┘
       │
       ▼
┌──────────────────┐
│  Notify Status   │
└──────────────────┘
```

## Technology Stack

### Infrastructure
- **Cloud Provider**: AWS EC2
- **OS**: Ubuntu 22.04 LTS
- **Containerization**: Docker + Docker Compose
- **Reverse Proxy**: Nginx 1.25
- **SSL/TLS**: Let's Encrypt (Certbot)

### Backend
- **Language**: Go 1.21
- **Framework**: Gin
- **Database**: PostgreSQL 15
- **Cache**: Redis 7
- **Authentication**: JWT (golang-jwt/jwt)
- **ORM**: GORM
- **WebSocket**: gorilla/websocket

### DevOps
- **CI/CD**: GitHub Actions
- **Monitoring**: Docker stats, health checks
- **Logging**: JSON logs, log rotation
- **Backup**: Automated cron jobs

## Performance Characteristics

### Response Times (Expected)
- Health check: < 10ms
- Authentication: < 100ms
- API queries: < 200ms
- File upload (1MB): < 500ms
- WebSocket latency: < 50ms

### Throughput (Single t3.medium)
- Concurrent connections: ~1000
- Requests per second: ~500
- WebSocket connections: ~100

### Storage
- Database: Grows with users (~100KB per user)
- Redis: ~10-50MB typical
- Uploads: User-dependent (limit 20MB per file)

## Cost Estimation (AWS)

### Basic Setup (Monthly)
```
EC2 t3.medium (on-demand):     ~$30
EBS 30GB storage:              ~$3
Data transfer (50GB):          ~$5
Elastic IP:                    Free (if attached)
────────────────────────────────────
Total:                         ~$38/month
```

### With Reserved Instance (1 year)
```
EC2 t3.medium (reserved):      ~$20
Other costs:                   ~$8
────────────────────────────────────
Total:                         ~$28/month
```

### Cost Optimization Tips
1. Use reserved instances (save 30-40%)
2. Use spot instances for dev/staging (save 70%)
3. Enable AWS Savings Plans
4. Use S3 for file storage instead of EBS
5. Set up auto-scaling to scale down during low traffic

## Disaster Recovery

### Backup Strategy
- **Frequency**: Daily (2 AM)
- **Retention**: 7 days
- **Location**: Same EC2 instance + copy to S3 (recommended)
- **Components**: Database, Redis, uploads, configs

### Recovery Time Objective (RTO)
- Manual restore: ~15 minutes
- Full redeployment: ~30 minutes

### Recovery Point Objective (RPO)
- Database: Last backup (max 24 hours)
- Files: Last backup (max 24 hours)

### Recovery Steps
1. Stop application
2. Restore volumes from backup
3. Start application
4. Verify functionality
5. Monitor for issues

---

**Document Version**: 1.0  
**Last Updated**: January 2026  
**Maintained By**: DevOps Team
