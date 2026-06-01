---
name: prisma-database-setup
description: Guides for configuring Prisma with different database providers (PostgreSQL, MySQL, SQLite, MongoDB, etc.). Use when setting up a new project, changing databases, or troubleshooting connection issues. Triggers on "configure postgres", "connect to mysql", "setup mongodb", "sqlite setup".
license: MIT
metadata:
  author: prisma
  version: "1.0.0"
---

This skill covers Prisma v7 database setup: supported providers and their driver adapters, the two-file configuration model (`schema.prisma` + `prisma.config.ts`), required driver adapter installation and instantiation for each database, and Prisma Client generation with mandatory explicit `output` path. MongoDB is not supported in v7.

## References

| File | Purpose |
|------|---------|
| references/postgresql.md | PostgreSQL setup |
| references/mysql.md | MySQL/MariaDB setup |
| references/sqlite.md | SQLite setup |
| references/mongodb.md | MongoDB (v6 only) |
| references/sqlserver.md | SQL Server setup |
| references/cockroachdb.md | CockroachDB setup |
| references/prisma-postgres.md | Prisma Postgres managed setup |
| references/prisma-client-setup.md | Client generation and adapter wiring |
