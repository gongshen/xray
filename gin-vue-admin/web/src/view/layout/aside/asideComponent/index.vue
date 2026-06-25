<template>
  <component
    :is="menuComponent"
    v-if="!routerInfo.hidden"
    :is-collapse="isCollapse"
    :theme="theme"
    :router-info="routerInfo"
  >
    <template v-if="visibleChildren.length">
      <AsideComponent
        v-for="item in visibleChildren"
        :key="item.name"
        :is-collapse="false"
        :router-info="item"
        :theme="theme"
      />
    </template>
  </component>
</template>

<script>
export default {
  name: 'AsideComponent',
}
</script>

<script setup>
import MenuItem from './menuItem.vue'
import AsyncSubmenu from './asyncSubmenu.vue'
import { computed } from 'vue'
const props = defineProps({
  routerInfo: {
    type: Object,
    default: () => ({}),
  },
  isCollapse: {
    default: function() {
      return false
    },
    type: Boolean
  },
  theme: {
    default: function() {
      return {}
    },
    type: Object
  }
})

const visibleChildren = computed(() => props.routerInfo.children?.filter(item => !item.hidden) || [])

const menuComponent = computed(() => {
  return visibleChildren.value.length ? AsyncSubmenu : MenuItem
})

</script>

