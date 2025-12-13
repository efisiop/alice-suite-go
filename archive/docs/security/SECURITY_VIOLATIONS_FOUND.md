# Security Violations Found - Data Routing Rules

**Date:** 2025-01-23  
**Status:** ⚠️ REQUIRES FIXING  
**Priority:** HIGH

---

## Summary

Found **7 security violations** where `user_id` is extracted from request body instead of JWT token. These violate the data routing rules and must be fixed.

---

## Violations Found

### 1. `HandleTrackActivity` - `internal/handlers/activity.go:136`
**Issue:** Uses `req.UserID` from request body  
**Risk:** Users can track activities for other users  
**Fix Required:** Extract `user_id` from JWT token

```go
// ❌ CURRENT (WRONG)
err := TrackActivity(req.UserID, req.EventType, req.BookID, data)

// ✅ SHOULD BE
claims, _ := auth.ValidateJWT(token)
err := TrackActivity(claims.UserID, req.EventType, req.BookID, data)
```

### 2. `HandleHelpRequests` (POST) - `internal/handlers/api.go:359`
**Issue:** Uses `req.UserID` from request body  
**Risk:** Users can create help requests for other users  
**Fix Required:** Extract `user_id` from JWT token

```go
// ❌ CURRENT (WRONG)
request, err := helpService.CreateHelpRequest(req.UserID, req.BookID, ...)

// ✅ SHOULD BE
claims, _ := auth.ValidateJWT(token)
request, err := helpService.CreateHelpRequest(claims.UserID, req.BookID, ...)
```

### 3. `HandleLookupWord` - `internal/handlers/api.go:452,457`
**Issue:** Uses `req.UserID` from request body  
**Risk:** Users can record lookups for other users  
**Fix Required:** Extract `user_id` from JWT token

```go
// ❌ CURRENT (WRONG)
if req.UserID != "" {
    dictionaryService.RecordLookup(req.UserID, ...)
}

// ✅ SHOULD BE
claims, _ := auth.ValidateJWT(token)
dictionaryService.RecordLookup(claims.UserID, ...)
```

### 4. `HandleAskAI` - `internal/handlers/api.go:496`
**Issue:** Uses `req.UserID` from request body  
**Risk:** Users can create AI interactions for other users  
**Fix Required:** Extract `user_id` from JWT token

```go
// ❌ CURRENT (WRONG)
interaction, err := aiService.AskAI(req.UserID, req.BookID, ...)

// ✅ SHOULD BE
claims, _ := auth.ValidateJWT(token)
interaction, err := aiService.AskAI(claims.UserID, req.BookID, ...)
```

### 5. `HandleCreateHelpRequest` - `internal/handlers/api.go:530`
**Issue:** Uses `req.UserID` from request body  
**Risk:** Users can create help requests for other users  
**Fix Required:** Extract `user_id` from JWT token

```go
// ❌ CURRENT (WRONG)
request, err := helpService.CreateHelpRequest(req.UserID, req.BookID, ...)

// ✅ SHOULD BE
claims, _ := auth.ValidateJWT(token)
request, err := helpService.CreateHelpRequest(claims.UserID, req.BookID, ...)
```

### 6. `HandleVerifyBookCode` - `internal/handlers/verification.go:28`
**Issue:** Uses `req.UserID` from request body  
**Risk:** Users can verify codes for other users  
**Fix Required:** Extract `user_id` from JWT token

```go
// ❌ CURRENT (WRONG)
bookID, err := auth.VerifyBookCode(req.Code, req.UserID)

// ✅ SHOULD BE
claims, _ := auth.ValidateJWT(token)
bookID, err := auth.VerifyBookCode(req.Code, claims.UserID)
```

### 7. `HandleCheckBookVerified` - `internal/handlers/verification.go:56,62`
**Issue:** Uses `user_id` from query parameter  
**Risk:** Users can check verification status for other users  
**Fix Required:** Extract `user_id` from JWT token

```go
// ❌ CURRENT (WRONG)
userID := r.URL.Query().Get("user_id")
verified, err := auth.CheckBookVerified(userID)

// ✅ SHOULD BE
claims, _ := auth.ValidateJWT(token)
verified, err := auth.CheckBookVerified(claims.UserID)
```

---

## Impact

### Security Risks
1. **Data Manipulation:** Users can create/modify data for other users
2. **Privacy Violation:** Users can access other users' verification status
3. **Data Integrity:** Incorrect user associations in database
4. **Audit Trail:** Incorrect tracking of who performed actions

### Business Impact
- Loss of user trust
- Data corruption
- Compliance violations
- Potential legal issues

---

## Fix Priority

**HIGH PRIORITY** - These violations allow users to:
- Track activities for other users
- Create help requests for other users
- Record vocabulary lookups for other users
- Create AI interactions for other users
- Verify book codes for other users
- Check verification status for other users

---

## Recommended Fixes

### Pattern to Follow

```go
func HandleEndpoint(w http.ResponseWriter, r *http.Request) {
    // 1. Extract token
    authHeader := r.Header.Get("Authorization")
    token, err := auth.ExtractTokenFromHeader(authHeader)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // 2. Validate token and get user_id
    claims, err := auth.ValidateJWT(token)
    if err != nil {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // 3. Use claims.UserID (NOT from request body)
    userID := claims.UserID
    
    // 4. Use userID in service calls
    result, err := service.DoSomething(userID, ...)
}
```

### Middleware Option

Alternatively, use middleware to extract user_id:

```go
// In middleware/auth.go
func RequireAuthWithUserID(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Extract and validate token
        claims, err := getClaimsFromRequest(r)
        if err != nil {
            http.Error(w, "Unauthorized", http.StatusUnauthorized)
            return
        }
        
        // Add user_id to request context
        ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
        r = r.WithContext(ctx)
        
        next(w, r)
    }
}

// In handler
userID := r.Context().Value("user_id").(string)
```

---

## Testing After Fix

After fixing, test:
1. ✅ User cannot track activity for another user
2. ✅ User cannot create help request for another user
3. ✅ User cannot record lookup for another user
4. ✅ User cannot create AI interaction for another user
5. ✅ User cannot verify code for another user
6. ✅ User cannot check verification for another user

---

## Related Documentation

- `ARCHITECTURE_DATA_ROUTING.md` - Full architecture rules
- `QUICK_REFERENCE_DATA_ROUTING.md` - Quick reference guide

---

**Status:** ⚠️ **REQUIRES IMMEDIATE FIX**

