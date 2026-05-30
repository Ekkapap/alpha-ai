---
name: golang-design-patterns
description: "Idiomatic Golang design patterns — functional options, constructors, error flow and cascading, resource management and lifecycle, graceful shutdown, resilience, architecture, dependency injection, data handling, and streaming. Apply when designing Go APIs, structuring applications, choosing between patterns, making design decisions, architectural choices, or production hardening."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "🏗️"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent AskUserQuestion
---

Go design patterns favor simplicity and explicitness: use functional options for optional config, constructor functions (`NewXxx`) for initialization, explicit resource cleanup with `defer`, context for cancellation and timeouts, and graceful shutdown with signal handling. Apply the smallest pattern that solves the problem — favor clean architecture or hexagonal only when the domain complexity justifies it.

## References

| File | Purpose |
|------|---------|
| references/architecture.md | Architecture patterns: functional options, constructors, graceful shutdown, resilience |
| references/clean-architecture.md | Clean architecture with layered structure |
| references/hexagonal-architecture.md | Hexagonal (ports & adapters) architecture |
| references/ddd.md | Domain-Driven Design patterns in Go |
| references/data-handling.md | Data handling, streaming, transformation patterns |
| references/resource-management.md | Resource lifecycle, cleanup, connection management |
