import './SavingStatus.sass'

type SavingStatusKind = 'saved' | 'saving' | 'error'

type SavingStatusProps = {
    status: SavingStatusKind
    className?: string
}

// SavingStatus renders the shared autosave badge with the dot-and-pill pattern.
export function SavingStatus({ status, className }: SavingStatusProps) {
    return (
        <span className={joinClassNames('saving-status', className)} data-status={status}>
            {statusText(status)}
        </span>
    )
}

function statusText(status: SavingStatusKind) {
    if (status === 'saving') {
        return 'Saving...'
    }

    if (status === 'error') {
        return 'Save failed'
    }

    return 'Saved'
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
