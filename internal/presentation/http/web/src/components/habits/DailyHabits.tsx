import axios from 'axios'
import { useEffect, useMemo, useRef, useState } from 'react'
import { Link } from 'react-router-dom'
import { GoogleIcon, GoogleIcons } from '../icons'
import { Button } from '../primitives'
import { Panel, PanelBody } from '../layout'
import { useApiClient } from '../../hooks/useApiClient'
import './DailyHabits.sass'

const timeHourOptions = Array.from({ length: 24 }, (_, index) => String(index).padStart(2, '0'))
const timeMinuteOptions = Array.from({ length: 60 }, (_, index) =>
    String(index).padStart(2, '0'),
)

export type DailyCompletableHabit = {
    id: string
    label: string
    icon: number
    archivedAt: string | null
    deletedAt?: string | null
}

export type DailyTimeHabit = {
    id: string
    label: string
    icon: number
    archivedAt: string | null
    deletedAt?: string | null
}

export type DailyMeasurableHabit = {
    id: string
    label: string
    icon: number
    step: number
    unit: string
    archivedAt: string | null
    deletedAt?: string | null
}

export type DailyCompletableHabitRecord = {
    habitId: string
    value: boolean
}

export type DailyTimeHabitRecord = {
    habitId: string
    value: number
}

export type DailyMeasurableHabitRecord = {
    habitId: string
    value: number
}

export type DailyHabitsData = {
    completable: DailyCompletableHabit[]
    time: DailyTimeHabit[]
    measurable: DailyMeasurableHabit[]
}

export type DailyHabitRecords = {
    completable: DailyCompletableHabitRecord[]
    time: DailyTimeHabitRecord[]
    measurable: DailyMeasurableHabitRecord[]
}

type DailyHabitsProps = {
    date: Date
    dateKey: string
    habits: DailyHabitsData
    records: DailyHabitRecords
}

type SavingStatus = 'saved' | 'saving' | 'error'

type HabitValues = {
    completable: Record<string, boolean>
    time: Record<string, string>
    measurable: Record<string, number>
}

type HoldTimers = {
    delayId?: number
    intervalId?: number
}

type HabitRow =
    | {
          kind: 'completable'
          id: string
          title: string
          icon: GoogleIcons
          archivedAt: string | null
          deletedAt?: string | null
          checked: boolean
      }
    | {
          kind: 'time'
          id: string
          title: string
          icon: GoogleIcons
          archivedAt: string | null
          deletedAt?: string | null
          value: string
      }
    | {
          kind: 'measurable'
          id: string
          title: string
          icon: GoogleIcons
          archivedAt: string | null
          deletedAt?: string | null
          value: number
          step: number
          unit: string
      }

// DailyHabits renders the day habits panel, keeps the local record state, and synchronizes changes with the API.
export function DailyHabits({ date, dateKey, habits, records }: DailyHabitsProps) {
    const apiClient = useApiClient()
    const [values, setValues] = useState<HabitValues>(() => buildInitialValues(habits, records))
    const [recordedHabitIds, setRecordedHabitIds] = useState<Record<string, boolean>>(() =>
        buildRecordedHabitIds(records),
    )
    const [status, setStatus] = useState<SavingStatus>('saved')
    const measurableTimers = useRef<Record<string, HoldTimers>>({})
    const measurableSaveTimers = useRef<Record<string, number | undefined>>({})
    const measurableClickSuppression = useRef<Record<string, boolean>>({})

    useEffect(() => {
        setValues(buildInitialValues(habits, records))
        setRecordedHabitIds(buildRecordedHabitIds(records))
        setStatus('saved')
        clearAllMeasurableTimers(measurableTimers.current)
        clearAllMeasurableSaveTimers(measurableSaveTimers.current)
        measurableClickSuppression.current = {}

        return () => {
            clearAllMeasurableTimers(measurableTimers.current)
            clearAllMeasurableSaveTimers(measurableSaveTimers.current)
            measurableClickSuppression.current = {}
        }
    }, [dateKey, habits, records])

    const rows = useMemo(() => buildHabitRows(habits, values, date), [date, habits, values])
    const checkedCount = useMemo(
        () => countMarkedHabits(rows, recordedHabitIds),
        [recordedHabitIds, rows],
    )
    const totalCount = rows.length

    function updateCompletableValue(habitId: string, value: boolean) {
        setStatus('saving')
        markHabitRecorded(habitId)
        setValues((current) => ({
            ...current,
            completable: {
                ...current.completable,
                [habitId]: value,
            },
        }))
    }

    function updateTimeValue(habitId: string, value: string) {
        setStatus('saving')
        markHabitRecorded(habitId)
        setValues((current) => ({
            ...current,
            time: {
                ...current.time,
                [habitId]: value,
            },
        }))
    }

    async function sendCompletableRecord(habitId: string, value: boolean) {
        setStatus('saving')

        try {
            await apiClient.post(`habits/completable/${habitId}/${dateKey}`, {
                value,
            })
            setStatus('saved')
        } catch (error) {
            setStatus('error')
            if (!axios.isAxiosError(error)) {
                return
            }
        }
    }

    async function sendTimeRecord(habitId: string, value: string) {
        const minutes = parseTimeValue(value)
        if (minutes === null) {
            return
        }

        setStatus('saving')

        try {
            await apiClient.post(`habits/time/${habitId}/${dateKey}`, {
                value: minutes,
            })
            setStatus('saved')
        } catch (error) {
            setStatus('error')
            if (!axios.isAxiosError(error)) {
                return
            }
        }
    }

    async function sendMeasurableRecord(habitId: string, value: number) {
        setStatus('saving')

        try {
            await apiClient.post(`habits/measurable/${habitId}/${dateKey}`, {
                value,
            })
            setStatus('saved')
        } catch (error) {
            setStatus('error')
            if (!axios.isAxiosError(error)) {
                return
            }
        }
    }

    function handleCompletableToggle(habitId: string, nextChecked: boolean) {
        updateCompletableValue(habitId, nextChecked)
        void sendCompletableRecord(habitId, nextChecked)
    }

    function handleTimeChange(habitId: string, nextValue: string) {
        updateTimeValue(habitId, nextValue)
        void sendTimeRecord(habitId, nextValue)
    }

    function handleMeasurableStep(habitId: string, step: number, delta: number) {
        setStatus('saving')
        markHabitRecorded(habitId)
        setValues((current) => {
            const currentValue = current.measurable[habitId] ?? 0
            const nextValue = Math.max(0, roundToPrecision(currentValue + delta, step))
            scheduleMeasurableSave(habitId, nextValue)

            return {
                ...current,
                measurable: {
                    ...current.measurable,
                    [habitId]: nextValue,
                },
            }
        })
    }

    function startMeasurableHold(habitId: string, step: number, delta: number) {
        clearMeasurableHoldTimers(habitId, measurableTimers)
        measurableClickSuppression.current[habitId] = true
        handleMeasurableStep(habitId, step, delta)

        const delayId = window.setTimeout(() => {
            const intervalId = window.setInterval(() => {
                handleMeasurableStep(habitId, step, delta)
            }, 100)

            measurableTimers.current[habitId] = { intervalId }
        }, 400)

        measurableTimers.current[habitId] = { delayId }
    }

    function markHabitRecorded(habitId: string) {
        setRecordedHabitIds((current) => {
            if (current[habitId]) {
                return current
            }

            return {
                ...current,
                [habitId]: true,
            }
        })
    }

    function stopMeasurableHold(habitId: string) {
        clearMeasurableHoldTimers(habitId, measurableTimers)
        window.setTimeout(() => {
            measurableClickSuppression.current[habitId] = false
        }, 0)
    }

    function scheduleMeasurableSave(habitId: string, value: number) {
        const timer = measurableSaveTimers.current[habitId]
        if (timer) {
            window.clearTimeout(timer)
        }

        measurableSaveTimers.current[habitId] = window.setTimeout(() => {
            void sendMeasurableRecord(habitId, value)
        }, 800)
    }

    return (
        <Panel className="daily-habits">
            <PanelBody className="daily-habits__body">
                <div className="daily-habits__header">
                    <div className="daily-habits__heading">
                        <h2 className="daily-habits__title">Habits</h2>
                        <p className="daily-habits__subtitle">
                            Track the habits that belong to {formatPageDate(date)}.
                        </p>
                    </div>

                    <div className="daily-habits__header-actions">
                        <Link className="button button--secondary" to="/habits">
                            Manage
                        </Link>
                    </div>
                </div>

                <div className="daily-habits__list">
                    {rows.map((habit) =>
                        habit.kind === 'completable' ? (
                            <article className="daily-habit" key={habit.id}>
                                <div className="daily-habit__icon">
                                    <GoogleIcon icon={habit.icon} size={28} />
                                </div>

                                <div className="daily-habit__content">
                                    <h3 className="daily-habit__title">{habit.title}</h3>
                                    <p className="daily-habit__meta">Completion</p>
                                </div>

                                <div className="daily-habit__control daily-habit__control--completion">
                                    <button
                                        aria-pressed={habit.checked}
                                        className={`daily-habit-checkbox ${habit.checked ? 'daily-habit-checkbox--checked' : ''}`}
                                        onClick={() => handleCompletableToggle(habit.id, !habit.checked)}
                                        type="button"
                                    >
                                        <span className="daily-habit-checkbox__box" aria-hidden="true">
                                            {habit.checked ? '✓' : ''}
                                        </span>
                                        <span className="daily-habit-checkbox__label">
                                            {habit.checked ? 'Done' : 'Pending'}
                                        </span>
                                    </button>
                                </div>
                            </article>
                        ) : habit.kind === 'measurable' ? (
                            <article className="daily-habit" key={habit.id}>
                                <div className="daily-habit__icon">
                                    <GoogleIcon icon={habit.icon} size={28} />
                                </div>

                                <div className="daily-habit__content">
                                    <h3 className="daily-habit__title">{habit.title}</h3>
                                    <p className="daily-habit__meta">
                                        Quantity · step {formatMeasuredValue(habit.step)} {habit.unit}
                                    </p>
                                </div>

                                <div className="daily-habit__control daily-habit__control--measurable">
                                    <div className="daily-habit-stepper">
                                        <Button
                                            aria-label={`Decrease ${habit.title}`}
                                            className="daily-habit-stepper__button"
                                            onClick={(event) => {
                                                if (
                                                    event.detail > 0 &&
                                                    measurableClickSuppression.current[habit.id]
                                                ) {
                                                    measurableClickSuppression.current[habit.id] = false
                                                    return
                                                }

                                                handleMeasurableStep(habit.id, habit.step, -habit.step)
                                            }}
                                            onPointerDown={() =>
                                                startMeasurableHold(habit.id, habit.step, -habit.step)
                                            }
                                            onPointerLeave={() => stopMeasurableHold(habit.id)}
                                            onPointerUp={() => stopMeasurableHold(habit.id)}
                                            onPointerCancel={() => stopMeasurableHold(habit.id)}
                                            type="button"
                                            variant="secondary"
                                        >
                                            -
                                        </Button>
                                        <span className="daily-habit-stepper__value">
                                            {formatMeasuredValue(habit.value, habit.step)} {habit.unit}
                                        </span>
                                        <Button
                                            aria-label={`Increase ${habit.title}`}
                                            className="daily-habit-stepper__button"
                                            onClick={(event) => {
                                                if (
                                                    event.detail > 0 &&
                                                    measurableClickSuppression.current[habit.id]
                                                ) {
                                                    measurableClickSuppression.current[habit.id] = false
                                                    return
                                                }

                                                handleMeasurableStep(habit.id, habit.step, habit.step)
                                            }}
                                            onPointerDown={() =>
                                                startMeasurableHold(habit.id, habit.step, habit.step)
                                            }
                                            onPointerLeave={() => stopMeasurableHold(habit.id)}
                                            onPointerUp={() => stopMeasurableHold(habit.id)}
                                            onPointerCancel={() => stopMeasurableHold(habit.id)}
                                            type="button"
                                            variant="secondary"
                                        >
                                            +
                                        </Button>
                                    </div>
                                </div>
                            </article>
                        ) : (
                            <article className="daily-habit" key={habit.id}>
                                <div className="daily-habit__icon">
                                    <GoogleIcon icon={habit.icon} size={28} />
                                </div>

                                <div className="daily-habit__content">
                                    <h3 className="daily-habit__title">{habit.title}</h3>
                                    <p className="daily-habit__meta">Time</p>
                                </div>

                                <div className="daily-habit__control daily-habit__control--time">
                                    <div className="daily-habit-time">
                                        <span className="daily-habit-time__icon" aria-hidden="true">
                                            <GoogleIcon icon={GoogleIcons.Schedule} size={18} />
                                        </span>
                                        <select
                                            aria-label={`Select hour for ${habit.title}`}
                                            className="daily-habit-time__select"
                                            onChange={(event) =>
                                                handleTimeChange(
                                                    habit.id,
                                                    joinTimeValue(
                                                        event.target.value,
                                                        getTimeMinutes(habit.value),
                                                    ),
                                                )
                                            }
                                            value={getTimeHours(habit.value)}
                                        >
                                            {timeHourOptions.map((hour) => (
                                                <option key={hour} value={hour}>
                                                    {hour}
                                                </option>
                                            ))}
                                        </select>

                                        <span className="daily-habit-time__separator" aria-hidden="true">
                                            :
                                        </span>

                                        <select
                                            aria-label={`Select minute for ${habit.title}`}
                                            className="daily-habit-time__select"
                                            onChange={(event) =>
                                                handleTimeChange(
                                                    habit.id,
                                                    joinTimeValue(
                                                        getTimeHours(habit.value),
                                                        event.target.value,
                                                    ),
                                                )
                                            }
                                            value={getTimeMinutes(habit.value)}
                                        >
                                            {timeMinuteOptions.map((minute) => (
                                                <option key={minute} value={minute}>
                                                    {minute}
                                                </option>
                                            ))}
                                        </select>
                                    </div>
                                </div>
                            </article>
                        ),
                        )}
                </div>

                <div className="daily-habits__footer">
                    <div className="daily-habits__count" aria-label="Checked habits count">
                        <span className="daily-habits__count-text">
                            {checkedCount} of {totalCount} marked
                        </span>

                        <span className="daily-habits__count-bar" aria-hidden="true">
                            <span
                                className="daily-habits__count-bar-fill"
                                style={{
                                    width: `${totalCount === 0 ? 0 : (checkedCount / totalCount) * 100}%`,
                                }}
                            />
                        </span>
                    </div>

                    <span className="daily-habits__status" data-status={status}>
                        {statusLabel(status)}
                    </span>
                </div>
            </PanelBody>
        </Panel>
    )
}

function buildInitialValues(habits: DailyHabitsData, records: DailyHabitRecords): HabitValues {
    return {
        completable: habits.completable.reduce<Record<string, boolean>>((accumulator, habit) => {
            const record = records.completable.find((item) => item.habitId === habit.id)
            accumulator[habit.id] = record?.value ?? false

            return accumulator
        }, {}),
        time: habits.time.reduce<Record<string, string>>((accumulator, habit) => {
            const record = records.time.find((item) => item.habitId === habit.id)
            accumulator[habit.id] = record ? formatTimeValue(record.value) : '00:00'

            return accumulator
        }, {}),
        measurable: habits.measurable.reduce<Record<string, number>>((accumulator, habit) => {
            const record = records.measurable.find((item) => item.habitId === habit.id)
            accumulator[habit.id] = record?.value ?? 0

            return accumulator
        }, {}),
    }
}

function buildHabitRows(habits: DailyHabitsData, values: HabitValues, date: Date): HabitRow[] {
    const pageDate = toCalendarDate(date)

    return [
        ...habits.completable
            .filter((habit) => isHabitVisible(habit.archivedAt, habit.deletedAt, pageDate))
            .map((habit) => ({
                kind: 'completable' as const,
                id: habit.id,
                title: habit.label,
                icon: habit.icon as GoogleIcons,
                archivedAt: habit.archivedAt,
                deletedAt: habit.deletedAt,
                checked: values.completable[habit.id] ?? false,
            })),
        ...habits.measurable
            .filter((habit) => isHabitVisible(habit.archivedAt, habit.deletedAt, pageDate))
            .map((habit) => ({
                kind: 'measurable' as const,
                id: habit.id,
                title: habit.label,
                icon: habit.icon as GoogleIcons,
                archivedAt: habit.archivedAt,
                deletedAt: habit.deletedAt,
                value: values.measurable[habit.id] ?? 0,
                step: habit.step,
                unit: habit.unit,
            })),
        ...habits.time
            .filter((habit) => isHabitVisible(habit.archivedAt, habit.deletedAt, pageDate))
            .map((habit) => ({
                kind: 'time' as const,
                id: habit.id,
                title: habit.label,
                icon: habit.icon as GoogleIcons,
                archivedAt: habit.archivedAt,
                deletedAt: habit.deletedAt,
                value: values.time[habit.id] ?? '',
            })),
    ]
}

function isHabitVisible(
    archivedAt: string | null | undefined,
    deletedAt: string | null | undefined,
    pageDate: Date,
) {
    return isAfterCalendarDate(archivedAt, pageDate) && isAfterCalendarDate(deletedAt, pageDate)
}

function isAfterCalendarDate(value: string | null | undefined, pageDate: Date) {
    if (!value) {
        return true
    }

    return toCalendarDate(new Date(value)).getTime() > pageDate.getTime()
}

function toCalendarDate(value: Date) {
    const calendarDate = new Date(value)
    calendarDate.setHours(0, 0, 0, 0)

    return calendarDate
}

function parseTimeValue(value: string) {
    if (!value) {
        return null
    }

    const [hoursText, minutesText] = value.split(':')
    const hours = Number(hoursText)
    const minutes = Number(minutesText)

    if (!Number.isInteger(hours) || !Number.isInteger(minutes)) {
        return null
    }

    return hours * 60 + minutes
}

function formatTimeValue(value: number) {
    const hours = Math.floor(value / 60)
    const minutes = value % 60

    return `${String(hours).padStart(2, '0')}:${String(minutes).padStart(2, '0')}`
}

function getTimeHours(value: string) {
    if (!value) {
        return '00'
    }

    const [hours = '00'] = value.split(':')
    return hours.padStart(2, '0')
}

function getTimeMinutes(value: string) {
    if (!value) {
        return '00'
    }

    const [, minutes = '00'] = value.split(':')
    return minutes.padStart(2, '0')
}

function joinTimeValue(hours: string, minutes: string) {
    return `${hours.padStart(2, '0')}:${minutes.padStart(2, '0')}`
}

function formatMeasuredValue(value: number, step?: number) {
    const precision = step ? resolvePrecision(step) : 0
    return Number(value.toFixed(precision)).toString()
}

function resolvePrecision(step: number) {
    const text = step.toString()
    if (text.includes('e-')) {
        const [, exponent] = text.split('e-')
        return Number(exponent)
    }

    const [, decimals = ''] = text.split('.')
    return decimals.length
}

function roundToPrecision(value: number, step: number) {
    return Number(value.toFixed(resolvePrecision(step)))
}

function clearAllMeasurableTimers(timers: Record<string, HoldTimers>) {
    Object.keys(timers).forEach((habitId) => {
        clearMeasurableHoldTimers(habitId, { current: timers })
    })
}

function clearAllMeasurableSaveTimers(timers: Record<string, number | undefined>) {
    Object.values(timers).forEach((timer) => {
        if (timer) {
            window.clearTimeout(timer)
        }
    })
}

function clearMeasurableHoldTimers(
    habitId: string,
    timersRef: { current: Record<string, HoldTimers> },
) {
    const timers = timersRef.current[habitId]
    if (!timers) {
        return
    }

    if (timers.delayId) {
        window.clearTimeout(timers.delayId)
    }

    if (timers.intervalId) {
        window.clearInterval(timers.intervalId)
    }

    delete timersRef.current[habitId]
}

function formatPageDate(date: Date) {
    return new Intl.DateTimeFormat('en-US', {
        weekday: 'long',
        day: 'numeric',
        month: 'long',
        year: 'numeric',
    }).format(date)
}

function statusLabel(status: SavingStatus) {
    if (status === 'saving') {
        return 'Saving...'
    }

    if (status === 'error') {
        return 'Save failed'
    }

    return 'Saved'
}

function buildRecordedHabitIds(records: DailyHabitRecords) {
    return [...records.completable, ...records.time, ...records.measurable].reduce<Record<string, boolean>>(
        (accumulator, record) => {
            accumulator[record.habitId] = true
            return accumulator
        },
        {},
    )
}

function countMarkedHabits(rows: HabitRow[], recordedHabitIds: Record<string, boolean>) {
    return rows.reduce((count, row) => count + (recordedHabitIds[row.id] ? 1 : 0), 0)
}
