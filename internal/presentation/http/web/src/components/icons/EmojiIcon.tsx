import type { HTMLAttributes } from 'react'
import { emojiIconNames } from './emojiIconCatalog'
import type { EmojiIcon as EmojiIconValue } from './emojiIconCatalog'
import './GoogleIcon.sass'

type EmojiIconProps = HTMLAttributes<HTMLSpanElement> & {
    icon: EmojiIconValue
    size: number
    title?: string
}

// EmojiIcon renders one emoji from the local icon catalog.
export function EmojiIcon({ icon, size, title, className, style, ...props }: EmojiIconProps) {
    const iconSize = `${size}px`

    return (
        <span
            aria-hidden={title ? undefined : true}
            aria-label={title}
            className={joinClassNames('google-icon', 'emoji-icon', className)}
            role={title ? 'img' : undefined}
            style={{ ...style, width: iconSize, height: iconSize, fontSize: iconSize }}
            {...props}
        >
            {emojiIconNames[icon]}
        </span>
    )
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
