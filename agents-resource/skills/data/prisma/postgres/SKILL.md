---
name: prisma-postgres
description: Prisma Postgres setup and operations guidance across Console, create-db CLI, Management API, and Management API SDK. Use when creating Prisma Postgres databases, working in Prisma Console, provisioning with create-db/create-pg/create-postgres, or integrating programmatic provisioning with service tokens or OAuth.
license: MIT
metadata:
  author: prisma
  version: "1.0.0"
---

Prisma Postgres is a managed PostgreSQL service integrated with the Prisma ecosystem, accessible via the Console UI, the `create-db` CLI for instant provisioning, the Management API for server-to-server control, and the `@prisma/management-api-sdk` for type-safe TypeScript integration. Temporary databases auto-delete after ~24 hours unless claimed via a claim URL.

## References

| File | Purpose |
|------|---------|
| references/create-db-cli.md | Instant database provisioning via CLI and programmatic API |
| references/management-api.md | Service token and OAuth API workflows |
| references/management-api-sdk.md | Typed SDK usage with token storage |
| references/console-and-connections.md | Console operations and connection setup |
