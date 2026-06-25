// responsive.js - Utility functions for responsive design

import { ref, onMounted, onUnmounted } from 'vue'
import { emitter } from '@/utils/bus.js'
import {
  bindEmitterHandler,
  bindWindowEvent,
} from '@/utils/eventLifecycle.mjs'
import { lockBodyScroll } from '@/utils/bodyScrollLock.mjs'

/**
 * Hook to detect mobile devices and handle responsive behavior
 * @returns {Object} - Mobile state and utility functions
 */
export function useResponsive() {
  const isMobile = ref(false)
  const isMenuVisible = ref(false)
  let disposeResize = null
  let disposeRouteChange = null
  let releaseBodyScrollLock = null

  // Check if current viewport is mobile
  const checkMobile = () => {
    const screenWidth = document.body.clientWidth
    const newIsMobile = screenWidth < 768
    
    if (isMobile.value !== newIsMobile) {
      isMobile.value = newIsMobile
      emitter.emit('mobileStateChanged', isMobile.value)
    }
  }
  const releaseMobileMenuScroll = () => {
    releaseBodyScrollLock?.()
    releaseBodyScrollLock = null
  }

  // Toggle mobile menu visibility
  const toggleMobileMenu = () => {
    isMenuVisible.value = !isMenuVisible.value
    emitter.emit('toggleMobileMenu', isMenuVisible.value)
    
    // Toggle mobile-visible class on aside element
    const asideElement = document.querySelector('.gva-aside')
    if (asideElement) {
      if (isMenuVisible.value) {
        asideElement.classList.add('mobile-visible')
        releaseMobileMenuScroll()
        releaseBodyScrollLock = lockBodyScroll()
      } else {
        asideElement.classList.remove('mobile-visible')
        releaseMobileMenuScroll()
      }
    } else if (!isMenuVisible.value) {
      releaseMobileMenuScroll()
    }
  }

  // Close mobile menu
  const closeMobileMenu = () => {
    if (isMenuVisible.value) {
      isMenuVisible.value = false
      emitter.emit('toggleMobileMenu', false)
      
      const asideElement = document.querySelector('.gva-aside')
      if (asideElement) {
        asideElement.classList.remove('mobile-visible')
      }
    }
    releaseMobileMenuScroll()
  }

  // Apply responsive table layout
  const makeTableResponsive = (tableRef) => {
    if (!tableRef || !tableRef.value) return
    
    if (isMobile.value) {
      tableRef.value.$el.classList.add('mobile-friendly-table')
    } else {
      tableRef.value.$el.classList.remove('mobile-friendly-table')
    }
  }

  onMounted(() => {
    checkMobile()
    disposeResize = bindWindowEvent(window, 'resize', checkMobile)
    
    // Close menu on route change
    disposeRouteChange = bindEmitterHandler(emitter, 'routeChange', closeMobileMenu)
  })

  onUnmounted(() => {
    closeMobileMenu()
    disposeResize?.()
    disposeResize = null
    disposeRouteChange?.()
    disposeRouteChange = null
  })

  return {
    isMobile,
    isMenuVisible,
    toggleMobileMenu,
    closeMobileMenu,
    makeTableResponsive
  }
}

/**
 * Apply responsive design to Element Plus tables
 * @param {Object} tableRef - Reference to the table component
 * @param {Array} columns - Table columns configuration
 */
export function applyResponsiveTable(tableRef, columns) {
  if (!tableRef || !tableRef.value) return
  
  const screenWidth = document.body.clientWidth
  const isMobile = screenWidth < 768
  
  if (isMobile) {
    // Add responsive class
    tableRef.value.$el.classList.add('mobile-friendly-table')
    
    // Add data-label attributes for mobile card view
    const headerCells = tableRef.value.$el.querySelectorAll('th .cell')
    headerCells.forEach((headerCell, index) => {
      const label = headerCell.textContent.trim()
      
      // Find corresponding body cells in the same column
      const columnCells = tableRef.value.$el.querySelectorAll(`tbody td:nth-child(${index + 1}) .cell`)
      columnCells.forEach(cell => {
        cell.setAttribute('data-label', label)
      })
    })
  }
}

// Utility function to adapt form layouts for mobile
export function adaptFormForMobile(formRef) {
  if (!formRef || !formRef.value) return
  
  const screenWidth = document.body.clientWidth
  const isMobile = screenWidth < 768
  
  if (isMobile) {
    // Change label position to top for mobile
    formRef.value.labelPosition = 'top'
  } else {
    // Restore default label position
    formRef.value.labelPosition = 'right'
  }
}
