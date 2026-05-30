---
name: golang-data-structures
description: "Golang data structures — slices (internals, capacity growth, preallocation, slices package), maps (internals, hash buckets, maps package), arrays, container/list/heap/ring, strings.Builder vs bytes.Buffer, generic collections, pointers (unsafe.Pointer, weak.Pointer), and copy semantics. Use when choosing or optimizing Go data structures, implementing generic containers, using container/ packages, unsafe or weak pointers, or questioning slice/map internals."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.1"
  openclaw:
    emoji: "🗃️"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

Go data structures: preallocate slices/maps with `make(T, 0, n)`, use `strings.Builder` for string building (not `bytes.Buffer`), prefer `slices`/`maps` packages (Go 1.21+) over manual iteration, use `container/heap` for priority queues, use `weak.Pointer[T]` (Go 1.24+) for caches, and follow the 6 valid `unsafe.Pointer` conversion patterns strictly.

## References

| File | Purpose |
|------|---------|
| references/slice-internals.md | Full slices package reference, growth mechanics, backing array aliasing |
| references/map-internals.md | Map hash table internals, maps package, concurrent access |
| references/generics.md | Generic collection patterns, type constraints |
| references/containers.md | container/list, container/heap, container/ring usage |
| references/pointers.md | unsafe.Pointer valid patterns, weak.Pointer, copy semantics |
