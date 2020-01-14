<template>
  <el-row class="app-container">
    <el-row>
      <span style="display:inline-block;width:60px;font-size:14px;margin-right:10px">时间戳</span>
      <el-input v-model="timeExchange.timestamp" style="width: 200px" :placeholder="timeExchange.placeholder" />
      <el-button type="primary" @click="timestampToDate">转换>></el-button>
      <el-input v-model="timeExchange.date" style="width: 200px" />
    </el-row>
    <el-row style="margin-top: 10px">
      <span style="display:inline-block;width:60px;font-size:14px;margin-right:10px">时间</span>
      <el-input v-model="dateExchange.date" style="width: 200px" :placeholder="dateExchange.placeholder" />
      <el-button type="primary" @click="dateToTimestamp">转换>></el-button>
      <el-input v-model="dateExchange.timestamp" style="width: 200px" />
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
      timeExchange: {
        date: '',
        timestamp: '',
        timer: null,
        placeholder: Date.parse(new Date()) / 1000
      },
      dateExchange: {
        date: '',
        timestamp: '',
        timer: null,
        placeholder: parseTime(new Date())
      },
      bytes: '',
      bytesUnit: 1,
      humanSize: ''
    }
  },
  computed: {

  },
  created() {
    this.timeExchange.timer = setInterval(() => {
      this.timeExchange.placeholder = Date.parse(new Date()) / 1000
    }, 1000)
    this.dateExchange.timer = setInterval(() => {
      this.dateExchange.placeholder = parseTime(new Date())
    }, 1000)
  },
  beforeDestroy() {
    clearTimeout(this.timeExchange.timer)
    clearTimeout(this.dateExchange.timer)
  },
  methods: {
    timestampToDate() {
      this.timeExchange.date = parseTime(this.timeExchange.timestamp)
    },
    dateToTimestamp() {
      this.dateExchange.timestamp = Date.parse(new Date(this.dateExchange.date)) / 1000
    },
    bytesToHumanSize() {
      this.humanSize = humanSize(this.bytes * this.bytesUnit)
    }
  }
}
</script>
<style lang="scss" scoped>
</style>
