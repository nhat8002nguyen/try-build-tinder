# Makefile for Tinder Clone Production Deployment

.PHONY: help deploy start stop restart logs monitor backup restore health clean ssl

# Default target
help:
	@echo "Tinder Clone - Production Deployment"
	@echo ""
	@echo "Available commands:"
	@echo "  make deploy          - Deploy application"
	@echo "  make start           - Start all services"
	@echo "  make stop            - Stop all services"
	@echo "  make restart         - Restart all services"
	@echo "  make logs            - View all logs (follow mode)"
	@echo "  make logs-backend    - View backend logs"
	@echo "  make logs-nginx      - View nginx logs"
	@echo "  make monitor         - Show monitoring dashboard"
	@echo "  make backup          - Create backup"
	@echo "  make restore         - Restore from backup (usage: make restore TIMESTAMP=20260125_120000)"
	@echo "  make health          - Check service health"
	@echo "  make ssl             - Setup SSL/TLS"
	@echo "  make clean           - Clean old Docker resources"
	@echo "  make ps              - Show running containers"
	@echo "  make shell           - Access backend shell"
	@echo "  make db              - Access PostgreSQL shell"
	@echo ""

# Deployment
deploy:
	@echo "Deploying application..."
	cd /opt/tinder-app/deploy/scripts && ./deploy.sh

# Service management
start:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml up -d

stop:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml down

restart:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart

restart-backend:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart backend

restart-nginx:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart nginx

# Logs
logs:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs -f

logs-backend:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs -f backend

logs-nginx:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs -f nginx

logs-postgres:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs -f postgres

logs-redis:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs -f redis

# Monitoring
monitor:
	cd /opt/tinder-app/deploy/scripts && ./monitor.sh

health:
	cd /opt/tinder-app/deploy/scripts && ./health-check.sh

ps:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml ps

stats:
	docker stats

# Backup and restore
backup:
	cd /opt/tinder-app/deploy/scripts && ./backup.sh

restore:
	@if [ -z "$(TIMESTAMP)" ]; then \
		echo "Error: TIMESTAMP required. Usage: make restore TIMESTAMP=20260125_120000"; \
		echo "Available backups:"; \
		ls -1 /opt/tinder-app/backups/ 2>/dev/null || echo "  No backups found"; \
		exit 1; \
	fi
	cd /opt/tinder-app/deploy/scripts && ./restore.sh $(TIMESTAMP)

# SSL
ssl:
	cd /opt/tinder-app/deploy/scripts && ./setup-ssl.sh

ssl-renew:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml run --rm certbot renew
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml restart nginx

ssl-check:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml run --rm certbot certificates

# Maintenance
clean:
	docker system prune -f
	@echo "Cleaned up unused Docker resources"

clean-all:
	docker system prune -a -f --volumes
	@echo "WARNING: Cleaned up all Docker resources including volumes!"

pull:
	cd /opt/tinder-app && git pull

update: pull deploy
	@echo "Application updated successfully"

# Database
db:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml exec postgres psql -U postgres -d tinder_clone

db-backup:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml exec postgres pg_dump -U postgres tinder_clone | gzip > /opt/tinder-app/backups/manual_backup_$(shell date +%Y%m%d_%H%M%S).sql.gz

# Shell access
shell:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml exec backend sh

redis:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml exec redis redis-cli

# Testing
test-health:
	@curl -f http://localhost/health && echo "✓ Health check passed" || echo "✗ Health check failed"

test-api:
	@echo "Testing API endpoints..."
	@curl -f http://localhost/health && echo "✓ Health endpoint OK" || echo "✗ Health endpoint FAILED"
	@echo ""

# Build
build:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml build --no-cache

rebuild: stop build start
	@echo "Application rebuilt and restarted"

# Development helpers (use with caution in production)
dev-logs:
	docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs --tail=100

dev-reset: stop clean start
	@echo "Development reset complete"
