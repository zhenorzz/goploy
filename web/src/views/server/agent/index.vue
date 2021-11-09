<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-date-picker
          v-model="datetimeRange"
          :shortcuts="shortcuts"
          type="datetimerange"
          range-separator="-"
          :start-placeholder="$t('startDate')"
          :end-placeholder="$t('endDate')"
          @change="datetimeChange"
        />
        <el-button icon="el-icon-refresh" @click="handleRefresh" />
      </el-row>
      <el-row></el-row>
    </el-row>
    <el-row class="chart-container" :gutter="10">
      <el-col
        v-for="(item, name) in chartNameMap"
        :key="name"
        :xs="24"
        :sm="24"
        :lg="12"
        style="margin-bottom: 10px"
      >
        <div
          :ref="name"
          style="height: 288px; border: solid 1px #e6e6e6; padding: 10px 0"
        ></div>
      </el-col>
    </el-row>
  </el-row>
</template>
<script lang="ts">
import * as echarts from 'echarts'
import { ElDatePicker } from 'element-plus'
import { ServerReport } from '@/api/server'
import dayjs, { Dayjs } from 'dayjs'
import { defineComponent } from 'vue'
import { deepClone, parseTime } from '@/utils'

const chartBaseOption = {
  title: {
    text: '',
    textStyle: {
      fontSize: 14,
    },
    padding: [5, 20],
  },
  tooltip: {
    trigger: 'axis',
  },
  xAxis: {
    type: 'time',
  },
  yAxis: {
    type: 'value',
  },
  legend: {
    bottom: 0,
    data: [],
  },
  series: [],
}

export default defineComponent({
  name: 'ServerAgent',
  data() {
    return {
      serverId: Number(this.$route.query.serverId),
      chartNameMap: {
        cpuChart: {
          type: 1,
          title: this.$t('serverPage.cpuUsage'),
          subtitle: '(%)',
        },
        ramChart: {
          type: 2,
          title: this.$t('serverPage.ramUsage'),
          subtitle: '(%)',
        },
        loadavgChart: {
          type: 3,
          title: this.$t('serverPage.loadavg'),
          subtitle: '',
        },
        tcpChart: {
          type: 4,
          title: this.$t('serverPage.tcp'),
          subtitle: '(count)',
        },
        pubNetChart: {
          type: 5,
          title: this.$t('serverPage.pubBandwidth'),
          subtitle: '(bit/s)',
        },
        loNetChart: {
          type: 6,
          title: this.$t('serverPage.loBandwidth'),
          subtitle: '(bit/s)',
        },
        diskUsageChart: {
          type: 7,
          title: this.$t('serverPage.diskUsage'),
          subtitle: '(%)',
        },
        diskIOChart: {
          type: 8,
          title: this.$t('serverPage.diskIO'),
          subtitle: '(count/s)',
        },
      } as Record<string, { type: number; title: string; subtitle: string }>,

      datetimeRange: [] as Dayjs[],
      shortcuts: [
        {
          text: this.$t('lastHour'),
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000)
            picker.emit('pick', [dayjs(start), dayjs(end)])
          },
        },
        {
          text: this.$t('last6Hours'),
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 6)
            picker.emit('pick', [dayjs(start), dayjs(end)])
          },
        },
        {
          text: this.$t('lastDay'),
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24)
            picker.emit('pick', [dayjs(start), dayjs(end)])
          },
        },
        {
          text: this.$t('lastWeek'),
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.emit('pick', [dayjs(start), dayjs(end)])
          },
        },
      ],
    }
  },
  activated() {
    this.serverId = Number(this.$route.query.serverId)
    this.handleRefresh()
  },
  methods: {
    datetimeChange(values: Date[]) {
      for (const key in this.chartNameMap) {
        this.report(key, values)
      }
    },
    handleRefresh() {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000)
      this.datetimeRange = [dayjs(start), dayjs(end)]
      this.datetimeChange(
        this.datetimeRange.map((datetime) => datetime.toDate())
      )
    },
    report(chartName: string, values: Date[]) {
      new ServerReport({
        serverId: this.serverId,
        type: this.chartNameMap[chartName].type,
        datetimeRange: values.map((value) => parseTime(value)).join(','),
      })
        .request()
        .then((response) => {
          echarts.dispose(this.$refs[chartName] as HTMLElement)
          let chart = echarts.init(this.$refs[chartName] as HTMLElement)
          let chartOption = deepClone(chartBaseOption)
          chartOption.title.text = this.chartNameMap[chartName].title
          chartOption.title.subtext = this.chartNameMap[chartName].subtitle
          for (const key in response.data.map) {
            chartOption.legend.data.push(key)
            const series = {
              name: key,
              type: 'line',
              symbol: 'none',
              smooth: true,
              data: response.data.map[key].map((item) => {
                return [new Date(item.reportTime), item.value]
              }),
            }
            chartOption.series.push(series)
          }
          chart.setOption(chartOption)
        })
    },
  },
})
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';

.chart-container {
  width: 100%;
  max-height: calc(100vh - 180px);
  overflow-y: auto;
  @include scrollBar();
}
</style>
