import axios from 'axios'
import { useEffect, useMemo, useState } from 'react'
import { useSearchParams } from 'react-router-dom'
import { DateSelector, type DateRange } from '../../components/date'
import { AppNavigation } from '../../components/navigation'
import { GoogleIcon } from '../../components/icons'
import { Message } from '../../components/primitives'
import { Page, PageHeader } from '../../components/layout'
import { useApiClient } from '../../hooks/useApiClient'
import './TransactionStatisticsPage.sass'

type OverviewResource = {
    expenses: number
    incomes: number
    netChange: number
}

type TransactionStatisticsResponse = {
    target: {
        overview: OverviewResource
    }
    baseline: {
        overview: OverviewResource
    }
}

type SummaryCard = {
    key: 'expenses' | 'incomes' | 'netChange'
    label: string
    icon: string
    accent: 'coral' | 'green' | 'gold'
    value: number
    baselineValue: number
}

const dateFormatter = new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'short',
    year: 'numeric',
})

const amountFormatter = new Intl.NumberFormat('ru-RU', {
    maximumFractionDigits: 2,
})

function startOfDay(date: Date) {
    const nextDate = new Date(date)
    nextDate.setHours(0, 0, 0, 0)
    return nextDate
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

function resolveRangeFromSearchParams(searchParams: URLSearchParams, fallbackRange: { from: Date; to: Date }) {
    const from = parseDateKey(searchParams.get('from') ?? '')
    const to = parseDateKey(searchParams.get('to') ?? '')

    if (!from || !to) {
        return fallbackRange
    }

    if (from.getTime() > to.getTime()) {
        return fallbackRange
    }

    return {
        from,
        to,
    }
}

function formatCurrency(value: number) {
    const normalized = amountFormatter.format(Math.abs(value))
    return `${normalized} ₽`
}

function formatSignedCurrency(value: number) {
    const normalized = amountFormatter.format(Math.abs(value))
    const sign = value < 0 ? '−' : '+'
    return `${sign}${normalized} ₽`
}

function formatPercentChange(value: number) {
    const rounded = Math.round(value)
    const sign = rounded > 0 ? '+' : rounded < 0 ? '−' : ''
    return `${sign}${amountFormatter.format(Math.abs(rounded))}%`
}

function buildSummaryCards(data: TransactionStatisticsResponse | null): SummaryCard[] {
    if (!data) {
        return [
            { key: 'expenses', label: 'Expenses', icon: 'south_east', accent: 'coral', value: 0, baselineValue: 0 },
            { key: 'incomes', label: 'Incomes', icon: 'north_east', accent: 'green', value: 0, baselineValue: 0 },
            { key: 'netChange', label: 'Net change', icon: 'compare_arrows', accent: 'gold', value: 0, baselineValue: 0 },
        ]
    }

    return [
        {
            key: 'expenses',
            label: 'Expenses',
            icon: 'south_east',
            accent: 'coral',
            value: data.target.overview.expenses,
            baselineValue: data.baseline.overview.expenses,
        },
        {
            key: 'incomes',
            label: 'Income',
            icon: 'north_east',
            accent: 'green',
            value: data.target.overview.incomes,
            baselineValue: data.baseline.overview.incomes,
        },
        {
            key: 'netChange',
            label: 'Net change',
            icon: 'compare_arrows',
            accent: 'gold',
            value: data.target.overview.netChange,
            baselineValue: data.baseline.overview.netChange,
        },
    ]
}

function calculatePercentDifference(value: number, baselineValue: number) {
    if (baselineValue === 0) {
        return value === 0 ? 0 : 100
    }

    return ((value - baselineValue) / baselineValue) * 100
}

function resolvePercentTone(cardKey: SummaryCard['key'], percentDifference: number) {
    if (cardKey === 'expenses') {
        if (percentDifference < 0) {
            return 'positive'
        }

        if (percentDifference > 0) {
            return 'negative'
        }

        return null
    }

    if (percentDifference > 0) {
        return 'positive'
    }

    if (percentDifference < 0) {
        return 'negative'
    }

    return null
}

function getRangeLabel(fromDate: Date, toDate: Date) {
    return `${dateFormatter.format(fromDate)} – ${dateFormatter.format(toDate)}`
}

// TransactionStatisticsPage renders the transaction statistics overview dashboard.
export function TransactionStatisticsPage() {
    const apiClient = useApiClient()
    const [searchParams, setSearchParams] = useSearchParams()
    const today = useMemo(() => startOfDay(new Date()), [])
    const defaultRange = useMemo(() => buildRange(30, today), [today])
    const range = useMemo(() => resolveRangeFromSearchParams(searchParams, defaultRange), [defaultRange, searchParams])
    const [statistics, setStatistics] = useState<TransactionStatisticsResponse | null>(null)
    const [isLoading, setIsLoading] = useState(true)
    const [loadError, setLoadError] = useState('')

    const summaryCards = useMemo(() => buildSummaryCards(statistics), [statistics])

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

        async function loadStatistics() {
            setIsLoading(true)
            setLoadError('')

            try {
                const response = await apiClient.get<TransactionStatisticsResponse>(
                    `transactions/statistics?from=${formatDateKey(range.from)}&to=${formatDateKey(range.to)}`,
                )

                if (!isActive) {
                    return
                }

                setStatistics(response.data)
            } catch (error) {
                if (!isActive) {
                    return
                }

                if (axios.isAxiosError(error) && error.response?.status === 400) {
                    setLoadError('Could not load transaction statistics.')
                    return
                }

                setLoadError('Could not load transaction statistics.')
            } finally {
                if (isActive) {
                    setIsLoading(false)
                }
            }
        }

        void loadStatistics()

        return () => {
            isActive = false
        }
    }, [apiClient, range.from, range.to])

    return (
        <Page className="transaction-statistics-page">
            <div className="transaction-statistics-page__shell">
                <PageHeader
                    eyebrow="Transactions"
                    title="Transaction statistics"
                    subtitle={`Overview for ${getRangeLabel(range.from, range.to)}. The first pass focuses on the three overview metrics only.`}
                    actions={
                        <DateSelector
                            className="transaction-statistics-page__date-selector"
                            mode="range"
                            value={{
                                startDate: range.from,
                                endDate: range.to,
                            }}
                            onChange={handleRangeChange}
                        />
                    }
                />

                {isLoading ? (
                    <Message variant="info">Loading transaction statistics...</Message>
                ) : loadError ? (
                    <Message variant="error">{loadError}</Message>
                ) : (
                    <div className="transaction-statistics-page__overview-grid">
                        {summaryCards.map((card) => {
                            const percentDifference = calculatePercentDifference(card.value, card.baselineValue)
                            const percentTone = resolvePercentTone(card.key, percentDifference)
                            const isNetChange = card.key === 'netChange'
                            const netChangeIsPositive = card.value >= 0

                            return (
                                <article
                                    className={`transaction-statistics-page__overview-card transaction-statistics-page__overview-card--${card.accent}`}
                                    key={card.key}
                                >
                                    <div className="transaction-statistics-page__card-head">
                                        <div className="transaction-statistics-page__icon-shell">
                                            <GoogleIcon icon={card.icon} size={14} />
                                        </div>
                                        <span className="transaction-statistics-page__card-title">
                                            {card.label}
                                        </span>
                                    </div>

                                    <strong
                                        className={[
                                            'transaction-statistics-page__card-value',
                                            isNetChange
                                                ? netChangeIsPositive
                                                    ? 'transaction-statistics-page__card-value--positive'
                                                    : 'transaction-statistics-page__card-value--negative'
                                                : undefined,
                                        ]
                                            .filter(Boolean)
                                            .join(' ')}
                                    >
                                        {isNetChange ? formatSignedCurrency(card.value) : formatCurrency(card.value)}
                                    </strong>

                                    <div
                                        className={[
                                            'transaction-statistics-page__card-delta',
                                        ]
                                            .filter(Boolean)
                                            .join(' ')}
                                    >
                                        <span
                                            className={[
                                                'transaction-statistics-page__card-delta-value',
                                                percentTone === 'positive'
                                                    ? 'transaction-statistics-page__card-delta-value--positive'
                                                    : undefined,
                                                percentTone === 'negative'
                                                    ? 'transaction-statistics-page__card-delta-value--negative'
                                                    : undefined,
                                            ]
                                                .filter(Boolean)
                                                .join(' ')}
                                        >
                                            {formatPercentChange(percentDifference)}
                                        </span>
                                        <span className="transaction-statistics-page__card-delta-label">
                                            vs. 3-period average
                                        </span>
                                    </div>
                                </article>
                            )
                        })}
                    </div>
                )}
            </div>

            <AppNavigation />
        </Page>
    )
}
