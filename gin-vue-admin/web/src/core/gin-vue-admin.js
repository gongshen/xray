/*
 * gin-vue-admin web框架组
 *
 * */
// 加载网站配置文件夹
import { register } from './global'

export default {
  install: (app) => {
    register(app)
    console.log(`
       欢迎使用 Gin-Vue-Admin
       当前版本:v2.5.5
       加群方式:微信：shouzi_1994 QQ群：622360840
       当前运行模式:本地资源
       外部资源加载:已关闭
       默认自动化文档地址:http://127.0.0.1:${import.meta.env.VITE_SERVER_PORT}/swagger/index.html
       默认前端文件运行地址:http://127.0.0.1:${import.meta.env.VITE_CLI_PORT}
       静态资源来源:本地打包文件
    `)
  }
}
