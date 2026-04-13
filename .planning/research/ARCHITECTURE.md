# Architecture Research

**Domain:** CLI Todo/Task Manager in Go
**Researched:** 2026-04-13
**Confidence:** HIGH

## Standard Architecture

### System Overview

```
┌─────────────────────────────────────────────────────────────┐
│                        CLI Layer                            │
├─────────────────────────────────────────────────────────────┤
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────┐        │
│  │   add   │  │  list   │  │  done   │  │ delete  │        │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────┬────┘        │
│       │            │            │            │              │
├───────┴────────────┴────────────┴────────────┴──────────────┤
│                     Command Layer                           │
│  ┌─────────────────────────────────────────────────────┐    │
│  │              flag.FlagSet (subcommands)              │    │
│  └─────────────────────────────────────────────────────┘    │
├─────────────────────────────────────────────────────────────┤
│                       Domain Layer                           │
│  ┌─────────────────────┐  ┌─────────────────────────────┐  │
│  │    Task struct      │  │   Custom error types        │  │
│  │   - ID               │  │   - ErrTaskNotFound         │  │
│  │   - Description      │  │   - ErrInvalidID            │  │
│  │   - Completed        │  │                             │  │
│  │   - CreatedAt        │  │                             │  │
│  └─────────────────────┘  └─────────────────────────────┘  │
├─────────────────────────────────────────────────────────────┤
│                      Persistence Layer                       │
│  ┌──────────────────────────────────────────────────────┐  │
│  │   JSONFileStore (~/.todo.json)                         │  │
│  │   - os/user for cross-platform home dir               │  │
│  │   - json.MarshalIndent / json.Unmarshal               │  │
│  └──────────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Responsibility | Typical Implementation |
|-----------|----------------|------------------------|
| `cmd/main.go` | Entry point, flag parsing, command routing | `flag.NewFlagSet`, switch on `os.Args[1]` |
| `internal/commands/*.go` | Command-specific logic, validation | Functions that delegate to TaskStore |
| `internal/task/task.go` | Task struct, business logic, custom errors | Value receiver methods |
| `internal/task/store.go` | JSON file persistence | `json.MarshalIndent`, `json.Unmarshal` |

## Recommended Project Structure

```
cli-todo/
├── cmd/
│   └── main.go           # Entry point, flag parsing, command routing
├── internal/
│   ├── task/
│   │   ├── task.go       # Task struct, methods, custom errors
│   │   └── store.go      # JSON file persistence layer
│   └── commands/
│       ├── add.go        # add command logic
│       ├── list.go       # list command logic
│       ├── done.go       # mark complete command logic
│       └── delete.go     # delete command logic
├── go.mod
└── go.sum
```

### Structure Rationale

- **`cmd/main.go`:** Keeps entry point minimal. Handles only flag parsing and routing, delegating all business logic to the `internal` packages. This makes the application testable and follows Go conventions.
- **`internal/task/`:** Encapsulates all domain logic. `task.go` holds the `Task` struct and its methods. `store.go` handles persistence, keeping the serialization format (JSON) isolated from domain logic.
- **`internal/commands/`:** Each command (add, list, done, delete) is a separate file. This follows the Go idiom of one concern per file and makes it easy to add new commands without touching existing code.
- **`internal/` prefix:** Go's `internal` package convention prevents external imports, enforcing that only `cmd/main.go` can access these packages.

## Architectural Patterns

### Pattern 1: Subcommand Routing via Flag Sets

**What:** Use `flag.NewFlagSet` with `flag.ExitOnError` for each subcommand, allowing standard Go flag parsing to handle argument validation.

**When to use:** CLI tools with multiple commands (add, list, delete).

**Trade-offs:**
- Pros: Standard library, no dependencies, consistent with Go idioms
- Cons: Limited to positional arguments after flags; more complex help formatting

**Example:**
```go
func main() {
    if len(os.Args) < 2 {
        fmt.Fprintln(os.Stderr, "usage: todo <add|list|done|delete> [args]")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "add":
        addCmd := flag.NewFlagSet("add", flag.ExitOnError)
        desc := addCmd.String("desc", "", "task description")
        addCmd.Parse(os.Args[2:])
        // ... validate and execute
    case "list":
        listCmd := flag.NewFlagSet("list", flag.ExitOnError)
        listCmd.Parse(os.Args[2:])
        // ...
    }
}
```

### Pattern 2: Custom Error Types for Domain Errors

**What:** Define sentinel errors (`var ErrTaskNotFound = errors.New("task not found")`) and custom error types for domain-specific failures.

**When to use:** When operations can fail with domain-specific reasons that callers need to handle.

**Trade-offs:**
- Pros: Callers can use `errors.Is()` to check specific failures; descriptive
- Cons: Requires discipline to create meaningful error messages

**Example:**
```go
// task.go
var (
    ErrTaskNotFound = errors.New("task not found")
    ErrInvalidID    = errors.New("invalid task ID")
    ErrEmptyDesc    = errors.New("task description cannot be empty")
)

type Task struct {
    ID          string    `json:"id"`
    Description string    `json:"description"`
    Completed   bool      `json:"completed"`
    CreatedAt   time.Time `json:"created_at"`
}
```

### Pattern 3: Repository Pattern for Persistence

**What:** Abstract data access behind an interface (`TaskStore`) with concrete implementation (`JSONFileStore`).

**When to use:** When persistence might change (file to database, etc.).

**Trade-offs:**
- Pros: Testable, flexible, clear separation of concerns
- Cons: Slight indirection; overkill for simple single-file storage

**Example:**
```go
// store.go
type TaskStore interface {
    Load() ([]Task, error)
    Save([]Task) error
}

type JSONFileStore struct {
    path string
}

func (s *JSONFileStore) Load() ([]Task, error) {
    data, err := os.ReadFile(s.path)
    if err != nil {
        if os.IsNotExist(err) {
            return []Task{}, nil
        }
        return nil, err
    }
    var tasks []Task
    err = json.Unmarshal(data, &tasks)
    return tasks, err
}
```

### Pattern 4: Functional Options for Configuration

**What:** Use `Option` functions to configure structs with defaults that can be overridden.

**When to use:** When initializing structs that may need configuration (e.g., store path).

**Trade-offs:**
- Pros: Flexible, readable, easy to add new options without breaking existing callers
- Cons: Slightly more code than direct field initialization

**Example:**
```go
type StoreOption func(*JSONFileStore)

func WithPath(path string) StoreOption {
    return func(s *JSONFileStore) {
        s.path = path
    }
}

func NewJSONFileStore(opts ...StoreOption) *JSONFileStore {
    s := &JSONFileStore{path: defaultPath()}
    for _, opt := range opts {
        opt(s)
    }
    return s
}
```

## Data Flow

### Request Flow

```
[User CLI Input]
    ↓
[main.go] → Parse os.Args, create FlagSet
    ↓
[Command Handler (e.g., add.go)] → Validate flags, build Task
    ↓
[TaskStore.Load()] → Read ~/.todo.json
    ↓
[Business Logic] → Modify tasks slice
    ↓
[TaskStore.Save()] → Write ~/.todo.json
    ↓
[Success/Error Output] → fmt.Println or fmt.Fprintln(os.Stderr)
```

### State Management

```
~/.todo.json (single source of truth)
    ↑           ↓
  Load()    Save()
    ↑           ↓
[JSONFileStore] ← [Command Handlers]
    ↑
[main.go] ← [User via CLI]
```

### Key Data Flows

1. **Add task:** User runs `todo add -desc "Buy milk"` → Flag parsing validates `-desc` required → New Task created with UUID → Store loads, appends, saves → Success message printed.

2. **List tasks:** User runs `todo list` → Flag parsing (no args needed) → Store loads all tasks → Format and print to stdout.

3. **Complete task:** User runs `todo done 1` → Parse task ID → Store loads tasks → Find task by ID → Set Completed=true → Store saves → Success message printed.

4. **Delete task:** User runs `todo delete 1` → Parse task ID → Store loads tasks → Filter out task with matching ID → Store saves → Success message printed.

## Scaling Considerations

| Scale | Architecture Adjustments |
|-------|--------------------------|
| 0-1k tasks | No changes needed; JSON file handles this well |
| 1k-10k tasks | Consider lazy loading (Load returns pointer, not slice) or background compaction |
| 10k-100k tasks | Switch to embedded database (SQLite via mattn/go-sqlite3) |
| 100k+ tasks | Not appropriate for single-file JSON storage; redesign with database backend |

### Scaling Priorities

1. **First bottleneck:** File size. At ~10MB (roughly 10k tasks), JSON parsing/serialization becomes slow. Mitigation: Switch to line-delimited JSON (jsonl) or SQLite.
2. **Second bottleneck:** Concurrent access. `os.Rename` is not atomic on Windows for existing files. Mitigation: Use file locking (golang.org/x/sys/windows) or switch to SQLite which handles this.

## Anti-Patterns

### Anti-Pattern 1: Global State via Package-Level Variables

**What people do:** Define `var tasks []Task` and `var store *JSONFileStore` at package level.

**Why it's wrong:** Makes the code untestable, creates hidden dependencies between commands, prevents multiple store instances.

**Do this instead:** Pass dependencies via function parameters or struct fields. Use functional options for configuration.

### Anti-Pattern 2: Mixing Presentation and Business Logic

**What people do:** Putting JSON marshaling directly in `main.go` or command handlers.

**Why it's wrong:** If JSON format changes, every command breaks. Hard to test business logic independently of I/O.

**Do this instead:** Keep `internal/task/` focused purely on the `Task` struct and business rules. `internal/task/store.go` handles serialization only. Commands only handle input validation and output formatting.

### Anti-Pattern 3: Silently Ignoring Errors

**What people do:** `data, _ := os.ReadFile(path)` or `json.Unmarshal(data, &tasks)` without checking errors.

**Why it's wrong:** File permission errors, corruption, or JSON syntax errors go unnoticed. Data loss can occur.

**Do this instead:** Always check and handle errors. Return error to caller; let `main.go` decide exit codes.

### Anti-Pattern 4: Incrementing Integer IDs

**What people do:** `nextID := len(tasks) + 1` and `task.ID = strconv.Itoa(nextID)`.

**Why it's wrong:** Race conditions when two processes add simultaneously. ID reuse after deletion causes confusion.

**Do this instead:** Use UUIDs (`crypto/rand` or github.com/google/uuid`) for unique, non-guessable IDs.

## Integration Points

### External Services

| Service | Integration Pattern | Notes |
|---------|---------------------|-------|
| Filesystem | Direct via `os` package | Use `os/user.Current()` for cross-platform home dir |
| Terminal | Direct via `fmt` and `os.Stderr` | Consider `github.com/fatih/color` for colored output (optional, not required) |

### Internal Boundaries

| Boundary | Communication | Notes |
|----------|---------------|-------|
| `main.go` → `commands/*` | Function calls with parsed flags | Commands return error or formatted output |
| `commands/*` → `task.Store` | Interface method calls | Commands do not know about JSON format |
| `task.Store` → filesystem | Direct `os` calls | Isolated in store.go |

## Sources

- [Go flag package documentation](https://pkg.go.dev/flag)
- [Go json package documentation](https://pkg.go.dev/encoding/json)
- [Go os/user package documentation](https://pkg.go.dev/os/user)
- [Go error handling guidelines](https://go.dev/blog/error-handling)
- [Internal packages Go blog](https://go.dev/blog/internal-packages)

---
*Architecture research for: CLI Todo/Task Manager in Go*
*Researched: 2026-04-13*
