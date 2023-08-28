import {reactive} from "vue";
import {
    getStatCharts,
    getStatRank,
} from '@/api/stat'

const chartData = reactive({data: [], data_axis: [], total: 0, rank: [], rank_axis: []})

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    const ans2 = await getStatRank(searchInfo)
    if (ans2.code === 0 && ans2.data != null) {
        chartData.rank = ans2.data.rank
        chartData.rank_axis =ans2.data.rank_axis
    }else {
        chartData.rank = []
        chartData.rank_axis = []
    }

    const ans = await getStatCharts(searchInfo)
    if (ans.code === 0 && ans.data != null && ans.data.data != null) {
        chartData.data = ans.data.data
        chartData.data_axis =ans.data.data_axis
        chartData.total = ans.data.data.reduce((total,value) => {
            return total+value
        },0)
    }else{
        chartData.data = []
        chartData.data_axis = []
        chartData.total = 0
    }
}