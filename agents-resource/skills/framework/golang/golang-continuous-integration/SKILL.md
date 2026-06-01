---
name: golang-continuous-integration
description: "Provides CI/CD pipeline configuration using GitHub Actions for Golang projects. Covers testing, linting, SAST, security scanning, code coverage, Dependabot, Renovate, GoReleaser, code review automation, and release pipelines. Use this whenever setting up CI for a Go project, configuring workflows, adding linters or security scanners, setting up Dependabot or Renovate, automating releases, or improving an existing CI pipeline. Also use when the user wants to add quality gates to their Go project."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "🚀"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - goreleaser
        - gh
    install:
      - kind: brew
        formula: goreleaser
        bins: [goreleaser]
      - kind: brew
        formula: gh
        bins: [gh]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch Bash(goreleaser:*) Bash(gh:*) AskUserQuestion
---

Production-grade CI/CD for Go projects via GitHub Actions: always run `go test -race -shuffle=on` in a Go version matrix, lint with golangci-lint, scan with govulncheck and gosec/CodeQL, automate dependency updates with Dependabot/Renovate, and release binaries with GoReleaser. Always search for the latest action versions before generating workflow YAML. Always add `permissions` blocks per job for least privilege.

## References

| File | Purpose |
|------|---------|
| assets/test.yml | Testing workflow with Go version matrix, race detection, coverage |
| assets/lint.yml | golangci-lint workflow |
| assets/integration.yml | Integration test workflow |
| assets/docker.yml | Multi-platform Docker image build and push |
| assets/dependabot.yml | Dependabot configuration |
| assets/dependabot-auto-merge.yml | Auto-merge for Dependabot PRs |
| assets/goreleaser-cli.yml | GoReleaser release pipeline for CLI tools |
| assets/goreleaser-lib.yml | GoReleaser release pipeline for libraries |
| assets/goreleaser-monorepo.yml | GoReleaser monorepo release |
| assets/codeql-config.yml | CodeQL configuration |
| assets/codecov.yml | Codecov configuration |
