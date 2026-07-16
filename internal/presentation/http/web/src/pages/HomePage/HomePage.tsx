import { Link } from 'react-router-dom'
import { Page, PageHeader, Panel, PanelBody, Section, SectionHeader } from '../../components/layout'
import { Button, Message, Metric, NavigationItem } from '../../components/primitives'
import './HomePage.sass'

export function HomePage() {
    return (
        <Page>
            <PageHeader
                eyebrow="Lifeline"
                title="Today"
                subtitle="Track habits, entries, and daily progress from one quiet workspace."
                actions={<Button type="button">+ Add entry</Button>}
            />

            <div className="home-layout">
                <Panel className="home-main-panel">
                    <PanelBody>
                        <Section>
                            <SectionHeader title="Habits" meta="3 of 5 completed" />

                            <div className="home-habit-list">
                                <article className="home-habit home-habit--selected">
                                    <div className="home-habit__icon">✓</div>
                                    <div>
                                        <div className="home-habit__name">Project work</div>
                                        <p className="home-habit__note">Completed today</p>
                                    </div>
                                    <div className="home-habit__check">✓</div>
                                </article>

                                <article className="home-habit">
                                    <div className="home-habit__icon">H₂O</div>
                                    <div>
                                        <div className="home-habit__name">Drink water</div>
                                        <p className="home-habit__note">Goal: 2000 ml</p>
                                    </div>
                                    <div className="home-habit__actions">
                                        <button className="home-icon-button" type="button">
                                            −
                                        </button>
                                        <span className="home-habit__value">1250</span>
                                        <button className="home-icon-button" type="button">
                                            +
                                        </button>
                                    </div>
                                </article>

                                <article className="home-habit">
                                    <div className="home-habit__icon">◷</div>
                                    <div>
                                        <div className="home-habit__name">Go to bed</div>
                                        <p className="home-habit__note">Time marker</p>
                                    </div>
                                    <div className="home-habit__actions">
                                        <span className="home-habit__value">00:30</span>
                                    </div>
                                </article>
                            </div>
                        </Section>

                        <Section>
                            <SectionHeader title="Daily summary" meta="Updated now" />

                            <div className="home-summary">
                                <Metric value="73%" label="Habits completed" />
                                <Metric value="8 days" label="Current streak" />
                                <Metric value="2" label="Open entries" />
                            </div>
                        </Section>
                    </PanelBody>

                    <nav className="home-navigation">
                        <NavigationItem icon="●" active>
                            Today
                        </NavigationItem>
                        <NavigationItem icon="▤">Journal</NavigationItem>
                        <NavigationItem icon="◷">Habits</NavigationItem>
                        <NavigationItem icon="⌁">Stats</NavigationItem>
                    </nav>
                </Panel>

                <aside className="home-aside">
                    <Panel>
                        <PanelBody>
                            <Section>
                                <SectionHeader title="Status" />
                                <Message variant="success" title="Session active">
                                    Your login and refresh tokens are stored locally.
                                </Message>
                            </Section>

                            <Section>
                                <SectionHeader title="Account" meta="Auth routes" />
                                <div className="home-link-list">
                                    <Link to="/login">Login</Link>
                                    <Link to="/signup">Create account</Link>
                                </div>
                            </Section>
                        </PanelBody>
                    </Panel>
                </aside>
            </div>
        </Page>
    )
}
