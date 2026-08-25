# Alice Suite agent-loop inventory

Status: proposed. Owner: user + agent. Updated: 2026-06-20.

## Research note

This inventory applies the loop-library idea to Alice Suite. The requested source,
`https://signals.forwardfuture.ai/loop-library/agents/`, could not be retrieved in
this environment: its page was absent from browser search and both normal and
proxy-bypassed HTTPS fetches timed out. Therefore this is **not** a verbatim
transcription of that page. It is a product-fit inventory grounded in Alice's
three-tier model and current Reader/Consultant flows.

## Operating rule

Every loop must specify a trigger, bounded action, observable feedback, a stop
condition, and a human escalation path. Do not ship autonomous actions that
change reader data, send messages, or resolve support requests without an
explicit user or consultant decision.

## Reusable skill entry

Use **Agent-loop design and governance** when proposing or improving an Alice
workflow with AI, automation, feedback, or escalation. Define the trigger,
bounded action, feedback, stop condition, ownership, privacy boundary, and
human decision point before implementation. The skill is listed in
`docs/HERMES_HANDOVER_BRIEFING.md` as a cross-cutting capability for all three
apps.

## Product loops

| Priority | Loop | Trigger → bounded action → feedback | Stop / escalation | Current fit |
| --- | --- | --- | --- | --- |
| P0 | Reading-context help | Reader selects section, page, or passage → asks AI for an explanation tied to that scope → reader can ask again, change scope, or rate usefulness | Answer delivered; escalate to human help on request | Existing Tier 2; stabilize the context seam first |
| P0 | Dictionary disambiguation | Reader taps a word → show definition, pronunciation, derivation, and example → reader marks it clear or asks for a contextual explanation | Definition viewed or AI/human help requested | Existing Tier 1; add lightweight helpful/not-helpful signal only |
| P0 | AI answer grounding and recovery | AI request → validate scope, source text, and response shape → show answer or a recoverable error | One bounded retry; offer human help after repeated failure | Required before wider AI use |
| P0 | Human-help handoff | Reader submits a question with reading context → create a help request → consultant claims, replies, and reader confirms/continues | Resolved by consultant or reopened by reader | Existing Tier 3; make ownership and status explicit |
| P0 | Help-request triage | New/changed request → classify urgency and route to an eligible consultant → consultant accepts or redirects | Assigned, resolved, or placed in a visible unassigned queue | Consultant workflow needs a deterministic queue, not autonomous replies |
| P1 | Consultant follow-up | Consultant response → reader sees it and can mark resolved, ask a follow-up, or request reassignment → consultant receives the result | Reader confirms resolution, request is reopened, or timeout prompts manual review | Needs request-status and notification design |
| P1 | Reading-resumption | Reading event → save last verified book/section/page → offer Continue Reading on the next visit → reader accepts or chooses another location | Reader resumes, dismisses, or changes book | Roadmap already identifies this as a small, high-value slice |
| P1 | Learning-check / quiz | Reader asks for a quiz after a bounded reading scope → answer → immediate explanation and optional retry | Score summary shown; return to book | Roadmap scope; retain book context and avoid grading ambiguity as fact |
| P1 | AI quality feedback | After an AI response → collect optional helpfulness/reason signal plus scope/model metadata → aggregate failures and sample responses for review | Feedback recorded; no automatic prompt change | Add only after reliable request IDs and privacy rules exist |
| P1 | Escalation recommendation | Low confidence, repeated failed AI attempts, reader explicitly dissatisfied, or safety-sensitive request → offer the human consultant path with captured context | Reader declines or request created | Recommendation only; never silently create a request |
| P2 | Reader wellbeing / friction detection | Repeated dictionary taps, rewinds, failed answers, or abandoned quiz → offer one unobtrusive aid (simplify, explain context, consultant) | Reader acts/dismisses; suppress repeated prompts | Needs careful consent and frequency caps; do not infer a reader's ability |
| P2 | Content-correction | Reader/consultant flags a wrong definition, misleading answer, or broken book text → create an internal review item with evidence → curator fixes source/prompt/test | Accepted, rejected with rationale, or deferred | Requires a curator role and audit trail |
| P2 | Knowledge-base improvement | Resolved human requests and repeated AI failures → de-identify and cluster → propose FAQ, glossary, or prompt-test updates → human approves | Approved update shipped or proposal rejected | Human approval is mandatory; never train on raw reader content by default |

## Operational loops

| Priority | Loop | Trigger → bounded action → feedback | Stop / escalation | Current fit |
| --- | --- | --- | --- | --- |
| P0 | Request safety and privacy | Every AI/help request → minimise context, enforce access control, and retain audit metadata → log policy outcome | Block/redact on violation; user-visible reason; security review for incidents | Required guardrail for all Tier 2/3 loops |
| P0 | Delivery verification | Code or prompt/config change → build, tests, template checks, and targeted manual flow → failures return to the change owner | Checks pass or change is blocked | Existing autoresearch pattern is the starting point |
| P1 | Production health and error recovery | Health/error/latency threshold → alert with request-safe diagnostics → operator investigates and confirms recovery | Metrics return to threshold or incident remains open | Define SLOs before automation |
| P1 | Prompt/model release evaluation | Candidate prompt/model → run a versioned evaluation set covering explanation, grounding, safety, and handoff → compare against current baseline | Human approves promotion or rollback | Do not let production feedback auto-promote a model |
| P1 | Cost and rate-limit control | AI usage crosses per-reader or service budget → throttle or present a clear fallback → record outcome | Budget window resets or operator adjusts policy | Needed before unattended growth |
| P2 | Backlog / product-discovery | Repeated feedback, failed flows, and consultant workload signals → create a reviewed product hypothesis → prioritize or close | Decision recorded in wiki/roadmap | Weekly human review; not an autonomous roadmap writer |

## Recommended implementation order

1. Finish the Reader AI context seam, then implement the **grounding and recovery** contract.
2. Complete the **human-help handoff**, triage, and follow-up states so Tier 3 is a closed loop.
3. Add reading resumption and quiz feedback as bounded reader-value loops.
4. Add safe telemetry, evaluation, cost, and health loops only after consent, retention, and access rules are written.
5. Defer behavioral nudges and knowledge-base learning until there is enough reviewed data and a designated human owner.

## Source docs

- `README.md` — physical-book companion and three-tier service model.
- `docs/UPGRADE_ROADMAP.md` — Reader, quiz, AI, and Consultant workflow priorities.
- `docs/AI_ASSISTANT_FLOW_INVENTORY.md` — current reading-context behavior and known constraints.
- `docs/READER_AI_CONTEXT_MODULE_DESIGN.md` — the required Reader-AI seam before extending Tier 2 behavior.
