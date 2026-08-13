# Habit Statistics Page

## Summary

`HabitStatisticsPage` displays habit activity for a selected inclusive date range.

The page follows the visual structure of `TransactionStatisticsPage`: it uses the shared page header, period presets, range date selector, navigation, loading overlay, error message, and a full-width statistics panel.

The page loads two independent resources:

- `GET /api/v1/habits` supplies habit metadata and is converted into an in-memory registry keyed by habit ID.
- `GET /api/v1/habits/statistics?from=YYYY-MM-DD&to=YYYY-MM-DD` supplies heatmap values for the selected range.

The registry contains measurable, time, and completable habits. Only measurable habits are rendered at the moment because the current statistics response only contains `measurableHeatmap` data.

## Data invariants

- The habit registry is a `Map<string, HabitRegistryEntry>` keyed by the habit UUID.
- Every registry entry retains its habit type. Measurable entries additionally retain `step` and `unit`.
- Heatmaps are matched to registry entries by `habitId`.
- Habit metadata, including the label, icon, and measurement unit, always comes from the registry rather than the statistics response.
- A measurable habit remains renderable when its heatmap is missing; in that case it has no nodes and zero non-zero records.
- The displayed record count includes only nodes whose value is not zero.
- `maxValue` from the statistics response is the upper bound used to calculate node color intensity.
- The statistics request is repeated whenever the selected date range changes.
- The habits request is made once when the page is mounted.
- Habit-loading and statistics-loading states are tracked separately. The loading overlay remains visible until both requests have settled.
- Habit-loading and statistics-loading errors are tracked separately and combined for display without one request clearing the other request's error.

## Date range invariants

- The default range contains 30 calendar days and ends today.
- `from` and `to` are inclusive.
- The selected range is stored in the URL as `from` and `to` query parameters using `YYYY-MM-DD`.
- Invalid, missing, or reversed URL ranges fall back to the default range.
- `Week`, `Month`, and `Year` represent 7, 30, and 365 inclusive days respectively.
- `Custom` opens the shared range `DateSelector`.
- The period selector and date selector always operate on the same range state.

## Heatmap invariants

- Heatmaps are built with HTML elements, not canvas or SVG.
- Every node represents one calendar date.
- Nodes preserve chronological order from the API response.
- Monday is weekday index 1 and Sunday is weekday index 7.
- A partial first week starts in the row or column matching the actual weekday of the first date; it is never shifted to Monday artificially.
- Node tooltips display the full date, the value rounded to at most one fractional digit, and the measurable habit unit.
- Tooltips are positioned in a fixed layer outside the scrolling container so they are not clipped.
- An open tooltip closes on node blur, pointer leave, any scroll event, or window resize.
- Zero-value nodes use a neutral gray-green color.
- Positive nodes start with a light accent tint and interpolate toward the primary accent color according to `value / maxValue`.
- Color intensity is clamped to the range from 0 to 1.
- Nodes use rounded corners and a 2 px gap.

## Layout threshold invariants

The number of nodes determines the calendar orientation.

### 30 nodes or fewer

- The heatmap uses seven weekday columns ordered Monday through Sunday.
- Weekday initials are displayed above the columns.
- Calendar weeks progress from top to bottom.
- Month labels are not displayed.

### More than 30 nodes

- The heatmap uses seven weekday rows ordered Monday through Sunday.
- Weekday initials are displayed to the left of the rows.
- Calendar weeks progress from left to right.
- Short month names are displayed above the node grid.
- Each month label is aligned with the first week column containing a date from that month.

## Responsive invariants

Responsive behavior is determined by orientation rather than viewport width:

- `orientation: portrait` is the mobile layout.
- `orientation: landscape` is the desktop layout.

### Portrait

- Heatmaps with 30 nodes or fewer expand their seven columns to fill the available panel width.
- Short-layout nodes remain square through `aspect-ratio`.
- Heatmaps with more than 30 nodes retain fixed-size nodes and may be wider than their container.
- Long heatmaps scroll horizontally inside their own container and must never increase the width of the page.

### Landscape

- Nodes retain their fixed size and do not stretch to fill available space.
- Habits with 30 nodes or fewer use content-sized blocks and may share a row.
- Short habit blocks wrap onto another row when there is insufficient horizontal space.
- Habits with more than 30 nodes always occupy a full row.
- Long heatmaps remain horizontally scrollable within their habit block.
- Multiple short habit blocks are distributed horizontally with `space-between`.

## Overflow invariants

- The page shell, panel, habit blocks, and scroll containers use `min-width: 0` where required to prevent intrinsic grid width from expanding the page.
- The outer heatmap panel clips overflow.
- Only the dedicated heatmap scroll container may scroll horizontally.
- Fixed-size long heatmaps must not change the document width in either orientation.
