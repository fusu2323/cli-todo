# CLI Todo Manager - Roadmap

## Core Value

A fast, reliable, local-first todo manager that just works. No accounts, no cloud sync.

## Phases

- [ ] **Phase 1: Data Layer & Core Logic** - Task struct, JSON persistence with mutex, all CRUD operations, category support, atomic writes
- [ ] **Phase 2: CLI Interface & Polish** - flag parsing, subcommands, help text, custom errors, error handling

---

## Phase Details

### Phase 1: Data Layer & Core Logic

**Goal**: Tasks persist reliably with category support and concurrent safety

**Depends on**: Nothing (first phase)

**Requirements**: TASK-01, TASK-02, TASK-03, TASK-04, TASK-05, CAT-01, CAT-02, DATA-01, DATA-02, DATA-03

**Success Criteria** (what must be TRUE):

1. User can add a task with a title and optional category, and it persists after application restart
2. User can list all tasks and see all previously added tasks
3. User can mark a task as complete by ID and see it marked on subsequent list operations
4. User can delete a task by ID and it no longer appears in list operations
5. First run (no existing file) returns an empty task list without error
6. Concurrent reads and writes don't corrupt the JSON data file

**Plans**: 3 plans

Plans:
- [x] 01-01-PLAN.md — Task model: go.mod, Task struct, NewTask constructor, UUID generation
- [x] 01-02-PLAN.md — JSONFileStore: mutex-protected CRUD, atomic writes, category filtering
- [x] 01-03-PLAN.md — Concurrency tests: concurrent access, atomic write verification

---

### Phase 2: CLI Interface & Polish

**Goal**: Users can interact with tasks through a clean CLI with clear feedback

**Depends on**: Phase 1

**Requirements**: ERR-01, ERR-02, ERR-03, CLI-01, CLI-02, CLI-03, CLI-04, CLI-05, CLI-06

**Success Criteria** (what must be TRUE):

1. User can run `todo add <title>` and see the task added successfully
2. User can run `todo add <title> -c <category>` to add a task with a category
3. User can run `todo list` to see all tasks formatted clearly
4. User can run `todo list -c <category>` to filter and see only tasks in that category
5. User can run `todo done <id>` to mark a task complete
6. User can run `todo delete <id>` to remove a task
7. User can run `todo help` or execute `todo` with no arguments to see usage information
8. Invalid task ID returns a clear "task not found" error message
9. Corrupted JSON file returns a clear error message without crashing
10. All error messages include context (wrapped errors with %w)

**Plans**: TBD

---

## Progress Table

| Phase | Plans Complete | Status | Completed |
|-------|----------------|--------|-----------|
| 1. Data Layer & Core Logic | 0/3 | Not started | - |
| 2. CLI Interface & Polish | 0/10 | Not started | - |

---

## Coverage

- **Total v1 requirements**: 19
- **Mapped to phases**: 19
- **Unmapped**: 0

| Requirement | Phase | Status |
|-------------|-------|--------|
| TASK-01 | Phase 1 | Pending |
| TASK-02 | Phase 1 | Pending |
| TASK-03 | Phase 1 | Pending |
| TASK-04 | Phase 1 | Pending |
| TASK-05 | Phase 1 | Pending |
| CAT-01 | Phase 1 | Pending |
| CAT-02 | Phase 1 | Pending |
| DATA-01 | Phase 1 | Pending |
| DATA-02 | Phase 1 | Pending |
| DATA-03 | Phase 1 | Pending |
| ERR-01 | Phase 2 | Pending |
| ERR-02 | Phase 2 | Pending |
| ERR-03 | Phase 2 | Pending |
| CLI-01 | Phase 2 | Pending |
| CLI-02 | Phase 2 | Pending |
| CLI-03 | Phase 2 | Pending |
| CLI-04 | Phase 2 | Pending |
| CLI-05 | Phase 2 | Pending |
| CLI-06 | Phase 2 | Pending |
