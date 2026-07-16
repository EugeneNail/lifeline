import axios from 'axios'
import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { Page, PageHeader, Panel, PanelBody, Section, SectionHeader } from '../../components/layout'
import { GoogleIcon, GoogleIcons, IconSelector } from '../../components/icons'
import { AppNavigation } from '../../components/navigation'
import { Button, Message, TextField } from '../../components/primitives'
import { useApiClient } from '../../hooks/useApiClient'
import './CreateHabitPage.sass'

type HabitType = 'completable' | 'measurable' | 'time'

type HabitTypeOption = {
    value: HabitType
    icon: string
    title: string
    description: string
    example: string
}

type CreateHabitFieldErrors = Partial<Record<'label' | 'icon' | 'step' | 'unit', string>>

const habitTypeEndpoints: Record<HabitType, string> = {
    completable: 'habits/completable',
    measurable: 'habits/measurable',
    time: 'habits/time',
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
        value: 'time',
        icon: '00:00',
        title: 'Time',
        description: 'Record hours and minutes.',
        example: 'Example: "Run - 00:30"',
    },
    {
        value: 'measurable',
        icon: '123',
        title: 'Quantity',
        description: 'Record a numeric value.',
        example: 'Example: "Water - 250 ml"',
    },
]

// CreateHabitPage renders the habit creation form.
export function CreateHabitPage() {
    const apiClient = useApiClient()
    const navigate = useNavigate()
    const [label, setLabel] = useState('')
    const [icon, setIcon] = useState<GoogleIcons>(GoogleIcons.Check)
    const [habitType, setHabitType] = useState<HabitType>('completable')
    const [unit, setUnit] = useState('')
    const [step, setStep] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [fieldErrors, setFieldErrors] = useState<CreateHabitFieldErrors>({})
    const [formError, setFormError] = useState('')

    function setCreateHabitErrors(error: unknown) {
        if (!axios.isAxiosError(error)) {
            setFormError('Could not create habit.')
            return
        }

        if (
            error.response?.status === 422 &&
            error.response.data &&
            typeof error.response.data === 'object'
        ) {
            const response = error.response.data as Record<string, string>
            setFieldErrors({
                label: response.label,
                icon: response.icon,
                step: response.step,
                unit: response.unit,
            })
            return
        }

        if (error.response?.status === 409) {
            setFormError('Habit limit exceeded.')
            return
        }

        setFormError('Could not create habit.')
    }

    function createPayload() {
        if (habitType !== 'measurable') {
            return { label, icon }
        }

        return {
            label,
            icon,
            step: Number(step),
            unit,
        }
    }

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        setIsSubmitting(true)
        setFieldErrors({})
        setFormError('')

        try {
            await apiClient.post<string>(habitTypeEndpoints[habitType], createPayload())
            navigate('/')
        } catch (error) {
            setCreateHabitErrors(error)
        } finally {
            setIsSubmitting(false)
        }
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
                                error={fieldErrors.label}
                                onChange={(event) => setLabel(event.target.value)}
                            />
                        </Section>

                        <Section>
                            <SectionHeader title="Icon" meta="Choose one" />
                            <IconSelector value={icon} onChange={setIcon} />
                            {fieldErrors.icon ? (
                                <Message variant="error">{fieldErrors.icon}</Message>
                            ) : null}
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
                                        error={fieldErrors.unit}
                                        onChange={(event) => setUnit(event.target.value)}
                                    />
                                    <TextField
                                        label="Step"
                                        name="step"
                                        placeholder="250"
                                        type="number"
                                        value={step}
                                        error={fieldErrors.step}
                                        onChange={(event) => setStep(event.target.value)}
                                    />
                                </div>
                            </Section>
                        ) : null}

                        {formError ? <Message variant="error">{formError}</Message> : null}

                        <div className="create-habit-actions">
                            <Button
                                type="submit"
                                disabled={isSubmitting}
                                aria-label={isSubmitting ? 'Creating habit' : 'Create habit'}
                            >
                                <GoogleIcon icon={icon} size={18} />
                                <span
                                    className="create-habit-submit-label"
                                    data-submitting={isSubmitting}
                                    aria-hidden="true"
                                >
                                    <span className="create-habit-submit-label__text">
                                        Create habit
                                    </span>
                                    <span className="create-habit-submit-label__text create-habit-submit-label__text--pending">
                                        Creating...
                                    </span>
                                </span>
                            </Button>
                            <Button
                                type="button"
                                variant="secondary"
                                disabled={isSubmitting}
                                onClick={() => navigate(-1)}
                            >
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
