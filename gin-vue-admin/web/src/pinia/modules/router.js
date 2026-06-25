import { asyncRouterHandle } from '@/utils/asyncRouter'
import { emitter } from '@/utils/bus.js'
import { asyncMenu } from '@/api/menu'
import { defineStore } from 'pinia'
import { ref } from 'vue'

const routerListArr = []
const notLayoutRouterArr = []
const nameMap = {}
const pendingNameMap = {}
let keepAliveCacheVersion = 0

const clearObject = (target) => {
  Object.keys(target).forEach((key) => {
    delete target[key]
  })
}

const resetRouterCache = (routeMap) => {
  keepAliveCacheVersion++
  routerListArr.length = 0
  notLayoutRouterArr.length = 0
  clearObject(nameMap)
  clearObject(pendingNameMap)
  clearObject(routeMap)
}

const routeNeedsKeepAlive = (route) => {
  return Boolean(route?.meta?.keepAlive || route?.children?.some(child => child?.meta?.keepAlive))
}

const formatRouter = (routes, routeMap) => {
  routes && routes.forEach(item => {
    item.meta = item.meta || {}
    if ((!item.children || item.children.every(ch => ch.hidden)) && item.name !== '404' && !item.hidden) {
      routerListArr.push({ label: item.meta.title, value: item.name })
    }
    item.meta.btns = item.btns
    item.meta.hidden = item.hidden
    if (item.meta.defaultMenu === true) {
      notLayoutRouterArr.push({
        ...item,
        path: `/${item.path}`,
      })
    } else {
      routeMap[item.name] = item
      if (item.children && item.children.length > 0) {
        formatRouter(item.children, routeMap)
      }
    }
  })
}

const resolveKeepAliveName = (routeName, routeRecord) => {
  if (!routeName || nameMap[routeName]) {
    return Promise.resolve(nameMap[routeName] || '')
  }

  if (pendingNameMap[routeName]) {
    return pendingNameMap[routeName]
  }

  if (typeof routeRecord?.component !== 'function') {
    nameMap[routeName] = routeName
    return Promise.resolve(routeName)
  }

  const cacheVersion = keepAliveCacheVersion
  const pendingPromise = routeRecord.component()
    .then((module) => {
      if (cacheVersion !== keepAliveCacheVersion) {
        return ''
      }
      const componentName = module?.default?.name || routeName
      nameMap[routeName] = componentName
      return componentName
    })
    .catch(() => '')
    .finally(() => {
      if (pendingNameMap[routeName] === pendingPromise) {
        delete pendingNameMap[routeName]
      }
    })

  pendingNameMap[routeName] = pendingPromise
  return pendingPromise
}

const collectResolvedKeepAliveNames = (history, routeMap) => {
  return Array.from(new Set((history || []).map((item) => {
    const routeRecord = routeMap[item.name]
    if (!routeNeedsKeepAlive(item) && !routeNeedsKeepAlive(routeRecord)) {
      return ''
    }
    return nameMap[item.name] || ''
  }).filter(Boolean)))
}

export const useRouterStore = defineStore('router', () => {
  const keepAliveRouters = ref([])
  const asyncRouterFlag = ref(0)
  let keepAliveResolveVersion = 0

  const asyncRouters = ref([])
  const routerList = ref(routerListArr)
  const routeMap = ({})
  let setAsyncRouterPromise = null

  const setKeepAliveRouters = (history = []) => {
    const version = ++keepAliveResolveVersion
    const pending = []

    history.forEach((item) => {
      const routeRecord = routeMap[item.name]
      if (!routeNeedsKeepAlive(item) && !routeNeedsKeepAlive(routeRecord)) {
        return
      }
      if (!nameMap[item.name]) {
        pending.push(resolveKeepAliveName(item.name, routeRecord))
      }
    })

    keepAliveRouters.value = collectResolvedKeepAliveNames(history, routeMap)

    if (pending.length) {
      Promise.all(pending).then(() => {
        if (version === keepAliveResolveVersion) {
          keepAliveRouters.value = collectResolvedKeepAliveNames(history, routeMap)
        }
      })
    }
  }
  emitter.on('setKeepAlive', setKeepAliveRouters)

  // Load dynamic routes from backend
  const buildAsyncRouter = async() => {
    const asyncRouterRes = await asyncMenu()
    const asyncRouter = asyncRouterRes.data?.menus || []
    resetRouterCache(routeMap)
    keepAliveResolveVersion++
    keepAliveRouters.value = []
    const baseRouter = [{
      path: '/layout',
      name: 'layout',
      component: 'view/layout/index.vue',
      meta: {
        title: '\u5e95\u5c42layout'
      },
      children: []
    }]
    if (!asyncRouter.some(item => item.name === 'Reload')) {
      asyncRouter.push({
        path: 'reload',
        name: 'Reload',
        hidden: true,
        meta: {
          title: '',
          closeTab: true,
        },
        component: 'view/error/reload.vue'
      })
    }
    formatRouter(asyncRouter, routeMap)
    baseRouter[0].children = asyncRouter
    if (notLayoutRouterArr.length !== 0) {
      baseRouter.push(...notLayoutRouterArr)
    }
    asyncRouterHandle(baseRouter)
    asyncRouters.value = baseRouter
    routerList.value = [...routerListArr]
    asyncRouterFlag.value++
    return true
  }

  const SetAsyncRouter = async() => {
    if (setAsyncRouterPromise) {
      return setAsyncRouterPromise
    }

    setAsyncRouterPromise = buildAsyncRouter()
    try {
      return await setAsyncRouterPromise
    } finally {
      setAsyncRouterPromise = null
    }
  }

  return {
    asyncRouters,
    routerList,
    keepAliveRouters,
    asyncRouterFlag,
    SetAsyncRouter,
    routeMap
  }
})
