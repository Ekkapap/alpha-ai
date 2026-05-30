---
name: best-practices
description: "Mantine UI library for React: 100+ components, hooks, forms, theming, dark mode, CSS modules, and Vite/TypeScript setup. Use when building React applications with Mantine components, configuring theming/dark mode, or working with Mantine hooks and forms. Keywords: Mantine, React, UI components, CSS modules, theming."
metadata:
  version: "8.3.17"
  release_date: "2026-03-14"
---

Mantine is a fully-featured React component library (100+ components, hooks, forms) with TypeScript, native dark mode via CSS variables, CSS modules with PostCSS, and excellent accessibility. Always wrap the app in `MantineProvider`, import `@mantine/core/styles.css`, configure PostCSS with `postcss-preset-mantine`, and never skip `form.key('path')` in uncontrolled forms.

## References

| File | Purpose |
|------|---------|
| references/getting-started.md | Installation, Vite setup, project structure |
| references/styling.md | MantineProvider, theme, CSS modules, style props, dark mode |
| references/components.md | Core UI component patterns |
| references/hooks.md | @mantine/hooks utility hooks |
| references/forms.md | @mantine/form, useForm, validation |
| references/testing.md | Vitest setup, custom render, mocking |
| references/eslint.md | eslint-config-mantine setup |
