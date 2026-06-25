import assert from 'node:assert/strict'
import { createRequire } from 'node:module'
import fs from 'node:fs'
import path from 'node:path'
import { loadConfigFromFile } from 'vite'

const source = fs.readFileSync(new URL('../vite.config.js', import.meta.url), 'utf8')
const require = createRequire(import.meta.url)

assert.doesNotMatch(
  source,
  /\['element-plus',\s*\[[^\]]*(?:node_modules\/element-plus|node_modules\/@element-plus)/,
  'Element Plus should not be split into its own manual chunk because it can create a vendor circular import during production build'
)

assert.match(
  source,
  /if \(normalizedId\.includes\('node_modules'\)\) {\s*return 'vendor'\s*}/,
  'Element Plus should fall back into the generic vendor chunk with its runtime dependencies'
)

process.env.NODE_ENV = 'production'
const result = await loadConfigFromFile(
  { command: 'build', mode: 'production' },
  'vite.config.js',
  process.cwd()
)
const manualChunks = result.config.build.rollupOptions.output.manualChunks
const chunkForPackage = (pkg) => manualChunks(require.resolve(pkg)) || 'app'

assert.equal(chunkForPackage('element-plus'), 'vendor')
assert.equal(chunkForPackage('@element-plus/icons-vue'), 'vendor')
assert.equal(chunkForPackage('@popperjs/core'), 'vendor')
assert.equal(chunkForPackage('dayjs'), 'vendor')
assert.equal(chunkForPackage('lodash'), 'vendor')
assert.equal(chunkForPackage('vue'), 'vue-vendor')
assert.equal(chunkForPackage('vue-router'), 'vue-vendor')
assert.equal(chunkForPackage('pinia'), 'vue-vendor')
assert.equal(chunkForPackage('echarts'), 'echarts')
assert.equal(chunkForPackage('zrender'), 'echarts')

const packages = [
  'vue',
  '@vue/runtime-core',
  '@vue/shared',
  'vue-router',
  'pinia',
  'element-plus',
  '@element-plus/icons-vue',
  '@popperjs/core',
  'dayjs',
  'lodash',
  'echarts',
  'zrender',
  'axios',
  'qs',
]
const knownPackages = new Set(packages)
const chunkEdges = new Set()
for (const pkg of packages) {
  const packageJsonPath = require.resolve(path.join(pkg, 'package.json'))
  const packageJson = require(packageJsonPath)
  const fromChunk = chunkForPackage(pkg)
  for (const dep of Object.keys({ ...(packageJson.dependencies || {}), ...(packageJson.peerDependencies || {}) })) {
    if (!knownPackages.has(dep)) {
      continue
    }
    const toChunk = chunkForPackage(dep)
    if (fromChunk !== toChunk) {
      chunkEdges.add(`${fromChunk}->${toChunk}`)
    }
  }
}

for (const edge of chunkEdges) {
  const [from, to] = edge.split('->')
  assert.ok(
    !chunkEdges.has(`${to}->${from}`),
    `manual chunks should not create a bidirectional package dependency between ${from} and ${to}`
  )
}

console.log('viteConfigManualChunks tests passed')
