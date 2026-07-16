# Lifeline Web Design System

`palette-preview.html` is the visual reference surface for the Lifeline web interface. For concrete examples, see this file: `internal/presentation/http/web/palette-preview.html`.

This document describes the current design rules that should be used for future web screens. It is no longer a migration plan.

## Color Tokens

All shared color roles live in `src/index.css` under `:root`.

- Background: `--background #F7F7FA` is the page-level canvas.
- Primary surfaces: `--surface #FFFFFF` is used for panels, cards, inputs, and swatches.
- Secondary surfaces: `--surface-secondary #F0F0F5` is used for nested rows, metrics, and low-emphasis controls.
- Selected surfaces: `--surface-selected #E8E8FF` is used for active rows and selected navigation states.
- Text hierarchy: `--text-primary #19191F`, `--text-secondary #65656F`, `--text-disabled #9696A0`, and `--text-on-primary #FFFFFF`.
- Primary accent: `--primary #5556D8`, `--primary-hover #4849C2`, `--primary-container #E8E8FF`, and `--on-primary-container #25256F`.
- Borders and focus: `--divider #E1E1E8`, `--outline #C7C7D0`, and `--outline-focused #5556D8`.
- Status roles use paired foreground/container tokens: success, warning, error, and info.

Use tokens instead of page-local color literals. New page-local colors should be introduced only when a new shared design role is intentionally being created.

## Shape, Shadow, And Depth

- `--radius-lg 24px` is for major panels.
- `--radius-md 16px` is for repeated rows, metrics, and messages.
- `--radius-sm 12px` is for smaller controls when needed.
- `--shadow 0 14px 40px rgba(25, 25, 31, 0.08)` is reserved for top-level panels.

Depth rules:

- Page background stays flat.
- Panels use border plus shadow.
- Rows and metrics sit inside panels without extra shadows.
- Selected and status states are communicated primarily through color, not elevation.

## Spacing And Layout

- Pages use `Page` from `src/components/layout`.
- Page width is constrained with `width: min(1180px, calc(100% - 32px))`.
- Desktop page padding is `40px 0 64px`; mobile top padding is `20px`.
- Dashboard-like layouts use a two-column grid: `1.25fr` main content and `0.75fr` side content with a minimum side width around `320px`.
- Panel body padding is `24px`, reduced to `18px` on small screens.
- Sections are separated by `28px` vertical rhythm and a divider.
- Repeated lists use compact gaps around `10px`.
- Form and button groups use `10px-18px` local spacing.

## Typography

- Font stack starts with `Inter`, then system UI fonts.
- Eyebrows use uppercase, 13px text, heavy weight, and positive letter spacing.
- Page titles use `clamp(32px, 6vw, 56px)`, line-height `1`, and tight letter spacing.
- Section titles are compact: 20px with slight negative letter spacing.
- Secondary text uses `--text-secondary`, usually 13px or 14px.
- Numeric metric values are heavy and tight, using 24px and high font weight.

## Implemented Components

### Layout Components

Location: `src/components/layout`.

- `Page`: constrained width and standard vertical padding.
- `PageHeader`: eyebrow, title, subtitle, and optional actions.
- `Panel`: top-level raised surface.
- `PanelBody`: standard panel content padding.
- `Section`: vertically separated content block inside a panel.
- `SectionHeader`: compact title row with optional metadata.

Use these before adding page-local wrappers.

### Primitive Components

Location: `src/components/primitives`.

- `Button`: standard action button with `primary`, `secondary`, and `danger` variants.
- `IconButton`: circular icon-only action button.
- `TextField`: labeled input with focus and error styling.
- `Message`: status block with `success`, `warning`, `error`, and `info` variants.
- `Metric`: compact value and label pair.
- `NavigationItem`: compact navigation target with optional icon and active state.

Use primitives instead of duplicating button, field, message, metric, or navigation styling inside page SASS.

## Page-Specific SASS Rules

Page SASS should only handle page-specific composition.

- Do not define page-local color literals unless a new shared token is being introduced.
- Do not repeat panel, card, input, button, message, metric, or navigation styling in page files.
- Do not use gradient page backgrounds.
- Keep BEM naming for page-specific classes and component classes.
- Prefer design tokens and shared components over new local styles.
- Keep page SASS small enough to show layout intent, not component implementation.

## Responsive Behavior

- At `860px`, multi-column layouts collapse to one column and page headers stack vertically.
- At `560px`, page width tightens, panel padding decreases, summaries become single-column, and dense row actions can move to a second row.
- Text should not overflow its container on mobile or desktop.
- Operational screens should stay dense, readable, and calm rather than using marketing-style hero composition.

## Current App Surfaces

- `HomePage` is the initial dashboard shell and should be the model for future app surfaces.
- `LoginPage` and `SignupPage` use the same layout and primitive components as the app surfaces.
- API client and routing behavior are separate from the visual system and should remain stable during visual changes.
