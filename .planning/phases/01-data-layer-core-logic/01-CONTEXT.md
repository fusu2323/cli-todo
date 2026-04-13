# Phase 01: Data Layer & Core Logic - Context

**Gathered:** 2026-04-13
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver a working data layer: Task struct, JSON file persistence with mutex protection, all CRUD operations (add/list/done/delete), category support, and atomic writes. No CLI interface yet — that's Phase 2.

</domain>

<decisions>
## Implementation Decisions

### Data Model
- **D-01:** Hex UUID via `crypto/rand` + hex encoding — 32-character hex string, no external dependencies. IDs like: `a3f1b2c8d4e5f6789012345678901234`
- **D-02:** `sync.Mutex` — simple single-lock approach for concurrent access protection (DATA-01). Appropriate for a single-user CLI tool.
- **D-03:** Concrete `JSONFileStore` struct with `Load()` and `Save()` methods — no interface abstraction for this phase.

### Atomic Writes (DATA-03)
- Write to temp file via `os.MkdirTemp`, then rename via `os.Rename`
- `os.Rename` is atomic on POSIX, cross-platform considerations handled by Go stdlib

### Category Support (CAT-01)
- Tasks have an optional `category` string field (empty string = uncategorized)

### JSON Structure
- Flat array of tasks: `[]Task`
- Human-readable via `json.MarshalIndent` with `"  "` indent

### File Location
- `~/.todo.json` via `os.UserHomeDir()` (cross-platform)

### Error Handling
- Phase 2 handles custom error types and wrapped errors
- Phase 1 uses basic Go error returns; DATA-02 (file-not-found → empty list) handled inline

### Claude's Discretion
- Task struct field names (e.g., `CreatedAt` vs `created_at` in JSON)
- Whether to include `UpdatedAt` timestamp on tasks
- Internal function organization within store.go and task.go
</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

- `.planning/research/ARCHITECTURE.md` — Project structure, anti-patterns (no global state, no integer IDs, repository pattern guidance)
- `.planning/research/STACK.md` — Go stdlib patterns, flag.FlagSet for subcommands, os.MkdirTemp for atomic writes
- `.planning/REQUIREMENTS.md` — All v1 requirements mapped to phases (TASK-01 through TASK-05, CAT-01, CAT-02, DATA-01, DATA-02, DATA-03)
- `.planning/ROADMAP.md` — Phase 1 goal and success criteria

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- No existing Go code yet — this is Phase 1 (greenfield)

### Established Patterns
- `flag.FlagSet` for subcommand routing (Phase 2 will use this)
- Repository pattern guidance from ARCHITECTURE.md (but concrete type chosen, not interface)

### Integration Points
- `internal/task/task.go` — Task struct lives here
- `internal/task/store.go` — JSONFileStore lives here
- Phase 2: `cmd/main.go` will import `internal/task` packages

</code_context>

<specifics>
## Specific Ideas

- Keep it simple — this is a learning project, not production infrastructure
- No need for `github.com/google/uuid` — pure Go `crypto/rand` is sufficient
- No need for interface abstraction — single implementation, concrete type is fine

</specifics>

<deferred>
## Deferred Ideas

### Reviewed Todos (not folded)
None — no pending todos matched Phase 1 scope.

</deferred>

---

*Phase: 01-data-layer-core-logic*
*Context gathered: 2026-04-13*
