import { useMemo } from 'react'
import { apiClient } from '../api/client'

// useApiClient returns the shared Axios client configured for the backend API.
export function useApiClient() {
    return useMemo(() => apiClient, [])
}
