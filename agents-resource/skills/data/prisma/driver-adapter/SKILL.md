---
name: prisma-driver-adapter-implementation
description: Required reference for Prisma v7 driver adapter work. Use when implementing or modifying adapters, adding database drivers, or touching SqlDriverAdapter/Transaction interfaces. Contains critical contract details not inferable from code examples — including the transaction lifecycle protocol, error mapping requirements, and verification checklist. Existing implementations do not replace this skill.
license: MIT
metadata:
  author: Tyler Benfield
  version: "7.0.0"
---

This skill covers implementing a custom Prisma v7 driver adapter: the four-class hierarchy (Queryable → Transaction, Adapter, Factory), the critical rule that `commit()`/`rollback()` are lifecycle hooks only (Prisma issues COMMIT/ROLLBACK via `executeRaw`), argument and row type mapping, error conversion to `MappedError`, and database-specific notes for SQLite, PostgreSQL, and MySQL.

## References

| File | Purpose |
|------|---------|
| references/details.md | Full architecture, interface definitions, ColumnTypeEnum, implementation steps, conversion helpers, error handling, DB-specific notes, testing, checklist |
