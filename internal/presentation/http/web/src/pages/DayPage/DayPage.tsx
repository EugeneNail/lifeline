import axios from 'axios'
import { useEffect, useMemo, useState } from 'react'
import { Navigate, useNavigate, useParams } from 'react-router-dom'
import { DateSelector } from '../../components/date'
import { AppNavigation } from '../../components/navigation'
import { DailyMood, type MoodValue } from '../../components/mood'
import { DailyJournal } from '../../components/journal'
import {
    DailyHabits,
    type DailyHabitsData,
    type DailyHabitRecords,
} from '../../components/habits/DailyHabits'
import { Page, PageHeader } from '../../components/layout'
import { Button, IconButton, Message } from '../../components/primitives'
import { GoogleIcon, GoogleIcons } from '../../components/icons'
import { useApiClient } from '../../hooks/useApiClient'
import './DayPage.sass'

type HabitsResponse = DailyHabitsData & {
    completable: Array<DailyHabitsData['completable'][number]>
    measurable: Array<DailyHabitsData['measurable'][number]>
    time: Array<DailyHabitsData['time'][number]>
}

type RecordsResponse = DailyHabitRecords

type MoodResponse = {
    mood: MoodValue
}

type JournalResponse = {
    note: string
}

function resolvePageDate(rawDate: string | undefined) {
    if (!rawDate) {
        return null
    }

    if (rawDate === 'today') {
        const today = new Date()
        today.setHours(0, 0, 0, 0)
        return today
    }

    const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(rawDate)
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

    return date
}

function formatDateKey(date: Date) {
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

type DayPageProps = {
    date?: string
}

// DayPage loads the selected day data and renders the daily habits workspace.
export function DayPage({ date: explicitDate }: DayPageProps) {
    const apiClient = useApiClient()
    const navigate = useNavigate()
    const params = useParams<{ date: string }>()
    const rawDate = explicitDate ?? params.date
    const pageDate = useMemo(() => resolvePageDate(rawDate), [rawDate])
    const dateKey = useMemo(() => (pageDate ? formatDateKey(pageDate) : ''), [pageDate])
    const [habits, setHabits] = useState<HabitsResponse | null>(null)
    const [records, setRecords] = useState<RecordsResponse | null>(null)
    const [mood, setMood] = useState<MoodValue | null>(null)
    const [journal, setJournal] = useState<string | null>(null)
    const [isDateSelectorOpen, setDateSelectorOpen] = useState(false)
    const [isLoading, setIsLoading] = useState(true)
    const [loadError, setLoadError] = useState('')
    const todayKey = useMemo(() => formatDateKey(startOfDay(new Date())), [])
    const isToday = dateKey === todayKey
    const pageTitle = pageDate ? formatPageTitle(pageDate) : ''
    const pageDateLabel = pageDate ? formatPageLabel(pageDate) : ''

    useEffect(() => {
        let isActive = true

        async function loadDayData() {
            if (!pageDate) {
                if (isActive) {
                    setIsLoading(false)
                }
                return
            }

            setIsLoading(true)
            setLoadError('')

            try {
                const [habitsResponse, recordsResponse, moodResponse, journalResponse] = await Promise.all([
                    apiClient.get<HabitsResponse>('habits'),
                    apiClient.get<RecordsResponse>(`habits/${dateKey}`),
                    apiClient.get<MoodResponse>(`moods/${dateKey}`).catch((error) => {
                        if (axios.isAxiosError(error) && error.response?.status === 404) {
                            return null
                        }

                        throw error
                    }),
                    apiClient.get<JournalResponse>(`journals/${dateKey}`).catch((error) => {
                        if (axios.isAxiosError(error) && error.response?.status === 404) {
                            return null
                        }

                        throw error
                    }),
                ])

                if (!isActive) {
                    return
                }

                setHabits(habitsResponse.data)
                setRecords(recordsResponse.data)
                setMood(moodResponse?.data.mood ?? null)
                setJournal(journalResponse?.data.note ?? null)
            } catch (error) {
                if (!isActive) {
                    return
                }

                if (axios.isAxiosError(error) && error.response?.status === 404) {
                    setLoadError('Could not load the selected day.')
                    return
                }

                setLoadError('Could not load the selected day.')
            } finally {
                if (isActive) {
                    setIsLoading(false)
                }
            }
        }

        void loadDayData()

        return () => {
            isActive = false
        }
    }, [apiClient, dateKey, pageDate])

    if (!pageDate) {
        return <Navigate replace to="/habits" />
    }

    const currentPageDate = pageDate

    function handleDateSelect(selectedDate: Date) {
        const selectedKey = formatDateKey(selectedDate)
        setDateSelectorOpen(false)

        if (selectedKey === todayKey) {
            navigate('/')
            return
        }

        navigate(`/dates/${selectedKey}`)
    }

    function handlePreviousDay() {
        navigateForDate(addDays(currentPageDate, -1), navigate, todayKey)
    }

    function handleNextDay() {
        navigateForDate(addDays(currentPageDate, 1), navigate, todayKey)
    }

    return (
        <Page className="day-page">
            <PageHeader
                eyebrow={isToday ? 'Today' : 'Day overview'}
                title={pageTitle}
                subtitle="Mood, journal, and habits for the selected day."
                actions={
                    <div className="day-page__header-actions">
                        <IconButton
                            aria-label="Previous day"
                            className="day-page__date-step-button"
                            type="button"
                            onClick={handlePreviousDay}
                        >
                            ←
                        </IconButton>

                        <Button
                            className="day-page__date-button"
                            variant="secondary"
                            type="button"
                            onClick={() => setDateSelectorOpen(true)}
                        >
                            <GoogleIcon icon={GoogleIcons.CalendarMonth} size={18} />
                            Change date
                        </Button>

                        <IconButton
                            aria-label="Next day"
                            className="day-page__date-step-button"
                            type="button"
                            onClick={handleNextDay}
                        >
                            →
                        </IconButton>
                    </div>
                }
            />

            <div className="day-page__content">
                <DailyMood dateKey={dateKey} dateLabel={pageDateLabel} initialMood={mood} />
                <DailyJournal dateKey={dateKey} dateLabel={pageDateLabel} initialNote={journal} />

                {isLoading ? (
                    <Message variant="info">Loading day data...</Message>
                ) : loadError ? (
                    <Message variant="error">{loadError}</Message>
                ) : habits && records ? (
                    <DailyHabits date={pageDate} dateKey={dateKey} habits={habits} records={records} />
                ) : null}
            </div>

            <AppNavigation />
            <DateSelector
                mode="single"
                open={isDateSelectorOpen}
                value={pageDate}
                onChange={handleDateSelect}
                onClose={() => setDateSelectorOpen(false)}
            />
        </Page>
    )
}

function formatPageTitle(date: Date) {
    return new Intl.DateTimeFormat('en-US', {
        weekday: 'long',
        day: 'numeric',
        month: 'long',
    }).format(date)
}

function formatPageLabel(date: Date) {
    return new Intl.DateTimeFormat('en-US', {
        weekday: 'long',
        month: 'long',
        day: 'numeric',
        year: 'numeric',
    }).format(date)
}

function startOfDay(date: Date) {
    const copy = new Date(date)
    copy.setHours(0, 0, 0, 0)

    return copy
}

function addDays(date: Date, delta: number) {
    const nextDate = new Date(date)
    nextDate.setDate(nextDate.getDate() + delta)
    nextDate.setHours(0, 0, 0, 0)

    return nextDate
}

function navigateForDate(date: Date, navigate: (path: string) => void, todayKey: string) {
    const selectedKey = formatDateKey(date)

    if (selectedKey === todayKey) {
        navigate('/')
        return
    }

    navigate(`/dates/${selectedKey}`)
}
