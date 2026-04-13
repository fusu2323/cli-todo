# Project Research Summary

**Project:** CLI Todo/Task Manager
**Domain:** Go CLI Application
**Researched:** 2026-04-13
**Confidence:** HIGH

## Executive Summary

This is a local-first CLI task manager built in Go using only the standard library. Experts build this class of tool with a clean separation between the CLI layer (flag parsing and command routing), the domain layer (Task struct and business logic), and the persistence layer (JSON file store). The recommended approach is to start with a minimal feature set (add, list, done, delete) backed by a robust persistence layer that handles first-run gracefully, uses atomic writes, and serializes concurrent access.

The key risks are data loss from race conditions and non-atomic writes, plus poor user experience from silent failures. Mitigation requires a sync.Mutex around all file operations, temp-file-plus-rename for writes, and explicit error output to stderr. The project layout follows Go conventions with `cmd/main.go` as entry point, `internal/task/` for domain logic, and `internal/commands/` for command handlers.

## Key Findings

### Recommended Stack

Go standard library only, no external dependencies. All technologies are production-grade and sufficient for the scope.

**Core technologies:**
- `flag` - CLI argument parsing with built-in -h/--help support
- `os` - File I/O, environment variables, cross-platform operations
- `encoding/json` - JSON serialization with `MarshalIndent` for readable output
- `path/filepath` - Cross-platform path handling (Windows vs Unix)
- `os/user` - Home directory lookup for `~/.todo.json` path
- `testing` - Built-in unit testing with `go test`
- `fmt` - Error formatting with `%w` for wrapped errors

### Expected Features

**Must have (table stakes):**
- Add task with title - core CRUD, generates UUID and creation timestamp
- List all tasks - plain text output showing completion status and ID
- Mark task complete - toggle completion, update completed_at
- Delete task - remove by ID, no confirmation in v1
- Persist to JSON file - atomic writes to prevent corruption
- Category filtering - simple `--category` flag for task triage

**Should have (competitive):**
- Status filtering (`--pending`, `--completed`) - enhances list command
- Task detail view - show full task info on request

**Defer (v2+):**
- Due dates - complex parsing, calendar tool territory
- Priority levels (P1/P2/P3) - cognitive overhead, ranking debates
- Cloud sync - violates local-first core value
- Interactive prompts - breaks scriptability
- Rich text/markdown - display complexity

### Architecture Approach

Four-layer architecture: CLI entry point (`cmd/main.go`) handles flag parsing and command routing via `flag.FlagSet`. Command handlers in `internal/commands/` delegate to the domain layer (`internal/task/`). The persistence layer (`internal/task/store.go`) uses a repository pattern with `JSONFileStore` implementing a `TaskStore` interface. Data flows from user input through main.go to command handlers to the store and back as formatted output.

**Major components:**
1. `cmd/main.go` - Entry point, flag parsing, command routing using `flag.NewFlagSet`
2. `internal/task/task.go` - Task struct, business logic, custom sentinel errors
3. `internal/task/store.go` - JSON persistence with `Load()`/`Save()` interface
4. `internal/commands/*.go` - Command handlers (add, list, done, delete)

### Critical Pitfalls

1. **File I/O race conditions** - Concurrent reads/writes corrupt JSON. Prevention: wrap all file access with `sync.Mutex` at the store level.

2. **Non-atomic writes** - Direct `os.WriteFile` can leave corrupted data on failure. Prevention: write to temp file then `os.Rename`.

3. **File-not-found on first run** - Program panics or errors when `tasks.json` does not exist. Prevention: check `os.ErrNotExist` in `Load()` and return empty slice.

4. **Error wrapping with `%s` instead of `%w`** - Error chains are lost, making debugging impossible. Prevention: always use `%w` for error wrapping.

5. **JSON unknown fields silently dropped** - Users editing JSON manually lose data when `Unmarshal` ignores unknown fields. Prevention: implement custom `UnmarshalJSON` that warns on unknown fields.

## Implications for Roadmap

Based on research, suggested phase structure:

### Phase 1: Foundation (Data Persistence)
**Rationale:** All CRUD operations depend on persistence. Must be correct before any command can work.
**Delivers:** Working Task struct, JSONFileStore with mutex protection, atomic writes, first-run handling
**Addresses:** Add, List, Done, Delete (base functionality)
**Avoids:** Race conditions (mutex), non-atomic writes (temp+rename), file-not-found panic, error wrapping mistakes

### Phase 2: Command Layer (CLI Commands)
**Rationale:** With persistence working, wire up the CLI layer with proper flag parsing
**Delivers:** `cmd/main.go` with FlagSet-based subcommand routing, `--help` support, category filtering
**Uses:** `flag.NewFlagSet` for subcommands, proper `flag.Parse()` timing
**Implements:** Commands delegate to TaskStore; output formatting in handlers

### Phase 3: Polish and Validation
**Rationale:** Core loop is working; add validation, error messages, and edge case handling
**Delivers:** Custom error types (`ErrTaskNotFound`, `ErrInvalidID`), output sanitization, empty state messaging, `errors.Is`/`errors.As` usage

### Phase Ordering Rationale

- **Phase 1 first** because all commands depend on persistence. Getting this right avoids the most painful bugs (data loss).
- **Phase 2 next** because CLI routing is straightforward once persistence is isolated behind an interface.
- **Phase 3 last** because error handling and UX are refinements, not blockers.

### Research Flags

Phases with standard patterns (skip research-phase):
- **Phase 1:** Well-documented persistence patterns (JSON file, mutex, atomic rename) - all from Go stdlib docs
- **Phase 2:** Flag parsing via `flag.NewFlagSet` is standard Go CLI pattern

Phases likely needing deeper research during planning:
- **Phase 3:** None identified - error handling and output formatting are straightforward

## Confidence Assessment

| Area | Confidence | Notes |
|------|------------|-------|
| Stack | HIGH | Go stdlib only, all technologies verified against official docs |
| Features | HIGH | Clear MVP definition, competitor analysis confirms market expectations |
| Architecture | HIGH | Standard Go project layout, established CLI patterns |
| Pitfalls | HIGH | All pitfalls verified with working code examples |

**Overall confidence:** HIGH

### Gaps to Address

- **Project name:** Not specified in research files. Confirm before roadmap creation.
- **File location:** `~/.todo.json` assumed, but `~/.config/todo/` following XDG conventions may be more appropriate for cross-platform. Validate during Phase 1 planning.

## Sources

### Primary (HIGH confidence)
- Go flag package docs (pkg.go.dev/flag) - Parse() timing, FlagSet patterns
- Go os package docs (pkg.go.dev/os) - file I/O, error handling
- Go encoding/json docs (pkg.go.dev/encoding/json) - Unmarshal behavior with unknown fields
- Go error handling guidelines (go.dev/blog/error-handling) - %w wrapping
- Internal packages Go blog (go.dev/blog/internal-packages) - internal/ package convention

### Secondary (MEDIUM confidence)
- XDG Base Directory Specification - file location conventions (referenced in PITFALLS.md security section)

---
*Research completed: 2026-04-13*
*Ready for roadmap: yes*
