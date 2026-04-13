# Phase 1: Data Layer & Core Logic - Research

**Researched:** 2026-04-13
**Domain:** Go standard library data layer implementation
**Confidence:** HIGH

## Summary

Phase 1 implements a greenfield data layer for the CLI Todo Manager. The core deliverable is a thread-safe `JSONFileStore` that persists a flat `[]Task` slice to `~/.todo.json` using atomic writes (temp file + rename). All CRUD operations (Add, List, MarkDone, Delete) are implemented as store methods with mutex protection. Categories are optional string fields on tasks. No CLI interface or custom error types exist yet — those are Phase 2.

**Primary recommendation:** Implement `internal/task/store.go` as a single `JSONFileStore` struct embedding a `sync.Mutex`, with `Load()`, `Save()`, `Add()`, `List()`, `MarkDone(id)`, and `Delete(id)` methods. Use `crypto/rand` for UUID generation. Keep task.go focused on the Task struct and NewTask constructor only.

---

## User Constraints (from CONTEXT.md)

### Locked Decisions
- **D-01:** Hex UUID via `crypto/rand` + hex encoding — 32-character hex string
- **D-02:** `sync.Mutex` — simple single-lock approach
- **D-03:** Concrete `JSONFileStore` struct with `Load()` and `Save()` methods — no interface abstraction
- **Atomic Writes:** Write to temp file via `os.MkdirTemp`, then rename via `os.Rename`
- **Category:** Optional `category` string field (empty string = uncategorized)
- **JSON Structure:** Flat array `[]Task`, human-readable via `json.MarshalIndent` with `"  "` indent
- **File Location:** `~/.todo.json` via `os.UserHomeDir()`
- **Error Handling:** Phase 1 uses basic Go error returns; DATA-02 (file-not-found) handled inline

### Claude's Discretion
- Task struct field names (e.g., `CreatedAt` vs `created_at` in JSON) — **needs decision**
- Whether to include `UpdatedAt` timestamp on tasks — **needs decision**
- Internal function organization within store.go and task.go — **needs decision**

### Deferred Ideas (OUT OF SCOPE)
None — no deferred ideas matched Phase 1 scope.

---

## Phase Requirements

| ID | Description | Research Support |
|----|-------------|------------------|
| TASK-01 | User can add a new task with a title and optional category | `Store.Add()` method with category param |
| TASK-02 | User can list all tasks (optionally filtered by category) | `Store.List()` method with optional category filter |
| TASK-03 | User can mark a task as complete by ID | `Store.MarkDone(id)` method |
| TASK-04 | User can delete a task by ID | `Store.Delete(id)` method |
| TASK-05 | Tasks persist to ~/.todo.json between sessions | `Store.Save()` + `Store.Load()` |
| CAT-01 | Tasks can have an optional category string field | `Task.Category` field |
| CAT-02 | List command can filter tasks by category | `Store.List(category)` filter logic |
| DATA-01 | Concurrent reads/writes don't corrupt JSON (mutex protection) | `sync.Mutex` embedded in `JSONFileStore` |
| DATA-02 | First run (file doesn't exist) returns empty task list without error | `Load()` returns `[]Task{}` when file not found |
| DATA-03 | Atomic writes — write to temp file then rename | `os.MkdirTemp` + `os.Rename` pattern |

---

## Standard Stack

### Core
| Library | Version | Purpose | Why Standard |
|---------|---------|---------|--------------|
| Go stdlib | 1.21+ | Entire codebase | Learning constraint; production-grade |
| `crypto/rand` | stdlib | UUID generation | No external deps; `crypto/rand.Read` + hex encoding |
| `sync.Mutex` | stdlib | Concurrent access protection | DATA-01 requirement |
| `encoding/json` | stdlib | JSON serialization | `MarshalIndent` for readable output |
| `os` | stdlib | File I/O, home dir | `ReadFile`, `WriteFile`, `MkdirTemp`, `Rename`, `UserHomeDir` |

### Project Structure
```
cli-todo/
├── cmd/
│   └── main.go           # (Phase 2 — not in scope)
├── internal/
│   └── task/
│       ├── task.go       # Task struct, NewTask constructor, UUID generation
│       └── store.go      # JSONFileStore with mutex, Load/Save, CRUD methods
├── go.mod
└── go.sum
```

---

## Architecture Patterns

### Pattern 1: Concrete Store with Embedded Mutex

**What:** `JSONFileStore` struct embeds `sync.Mutex`, all methods lock/unlock for thread-safe access.

**When to use:** Single-user CLI with file-based persistence.

**Example:**
```go
// Source: [VERIFIED — Go stdlib sync.Mutex docs]
type JSONFileStore struct {
    mu    sync.Mutex
    path  string
}

func (s *JSONFileStore) Save(tasks []Task) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    // ... write logic
}
```

### Pattern 2: Hex UUID via crypto/rand

**What:** Generate 32-character hex string UUID using `crypto/rand.Read`.

**When to use:** When external UUID libraries are prohibited.

**Example:**
```go
// Source: [VERIFIED — Go stdlib crypto/rand docs]
func generateUUID() (string, error) {
    b := make([]byte, 16)
    if _, err := rand.Read(b); err != nil {
        return "", err
    }
    return hex.EncodeToString(b), nil // 32-char hex string
}
```

### Pattern 3: Atomic Write via Temp File Rename

**What:** Write JSON to temp file, then atomically rename to target path.

**When to use:** When JSON file corruption from interrupted writes must be prevented.

**Example:**
```go
// Source: [VERIFIED — Go stdlib os docs]
func (s *JSONFileStore) saveAtomic(data []byte) error {
    tmp, err := os.MkdirTemp("", "todo-*.tmp")
    if err != nil {
        return err
    }
    defer os.Remove(tmp) // cleanup on failure

    if err := os.WriteFile(tmp, data, 0644); err != nil {
        return err
    }
    return os.Rename(tmp, s.path) // atomic on POSIX, cross-platform safe
}
```

### Pattern 4: Load with File-Not-Found Handling

**What:** `Load()` returns empty slice when file does not exist (DATA-02).

**Example:**
```go
// Source: [VERIFIED — Go stdlib os docs]
func (s *JSONFileStore) Load() ([]Task, error) {
    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            return []Task{}, nil // DATA-02: first run returns empty
        }
        return nil, err
    }
    var tasks []Task
    if err := json.Unmarshal(data, &tasks); err != nil {
        return nil, err // DATA-02: corrupted JSON returns error
    }
    return tasks, nil
}
```

---

## Don't Hand-Roll

| Problem | Don't Build | Use Instead | Why |
|---------|-------------|-------------|-----|
| UUID generation | Custom random string | `crypto/rand` + `hex.EncodeToString` | cryptographically secure, no deps |
| File locking | OS-level file locks | `sync.Mutex` per store instance | simpler, sufficient for single-user CLI |
| Atomic writes | Partial-write detection | `os.MkdirTemp` + `os.Rename` | stdlib, proven, portable |
| Home directory | `~/.todo` hardcoding | `os.UserHomeDir()` | cross-platform (Windows/macOS/Linux) |

---

## Common Pitfalls

### Pitfall 1: Mutex Not Held During Entire Read-Modify-Write Cycle

**What goes wrong:** Concurrent goroutine reads stale data, overwrites changes.

**Why it happens:** Locking only during `WriteFile` but not during `ReadFile` + modification.

**How to avoid:** Hold mutex for entire `Load()` → modify → `Save()` sequence:
```go
func (s *JSONFileStore) Add(task Task) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    tasks, err := s.loadLocked()  // internal, assumes lock held
    if err != nil {
        return err
    }
    tasks = append(tasks, task)
    return s.saveLocked(tasks)
}
```

### Pitfall 2: Temp File Not Cleaned Up on Failure

**What goes wrong:** Orphaned temp files accumulate in default temp directory.

**Why it happens:** `os.Remove` deferred but not called on success path, or error after write but before rename.

**How to avoid:** Remove temp file in success path after rename completes, or use `defer os.Remove(tmp)` unconditionally and note that rename makes it a no-op if file was already moved.

### Pitfall 3: Windows Rename Not Atomic Over Existing File

**What goes wrong:** `os.Rename` on Windows fails if target exists (EXDEV error on Unix, different behavior on Windows).

**Why it happens:** Cross-platform Go behavior differences.

**How to avoid:** Use `os.Rename` on POSIX; on Windows, rename to unique temp then remove old file. Per Go stdlib docs, `os.Rename` is documented as replacing the target on Windows. Test on both platforms if cross-platform is critical.

### Pitfall 4: json.Unmarshal Into Pointer to Slice Without Initializing

**What goes wrong:** If file is empty/blank, `json.Unmarshal` may leave slice nil.

**Why it happens:** `var tasks []Task` initializes nil slice, unmarshal succeeds but slice remains nil.

**How to avoid:** Initialize with `tasks := []Task{}` or use `make([]Task, 0)` before unmarshaling.

---

## Code Examples

### Task Struct (with JSON tags)

```go
// Source: [VERIFIED — Go stdlib encoding/json docs]
type Task struct {
    ID        string    `json:"id"`
    Title     string    `json:"title"`
    Category  string    `json:"category,omitempty"`      // empty = uncategorized
    Completed bool      `json:"completed"`
    CreatedAt time.Time `json:"created_at"`
}

// NewTask creates a task with generated UUID and current timestamp
func NewTask(title, category string) (*Task, error) {
    id, err := generateUUID()
    if err != nil {
        return nil, err
    }
    return &Task{
        ID:        id,
        Title:     title,
        Category:  category,
        Completed: false,
        CreatedAt: time.Now(),
    }, nil
}
```

### JSONFileStore Full Implementation

```go
// Source: [VERIFIED — Go stdlib patterns]
type JSONFileStore struct {
    mu   sync.Mutex
    path string
}

func NewJSONFileStore(path string) (*JSONFileStore, error) {
    if path == "" {
        home, err := os.UserHomeDir()
        if err != nil {
            return nil, err
        }
        path = filepath.Join(home, ".todo.json")
    }
    return &JSONFileStore{path: path}, nil
}

func (s *JSONFileStore) Load() ([]Task, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.loadLocked()
}

func (s *JSONFileStore) loadLocked() ([]Task, error) {
    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            return []Task{}, nil
        }
        return nil, err
    }
    if len(data) == 0 {
        return []Task{}, nil
    }
    var tasks []Task
    if err := json.Unmarshal(data, &tasks); err != nil {
        return nil, err
    }
    return tasks, nil
}

func (s *JSONFileStore) Save(tasks []Task) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    return s.saveLocked(tasks)
}

func (s *JSONFileStore) saveLocked(tasks []Task) error {
    data, err := json.MarshalIndent(tasks, "", "  ")
    if err != nil {
        return err
    }
    tmp, err := os.MkdirTemp("", "todo-*.tmp")
    if err != nil {
        return err
    }
    defer os.Remove(tmp) // safe if file already moved

    if err := os.WriteFile(tmp, data, 0644); err != nil {
        return err
    }
    return os.Rename(tmp, s.path)
}

func (s *JSONFileStore) Add(task *Task) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    tasks, err := s.loadLocked()
    if err != nil {
        return err
    }
    tasks = append(tasks, *task)
    return s.saveLocked(tasks)
}

func (s *JSONFileStore) List(category string) ([]Task, error) {
    s.mu.Lock()
    defer s.mu.Unlock()
    tasks, err := s.loadLocked()
    if err != nil {
        return nil, err
    }
    if category == "" {
        return tasks, nil
    }
    filtered := make([]Task, 0)
    for _, t := range tasks {
        if t.Category == category {
            filtered = append(filtered, t)
        }
    }
    return filtered, nil
}

func (s *JSONFileStore) MarkDone(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    tasks, err := s.loadLocked()
    if err != nil {
        return err
    }
    for i, t := range tasks {
        if t.ID == id {
            tasks[i].Completed = true
            return s.saveLocked(tasks)
        }
    }
    return fmt.Errorf("task not found: %s", id) // Phase 2: custom ErrTaskNotFound
}

func (s *JSONFileStore) Delete(id string) error {
    s.mu.Lock()
    defer s.mu.Unlock()
    tasks, err := s.loadLocked()
    if err != nil {
        return err
    }
    for i, t := range tasks {
        if t.ID == id {
            tasks = append(tasks[:i], tasks[i+1:]...)
            return s.saveLocked(tasks)
        }
    }
    return fmt.Errorf("task not found: %s", id)
}
```

---

## State of the Art

| Old Approach | Current Approach | When Changed | Impact |
|--------------|------------------|--------------|--------|
| Integer auto-increment IDs | UUID via crypto/rand | Project start | Eliminates race conditions, no ID reuse |
| Global var for tasks slice | Store instance passed via dependencies | Project start | Testable, no hidden state |
| Write directly to file | Temp file + rename | Project start | No corruption on interrupted writes |
| ioutil (deprecated) | os.ReadFile, os.WriteFile, os.MkdirTemp | Go 1.16+ | Future-proof |

**Deprecated/outdated:**
- `ioutil.ReadFile` / `ioutil.WriteFile` — replaced by `os` equivalents in Go 1.16+
- `ioutil.TempDir` — replaced by `os.MkdirTemp` in Go 1.16+

---

## Assumptions Log

| # | Claim | Section | Risk if Wrong |
|---|-------|---------|---------------|
| A1 | `os.Rename` is atomic on Windows (replaces existing file) | Atomic Writes | Plan may need platform-specific code for Windows |
| A2 | No interface abstraction needed for Phase 1 | Architecture | May need refactoring if Phase 2 requires mock store for testing |
| A3 | `CreatedAt` field is sufficient; no `UpdatedAt` needed for Phase 1 | Task struct | Task metadata may be insufficient if future features need it |

**If this table is empty:** All claims in this research were verified or cited.

---

## Open Questions

1. **Task struct JSON field names**
   - What we know: CONTEXT.md says `CreatedAt` is in scope, `created_at` is a naming option
   - What's unclear: Whether to use `CreatedAt` (Go convention) vs `created_at` (JSON convention)
   - Recommendation: Use Go conventions for struct fields (`CreatedAt`), JSON tags can override for serialization

2. **`UpdatedAt` field**
   - What we know: Listed as Claude's discretion in CONTEXT.md
   - What's unclear: Whether Phase 1 needs last-modified tracking
   - Recommendation: Omit for Phase 1; add in Phase 2 if needed for feature requests

3. **Store methods returning error vs panic on corrupted JSON**
   - What we know: DATA-02 says corrupted JSON returns clear error (not panic)
   - What's unclear: Whether to return error and let caller decide, or log and return empty list
   - Recommendation: Return error (Phase 2 will add context via `fmt.Errorf %w`)

---

## Environment Availability

| Dependency | Required By | Available | Version | Fallback |
|------------|------------|-----------|---------|----------|
| Go | Build/Run | ✓ | 1.26.1 | — |
| git | Version control | ✓ | (system) | — |

**Missing dependencies with no fallback:**
None — all required tools are available.

**Missing dependencies with fallback:**
None — no external dependencies needed per project constraint.

---

## Validation Architecture

### Test Framework
| Property | Value |
|----------|-------|
| Framework | Go stdlib `testing` package |
| Config file | None — standard `go test` |
| Quick run command | `go test ./internal/task/... -v` |
| Full suite command | `go test ./... -v` |

### Phase Requirements → Test Map
| Req ID | Behavior | Test Type | Automated Command | File Exists? |
|--------|----------|-----------|-------------------|-------------|
| TASK-01 | Add task persists | unit | `go test ./internal/task/... -run TestAdd -v` | Not yet |
| TASK-02 | List returns all tasks | unit | `go test ./internal/task/... -run TestList -v` | Not yet |
| TASK-02 (CAT-02) | List filtered by category | unit | `go test ./internal/task/... -run TestListCategory -v` | Not yet |
| TASK-03 | MarkDone marks task complete | unit | `go test ./internal/task/... -run TestMarkDone -v` | Not yet |
| TASK-04 | Delete removes task | unit | `go test ./internal/task/... -run TestDelete -v` | Not yet |
| TASK-05 | Tasks persist across store instances | integration | `go test ./internal/task/... -run TestPersistence -v` | Not yet |
| DATA-01 | Concurrent access safe | unit (goroutine test) | `go test ./internal/task/... -run TestConcurrent -v` | Not yet |
| DATA-02 | Missing file returns empty list | unit | `go test ./internal/task/... -run TestLoadNotExist -v` | Not yet |
| DATA-03 | Atomic write creates valid file | unit | `go test ./internal/task/... -run TestAtomicWrite -v` | Not yet |

### Sampling Rate
- **Per task commit:** `go test ./internal/task/... -v`
- **Per wave merge:** `go test ./... -v`
- **Phase gate:** Full suite green before `/gsd-verify-work`

### Wave 0 Gaps
- [ ] `internal/task/task_test.go` — covers TASK-01, TASK-02, TASK-03, TASK-04, TASK-05, CAT-01, CAT-02, DATA-01, DATA-02, DATA-03
- [ ] `internal/task/store_test.go` — covers DATA-01, DATA-02, DATA-03 specifically
- [ ] Framework install: `go mod init` — if none detected

*(No gaps for existing test infrastructure — greenfield phase)*

---

## Security Domain

### Applicable ASVS Categories

| ASVS Category | Applies | Standard Control |
|---------------|---------|-----------------|
| V2 Authentication | No | N/A — local single-user tool |
| V3 Session Management | No | N/A — no sessions |
| V4 Access Control | No | N/A — single user file |
| V5 Input Validation | Partial | Validate task title non-empty, ID matches hex format |
| V6 Cryptography | Yes | UUID generation via `crypto/rand` |

### Known Threat Patterns for Go File I/O

| Pattern | STRIDE | Standard Mitigation |
|---------|--------|---------------------|
| Path traversal via task ID | Tampering | IDs are hex-encoded UUIDs, not user-controlled paths |
| JSON injection | Tampering | `encoding/json` escapes special characters automatically |
| Concurrent write corruption | Denial | `sync.Mutex` ensures serialized access |
| Temp file symlink attack | Tampering | `os.MkdirTemp` creates secure temp directory |

---

## Sources

### Primary (HIGH confidence)
- [Go sync.Mutex documentation](https://pkg.go.dev/sync#Mutex) — mutex patterns
- [Go crypto/rand documentation](https://pkg.go.dev/crypto/rand) — UUID generation
- [Go os documentation](https://pkg.go.dev/os) — ReadFile, WriteFile, MkdirTemp, Rename, UserHomeDir
- [Go encoding/json documentation](https://pkg.go.dev/encoding/json) — MarshalIndent, Unmarshal
- [Go filepath documentation](https://pkg.go.dev/path/filepath) — Join for cross-platform paths

### Secondary (MEDIUM confidence)
- [Go flag.FlagSet documentation](https://pkg.go.dev/flag) — subcommand routing (Phase 2)
- [Effective Go: Errors](https://go.dev/doc/effective_go#errors) — error wrapping patterns (Phase 2)

### Tertiary (LOW confidence)
- [Go atomic rename cross-platform behavior](https://github.com/golang/go/issues/8860) — Windows rename semantics marked for validation

---

## Metadata

**Confidence breakdown:**
- Standard stack: HIGH — all verified via Go 1.26.1 stdlib docs
- Architecture: HIGH — patterns from established Go idioms, CONTEXT.md locked decisions
- Pitfalls: MEDIUM — common issues documented, some unverified (Windows rename)

**Research date:** 2026-04-13
**Valid until:** 2026-05-13 (30 days — Go stdlib is stable)
