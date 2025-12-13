# Deployment Guide

**Date:** 2025-01-23  
**Purpose:** Guide for deploying Alice Suite Go application

---

## Quick Start

### 1. Build the Application

```bash
cd /Users/efisiopittau/Project_1/alice-suite-go
go build -o alice-suite-server ./cmd/server
```

### 2. Ensure Database Exists

```bash
# Database should be at: data/alice-suite.db
# If it doesn't exist, run migrations first
./alice-suite-server migrate  # If migrate command exists
```

### 3. Run the Server

```bash
# Default port: 8080
./alice-suite-server

# Or specify port
PORT=3000 ./alice-suite-server

# Or specify database path
DB_PATH=/path/to/database.db ./alice-suite-server
```

### 4. Access Applications

- **Reader App:** http://localhost:8080/
- **Consultant Dashboard:** http://localhost:8080/consultant/login
- **Health Check:** http://localhost:8080/health
- **API:** http://localhost:8080/rest/v1/

---

## Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `PORT` | `8080` | Server port |
| `DB_PATH` | `data/alice-suite.db` | Database file path |
| `JWT_SECRET` | (default secret) | JWT signing secret (change in production!) |

---

## Production Deployment

### 1. Set Production Environment Variables

```bash
export PORT=8080
export DB_PATH=/var/lib/alice-suite/alice-suite.db
export JWT_SECRET="your-secure-random-secret-key-here"
```

### 2. Build for Production

```bash
# Build optimized binary
go build -ldflags="-s -w" -o alice-suite-server ./cmd/server

# Or build for specific OS/architecture
GOOS=linux GOARCH=amd64 go build -o alice-suite-server-linux ./cmd/server
```

### 3. Create Systemd Service (Linux)

Create `/etc/systemd/system/alice-suite.service`:

```ini
[Unit]
Description=Alice Suite Go Server
After=network.target

[Service]
Type=simple
User=alice-suite
WorkingDirectory=/opt/alice-suite-go
ExecStart=/opt/alice-suite-go/alice-suite-server
Restart=always
RestartSec=5
Environment="PORT=8080"
Environment="DB_PATH=/var/lib/alice-suite/alice-suite.db"
Environment="JWT_SECRET=your-secret-key"

[Install]
WantedBy=multi-user.target
```

Enable and start:

```bash
sudo systemctl enable alice-suite
sudo systemctl start alice-suite
sudo systemctl status alice-suite
```

### 4. Reverse Proxy (Nginx)

Create `/etc/nginx/sites-available/alice-suite`:

```nginx
server {
    listen 80;
    server_name alice-suite.example.com;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # SSE support
    location /api/realtime/events {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_http_version 1.1;
        proxy_set_header Connection "";
        proxy_buffering off;
        proxy_cache off;
    }
}
```

Enable:

```bash
sudo ln -s /etc/nginx/sites-available/alice-suite /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 5. HTTPS Setup (Let's Encrypt)

```bash
sudo certbot --nginx -d alice-suite.example.com
```

---

## Database Setup

### Initial Migration

```bash
# Run migrations
cd /Users/efisiopittau/Project_1/alice-suite-go
go run ./cmd/migrate
```

### Database Location

- **Development:** `data/alice-suite.db` (relative to project root)
- **Production:** `/var/lib/alice-suite/alice-suite.db` (or custom path)

### Backup

```bash
# Simple backup
cp data/alice-suite.db data/alice-suite.db.backup

# Or use SQLite backup command
sqlite3 data/alice-suite.db ".backup data/alice-suite.db.backup"
```

---

## Monitoring

### Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "ok",
  "message": "Alice Suite Reader API - Physical Book Companion",
  "version": "1.0.0"
}
```

### Logs

Logs are written to stdout/stderr. For production, redirect to log files:

```bash
./alice-suite-server >> /var/log/alice-suite/app.log 2>&1
```

Or use systemd journal:

```bash
journalctl -u alice-suite -f
```

---

## Troubleshooting

### Server Won't Start

1. Check if port is already in use:
   ```bash
   lsof -i :8080
   ```

2. Check database file exists and is readable:
   ```bash
   ls -l data/alice-suite.db
   ```

3. Check permissions:
   ```bash
   chmod 644 data/alice-suite.db
   ```

### Database Errors

1. Check foreign keys are enabled:
   ```sql
   PRAGMA foreign_keys;
   ```

2. Verify schema:
   ```bash
   sqlite3 data/alice-suite.db ".schema"
   ```

### Real-time Features Not Working

1. Check SSE endpoint:
   ```bash
   curl -N "http://localhost:8080/api/realtime/events?token=YOUR_TOKEN"
   ```

2. Check browser console for errors
3. Verify token is valid

---

## File Structure

```
alice-suite-go/
├── alice-suite-server          # Compiled binary
├── cmd/server/main.go          # Server entry point
├── internal/
│   ├── handlers/              # HTTP handlers
│   ├── templates/             # HTML templates
│   ├── static/                # CSS, JS, images
│   ├── database/              # Database layer
│   ├── realtime/              # Real-time features
│   └── query/                 # Query parsing
├── data/
│   └── alice-suite.db         # SQLite database
└── migrations/                # Database migrations
```

---

## Single Binary Deployment

The entire application is self-contained in a single binary:

```bash
# Copy binary and database to server
scp alice-suite-server user@server:/opt/alice-suite-go/
scp data/alice-suite.db user@server:/opt/alice-suite-go/data/

# Run on server
cd /opt/alice-suite-go
./alice-suite-server
```

No additional dependencies required!

---

## Performance Tuning

### Database Optimization

```sql
-- Enable WAL mode for better concurrency
PRAGMA journal_mode=WAL;

-- Increase cache size
PRAGMA cache_size=10000;

-- Enable foreign keys
PRAGMA foreign_keys=ON;
```

### Server Configuration

- Adjust `PORT` for your environment
- Set `JWT_SECRET` to a secure random string
- Use reverse proxy (Nginx) for SSL termination
- Enable gzip compression in Nginx

---

## Security Checklist

- [ ] Change `JWT_SECRET` from default
- [ ] Use HTTPS in production
- [ ] Set proper file permissions
- [ ] Enable firewall rules
- [ ] Regular database backups
- [ ] Monitor logs for errors
- [ ] Keep Go version updated
- [ ] Review CORS settings

---

## Support

For issues or questions:
1. Check logs: `journalctl -u alice-suite -n 100`
2. Check health endpoint: `curl http://localhost:8080/health`
3. Review migration guide: `MIGRATION_TO_GO_COMPLETE.md`

---

**Deployment Status:** Ready

