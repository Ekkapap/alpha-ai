---
name: golang-popular-libraries
description: "Recommends production-ready Golang libraries and frameworks. Apply when the user asks for library suggestions, wants to compare alternatives, or needs to choose a library for a specific task. Also apply when the AI agent is about to add a new dependency — ensures vetted, production-ready libraries are chosen."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.4"
  openclaw:
    emoji: "📚"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch WebSearch AskUserQuestion
---

Recommend production-ready Go libraries by checking the standard library first, then vetted third-party options. Always check maintenance status and license before recommending. Avoid libraries that wrap stdlib without adding value or have large dependency footprints for simple needs. Always ask the user before adding any new dependency (including via `go get`).

## References

| File | Purpose |
|------|---------|
| references/stdlib.md | Standard library v2 packages, promoted x/exp packages, golang.org/x extensions |
| references/libraries.md | Vetted third-party libraries by category (web, database, testing, logging, messaging) |
| references/tools.md | Development tools: debugging, linting, testing, dependency management |
