# Phase 2: CLI Interface & Polish - Research

**Researched:** 2026-04-14
**Domain:** Go standard library CLI interface implementation
**Confidence:** HIGH

## Summary

Phase 2 implements a CLI interface on top of Phase 1's data layer. The deliverable is `cmd/main.go` with subcommand routing via `flag.FlagSet`, custom error types with proper error wrapping, and auto-generated help text. Key technical challenges are: (1) wiring `ErrTaskNotFound` into `MarkDone`/`Delete` with `%w` wrapping, (2) wrapping JSON parse errors with context, (3) routing subcommands (add/list/done/delete) via separate FlagSets, and (4) printing help when no arguments provided.

**Primary recommendation:** Create `cmd/main.go` that uses a top-level FlagSet to parse the subcommand name, then delegates to subcommand-specific FlagSets. Use `fmt.Errorf` with `%w` to wrap `ErrTaskNotFound` and JSON errors. Print plain text errors to stderr and exit 1 on any error.

---

## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Plain stderr errors, exit 1 on error. No colored output.
- **D-02:** Use `flag.FlagSet.PrintDefaults()` for auto-generated help.
- **D-03:** Plain text list format: `[x] Title @category`
- **D-04:** Show help automatically when no arguments provided.
- **D-05:** Use `fmt.Errorf` with `%w` for wrapped errors.
- **D-06:** Use `flag.FlagSet` for subcommand routing.

### Claude's Discretion
- Specific flag names and short flags (e.g., `-c` for category, `--help` is automatic from flag package)
- Exact error message wording for "task not found" and "corrupted file"
- Whether `done` and `delete` commands should print success feedback or be silent

### Deferred Ideas (OUT OF SCOPE)
None.

---

## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| ERR-01 | Invalid task ID returns "task not found" error | `ErrTaskNotFound` wired with `%w` wrapping in MarkDone/Delete |
| ERR-02 | Corrupted JSON file returns clear error message | JSON parse error wrapped with context via `fmt.Errorf` |
| ERR-03 | All errors wrapped with context using %w pattern | `fmt.Errorf` with `%w` verb throughout store.go |
| CLI-01 | `todo add <title> [-c category]` adds a task | FlagSet for "add" subcommand with title arg and -c flag |
| CLI-02 | `todo list [-c category]` lists tasks | FlagSet for "list" subcommand with optional -c flag |
| CLI-03 | `todo done <id>` marks task complete | FlagSet for "done" subcommand with id arg |
| CLI-04 | `todo delete <id>` removes task | FlagSet for "delete" subcommand with id arg |
| CLI-05 | `todo help` shows usage information | FlagSet.PrintDefaults() + custom preamble |
| CLI-06 | No arguments shows help | Check len(os.Args) < 2 and print help before parsing |

---

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go stdlib | 1.21+ | Entire codebase | Learning constraint; production-grade |
| `flag` | stdlib | CLI argument parsing | Built-in, handles -h/--help automatically |
| `fmt` | stdlib | Error formatting | `fmt.Errorf` with `%w` for error wrapping |
| `os` | stdlib | stderr, exit codes | `fmt.Fprintln(os.Stderr, ...)` for errors |

### Project Structure
```
cli-todo/
├── cmd/
│   └── main.go           # Phase 2: CLI entry point with subcommand routing
├── internal/
│   └── task/
│       ├── task.go       # Phase 1: Task struct, NewTask constructor
│       ├── store.go      # Phase 1: JSONFileStore with Add/List/MarkDone/Delete
│       └── errors.go     # Phase 2: ErrTaskNotFound, ErrCorruptedFile definitions
├── go.mod
└── go.sum
```

---

## Architecture Patterns

### Pattern 1: FlagSet Subcommand Routing

**What:** Use `flag.NewFlagSet` for each subcommand, parse `os.Args[1:]` by identifying subcommand first, then delegate to subcommand's FlagSet.

**When to use:** CLI tools with multiple subcommands that each have their own flags.

**Verified:** Go stdlib `flag` package behavior via `go doc flag.NewFlagSet` and `go doc flag.FlagSet.Parse`.

**Example:**
```go
func main() {
    if len(os.Args) < 2 {
        printGlobalHelp()
        os.Exit(0)
    }

    switch os.Args[1] {
    case "add":
        handleAdd()
    case "list":
        handleList()
    case "done":
        handleDone()
    case "delete":
        handleDelete()
    case "help":
        printGlobalHelp()
    default:
        fmt.Fprintln(os.Stderr, "unknown subcommand:", os.Args[1])
        printGlobalHelp()
        os.Exit(1)
    }
}

func handleAdd() {
    fs := flag.NewFlagSet("add", flag.ExitOnError)
    category := fs.String("c", "", "category")
    fs.Parse(os.Args[2:])
    if fs.NArg() < 1 {
        fs.Usage()
        os.Exit(1)
    }
    title := fs.Arg(0)
    // ... call store.Add with title and *category
}
```

### Pattern 2: Error Wrapping with %w

**What:** Use `fmt.Errorf` with `%w` verb to wrap sentinel errors so `errors.Is()` works.

**When to use:** When you need to add context to an error while preserving the ability to check for specific error types.

**Verified:** `go doc fmt.Errorf` confirms `%w` creates error with Unwrap method.

**Example:**
```go
// In store.go — MarkDone and Delete currently use:
return fmt.Errorf("task not found: %s", id)  // Phase 1 style

// Phase 2 style — wrap with ErrTaskNotFound:
return fmt.Errorf("task not found: %s: %w", id, ErrTaskNotFound)

// Then CLI can check:
if errors.Is(err, task.ErrTaskNotFound) {
    fmt.Fprintln(os.Stderr, "task not found")
    os.Exit(1)
}
```

### Pattern 3: JSON Error Wrapping

**What:** Wrap `json.Unmarshal` errors with context so users see "corrupted file" not raw syntax error.

**When to use:** When loading JSON file that may be corrupted.

**Example:**
```go
// In store.go loadLocked():
if err := json.Unmarshal(data, &tasks); err != nil {
    return nil, fmt.Errorf("corrupted todo file: %w", err)
}
```

### Pattern 4: Auto-Help via PrintDefaults

**What:** Use `flag.FlagSet.PrintDefaults()` for auto-generated flag help text.

**When to use:** When you want help without writing it manually.

**Verified:** `flag.FlagSet.PrintDefaults()` outputs standard Go flag formatting.

**Example:**
```go
func printAddHelp() {
    fmt.Println("Usage: todo add <title> [-c category]")
    fmt.Println()
    fmt.Println("Flags:")
    fs := flag.NewFlagSet("add", flag.ContinueOnError)
    fs.String("c", "", "category")
    fs.PrintDefaults()
}
```

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| CLI flag parsing | Custom argument parser | `flag.FlagSet` | Built-in, handles -h/--help, portable |
| Error wrapping | `errors.New` for wrapped errors | `fmt.Errorf` with `%w` | `%w` enables `errors.Is()` checks |
| Help text | Hand-written help strings | `PrintDefaults()` | Auto-updates when flags change |

---

## Common Pitfalls

### Pitfall 1: Parsing Flags After Subcommand Without Slicing Correctly

**What goes wrong:** Flags like `-c category` don't parse correctly because the subcommand name is still in the args.

**Why it happens:** `os.Args` includes program name and subcommand. Subcommand's FlagSet must parse `os.Args[2:]` (skipging program name and subcommand).

**How to avoid:** Always pass `os.Args[2:]` to subcommand FlagSet.Parse(), not `os.Args[1:]`.

### Pitfall 2: Exit Code 0 for Errors

**What goes wrong:** Error messages go to stderr but the program exits with 0 (success).

**Why it happens:** Forgetting to call `os.Exit(1)` after error handling.

**How to avoid:** Always exit with 1 when an error occurs (per D-01).

### Pitfall 3: Missing Args Check Before Parsing

**What goes wrong:** `fs.Parse()` panics or returns error when no arguments provided to subcommand that requires them.

**Why it happens:** Not checking `fs.NArg()` before accessing `fs.Arg(0)`.

**How to avoid:** Validate required positional args after `fs.Parse()` returns:
```go
if fs.NArg() < 1 {
    fs.Usage()
    os.Exit(1)
}
title := fs.Arg(0)
```

### Pitfall 4: ErrHelp Not Handled

**What goes wrong:** When user passes `-h` or `--help`, FlagSet.Parse returns `flag.ErrHelp` but program continues.

**Why it happens:** Not checking the return value of `fs.Parse()`.

**How to avoid:** Check and exit when `flag.ErrHelp` is returned (this is automatic with `flag.ExitOnError` error handling).

---

## Code Examples

### cmd/main.go Skeleton

```go
// Source: [VERIFIED — Go stdlib flag docs]
package main

import (
    "errors"
    "flag"
    "fmt"
    "os"

    "github.com/user/cli-todo/internal/task"
)

func main() {
    if len(os.Args) < 2 {
        printGlobalHelp()
        os.Exit(0)
    }

    store, err := task.NewJSONFileStore("")
    if err != nil {
        fmt.Fprintln(os.Stderr, "failed to initialize store:", err)
        os.Exit(1)
    }

    switch os.Args[1] {
    case "add":
        handleAdd(store)
    case "list":
        handleList(store)
    case "done":
        handleDone(store)
    case "delete":
        handleDelete(store)
    case "help":
        printGlobalHelp()
    default:
        fmt.Fprintln(os.Stderr, "unknown subcommand:", os.Args[1])
        printGlobalHelp()
        os.Exit(1)
    }
}
```

### handleAdd Implementation

```go
// Source: [VERIFIED — Go stdlib flag docs]
func handleAdd(store *task.JSONFileStore) {
    fs := flag.NewFlagSet("add", flag.ExitOnError)
    category := fs.String("c", "", "category")
    fs.Usage = func() {
        fmt.Println("Usage: todo add <title> [-c category]")
        fmt.Println()
        fmt.Println("Flags:")
        fs.PrintDefaults()
    }
    if err := fs.Parse(os.Args[2:]); err != nil {
        if errors.Is(err, flag.ErrHelp) {
            os.Exit(0)
        }
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
    if fs.NArg() < 1 {
        fs.Usage()
        os.Exit(1)
    }
    title := fs.Arg(0)

    t, err := task.NewTask(title, *category)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error creating task:", err)
        os.Exit(1)
    }
    if err := store.Add(t); err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
    // No success output per D-01 (silent success)
}
```

### handleList Implementation (D-03 format)

```go
// Source: [VERIFIED — D-03 decision in CONTEXT.md]
func handleList(store *task.JSONFileStore) {
    fs := flag.NewFlagSet("list", flag.ExitOnError)
    category := fs.String("c", "", "category")
    fs.Usage = func() {
        fmt.Println("Usage: todo list [-c category]")
        fmt.Println()
        fmt.Println("Flags:")
        fs.PrintDefaults()
    }
    if err := fs.Parse(os.Args[2:]); err != nil {
        if errors.Is(err, flag.ErrHelp) {
            os.Exit(0)
        }
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }

    tasks, err := store.List(*category)
    if err != nil {
        fmt.Fprintln(os.Stderr, "error:", err)
        os.Exit(1)
    }
    for _, t := range tasks {
        check := " "
        if t.Completed {
            check = "x"
        }
        if t.Category != "" {
            fmt.Printf("[%s] %s @%s\n", check, t.Title, t.Category)
        } else {
            fmt.Printf("[%s] %s\n", check, t.Title)
        }
    }
}
```

### Error Wrapping in store.go (Phase 2 Changes)

```go
// Source: [VERIFIED — go doc fmt.Errorf]
// MarkDone — change line 138 from:
//   return fmt.Errorf("task not found: %s", id)
// To:
return fmt.Errorf("task not found: %s: %w", id, ErrTaskNotFound)

// Delete — change line 156 from:
//   return fmt.Errorf("task not found: %s", id)
// To:
return fmt.Errorf("task not found: %s: %w", id, ErrTaskNotFound)

// loadLocked — change line 53-54 from:
//   return nil, err
// To:
return nil, fmt.Errorf("corrupted todo file: %w", err)
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| `fmt.Errorf` without `%w` | `fmt.Errorf` with `%w` for wrapped errors | Phase 2 | Enables `errors.Is()` checks |
| Plain error strings | Sentinel errors (`ErrTaskNotFound`) | Phase 2 | Programmatic error detection |
| No CLI interface | FlagSet-based subcommand CLI | Phase 2 | User-facing interface |

---

## Assumptions Log

> List all claims tagged `[ASSUMED]` in this research. The planner and discuss-phase use this
> section to identify decisions that need user confirmation before execution.

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | `fmt.Fprintln(os.Stderr, ...)` is acceptable for plain error output | Common Pitfalls | Works on all platforms, no additional deps |
| A2 | `flag.ExitOnError` error handling is appropriate for subcommand FlagSets | Architecture | Could use `ContinueOnError` and handle manually, but `ExitOnError` is simpler |
| A3 | No success message for `done`/`delete` commands (silent success) | Architecture | Per D-01 "plain stderr errors" — but silent success was Claude's discretion |

**If this table is empty:** All claims in this research were verified or cited — no user confirmation needed.

---

## Open Questions

1. **Silent success for done/delete**
   - What we know: D-01 says "plain stderr errors, exit 1 on error". Silent success is implied but not explicitly stated.
   - What's unclear: Should `done` and `delete` commands print "done" or "deleted" on success, or be completely silent?
   - Recommendation: Silent success (no output on success) — aligns with minimal UX and D-01's plain approach.

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | Build/Run | Yes | 1.26.1 | — |
| git | Version control | Yes | (system) | — |

**Missing dependencies with no fallback:**
None — pure Go stdlib project.

**Missing dependencies with fallback:**
None — no external dependencies needed per project constraint.

---

## Sources

### Primary (HIGH confidence)
- [Go flag package docs](https://pkg.go.dev/flag) — `NewFlagSet`, `Parse`, `PrintDefaults` behavior verified via `go doc`
- [Go fmt package docs](https://pkg.go.dev/fmt) — `Errorf` with `%w` verified via `go doc`
- [Effective Go: Errors](https://go.dev/doc/effective_go#errors) — error wrapping patterns

### Secondary (MEDIUM confidence)
- [Go stdlib flag source](https://cs.opensource.google/go/go/+/refs/tags/go1.26.1:src/flag/flag.go) — FlagSet implementation (training knowledge, not directly verified)

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all verified via Go 1.26.1 stdlib `go doc`
- Architecture: HIGH — FlagSet subcommand patterns well-established in Go ecosystem
- Pitfalls: HIGH — common CLI mistakes, all verifiable from stdlib behavior

**Research date:** 2026-04-14
**Valid until:** 2026-05-14 (30 days — Go stdlib is stable)
