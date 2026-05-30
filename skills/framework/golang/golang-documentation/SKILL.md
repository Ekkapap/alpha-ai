---
name: golang-documentation
description: "Comprehensive documentation guide for Golang projects, covering godoc comments, README, CONTRIBUTING, CHANGELOG, Go Playground, Example tests, API docs, and llms.txt. Use when writing or reviewing doc comments, documentation, adding code examples, setting up doc sites, or discussing documentation best practices. Triggers for both libraries and applications/CLIs."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "📝"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch
---

Go documentation is a first-class deliverable: write godoc comments for all exported symbols (package, types, functions, constants), include runnable Example functions, maintain README with installation/usage/API overview, CONTRIBUTING with development setup, and CHANGELOG following Keep a Changelog format. Libraries also need `llms.txt` for AI-assisted usage.

## References

| File | Purpose |
|------|---------|
| references/code-comments.md | godoc comment style, exported symbol requirements, Example functions |
| references/application.md | Application/CLI documentation: README, man pages, help text |
| references/library.md | Library documentation: README, llms.txt, API docs, Go Playground |
| references/project-docs.md | CONTRIBUTING, CHANGELOG, ADRs, project-level docs |
