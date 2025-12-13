# Forward Pipeline Execution Guide

## Overview
This guide walks you through executing the forward pipeline to generate requirements and specifications for the Alice Suite Go rebuild.

## Prerequisites
✅ Cursor restarted (agents are now available)  
✅ Working directory: `/Users/efisiopittau/Project_1/alice-suite-go`  
✅ Recovered brief available: `ALICE_SUITE_RECOVERED_BRIEF.md`

---

## Step 1: Idea Scaffolder

**Goal:** Refine the recovered project brief into a clear, actionable project concept.

**Command to run in Cursor chat (Cmd+L):**
```
@idea-scaffolder Take this recovered and refined project brief for Alice Suite and refine it into a clear, actionable project concept. Focus on making each step seamless and working perfectly before expanding goals, graphics, or new features. Prioritize functionality and user experience over visual enhancements.

[Paste the content from ALICE_SUITE_RECOVERED_BRIEF.md here, or reference it]
```

**What to expect:**
- The agent will ask clarifying questions OR generate a refined concept directly
- Output will be a refined project brief ready for requirements analysis

**Save the output as:** `REFINED_CONCEPT.md`

---

## Step 2: Requirements Analyst

**Goal:** Convert the refined concept into detailed, testable requirements.

**Command to run in Cursor chat:**
```
@requirements-analyst Create detailed, implementable requirements from this refined Alice Suite concept. Include:
- User Stories with acceptance criteria for the reader app
- Functional requirements for sign-up, authorization, login, welcome page, main reader page, and word definitions
- Non-functional requirements (performance, security, usability)
- Constraints and assumptions
- Success criteria

Focus on the reader app specifications only (consultant dashboard will come later).
```

**What to expect:**
- User stories in format: "As a [user], I want [goal] so that [benefit]"
- Detailed functional requirements
- Non-functional requirements
- Edge cases and constraints

**Save the output as:** `REQUIREMENTS.md`

---

## Step 3: Specification Writer

**Goal:** Create technical specifications for implementing the requirements.

**Command to run in Cursor chat:**
```
@specification-writer Create comprehensive technical specifications for implementing Alice Suite based on these requirements. Include:
- Technical architecture (Go backend, SQLite database, frontend approach)
- Data models and database schema
- API endpoint specifications
- Authentication and authorization flow
- Implementation phases (prioritizing seamless functionality)
- Testing strategy
- Technology stack details

Remember: Backend is Go, database is SQLite, focus on seamless functionality before visual enhancements.
```

**What to expect:**
- System architecture design
- Database schema
- API specifications
- Implementation plan
- Technology stack recommendations

**Save the output as:** `TECHNICAL_SPECIFICATIONS.md`

---

## Step 4: Review & Build

After completing all three steps:

1. **Review all outputs:**
   - `REFINED_CONCEPT.md`
   - `REQUIREMENTS.md`
   - `TECHNICAL_SPECIFICATIONS.md`

2. **Make any necessary adjustments** based on your vision

3. **Begin implementation** following the technical specifications

---

## Tips

- **Be specific:** When calling agents, include context about the Go/SQLite migration and focus on seamless functionality
- **Save outputs:** Keep all intermediate documents for reference
- **Iterate if needed:** You can run agents multiple times to refine outputs
- **Reference the brief:** Always remind agents about the development philosophy (seamless functionality first)

---

## Quick Reference

**Agent Names:**
- `@idea-scaffolder` - Refines concepts
- `@requirements-analyst` - Creates requirements
- `@specification-writer` - Creates technical specs

**Key Documents:**
- `ALICE_SUITE_RECOVERED_BRIEF.md` - Starting point (recovered brief)
- `REFINED_CONCEPT.md` - After Step 1
- `REQUIREMENTS.md` - After Step 2
- `TECHNICAL_SPECIFICATIONS.md` - After Step 3

---

**Ready to start? Begin with Step 1 above!**



