import { useEffect, useState } from 'react'
import { createPortal } from 'react-dom'
import { Link, NavLink } from 'react-router-dom'
import { GoogleIcon } from '../icons'
import './AppNavigation.sass'

type NavigationLink = {
    icon: string
    label: string
    to: string
    end?: boolean
}

const primaryNavigationLinks: NavigationLink[] = [
    { icon: 'today', label: 'Today', to: '/today', end: true },
    { icon: 'menu_book', label: 'Journals', to: '/journals', end: true },
    { icon: 'receipt_long', label: 'Transactions', to: '/transactions', end: true },
    { icon: 'task_alt', label: 'Habits', to: '/habits', end: true },
]

const secondaryNavigationLinks: NavigationLink[] = [
    { icon: 'monitoring', label: 'Transaction statistics', to: '/transactions/statistics', end: true },
    { icon: 'grid_on', label: 'Habit statistics', to: '/habits/statistics', end: true },
    { icon: 'settings', label: 'Settings', to: '/settings', end: true },
]

// AppNavigation renders the desktop sidebar and the mobile bottom navigation with an additional full-screen menu.
export function AppNavigation() {
    const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false)

    useEffect(() => {
        if (!isMobileMenuOpen) {
            return
        }

        const previousOverflow = document.body.style.overflow

        function handleKeyDown(event: KeyboardEvent) {
            if (event.key === 'Escape') {
                setIsMobileMenuOpen(false)
            }
        }

        document.body.style.overflow = 'hidden'
        window.addEventListener('keydown', handleKeyDown)

        const portraitMediaQuery = window.matchMedia('(orientation: portrait)')

        function handleOrientationChange(event: MediaQueryListEvent) {
            if (!event.matches) {
                setIsMobileMenuOpen(false)
            }
        }

        portraitMediaQuery.addEventListener('change', handleOrientationChange)

        return () => {
            document.body.style.overflow = previousOverflow
            window.removeEventListener('keydown', handleKeyDown)
            portraitMediaQuery.removeEventListener('change', handleOrientationChange)
        }
    }, [isMobileMenuOpen])

    return (
        <>
            <aside className="app-navigation" aria-label="Application navigation">
                <Link className="app-navigation__brand" to="/today">
                    <span className="app-navigation__brand-mark">L</span>
                    <span className="app-navigation__brand-label">Lifeline</span>
                </Link>

                <div className="app-navigation__menus">
                    <nav className="app-navigation__primary-nav" aria-label="Primary navigation">
                        {primaryNavigationLinks.map((link) => (
                            <NavigationItem key={link.to} link={link} />
                        ))}
                        <button
                            aria-controls="app-navigation-mobile-menu"
                            aria-expanded={isMobileMenuOpen}
                            aria-label={isMobileMenuOpen ? 'Close additional navigation' : 'Open additional navigation'}
                            className="app-navigation__burger"
                            type="button"
                            onClick={() => setIsMobileMenuOpen((isOpen) => !isOpen)}
                        >
                            <GoogleIcon
                                icon={isMobileMenuOpen ? 'close' : 'menu'}
                                size={24}
                            />
                            <span className="app-navigation__label">More</span>
                        </button>
                    </nav>

                    <nav className="app-navigation__secondary-nav" aria-label="Additional navigation">
                        {secondaryNavigationLinks.map((link) => (
                            <NavigationItem key={link.to} link={link} />
                        ))}
                    </nav>
                </div>

                <div className="app-navigation__footer">
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

            {isMobileMenuOpen
                ? createPortal(
                      <section
                          aria-label="Additional navigation"
                          aria-modal="true"
                          className="app-navigation__mobile-menu"
                          id="app-navigation-mobile-menu"
                          role="dialog"
                      >
                          <div className="app-navigation__mobile-menu-content">
                              <p className="app-navigation__mobile-menu-eyebrow">Navigation</p>
                              <h2 className="app-navigation__mobile-menu-title">More</h2>
                              <nav
                                  aria-label="Mobile additional navigation"
                                  className="app-navigation__mobile-secondary-nav"
                              >
                                  {secondaryNavigationLinks.map((link) => (
                                      <NavigationItem
                                          key={link.to}
                                          link={link}
                                          mobile
                                          onNavigate={() => setIsMobileMenuOpen(false)}
                                      />
                                  ))}
                              </nav>
                          </div>
                      </section>,
                      document.body,
                  )
                : null}
        </>
    )
}

type NavigationItemProps = {
    link: NavigationLink
    mobile?: boolean
    onNavigate?: () => void
}

function NavigationItem({ link, mobile = false, onNavigate }: NavigationItemProps) {
    return (
        <NavLink
            className={({ isActive }) =>
                joinClassNames(
                    mobile ? 'app-navigation__mobile-menu-item' : 'app-navigation__item',
                    isActive
                        ? mobile
                            ? 'app-navigation__mobile-menu-item--active'
                            : 'app-navigation__item--active'
                        : undefined,
                )
            }
            end={link.end}
            to={link.to}
            onClick={onNavigate}
        >
            <span className="app-navigation__icon" aria-hidden="true">
                <GoogleIcon icon={link.icon} size={20} />
            </span>
            <span className="app-navigation__label">{link.label}</span>
        </NavLink>
    )
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
