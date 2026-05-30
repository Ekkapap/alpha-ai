---
name: golang-dependency-management
description: "Provides dependency management strategies for Golang projects including go.mod management, installing/upgrading packages, semantic versioning, Minimal Version Selection, vulnerability scanning, outdated dependency tracking, dependency size analysis, automated updates with Dependabot/Renovate, conflict resolution, and dependency graph visualization. Use this skill whenever adding, removing, updating, or auditing Go dependencies, resolving version conflicts, setting up automated dependency updates, analyzing binary size, or working with go.work workspaces."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "📦"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - govulncheck
    install:
      - kind: go
        package: golang.org/x/vuln/cmd/govulncheck@latest
        bins: [govulncheck]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent Bash(govulncheck:*) AskUserQuestion
---

Go dependency management: always ask before adding a new dependency (AI agents MUST confirm with user before `go get`), prefer the standard library, check maintenance status and license before recommending external packages. Use `govulncheck` for vulnerability scanning, `go mod tidy` in CI (fail on diff), and Dependabot/Renovate for automated updates.

## References

| File | Purpose |
|------|---------|
| references/versioning.md | go.mod semantics, Minimal Version Selection, semantic versioning |
| references/auditing.md | Dependency auditing, govulncheck, license checking |
| references/automated-updates.md | Dependabot and Renovate configuration |
| references/conflicts.md | Version conflict resolution patterns |
| references/visualization.md | Dependency graph visualization and analysis |
| references/workspaces.md | go.work for multi-module development |
