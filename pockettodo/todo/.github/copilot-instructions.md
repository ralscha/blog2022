# Copilot Instructions for Todo App

## Project Overview

This is an **Ionic Angular** mobile/web application using Angular 20+ with standalone components architecture. It's a todo app.

## Architecture Patterns

### Standalone Components Architecture

- All components use Angular standalone components (no NgModules)
- Import Ionic components directly in each component: `import {IonHeader, IonToolbar, IonTitle, IonContent} from '@ionic/angular/standalone'`
- Components follow the pattern: `selector: 'app-{name}'`, `templateUrl`, `styleUrl`, and explicit `imports` array

### Routing & Navigation

- Uses Angular router with lazy loading: `loadComponent: () => import('./path/component').then((m) => m.ComponentName)`
- Route configuration in `src/app/app.routes.ts`
- Main bootstrap in `src/main.ts` with IonicRouteStrategy and PreloadAllModules

### File Structure Conventions

- Pages use `.page.ts` suffix (e.g., `home.page.ts`)
- Components use `.component.ts` suffix
- Page styles use `.page.css` suffix
- Each page/component has its own directory with `.ts`, `.html`, and `.css` files

## Development Workflows

### Build & Serve Commands

```bash
npm start          # Development server
npm run build      # Production build
npm run watch      # Development build with file watching
```

### Ionic-Specific Commands

- Use `ionic serve` for enhanced development experience with live reload

## Styling Architecture

### SCSS Configuration

- Global styles in `src/styles.scss` with Ionic CSS imports
- Ionic variables in `src/variables.scss` (currently minimal)
- Component-specific styles use `.css` files (not `.scss` by default)

### Ionic CSS Framework

- Import core Ionic CSS: `@import "@ionic/angular/css/core.css"`
- Use Ionic's built-in utility classes for layout and typography
- Follow Ionic's design system for consistent mobile UI

## Key Dependencies

### Core Stack

- **Angular 20.1.2** - Latest Angular with modern features
- **Ionic Angular 8.6.5** - Mobile-first UI framework

### TypeScript Configuration

- Strict mode enabled with comprehensive type checking
- Target: ES2022 with modern JavaScript features
- Experimental decorators enabled for Angular

## Component Development Patterns

### Ionic Component Imports

```typescript
// Always import specific Ionic components needed
import {IonHeader, IonToolbar, IonTitle, IonContent} from '@ionic/angular/standalone';

@Component({
  selector: 'app-example',
  templateUrl: './example.page.html',
  styleUrl: './example.page.css',
  imports: [IonHeader, IonToolbar, IonTitle, IonContent], // Explicit imports required
})
```

### Page Template Structure

```html
<!-- Standard Ionic page structure -->
<ion-header [translucent]="true">
  <ion-toolbar>
    <ion-title>Page Title</ion-title>
  </ion-toolbar>
</ion-header>

<ion-content [fullscreen]="true">
  <!-- Collapsible header for better UX -->
  <ion-header collapse="condense">
    <ion-toolbar>
      <ion-title size="large">Page Title</ion-title>
    </ion-toolbar>
  </ion-header>

  <!-- Page content -->
</ion-content>
```

## Environment & Configuration

- Environment switching via `src/environments/` files
- Angular build replaces `environment.ts` with `environment.prod.ts` for production
- Ionic configuration in `ionic.config.json`

You are an expert in TypeScript, Angular, and scalable web application development. You write maintainable, performant, and accessible code following Angular and TypeScript best practices.

## TypeScript Best Practices

- Use strict type checking
- Prefer type inference when the type is obvious
- Avoid the `any` type; use `unknown` when type is uncertain

## Angular Best Practices

- Always use standalone components over NgModules
- Do NOT set `standalone: true` inside the `@Component`, `@Directive` and `@Pipe` decorators
- Use signals for state management
- Implement lazy loading for feature routes
- Use `NgOptimizedImage` for all static images.
- Do NOT use the `@HostBinding` and `@HostListener` decorators. Put host bindings inside the `host` object of the `@Component` or `@Directive` decorator instead

## Components

- Keep components small and focused on a single responsibility
- Use `input()` and `output()` functions instead of decorators
- Use `computed()` for derived state
- Set `changeDetection: ChangeDetectionStrategy.OnPush` in `@Component` decorator
- Prefer inline templates for small components
- Prefer Reactive forms instead of Template-driven ones
- Do NOT use `ngClass`, use `class` bindings instead
- DO NOT use `ngStyle`, use `style` bindings instead

## State Management

- Use signals for local component state
- Use `computed()` for derived state
- Keep state transformations pure and predictable
- Do NOT use `mutate` on signals, use `update` or `set` instead

## Templates

- Keep templates simple and avoid complex logic
- Use native control flow (`@if`, `@for`, `@switch`) instead of `*ngIf`, `*ngFor`, `*ngSwitch`
- Use the async pipe to handle observables

## Services

- Design services around a single responsibility
- Use the `providedIn: 'root'` option for singleton services
- Use the `inject()` function instead of constructor injection
