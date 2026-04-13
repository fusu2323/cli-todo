# Stack Research

**Domain:** CLI Todo/Task Manager in Go
**Researched:** 2026-04-13
**Confidence:** HIGH

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

No external packages required. Go standard library only:

```bash
# No packages to install — everything is in stdlib
go mod init github.com/yourname/todo
```

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

**If simple single-command interface (add, list, done):**
- Use flat `flag` parsing in `main()`
- Single JSON file `~/.todo.json`
- Direct `os.ReadFile`/`os.WriteFile` pattern

**If subcommand interface (todo add, todo list, todo done):**
- Parse base command first, then subcommand flags
- Consider command pattern with separate handler functions
- May need `flag.FlagSet` for subcommand-level flag parsing

**If data grows beyond few hundred tasks:**
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

---
*Stack research for: CLI Todo/Task Manager in Go*
*Researched: 2026-04-13*
