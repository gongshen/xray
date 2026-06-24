import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./index.vue', import.meta.url), 'utf8')
const styleSource = fs.readFileSync(new URL('../../style/newLogin.scss', import.meta.url), 'utf8')

const logoImage = source.match(/<img[^>]*class="login_panel_form_title_logo"[^>]*>/)
const captchaButton = source.match(/<button[^>]*class="captcha-refresh"[^>]*@click="loginVerify\(\)"[^>]*>/)
const captchaImage = source.match(/<img[^>]*:src="picPath"[^>]*>/)
const captchaRefreshStyle = styleSource.match(/\.captcha-refresh\s*\{([\s\S]*?)\n\s*\}/)

assert.ok(logoImage, 'login logo should be rendered')
assert.match(logoImage[0], /alt="应用标识"/, 'login logo should have alt text')

assert.ok(captchaButton, 'captcha refresh should use a native button')
assert.match(captchaButton[0], /type="button"/, 'captcha refresh button should not submit the login form')
assert.match(captchaButton[0], /aria-label="刷新验证码"/, 'captcha refresh button should have an accessible name')

assert.ok(captchaImage, 'captcha image should be rendered inside the refresh button')
assert.match(captchaImage[0], /alt="验证码图片"/, 'captcha image should have descriptive alt text')
assert.doesNotMatch(source, /<img[^>]*@click="loginVerify\(\)"/, 'captcha image should not be directly clickable')

assert.ok(captchaRefreshStyle, 'captcha refresh button should have reset styles')
for (const declaration of [
  'border: 0;',
  'background: transparent;',
  'padding: 0;',
  'cursor: pointer;',
  'width: 100%;',
  'height: 100%;',
  'display: block;'
]) {
  assert.ok(
    captchaRefreshStyle[1].includes(declaration),
    `captcha refresh style should include ${declaration}`
  )
}

console.log('loginAccessibility tests passed')
