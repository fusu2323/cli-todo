# Requirements: CLI Todo Manager

**Defined:** 2026-04-13
**Core Value:** A fast, reliable, local-first todo manager that just works. No accounts, no cloud sync.

## v1 Requirements

### Task Management

- [ ] **TASK-01**: User can add a new task with a title and optional category
- [ ] **TASK-02**: User can list all tasks (optionally filtered by category)
- [ ] **TASK-03**: User can mark a task as complete by ID
- [ ] **TASK-04**: User can delete a task by ID
- [ ] **TASK-05**: Tasks persist to ~/.todo.json between sessions

### Categories

- [ ] **CAT-01**: Tasks can have an optional category string field
- [ ] **CAT-02**: List command can filter tasks by category

### Data Integrity

- [ ] **DATA-01**: Concurrent reads/writes don't corrupt the JSON file (mutex protection)
- [ ] **DATA-02**: First run (file doesn't exist) returns empty task list without error
- [ ] **DATA-03**: Atomic writes — write to temp file then rename to prevent corruption

### Error Handling

- [ ] **ERR-01**: Invalid task ID returns clear "task not found" error
- [ ] **ERR-02**: Corrupted JSON file returns clear error message (not panic)
- [ ] **ERR-03**: All errors wrapped with context using %w pattern

### CLI Interface

- [ ] **CLI-01**: `todo add <title> [-c category]` adds a task
- [ ] **CLI-02**: `todo list [-c category]` lists tasks
- [ ] **CLI-03**: `todo done <id>` marks task complete
- [ ] **CLI-04**: `todo delete <id>` removes task
- [ ] **CLI-05**: `todo help` shows usage information
- [ ] **CLI-06**: Help text shown when no arguments provided

## v2 Requirements

Deferred to future release. Tracked but not in current roadmap.

### Due Dates

- **DUE-01**: Tasks can have an optional due date
- **DUE-02**: List command can filter tasks by due date

### Priority

- **PRIO-01**: Tasks can have a priority (high/medium/low)
- **PRIO-02**: List command sorts by priority

## Out of Scope

Explicitly excluded. Documented to prevent scope creep.

| Feature | Reason |
|---------|--------|
| Cloud sync / multi-device | Against "local-first" core value |
| Interactive prompts | CLI-first design — flags and args only |
| Rich task fields (description, notes) | Keep v1 minimal — extend via v2 |
| Due dates / priorities | Add complexity; defer to v2 |
| Multiple todo files | Single file at ~/.todo.json is simpler |

## Traceability

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

**Coverage:**
- v1 requirements: 19 total
- Mapped to phases: 19
- Unmapped: 0 ✓

---
*Requirements defined: 2026-04-13*
*Last updated: 2026-04-13 after initial definition*
