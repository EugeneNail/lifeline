import { useEffect, useMemo, useRef, useState, type MouseEvent } from 'react'
import './HabitChart.sass'

export type HabitChartHabitType = 'measurable' | 'time' | 'completable'

export type HabitChartSeries = {
    nodes: Array<{ date: string; value: number }>
    minValue: number
    maxValue: number
}

type HabitChartProps = {
    habitType: HabitChartHabitType
    series: HabitChartSeries
}

type CanvasSize = {
    width: number
    height: number
    pixelRatio: number
}

type ChartLayout = {
    left: number
    right: number
    plotTop: number
    plotBottom: number
    labelFontSize: number
    xLabelFontSize: number
}

const plotTopPadding = 16
const xLabelHeight = 44
const yLabelGap = 14
const xLabelGap = 12
const plotRightPadding = 12
const chartValueFormatter = new Intl.NumberFormat('en-GB', {
    maximumFractionDigits: 1,
})
const shortMonthFormatter = new Intl.DateTimeFormat('en-GB', {
    month: 'short',
})

// HabitChart renders a responsive canvas bar chart with a date-linked hover guide.
export function HabitChart({ habitType, series }: HabitChartProps) {
    const canvasRef = useRef<HTMLCanvasElement>(null)
    const containerRef = useRef<HTMLDivElement>(null)
    const [canvasSize, setCanvasSize] = useState<CanvasSize>({
        width: 0,
        height: 0,
        pixelRatio: 1,
    })
    const [isPortrait, setIsPortrait] = useState(
        () => typeof window !== 'undefined' && window.matchMedia('(orientation: portrait)').matches,
    )
    const [hoveredNodeIndex, setHoveredNodeIndex] = useState<number | null>(null)
    const chartAspectRatio = series.nodes.length > 90 ? (isPortrait ? 3 : 4) : 2
    const values = useMemo(
        () => series.nodes.map((node) => normalizeValue(node.value, habitType)),
        [habitType, series.nodes],
    )
    const bounds = { min: 0, max: series.maxValue }

    useEffect(() => {
        const mediaQuery = window.matchMedia('(orientation: portrait)')
        const updateOrientation = () => setIsPortrait(mediaQuery.matches)

        updateOrientation()
        mediaQuery.addEventListener('change', updateOrientation)

        return () => mediaQuery.removeEventListener('change', updateOrientation)
    }, [])

    useEffect(() => {
        const container = containerRef.current
        if (!container) {
            return
        }
        const element = container

        function updateCanvasSize() {
            const width = element.clientWidth
            setCanvasSize({
                width,
                height: width / chartAspectRatio,
                pixelRatio: window.devicePixelRatio || 1,
            })
        }

        updateCanvasSize()
        const resizeObserver = new ResizeObserver(updateCanvasSize)
        resizeObserver.observe(element)

        return () => resizeObserver.disconnect()
    }, [chartAspectRatio])

    useEffect(() => {
        const canvas = canvasRef.current
        if (!canvas || canvasSize.width <= 0 || canvasSize.height <= 0) {
            return
        }

        const context = canvas.getContext('2d')
        if (!context) {
            return
        }

        const layout = getChartLayout(context, canvasSize.width, canvasSize.height, bounds, habitType)
        const scale = canvasSize.pixelRatio
        canvas.width = Math.round(canvasSize.width * scale)
        canvas.height = Math.round(canvasSize.height * scale)
        context.setTransform(scale, 0, 0, scale, 0, 0)
        context.clearRect(0, 0, canvasSize.width, canvasSize.height)

        drawGrid(context, canvasSize.width, layout)
        drawYAxisLabels(context, bounds, habitType, layout)
        drawBars(context, values, bounds, canvasSize.width, layout)
        drawXAxisLabels(context, series.nodes, canvasSize.width, layout)

        if (hoveredNodeIndex !== null && values.length > 0) {
            drawHoverLine(
                context,
                getNodeX(hoveredNodeIndex, values.length, layout.left, layout.right, canvasSize.width),
                layout,
            )
        }
    }, [bounds, canvasSize, habitType, hoveredNodeIndex, series.nodes, values])

    function handleMouseMove(event: MouseEvent<HTMLCanvasElement>) {
        if (values.length === 0 || canvasSize.width <= 0) {
            return
        }

        const bounds = event.currentTarget.getBoundingClientRect()
        const x = event.clientX - bounds.left
        const layout = getChartLayout(
            null,
            canvasSize.width,
            canvasSize.height,
            { min: 0, max: series.maxValue },
            habitType,
        )
        const plotWidth = canvasSize.width - layout.left - layout.right
        const slotWidth = plotWidth / values.length
        const index = Math.min(
            Math.max(Math.round((x - layout.left) / slotWidth - .5), 0),
            values.length - 1,
        )
        setHoveredNodeIndex(index)
    }

    return (
        <div className="habit-chart" ref={containerRef}>
            <canvas
                    aria-label={`Habit bar chart for ${habitType}`}
                className="habit-chart__canvas"
                height={Math.round(canvasSize.height * canvasSize.pixelRatio)}
                ref={canvasRef}
                role="img"
                style={{ aspectRatio: `${chartAspectRatio} / 1` }}
                width={Math.round(canvasSize.width * canvasSize.pixelRatio)}
                onMouseLeave={() => setHoveredNodeIndex(null)}
                onMouseMove={handleMouseMove}
            />
        </div>
    )
}

function normalizeValue(value: number, habitType: HabitChartHabitType) {
    if (habitType === 'completable') {
        return value > 0 ? 1 : 0
    }

    if (habitType === 'time') {
        return Math.min(Math.max(value, 0), 1439)
    }

    return Math.max(value, 0)
}

function getChartLayout(
    context: CanvasRenderingContext2D | null,
    width: number,
    height: number,
    bounds: { min: number; max: number },
    habitType: HabitChartHabitType,
): ChartLayout {
    const labelFontSize = Math.max(10, Math.min(14, width / 55))
    const xLabelFontSize = Math.max(9, Math.min(12, width / 70))
    const labels = getYAxisLabels(bounds, habitType)

    if (context) {
        context.font = `${labelFontSize}px sans-serif`
    }

    const widestLabel = context
        ? Math.max(...labels.map((label) => context.measureText(label).width))
        : labels.reduce((widest, label) => Math.max(widest, label.length * labelFontSize * .6), 0)

    return {
        left: widestLabel + yLabelGap + 6,
        right: plotRightPadding,
        plotTop: plotTopPadding,
        plotBottom: height - xLabelHeight,
        labelFontSize,
        xLabelFontSize,
    }
}

function getYAxisLabels(bounds: { min: number; max: number }, habitType: HabitChartHabitType) {
    return [0, .25, .5, .75, 1].map((ratio) => {
        const value = bounds.max - (bounds.max - bounds.min) * ratio
        return formatValue(value, habitType)
    })
}

function drawGrid(
    context: CanvasRenderingContext2D,
    width: number,
    layout: ChartLayout,
) {
    const rowHeight = (layout.plotBottom - layout.plotTop) / 4
    context.lineWidth = 1

    for (let row = 0; row < 5; row += 1) {
        const y = layout.plotTop + rowHeight * row
        context.beginPath()
        context.moveTo(layout.left, y)
        context.lineTo(width - layout.right, y)
        context.strokeStyle = 'rgba(66, 104, 88, .14)'
        context.stroke()
    }
}

function drawYAxisLabels(
    context: CanvasRenderingContext2D,
    bounds: { min: number; max: number },
    habitType: HabitChartHabitType,
    layout: ChartLayout,
) {
    const labels = getYAxisLabels(bounds, habitType)
    const rowHeight = (layout.plotBottom - layout.plotTop) / 4

    context.fillStyle = 'rgb(104, 115, 108)'
    context.font = `${layout.labelFontSize}px sans-serif`
    context.textAlign = 'right'
    context.textBaseline = 'middle'

    labels.forEach((label, index) => {
        context.fillText(label, layout.left - yLabelGap, layout.plotTop + rowHeight * index)
    })
}

function drawBars(
    context: CanvasRenderingContext2D,
    values: number[],
    bounds: { min: number; max: number },
    width: number,
    layout: ChartLayout,
) {
    if (values.length === 0) {
        return
    }

    const range = bounds.max - bounds.min || 1
    const plotHeight = layout.plotBottom - layout.plotTop
    const plotWidth = width - layout.left - layout.right
    const slotWidth = plotWidth / values.length
    const barWidth = Math.max(2, slotWidth * .72)

    context.fillStyle = 'rgb(66, 104, 88)'
    values.forEach((value, index) => {
        const normalizedValue = Math.max((value - bounds.min) / range, 0)
        const barHeight = normalizedValue * plotHeight
        if (barHeight <= 0) {
            return
        }

        const x = layout.left + index * slotWidth + (slotWidth - barWidth) / 2
        const y = layout.plotBottom - barHeight
        drawRoundedBar(context, x, y, barWidth, barHeight)
    })
}

function drawRoundedBar(
    context: CanvasRenderingContext2D,
    x: number,
    y: number,
    width: number,
    height: number,
) {
    const radius = Math.min(2, width / 2, height)

    context.beginPath()
    context.moveTo(x, y + radius)
    context.quadraticCurveTo(x, y, x + radius, y)
    context.lineTo(x + width - radius, y)
    context.quadraticCurveTo(x + width, y, x + width, y + radius)
    context.lineTo(x + width, y + height)
    context.lineTo(x, y + height)
    context.closePath()
    context.fill()
}

function drawXAxisLabels(
    context: CanvasRenderingContext2D,
    nodes: HabitChartSeries['nodes'],
    width: number,
    layout: ChartLayout,
) {
    if (nodes.length === 0) {
        return
    }

    context.fillStyle = 'rgb(104, 115, 108)'
    context.font = `${layout.xLabelFontSize}px sans-serif`
    context.textAlign = 'center'
    context.textBaseline = 'top'

    for (let labelIndex = 0; labelIndex < 7; labelIndex += 1) {
        const nodeIndex = nodes.length === 1
            ? 0
            : Math.round((nodes.length - 1) * labelIndex / 6)
        const x = getNodeX(nodeIndex, nodes.length, layout.left, layout.right, width)
        const date = parseDate(nodes[nodeIndex]?.date)
        if (!date) {
            continue
        }

        context.fillText(String(date.getDate()), x, layout.plotBottom + xLabelGap)
        context.fillText(shortMonthFormatter.format(date).toLowerCase(), x, layout.plotBottom + xLabelGap + layout.xLabelFontSize + 2)
    }
}

function drawHoverLine(
    context: CanvasRenderingContext2D,
    x: number,
    layout: ChartLayout,
) {
    context.beginPath()
    context.moveTo(x, layout.plotTop)
    context.lineTo(x, layout.plotBottom)
    context.lineWidth = 1
    context.strokeStyle = 'rgba(66, 104, 88, .6)'
    context.stroke()
}

function getNodeX(index: number, nodeCount: number, left: number, right: number, width: number) {
    const plotWidth = width - left - right
    return left + plotWidth * (index + .5) / nodeCount
}

function parseDate(dateValue?: string) {
    if (!dateValue) {
        return null
    }

    const date = new Date(`${dateValue}T00:00:00`)
    return Number.isNaN(date.getTime()) ? null : date
}

function formatValue(value: number, habitType: HabitChartHabitType) {
    if (habitType === 'time') {
        const minutes = Math.min(Math.max(Math.round(value), 0), 1439)
        return `${String(Math.floor(minutes / 60)).padStart(2, '0')}:${String(minutes % 60).padStart(2, '0')}`
    }

    return chartValueFormatter.format(value)
}
