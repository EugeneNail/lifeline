import axios, { AxiosError, AxiosHeaders } from 'axios'
import type { InternalAxiosRequestConfig } from 'axios'
import {
    clearAuthTokens,
    getLoginToken,
    getRefreshToken,
    setLoginToken,
} from './auth-tokens'

type RetriableRequestConfig = {
    _retry?: boolean
}

type AuthenticatedRequestConfig = InternalAxiosRequestConfig & RetriableRequestConfig

const guestRoutes = new Set(['users/login', 'users/register', 'users/refresh'])

const apiBaseURL = '/api/v1/'

const sharedHeaders = {
    'Content-Type': 'application/json',
}

export const apiClient = axios.create({
    baseURL: apiBaseURL,
    headers: sharedHeaders,
})

const refreshClient = axios.create({
    baseURL: apiBaseURL,
    headers: sharedHeaders,
})

function normalizeRoute(url?: string) {
    return (url ?? '').replace(/^\/+/, '').replace(/^api\/v1\//, '')
}

function isGuestRoute(url?: string) {
    return guestRoutes.has(normalizeRoute(url))
}

function setAuthorizationHeader(config: InternalAxiosRequestConfig, token: string) {
    const headers = AxiosHeaders.from(config.headers)
    headers.set('Authorization', `Bearer ${token}`)
    config.headers = headers
}

function redirectToLogin() {
    clearAuthTokens()
    window.location.replace('/login')
}

let refreshInFlight: Promise<string | null> | null = null

async function refreshLoginToken() {
    const refreshToken = getRefreshToken()
    if (!refreshToken) {
        return null
    }

    const response = await refreshClient.post<string>('users/refresh', {
        refreshToken,
    })

    if (typeof response.data !== 'string' || response.data.length === 0) {
        return null
    }

    setLoginToken(response.data)
    return response.data
}

function getRefreshPromise() {
    if (refreshInFlight === null) {
        refreshInFlight = refreshLoginToken().finally(() => {
            refreshInFlight = null
        })
    }

    return refreshInFlight
}

apiClient.interceptors.request.use((config) => {
    if (isGuestRoute(config.url)) {
        return config
    }

    const loginToken = getLoginToken()
    if (loginToken) {
        setAuthorizationHeader(config, loginToken)
    }

    return config
})

apiClient.interceptors.response.use(
    (response) => response,
    async (error: AxiosError) => {
        const status = error.response?.status
        const request = error.config as AuthenticatedRequestConfig | undefined

        if (status !== 401 || !request || isGuestRoute(request.url)) {
            return Promise.reject(error)
        }

        if (request._retry) {
            redirectToLogin()
            return Promise.reject(error)
        }

        request._retry = true

        try {
            const refreshedLoginToken = await getRefreshPromise()
            if (!refreshedLoginToken) {
                redirectToLogin()
                return Promise.reject(error)
            }

            return apiClient(request)
        } catch (refreshError) {
            if (axios.isAxiosError(refreshError) && refreshError.response?.status === 401) {
                redirectToLogin()
                return Promise.reject(error)
            }

            redirectToLogin()
            return Promise.reject(refreshError)
        }
    },
)
