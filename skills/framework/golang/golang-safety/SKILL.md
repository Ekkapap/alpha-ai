---
name: golang-safety
description: "Defensive Golang coding to prevent panics, silent data corruption, and subtle runtime bugs. Use whenever writing or reviewing Go code that involves nil-prone types (pointers, interfaces, maps, slices, channels), numeric conversions, resource lifecycle (defer in loops), or defensive copying. Also triggers on questions about nil panics, append aliasing, map concurrent access, float comparison, or zero-value design."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.1"
  openclaw:
    emoji: "🛡️"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

Defensive Go: guard against nil interface trap (non-nil interface holding nil pointer), slice aliasing from shared backing arrays after append, concurrent map read/write hard crash, defer in loops accumulating until function return, and integer overflow in numeric conversions. Use `go test -race` in CI to catch data races automatically.

## References

| File | Purpose |
|------|---------|
| references/nil-safety.md | Nil pointer patterns, nil interface trap, nil error comparison |
| references/slice-map-safety.md | Slice aliasing/defensive copies, concurrent map access, append traps |
