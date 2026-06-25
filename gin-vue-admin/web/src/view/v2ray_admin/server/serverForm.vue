<template>
  <div>
    <div class="gva-form-box">
      <el-form :model="formData" ref="elFormRef" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="服务器ip:" prop="ip">
          <el-input v-model="formData.ip" :clearable="true" :required="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="端口:" prop="port">
          <el-input v-model.number="formData.port" :clearable="true" :required="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="统计端口:" prop="stat_port">
          <el-input v-model.number="formData.stat_port" :clearable="true" placeholder="默认: 56611" />
        </el-form-item>
        <el-form-item label="备注:" prop="remark">
          <el-input v-model="formData.remark" :clearable="false" placeholder="请输入" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Server'
}
</script>

<script setup>
import {
  createServer,
  updateServer,
  findServer
} from '@/api/server'

import { useRoute, useRouter } from "vue-router"
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, reactive } from 'vue'
import { buildServerPortChangeReminder } from './serverPortReminder.mjs'
const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
            ip: '',
            port: 80,
            remark: '',
            stat_port: 0,
        })
const originalFormData = ref(null)
// 验证规则
const rule = reactive({})
const elFormRef = ref()

// 初始化方法
const init = async () => {
 // 建议通过url传参获取目标数据ID 调用 find方法进行查询数据操作 从而决定本页面是create还是update 以下为id作为url参数示例
    if (route.query.id) {
      const res = await findServer({ ID: route.query.id })
      if (res.code === 0) {
        formData.value = res.data.reserver
        originalFormData.value = { ...res.data.reserver }
        type.value = 'update'
      }
    } else {
      type.value = 'create'
      originalFormData.value = null
    }
}

init()
const confirmPortChangeReminder = async () => {
  const reminder = buildServerPortChangeReminder(originalFormData.value, formData.value)
  if (!reminder) {
    return true
  }

  try {
    await ElMessageBox.confirm(reminder, '端口变更提醒', {
      confirmButtonText: '已同步，继续保存',
      cancelButtonText: '返回修改',
      type: 'warning'
    })
    return true
  } catch {
    return false
  }
}

// 保存按钮
const save = async() => {
      elFormRef.value?.validate( async (valid) => {
          if (!valid) return
            if (!(await confirmPortChangeReminder())) return
             let res
           switch (type.value) {
             case 'create':
               res = await createServer(formData.value)
               break
             case 'update':
               res = await updateServer(formData.value)
               break
             default:
               res = await createServer(formData.value)
               break
           }
            if (res.code === 0) {
              originalFormData.value = { ...formData.value }
              ElMessage({
                type: 'success',
               message: '创建/更改成功'
             })
           }
       })
}

// 返回按钮
const back = () => {
    router.go(-1)
}

</script>

<style>
</style>
