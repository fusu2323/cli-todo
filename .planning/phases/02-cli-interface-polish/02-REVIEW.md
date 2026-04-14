---
phase: 02-cli-interface-polish
reviewed: 2026-04-14T00:00:00Z
depth: standard
files_reviewed: 2
files_reviewed_list:
  - internal/task/store.go
  - cmd/main.go
findings:
  critical: 0
  warning: 1
  info: 3
  total: 4
status: issues_found
---

# Phase 02: Code Review Report

**Reviewed:** 2026-04-14
**Depth:** standard
**Files Reviewed:** 2
**Status:** issues_found

## Summary

Reviewed `internal/task/store.go` (161 lines) and `cmd/main.go` (161 lines) for bugs, security issues, and code quality. Both files demonstrate solid Go stdlib compliance with proper use of `flag`, `os`, `encoding/json`, `fmt`, and `errors` packages. Error handling is correct (use of `%w` wrapping, `errors.Is` checks, atomic writes via temp file rename pattern). `go vet` reports no issues and all 12 tests pass.

One warning identified regarding redundant error messages. Three informational items noted for consideration.

## Warnings

### WR-01: Redundant error message text in MarkDone and Delete

**File:** `internal/task/store.go:138` and `internal/task/store.go:156`
**Issue:** Error messages prepend "task not found: %s: " before wrapping `ErrTaskNotFound`, but `ErrTaskNotFound` already contains the text "task not found". This produces verbose messages like:

```
task not found: abc123: task not found
```

**Fix:**
```go
// Line 138 - MarkDone
return fmt.Errorf("task %s: %w", id, ErrTaskNotFound)

// Line 156 - Delete
return fmt.Errorf("task %s: %w", id, ErrTaskNotFound)
```

This preserves `errors.Is` functionality (callers check `errors.Is(err, task.ErrTaskNotFound)`) while producing cleaner error messages.

## Info

### IN-01: Direct os.Args usage in handlers (testing trade-off)

**File:** `cmd/main.go:34,61,89,115`
**Issue:** Each handler calls `fs.Parse(os.Args[2:])` directly rather than accepting `[]string` as a parameter. This couples handlers to global state and makes unit testing harder.

Per CLAUDE.md conventions: "Global flag parsing in init() ... makes testing difficult; parse in main() or explicit function". The current approach is one step removed from `init()` parsing but still couples to `os.Args`.

**Fix (if testing becomes a priority):**
```go
func handleAdd(store *task.JSONFileStore, args []string) {
    fs := flag.NewFlagSet("add", flag.ExitOnError)
    // ... flag setup ...
    fs.Parse(args)  // instead of os.Args[2:]
```

For a learning project prioritizing simplicity, the current approach is acceptable.

### IN-02: Extra positional arguments silently ignored

**File:** `cmd/main.go:35-37` (and similar in other handlers)
**Issue:** After `fs.Parse(os.Args[2:])`, the handler only checks `fs.NArg() < 1` but does not warn if `fs.NArg() > 1`. Users passing `todo add "title" extra-arg` will have the extra argument silently ignored.

**Fix (optional improvement):**
```go
if fs.NArg() < 1 {
    fs.Usage()
    os.Exit(1)
}
if fs.NArg() > 1 {
    fmt.Fprintln(os.Stderr, "error: too many arguments")
    fs.Usage()
    os.Exit(1)
}
```

### IN-03: Task struct dereference in handleAdd is correct but could clarify

**File:** `cmd/main.go:39`
**Issue:** `task.NewTask` returns `*Task` and `store.Add` accepts `*Task`, but the code dereferences with `*task` then re-references with `&*task` at call site. This is a no-op but potentially confusing.

**Fix (for clarity):**
```go
// Current
t, err := task.NewTask(fs.Arg(0), *category)
// ... later
store.Add(t)

// Where store.Add signature is: func (s *JSONFileStore) Add(task *Task) error
```

Pass `t` directly since it is already `*Task`.

---

_Reviewed: 2026-04-14_
_Reviewer: Claude (gsd-code-reviewer)_
_Depth: standard_
