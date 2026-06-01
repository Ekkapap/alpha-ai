---
name: golang-security
description: "Security best practices and vulnerability prevention for Golang. Covers injection (SQL, command, XSS), cryptography, filesystem safety, network security, cookies, secrets management, memory safety, and logging. Apply when writing, reviewing, or auditing Go code for security, or when working on any risky code involving crypto, I/O, secrets management, user input handling, or authentication. Includes configuration of security tools."
user-invocable: true
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.3"
  openclaw:
    emoji: "🔒"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
        - govulncheck
    install:
      - kind: go
        package: golang.org/x/vuln/cmd/govulncheck@latest
        bins: [govulncheck]
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent WebFetch WebSearch AskUserQuestion
---

Go security: always parameterize SQL and shell commands (never concatenate user input), use `crypto/rand` not `math/rand`, validate and sanitize all file paths, rotate secrets via environment variables never hardcoded, use `SecureCompare` for token comparison, and run `govulncheck` and `gosec` in CI. Audit third-party dependencies and implement threat modeling for sensitive flows.

## References

| File | Purpose |
|------|---------|
| references/injection.md | SQL injection, command injection, XSS prevention |
| references/cryptography.md | Correct use of crypto/rand, AES-GCM, bcrypt, HMAC |
| references/filesystem.md | Path traversal prevention, safe file operations |
| references/network.md | TLS configuration, certificate validation, timeouts |
| references/cookies.md | Secure, HttpOnly, SameSite cookie settings |
| references/secrets.md | Secrets management, env vars, vault integration |
| references/memory-safety.md | Buffer handling, unsafe pointer safety |
| references/logging.md | Avoiding PII in logs, structured security events |
| references/third-party.md | Dependency auditing, supply chain security |
| references/threat-modeling.md | STRIDE methodology for Go services |
| references/architecture.md | Security architecture patterns |
| references/checklist.md | Pre-deployment security checklist |
