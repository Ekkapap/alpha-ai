---
name: electric-db
description: ElectricSQL + TanStack DB real-time sync. Use when replacing polling with Postgres-to-client streaming, reactive collections, or optimistic updates.
---

reference-path = ./data/electric-db/

ElectricSQL syncs Postgres data to the client via HTTP streaming; TanStack DB receives it as reactive collections and performs joins on the client side. Use when replacing polling-based data fetching with real-time reactive state, or when implementing optimistic updates with sub-second latency.

## References

| File | Purpose |
|------|---------|
| references/details.md | Core concepts (Shape, Collection, useQuery), best practices, full example, comparison table, limitations |
