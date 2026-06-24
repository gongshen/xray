import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./chatTable.vue', import.meta.url), 'utf8')

const selectChatButton = source.match(/<button[^>]*class="history-select"[^>]*@click="selectChat\(index\)"[^>]*>/)
const deleteChatButton = source.match(/<button[^>]*class="delete-icon"[^>]*@click\.stop="deleteChat\(index\)"[^>]*>/)
const sendButton = source.match(/<el-button[^>]*class="send-button"[^>]*@click="handleQueryTable"[^>]*>/)

assert.ok(selectChatButton, 'chat history selection should use a native button element')
assert.match(selectChatButton[0], /type="button"/, 'chat history selection button should not submit forms')
assert.match(selectChatButton[0], /:aria-current="currentChatIndex === index \? 'true' : undefined"/, 'active chat history item should expose current state')

assert.ok(deleteChatButton, 'chat history delete action should use a native button element')
assert.match(deleteChatButton[0], /type="button"/, 'chat history delete button should not submit forms')
assert.match(deleteChatButton[0], /:aria-label="`删除对话：\$\{chat\.title \|\| '新对话'\}`"/, 'chat history delete button should include the chat title in its accessible name')
assert.match(deleteChatButton[0], /@click\.stop="deleteChat\(index\)"/, 'chat history delete button should keep stop propagation behavior')

assert.doesNotMatch(source, /<div[^>]*class="history-item"[^>]*@click="selectChat\(index\)"/, 'chat history row should not be a clickable div')
assert.doesNotMatch(source, /<el-icon[^>]*class="delete-icon"[^>]*@click\.stop="deleteChat\(index\)"/, 'chat history delete action should not be a clickable icon component')

assert.ok(sendButton, 'chat send action should keep the Element Plus button')
assert.match(sendButton[0], /aria-label="发送问题"/, 'icon-only chat send button should have an accessible name')
assert.match(sendButton[0], /:disabled="!form\.chat \|\| !form\.dbname"/, 'chat send button should keep its disabled guard')

console.log('chatTableAccessibility tests passed')
