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
      <el-row>
        <el-button icon="el-icon-bell" @click="handleAdd">
          {{ $t('serverPage.addMonitor') }}
        </el-button>
        <el-button icon="el-icon-tickets" @click="handleGetList">
          {{ $t('serverPage.monitorList') }}
        </el-button>
      </el-row>
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
    <el-dialog
      v-model="monitorDialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('setting')"
      :close-on-click-modal="false"
    >
      <el-form
        ref="form"
        v-loading="formProps.loading"
        :rules="formRules"
        :model="formData"
        label-width="100px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item :label="$t('serverPage.item')" prop="item">
          <el-select
            v-model="formData.item"
            style="width: 100%"
            filterable
            placeholder="Please start goploy agent"
          >
            <el-option
              v-for="item in itemOptions"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('serverPage.formula')" prop="value">
          <el-row>
            <el-col :span="8">
              <el-select
                v-model="formData.formula"
                style="width: 100%"
                filterable
              >
                <el-option label="avg" value="avg" />
                <el-option label="max" value="max" />
                <el-option label="min" value="min" />
              </el-select>
            </el-col>
            <el-col :span="8">
              <el-select
                v-model="formData.operator"
                style="width: 100%"
                filterable
              >
                <el-option label=">=" value=">=" />
                <el-option label=">" value=">" />
                <el-option label="<=" value="<=" />
                <el-option label="<" value="<" />
                <el-option label="!=" value="!=" />
              </el-select>
            </el-col>
            <el-col :span="8">
              <el-input v-model="formData.value" placeholder="value" />
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item :label="$t('serverPage.cycle')" prop="cycle">
          <el-row>
            <el-col :span="12">
              <el-select
                v-model="formData.groupCycle"
                style="width: 100%"
                filterable
              >
                <el-option label="1 min" :value="1" />
                <el-option label="5 min" :value="5" />
                <el-option label="15 min" :value="15" />
                <el-option label="30 min" :value="30" />
                <el-option label="60 min" :value="60" />
              </el-select>
            </el-col>
            <el-col :span="12">
              <el-select
                v-model="formData.lastCycle"
                style="width: 100%"
                filterable
              >
                <el-option label="Last 1 cycle" :value="1" />
                <el-option label="Last 3 cycles" :value="3" />
                <el-option label="Last 5 cycles" :value="5" />
                <el-option label="Last 10 cycles" :value="10" />
                <el-option label="Last 15 cycles" :value="15" />
                <el-option label="Last 30 cycles" :value="30" />
                <el-option label="Last 60 cycles" :value="60" />
                <el-option label="Last 90 cycles" :value="90" />
                <el-option label="Last 120 cycles" :value="120" />
                <el-option label="Last 180 cycles" :value="180" />
              </el-select>
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item :label="$t('serverPage.vaildPeriod')">
          <el-row>
            <el-col :span="12">
              <el-time-select
                v-model="formData.startTime"
                style="width: 100%"
                placeholder="Start time"
                start="00:00"
                step="01:00"
                end="23:00"
              >
              </el-time-select>
            </el-col>
            <el-col :span="12">
              <el-time-select
                v-model="formData.endTime"
                :min-time="formData.startTime"
                style="width: 100%"
                placeholder="End time"
                start="00:59"
                step="01:00"
                end="23:59"
              >
              </el-time-select>
            </el-col>
          </el-row>
        </el-form-item>
        <el-form-item :label="$t('serverPage.silentCycle')">
          <el-select
            v-model="formData.silentCycle"
            style="width: 100%"
            filterable
          >
            <el-option label="5 min" :value="5" />
            <el-option label="10 min" :value="10" />
            <el-option label="15 min" :value="15" />
            <el-option label="30 min" :value="30" />
            <el-option label="60 min" :value="60" />
            <el-option label="3 hour" :value="180" />
            <el-option label="6 hour" :value="360" />
            <el-option label="12 hour" :value="720" />
            <el-option label="24 hour" :value="1440" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('notice')" prop="notifyTarget">
          <el-row>
            <el-col :span="8">
              <el-select v-model="formData.notifyType" style="width: 100%">
                <el-option :label="$t('webhookOption[1]')" :value="1" />
                <el-option :label="$t('webhookOption[2]')" :value="2" />
                <el-option :label="$t('webhookOption[3]')" :value="3" />
                <el-option :label="$t('webhookOption[255]')" :value="255" />
              </el-select>
            </el-col>
            <el-col :span="16">
              <el-input
                v-model.trim="formData.notifyTarget"
                style="flex: 1"
                autocomplete="off"
                placeholder="webhook"
              />
            </el-col>
          </el-row>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
        <el-button
          :disabled="formProps.disabled"
          type="primary"
          @click="submit"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
    <el-dialog
      v-model="monitorListDialogVisible"
      :title="$t('manage')"
      :close-on-click-modal="false"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <el-table
        v-loading="tableLoading"
        border
        stripe
        highlight-current-row
        :data="tableData"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column
          prop="item"
          :label="$t('serverPage.item')"
          width="150"
          show-overflow-tooltip
        />
        <el-table-column :label="$t('serverPage.formula')">
          <template #default="scope">
            {{ scope.row.formula }}
            {{ scope.row.operator }}
            {{ scope.row.value }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('serverPage.cycle')" width="200">
          <template #default="scope">
            {{ scope.row.startTime }}-{{ scope.row.endTime }} group by
            {{ scope.row.groupCycle }} min last {{ scope.row.lastCycle }} cycle,
            silent in {{ scope.row.silentCycle }} cycle.
          </template>
        </el-table-column>
        <el-table-column prop="notifyType" :label="$t('notice')" width="70">
          <template #default="scope">
            <el-tooltip
              class="item"
              effect="dark"
              :content="scope.row.notifyTarget"
              placement="top"
            >
              <span v-if="scope.row.notifyType === 1">
                {{ $t('webhookOption[1]') }}
              </span>
              <span v-else-if="scope.row.notifyType === 2">
                {{ $t('webhookOption[2]') }}
              </span>
              <span v-else-if="scope.row.notifyType === 3">
                {{ $t('webhookOption[3]') }}
              </span>
              <span v-else-if="scope.row.notifyType === 255">
                {{ $t('webhookOption[255]') }}
              </span>
            </el-tooltip>
          </template>
        </el-table-column>
        <el-table-column
          prop="insertTime"
          :label="$t('insertTime')"
          width="135"
          align="center"
        />
        <el-table-column
          prop="updateTime"
          :label="$t('updateTime')"
          width="135"
          align="center"
        />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="120"
          align="center"
          fixed="right"
        >
          <template #default="scope">
            <el-button
              type="primary"
              icon="el-icon-edit"
              @click="handleEdit(scope.row)"
            />
            <el-button
              type="danger"
              icon="el-icon-delete"
              @click="handleDelete(scope.row)"
            />
          </template>
        </el-table-column>
      </el-table>
      <template #footer>
        <el-button @click="monitorListDialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
      </template>
    </el-dialog>
  </el-row>
</template>
<script lang="ts">
import * as echarts from 'echarts'
import { ElMessageBox, ElMessage, ElDatePicker } from 'element-plus'
import {
  ServerMonitorData,
  ServerReport,
  ServerMonitorList,
  ServerMonitorAdd,
  ServerMonitorEdit,
  ServerMonitorDelete,
} from '@/api/server'
import Validator from 'async-validator'
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
      monitorDialogVisible: false,
      monitorListDialogVisible: false,
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
      tableData: [] as ServerMonitorData['datagram'][],
      tableLoading: false,
      tempFormData: {},
      formData: {
        id: 0,
        serverId: 0,
        item: '',
        formula: 'avg',
        operator: '>=',
        value: '',
        groupCycle: 60,
        lastCycle: 1,
        silentCycle: 1440,
        startTime: '00:00',
        endTime: '23:59',
        notifyType: 1,
        notifyTarget: '',
      },
      formProps: { loading: false, disabled: false },
      formRules: {
        item: [{ required: true, message: 'Target required', trigger: 'blur' }],
        value: [{ required: true, message: 'Value required', trigger: 'blur' }],
        notifyTarget: [{ required: true, message: 'Webhook required' }],
      },
      itemOptions: [] as string[],
    }
  },

  activated() {
    this.formData.serverId = this.serverId = Number(this.$route.query.serverId)
    if (!this.serverId) {
      this.$store.dispatch('tagsView/delView', this.$route).then(() => {
        this.$router.push('/server/index')
      })
    } else {
      this.storeFormData()
      this.handleRefresh()
    }
  },

  methods: {
    datetimeChange(values: Date[]) {
      for (const key in this.chartNameMap) {
        this.report(key, values)
      }
    },

    handleGetList() {
      this.monitorListDialogVisible = true
      this.getList()
    },

    handleAdd() {
      this.restoreFormData()
      this.monitorDialogVisible = true
    },

    handleEdit(data: ServerMonitorData['datagram']) {
      this.formData = Object.assign({}, data)
      this.monitorDialogVisible = true
    },

    handleDelete(data: ServerMonitorData['datagram']) {
      ElMessageBox.confirm(
        this.$t('serverPage.removeMonitorTips', {
          item: data.item,
        }),
        this.$t('tips'),
        {
          confirmButtonText: this.$t('confirm'),
          cancelButtonText: this.$t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          new ServerMonitorDelete({ id: data.id }).request().then(() => {
            ElMessage.success('Success')
            this.getList()
          })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },

    submit() {
      ;(this.$refs.form as Validator).validate((valid: boolean) => {
        if (valid) {
          if (this.formData.id === 0) {
            this.add()
          } else {
            this.edit()
          }
        } else {
          return false
        }
      })
    },

    add() {
      this.formProps.disabled = true
      new ServerMonitorAdd(this.formData)
        .request()
        .then(() => {
          ElMessage.success('Success')
        })
        .finally(() => {
          this.formProps.disabled = this.monitorDialogVisible = false
        })
    },

    edit() {
      this.formProps.disabled = true
      new ServerMonitorEdit(this.formData)
        .request()
        .then(() => {
          ElMessage.success('Success')
          this.getList()
        })
        .finally(() => {
          this.formProps.disabled = this.monitorDialogVisible = false
        })
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

    getList() {
      this.tableLoading = true
      new ServerMonitorList({ serverId: this.serverId })
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
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
            this.itemOptions.push(key)
            this.itemOptions = this.itemOptions.filter(
              (e, i, a) => i === a.indexOf(e)
            )
            this.itemOptions.sort()
            const series = {
              name: key,
              type: 'line',
              symbol: 'none',
              smooth: true,
              data: response.data.map[key].map(
                (item: { reportTime: string; value: string }) => {
                  return [new Date(item.reportTime), item.value]
                }
              ),
            }
            chartOption.series.push(series)
          }
          chart.setOption(chartOption)
        })
    },

    storeFormData() {
      this.tempFormData = JSON.parse(JSON.stringify(this.formData))
    },

    restoreFormData() {
      this.formData = JSON.parse(JSON.stringify(this.tempFormData))
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
