import { GoogleIcon, GoogleIcons } from './GoogleIcon'
import './IconSelector.sass'

type IconOption = {
    icon: GoogleIcons
}

type IconSection = {
    name: string
    icons: IconOption[]
}

const iconSections: IconSection[] = [
    {
        name: 'Daily basics',
        icons: [
            { icon: GoogleIcons.Check },
            { icon: GoogleIcons.Water },
            { icon: GoogleIcons.Schedule },
        ],
    },
    {
        name: 'Progress',
        icons: [
            { icon: GoogleIcons.Numbers },
            { icon: GoogleIcons.Fitness },
        ],
    },
]

type IconSelectorProps = {
    value: GoogleIcons
    onChange: (icon: GoogleIcons) => void
}

// IconSelector renders grouped clickable Google icons for choosing a habit icon.
export function IconSelector({ value, onChange }: IconSelectorProps) {
    return (
        <div className="icon-selector">
            {iconSections.map((section) => (
                <section className="icon-selector__section" key={section.name}>
                    <h3 className="icon-selector__section-title">{section.name}</h3>

                    <div className="icon-selector__grid">
                        {section.icons.map((option) => (
                            <button
                                aria-label={`Select icon ${option.icon}`}
                                aria-pressed={value === option.icon}
                                className="icon-selector__option"
                                key={option.icon}
                                onClick={() => onChange(option.icon)}
                                type="button"
                            >
                                <GoogleIcon icon={option.icon} size={28} />
                            </button>
                        ))}
                    </div>
                </section>
            ))}
        </div>
    )
}
