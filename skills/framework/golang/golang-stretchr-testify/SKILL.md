---
name: golang-stretchr-testify
description: "Comprehensive guide to stretchr/testify for Golang testing. Covers assert, require, mock, and suite packages in depth. Use whenever writing tests with testify, creating mocks, setting up test suites, or choosing between assert and require. Essential for testify assertions, mock expectations, argument matchers, call verification, suite lifecycle, and advanced patterns like Eventually, JSONEq, and custom matchers. Trigger on any Go test file importing testify."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.3"
  openclaw:
    emoji: "✅"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - gotests
    install:
      - kind: go
        package: github.com/cweill/gotests/...@latest
        bins: [gotests]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch mcp__context7__resolve-library-id mcp__context7__query-docs
---

stretchr/testify: use `assert` when tests should continue on failure, `require` when remaining assertions are meaningless after failure. Create mocks with `mock.Mock` embedded in a struct, set expectations with `On().Return()`, verify with `AssertExpectations(t)`. Use `suite.Suite` for shared setup/teardown across multiple tests.

## References

| File | Purpose |
|------|---------|
| references/mock.md | Mock creation, expectations, argument matchers, call verification |
| references/helpers.md | Helper functions: Eventually, JSONEq, YAMLEq, custom matchers |
| references/http-testing.md | HTTP handler testing with testify |
| references/integration-testing.md | Integration test patterns and fixtures |
| references/mocking.md | Interface mocking strategies and patterns |
