const loginTokenKey = 'lifeline.auth.loginToken'
const refreshTokenKey = 'lifeline.auth.refreshToken'

function getStorage() {
    if (typeof window === 'undefined') {
        return null
    }

    return window.localStorage
}

// getLoginToken returns the stored login token or null when it is absent.
export function getLoginToken() {
    return getStorage()?.getItem(loginTokenKey) ?? null
}

// getRefreshToken returns the stored refresh token or null when it is absent.
export function getRefreshToken() {
    return getStorage()?.getItem(refreshTokenKey) ?? null
}

// setAuthTokens stores both login and refresh tokens in localStorage.
export function setAuthTokens(loginToken: string, refreshToken: string) {
    const storage = getStorage()
    if (!storage) {
        return
    }

    storage.setItem(loginTokenKey, loginToken)
    storage.setItem(refreshTokenKey, refreshToken)
}

// setLoginToken stores the refreshed login token in localStorage.
export function setLoginToken(loginToken: string) {
    const storage = getStorage()
    if (!storage) {
        return
    }

    storage.setItem(loginTokenKey, loginToken)
}

// clearAuthTokens removes both tokens from localStorage.
export function clearAuthTokens() {
    const storage = getStorage()
    if (!storage) {
        return
    }

    storage.removeItem(loginTokenKey)
    storage.removeItem(refreshTokenKey)
}
