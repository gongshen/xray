<template>
  <stat-form-base :api="statApi" />
</template>

<script>
export default {
  name: 'Stat'
}
</script>

<script setup>
import * as statApi from '@/api/v2ray_stat'
import StatFormBase from '@/view/v2ray_shared/stat/StatFormBase.vue'
</script>
