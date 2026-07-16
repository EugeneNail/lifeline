import { NavLink } from 'react-router-dom'
import type { ReactNode } from 'react'
import './AppNavigation.sass'

type NavigationLink = {
    icon: ReactNode
    label: string
    to: string
    end?: boolean
}

const navigationLinks: NavigationLink[] = [
    { icon: '●', label: 'Today', to: '/', end: true },
    { icon: '▤', label: 'Journal', to: '/journal' },
    { icon: '◷', label: 'Habits', to: '/habits' },
    { icon: '⌁', label: 'Stats', to: '/stats' },
]

// AppNavigation renders the primary application navigation adapted to screen orientation.
export function AppNavigation() {
    return (
        <nav className="app-navigation" aria-label="Primary navigation">
            {navigationLinks.map((link) => (
                <NavLink
                    className={({ isActive }) =>
                        joinClassNames(
                            'app-navigation__item',
                            isActive ? 'app-navigation__item--active' : undefined,
                        )
                    }
                    end={link.end}
                    key={link.to}
                    to={link.to}
                >
                    <span className="app-navigation__icon">{link.icon}</span>
                    <span className="app-navigation__label">{link.label}</span>
                </NavLink>
            ))}
        </nav>
    )
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
