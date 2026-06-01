---
name: tanstack-store
description: Framework-agnostic, immutable reactive data store with framework adapters for React, Vue, Solid, Angular, and Svelte.
---

TanStack Store is a lightweight signals-like reactive store providing `Store` for state, `Derived` for computed values, `Effect` for side effects, and `batch` for atomic updates. Framework adapters (React, Vue, Solid, Angular, Svelte) expose reactive hooks; always call `mount()` on Derived/Effect instances and clean up the returned unmount function.

## References

| File | Purpose |
|------|---------|
| references/details.md | API reference, React integration, framework adapters, best practices, common pitfalls |
