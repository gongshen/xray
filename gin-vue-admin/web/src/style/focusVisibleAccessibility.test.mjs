import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./main.scss', import.meta.url), 'utf8')

assert.match(
  source,
  /:focus-visible\s*\{/,
  'global styles should provide a visible keyboard focus fallback'
)

assert.match(
  source,
  /outline:\s*2px\s+solid\s+#409eff/,
  'global focus-visible style should restore a clear outline'
)

assert.match(
  source,
  /outline-offset:\s*2px/,
  'global focus-visible style should keep the outline from covering control borders'
)

console.log('focusVisibleAccessibility tests passed')
