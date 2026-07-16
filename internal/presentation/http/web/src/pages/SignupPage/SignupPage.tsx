import axios from 'axios'
import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { setAuthTokens } from '../../api/auth-tokens'
import { Page, PageHeader, Panel, PanelBody } from '../../components/layout'
import { Button, TextField } from '../../components/primitives'
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
        <Page className="signup-page">
            <PageHeader
                eyebrow="Lifeline"
                title="Create your account"
                subtitle="Set up a workspace for habits, entries, and the rest of your daily tracking flow."
            />

            <Panel>
                <PanelBody>
                    <form className="signup-form" onSubmit={handleSubmit}>
                        <TextField
                            label="Email"
                            type="email"
                            name="email"
                            autoComplete="email"
                            placeholder="you@example.com"
                            value={email}
                            error={fieldErrors.email}
                            onChange={(event) => setEmail(event.target.value)}
                        />

                        <TextField
                            label="Password"
                            type="password"
                            name="password"
                            autoComplete="new-password"
                            placeholder="Create a strong password"
                            value={password}
                            error={fieldErrors.password}
                            onChange={(event) => setPassword(event.target.value)}
                        />

                        <TextField
                            label="Confirm password"
                            type="password"
                            name="passwordConfirmation"
                            autoComplete="new-password"
                            placeholder="Repeat the password"
                            value={passwordConfirmation}
                            error={fieldErrors.passwordConfirmation}
                            onChange={(event) => setPasswordConfirmation(event.target.value)}
                        />

                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? 'Signing up...' : 'Sign up'}
                        </Button>

                        <p className="signup-switch">
                            Already have an account? <Link to="/login">Sign in</Link>
                        </p>
                    </form>
                </PanelBody>
            </Panel>
        </Page>
    )
}
