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

const isMobile = ref(false)
const isMenuVisible = ref(false)
let disposeResize = null
let disposeRouteChange = null

const toggleMobileMenu = () => {
  isMenuVisible.value = !isMenuVisible.value
  emitter.emit('toggleMobileMenu', isMenuVisible.value)
  
  // Toggle mobile-visible class on aside element
  const asideElement = document.querySelector('.gva-aside')
  if (asideElement) {
    if (isMenuVisible.value) {
      asideElement.classList.add('mobile-visible')
    } else {
      asideElement.classList.remove('mobile-visible')
    }
  }
}

const checkMobile = () => {
  const screenWidth = document.body.clientWidth
  isMobile.value = screenWidth < 768
}
const closeMenuOnRouteChange = () => {
  if (isMenuVisible.value) {
    toggleMobileMenu()
  }
}

onMounted(() => {
  checkMobile()
  disposeResize = bindWindowEvent(window, 'resize', checkMobile)
  
  // Listen for route changes to close menu
  disposeRouteChange = bindEmitterHandler(emitter, 'routeChange', closeMenuOnRouteChange)
})

onUnmounted(() => {
  disposeResize?.()
  disposeResize = null
  disposeRouteChange?.()
  disposeRouteChange = null
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
