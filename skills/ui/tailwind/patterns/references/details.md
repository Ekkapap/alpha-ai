# Tailwind CSS Development Patterns

## Overview

Expert guide for building modern, responsive user interfaces with Tailwind CSS utility-first framework. Covers v4.1+ features including CSS-first configuration, custom utilities, and enhanced developer experience.

## Instructions

1. **Start Mobile-First**: Write base styles for mobile, add responsive prefixes for larger screens
2. **Use Design Tokens**: Leverage Tailwind's spacing, color, and typography scales
3. **Compose Utilities**: Combine multiple utilities for complex styles
4. **Extract Components**: Create reusable component classes for repeated patterns
5. **Configure Theme**: Customize design tokens in tailwind.config.js
6. **Optimize for Production**: Ensure content paths are configured for CSS purging
7. **Test Responsive**: Verify layouts at all breakpoint sizes

## Constraints and Warnings

- **Class Proliferation**: Long class strings can reduce readability; extract components when needed
- **Purge Configuration**: Must configure content paths correctly for production builds
- **Arbitrary Values**: Use sparingly; prefer design tokens for consistency
- **Specificity Issues**: Avoid `@apply` with complex selectors
- **Dark Mode**: Requires proper configuration (class or media strategy)
- **JIT Mode**: Some dynamic patterns may not be detected; use safelist if needed

## Core Concepts

### Utility-First Approach

```html
<button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Click me</button>
```

### Responsive Design

Mobile-first breakpoints: `sm:` (640px+), `md:` (768px+), `lg:` (1024px+), `xl:` (1280px+), `2xl:` (1536px+)

```html
<div class="w-full md:w-1/2 lg:w-1/3">...</div>
```

## Layout Utilities

### Flexbox Layouts

```html
<div class="flex items-center justify-between">...</div>
<div class="flex flex-col md:flex-row gap-4"><div class="flex-1">Item 1</div></div>
<div class="flex items-center justify-center min-h-screen">Centered</div>
<div class="flex flex-col gap-4">...</div>
```

### Grid Layouts

```html
<div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">...</div>
<div class="grid grid-cols-[repeat(auto-fit,minmax(250px,1fr))] gap-4">...</div>
```

### Container

```html
<div class="container mx-auto px-4 max-w-7xl">...</div>
```

## Spacing, Typography, Colors

### Spacing

```html
<div class="p-4 md:p-8 lg:p-12">Responsive padding</div>
<div class="space-y-4">Vertical stack with gap</div>
<div class="px-4 py-8">Axis-based</div>
```

### Typography

```html
<h1 class="text-4xl font-bold">Large Heading</h1>
<h1 class="text-2xl md:text-4xl lg:text-6xl font-bold">Responsive Heading</h1>
<p class="leading-relaxed tracking-wide">Body text</p>
```

### Colors

```html
<div class="bg-blue-500">Blue background</div>
<div class="bg-gradient-to-r from-blue-500 to-purple-600">Gradient</div>
<div class="bg-blue-500 bg-opacity-50">Semi-transparent</div>
```

## Interactive States

```html
<button class="bg-blue-500 hover:bg-blue-700 transition">Hover</button>
<input class="border border-gray-300 focus:border-blue-500 focus:ring-2 focus:ring-blue-200 outline-none">
<button class="bg-blue-500 active:bg-blue-800 disabled:opacity-50 disabled:cursor-not-allowed">Button</button>
<div class="group"><img class="group-hover:opacity-75" /><p class="group-hover:text-blue-600">...</p></div>
```

## Component Patterns

### Card Component

```html
<div class="bg-white rounded-lg shadow-lg overflow-hidden">
  <img class="w-full h-48 object-cover" src="image.jpg" alt="Card image" />
  <div class="p-6">
    <h3 class="text-xl font-bold mb-2">Card Title</h3>
    <p class="text-gray-700 mb-4">Card description.</p>
    <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded">Action</button>
  </div>
</div>
```

### Navigation Bar

```html
<nav class="bg-white shadow-lg">
  <div class="container mx-auto px-4">
    <div class="flex justify-between items-center h-16">
      <a href="#" class="text-xl font-bold text-gray-800">Logo</a>
      <div class="hidden md:flex space-x-8">
        <a href="#" class="text-gray-700 hover:text-blue-600 transition">Home</a>
        <a href="#" class="text-gray-700 hover:text-blue-600 transition">About</a>
      </div>
    </div>
  </div>
</nav>
```

### Form Elements

```html
<form class="space-y-6 max-w-md mx-auto">
  <div>
    <label class="block text-sm font-medium text-gray-700 mb-2">Email</label>
    <input type="email" class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent" />
  </div>
  <button type="submit" class="w-full bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg hover:bg-blue-700 transition">Sign In</button>
</form>
```

### Modal/Dialog

```html
<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4">
  <div class="bg-white rounded-lg shadow-xl max-w-md w-full p-6">
    <div class="flex justify-between items-center mb-4">
      <h3 class="text-xl font-bold">Modal Title</h3>
      <button class="text-gray-500 hover:text-gray-700">✕</button>
    </div>
    <p class="text-gray-700 mb-6">Modal content.</p>
    <div class="flex justify-end space-x-4">
      <button class="px-4 py-2 text-gray-600 hover:text-gray-800">Cancel</button>
      <button class="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700">Confirm</button>
    </div>
  </div>
</div>
```

## Dark Mode

```html
<div class="bg-white dark:bg-gray-900 text-gray-900 dark:text-white">
  <p class="text-gray-600 dark:text-gray-400">Description</p>
</div>
```

## Animations & Transitions

```html
<button class="bg-blue-500 hover:bg-blue-700 transition duration-300">Smooth</button>
<div class="transform hover:scale-110 transition duration-300">Scale on hover</div>
<div class="animate-spin">Spinning</div>
<div class="animate-pulse">Pulsing</div>
<div class="transform transition-transform motion-reduce:transition-none">Respects prefers-reduced-motion</div>
```

## Accessibility

```html
<button class="focus:outline-none focus:ring-4 focus:ring-blue-500 focus:ring-offset-2">Accessible Button</button>
<a href="#main-content" class="sr-only focus:not-sr-only focus:absolute focus:top-4 focus:left-4">Skip to main content</a>
<button aria-label="Close dialog" class="p-2">✕</button>
```

## Configuration

### CSS-First (v4.1+)

```css
@import "tailwindcss";

@theme {
  --color-brand-500: #3b82f6;
  --font-display: "Inter", system-ui, sans-serif;
  --spacing-128: 32rem;
  --breakpoint-3xl: 1920px;
}

@utility content-auto {
  content-visibility: auto;
}
```

### Vite Integration (v4.1+)

```javascript
// vite.config.ts
import tailwindcss from '@tailwindcss/vite'
export default defineConfig({ plugins: [tailwindcss()] })
```

### Legacy JS Config

```javascript
export default {
  content: ["./index.html", "./src/**/*.{js,ts,jsx,tsx,vue,svelte}"],
  theme: { extend: { colors: { primary: { 500: '#3b82f6' } } } },
}
```

## Container Queries

```html
<div class="@container">
  <div class="@lg:text-xl @2xl:text-2xl">Text size based on container</div>
</div>
```

## React/JSX Pattern

```tsx
function Button({ variant = 'primary', size = 'md', children }) {
  const variantClasses = {
    primary: 'bg-blue-600 text-white hover:bg-blue-700',
    secondary: 'bg-gray-200 text-gray-800 hover:bg-gray-300',
  }
  const sizeClasses = { sm: 'px-3 py-1 text-sm', md: 'px-4 py-2 text-base', lg: 'px-6 py-3 text-lg' }
  return <button className={`font-semibold rounded transition ${variantClasses[variant]} ${sizeClasses[size]}`}>{children}</button>
}
```
