# Alice Suite - Detailed Requirements

**Created:** 2025-01-20  
**Status:** Requirements for Reader App (Phase 1-4)  
**Based on:** PROJECT_CONCEPT.md and current database structure

---

## 📋 Overview

This document defines detailed, testable requirements for the Alice Suite reader application. The requirements are organized by user stories, functional requirements, and non-functional requirements.

**Key Context:**
- **Database Structure:** Pages → Sections (40 words average per section)
- **Glossary:** 1,209 terms linked to sections via `glossary_section_links` table
- **Content Scope:** First 3 chapters (test ground), expandable to full book
- **Technology:** Go backend, SQLite database, web-based frontend

---

## 👥 User Stories

### US-1: User Registration
**As a** student reader  
**I want to** create an account with my name, email, and password  
**So that** I can access the companion app and track my reading progress

**Acceptance Criteria:**
- [ ] User can enter: First name, Last name, Email, Password
- [ ] Email must be valid format and unique
- [ ] Password must meet minimum security requirements (8+ characters)
- [ ] Account is created successfully
- [ ] User receives confirmation message
- [ ] Duplicate email shows appropriate error message

---

### US-2: User Login
**As a** registered student  
**I want to** log in with my email and password  
**So that** I can access my reading progress and continue where I left off

**Acceptance Criteria:**
- [ ] User can enter email and password
- [ ] Valid credentials grant access
- [ ] Invalid credentials show error message
- [ ] Session is maintained after login
- [ ] User is redirected to main reader interface

---

### US-3: Book Access Authorization
**As a** student  
**I want to** enter a unique book verification code  
**So that** I can access the Alice book content in the app

**Acceptance Criteria:**
- [ ] User can enter verification code
- [ ] Valid code grants access to book
- [ ] Invalid code shows error message
- [ ] Code can only be used once (or as per business rules)
- [ ] Access is remembered for future sessions

---

### US-4: Welcome Screen (First Time)
**As a** first-time user  
**I want to** see a welcome screen explaining how to use the app  
**So that** I understand how the companion app works with my physical book

**Acceptance Criteria:**
- [ ] Welcome screen appears only on first login
- [ ] Welcome screen explains:
  - How to use the app with physical book
  - How to look up words
  - How to ask for help
- [ ] User can dismiss welcome screen
- [ ] Welcome screen accessible from menu anytime

---

### US-5: Page Selection
**As a** student reading the physical book  
**I want to** select which page I'm currently reading  
**So that** the app shows me the corresponding content

**Acceptance Criteria:**
- [ ] Main screen displays page selection interface
- [ ] User can enter page number (1-100 for full book, 1-17 for test chapters)
- [ ] App displays page content with sections
- [ ] Page number is validated (within book range)
- [ ] Invalid page number shows error message

---

### US-6: Section Selection
**As a** student reading a specific page  
**I want to** select which section I'm reading on that page  
**So that** I can see the exact text I'm reading in the physical book

**Acceptance Criteria:**
- [ ] Page content shows all sections for that page
- [ ] Each section displays: section number and word count (~40 words)
- [ ] User can click/tap a section to view its content
- [ ] Selected section is highlighted
- [ ] Section content displays clearly formatted text

---

### US-7: Word Lookup (Glossary Priority)
**As a** student reading a section  
**I want to** look up words I don't understand  
**So that** I can continue reading with comprehension

**Acceptance Criteria:**
- [ ] Words in section text are clickable/tappable
- [ ] Clicking a word shows definition popup/modal
- [ ] If word is in glossary (1,209 terms), show glossary definition first
- [ ] Glossary definition includes: term, definition, example (if available)
- [ ] If word not in glossary, show "Definition not available" or external dictionary
- [ ] User can close definition and continue reading
- [ ] Lookup is recorded for analytics

---

### US-8: Glossary Term Highlighting
**As a** student reading a section  
**I want to** see which words have glossary definitions  
**So that** I know which words I can look up for help

**Acceptance Criteria:**
- [ ] Words that appear in glossary are visually highlighted (e.g., yellow underline)
- [ ] Hovering over highlighted word shows tooltip with definition
- [ ] Highlighting is based on `glossary_section_links` table
- [ ] Only terms linked to current section are highlighted

---

### US-9: Last Page Memory
**As a** returning student  
**I want to** resume reading from my last page  
**So that** I don't have to remember where I was

**Acceptance Criteria:**
- [ ] App saves last page number user viewed
- [ ] App saves last section number user viewed
- [ ] On login, app shows last page/section automatically
- [ ] User can navigate to different pages if needed
- [ ] Progress is saved per user per book

---

### US-10: AI Help (Tier 2)
**As a** student who needs more than a definition  
**I want to** ask AI questions about what I'm reading  
**So that** I can understand complex passages

**Acceptance Criteria:**
- [ ] User can select text from current section
- [ ] User can type a question about the selected text
- [ ] AI provides contextual answer based on:
  - Selected text
  - Current page/section
  - Chapter context
- [ ] Answer is clear and easy to understand
- [ ] Interaction is saved to database
- [ ] User can ask follow-up questions

---

### US-11: Human Help Request (Tier 3)
**As a** student still confused after AI help  
**I want to** request help from a human consultant  
**So that** I can get personalized assistance

**Acceptance Criteria:**
- [ ] User can click "Ask for Help" button
- [ ] User can type their question
- [ ] User can specify context (current page/section)
- [ ] Request is submitted and saved
- [ ] User sees confirmation message
- [ ] Request status is trackable (pending, assigned, resolved)
- [ ] User receives notification when consultant responds

---

### US-12: Reading Progress Tracking
**As a** student  
**I want to** see my reading progress  
**So that** I know how much I've read

**Acceptance Criteria:**
- [ ] App tracks pages viewed
- [ ] App tracks sections viewed
- [ ] User can see progress dashboard
- [ ] Progress shows: pages read, sections read, percentage complete
- [ ] Progress is saved per user per book

---

### US-13: Vocabulary Tracking
**As a** student  
**I want to** see words I've looked up  
**So that** I can review vocabulary I'm learning

**Acceptance Criteria:**
- [ ] App tracks all words user has looked up
- [ ] User can view vocabulary list
- [ ] Vocabulary list shows: word, definition, date looked up, context
- [ ] Words are organized by date or alphabetically
- [ ] User can filter/search vocabulary list

---

## 🔧 Functional Requirements

### FR-1: Authentication System
1. **User Registration**
   - Fields: First name, Last name, Email, Password
   - Email validation (format check)
   - Password requirements: minimum 8 characters
   - Password hashing using bcrypt
   - Duplicate email prevention

2. **User Login**
   - Email/password authentication
   - Session management (JWT tokens recommended)
   - Password verification
   - Error handling for invalid credentials

3. **Authorization**
   - Book access via verification codes
   - Code validation and tracking
   - One-time use or multi-use as per business rules

### FR-2: Content Display System
1. **Page Structure**
   - Pages organized by page number (1-100 for full book)
   - Each page contains multiple sections
   - Average 40 words per section
   - Page word count: 180-200 words (4-5 sections per page)

2. **Section Structure**
   - Sections numbered within each page (1, 2, 3, etc.)
   - Each section: ~40 words average (range: 35-45 words)
   - Section content displayed as formatted text
   - Chapter titles appear on first page of chapter only

3. **Navigation**
   - Page selection interface
   - Section selection within page
   - Previous/Next page navigation
   - Jump to specific page number

### FR-3: Glossary Integration
1. **Glossary Database**
   - 1,209 glossary terms stored in `alice_glossary` table
   - Terms linked to sections via `glossary_section_links` junction table
   - Each link includes: glossary_id, section_id, page_number, section_number

2. **Word Lookup**
   - Priority: Glossary definitions first
   - Fallback: External dictionary API (if not in glossary)
   - Case-insensitive matching
   - Multi-word term support (e.g., "waistcoat-pocket")

3. **Visual Indicators**
   - Glossary terms highlighted in section text
   - Hover tooltips show definitions
   - Glossary list displayed below each section

### FR-4: AI Assistance
1. **AI Integration**
   - Moonshot AI (Kimi) or Anthropic API
   - Context-aware prompts including:
     - Selected text
     - Current page/section
     - Chapter information
   - Multiple interaction types: explain, quiz, simplify, definition, chat

2. **Interaction Storage**
   - All AI interactions saved to database
   - User can view interaction history
   - Context preserved for follow-up questions

### FR-5: Help Request System
1. **Request Creation**
   - User submits help request with question
   - Context automatically included (page/section)
   - Status tracking: pending → assigned → resolved

2. **Consultant Assignment**
   - Requests assigned to available consultants
   - Consultant can view assigned requests
   - Consultant can respond to requests

### FR-6: Progress Tracking
1. **Reading Progress**
   - Track last page viewed
   - Track last section viewed
   - Calculate percentage complete
   - Store in `reading_progress` table

2. **Vocabulary Progress**
   - Track all word lookups
   - Store in `vocabulary_lookups` table
   - Include: word, definition, context, timestamp

---

## 🚫 Non-Functional Requirements

### NFR-1: Performance
- Page load time: < 2 seconds
- Word lookup response: < 500ms
- AI response time: < 5 seconds
- Database queries optimized with indexes

### NFR-2: Security
- Password hashing (bcrypt)
- SQL injection prevention (parameterized queries)
- Session management (secure tokens)
- Input validation and sanitization

### NFR-3: Usability
- Simple, intuitive interface
- Mobile-responsive design
- Clear error messages
- Accessible design (WCAG 2.1 AA minimum)

### NFR-4: Reliability
- Database backup strategy
- Error handling and logging
- Graceful degradation if services unavailable

### NFR-5: Scalability
- Support for multiple books (future)
- Support for multiple users concurrently
- Database schema supports expansion

---

## 📊 Database Schema Requirements

### Current Structure (Must Be Maintained)

1. **Pages Table**
   - `id`, `book_id`, `page_number`, `chapter_id`, `chapter_title`, `content`, `word_count`
   - Page numbers: 1-100 (for full book)
   - Average 180-200 words per page

2. **Sections Table**
   - `id`, `page_id`, `page_number`, `section_number`, `content`, `word_count`
   - Average 40 words per section (range: 35-45)
   - 4-5 sections per page

3. **Glossary Tables**
   - `alice_glossary`: 1,209 terms with definitions
   - `glossary_section_links`: Links terms to sections
   - Supports efficient lookup by section

---

## ✅ Success Criteria

### Phase 1 Complete When:
- [ ] User can register and login successfully
- [ ] User can enter book code and access content
- [ ] User can select page and section
- [ ] User can look up words (glossary priority)
- [ ] App remembers last page
- [ ] All features work with first 3 chapters
- [ ] Zero critical bugs

### Overall Success When:
- [ ] All user stories implemented and tested
- [ ] Performance meets NFR requirements
- [ ] Security requirements met
- [ ] User testing shows positive feedback
- [ ] System ready for production use

---

## 🔄 Constraints and Assumptions

### Constraints
- Backend: Go language
- Database: SQLite (file-based)
- Content: First 3 chapters initially (test ground)
- Sections: Average 40 words (must be maintained)

### Assumptions
- Users have access to physical book
- Users have internet connection
- AI service (Moonshot/Anthropic) is available
- Consultants available for Tier 3 help

---

## 📝 Notes

- **Database Structure:** Current page/section structure with 40-word average must be preserved
- **Glossary Integration:** 1,209 terms already linked to sections - this linking must be maintained
- **Progressive Enhancement:** Start with basic functionality, add features incrementally
- **Testing:** Each feature must be thoroughly tested before moving to next phase

---

**Next Step:** Use these requirements to create TECHNICAL_SPECIFICATIONS.md


