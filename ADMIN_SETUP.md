# How to start an Administrator account

Follow these steps in your project folder. Use your normal terminal (e.g. Terminal.app or the one in Cursor).

---

## Step 1: Open the project folder in the terminal

```bash
cd /Users/efisiopittau/Project_1/alice-suite-go
```

---

## Step 2: Run database migrations (if you haven’t already)

This updates the database so it allows the `administrator` role (migration 013).

```bash
go run ./cmd/migrate
```

You should see lines like:
- `Running migration: 001_initial_schema.sql`
- …
- `Running migration: 013_add_administrator_role.sql`
- `All migrations completed successfully!`

---

## Step 3: Create the test users (including the admin)

This creates the reader, consultant, and **administrator** test accounts.

```bash
go run ./cmd/init-users
```

You should see something like:
- `Created administrator user: admin@example.com (Password: admin123)`
- or `Administrator user already exists: admin@example.com` if it was already created.

At the end it will list:
- **Administrator: admin@example.com / admin123**

---

## Step 4: Start the server

```bash
make start
```

Or, if you prefer:

```bash
go build -o bin/server ./cmd/server
./bin/server
```

The server usually runs at **http://127.0.0.1:8080** (or the port shown in the terminal).

---

## Step 5: Log in as administrator

1. In your browser, open: **http://127.0.0.1:8080/admin/login**
2. Use:
   - **Email:** `admin@example.com`
   - **Password:** `admin123`
3. Click **Login**.

You should be redirected to the Administrator dashboard at **http://127.0.0.1:8080/admin**.

---

## Summary of commands (copy-paste)

From the project root:

```bash
cd /Users/efisiopittau/Project_1/alice-suite-go
go run ./cmd/migrate
go run ./cmd/init-users
make start
```

Then open **http://127.0.0.1:8080/admin/login** and sign in with **admin@example.com** / **admin123**.

---

## If something goes wrong

- **“Database not found”** when running `init-users`  
  Run `go run ./cmd/migrate` first (Step 2).

- **“role IN ('reader', 'consultant')” or similar error when running init-users**  
  Migration 013 did not run. Run `go run ./cmd/migrate` again so that 013_add_administrator_role.sql is applied.

- **“This account is not an administrator” on login**  
  The user in the database is not an administrator. Run `go run ./cmd/init-users` again so the admin user is created (or updated if you changed the code to create it).

- **Different database path**  
  If you use a custom DB path, set it before running migrate and init-users, for example:
  ```bash
  export DB_PATH=/path/to/your/alice-suite.db
  go run ./cmd/migrate
  go run ./cmd/init-users
  ```
