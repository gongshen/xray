import {
    getStatCharts,
} from '@/api/v2ray_stat';
import { reactive } from "vue";
import {
    applyTrendChartResponse,
    createChartDataState,
} from '../../v2ray_admin/stat/statChartData.mjs'

const chartData = reactive(createChartDataState())

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    // 获取图表数据
    try {
        const ans = await getStatCharts({ ...searchInfo })
        if (!applyTrendChartResponse(chartData, ans)) {
            console.warn('图表数据格式不正确或为空:', ans)
        }
    } catch (error) {
        console.error('获取图表数据失败:', error)
        applyTrendChartResponse(chartData, null)
    }
}
