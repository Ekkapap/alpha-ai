---
name: tanstack-query-best-practices
description: TanStack Query (React Query) best practices for data fetching, caching, mutations, and server state management. Activate when building data-driven React applications with server state.
---

TanStack Query manages server state in React applications through a rules-based system covering query keys, caching strategies, mutations with optimistic updates, SSR hydration, prefetching, and offline support. Rules are organized by priority from CRITICAL (query keys, caching) down to LOW (performance, offline) and live in individual files under `rules/`.

## References

| File | Purpose |
|------|---------|
| rules/qk-array-structure.md | Always use arrays for query keys |
| rules/qk-factory-pattern.md | Query key factories for complex apps |
| rules/qk-hierarchical-organization.md | Hierarchical key organization |
| rules/qk-include-dependencies.md | Include all query variables in key |
| rules/qk-serializable.md | Ensure keys are JSON-serializable |
| rules/cache-stale-time.md | staleTime based on data volatility |
| rules/cache-gc-time.md | gcTime for inactive query retention |
| rules/cache-invalidation.md | Targeted invalidation patterns |
| rules/cache-placeholder-vs-initial.md | Placeholder vs initial data |
| rules/mut-invalidate-queries.md | Invalidate related queries after mutations |
| rules/mut-optimistic-updates.md | Optimistic UI updates |
| rules/mut-mutation-state.md | Cross-component mutation state tracking |
| rules/err-error-boundaries.md | Error boundaries with query reset |
| rules/pf-intent-prefetch.md | Prefetch on hover/focus |
| rules/inf-page-params.md | Infinite query pagination |
| rules/ssr-dehydration.md | SSR dehydrate/hydrate pattern |
| rules/parallel-use-queries.md | Dynamic parallel queries |
| rules/query-cancellation.md | Query cancellation |
| rules/perf-select-transform.md | Select to transform/filter data |
| rules/network-mode.md | Offline network mode |
| rules/persist-queries.md | Query persistence for offline |
