---
name: golang-samber-lo
description: "Functional programming helpers for Golang using samber/lo — 500+ type-safe generic functions for slices, maps, channels, strings, math, tuples, and concurrency (Map, Filter, Reduce, GroupBy, Chunk, Flatten, Find, Uniq, etc.). Core immutable package (lo), concurrent variants (lo/parallel aka lop), in-place mutations (lo/mutable aka lom), lazy iterators (lo/it aka loi for Go 1.23+), and experimental SIMD (lo/exp/simd). Apply when using or adopting samber/lo, when the codebase imports github.com/samber/lo, or when implementing functional-style data transformations in Go. Not for streaming pipelines (→ See golang-samber-ro skill)."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.0.3"
  openclaw:
    emoji: "🧰"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) mcp__context7__resolve-library-id mcp__context7__query-docs AskUserQuestion
---

samber/lo provides 500+ type-safe generic collection functions for Go: `lo.Map`, `lo.Filter`, `lo.Reduce`, `lo.GroupBy`, `lo.Chunk`, `lo.Flatten`, `lo.Find`, `lo.Uniq`. Use `lop` for concurrent transforms, `lom` for in-place mutations, `loi` for lazy Go 1.23+ iterators. Prefer `lo` over manual for-loops for declarative, nil-safe collection transforms.

## References

| File | Purpose |
|------|---------|
| references/api-reference.md | Full API: all function signatures, type constraints, return types |
| references/package-guide.md | Package selection guide: lo vs lop vs lom vs loi vs lo/exp/simd |
| references/advanced-patterns.md | Advanced patterns: function composition, pipeline building, custom iterators |
