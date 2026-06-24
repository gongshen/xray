<template>
  <transition name="fade">
    <button
      v-if="visible"
      type="button"
      class="responsive-overlay"
      aria-label="关闭移动端菜单"
      @click="closeOverlay"
    />
  </transition>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { emitter } from '@/utils/bus.js'
import { bindEmitterHandler } from '@/utils/eventLifecycle.mjs'

const visible = ref(false)
let disposeToggleMobileMenu = null

const closeOverlay = () => {
  emitter.emit('toggleMobileMenu', false)
}
const handleToggleMobileMenu = (isVisible) => {
  visible.value = isVisible
}

onMounted(() => {
  disposeToggleMobileMenu = bindEmitterHandler(emitter, 'toggleMobileMenu', handleToggleMobileMenu)
})

onUnmounted(() => {
  disposeToggleMobileMenu?.()
  disposeToggleMobileMenu = null
})
</script>

<style lang="scss" scoped>
.responsive-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
  padding: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
