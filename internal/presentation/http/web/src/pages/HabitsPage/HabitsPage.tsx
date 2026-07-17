import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { HabitCard } from '../../components/habits'
import { Page, PageHeader, Panel, PanelBody } from '../../components/layout'
import { AppNavigation } from '../../components/navigation'
import { GoogleIcons } from '../../components/icons'
import { Message } from '../../components/primitives'
import { useApiClient } from '../../hooks/useApiClient'
import type { CompletableHabit } from './CompletableHabit'
import type { MeasurableHabit } from './MeasurableHabit'
import type { TimeHabit } from './TimeHabit'
import './HabitsPage.sass'

type HabitsResponse = {
    measurable: MeasurableHabit[]
    time: TimeHabit[]
    completable: CompletableHabit[]
}

type HabitCardModel = {
    uuid: string
    type: 'completable' | 'measurable' | 'time'
    icon: GoogleIcons
    title: string
    typeLabel: string
    enabled: boolean
}

function mapHabitsResponse(response: HabitsResponse) {
    return [
        ...response.measurable.map(mapMeasurableHabit),
        ...response.time.map(mapTimeHabit),
        ...response.completable.map(mapCompletableHabit),
    ]
}

function mapMeasurableHabit(habit: MeasurableHabit): HabitCardModel {
    return {
        uuid: habit.id,
        type: 'measurable',
        icon: habit.icon as GoogleIcons,
        title: habit.label,
        typeLabel: 'Quantity',
        enabled: habit.archivedAt === null,
    }
}

function mapTimeHabit(habit: TimeHabit): HabitCardModel {
    return {
        uuid: habit.id,
        type: 'time',
        icon: habit.icon as GoogleIcons,
        title: habit.label,
        typeLabel: 'Time',
        enabled: habit.archivedAt === null,
    }
}

function mapCompletableHabit(habit: CompletableHabit): HabitCardModel {
    return {
        uuid: habit.id,
        type: 'completable',
        icon: habit.icon as GoogleIcons,
        title: habit.label,
        typeLabel: 'Completion',
        enabled: habit.archivedAt === null,
    }
}

// HabitsPage renders the habit management dashboard with data loaded from the API.
export function HabitsPage() {
    const apiClient = useApiClient()
    const [habits, setHabits] = useState<HabitCardModel[]>([])
    const [isLoading, setIsLoading] = useState(true)
    const [loadError, setLoadError] = useState('')

    useEffect(() => {
        let isActive = true

        async function loadHabits() {
            setIsLoading(true)
            setLoadError('')

            try {
                const response = await apiClient.get<HabitsResponse>('habits')
                if (!isActive) {
                    return
                }

                setHabits(mapHabitsResponse(response.data))
            } catch {
                if (!isActive) {
                    return
                }

                setLoadError('Could not load habits.')
            } finally {
                if (isActive) {
                    setIsLoading(false)
                }
            }
        }

        void loadHabits()

        return () => {
            isActive = false
        }
    }, [apiClient])

    function handleToggle(uuid: string, nextEnabled: boolean) {
        setHabits((current) =>
            current.map((habit) =>
                habit.uuid === uuid ? { ...habit, enabled: nextEnabled } : habit,
            ),
        )
    }

    return (
        <Page className="habits-page">
            <PageHeader
                eyebrow="Habits"
                title="Habits"
                subtitle="Manage recurring habits, keep them active, and open each one for editing."
                actions={
                    <Link className="button button--primary" to="/habits/new">
                        + Add Habit
                    </Link>
                }
            />

            <Panel>
                <PanelBody>
                    {isLoading ? (
                        <Message variant="info">Loading habits...</Message>
                    ) : loadError ? (
                        <Message variant="error">{loadError}</Message>
                    ) : (
                        <div className="habits-grid">
                            {habits.map((habit) => (
                                <HabitCard
                                    enabled={habit.enabled}
                                    icon={habit.icon}
                                    key={habit.uuid}
                                    onToggle={(nextEnabled) => handleToggle(habit.uuid, nextEnabled)}
                                    type={habit.type}
                                    title={habit.title}
                                    typeLabel={habit.typeLabel}
                                    uuid={habit.uuid}
                                />
                            ))}
                        </div>
                    )}
                </PanelBody>
            </Panel>

            <AppNavigation />
        </Page>
    )
}
