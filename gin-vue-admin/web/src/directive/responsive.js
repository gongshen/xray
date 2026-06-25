// responsive.js - Custom directives for responsive design

import {
  bindElementResizeHandler,
  unbindElementResizeHandler,
} from './responsiveDirectiveHandlers.mjs'

const RESPONSIVE_TABLE_KEY = 'responsive-table'
const RESPONSIVE_FORM_KEY = 'responsive-form'
const RESPONSIVE_TABLE_LABEL_TIMERS = new WeakMap()

/**
 * Directive to make tables responsive on mobile devices
 * Usage: v-responsive-table
 */
const responsiveTable = {
  mounted(el) {
    makeTableResponsive(el)
    bindElementResizeHandler(el, RESPONSIVE_TABLE_KEY, () => makeTableResponsive(el))
  },
  unmounted(el) {
    clearResponsiveTableTimer(el)
    unbindElementResizeHandler(el, RESPONSIVE_TABLE_KEY)
  }
}

/**
 * Helper function to transform tables for mobile view
 * @param {HTMLElement} el - The table element
 */
function clearResponsiveTableTimer(el) {
  const timer = RESPONSIVE_TABLE_LABEL_TIMERS.get(el)
  if (timer) {
    clearTimeout(timer)
    RESPONSIVE_TABLE_LABEL_TIMERS.delete(el)
  }
}

function makeTableResponsive(el) {
  const screenWidth = document.body.clientWidth
  const isMobile = screenWidth < 768
  
  if (!el || !el.classList.contains('el-table')) return

  clearResponsiveTableTimer(el)

  if (isMobile) {
    // Add responsive class
    el.classList.add('mobile-friendly-table')
    
    // Add data-label attributes for mobile card view
    const labelTimer = setTimeout(() => {
      RESPONSIVE_TABLE_LABEL_TIMERS.delete(el)
      const headerCells = el.querySelectorAll('th .cell')
      const rows = el.querySelectorAll('tbody tr')
      
      rows.forEach(row => {
        const cells = row.querySelectorAll('td .cell')
        cells.forEach((cell, index) => {
          if (headerCells[index]) {
            const label = headerCells[index].textContent.trim()
            cell.setAttribute('data-label', label)
          }
        })
      })
    }, 100)
    RESPONSIVE_TABLE_LABEL_TIMERS.set(el, labelTimer)
  } else {
    el.classList.remove('mobile-friendly-table')
  }
}

/**
 * Directive to adapt form layouts for mobile
 * Usage: v-responsive-form
 */
const responsiveForm = {
  mounted(el, binding) {
    const form = el.querySelector('.el-form')
    if (!form) return
    
    adaptFormForMobile(form, binding.value)
    bindElementResizeHandler(el, RESPONSIVE_FORM_KEY, () => adaptFormForMobile(form, binding.value))
  },
  unmounted(el) {
    unbindElementResizeHandler(el, RESPONSIVE_FORM_KEY)
  }
}

/**
 * Helper function to adapt forms for mobile
 * @param {HTMLElement} form - The form element
 * @param {Object} options - Configuration options
 */
function adaptFormForMobile(form, options = {}) {
  const screenWidth = document.body.clientWidth
  const isMobile = screenWidth < 768
  
  if (isMobile) {
    // Add responsive class
    form.classList.add('mobile-form')
    
    // Adjust label position
    const formItems = form.querySelectorAll('.el-form-item')
    formItems.forEach(item => {
      const label = item.querySelector('.el-form-item__label')
      const content = item.querySelector('.el-form-item__content')
      
      if (label && content) {
        label.style.float = 'none'
        label.style.display = 'block'
        label.style.textAlign = 'left'
        label.style.padding = '0 0 8px'
        content.style.marginLeft = '0'
      }
    })
  } else {
    form.classList.remove('mobile-form')
    
    // Restore default styles if not mobile
    const formItems = form.querySelectorAll('.el-form-item')
    formItems.forEach(item => {
      const label = item.querySelector('.el-form-item__label')
      const content = item.querySelector('.el-form-item__content')
      
      if (label && content) {
        label.style.float = ''
        label.style.display = ''
        label.style.textAlign = ''
        label.style.padding = ''
        content.style.marginLeft = ''
      }
    })
  }
}

export default {
  install(app) {
    app.directive('responsive-table', responsiveTable)
    app.directive('responsive-form', responsiveForm)
  }
}
