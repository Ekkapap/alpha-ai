---
name: golang-observability
description: "Golang everyday observability — the always-on signals in production. Covers structured logging with slog, Prometheus metrics, OpenTelemetry distributed tracing, continuous profiling with pprof/Pyroscope, server-side RUM event tracking, alerting, and Grafana dashboards. Apply when instrumenting Go services for production monitoring, setting up metrics or alerting, adding OpenTelemetry tracing, correlating logs with traces, migrating legacy loggers (zap/logrus/zerolog) to slog, adding observability to new features, or implementing GDPR/CCPA-compliant tracking with Customer Data Platforms (CDP). Not for temporary deep-dive performance investigation (→ See golang-benchmark and golang-performance skills)."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.3"
  openclaw:
    emoji: "📡"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch WebSearch AskUserQuestion
---

Go observability: use `slog` for structured logging with context-propagated trace IDs, Prometheus for metrics (counter/gauge/histogram), OpenTelemetry for distributed tracing with automatic context propagation, and Pyroscope for continuous profiling. Define SLOs before writing alert rules; correlate logs with traces using trace ID fields. Migrate legacy loggers (zap/logrus/zerolog) to slog with samber/slog-* adapters.

## References

| File | Purpose |
|------|---------|
| references/logging.md | slog setup, log levels, structured attributes, samber/slog-* plugins, migration from zap/logrus |
| references/tracing.md | OpenTelemetry setup, span creation, context propagation, sampling |
| references/metrics.md | Prometheus client_golang, counter/gauge/histogram patterns, /metrics endpoint |
| references/dashboards.md | Grafana dashboard design for Go services |
| references/alerting.md | Prometheus alerting rules, SLO-based alerts |
| references/profiling.md | Pyroscope continuous profiling, pprof endpoint setup |
| references/rum.md | Server-side RUM event tracking, CDP integration |
