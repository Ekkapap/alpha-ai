---
name: golang-samber-mo
description: "Monadic types for Golang using samber/mo — Option, Result, Either, Future, IO, Task, and State types for type-safe nullable values, error handling, and functional composition with pipeline sub-packages. Apply when using or adopting samber/mo, when the codebase imports `github.com/samber/mo`, or when considering functional programming patterns as a safety design for Golang."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.0.3"
  openclaw:
    emoji: "🎭"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch mcp__context7__resolve-library-id mcp__context7__query-docs AskUserQuestion
---

samber/mo provides monadic types for Go: `Option[T]` for nullable values without nil panics, `Result[T]` for explicit error handling, `Either[L, R]` for two-outcome types, `Future[T]` for async operations, and pipeline sub-packages for functional composition. Use to make impossible states unrepresentable and errors composable.

## References

| File | Purpose |
|------|---------|
| references/option.md | Option[T] API: Some, None, IsPresent, Get, OrElse, Map |
| references/result.md | Result[T] API: Ok, Err, IsOk, Unwrap, MapErr, FlatMap |
| references/either.md | Either[L, R] API: Left, Right, IsLeft, Fold |
| references/monads-guide.md | When to use each monad type, composition patterns |
| references/pipelines.md | Pipeline sub-packages for chaining operations |
| references/advanced-types.md | Future, IO, Task, State types |
