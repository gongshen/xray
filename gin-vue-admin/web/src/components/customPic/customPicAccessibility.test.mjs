import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')

const userAvatar = source.match(/<el-avatar[^>]*v-if="userStore\.userInfo\.headerImg"[^>]*>/)
const defaultAvatar = source.match(/<el-avatar[^>]*v-else[^>]*>/)
const userImage = source.match(/<img[^>]*v-if="userStore\.userInfo\.headerImg"[^>]*class="avatar"[^>]*>/)
const defaultImage = source.match(/<img[^>]*v-else[^>]*class="avatar"[^>]*>/)
const fileImage = source.match(/<img[^>]*:src="file"[^>]*class="file"[^>]*>/)

assert.ok(userAvatar, 'avatar mode should render a user avatar when an image exists')
assert.match(userAvatar[0], /alt="用户头像"/, 'user avatar should have alt text')

assert.ok(defaultAvatar, 'avatar mode should render a default avatar fallback')
assert.match(defaultAvatar[0], /alt="默认用户头像"/, 'default avatar should have alt text')

assert.ok(userImage, 'img mode should render a user avatar image when an image exists')
assert.match(userImage[0], /alt="用户头像"/, 'user avatar image should have alt text')

assert.ok(defaultImage, 'img mode should render a default avatar image fallback')
assert.match(defaultImage[0], /alt="默认用户头像"/, 'default avatar image should have alt text')

assert.ok(fileImage, 'file mode should render a file preview image')
assert.match(fileImage[0], /alt="文件预览"/, 'file preview image should have alt text')

console.log('customPicAccessibility tests passed')
