---
phase: "02-cli-interface-polish"
plan: "01"
subsystem: "task"
tags:
  - "error-handling"
  - "phase-02"
  - "store"
dependency_graph:
  requires: []
  provides:
    - "ErrTaskNotFound"
    - "corrupted todo file error"
  affects:
    - "internal/task/store.go"
tech_stack:
  added:
    - "fmt.Errorf with %w error wrapping"
  patterns:
    - "sentinel error wrapping with context"
    - "json.Unmarshal error wrapping with context"
key_files:
  created: []
  modified:
    - "internal/task/store.go"
decisions:
  - "Use fmt.Errorf with %w for error wrapping to enable errors.Is() checks"
  - "Wrap ErrTaskNotFound with task ID context in MarkDone and Delete"
  - "Wrap json.Unmarshal errors with 'corrupted todo file' context"
metrics:
  duration: "less than 1 minute"
  completed_date: "2026-04-14T07:17:00Z"
---

# Phase 02 Plan 01 Summary: Error Wrapping with %w

## One-liner

Error-wrapped MarkDone/Delete/loadLocked methods with sentinel error support via fmt.Errorf %w pattern.

## Completed Tasks

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Wire ErrTaskNotFound with %w in MarkDone and Delete | a97a662 | internal/task/store.go |
| 2 | Wrap JSON parse error with context in loadLocked | a97a662 | internal/task/store.go |

## Changes Made

**internal/task/store.go** (3 lines modified)

- Line 138 (MarkDone): `fmt.Errorf("task not found: %s: %w", id, ErrTaskNotFound)`
- Line 156 (Delete): `fmt.Errorf("task not found: %s: %w", id, ErrTaskNotFound)`
- Line 54 (loadLocked): `fmt.Errorf("corrupted todo file: %w", err)`

## Verification

- `go vet ./internal/task/...` passes with no warnings
- `go test ./internal/task/...` passes (13 tests, no regressions)
- `grep "ErrTaskNotFound" store.go` confirms %w wrapping at both call sites
- `grep "corrupted todo file" store.go` confirms error context wrapping

## Success Criteria

- [x] MarkDone returns `fmt.Errorf("task not found: %s: %w", id, ErrTaskNotFound)`
- [x] Delete returns `fmt.Errorf("task not found: %s: %w", id, ErrTaskNotFound)`
- [x] loadLocked returns `fmt.Errorf("corrupted todo file: %w", err)` on JSON parse error
- [x] go vet passes
- [x] go test passes

## Deviations from Plan

None - plan executed exactly as written.

## Threat Flags

None - error wrapping is a correctness improvement with no new attack surface.
