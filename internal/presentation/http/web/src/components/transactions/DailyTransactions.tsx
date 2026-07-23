import { Link } from 'react-router-dom'
import './DailyTransactions.sass'

export type DailyTransaction = {
    id: string
    money: number
    date: string
    direction: number
    category: number
    description: string
}

type DailyTransactionsProps = {
    dateKey: string
    dateLabel: string
    transactions: DailyTransaction[]
}

type TransactionCategory = {
    title: string
    icon: string
}

const transactionCategories: Record<number, TransactionCategory> = {
    1: { title: 'Bills', icon: '💡' },
    2: { title: 'Food', icon: '🍽️' },
    3: { title: 'Transport', icon: '🚕' },
    4: { title: 'Household', icon: '🏠' },
    5: { title: 'Entertainment', icon: '🎬' },
    6: { title: 'Personal items', icon: '🎒' },
    7: { title: 'Health', icon: '🩺' },
    8: { title: 'Work', icon: '💼' },
    9: { title: 'Debt', icon: '💸' },
    10: { title: 'Investments', icon: '📈' },
    11: { title: 'Gifts', icon: '🎁' },
    12: { title: 'Other', icon: '✨' },
}

// DailyTransactions renders the expenses section for the selected day and returns a transaction list view.
export function DailyTransactions({ dateKey, dateLabel, transactions }: DailyTransactionsProps) {
    const totalMoney = transactions.reduce((sum, transaction) => {
        return sum + (transaction.direction === 2 ? transaction.money : -transaction.money)
    }, 0)

    return (
        <section className="daily-transactions" id="transactions">
            <div className="daily-transactions__header">
                <div>
                    <h2 className="daily-transactions__title">Expenses</h2>
                    <p className="daily-transactions__subtitle">
                        Track the spending that belongs to {dateLabel}.
                    </p>
                </div>

                <Link className="button button--primary daily-transactions__add-button" to={`/transactions/new?date=${dateKey}`}>
                    + Add transaction
                </Link>
            </div>

            <div className="daily-transactions__summary">
                <div className="daily-transactions__summary-cell">
                    <small>Expenses for the day</small>
                    <strong>{formatSummaryAmount(totalMoney)}</strong>
                </div>
                <div className="daily-transactions__summary-cell">
                    <small>Transactions</small>
                    <strong>{transactions.length}</strong>
                </div>
            </div>

            <div className="daily-transactions__list">
                {transactions.map((transaction) => {
                    const category = transactionCategories[transaction.category] ?? transactionCategories[12]
                    const description = transaction.description.trim()
                    const title = description || category.title

                    return (
                        <article className="daily-transactions__item" key={transaction.id}>
                            <div className="daily-transactions__icon">{category.icon}</div>

                            <div className="daily-transactions__content">
                                <h3 className="daily-transactions__item-title">{title}</h3>
                                {description ? <p className="daily-transactions__item-meta">{category.title}</p> : null}
                            </div>

                            <strong className="daily-transactions__amount">{formatTransactionAmount(transaction.money, transaction.direction)}</strong>

                            <Link
                                aria-label="Open transaction"
                                className="daily-transactions__menu"
                                to={`/transactions/${transaction.id}`}
                            >
                                ⋮
                            </Link>
                        </article>
                    )
                })}
            </div>
        </section>
    )
}

function formatTransactionAmount(value: number, direction: number) {
    const sign = direction === 2 ? '+' : '−'
    const normalized = new Intl.NumberFormat('ru-RU', {
        maximumFractionDigits: 2,
    }).format(value)

    return `${sign}${normalized} ₽`
}

function formatSummaryAmount(value: number) {
    const sign = value > 0 ? '+' : value < 0 ? '−' : ''
    const normalized = new Intl.NumberFormat('ru-RU', {
        maximumFractionDigits: 2,
    }).format(Math.abs(value))

    return `${sign}${normalized} ₽`
}
