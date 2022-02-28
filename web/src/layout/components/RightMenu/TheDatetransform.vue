<template>
  <el-row v-show="modelValue === 'time'">
    <el-button
      style="margin-left: 80px"
      type="primary"
      @click="timestamp('now')"
    >
      {{ $t('now') }}
    </el-button>
    <el-button type="primary" @click="timestamp('today')">
      {{ $t('today') }}
    </el-button>
    <el-button type="primary" @click="timestamp('m1d')">
      {{ $t('m1d') }}
    </el-button>
    <el-button type="primary" @click="timestamp('p1d')">
      {{ $t('p1d') }}
    </el-button>
    <el-row style="margin-top: 10px" type="flex" align="middle">
      <span style="width: 70px; font-size: 14px; margin-right: 10px">
        Timestamp
      </span>
      <el-input
        v-model="timeExchange.timestamp"
        style="width: 200px"
        :placeholder="timeExchange.placeholder"
        clearable
        @keyup.enter="timestampToDate"
      />
      <el-button type="primary" @click="timestampToDate">>></el-button>
      <el-input v-model="timeExchange.date" style="width: 200px" />
    </el-row>
    <el-row style="margin-top: 10px" type="flex" align="middle">
      <span style="width: 70px; font-size: 14px; margin-right: 10px">
        Date
      </span>
      <el-input
        v-model="dateExchange.date"
        style="width: 200px"
        :placeholder="dateExchange.placeholder"
        clearable
        @keyup.enter="dateToTimestamp"
      />
      <el-button type="primary" @click="dateToTimestamp">>></el-button>
      <el-input v-model="dateExchange.timestamp" style="width: 200px" />
    </el-row>
  </el-row>
</template>

<script lang="ts" setup>
import { parseTime } from '@/utils'
import { onUnmounted, reactive } from 'vue'

defineProps({
  modelValue: {
    type: String,
    default: '',
  },
})

const timeExchange = reactive({
  date: parseTime(new Date().getTime()),
  timestamp: '',
  timer: setInterval(() => {
    timeExchange.placeholder = String(Math.round(Date.now() / 1000))
  }, 1000),
  placeholder: String(Math.round(Date.now() / 1000)),
})

const dateExchange = reactive({
  date: '',
  timestamp: Math.round(Date.now() / 1000),
  timer: setInterval(() => {
    dateExchange.placeholder = parseTime(new Date().getTime())
  }, 1000),
  placeholder: parseTime(new Date().getTime()),
})

const timestamp = (value: string) => {
  let ts = 0
  switch (value) {
    case 'now':
      ts = Math.round(Date.now() / 1000)
      break
    case 'today':
      ts = Math.round(
        new Date(new Date().setHours(0, 0, 0, 0)).getTime() / 1000
      )
      break
    case 'm1d':
      ts =
        timeExchange.timestamp !== ''
          ? parseInt(timeExchange.timestamp) - 86400
          : Math.round(Date.now() / 1000) - 86400
      break
    case 'p1d':
      ts =
        timeExchange.timestamp !== ''
          ? parseInt(timeExchange.timestamp) + 86400
          : Math.round(Date.now() / 1000) + 86400
      break
    default:
      ts = Math.round(Date.now() / 1000)
  }
  timeExchange.timestamp = String(ts)
  timeExchange.date = parseTime(ts)
}

const timestampToDate = () => {
  timeExchange.date = parseTime(Number(timeExchange.timestamp))
}

const dateToTimestamp = () => {
  dateExchange.timestamp = new Date(dateExchange.date).getTime() / 1000
}

onUnmounted(() => {
  clearTimeout(timeExchange.timer)
  clearTimeout(dateExchange.timer)
})
</script>
