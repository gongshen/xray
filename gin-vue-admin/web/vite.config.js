import legacyPlugin from '@vitejs/plugin-legacy'
import vuePlugin from '@vitejs/plugin-vue'
import * as dotenv from 'dotenv'
import * as fs from 'node:fs'
import * as path from 'node:path'
import AutoImport from 'unplugin-auto-import/vite'
import { ElementPlusResolver } from 'unplugin-vue-components/resolvers'
import Components from 'unplugin-vue-components/vite'
import Banner from 'vite-plugin-banner'
import webConfig from './src/core/config'
import { viteLogo } from './vitePlugin/viteLogo'

const vendorChunkRules = [
  ['vue-vendor', ['node_modules/vue', 'node_modules/@vue', 'node_modules/vue-router', 'node_modules/pinia']],
  ['echarts', ['node_modules/echarts', 'node_modules/zrender']],
  ['http-vendor', ['node_modules/axios', 'node_modules/qs']],
]

const manualChunks = (id) => {
  const normalizedId = id.split(path.sep).join('/')
  const matched = vendorChunkRules.find(([, packages]) => packages.some((packageName) => normalizedId.includes(packageName)))
  if (matched) {
    return matched[0]
  }
  if (normalizedId.includes('node_modules')) {
    return 'vendor'
  }
}

const injectElementTheme = (source, filename = '') => {
  const normalizedFilename = filename.replace(/\\/g, '/')
  if (normalizedFilename.endsWith('/node_modules/element-plus/theme-chalk/src/common/var.scss')) {
    return source
  }
  if (normalizedFilename.includes('/node_modules/element-plus/theme-chalk/')) {
    return `@use "@/style/element/index.scss" as *;\n${source}`
  }
  return source
}

const createDevInspectorPlugins = async (enabled) => {
  if (!enabled) {
    return []
  }
  const [
    { default: GvaPositionServer },
    { default: GvaPosition },
  ] = await Promise.all([
    import('./vitePlugin/codeServer/index.js'),
    import('./vitePlugin/gvaPosition/index.js'),
  ])
  return [GvaPositionServer(), GvaPosition()]
}

const createFullImportPlugin = async () => {
  const { default: fullImportPlugin } = await import('./vitePlugin/fullImport/fullImport.js')
  return fullImportPlugin()
}

// @see https://cn.vitejs.dev/config/
export default async ({
  command,
  mode
}) => {
  const NODE_ENV = process.env.NODE_ENV || 'development'
  const isServe = command === 'serve'
  const envFiles = [
    `.env.${NODE_ENV}`
  ]
  for (const file of envFiles) {
    if (!fs.existsSync(file)) {
      continue
    }
    const envConfig = dotenv.parse(fs.readFileSync(file))
    for (const k in envConfig) {
      process.env[k] = envConfig[k]
    }
  }

  viteLogo(process.env, isServe && webConfig.showViteLogo)

  const timestamp = Date.now()

  const optimizeDeps = isServe ? {
    include: [
      'vue',
      'vue-router',
      'pinia',
      'axios',
      'element-plus',
      '@element-plus/icons-vue',
    ]
  } : {}

  const alias = {
    '@': path.resolve(__dirname, './src'),
    'vue$': 'vue/dist/vue.runtime.esm-bundler.js',
  }

  const esbuild = {}

  const config = {
    base: './', // index.html file location
    root: './', // imported resource root
    resolve: {
      alias,
    },
    define: {
      'process.env': {}
    },
    server: {
      // Set to false when using docker-compose development mode.
      open: true,
      port: process.env.VITE_CLI_PORT,
      proxy: {
        [process.env.VITE_BASE_API]: {
          target: `${process.env.VITE_BASE_PATH}:${process.env.VITE_SERVER_PORT}/`,
          changeOrigin: true,
          rewrite: path => path.replace(new RegExp('^' + process.env.VITE_BASE_API), ''),
        }
      },
    },
    build: {
      minify: 'terser',
      manifest: false,
      sourcemap: false,
      outDir: 'dist',
      chunkSizeWarningLimit: 1000,
      reportCompressedSize: false,
      terserOptions: {
        compress: {
          drop_console: true,
          drop_debugger: true,
        },
        format: {
          comments: false,
        }
      },
      rollupOptions: {
        output: {
          manualChunks,
        }
      },
      rolldownOptions: {
        checks: {
          invalidAnnotation: false,
          pluginTimings: false,
        },
      },
    },
    esbuild,
    optimizeDeps,
    plugins: [
      ...(await createDevInspectorPlugins(isServe)),
      legacyPlugin({
        targets: ['Android > 39', 'Chrome >= 60', 'Safari >= 10.1', 'iOS >= 10.3', 'Firefox >= 54', 'Edge >= 15'],
      }),
      vuePlugin(),
      [Banner(`\n Build based on gin-vue-admin \n Time : ${timestamp}`)]
    ],
    css: {
      preprocessorOptions: {
        scss: {
          additionalData: injectElementTheme,
        }
      }
    },
  }

  if (NODE_ENV === 'development') {
    config.plugins.push(
      await createFullImportPlugin()
    )
  } else {
    config.plugins.push(AutoImport({
      resolvers: [ElementPlusResolver()]
    }),
    Components({
      resolvers: [ElementPlusResolver({
        importStyle: 'sass'
      })]
    }))
  }
  return config
}
