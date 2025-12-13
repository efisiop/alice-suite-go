# How to Open the Book Viewer

## Option 1: Browser (if proxy allows)
1. Open your browser
2. Go to: http://localhost:8080/viewer.html
   OR: http://127.0.0.1:8080/viewer.html

## Option 2: Bypass Proxy for Localhost
If your corporate proxy blocks localhost:
1. Open browser settings
2. Find proxy settings
3. Add "localhost" and "127.0.0.1" to bypass list
4. Then try: http://localhost:8080/viewer.html

## Option 3: DB Browser for SQLite (Recommended if proxy blocks)
1. Download: https://sqlitebrowser.org/
2. Install DB Browser for SQLite
3. Open: /Users/efisiopittau/Project_1/alice-suite-go/data/alice-suite.db
4. Go to "Browse Data" tab
5. Select "sections" table to see all book content

## Option 4: Command Line Viewer
Run: `go run cmd/viewer/main.go`

## Check Server Status
The server should be running. To verify:
- Check if port 8080 is in use: `lsof -ti:8080`
- To restart: `go run cmd/reader/main.go`
