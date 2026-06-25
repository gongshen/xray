import assert from 'node:assert/strict'
import {
  clearCloseAllHistoryFlag,
  markCloseAllHistory,
  readActiveValue,
  readHistoryStorage,
  shouldCloseAllHistory,
  writeActiveValue,
  writeHistoryStorage,
} from './historyStorage.mjs'

const originalStorage = Object.getOwnPropertyDescriptor(globalThis, 'sessionStorage')

const setStorage = (storage) => {
  Object.defineProperty(globalThis, 'sessionStorage', {
    configurable: true,
    value: storage,
  })
}

function createStorage() {
  const values = new Map()
  return {
    values,
    getItem(key) {
      return values.has(key) ? values.get(key) : null
    },
    setItem(key, value) {
      values.set(key, String(value))
    },
    removeItem(key) {
      values.delete(key)
    },
  }
}

try {
  const fallback = [{ name: 'home' }]
  const storage = createStorage()
  setStorage(storage)

  assert.deepEqual(readHistoryStorage(fallback), fallback)
  storage.setItem('historys', '{bad-json')
  assert.deepEqual(readHistoryStorage(fallback), fallback)
  storage.setItem('historys', JSON.stringify({ name: 'bad' }))
  assert.deepEqual(readHistoryStorage(fallback), fallback)

  writeHistoryStorage([{ name: 'dashboard' }])
  assert.deepEqual(readHistoryStorage(fallback), [{ name: 'dashboard' }])

  assert.equal(readActiveValue(), '')
  writeActiveValue('route-key')
  assert.equal(readActiveValue(), 'route-key')

  assert.equal(shouldCloseAllHistory(), false)
  markCloseAllHistory()
  assert.equal(storage.getItem('needCloseAll'), 'true')
  assert.equal(shouldCloseAllHistory(), true)
  clearCloseAllHistoryFlag()
  assert.equal(shouldCloseAllHistory(), false)

  setStorage({
    getItem() { throw new Error('blocked') },
    setItem() { throw new Error('blocked') },
    removeItem() { throw new Error('blocked') },
  })
  assert.deepEqual(readHistoryStorage(fallback), fallback)
  assert.equal(readActiveValue(), '')
  assert.equal(shouldCloseAllHistory(), false)
  assert.doesNotThrow(() => writeHistoryStorage([{ name: 'safe' }]))
  assert.doesNotThrow(() => writeActiveValue('safe'))
  assert.doesNotThrow(() => markCloseAllHistory())
  assert.doesNotThrow(() => clearCloseAllHistoryFlag())
} finally {
  if (originalStorage) {
    Object.defineProperty(globalThis, 'sessionStorage', originalStorage)
  } else {
    delete globalThis.sessionStorage
  }
}

console.log('historyStorage tests passed')
