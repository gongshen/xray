import assert from 'node:assert/strict'
import fs from 'node:fs'

function escapeRegExp(value) {
  return String(value).replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
}

function findRawHexColors(source) {
  return source.match(/#[0-9a-fA-F]{3,8}\b/g) || []
}

let tokenModule
try {
  tokenModule = await import('./designTokens.mjs')
} catch (error) {
  assert.fail(`expected shared designTokens.mjs module: ${error.message}`)
}

const {
  chartPalette,
  progressThresholdColors,
  shadows,
  uiColors,
} = tokenModule

assert.equal(uiColors.brandPrimary, '#1e40af')
assert.equal(uiColors.brandSecondary, '#3b82f6')
assert.equal(uiColors.accent, '#f59e0b')
assert.equal(uiColors.pageBg, '#f8fafc')
assert.equal(uiColors.textMuted, '#475569')
assert.equal(progressThresholdColors[1].color, uiColors.trafficUp)
assert.equal(chartPalette.trendLineStart, uiColors.brandSecondary)

const basics = fs.readFileSync('src/style/basics.scss', 'utf8')
const cssTokens = {
  'gva-color-brand-primary': uiColors.brandPrimary,
  'gva-color-brand-secondary': uiColors.brandSecondary,
  'gva-color-accent': uiColors.accent,
  'gva-color-page-bg': uiColors.pageBg,
  'gva-color-panel-bg': uiColors.panelBg,
  'gva-color-panel-muted-bg': uiColors.panelMutedBg,
  'gva-color-text-strong': uiColors.textStrong,
  'gva-color-text-regular': uiColors.textRegular,
  'gva-color-text-muted': uiColors.textMuted,
  'gva-color-border-subtle': uiColors.borderSubtle,
  'gva-color-border-muted': uiColors.borderMuted,
  'gva-color-traffic-down': uiColors.trafficDown,
  'gva-color-traffic-up': uiColors.trafficUp,
  'gva-shadow-panel': shadows.panel,
  'gva-shadow-panel-hover': shadows.panelHover,
}

for (const [token, value] of Object.entries(cssTokens)) {
  assert.match(
    basics,
    new RegExp(`--${escapeRegExp(token)}:\\s*${escapeRegExp(value)};`, 'i'),
    `missing CSS token --${token}`
  )
}

const vueTargets = [
  'src/view/system/state.vue',
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]

for (const file of vueTargets) {
  const source = fs.readFileSync(file, 'utf8')
  assert.deepEqual(findRawHexColors(source), [], `${file} still has raw hex colors`)
  assert.match(source, /var\(--gva-color-/, `${file} should use shared CSS color tokens`)
}

for (const file of [
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]) {
  const source = fs.readFileSync(file, 'utf8')
  assert.match(
    source,
    /@media\s*\(prefers-reduced-motion:\s*reduce\)/,
    `${file} should disable chart entrance animation for reduced motion`
  )
}

const stateSource = fs.readFileSync('src/view/system/state.vue', 'utf8')
assert.match(stateSource, /progressThresholdColors/, 'state.vue should use shared progress colors')

const chartSource = fs.readFileSync('src/view/v2ray_admin/stat/statChartOptions.mjs', 'utf8')
assert.match(chartSource, /chartPalette/, 'statChartOptions should use shared chart palette')
assert.deepEqual(
  findRawHexColors(chartSource),
  [],
  'statChartOptions.mjs should not contain raw hex colors'
)

console.log('designTokens tests passed')
