import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./person.vue', import.meta.url), 'utf8')

const phoneButton = source.match(/<button[^>]*class="link-action"[^>]*@click="changePhoneFlag = true"[^>]*>/)
const emailButton = source.match(/<button[^>]*class="link-action"[^>]*@click="changeEmailFlag = true"[^>]*>/)
const passwordButton = source.match(/<button[^>]*class="link-action"[^>]*@click="showPassword = true"[^>]*>/)

assert.ok(phoneButton, 'phone change action should use a native button element')
assert.match(phoneButton[0], /type="button"/, 'phone change button should not submit forms')
assert.match(phoneButton[0], /aria-label="修改密保手机"/, 'phone change button should have an explicit accessible name')

assert.ok(emailButton, 'email change action should use a native button element')
assert.match(emailButton[0], /type="button"/, 'email change button should not submit forms')
assert.match(emailButton[0], /aria-label="修改密保邮箱"/, 'email change button should have an explicit accessible name')

assert.ok(passwordButton, 'password change action should use a native button element')
assert.match(passwordButton[0], /type="button"/, 'password change button should not submit forms')
assert.match(passwordButton[0], /aria-label="修改个人密码"/, 'password change button should have an explicit accessible name')

assert.doesNotMatch(source, /<a[^>]*href="javascript:void\(0\)"[^>]*@click="(changePhoneFlag = true|changeEmailFlag = true|showPassword = true)"/, 'account actions should not use javascript links for button behavior')

console.log('personAccountActionsAccessibility tests passed')
