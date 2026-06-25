<template>
  <div>
    <el-row :gutter="10">
      <el-col :span="12">
        <el-card>
          <template #header>
            <el-divider>Gin-Vue-Admin</el-divider>
          </template>
          <div>
            <el-row>
              <el-col :span="8" :offset="8">
                <img
                  class="org-img dom-center"
                  :src="logo"
                  alt="gin-vue-admin" decoding="async" loading="lazy"
                >
              </el-col>
            </el-row>
            <el-row :gutter="10">
              <el-col
                v-for="item in badges"
                :key="item.label"
                :span="8"
              >
                <div class="local-badge dom-center">
                  <span class="local-badge-label">{{ item.label }}</span>
                  <span class="local-badge-value">{{ item.value }}</span>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
        <el-card style="margin-top: 20px">
          <div>
            <el-row>
              <el-col :span="8" :offset="8">
                <img
                  class="org-img dom-center"
                  :src="flippedAurora"
                  alt="flipped-aurora" decoding="async" loading="lazy"
                >
              </el-col>
            </el-row>
            <el-row style="margin-left: 40px" :gutter="20">
              <el-col v-for="item in members" :key="item.name" :span="8">
                <div class="member-item">
                  <img class="avatar-img" :src="item.avatar" :alt="item.name" decoding="async" loading="lazy">
                  <span class="author-name">{{ item.name }}</span>
                </div>
              </el-col>
            </el-row>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div>更新记录</div>
          </template>
          <div>
            <el-timeline>
              <el-timeline-item
                v-for="item in dataTimeline"
                :key="item.from || item.title"
                :timestamp="item.from"
                placement="top"
              >
                <el-card>
                  <h4>{{ item.title }}</h4>
                  <p>{{ item.message }}</p>
                </el-card>
              </el-timeline-item>
            </el-timeline>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script>
export default {
  name: 'About',
}
</script>

<script setup>
import logo from '@/assets/logo.png'
import flippedAurora from '@/assets/flipped-aurora.png'
import docs from '@/assets/docs.png'
import kefu from '@/assets/kefu.png'

const badges = [
  { label: 'Resource', value: 'Local' },
  { label: 'Network', value: 'None' },
  { label: 'Build', value: 'Vite' },
]

const members = [
  { name: 'Admin', avatar: logo },
  { name: 'Docs', avatar: docs },
  { name: 'Support', avatar: kefu },
]

const dataTimeline = [
  {
    from: 'Local',
    title: '本地资源',
    message: '关于页已移除 GitHub badge、成员头像和提交记录接口请求。',
  },
  {
    from: 'Local',
    title: '静态展示',
    message: '页面内容使用项目内置图片和静态数据展示。',
  },
  {
    from: 'Local',
    title: '离线可用',
    message: '无公网访问能力时页面也可以正常渲染。',
  },
]
</script>

<style scoped>
.avatar-img {
  float: left;
  height: 40px;
  width: 40px;
  border-radius: 50%;
  -webkit-border-radius: 50%;
  -moz-border-radius: 50%;
  margin-top: 15px;
}

.org-img {
  height: 150px;
  width: 150px;
}

.author-name {
  float: left;
  line-height: 65px !important;
  margin-left: 10px;
  color: darkblue;
  font-family: "Lucida Sans", "Lucida Sans Regular", "Lucida Grande",
    "Lucida Sans Unicode", Geneva, Verdana, sans-serif;
}

.dom-center {
  margin-left: 50%;
  transform: translateX(-50%);
}

.local-badge {
  display: inline-flex;
  height: 22px;
  overflow: hidden;
  border-radius: 3px;
  color: #fff;
  font-size: 12px;
  line-height: 22px;
  white-space: nowrap;
}

.local-badge-label {
  padding: 0 8px;
  background: #555;
}

.local-badge-value {
  padding: 0 8px;
  background: #409eff;
}

.member-item {
  min-height: 70px;
}
</style>
