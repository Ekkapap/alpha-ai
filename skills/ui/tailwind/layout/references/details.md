# Tailwind CSS Advanced Layout Techniques

## CSS Grid Mastery

### Complex Grid Layouts

```html
<!-- Holy Grail Layout -->
<div class="grid min-h-screen grid-rows-[auto_1fr_auto]">
  <header class="bg-white shadow">Header</header>
  <div class="grid grid-cols-[250px_1fr_300px]">
    <aside class="bg-gray-50 p-4">Sidebar</aside>
    <main class="p-6">Main Content</main>
    <aside class="bg-gray-50 p-4">Right Sidebar</aside>
  </div>
  <footer class="bg-gray-800 text-white">Footer</footer>
</div>

<!-- Responsive Holy Grail -->
<div class="grid min-h-screen grid-rows-[auto_1fr_auto]">
  <header>Header</header>
  <div class="grid grid-cols-1 md:grid-cols-[250px_1fr] lg:grid-cols-[250px_1fr_300px]">
    <aside class="order-2 md:order-1">Sidebar</aside>
    <main class="order-1 md:order-2">Main</main>
    <aside class="order-3 hidden lg:block">Right</aside>
  </div>
  <footer>Footer</footer>
</div>
```

### Grid Template Areas

```css
@utility grid-areas-dashboard {
  grid-template-areas:
    "header header header"
    "nav main aside"
    "nav footer footer";
}

@utility area-header { grid-area: header; }
@utility area-nav { grid-area: nav; }
@utility area-main { grid-area: main; }
@utility area-aside { grid-area: aside; }
@utility area-footer { grid-area: footer; }
```

### Auto-Fill and Auto-Fit Grids

```html
<!-- Auto-fill: Creates as many tracks as fit, even empty ones -->
<div class="grid grid-cols-[repeat(auto-fill,minmax(250px,1fr))] gap-6">...</div>

<!-- Auto-fit: Collapses empty tracks -->
<div class="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-6">...</div>

<!-- Edge case where container is smaller than minmax min -->
<div class="grid grid-cols-[repeat(auto-fill,minmax(min(100%,300px),1fr))] gap-4">...</div>
```

### Subgrid

```css
@utility subgrid-cols { grid-template-columns: subgrid; }
@utility subgrid-rows { grid-template-rows: subgrid; }
```

```html
<div class="grid grid-cols-4 gap-4">
  <div class="col-span-2 grid subgrid-cols gap-4">
    <div>Aligned to parent column 1</div>
    <div>Aligned to parent column 2</div>
  </div>
</div>
```

## Advanced Flexbox Patterns

### Space Distribution

```html
<div class="flex justify-between">...</div>  <!-- Equal spacing, first/last at edges -->
<div class="flex justify-around">...</div>   <!-- Equal spacing including edges -->
<div class="flex justify-evenly">...</div>   <!-- Double space between items vs edges -->
```

### Flexible Item Sizing

```html
<div class="flex"><div class="flex-1">1/3</div><div class="flex-1">1/3</div><div class="flex-1">1/3</div></div>
<div class="flex"><div class="flex-[2]">2/4</div><div class="flex-1">1/4</div><div class="flex-1">1/4</div></div>

<!-- Fixed + flexible, prevent text overflow -->
<div class="flex min-w-0">
  <div class="shrink-0">Icon</div>
  <div class="min-w-0 truncate">Very long text that should truncate</div>
</div>
```

### Masonry-Like with Flexbox

```html
<div class="flex flex-col flex-wrap h-[800px] gap-4">
  <div class="w-[calc(33.333%-1rem)] h-48">Item 1</div>
  <div class="w-[calc(33.333%-1rem)] h-64">Item 2</div>
</div>
```

## Container Queries

```html
<div class="@container">
  <div class="flex flex-col @md:flex-row @lg:grid @lg:grid-cols-3 gap-4">...</div>
</div>

<!-- Named containers -->
<div class="@container/sidebar">
  <nav class="@[200px]/sidebar:flex-col @[300px]/sidebar:flex-row">Navigation</nav>
</div>

<!-- Container query units -->
<div class="@container">
  <h1 class="text-[5cqw]">Scales with container width</h1>
</div>
```

## Position and Layering

### Sticky Positioning

```html
<header class="sticky top-0 z-50 bg-white/80 backdrop-blur-sm border-b">Navigation</header>
<aside class="sticky top-20 h-[calc(100vh-5rem)] overflow-auto">Sidebar content</aside>

<!-- Sticky table header with frozen column -->
<div class="overflow-auto max-h-96">
  <table>
    <thead class="sticky top-0 bg-white shadow">
      <tr><th class="sticky left-0 bg-white z-10">Corner cell</th></tr>
    </thead>
  </table>
</div>
```

### Fixed Elements

```html
<!-- Mobile bottom navigation -->
<nav class="fixed bottom-0 inset-x-0 z-50 bg-white border-t md:hidden">
  <div class="flex justify-around py-2">...</div>
</nav>

<!-- Floating action button -->
<button class="fixed bottom-6 right-6 z-40 rounded-full bg-brand-500 p-4 shadow-lg">+</button>
```

### Z-Index Management

```css
@theme {
  --z-dropdown: 100; --z-sticky: 200; --z-fixed: 300;
  --z-modal-backdrop: 400; --z-modal: 500;
  --z-popover: 600; --z-tooltip: 700; --z-toast: 800;
}
```

## Overflow and Scrolling

### Custom Scrollbars

```css
@utility scrollbar-thin { scrollbar-width: thin; }
@utility scrollbar-none { scrollbar-width: none; -ms-overflow-style: none; }
```

### Scroll Snap

```html
<!-- Horizontal carousel -->
<div class="flex snap-x snap-mandatory overflow-x-auto gap-4 pb-4">
  <div class="snap-start shrink-0 w-80">Card 1</div>
  <div class="snap-start shrink-0 w-80">Card 2</div>
</div>

<!-- Full-page sections -->
<div class="h-screen snap-y snap-mandatory overflow-y-auto">
  <section class="h-screen snap-start">Section 1</section>
  <section class="h-screen snap-start">Section 2</section>
</div>
```

### Scroll Margin for Anchors

```html
<section id="about" class="scroll-mt-20"><!-- Offset for fixed header --></section>
```

## Aspect Ratio and Object Fit

```html
<div class="aspect-video"><video class="h-full w-full object-cover">...</video></div>
<div class="aspect-square rounded-full overflow-hidden"><img class="h-full w-full object-cover" /></div>
<div class="aspect-[4/3]">4:3 content</div>

<!-- Focus specific part of image -->
<div class="h-64 overflow-hidden">
  <img class="h-full w-full object-cover object-top" src="portrait.jpg" />
</div>
```

## Advanced Spacing

### Logical Properties

```html
<div class="ps-4 pe-6 ms-auto">Padding/margin that respect text direction (LTR/RTL)</div>
```

### Space Between with Dividers

```html
<ul class="divide-y divide-gray-200"><li class="py-4">Item 1</li><li class="py-4">Item 2</li></ul>
<div class="flex divide-x divide-gray-200"><div class="px-4">Section 1</div><div class="px-4">Section 2</div></div>
```

### Negative Margins for Bleeds

```html
<!-- Full-bleed image in padded container -->
<article class="px-6">
  <p>Padded content</p>
  <img src="hero.jpg" class="-mx-6 w-[calc(100%+3rem)]" />
</article>
```

## Multi-Column Layout

```html
<div class="columns-1 sm:columns-2 lg:columns-3 gap-8"><p>Content flows across columns...</p></div>
<div class="columns-[300px] gap-6"><p>As many 300px columns as fit</p></div>

<!-- Prevent breaks inside card -->
<div class="columns-2"><div class="break-inside-avoid mb-4">Card that stays together</div></div>
```

## Responsive Patterns

### Fluid Sizing with Clamp

```html
<section class="py-[clamp(2rem,5vw,6rem)] px-[clamp(1rem,3vw,4rem)]">Responsive section</section>
<div class="mx-auto w-full max-w-[clamp(300px,90vw,1200px)]">Responsive container</div>
```

### Breakpoint-Based Visibility

```html
<button class="md:hidden">Menu (mobile only)</button>
<ul class="hidden md:flex gap-4"><li>Desktop nav</li></ul>
```

## Print Styles

```html
<nav class="print:hidden">Navigation</nav>
<div class="hidden print:block">Print-only content</div>
<div class="print:break-inside-avoid">Keep together on one page</div>
<div class="print:break-before-page">Start on new page</div>
```

## Best Practices

1. **Use Grid for 2D layouts** (`grid grid-cols-3`), **Flexbox for 1D** (`flex items-center`)
2. **Handle flex overflow**: always use `min-w-0` on flex containers with truncating children
3. **Use max-w-prose** for reading content, **container** for page sections
4. **Prefer container queries** over media queries for component-level responsiveness
5. **Test all breakpoints** systematically
