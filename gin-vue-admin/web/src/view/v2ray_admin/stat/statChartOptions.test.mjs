import assert from 'node:assert/strict'
import {
  buildRankChartOptions,
  buildTrendChartOptions,
} from './statChartOptions.mjs'

const trendOptions = buildTrendChartOptions({
  data_axis: ['2026-06-22', '2026-06-23'],
  data: [1024, 1024 * 1024],
})

assert.equal(trendOptions.xAxis.type, 'category')
assert.deepEqual(trendOptions.xAxis.data, ['2026-06-22', '2026-06-23'])
assert.equal(trendOptions.yAxis.type, 'value')
assert.equal(trendOptions.series[0].type, 'line')
assert.deepEqual(trendOptions.series[0].data, [1024, 1024 * 1024])
assert.match(
  trendOptions.tooltip.formatter([
    { name: '2026-06-23', marker: '*', seriesName: 'Traffic', value: 1024 },
  ]),
  /1.0 KB/
)

const rankOptions = buildRankChartOptions({
  rank_axis: ['alice', 'bob'],
  rank: [100 * 1024 * 1024, 1024 * 1024 * 1024],
})

assert.equal(rankOptions.xAxis.type, 'value')
assert.equal(rankOptions.yAxis.type, 'category')
assert.deepEqual(rankOptions.yAxis.data, ['alice', 'bob'])
assert.equal(rankOptions.series[0].type, 'bar')
assert.deepEqual(rankOptions.series[0].data, [100 * 1024 * 1024, 1024 * 1024 * 1024])
assert.match(
  rankOptions.tooltip.formatter([
    { name: 'bob', marker: '*', seriesName: 'Traffic', value: 1024 * 1024 * 1024 },
  ]),
  /1.0 GB/
)

assert.deepEqual(buildTrendChartOptions({}).xAxis.data, [])
assert.deepEqual(buildTrendChartOptions({}).series[0].data, [])
assert.deepEqual(buildRankChartOptions({}).yAxis.data, [])
assert.deepEqual(buildRankChartOptions({}).series[0].data, [])

console.log('statChartOptions tests passed')
