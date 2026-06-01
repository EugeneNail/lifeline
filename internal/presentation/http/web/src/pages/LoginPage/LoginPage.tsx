import axios from 'axios'
import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { setAuthTokens } from '../../api/auth-tokens'
import { useApiClient } from '../../hooks/useApiClient'
import './LoginPage.sass'

type LoginResponse = {
    loginToken: string
    refreshToken: string
}

type LoginFieldErrors = Partial<Record<'email' | 'password', string>>

export function LoginPage() {
    const apiClient = useApiClient()
    const navigate = useNavigate()
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [fieldErrors, setFieldErrors] = useState<LoginFieldErrors>({})

    function setLoginFieldErrors(error: unknown) {
        if (!axios.isAxiosError(error)) {
            return false
        }

        if (error.response?.data && typeof error.response.data === 'object') {
            const response = error.response.data as Record<string, string>
            setFieldErrors({
                email: response.email,
                password: response.password,
            })
            return true
        }

        return false
    }

    async function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()
        setIsSubmitting(true)
        setFieldErrors({})

        try {
            const response = await apiClient.post<LoginResponse>('users/login', {
                email,
                password,
            })

            setAuthTokens(response.data.loginToken, response.data.refreshToken)
            navigate('/')
        } catch (error) {
            setLoginFieldErrors(error)
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <main className="login-page">
            <section className="login-shell">
                <div className="login-brand">
                    <p className="login-kicker">Lifeline</p>
                    <h1>Welcome back</h1>
                    <p className="login-lead">
                        Sign in to continue tracking habits, entries, and daily progress.
                    </p>
                </div>

                <form className="login-panel" onSubmit={handleSubmit}>
                    <label className="login-field">
                        <span>Email</span>
                        <input
                            className="login-input"
                            type="email"
                            name="email"
                            autoComplete="email"
                            placeholder="you@example.com"
                            value={email}
                            onChange={(event) => setEmail(event.target.value)}
                        />
                        <span className="login-field-error" aria-live="polite">
                            {fieldErrors.email || '\u00A0'}
                        </span>
                    </label>

                    <label className="login-field">
                        <span>Password</span>
                        <input
                            className="login-input"
                            type="password"
                            name="password"
                            autoComplete="current-password"
                            placeholder="Enter your password"
                            value={password}
                            onChange={(event) => setPassword(event.target.value)}
                        />
                        <span className="login-field-error" aria-live="polite">
                            {fieldErrors.password || '\u00A0'}
                        </span>
                    </label>

                    <button className="login-button" type="submit" disabled={isSubmitting}>
                        {isSubmitting ? 'Signing in...' : 'Sign in'}
                    </button>

                    <p className="login-switch">
                        No account yet? <Link to="/signup">Create one</Link>
                    </p>
                </form>
            </section>
        </main>
    )
}
