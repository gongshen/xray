<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" role="search" aria-label="用户绑定筛选" @keyup.enter="onSubmit">
      <el-form-item label="创建时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间" aria-label="绑定开始日期"></el-date-picker>
      <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束时间" aria-label="绑定结束日期"></el-date-picker>
      </el-form-item>
      <el-form-item role="group" aria-label="绑定查询操作">
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
      </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        v-loading="tableLoading"
        >
        <el-table-column type="selection" width="55" />
        <el-table-column align="left" label="用户名" width="120">
            <template #default="scope">{{ scope.row.user.nickName }}</template>
        </el-table-column>
        <el-table-column align="left" label="服务器ip" width="140">
            <template #default="scope">{{ scope.row.server.ip }}</template>
        </el-table-column>
        <el-table-column align="left" label="日期" width="180">
            <template #default="scope">{{ formatDate(scope.row.CreatedAt) }}</template>
        </el-table-column>
          <el-table-column align="left" label="按钮组">
            <template #default="scope">
            <el-button type="primary" link icon="share" class="table-button" @click="shareBindingFunc(scope.row)">分享</el-button>
            </template>
        </el-table-column>
        </el-table>
        <div class="gva-pagination">
            <el-pagination
            layout="total, sizes, prev, pager, next, jumper"
            :current-page="page"
            :page-size="pageSize"
            :page-sizes="[10, 30, 50, 100]"
            :total="total"
            @current-change="handleCurrentChange"
            @size-change="handleSizeChange"
            />
        </div>
    </div>
    <el-dialog 
      v-model="shareFormVisible" 
      :before-close="closeShareDialog" 
      title="分享配置"
      width="90%"
      :style="{ maxWidth: '600px' }"
      class="share-dialog"
      center
    >
      <div v-loading="shareLoading" class="share-container">
        <!-- 第一个配置 -->
        <div class="config-section">
          <div class="config-header">
            <div class="config-title">
              <el-icon class="title-icon"><Monitor /></el-icon>
              <span>Shadowrocket / Qv2ray / V2rayXS</span>
            </div>
            <el-tag type="primary" size="small">通用配置</el-tag>
          </div>
          
          <div class="config-content">
            <div class="qr-container">
              <div class="qr-wrapper">
                <img :src="shareInfo.share1" alt="配置二维码" class="qr-image"/>
                <div class="qr-overlay">
                  <el-icon class="qr-icon"><Picture /></el-icon>
                </div>
              </div>
            </div>
            
            <div class="config-actions">
              <div class="button-group">
                <el-button 
                  type="primary" 
                  class="copy-btn" 
                  @click="handleCopy('config1')"
                >
                  <el-icon><DocumentCopy /></el-icon>
                  复制配置链接
                </el-button>
                <el-button 
                  type="success" 
                  plain 
                  class="qr-btn"
                  @click="downloadQR(shareInfo.share1, 'shadowrocket-config')"
                >
                  <el-icon><Download /></el-icon>
                  下载二维码
                </el-button>
              </div>
            </div>
          </div>
        </div>

        <el-divider class="section-divider">
          <el-icon><Connection /></el-icon>
        </el-divider>

        <!-- 第二个配置 -->
        <div class="config-section">
          <div class="config-header">
            <div class="config-title">
              <el-icon class="title-icon"><Cellphone /></el-icon>
              <span>V2rayN / V2rayNG / V2rayXS</span>
            </div>
            <el-tag type="success" size="small">移动端配置</el-tag>
          </div>
          
          <div class="config-content">
            <div class="qr-container">
              <div class="qr-wrapper">
                <img :src="shareInfo.share2" alt="配置二维码" class="qr-image"/>
                <div class="qr-overlay">
                  <el-icon class="qr-icon"><Picture /></el-icon>
                </div>
              </div>
            </div>
            
            <div class="config-actions">
              <div class="button-group">
                <el-button 
                  type="primary" 
                  class="copy-btn" 
                  @click="handleCopy('config2')"
                >
                  <el-icon><DocumentCopy /></el-icon>
                  复制配置链接
                </el-button>
                <el-button 
                  type="success" 
                  plain 
                  class="qr-btn"
                  @click="downloadQR(shareInfo.share2, 'v2ray-config')"
                >
                  <el-icon><Download /></el-icon>
                  下载二维码
                </el-button>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeShareDialog" size="large">
            <el-icon><Close /></el-icon>
            关闭
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script>
export default {
  name: 'Binding'
}
</script>

<script setup>
import {
  getBindingList,
  shareBinding
} from '@/api/v2ray_binding'
import QRCode from 'qrcode'

import { formatDate } from '@/utils/format'
import { ElMessage } from 'element-plus'
import { ref } from 'vue'
import { 
  Monitor, 
  Cellphone, 
  Picture, 
  DocumentCopy, 
  Download, 
  Connection, 
  Close 
} from '@element-plus/icons-vue'
import {
  buildQrDownloadName,
  createShareDialogInfo,
  getCopySuccessMessage,
  getShareLink,
} from '../../v2ray_admin/binding/bindingShare.mjs'

// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const tableLoading = ref(false)
const searchInfo = ref({})
const shareInfo = ref({
  share1: '',
  share1_link: '',
  share2: '',
  share2_link: '',
})
const emptyShareInfo = () => ({
  share1: '',
  share1_link: '',
  share2: '',
  share2_link: '',
})

// 重置
const onReset = () => {
  searchInfo.value = {}
  getTableData()
}

// 搜索
const onSubmit = () => {
  page.value = 1
  pageSize.value = 10
  getTableData()
}

// 分页
const handleSizeChange = (val) => {
  pageSize.value = val
  getTableData()
}

// 修改页面容量
const handleCurrentChange = (val) => {
  page.value = val
  getTableData()
}

// 查询
const getTableData = async() => {
  tableLoading.value = true
  try {
    const table = await getBindingList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
    if (table.code === 0) {
      tableData.value = table.data?.list || []
      total.value = table.data?.total || 0
      page.value = table.data?.page || 1
      pageSize.value = table.data?.pageSize || 10
    } else {
      ElMessage.error(table.msg || '获取数据失败')
      tableData.value = []
      total.value = 0
    }
  } catch (error) {
    ElMessage.error('网络请求失败')
    tableData.value = []
    total.value = 0
  } finally {
    tableLoading.value = false
  }
}

getTableData()

// 分享标记
const shareFormVisible = ref(false)
const shareLoading = ref(false)

const shareBindingFunc = async(row) => {
  shareLoading.value = true
  shareInfo.value = emptyShareInfo()
  shareFormVisible.value = true
  try {
    const res = await shareBinding({ ID: row.ID })
    if (res.code === 0) {
      shareInfo.value = await createShareDialogInfo(res.data, QRCode.toDataURL)
    } else {
      shareFormVisible.value = false
      ElMessage.error(res.msg || '获取分享配置失败')
    }
  } catch (error) {
    shareFormVisible.value = false
    ElMessage.error('获取分享配置失败')
  } finally {
    shareLoading.value = false
  }
}

const closeShareDialog = () => {
  shareFormVisible.value = false
}

// 处理复制操作
const handleCopy = async (configType) => {
  const textToCopy = getShareLink(shareInfo.value, configType)
  const successMessage = getCopySuccessMessage(configType)

  if (!textToCopy) {
    ElMessage({
      type: 'warning',
      message: '暂无可复制的配置',
      duration: 2000
    })
    return
  }

  try {
    await navigator.clipboard.writeText(textToCopy)
    ElMessage({
      type: 'success',
      message: successMessage,
      duration: 2000
    })
  } catch (err) {
    // 如果现代API失败，使用传统方法
    const textArea = document.createElement('textarea')
    textArea.value = textToCopy
    document.body.appendChild(textArea)
    textArea.select()
    try {
      document.execCommand('copy')
      ElMessage({
        type: 'success',
        message: successMessage,
        duration: 2000
      })
    } catch (fallbackErr) {
      ElMessage({
        type: 'error',
        message: '复制失败，请手动复制',
        duration: 2000
      })
    }
    document.body.removeChild(textArea)
  }
}

// 下载二维码
const downloadQR = (dataUrl, filename) => {
  if (!dataUrl) {
    ElMessage({
      type: 'warning',
      message: '暂无可下载的二维码',
      duration: 2000
    })
    return
  }

  const link = document.createElement('a')
  link.download = buildQrDownloadName(filename)
  link.href = dataUrl
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  
  ElMessage({
    type: 'success',
    message: '二维码已下载',
    duration: 2000
  })
}

</script>

<style scoped>
/* 分享弹窗样式 */
:deep(.share-dialog) {
  .el-dialog {
    margin: 0 auto;
    border-radius: 8px;
    overflow: hidden;
    width: 100% !important;
    max-width: 600px !important;
  }
  
  .el-dialog__header {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    color: white;
    padding: 20px 24px;
    border-radius: 8px 8px 0 0;
    margin: 0;
  }
  
  .el-dialog__title {
    font-size: 18px;
    font-weight: 600;
  }
  
  .el-dialog__body {
    padding: 0;
    margin: 0;
    width: 100%;
    box-sizing: border-box;
  }
  
  .el-dialog__footer {
    padding: 0;
    margin: 0;
    width: 100%;
    box-sizing: border-box;
  }
}

.share-container {
  padding: 24px;
  background: #f8f9ff;
  width: 100%;
  box-sizing: border-box;
}

.config-section {
  background: white;
  border-radius: 16px;
  padding: 24px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.config-section:hover {
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.12);
}

.config-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 2px solid #f0f0f0;
}

.config-title {
  display: flex;
  align-items: center;
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.title-icon {
  margin-right: 8px;
  font-size: 20px;
  color: #409eff;
}

.config-content {
  display: flex;
  gap: 20px;
  align-items: center;
  width: 100%;
  box-sizing: border-box;
  min-height: 180px;
}

.qr-container {
  flex-shrink: 0;
  display: flex;
  align-items: center;
  justify-content: center;
}

.qr-wrapper {
  position: relative;
  width: 180px;
  height: 180px;
  border-radius: 16px;
  overflow: hidden;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.qr-wrapper:hover {
  transform: scale(1.05);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
}

.qr-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.qr-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(64, 158, 255, 0.1);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.qr-wrapper:hover .qr-overlay {
  opacity: 1;
}

.qr-icon {
  font-size: 32px;
  color: white;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.3));
}

.config-actions {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 0;
  height: 180px;
}

.button-group {
  display: flex;
  flex-direction: column;
  gap: 16px;
  width: 100%;
  max-width: 280px;
}

.copy-btn, .qr-btn {
  width: 100%;
  height: 48px;
  font-size: 14px;
  font-weight: 500;
  border-radius: 12px;
  transition: all 0.3s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  box-sizing: border-box;
  padding: 12px 16px;
  margin: 0;
  border: 1px solid transparent;
  line-height: 1;
}

.copy-btn {
  background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
  border: 1px solid transparent;
  color: white;
}

.qr-btn {
  background: white;
  border: 1px solid #67c23a;
  color: #67c23a;
}

.copy-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(64, 158, 255, 0.3);
}

.qr-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(103, 194, 58, 0.3);
}

.section-divider {
  margin: 32px 0;
  font-size: 18px;
  color: #909399;
}

.usage-tips {
  margin-top: 24px;
}

.tips-content p {
  margin: 8px 0;
  font-size: 14px;
  line-height: 1.6;
}

.dialog-footer {
  padding: 16px 24px;
  background: #f8f9ff;
  border-radius: 0 0 8px 8px;
  text-align: center;
  margin: 0;
}

.dialog-footer .el-button {
  min-width: 120px;
  height: 40px;
  border-radius: 20px;
  font-weight: 500;
}

/* 移动端优化 */
@media screen and (max-width: 768px) {
  :deep(.share-dialog) {
    width: 95% !important;
    margin: 0 auto;
  }
  
  .share-container {
    padding: 16px;
  }
  
  .config-section {
    padding: 16px;
  }
  
  .config-content {
    flex-direction: column;
    gap: 16px;
    text-align: center;
    align-items: center;
  }
  
  .qr-wrapper {
    width: 160px;
    height: 160px;
  }
  
  .config-header {
    flex-direction: column;
    gap: 8px;
    text-align: center;
  }
  
  .config-title {
    justify-content: center;
  }
  
  .config-actions {
    width: 100%;
    max-width: 300px;
  }
}
</style>
