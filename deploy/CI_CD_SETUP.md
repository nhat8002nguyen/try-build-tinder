# GitHub Actions CI/CD Setup Guide

This guide explains how to set up automated deployments using GitHub Actions.

## Overview

The CI/CD pipeline automatically:
1. Runs tests on every push
2. Builds Docker image
3. Deploys to EC2 when pushing to `main` branch
4. Verifies deployment health

## Prerequisites

- GitHub repository with your code
- EC2 instance already set up and configured
- Application already deployed manually at least once

## Setup Instructions

### 1. Generate SSH Key for GitHub Actions

On your local machine:

```bash
# Generate a new SSH key pair
ssh-keygen -t rsa -b 4096 -f github-actions-key -N ""

# This creates:
# - github-actions-key (private key)
# - github-actions-key.pub (public key)
```

### 2. Add Public Key to EC2

```bash
# Copy the public key
cat github-actions-key.pub

# SSH into your EC2 instance
ssh -i your-ec2-key.pem ubuntu@your-ec2-ip

# Add the public key to authorized_keys
echo "YOUR_PUBLIC_KEY_CONTENT" >> ~/.ssh/authorized_keys

# Set proper permissions
chmod 600 ~/.ssh/authorized_keys
```

### 3. Add Secrets to GitHub

Go to your GitHub repository:
1. Click **Settings** → **Secrets and variables** → **Actions**
2. Click **New repository secret**

Add these secrets:

#### EC2_SSH_KEY
```
The PRIVATE key content from github-actions-key file
(entire file content including BEGIN and END lines)
```

#### EC2_HOST
```
Your EC2 public IP or domain
Example: 54.123.45.67 or ec2-54-123-45-67.compute-1.amazonaws.com
```

#### EC2_USER
```
ubuntu
(or whatever user you use to SSH into EC2)
```

#### DOMAIN_NAME (Optional)
```
your-domain.com
(only if you have SSL configured)
```

### 4. Verify Secrets

Your GitHub repository should have these secrets:
- ✅ EC2_SSH_KEY
- ✅ EC2_HOST
- ✅ EC2_USER
- ✅ DOMAIN_NAME (optional)

## Workflow Triggers

The deployment workflow runs when:

### Automatic Triggers
- Push to `main` branch
- Changes in `backend/`, `deploy/`, or workflow file

### Manual Trigger
You can also trigger manually:
1. Go to **Actions** tab in GitHub
2. Select **Deploy to Production**
3. Click **Run workflow**

## Workflow Stages

### 1. Test Stage
```yaml
- Checks out code
- Sets up Go 1.21
- Runs all tests
- Fails if any test fails
```

### 2. Build Stage
```yaml
- Builds Docker image
- Tags with commit SHA and 'latest'
- Uploads image as artifact
```

### 3. Deploy Stage
```yaml
- Downloads Docker image
- Uploads to EC2
- Creates backup of current deployment
- Deploys new version
- Runs health checks
- Cleans up old resources
```

## Monitoring Deployments

### View Deployment Status

1. Go to **Actions** tab in your GitHub repository
2. Click on the latest workflow run
3. View logs for each stage

### Successful Deployment

You should see:
```
✓ Deployment successful!
✓ Deployment verification passed!
🚀 Deployment to production completed successfully!
```

### Failed Deployment

If deployment fails:
1. Check the workflow logs
2. SSH into EC2 to investigate
3. Restore from backup if needed

## Rollback Process

If you need to rollback:

### Method 1: Revert Git Commit
```bash
git revert HEAD
git push origin main
# This triggers a new deployment with previous code
```

### Method 2: Manual Restore on EC2
```bash
ssh -i your-key.pem ubuntu@your-ec2-ip
cd /opt/tinder-app/deploy/scripts
./restore.sh <backup-timestamp>
```

### Method 3: Redeploy Previous Commit
```bash
# Find the commit hash you want to deploy
git log

# Push that commit to main
git checkout <commit-hash>
git push origin HEAD:main --force
# ⚠️ Use force push carefully!
```

## Testing the Pipeline

### Test Without Deploying

To test changes without auto-deploying:

1. Create a branch:
```bash
git checkout -b test-ci
```

2. Push changes:
```bash
git push origin test-ci
```

3. This will run tests but NOT deploy

### Test Deployment to Staging

Modify workflow to add a staging environment:

```yaml
deploy-staging:
  name: Deploy to Staging
  if: github.ref == 'refs/heads/develop'
  # ... same as deploy but with different secrets
```

## Security Best Practices

### 1. Rotate SSH Keys Regularly
```bash
# Generate new key
ssh-keygen -t rsa -b 4096 -f new-github-key

# Update on EC2 and GitHub secrets
```

### 2. Restrict SSH Key Access

On EC2, create a deployment-specific user:
```bash
sudo adduser github-deploy
sudo usermod -aG docker github-deploy

# Switch authorized_keys to this user
sudo su - github-deploy
mkdir -p ~/.ssh
echo "PUBLIC_KEY" > ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys
```

Update `EC2_USER` secret to `github-deploy`

### 3. Use Environment-Specific Secrets

For production vs staging:
- Use GitHub Environments
- Set up different secrets per environment
- Require approvals for production

### 4. Enable Branch Protection

1. Go to **Settings** → **Branches**
2. Add rule for `main` branch:
   - Require pull request reviews
   - Require status checks to pass
   - Include administrators

## Customization

### Add Slack Notifications

Add this step at the end of deploy job:

```yaml
- name: Notify Slack
  if: always()
  uses: slackapi/slack-github-action@v1
  with:
    webhook-url: ${{ secrets.SLACK_WEBHOOK }}
    payload: |
      {
        "text": "Deployment ${{ job.status }}: ${{ github.repository }}"
      }
```

### Add Discord Notifications

```yaml
- name: Notify Discord
  if: always()
  uses: sarisia/actions-status-discord@v1
  with:
    webhook: ${{ secrets.DISCORD_WEBHOOK }}
    status: ${{ job.status }}
```

### Run Database Migrations

Add before starting services:

```yaml
- name: Run migrations
  run: |
    ssh ${SSH_USER}@${SSH_HOST} << 'ENDSSH'
      cd /opt/tinder-app
      docker compose -f deploy/docker-compose.prod.yml exec -T backend \
        ./main migrate up
    ENDSSH
```

### Add Performance Tests

```yaml
- name: Run performance tests
  run: |
    # Install k6 or similar
    k6 run performance-tests.js
```

## Troubleshooting

### SSH Connection Failed
```
Error: Permission denied (publickey)
```

**Solution**:
1. Verify public key is in `~/.ssh/authorized_keys` on EC2
2. Check `EC2_SSH_KEY` secret contains full private key
3. Verify `EC2_HOST` and `EC2_USER` are correct

### Docker Image Load Failed
```
Error: Cannot load Docker image
```

**Solution**:
1. Check if EC2 has enough disk space: `df -h`
2. Prune old images: `docker system prune -a`
3. Check artifact upload/download succeeded

### Health Check Failed
```
Error: Health check failed after deployment
```

**Solution**:
1. SSH into EC2 and check logs:
```bash
docker compose -f /opt/tinder-app/deploy/docker-compose.prod.yml logs
```
2. Verify environment variables in `.env.production`
3. Check if database is accessible

### Workflow Permissions
```
Error: Resource not accessible by integration
```

**Solution**:
1. Go to **Settings** → **Actions** → **General**
2. Set "Workflow permissions" to "Read and write permissions"

## Advanced: Matrix Builds

Test multiple Go versions:

```yaml
test:
  strategy:
    matrix:
      go-version: ['1.20', '1.21', '1.22']
  steps:
    - uses: actions/setup-go@v5
      with:
        go-version: ${{ matrix.go-version }}
```

## Cost Optimization

### Reduce Build Time

1. **Cache Go modules**: Already implemented
2. **Cache Docker layers**:
```yaml
- uses: docker/build-push-action@v5
  with:
    cache-from: type=gha
    cache-to: type=gha,mode=max
```

### Reduce EC2 Costs

1. Use EC2 Auto Scaling
2. Consider AWS Fargate for containers
3. Use spot instances for dev/staging

## Monitoring & Alerts

### GitHub Actions Monitoring

Set up email notifications:
1. Go to **Settings** → **Notifications**
2. Enable "Actions" notifications

### CloudWatch Integration

Add CloudWatch logs to workflow:

```yaml
- name: Send logs to CloudWatch
  uses: aws-actions/configure-aws-credentials@v4
  with:
    aws-access-key-id: ${{ secrets.AWS_ACCESS_KEY_ID }}
    aws-secret-access-key: ${{ secrets.AWS_SECRET_ACCESS_KEY }}
    aws-region: us-east-1
```

---

## Quick Reference

### Workflow File Location
```
.github/workflows/deploy.yml
```

### Required Secrets
- `EC2_SSH_KEY`
- `EC2_HOST`
- `EC2_USER`
- `DOMAIN_NAME` (optional)

### Manual Trigger
Actions → Deploy to Production → Run workflow

### View Logs
Actions → Click workflow run → Click job → View logs

---

**Last Updated**: January 2026
