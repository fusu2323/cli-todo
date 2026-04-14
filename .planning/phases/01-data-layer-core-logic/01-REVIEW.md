---
status: issues_found
findings:
  critical: 0
  warning: 2
  minor: 3
---

# Phase 01: Code Review Report

**Reviewed:** 2026-04-14T00:00:00Z
**Depth:** standard
**Files Reviewed:** 5
- internal/task/task.go
- internal/task/task_test.go
- internal/task/store.go
- internal/task/store_test.go
- go.mod

**Status:** issues_found

## Summary

The Phase 1 data layer implementation is well-structured with solid fundamentals: atomic writes via temp file rename, mutex-protected concurrent access, cryptographically secure UUID generation, and comprehensive test coverage including concurrent access tests. No critical bugs or security vulnerabilities were found. However, there are 2 warnings and 3 minor issues to address before Phase 2.

---

## Warnings

### W-01: Defined error constant `ErrTaskNotFound` is not used

**File:** `internal/task/task.go:138` and `internal/task/task.go:156`
**Issue:** `ErrTaskNotFound` is declared on line 160 but never used. Both `MarkDone` and `Delete` use `fmt.Errorf("task not found: %s", id)` instead of the exported error constant. Callers cannot use `errors.Is()` to check for this specific error condition.
**Fix:**
```go
// MarkDone (line 138)
return ErrTaskNotFound

// Delete (line 156)
return ErrTaskNotFound
```
Note: The comment on line 138 says "Phase 2: custom ErrTaskNotFound" - but the constant is already defined. Either use it now or remove it.

---

### W-02: Test code ignores errors from `NewTask` and `store.Add`

**File:** `internal/task/task_test.go:73-74`, `internal/task/task_test.go:134`
**File:** `internal/task/store_test.go:399`, `internal/task/store_test.go:411-412`
**Issue:** Multiple test cases silently discard errors using `_`:
```go
task1, _ := NewTask("Task 1", "work")   // line 73
task2, _ := NewTask("Task 2", "home")   // line 74
task2, _ := NewTask("Second Task", "testing")  // line 134
```
If `NewTask` returns an error (e.g., `crypto/rand` failure), the test continues with a nil task, causing a panic when dereferenced later.
**Fix:**
```go
task1, err := NewTask("Task 1", "work")
if err != nil {
    t.Fatalf("NewTask failed: %v", err)
}
```

---

## Minor Issues

### M-01: Task `Title` has no validation

**File:** `internal/task/store.go:30`
**Issue:** `NewTask` accepts any title string, including empty strings. Empty titles may cause confusing behavior in the UI.
**Fix:** Consider adding validation:
```go
func NewTask(title, category string) (*Task, error) {
    if strings.TrimSpace(title) == "" {
        return nil, errors.New("title cannot be empty")
    }
    // ...
}
```

### M-02: Regex compilation in tests repeated on each call

**File:** `internal/task/task_test.go:23`, `internal/task/task_test.go:67`
**Issue:** `regexp.MustCompile()` is called inside test functions. While not a performance concern for tests, it could mask compilation errors if the pattern is invalid.
**Fix:** Move compiled regex to package-level variables:
```go
var hex32Pattern = regexp.MustCompile(`^[a-f0-9]{32}$`)
// Then use hex32Pattern.MatchString(id)
```

### M-03: `filepath.Dir` behavior on Windows with forward-slash paths

**File:** `internal/task/task.go:75`
**Issue:** On Windows, `filepath.Dir("/custom/path/todo.json")` returns `/custom/path` which is treated as a relative path. While not a practical issue since paths are constructed correctly, it could cause confusion if custom paths are passed with forward slashes on Windows.
**Fix:** This is not currently causing bugs but be aware if testing on Windows with custom paths.

---

## Positive Findings

- **Atomic writes**: `saveLocked` correctly uses `os.CreateTemp` + `os.Rename` for crash-safe writes
- **Concurrent access**: `sync.Mutex` properly protects all operations; `defer Unlock()` pattern is correct
- **Cryptographic UUID**: Uses `crypto/rand` (not `math/rand`) for secure ID generation
- **Error wrapping**: Uses `%w` pattern (via `fmt.Errorf`) for proper error chain
- **Resource cleanup**: `defer os.Remove(tmp)` ensures temp files are cleaned up even on error
- **Test coverage**: Includes concurrent access tests (`TestConcurrentAccess`) and atomic write tests (`TestAtomicWrite`)

---

_Reviewed: 2026-04-14_
_Reviewer: Claude (gsd-code-reviewer)_
_Depth: standard_
