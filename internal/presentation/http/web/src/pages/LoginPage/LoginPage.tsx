import axios from 'axios'
import { useState } from 'react'
import type { FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { setAuthTokens } from '../../api/auth-tokens'
import { Page, PageHeader, Panel, PanelBody } from '../../components/layout'
import { Button, TextField } from '../../components/primitives'
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
        <Page className="login-page">
            <PageHeader
                eyebrow="Lifeline"
                title="Welcome back"
                subtitle="Sign in to continue tracking habits, journals, and daily progress."
            />

            <Panel>
                <PanelBody>
                    <form className="login-form" onSubmit={handleSubmit}>
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
                            autoComplete="current-password"
                            placeholder="Enter your password"
                            value={password}
                            error={fieldErrors.password}
                            onChange={(event) => setPassword(event.target.value)}
                        />

                        <Button type="submit" disabled={isSubmitting}>
                            {isSubmitting ? 'Signing in...' : 'Sign in'}
                        </Button>

                        <p className="login-switch">
                            No account yet? <Link to="/signup">Create one</Link>
                        </p>
                    </form>
                </PanelBody>
            </Panel>
        </Page>
    )
}
