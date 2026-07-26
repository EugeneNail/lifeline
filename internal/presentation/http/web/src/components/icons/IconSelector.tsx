import { EmojiIcon } from './EmojiIcon'
import { emojiIconSections } from './emojiIconCatalog'
import type { EmojiIcon as EmojiIcons } from './emojiIconCatalog'
import './IconSelector.sass'

type IconSelectorProps = {
    value: EmojiIcons
    onChange: (icon: EmojiIcons) => void
}

// IconSelector renders grouped clickable emoji icons for choosing a habit icon.
export function IconSelector({ value, onChange }: IconSelectorProps) {
    return (
        <div className="icon-selector">
            {emojiIconSections.map((section) => (
                <section className="icon-selector__section" key={section.name}>
                    <h3 className="icon-selector__section-title">{section.name}</h3>

                    <div className="icon-selector__grid">
                        {section.icons.map(({ icon, label }) => (
                            <button
                                aria-label={`Select icon ${section.name}: ${label}`}
                                aria-pressed={value === icon}
                                className="icon-selector__option"
                                key={icon}
                                onClick={() => onChange(icon)}
                                type="button"
                            >
                                <EmojiIcon icon={icon} size={28} />
                            </button>
                        ))}
                    </div>
                </section>
            ))}
        </div>
    )
}
