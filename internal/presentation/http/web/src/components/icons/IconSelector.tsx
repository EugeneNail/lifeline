import { GoogleIcon } from './GoogleIcon'
import { googleIconNames, googleIconSections } from './googleIconCatalog'
import type { GoogleIcons } from './googleIconCatalog'
import './IconSelector.sass'

type IconSelectorProps = {
    value: GoogleIcons
    onChange: (icon: GoogleIcons) => void
}

// IconSelector renders grouped clickable Google icons for choosing a habit icon.
export function IconSelector({ value, onChange }: IconSelectorProps) {
    return (
        <div className="icon-selector">
            {googleIconSections.map((section) => (
                <section className="icon-selector__section" key={section.name}>
                    <h3 className="icon-selector__section-title">{section.name}</h3>

                    <div className="icon-selector__grid">
                        {section.icons.map((icon) => (
                            <button
                                aria-label={`Select icon ${googleIconNames[icon]}`}
                                aria-pressed={value === icon}
                                className="icon-selector__option"
                                key={icon}
                                onClick={() => onChange(icon)}
                                type="button"
                            >
                                <GoogleIcon icon={icon} size={28} />
                            </button>
                        ))}
                    </div>
                </section>
            ))}
        </div>
    )
}
