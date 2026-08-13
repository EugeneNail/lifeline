import type { DateRange } from './DateSelector'
import './PeriodSelector.sass'

type PeriodSelectorProps = {
    ariaLabel?: string
    onChange: (range: DateRange) => void
    onCustomClick: () => void
    today: Date
    value: DateRange
}

type PeriodPreset = {
    days: 7 | 30 | 365
    label: 'Week' | 'Month' | 'Year'
}

const periodPresets: PeriodPreset[] = [
    { days: 7, label: 'Week' },
    { days: 30, label: 'Month' },
    { days: 365, label: 'Year' },
]

const millisecondsPerDay = 24 * 60 * 60 * 1000

function startOfDay(date: Date) {
    const normalizedDate = new Date(date)
    normalizedDate.setHours(0, 0, 0, 0)
    return normalizedDate
}

function getDateNumber(date: Date) {
    return Date.UTC(date.getFullYear(), date.getMonth(), date.getDate())
}

function getRangeLength(range: DateRange) {
    return Math.floor(
        (getDateNumber(range.endDate) - getDateNumber(range.startDate)) / millisecondsPerDay,
    ) + 1
}

function buildPresetRange(days: number, today: Date): DateRange {
    const endDate = startOfDay(today)
    const startDate = new Date(endDate)
    startDate.setDate(startDate.getDate() - (days - 1))

    return { startDate, endDate }
}

// PeriodSelector renders compact range presets and opens the custom date picker on request.
export function PeriodSelector({
    ariaLabel = 'Statistics period',
    onChange,
    onCustomClick,
    today,
    value,
}: PeriodSelectorProps) {
    const rangeLength = getRangeLength(value)
    const activePreset = periodPresets.find((preset) => preset.days === rangeLength)

    return (
        <div aria-label={ariaLabel} className="period-selector">
            {periodPresets.map((preset) => (
                <button
                    aria-pressed={activePreset?.days === preset.days}
                    className="period-selector__button"
                    key={preset.days}
                    type="button"
                    onClick={() => onChange(buildPresetRange(preset.days, today))}
                >
                    {preset.label}
                </button>
            ))}
            <button
                aria-pressed={activePreset === undefined}
                className="period-selector__button"
                type="button"
                onClick={onCustomClick}
            >
                Custom
            </button>
        </div>
    )
}
