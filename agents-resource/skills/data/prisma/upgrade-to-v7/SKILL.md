---
name: prisma-upgrade-v7
description: Complete migration guide from Prisma ORM v6 to v7 covering all breaking changes. Use when upgrading Prisma versions, encountering v7 errors, or migrating existing projects. Triggers on "upgrade to prisma 7", "prisma 7 migration", "prisma-client generator", "driver adapter required".
license: MIT
metadata:
  author: prisma
  version: "7.0.0"
---

This skill covers the full Prisma v6→v7 migration: ESM-only module format, generator provider rename (`prisma-client-js` → `prisma-client`) with mandatory `output`, required driver adapter installation and client instantiation, `prisma.config.ts` replacing env-based datasource config, explicit dotenv loading, removal of middleware/metrics, and CLI flag changes.

## References

| File | Purpose |
|------|---------|
| references/esm-support.md | ESM module configuration |
| references/schema-changes.md | Generator and schema updates |
| references/driver-adapters.md | Required driver adapter setup |
| references/prisma-config.md | New configuration file |
| references/env-variables.md | Environment variable loading |
| references/removed-features.md | Middleware, metrics, and CLI flags |
| references/accelerate-users.md | Special handling for Accelerate |
