<template>
  <div class="server-state">
    <!-- 服务器选择器 -->
    <div class="server-selector">
      <el-select v-model="selectedServerId" placeholder="选择代理服务器" @change="onServerChange" class="server-select">
        <el-option
          v-for="server in serverList"
          :key="server.ID"
          :label="`${server.remark || server.ip} (${server.ip})`"
          :value="server.ID"
        />
      </el-select>
      <el-button type="primary" :icon="Refresh" @click="refreshData" :loading="loading">刷新</el-button>
      <el-button type="danger" :icon="RefreshRight" @click="restartVPS" :loading="restartLoading">重启服务器</el-button>
    </div>

    <template v-if="currentServer">
      <!-- 服务器基本信息 + 磁盘 -->
      <el-row :gutter="15" class="system_state">
        <el-col :xs="24" :sm="24" :md="12" :lg="12">
          <el-card class="card_item">
            <template #header>
              <div class="card-header">
                <span>服务器信息</span>
                <el-tag :type="isOnline ? 'success' : 'danger'" size="small">
                  {{ isOnline ? '在线' : '离线' }}
                </el-tag>
              </div>
            </template>
            <div class="info-list">
              <div class="info-row">
                <span class="info-label">IP地址:</span>
                <span class="info-value">{{ currentServer.ip }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">备注:</span>
                <span class="info-value">{{ currentServer.remark || '-' }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">端口:</span>
                <span class="info-value">{{ currentServer.port }}</span>
              </div>
              <div class="info-row">
                <span class="info-label">流量重置日:</span>
                <span class="info-value">每月{{ currentServer.reset_date }}号</span>
              </div>
              <div class="info-row">
                <span class="info-label">更新时间:</span>
                <span class="info-value">{{ formatTime(currentServer.sysinfo_at) }}</span>
              </div>
            </div>
          </el-card>
        </el-col>

        <!-- 磁盘信息 -->
        <el-col :xs="24" :sm="24" :md="12" :lg="12">
          <el-card class="card_item">
            <template #header>
              <div>磁盘</div>
            </template>
            <div class="metric-container">
              <div class="metric-info">
                <div class="info-row">
                  <span class="info-label">总量:</span>
                  <span class="info-value">{{ formatSize(currentServer.disk_total) }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">已用:</span>
                  <span class="info-value">{{ formatSize(currentServer.disk_used) }}</span>
                </div>
              </div>
              <div class="metric-progress">
                <el-progress
                  type="dashboard"
                  :percentage="diskPercent"
                  :color="colors"
                  :width="progressWidth"
                />
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>

      <!-- CPU + 内存 -->
      <el-row :gutter="15" class="system_state">
        <el-col :xs="24" :sm="24" :md="12" :lg="12">
          <el-card class="card_item">
            <template #header>
              <div>CPU</div>
            </template>
            <div class="cpu-container">
              <el-progress
                type="dashboard"
                :percentage="cpuPercent"
                :color="colors"
                :width="progressWidth"
              />
              <div class="cpu-text">使用率: {{ cpuPercent }}%</div>
            </div>
          </el-card>
        </el-col>

        <!-- 内存信息 -->
        <el-col :xs="24" :sm="24" :md="12" :lg="12">
          <el-card class="card_item">
            <template #header>
              <div>内存</div>
            </template>
            <div class="metric-container">
              <div class="metric-info">
                <div class="info-row">
                  <span class="info-label">总量:</span>
                  <span class="info-value">{{ formatSize(currentServer.mem_total) }}</span>
                </div>
                <div class="info-row">
                  <span class="info-label">已用:</span>
                  <span class="info-value">{{ formatSize(currentServer.mem_used) }}</span>
                </div>
              </div>
              <div class="metric-progress">
                <el-progress
                  type="dashboard"
                  :percentage="memPercent"
                  :color="colors"
                  :width="progressWidth"
                />
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </template>

    <!-- 无数据提示 -->
    <el-empty v-if="!currentServer && !loading" description="请选择一个代理服务器" />
  </div>
</template>

<script setup>
import { getAllServerApi, restartVPSApi } from '@/api/server'
import { onUnmounted, onMounted, ref, computed } from 'vue'
import { Refresh, RefreshRight } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'

const timer = ref(null)
const loading = ref(false)
const restartLoading = ref(false)
const serverList = ref([])
const selectedServerId = ref(null)
const isMobile = ref(false)
const offlineThresholdSeconds = 10 * 60

const colors = ref([
  { color: '#5cb87a', percentage: 20 },
  { color: '#e6a23c', percentage: 40 },
  { color: '#f56c6c', percentage: 80 }
])

// 响应式进度条宽度
const progressWidth = computed(() => isMobile.value ? 100 : 120)

// 检测移动端
const checkMobile = () => {
  isMobile.value = window.innerWidth <= 768
}

// 当前选中的服务器
const currentServer = computed(() => {
  if (!selectedServerId.value) return null
  return serverList.value.find(s => s.ID === selectedServerId.value)
})

// 计算各项百分比
const diskPercent = computed(() => {
  if (!currentServer.value || !currentServer.value.disk_total) return 0
  return Math.round((currentServer.value.disk_used / currentServer.value.disk_total) * 100)
})

const memPercent = computed(() => {
  if (!currentServer.value || !currentServer.value.mem_total) return 0
  return Math.round((currentServer.value.mem_used / currentServer.value.mem_total) * 100)
})

const cpuPercent = computed(() => {
  if (!currentServer.value) return 0
  return Math.round(currentServer.value.cpu_percent || 0)
})

// 判断服务器是否在线 (10分钟内有更新)
const isOnline = computed(() => {
  if (!currentServer.value || !currentServer.value.sysinfo_at) return false
  const now = Math.floor(Date.now() / 1000)
  return (now - currentServer.value.sysinfo_at) < offlineThresholdSeconds
})

// 格式化时间
const formatTime = (timestamp) => {
  if (!timestamp) return '-'
  const date = new Date(timestamp * 1000)
  return date.toLocaleString('zh-CN')
}

// 格式化字节
const formatBytes = (bytes) => {
  if (!bytes) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let i = 0
  while (bytes >= 1024 && i < units.length - 1) {
    bytes /= 1024
    i++
  }
  return `${bytes.toFixed(2)} ${units[i]}`
}

// 格式化MB为GB (简化显示)
const formatSize = (mb) => {
  if (!mb) return '0 GB'
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(1)} GB`
  }
  return `${mb} MB`
}

// 加载服务器列表
const loadServerList = async () => {
  loading.value = true
  try {
    const { data } = await getAllServerApi()
    serverList.value = data.srvs || []
    // 默认选中第一个服务器
    if (serverList.value.length > 0 && !selectedServerId.value) {
      selectedServerId.value = serverList.value[0].ID
    }
  } finally {
    loading.value = false
  }
}

// 刷新数据
const refreshData = () => {
  loadServerList()
}

// 服务器切换
const onServerChange = () => {
  // 切换服务器时可以做一些额外操作
}

// 重启VPS服务器
const restartVPS = () => {
  if (!currentServer.value) {
    ElMessage.warning('请先选择服务器')
    return
  }
  
  ElMessageBox.confirm('确定要重启该服务器吗？重启可能需要几分钟时间。', '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    restartLoading.value = true
    try {
      const res = await restartVPSApi({ ID: currentServer.value.ID })
      if (res.code === 0) {
        ElMessage.success('重启成功')
      } else {
        ElMessage.error(res.msg || '重启失败')
      }
    } catch (error) {
      ElMessage.error('重启失败: ' + (error.message || '未知错误'))
    } finally {
      restartLoading.value = false
    }
  }).catch(() => {
    // 用户取消
  })
}

// 初始化
onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  loadServerList()
})

// 定时刷新 (每30秒)
timer.value = setInterval(() => {
  loadServerList()
}, 1000 * 30)

onUnmounted(() => {
  clearInterval(timer.value)
  timer.value = null
  window.removeEventListener('resize', checkMobile)
})
</script>

<script>
export default {
  name: 'State',
}
</script>

<style scoped>
.server-state {
  padding: 10px;
}

.server-selector {
  display: flex;
  gap: 10px;
  margin-bottom: 15px;
  align-items: center;
  flex-wrap: wrap;
}

.server-select {
  flex: 1;
  min-width: 200px;
  max-width: 400px;
}

.system_state {
  margin-bottom: 15px;
  display: flex;
}

/* 让同一行的卡片高度一致 */
.system_state > .el-col {
  display: flex;
}

.card_item {
  width: 100%;
  margin-bottom: 15px;
}

/* 卡片内容区域使用flex布局 */
.card_item :deep(.el-card__body) {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 140px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

/* 信息列表样式 */
.info-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 0;
  border-bottom: 1px solid #f0f0f0;
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  color: #666;
  flex-shrink: 0;
}

.info-value {
  color: #333;
  font-weight: 500;
  text-align: right;
  word-break: break-all;
}

/* 指标容器 (磁盘/内存) */
.metric-container {
  display: flex;
  align-items: center;
  justify-content: space-around;
  gap: 20px;
  width: 100%;
}

.metric-info {
  flex: 1;
  max-width: 200px;
}

.metric-progress {
  flex-shrink: 0;
}

/* CPU容器 */
.cpu-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 100%;
}

.cpu-text {
  margin-top: 10px;
  font-size: 14px;
  color: #666;
}

/* 移动端适配 - 使用rem实现动态缩放 */
@media screen and (max-width: 768px) {
  .server-state {
    padding: 0.5rem;
  }

  .server-selector {
    flex-direction: column;
    align-items: stretch;
    gap: 0.625rem;
    margin-bottom: 1rem;
  }

  .server-select {
    max-width: none;
    width: 100%;
  }
  
  /* 按钮容器 */
  .server-selector {
    .el-button {
      width: 100%;
    }
  }

  .system_state {
    display: block;
    margin-bottom: 1rem;
  }

  .system_state > .el-col {
    display: block;
  }

  .card_item {
    margin-bottom: 1rem;
  }

  .card_item :deep(.el-card__body) {
    min-height: auto;
    padding: 0.75rem;
  }

  .metric-container {
    flex-direction: column;
    text-align: center;
    gap: 1rem;
  }

  .metric-info {
    width: 100%;
    max-width: none;
    margin-bottom: 0.75rem;
  }

  .info-row {
    font-size: 0.875rem;
    padding: 0.25rem 0;
  }

  .info-label {
    font-size: 0.8125rem;
  }

  .info-value {
    font-size: 0.8125rem;
  }

  .cpu-text {
    font-size: 0.875rem;
    margin-top: 0.625rem;
  }
}
</style>
