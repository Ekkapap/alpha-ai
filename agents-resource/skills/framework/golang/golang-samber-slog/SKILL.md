---
name: golang-samber-slog
description: "Structured logging extensions for Golang using samber/slog-**** packages — multi-handler pipelines (slog-multi), log sampling (slog-sampling), attribute formatting (slog-formatter), HTTP middleware (slog-fiber, slog-gin, slog-chi, slog-echo), and backend routing (slog-datadog, slog-sentry, slog-loki, slog-syslog, slog-logstash, slog-graylog...). Apply when using or adopting slog, or when the codebase already imports any github.com/samber/slog-* package."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.0.3"
  openclaw:
    emoji: "🪵"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch mcp__context7__resolve-library-id mcp__context7__query-docs AskUserQuestion
---

samber/slog-* packages extend Go's built-in `slog` with multi-handler fanout (slog-multi), log sampling to drop noise (slog-sampling), attribute formatting/PII stripping (slog-formatter), HTTP request middleware for Fiber/Gin/Chi/Echo (slog-*), and backend routing to Datadog/Sentry/Loki/Logstash/Graylog. Compose handlers in a pipeline: sample → format → route.

## References

| File | Purpose |
|------|---------|
| references/backend-handlers.md | Backend handler packages: slog-datadog, slog-sentry, slog-loki, slog-graylog |
| references/http-middlewares.md | HTTP middleware packages: slog-fiber, slog-gin, slog-chi, slog-echo |
| references/pipeline-patterns.md | Multi-handler pipeline composition with slog-multi |
| references/sampling-strategies.md | Log sampling strategies with slog-sampling |
