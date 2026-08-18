import type { CSSProperties, FocusEvent, MouseEvent } from 'react'
import { EmojiIcon, type GoogleIcons } from '../icons'

export type HabitHeatmapHabit = {
    id: string
    label: string
    icon: number
    type: 'measurable' | 'time' | 'completable'
    unit?: string
}

export type HabitHeatmapSeries = {
    nodes: Array<{ date: string; value: number }>
    maxValue: number
}

type HeatmapLayout = 'weekday-columns' | 'weekday-rows'

type HeatmapTooltipHandlers = {
    onShow: (element: HTMLElement, dateLabel: string, valueLabel: string) => void
    onHide: () => void
}

type HabitHeatmapProps = {
    habit: HabitHeatmapHabit
    series: HabitHeatmapSeries
    tooltip: HeatmapTooltipHandlers
}

type HeatmapMonthLabel = {
    column: number
    key: string
    label: string
}

const weekdayLetters = ['M', 'T', 'W', 'T', 'F', 'S', 'S']
const heatmapDateFormatter = new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
})
const heatmapValueFormatter = new Intl.NumberFormat('en-GB', {
    maximumFractionDigits: 1,
})
const heatmapMonthFormatter = new Intl.DateTimeFormat('en-GB', {
    month: 'short',
})

// HabitHeatmap renders one habit's daily values in a responsive weekday grid.
export function HabitHeatmap({ habit, series, tooltip }: HabitHeatmapProps) {
    const nodes = series.nodes
    const nonZeroRecords = nodes.filter((node) => node.value !== 0).length
    const firstDateKey = nodes[0]?.date ?? ''
    const heatmapLayout: HeatmapLayout = nodes.length <= 30
        ? 'weekday-columns'
        : 'weekday-rows'
    const lastNodePosition = nodes.length > 0
        ? getHeatmapPosition(
            nodes[nodes.length - 1].date,
            nodes.length - 1,
            firstDateKey,
            heatmapLayout,
        )
        : null
    const heatmapColumnCount = heatmapLayout === 'weekday-columns'
        ? 7
        : Math.max(lastNodePosition?.column ?? 1, 1)
    const monthLabels = heatmapLayout === 'weekday-rows'
        ? buildHeatmapMonthLabels(nodes, firstDateKey)
        : []

    return (
        <article className={`habit-statistics-page__habit habit-statistics-page__habit--${heatmapLayout}`}>
            <header className="habit-statistics-page__habit-header">
                <div className="habit-statistics-page__habit-icon">
                    <EmojiIcon icon={habit.icon as GoogleIcons} size={22} />
                </div>
                <div className="habit-statistics-page__habit-copy">
                    <h2>{habit.label}</h2>
                    <p>{nonZeroRecords} non-zero records</p>
                </div>
            </header>

            <div className="habit-statistics-page__heatmap-scroll">
                <div
                    className={`habit-statistics-page__heatmap-frame habit-statistics-page__heatmap-frame--${heatmapLayout}`}
                    style={{ '--heatmap-columns': heatmapColumnCount } as CSSProperties}
                >
                    <div
                        aria-hidden="true"
                        className={`habit-statistics-page__weekday-labels habit-statistics-page__weekday-labels--${heatmapLayout}`}
                    >
                        {weekdayLetters.map((letter, index) => (
                            <span key={`${letter}-${index}`}>{letter}</span>
                        ))}
                    </div>
                    <div
                        aria-label={`${habit.label} heatmap`}
                        className={`habit-statistics-page__heatmap habit-statistics-page__heatmap--${heatmapLayout}`}
                        role="img"
                    >
                        {nodes.map((node, index) => {
                            const position = getHeatmapPosition(
                                node.date,
                                index,
                                firstDateKey,
                                heatmapLayout,
                            )
                            const valueLabel = getHeatmapValueLabel(
                                node.value,
                                habit.type,
                                habit.unit,
                            )
                            const showTooltip = (event: MouseEvent<HTMLSpanElement> | FocusEvent<HTMLSpanElement>) => {
                                tooltip.onShow(
                                    event.currentTarget,
                                    formatHeatmapDate(node.date),
                                    valueLabel,
                                )
                            }

                            return (
                                <span
                                    aria-label={`${node.date}: ${valueLabel}`}
                                    className="habit-statistics-page__heatmap-node"
                                    key={node.date}
                                    style={{
                                        backgroundColor: getHeatmapNodeColor(
                                            node.value,
                                            series.maxValue,
                                            habit.type,
                                        ),
                                        gridColumn: position.column,
                                        gridRow: position.row,
                                    }}
                                    tabIndex={0}
                                    onBlur={tooltip.onHide}
                                    onFocus={showTooltip}
                                    onMouseEnter={showTooltip}
                                    onMouseLeave={tooltip.onHide}
                                />
                            )
                        })}
                    </div>
                    {heatmapLayout === 'weekday-rows' ? (
                        <div
                            aria-hidden="true"
                            className="habit-statistics-page__month-labels"
                            style={{ gridTemplateColumns: `repeat(${heatmapColumnCount}, 22px)` }}
                        >
                            {monthLabels.map((month) => (
                                <span key={month.key} style={{ gridColumn: month.column }}>
                                    {month.label}
                                </span>
                            ))}
                        </div>
                    ) : null}
                </div>
            </div>
        </article>
    )
}

function parseDateKey(dateKey: string) {
    const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(dateKey)
    if (!match) {
        return null
    }

    const year = Number(match[1])
    const month = Number(match[2]) - 1
    const day = Number(match[3])
    const date = new Date(year, month, day)

    if (date.getFullYear() !== year || date.getMonth() !== month || date.getDate() !== day) {
        return null
    }

    date.setHours(0, 0, 0, 0)
    return date
}

function formatHeatmapDate(dateKey: string) {
    const date = parseDateKey(dateKey)
    return date ? heatmapDateFormatter.format(date) : dateKey
}

function getHeatmapPosition(
    dateKey: string,
    index: number,
    firstDateKey: string,
    layout: HeatmapLayout,
) {
    const firstDate = parseDateKey(firstDateKey)
    const date = parseDateKey(dateKey)

    if (!firstDate || !date) {
        const week = Math.floor(index / 7) + 1
        const weekday = (index % 7) + 1
        return layout === 'weekday-columns'
            ? { column: weekday, row: week }
            : { column: week, row: weekday }
    }

    const firstDayOffset = (firstDate.getDay() + 6) % 7
    const dayOffset = Math.round((Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()) - Date.UTC(firstDate.getFullYear(), firstDate.getMonth(), firstDate.getDate())) / (24 * 60 * 60 * 1000))
    const position = firstDayOffset + dayOffset
    const week = Math.floor(position / 7) + 1
    const weekday = (position % 7) + 1

    return layout === 'weekday-columns'
        ? { column: weekday, row: week }
        : { column: week, row: weekday }
}

function buildHeatmapMonthLabels(nodes: HabitHeatmapSeries['nodes'], firstDateKey: string) {
    const labels: HeatmapMonthLabel[] = []
    let previousMonthKey = ''

    nodes.forEach((node, index) => {
        const date = parseDateKey(node.date)
        if (!date) {
            return
        }

        const monthKey = `${date.getFullYear()}-${date.getMonth()}`
        if (monthKey === previousMonthKey) {
            return
        }

        previousMonthKey = monthKey
        labels.push({
            column: getHeatmapPosition(node.date, index, firstDateKey, 'weekday-rows').column,
            key: monthKey,
            label: heatmapMonthFormatter.format(date),
        })
    })

    return labels
}

function interpolateHeatmapColor(value: number, maxValue: number) {
    const emptyColor = [222, 227, 220]
    const nonZeroBaseColor = [190, 208, 199]
    const accentColor = [66, 104, 88]

    if (value <= 0 || maxValue <= 0) {
        return `rgb(${emptyColor.join(', ')})`
    }

    const intensity = Math.min(Math.max(value / maxValue, 0), 1)
    const channels = nonZeroBaseColor.map((channel, index) =>
        Math.round(channel + (accentColor[index] - channel) * intensity),
    )

    return `rgb(${channels.join(', ')})`
}

function getHeatmapNodeColor(value: number, maxValue: number, habitType: HabitHeatmapHabit['type']) {
    if (habitType === 'completable') {
        return interpolateHeatmapColor(value > 0 ? 1 : 0, 1)
    }

    return interpolateHeatmapColor(value, maxValue)
}

function formatTimeValue(value: number) {
    const normalizedMinutes = Math.min(Math.max(Math.round(value), 0), 1439)
    return `${String(Math.floor(normalizedMinutes / 60)).padStart(2, '0')}:${String(normalizedMinutes % 60).padStart(2, '0')}`
}

function getHeatmapValueLabel(value: number, habitType: HabitHeatmapHabit['type'], unit?: string) {
    if (habitType === 'completable') {
        return value > 0 ? 'Completed' : 'Not completed'
    }

    if (habitType === 'time') {
        return formatTimeValue(value)
    }

    return `${heatmapValueFormatter.format(value)} ${unit ?? ''}`.trim()
}
