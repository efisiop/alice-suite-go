# Baby Steps: Start Alice on Your Mac

Do these steps **in order**. Copy each line into Terminal and press Enter. Wait for it to finish before going to the next step.

---

## Step 1: Open Terminal

- Press **Cmd + Space**, type **Terminal**, press **Enter**.
- A white or black window opens. That’s Terminal.

---

## Step 2: Tell Go to work with your network (fix the certificate error)

Copy this **whole line** and paste it into Terminal. Press **Enter**.

```
export GODEBUG=x509negativeserial=1
```

**You should see:** Nothing special—just a new line. That’s OK.

---

## Step 3: Go to the project folder

Copy this line and press **Enter**.

```
cd /Users/efisiopittau/Project_1/alice-suite-go
```

**You should see:** The line at the top of Terminal might now end with `alice-suite-go %`.

---

## Step 3b: Fix go.mod (do this once if you ever saw “updates to go.mod needed”)

Copy this line and press **Enter**.

```
go mod tidy
```

**You should see:** A few lines about downloading, then the prompt comes back. No error.

(If you never saw “go mod tidy” or “updates to go.mod needed”, you can skip this step. The start script also runs it for you.)

---

## Step 4: Run the database setup (first time only)

Copy this line and press **Enter**.

```
go run ./cmd/migrate
```

**You should see:** Some lines about “migration” or “database”. No line saying “Error” or “failed”.

If you see **TLS** or **x509** or **negative serial number**: go back to Step 2 in **this same Terminal window** and run Step 2 again, then try Step 4 again.

---

## Step 5: Create test users (first time only)

Copy this line and press **Enter**.

```
go run ./cmd/init-users
```

**You should see:** Something like “user already exists” or “created”. No red error.

---

## Step 6: Start the server

Copy this line and press **Enter**.

```
./start_dev_server.sh
```

**You should see:**  
- “Building server…”  
- “✅ Server built successfully”  
- “Starting server on port 8080…”  
- “Access at: http://127.0.0.1:8080/reader/login”  
- Then the window will **not** return to the prompt; that’s normal. The server is running.

**Do not close this Terminal window** while you use the app.

---

## Step 7: Open the app in your browser

1. Open **Safari** (or Chrome).
2. Click the address bar at the top.
3. Type or paste: **http://127.0.0.1:8080/reader/login**
4. Press **Enter**.

**You should see:** A login page for the Alice Reader app.

---

## Step 8: Log in (test account)

On the login page use:

- **Email:** `reader@example.com`
- **Password:** `reader123`
- **Verification code:** `ALICE2024`

Then click **Login**.

---

## When you want to stop the server

1. Click the **Terminal** window (where the server is running).
2. Press **Ctrl + C** once.
3. The server stops and you get the prompt back.

---

## If something goes wrong

- **“command not found: go”**  
  You need to install Go: https://go.dev/dl/ — get the **macOS** installer, run it, then close and reopen Terminal and start from Step 2.

- **Same TLS / x509 error in Step 4 or 6**  
  In the **same** Terminal window, run Step 2 again, then try the step that failed again.

- **“Permission denied” on start_dev_server.sh**  
  Run this once, then try Step 6 again:
  ```
  chmod +x /Users/efisiopittau/Project_1/alice-suite-go/start_dev_server.sh
  ```

- **Port 8080 already in use**  
  Something else is using that port. Restart your Mac and do from Step 2 again, or tell me and we can use another port.

- **“go: updates to go.mod needed; to update it: go mod tidy”**  
  In the same Terminal (after Step 2 and Step 3), run: `go mod tidy` then press Enter. Then try Step 6 again.

- **“make setup” says setup_safari_proxy.sh not found**  
  That’s OK. Use `make setup` only if you need to configure Safari’s proxy bypass; otherwise you can ignore it. The script lives in `archive/scripts/setup_safari_proxy.sh` and the Makefile now looks there.

---

**Summary:**  
Steps 2 and 3 every time you open a **new** Terminal.  
Steps 4 and 5 only the **first** time (or if you reset the app).  
Step 6 to start, Step 7 to open in the browser, Ctrl+C to stop.
