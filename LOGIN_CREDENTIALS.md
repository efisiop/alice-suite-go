# Login Credentials

**Created:** 2025-01-23  
**Purpose:** Quick reference for test user credentials

---

## Test Users

### Reader Account
- **Email:** `reader@example.com`
- **Password:** `reader123`
- **Role:** Reader
- **Verification Required:** Yes (use code below)

### Consultant Account
- **Email:** `consultant@example.com`
- **Password:** `consultant123`
- **Role:** Consultant
- **Verification Required:** No

---

## Verification Code

- **Code:** `ALICE2024`
- **Book:** Alice in Wonderland
- **Usage:** Enter this code on the verification page after registering/logging in as a reader

---

## Access URLs

- **Reader Login:** http://localhost:8080/login
- **Consultant Login:** http://localhost:8080/consultant/login
- **Reader Dashboard:** http://localhost:8080/reader (requires login + verification)
- **Consultant Dashboard:** http://localhost:8080/consultant (requires login)

---

## Quick Start

1. **Start the server:**
   ```bash
   ./start.sh
   ```

2. **Login as Reader:**
   - Go to http://localhost:8080/login
   - Email: `reader@example.com`
   - Password: `reader123`
   - After login, verify with code: `ALICE2024`

3. **Login as Consultant:**
   - Go to http://localhost:8080/consultant/login
   - Email: `consultant@example.com`
   - Password: `consultant123`

---

## Creating New Users

To create additional test users, run:

```bash
go run ./cmd/init-users
```

This will create the default test users if they don't already exist.

---

**Note:** These are test credentials. Change passwords in production!

