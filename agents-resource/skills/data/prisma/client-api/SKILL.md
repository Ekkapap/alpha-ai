---
name: prisma-client-api
description: Prisma Client API reference covering model queries, filters, operators, and client methods. Use when writing database queries, using CRUD operations, filtering data, or configuring Prisma Client. Triggers on "prisma query", "findMany", "create", "update", "delete", "$transaction".
license: MIT
metadata:
  author: prisma
  version: "7.0.0"
---

This skill covers the Prisma Client v7 API: instantiation with driver adapters, all model query methods (find, create, update, delete, upsert, aggregate, groupBy), query options (where, select, include, omit, orderBy, pagination, cursor), filter and relation operators, transaction patterns, raw SQL methods, and client lifecycle methods.

## References

| File | Purpose |
|------|---------|
| references/constructor.md | PrismaClient constructor options |
| references/model-queries.md | CRUD operations |
| references/query-options.md | select, include, omit, where, orderBy |
| references/filters.md | Filter conditions and operators |
| references/relations.md | Relation queries and nested operations |
| references/transactions.md | Transaction API |
| references/raw-queries.md | $queryRaw, $executeRaw |
| references/client-methods.md | $connect, $disconnect, $on, $extends |
