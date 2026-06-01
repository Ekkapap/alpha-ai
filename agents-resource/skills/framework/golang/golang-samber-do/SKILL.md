---
name: golang-samber-do
description: "Implements dependency injection in Golang using samber/do. Apply this skill when working with dependency injection, setting up service containers, managing service lifecycles, or when you see code using github.com/samber/do/v2. Also use when refactoring manual dependency injection, implementing health checks, graceful shutdown, or organizing services into scopes/modules."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.3"
  openclaw:
    emoji: "💉"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch mcp__context7__resolve-library-id mcp__context7__query-docs
---

samber/do v2 is a type-safe dependency injection container for Go: register providers at the composition root, resolve via `do.MustInvoke[T](injector)`, support lazy initialization and singletons, implement `do.Healthcheckable` and `do.Shutdownable` interfaces for production lifecycle management. Keep the injector only at `main()`.

## References

| File | Purpose |
|------|---------|
| references/advanced.md | Advanced patterns: scopes, named services, modules, provider chaining |
| references/testing.md | Testing with samber/do: overriding providers, mock injection |
