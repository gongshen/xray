import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./warningBar.vue', import.meta.url), 'utf8')

const linkedWarning = source.match(/<button[^>]*v-if="href"[^>]*class="warning-bar can-click"[^>]*>/)
const staticWarning = source.match(/<div[^>]*v-else[^>]*class="warning-bar"[^>]*>/)

assert.ok(linkedWarning, 'linked warning bar should use a native button')
assert.match(linkedWarning[0], /type="button"/, 'linked warning bar button should not submit forms')
assert.match(linkedWarning[0], /:aria-label="warningLinkLabel"/, 'linked warning bar button should have an accessible name')
assert.match(linkedWarning[0], /@click="open"/, 'linked warning bar should keep the open handler')

assert.ok(staticWarning, 'static warning bar should render as non-interactive content')
assert.doesNotMatch(staticWarning[0], /@click=/, 'static warning bar should not bind click behavior')
assert.doesNotMatch(source, /<div[^>]*class="warning-bar"[^>]*@click="open"/, 'warning bar should not be a clickable div')
assert.match(source, /const warningLinkLabel = computed\(\(\) => .*\$\{prop\.title\}/, 'warning bar accessible name should stay derived from the title prop')
assert.match(source, /import \{ openExternalUrl \} from '@\/utils\/openExternalUrl'/, 'warning bar should use the shared external URL helper')
assert.match(source, /openExternalUrl\(prop\.href\)/, 'warning bar should open links through the shared external URL helper')
assert.doesNotMatch(source, /window\.open/, 'warning bar should not call window.open directly')

console.log('warningBarAccessibility tests passed')
