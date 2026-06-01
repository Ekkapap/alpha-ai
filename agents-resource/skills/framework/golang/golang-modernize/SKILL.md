---
name: golang-modernize
description: "Continuously modernize Golang code to use the latest language features, standard library improvements, and idiomatic patterns. Use this skill whenever writing, reviewing, or refactoring Go code to ensure it leverages modern Go idioms. Also use when the user asks about Go upgrades, migration, modernization, deprecation, or when modernize linter reports issues. Also covers tooling modernization: linters, SAST, AI-powered code review in CI, and modern development practices. Trigger this skill proactively when you notice old-style Go patterns that have modern replacements."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.3"
  openclaw:
    emoji: "🔄"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch WebSearch AskUserQuestion
---

Go modernization: replace deprecated APIs and patterns with modern equivalents across language features (range-over-int, min/max, any, Go 1.23+ iterators), standard library (slices, maps, cmp, slog, t.Context, b.Loop), and tooling (golangci-lint v2, govulncheck, PGO). Prioritize correctness/safety fixes first, readability second, gradual improvements third.

## References

| File | Purpose |
|------|---------|
| references/versions.md | Per-version modernization guide: what changed in Go 1.18–1.24+, migration patterns |
| references/tooling.md | Tooling modernization: golangci-lint v2, govulncheck, PGO, CI pipeline |
