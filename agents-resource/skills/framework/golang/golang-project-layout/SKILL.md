---
name: golang-project-layout
description: "Provides a guide for setting up Golang project layouts and workspaces. Use this whenever starting a new Go project, organizing an existing codebase, setting up a monorepo with multiple packages, creating CLI tools with multiple main packages, or deciding on directory structure. Apply this for any Go project initialization or restructuring work."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.3"
  openclaw:
    emoji: "📁"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent AskUserQuestion
---

Go project layout: right-size to the problem — a script stays flat, a service gets layers only when justified by complexity. Always ask the developer about their architecture preference (clean, hexagonal, DDD, flat) before proposing structure. Never over-structure small projects. For monorepos and multi-module workspaces, use `go.work`.

## References

| File | Purpose |
|------|---------|
| references/directory-layouts.md | Standard layouts for CLI, library, service, and monorepo projects |
| references/config.md | Configuration management patterns and file placement |
| references/testing-layout.md | Test file organization, testdata directories, integration test setup |
| references/workspaces.md | go.work setup for monorepos and multi-module development |
