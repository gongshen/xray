import assert from 'node:assert/strict'
import fs from 'node:fs'

const source = fs.readFileSync(new URL('./clipboard.js', import.meta.url), 'utf8')
const module = { exports: {} }
Function('module', source.replace(/export const /g, 'const ') + '; module.exports = { copyTextToClipboard }')(module)

const { copyTextToClipboard } = module.exports
const originalNavigator = Object.getOwnPropertyDescriptor(globalThis, 'navigator')
const originalDocument = Object.getOwnPropertyDescriptor(globalThis, 'document')

const setGlobal = (key, value) => {
  Object.defineProperty(globalThis, key, {
    configurable: true,
    value,
  })
}

try {
  const clipboardWrites = []
  setGlobal('navigator', {
    clipboard: {
      writeText: async(value) => clipboardWrites.push(value),
    },
  })
  setGlobal('document', undefined)

  assert.equal(await copyTextToClipboard('hello'), true)
  assert.deepEqual(clipboardWrites, ['hello'])
  assert.equal(await copyTextToClipboard(''), false)

  let appended = null
  let removed = null
  const fakeTextArea = {
    style: {},
    setAttribute(name, value) {
      this[name] = value
    },
    select() {
      this.selected = true
    },
    setSelectionRange(start, end) {
      this.selection = [start, end]
    },
  }
  setGlobal('navigator', {
    clipboard: {
      writeText: async() => {
        throw new Error('denied')
      },
    },
  })
  setGlobal('document', {
    body: {
      appendChild(element) {
        appended = element
      },
      removeChild(element) {
        removed = element
      },
    },
    createElement(tagName) {
      assert.equal(tagName, 'textarea')
      return fakeTextArea
    },
    execCommand(command) {
      assert.equal(command, 'copy')
      return true
    },
  })

  assert.equal(await copyTextToClipboard(123), true)
  assert.equal(appended, fakeTextArea)
  assert.equal(removed, fakeTextArea)
  assert.equal(fakeTextArea.value, '123')
  assert.deepEqual(fakeTextArea.selection, [0, 3])

  setGlobal('navigator', undefined)
  setGlobal('document', undefined)
  assert.equal(await copyTextToClipboard('no-dom'), false)
} finally {
  if (originalNavigator) {
    Object.defineProperty(globalThis, 'navigator', originalNavigator)
  } else {
    delete globalThis.navigator
  }
  if (originalDocument) {
    Object.defineProperty(globalThis, 'document', originalDocument)
  } else {
    delete globalThis.document
  }
}

console.log('clipboard tests passed')
