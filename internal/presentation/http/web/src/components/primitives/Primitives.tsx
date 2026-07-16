import type {
    ButtonHTMLAttributes,
    HTMLAttributes,
    InputHTMLAttributes,
    ReactNode,
} from 'react'
import './Primitives.sass'

type ButtonVariant = 'primary' | 'secondary' | 'danger'

type ButtonProps = ButtonHTMLAttributes<HTMLButtonElement> & {
    variant?: ButtonVariant
}

// Button renders the standard action control with primary, secondary, or danger styling.
export function Button({ variant = 'primary', className, ...props }: ButtonProps) {
    return (
        <button
            className={joinClassNames('button', `button--${variant}`, className)}
            {...props}
        />
    )
}

type TextFieldProps = InputHTMLAttributes<HTMLInputElement> & {
    label: ReactNode
    error?: ReactNode
}

// TextField renders a labeled text input with the standard focus and error presentation.
export function TextField({ label, error, className, id, ...props }: TextFieldProps) {
    const fieldId = id ?? props.name

    return (
        <label className="text-field">
            <span className="text-field__label">{label}</span>
            <input
                className={joinClassNames('text-field__input', className)}
                id={fieldId}
                aria-invalid={error ? true : undefined}
                {...props}
            />
            {error ? (
                <span className="text-field__error" aria-live="polite">
                    {error}
                </span>
            ) : null}
        </label>
    )
}

type MessageVariant = 'success' | 'warning' | 'error' | 'info'

type MessageProps = HTMLAttributes<HTMLDivElement> & {
    variant: MessageVariant
    title?: ReactNode
}

// Message renders a status block for success, warning, error, or informational feedback.
export function Message({ variant, title, children, className, ...props }: MessageProps) {
    return (
        <div className={joinClassNames('message', `message--${variant}`, className)} {...props}>
            {title ? <strong className="message__title">{title}</strong> : null}
            {children}
        </div>
    )
}

type MetricProps = HTMLAttributes<HTMLDivElement> & {
    value: ReactNode
    label: ReactNode
}

// Metric renders a compact value and label pair for dashboard summaries.
export function Metric({ value, label, className, ...props }: MetricProps) {
    return (
        <div className={joinClassNames('metric', className)} {...props}>
            <div className="metric__value">{value}</div>
            <div className="metric__label">{label}</div>
        </div>
    )
}

type NavigationItemProps = HTMLAttributes<HTMLDivElement> & {
    icon?: ReactNode
    active?: boolean
}

// NavigationItem renders a compact navigation target with optional icon and active state.
export function NavigationItem({
    icon,
    active = false,
    children,
    className,
    ...props
}: NavigationItemProps) {
    return (
        <div
            className={joinClassNames(
                'navigation-item',
                active ? 'navigation-item--active' : undefined,
                className,
            )}
            {...props}
        >
            {icon ? <span className="navigation-item__icon">{icon}</span> : null}
            {children}
        </div>
    )
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
