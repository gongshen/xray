import assert from 'node:assert/strict'
import {
  calculatePercent,
  formatBytes,
  formatSize,
  formatUserOption,
  getTrafficAnalysisTotals,
  isServerOnline,
  parseClockMinute,
  rowTrafficAnalysisTargets,
  targetNames,
  todayCompact,
  validateTrafficAnalysisQuery,
} from './stateHelpers.mjs'

assert.equal(formatBytes(0), '0 B')
assert.equal(formatBytes(1024), '1.00 KB')
assert.equal(formatBytes(1536), '1.50 KB')
assert.equal(formatBytes(1024 ** 3), '1.00 GB')
assert.equal(formatBytes(Number.NaN), '0 B')

assert.equal(formatSize(0), '0 GB')
assert.equal(formatSize(512), '512 MB')
assert.equal(formatSize(1536), '1.5 GB')
assert.equal(formatSize(Number.NaN), '0 GB')

assert.equal(todayCompact(new Date(2026, 5, 23)), '20260623')
assert.equal(parseClockMinute('8:10'), 490)
assert.equal(parseClockMinute('09:00'), 540)
assert.equal(parseClockMinute('24:00'), null)
assert.equal(parseClockMinute('1:2'), null)

assert.equal(calculatePercent(25, 100), 25)
assert.equal(calculatePercent(2, 3), 67)
assert.equal(calculatePercent(5, 0), 0)
assert.equal(calculatePercent('bad', 10), 0)

assert.deepEqual(
  getTrafficAnalysisTotals([
    { down: 1, up: 2, total: 3 },
    { down: 4 },
  ]),
  { down: 5, up: 2, total: 3 }
)

assert.deepEqual(
  targetNames([
    { target: 'Example.com' },
    { target: 'Example.com' },
    { target: '' },
    { target: 'api.example.com' },
  ]),
  ['Example.com', 'api.example.com']
)
assert.deepEqual(
  rowTrafficAnalysisTargets({
    targets: [
      { target: 'Example.com' },
      { target: ' example.com ' },
      { target: 'API.EXAMPLE.COM' },
    ],
  }),
  ['example.com', 'api.example.com']
)

assert.equal(formatUserOption({ ID: 7, nickName: 'Alice', userName: 'alice' }), 'Alice (ID:7)')
assert.equal(formatUserOption({ ID: 8, userName: 'bob' }), 'bob (ID:8)')

assert.equal(isServerOnline(1000, 600, 1500), true)
assert.equal(isServerOnline(1000, 600, 1700), false)
assert.equal(isServerOnline(0, 600, 1700), false)

const validQuery = {
  currentServer: { ID: 1 },
  form: {
    user_tag: '7',
    date: '20260623',
    start: '8:10',
    end: '9:00',
  },
}
assert.equal(validateTrafficAnalysisQuery(validQuery), '')
assert.equal(validateTrafficAnalysisQuery({ ...validQuery, currentServer: null }), '服务器信息无效')
assert.equal(validateTrafficAnalysisQuery({ ...validQuery, form: { ...validQuery.form, user_tag: '' } }), '请选择用户')
assert.equal(validateTrafficAnalysisQuery({ ...validQuery, form: { ...validQuery.form, date: '2026-06-23' } }), '日期格式应为 20260617')
assert.equal(validateTrafficAnalysisQuery({ ...validQuery, form: { ...validQuery.form, start: '8:1' } }), '时间格式应为 8:10 或 09:00')
assert.equal(validateTrafficAnalysisQuery({ ...validQuery, form: { ...validQuery.form, start: '9:00', end: '8:10' } }), '结束时间不能早于开始时间')
assert.equal(validateTrafficAnalysisQuery({ ...validQuery, form: { ...validQuery.form, start: '8:00', end: '10:30' } }), '时间范围不能超过2小时')

console.log('stateHelpers tests passed')
