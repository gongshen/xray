import {
    getStatCharts,
    getStatRank,
} from '@/api/stat';
import { reactive } from "vue";

const chartData = reactive({data: [], data_axis: [], total: 0, rank: [], rank_axis: []})

export const useChartData = () => {
    return chartData
}

export const setChartData = async (searchInfo) => {
    const ans2 = await getStatRank(searchInfo)
    if (ans2.code === 0 && ans2.data != null) {
        // 限制排行榜只显示前10名（流量最大的）
        const maxRankItems = 10
        // 后端数据是从小到大排序，所以取最后10条数据（流量最大的）
        const rankData = ans2.data.rank || []
        const rankAxisData = ans2.data.rank_axis || []
        
        if (rankData.length > maxRankItems) {
            chartData.rank = rankData.slice(-maxRankItems)
            chartData.rank_axis = rankAxisData.slice(-maxRankItems)
        } else {
            chartData.rank = rankData
            chartData.rank_axis = rankAxisData
        }
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