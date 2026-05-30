---
name: golang-concurrency
description: "Golang concurrency patterns. Use when writing or reviewing concurrent Go code involving goroutines, channels, select, locks, sync primitives, errgroup, singleflight, worker pools, or fan-out/fan-in pipelines. Also triggers when you detect goroutine leaks, race conditions, channel ownership issues, or need to choose between channels and mutexes."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "⚡"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent AskUserQuestion
---

Go concurrency: every goroutine must have a clear exit mechanism (context, done channel, WaitGroup); only the sender closes a channel; always include `ctx.Done()` in select; use `errgroup` over `sync.WaitGroup` when errors matter; use `errgroup.SetLimit(n)` for bounded workers; detect leaks in tests with `goleak`. Use channels for ownership transfer, mutexes for shared struct state, atomics for simple counters.

## References

| File | Purpose |
|------|---------|
| references/channels-and-select.md | Channel/select patterns, direction types, buffering, ownership |
| references/sync-primitives.md | Mutex, RWMutex, atomic, sync.Pool, sync.Once, WaitGroup, errgroup, singleflight |
| references/pipelines.md | Fan-out/fan-in, bounded workers, generator chains, Go 1.23+ iterators, samber/ro |
