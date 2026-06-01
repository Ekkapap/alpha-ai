---
name: prisma-orm-v7-skills
description: Key facts and breaking changes for upgrading to Prisma ORM 7. Consider version 7 changes before generation or troubleshooting
---

Prisma ORM v7 is ESM-only, requires driver adapters for all databases, replaces env-based datasource config with `prisma.config.ts`, and removes automatic seeding and client generation from migrate commands. MongoDB is not supported in v7; use v6 for MongoDB projects.

## References

| File | Purpose |
|------|---------|
| references/details.md | Full breaking changes, upgrade checklist, removed env vars, ESM setup |

