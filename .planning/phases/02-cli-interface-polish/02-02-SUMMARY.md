---
phase: "02-cli-interface-polish"
plan: "02"
subsystem: "cli"
tags: ["cli", "flag", "subcommand"]
dependency_graph:
  requires: ["02-01"]
  provides: ["CLI-01", "CLI-02", "CLI-03", "CLI-04", "CLI-05", "CLI-06", "ERR-01", "ERR-02", "ERR-03"]
  affects: ["cmd/main.go"]
tech_stack:
  added: ["flag.FlagSet"]
  patterns: ["subcommand routing via switch", "flag.ExitOnError", "errors.Is() for sentinel error checking"]
key_files:
  created:
    - path: "cmd/main.go"
      lines: 160
      exports: "main, handleAdd, handleList, handleDone, handleDelete, printGlobalHelp"
decisions:
  - "Used flag.NewFlagSet for each subcommand to enable proper flag parsing per os.Args slicing"
  - "Silent success for add/done/delete per D-01 constraint"
  - "Global help on no args, subcommand help via fs.Usage()"
---

# Phase 02 Plan 02: CLI Entry Point with Subcommand Routing Summary

## One-liner

CLI entry point at `cmd/main.go` with flag.FlagSet subcommand routing for add/list/done/delete commands.

## Tasks Completed

| Task | Name | Commit | Files |
|------|------|--------|-------|
| 1 | Create cmd/main.go with subcommand routing | 6b59bab | cmd/main.go |

## What Was Built

Created `cmd/main.go` (160 lines) with:

- **Subcommand routing**: `switch os.Args[1]` routes to add/list/done/delete/help
- **handleAdd**: Parses `todo add <title> [-c category]`, silent on success
- **handleList**: Parses `todo list [-c category]`, outputs `[x] Title @category`
- **handleDone**: Parses `todo done <id>`, checks `errors.Is(err, task.ErrTaskNotFound)`
- **handleDelete**: Parses `todo delete <id>`, checks `errors.Is(err, task.ErrTaskNotFound)`
- **Global help**: Shown when no args, unknown subcommand prints error + help

## Verification

- `go build ./cmd/...` - PASSED
- `go vet ./...` - PASSED (no warnings)
- `go test ./...` - PASSED (internal/task tests pass)

## Success Criteria

- [x] `todo add "Buy milk"` creates a task (silent)
- [x] `todo add "Work out" -c fitness` creates task with category
- [x] `todo list` shows formatted output `[x] Title @category`
- [x] `todo list -c fitness` filters by category
- [x] `todo done <id>` marks task complete (silent)
- [x] `todo delete <id>` removes task (silent)
- [x] `todo help` shows usage information
- [x] `todo` (no args) shows help
- [x] Invalid ID returns "task not found" to stderr, exit 1
- [x] `go build` passes

## Deviations from Plan

None - plan executed exactly as written.

## Threat Surface Scan

No new security surface introduced. Trust boundary between user input and store already mitigated by atomic writes (DATA-03) in store.go.

## Self-Check: PASSED

- cmd/main.go exists at path
- Commit 6b59bab found in git log
- go build passes
- go vet passes
- go test passes
