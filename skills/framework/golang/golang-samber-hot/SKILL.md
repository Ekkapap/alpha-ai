---
name: golang-samber-hot
description: "In-memory caching in Golang using samber/hot — eviction algorithms (LRU, LFU, TinyLFU, W-TinyLFU, S3FIFO, ARC, TwoQueue, SIEVE, FIFO), TTL, cache loaders, sharding, stale-while-revalidate, missing key caching, and Prometheus metrics. Apply when using or adopting samber/hot, when the codebase imports github.com/samber/hot, or when the project repeatedly loads the same medium-to-low cardinality resources at high frequency and needs to reduce latency or backend pressure."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.0.3"
  openclaw:
    emoji: "🔥"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch mcp__context7__resolve-library-id mcp__context7__query-docs AskUserQuestion
---

samber/hot is a type-safe in-memory caching library for Go 1.22+ with 9 eviction algorithms. Default to `hot.WTinyLFU` for general-purpose workloads; switch only when profiling shows high miss rates. Supports TTL, loader chains with singleflight deduplication, sharding for high-concurrency, stale-while-revalidate, and Prometheus metrics.

## References

| File | Purpose |
|------|---------|
| references/algorithm-guide.md | Algorithm comparison, benchmarks, decision tree |
| references/api-reference.md | Full API: cache creation, TTL, loaders, sharding, metrics |
| references/production-patterns.md | Production usage patterns: stale-while-revalidate, missing key caching, warm-up |
