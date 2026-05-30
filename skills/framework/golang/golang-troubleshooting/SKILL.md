---
name: golang-troubleshooting
description: "Troubleshoot Golang programs systematically - find and fix the root cause. Use when encountering bugs, crashes, deadlocks, or unexpected behavior in Go code. Covers debugging methodology, common Go pitfalls, test-driven debugging, pprof setup and capture, Delve debugger, race detection, GODEBUG tracing, and production debugging. Start here for any 'something is wrong' situation. Not for interpreting profiles or benchmarking (see golang-benchmark skill) or applying optimization patterns (see golang-performance skill)."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "🔍"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - dlv
    install:
      - kind: go
        package: github.com/go-delve/delve/cmd/dlv@latest
        bins: [dlv]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Bash(dlv:*) Agent WebFetch WebSearch AskUserQuestion
---

Systematic Go debugging: reproduce the bug with a failing test first, then diagnose. Use `-race` for concurrency bugs, `dlv` for interactive debugging, pprof for performance, `GODEBUG` flags for runtime tracing, and structured logs for production. Work from evidence — instrument, reproduce, trace root causes — never guess.

## References

| File | Purpose |
|------|---------|
| references/methodology.md | Systematic debugging workflow: reproduce → isolate → fix → verify |
| references/common-go-bugs.md | Common Go pitfalls and their root causes |
| references/compilation.md | Build errors, module issues, import cycles |
| references/concurrency-debug.md | Race conditions, deadlocks, goroutine leaks with dlv and goleak |
| references/diagnostic-tools.md | pprof setup, dlv commands, GODEBUG flags |
| references/testing-debug.md | Test failures, flaky tests, test isolation |
| references/performance-debug.md | Performance investigation workflow |
| references/production-debug.md | Production debugging: pprof endpoints, log correlation |
| references/pprof.md | pprof capture, analysis, interpretation |
| references/code-review-flags.md | Code review red flags to look for |
