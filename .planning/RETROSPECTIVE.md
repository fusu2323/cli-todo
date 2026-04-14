# Project Retrospective

*A living document updated after each milestone. Lessons feed forward into future planning.*

## Milestone: v1.0 — MVP

**Shipped:** 2026-04-14
**Phases:** 2 | **Plans:** 5 | **Sessions:** ~3

### What Was Built
- Task domain model with cryptographically secure UUID (crypto/rand + hex encoding)
- Thread-safe JSON persistence layer with sync.Mutex + atomic writes (CreateTemp + Rename)
- Concurrent access safety verified via stress testing (200 concurrent ops, no data corruption)
- Error wrapping with %w pattern for sentinel error support (errors.Is checks in CLI)
- CLI entry point with flag.FlagSet subcommand routing for add/list/done/delete/help

### What Worked
- Go stdlib-only constraint kept dependencies at zero — appropriate for a learning project
- Phase 1 (data layer) gave solid foundation before Phase 2 (CLI) — clear dependency order
- Atomic writes via CreateTemp+Rename provided cross-platform safety without external deps
- sync.Mutex single-lock approach was simple and sufficient for single-user CLI

### What Was Inefficient
- GSD workflow overhead for a 2-phase project — lots of ceremony for small scope
- No milestone audit run before completion (skipped for speed given small scope)
- RETROSPECTIVE.md created post-hoc instead of during execution

### Patterns Established
- Task struct with JSON tags, omitempty on optional Category field
- loadLocked/saveLocked private methods for nested lock calls (avoids lock inversion)
- Silent success for mutations (add/done/delete), errors to stderr with exit 1
- flag.NewFlagSet per subcommand for proper os.Args slicing

### Key Lessons
1. Go stdlib is sufficient for CLI apps — flag, os, encoding/json, sync cover all needs
2. Phase ordering matters: data layer before CLI gives a testable API to build against
3. Atomic writes (CreateTemp+Rename) are essential for JSON file corruption prevention

---

## Cross-Milestone Trends

### Process Evolution

| Milestone | Sessions | Phases | Key Change |
|-----------|----------|--------|------------|
| v1.0 | ~3 | 2 | Initial delivery — no prior milestone |

### Cumulative Quality

| Milestone | Tests | Zero-Dep Additions |
|-----------|-------|-------------------|
| v1.0 | 13 | 0 (stdlib only) |

### Top Lessons (Verified Across Milestones)

1. Go stdlib covers CLI essentials — no external dependencies needed for v1
2. Data layer first, CLI second provides testable API foundation
3. Atomic writes prevent JSON corruption in concurrent scenarios
