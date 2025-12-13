# Migration Plan: React/TypeScript/Supabase → Go/SQLite

**Created:** 2025-01-08  
**Source:** `/Users/efisiopittau/Project_1/alice-suite`  
**Destination:** `/Users/efisiopittau/Project_1/alice-suite-go`

---

## 🎯 Migration Goals

1. **Fresh Start** - Build new codebase in Go from scratch
2. **Simplified Stack** - Go backend + SQLite (no cloud dependencies)
3. **Physical Book Companion** - Emphasize companion nature throughout
4. **Test Ground** - Start with first 3 chapters only
5. **Streamlined UX** - Simplified authentication and access flow

---

## 📋 Migration Phases

### Phase 1: Project Setup ✅
- [x] Create new directory: `alice-suite-go`
- [ ] Initialize Go module
- [ ] Set up project structure
- [ ] Create basic README and documentation

### Phase 2: Database Setup
- [ ] Design SQLite schema (based on current Supabase schema)
- [ ] Create migration scripts
- [ ] Load first 3 chapters of Alice in Wonderland
- [ ] Set up Alice glossary (first 3 chapters)

### Phase 3: Core Backend (Go)
- [ ] Authentication service (streamlined)
- [ ] Book/content service (first 3 chapters)
- [ ] Dictionary/glossary service
- [ ] AI service integration
- [ ] Help request service
- [ ] Progress tracking service

### Phase 4: API Layer
- [ ] REST API endpoints (or GraphQL if preferred)
- [ ] Request/response handlers
- [ ] Middleware (auth, logging, etc.)
- [ ] Error handling

### Phase 5: Frontend (To Be Determined)
- [ ] Decide: Web frontend (HTML/JS) or Go-based templating
- [ ] Implement reader interface
- [ ] Implement consultant dashboard
- [ ] Ensure physical book companion messaging

### Phase 6: Testing & Refinement
- [ ] Test with first 3 chapters
- [ ] Fix bugs
- [ ] Streamline authentication flow
- [ ] Performance optimization
- [ ] Documentation

### Phase 7: Expansion
- [ ] Add remaining chapters (full book)
- [ ] Enhanced features
- [ ] Production deployment

---

## 🔄 Data Migration

### From Supabase to SQLite

**Tables to Migrate:**
- users
- books
- book_sections (first 3 chapters only)
- alice_glossary (first 3 chapters only)
- user_interactions
- help_requests
- feedback
- consultant_assignments
- verification_codes

**Migration Script:** To be created in Go

---

## 📝 Key Differences

### Current Stack → New Stack

| Component | Current | New |
|-----------|---------|-----|
| Backend | Supabase (PostgreSQL) | Go + SQLite |
| Frontend | React + TypeScript | To be determined |
| Auth | Supabase Auth | Go-based auth |
| Real-time | WebSocket/Socket.io | To be determined (Go) |
| AI | Edge Functions | Go HTTP client |
| Database | PostgreSQL (cloud) | SQLite (local) |

---

## 🎯 Success Criteria

Migration is successful when:
- ✅ New Go codebase is functional
- ✅ First 3 chapters work seamlessly
- ✅ All core features work without bugs
- ✅ Authentication is streamlined
- ✅ Physical book companion nature is clear
- ✅ Performance is good
- ✅ Ready to expand to full book

---

**This is a fresh start. Build it right!**

