---
name: golang-dependency-injection
description: "Golang dependency injection patterns — manual constructor injection, google/wire, uber-go/dig, uber-go/fx, and samber/do. Use when wiring services, choosing a DI approach, or testing with mocks. Triggers when code uses global variables for services, init() for setup, or when the project has 10+ interconnected services."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "💉"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent AskUserQuestion
---

Dependency injection in Go: inject via constructors (never globals or `init()`), define interfaces where consumed (not where implemented), keep the DI container only at `main()` composition root. Use manual injection for <10 services, a library (samber/do, uber-go/dig/fx, google/wire) for 10+ services. Deep dependency chains signal design problems.

## References

| File | Purpose |
|------|---------|
| references/manual-di.md | Manual constructor injection with complete application example |
| references/samber-do.md | samber/do container usage patterns |
| references/uber-dig-fx.md | uber-go/dig and uber-go/fx patterns |
| references/google-wire.md | google/wire compile-time DI |
