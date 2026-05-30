---
name: golang-testing
description: "Provides a comprehensive guide for writing production-ready Golang tests. Covers table-driven tests, test suites with testify, mocks, unit tests, integration tests, benchmarks, code coverage, parallel tests, fuzzing, fixtures, goroutine leak detection with goleak, snapshot testing, memory leaks, CI with GitHub Actions, and idiomatic naming conventions. Use this whenever writing tests, asking about testing patterns or setting up CI for Go projects. Essential for ANY test-related conversation in Go."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "🧪"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - gotests
    install:
      - kind: go
        package: github.com/cweill/gotests/gotests@latest
        bins: [gotests]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent Bash(gotests:*) AskUserQuestion
---

Go tests are executable specifications: write table-driven tests for parametric cases, use `t.Parallel()` for independent tests, mock at interface boundaries, always run with `-race` in CI, use `goleak` for goroutine leak detection, and write integration tests with `//go:build integration` build tags. Tests MUST be in the same package (white-box) or `_test` package (black-box).

## References

| File | Purpose |
|------|---------|
| references/mocking.md | Interface mock strategies, testify/mock, hand-rolled stubs |
| references/http-testing.md | HTTP handler testing with httptest, recording, assertions |
| references/integration-testing.md | Integration test setup, build tags, test containers |
| references/helpers.md | Helper utilities: t.Helper, testdata, golden files, fixtures |
