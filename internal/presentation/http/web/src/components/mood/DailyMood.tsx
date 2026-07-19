import axios from 'axios'
import { useEffect, useRef, useState } from 'react'
import { Panel, PanelBody } from '../layout'
import { SavingStatus } from '../primitives'
import { useApiClient } from '../../hooks/useApiClient'
import './DailyMood.sass'

type MoodValue = 1 | 2 | 3 | 4 | 5

type MoodOption = {
    value: MoodValue
    emoji: string
    label: string
}

type DailyMoodProps = {
    dateKey: string
}

type MoodSavingStatus = 'saved' | 'saving' | 'error'

const moodOptions: MoodOption[] = [
    { value: 1, emoji: '😞', label: 'Awful' },
    { value: 2, emoji: '😕', label: 'Bad' },
    { value: 3, emoji: '😐', label: 'Okay' },
    { value: 4, emoji: '🙂', label: 'Good' },
    { value: 5, emoji: '😄', label: 'Great' },
]

// DailyMood renders the day mood panel and saves the selected mood independently from other panels.
export function DailyMood({ dateKey }: DailyMoodProps) {
    const apiClient = useApiClient()
    const [selectedMood, setSelectedMood] = useState<MoodValue | null>(null)
    const [status, setStatus] = useState<MoodSavingStatus>('saved')
    const saveGeneration = useRef(0)

    useEffect(() => {
        setSelectedMood(null)
        setStatus('saved')
        saveGeneration.current += 1
    }, [dateKey])

    function handleSelectMood(mood: MoodValue) {
        setSelectedMood(mood)
        setStatus('saving')

        saveGeneration.current += 1
        const generation = saveGeneration.current

        void saveMood(mood, generation)
    }

    async function saveMood(mood: MoodValue, generation: number) {
        try {
            await apiClient.post(`moods/${dateKey}`, {
                mood,
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
        <Panel className="daily-mood">
            <PanelBody className="daily-mood__body">
                <div className="daily-mood__header">
                    <div className="daily-mood__heading">
                        <h2 className="daily-mood__title">How did the day go?</h2>
                        <p className="daily-mood__subtitle">
                            Choose one mood. It saves immediately after selection.
                        </p>
                    </div>
                    <SavingStatus status={status} className="daily-mood__status" />
                </div>

                <div className="daily-mood__choices" role="radiogroup" aria-label="Daily mood">
                    {moodOptions.map((option) => (
                        <button
                            aria-checked={selectedMood === option.value}
                            aria-label={option.label}
                            className={joinClassNames(
                                'daily-mood__choice',
                                selectedMood === option.value
                                    ? 'daily-mood__choice--selected'
                                    : undefined,
                            )}
                            key={option.value}
                            onClick={() => handleSelectMood(option.value)}
                            role="radio"
                            type="button"
                        >
                            <span className="daily-mood__emoji" aria-hidden="true">
                                {option.emoji}
                            </span>
                            <span className="daily-mood__label">{option.label}</span>
                        </button>
                    ))}
                </div>
            </PanelBody>
        </Panel>
    )
}

function joinClassNames(...classNames: Array<string | undefined>) {
    return classNames.filter(Boolean).join(' ')
}
