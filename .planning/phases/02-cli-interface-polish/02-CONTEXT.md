# Phase 02: CLI Interface & Polish - Context

**Gathered:** 2026-04-14
**Status:** Ready for planning

<domain>
## Phase Boundary

Deliver a working CLI interface with subcommands (add, list, done, delete), custom error types with wrapped errors (ERR-01, ERR-02, ERR-03), and help text (CLI-05, CLI-06). This phase builds on Phase 1's data layer.

</domain>

<decisions>
## Implementation Decisions

### Error Presentation (ERR-01, ERR-02, ERR-03)
- **D-01:** Plain stderr — simple plain text errors on stderr, exit 1 on error. No colored output. Keeps things simple and cross-platform.

### Help Text Style (CLI-05)
- **D-02:** Go flag default — use `flag.FlagSet.PrintDefaults()` for auto-generated help from flag definitions. No custom help text authoring needed.

### List Output Format (CLI-02)
- **D-03:** Plain text list — format: `[x] Title @category`. Simple, readable, no table formatting or JSON output.

### No-Arguments Behavior (CLI-06)
- **D-04:** Show help on no args — when user runs `todo` with no arguments, display help automatically. Friendly UX that helps users discover the interface.

### Error Wrapping (ERR-03)
- **D-05:** Use `fmt.Errorf` with `%w` for all wrapped errors. Phase 1 left `ErrTaskNotFound` as a placeholder — Phase 2 wires it up properly.

### CLI Structure
- **D-06:** Use `flag.FlagSet` for subcommand routing (carried from Phase 1 decision). Each subcommand has its own FlagSet.

### Claude's Discretion
- Specific flag names and short flags (e.g., `-c` for category, `--help` is automatic from flag package)
- Exact error message wording for "task not found" and "corrupted file"
- Whether `done` and `delete` commands should print success feedback or be silent

</decisions>

<canonical_refs>
## Canonical References

**Downstream agents MUST read these before planning or implementing.**

- `.planning/REQUIREMENTS.md` — ERR-01, ERR-02, ERR-03, CLI-01 through CLI-06
- `.planning/ROADMAP.md` — Phase 2 goal and success criteria
- `.planning/phases/01-data-layer-core-logic/01-CONTEXT.md` — Phase 1 decisions (flag.FlagSet, ErrTaskNotFound placeholder)

</canonical_refs>

<code_context>
## Existing Code Insights

### Reusable Assets
- `internal/task/store.go` — `JSONFileStore` struct with `Add`, `List`, `MarkDone`, `Delete` methods
- `internal/task/store.go` — `ErrTaskNotFound` variable already defined (placeholder to wire up)
- `internal/task/task.go` — `Task` struct with `ID`, `Title`, `Category`, `Completed`, `CreatedAt` fields

### Established Patterns
- `flag.FlagSet` for subcommand routing (Phase 1 decision)
- Go standard library only — no external dependencies
- `sync.Mutex` for concurrent access protection (Phase 1)

### Integration Points
- `cmd/main.go` (to be created) — will import `internal/task` packages
- CLI commands call `JSONFileStore` methods: `Add`, `List`, `MarkDone`, `Delete`
- `flag.FlagSet` subcommands: `add`, `list`, `done`, `delete`

</code_context>

<specifics>
## Specific Ideas

No specific examples or references provided — open to standard approaches.

</specifics>

<deferred>
## Deferred Ideas

None — discussion stayed within phase scope.

</deferred>

---

*Phase: 02-cli-interface-polish*
*Context gathered: 2026-04-14*
