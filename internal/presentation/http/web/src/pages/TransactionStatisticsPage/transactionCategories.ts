export type TransactionCategory = {
    id: number
    name: string
}

export const transactionCategories: TransactionCategory[] = [
    { id: 1, name: 'Bills' },
    { id: 2, name: 'Food' },
    { id: 3, name: 'Transport' },
    { id: 4, name: 'Household' },
    { id: 5, name: 'Entertainment' },
    { id: 6, name: 'Personal items' },
    { id: 7, name: 'Health' },
    { id: 8, name: 'Work' },
    { id: 9, name: 'Debt' },
    { id: 10, name: 'Investments' },
    { id: 11, name: 'Gifts' },
    { id: 12, name: 'Other' },
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
