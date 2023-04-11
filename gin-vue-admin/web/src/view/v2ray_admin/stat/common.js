import {reactive} from "vue";
import {
    getStatCharts,
} from '@/api/stat'

const chartData = reactive({data: [], data_axis: []})

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    const ans = await getStatCharts({ ...searchInfo })
    if (ans.code === 0 && ans.data != null) {
        chartData.data = ans.data.data
        chartData.data_axis =ans.data.data_axis
    }else{
        chartData.data = []
        chartData.data_axis = []
    }
    console.log("222:",chartData.data,chartData.data_axis)
}