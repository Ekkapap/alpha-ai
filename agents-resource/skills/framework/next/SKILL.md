---
name: next-best-practices
description: Next.js best practices - file conventions, RSC boundaries, data patterns, async APIs, metadata, error handling, route handlers, image/font optimization, bundling
user-invocable: false
---

Next.js best practices covering file conventions, RSC boundaries, async APIs (Next.js 15+), directives, functions, error handling, data patterns, route handlers, metadata, image/font optimization, bundling, scripts, hydration errors, Suspense boundaries, parallel/intercepting routes, self-hosting, and debug tricks. All detail lives in `references/`.

## References

| File | Purpose |
|------|---------|
| references/file-conventions.md | Project structure, special files, route segments, parallel/intercepting routes, middleware rename |
| references/rsc-boundaries.md | Async client component detection, non-serializable props, Server Action exceptions |
| references/async-patterns.md | Async params/searchParams/cookies/headers, migration codemod |
| references/runtime-selection.md | Node.js vs Edge runtime selection |
| references/directives.md | use client, use server, use cache |
| references/functions.md | Navigation hooks, server functions, generate functions |
| references/error-handling.md | error.tsx, redirect, notFound, forbidden, unstable_rethrow |
| references/data-patterns.md | Server Components vs Actions vs Route Handlers, waterfall avoidance |
| references/route-handlers.md | route.ts basics, GET conflicts, when to use vs Server Actions |
| references/metadata.md | Static/dynamic metadata, OG images, file-based conventions |
| references/image.md | next/image, remote images, responsive sizes, blur placeholders |
| references/font.md | next/font setup, Google Fonts, local fonts, Tailwind integration |
| references/bundling.md | Server-incompatible packages, CSS imports, ESM/CommonJS |
| references/scripts.md | next/script vs native, inline scripts, Google Analytics |
| references/hydration-error.md | Common causes, debugging, fixes |
| references/suspense-boundaries.md | CSR bailout with useSearchParams/usePathname |
| references/parallel-routes.md | Modal patterns, @slot, interceptors, default.tsx |
| references/self-hosting.md | Docker standalone, cache handlers, multi-instance ISR |
| references/debug-tricks.md | MCP endpoint, --debug-build-paths |
