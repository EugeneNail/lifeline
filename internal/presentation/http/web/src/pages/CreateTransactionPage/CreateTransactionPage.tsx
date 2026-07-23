import axios from 'axios'
import type {FormEvent, KeyboardEvent} from 'react'
import {useEffect, useMemo, useState} from 'react'
import {useNavigate, useSearchParams} from 'react-router-dom'
import {DateSelector} from '../../components/date'
import {AppNavigation} from '../../components/navigation'
import {GoogleIcon, GoogleIcons} from '../../components/icons'
import {Button, Message} from '../../components/primitives'
import {Page, PageHeader, Panel, PanelBody, Section, SectionHeader} from '../../components/layout'
import {useApiClient} from '../../hooks/useApiClient'
import './CreateTransactionPage.sass'

type TransactionDirection = 'expense' | 'income'

type TransactionCategory = {
    value: number
    icon: string
    title: string
}

const transactionCategories: TransactionCategory[] = [
    {value: 1, icon: '💡', title: 'Bills'},
    {value: 2, icon: '🍽️', title: 'Food'},
    {value: 3, icon: '🚕', title: 'Transport'},
    {value: 4, icon: '🏠', title: 'Household'},
    {value: 5, icon: '🎬', title: 'Entertainment'},
    {value: 6, icon: '🎒', title: 'Personal items'},
    {value: 7, icon: '🩺', title: 'Health'},
    {value: 8, icon: '💼', title: 'Work'},
    {value: 9, icon: '💸', title: 'Debt'},
    {value: 10, icon: '📈', title: 'Investments'},
    {value: 11, icon: '🎁', title: 'Gifts'},
    {value: 12, icon: '✨', title: 'Other'},
]

function startOfDay(date: Date) {
    const nextDate = new Date(date)
    nextDate.setHours(0, 0, 0, 0)
    return nextDate
}

function resolveQueryDate(rawDate: string | null) {
    if (!rawDate) {
        return null
    }

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

function formatDateLabel(date: Date) {
    return new Intl.DateTimeFormat('en-US', {
        weekday: 'long',
        day: 'numeric',
        month: 'long',
        year: 'numeric',
    }).format(date)
}

function formatDateFieldValue(date: Date) {
    return `${date.getFullYear()}-${String(date.getMonth() + 1).padStart(2, '0')}-${String(date.getDate()).padStart(2, '0')}`
}

function parseDecimal(value: string) {
    return value.replace(/[^\d.,\s]/g, '').replace(/([.,].*)[.,]/g, '$1')
}

// CreateTransactionPage renders the transaction creation form and date selector.
export function CreateTransactionPage() {
    const apiClient = useApiClient()
    const navigate = useNavigate()
    const [searchParams] = useSearchParams()
    const queryDateValue = searchParams.get('date')
    const today = useMemo(() => startOfDay(new Date()), [])
    const [direction, setDirection] = useState<TransactionDirection>('expense')
    const [selectedDate, setSelectedDate] = useState<Date>(() => {
        return resolveQueryDate(queryDateValue) ?? today
    })
    const [selectedCategory, setSelectedCategory] = useState<number>(1)
    const [amount, setAmount] = useState('')
    const [description, setDescription] = useState('')
    const [isDateSelectorOpen, setDateSelectorOpen] = useState(false)
    const [isSubmitting, setIsSubmitting] = useState(false)
    const [formError, setFormError] = useState('')

    function handleAmountKeyDown(event: KeyboardEvent<HTMLInputElement>) {
        if (event.key === '-') {
            event.preventDefault()
        }
    }

    function handleSubmit(event: FormEvent<HTMLFormElement>) {
        event.preventDefault()

        const normalizedAmount = Number(amount.replace(/\s+/g, '').replace(',', '.'))
        if (!Number.isFinite(normalizedAmount) || normalizedAmount <= 0) {
            setFormError('Enter a valid amount.')
            return
        }

        setIsSubmitting(true)
        setFormError('')

        void apiClient
            .post('transactions', {
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
                if (axios.isAxiosError(error) && error.response?.status === 422) {
                    const response = error.response.data
                    if (response && typeof response === 'object') {
                        setFormError(Object.values(response as Record<string, string>).join(' '))
                        return
                    }
                }

                setFormError('Could not create transaction.')
            })
            .finally(() => {
                setIsSubmitting(false)
            })
    }

    const selectedDateLabel = selectedDate.getTime() === today.getTime() ? 'Today' : formatDateLabel(selectedDate)

    useEffect(() => {
        const queryDate = resolveQueryDate(queryDateValue)
        if (queryDate) {
            setSelectedDate(queryDate)
        }
    }, [queryDateValue])

    return (
        <Page className="create-transaction-page">
            <PageHeader
                eyebrow="Transactions"
                title="New transaction"
                subtitle="Record an expense or income. Enter the amount as a positive value and choose the date separately."
                actions={
                    <Button type="button" variant="secondary" onClick={() => navigate(-1)}>
                        Back
                    </Button>
                }
            />

            <Panel>
                <PanelBody>
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
                        </Section>

                        <Section>
                            <Button
                                className="create-transaction-page__date-button"
                                type="button"
                                variant="secondary"
                                onClick={() => setDateSelectorOpen(true)}
                            >
                                <GoogleIcon icon={GoogleIcons.CalendarMonth} size={18}/>
                                {selectedDateLabel}
                            </Button>
                            <div className="create-transaction-page__helper-text">
                                {formatDateFieldValue(selectedDate)}
                            </div>
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
                        </Section>

                        <Section>
                            <SectionHeader title="Description" meta={`${description.length} / 32`}/>
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
                        </div>
                    </form>
                </PanelBody>
            </Panel>

            <AppNavigation/>

            <DateSelector
                mode="single"
                open={isDateSelectorOpen}
                value={selectedDate}
                onChange={(date) => {
                    setSelectedDate(startOfDay(date))
                    setDateSelectorOpen(false)
                }}
                onClose={() => setDateSelectorOpen(false)}
            />
        </Page>
    )
}
