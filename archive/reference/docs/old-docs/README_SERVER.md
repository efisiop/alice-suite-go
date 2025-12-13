# 🚀 Starting the Alice Suite Go Server

## Quick Start

The server is now running! Here's how to access it:

### Option 1: Use the Start Script
```bash
cd /Users/efisiopittau/Project_1/alice-suite-go
./START_SERVER.sh
```

This will:
- Start the server on port 8080
- Open your browser automatically
- Show server logs

### Option 2: Manual Start
```bash
cd /Users/efisiopittau/Project_1/alice-suite-go

# Start the server
go run cmd/reader/main.go
```

Then open your browser to: **http://localhost:8080**

---

## 🌐 Access Points

### Web Interface
- **Main Page:** http://localhost:8080
- Interactive test page with all API endpoints

### API Endpoints
- **Health Check:** http://localhost:8080/api/health
- **Get Books:** http://localhost:8080/api/books
- **Get Chapters:** http://localhost:8080/api/chapters?book_id=alice-in-wonderland
- **Register:** POST http://localhost:8080/api/auth/register
- **Login:** POST http://localhost:8080/api/auth/login
- **Lookup Word:** POST http://localhost:8080/api/dictionary/lookup

---

## 📋 What You'll See

When you open http://localhost:8080, you'll see:

1. **Status Indicator** - Shows if API is running
2. **Health Check** - Test API connectivity
3. **Get Books** - View available books
4. **Register User** - Create a new account
5. **Login** - Authenticate and get token
6. **Get Chapters** - View first 3 chapters
7. **Lookup Word** - Search glossary terms

---

## 🛑 Stopping the Server

Press `Ctrl+C` in the terminal where the server is running, or:

```bash
# Find and kill the server process
lsof -ti:8080 | xargs kill -9
```

---

## 📊 Server Logs

View server logs:
```bash
tail -f /tmp/alice-server.log
```

---

## ✅ Server Status

The server is currently **RUNNING** on port 8080!

**Open your browser and navigate to:**
# http://localhost:8080

You should see the interactive test page where you can test all API endpoints!



