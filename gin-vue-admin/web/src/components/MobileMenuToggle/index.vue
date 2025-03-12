<template>
  <div class="mobile-menu-toggle" v-if="isMobile" @click="toggleMobileMenu">
    <el-icon size="24">
      <component :is="isMenuVisible ? 'Close' : 'Menu'" />
    </el-icon>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { emitter } from '@/utils/bus.js'

const isMobile = ref(false)
const isMenuVisible = ref(false)

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

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  
  // Listen for route changes to close menu
  emitter.on('routeChange', () => {
    if (isMenuVisible.value) {
      toggleMobileMenu()
    }
  })
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
  emitter.off('routeChange')
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
    background: var(--el-color-primary);
    color: white;
    align-items: center;
    justify-content: center;
    box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
    cursor: pointer;
  }
}
</style>
