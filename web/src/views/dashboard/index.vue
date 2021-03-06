<template>
  <div class="app-container">
    <el-row>
      <el-row type="flex" justify="end">
        <el-date-picker v-model="date" type="date" value-format="yyyy-MM-dd" placeholder="选择日期" @change="get" />
      </el-row>
      <h1 class="app-title">数据总览</h1>
      <el-row :gutter="20">
        <el-col :span="6">
          <el-row class="data-list">
            <el-row class="data-list-title">提交量</el-row>
            <countTo class="data-list-number" :start-val="0" :end-val="commitNumber" :duration="3000" />
          </el-row>
        </el-col>
        <el-col :span="6">
          <el-row class="data-list">
            <el-row class="data-list-title">发布量</el-row>
            <countTo class="data-list-number" :start-val="0" :end-val="deployNumber" :duration="3000" />
          </el-row>
        </el-col>
        <el-col :span="6">
          <el-row class="data-list">
            <el-row class="data-list-title">回滚量</el-row>
            <countTo class="data-list-number" :start-val="0" :end-val="backNumber" :duration="3000" />
          </el-row>
        </el-col>
        <el-col :span="6">
          <el-row class="data-list">
            <el-row class="data-list-title">失败量</el-row>
            <countTo class="data-list-number" :start-val="0" :end-val="failNumber" :duration="3000" />
          </el-row>
        </el-col>
      </el-row>
      <h1 class="app-title" style="margin-top:30px;">数据趋势</h1>
      <el-row>
        <div id="data-chart-view" style="height:350px;" />
      </el-row>
    </el-row>
  </div>
</template>

<script>
import echarts from 'echarts'
import countTo from 'vue-count-to'
import { parseTime } from '@/utils'
import { get } from '@/api/home'

export default {
  components: { countTo },
  data() {
    return {
      date: parseTime(new Date(), '{y}-{m}-{d}'),
      commitNumber: 0,
      deployNumber: 0,
      failNumber: 0,
      backNumber: 0,
      setOption: {
        tooltip: {
          axisPointer: {
            type: 'shadow',
            shadowStyle: {
              color: 'rgba(204,213,216,0.5)',
              type: 'default'
            }
          },
          trigger: 'axis'
        },
        grid: {
          borderWidth: 0,
          x: 60,
          y: 50,
          x2: 50,
          y2: 100,
          boundaryGap: false
        },
        legend: {
          data: ['提交量', '发布量', '回滚量', '失败量']
        },
        xAxis: [{
          type: 'category',
          axisLine: {
            show: true,
            lineStyle: {
              color: '#e6e6e6',
              width: 1,
              type: 'solid'
            }
          },
          axisLabel: {
            show: true,
            textStyle: {
              color: '#666'
            }
          },
          axisTick: {
            show: false
          },
          splitArea: {
            show: false
          },
          splitLine: {
            show: true,
            lineStyle: {
              color: '#e6e6e6',
              type: 'dotted'
            }
          },
          data: ['0:00', '1:00', '2:00', '3:00', '4:00', '5:00', '6:00', '7:00', '8:00', '9:00', '10:00', '11:00', '12:00', '13:00', '14:00', '15:00', '16:00', '17:00', '18:00', '19:00', '20:00', '21:00', '22:00', '23:00']
        }],
        yAxis: [{
          axisLine: {
            show: true,
            lineStyle: {
              color: '#e6e6e6',
              width: 0,
              type: 'solid'
            }
          },
          splitLine: {
            show: true,
            lineStyle: {
              color: '#e6e6e6',
              width: 1,
              type: 'dotted'
            }
          },
          axisLabel: {
            show: true,
            textStyle: {
              color: '#666'
            }
          },
          type: 'value',
          name: ''
        }],
        series: [
          {
            name: '提交量',
            type: 'line',
            color: '#916DE3',
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
          },
          {
            name: '发布量',
            type: 'line',
            color: '#87D3AC',
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
          },
          {
            name: '失败量',
            type: 'line',
            color: '#A9ABAD',
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
          },
          {
            name: '回滚量',
            type: 'line',
            color: '#5092E1',
            data: [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0]
          }
        ]
      }
    }
  },
  mounted() {
    this.get()
  },
  methods: {
    drawLine() {
      this.charts = echarts.init(document.getElementById('data-chart-view'))
      this.charts.setOption(this.setOption)
    },
    get() {
      get(this.date).then((response) => {
        const charts = response.data.charts
        charts.forEach((element) => {
          this.setOption.series[0].data[element.hour] = element.commitNumber
          this.setOption.series[1].data[element.hour] = element.deployNumber
          this.setOption.series[2].data[element.hour] = element.failNumber
          this.setOption.series[3].data[element.hour] = element.backNumber
          this.commitNumber += element.commitNumber
          this.deployNumber += element.deployNumber
          this.failNumber += element.failNumber
          this.backNumber += element.backNumber
        })
        this.drawLine()
      })
    }
  }
}
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
.data-list {
  border: 1px solid #ccc;
  border-radius: 2px;
  text-align: center;
  background: #fff;
  height: 95px;
  &-title {
    padding: 10px 0;
    font-size: 14px;
    color: #222;
  }
  &-number {
    font-size: 30px;
    color: #5092e1;
  }
}
</style>
