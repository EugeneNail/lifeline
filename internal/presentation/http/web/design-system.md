# Lifeline Web Design System Direction

## Palette Preview Structural Summary

The file `palette-preview.html` is the reference surface for the future Lifeline web interface. За конкретными примерами смотреть вот этот файл: `internal/presentation/http/web/palette-preview.html`.

### Color Roles

- Background: `--background #F7F7FA` is the page-level canvas. It should remain quiet and low contrast.
- Primary surfaces: `--surface #FFFFFF` is used for panels, cards, inputs, and swatches.
- Secondary surfaces: `--surface-secondary #F0F0F5` is used for nested rows, metrics, and low-emphasis controls.
- Selected surfaces: `--surface-selected #E8E8FF` is used for active habit rows and selected navigation states.
- Text hierarchy: `--text-primary #19191F`, `--text-secondary #65656F`, `--text-disabled #9696A0`, and `--text-on-primary #FFFFFF`.
- Primary accent: `--primary #5556D8` with `--primary-hover #4849C2`, `--primary-container #E8E8FF`, and `--on-primary-container #25256F`.
- Borders: `--divider #E1E1E8`, `--outline #C7C7D0`, and `--outline-focused #5556D8`.
- Status colors use paired foreground/container roles: success, warning, error, and info.

### Shape, Shadow, And Depth

- Radius scale:
  - `--radius-lg 24px` for major panels.
  - `--radius-md 16px` for repeated rows, metrics, and messages.
  - `--radius-sm 12px` for smaller controls when needed.
- Shadow: `--shadow 0 14px 40px rgba(25, 25, 31, 0.08)` is reserved for top-level panels, not for every nested element.
- Depth model:
  - Page background is flat.
  - Panels are raised with border and shadow.
  - Rows and metrics sit inside panels without extra shadows.
  - Selected and status states are communicated primarily through color, not elevation.

### Spacing And Layout

- Page width is constrained with `width: min(1180px, calc(100% - 32px))`.
- Desktop page padding is `40px 0 64px`; mobile page top padding is `20px`.
- The main dashboard layout is a two-column grid: `1.25fr` content and `0.75fr` aside with a minimum aside width of `320px`.
- Panel body padding is `24px`, reduced to `18px` on small screens.
- Sections are separated by `28px` vertical rhythm and a divider line.
- Repeated lists use compact grid gaps around `10px`.
- Form and button groups use `10px-16px` local spacing.

### Typography

- Font stack starts with `Inter`, then system UI fonts.
- Eyebrows use uppercase, 13px text, heavy weight, and positive letter spacing.
- Page `h1` is large and compressed: `clamp(32px, 6vw, 56px)`, line-height `1`, and negative letter spacing.
- Section titles are compact: 20px with slight negative letter spacing.
- Secondary text uses `--text-secondary`, usually 13px or 14px.
- Numeric metric values are heavy and tight, using 24px and high font weight.

### Components

- Buttons:
  - Base radius is 14px.
  - Primary buttons use `--primary` and white text.
  - Secondary buttons use `--surface-secondary`.
  - Danger buttons use error container and error foreground.
- Icon buttons are circular, 34px square, and use primary color on white.
- Habit rows are 68px minimum height, use a three-column grid, and can switch to two-column/mobile layout.
- Payment rows mirror habit rows but use white surface and divider border.
- Metrics are simple nested blocks with secondary surface and medium radius.
- Fields use outlined surfaces, 14px radius, and a focus ring based on the primary container.
- Messages use status container colors and status foreground colors.
- Bottom navigation is embedded into the panel, uses a four-column grid, and highlights active items with primary container color.
- Swatches are compact cards for documenting color roles.

### Responsive Behavior

- At `860px`, the main layout collapses to one column and the header stacks vertically.
- At `560px`, page width tightens, panel padding decreases, summary/swatches become single-column, and habit actions move to a second row.
- The design should favor dense, readable operational UI rather than marketing hero layouts.

## Current `src` Directory Summary

- Routing is minimal: `/`, `/login`, `/signup`, and fallback to `/`.
- `HomePage` is a centered session card with two navigation links.
- `LoginPage` and `SignupPage` use separate page-specific SASS files with duplicated auth shell, brand panel, form field, input, button, error, and switch styles.
- Current global CSS still carries a warm radial/linear gradient background and startup-style visual language.
- Current auth pages use dark slate brand panels, translucent white panels, large 28px radius, blur, and heavier decorative depth.
- API client logic is isolated and can remain unchanged while the visual system is replaced.

## Template-Based Migration Approach

### 1. Extract Design Tokens

Create a shared stylesheet, for example `src/styles/tokens.sass` or `src/styles/design-tokens.css`, that defines the palette preview variables at `:root`.

Move these roles first:

- Surface tokens.
- Text tokens.
- Primary tokens.
- Border/focus tokens.
- Status tokens.
- Radius and shadow tokens.

Then remove hard-coded color values from page SASS files by replacing them with variables.

### 2. Normalize Global Page Foundation

Update `src/index.css` to match the preview:

- Use the same font stack.
- Use `background: var(--background)` instead of radial gradients.
- Set `color: var(--text-primary)`.
- Keep reset rules for `box-sizing`, `button`, `input`, and links.
- Avoid global decorative backgrounds; pages should compose from panels and surfaces.

### 3. Introduce Layout Templates

Create reusable layout components before rewriting pages:

- `Page`: constrained width and vertical padding.
- `PageHeader`: eyebrow, title, subtitle, optional actions.
- `Panel`: surface, border, radius, and optional shadow.
- `PanelBody`: consistent internal padding.
- `Section`: vertical rhythm and divider behavior.
- `SectionHeader`: title and optional metadata.

These templates should mirror `.page`, `.page-header`, `.panel`, `.panel-body`, `.section`, and `.section-header` from the preview.

### 4. Introduce Primitive Components

Add small reusable primitives:

- `Button` with `primary`, `secondary`, and `danger` variants.
- `Field` or `TextField` with label, input, error slot, and focus styling.
- `Message` with `success`, `warning`, `error`, and `info` variants.
- `Metric`.
- `NavigationItem`.

This removes duplicated `login-*` and `signup-*` style blocks and makes future dashboard screens consistent.

### 5. Rebuild Auth Pages On The Template

Convert `LoginPage` and `SignupPage` from bespoke two-panel marketing layouts to quieter application screens:

- Use `Page` and `Panel` instead of full-screen gradient auth shells.
- Keep forms centered or constrained, but use preview button/input/status styles.
- Use the same `Field` component for email, password, and password confirmation.
- Use status/error color roles for field errors.
- Keep route links, API behavior, and token logic unchanged.

### 6. Rebuild Home Page As The First App Surface

Replace the current session card with an operational dashboard shell based on the preview:

- Page header with the current day and primary action.
- Main panel for habits.
- Optional side panel for states, stats, or account metadata.
- Bottom or local navigation modeled after `.navigation`.

### 7. Keep Page-Specific SASS Thin

After templates and primitives exist, page SASS should only handle page-specific composition.

Rules:

- No page-local color literals unless a new token is intentionally introduced.
- No repeated card/panel/input/button definitions in pages.
- No gradient page backgrounds.
- Keep BEM naming for page-specific classes and component classes.

### 8. Migration Order

Recommended order:

1. Add shared design tokens and global base styles.
2. Add layout templates.
3. Add button, field, and message primitives.
4. Convert `LoginPage`.
5. Convert `SignupPage`.
6. Convert `HomePage` into the initial dashboard shell.
7. Remove obsolete `App.css` startup styles and duplicated auth SASS blocks.

This order keeps API and routing behavior stable while replacing the visual layer incrementally.
