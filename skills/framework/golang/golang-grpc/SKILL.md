---
name: golang-grpc
description: "Provides gRPC usage guidelines, protobuf organization, and production-ready patterns for Golang microservices. Use when implementing, reviewing, or debugging gRPC servers/clients, writing proto files, setting up interceptors, handling gRPC errors with status codes, configuring TLS/mTLS, testing with bufconn, or working with streaming RPCs."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.3"
  openclaw:
    emoji: "🌐"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - protoc
    install:
      - kind: brew
        formula: protobuf
        bins: [protoc]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch mcp__context7__resolve-library-id mcp__context7__query-docs Bash(protoc:*) AskUserQuestion
---

Go gRPC services: organize proto files in `api/proto/`, use status codes (not raw errors), always set deadlines via context, chain interceptors for logging/auth/recovery, configure TLS/mTLS for production, test with `bufconn` in-process, handle graceful shutdown with `grpcServer.GracefulStop()`. Never swallow errors — wrap with `status.Errorf`.

## References

| File | Purpose |
|------|---------|
| references/protoc-reference.md | protoc command reference, plugin setup, buf toolchain |
| references/testing.md | bufconn-based in-process testing, mock patterns |
