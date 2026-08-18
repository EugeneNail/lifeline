import { useEffect, useMemo, useRef, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { DateSelector, PeriodSelector, type DateRange, type DateSelectorHandle } from '../../components/date'
import { HabitChart, HabitHeatmap, type HabitChartHabitType } from '../../components/charts'
import { EmojiIcon, type GoogleIcons } from '../../components/icons'
import { GoogleIcon } from '../../components/icons'
import { Page, PageHeader } from '../../components/layout'
import { AppNavigation } from '../../components/navigation'
import { Message } from '../../components/primitives'
import { useApiClient } from '../../hooks/useApiClient'
import './HabitStatisticsPage.sass'

type NodeResource = {
    date: string
    value: number
}

type MeasurableHabitSeriesResource = {
    habitId: string
    nodes: NodeResource[]
    minValue: number
    maxValue: number
}

type TimeHabitSeriesResource = {
    habitId: string
    nodes: NodeResource[]
    minValue: number
    maxValue: number
}

type CompletableSeriesResource = {
    habitId: string
    nodes: NodeResource[]
    minValue: number
    maxValue: number
}

type HabitSeriesResource = {
    habitId: string
    nodes: NodeResource[]
    minValue: number
    maxValue: number
}

type HabitStatisticsResponse = {
    measurableSeries: MeasurableHabitSeriesResource[]
    timeSeries: TimeHabitSeriesResource[]
    completableSeries: CompletableSeriesResource[]
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

const dateFormatter = new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
})

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
    const [areHeatmapsExpanded, setAreHeatmapsExpanded] = useState(true)
    const [areChartsExpanded, setAreChartsExpanded] = useState(true)
    const [heatmapTooltip, setHeatmapTooltip] = useState<HeatmapTooltip | null>(null)
    const dateSelectorRef = useRef<DateSelectorHandle>(null)
    const statisticsHabits = useMemo(() => {
        const habits = [...habitRegistry.values()]

        return [
            ...habits.filter((habit) => habit.type === 'completable'),
            ...habits.filter((habit) => habit.type === 'measurable'),
            ...habits.filter((habit) => habit.type === 'time'),
        ]
    }, [habitRegistry])
    const chartHabits = statisticsHabits
    const heatmapsByHabitID = useMemo(() => {
        const heatmaps = new Map<string, HabitSeriesResource>()

        for (const heatmap of statistics?.completableSeries ?? []) {
            heatmaps.set(heatmap.habitId, {
                ...heatmap,
                minValue: 0,
                maxValue: 1,
            })
        }
        for (const heatmap of statistics?.measurableSeries ?? []) {
            heatmaps.set(heatmap.habitId, heatmap)
        }
        for (const heatmap of statistics?.timeSeries ?? []) {
            heatmaps.set(heatmap.habitId, heatmap)
        }

        return heatmaps
    }, [statistics])
    const isLoading = isStatisticsLoading || areHabitsLoading
    const loadError = habitsLoadError || statisticsLoadError

    function getHabitTypeLabel(type: HabitRegistryEntry['type']) {
        if (type === 'completable') {
            return 'Completable habit'
        }

        if (type === 'measurable') {
            return 'Measurable habit'
        }

        return 'Time habit'
    }

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
                    <>
                        <button
                            aria-expanded={areHeatmapsExpanded}
                            className="habit-statistics-page__heatmaps-heading"
                            type="button"
                            onClick={() => setAreHeatmapsExpanded((expanded) => !expanded)}
                        >
                            <div>
                                <p className="habit-statistics-page__charts-eyebrow">ACTIVITY OVER TIME</p>
                                <h2>Habit heatmaps</h2>
                            </div>
                            <GoogleIcon icon="chevron_right" size={20} />
                        </button>
                        {areHeatmapsExpanded ? (
                            <section className="habit-statistics-page__heatmap-panel">
                                {statisticsHabits.map((habit) => {
                                const heatmap = heatmapsByHabitID.get(habit.id)
                                return (
                                    <HabitHeatmap
                                        habit={habit}
                                        key={habit.id}
                                        series={{
                                            nodes: heatmap?.nodes ?? [],
                                            maxValue: heatmap?.maxValue ?? 0,
                                        }}
                                        tooltip={{
                                            onHide: () => setHeatmapTooltip(null),
                                            onShow: showHeatmapTooltip,
                                        }}
                                    />
                                )
                                })}
                            </section>
                        ) : null}

                        <section className="habit-statistics-page__charts-panel">
                            <button
                                aria-expanded={areChartsExpanded}
                                className="habit-statistics-page__charts-heading"
                                type="button"
                                onClick={() => setAreChartsExpanded((expanded) => !expanded)}
                            >
                                <div>
                                    <p className="habit-statistics-page__charts-eyebrow">TRENDS OVER TIME</p>
                                    <h2>Habit charts</h2>
                                </div>
                                <GoogleIcon icon="chevron_right" size={20} />
                            </button>

                            {areChartsExpanded ? (
                                <div className="habit-statistics-page__charts-list">
                                {chartHabits.map((habit) => {
                                    const series = heatmapsByHabitID.get(habit.id) ?? {
                                        habitId: habit.id,
                                        nodes: [],
                                        minValue: 0,
                                        maxValue: 0,
                                    }
                                    return (
                                    <div
                                        className={series.nodes.length <= 90
                                            ? 'habit-statistics-page__chart-item habit-statistics-page__chart-item--grid'
                                            : 'habit-statistics-page__chart-item habit-statistics-page__chart-item--wide'}
                                        key={habit.id}
                                    >
                                            <div className="habit-statistics-page__chart-trigger">
                                                <span className="habit-statistics-page__chart-icon">
                                                    <EmojiIcon icon={habit.icon as GoogleIcons} size={20} />
                                                </span>
                                                <span className="habit-statistics-page__chart-copy">
                                                    <strong>{habit.label}</strong>
                                                    <small>{getHabitTypeLabel(habit.type)}</small>
                                                </span>
                                            </div>
                                            <div className="habit-statistics-page__chart-content">
                                                <HabitChart
                                                    habitType={habit.type as HabitChartHabitType}
                                                    series={series}
                                                />
                                            </div>
                                        </div>
                                    )
                                })}
                                </div>
                            ) : null}
                        </section>
                    </>
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
