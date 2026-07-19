import axios from 'axios'
import { useEffect, useMemo, useState } from 'react'
import { Navigate, useParams } from 'react-router-dom'
import { AppNavigation } from '../../components/navigation'
import { DailyMood, type MoodValue } from '../../components/mood'
import { DailyJournal } from '../../components/journal'
import {
    DailyHabits,
    type DailyHabitsData,
    type DailyHabitRecords,
} from '../../components/habits/DailyHabits'
import { Page } from '../../components/layout'
import { Message } from '../../components/primitives'
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
    const params = useParams<{ date: string }>()
    const rawDate = explicitDate ?? params.date
    const pageDate = useMemo(() => resolvePageDate(rawDate), [rawDate])
    const dateKey = useMemo(() => (pageDate ? formatDateKey(pageDate) : ''), [pageDate])
    const [habits, setHabits] = useState<HabitsResponse | null>(null)
    const [records, setRecords] = useState<RecordsResponse | null>(null)
    const [mood, setMood] = useState<MoodValue | null>(null)
    const [isLoading, setIsLoading] = useState(true)
    const [loadError, setLoadError] = useState('')

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
                const [habitsResponse, recordsResponse] = await Promise.all([
                    apiClient.get<HabitsResponse>('habits'),
                    apiClient.get<RecordsResponse>(`habits/${dateKey}`),
                ])
                const moodResponse = await apiClient.get<MoodResponse>(`moods/${dateKey}`).catch(
                    (error) => {
                        if (axios.isAxiosError(error) && error.response?.status === 404) {
                            return null
                        }

                        throw error
                    },
                )

                if (!isActive) {
                    return
                }

                setHabits(habitsResponse.data)
                setRecords(recordsResponse.data)
                setMood(moodResponse?.data.mood ?? null)
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

    return (
        <Page className="day-page">
            <div className="day-page__content">
                <DailyMood dateKey={dateKey} initialMood={mood} />
                <DailyJournal dateKey={dateKey} />

                {isLoading ? (
                    <Message variant="info">Loading day data...</Message>
                ) : loadError ? (
                    <Message variant="error">{loadError}</Message>
                ) : habits && records ? (
                    <DailyHabits date={pageDate} dateKey={dateKey} habits={habits} records={records} />
                ) : null}
            </div>

            <AppNavigation />
        </Page>
    )
}
