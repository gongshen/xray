import assert from 'node:assert/strict'
import {
  createDefaultTrafficSearchRange,
  formatFlow,
  getDateRangeText,
  getTrafficTagType,
  normalizeDateOnlyToUtcIso,
  normalizeTrafficSearchRange,
  parseTrafficBytes,
} from './statTraffic.mjs'

assert.equal(formatFlow(null), '0 B')
assert.equal(formatFlow(0), '0 B')
assert.equal(formatFlow(512), '512.0 B')
assert.equal(formatFlow(1024), '1.0 KB')
assert.equal(formatFlow(1024 * 1024), '1.0 MB')
assert.equal(formatFlow(1024 * 1024 * 1024), '1.0 GB')

assert.equal(parseTrafficBytes('726.27MB'), 726.27 * 1024 * 1024)
assert.equal(parseTrafficBytes('1.5 GB'), 1.5 * 1024 * 1024 * 1024)
assert.equal(parseTrafficBytes('bad input'), null)
assert.equal(parseTrafficBytes(''), null)

assert.equal(getTrafficTagType('9 MB'), 'info')
assert.equal(getTrafficTagType('10 MB'), 'success')
assert.equal(getTrafficTagType('100 MB'), 'warning')
assert.equal(getTrafficTagType('1 GB'), 'danger')
assert.equal(getTrafficTagType('bad input'), 'info')

assert.equal(
  normalizeDateOnlyToUtcIso(new Date(2026, 5, 23, 15, 42)),
  '2026-06-23T00:00:00.000Z'
)
assert.equal(normalizeDateOnlyToUtcIso(''), '')
assert.equal(normalizeDateOnlyToUtcIso(null), null)

assert.deepEqual(
  createDefaultTrafficSearchRange({ now: new Date('2026-06-23T12:00:00.000Z') }),
  {
    startCreatedAt: '2026-05-24T12:00:00.000Z',
    endCreatedAt: '2026-06-23T12:00:00.000Z',
  }
)
assert.deepEqual(
  createDefaultTrafficSearchRange({ now: new Date('2026-06-23T12:00:00.000Z'), days: 7 }),
  {
    startCreatedAt: '2026-06-16T12:00:00.000Z',
    endCreatedAt: '2026-06-23T12:00:00.000Z',
  }
)

const searchRange = {
  startCreatedAt: new Date(2026, 5, 23, 15, 42),
  endCreatedAt: new Date(2026, 5, 24, 8, 10),
  tag: 'user-1',
}
const normalizedRange = normalizeTrafficSearchRange(searchRange)
assert.notEqual(normalizedRange, searchRange)
assert.deepEqual(normalizedRange, {
  startCreatedAt: '2026-06-23T00:00:00.000Z',
  endCreatedAt: '2026-06-24T00:00:00.000Z',
  tag: 'user-1',
})

assert.equal(getDateRangeText({}), '最近 30 天')
assert.equal(
  getDateRangeText({
    startCreatedAt: '2026-06-23T00:00:00.000Z',
    endCreatedAt: '2026-06-23T00:00:00.000Z',
  }),
  '2026/6/23'
)
assert.equal(
  getDateRangeText({
    startCreatedAt: '2026-06-01T00:00:00.000Z',
    endCreatedAt: '2026-06-23T00:00:00.000Z',
  }),
  '2026/6/1 - 2026/6/23'
)

console.log('statTraffic tests passed')
