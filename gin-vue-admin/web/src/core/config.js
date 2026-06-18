/**
 * 网站配置文件
 */

const config = {
  appName: 'Gin-Vue-Admin',
  appLogo: new URL('../assets/logo_login.png', import.meta.url).href,
  showViteLogo: true
}

export const viteLogo = (env) => {
  if (config.showViteLogo) {
    const chalk = require('chalk')
    console.log(
      chalk.green(
        `> 欢迎使用Gin-Vue-Admin`
      )
    )
    console.log(
      chalk.green(
        `> 当前版本:v2.5.5`
      )
    )
    console.log(
      chalk.green(
        `> 加群方式:微信：shouzi_1994 QQ群：622360840`
      )
    )
    console.log(
      chalk.green(
        `> 当前运行模式:本地资源`
      )
    )
    console.log(
      chalk.green(
        `> 外部资源加载:已关闭`
      )
    )
    console.log(
      chalk.green(
        `> 默认自动化文档地址:http://127.0.0.1:${env.VITE_SERVER_PORT}/swagger/index.html`
      )
    )
    console.log(
      chalk.green(
        `> 默认前端文件运行地址:http://127.0.0.1:${env.VITE_CLI_PORT}`
      )
    )
    console.log(
      chalk.green(
        `> 静态资源来源:本地打包文件`
      )
    )
    console.log('\n')
  }
}

export default config
