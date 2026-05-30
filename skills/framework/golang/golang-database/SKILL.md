---
name: golang-database
description: "Comprehensive guide for Go database access. Covers parameterized queries, struct scanning, NULLable column handling, error patterns, transactions, isolation levels, SELECT FOR UPDATE, connection pool, batch processing, context propagation, and migration tooling. Use this skill whenever writing, reviewing, or debugging Golang code that interacts with PostgreSQL, MariaDB, MySQL, or SQLite. Also triggers for database testing or any question about database/sql, sqlx, pgx, or SQL queries in Golang. This skill explicitly does NOT generate database schemas or migration SQL."
user-invocable: false
license: MIT
compatibility: Designed for Claude Code or similar AI coding agents, and for projects using Golang.
metadata:
  author: samber
  version: "1.1.2"
  openclaw:
    emoji: "🗄️"
    homepage: https://github.com/samber/cc-skills-golang
    requires:
      bins:
        - go
    install: []
allowed-tools: Read Edit Write Glob Grep Bash(go:*) Bash(golangci-lint:*) Bash(git:*) Agent AskUserQuestion
---

Go database access uses `sqlx` or `pgx` over raw `database/sql` — never ORMs. Always use parameterized queries, pass context to all operations, handle `sql.ErrNoRows` explicitly, `defer rows.Close()` immediately, configure connection pool settings, wrap multi-statement writes in transactions, and use external tools (golang-migrate) for migrations. Never generate schemas or migrations.

## References

| File | Purpose |
|------|---------|
| references/scanning.md | Struct scanning with sqlx/pgx, NULLable column handling |
| references/transactions.md | Transaction patterns, isolation levels, SELECT FOR UPDATE |
| references/performance.md | Connection pool config, batch operations, query optimization |
| references/testing.md | Database testing patterns, test containers, mocking |
