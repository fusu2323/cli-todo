# Pitfalls Research

**Domain:** CLI Todo/Task Manager in Go
**Researched:** 2026-04-13
**Confidence:** HIGH

## Critical Pitfalls

### Pitfall 1: File I/O Race Conditions

**What goes wrong:**
Concurrent reads and writes to the JSON file cause data corruption or lost updates. User A reads the file, User B reads the file, User A writes, User B writes — User A's changes are lost. The JSON file can end up truncated or contain malformed data.

**Why it happens:**
Go's `os.Open` and `json.Marshal/Unmarshal` are not safe for concurrent access to the same file. Without synchronization, simultaneous operations race on the file descriptor and underlying bytes. The Go memory model does not protect against concurrent file access without explicit synchronization primitives.

**How to avoid:**
Use a `sync.Mutex` to serialize all file access, or use `os.OpenFile` with the `os.O_EXCL` flag for exclusive locking. For a simple CLI, a package-level mutex protecting read/write operations is sufficient:

```go
var fileMu sync.Mutex

func saveTasks(tasks []Task) error {
    fileMu.Lock()
    defer fileMu.Unlock()
    // safe to access file now
}
```

Alternatively, use `os.Rename` to atomically replace the old file after writing to a temp file.

**Warning signs:**
- "unexpected end of JSON input" errors
- Tasks disappearing after concurrent invocations
- File size grows to 0 bytes occasionally

**Phase to address:**
Phase 1 — foundational data model and persistence layer.

---

### Pitfall 2: JSON Unmarshaling with Unknown Fields

**What goes wrong:**
After releasing v1, users add custom fields to the JSON file (e.g., `"priority": "high"`). A later version using a typed struct with `json.Unmarshal` silently drops those fields, losing user data on the next save.

**Why it happens:**
By default, Go's `json.Unmarshal` ignores JSON properties that have no corresponding field in the target struct. This is the "unknown fields" problem — it seems safe until users hand-edit their todo file or migrate data between versions.

**How to avoid:**
Use `map[string]any` for flexible storage, or implement a custom `UnmarshalJSON` method that warns on unknown fields:

```go
func (t *Task) UnmarshalJSON(data []byte) error {
    var m map[string]any
    if err := json.Unmarshal(data, &m); err != nil {
        return err
    }
    // Check for unknown fields
    for k := range m {
        if k != "id" && k != "title" && k != "done" {
            fmt.Fprintf(os.Stderr, "warning: unknown field %q ignored\n", k)
        }
    }
    // Then unmarshal into typed struct
    type taskAlias Task
    return json.Unmarshal(data, (*taskAlias)(t))
}
```

Alternatively, use `map[string]any` as an intermediate representation, then validate and copy known fields to typed structs.

**Warning signs:**
- Users reporting lost metadata after updates
- No warning when loading files with extra fields
- Struct tags do not include `,unknownkeyphrase=ignore` behavior is silent

**Phase to address:**
Phase 1 — data persistence design.

---

### Pitfall 3: Error Wrapping Antipatterns

**What goes wrong:**
Error messages lose context, making debugging impossible. Instead of "failed to save tasks: permission denied", users see just "permission denied" with no indication of what operation failed.

**Why it happens:**
Using `fmt.Errorf` with `%s` or `%v` converts the error to a string, losing the underlying error type. Using `%w` wraps the error properly, preserving it for `errors.Is` and `errors.As`:

```go
// WRONG — loses error chain
return fmt.Errorf("failed to save tasks: %s", err)

// CORRECT — wraps the error
return fmt.Errorf("failed to save tasks: %w", err)
```

**How to avoid:**
Always use `%w` when wrapping errors. Reserve `%s` only for non-error values. Use `errors.Is` and `errors.As` to inspect wrapped errors:

```go
if err := saveTasks(tasks); err != nil {
    if errors.Is(err, os.ErrPermission) {
        fmt.Fprintln(os.Stderr, "Error: cannot write to todo file — permission denied")
        os.Exit(1)
    }
    return fmt.Errorf("add command: %w", err)
}
```

**Warning signs:**
- Error messages that do not indicate which operation failed
- Testing errors with `==` instead of `errors.Is`
- `fmt.Println(err)` instead of `fmt.Fprintln(os.Stderr, err)`

**Phase to address:**
Phase 1 — error handling conventions from the start.

---

### Pitfall 4: flag.Parse() Timing

**What goes wrong:**
`flag.Parse()` is called before all flags are registered, causing panic ("flag already declared") or silently ignoring flags that appear after the parse call. Users' flags do not work as expected.

**Why it happens:**
`flag.Parse()` consumes `os.Args[1:]` immediately. Any flag registered after `Parse()` is called is not recognized. Additionally, the `flag` package has a single global namespace — declaring the same flag twice causes a panic.

**How to avoid:**
Register all flags in `init()` or at the top of `main()` before any call to `flag.Parse()`:

```go
var (
    doneFlag = flag.Bool("done", false, "mark task as done")
    listFlag = flag.Bool("list", false, "list all tasks")
)

func main() {
    flag.Parse()
    // now use *doneFlag, *listFlag
}
```

For subcommands, consider using a flag subpackage (e.g., `github.com/urfave/cli`) or manually parse arguments.

**Warning signs:**
- Panic with "flag already declared"
- Flags defined in functions other than `main` or `init`
- `flag.PanicOnError` behavior in tests

**Phase to address:**
Phase 2 — command structure and argument parsing.

---

### Pitfall 5: Not Handling File-Not-Found on First Run

**What goes wrong:**
The program panics or exits with an unhelpful error on first run because the JSON file does not exist yet. The expected behavior is to create an empty task list and continue silently.

**Why it happens:**
`os.Open("tasks.json")` returns an error when the file does not exist, but many tutorials skip error handling for this case, assuming the file always exists. On first run, `err != nil` and the program fails.

**How to avoid:**
Check for `os.ErrNotExist` explicitly and handle it as a signal to initialize an empty list:

```go
func loadTasks() ([]Task, error) {
    data, err := os.ReadFile("tasks.json")
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            return []Task{}, nil // first run — empty list is fine
        }
        return nil, fmt.Errorf("load tasks: %w", err)
    }
    var tasks []Task
    if err := json.Unmarshal(data, &tasks); err != nil {
        return nil, fmt.Errorf("parse tasks: %w", err)
    }
    return tasks, nil
}
```

**Warning signs:**
- "no such file or directory" error on first invocation
- No error handling around `os.Open` or `os.ReadFile`
- Test failures when running without a pre-existing data file

**Phase to address:**
Phase 1 — data persistence layer must handle first-run gracefully.

---

### Pitfall 6: Reading File Before Writing

**What goes wrong:**
The program reads the file, modifies the slice in memory, then writes back — but if the write fails (disk full, permissions), the in-memory changes are lost. User loses work.

**Why it happens:**
The read-modify-write cycle is not atomic. A failure during write leaves the file in an indeterminate state (could be empty, partial, or old content) while the program has already moved on with the new data in memory.

**How to avoid:**
Write to a temporary file first, then atomically rename over the original:

```go
func saveTasksAtomic(tasks []Task) error {
    tmp, err := os.CreateTemp("", "tasks-*.tmp")
    if err != nil {
        return fmt.Errorf("create temp file: %w", err)
    }
    defer os.Remove(tmp.Name())

    if err := json.NewEncoder(tmp).Encode(tasks); err != nil {
        tmp.Close()
        return fmt.Errorf("encode tasks: %w", err)
    }
    if err := tmp.Close(); err != nil {
        return fmt.Errorf("close temp file: %w", err)
    }
    if err := os.Rename(tmp.Name(), "tasks.json"); err != nil {
        return fmt.Errorf("rename temp file: %w", err)
    }
    return nil
}
```

**Warning signs:**
- Data loss when disk is full
- Corrupted JSON after write failures
- No backup of previous state

**Phase to address:**
Phase 1 — data persistence layer.

---

## Technical Debt Patterns

| Shortcut | Immediate Benefit | Long-term Cost | When Acceptable |
|----------|-------------------|----------------|-----------------|
| Using `json.RawMessage` for everything | Simple, no type definitions | No validation, silent failures on bad data | Never in production |
| Ignoring `flag.Parse()` errors | Saves 3 lines of code | Panics crash the program silently | Never |
| Global `[]Task` variable | No passing through functions | Untestable, race-prone | Only in throwaway prototypes |
| Using `ioutil.ReadFile` instead of `os.ReadFile` | None — `ioutil` is deprecated | Will break in future Go versions | Replace immediately |
| Skipping `json.MarshalIndent` in favor of `Marshal` | Slightly smaller file | Human-editing the file becomes harder | Only for high-volume writes where size matters |

## Integration Gotchas

This project is a standalone CLI with no external integrations. No integration gotchas apply.

## Performance Traps

| Trap | Symptoms | Prevention | When It Breaks |
|------|----------|------------|----------------|
| Loading entire file into memory on every command | Slow startup with large task files | File is small; not a concern for todo app | Breaks at ~100k tasks (multi-MB JSON) — not a realistic concern for a CLI todo app |
| Inefficient JSON encoding in tight loops | CPU spike during `saveTasks` | Use `json.NewEncoder` with a buffer | Only with thousands of saves per second — unrealistic for CLI |
| No file locking for concurrent access | Data loss, corrupted file | Use `sync.Mutex` around file operations | When user runs multiple instances simultaneously |

## Security Mistakes

| Mistake | Risk | Prevention |
|---------|------|------------|
| Storing tasks in a world-readable location | Other users on the same system can read private tasks | Use `$XDG_DATA_HOME` or `$HOME/.local/share` — respect XDG base directory spec |
| No sanitization of task content | Task content is printed to terminal — potential for output injection if content contains terminal control characters | Sanitize output before printing, or use a safe print method that treats content as text, not terminal codes |
| File permissions too permissive (0o777) | Other users can modify the task file | Create file with 0o600 (`-rw-------`) — owner read/write only |

## UX Pitfalls

| Pitfall | User Impact | Better Approach |
|---------|-------------|-----------------|
| Silent failures on save | User thinks task was added, but it was lost | Print error to stderr and exit non-zero when save fails |
| No indication when task list is empty | User does not know if the app works | Print a friendly message: "No tasks yet. Add one with `todo add <task>`" |
| Unhelpful error messages | User cannot diagnose problems | Include the operation in error: "failed to load tasks: permission denied" not just "permission denied" |
| Mixing stdout and stderr | User cannot redirect output cleanly | Command output goes to stdout, errors and progress to stderr |
| No `--help` or `-h` flag | User cannot discover available commands | At minimum, print usage when no arguments given or when `--help` is passed |

## "Looks Done But Isn't" Checklist

- [ ] **Persistence:** First run does not panic — empty file is created if missing
- [ ] **Persistence:** Write uses atomic rename (temp file + os.Rename), not direct write
- [ ] **Persistence:** Concurrent invocations do not corrupt the file (mutex or exclusive lock)
- [ ] **JSON handling:** Unknown fields in the JSON file are detected and warned about, not silently dropped
- [ ] **Error handling:** All errors are wrapped with `%w`, never `%s` or `%v`
- [ ] **Error handling:** `errors.Is` and `errors.As` are used for error inspection, not `==`
- [ ] **Flags:** All flags are registered before `flag.Parse()` is called
- [ ] **Flags:** Duplicate flag registration causes a clear error, not a panic
- [ ] **Output:** Errors go to stderr, normal output to stdout
- [ ] **Output:** Terminal control characters in task content are sanitized before printing

## Recovery Strategies

| Pitfall | Recovery Cost | Recovery Steps |
|---------|---------------|----------------|
| Corrupted JSON file | MEDIUM | Rename corrupted file to `tasks.json.bak`, create new empty `tasks.json`. User loses recent changes. |
| Lost updates due to race condition | HIGH | No automatic recovery. Prevent with mutex. User must re-add lost tasks. |
| Silent data loss on failed save | MEDIUM | Only recoverable if user has backup. Print error and exit non-zero to prevent silent failures. |

## Pitfall-to-Phase Mapping

| Pitfall | Prevention Phase | Verification |
|---------|------------------|--------------|
| File-not-found on first run | Phase 1: Data persistence | Run `todo list` with no existing file — should show empty list, not error |
| File I/O race conditions | Phase 1: Data persistence | Run multiple concurrent `todo add` commands, verify all tasks present |
| JSON unknown fields silently dropped | Phase 1: Data persistence | Add unknown field to JSON file, load and resave — verify warning printed and field preserved or explicitly handled |
| Error wrapping with %s instead of %w | Phase 1: Error handling | Inspect all `fmt.Errorf` calls — none should use %s or %v for errors |
| flag.Parse() timing | Phase 2: Command structure | Verify all flags work when called in any order |
| Atomic write (temp + rename) | Phase 1: Data persistence | Simulate write failure mid-write (e.g., disk full) — verify original file intact |
| File permissions too permissive | Phase 1: Data persistence | Check file permissions — should be 0o600 |

## Sources

- Go Standard Library: `flag` package documentation — Parse() must be called after all flag registrations
- Go Standard Library: `os` package — file operations and error handling patterns
- Go Standard Library: `json` package — Unmarshal behavior with unknown fields
- Go Blog: "Error handling in Go" — covers %w wrapping and error chains
- XDG Base Directory Specification — standard location for user data files on Unix systems
- OWASP: "Command Injection" — sanitizing output for terminal safety

---

*Pitfalls research for: CLI Todo/Task Manager in Go*
*Researched: 2026-04-13*
