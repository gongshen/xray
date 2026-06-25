import { devWarn } from '@/utils/devLogger'

const normalizeModules = (modules) => {
  return Object.entries(modules).reduce((result, [key, loader]) => {
    result[key.replace('../', '')] = loader
    return result
  }, {})
}

const viewModules = normalizeModules(import.meta.glob('../view/**/*.vue'))
const pluginModules = normalizeModules(import.meta.glob('../plugin/**/*.vue'))
const fallbackComponent = viewModules['view/error/index.vue']

export const asyncRouterHandle = (asyncRouter) => {
  asyncRouter.forEach(item => {
    if (item.component) {
      const moduleGroup = item.component.split('/')[0]
      if (moduleGroup === 'view') {
        item.component = dynamicImport(viewModules, item.component)
      } else if (moduleGroup === 'plugin') {
        item.component = dynamicImport(pluginModules, item.component)
      }
    } else {
      delete item.component
    }
    if (item.children) {
      asyncRouterHandle(item.children)
    }
  })
}

function dynamicImport(dynamicViewsModules, component) {
  if (!dynamicViewsModules[component]) {
    devWarn('[asyncRouter] component not found: ' + component)
    return fallbackComponent
  }

  return dynamicViewsModules[component]
}
