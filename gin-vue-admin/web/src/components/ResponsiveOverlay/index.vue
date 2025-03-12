<template>
  <transition name="fade">
    <div class="responsive-overlay" v-if="visible" @click="closeOverlay"></div>
  </transition>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { emitter } from '@/utils/bus.js'

const visible = ref(false)

const closeOverlay = () => {
  emitter.emit('toggleMobileMenu', false)
}

onMounted(() => {
  emitter.on('toggleMobileMenu', (isVisible) => {
    visible.value = isVisible
  })
})

onUnmounted(() => {
  emitter.off('toggleMobileMenu')
})
</script>

<style lang="scss" scoped>
.responsive-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 999;
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
