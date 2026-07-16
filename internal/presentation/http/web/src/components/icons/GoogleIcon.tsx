import type { HTMLAttributes } from 'react'
import './GoogleIcon.sass'

// GoogleIcons identifies the supported Material Symbols placeholder icons.
export const GoogleIcons = {
    Check: 1,
    Water: 2,
    Numbers: 3,
    Schedule: 4,
    Fitness: 5,
} as const

export type GoogleIcons = (typeof GoogleIcons)[keyof typeof GoogleIcons]

const googleIconNames: Record<GoogleIcons, string> = {
    [GoogleIcons.Check]: 'check_circle',
    [GoogleIcons.Water]: 'water_drop',
    [GoogleIcons.Numbers]: 'pin',
    [GoogleIcons.Schedule]: 'schedule',
    [GoogleIcons.Fitness]: 'exercise',
}

type GoogleIconProps = HTMLAttributes<HTMLSpanElement> & {
    icon: GoogleIcons
    size: number
    title?: string
}

// GoogleIcon renders one Material Symbols icon from the local GoogleIcons enum.
export function GoogleIcon({ icon, size, title, className, style, ...props }: GoogleIconProps) {
    const iconSize = `${size}px`

    return (
        <span
            aria-hidden={title ? undefined : true}
            aria-label={title}
            className={joinClassNames('google-icon', 'material-symbols-rounded', className)}
            role={title ? 'img' : undefined}
            style={{ ...style, width: iconSize, height: iconSize, fontSize: iconSize }}
            {...props}
        >
            {googleIconNames[icon]}
        </span>
    )
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
