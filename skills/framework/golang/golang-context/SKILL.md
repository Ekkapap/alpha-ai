---
name: golang-context
description: "Idiomatic context.Context usage in Golang — creation, propagation, cancellation, timeouts, deadlines, context values, and cross-service tracing. Apply when working with context.Context in any Go code."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.1"
  openclaw:
    emoji: "🔗"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

`context.Context` is Go's mechanism for propagating cancellation, deadlines, and request-scoped values across API boundaries. Always pass as the first parameter named `ctx`, never store in a struct, never pass `nil` (use `context.TODO()`), always `defer cancel()` immediately after `WithCancel`/`WithTimeout`, and always propagate the same context through the entire call chain.

## References

| File | Purpose |
|------|---------|
| references/cancellation.md | WithCancel, WithTimeout, WithDeadline, AfterFunc, WithoutCancel |
| references/values-tracing.md | Safe context value patterns, trace context propagation, correlation IDs |
| references/http-services.md | HTTP handler context, middleware, client patterns, database *Context methods |
