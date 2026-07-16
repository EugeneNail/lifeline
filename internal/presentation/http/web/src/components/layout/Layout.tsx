import type { HTMLAttributes, ReactNode } from 'react'
import './Layout.sass'

type PageProps = HTMLAttributes<HTMLElement>

// Page constrains page width and applies the standard vertical rhythm.
export function Page({ className, ...props }: PageProps) {
    return <main className={joinClassNames('page', className)} {...props} />
}

type PageHeaderProps = {
    eyebrow?: ReactNode
    title: ReactNode
    subtitle?: ReactNode
    actions?: ReactNode
} & Omit<HTMLAttributes<HTMLElement>, 'title'>

// PageHeader renders the standard page heading with optional eyebrow, subtitle, and actions.
export function PageHeader({
    eyebrow,
    title,
    subtitle,
    actions,
    className,
    ...props
}: PageHeaderProps) {
    return (
        <header className={joinClassNames('page-header', className)} {...props}>
            <div>
                {eyebrow ? <p className="page-header__eyebrow">{eyebrow}</p> : null}
                <h1 className="page-header__title">{title}</h1>
                {subtitle ? <p className="page-header__subtitle">{subtitle}</p> : null}
            </div>

            {actions ? <div className="page-header__actions">{actions}</div> : null}
        </header>
    )
}

type PanelProps = HTMLAttributes<HTMLElement>

// Panel renders the standard raised surface for page-level content groups.
export function Panel({ className, ...props }: PanelProps) {
    return <section className={joinClassNames('panel', className)} {...props} />
}

type PanelBodyProps = HTMLAttributes<HTMLDivElement>

// PanelBody applies the standard inner padding for panel content.
export function PanelBody({ className, ...props }: PanelBodyProps) {
    return <div className={joinClassNames('panel-body', className)} {...props} />
}

type SectionProps = HTMLAttributes<HTMLElement>

// Section renders a vertically separated content block inside a panel.
export function Section({ className, ...props }: SectionProps) {
    return <section className={joinClassNames('section', className)} {...props} />
}

type SectionHeaderProps = {
    title: ReactNode
    meta?: ReactNode
} & Omit<HTMLAttributes<HTMLDivElement>, 'title'>

// SectionHeader renders the standard compact title row with optional metadata.
export function SectionHeader({ title, meta, className, ...props }: SectionHeaderProps) {
    return (
        <div className={joinClassNames('section-header', className)} {...props}>
            <h2 className="section-header__title">{title}</h2>
            {meta ? <span className="section-header__meta">{meta}</span> : null}
        </div>
    )
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
