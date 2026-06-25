import fs from 'node:fs'
import path from 'node:path'

const root = new URL('./', import.meta.url)
const failures = []

function walk(dirUrl, predicate) {
  const entries = fs.readdirSync(dirUrl, { withFileTypes: true })
  return entries.flatMap((entry) => {
    const url = new URL(`${entry.name}${entry.isDirectory() ? '/' : ''}`, dirUrl)
    if (entry.isDirectory()) return walk(url, predicate)
    if (!entry.isFile()) return []
    return predicate(entry.name) ? [url] : []
  })
}

function lineOf(source, index) {
  return source.slice(0, index).split(/\r?\n/).length
}

function relativePath(fileUrl) {
  return path.relative(path.dirname(new URL(import.meta.url).pathname), fileUrl.pathname).replace(/\\/g, '/')
}

function nearestOpenFormItemHasLabel(source, index) {
  const before = source.slice(0, index)
  const open = before.lastIndexOf('<el-form-item')
  const close = before.lastIndexOf('</el-form-item>')
  if (open < 0 || close > open) return false
  const tagEnd = source.indexOf('>', open)
  if (tagEnd < 0 || tagEnd > index) return false
  const attrs = source.slice(open, tagEnd)
  return /\blabel=|:label=/.test(attrs) && !/label=""/.test(attrs)
}

for (const fileUrl of walk(root, (name) => name.endsWith('.vue'))) {
  const relative = relativePath(fileUrl)
  const source = fs.readFileSync(fileUrl, 'utf8').replace(/<!--[\s\S]*?-->/g, '')

  for (const match of source.matchAll(/<img\b([^>]*)>/g)) {
    const attrs = match[1]
    if (!/\balt=/.test(attrs) || /\balt(?:\s|>|="")/.test(attrs)) {
      failures.push(`${relative}:${lineOf(source, match.index)} image must have non-empty alt text`)
    }
  }

  for (const match of source.matchAll(/<(div|span|img|el-image|el-row|el-icon)\b([^>]*)>/g)) {
    const [, tag, attrs] = match
    if (/@click=/.test(attrs)) {
      failures.push(`${relative}:${lineOf(source, match.index)} ${tag} must not bind click directly`)
    }
  }

  for (const match of source.matchAll(/(?:v-for="\([^)]*,\s*(?:index|key)\)[^"]*"[\s\S]{0,160}:key="(?:index|key)"|:key="(?:index|key|\$index|scope\.\$index)")/g)) {
    failures.push(`${relative}:${lineOf(source, match.index)} list rendering must use a stable key instead of index/key`)
  }

  for (const match of source.matchAll(/<button\b([\s\S]*?)(?:>([\s\S]*?)<\/button>|\/>)/g)) {
    const attrs = match[1]
    const body = match[2] || ''
    const visibleBody = body
      .replace(/<[^>]+>/g, ' ')
      .replace(/\{\{[\s\S]*?\}\}/g, 'value')
      .replace(/\s+/g, ' ')
      .trim()
    const hasName = /\baria-label=/.test(attrs) || /\baria-labelledby=/.test(attrs) || visibleBody.length > 0
    if (!/\btype=/.test(attrs)) {
      failures.push(`${relative}:${lineOf(source, match.index)} button must declare type`)
    }
    if (!hasName) {
      failures.push(`${relative}:${lineOf(source, match.index)} button must have text or an accessible name`)
    }
  }

  for (const match of source.matchAll(/<el-(input(?:-number)?|select|date-picker|cascader|switch)\b([\s\S]*?)>/g)) {
    const attrs = match[2]
    const hasDirectName = /\baria-label=|\baria-labelledby=|\bplaceholder=|:placeholder=|\bactive-text=|\binactive-text=/.test(attrs)
    const hasFormLabel = nearestOpenFormItemHasLabel(source, match.index)
    if (!hasDirectName && !hasFormLabel) {
      failures.push(`${relative}:${lineOf(source, match.index)} form control must have a label, placeholder, or aria-label`)
    }
  }

  for (const match of source.matchAll(/<el-table(?!-)\b([\s\S]*?)>/g)) {
    const attrs = match[1]
    if (!/\baria-label=/.test(attrs)) {
      failures.push(`${relative}:${lineOf(source, match.index)} el-table must have aria-label`)
    }
    if (!/(?:^|\s)(?::empty-text|empty-text)=/.test(attrs)) {
      failures.push(`${relative}:${lineOf(source, match.index)} el-table must define empty text`)
    }
  }

  for (const match of source.matchAll(/<el-pagination\b([\s\S]*?)>/g)) {
    if (!/\baria-label=/.test(match[1])) {
      failures.push(`${relative}:${lineOf(source, match.index)} el-pagination must have aria-label`)
    }
  }

  for (const match of source.matchAll(/<[^>]+\bv-loading="([^"]+)"[^>]*>/g)) {
    const tag = match[0]
    if (!/\baria-busy=|:aria-busy=|\belement-loading-text=|:element-loading-text=/.test(tag)) {
      failures.push(`${relative}:${lineOf(source, match.index)} loading region must expose busy state or loading text`)
    }
  }

  if (/href="(?:javascript|#)/.test(source)) {
    failures.push(`${relative} must not use javascript or hash pseudo-links`)
  }
  if (/\.native(?:=|\b)/.test(source)) {
    failures.push(`${relative} must not use Vue 2 native event modifiers`)
  }
  if (/JSON\.parse\(sessionStorage\.getItem\('historys'\)\)|sessionStorage\.setItem\('historys'/.test(source)) {
    failures.push(`${relative} history storage should use historyStorage helpers`)
  }
  if (/::v-deep|\/deep\/|>>>/.test(source)) {
    failures.push(`${relative} must use Vue 3 :deep() syntax for deep selectors`)
  }
  for (const style of source.matchAll(/<style\b([^>]*)>([\s\S]*?)<\/style>/g)) {
    if (!/\bscoped\b/.test(style[1]) && /:deep\(/.test(style[2])) {
      failures.push(`${relative}:${lineOf(source, style.index)} non-scoped styles must not use Vue-only :deep() selectors`)
    }
  }
  if (/(?:aria-label|empty-text|alt)="[^"]*\?\?[^"]*"/.test(source)) {
    failures.push(`${relative} contains likely mojibake in accessible text`)
  }
}

for (const fileUrl of walk(root, (name) => /\.(vue|js|mjs|scss)$/.test(name) && !name.endsWith('.test.mjs'))) {
  const relative = relativePath(fileUrl)
  const source = fs.readFileSync(fileUrl, 'utf8')
  if (/transition:\s*all\b/.test(source)) {
    failures.push(`${relative} must not use transition: all`)
  }
  if (relative !== 'utils/clipboard.js' && /navigator\.clipboard|document\.execCommand|createElement\(['"]textarea['"]\)/.test(source)) {
    failures.push(`${relative} clipboard access should go through utils/clipboard`)
  }
  if (relative !== 'utils/bodyScrollLock.mjs' && /document\.body\.style\.overflow\s*=/.test(source)) {
    failures.push(`${relative} body scroll locking should use utils/bodyScrollLock`)
  }
  if (/:hover\s*\{[^}]*transform:\s*(?:scale|translateY\(-[24]px)/s.test(source)) {
    failures.push(`${relative} hover states must not move or scale controls`)
  }
  for (const match of source.matchAll(/border-radius:\s*(?:1[0-9]|[2-9][0-9])px/g)) {
    failures.push(`${relative}:${lineOf(source, match.index)} custom radius should stay at 8px or less`)
  }
}

if (failures.length) {
  console.error(failures.join('\n'))
  process.exit(1)
}

console.log('frontendAccessibility tests passed')
