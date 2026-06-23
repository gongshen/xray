import assert from 'node:assert/strict'
import fs from 'node:fs'

const files = [
  'src/view/v2ray_admin/stat/stat.vue',
  'src/view/v2ray/stat/stat.vue',
]

for (const file of files) {
  const source = fs.readFileSync(file, 'utf8')
  assert.doesNotMatch(source, /\sicon="(?:search|refresh)"/i, `${file} should not use string action icons`)
  assert.match(source, /:icon="Search"/, `${file} should bind the Search icon component`)
  assert.match(source, /:icon="Refresh"/, `${file} should bind the Refresh icon component`)
  assert.match(
    source,
    /import\s*\{[^}]*Refresh[^}]*Search[^}]*\}\s*from '@element-plus\/icons-vue'|import\s*\{[^}]*Search[^}]*Refresh[^}]*\}\s*from '@element-plus\/icons-vue'/s,
    `${file} should import Search and Refresh icons`
  )
}

console.log('statActionIcons tests passed')
