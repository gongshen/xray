export function createChartDataState({ includeRank = false } = {}) {
  const state = {
    data: [],
    data_axis: [],
    total: 0,
  }

  if (includeRank) {
    state.rank = []
    state.rank_axis = []
  }

  return state
}

function sumFiniteValues(values = []) {
  return values.reduce((total, value) => {
    const number = Number(value)
    return Number.isFinite(number) ? total + number : total
  }, 0)
}

export function applyTrendChartResponse(target, response) {
  const data = response?.data?.data

  if (response?.code !== 0 || !Array.isArray(data)) {
    target.data = []
    target.data_axis = []
    target.total = 0
    return false
  }

  target.data = data
  target.data_axis = Array.isArray(response.data.data_axis) ? response.data.data_axis : []
  target.total = sumFiniteValues(data)
  return true
}

export function applyRankChartResponse(target, response, { maxRankItems = 10 } = {}) {
  if (response?.code !== 0 || response?.data == null) {
    target.rank = []
    target.rank_axis = []
    return false
  }

  const rank = Array.isArray(response.data.rank) ? response.data.rank : []
  const rankAxis = Array.isArray(response.data.rank_axis) ? response.data.rank_axis : []

  target.rank = rank.length > maxRankItems ? rank.slice(-maxRankItems) : rank
  target.rank_axis = rankAxis.length > maxRankItems ? rankAxis.slice(-maxRankItems) : rankAxis
  return true
}
