import {reactive} from "vue";
import {
    getStatCharts,
} from '@/api/stat'

const chartData = reactive({data: [100,200,300,400,500,600,700,800,900,1000,1100,1200], data_axis: [202201,202202,202203,202204,202205,202206,202207,202208,202209,202210,202211,202212]})

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
    console.log("222:",chartData.data)
}