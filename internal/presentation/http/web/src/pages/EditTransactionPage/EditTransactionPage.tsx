import axios from 'axios'
import type { FormEvent, KeyboardEvent } from 'react'
import { useEffect, useMemo, useState } from 'react'
import { Navigate, useNavigate, useParams } from 'react-router-dom'
import { DateSelector } from '../../components/date'
import { AppNavigation } from '../../components/navigation'
import { Button, Message } from '../../components/primitives'
import { Page, PageHeader, Panel, PanelBody, Section, SectionHeader } from '../../components/layout'
import { useApiClient } from '../../hooks/useApiClient'
import './../CreateTransactionPage/CreateTransactionPage.sass'

type TransactionDirection = 'expense' | 'income'

type TransactionCategory = {
    value: number
    icon: string
    title: string
}

type TransactionResponse = {
    id: string
    money: number
    date: string
    direction: number
    category: number
    description: string
}

type EditTransactionFieldErrors = Partial<Record<'money' | 'date' | 'direction' | 'category' | 'description', string>>

const transactionCategories: TransactionCategory[] = [
    { value: 1, icon: '💡', title: 'Bills' },
    { value: 2, icon: '🍽️', title: 'Food' },
    { value: 3, icon: '🚕', title: 'Transport' },
    { value: 4, icon: '🏠', title: 'Household' },
    { value: 5, icon: '🎬', title: 'Entertainment' },
    { value: 6, icon: '🎒', title: 'Personal items' },
    { value: 7, icon: '🩺', title: 'Health' },
    { value: 8, icon: '💼', title: 'Work' },
    { value: 9, icon: '💸', title: 'Debt' },
    { value: 10, icon: '📈', title: 'Investments' },
    { value: 11, icon: '🎁', title: 'Gifts' },
    { value: 12, icon: '✨', title: 'Other' },
]

function startOfDay(date: Date) {
    const nextDate = new Date(date)
    nextDate.setHours(0, 0, 0, 0)
    return nextDate
}

function resolveDateFieldValue(rawDate: string) {
    const match = /^(\d{4})-(\d{2})-(\d{2})$/.exec(rawDate)
    if (!match) {
        return null
    }

    const year = Number(match[1])
    const month = Number(match[2]) - 1
    const day = Number(match[3])
    const date = new Date(year, month, day)

    if (
        date.getFullYear() !== year ||
        date.getMonth() !== month ||
        date.getDate() !== day
    ) {
        return null
    }

    return startOfDay(date)
}

function formatDateFieldValue(date: Date) {
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function parseDecimal(value: string) {
    return value.replace(/[^\d.,\s]/g, '').replace(/([.,].*)[.,]/g, '$1')
}

// EditTransactionPage renders the transaction edit form and loads the transaction state from the API.
export function EditTransactionPage() {
    const apiClient = useApiClient()
    const navigate = useNavigate()
    const params = useParams<{ id: string }>()
    const transactionId = params.id
    const today = useMemo(() => startOfDay(new Date()), [])
    const [direction, setDirection] = useState<TransactionDirection>('expense')
    const [selectedDate, setSelectedDate] = useState<Date>(today)
    const [selectedCategory, setSelectedCategory] = useState<number>(1)
    const [amount, setAmount] = useState('')
    const [description, setDescription] = useState('')
    const [isLoading, setIsLoading] = useState(true)
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [fieldErrors, setFieldErrors] = useState<EditTransactionFieldErrors>({})
    const [formError, setFormError] = useState('')
    const [notFound, setNotFound] = useState(false)

    useEffect(() => {
        let isActive = true

        async function loadTransaction() {
            if (!transactionId) {
                if (isActive) {
                    setNotFound(true)
                }
                return
            }

            setIsLoading(true)
            setFormError('')
            setFieldErrors({})

            try {
                const response = await apiClient.get<TransactionResponse>(`transactions/${transactionId}`)
                if (!isActive) {
                    return
                }

                const transaction = response.data
                const loadedDate = resolveDateFieldValue(transaction.date)
                if (!loadedDate) {
                    setNotFound(true)
                    return
                }

                setDirection(transaction.direction === 2 ? 'income' : 'expense')
                setSelectedDate(loadedDate)
                setSelectedCategory(transaction.category)
                setAmount(String(transaction.money))
                setDescription(transaction.description)
            } catch (error) {
                if (!isActive) {
                    return
                }

                if (axios.isAxiosError(error) && error.response?.status === 404) {
                    setNotFound(true)
                    return
                }

                setFormError('Could not load transaction.')
            } finally {
                if (isActive) {
                    setIsLoading(false)
                }
            }
        }

        void loadTransaction()

        return () => {
            isActive = false
        }
    }, [apiClient, transactionId])

    if (notFound) {
        return <Navigate replace to="/" />
    }

    if (!transactionId) {
        return <Navigate replace to="/" />
    }

    function handleAmountKeyDown(event: KeyboardEvent<HTMLInputElement>) {
        if (event.key === '-') {
            event.preventDefault()
        }
    }

    function setEditTransactionErrors(error: unknown) {
        if (!axios.isAxiosError(error)) {
            setFormError('Could not update transaction.')
            return
        }

        if (
            error.response?.status === 422 &&
            error.response.data &&
            typeof error.response.data === 'object'
        ) {
            const response = error.response.data as Record<string, string>
            setFieldErrors({
                money: response.money,
                date: response.date,
                direction: response.direction,
                category: response.category,
                description: response.description,
            })
            return
        }

        if (error.response?.status === 404) {
            setFormError('Transaction not found.')
            return
        }

        setFormError('Could not update transaction.')
    }

    function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()

        const normalizedAmount = Number(amount.replace(/\s+/g, '').replace(',', '.'))
        if (!Number.isFinite(normalizedAmount) || normalizedAmount <= 0) {
            setFormError('Enter a valid amount.')
            return
        }

        setIsSubmitting(true)
        setFieldErrors({})
        setFormError('')

        void apiClient
            .put(`transactions/${transactionId}`, {
                money: normalizedAmount,
                date: formatDateFieldValue(selectedDate),
                direction: direction === 'income' ? 2 : 1,
                category: selectedCategory,
                description,
            })
            .then(() => {
                navigate(`/dates/${formatDateFieldValue(selectedDate)}#transactions`)
            })
            .catch((error) => {
                setEditTransactionErrors(error)
            })
            .finally(() => {
                setIsSubmitting(false)
            })
    }

    return (
        <Page className="create-transaction-page edit-transaction-page">
            <PageHeader
                eyebrow="Transactions"
                title="Edit transaction"
                subtitle="Update the amount, direction, category, description, or the transaction date."
                actions={
                    <Button type="button" variant="secondary" onClick={() => navigate(-1)}>
                        Back
                    </Button>
                }
            />

            <Panel>
                <PanelBody>
                    {isLoading ? (
                        <Message variant="info">Loading transaction...</Message>
                    ) : (
                        <form className="create-transaction-form" onSubmit={handleSubmit}>
                            <Section>
                                <div className="create-transaction-page__direction-switch" role="tablist" aria-label="Transaction direction">
                                    <button
                                        aria-pressed={direction === 'expense'}
                                        className="create-transaction-page__direction-option"
                                        data-direction="expense"
                                        type="button"
                                        onClick={() => setDirection('expense')}
                                    >
                                        ↓ Expense
                                    </button>
                                    <button
                                        aria-pressed={direction === 'income'}
                                        className="create-transaction-page__direction-option"
                                        data-direction="income"
                                        type="button"
                                        onClick={() => setDirection('income')}
                                    >
                                        ↑ Income
                                    </button>
                                </div>
                            </Section>

                            <Section>
                                <div className="create-transaction-page__amount-shell">
                                    <input
                                        aria-label="Transaction amount"
                                        className="create-transaction-page__amount-input"
                                        inputMode="decimal"
                                        name="amount"
                                        placeholder="0"
                                        type="text"
                                        value={amount}
                                        onChange={(event) => setAmount(parseDecimal(event.target.value))}
                                        onKeyDown={handleAmountKeyDown}
                                    />
                                    <div className="create-transaction-page__currency">
                                        <span className="create-transaction-page__currency-badge">₽</span>
                                        <span className="create-transaction-page__currency-label">
                                            RUB
                                        </span>
                                    </div>
                                </div>
                                <div className="create-transaction-page__helper-text">
                                    TODO: replace Ruble with a different currency label.
                                </div>
                                {fieldErrors.money ? <Message variant="error">{fieldErrors.money}</Message> : null}
                            </Section>

                            <Section>
                                <DateSelector
                                    className="create-transaction-page__date-selector"
                                    mode="single"
                                    value={selectedDate}
                                    onChange={(date) => setSelectedDate(startOfDay(date))}
                                />
                                <div className="create-transaction-page__helper-text">
                                    {formatDateFieldValue(selectedDate)}
                                </div>
                                {fieldErrors.date ? <Message variant="error">{fieldErrors.date}</Message> : null}
                            </Section>

                            <Section>
                                <div className="create-transaction-page__category-grid">
                                    {transactionCategories.map((category) => (
                                        <button
                                            aria-pressed={selectedCategory === category.value}
                                            className="create-transaction-page__category-tile"
                                            key={category.value}
                                            type="button"
                                            onClick={() => setSelectedCategory(category.value)}
                                        >
                                            <span className="create-transaction-page__category-icon">
                                                {category.icon}
                                            </span>
                                            <span className="create-transaction-page__category-title">
                                                {category.title}
                                            </span>
                                        </button>
                                    ))}
                                </div>
                                {fieldErrors.category ? <Message variant="error">{fieldErrors.category}</Message> : null}
                            </Section>

                            <Section>
                                <SectionHeader title="Description" meta={`${description.length} / 32`} />
                                <textarea
                                    aria-label="Transaction description"
                                    className="create-transaction-page__description-input"
                                    maxLength={32}
                                    name="description"
                                    placeholder="Coffee with the team"
                                    rows={4}
                                    value={description}
                                    onChange={(event) => setDescription(event.target.value)}
                                />
                                {fieldErrors.description ? <Message variant="error">{fieldErrors.description}</Message> : null}
                            </Section>

                            {formError ? (
                                <Message variant="error">{formError}</Message>
                            ) : null}

                            <div className="create-transaction-page__actions">
                                <Button
                                    className="create-transaction-page__submit-button"
                                    disabled={isSubmitting}
                                    type="submit"
                                >
                                    Save transaction
                                </Button>
                                <Button
                                    className="create-transaction-page__submit-button"
                                    type="button"
                                    variant="secondary"
                                    disabled
                                >
                                    Delete transaction
                                </Button>
                            </div>
                        </form>
                    )}
                </PanelBody>
            </Panel>
            <AppNavigation />
        </Page>
    )
}
