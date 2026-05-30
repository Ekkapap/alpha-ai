---
name: golang-samber-ro
description: "Reactive streams and event-driven programming in Golang using samber/ro — ReactiveX implementation with 150+ type-safe operators, cold/hot observables, 5 subject types (Publish, Behavior, Replay, Async, Unicast), declarative pipelines via Pipe, 40+ plugins (HTTP, cron, fsnotify, JSON, logging), automatic backpressure, error propagation, and Go context integration. Apply when using or adopting samber/ro, when the codebase imports github.com/samber/ro, or when building asynchronous event-driven pipelines, real-time data processing, streams, or reactive architectures in Go. Not for finite slice transforms (-> See golang-samber-lo skill)."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.0.3"
  openclaw:
    emoji: "👁️"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent mcp__context7__resolve-library-id mcp__context7__query-docs AskUserQuestion
---

samber/ro is a ReactiveX implementation for Go: create cold Observables (lazy, restart per subscriber) or hot Subjects (Publish, Behavior, Replay, Async, Unicast), compose with 150+ operators via `Pipe()`, and integrate Go context for cancellation. Automatic backpressure and error propagation with 40+ plugins for HTTP, cron, fsnotify, JSON, and logging sources.

## References

| File | Purpose |
|------|---------|
| references/subjects-guide.md | Hot observable subjects: Publish, Behavior, Replay, Async, Unicast |
| references/operators-guide.md | All 150+ operators with signatures and marble diagrams |
| references/pipeline-patterns.md | Declarative pipeline composition with Pipe() |
| references/backend-handlers.md | Integration with HTTP handlers and backend services |
| references/http-middlewares.md | HTTP middleware patterns |
| references/sampling-strategies.md | Backpressure, throttle, debounce, sample strategies |
