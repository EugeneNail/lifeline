import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Page, PageHeader, Panel, PanelBody } from '../../components/layout'
import { AppNavigation } from '../../components/navigation'
import { HabitCard } from '../../components/habits'
import { GoogleIcons } from '../../components/icons'
import './HabitsPage.sass'

type HabitItem = {
    uuid: string
    icon: GoogleIcons
    title: string
    typeLabel: string
    enabled: boolean
}

const completionHabits: HabitItem[] = [
    {
        uuid: '018f1b7a-0000-7000-8000-000000000001',
        icon: GoogleIcons.Check,
        title: 'Work on project',
        typeLabel: 'Completion',
        enabled: true,
    },
    {
        uuid: '018f1b7a-0000-7000-8000-000000000002',
        icon: GoogleIcons.Check,
        title: 'Read a book',
        typeLabel: 'Completion',
        enabled: false,
    },
]

const timeHabits: HabitItem[] = [
    {
        uuid: '018f1b7a-0000-7000-8000-000000000003',
        icon: GoogleIcons.Schedule,
        title: 'Go to bed',
        typeLabel: 'Time',
        enabled: true,
    },
    {
        uuid: '018f1b7a-0000-7000-8000-000000000004',
        icon: GoogleIcons.Schedule,
        title: 'Morning walk',
        typeLabel: 'Time',
        enabled: false,
    },
]

const measurableHabits: HabitItem[] = [
    {
        uuid: '018f1b7a-0000-7000-8000-000000000005',
        icon: GoogleIcons.Water,
        title: 'Drink water',
        typeLabel: 'Quantity',
        enabled: true,
    },
    {
        uuid: '018f1b7a-0000-7000-8000-000000000006',
        icon: GoogleIcons.Numbers,
        title: 'Study pages',
        typeLabel: 'Quantity',
        enabled: true,
    },
]

// HabitsPage renders the habit management dashboard with grouped habit cards.
export function HabitsPage() {
    const [habitStates, setHabitStates] = useState<Record<string, boolean>>(() =>
        Object.fromEntries(
            [...completionHabits, ...timeHabits, ...measurableHabits].map((habit) => [
                habit.uuid,
                habit.enabled,
            ]),
        ),
    )

    function handleToggle(uuid: string, nextEnabled: boolean) {
        setHabitStates((current) => ({
            ...current,
            [uuid]: nextEnabled,
        }))
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
                    <div className="habits-grid">
                        {[...completionHabits, ...timeHabits, ...measurableHabits].map((habit) => (
                            <HabitCard
                                enabled={habitStates[habit.uuid]}
                                icon={habit.icon}
                                key={habit.uuid}
                                onToggle={(nextEnabled) => handleToggle(habit.uuid, nextEnabled)}
                                title={habit.title}
                                typeLabel={habit.typeLabel}
                                uuid={habit.uuid}
                            />
                        ))}
                    </div>
                </PanelBody>
            </Panel>

            <AppNavigation />
        </Page>
    )
}
