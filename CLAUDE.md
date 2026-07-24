# CLAUDE.md - Development Guide for Agentic-Init (ainit)

## CLI & TUI Context
- Binary Name: `ainit`
- Program Name: Agentic-Init
- Language: Go (Golang) 1.21+
- Frameworks: Bubble Tea (`github.com/charmbracelet/bubbletea`), Lipgloss, Bubbles

## Build & Run Commands
- Install Dependencies: `go mod tidy`
- Build Binary: `go build -o bin/ainit ./cmd/ainit`
- Run Binary: `./bin/ainit`
- Run from Source: `go run ./cmd/ainit`
- Run Tests: `go test -v ./...`

## Coding & Commit Standards
- Follow TDD First (write failing test spec, then implement, then pass).
- Atomic Micro-Commits with Conventional Commits format (`feat:`, `fix:`, `docs:`).
- Keep dependencies minimal and idiomatic to Go.
