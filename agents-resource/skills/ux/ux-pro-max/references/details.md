# UI/UX Pro Max — Detailed Reference

This file contains the full quick reference rules, usage workflow, search commands, and pre-delivery checklist extracted from the main skill body.

## Quick Reference Rules by Priority

### 1. Accessibility (CRITICAL)

- `color-contrast` — Minimum 4.5:1 for normal text (3:1 large text)
- `focus-states` — Visible focus rings (2–4px) on interactive elements
- `alt-text` — Descriptive alt text for meaningful images
- `aria-labels` — aria-label for icon-only buttons
- `keyboard-nav` — Tab order matches visual order; full keyboard support
- `form-labels` — Use label with for attribute
- `skip-links` — Skip to main content for keyboard users
- `heading-hierarchy` — Sequential h1→h6, no level skip
- `color-not-only` — Don't convey info by color alone (add icon/text)
- `dynamic-type` — Support system text scaling (Apple Dynamic Type, MD)
- `reduced-motion` — Respect prefers-reduced-motion
- `voiceover-sr` — Meaningful accessibilityLabel/accessibilityHint; logical reading order
- `escape-routes` — Provide cancel/back in modals and multi-step flows
- `keyboard-shortcuts` — Preserve system and a11y shortcuts

### 2. Touch & Interaction (CRITICAL)

- `touch-target-size` — Min 44×44pt (Apple) / 48×48dp (Material)
- `touch-spacing` — Minimum 8px/8dp gap between touch targets
- `hover-vs-tap` — Click/tap for primary interactions; don't rely on hover alone
- `loading-buttons` — Disable button during async; show spinner or progress
- `error-feedback` — Clear error messages near problem
- `cursor-pointer` — Add cursor-pointer to clickable elements (Web)
- `gesture-conflicts` — Avoid horizontal swipe on main content
- `tap-delay` — Use touch-action: manipulation to reduce 300ms delay
- `standard-gestures` — Use platform standard gestures consistently
- `system-gestures` — Don't block system gestures (Control Center, back swipe)
- `press-feedback` — Visual feedback on press (ripple/highlight)
- `haptic-feedback` — Use haptic for confirmations; avoid overuse
- `safe-area-awareness` — Keep targets away from notch, Dynamic Island, gesture bar

### 3. Performance (HIGH)

- `image-optimization` — WebP/AVIF, responsive images, lazy load non-critical
- `image-dimension` — Declare width/height or use aspect-ratio (CLS prevention)
- `font-loading` — font-display: swap/optional; avoid FOIT
- `critical-css` — Prioritize above-the-fold CSS
- `lazy-loading` — Lazy load non-hero components
- `bundle-splitting` — Split by route/feature; reduce initial load and TTI
- `third-party-scripts` — Load async/defer; audit unnecessary ones
- `virtualize-lists` — Virtualize lists with 50+ items
- `main-thread-budget` — Keep per-frame work under ~16ms
- `progressive-loading` — Skeleton screens for >1s operations
- `input-latency` — Keep input latency under ~100ms
- `debounce-throttle` — Debounce/throttle for high-frequency events
- `offline-support` — Provide offline state messaging and fallback

### 4. Style Selection (HIGH)

- `style-match` — Match style to product type
- `consistency` — Use same style across all pages
- `no-emoji-icons` — Use SVG icons, not emojis
- `color-palette-from-product` — Choose palette from product/industry
- `platform-adaptive` — Respect platform idioms (iOS HIG vs Material)
- `elevation-consistent` — Consistent elevation/shadow scale
- `dark-mode-pairing` — Design light/dark variants together
- `icon-style-consistent` — One icon set/visual language across product
- `primary-action` — One primary CTA per screen

### 5. Layout & Responsive (HIGH)

- `viewport-meta` — width=device-width initial-scale=1 (never disable zoom)
- `mobile-first` — Design mobile-first
- `breakpoint-consistency` — Systematic breakpoints (375/768/1024/1440)
- `readable-font-size` — Minimum 16px body text on mobile
- `line-length-control` — Mobile 35–60 chars; desktop 60–75 chars
- `horizontal-scroll` — No horizontal scroll on mobile
- `spacing-scale` — 4pt/8dp incremental spacing system
- `z-index-management` — Define layered z-index scale
- `fixed-element-offset` — Fixed navbar must reserve safe padding

### 6. Typography & Color (MEDIUM)

- `line-height` — 1.5–1.75 for body text
- `font-scale` — Consistent type scale (12/14/16/18/24/32)
- `color-semantic` — Semantic color tokens (primary, secondary, error, surface)
- `color-dark-mode` — Desaturated/lighter tonal variants, not inverted colors
- `color-accessible-pairs` — Meet 4.5:1 (AA) or 7:1 (AAA)
- `number-tabular` — Tabular figures for data columns, prices, timers

### 7. Animation (MEDIUM)

- `duration-timing` — 150–300ms for micro-interactions; ≤400ms complex
- `transform-performance` — Use transform/opacity only; avoid width/height animation
- `excessive-motion` — Animate 1–2 key elements per view max
- `easing` — ease-out for entering, ease-in for exiting
- `motion-meaning` — Every animation expresses cause-effect
- `exit-faster-than-enter` — Exit ~60–70% of enter duration
- `stagger-sequence` — 30–50ms stagger for list/grid entrance
- `interruptible` — Animations must be interruptible by user
- `no-blocking-animation` — Never block user input during animation

### 8. Forms & Feedback (MEDIUM)

- `input-labels` — Visible label per input (not placeholder-only)
- `error-placement` — Show error below the related field
- `submit-feedback` — Loading then success/error state
- `required-indicators` — Mark required fields
- `empty-states` — Helpful message and action when no content
- `toast-dismiss` — Auto-dismiss toasts in 3–5s
- `confirmation-dialogs` — Confirm before destructive actions
- `inline-validation` — Validate on blur, not keystroke
- `input-type-keyboard` — Semantic input types for mobile keyboard
- `undo-support` — Allow undo for destructive/bulk actions
- `error-recovery` — Error messages include clear recovery path
- `form-autosave` — Auto-save drafts for long forms
- `focus-management` — After submit error, focus first invalid field

### 9. Navigation Patterns (HIGH)

- `bottom-nav-limit` — Bottom nav max 5 items with labels
- `back-behavior` — Predictable and consistent back navigation
- `deep-linking` — All key screens reachable via deep link/URL
- `nav-label-icon` — Both icon and text label; not icon-only
- `nav-state-active` — Current location visually highlighted
- `modal-escape` — Clear close/dismiss affordance on modals
- `state-preservation` — Back navigation restores scroll position and state
- `adaptive-navigation` — Sidebar on ≥1024px; bottom/top on mobile

### 10. Charts & Data (LOW)

- `chart-type` — Match chart type to data (trend→line, comparison→bar)
- `color-guidance` — Accessible palettes; avoid red/green only for colorblind
- `pattern-texture` — Supplement color with patterns/shapes
- `legend-visible` — Always show legend near chart
- `tooltip-on-interact` — Tooltips on hover/tap with exact values
- `axis-labels` — Label axes with units; no truncated labels
- `responsive-chart` — Reflow or simplify on small screens
- `empty-data-state` — Meaningful empty state, not blank chart
- `touch-target-chart` — ≥44pt tap area for chart elements

---

## How to Use This Skill

### Step 1: Analyze User Requirements

Extract: product type, target audience, style keywords, stack (React Native for this project).

### Step 2: Generate Design System (REQUIRED)

```bash
python3 skills/ui-ux-pro-max/scripts/search.py "<product_type> <industry> <keywords>" --design-system [-p "Project Name"]
```

Searches domains in parallel (product, style, color, landing, typography) with reasoning from `ui-reasoning.csv`.

**Persist for hierarchical retrieval:**
```bash
python3 skills/ui-ux-pro-max/scripts/search.py "<query>" --design-system --persist -p "Project Name"
# Creates design-system/MASTER.md + design-system/pages/{page}.md
```

### Step 3: Supplement with Detailed Searches

```bash
python3 skills/ui-ux-pro-max/scripts/search.py "<keyword>" --domain <domain> [-n <max_results>]
```

| Domain | Use For |
|--------|---------|
| `product` | Product type patterns |
| `style` | UI styles, visual effects |
| `typography` | Font pairings |
| `color` | Color palettes by industry |
| `chart` | Chart type recommendations |
| `ux` | UX best practices |
| `google-fonts` | Individual Google Fonts |
| `react` | React/Next.js performance |
| `web` | App interface guidelines (iOS/Android) |

### Step 4: Stack Guidelines (React Native)

```bash
python3 skills/ui-ux-pro-max/scripts/search.py "<keyword>" --stack react-native
```

---

## Common Sticking Points

| Problem | Resolution |
|---------|------------|
| Can't decide on style/color | Re-run `--design-system` with different keywords |
| Dark mode contrast issues | Quick ref §6: `color-dark-mode` + `color-accessible-pairs` |
| Animations feel unnatural | §7: `spring-physics` + `easing` + `exit-faster-than-enter` |
| Form UX is poor | §8: `inline-validation` + `error-clarity` + `focus-management` |
| Navigation feels confusing | §9: `nav-hierarchy` + `bottom-nav-limit` + `back-behavior` |
| Layout breaks on small screens | §5: `mobile-first` + `breakpoint-consistency` |
| Performance / jank | §3: `virtualize-lists` + `main-thread-budget` + `debounce-throttle` |

---

## Pre-Delivery Checklist (App UI)

### Visual Quality
- [ ] No emojis used as icons (use SVG instead)
- [ ] All icons from consistent family and style
- [ ] Official brand assets with correct proportions
- [ ] Semantic theme tokens (no ad-hoc hardcoded colors)

### Interaction
- [ ] All tappable elements provide clear pressed feedback
- [ ] Touch targets ≥44x44pt iOS / ≥48x48dp Android
- [ ] Micro-interaction timing 150–300ms with native easing
- [ ] Disabled states visually clear and non-interactive
- [ ] Screen reader focus order matches visual order
- [ ] No nested/conflicting gesture regions

### Light/Dark Mode
- [ ] Primary text contrast ≥4.5:1 in both modes
- [ ] Secondary text contrast ≥3:1 in both modes
- [ ] Dividers/borders distinguishable in both modes
- [ ] Modal scrim 40–60% black opacity
- [ ] Both themes tested before delivery

### Layout
- [ ] Safe areas respected (headers, tab bars, CTA bars)
- [ ] Scroll content not hidden behind fixed/sticky bars
- [ ] Verified on small phone, large phone, and tablet (portrait + landscape)
- [ ] 4/8dp spacing rhythm maintained
- [ ] Long-form text readable on large devices

### Accessibility
- [ ] All meaningful images/icons have accessibility labels
- [ ] Form fields have labels, hints, clear error messages
- [ ] Color not the only indicator
- [ ] Reduced motion and dynamic text size supported
- [ ] Accessibility traits/roles/states announced correctly
