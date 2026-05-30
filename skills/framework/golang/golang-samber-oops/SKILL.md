---
name: golang-samber-oops
description: "Structured error handling in Golang with samber/oops — error builders, stack traces, error codes, error context, error wrapping, error attributes, user-facing vs developer messages, panic recovery, and logger integration. Apply when using or adopting samber/oops, or when the codebase already imports github.com/samber/oops."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "💥"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch mcp__context7__resolve-library-id mcp__context7__query-docs
---

samber/oops provides a builder API for structured errors with stack traces, error codes, HTTP status, user-facing messages, and arbitrary attributes. Use `oops.Code("not_found").In("user_service").With("user_id", id).Errorf("user not found")` pattern. Integrates with slog, zerolog, and logrus via plugin packages.

## References

| File | Purpose |
|------|---------|
| references/advanced.md | Advanced patterns: error classification, logger integration, panic recovery |
| references/operators-guide.md | Full builder operator API: Code, In, Tags, With, Trace, User, Tenant, Request |
| references/patterns.md | Common usage patterns by scenario |
| references/plugin-ecosystem.md | Logger plugins: slog, zerolog, logrus integration |
