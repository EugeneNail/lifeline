import { Link } from 'react-router-dom'
import type { GoogleIcons } from '../icons'
import { GoogleIcon } from '../icons'
import { Switch } from '../primitives'
import './HabitCard.sass'

type HabitCardProps = {
    uuid: string
    type: 'completable' | 'measurable' | 'time'
    icon: GoogleIcons
    title: string
    typeLabel: string
    enabled: boolean
    onToggle: (enabled: boolean) => void
}

// HabitCard renders the management row variant with icon, content, type, switch, and edit action.
export function HabitCard({
    uuid,
    type,
    icon,
    title,
    typeLabel,
    enabled,
    onToggle,
}: HabitCardProps) {
    return (
        <article className="habit-card">
            <div className="habit-card__icon">
                <GoogleIcon icon={icon} size={22} />
            </div>

            <div className="habit-card__content">
                <h3 className="habit-card__title">{title}</h3>
            </div>

            <div className="habit-card__controls">
                <span className="habit-card__type">{typeLabel}</span>
                <Switch checked={enabled} label={`Enable ${title}`} onChange={onToggle} />
            </div>

            <Link
                className="habit-card__more"
                aria-label={`Edit ${title}`}
                to={`/habits/${type}/${uuid}`}
            >
                ⋮
            </Link>
        </article>
    )
}
