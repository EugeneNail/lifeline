export type TransactionCategory = {
    id: number
    icon: string
    name: string
}

export const transactionCategories: TransactionCategory[] = [
    { id: 1, icon: '💡', name: 'Bills' },
    { id: 2, icon: '🍽️', name: 'Food' },
    { id: 3, icon: '🚕', name: 'Transport' },
    { id: 4, icon: '🏠', name: 'Household' },
    { id: 5, icon: '🎬', name: 'Entertainment' },
    { id: 6, icon: '🎒', name: 'Personal items' },
    { id: 7, icon: '🩺', name: 'Health' },
    { id: 8, icon: '💼', name: 'Work' },
    { id: 9, icon: '💸', name: 'Debt' },
    { id: 10, icon: '📈', name: 'Investments' },
    { id: 11, icon: '🎁', name: 'Gifts' },
    { id: 12, icon: '✨', name: 'Other' },
]

const allTransactionCategoryIds = transactionCategories.map((category) => category.id)
const transactionCategoriesStorageKey = 'lifeline.transactionStatistics.categories'

// readStoredTransactionCategories returns the stored selection or all categories when storage contains no valid value.
export function readStoredTransactionCategories() {
    try {
        const rawCategories = window.localStorage.getItem(transactionCategoriesStorageKey)
        if (rawCategories === null) {
            return [...allTransactionCategoryIds]
        }

        const categories: unknown = JSON.parse(rawCategories)
        if (
            !Array.isArray(categories) ||
            !categories.every((category) => (
                typeof category === 'number' && allTransactionCategoryIds.includes(category)
            ))
        ) {
            return [...allTransactionCategoryIds]
        }

        return [...new Set(categories)]
    } catch {
        return [...allTransactionCategoryIds]
    }
}

// storeTransactionCategories saves the category selection when browser storage is available.
export function storeTransactionCategories(categories: number[]) {
    try {
        window.localStorage.setItem(transactionCategoriesStorageKey, JSON.stringify(categories))
    } catch {
        // The in-memory selection remains usable when browser storage is unavailable.
    }
}
