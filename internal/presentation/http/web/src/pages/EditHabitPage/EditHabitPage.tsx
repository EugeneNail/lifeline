import axios from 'axios'
import { useEffect, useMemo, useState } from 'react'
import type { FormEvent } from 'react'
import { Navigate, useNavigate, useParams } from 'react-router-dom'
import { AppNavigation } from '../../components/navigation'
import { GoogleIcon, GoogleIcons, IconSelector } from '../../components/icons'
import { Button, Message, TextField } from '../../components/primitives'
import { Page, PageHeader, Panel, PanelBody, Section, SectionHeader } from '../../components/layout'
import { useApiClient } from '../../hooks/useApiClient'
import type { CompletableHabit } from './CompletableHabit'
import type { MeasurableHabit } from './MeasurableHabit'
import type { TimeHabit } from './TimeHabit'
import './EditHabitPage.sass'

type HabitType = 'completable' | 'measurable' | 'time'

type HabitTypeConfig = {
    title: string
    subtitle: string
    endpoint: string
}

type HabitResponse = CompletableHabit | MeasurableHabit | TimeHabit

type EditHabitFieldErrors = Partial<Record<'label' | 'icon' | 'step' | 'unit', string>>

const habitTypeConfigs: Record<HabitType, HabitTypeConfig> = {
    completable: {
        title: 'Edit completion habit',
        subtitle: 'Update the habit details and decide whether to archive or remove it.',
        endpoint: 'habits/completable',
    },
    measurable: {
        title: 'Edit quantity habit',
        subtitle: 'Update the habit details, measurement settings, or manage its lifecycle.',
        endpoint: 'habits/measurable',
    },
    time: {
        title: 'Edit time habit',
        subtitle: 'Update the habit details and manage its lifecycle.',
        endpoint: 'habits/time',
    },
}

function isHabitType(value: string | undefined): value is HabitType {
    return value === 'completable' || value === 'measurable' || value === 'time'
}

function mapResponseToState(response: HabitResponse) {
    return {
        label: response.label,
        icon: response.icon as GoogleIcons,
        step: 'step' in response ? String(response.step) : '',
        unit: 'unit' in response ? response.unit : '',
    }
}

function createUpdatePayload(
    habitType: HabitType,
    label: string,
    icon: GoogleIcons,
    step: string,
    unit: string,
) {
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

// EditHabitPage renders the habit edit form and loads the habit state from the API.
export function EditHabitPage() {
    const apiClient = useApiClient()
    const navigate = useNavigate()
    const params = useParams<{ type: string; id: string }>()
    const habitType = params.type
    const habitId = params.id
    const [label, setLabel] = useState('')
    const [icon, setIcon] = useState<GoogleIcons>(GoogleIcons.Favorite)
    const [unit, setUnit] = useState('')
    const [step, setStep] = useState('')
    const [isLoading, setIsLoading] = useState(true)
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [fieldErrors, setFieldErrors] = useState<EditHabitFieldErrors>({})
    const [formError, setFormError] = useState('')
    const [notFound, setNotFound] = useState(false)

    const habitConfig = useMemo(() => {
        if (!isHabitType(habitType)) {
            return null
        }

        return habitTypeConfigs[habitType]
    }, [habitType])

    useEffect(() => {
        let isActive = true

        async function loadHabit() {
            if (!habitConfig || !habitId) {
                if (isActive) {
                    setNotFound(true)
                }
                return
            }

            setIsLoading(true)
            setFormError('')
            setFieldErrors({})

            try {
                const response = await apiClient.get<HabitResponse>(`${habitConfig.endpoint}/${habitId}`)
                if (!isActive) {
                    return
                }

                const mapped = mapResponseToState(response.data)
                setLabel(mapped.label)
                setIcon(mapped.icon)
                setStep(mapped.step)
                setUnit(mapped.unit)
            } catch (error) {
                if (!isActive) {
                    return
                }

                if (axios.isAxiosError(error) && error.response?.status === 404) {
                    setNotFound(true)
                    return
                }

                setFormError('Could not load habit.')
            } finally {
                if (isActive) {
                    setIsLoading(false)
                }
            }
        }

        void loadHabit()

        return () => {
            isActive = false
        }
    }, [apiClient, habitConfig, habitId])

    if (notFound) {
        return <Navigate replace to="/habits" />
    }

    if (!habitConfig || !habitId) {
        return <Navigate replace to="/habits" />
    }

    const resolvedHabitType = habitType as HabitType
    const resolvedHabitConfig = habitConfig
    const resolvedHabitId = habitId

    function setEditHabitErrors(error: unknown) {
        if (!axios.isAxiosError(error)) {
            setFormError('Could not update habit.')
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

        if (error.response?.status === 403) {
            setFormError('This habit belongs to another user.')
            return
        }

        if (error.response?.status === 409) {
            setFormError('Habit is archived.')
            return
        }

        if (error.response?.status === 410) {
            setFormError('Habit is deleted.')
            return
        }

        setFormError('Could not update habit.')
    }

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        setIsSubmitting(true)
        setFieldErrors({})
        setFormError('')

        try {
            await apiClient.put(`${resolvedHabitConfig.endpoint}/${resolvedHabitId}`, createUpdatePayload(
                resolvedHabitType,
                label,
                icon,
                step,
                unit,
            ))
            navigate('/habits')
        } catch (error) {
            setEditHabitErrors(error)
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <Page className="edit-habit-page">
            <PageHeader
                eyebrow="Habits"
                title={resolvedHabitConfig.title}
                subtitle={resolvedHabitConfig.subtitle}
                actions={
                    <Button type="button" variant="secondary" onClick={() => navigate(-1)}>
                        Back
                    </Button>
                }
            />

            <Panel>
                <PanelBody>
                    {isLoading ? (
                        <Message variant="info">Loading habit...</Message>
                    ) : (
                        <form className="edit-habit-form" onSubmit={handleSubmit}>
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

                            <div className="edit-habit-actions">
                                <Button
                                    type="submit"
                                    disabled={isSubmitting}
                                    aria-label={isSubmitting ? 'Saving habit' : 'Save habit'}
                                >
                                    <GoogleIcon icon={icon} size={18} />
                                    <span
                                        className="edit-habit-submit-label"
                                        data-submitting={isSubmitting}
                                        aria-hidden="true"
                                    >
                                        <span className="edit-habit-submit-label__text">
                                            Save habit
                                        </span>
                                        <span className="edit-habit-submit-label__text edit-habit-submit-label__text--pending">
                                            Saving...
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
                                <Button type="button" variant="secondary" disabled>
                                    Archive
                                </Button>
                                <Button type="button" variant="danger" disabled>
                                    Delete
                                </Button>
                            </div>
                        </form>
                    )}
                </PanelBody>
            </Panel>

            <AppNavigation />
        </Page>
    )
}
