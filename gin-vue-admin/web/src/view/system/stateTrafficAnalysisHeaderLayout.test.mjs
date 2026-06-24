import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync('src/view/system/state.vue', 'utf8')

function ruleFor(selector) {
  const escaped = selector.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  const match = source.match(new RegExp(`${escaped}\\s*\\{[^}]*\\}`))
  assert.ok(match, `${selector} style rule should exist`)
  return match[0]
}

const headerRule = ruleFor('.analysis-header')
const titleRule = ruleFor('.analysis-title')
const subtitleRule = ruleFor('.analysis-subtitle')

assert.match(
  headerRule,
  /flex-wrap:\s*wrap;/,
  'traffic analysis dialog header should allow wrapping on narrow screens'
)
assert.match(
  titleRule,
  /word-break:\s*break-word;/,
  'traffic analysis title should wrap long server remarks'
)
assert.match(
  subtitleRule,
  /word-break:\s*break-word;/,
  'traffic analysis subtitle should wrap long server addresses and metadata'
)

console.log('stateTrafficAnalysisHeaderLayout tests passed')
