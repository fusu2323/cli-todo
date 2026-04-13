# Feature Research

**Domain:** CLI Todo/Task Manager
**Researched:** 2026-04-13
**Confidence:** HIGH

## Feature Landscape

### Table Stakes (Users Expect These)

Features users assume exist. Missing these = product feels incomplete.

| Feature | Why Expected | Complexity | Notes |
|---------|--------------|------------|-------|
| Add task with title | Core CRUD operation - users must be able to create tasks | LOW | Simple string input, generates UUID, stores creation timestamp |
| List all tasks | Users need to see what they have outstanding | LOW | Plain text output, show completion status, task ID prominently |
| Mark task complete | Basic task lifecycle - todos get done | LOW | Toggle completion status, update completed_at timestamp |
| Delete task | Remove unwanted tasks | LOW | Remove by ID, no confirmation in v1 (can add later) |
| Persist to JSON file | Data must survive restarts | LOW | Single `todos.json` file, atomic writes to prevent corruption |
| View task details | Understand what a task is about | LOW | Show title, creation date, completion status, category if set |

### Differentiators (Competitive Advantage)

Features that set the product apart. Not required, but valuable.

| Feature | Value Proposition | Complexity | Notes |
|---------|-------------------|------------|-------|
| Category/tag filtering | Fast task triage without full list scan | LOW | Simple `--category work` flag, no complex tagging system |
| Local-first architecture | Full offline capability, no account needed | LOW | JSON file in config dir, respects XDG conventions on Unix |
| Sub-second performance | CLI tools must feel instant | LOW | No database, direct file read/write, minimal dependencies |
| Filter by status | View only incomplete or completed tasks | LOW | `--pending` and `--completed` flags |

### Anti-Features (Commonly Requested, Often Problematic)

Features that seem good but create problems.

| Feature | Why Requested | Why Problematic | Alternative |
|---------|---------------|-----------------|-------------|
| Interactive prompts | Feels modern, wizard-like | Breaks scriptability, harder to pipe into other tools, slows down power users | Keep CLI flags simple and script-friendly |
| Due dates | Users want to schedule | Adds validation complexity, date parsing edge cases, "today/tomorrow/tuesday" parsing is a rabbit hole | Defer to calendar tools; this is a todo manager, not a calendar |
| Priority levels (P1, P2, P3) | Users want to rank tasks | Artificial ranking debates, inconsistent usage, adds cognitive overhead | Simple ordering (explicit position field) is more honest |
| Cloud sync | Users want access everywhere | Sync conflicts, auth complexity, server costs, not local-first | Respect the local-first core value |
| Rich text descriptions | Users want detailed tasks | Markdown parsing complexity, display inconsistencies across terminals | Plain text is fine for CLI; descriptions are plain text |

## Feature Dependencies

```
[Add Task]
    └──requires──> [Persist to JSON]

[List Tasks]
    └──requires──> [Persist to JSON]

[Mark Complete]
    └──requires──> [Persist to JSON]

[Delete Task]
    └──requires──> [Persist to JSON]

[Category Filter]
    └──enhances──> [List Tasks]

[Status Filter (pending/completed)]
    └──enhances──> [List Tasks]
```

### Dependency Notes

- **All CRUD operations require [Persist to JSON]:** The file storage is the foundation. Everything builds on read/write operations.
- **[Category Filter] enhances [List Tasks]:** Filtering is a modifier on list, not standalone.
- **[Status Filter] enhances [List Tasks]:** Same - filters are modifiers.

## MVP Definition

### Launch With (v1)

Minimum viable product - what is needed to validate the concept.

- [ ] Add task with title - essential, no todos without it
- [ ] List all tasks - essential, need to see what exists
- [ ] Mark task complete - essential, todos are for completing
- [ ] Delete task - essential, cleanup unwanted tasks
- [ ] Persist to JSON file - essential, data must survive restarts
- [ ] Category filtering - low complexity, high utility for task triage

### Add After Validation (v1.x)

Features to add once core is working.

- [ ] Task detail view - show full task info
- [ ] Bulk complete/delete - operate on multiple tasks
- [ ] Multiple category assignment - one category per task is limiting

### Future Consideration (v2+)

Features to defer until product-market fit is established.

- [ ] Due dates - complex parsing, calendar tool territory
- [ ] Priority levels - adds cognitive overhead, ranking debates
- [ ] Cloud sync - violates local-first core value
- [ ] Interactive prompts - breaks scriptability
- [ ] Rich text/markdown in descriptions - display complexity

## Feature Prioritization Matrix

| Feature | User Value | Implementation Cost | Priority |
|---------|------------|---------------------|----------|
| Add task | HIGH | LOW | P1 |
| List tasks | HIGH | LOW | P1 |
| Mark complete | HIGH | LOW | P1 |
| Delete task | HIGH | LOW | P1 |
| JSON persistence | HIGH | LOW | P1 |
| Category filtering | MEDIUM | LOW | P1 |
| Status filtering (pending/completed) | MEDIUM | LOW | P2 |
| Task detail view | MEDIUM | LOW | P2 |
| Bulk operations | MEDIUM | MEDIUM | P3 |
| Multiple categories | MEDIUM | MEDIUM | P3 |
| Due dates | MEDIUM | HIGH | P3 |
| Priority levels | LOW | MEDIUM | P3 |
| Interactive prompts | LOW | MEDIUM | P3 |

**Priority key:**
- P1: Must have for launch
- P2: Should have, add when possible
- P3: Nice to have, future consideration

## Competitor Feature Analysis

| Feature | todo.txt | Taskwarrior | taskell | Our Approach |
|---------|----------|--------------|---------|--------------|
| Add task | LOW (manual file edit) | HIGH (command) | MEDIUM (org-mode) | HIGH - simple command |
| List tasks | LOW (grep/cat) | HIGH (rich output) | HIGH (visual board) | HIGH - plain text list |
| Mark complete | LOW (manual edit) | HIGH (done) | HIGH (move column) | HIGH (toggle command) |
| Delete task | LOW (manual edit) | HIGH (remove) | HIGH (delete row) | HIGH (remove command) |
| Persistence | File-based | SQLite | File-based | JSON file |
| Categories/tags | LOW (manual) | HIGH (tags) | LOW (list-based) | MEDIUM - simple string field |
| Due dates | LOW (manual) | HIGH (due:) | LOW (future) | Out of scope |
| Priority | LOW (manual A/B/C) | HIGH (priority:) | LOW (manual order) | Out of scope |

**Analysis:**
- **todo.txt:** Unix philosophy, but requires manual file editing - too primitive
- **Taskwarrior:** Feature-rich but steep learning curve, complex configuration
- **taskell:** Markdown/org-mode based, good for visual thinkers but slow

Our approach: Command-line simplicity like todo.txt, but with actual commands. Not as powerful as Taskwarrior, but easier to learn. Local-first JSON persistence is simpler than Taskwarrior's SQLite.

## Sources

- [todo.txt format](https://github.com/todotxt/todo.txt) - simple file-based approach
- [Taskwarrior](https://taskwarrior.org/) - feature-rich CLI task manager
- [taskell](https://taskell.app/) - hierarchical task management
- Go standard library `encoding/json` for JSON handling
- [spf13/cobra](https://github.com/spf13/cobra) for CLI framework (if needed)

---
*Feature research for: Go CLI Todo Manager*
*Researched: 2026-04-13*
