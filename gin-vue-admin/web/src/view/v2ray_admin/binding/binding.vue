<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo" class="demo-form-inline" @keyup.enter="onSubmit">
      <el-form-item label="创建时间">
      <el-date-picker v-model="searchInfo.startCreatedAt" type="date" placeholder="开始时间"></el-date-picker>
       —
      <el-date-picker v-model="searchInfo.endCreatedAt" type="date" placeholder="结束时间"></el-date-picker>
      </el-form-item>
        <el-form-item label="服务器ip">
         <el-select v-model="searchInfo.server_id" clearable filterable style="width:194px">
           <el-option v-for="item in srvs" :key="item.ID" :value="item.ID" :label="item.ip" />
         </el-select>
         <el-icon class="auto-icon" @click="getSrvs"><refresh /></el-icon>
        </el-form-item>
        <el-form-item label="用户名">
         <el-select v-model="searchInfo.user_id" clearable filterable style="width:194px">
           <el-option v-for="item in users" :key="item.ID" :value="item.ID" :label="item.nickName" />
         </el-select>
          <el-icon class="auto-icon" @click="getUsers"><refresh /></el-icon>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSubmit">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="gva-table-box">
        <div class="gva-btn-list">
            <el-button type="primary" icon="plus" @click="openDialog">新增</el-button>
            <el-popover v-model:visible="deleteVisible" placement="top" width="160">
            <p>确定要删除吗？</p>
            <div style="text-align: right; margin-top: 8px;">
                <el-button type="primary" link @click="deleteVisible = false">取消</el-button>
                <el-button type="primary" @click="onDelete">确定</el-button>
            </div>
            <template #reference>
                <el-button icon="delete" style="margin-left: 10px;" :disabled="!multipleSelection.length" @click="deleteVisible = true">删除</el-button>
            </template>
            </el-popover>
        </div>
        <el-table
        ref="multipleTable"
        style="width: 100%"
        tooltip-effect="dark"
        :data="tableData"
        row-key="ID"
        @selection-change="handleSelectionChange"
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
        <el-table-column align="left" label="是否限流" width="100">
            <template #default="scope">{{ scope.row.is_limited ? '是' : '否' }}</template>
        </el-table-column>
        <el-table-column align="left" label="按钮组">
            <template #default="scope">
            <el-button type="primary" link icon="share" class="table-button" @click="shareBindingFunc(scope.row)">分享</el-button>
            <el-button v-if="scope.row.is_limited" type="primary" link icon="unlock" class="table-button" @click="removeLimitedFunc(scope.row)">解除限流</el-button>
            <el-button type="primary" link icon="delete" @click="deleteRow(scope.row)">删除</el-button>
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
      <div class="share-container">
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
              <el-button 
                type="primary" 
                class="copy-btn" 
                :data-clipboard-text="shareInfo.share1_link"
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
              <el-button 
                type="primary" 
                class="copy-btn" 
                :data-clipboard-text="shareInfo.share2_link"
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

        <!-- 使用说明 -->
        <div class="usage-tips">
          <el-alert
            title="使用说明"
            type="info"
            :closable="false"
            show-icon
          >
            <template #default>
              <div class="tips-content">
                <p><strong>扫码导入：</strong>使用对应客户端扫描二维码即可自动导入配置</p>
                <p><strong>链接导入：</strong>复制配置链接到客户端中手动导入</p>
                <p><strong>客户端推荐：</strong>iOS推荐Shadowrocket，Android推荐V2rayNG</p>
              </div>
            </template>
          </el-alert>
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
    <el-dialog v-model="dialogFormVisible" :before-close="closeDialog" title="弹窗操作">
      <el-form :model="formData" label-position="right" ref="elFormRef" :rules="rule" label-width="80px">
        <el-form-item label="服务器ip:">
          <el-select v-model="formData.server_id" clearable filterable style="width:194px">
            <el-option v-for="item in srvs" :key="item.ID" :value="item.ID" :label="item.ip" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户名:">
          <el-select v-model="formData.user_id" clearable filterable style="width:194px">
            <el-option v-for="item in users" :key="item.ID" :value="item.ID" :label="item.nickName" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="closeDialog">取 消</el-button>
          <el-button type="primary" @click="enterDialog">确 定</el-button>
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
  createBinding,
  deleteBinding,
  deleteBindingByIds,
  updateBinding,
  findBinding,
  getBindingList,
  shareBinding,
  removeLimited
} from '@/api/binding'
import { getAllServerApi } from '@/api/server'
import { getAllUserApi } from '@/api/user'
import  QRCode  from 'qrcode'
import ClipboardJS from 'clipboard';

// 全量引入格式化工具 请按需保留
import { getDictFunc, formatDate, formatBoolean, filterDict } from '@/utils/format'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive, onMounted } from 'vue'
import { 
  Monitor, 
  Cellphone, 
  Picture, 
  DocumentCopy, 
  Download, 
  Connection, 
  Close 
} from '@element-plus/icons-vue'

const clipboard = new ClipboardJS('.btn');

// 自动化生成的字典（可能为空）以及字段
const formData = ref({
        server_id: '',
        user_id: '',
        alter_id: 64,
        })

// 验证规则
const rule = reactive({})

const elFormRef = ref()


// =========== 表格控制部分 ===========
const page = ref(1)
const total = ref(0)
const pageSize = ref(10)
const tableData = ref([])
const searchInfo = ref({})
const shareInfo = ref({
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
  const table = await getBindingList({ page: page.value, pageSize: pageSize.value, ...searchInfo.value })
  if (table.code === 0) {
    tableData.value = table.data.list
    total.value = table.data.total
    page.value = table.data.page
    pageSize.value = table.data.pageSize
  }
}

getTableData()

// ============== 表格控制部分结束 ===============

// 获取需要的字典 可能为空 按需保留
const setOptions = async () =>{
}

// 获取需要的字典 可能为空 按需保留
setOptions()


// 多选数据
const multipleSelection = ref([])
// 多选
const handleSelectionChange = (val) => {
    multipleSelection.value = val
}

// 删除行
const deleteRow = (row) => {
    ElMessageBox.confirm('确定要删除吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
    }).then(() => {
            deleteBindingFunc(row)
        })
    }


// 批量删除控制标记
const deleteVisible = ref(false)

// 多选删除
const onDelete = async() => {
      const ids = []
      if (multipleSelection.value.length === 0) {
        ElMessage({
          type: 'warning',
          message: '请选择要删除的数据'
        })
        return
      }
      multipleSelection.value &&
        multipleSelection.value.map(item => {
          ids.push(item.ID)
        })
      const res = await deleteBindingByIds({ ids })
      if (res.code === 0) {
        ElMessage({
          type: 'success',
          message: '删除成功'
        })
        if (tableData.value.length === ids.length && page.value > 1) {
          page.value--
        }
        deleteVisible.value = false
        getTableData()
      }
    }

// 行为控制标记（弹窗内部需要增还是改）
const type = ref('')

// 分享标记
const shareFormVisible = ref(false)

const shareBindingFunc = async(row) => {
    const res = await shareBinding({ ID: row.ID })
    if (res.code === 0) {
      console.log("aaa:", res.data.share1)
        shareInfo.value.share1_link = res.data.share1
        shareInfo.value.share2_link = res.data.share2
        QRCode.toDataURL(res.data.share1)
            .then((url) => {
                shareInfo.value.share1 = url
            })
        QRCode.toDataURL(res.data.share2)
            .then((url) => {
                shareInfo.value.share2 = url
            })
        shareFormVisible.value = true
    }
}

const closeShareDialog = () => {
  shareFormVisible.value = false
}

// 处理复制操作
const handleCopy = (configType) => {
  ElMessage({
    type: 'success',
    message: `${configType === 'config1' ? 'Shadowrocket' : 'V2rayN'} 配置已复制到剪贴板`,
    duration: 2000
  })
}

// 下载二维码
const downloadQR = (dataUrl, filename) => {
  const link = document.createElement('a')
  link.download = `${filename}.png`
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

// 删除行
const deleteBindingFunc = async (row) => {
    const res = await deleteBinding({ ID: row.ID })
    if (res.code === 0) {
        ElMessage({
                type: 'success',
                message: '删除成功'
            })
            if (tableData.value.length === 1 && page.value > 1) {
            page.value--
        }
        getTableData()
    }
}

// 解除限流
const removeLimitedFunc = async(row) => {
  const res = await removeLimited({ ID: row.ID })
  if (res.code === 0) {
    ElMessage({
      type: 'success',
      message: '解除限流成功'
    })
    getTableData()
  } else {
    ElMessage({
      type: 'error',
      message: '解除限流失败: ' + res.msg
    })
  }
}

// 弹窗控制标记
const dialogFormVisible = ref(false)

// 打开弹窗
const openDialog = () => {
    type.value = 'create'
    dialogFormVisible.value = true
}

// 关闭弹窗
const closeDialog = () => {
    dialogFormVisible.value = false
    formData.value = {
        server_id: '',
        user_id: '',
        }
}
// 弹窗确定
const enterDialog = async () => {
     elFormRef.value?.validate( async (valid) => {
             if (!valid) return
              let res
              switch (type.value) {
                case 'create':
                  res = await createBinding(formData.value)
                  break
                case 'update':
                  res = await updateBinding(formData.value)
                  break
                default:
                  res = await createBinding(formData.value)
                  break
              }
              if (res.code === 0) {
                ElMessage({
                  type: 'success',
                  message: '创建/更改成功'
                })
                closeDialog()
                getTableData()
              }
      })
}

const srvs = ref([])
const getSrvs = async() => {
  const res = await getAllServerApi()
  if (res.code === 0) {
    srvs.value = res.data.srvs
  }
}

const users = ref([])
const getUsers = async() => {
  const res = await getAllUserApi()
  if (res.code === 0) {
    users.value = res.data.users
  }
}

const init = () => {
  getSrvs()
  getUsers()
}
init()

</script>

<style scoped>
/* 分享弹窗样式 */
:deep(.share-dialog) {
  .el-dialog {
    margin: 0 auto;
    border-radius: 8px;
    overflow: hidden;
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
  }
  
  .el-dialog__footer {
    padding: 0;
    margin: 0;
  }
}

.share-container {
  padding: 24px;
  background: #f8f9ff;
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
  gap: 24px;
  align-items: flex-start;
  width: 100%;
}

.qr-container {
  flex-shrink: 0;
}

.qr-wrapper {
  position: relative;
  width: 200px;
  height: 200px;
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
  flex-direction: column;
  gap: 12px;
  min-width: 0; /* 防止flex子元素溢出 */
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
}

.copy-btn {
  background: linear-gradient(135deg, #409eff 0%, #67c23a 100%);
  border: none;
  color: white;
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

/* 原有的移动端优化样式 */
@media screen and (max-width: 768px) {
  .gva-search-box {
    .el-form-item {
      /* 服务器ip和用户名的选择器容器 */
      &:has(.el-select) {
        position: relative;
        
        .el-form-item__content {
          position: relative;
          padding-right: 2rem;
        }
        
        .auto-icon {
          position: absolute;
          right: 0;
          top: 50%;
          transform: translateY(-50%);
          font-size: 1.25rem;
          color: #409eff;
          cursor: pointer;
          z-index: 10;
        }
      }
    }
  }
  
  /* 表格横向滚动优化 - 确保按钮固定 */
  .gva-table-box {
    position: relative;
    
    /* 按钮区域固定，不随表格滚动 */
    .gva-btn-list {
      position: relative;
      z-index: 2;
      background: #fff;
      padding: 0.5rem;
      border-bottom: 1px solid #ebeef5;
    }
    
    /* 表格容器可滚动 */
    .el-table {
      overflow-x: auto;
      -webkit-overflow-scrolling: touch;
    }
    
    /* 添加渐变提示 */
    &::before {
      content: '';
      position: absolute;
      right: 0;
      top: 4rem; /* 调整位置，避免覆盖按钮 */
      bottom: 3rem;
      width: 2rem;
      background: linear-gradient(to left, rgba(255,255,255,0.9), transparent);
      pointer-events: none;
      z-index: 1;
    }
  }
}
</style>
