<template>
  <div class="server-state">
    <!-- 服务器选择器 -->
    <div class="server-selector" role="group" aria-label="服务器状态操作">
      <el-select v-model="selectedServerId" placeholder="选择代理服务器" aria-label="代理服务器" @change="onServerChange" class="server-select">
        <el-option
          v-for="server in serverList"
          :key="server.ID"
          :label="`${server.remark || server.ip} (${server.ip})`"
          :value="server.ID"
        />
      </el-select>
      <el-button type="primary" :icon="Refresh" @click="refreshData" :loading="loading">刷新</el-button>
      <el-button type="primary" :icon="Search" @click="openTrafficAnalysis" :disabled="!currentServer">流量分析</el-button>
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
                  aria-label="磁盘使用率"
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
                aria-label="CPU 使用率"
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
                  aria-label="内存使用率"
                />
              </div>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </template>

    <!-- 无数据提示 -->
    <el-empty v-if="!currentServer && !loading" description="请选择一个代理服务器" />

    <el-dialog v-model="trafficAnalysisVisible" :before-close="closeTrafficAnalysis" title="用户流量分析" :width="isMobile ? '95%' : '1120px'" class="traffic-analysis-dialog">
      <div class="analysis-header">
        <div>
          <div class="analysis-title">{{ currentServer?.remark || currentServer?.ip || '-' }}</div>
          <div class="analysis-subtitle">服务器：{{ currentServer?.ip || '-' }} · Stat 端口：{{ currentServer?.stat_port || 56611 }}</div>
        </div>
        <el-tag :type="isOnline ? 'success' : 'danger'" size="small">{{ isOnline ? '在线' : '离线' }}</el-tag>
      </div>

      <el-form :model="trafficAnalysisForm" class="analysis-form" label-position="top" @keyup.enter="queryTrafficAnalysis">
        <el-form-item label="用户">
          <el-select
            v-model="trafficAnalysisForm.user_tag"
            clearable
            filterable
            :loading="trafficUserLoading"
            placeholder="请选择用户"
            aria-label="流量分析用户"
          >
            <el-option
              v-for="item in users"
              :key="item.ID"
              :value="item.ID"
              :label="formatUserOption(item)"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="日期">
          <el-date-picker
            v-model="trafficAnalysisForm.date"
            type="date"
            value-format="YYYYMMDD"
            format="YYYY-MM-DD"
            clearable
            placeholder="选择日期"
            aria-label="流量分析日期"
          />
        </el-form-item>
        <el-form-item label="开始">
          <el-input v-model="trafficAnalysisForm.start" clearable placeholder="8:10" aria-label="流量分析开始时间" />
        </el-form-item>
        <el-form-item label="结束">
          <el-input v-model="trafficAnalysisForm.end" clearable placeholder="9:00" aria-label="流量分析结束时间" />
        </el-form-item>
        <el-form-item class="analysis-action">
          <el-button type="primary" :icon="Search" :loading="trafficAnalysisLoading" @click="queryTrafficAnalysis">查询</el-button>
        </el-form-item>
      </el-form>

      <div v-if="trafficAnalysisResult" class="analysis-summary" role="status" aria-live="polite" aria-label="流量分析摘要">
        <div class="summary-item">
          <span>时间范围</span>
          <strong>{{ trafficAnalysisResult.start_time }} ~ {{ trafficAnalysisResult.end_time }}</strong>
        </div>
        <div class="summary-item">
          <span>总流量</span>
          <strong>{{ formatBytes(trafficAnalysisTotals.total) }}</strong>
        </div>
        <div class="summary-item">
          <span>下行 / 上行</span>
          <strong>{{ formatBytes(trafficAnalysisTotals.down) }} / {{ formatBytes(trafficAnalysisTotals.up) }}</strong>
        </div>
        <div class="summary-item">
          <span>日志匹配 / 分钟数</span>
          <strong>{{ trafficAnalysisResult.access_log_matched || 0 }} / {{ trafficAnalysisRows.length }}</strong>
        </div>
      </div>

      <el-table
        v-loading="trafficAnalysisLoading"
        :aria-busy="trafficAnalysisLoading"
        aria-label="用户流量分钟明细表"
        :data="trafficAnalysisRows"
        border
        height="420"
        empty-text="暂无数据"
        class="analysis-table"
      >
        <el-table-column label="分钟" prop="minute" width="160" />
        <el-table-column label="采集次数" prop="events" width="90" />
        <el-table-column label="下行" width="110">
          <template #default="scope">{{ formatBytes(scope.row.down) }}</template>
        </el-table-column>
        <el-table-column label="上行" width="110">
          <template #default="scope">{{ formatBytes(scope.row.up) }}</template>
        </el-table-column>
        <el-table-column label="总流量" width="110">
          <template #default="scope">{{ formatBytes(scope.row.total) }}</template>
        </el-table-column>
        <el-table-column label="访问目标" min-width="420">
          <template #default="scope">
            <div v-if="targetNames(scope.row.targets).length" class="target-list">
              <el-tag v-for="target in targetNames(scope.row.targets)" :key="target" size="small" effect="plain">
                {{ target }}
              </el-tag>
              <el-button
                type="primary"
                size="small"
                link
                class="target-classify-button"
                :icon="MagicStick"
                :loading="targetClassificationLoading && targetClassificationMinute === scope.row.minute"
                :disabled="targetClassificationLoading && targetClassificationMinute !== scope.row.minute"
                @click.stop="classifyTrafficTargets(scope.row)"
              >
                分类
              </el-button>
            </div>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="targetClassificationResult" class="classification-result" role="status" aria-live="polite">
        <div class="classification-title">
          访问目标分类
          <span v-if="targetClassificationMinute">· {{ targetClassificationMinute }}</span>
          <span v-if="targetClassificationTargetCount">· {{ targetClassificationTargetCount }} 个目标</span>
        </div>
        <pre>{{ targetClassificationResult }}</pre>
      </div>
    </el-dialog>
  </div>
</template>

<script setup>
import { getAllServerApi, restartVPSApi, analyzeUserTrafficApi, classifyTrafficTargetsApi } from '@/api/server'
import { getAllUserApi } from '@/api/user'
import { onUnmounted, onMounted, ref, computed, reactive } from 'vue'
import { MagicStick, Refresh, RefreshRight, Search } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { progressThresholdColors } from '@/style/designTokens.mjs'
import { bindWindowEvent } from '@/utils/eventLifecycle.mjs'
import {
  calculatePercent,
  formatBytes,
  formatSize,
  formatTime,
  formatUserOption,
  getTrafficAnalysisTotals,
  isServerOnline,
  rowTrafficAnalysisTargets,
  targetNames,
  todayCompact,
  validateTrafficAnalysisQuery,
} from './stateHelpers.mjs'

const timer = ref(null)
let disposeResize = null
const loading = ref(false)
const restartLoading = ref(false)
const serverList = ref([])
const selectedServerId = ref(null)
const isMobile = ref(false)
const offlineThresholdSeconds = 10 * 60
const trafficAnalysisVisible = ref(false)
const trafficAnalysisLoading = ref(false)
const trafficUserLoading = ref(false)
const targetClassificationLoading = ref(false)
const targetClassificationResult = ref('')
const targetClassificationMinute = ref('')
const targetClassificationTargetCount = ref(0)
const trafficAnalysisRows = ref([])
const trafficAnalysisResult = ref(null)
const users = ref([])
const trafficAnalysisForm = reactive({
  user_tag: '',
  date: '',
  start: '',
  end: '',
})

const colors = progressThresholdColors

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
  return calculatePercent(currentServer.value?.disk_used, currentServer.value?.disk_total)
})

const memPercent = computed(() => {
  return calculatePercent(currentServer.value?.mem_used, currentServer.value?.mem_total)
})

const cpuPercent = computed(() => {
  if (!currentServer.value) return 0
  return Math.round(currentServer.value.cpu_percent || 0)
})

const trafficAnalysisTotals = computed(() => {
  return getTrafficAnalysisTotals(trafficAnalysisRows.value)
})

// 判断服务器是否在线 (10分钟内有更新)
const isOnline = computed(() => {
  return isServerOnline(currentServer.value?.sysinfo_at, offlineThresholdSeconds)
})

const openTrafficAnalysis = () => {
  if (!currentServer.value) {
    ElMessage.warning('请先选择服务器')
    return
  }
  if (users.value.length === 0) {
    getUsers()
  }
  trafficAnalysisForm.user_tag = ''
  trafficAnalysisForm.date = todayCompact()
  trafficAnalysisForm.start = ''
  trafficAnalysisForm.end = ''
  trafficAnalysisRows.value = []
  trafficAnalysisResult.value = null
  targetClassificationResult.value = ''
  targetClassificationMinute.value = ''
  targetClassificationTargetCount.value = 0
  trafficAnalysisVisible.value = true
}

const closeTrafficAnalysis = () => {
  trafficAnalysisVisible.value = false
  trafficAnalysisLoading.value = false
  targetClassificationLoading.value = false
  trafficAnalysisRows.value = []
  trafficAnalysisResult.value = null
  targetClassificationResult.value = ''
  targetClassificationMinute.value = ''
  targetClassificationTargetCount.value = 0
}

const queryTrafficAnalysis = async () => {
  const message = validateTrafficAnalysisQuery({
    currentServer: currentServer.value,
    form: trafficAnalysisForm,
  })
  if (message) {
    ElMessage.warning(message)
    return
  }
  trafficAnalysisLoading.value = true
  targetClassificationResult.value = ''
  targetClassificationMinute.value = ''
  targetClassificationTargetCount.value = 0
  try {
    const res = await analyzeUserTrafficApi({
      server_id: currentServer.value.ID,
      user_tag: String(trafficAnalysisForm.user_tag || '').trim(),
      date: String(trafficAnalysisForm.date || '').trim(),
      start: String(trafficAnalysisForm.start || '').trim(),
      end: String(trafficAnalysisForm.end || '').trim(),
    })
    if (res?.code === 0) {
      trafficAnalysisResult.value = res.data.analysis || {}
      trafficAnalysisRows.value = Array.isArray(trafficAnalysisResult.value.rows) ? trafficAnalysisResult.value.rows : []
      if (trafficAnalysisRows.value.length === 0) {
        ElMessage.info('该时间段没有匹配到流量明细')
      }
    }
  } finally {
    trafficAnalysisLoading.value = false
  }
}

const classifyTrafficTargets = async (row) => {
  const targets = rowTrafficAnalysisTargets(row)
  if (targets.length === 0) {
    ElMessage.warning('当前行没有可分类的域名/IP')
    return
  }
  targetClassificationLoading.value = true
  targetClassificationResult.value = ''
  targetClassificationMinute.value = row?.minute || ''
  targetClassificationTargetCount.value = targets.length
  try {
    const res = await classifyTrafficTargetsApi({ targets })
    if (res?.code === 0) {
      targetClassificationResult.value = res.data.classification?.result || ''
      if (!targetClassificationResult.value) {
        ElMessage.warning('分类结果为空')
      }
    }
  } finally {
    targetClassificationLoading.value = false
  }
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

const getUsers = async() => {
  trafficUserLoading.value = true
  try {
    const res = await getAllUserApi()
    if (res.code === 0) {
      users.value = res.data.users || []
    }
  } finally {
    trafficUserLoading.value = false
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
  disposeResize = bindWindowEvent(window, 'resize', checkMobile)
  loadServerList()
  getUsers()
})

// 定时刷新 (每30秒)
timer.value = setInterval(() => {
  loadServerList()
}, 1000 * 30)

onUnmounted(() => {
  clearInterval(timer.value)
  timer.value = null
  disposeResize?.()
  disposeResize = null
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

.analysis-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 16px;
  padding: 14px 16px;
  margin-bottom: 16px;
  background: var(--gva-color-panel-muted-bg);
  border: 1px solid var(--gva-color-border-subtle);
  border-radius: 6px;
}

.analysis-title {
  font-size: 15px;
  font-weight: 600;
  color: var(--gva-color-text-strong);
  word-break: break-word;
}

.analysis-subtitle {
  margin-top: 4px;
  font-size: 13px;
  color: var(--gva-color-text-regular);
  word-break: break-word;
}

.analysis-form {
  display: grid;
  grid-template-columns: minmax(180px, 1.5fr) minmax(160px, 1fr) minmax(110px, 0.7fr) minmax(110px, 0.7fr) auto;
  gap: 12px;
  align-items: end;
  margin-bottom: 12px;
}

.analysis-form :deep(.el-form-item) {
  margin-bottom: 0;
}

.analysis-form :deep(.el-select),
.analysis-form :deep(.el-date-editor.el-input),
.analysis-form :deep(.el-input) {
  width: 100%;
}

.analysis-action {
  min-width: 86px;
}

.analysis-summary {
  display: grid;
  grid-template-columns: minmax(260px, 2fr) repeat(3, minmax(150px, 1fr));
  gap: 10px;
  margin: 12px 0;
}

.summary-item {
  padding: 10px 12px;
  background: var(--gva-color-panel-muted-bg);
  border: 1px solid var(--gva-color-border-subtle);
  border-radius: 6px;
}

.summary-item span {
  display: block;
  margin-bottom: 4px;
  font-size: 12px;
  color: var(--gva-color-text-muted);
}

.summary-item strong {
  font-size: 13px;
  font-weight: 600;
  color: var(--gva-color-text-strong);
}

.target-list {
  display: flex;
  flex-wrap: wrap;
  align-items: center;
  gap: 6px;
  padding: 4px 0;
}

.target-classify-button {
  position: relative;
  margin-left: 8px;
  padding-left: 12px;
}

.target-classify-button::before {
  position: absolute;
  left: 0;
  top: 50%;
  width: 1px;
  height: 14px;
  content: '';
  background: var(--gva-color-divider);
  transform: translateY(-50%);
}

.classification-result {
  margin-top: 12px;
  padding: 12px;
  background: var(--gva-color-panel-muted-bg);
  border: 1px solid var(--gva-color-border-subtle);
  border-radius: 6px;
}

.classification-title {
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--gva-color-text-strong);
}

.classification-result pre {
  margin: 0;
  color: var(--gva-color-text-strong);
  font-family: inherit;
  font-size: 13px;
  line-height: 1.7;
  white-space: pre-wrap;
  word-break: break-word;
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
  border-bottom: 1px solid var(--gva-color-border-muted);
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  color: var(--gva-color-text-regular);
  flex-shrink: 0;
}

.info-value {
  color: var(--gva-color-text-strong);
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
  color: var(--gva-color-text-regular);
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

  .analysis-header {
    align-items: flex-start;
    padding: 0.75rem;
  }

  .analysis-form,
  .analysis-summary {
    grid-template-columns: 1fr;
  }

  .analysis-action {
    min-width: 0;
  }

  .analysis-action .el-button {
    width: 100%;
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
