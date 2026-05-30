---
name: golang-structs-interfaces
description: 'Golang struct and interface design patterns — composition, embedding, type assertions, type switches, interface segregation, dependency injection via interfaces, struct field tags, and pointer vs value receivers. Use this skill when designing Go types, defining or implementing interfaces, embedding structs or interfaces, writing type assertions or type switches, adding struct field tags for JSON/YAML/DB serialization, or choosing between pointer and value receivers. Also use when the user asks about "accept interfaces, return structs", compile-time interface checks, or composing small interfaces into larger ones.'
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.3"
  openclaw:
    emoji: "🧩"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent AskUserQuestion
---

Go type design: keep interfaces small (1-3 methods), define them where consumed (not where implemented), accept interfaces and return concrete structs, never create interfaces prematurely. Use compile-time checks (`var _ Interface = (*Type)(nil)`), safe comma-ok type assertions, consistent pointer or value receivers across all methods of a type.

## References

| File | Purpose |
|------|---------|
| references/details.md | Full rules: interface principles, zero value design, generics vs any, standard library interfaces, type assertions, embedding, field tags, noCopy, pointer vs value receivers, common mistakes |
