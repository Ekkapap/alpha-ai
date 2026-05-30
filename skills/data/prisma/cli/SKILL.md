---
name: prisma-cli
description: Prisma CLI commands reference covering all available commands, options, and usage patterns. Use when running Prisma CLI commands, setting up projects, generating client, running migrations, or managing databases. Triggers on "prisma init", "prisma generate", "prisma migrate", "prisma db", "prisma studio".
license: MIT
metadata:
  author: prisma
  version: "7.0.0"
---

This skill covers the full Prisma CLI surface for v7: project initialization, client generation (with `--watch`/`--bun` flags), the local `prisma dev` database, database operations (`pull`, `push`, `seed`, `execute`), migration commands for dev and production, utility commands (`studio`, `validate`, `format`, `debug`), and v7 breaking changes including the new `prisma.config.ts` and removed flags.

## References

| File | Purpose |
|------|---------|
| references/init.md | Project initialization |
| references/generate.md | Client generation |
| references/dev.md | Local development database |
| references/db-pull.md | Database introspection |
| references/db-push.md | Schema push |
| references/db-seed.md | Database seeding |
| references/db-execute.md | Raw SQL execution |
| references/migrate-dev.md | Development migrations |
| references/migrate-deploy.md | Production migrations |
| references/migrate-reset.md | Database reset |
| references/migrate-status.md | Migration status |
| references/migrate-resolve.md | Migration resolution |
| references/migrate-diff.md | Schema diffing |
| references/studio.md | Database GUI |
| references/validate.md | Schema validation |
| references/format.md | Schema formatting |
| references/debug.md | Debug info |
