import axios from 'axios'
import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { setAuthTokens } from '../../api/auth-tokens'
import { useApiClient } from '../../hooks/useApiClient'
import './SignupPage.sass'

type SignupResponse = string
type LoginResponse = {
    loginToken: string
    refreshToken: string
}

type SignupFieldErrors = Partial<Record<'email' | 'password' | 'passwordConfirmation', string>>

export function SignupPage() {
    const apiClient = useApiClient()
    const navigate = useNavigate()
    const [email, setEmail] = useState('')
    const [password, setPassword] = useState('')
    const [passwordConfirmation, setPasswordConfirmation] = useState('')
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [fieldErrors, setFieldErrors] = useState<SignupFieldErrors>({})

    function setRegistrationFieldErrors(error: unknown) {
        if (!axios.isAxiosError(error)) {
            return false
        }

        if (error.response?.status === 409) {
            setFieldErrors({
                email: 'Email already taken',
            })
            return true
        }

        if (error.response?.data && typeof error.response.data === 'object') {
            const response = error.response.data as Record<string, string>
            setFieldErrors({
                email: response.email,
                password: response.password,
                passwordConfirmation: response.passwordConfirmation,
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
            await apiClient.post<SignupResponse>('users/register', {
                email,
                password,
                passwordConfirmation,
            })

            const loginResponse = await apiClient.post<LoginResponse>('users/login', {
                email,
                password,
            })

            setAuthTokens(loginResponse.data.loginToken, loginResponse.data.refreshToken)
            navigate('/')
        } catch (error) {
            setRegistrationFieldErrors(error)
        } finally {
            setIsSubmitting(false)
        }
    }

    return (
        <main className="signup-page">
            <section className="signup-shell">
                <div className="signup-brand">
                    <p className="signup-kicker">Lifeline</p>
                    <h1>Create your account</h1>
                    <p className="signup-lead">
                        Set up a workspace for habits, entries, and the rest of your daily
                        tracking flow.
                    </p>
                </div>

                <form className="signup-panel" onSubmit={handleSubmit}>
                    <label className="signup-field">
                        <span>Email</span>
                        <input
                            className="signup-input"
                            type="email"
                            name="email"
                            autoComplete="email"
                            placeholder="you@example.com"
                            value={email}
                            onChange={(event) => setEmail(event.target.value)}
                        />
                        <span className="signup-field-error" aria-live="polite">
                            {fieldErrors.email || '\u00A0'}
                        </span>
                    </label>

                    <label className="signup-field">
                        <span>Password</span>
                        <input
                            className="signup-input"
                            type="password"
                            name="password"
                            autoComplete="new-password"
                            placeholder="Create a strong password"
                            value={password}
                            onChange={(event) => setPassword(event.target.value)}
                        />
                        <span className="signup-field-error" aria-live="polite">
                            {fieldErrors.password || '\u00A0'}
                        </span>
                    </label>

                    <label className="signup-field">
                        <span>Confirm password</span>
                        <input
                            className="signup-input"
                            type="password"
                            name="passwordConfirmation"
                            autoComplete="new-password"
                            placeholder="Repeat the password"
                            value={passwordConfirmation}
                            onChange={(event) => setPasswordConfirmation(event.target.value)}
                        />
                        <span className="signup-field-error" aria-live="polite">
                            {fieldErrors.passwordConfirmation || '\u00A0'}
                        </span>
                    </label>

                    <button className="signup-button" type="submit" disabled={isSubmitting}>
                        {isSubmitting ? 'Signing up...' : 'Sign up'}
                    </button>

                    <p className="signup-switch">
                        Already have an account? <Link to="/login">Sign in</Link>
                    </p>
                </form>
            </section>
        </main>
    )
}
