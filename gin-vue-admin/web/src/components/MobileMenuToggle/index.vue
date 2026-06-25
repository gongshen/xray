<template>
  <button
    v-if="isMobile"
    type="button"
    class="mobile-menu-toggle"
    :aria-expanded="isMenuVisible"
    :aria-label="isMenuVisible ? '关闭移动端菜单' : '打开移动端菜单'"
    @click="toggleMobileMenu"
  >
    <el-icon size="24">
      <component :is="isMenuVisible ? 'Close' : 'Menu'" />
    </el-icon>
  </button>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { emitter } from '@/utils/bus.js'
import {
  bindEmitterHandler,
  bindWindowEvent,
} from '@/utils/eventLifecycle.mjs'
import { lockBodyScroll } from '@/utils/bodyScrollLock.mjs'

const isMobile = ref(false)
const isMenuVisible = ref(false)
let disposeResize = null
let disposeRouteChange = null
let disposeToggleMobileMenu = null
let releaseBodyScrollLock = null

const releaseMobileMenuScroll = () => {
  releaseBodyScrollLock?.()
  releaseBodyScrollLock = null
}

const syncAsideVisibility = (visible) => {
  const asideElement = document.querySelector('.gva-aside')
  if (asideElement) {
    asideElement.classList.toggle('mobile-visible', visible)
  }

  if (visible) {
    releaseMobileMenuScroll()
    releaseBodyScrollLock = lockBodyScroll()
  } else {
    releaseMobileMenuScroll()
  }
}

const setMobileMenuVisible = (visible, shouldEmit = false) => {
  const nextVisible = Boolean(visible)
  if (isMenuVisible.value !== nextVisible) {
    isMenuVisible.value = nextVisible
    syncAsideVisibility(nextVisible)
  }
  if (shouldEmit) {
    emitter.emit('toggleMobileMenu', nextVisible)
  }
}

const toggleMobileMenu = () => {
  setMobileMenuVisible(!isMenuVisible.value, true)
}

const checkMobile = () => {
  const screenWidth = document.body.clientWidth
  isMobile.value = screenWidth < 768
}
const closeMenuOnRouteChange = () => {
  setMobileMenuVisible(false, true)
}
const handleToggleMobileMenu = (visible) => {
  setMobileMenuVisible(visible)
}

onMounted(() => {
  checkMobile()
  disposeResize = bindWindowEvent(window, 'resize', checkMobile)
  
  // Listen for route changes and overlay close actions to close menu
  disposeRouteChange = bindEmitterHandler(emitter, 'routeChange', closeMenuOnRouteChange)
  disposeToggleMobileMenu = bindEmitterHandler(emitter, 'toggleMobileMenu', handleToggleMobileMenu)
})

onUnmounted(() => {
  setMobileMenuVisible(false, true)
  disposeResize?.()
  disposeResize = null
  disposeRouteChange?.()
  disposeRouteChange = null
  disposeToggleMobileMenu?.()
  disposeToggleMobileMenu = null
})
</script>

<style lang="scss" scoped>
.mobile-menu-toggle {
  display: none;
  
  @media screen and (max-width: 768px) {
    display: flex;
    position: fixed;
    top: 10px;
    left: 10px;
    z-index: 1001;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    border: 0;
    background: var(--el-color-primary);
    color: white;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    cursor: pointer;
    padding: 0;
  }
}
</style>
