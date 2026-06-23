import {
    getStatCharts,
    getStatRank,
} from '@/api/stat';
import { reactive } from "vue";
import {
    applyRankChartResponse,
    applyTrendChartResponse,
    createChartDataState,
} from './statChartData.mjs'

const chartData = reactive(createChartDataState({includeRank: true}))

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    const ans2 = await getStatRank(searchInfo)
    applyRankChartResponse(chartData, ans2)

    const ans = await getStatCharts(searchInfo)
    applyTrendChartResponse(chartData, ans)
}
