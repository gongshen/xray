import {
    getStatCharts,
    getStatRank,
} from '@/api/stat';
import { reactive } from "vue";
import {
    createChartDataState,
    loadRankChartData,
    loadTrendChartData,
} from './statChartData.mjs'

const chartData = reactive(createChartDataState({includeRank: true}))
let chartDataRequestId = 0

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    const requestId = ++chartDataRequestId
    const nextChartData = createChartDataState({includeRank: true})
    const [rankLoaded, trendLoaded] = await Promise.all([
        loadRankChartData(nextChartData, getStatRank, searchInfo),
        loadTrendChartData(nextChartData, getStatCharts, searchInfo),
    ])
    if (requestId !== chartDataRequestId) {
        return false
    }

    Object.assign(chartData, nextChartData)
    return rankLoaded && trendLoaded
}
