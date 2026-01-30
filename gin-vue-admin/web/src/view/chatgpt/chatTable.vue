<template>
  <div class="chat-container">
    <!-- Warning bar at the top -->
    <warning-bar title="使用GPT-3.5模型，存在一定不稳定性，成功率为50%左右，使用GPT-4可以极大提升成功率，但是费用较高。" />
    
    <!-- SK Input Section -->
    <div v-if="!chatToken" class="sk-container">
      <div class="sk-input-wrapper">
        <el-input v-model="skObj.sk" class="query-ipt" placeholder="请输入您的ChatGpt SK" clearable />
        <el-button type="primary" @click="save">保存</el-button>
      </div>
      <div class="secret">
        <el-empty description="请到gpt网站获取您的sk：https://platform.openai.com/account/api-keys" />
      </div>
    </div>
    
    <!-- Main Chat Interface -->
    <div v-else class="chat-interface">
      <!-- Left Sidebar - Chat History -->
      <div class="chat-sidebar">
        <div class="sidebar-header">
          <el-button type="primary" class="new-chat-btn" @click="startNewChat">
            <el-icon><Plus /></el-icon> 新建对话
          </el-button>
        </div>
        <div class="chat-history">
          <div 
            v-for="(chat, index) in chatHistory" 
            :key="index" 
            class="history-item"
            :class="{ 'active': currentChatIndex === index }"
            @click="selectChat(index)"
          >
            <span class="history-title">{{ chat.title || '新对话' }}</span>
            <el-icon class="delete-icon" @click.stop="deleteChat(index)"><Delete /></el-icon>
          </div>
        </div>
        <div class="sidebar-footer">
          <el-popover placement="top" width="160">
            <p>确定要删除当前SK吗？</p>
            <div style="text-align: right; margin-top: 8px;">
              <el-button type="primary" @click="deleteSK">确定</el-button>
            </div>
            <template #reference>
              <el-button type="danger" class="delete-sk-btn">删除SK</el-button>
            </template>
          </el-popover>
        </div>
      </div>
      
      <!-- Main Chat Area -->
      <div class="chat-main">
        <!-- Chat Messages -->
        <div class="chat-messages" ref="chatMessagesRef">
          <div v-if="tableData.length === 0" class="empty-chat">
            <div class="empty-chat-content">
              <h2>万用表格 AI 助手</h2>
              <p>您可以询问我任何关于数据库的问题，我会帮您查询并生成相应的SQL语句。</p>
            </div>
          </div>
          
          <div v-else class="message-container">
            <div v-for="(item, index) in messages" :key="index" class="message" :class="item.role">
              <div class="message-avatar">
                <el-avatar :icon="item.role === 'user' ? 'UserFilled' : 'Service'" :size="36"></el-avatar>
              </div>
              <div class="message-content">
                <div v-if="item.role === 'assistant' && item.sql" class="sql-block">
                  <div class="sql-header">
                    <span>生成的SQL</span>
                    <el-button type="text" size="small" @click="copySql(item.sql)">复制</el-button>
                  </div>
                  <pre>{{ item.sql }}</pre>
                </div>
                <div v-if="item.role === 'assistant' && item.data && item.data.length">
                  <el-table
                    :data="item.data"
                    style="width: 100%"
                    tooltip-effect="dark"
                    max-height="300px"
                  >
                    <el-table-column
                      v-for="(col, colIndex) in Object.keys(item.data[0])"
                      :key="colIndex"
                      :prop="col"
                      :label="col"
                      min-width="120"
                      show-overflow-tooltip
                    />
                  </el-table>
                </div>
                <div v-else>
                  {{ item.content }}
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Input Area -->
        <div class="chat-input-area">
          <div class="input-container">
            <el-form :model="form" class="chat-form">
              <div class="model-selector">
                <el-select v-model="form.dbname" placeholder="请选择数据库" size="large">
                  <el-option
                    v-for="(item, index) in dbArr"
                    :key="index"
                    :label="item.database"
                    :value="item.database"
                  />
                </el-select>
              </div>
              <div class="input-wrapper">
                <el-input
                  v-model="form.chat"
                  :autosize="{ minRows: 1, maxRows: 4 }"
                  type="textarea"
                  clearable
                  placeholder="输入您的问题..."
                  @keyup.enter.native="handleQueryTable"
                />
                <el-button 
                  type="primary" 
                  class="send-button" 
                  :disabled="!form.chat || !form.dbname" 
                  @click="handleQueryTable"
                >
                  <el-icon><Position /></el-icon>
                </el-button>
              </div>
            </el-form>
          </div>
          <div class="input-footer">
            <p class="disclaimer">万用表格可以帮您查询数据库并生成SQL。请确保您的问题清晰明确。</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import WarningBar from '@/components/warningBar/warningBar.vue'
import { getTableApi,
  createSKApi,
  getSKApi,
  deleteSKApi } from '@/api/chatgpt'
import { ref, reactive, onMounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { Delete, Plus, Position, UserFilled, Service } from '@element-plus/icons-vue'

const chatToken = ref(null)
const skObj = reactive({
  sk: '',
})
const sql = ref("")
const messages = ref([])
const chatMessagesRef = ref(null)
const chatHistory = ref([])
const currentChatIndex = ref(0)

// 获取SK
const getSK = async() => {
  const res = await getSKApi()
  chatToken.value = res.data.ok
}

// 初始化
onMounted(() => {
  getSK()
})

// 保存SK
const save = async() => {
  if (!skObj.sk) {
    ElMessage.warning('请输入有效的SK')
    return
  }
  const res = await createSKApi(skObj)
  if (res.code === 0) {
    ElMessage.success('SK保存成功')
    await getSK()
  }
}

// 删除SK
const deleteSK = async() => {
  const res = await deleteSKApi()
  if (res.code === 0) {
    ElMessage.success('SK删除成功')
    await getSK()
  }
}

const form = ref({
  dbname: '',
  chat: '',
})
const dbArr = ref([])
const tableData = ref([])

// 查询表格
const handleQueryTable = async() => {
  if (!form.value.chat || !form.value.dbname) {
    ElMessage.warning('请输入问题并选择数据库')
    return
  }
  
  // 添加用户消息
  messages.value.push({
    role: 'user',
    content: form.value.chat
  })
  
  // 滚动到底部
  await nextTick()
  scrollToBottom()
  
  // 添加当前对话到历史记录
  if (chatHistory.value.length === 0) {
    chatHistory.value.push({
      title: form.value.chat.substring(0, 20) + (form.value.chat.length > 20 ? '...' : ''),
      messages: [...messages.value]
    })
  } else {
    chatHistory.value[currentChatIndex.value].messages = [...messages.value]
    chatHistory.value[currentChatIndex.value].title = form.value.chat.substring(0, 20) + (form.value.chat.length > 20 ? '...' : '')
  }
  
  // 发送请求
  const res = await getTableApi(form.value)
  if (res.code === 0) {
    tableData.value = res.data.results || []
    sql.value = res.data.sql
    
    // 添加AI回复
    messages.value.push({
      role: 'assistant',
      content: tableData.value.length ? `已为您查询到${tableData.value.length}条结果` : '未查询到相关数据',
      sql: res.data.sql,
      data: tableData.value
    })
    
    // 更新历史记录
    chatHistory.value[currentChatIndex.value].messages = [...messages.value]
    
    // 滚动到底部
    await nextTick()
    scrollToBottom()
  } else {
    // 添加错误消息
    messages.value.push({
      role: 'assistant',
      content: '查询失败: ' + (res.msg || '未知错误'),
    })
  }
  
  // 清空输入框
  form.value.chat = ''
}

// 滚动到底部
const scrollToBottom = () => {
  if (chatMessagesRef.value) {
    chatMessagesRef.value.scrollTop = chatMessagesRef.value.scrollHeight
  }
}

// 复制SQL
const copySql = (sql) => {
  navigator.clipboard.writeText(sql)
    .then(() => {
      ElMessage.success('SQL已复制到剪贴板')
    })
    .catch(() => {
      ElMessage.error('复制失败')
    })
}

// 新建对话
const startNewChat = () => {
  messages.value = []
  tableData.value = []
  form.value.chat = ''
  sql.value = ''
  
  chatHistory.value.push({
    title: '新对话',
    messages: []
  })
  currentChatIndex.value = chatHistory.value.length - 1
}

// 选择对话
const selectChat = (index) => {
  currentChatIndex.value = index
  messages.value = [...chatHistory.value[index].messages]
  tableData.value = messages.value.filter(m => m.role === 'assistant' && m.data).length > 0 
    ? messages.value.filter(m => m.role === 'assistant' && m.data)[0].data || []
    : []
}

// 删除对话
const deleteChat = (index) => {
  if (chatHistory.value.length === 1) {
    ElMessage.warning('至少保留一个对话')
    return
  }
  
  chatHistory.value.splice(index, 1)
  
  if (currentChatIndex.value === index) {
    currentChatIndex.value = 0
    messages.value = [...chatHistory.value[0].messages]
  } else if (currentChatIndex.value > index) {
    currentChatIndex.value--
  }
}

// 监听消息变化，自动滚动到底部
watch(messages, () => {
  nextTick(() => {
    scrollToBottom()
  })
})
</script>

<style scoped lang="scss">
.chat-container {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  width: 100%;
  background-color: #f9f9f9;
}

.sk-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 40px 20px;
  
  .sk-input-wrapper {
    display: flex;
    width: 100%;
    max-width: 600px;
    margin-bottom: 20px;
    
    .query-ipt {
      flex: 1;
      margin-right: 10px;
    }
  }
  
  .secret {
    width: 100%;
    max-width: 600px;
    padding: 30px;
    margin-top: 20px;
    background: #F5F5F5;
    border-radius: 8px;
  }
}

.chat-interface {
  display: flex;
  height: 100%;
  overflow: hidden;
}

.chat-sidebar {
  width: 260px;
  background-color: #202123;
  color: white;
  display: flex;
  flex-direction: column;
  height: 100%;
  
  .sidebar-header {
    padding: 10px;
    border-bottom: 1px solid #4d4d4f;
    
    .new-chat-btn {
      width: 100%;
      background-color: #343541;
      border: 1px solid #565869;
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      
      &:hover {
        background-color: #444654;
      }
    }
  }
  
  .chat-history {
    flex: 1;
    overflow-y: auto;
    padding: 10px;
    
    .history-item {
      padding: 10px;
      margin-bottom: 5px;
      border-radius: 5px;
      cursor: pointer;
      display: flex;
      justify-content: space-between;
      align-items: center;
      
      &:hover {
        background-color: #343541;
      }
      
      &.active {
        background-color: #343541;
      }
      
      .history-title {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        flex: 1;
      }
      
      .delete-icon {
        opacity: 0;
        transition: opacity 0.2s;
      }
      
      &:hover .delete-icon {
        opacity: 1;
      }
    }
  }
  
  .sidebar-footer {
    padding: 10px;
    border-top: 1px solid #4d4d4f;
    
    .delete-sk-btn {
      width: 100%;
      background-color: transparent;
      border: 1px solid #565869;
      color: #ff4d4f;
    }
  }
}

.chat-main {
  flex: 1;
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #ffffff;
}

.chat-messages {
  flex: 1;
  overflow-y: auto;
  padding: 20px;
  
  .empty-chat {
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: center;
    
    .empty-chat-content {
      text-align: center;
      max-width: 600px;
      
      h2 {
        font-size: 2rem;
        margin-bottom: 20px;
        color: #343541;
      }
      
      p {
        color: #6e6e80;
        font-size: 1rem;
        line-height: 1.5;
      }
    }
  }
  
  .message-container {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }
  
  .message {
    display: flex;
    padding: 20px;
    
    &.user {
      background-color: #f7f7f8;
    }
    
    &.assistant {
      background-color: #ffffff;
    }
    
    .message-avatar {
      margin-right: 16px;
      flex-shrink: 0;
    }
    
    .message-content {
      flex: 1;
      
      .sql-block {
        background-color: #f2f2f2;
        border-radius: 6px;
        margin: 10px 0;
        overflow: hidden;
        
        .sql-header {
          background-color: #e6e6e6;
          padding: 8px 12px;
          display: flex;
          justify-content: space-between;
          align-items: center;
          font-weight: bold;
        }
        
        pre {
          padding: 12px;
          margin: 0;
          overflow-x: auto;
          font-family: 'Courier New', Courier, monospace;
        }
      }
    }
  }
}

.chat-input-area {
  padding: 20px;
  border-top: 1px solid #e6e6e6;
  
  .input-container {
    max-width: 800px;
    margin: 0 auto;
    
    .chat-form {
      display: flex;
      flex-direction: column;
      gap: 10px;
      
      .model-selector {
        width: 100%;
      }
      
      .input-wrapper {
        display: flex;
        align-items: flex-end;
        width: 100%;
        border: 1px solid #e6e6e6;
        border-radius: 8px;
        padding: 8px;
        background-color: #ffffff;
        box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
        
        .el-textarea {
          flex: 1;
          margin-right: 10px;
          
          :deep(.el-textarea__inner) {
            border: none;
            padding: 0;
            resize: none;
            
            &:focus {
              box-shadow: none;
            }
          }
        }
        
        .send-button {
          border-radius: 6px;
          padding: 8px 16px;
        }
      }
    }
  }
  
  .input-footer {
    max-width: 800px;
    margin: 10px auto 0;
    text-align: center;
    
    .disclaimer {
      font-size: 12px;
      color: #8e8ea0;
    }
  }
}
</style>
