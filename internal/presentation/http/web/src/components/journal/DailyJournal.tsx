import axios from 'axios'
import { useEffect, useRef, useState } from 'react'
import { Panel, PanelBody } from '../layout'
import { SavingStatus } from '../primitives'
import { useApiClient } from '../../hooks/useApiClient'
import './DailyJournal.sass'

type JournalSavingStatus = 'saved' | 'saving' | 'error'

type DailyJournalProps = {
    dateKey: string
    dateLabel: string
    initialNote: string | null
}

// DailyJournal renders the day journal panel and saves the note independently from habits.
export function DailyJournal({ dateKey, dateLabel, initialNote }: DailyJournalProps) {
    const apiClient = useApiClient()
    const [note, setNote] = useState('')
    const [status, setStatus] = useState<JournalSavingStatus>('saved')
    const saveTimer = useRef<number | undefined>(undefined)
    const saveGeneration = useRef(0)

    useEffect(() => {
        setNote(initialNote ?? '')
        setStatus('saved')
        clearPendingSave(saveTimer.current)
        saveGeneration.current += 1

        return () => {
            clearPendingSave(saveTimer.current)
        }
    }, [dateKey, initialNote])

    function handleChange(nextNote: string) {
        setNote(nextNote)
        setStatus('saving')

        saveGeneration.current += 1
        const generation = saveGeneration.current

        clearPendingSave(saveTimer.current)
        saveTimer.current = window.setTimeout(() => {
            void saveNote(nextNote, generation)
        }, 1000)
    }

    async function saveNote(nextNote: string, generation: number) {
        try {
            await apiClient.post(`journals/${dateKey}`, {
                note: nextNote,
            })

            if (generation === saveGeneration.current) {
                setStatus('saved')
            }
        } catch (error) {
            if (generation === saveGeneration.current) {
                setStatus('error')
            }

            if (!axios.isAxiosError(error)) {
                return
            }
        }
    }

    return (
        <Panel className="daily-journal">
            <PanelBody className="daily-journal__body">
                <div className="daily-journal__header">
                    <div className="daily-journal__heading">
                        <h2 className="daily-journal__title">Journal</h2>
                        <p className="daily-journal__subtitle">
                            Capture the notes that belong to {dateLabel}.
                        </p>
                    </div>
                    <SavingStatus status={status} className="daily-journal__status" />
                </div>

                <textarea
                    className="daily-journal__input"
                    maxLength={10000}
                    placeholder="What happened today?"
                    value={note}
                    onChange={(event) => handleChange(event.target.value)}
                />

                <div className="daily-journal__footer">
                    <span className="daily-journal__count">
                        {formatCount(note)} / 10 000
                    </span>
                    <div className="daily-journal__footer-spacer" aria-hidden="true" />
                </div>
            </PanelBody>
        </Panel>
    )
}

function clearPendingSave(timerId: number | undefined) {
    if (timerId === undefined) {
        return
    }

    window.clearTimeout(timerId)
}

function formatCount(note: string) {
    return Array.from(note)
        .length.toString()
        .replace(/\B(?=(\d{3})+(?!\d))/g, ' ')
}
