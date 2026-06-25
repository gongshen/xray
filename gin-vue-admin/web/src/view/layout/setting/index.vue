<template>
  <div>
    <el-button type="primary" class="drawer-container" :icon="$gvaIcons.Setting" aria-label="打开系统配置" @click="showSettingDrawer" />
    <el-drawer
      v-model="drawer"
      title="系统配置"
      :direction="direction"
      :before-close="handleClose"
    >
      <div class="setting_body">
        <div class="setting_card">
          <div class="setting_content">
            <div class="theme-box" role="group" aria-label="主题模式">
              <button
                class="item"
                type="button"
                :aria-pressed="userStore.mode === 'light'"
                aria-label="切换为简约白主题"
                @click="changeMode('light')"
              >
                <div class="item-top">
                  <el-icon v-if="userStore.mode === 'light'" class="check">
                    <check />
                  </el-icon>
                  <img :src="themeLight" alt="light theme" decoding="async" loading="lazy">
                </div>
                <span class="item-label">
                  简约白
                </span>
              </button>
              <button
                class="item"
                type="button"
                :aria-pressed="userStore.mode === 'dark'"
                aria-label="切换为商务黑主题"
                @click="changeMode('dark')"
              >
                <div class="item-top">
                  <el-icon v-if="userStore.mode === 'dark'" class="check">
                    <check />
                  </el-icon>
                  <img :src="themeDark" alt="dark theme" decoding="async" loading="lazy">
                </div>
                <span class="item-label">
                  商务黑
                </span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </el-drawer>

  </div>
</template>

<script>
export default {
  name: 'Setting',
}
</script>

<script setup>
import { ref } from 'vue'
import { useUserStore } from '@/pinia/modules/user'
import themeLight from '@/assets/theme-light.svg'
import themeDark from '@/assets/theme-dark.svg'
const drawer = ref(false)
const direction = ref('rtl')

const userStore = useUserStore()

const handleClose = () => {
  drawer.value = false
}
const showSettingDrawer = () => {
  drawer.value = true
}
const changeMode = (e) => {
  if (e === null) {
    userStore.changeSideMode('dark')
    return
  }
  userStore.changeSideMode(e)
}

</script>

<style lang="scss" scoped>
.drawer-container {
  transition: right 0.2s ease;
  &:hover{
    right: 0
  }
  position: fixed;
  right: -20px;
  bottom: 15%;
  height: 40px;
  width: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 999;
  color: #fff;
  border-radius: 4px 0 0 4px;
  cursor: pointer;
  -webkit-box-shadow: inset 0 0 6px rgba(0 ,0 ,0, 10%);
}
.setting_body{
  padding: 20px;
  .setting_card{
    margin-bottom: 20px;
  }
  .setting_content{
    margin-top: 20px;
    display: flex;
    flex-direction: column;
    >.theme-box{
     display: flex;
    }
    >.color-box{
      div{
        display: flex;
        flex-direction: column;
      }
    }
    .item{
      border: 0;
      background: transparent;
      color: inherit;
      padding: 0;
      cursor: pointer;
      display: flex;
      align-items: center;
      justify-content: center;
      flex-direction: column;
      margin-right: 20px;
      .item-top{
        position: relative;
      }
      .check{
        position: absolute;
        font-size: 20px;
        color: #00afff;
        right:10px;
        bottom: 10px;
      }
      .item-label{
        text-align: center;
        font-size: 12px;
      }
    }
  }
}

</style>
