<template>
  <el-row class="app-container">
    <el-row>
      <span style="font-size:14px;margin-right:10px">时间戳</span>
      <el-input v-model="timestamp" style="width: 200px" />
      <el-button type="primary" @click="timestampToDate">转换>></el-button>
      <el-input v-model="date" style="width: 200px" />
      <span style="font-size:14px;margin-left:10px;">北京时间</span>
    </el-row>
  </el-row>
</template>
<script>
import { parseTime } from '@/utils'
export default {

  data() {
    return {
      date: '',
      timestamp: ''
    }
  },
  computed: {

  },
  created() {

  },
  methods: {
    timestampToDate() {
      this.date = parseTime(this.timestamp)
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
