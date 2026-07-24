# CLAUDE.md - Development & Agent Context Guide

> ⚠️ **CRITICAL AGENT INSTRUCTION**: This repository uses [`AGENTS.md`](./AGENTS.md) as the Single Source of Truth for all coding, commit, TDD, and architectural conventions.
> Claude Code MUST strictly read, respect, and obey all instructions in [`AGENTS.md`](./AGENTS.md).

## Quick Reference
- Primary Agent Rules: See [`AGENTS.md`](./AGENTS.md)
- Build: `make build` (`go build -o bin/ainit ./cmd/ainit`)
- Run TUI: `make run` (`./bin/ainit`)
- Test: `make test` (`go test -v ./...`)
- Lint & Tidy: `make fmt && make tidy`

## Agent Execution Guidelines
1. Always write failing test specs prior to writing feature code (TDD First).
2. Execute local sandbox tests (`make test`) before performing git commits.
3. Every commit MUST strictly follow Conventional Commits and atomic micro-commit principles.
