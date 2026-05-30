---
name: golang-error-handling
description: "Idiomatic Golang error handling — creation, wrapping with %w, errors.Is/As, errors.Join, custom error types, sentinel errors, panic/recover, the single handling rule, structured logging with slog, HTTP request logging middleware, and samber/oops for production errors. Built to make logs usable at scale with log aggregation 3rd-party tools. Apply when creating, wrapping, inspecting, or logging errors in Go code."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.1"
  openclaw:
    emoji: "⚠️"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

Go error handling: always check returned errors, wrap with `fmt.Errorf("{context}: %w", err)`, use `errors.Is`/`errors.As` for inspection, enforce the single-handling rule (log OR return, never both), keep error messages lowercase with no punctuation, keep messages low-cardinality (no interpolated IDs), and use `slog` for structured logging. Use `samber/oops` for production errors needing stack traces.

## References

| File | Purpose |
|------|---------|
| references/error-creation.md | Sentinel errors, custom error types, message style |
| references/error-wrapping.md | %w vs %v, errors.Is/As, errors.Join |
| references/error-handling.md | Single handling rule, panic/recover, slog, samber/oops |
