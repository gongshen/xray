<template>
  <button
    v-if="href"
    class="warning-bar can-click"
    type="button"
    :aria-label="warningLinkLabel"
    @click="open"
  >
    <el-icon>
      <warning-filled />
    </el-icon>
    <span>
      {{ title }}
    </span>
  </button>
  <div
    v-else
    class="warning-bar"
  >
    <el-icon>
      <warning-filled />
    </el-icon>
    <span>
      {{ title }}
    </span>
  </div>
</template>
<script setup>
import { WarningFilled } from '@element-plus/icons-vue'
import { computed } from 'vue'
import { openExternalUrl } from '@/utils/openExternalUrl'
const prop = defineProps({
  title: {
    type: String,
    default: ''
  },
  href: {
    type: String,
    default: ''
  }
})


const warningLinkLabel = computed(() => `\u6253\u5f00\u63d0\u793a\u94fe\u63a5\uff1a${prop.title}`)

const open = () => {
  openExternalUrl(prop.href)
}
</script>
<style lang="scss" scoped>
.warning-bar{
    border: 0;
    background-color: #FFF5ED;
    color: inherit;
    font-size: 14px;
    font-family: inherit;
    padding: 6px 14px;
    display: flex;
    align-items: center;
    border-radius: 2px;
    .el-icon{
      font-size: 18px;
      color: #ED6A0C;
    }
    margin-bottom: 12px;
    span{
        line-height: 22px;
        color:#F67207;
        margin-left: 8px;
    }
}
.can-click{
  cursor: pointer;
}
</style>
