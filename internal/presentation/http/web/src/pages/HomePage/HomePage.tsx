import { Link } from 'react-router-dom'
import './HomePage.sass'

export function HomePage() {
    return (
        <main className="home-page">
            <section className="home-card">
                <p className="home-kicker">Lifeline</p>
                <h1>Session active</h1>
                <p className="home-lead">
                    Your login and refresh tokens are stored locally. Continue to the app or
                    return to the auth screens.
                </p>

                <div className="home-actions">
                    <Link to="/signup">Signup</Link>
                    <Link to="/login">Login</Link>
                </div>
            </section>
        </main>
    )
}
