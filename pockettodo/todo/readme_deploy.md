# Production Deployment Guide

This guide covers deploying the PocketTodo application to a Debian 12 server with Caddy web server using the domain `todo.rasc.ch`.

## Prerequisites

- Debian 12 server with root/sudo access
- Domain `todo.rasc.ch` pointing to your server IP
- Caddy web server installed
- Basic knowledge of Linux command line

## Deployment Architecture

```
Internet → Caddy (Reverse Proxy) → Angular App (Static Files)
                                 → PocketBase API (:8090)
```

## 1. Server Preparation

### Update System

```bash
sudo apt update && sudo apt upgrade -y
```

### Install Required Packages

```bash
# Install essential tools
sudo apt install -y curl wget unzip git htop nano systemd

# Install Node.js (for building Angular app)
curl -fsSL https://deb.nodesource.com/setup_20.x | sudo -E bash -
sudo apt install -y nodejs

# Verify installations
node --version
npm --version
caddy version
```

### Create Application User

```bash
# Create dedicated user for the application
sudo useradd -m -s /bin/bash todo
sudo usermod -aG sudo todo

# Switch to application user
sudo su - todo
```

## 2. Deploy PocketBase

### Download and Setup PocketBase

```bash
# Create application directory
mkdir -p /home/todo/pocketbase
cd /home/todo/pocketbase

# Download PocketBase (check for latest version at https://github.com/pocketbase/pocketbase/releases)
wget https://github.com/pocketbase/pocketbase/releases/download/v0.22.0/pocketbase_0.22.0_linux_amd64.zip

# Extract
unzip pocketbase_0.22.0_linux_amd64.zip
chmod +x pocketbase

# Clean up
rm pocketbase_0.22.0_linux_amd64.zip

# Create data directory
mkdir -p pb_data
```

### Create PocketBase Systemd Service

```bash
sudo nano /etc/systemd/system/pocketbase.service
```

Add the following content:

```ini
[Unit]
Description=PocketBase
After=network.target

[Service]
Type=simple
User=todo
Group=todo
WorkingDirectory=/home/todo/pocketbase
ExecStart=/home/todo/pocketbase/pocketbase serve --http=127.0.0.1:8090
Restart=always
RestartSec=5
StandardOutput=journal
StandardError=journal
SyslogIdentifier=pocketbase

# Security settings
NoNewPrivileges=yes
PrivateTmp=yes
ProtectSystem=strict
ReadWritePaths=/home/todo/pocketbase
ProtectHome=yes

[Install]
WantedBy=multi-user.target
```

### Start and Enable PocketBase Service

```bash
# Reload systemd and start service
sudo systemctl daemon-reload
sudo systemctl enable pocketbase
sudo systemctl start pocketbase

# Check status
sudo systemctl status pocketbase

# View logs
sudo journalctl -u pocketbase -f
```

## 3. Deploy Angular Application

### Build Application Locally

On your development machine:

```bash
# Navigate to your project
cd /path/to/your/todo/project

# Install dependencies (if not already done)
npm install

# Update production environment
nano src/environments/environment.prod.ts
```

Update the production environment:

```typescript
export const environment = {
  production: true,
  pocketbaseUrl: 'https://todo.rasc.ch/api'
};
```

Build for production:

```bash
# Build the application
npm run build

# This creates a dist/ folder with static files
```

### Upload to Server

```bash
# Create web directory on server
ssh todo@your-server-ip
mkdir -p /home/todo/www

# From your local machine, upload the built files
scp -r dist/* todo@your-server-ip:/home/todo/www/

# Or using rsync (recommended)
rsync -avz --delete dist/ todo@your-server-ip:/home/todo/www/
```

### Alternative: Build on Server

```bash
# On the server as todo user
cd /home/todo
git clone https://github.com/yourusername/your-todo-repo.git app
cd app

# Install dependencies
npm install

# Build for production
npm run build

# Copy built files to web directory
cp -r dist/* /home/todo/www/
```

## 4. Configure Caddy

### Create Caddyfile

```bash
sudo nano /etc/caddy/Caddyfile
```

Add the following configuration:

```caddyfile
# PocketTodo Application
todo.rasc.ch {
    # Serve Angular static files
    root * /home/todo/www

    # Enable gzip compression
    encode gzip

    # Security headers
    header {
        # HSTS
        Strict-Transport-Security "max-age=31536000; includeSubDomains; preload"
        # Prevent clickjacking
        X-Frame-Options "DENY"
        # Prevent MIME type sniffing
        X-Content-Type-Options "nosniff"
        # XSS Protection
        X-XSS-Protection "1; mode=block"
        # Content Security Policy
        Content-Security-Policy "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; font-src 'self'; connect-src 'self' https://todo.rasc.ch/api"
    }

    # Handle PocketBase API requests
    handle_path /api/* {
        reverse_proxy 127.0.0.1:8090
    }

    # Handle Angular routing (SPA)
    try_files {path} /index.html

    # Cache static assets
    @static {
        file
        path *.js *.css *.png *.jpg *.jpeg *.gif *.ico *.svg *.woff *.woff2
    }
    header @static Cache-Control "public, max-age=31536000"

    # Logging
    log {
        output file /var/log/caddy/todo.rasc.ch.log
        format json
    }
}
```

### Test and Reload Caddy Configuration

```bash
# Test configuration
sudo caddy validate --config /etc/caddy/Caddyfile

# Reload Caddy
sudo systemctl reload caddy

# Check Caddy status
sudo systemctl status caddy

# View Caddy logs
sudo journalctl -u caddy -f
```

## 5. Configure PocketBase

### Initial Setup

1. Open your browser and navigate to `https://todo.rasc.ch/api/_/`
2. Create an admin account
3. Set up the database schema as described in `readme_pocketbase.md`

### Configure Settings

#### Application Settings

- **Application name**: PocketTodo
- **Application URL**: `https://todo.rasc.ch`

#### CORS Settings

- **Allowed origins**: `https://todo.rasc.ch`

#### Email Settings

Configure SMTP for password reset functionality:

- **SMTP server**: Your email provider's SMTP server
- **Port**: 587 (TLS) or 465 (SSL)
- **Username**: Your email address
- **Password**: Your email password or app-specific password

## 6. SSL and Security

### Verify SSL Certificate

```bash
# Check SSL certificate
openssl s_client -connect todo.rasc.ch:443 -servername todo.rasc.ch

# Check certificate expiry
echo | openssl s_client -connect todo.rasc.ch:443 2>/dev/null | openssl x509 -noout -dates
```

### Firewall Configuration

```bash
# Install and configure UFW
sudo apt install ufw
sudo ufw default deny incoming
sudo ufw default allow outgoing
sudo ufw allow ssh
sudo ufw allow 80/tcp
sudo ufw allow 443/tcp
sudo ufw enable

# Check status
sudo ufw status
```

## 7. Monitoring and Maintenance

### System Monitoring

```bash
# Check system resources
htop
df -h
free -h

# Check service status
sudo systemctl status pocketbase caddy

# View application logs
sudo journalctl -u pocketbase -f
sudo journalctl -u caddy -f
tail -f /var/log/caddy/todo.rasc.ch.log
```

### Log Rotation

```bash
# Configure logrotate for application logs
sudo nano /etc/logrotate.d/todo-app
```

Add the following:

```
/var/log/caddy/todo.rasc.ch.log {
    daily
    rotate 30
    compress
    delaycompress
    missingok
    create 0644 caddy caddy
    postrotate
        systemctl reload caddy
    endscript
}
```

## 8. Backup Strategy

### Automated Backup Script

```bash
# Create backup directory
sudo mkdir -p /home/todo/backups
sudo chown todo:todo /home/todo/backups

# Create backup script
nano /home/todo/backup.sh
```

Add the following content:

```bash
#!/bin/bash

# Backup script for PocketTodo application
# Run daily via cron

BACKUP_DIR="/home/todo/backups"
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_NAME="pockettodo_backup_$DATE"

# Create backup directory
mkdir -p "$BACKUP_DIR/$BACKUP_NAME"

echo "Starting backup at $(date)"

# Backup PocketBase data
echo "Backing up PocketBase data..."
cp -r /home/todo/pocketbase/pb_data "$BACKUP_DIR/$BACKUP_NAME/"

# Backup Angular application files
echo "Backing up web files..."
cp -r /home/todo/www "$BACKUP_DIR/$BACKUP_NAME/"

# Backup Caddy configuration
echo "Backing up Caddy configuration..."
cp /etc/caddy/Caddyfile "$BACKUP_DIR/$BACKUP_NAME/"

# Create compressed archive
echo "Creating compressed archive..."
cd "$BACKUP_DIR"
tar -czf "$BACKUP_NAME.tar.gz" "$BACKUP_NAME"
rm -rf "$BACKUP_NAME"

# Keep only last 7 days of backups
echo "Cleaning old backups..."
find "$BACKUP_DIR" -name "pockettodo_backup_*.tar.gz" -mtime +7 -delete

echo "Backup completed at $(date)"
echo "Backup saved as: $BACKUP_DIR/$BACKUP_NAME.tar.gz"

# Optional: Upload to remote storage (S3, etc.)
# aws s3 cp "$BACKUP_DIR/$BACKUP_NAME.tar.gz" s3://your-backup-bucket/
```

Make the script executable:

```bash
chmod +x /home/todo/backup.sh
```

### Configure Cron for Automated Backups

```bash
# Edit crontab
crontab -e

# Add daily backup at 3 AM
0 3 * * * /home/todo/backup.sh >> /home/todo/backup.log 2>&1
```

### Manual Backup Commands

```bash
# Quick manual backup
cd /home/todo
tar -czf "manual_backup_$(date +%Y%m%d_%H%M%S).tar.gz" \
    pocketbase/pb_data \
    www \
    /etc/caddy/Caddyfile

# Export PocketBase data
cd /home/todo/pocketbase
./pocketbase admin export backup_$(date +%Y%m%d_%H%M%S).zip
```

### Restore from Backup

```bash
# Stop services
sudo systemctl stop pocketbase caddy

# Restore PocketBase data
cd /home/todo/backups
tar -xzf pockettodo_backup_YYYYMMDD_HHMMSS.tar.gz
cp -r pockettodo_backup_YYYYMMDD_HHMMSS/pb_data/* /home/todo/pocketbase/pb_data/

# Restore web files
cp -r pockettodo_backup_YYYYMMDD_HHMMSS/www/* /home/todo/www/

# Restore Caddy config
sudo cp pockettodo_backup_YYYYMMDD_HHMMSS/Caddyfile /etc/caddy/

# Set correct permissions
sudo chown -R todo:todo /home/todo/pocketbase/pb_data
sudo chown -R todo:todo /home/todo/www
sudo chown root:root /etc/caddy/Caddyfile

# Start services
sudo systemctl start pocketbase caddy
sudo systemctl status pocketbase caddy
```

## 9. Deployment Updates

### Update Angular Application

```bash
# Build new version locally
npm run build

# Upload to server
rsync -avz --delete dist/ todo@your-server-ip:/home/todo/www/

# No restart needed - static files are served directly
```

### Update PocketBase

```bash
# Download new version
cd /home/todo/pocketbase
wget https://github.com/pocketbase/pocketbase/releases/download/vX.X.X/pocketbase_X.X.X_linux_amd64.zip

# Stop service
sudo systemctl stop pocketbase

# Backup current version
cp pocketbase pocketbase.backup

# Extract new version
unzip pocketbase_X.X.X_linux_amd64.zip
chmod +x pocketbase

# Start service
sudo systemctl start pocketbase

# Check status
sudo systemctl status pocketbase
```

## 10. Troubleshooting

### Common Issues

#### PocketBase Connection Issues

```bash
# Check if PocketBase is running
sudo systemctl status pocketbase
sudo netstat -tlnp | grep 8090

# Check logs
sudo journalctl -u pocketbase -n 50
```

#### Caddy Issues

```bash
# Test configuration
sudo caddy validate --config /etc/caddy/Caddyfile

# Check logs
sudo journalctl -u caddy -n 50
tail -f /var/log/caddy/todo.rasc.ch.log
```

#### SSL Certificate Issues

```bash
# Force certificate renewal
sudo caddy reload --config /etc/caddy/Caddyfile

# Check certificate status
curl -I https://todo.rasc.ch
```

#### Permission Issues

```bash
# Fix file permissions
sudo chown -R todo:todo /home/todo/
sudo chmod -R 755 /home/todo/www
sudo chmod +x /home/todo/pocketbase/pocketbase
```

### Performance Monitoring

```bash
# Monitor system resources
htop
iotop
nethogs

# Check disk usage
df -h
du -sh /home/todo/*

# Monitor application performance
curl -w "@curl-format.txt" -o /dev/null -s https://todo.rasc.ch
```

### Health Check Script

```bash
# Create health check script
nano /home/todo/health-check.sh
```

Add the following:

```bash
#!/bin/bash

echo "=== PocketTodo Health Check ==="
echo "Date: $(date)"

# Check PocketBase
if systemctl is-active --quiet pocketbase; then
    echo "✓ PocketBase service is running"
else
    echo "✗ PocketBase service is not running"
fi

# Check Caddy
if systemctl is-active --quiet caddy; then
    echo "✓ Caddy service is running"
else
    echo "✗ Caddy service is not running"
fi

# Check web accessibility
if curl -s -f https://todo.rasc.ch > /dev/null; then
    echo "✓ Website is accessible"
else
    echo "✗ Website is not accessible"
fi

# Check API accessibility
if curl -s -f https://todo.rasc.ch/api/health > /dev/null; then
    echo "✓ API is accessible"
else
    echo "✗ API is not accessible"
fi

# Check disk space
DISK_USAGE=$(df /home/todo | awk 'NR==2 {print $5}' | sed 's/%//')
if [ $DISK_USAGE -lt 80 ]; then
    echo "✓ Disk usage is acceptable ($DISK_USAGE%)"
else
    echo "⚠ Disk usage is high ($DISK_USAGE%)"
fi

echo "=========================="
```

Make it executable and run:

```bash
chmod +x /home/todo/health-check.sh
./health-check.sh
```

## 11. Security Hardening

### Additional Security Measures

```bash
# Install fail2ban for SSH protection
sudo apt install fail2ban

# Configure automatic security updates
sudo apt install unattended-upgrades
sudo dpkg-reconfigure unattended-upgrades

# Set up log monitoring
sudo apt install logwatch
```

### Regular Security Updates

```bash
# Create update script
nano /home/todo/update.sh
```

Add:

```bash
#!/bin/bash
sudo apt update
sudo apt list --upgradable
sudo apt upgrade -y
sudo apt autoremove -y
```

Schedule weekly updates:

```bash
# Add to crontab
0 2 * * 1 /home/todo/update.sh >> /home/todo/update.log 2>&1
```

This deployment guide provides a robust, secure, and maintainable setup for your PocketTodo application on Debian 12 with Caddy.
