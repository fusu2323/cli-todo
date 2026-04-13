# CLI Todo Manager

## What This Is

A command-line todo/task manager for personal productivity. Users manage tasks via terminal commands: add tasks, mark them complete, delete them, and list all tasks. Tasks support categories/tags for organization. Data persists to a JSON file. Primary goal is learning Go fundamentals through a practical project.

## Core Value

A fast, reliable, local-first todo manager that just works. No accounts, no cloud sync — just a simple tool to track what you need to do.

## Requirements

### Validated

(None yet — ship to validate)

### Active

- [ ] User can add a new task with a title and optional category
- [ ] User can list all tasks (optionally filtered by category)
- [ ] User can mark a task as complete
- [ ] User can delete a task
- [ ] Tasks persist to ~/.todo.json between sessions
- [ ] Proper error handling for file I/O, invalid input, missing tasks

### Out of Scope

- Due dates — adds complexity; defer to v2
- Priorities — adds sorting complexity; defer to v2
- Cloud sync or multi-device — against "local-first" core value
- Interactive mode or prompts — CLI-first design
- Tags/categories beyond simple string matching — basic filter is enough for v1

## Context

This is a learning project. The user wants to deepen Go understanding by building a real, working tool. Focus should be on:
- Go standard library patterns (flag for args, json for serialization, os for file I/O)
- Error handling best practices (wrap errors, sentinel errors, custom error types)
- Clean code organization (main.go, task.go, storage.go, commands as functions)
- Testing with Go's standard testing package

The user is comfortable with programming basics and is learning Go specifically.

## Constraints

- **Tech Stack**: Pure Go with only standard library (no external deps except maybe testing libs) — maximizes learning
- **Persistence**: Single JSON file at ~/.todo.json — simple, zero configuration
- **Compatibility**: Linux/macOS/Windows — use path.Join, os/user for cross-platform paths
- **No external dependencies**: Only Go standard library — this is a learning constraint

## Key Decisions

| Decision | Rationale | Outcome |
|----------|-----------|---------|
| flag (standard lib) over cobra | Less magic, more Go-idiomatic, better for learning how CLI args actually work | — Pending |
| JSON over CSV | Extensible, standard library support, easier to add fields later | — Pending |
| ~/.todo.json storage | Cross-platform home dir, persistent, hidden file keeps dir clean | — Pending |
| Categories/tags | Simple string field, not full tagging system — enough for v1 organization | — Pending |

## Evolution

This document evolves at phase transitions and milestone boundaries.

**After each phase transition** (via `/gsd-transition`):
1. Requirements invalidated? → Move to Out of Scope with reason
2. Requirements validated? → Move to Validated with phase reference
3. New requirements emerged? → Add to Active
4. Decisions to log? → Add to Key Decisions
5. "What This Is" still accurate? → Update if drifted

**After each milestone** (via `/gsd-complete-milestone`):
1. Full review of all sections
2. Core Value check — still the right priority?
3. Audit Out of Scope — reasons still valid?
4. Update Context with current state
