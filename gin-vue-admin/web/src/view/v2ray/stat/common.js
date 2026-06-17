import {
    getStatCharts,
} from '@/api/v2ray_stat';
import { reactive } from "vue";

const chartData = reactive({data: [], data_axis: [], total: 0})

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    // 获取图表数据
    try {
        const ans = await getStatCharts({ ...searchInfo })
        if (ans.code === 0 && ans.data != null && ans.data.data != null) {
            chartData.data = ans.data.data
            chartData.data_axis = ans.data.data_axis
            chartData.total = chartData.data.reduce((total, value) => {
                return total + value
            }, 0)
        } else {
            console.warn('图表数据格式不正确或为空:', ans)
            chartData.data = []
            chartData.data_axis = []
            chartData.total = 0
        }
    } catch (error) {
        console.error('获取图表数据失败:', error)
        chartData.data = []
        chartData.data_axis = []
        chartData.total = 0
    }
}
