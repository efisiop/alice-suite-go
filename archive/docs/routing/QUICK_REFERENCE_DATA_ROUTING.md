# Quick Reference: Data Routing Rules

**For:** Developers implementing new features  
**See:** `ARCHITECTURE_DATA_ROUTING.md` for full details

---

## 🎯 Quick Rules

### Reader Endpoints
```go
// ✅ DO THIS
claims, _ := auth.ValidateJWT(token)
userID := claims.UserID  // From token, NOT request body
query := `SELECT * FROM table WHERE user_id = ?`
db.Query(query, userID)
```

### Consultant Endpoints
```go
// ✅ DO THIS
// Use middleware.RequireConsultant
query := `
    SELECT i.*, u.first_name, u.last_name, u.email
    FROM interactions i
    JOIN users u ON i.user_id = u.id
    WHERE u.role = 'reader'
`
```

---

## 🚫 Never Do This

```go
// ❌ NEVER trust user_id from request body
userID := req.UserID

// ❌ NEVER query without user_id filter (reader endpoints)
query := `SELECT * FROM interactions`

// ❌ NEVER query without role filter (consultant endpoints)
query := `SELECT * FROM interactions`  // Missing WHERE u.role = 'reader'
```

---

## 📋 Checklist

- [ ] Reader endpoint extracts `user_id` from token
- [ ] Reader endpoint filters by `user_id = claims.UserID`
- [ ] Consultant endpoint uses `middleware.RequireConsultant`
- [ ] Consultant endpoint filters by `role = 'reader'`
- [ ] Database query joins with `users` table
- [ ] No `user_id` in request body is trusted

---

## 🔍 Route Patterns

| Pattern | User Type | Protection | Data Scope |
|---------|-----------|------------|------------|
| `/reader/*` | Reader | RequireAuth | Own data only |
| `/consultant/*` | Consultant | RequireConsultant | All reader data |
| `/api/*` | Reader | RequireAuth | Own data only |
| `/api/consultant/*` | Consultant | RequireConsultant | All reader data |

---

## 🗄️ Database Tables

All user-specific tables MUST have `user_id`:
- `interactions`
- `reading_progress`
- `help_requests`
- `vocabulary_lookups`
- `ai_interactions`
- `reading_stats`

---

## 🔐 Security Rules

1. **Always** extract `user_id` from JWT token
2. **Never** trust `user_id` from request body
3. **Always** filter by `user_id` for reader endpoints
4. **Always** filter by `role = 'reader'` for consultant endpoints
5. **Always** join with `users` table to get role

---

**Questions?** See `ARCHITECTURE_DATA_ROUTING.md` for detailed examples.

