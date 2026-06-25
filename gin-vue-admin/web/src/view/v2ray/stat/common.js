import {
    getStatCharts,
} from '@/api/v2ray_stat';
import { devWarn } from '@/utils/devLogger'
import { reactive } from "vue";
import {
    createChartDataState,
    loadTrendChartData,
} from '../../v2ray_admin/stat/statChartData.mjs'

const chartData = reactive(createChartDataState())
let chartDataRequestId = 0

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    const requestId = ++chartDataRequestId
    const nextChartData = createChartDataState()
    const loaded = await loadTrendChartData(nextChartData, getStatCharts, { ...searchInfo })
    if (requestId !== chartDataRequestId) {
        return false
    }

    Object.assign(chartData, nextChartData)
    if (!loaded) {
        devWarn('图表数据格式不正确或为空')
    }
    return loaded
}
