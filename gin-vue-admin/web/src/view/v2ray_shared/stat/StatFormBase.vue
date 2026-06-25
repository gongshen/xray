<template>
  <div>
    <div class="gva-form-box">
      <el-form ref="elFormRef" :model="formData" label-position="right" :rules="rule" label-width="80px">
        <el-form-item label="流量分类:" prop="category">
          <el-input v-model="formData.category" :clearable="false" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="标签:" prop="tag">
          <el-input v-model="formData.tag" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="下行流量:" prop="down">
          <el-input v-model="formData.down" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="上行流量:" prop="up">
          <el-input v-model="formData.up" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item label="总流量:" prop="total">
          <el-input v-model="formData.total" :clearable="true" placeholder="请输入" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="save">保存</el-button>
          <el-button type="primary" @click="back">返回</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup>
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { ref, reactive } from 'vue'

const props = defineProps({
  api: {
    type: Object,
    required: true,
  }
})

const route = useRoute()
const router = useRouter()

const type = ref('')
const formData = ref({
  category: '',
  tag: '',
  down: '',
  up: '',
  total: '',
})

const rule = reactive({
  category: [{
    required: true,
    message: '',
    trigger: ['input', 'blur'],
  }],
})

const elFormRef = ref()

const init = async() => {
  if (route.query.id) {
    const res = await props.api.findStat({ ID: route.query.id })
    if (res.code === 0) {
      formData.value = res.data.restat
      type.value = 'update'
    }
  } else {
    type.value = 'create'
  }
}

init()

const save = async() => {
  elFormRef.value?.validate(async(valid) => {
    if (!valid) return
    let res
    switch (type.value) {
      case 'create':
        res = await props.api.createStat(formData.value)
        break
      case 'update':
        res = await props.api.updateStat(formData.value)
        break
      default:
        res = await props.api.createStat(formData.value)
        break
    }
    if (res.code === 0) {
      ElMessage({
        type: 'success',
        message: '创建/更改成功'
      })
    }
  })
}

const back = () => {
  router.go(-1)
}
</script>
