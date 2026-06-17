import assert from 'node:assert/strict'
import { buildServerPortChangeReminder } from './serverPortReminder.mjs'

assert.equal(buildServerPortChangeReminder(null, { port: 80, stat_port: 56611 }), '')

assert.equal(
  buildServerPortChangeReminder(
    { port: 80, stat_port: 56611 },
    { port: '80', stat_port: '56611' }
  ),
  ''
)

assert.match(
  buildServerPortChangeReminder(
    { port: 80, stat_port: 56611 },
    { port: 443, stat_port: 56611 }
  ),
  /Xray 程序配置中的端口/
)

assert.match(
  buildServerPortChangeReminder(
    { port: 80, stat_port: 56611 },
    { port: 80, stat_port: 56612 }
  ),
  /stat 的端口/
)

const bothChanged = buildServerPortChangeReminder(
  { port: 80, stat_port: 56611 },
  { port: 443, stat_port: 56612 }
)
assert.match(bothChanged, /Xray 程序配置中的端口/)
assert.match(bothChanged, /stat 的端口/)

console.log('serverPortReminder tests passed')
