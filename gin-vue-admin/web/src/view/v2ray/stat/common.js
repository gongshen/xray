import {
    getStatCharts,
} from '@/api/v2ray_stat';
import { reactive } from "vue";
import {
    createChartDataState,
    loadTrendChartData,
} from '../../v2ray_admin/stat/statChartData.mjs'

const chartData = reactive(createChartDataState())

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    // 获取图表数据
    const loaded = await loadTrendChartData(chartData, getStatCharts, { ...searchInfo })
    if (!loaded) {
        console.warn('图表数据格式不正确或为空')
    }
}
