import { spawn } from 'node:child_process'

const requestOrigin = 'http://localhost'
const unsafeShellChars = /[\r\n<>"|&^%!]/

export default function GvaPositionServer() {
  return {
    name: 'gva-position-server',
    apply: 'serve',
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        const url = new URL(req.url || '/', requestOrigin)
        if (url.pathname !== '/gvaPositionCode') {
          next()
          return
        }

        const target = parseEditorTarget(url.searchParams.get('filePath'))
        if (target) {
          openInEditor(target)
        }

        res.statusCode = 204
        res.end()
      })
    }
  }
}

export const parseEditorTarget = (rawPath) => {
  if (!rawPath || rawPath === 'null' || unsafeShellChars.test(rawPath)) {
    return null
  }

  const lineSeparator = rawPath.lastIndexOf(':')
  const line = lineSeparator > -1 ? rawPath.slice(lineSeparator + 1) : ''
  if (!/^\d+$/.test(line)) {
    return null
  }

  const filePath = rawPath.slice(0, lineSeparator)
  if (!filePath) {
    return null
  }

  return {
    filePath,
    line,
    editorTarget: rawPath,
  }
}

export const getEditorLaunch = (target, editor = process.env.VITE_EDITOR, platform = process.platform) => {
  if (!target) {
    return null
  }

  if (editor === 'webstorm') {
    return {
      command: platform === 'win32' ? 'webstorm64.exe' : 'webstorm64',
      args: ['--line', target.line, target.filePath],
      shell: false,
    }
  }

  return {
    command: platform === 'win32' ? 'cmd.exe' : 'code',
    args: platform === 'win32'
      ? ['/d', '/s', '/c', 'code', '-r', '-g', target.editorTarget]
      : ['-r', '-g', target.editorTarget],
    shell: false,
  }
}

const openInEditor = (target) => {
  const launch = getEditorLaunch(target)
  if (!launch) {
    return
  }

  const child = spawn(launch.command, launch.args, {
    detached: true,
    stdio: 'ignore',
    windowsHide: true,
    shell: launch.shell,
  })
  child.on('error', () => {})
  child.unref()
}
