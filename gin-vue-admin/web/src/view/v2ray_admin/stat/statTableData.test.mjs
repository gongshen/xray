import assert from 'node:assert/strict'
import { normalizeStatTableResponse } from './statTableData.mjs'

const rows = [{ ID: 1 }, { ID: 2 }]
const success = normalizeStatTableResponse(
  {
    code: 0,
    data: {
      list: rows,
      total: 42,
      page: 3,
      pageSize: 50,
    },
  },
  { page: 1, pageSize: 10 }
)

assert.equal(success.ok, true)
assert.deepEqual(success.list, rows)
assert.notEqual(success.list, rows, 'list should be copied to trigger Vue updates')
assert.equal(success.total, 42)
assert.equal(success.page, 3)
assert.equal(success.pageSize, 50)
assert.equal(success.message, '')

assert.deepEqual(
  normalizeStatTableResponse({ code: 0, data: {} }, { page: 2, pageSize: 30 }),
  {
    ok: true,
    list: [],
    total: 0,
    page: 2,
    pageSize: 30,
    message: '',
  }
)

assert.deepEqual(
  normalizeStatTableResponse({ code: 7, msg: 'bad query' }, { page: 4, pageSize: 100 }),
  {
    ok: false,
    list: [],
    total: 0,
    page: 4,
    pageSize: 100,
    message: 'bad query',
  }
)

assert.deepEqual(
  normalizeStatTableResponse(null, { page: 5, pageSize: 10 }),
  {
    ok: false,
    list: [],
    total: 0,
    page: 5,
    pageSize: 10,
    message: '',
  }
)

console.log('statTableData tests passed')
