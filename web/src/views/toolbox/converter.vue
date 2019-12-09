<template>
  <el-row class="app-container">
    <el-row>
      <span style="display:inline-block;width:60px;font-size:14px;margin-right:10px">时间戳</span>
      <el-input v-model="timestamp" style="width: 200px" />
      <el-button type="primary" @click="timestampToDate">转换>></el-button>
      <el-input v-model="date" style="width: 200px" />
    </el-row>
    <el-row style="margin-top: 10px">
      <span style="display:inline-block;width:60px;font-size:14px;margin-right:10px">字节</span>
      <el-input v-model="bytes" style="width: 200px" />
      <el-select v-model="bytesUnit" style="width: 70px">
        <el-option :value="1" label="B" />
        <el-option :value="1*1024" label="KB" />
        <el-option :value="1024*1024" label="MB" />
      </el-select>
      <el-button type="primary" @click="bytesToHumanSize">转换>></el-button>
      <el-input v-model="humanSize" style="width: 200px" />
    </el-row>
  </el-row>
</template>
<script>
import { parseTime, humanSize } from '@/utils'
export default {

  data() {
    return {
      date: '',
      timestamp: '',
      bytes: '',
      bytesUnit: 1,
      humanSize: ''
    }
  },
  computed: {

  },
  created() {

  },
  methods: {
    timestampToDate() {
      this.date = parseTime(this.timestamp)
    },
    bytesToHumanSize() {
      console.log(this.bytes * this.bytesUnit > 9007199254740992)
      this.humanSize = humanSize(this.bytes * this.bytesUnit)
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
