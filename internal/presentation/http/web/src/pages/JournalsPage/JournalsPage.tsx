import axios from 'axios'
import { useEffect, useMemo, useState } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { DateSelector, type DateRange } from '../../components/date'
import { GoogleIcon } from '../../components/icons'
import { AppNavigation } from '../../components/navigation'
import { useApiClient } from '../../hooks/useApiClient'
import './JournalsPage.sass'

type JournalsResponse = {
    journals: JournalResource[]
}

type JournalResource = {
    date: string
    note: string
    createdAt: string
    updatedAt: string
}

type JournalEntry = {
    id: string
    date: Date
    note: string
    createdAt: Date
    updatedAt: Date
}

type RangePreset = 7 | 15 | 30 | 90 | 365

const rangePresets: Array<{ label: string; days: RangePreset }> = [
    { label: '7 days', days: 7 },
    { label: '15 days', days: 15 },
    { label: '30 days', days: 30 },
    { label: '90 days', days: 90 },
    { label: '1 year', days: 365 },
]

const dateFormatter = new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
})

const weekdayFormatter = new Intl.DateTimeFormat('en-GB', {
    weekday: 'long',
})

const MS_PER_DAY = 24 * 60 * 60 * 1000

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

function buildRange(days: number, endDate: Date) {
    const end = startOfDay(endDate)
    return {
        startDate: addDays(end, -(days - 1)),
        endDate: end,
    }
}

function getRangeLength(range: DateRange) {
    const from = startOfDay(range.startDate)
    const to = startOfDay(range.endDate)
    return Math.floor((to.getTime() - from.getTime()) / MS_PER_DAY) + 1
}

function formatSelectedRange(range: DateRange) {
    return `${dateFormatter.format(range.startDate)} – ${dateFormatter.format(range.endDate)}`
}

function resolveRangeFromSearchParams(searchParams: URLSearchParams, fallbackRange: DateRange) {
    const from = parseDateKey(searchParams.get('from') ?? '')
    const to = parseDateKey(searchParams.get('to') ?? '')

    if (!from || !to) {
        return fallbackRange
    }

    if (from.getTime() > to.getTime()) {
        return fallbackRange
    }

    return {
        startDate: from,
        endDate: to,
    }
}

function sortJournalEntries(entries: JournalEntry[]) {
    return [...entries].sort((left, right) => {
        const dateDifference = right.date.getTime() - left.date.getTime()
        if (dateDifference !== 0) {
            return dateDifference
        }

        return right.createdAt.getTime() - left.createdAt.getTime()
    })
}

function mapJournalEntry(resource: JournalResource): JournalEntry | null {
    const date = parseDateKey(resource.date)
    if (!date) {
        return null
    }

    const createdAt = new Date(resource.createdAt)
    const updatedAt = new Date(resource.updatedAt)

    if (Number.isNaN(createdAt.getTime()) || Number.isNaN(updatedAt.getTime())) {
        return null
    }

    return {
        id: resource.date,
        date,
        note: resource.note,
        createdAt,
        updatedAt,
    }
}

function splitJournalNote(note: string) {
    const normalized = note.replace(/\r\n/g, '\n').replace(/\r/g, '\n')
    if (normalized.trim().length === 0) {
        return []
    }

    return normalized
        .split(/\n{2,}/)
        .map((paragraph) => paragraph.trim())
        .filter((paragraph) => paragraph.length > 0)
}

function renderJournalParagraph(paragraph: string, paragraphIndex: number) {
    return paragraph.split('\n').map((line, index) => (
        <span key={`${paragraphIndex}-${index}`}>
            {index > 0 ? <br /> : null}
            {line}
        </span>
    ))
}

// JournalsPage renders the journal archive as a reading-first page with preset and custom date ranges.
export function JournalsPage() {
    const apiClient = useApiClient()
    const today = useMemo(() => startOfDay(new Date()), [])
    const defaultRange = useMemo(() => buildRange(30, today), [today])
    const [searchParams, setSearchParams] = useSearchParams()
    const [isDateSelectorOpen, setDateSelectorOpen] = useState(false)
    const [journals, setJournals] = useState<JournalEntry[]>([])
    const [isLoading, setIsLoading] = useState(true)
    const [loadError, setLoadError] = useState('')
    const [scrollProgress, setScrollProgress] = useState(0)

    const range = useMemo(
        () => resolveRangeFromSearchParams(searchParams, defaultRange),
        [defaultRange, searchParams],
    )

    const rangeLength = getRangeLength(range)
    const activePreset = rangePresets.find((preset) => preset.days === rangeLength)
    const selectedRangeLabel = formatSelectedRange(range)
    const fromKey = formatDateKey(range.startDate)
    const toKey = formatDateKey(range.endDate)

    useEffect(() => {
        let isActive = true

        async function loadJournals() {
            setIsLoading(true)
            setLoadError('')

            try {
                const response = await apiClient.get<JournalsResponse>(
                    `journals?from=${fromKey}&to=${toKey}`,
                )

                if (!isActive) {
                    return
                }

                const mappedEntries = response.data.journals
                    .map(mapJournalEntry)
                    .filter((entry): entry is JournalEntry => entry !== null)

                setJournals(sortJournalEntries(mappedEntries))
            } catch (error) {
                if (!isActive) {
                    return
                }

                if (axios.isAxiosError(error) && error.response?.status === 401) {
                    setLoadError('Sign in to read journals.')
                } else {
                    setLoadError('Could not load journals.')
                }
            } finally {
                if (isActive) {
                    setIsLoading(false)
                }
            }
        }

        void loadJournals()

        return () => {
            isActive = false
        }
    }, [apiClient, fromKey, toKey])

    useEffect(() => {
        let frameId = 0

        function updateScrollProgress() {
            const scrollableHeight = document.documentElement.scrollHeight - window.innerHeight
            if (scrollableHeight <= 0) {
                setScrollProgress(1)
                return
            }

            const nextProgress = Math.min(1, Math.max(0, window.scrollY / scrollableHeight))
            setScrollProgress(nextProgress)
        }

        function scheduleScrollProgressUpdate() {
            window.cancelAnimationFrame(frameId)
            frameId = window.requestAnimationFrame(updateScrollProgress)
        }

        scheduleScrollProgressUpdate()

        window.addEventListener('scroll', scheduleScrollProgressUpdate, { passive: true })
        window.addEventListener('resize', scheduleScrollProgressUpdate, { passive: true })

        return () => {
            window.removeEventListener('scroll', scheduleScrollProgressUpdate)
            window.removeEventListener('resize', scheduleScrollProgressUpdate)
            window.cancelAnimationFrame(frameId)
        }
    }, [journals.length, isLoading, loadError, range])

    function handlePresetSelect(days: RangePreset) {
        const nextRange = buildRange(days, today)
        setSearchParams(
            {
                from: formatDateKey(nextRange.startDate),
                to: formatDateKey(nextRange.endDate),
            },
            { replace: true },
        )
        setDateSelectorOpen(false)
    }

    function handleCustomRangeChange(nextRange: DateRange) {
        const normalizedRange = {
            startDate: startOfDay(nextRange.startDate),
            endDate: startOfDay(nextRange.endDate),
        }

        setSearchParams(
            {
                from: formatDateKey(normalizedRange.startDate),
                to: formatDateKey(normalizedRange.endDate),
            },
            { replace: true },
        )
        setDateSelectorOpen(false)
    }

    return (
        <main className="journals-page">
            <div className="journals-page__progress" aria-hidden="true">
                <div
                    className="journals-page__progress-fill"
                    style={{ width: `${scrollProgress * 100}%` }}
                />
            </div>

            <div className="journals-page__shell">
                <header className="journals-page__header">
                    <p className="journals-page__eyebrow">Journals</p>
                    <h1 className="journals-page__title">Journals</h1>
                    <p className="journals-page__subtitle">
                        Reading view for the selected date range.
                    </p>
                </header>

                <div className="journals-page__controls" aria-label="Journal date ranges">
                    <div className="journals-page__preset-list">
                        {rangePresets.map((preset) => (
                            <button
                                aria-pressed={activePreset?.days === preset.days}
                                className={
                                    activePreset?.days === preset.days
                                        ? 'journals-page__preset-button journals-page__preset-button--active'
                                        : 'journals-page__preset-button'
                                }
                                key={preset.days}
                                type="button"
                                onClick={() => handlePresetSelect(preset.days)}
                            >
                                {preset.label}
                            </button>
                        ))}
                    </div>

                    <button
                        aria-pressed={!activePreset}
                        className={
                            !activePreset
                                ? 'journals-page__custom-button journals-page__custom-button--active'
                                : 'journals-page__custom-button'
                        }
                        type="button"
                        onClick={() => setDateSelectorOpen(true)}
                    >
                        Custom range
                    </button>
                </div>

                <p className="journals-page__range-label">{selectedRangeLabel}</p>

                {loadError ? (
                    <p className="journals-page__status journals-page__status--error">
                        {loadError}
                    </p>
                ) : journals.length === 0 ? (
                    <p className="journals-page__status">No journal entries in this range.</p>
                ) : null}

                <section className="journals-page__entries" aria-label="Journal entries">
                    {journals.map((journal) => (
                        <article className="journals-page__entry" id={journal.id} key={journal.id}>
                            <header className="journals-page__entry-header">
                                <a className="journals-page__date-link" href={`#${journal.id}`}>
                                    {dateFormatter.format(journal.date)}
                                </a>
                                <span className="journals-page__weekday">
                                    {weekdayFormatter.format(journal.date)}
                                </span>
                                <Link
                                    aria-label="Edit journal"
                                    className="journals-page__edit-link"
                                    to={`/dates/${journal.id}`}
                                >
                                    <GoogleIcon icon="edit_square" size={18} />
                                </Link>
                            </header>

                            <div className="journals-page__note">
                                {splitJournalNote(journal.note).map((paragraph, index) => (
                                    <p className="journals-page__paragraph" key={`${journal.id}-${index}`}>
                                        {renderJournalParagraph(paragraph, index)}
                                    </p>
                                ))}
                            </div>
                        </article>
                    ))}
                </section>
            </div>

            <AppNavigation />
            {isLoading ? (
                <div className="journals-page__loading-overlay" role="status" aria-live="polite">
                    <div className="journals-page__loading-box">
                        <div className="journals-page__loading-spinner" aria-hidden="true" />
                        <p className="journals-page__loading-text">Loading journals...</p>
                    </div>
                </div>
            ) : null}
            <DateSelector
                mode="range"
                open={isDateSelectorOpen}
                value={{
                    startDate: range.startDate,
                    endDate: range.endDate,
                }}
                onChange={handleCustomRangeChange}
                onClose={() => setDateSelectorOpen(false)}
            />
        </main>
    )
}
