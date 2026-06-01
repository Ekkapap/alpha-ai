---
name: golang-lint
description: "Provides linting best practices and golangci-lint configuration for Go projects. Covers running linters, configuring .golangci.yml, suppressing warnings with nolint directives, interpreting lint output, and managing linter settings. Use this skill whenever the user runs linters, configures golangci-lint, asks about lint warnings or suppressions, sets up code quality tooling, or asks which linters to enable for a Go project. Also use when the user mentions golangci-lint, go vet, staticcheck, revive, or any Go linting tool."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.1"
  openclaw:
    emoji: "🧹"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - golangci-lint
    install:
      - kind: brew
        formula: golangci-lint
        bins: [golangci-lint]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

`golangci-lint` is the standard Go linting tool aggregating 100+ linters. Every project must have `.golangci.yml` as the source of truth for enabled linters. Run `golangci-lint run --fix ./...` during development and always in CI; use `//nolint:lintername // reason` sparingly with mandatory justification comments.

## References

| File | Purpose |
|------|---------|
| references/linter-reference.md | Available linters, recommended set, configuration patterns |
| references/nolint-directives.md | nolint directive syntax, scope, justification requirements |
| assets/.golangci.yml | Production-ready configuration with 33 linters enabled |
