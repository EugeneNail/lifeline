import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { Button, IconButton } from '../primitives'
import './DateSelector.sass'

export type DateRange = {
    startDate: Date
    endDate: Date
}

type DateSelectorCommonProps = {
    className?: string
    step?: number
}

type DateSelectorSingleProps = DateSelectorCommonProps & {
    mode: 'single'
    value: Date | null
    onChange: (date: Date) => void
}

type DateSelectorRangeProps = DateSelectorCommonProps & {
    mode: 'range'
    value: DateRange | null
    onChange: (range: DateRange) => void
}

type DateSelectorProps = DateSelectorSingleProps | DateSelectorRangeProps

type ViewMode = 'days' | 'months' | 'years'

const monthNames = Array.from({ length: 12 }, (_, index) =>
    new Intl.DateTimeFormat('en-US', { month: 'long' }).format(new Date(2026, index, 1)),
)

const weekdayNames = Array.from({ length: 7 }, (_, index) =>
    new Intl.DateTimeFormat('en-US', { weekday: 'short' }).format(new Date(2026, 7, 2 + index)),
)

const singleDateLabelFormatter = new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
})

const rangeMonthYearLabelFormatter = new Intl.DateTimeFormat('en-GB', {
    month: 'long',
    year: 'numeric',
})

const rangeStartMonthLabelFormatter = new Intl.DateTimeFormat('en-GB', {
    day: 'numeric',
    month: 'long',
})

// DateSelector renders an inline date switcher and an attached calendar overlay.
export function DateSelector(props: DateSelectorProps) {
    const [isOpen, setIsOpen] = useState(false)
    const [viewMode, setViewMode] = useState<ViewMode>('days')
    const [visibleDate, setVisibleDate] = useState(() =>
        props.mode === 'single'
            ? resolveInitialVisibleDate(props.value, 'single')
            : resolveInitialVisibleDate(props.value, 'range'),
    )
    const [rangeAnchor, setRangeAnchor] = useState<Date | null>(null)

    useEffect(() => {
        if (!isOpen) {
            setRangeAnchor(null)
            return
        }

        setViewMode('days')
        setVisibleDate(
            props.mode === 'single'
                ? resolveInitialVisibleDate(props.value, 'single')
                : resolveInitialVisibleDate(props.value, 'range'),
        )
        setRangeAnchor(null)
    }, [isOpen, props.mode, props.value])

    useEffect(() => {
        if (!isOpen) {
            return
        }

        function handleKeyDown(event: KeyboardEvent) {
            if (event.key === 'Escape') {
                setIsOpen(false)
            }
        }

        window.addEventListener('keydown', handleKeyDown)

        return () => {
            window.removeEventListener('keydown', handleKeyDown)
        }
    }, [isOpen])

    const step = props.step ?? (props.mode === 'range' && props.value ? getRangeLength(props.value) : 1)
    const label = props.mode === 'single' ? formatSingleDateLabel(props.value) : formatRangeDateLabel(props.value)

    function handleShift(direction: number) {
        const delta = step * direction

        if (props.mode === 'single') {
            const currentDate = normalizeDate(props.value ?? new Date())
            props.onChange(addDays(currentDate, delta))
            return
        }

        const currentRange = props.value ?? {
            startDate: normalizeDate(new Date()),
            endDate: normalizeDate(new Date()),
        }

        props.onChange({
            startDate: addDays(normalizeDate(currentRange.startDate), delta),
            endDate: addDays(normalizeDate(currentRange.endDate), delta),
        })
    }

    if (!label) {
        return null
    }

    return (
        <>
            <div className={joinClassNames('date-selector', props.className)}>
                <button
                    aria-label={props.mode === 'single' ? 'Previous day' : 'Previous range'}
                    className="date-selector__step-button"
                    type="button"
                    onClick={() => handleShift(-1)}
                >
                    ‹
                </button>
                <button
                    aria-label={props.mode === 'single' ? 'Open date selector' : 'Open date selector'}
                    className="date-selector__value-button"
                    type="button"
                    onClick={() => setIsOpen(true)}
                >
                    <span className="date-selector__value-label">{label}</span>
                </button>
                <button
                    aria-label={props.mode === 'single' ? 'Next day' : 'Next range'}
                    className="date-selector__step-button"
                    type="button"
                    onClick={() => handleShift(1)}
                >
                    ›
                </button>
            </div>

            {isOpen
                ? createPortal(
                      <div
                          className="date-selector__overlay"
                          onMouseDown={(event) => {
                              if (event.target === event.currentTarget) {
                                  setIsOpen(false)
                              }
                          }}
                          role="presentation"
                      >
                          <section className="date-selector__surface" role="dialog" aria-modal="true">
                              <header className="date-selector__header">
                                  <div className="date-selector__header-main">
                                      <Button
                                          className="date-selector__mode-button"
                                          variant="secondary"
                                          type="button"
                                          onClick={() => setViewMode('months')}
                                      >
                                          {monthNames[visibleDate.getMonth()]}
                                      </Button>
                                      <Button
                                          className="date-selector__mode-button"
                                          variant="secondary"
                                          type="button"
                                          onClick={() => setViewMode('years')}
                                      >
                                          {visibleDate.getFullYear()}
                                      </Button>
                                  </div>

                                  <div className="date-selector__header-actions">
                                      <IconButton
                                          aria-label="Previous period"
                                          className="date-selector__nav-button"
                                          type="button"
                                          onClick={() => shiftVisibleDate(setVisibleDate, viewMode, -1)}
                                      >
                                          ←
                                      </IconButton>
                                      <IconButton
                                          aria-label="Next period"
                                          className="date-selector__nav-button"
                                          type="button"
                                          onClick={() => shiftVisibleDate(setVisibleDate, viewMode, 1)}
                                      >
                                          →
                                      </IconButton>
                                      <IconButton
                                          aria-label="Close date selector"
                                          className="date-selector__close-button"
                                          type="button"
                                          onClick={() => setIsOpen(false)}
                                      >
                                          ×
                                      </IconButton>
                                  </div>
                              </header>

                              <div className="date-selector__toolbar">
                                  <p className="date-selector__eyebrow">
                                      {props.mode === 'single' ? 'Single date' : 'Date range'}
                                  </p>
                                  <p className="date-selector__subtitle">
                                      {viewMode === 'days'
                                          ? 'Choose a day, or switch month and year from the buttons above.'
                                          : viewMode === 'months'
                                            ? 'Choose a month to continue.'
                                            : 'Choose a year between 2000 and 2099.'}
                                  </p>
                              </div>

                              {viewMode === 'days' ? (
                                  <div className="date-selector__calendar">
                                      <div className="date-selector__weekdays">
                                          {weekdayNames.map((weekday) => (
                                              <span className="date-selector__weekday" key={weekday}>
                                                  {weekday}
                                              </span>
                                          ))}
                                      </div>

                                      <div className="date-selector__days">
                                          {buildMonthCells(visibleDate).map((day, index) =>
                                              day ? (
                                                  <button
                                                      aria-pressed={isDaySelected(day, props.value, rangeAnchor)}
                                                      className={joinClassNames(
                                                          'date-selector__day',
                                                          isDaySelected(day, props.value, rangeAnchor)
                                                              ? 'date-selector__day--selected'
                                                              : undefined,
                                                          isDayInRange(day, props.value, rangeAnchor)
                                                              ? 'date-selector__day--in-range'
                                                              : undefined,
                                                          isSameDay(day, normalizeDate(new Date()))
                                                              ? 'date-selector__day--today'
                                                              : undefined,
                                                      )}
                                                      key={formatDayKey(day)}
                                                      type="button"
                                                      onClick={() =>
                                                          handleDayClick(
                                                              props,
                                                              day,
                                                              rangeAnchor,
                                                              setRangeAnchor,
                                                              () => setIsOpen(false),
                                                          )
                                                      }
                                                  >
                                                      {day.getDate()}
                                                  </button>
                                              ) : (
                                                  <span
                                                      className="date-selector__day date-selector__day--empty"
                                                      key={`empty-${index}`}
                                                  />
                                              ),
                                          )}
                                      </div>
                                  </div>
                              ) : viewMode === 'months' ? (
                                  <div className="date-selector__picker-grid date-selector__picker-grid--months">
                                      {monthNames.map((monthName, index) => (
                                          <button
                                              aria-pressed={visibleDate.getMonth() === index}
                                              className={joinClassNames(
                                                  'date-selector__picker-button',
                                                  visibleDate.getMonth() === index
                                                      ? 'date-selector__picker-button--selected'
                                                      : undefined,
                                              )}
                                              key={monthName}
                                              type="button"
                                              onClick={() => {
                                                  const nextDate = new Date(visibleDate)
                                                  nextDate.setMonth(index, 1)
                                                  setVisibleDate(normalizeDate(nextDate))
                                                  setViewMode('days')
                                              }}
                                          >
                                              {monthName}
                                          </button>
                                      ))}
                                  </div>
                              ) : (
                                  <div className="date-selector__picker-grid date-selector__picker-grid--years">
                                      {Array.from({ length: 100 }, (_, index) => 2000 + index).map((year) => (
                                          <button
                                              aria-pressed={visibleDate.getFullYear() === year}
                                              className={joinClassNames(
                                                  'date-selector__picker-button',
                                                  visibleDate.getFullYear() === year
                                                      ? 'date-selector__picker-button--selected'
                                                      : undefined,
                                              )}
                                              key={year}
                                              type="button"
                                              onClick={() => {
                                                  const nextDate = new Date(visibleDate)
                                                  nextDate.setFullYear(year)
                                                  setVisibleDate(normalizeDate(nextDate))
                                                  setViewMode('days')
                                              }}
                                          >
                                              {year}
                                          </button>
                                      ))}
                                  </div>
                              )}
                          </section>
                      </div>,
                      document.body,
                  )
                : null}
        </>
    )
}

function handleDayClick(
    props: DateSelectorProps,
    day: Date,
    rangeAnchor: Date | null,
    setRangeAnchor: (date: Date | null) => void,
    close: () => void,
) {
    if (props.mode === 'single') {
        props.onChange(normalizeDate(day))
        close()
        return
    }

    const nextDay = normalizeDate(day)

    if (!rangeAnchor) {
        setRangeAnchor(nextDay)
        return
    }

    const startDate = rangeAnchor < nextDay ? rangeAnchor : nextDay
    const endDate = rangeAnchor < nextDay ? nextDay : rangeAnchor
    props.onChange({
        startDate,
        endDate,
    })
    setRangeAnchor(null)
    close()
}

function isDaySelected(day: Date, selectedRange: Date | DateRange | null, rangeAnchor: Date | null) {
    if (!selectedRange) {
        return false
    }

    if (rangeAnchor && isSameDay(day, rangeAnchor)) {
        return true
    }

    if (selectedRange instanceof Date) {
        return isSameDay(day, selectedRange)
    }

    return isSameDay(day, selectedRange.startDate) || isSameDay(day, selectedRange.endDate)
}

function isDayInRange(day: Date, selectedRange: Date | DateRange | null, rangeAnchor: Date | null) {
    if (rangeAnchor) {
        return isSameDay(day, rangeAnchor)
    }

    if (!selectedRange || selectedRange instanceof Date) {
        return false
    }

    const normalizedDay = normalizeDate(day)
    const startDate = normalizeDate(selectedRange.startDate)
    const endDate = normalizeDate(selectedRange.endDate)

    return normalizedDay >= startDate && normalizedDay <= endDate
}

function shiftVisibleDate(
    setVisibleDate: (value: Date | ((current: Date) => Date)) => void,
    viewMode: ViewMode,
    delta: number,
) {
    setVisibleDate((current) => {
        const nextDate = new Date(current)

        if (viewMode === 'days') {
            nextDate.setMonth(nextDate.getMonth() + delta)
        } else if (viewMode === 'months') {
            nextDate.setFullYear(nextDate.getFullYear() + delta)
        } else {
            nextDate.setFullYear(nextDate.getFullYear() + delta * 10)
        }

        return normalizeDate(nextDate)
    })
}

function buildMonthCells(date: Date) {
    const year = date.getFullYear()
    const month = date.getMonth()
    const firstDay = normalizeDate(new Date(year, month, 1))
    const firstWeekday = firstDay.getDay()
    const daysInMonth = new Date(year, month + 1, 0).getDate()
    const cells: Array<Date | null> = []

    for (let index = 0; index < firstWeekday; index += 1) {
        cells.push(null)
    }

    for (let day = 1; day <= daysInMonth; day += 1) {
        cells.push(normalizeDate(new Date(year, month, day)))
    }

    return cells
}

function resolveInitialVisibleDate(value: Date | null, mode: 'single'): Date
function resolveInitialVisibleDate(value: DateRange | null, mode: 'range'): Date
function resolveInitialVisibleDate(
    value: Date | DateRange | null,
    mode: DateSelectorProps['mode'],
) {
    if (!value) {
        return normalizeDate(new Date())
    }

    if (mode === 'single') {
        return normalizeDate(value as Date)
    }

    return normalizeDate((value as DateRange).startDate)
}

function normalizeDate(date: Date) {
    const normalized = new Date(date)
    normalized.setHours(0, 0, 0, 0)

    return normalized
}

function isSameDay(left: Date, right: Date) {
    return normalizeDate(left).getTime() === normalizeDate(right).getTime()
}

function formatDayKey(date: Date) {
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function getRangeLength(range: DateRange) {
    const from = normalizeDate(range.startDate)
    const to = normalizeDate(range.endDate)
    return Math.floor((to.getTime() - from.getTime()) / (24 * 60 * 60 * 1000)) + 1
}

function addDays(date: Date, days: number) {
    const nextDate = new Date(date)
    nextDate.setDate(nextDate.getDate() + days)
    return normalizeDate(nextDate)
}

function formatSingleDateLabel(date: Date | null) {
    if (!date) {
        return 'Select date'
    }

    return singleDateLabelFormatter.format(normalizeDate(date))
}

function formatRangeDateLabel(range: DateRange | null) {
    if (!range) {
        return 'Select range'
    }

    const startDate = normalizeDate(range.startDate)
    const endDate = normalizeDate(range.endDate)

    if (startDate.getTime() === endDate.getTime()) {
        return singleDateLabelFormatter.format(startDate)
    }

    if (
        startDate.getFullYear() === endDate.getFullYear() &&
        startDate.getMonth() === endDate.getMonth()
    ) {
        return `${startDate.getDate()}–${endDate.getDate()} ${rangeMonthYearLabelFormatter.format(endDate)}`
    }

    if (startDate.getFullYear() === endDate.getFullYear()) {
        return `${rangeStartMonthLabelFormatter.format(startDate)} – ${singleDateLabelFormatter.format(endDate)}`
    }

    return `${singleDateLabelFormatter.format(startDate)} – ${singleDateLabelFormatter.format(endDate)}`
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
