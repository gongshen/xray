import chalk from 'chalk'

export const viteLogo = (env = {}, enabled = true) => {
  if (!enabled) {
    return
  }

  const lines = [
    '> Welcome to Gin-Vue-Admin',
    '> Version: v2.5.5',
    '> Runtime mode: local resources',
    '> External resources: disabled',
    `> Swagger: http://127.0.0.1:${env.VITE_SERVER_PORT}/swagger/index.html`,
    `> Web: http://127.0.0.1:${env.VITE_CLI_PORT}`,
    '> Static resources: local bundle',
  ]

  lines.forEach((line) => {
    console.log(chalk.green(line))
  })
  console.log('')
}
