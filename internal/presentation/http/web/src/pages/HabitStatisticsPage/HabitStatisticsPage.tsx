import { useEffect, useMemo, useRef, useState, type CSSProperties } from 'react'
import { useSearchParams } from 'react-router-dom'
import { DateSelector, PeriodSelector, type DateRange, type DateSelectorHandle } from '../../components/date'
import { EmojiIcon, type GoogleIcons } from '../../components/icons'
import { Page, PageHeader } from '../../components/layout'
import { AppNavigation } from '../../components/navigation'
import { Message } from '../../components/primitives'
import { useApiClient } from '../../hooks/useApiClient'
import './HabitStatisticsPage.sass'

type HeatmapNodeResource = {
    date: string
    value: number
}

type MeasurableHabitHeatmapResource = {
    habitId: string
    nodes: HeatmapNodeResource[]
    maxValue: number
}

type HabitStatisticsResponse = {
    measurableHeatmap: MeasurableHabitHeatmapResource[]
}

type HabitResource = {
    id: string
    label: string
    icon: number
    archivedAt: string | null
}

type MeasurableHabitResource = HabitResource & {
    step: number
    unit: string
}

type HabitsResponse = {
    measurable: MeasurableHabitResource[]
    time: HabitResource[]
    completable: HabitResource[]
}

type HabitRegistryEntry = HabitResource & {
    type: 'measurable' | 'time' | 'completable'
    step?: number
    unit?: string
}

type HeatmapTooltip = {
    dateLabel: string
    left: number
    top: number
    valueLabel: string
}

type HeatmapLayout = 'weekday-columns' | 'weekday-rows'

type HeatmapMonthLabel = {
    column: number
    key: string
    label: string
}

const dateFormatter = new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
})

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

const weekdayLetters = ['M', 'T', 'W', 'T', 'F', 'S', 'S']

function startOfDay(date: Date) {
    const normalizedDate = new Date(date)
    normalizedDate.setHours(0, 0, 0, 0)
    return normalizedDate
}

function addDays(date: Date, days: number) {
    const nextDate = new Date(date)
    nextDate.setDate(nextDate.getDate() + days)
    return startOfDay(nextDate)
}

function buildRange(days: number, endDate: Date) {
    const end = startOfDay(endDate)
    return {
        from: addDays(end, -(days - 1)),
        to: end,
    }
}

function formatDateKey(date: Date) {
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
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

    if (
        date.getFullYear() !== year ||
        date.getMonth() !== month ||
        date.getDate() !== day
    ) {
        return null
    }

    return startOfDay(date)
}

function resolveRangeFromSearchParams(
    searchParams: URLSearchParams,
    fallbackRange: { from: Date; to: Date },
) {
    const from = parseDateKey(searchParams.get('from') ?? '')
    const to = parseDateKey(searchParams.get('to') ?? '')

    if (!from || !to || from.getTime() > to.getTime()) {
        return fallbackRange
    }

    return { from, to }
}

function getRangeLabel(fromDate: Date, toDate: Date) {
    return `${dateFormatter.format(fromDate)} – ${dateFormatter.format(toDate)}`
}

function formatHeatmapDate(dateKey: string) {
    const date = parseDateKey(dateKey)
    return date ? heatmapDateFormatter.format(date) : dateKey
}

function buildHabitRegistry(response: HabitsResponse) {
    const registry = new Map<string, HabitRegistryEntry>()

    response.measurable.forEach((habit) => {
        registry.set(habit.id, { ...habit, type: 'measurable' })
    })
    response.time.forEach((habit) => {
        registry.set(habit.id, { ...habit, type: 'time' })
    })
    response.completable.forEach((habit) => {
        registry.set(habit.id, { ...habit, type: 'completable' })
    })

    return registry
}

function getHeatmapPosition(
    dateKey: string,
    index: number,
    firstDateKey: string,
    layout: HeatmapLayout,
) {
    const firstDate = parseDateKey(firstDateKey)
    const date = parseDateKey(dateKey)
    const fallbackPosition = index

    if (!firstDate || !date) {
        const week = Math.floor(fallbackPosition / 7) + 1
        const weekday = (fallbackPosition % 7) + 1
        return layout === 'weekday-columns'
            ? { column: weekday, row: week }
            : { column: week, row: weekday }
    }

    const firstDayOffset = (firstDate.getDay() + 6) % 7
    const dayOffset = Math.round(
        (Date.UTC(date.getFullYear(), date.getMonth(), date.getDate()) -
            Date.UTC(firstDate.getFullYear(), firstDate.getMonth(), firstDate.getDate())) /
            (24 * 60 * 60 * 1000),
    )
    const position = firstDayOffset + dayOffset

    const week = Math.floor(position / 7) + 1
    const weekday = (position % 7) + 1

    return layout === 'weekday-columns'
        ? { column: weekday, row: week }
        : { column: week, row: weekday }
}

function buildHeatmapMonthLabels(
    nodes: HeatmapNodeResource[],
    firstDateKey: string,
) {
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

    const intensity = maxValue > 0 ? Math.min(Math.max(value / maxValue, 0), 1) : 0
    const channels = nonZeroBaseColor.map((channel, index) =>
        Math.round(channel + (accentColor[index] - channel) * intensity),
    )

    return `rgb(${channels.join(', ')})`
}

// HabitStatisticsPage renders the habit statistics controls and loads data for the selected range.
export function HabitStatisticsPage() {
    const apiClient = useApiClient()
    const [searchParams, setSearchParams] = useSearchParams()
    const today = useMemo(() => startOfDay(new Date()), [])
    const defaultRange = useMemo(() => buildRange(30, today), [today])
    const range = useMemo(
        () => resolveRangeFromSearchParams(searchParams, defaultRange),
        [defaultRange, searchParams],
    )
    const [statistics, setStatistics] = useState<HabitStatisticsResponse | null>(null)
    const [habitRegistry, setHabitRegistry] = useState<Map<string, HabitRegistryEntry>>(new Map())
    const [isStatisticsLoading, setIsStatisticsLoading] = useState(true)
    const [areHabitsLoading, setAreHabitsLoading] = useState(true)
    const [statisticsLoadError, setStatisticsLoadError] = useState('')
    const [habitsLoadError, setHabitsLoadError] = useState('')
    const [heatmapTooltip, setHeatmapTooltip] = useState<HeatmapTooltip | null>(null)
    const dateSelectorRef = useRef<DateSelectorHandle>(null)
    const measurableHabits = useMemo(
        () => [...habitRegistry.values()].filter((habit) => habit.type === 'measurable'),
        [habitRegistry],
    )
    const heatmapsByHabitID = useMemo(
        () => new Map(
            (statistics?.measurableHeatmap ?? []).map((heatmap) => [heatmap.habitId, heatmap]),
        ),
        [statistics],
    )
    const isLoading = isStatisticsLoading || areHabitsLoading
    const loadError = habitsLoadError || statisticsLoadError

    useEffect(() => {
        if (!heatmapTooltip) {
            return
        }

        function hideHeatmapTooltip() {
            setHeatmapTooltip(null)
        }

        window.addEventListener('scroll', hideHeatmapTooltip, true)
        window.addEventListener('resize', hideHeatmapTooltip)

        return () => {
            window.removeEventListener('scroll', hideHeatmapTooltip, true)
            window.removeEventListener('resize', hideHeatmapTooltip)
        }
    }, [heatmapTooltip])

    function showHeatmapTooltip(
        element: HTMLElement,
        dateLabel: string,
        valueLabel: string,
    ) {
        const bounds = element.getBoundingClientRect()
        const horizontalPadding = 70

        setHeatmapTooltip({
            dateLabel,
            left: Math.min(
                Math.max(bounds.left + bounds.width / 2, horizontalPadding),
                window.innerWidth - horizontalPadding,
            ),
            top: bounds.top - 8,
            valueLabel,
        })
    }

    function handleRangeChange(nextRange: DateRange) {
        setSearchParams(
            {
                from: formatDateKey(startOfDay(nextRange.startDate)),
                to: formatDateKey(startOfDay(nextRange.endDate)),
            },
            { replace: true },
        )
    }

    useEffect(() => {
        const fromKey = formatDateKey(range.from)
        const toKey = formatDateKey(range.to)

        if (searchParams.get('from') === fromKey && searchParams.get('to') === toKey) {
            return
        }

        setSearchParams(
            {
                from: fromKey,
                to: toKey,
            },
            { replace: true },
        )
    }, [range.from, range.to, searchParams, setSearchParams])

    useEffect(() => {
        let isActive = true

        async function loadHabits() {
            setAreHabitsLoading(true)
            setHabitsLoadError('')

            try {
                const response = await apiClient.get<HabitsResponse>('habits')
                if (!isActive) {
                    return
                }

                setHabitRegistry(buildHabitRegistry(response.data))
            } catch {
                if (!isActive) {
                    return
                }

                setHabitsLoadError('Could not load habits.')
            } finally {
                if (isActive) {
                    setAreHabitsLoading(false)
                }
            }
        }

        void loadHabits()

        return () => {
            isActive = false
        }
    }, [apiClient])

    useEffect(() => {
        let isActive = true

        async function loadStatistics() {
            setIsStatisticsLoading(true)
            setStatisticsLoadError('')

            try {
                const requestParams = new URLSearchParams({
                    from: formatDateKey(range.from),
                    to: formatDateKey(range.to),
                })
                const response = await apiClient.get<HabitStatisticsResponse>(
                    `habits/statistics?${requestParams.toString()}`,
                )

                if (!isActive) {
                    return
                }

                setStatistics(response.data)
            } catch {
                if (!isActive) {
                    return
                }

                setStatisticsLoadError('Could not load habit statistics.')
            } finally {
                if (isActive) {
                    setIsStatisticsLoading(false)
                }
            }
        }

        void loadStatistics()

        return () => {
            isActive = false
        }
    }, [apiClient, range.from, range.to])

    return (
        <Page className="habit-statistics-page">
            <div className="habit-statistics-page__shell">
                <PageHeader
                    eyebrow="Habits"
                    title="Habit statistics"
                    subtitle={`Progress and consistency overview for ${getRangeLabel(range.from, range.to)}.`}
                />

                <div className="habit-statistics-page__filters">
                    <PeriodSelector
                        ariaLabel="Habit statistics period"
                        today={today}
                        value={{ startDate: range.from, endDate: range.to }}
                        onChange={handleRangeChange}
                        onCustomClick={() => dateSelectorRef.current?.open()}
                    />
                    <DateSelector
                        className="habit-statistics-page__date-selector"
                        mode="range"
                        ref={dateSelectorRef}
                        value={{
                            startDate: range.from,
                            endDate: range.to,
                        }}
                        onChange={handleRangeChange}
                    />
                </div>

                {loadError ? <Message variant="error">{loadError}</Message> : null}

                {!loadError && !isLoading ? (
                    <section className="habit-statistics-page__heatmap-panel">
                        {measurableHabits.map((habit) => {
                            const heatmap = heatmapsByHabitID.get(habit.id)
                            const nodes = heatmap?.nodes ?? []
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
                                <article
                                    className={`habit-statistics-page__habit habit-statistics-page__habit--${heatmapLayout}`}
                                    key={habit.id}
                                >
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
                                            style={{
                                                '--heatmap-columns': heatmapColumnCount,
                                            } as CSSProperties}
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
                                                    const valueLabel = `${heatmapValueFormatter.format(node.value)} ${habit.unit ?? ''}`.trim()

                                                    return (
                                                        <span
                                                            aria-label={`${node.date}: ${valueLabel}`}
                                                            className="habit-statistics-page__heatmap-node"
                                                            key={node.date}
                                                            style={{
                                                                backgroundColor: interpolateHeatmapColor(
                                                                    node.value,
                                                                    heatmap?.maxValue ?? 0,
                                                                ),
                                                                gridColumn: position.column,
                                                                gridRow: position.row,
                                                            }}
                                                            tabIndex={0}
                                                            onBlur={() => setHeatmapTooltip(null)}
                                                            onFocus={(event) => showHeatmapTooltip(
                                                                event.currentTarget,
                                                                formatHeatmapDate(node.date),
                                                                valueLabel,
                                                            )}
                                                            onMouseEnter={(event) => showHeatmapTooltip(
                                                                event.currentTarget,
                                                                formatHeatmapDate(node.date),
                                                                valueLabel,
                                                            )}
                                                            onMouseLeave={() => setHeatmapTooltip(null)}
                                                        />
                                                    )
                                                })}
                                            </div>
                                            {heatmapLayout === 'weekday-rows' ? (
                                                <div
                                                    aria-hidden="true"
                                                    className="habit-statistics-page__month-labels"
                                                    style={{
                                                        gridTemplateColumns: `repeat(${heatmapColumnCount}, 22px)`,
                                                    }}
                                                >
                                                    {monthLabels.map((month) => (
                                                        <span
                                                            key={month.key}
                                                            style={{ gridColumn: month.column }}
                                                        >
                                                            {month.label}
                                                        </span>
                                                    ))}
                                                </div>
                                            ) : null}
                                        </div>
                                    </div>
                                </article>
                            )
                        })}
                    </section>
                ) : null}
            </div>

            <AppNavigation />
            {heatmapTooltip ? (
                <div
                    className="habit-statistics-page__heatmap-tooltip"
                    role="tooltip"
                    style={{
                        left: heatmapTooltip.left,
                        top: heatmapTooltip.top,
                    }}
                >
                    <strong>{heatmapTooltip.dateLabel}</strong>
                    <span>{heatmapTooltip.valueLabel}</span>
                </div>
            ) : null}
            {isLoading ? (
                <div
                    aria-live="polite"
                    className="habit-statistics-page__loading-overlay"
                    role="status"
                >
                    <div className="habit-statistics-page__loading-box">
                        <div className="habit-statistics-page__loading-spinner" aria-hidden="true" />
                        <p className="habit-statistics-page__loading-text">
                            Loading habit statistics...
                        </p>
                    </div>
                </div>
            ) : null}
        </Page>
    )
}
