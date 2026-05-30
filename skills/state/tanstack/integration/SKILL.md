---
name: tanstack-integration-best-practices
description: Best practices for integrating TanStack Query with TanStack Router and TanStack Start. Patterns for full-stack data flow, SSR, and caching coordination.
---

TanStack Integration covers the patterns for combining TanStack Query, Router, and Start: passing QueryClient through router context, using `ensureQueryData` in loaders with `useSuspenseQuery` in components, coordinating cache invalidation across mutations, and using `setupRouterSsrQueryIntegration` for automatic SSR dehydration/hydration. Rules live in `rules/`.

## References

| File | Purpose |
|------|---------|
| rules/setup-query-client-context.md | Pass QueryClient through router context |
| rules/flow-loader-query-pattern.md | Loaders with ensureQueryData |
| rules/cache-single-source.md | Let TanStack Query own caching |
| rules/ssr-dehydrate-hydrate.md | Automatic SSR integration |
