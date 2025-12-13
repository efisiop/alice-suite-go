# Testing the Development Server

## Quick Test Commands

### 1. Check if server is running:
```bash
make check
```

### 2. Build the server:
```bash
make build
```

### 3. Start the server (recommended):
```bash
make start
```
This will:
- ✅ Check if server is already running
- ✅ Configure Safari proxy bypass automatically
- ✅ Build the server
- ✅ Start the server on port 8080

### 4. Test in Safari:
After running `make start`, open Safari and go to:
```
http://127.0.0.1:8080/reader/login
```

### 5. Test with curl (bypasses proxy):
```bash
curl --noproxy '*' http://127.0.0.1:8080/health
```

## Troubleshooting

### Server won't start:
```bash
# Check if port 8080 is already in use
lsof -i :8080

# Stop any existing server
make stop

# Try starting again
make start
```

### Safari still can't connect:
```bash
# Manually configure proxy bypass
make proxy-setup

# Or manually:
sudo networksetup -setproxybypassdomains 'Wi-Fi' '127.0.0.1' 'localhost' '*.local'
```

### Check proxy configuration:
```bash
# Check current bypass domains
networksetup -getproxybypassdomains 'Wi-Fi'

# Should include: 127.0.0.1, localhost, *.local
```

