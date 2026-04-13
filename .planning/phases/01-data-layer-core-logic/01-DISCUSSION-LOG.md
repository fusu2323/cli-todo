# Phase 01: Data Layer & Core Logic - Discussion Log

> **Audit trail only.** Do not use as input to planning, research, or execution agents.
> Decisions are captured in CONTEXT.md — this log preserves the alternatives considered.

**Date:** 2026-04-13
**Phase:** 01-data-layer-core-logic
**Areas discussed:** UUID generation, Mutex strategy, Store interface

---

## UUID Generation

| Option | Description | Selected |
|--------|-------------|----------|
| Hex UUID | crypto/rand → 32 hex chars. Simple, no deps, Go-idiomatic. IDs like: a3f1b2c8d4e5f6789012345678901234 | ✓ |
| RFC 4122 UUID | crypto/rand + encoding/binary → version 4 compliant UUID. Slightly more code, but proper UUID format. IDs like: 550e8400-e29b-41d4-a716-446655440000 | |
| Timestamp + random | task_<unixnano>_<8hexchars>. Human-readable, sortable by creation time. IDs like: task_1672531200_a3f1b2c8 | |
| You decide | I trust the standard approach — pick whichever is most Go-idiomatic | |

**User's choice:** Hex UUID (Recommended)
**Notes:** Simplicity preferred — hex encoding with crypto/rand is sufficient for a todo app.

---

## Mutex Strategy

| Option | Description | Selected |
|--------|-------------|----------|
| sync.Mutex | Single lock for reads+writes. Simpler code, fine for a CLI todo app (low concurrency). One task at a time. | ✓ |
| sync.RWMutex | Separate read lock vs write lock. Better throughput if many concurrent reads. Slightly more complex. | |
| You decide | Use whichever is more standard for a single-user CLI tool | |

**User's choice:** sync.Mutex (Recommended)
**Notes:** Appropriate for a single-user CLI tool where concurrency is minimal.

---

## Store Interface

| Option | Description | Selected |
|--------|-------------|----------|
| Concrete type | JSONFileStore struct with Load/Save methods. Simpler, less indirection, no interface overhead. | ✓ |
| Interface (Repository pattern) | TaskStore interface with Load/Save, JSONFileStore implements it. Better for testing, future flexibility. Slight overhead. | |
| You decide | Go with whatever is more idiomatic for a learning project | |

**User's choice:** Concrete type (Recommended)
**Notes:** No need for interface abstraction — single implementation, concrete type is fine for learning project.

---

## Deferred Ideas

None — discussion stayed within phase scope.
