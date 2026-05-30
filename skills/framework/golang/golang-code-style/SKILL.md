---
name: golang-code-style
description: "Golang code style, formatting and conventions. Use when writing code, reviewing style, configuring linters, writing comments, or establishing project standards."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.1"
  openclaw:
    emoji: "🎨"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

Go code style focuses on clarity over cleverness: use early returns to reduce nesting, `switch` over if-else chains, `:=` for non-zero values and `var` for zero-value initialization, composite literals with field names, functions with ≤4 parameters, and `range` for iteration. Lines beyond ~120 chars must be broken at semantic boundaries. Many rules are enforced automatically by `gofmt`, `gofumpt`, `goimports`, `gocritic`, `revive`.

## References

| File | Purpose |
|------|---------|
| references/details.md | Additional rules: complex conditions, init scope, value vs pointer arguments, code organization order |
