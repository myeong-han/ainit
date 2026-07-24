# CLAUDE.md - Development Guide for Agentic-Init (ainit)

## CLI & TUI Context
- Binary Name: `ainit`
- Program Name: Agentic-Init
- Project Goal: Harness engineering TUI tool for initialization, architecture design, and release.

## Build & Test Commands
- Build: `go build -o bin/ainit ./cmd/ainit` (or package manager equivalent)
- Test: `go test ./...`
- Lint: `golangci-lint run`

## Coding Standards
- Follow idiomatic design for TUI components.
- Keep dependencies minimal.
- Use explicit error handling and status indicators in TUI.
