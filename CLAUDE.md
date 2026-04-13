<!-- GSD:project-start source:PROJECT.md -->
## Project

**CLI Todo Manager**

A command-line todo/task manager for personal productivity. Users manage tasks via terminal commands: add tasks, mark them complete, delete them, and list all tasks. Tasks support categories/tags for organization. Data persists to a JSON file. Primary goal is learning Go fundamentals through a practical project.

**Core Value:** A fast, reliable, local-first todo manager that just works. No accounts, no cloud sync — just a simple tool to track what you need to do.

### Constraints

- **Tech Stack**: Pure Go with only standard library (no external deps except maybe testing libs) — maximizes learning
- **Persistence**: Single JSON file at ~/.todo.json — simple, zero configuration
- **Compatibility**: Linux/macOS/Windows — use path.Join, os/user for cross-platform paths
- **No external dependencies**: Only Go standard library — this is a learning constraint
<!-- GSD:project-end -->

<!-- GSD:stack-start source:research/STACK.md -->
## Technology Stack

## Recommended Stack
### Core Technologies
| Technology | Version | Purpose | Why Recommended |
|------------|---------|---------|-----------------|
| Go standard library | 1.21+ | Entire codebase | Learning constraint; `flag`, `os`, `encoding/json`, `testing` are production-grade |
| `flag` | stdlib | CLI argument parsing | Built-in, no dependencies, handles -h/--help automatically |
| `encoding/json` | stdlib | JSON serialization | Standard for config/storage files, `MarshalIndent` for readable output |
| `os` | stdlib | File I/O, environment | `os.OpenFile`, `os.ReadFile`, `os.WriteFile`, `os.UserHomeDir` |
| `path/filepath` | stdlib | Cross-platform paths | Handles Windows backslashes vs Unix forward slashes |
| `os/user` | stdlib | Home directory lookup | `user.Current()` for `~/.todo.json` path resolution |
| `testing` | stdlib | Unit testing | Built-in `go test`, `t.Fatal` patterns, benchmark support |
| `fmt` | stdlib | Error formatting | `fmt.Errorf` with `%w` for error wrapping |
### Supporting Libraries
| Library | Version | Purpose | When to Use |
|---------|---------|---------|-------------|
| (none) | - | - | Standard library only constraint |
### Development Tools
| Tool | Purpose | Notes |
|------|---------|-------|
| `go test -v` | Run tests with verbose output | Standard Go testing |
| `go test -cover` | Check test coverage | `go tool cover` for detailed reports |
| `go vet` | Static analysis | Catches common errors before runtime |
| `gofmt` | Code formatting | Enforce consistent style |
| `go run .` | Run without installing | During development |
## Installation
# No packages to install — everything is in stdlib
## Alternatives Considered
| Recommended | Alternative | When to Use Alternative |
|-------------|-------------|-------------------------|
| `flag` stdlib | `spf13/cobra` | When project grows beyond simple CLI (many subcommands) |
| `flag` stdlib | `cli/cli` | When building GitHub-style interactive CLI |
| `encoding/json` | `goccy/go-json` | When performance is critical (rare for todo app) |
| `encoding/json` | `github.com/mitchellh/go-homedir` | When targeting Go < 1.12 (homedir expansion) |
## What NOT to Use
| Avoid | Why | Use Instead |
|-------|-----|-------------|
| External CLI parsing libraries | Violates learning constraint; `flag` is sufficient for simple commands | `flag` |
| `ioutil` functions | Deprecated in Go 1.16+; `os` equivalents preferred | `os.ReadFile`, `os.WriteFile`, `os.MkdirTemp` |
| `errors.New` for wrapped errors | Cannot be unwrapped with `errors.Is`/`errors.As` | `fmt.Errorf` with `%w` |
| Global flag parsing in `init()` | Makes testing difficult; parse in `main()` or explicit function | Parse flags where needed |
## Stack Patterns by Variant
- Use flat `flag` parsing in `main()`
- Single JSON file `~/.todo.json`
- Direct `os.ReadFile`/`os.WriteFile` pattern
- Parse base command first, then subcommand flags
- Consider command pattern with separate handler functions
- May need `flag.FlagSet` for subcommand-level flag parsing
- Add pagination to list command (`-page`, `-limit` flags)
- Consider in-memory cache of JSON file (avoid re-reading on every command)
## Version Compatibility
| Feature | Minimum Go Version | Notes |
|---------|-------------------|-------|
| `os.UserHomeDir` | 1.12 | Simpler than `os/user.Current()` |
| `flag.FlagSet` | 1.0 | For subcommand support |
| `%w` error wrapping | 1.13 | `fmt.Errorf` with `%w` |
| `errors.Is`, `errors.As` | 1.13 | Error hierarchy handling |
| `os.MkdirTemp` | 1.16 | Replaces `ioutil.TempDir` |
## Sources
- [Go flag package docs](https://pkg.go.dev/flag) — flag parsing patterns
- [Go os package docs](https://pkg.go.dev/os) — file I/O operations
- [Go encoding/json docs](https://pkg.go.dev/encoding/json) — JSON handling
- [Effective Go: Errors](https://go.dev/doc/effective_go#errors) — error handling patterns
- [Go testing package docs](https://pkg.go.dev/testing) — testing patterns
<!-- GSD:stack-end -->

<!-- GSD:conventions-start source:CONVENTIONS.md -->
## Conventions

Conventions not yet established. Will populate as patterns emerge during development.
<!-- GSD:conventions-end -->

<!-- GSD:architecture-start source:ARCHITECTURE.md -->
## Architecture

Architecture not yet mapped. Follow existing patterns found in the codebase.
<!-- GSD:architecture-end -->

<!-- GSD:skills-start source:skills/ -->
## Project Skills

No project skills found. Add skills to any of: `.claude/skills/`, `.agents/skills/`, `.cursor/skills/`, or `.github/skills/` with a `SKILL.md` index file.
<!-- GSD:skills-end -->

<!-- GSD:workflow-start source:GSD defaults -->
## GSD Workflow Enforcement

Before using Edit, Write, or other file-changing tools, start work through a GSD command so planning artifacts and execution context stay in sync.

Use these entry points:
- `/gsd-quick` for small fixes, doc updates, and ad-hoc tasks
- `/gsd-debug` for investigation and bug fixing
- `/gsd-execute-phase` for planned phase work

Do not make direct repo edits outside a GSD workflow unless the user explicitly asks to bypass it.
<!-- GSD:workflow-end -->



<!-- GSD:profile-start -->
## Developer Profile

> Profile not yet configured. Run `/gsd-profile-user` to generate your developer profile.
> This section is managed by `generate-claude-profile` -- do not edit manually.
<!-- GSD:profile-end -->
