import assert from 'node:assert/strict'
import {
  applyRankChartResponse,
  applyTrendChartResponse,
  createChartDataState,
  loadRankChartData,
  loadTrendChartData,
} from './statChartData.mjs'

assert.deepEqual(createChartDataState(), {
  data: [],
  data_axis: [],
  total: 0,
})
assert.deepEqual(createChartDataState({ includeRank: true }), {
  data: [],
  data_axis: [],
  total: 0,
  rank: [],
  rank_axis: [],
})

{
  const target = createChartDataState()
  const applied = applyTrendChartResponse(target, {
    code: 0,
    data: {
      data: [1024, 2048, 'bad'],
      data_axis: ['2026-06-22', '2026-06-23', '2026-06-24'],
    },
  })

  assert.equal(applied, true)
  assert.deepEqual(target.data, [1024, 2048, 'bad'])
  assert.deepEqual(target.data_axis, ['2026-06-22', '2026-06-23', '2026-06-24'])
  assert.equal(target.total, 3072)
}

{
  const target = { data: [1], data_axis: ['x'], total: 1 }
  const applied = applyTrendChartResponse(target, { code: 1, data: null })

  assert.equal(applied, false)
  assert.deepEqual(target, { data: [], data_axis: [], total: 0 })
}

{
  const target = createChartDataState({ includeRank: true })
  const rank = Array.from({ length: 12 }, (_, index) => index + 1)
  const rank_axis = rank.map(value => `user-${value}`)
  const applied = applyRankChartResponse(target, {
    code: 0,
    data: {
      rank,
      rank_axis,
    },
  })

  assert.equal(applied, true)
  assert.deepEqual(target.rank, [3, 4, 5, 6, 7, 8, 9, 10, 11, 12])
  assert.deepEqual(target.rank_axis, [
    'user-3',
    'user-4',
    'user-5',
    'user-6',
    'user-7',
    'user-8',
    'user-9',
    'user-10',
    'user-11',
    'user-12',
  ])
}

{
  const target = { rank: [1], rank_axis: ['x'] }
  const applied = applyRankChartResponse(target, { code: 0, data: null })

  assert.equal(applied, false)
  assert.deepEqual(target, { rank: [], rank_axis: [] })
}

{
  const target = createChartDataState()
  const calls = []
  const applied = await loadTrendChartData(target, async (searchInfo) => {
    calls.push(searchInfo)
    return {
      code: 0,
      data: {
        data: [1, 2],
        data_axis: ['a', 'b'],
      },
    }
  }, { server_ip: '127.0.0.1' })

  assert.equal(applied, true)
  assert.deepEqual(calls, [{ server_ip: '127.0.0.1' }])
  assert.deepEqual(target, { data: [1, 2], data_axis: ['a', 'b'], total: 3 })
}

{
  const target = { data: [1], data_axis: ['x'], total: 1 }
  const applied = await loadTrendChartData(target, async () => {
    throw new Error('network')
  }, {})

  assert.equal(applied, false)
  assert.deepEqual(target, { data: [], data_axis: [], total: 0 })
}

{
  const target = createChartDataState({ includeRank: true })
  const applied = await loadRankChartData(target, async () => {
    throw new Error('network')
  }, {})

  assert.equal(applied, false)
  assert.deepEqual(target.rank, [])
  assert.deepEqual(target.rank_axis, [])
}

console.log('statChartData tests passed')
