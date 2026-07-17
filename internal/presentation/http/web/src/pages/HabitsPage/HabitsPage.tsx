import { useState } from 'react'
import { Link } from 'react-router-dom'
import { Page, PageHeader, Panel, PanelBody, Section, SectionHeader } from '../../components/layout'
import { AppNavigation } from '../../components/navigation'
import { HabitCard } from '../../components/habits'
import { GoogleIcons } from '../../components/icons'
import './HabitsPage.sass'

type HabitItem = {
    uuid: string
    icon: GoogleIcons
    title: string
    description: string
    typeLabel: string
    enabled: boolean
}

const completionHabits: HabitItem[] = [
    {
        uuid: '018f1b7a-0000-7000-8000-000000000001',
        icon: GoogleIcons.Check,
        title: 'Work on project',
        description: 'Mark whether it is done or not.',
        typeLabel: 'Completion',
        enabled: true,
    },
    {
        uuid: '018f1b7a-0000-7000-8000-000000000002',
        icon: GoogleIcons.Check,
        title: 'Read a book',
        description: 'Mark whether it is done or not.',
        typeLabel: 'Completion',
        enabled: false,
    },
]

const timeHabits: HabitItem[] = [
    {
        uuid: '018f1b7a-0000-7000-8000-000000000003',
        icon: GoogleIcons.Schedule,
        title: 'Go to bed',
        description: 'Record hours and minutes.',
        typeLabel: 'Time',
        enabled: true,
    },
    {
        uuid: '018f1b7a-0000-7000-8000-000000000004',
        icon: GoogleIcons.Schedule,
        title: 'Morning walk',
        description: 'Record hours and minutes.',
        typeLabel: 'Time',
        enabled: false,
    },
]

const measurableHabits: HabitItem[] = [
    {
        uuid: '018f1b7a-0000-7000-8000-000000000005',
        icon: GoogleIcons.Water,
        title: 'Drink water',
        description: 'Step 250 · Unit ml',
        typeLabel: 'Quantity',
        enabled: true,
    },
    {
        uuid: '018f1b7a-0000-7000-8000-000000000006',
        icon: GoogleIcons.Numbers,
        title: 'Study pages',
        description: 'Step 5 · Unit pages',
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
                    <Section>
                        <SectionHeader title="Completion habits" />
                        <div className="habits-section__grid">
                            {completionHabits.map((habit) => (
                                <HabitCard
                                    description={habit.description}
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
                    </Section>

                    <Section>
                        <SectionHeader title="Time habits" />
                        <div className="habits-section__grid">
                            {timeHabits.map((habit) => (
                                <HabitCard
                                    description={habit.description}
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
                    </Section>

                    <Section>
                        <SectionHeader title="Quantity habits" />
                        <div className="habits-section__grid">
                            {measurableHabits.map((habit) => (
                                <HabitCard
                                    description={habit.description}
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
                    </Section>
                </PanelBody>
            </Panel>

            <AppNavigation />
        </Page>
    )
}
