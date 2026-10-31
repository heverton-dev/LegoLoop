# ==============================================================================
# LEGOLOOP AUTONOMOUS IMPLEMENTATION PROMPT (TELEGRAPHIC / IMPERATIVE)
# ==============================================================================

TARGET DIRECTORY: C:\Users\trcnologia\Desktop\LegoLoop
WORKSPACE CONTEXT: Composable Looping & Harness Fabric (AIDD + Symbiont + 9Router)

EXECUTE STRICT AUTONOMOUS IMPLEMENTATION PIPELINE ACCORDING TO SPEC.md AND PLAN.md.
ZERO TOLERANCE FOR PLACEHOLDERS, STUBS, TODO COMMENTS, OR MOCKED CODE.

--------------------------------------------------------------------------------
1. CONTEXT INGESTION & BOOTSTRAP (MANDATORY READ FIRST)
--------------------------------------------------------------------------------
READ AND INGEST IN PARALLEL:
- C:\Users\trcnologia\Desktop\LegoLoop\SPEC.md (The 4 Engineering Pillars: Harness, Looping, Graph, Router)
- C:\Users\trcnologia\Desktop\LegoLoop\PLAN.md (Tickets TICK-101 to TICK-503)
- C:\Users\trcnologia\Desktop\LegoLoop\AGENTS.md (Inviolable Rules: Zero Stubs, Caveman Economy, Strict TDD)
- C:\Users\trcnologia\Desktop\LegoLoop\docs\specs\*.md (Architectural Sub-specs)
- C:\Users\trcnologia\Desktop\LegoLoop\tests\suites\* (Sanitizer & Gateway Reference Tests)

--------------------------------------------------------------------------------
2. MASTER EXECUTION LOOP: ADDY OSMANI BUILD-AUTO CYCLE
--------------------------------------------------------------------------------
FOR EACH EPIC (EP-01 to EP-05) AND EACH TICKET (TICK-101 to TICK-503):

[STEP 2.1 - /build auto]
- Implement real, production-ready Go code (and accompanying CLI/UI files).
- Follow Clean Architecture / Vertical Slice Architecture (VSA).
- Enforce strict typing, zero global mutable state, idiomatic Go/TypeScript.

[STEP 2.2 - /test]
- Enforce strict Red-Green TDD.
- Execute: `go test -v -race ./...` (Target: >=85% unit & integration test coverage).
- Zero skipped tests. All assertions must validate real behavior.

[STEP 2.3 - /review]
- Execute static analysis and strict code review.
- Eliminate all stubs, unused imports, unchecked errors, and memory leaks.
- Ensure all handoffs match strict JSON schemas defined in `SPEC.md`.

[STEP 2.4 - /webperf]
- Benchmark startup time: Must be <100ms.
- Benchmark RAM footprint: Core daemon must remain <45MB in idle.
- Verify zero GPU lock, non-blocking I/O in scanner worker pools.

[STEP 2.5 - /code-simplify]
- Refactor for maximum clarity, readability, and DRY compliance.
- Keep functions small (<35 lines), explicit signatures, minimal nesting.

--------------------------------------------------------------------------------
3. GATEWAY AUDIT: CLOUDFLARE SECURITY PASS & CIRCUIT BREAKER
--------------------------------------------------------------------------------
TRIGGER: Execute deep security audit via Cloudflare Security standards on all code, APIs, and CLI handlers.
INSPECT:
- Command injection vulnerabilities (safe subprocess execution).
- Path traversal & symlink escape vulnerabilities in Windows/Linux.
- API authentication, CORS, rate-limiting, and header sanitization.
- JSON schema validation safety (no prototype pollution, no unvalidated maps).

CIRCUIT BREAKER RULE:
IF Cloudflare Security Pass discovers ANY vulnerability, warning, bug, or lint defect:
  1. DO NOT PROCEED TO NEXT EPIC.
  2. IMMEDIATELY invoke Addy Osmani `/plan` skill on the discovered defect list.
  3. Generate remediation nano-tickets.
  4. ROUTE BACK to [STEP 2.1 /build auto].
  5. RE-EXECUTE [build -> test -> review -> webperf -> code-simplify -> security pass].
  6. TARGET: ABSOLUTE ZERO DEFECTS (0 errors, 0 warnings, 0 security flags).

--------------------------------------------------------------------------------
4. UI/UX POST-HARDENING: /IMPECCABLE SUITE
--------------------------------------------------------------------------------
UPON 100% GREEN SECURITY PASS & COMPLETE TICKET IMPLEMENTATION:
TRIGGER: Execute `/impeccable` design and interface polishing cycle:
- `/impeccable init`: Inspect Looping Studio Canvas (`index.html` / CSS / components).
- `/impeccable document`: Extract design tokens, color palette, typography, grid spacing.
- `/impeccable audit`: Eliminate visual clutter, reduce cognitive load, align node badges.
- `/impeccable polish`: Refine drag-and-drop connectors, hover states, status pills, dark mode.
- `/impeccable verify`: Responsive validation across desktop viewports with zero layout shift.

--------------------------------------------------------------------------------
5. FINALIZATION & FACTUAL HANDOFF
--------------------------------------------------------------------------------
- Verify the entire test suite passes: `go test -v ./...`
- Execute: `faz-session --format all --update`
- Generate final delivery report with exact file links and verified benchmark metrics.
