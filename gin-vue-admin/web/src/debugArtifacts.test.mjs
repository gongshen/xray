import assert from 'node:assert/strict'
import fs from 'node:fs'

const files = [
  'src/utils/image.js',
  'src/utils/positionToCode.js',
  'src/view/example/breakpoint/breakpoint.vue',
  'src/view/example/upload/upload.vue',
  'src/view/person/person.vue',
]

for (const file of files) {
  const source = fs.readFileSync(file, 'utf8')
  assert.equal(source.includes('debugger'), false, `${file} contains debugger`)
  assert.equal(source.includes('console.log('), false, `${file} contains console.log`)
}

console.log('debugArtifacts tests passed')
