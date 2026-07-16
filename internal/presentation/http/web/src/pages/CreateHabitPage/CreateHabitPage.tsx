import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { Page, PageHeader, Panel, PanelBody, Section, SectionHeader } from '../../components/layout'
import { GoogleIcon, GoogleIcons, IconSelector } from '../../components/icons'
import { AppNavigation } from '../../components/navigation'
import { Button, TextField } from '../../components/primitives'
import './CreateHabitPage.sass'

type HabitType = 'completable' | 'measurable' | 'time'

type HabitTypeOption = {
    value: HabitType
    icon: string
    title: string
    description: string
    example: string
}

const habitTypes: HabitTypeOption[] = [
    {
        value: 'completable',
        icon: '✓',
        title: 'Completion',
        description: 'Mark whether it is done or not.',
        example: 'Example: "Brush teeth"',
    },
    {
        value: 'measurable',
        icon: '123',
        title: 'Quantity',
        description: 'Record a numeric value.',
        example: 'Example: "Water - 250 ml"',
    },
    {
        value: 'time',
        icon: '00:00',
        title: 'Time',
        description: 'Record hours and minutes.',
        example: 'Example: "Run - 00:30"',
    },
]

// CreateHabitPage renders the habit creation form.
export function CreateHabitPage() {
    const navigate = useNavigate()
    const [label, setLabel] = useState('')
    const [icon, setIcon] = useState<GoogleIcons>(GoogleIcons.Check)
    const [habitType, setHabitType] = useState<HabitType>('completable')
    const [unit, setUnit] = useState('')
    const [step, setStep] = useState('')

    function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
    }

    return (
        <Page className="create-habit-page">
            <PageHeader
                eyebrow="Habits"
                title="New habit"
                subtitle="Choose how this habit should be tracked."
                actions={
                    <Button type="button" variant="secondary" onClick={() => navigate(-1)}>
                        Back
                    </Button>
                }
            />

            <Panel>
                <PanelBody>
                    <form className="create-habit-form" onSubmit={handleSubmit}>
                        <Section>
                            <SectionHeader title="Habit details" />
                            <TextField
                                label="Habit name"
                                name="label"
                                placeholder="Drink water"
                                value={label}
                                onChange={(event) => setLabel(event.target.value)}
                            />
                        </Section>

                        <Section>
                            <SectionHeader title="Icon" meta="Choose one" />
                            <IconSelector value={icon} onChange={setIcon} />
                        </Section>

                        <Section>
                            <SectionHeader title="Tracking type" />
                            <div className="habit-type-list">
                                {habitTypes.map((option) => (
                                    <button
                                        aria-pressed={habitType === option.value}
                                        className="habit-type-card"
                                        key={option.value}
                                        onClick={() => setHabitType(option.value)}
                                        type="button"
                                    >
                                        <span className="habit-type-card__icon">{option.icon}</span>
                                        <span className="habit-type-card__content">
                                            <strong>{option.title}</strong>
                                            <span>{option.description}</span>
                                            <em>{option.example}</em>
                                        </span>
                                    </button>
                                ))}
                            </div>
                        </Section>

                        {habitType === 'measurable' ? (
                            <Section>
                                <SectionHeader title="Measurement" />
                                <div className="measurement-fields">
                                    <TextField
                                        label="Unit"
                                        name="unit"
                                        placeholder="ml"
                                        value={unit}
                                        onChange={(event) => setUnit(event.target.value)}
                                    />
                                    <TextField
                                        label="Step"
                                        name="step"
                                        placeholder="250"
                                        type="number"
                                        value={step}
                                        onChange={(event) => setStep(event.target.value)}
                                    />
                                </div>
                            </Section>
                        ) : null}

                        <div className="create-habit-actions">
                            <Button type="submit">
                                <GoogleIcon icon={icon} size={18} />
                                Create habit
                            </Button>
                            <Button type="button" variant="secondary" onClick={() => navigate(-1)}>
                                Back
                            </Button>
                        </div>
                    </form>
                </PanelBody>
            </Panel>

            <AppNavigation />
        </Page>
    )
}
