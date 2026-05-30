---
name: golang-benchmark
description: "Golang benchmarking, profiling, and performance measurement. Use when writing, running, or comparing Go benchmarks, profiling hot paths with pprof, interpreting CPU/memory/trace profiles, analyzing results with benchstat, setting up CI benchmark regression detection, or investigating production performance with Prometheus runtime metrics. Also use when the developer needs deep analysis on a specific performance indicator - this skill provides the measurement methodology, while golang-performance provides the optimization patterns."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "📊"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - benchstat
    install:
      - kind: go
        package: golang.org/x/perf/cmd/benchstat@latest
        bins: [benchstat]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch Bash(benchstat:*) Bash(benchdiff:*) Bash(cob:*) Bash(gobenchdata:*) Bash(curl:*) mcp__context7__resolve-library-id mcp__context7__query-docs WebSearch AskUserQuestion
---

Go benchmarking skill covering the full measurement workflow: write benchmarks with `b.Loop()` (Go 1.24+), run with `-benchmem -count=N`, profile with pprof (CPU, memory, goroutine), compare statistically with `benchstat`, detect regressions in CI with benchdiff/cob/gobenchdata, and investigate production performance with Prometheus runtime metrics. Never optimize without measuring first.

## References

| File | Purpose |
|------|---------|
| references/pprof.md | Interactive/non-interactive CPU, memory, goroutine profile analysis |
| references/benchstat.md | Statistical benchmark comparison with confidence intervals |
| references/trace.md | Execution tracer for goroutine scheduling and GC phases |
| references/tools.md | fieldalignment, GODEBUG, fgprof, race detector |
| references/compiler-analysis.md | Escape analysis, inlining, SSA, assembly output |
| references/ci-regression.md | CI benchmark regression detection (benchdiff, cob, gobenchdata) |
| references/investigation-session.md | Production performance troubleshooting with Prometheus |
| references/prometheus-go-metrics.md | Go runtime metrics exposed via prometheus/client_golang |
