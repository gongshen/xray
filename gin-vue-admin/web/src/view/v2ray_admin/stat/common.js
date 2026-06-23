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

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    await loadRankChartData(chartData, getStatRank, searchInfo)
    await loadTrendChartData(chartData, getStatCharts, searchInfo)
}
