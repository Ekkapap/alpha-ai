---
name: golang-naming
description: "Go (Golang) naming conventions — covers packages, constructors, structs, interfaces, constants, enums, errors, booleans, receivers, getters/setters, functional options, acronyms, test functions, and subtest names. Use this skill when writing new Go code, reviewing or refactoring, choosing between naming alternatives (New vs NewTypeName, isConnected vs connected, ErrNotFound vs NotFoundError, StatusReady vs StatusUnknown at iota 0), debating Go package names (utils/helpers anti-patterns), or asking about Go naming best practices. Also trigger when the user mentions MixedCaps vs snake_case, ALL_CAPS constants, Get-prefix on getters, or error string casing. Do NOT use for general Go implementation questions that don't involve naming decisions."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.1"
  openclaw:
    emoji: "🏷️"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent
---

Go naming uses MixedCaps (never underscores or ALL_CAPS), with capitalization as the visibility control. Key rules: `New` for single-type constructors, acronyms fully capitalized (URL not Url), no `Get` prefix on getters, error strings lowercase with no punctuation, boolean names state the condition directly, avoid `utils`/`helpers`/`common` package names.

## References

| File | Purpose |
|------|---------|
| references/identifiers.md | Variables, constants, acronyms, receivers, booleans |
| references/packages-files.md | Package naming rules, file naming, anti-patterns (utils/helpers) |
| references/functions-methods.md | Constructor patterns, getter/setter rules, functional options |
| references/types-errors.md | Struct, interface, enum/iota, sentinel error, custom error type naming |
| references/testing.md | Test function names, subtest names, benchmark names |
