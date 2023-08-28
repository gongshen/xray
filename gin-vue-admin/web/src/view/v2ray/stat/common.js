import {reactive} from "vue";
import {
    getStatCharts,
} from '@/api/v2ray_stat'

const chartData = reactive({data: [], data_axis: [], total: 0})

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    const ans = await getStatCharts({ ...searchInfo })
    if (ans.code === 0 && ans.data != null && ans.data.data != null) {
        chartData.data = ans.data.data
        chartData.data_axis =ans.data.data_axis
        chartData.total = chartData.data.reduce((total,value) => {
            return total+value
        },0)
    }else{
        chartData.data = []
        chartData.data_axis = []
        chartData.total = 0
    }
}