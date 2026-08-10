import { Link, NavLink } from 'react-router-dom'
import './AppNavigation.sass'

type NavigationLink = {
    icon: string
    label: string
    to: string
    end?: boolean
}

const navigationLinks: NavigationLink[] = [
    { icon: '⌂', label: 'Today', to: '/', end: true },
    { icon: '✎', label: 'Diary', to: '/journals' },
    { icon: '✓', label: 'Habits', to: '/habits' },
    { icon: '↕', label: 'Transactions', to: '/transactions/statistics' },
]

// AppNavigation renders the fixed application sidebar used across the web app.
export function AppNavigation() {
    return (
        <aside className="app-navigation" aria-label="Application sidebar">
            <Link className="app-navigation__brand" to="/">
                <span className="app-navigation__brand-mark">L</span>
                <span className="app-navigation__brand-label">Lifeline</span>
            </Link>

            <nav className="app-navigation__primary-nav" aria-label="Primary navigation">
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
                        <span className="app-navigation__icon" aria-hidden="true">
                            {link.icon}
                        </span>
                        <span className="app-navigation__label">{link.label}</span>
                    </NavLink>
                ))}
            </nav>

            <div className="app-navigation__footer">
                <Link className="app-navigation__settings-link" to="/settings">
                    <span className="app-navigation__settings-icon" aria-hidden="true">
                        ⚙
                    </span>
                    <span>Settings</span>
                </Link>

                <div className="app-navigation__profile">
                    <div className="app-navigation__profile-main">
                        <div className="app-navigation__avatar">PU</div>
                        <div className="app-navigation__profile-text">
                            <strong>Placeholder User</strong>
                            <span>Personal space</span>
                        </div>
                    </div>
                </div>
            </div>
        </aside>
    )
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
