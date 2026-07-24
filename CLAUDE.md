# CLAUDE.md - Development Guide for Agentic-Init (ainit)

## CLI & TUI Context
- Binary Name: `ainit`
- Program Name: Agentic-Init
- Language: Go (Golang) 1.21+
- Frameworks: Bubble Tea (`github.com/charmbracelet/bubbletea`), Lipgloss, Bubbles

## Build & Run Commands (via Makefile)
- `make build`: Build `bin/ainit` binary
- `make run`: Build and execute TUI
- `make test`: Run all unit tests
- `make fmt`: Format Go code
- `make tidy`: Clean and sync `go.mod`
- `make install`: Install `ainit` to `$GOPATH/bin`
- `make clean`: Remove `bin/` directory

## Coding & Commit Standards
- Follow TDD First (write failing test spec, then implement, then pass).
- Atomic Micro-Commits with Conventional Commits format (`feat:`, `fix:`, `docs:`).
- Keep dependencies minimal and idiomatic to Go.
