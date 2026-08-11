import { useEffect, useRef, useState } from 'react'
import {
    readStoredTransactionCategories,
    storeTransactionCategories,
    transactionCategories,
    type TransactionCategory,
} from './transactionCategories'
import './TransactionCategorySelector.sass'

type TransactionCategorySelectorProps = {
    onChange: (categories: number[]) => void
}

function getSelectionLabel(selectedCategoryIds: number[]) {
    if (selectedCategoryIds.length === 0) {
        return 'All categories'
    }

    const selectedCategories = selectedCategoryIds
        .map((categoryId) => transactionCategories.find((category) => category.id === categoryId))
        .filter((category): category is TransactionCategory => category !== undefined)

    if (selectedCategories.length === 1) {
        return selectedCategories[0].name
    }

    if (selectedCategories.length === 2) {
        const [firstCategory, secondCategory] = selectedCategories
        const secondName = secondCategory.name.length > 7
            ? `${secondCategory.name.slice(0, 7)}...`
            : secondCategory.name

        return `${firstCategory.name}, ${secondName}`
    }

    return `${selectedCategories.length} categories`
}

// TransactionCategorySelector renders a compact multi-select and applies category changes after a short delay.
export function TransactionCategorySelector({ onChange }: TransactionCategorySelectorProps) {
    const [isOpen, setIsOpen] = useState(false)
    const [selectedCategoryIds, setSelectedCategoryIds] = useState<number[]>(readStoredTransactionCategories)
    const rootRef = useRef<HTMLDivElement>(null)
    const updateTimeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null)

    useEffect(() => {
        return () => {
            if (updateTimeoutRef.current !== null) {
                clearTimeout(updateTimeoutRef.current)
            }
        }
    }, [])

    useEffect(() => {
        if (!isOpen) {
            return
        }

        function handlePointerDown(event: PointerEvent) {
            if (rootRef.current?.contains(event.target as Node)) {
                return
            }

            setIsOpen(false)
        }

        function handleKeyDown(event: KeyboardEvent) {
            if (event.key === 'Escape') {
                setIsOpen(false)
            }
        }

        document.addEventListener('pointerdown', handlePointerDown)
        window.addEventListener('keydown', handleKeyDown)

        return () => {
            document.removeEventListener('pointerdown', handlePointerDown)
            window.removeEventListener('keydown', handleKeyDown)
        }
    }, [isOpen])

    function handleCategoryChange(categoryId: number) {
        const nextCategoryIds = selectedCategoryIds.includes(categoryId)
            ? selectedCategoryIds.filter((selectedCategoryId) => selectedCategoryId !== categoryId)
            : [...selectedCategoryIds, categoryId].sort((left, right) => left - right)

        setSelectedCategoryIds(nextCategoryIds)
        storeTransactionCategories(nextCategoryIds)

        if (updateTimeoutRef.current !== null) {
            clearTimeout(updateTimeoutRef.current)
        }

        updateTimeoutRef.current = setTimeout(() => {
            onChange(nextCategoryIds)
            updateTimeoutRef.current = null
        }, 1000)
    }

    return (
        <div className="transaction-category-selector" ref={rootRef}>
            <button
                aria-expanded={isOpen}
                aria-haspopup="listbox"
                className="transaction-category-selector__button"
                type="button"
                onClick={() => setIsOpen((currentValue) => !currentValue)}
            >
                <span className="transaction-category-selector__button-label">
                    {getSelectionLabel(selectedCategoryIds)}
                </span>
                <span aria-hidden="true" className="transaction-category-selector__chevron">⌄</span>
            </button>

            {isOpen ? (
                <div
                    aria-label="Transaction categories"
                    className="transaction-category-selector__menu"
                    role="listbox"
                >
                    {transactionCategories.map((category) => (
                        <label className="transaction-category-selector__option" key={category.id}>
                            <input
                                checked={selectedCategoryIds.includes(category.id)}
                                type="checkbox"
                                onChange={() => handleCategoryChange(category.id)}
                            />
                            <span>{category.name}</span>
                        </label>
                    ))}
                </div>
            ) : null}
        </div>
    )
}
