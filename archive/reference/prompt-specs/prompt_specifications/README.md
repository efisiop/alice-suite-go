# ALICE READER APP - PROMPT SPECIFICATIONS

📁 **Complete Technical Blueprint for the Alice Reader App Enhancement**

---

## 📋 MASTER DOCUMENT

| Document | Purpose | Priority |
|----------|---------|----------|
| [**MASTER_BLUEPRINT.md**](MASTER_BLUEPRINT.md) | Complete CPB specifications | 🟢 CRITICAL |

---

## 🗂️ SPECIFICATION CATEGORIES

### 📊 **IMPLEMENTATION DELIVERABLES**
| Document | Focus | Metrics Covered |
|----------|--------|-----------------|
| [01_deliverables.md](01_deliverables.md) | Phase-by-phase deliverables | Quantifiable success criteria for all 4 phases |

### 🚨 **ERROR HANDLING**
| Document | Focus | Coverage |
|----------|--------|----------|
| [02_failure_states.md](02_failure_states.md) | Failure recovery procedures | All failure categories with graceful fallbacks |

### 📡 **API SPECIFICATIONS**
| Document | Focus | Standards |
|----------|--------|-----------|
| [03_schemas.md](03_schemas.md) | Response structure | Complete API output schemas |

### ✅ **SUCCESS EXAMPLES**
| Document | Focus | Examples |
|----------|--------|----------|
| [04_examples.md](04_examples.md) | Implementation examples | Real user journey demonstrations |

### 🔗 **AI ENHANCEMENT**
| Document | Focus | Integration |
|----------|--------|-------------|
| [05_rag_integration.md](05_rag_integration.md) | External knowledge | RAG system requirements |

---

## 🎯 READER'S GUIDE

### **IMMEDIATE ACTION ITEMS**
1. Start with [**MASTER_BLUEPRINT.md**](MASTER_BLUEPRINT.md) for complete overview
2. Review **01_deliverables.md** for specific technical requirements
3. Study **02_failure_states.md** for error handling implementation
4. Use **03_schemas.md** for API development standards

### **REFERENCE MATERIAL**
- **04_examples.md** - See how features should actually work
- **05_rag_integration.md** - Understand AI enhancement approach

---

## 🚨 CORE AUTHORITY REMINDERS

### **NON-NEGOTIABLE CONSTRAINTS**
- ✅ **Database:** Must use existing SQLite at `/Users/efisiopittau/Project_1/alice-suite-go/data/alice-suite.db`
- ✅ **Backend:** Go 1.24+ with standard net/http library only
- ✅ **Glossary:** 1,209 existing terms take absolute priority
- ✅ **Typography:** Georgia serif 1.1rem, 1.8 line-height (IMMUTABLE)
- ✅ **Colors:** Purple (#6a51ae) primary, Pink (#ff6b8b) secondary
- ✅ **Response Times:** Dictionary <500ms, AI <1.5s, Page load <200ms

### **SUCCESS METRICS (MUST ACHIEVE 5/5)**
- **F-CG (Fidelity to Core Goal):** Perfect alignment with user's passionate vision
- **A-C (Actionable Conciseness):** All specifications immediately implementable
- **F-D (Format Determinism):** Consistent, predictable outputs across all systems

---

## 🏗️ PROJECT STRUCTURE VISUALIZATION

```
alice-suite-go/
├── prompt_specifications/
│   ├── MASTER_BLUEPRINT.md          📋 Complete specifications overview
│   ├── 01_deliverables.md           📊 Phase-by-phase requirements
│   ├── 02_failure_states.md         🚨 Error handling procedures
│   ├── 03_schemas.md               📡 API response standards
│   ├── 04_examples.md              ✅ Implementation examples
│   ├── 05_rag_integration.md       🔗 AI knowledge enhancement
│   └── README.md                   📖 This file
├── data/
│   ├── alice-suite.db              💾 Your existing SQLite database
│   └── (1,209 glossary terms)      📝 The heart of the system
├── server.go                       ⚡ Your Go backend (extend only)
└── static/                         🎨 Frontend code (enhance only)
```

---

## 🎯 COMMANDS FOR DEVELOPMENT

### **WORKING WITH SPECIFICATIONS**
```bash
# View complete specifications
cat MASTER_BLUEPRINT.md

# Search for specific requirements
grep -r "response_time" ./            # Response time requirements
grep -r "failure_state" ./            # Error handling

# Validate against examples
diff 04_examples.md your_implementation.md
```

---

## 🚨 PROJECT STATUS: IMPLEMENTATION READY

✅ **All specifications complete**
✅ **All constraints documented**
✅ **All examples provided**
✅ **All schemas defined**
✅ **All failure states covered**

**Next Action:** Technical implementation using the MASTER_BLUEPRINT specifications.

---

*Created: 2025-11-21*
*Authority: Project Manager Partner Specifications*
*Core Vision: "Building the exact app the user is passionate about"*