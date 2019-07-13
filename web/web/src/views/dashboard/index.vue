<template>
  <div class="dashboard-container">
    <el-row>
      <div id="data-chart-view" style="height:500px;" />
    </el-row>
    <div v-for="(item, index) in items" :key="index">
      {{ item.text }}
    </div>
  </div>
</template>

<script>
import echarts from 'echarts'
export default {
  name: 'Dashboard',
  data() {
    return {
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
          data: ['排队量', '已接入会话量', '未接入会话量', '会话总量']
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
            name: '排队量',
            type: 'line',
            color: '#916DE3',
            data: [5682, 9186, 3983, 3368, 3447, 4421, 5709, 4469, 4068, 6647, 3560, 8183, 9977, 428, 9915, 5005, 3306, 8897, 9005, 5149, 6116, 5351, 2278, 5621]
          },
          {
            name: '已接入会话量',
            type: 'line',
            color: '#87D3AC',
            data: [8538, 7857, 5171, 8092, 1461, 2495, 7826, 6830, 5665, 9875, 4546, 826, 8378, 9471, 8789, 2772, 2598, 6524, 4125, 1522, 4108, 8714, 3269, 126]
          },
          {
            name: '未接入会话量',
            type: 'line',
            color: '#A9ABAD',
            data: [8636, 5925, 9122, 6340, 6989, 8269, 2090, 8176, 5723, 6039, 1555, 6236, 4401, 6399, 1902, 8504, 5127, 8881, 9917, 771, 2004, 3195, 4144, 6262]
          },
          {
            name: '会话总量',
            type: 'line',
            color: '#5092E1',
            data: [4505, 7931, 1959, 6327, 2985, 8594, 617, 4592, 458, 2679, 1448, 4652, 277, 7397, 7949, 9299, 3667, 7489, 6163, 5007, 8841, 8339, 6808, 7868]
          }
        ]
      },
      items: {
      }
    }
  },
  created() {
    this.items[1] = {
      text: 'dafjls'
    }
    this.items[9999] = {
      text: 'cvbxcvb'
    }
    this.items[2] = {
      text: 'asd'
    }
    console.log(this.items)
  },
  mounted() {
    this.drawLine()
  },
  methods: {
    drawLine() {
      const charts = echarts.init(document.getElementById('data-chart-view'))
      charts.setOption(this.setOption)
    },
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          alert('submit!')
        } else {
          console.log('error submit!!')
          return false
        }
      })
    },
    resetForm(formName) {
      this.$refs[formName].resetFields()
    }
  }
}
</script>

<style lang="scss" scoped>
.dashboard {
  &-container {
    margin: 30px;
  }
  &-text {
    font-size: 30px;
    line-height: 46px;
  }
}
</style>
