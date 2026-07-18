import type { HTMLAttributes } from 'react'
import { googleIconNames } from './googleIconCatalog'
import type { GoogleIcons } from './googleIconCatalog'
import './GoogleIcon.sass'

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
