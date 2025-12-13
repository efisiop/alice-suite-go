# ALICE READER APP - CONCISE PROMPT BLUEPRINT (COMPLETE)

**REF:** ALICE-CODE-2025-11-21-CPB-FINAL
**VERSION:** 1.0
**AUTHORITY:** Project Manager Partner Specifications
**SESSION DATE:** 2025-11-21
**CORE AUTHORITY:** Original User Request - "Build the right app I'm passionate about"

---

## 👑 PRIMARY AUTHORITY STATEMENT

**ILLUMINATING GOAL:** Build the exact Alice Reader app the user is passionate about - nothing more, nothing less.
**IMMUTABLE SCOPE:** Work only with the existing Alice database (1,209 glossary terms, SQLite, Go backend).
**CONTENT VISION:** Support physical book reading with gentle, contextual assistance.

---

## 🎯 CRITICAL SUCCESS METRICS (MUST ACHIEVE 5/5)

| Metric | Definition | Success Criteria |
|--------|------------|------------------|
| **F-CG** (Fidelity to Core Goal) | Alignment with user's passionate vision | ✓ Zero feature bloat, ✓ Focus on physical book support, ✓ Reading flow preservation |
| **A-C** (Actionable Conciseness) | Immediate implementation capability | ✓ Working Go/JS code provided, ✓ Clear business logic, ✓ No ambiguity |
| **F-D** (Format Determinism) | Consistent, predictable outputs | ✓Standardized API responses, ✓ Error handling schemas, ✓ Measurable deliverables |

---

## 🔒 NON-NEGOTIABLE TECHNICAL CONSTRAINTS

### DATABASE INTEGRATION (MUST)
- **File Path:** `/Users/efisiopittau/Project_1/alice-suite-go/data/alice-suite.db`
- **Core Tables:** books, pages, sections, alice_glossary, glossary_section_links
- **Treasured Data:** 1,209 glossary terms already linked to sections
- **Intelligent Linking:** Preserve existing priority lookup system
- **Schema Preservation:** ZERO modifications allowed

### BACKEND TECHNOLOGY (MUST)
- **Language:** Go 1.24+ (already implemented)
- **Framework:** Standard Go net/http library only
- **Database Driver:** mattn/go-sqlite3
- **API Extension:** Build upon existing endpoints, never replace
- **Authentication:** Extend current bcrypt system, support JWT-ready architecture

### EXISTING API ENDPOINTS (MUST MAINTAIN)
- `/api/auth/register` & `/api/auth/login` - Authentication
- `/api/books`, `/api/pages`, `/api/sections` - Content access
- `/api/dictionary/lookup` - Word lookup (glossary priority)
- `/api/ai/ask` - AI assistance integration
- `/api/help/request` - Tier 3 consultant help

---

## 📖 USER EXPERIENCE FOUNDATIONS (PRESERVE AT ALL COSTS)

### READING EXPERIENCE CORE
- **3-Step Progressive Flow:** Page Input → Section Selection → Content Interaction
- **Visual Ratio:** 75% content space, 25% sidebar
- **Typography:** Georgia serif 1.1rem, 1.8 line-height (IMMUTABLE)
- **Interaction Pattern:** Hover triggers for word definitions
- **Context Preservation:** Current page/section always visible

### SMART CONTEXT INTELLIGENCE
- **Word Threshold:** 1-5 words → dictionary; 6+ words → AI assistant
- **Glossary Priority:** Alice database terms ALWAYS before external sources
- **Multi-word Recognition:** Phrasal verb detection and explanation
- **Position Awareness:** Context-aware definitions based on reading location

### VISUAL DESIGN LANGUAGE (IMMUTABLE)
- **Colors:** Purple (#6a51ae) primary, Pink (#ff6b8b) secondary
- **Special Vocabulary:** Orange highlights for Alice-specific terms
- **Animations:** 300ms fade-in transitions throughout
- **Mobile-First:** Responsive design, no loading screens >1.5s

---

## 🏗️ 4-PHASE IMPLEMENTATION BLUEPRINT

### PHASE 1: Database Accessibility & Core Reading [COMPLETE]
**Deliverable:** 100% SQLite accessibility + 3-step reading interface
**Quantifiable:** Sub-500ms response times, 100% glossary integration
**Example Success:** User navigates to page 67 → selects Section 3 → reads content in 200ms

### PHASE 2: User Administration [COMPLETE]
**Deliverable:** Authentication + access codes + progress tracking
**Quantifiable:** >98% login success, real-time position saving
**Example Success:** User registers → enters book code → continues from last position

### PHASE 3: Enhanced Features & AI Integration [COMPLETE]
**Deliverable:** Complete AI assistant + vocabulary history + analytics
**Quantifiable:** <1.5s AI responses, >90% explanation satisfaction
**Example Success:** "dream-like quality" selection → AI explanation with context

### PHASE 4: Consultant Dashboard & Tier 3 System [COMPLETE]
**Deliverable:** Consultant workflow + two-way communication + analytics
**Quantifiable:** <30s assignment time, >85% resolution satisfaction
**Example Success:** Help request submitted → consultant assigned → resolution session

---

## 🚨 FAILURE STATE DEFINITIONS (COMPLETE)

### CRITICAL FAILURES (App Cannot Continue)
- **F1-001 Database Connection Failure:** Graceful degradation, offline mode available
- **F1-002 Authentication Failure:** Anonymous reading mode with local storage
- **F1-003 Complete API Failure:** Client-side cached content with basic navigation

### FEATURE FAILURES (Reduced Functionality)
- **AI Service Unavailable:** Fallback to enhanced dictionary mode
- **External Dictionary API Failure:** Alice glossary terms prioritized
- **Progress Sync Failure:** Dual storage (server + localStorage)

### INPUT VALIDATION FAILURES (User Retry)
- **F3-001 Invalid Page Number:** Clear range suggestions with alternatives
- **F3-002 Invalid Access Code:** Specific error messages with help options
- **Format Violations:** Real-time validation with examples

### TIMEOUT & PERFORMANCE FAILURES
- **Dictionary Lookup >2s:** Loading indicators + retry options
- **Page Load >3s:** Skeleton loaders with progress indication
- **AI Response >1.5s:** Contextual loading messages

---

## 📊 STRUCTURED OUTPUT SCHEMAS (COMPLETE)

### Base Response Wrapper
```json
{
  "success": true,
  "timestamp": "2025-11-21T14:30:00Z",
  "request_id": "req_abc123def456",
  "version": "v1.0",
  "data": {},
  "metadata": {
    "processing_time_ms": 45,
    "cache_hit": false,
    "rate_limit_remaining": 1000
  },
  "error": null
}
```

### Key APIs Schemas
- **Book Pages:** Navigation, sections, glossary terms highlighted
- **Dictionary Lookup:** Alice-first definitions with context
- **AI Assistant:** Context-aware explanations with 5 interaction types
- **User Management:** Authentication, progress, preferences
- **Help System:** Context capture, consultant assignment, communication

---

## ✅ SUCCESSFUL EXAMPLE OUTPUTS (COMPLETE)

### Phase 1 Success Story
```
User Journey: Types "67" → Sees 4 sections → Clicks Section 3
Performance: 150ms page display, 200ms section content
Visual: 75% content space with Georgia serif typography
Interaction: Hover word "Antipathies" → Definition appeared in 285ms
```

### Phase 2 Success Story
```
Registration: User enters email/password → Account created in 380ms
Access Control: Book code "ALICE001" → Verified in 234ms
Progress: Reading position saved → Continued from exact location
```

### Phase 3 Success Story
```
AI Assistance: "curious dream-like quality" selected → 923ms response
Explanation: Victorian context + Carroll's intentions + modern explanation
Follow-up: Generated quiz question + related discussion topics
```

### Phase 4 Success Story
```
Help Request: Submitted with full context → Assigned in 18 seconds
Consultant Response: Jane Smith (M.Ed) provided comprehensive explanation
Resolution: Student achieved emotional understanding of Alice's confusion
```

---

## 🔗 EXTERNAL CONTEXT (RAG) INTEGRATION (COMPLETE)

### Knowledge Bases
- **KB1 Literary Analysis:** Victorian literature frameworks, Carroll biography
- **KB2 Educational Strategies:** Age-appropriate teaching methods, comprehension strategies
- **KB3 Historical Context:** Victorian customs, educational systems, logical debates
- **KB4 Reading Support:** Difficulty analysis, dyslexia-friendly approaches

### Context Scoring System
```python
Relevance = (Query Similarity × 0.25) + (Context Similarity × 0.20) +
           (Reading Level Match × 0.15) + (Length Appropriateness × 0.10) +
           (Educational Value × 0.20) + (Curriculum Importance × 0.10)
```

### Integration Triggers
- AI detects analytical questions
- Victorian terminology confusion
- Reading comprehension difficulties
- Symbolism/thematic analysis requests

---

## 💡 CONFIRMED NEGATIVE BEHAVIORS (AVOID)

### ABSOLUTELY PROHIBITED
- ✗ Database schema modifications
- ✗ New programming languages or frameworks
- ✗ Features that interrupt reading
- ✗ Typography changes without justification
- ✗ Loading screens >1.5 seconds
- ✗ Alice theme color palette changes
- ✗ Technical jargon in user explanations
- ✗ Unnecessary features beyond core vision

### STRONGLY DISCOURAGED
- ✗ Complex animations that slow reading
- ✗ Multi-page application navigation
- ✗ Heavy JavaScript frameworks
- ✗ Dependencies requiring complex setup
- ✗ Features requiring extensive tutorials

---

## 📋 PROJECT MANAGER PARTNER INSTRUCTIONS

### Your Role
**PROJECT MANAGER PARTNER** for Alice Reader App enhancement

### Your Authority
- Interpret and translate user passion into precise technical specifications
- Ensure zero deviation from existing codebase constraints
- Maintain strict adherence to proven UX patterns
- Preserve the "Gentle Exploration" philosophy

### Your Constraints
- Work ONLY within existing technology stack (Go, SQLite, vanilla JS)
- Build upon current API endpoints - never replace them
- Maintain 1,209 glossary terms as primary knowledge source
- Preserve typography, colors, and layout specifications
- Keep response times under specified thresholds

### Your Success Criteria
- Deliver working code that implements specifications exactly
- Provide structured deliverables matching output schemas
- Ensure all failures have graceful recovery procedures
- Validate successful example outputs match real implementation
- Maintain F-CG = 5/5, A-C = 5/5, F-D = 5/5 across all features

### Your Communication Protocol
- Use technical specifications language, not user-facing explanations
- Reference specific file paths and line numbers
- Provide measurable outcomes for each deliverable
- Include fallback procedures for all failure states
- Timestamp all work with session date: 2025-11-21

---

## 📦 IMMEDIATE NEXT ACTIONS

1. **IMPLEMENT COMPLETE SPECIFICATIONS**: Use all detailed technical documentation
2. **VALIDATE AGAINST SUCCESS EXAMPLES**: Ensure real outputs match provided examples
3. **TEST FAILURE RECOVERY**: Verify all failure states have proper fallbacks
4. **PERFORMANCE VALIDATION**: Confirm all response time requirements are met
5. **USER EXPERIENCE CONFIRMATION**: Guarantee reading flow is never interrupted

---

## 🎯 FINAL AUTHORITY STATEMENT

**This Concise Prompt Blueprint serves as the complete technical specification for building the exact Alice Reader app the user is passionate about. All requirements, constraints, deliverables, and success criteria are defined within. Implementation must adhere strictly to these specifications while maintaining zero deviation from the user's vision of a gentle, contextual reading companion for physical books.**

**PROJECT STATUS:** Ready for immediate implementation
**DELIVERABILITY:** All specifications complete with examples
**IMPLEMENTATION AUTHORITY:** Proceed with full technical specifications

---

**END OF ALICE READER APP - CONCISE PROMPT BLUEPRINT**
**Session Complete: 2025-11-21**
**Next Action: Technical Implementation**