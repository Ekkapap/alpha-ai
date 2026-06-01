---
name: golang-performance
description: "Golang performance optimization patterns and methodology - if X bottleneck, then apply Y. Covers allocation reduction, CPU efficiency, memory layout, GC tuning, pooling, caching, and hot-path optimization. Use when profiling or benchmarks have identified a bottleneck and you need the right optimization pattern to fix it. Also use when performing performance code review to suggest improvements or benchmarks that could help identify quick performance gains. Not for measurement methodology (see golang-benchmark skill) or debugging workflow (see golang-troubleshooting skill)."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "🏎️"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - benchstat
    install:
      - kind: go
        package: golang.org/x/perf/cmd/benchstat@latest
        bins: [benchstat]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch Bash(benchstat:*) Bash(fieldalignment:*) Bash(staticcheck:*) Bash(curl:*) Bash(fgprof:*) Bash(perf:*) WebSearch AskUserQuestion
---

Go performance: profile before optimizing (pprof), reduce allocations first (biggest ROI), optimize one thing at a time with benchmark baseline → improve → compare with benchstat. Diagnose external bottlenecks (slow DB, API) before touching Go code. Use the decision tree: alloc_objects high → memory optimization, CPU-bound → CPU optimization, GC pauses → runtime tuning, blocked goroutines → I/O optimization.

## References

| File | Purpose |
|------|---------|
| references/memory.md | Allocation patterns, sync.Pool, struct alignment, backing array leaks |
| references/cpu.md | Inlining, cache locality, false sharing, reflection avoidance |
| references/io-networking.md | HTTP transport config, streaming, JSON performance, batch operations |
| references/runtime.md | GOGC, GOMEMLIMIT, GC diagnostics, GOMAXPROCS, PGO |
| references/caching.md | Algorithmic complexity, compiled patterns, singleflight, work avoidance |
| references/observability.md | Prometheus metrics, continuous profiling, alerting rules |
