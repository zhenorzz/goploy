<template>
  <el-row class="app-container" style="flex-wrap: nowrap">
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
        <el-button :icon="Refresh" @click="handleRefresh" />
      </el-row>
      <el-row>
        <Button
          :icon="Bell"
          :permissions="[pms.AddServerWarningRule]"
          @click="handleAdd"
        >
          {{ $t('serverPage.addMonitor') }}
        </Button>
        <el-button :icon="Tickets" @click="handleGetList">
          {{ $t('serverPage.monitorList') }}
        </el-button>
      </el-row>
    </el-row>
    <el-scrollbar class="chart-container">
      <el-row style="width: 100%" :gutter="10">
        <el-col
          v-for="(_, name) in chartNameMap"
          :key="name"
          :xs="24"
          :sm="24"
          :lg="12"
          style="margin-bottom: 10px"
        >
          <div
            :ref="
              (el) => {
                if (el) {
                  chartRefs[name] = el
                }
              }
            "
            style="
              height: 288px;
              padding-top: 10px;
              border: 1px solid var(--el-border-color);
            "
          ></div>
        </el-col>
      </el-row>
    </el-scrollbar>
    <el-dialog
      v-model="monitorDialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('setting')"
      :close-on-click-modal="false"
    >
      <el-form
        ref="form"
        v-loading="formProps.loading"
        :model="formData"
        label-width="100px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item
          :label="$t('serverPage.item')"
          prop="item"
          :rules="[
            { required: true, message: 'Target required', trigger: 'blur' },
          ]"
        >
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
        <el-form-item
          :label="$t('serverPage.formula')"
          prop="value"
          :rules="[
            { required: true, message: 'Value required', trigger: 'blur' },
          ]"
        >
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
        <el-form-item :label="$t('serverPage.validPeriod')">
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
        <el-form-item
          :label="$t('notice')"
          prop="notifyTarget"
          :rules="[{ required: true, message: 'Webhook required' }]"
        >
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
        <el-button @click="monitorDialogVisible = false">
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
          width="160"
          align="center"
        />
        <el-table-column
          prop="updateTime"
          :label="$t('updateTime')"
          width="160"
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
            <Button
              type="primary"
              :icon="Edit"
              :permissions="[pms.EditServerWarningRule]"
              @click="handleEdit(scope.row)"
            />
            <Button
              type="danger"
              :icon="Delete"
              :permissions="[pms.DeleteServerWarningRule]"
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
export default { name: 'ServerAgent' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Refresh, Bell, Tickets, Edit, Delete } from '@element-plus/icons-vue'
import * as echarts from 'echarts/core'
import { LineChart } from 'echarts/charts'
import 'echarts/theme/dark-mushroom.js'
import { CanvasRenderer } from 'echarts/renderers'
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
} from 'echarts/components'
import type { ElForm, ElDatePicker } from 'element-plus'
import {
  ServerMonitorData,
  ServerReport,
  ServerMonitorList,
  ServerMonitorAdd,
  ServerMonitorEdit,
  ServerMonitorDelete,
} from '@/api/server'
import dayjs, { Dayjs } from 'dayjs'
import { ref, onActivated, ComponentPublicInstance } from 'vue'
import { deepClone, parseTime } from '@/utils'
import { useI18n } from 'vue-i18n'
import { useDark } from '@vueuse/core'
import { useRoute, useRouter } from 'vue-router'
import { useStore } from 'vuex'
const isDark = useDark()
const { t } = useI18n()
const store = useStore()
const route = useRoute()
const router = useRouter()
let serverId = Number(route.query.serverId)
const monitorDialogVisible = ref(false)
const monitorListDialogVisible = ref(false)
echarts.use([
  LineChart,
  CanvasRenderer,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent,
])
const chartRefs = ref<Record<string, Element | ComponentPublicInstance>>({})
const chartNameMap: Record<
  string,
  { type: number; title: string; subtitle: string }
> = {
  cpuChart: {
    type: 1,
    title: t('serverPage.cpuUsage'),
    subtitle: '(%)',
  },
  ramChart: {
    type: 2,
    title: t('serverPage.ramUsage'),
    subtitle: '(%)',
  },
  loadavgChart: {
    type: 3,
    title: t('serverPage.loadavg'),
    subtitle: '',
  },
  tcpChart: {
    type: 4,
    title: t('serverPage.tcp'),
    subtitle: '(count)',
  },
  pubNetChart: {
    type: 5,
    title: t('serverPage.pubBandwidth'),
    subtitle: '(bit/s)',
  },
  loNetChart: {
    type: 6,
    title: t('serverPage.loBandwidth'),
    subtitle: '(bit/s)',
  },
  diskUsageChart: {
    type: 7,
    title: t('serverPage.diskUsage'),
    subtitle: '(%)',
  },
  diskIOChart: {
    type: 8,
    title: t('serverPage.diskIO'),
    subtitle: '(count/s)',
  },
}
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
const shortcuts = [
  {
    text: t('lastHour'),
    onClick(picker: typeof ElDatePicker) {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000)
      picker.emit('pick', [dayjs(start), dayjs(end)])
    },
  },
  {
    text: t('last6Hours'),
    onClick(picker: typeof ElDatePicker) {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 6)
      picker.emit('pick', [dayjs(start), dayjs(end)])
    },
  },
  {
    text: t('lastDay'),
    onClick(picker: typeof ElDatePicker) {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24)
      picker.emit('pick', [dayjs(start), dayjs(end)])
    },
  },
  {
    text: t('lastWeek'),
    onClick(picker: typeof ElDatePicker) {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      picker.emit('pick', [dayjs(start), dayjs(end)])
    },
  },
]
const datetimeRange = ref<Dayjs[]>([])
const tableData = ref<ServerMonitorData[]>([])
const tableLoading = ref(false)
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = {
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
}
const formData = ref(tempFormData)
const formProps = ref({ loading: false, disabled: false })
const itemOptions = ref<string[]>([])
onActivated(() => {
  formData.value.serverId = serverId = Number(route.query.serverId)
  if (!serverId) {
    store.dispatch('tagsView/delView', route).then(() => {
      router.push('/server/index')
    })
  } else {
    handleRefresh()
  }
})

function datetimeChange(values: Date[]) {
  for (const key in chartNameMap) {
    report(key, values)
  }
}

function handleGetList() {
  monitorListDialogVisible.value = true
  getList()
}

function handleAdd() {
  restoreFormData()
  monitorDialogVisible.value = true
}

function handleEdit(data: ServerMonitorData) {
  formData.value = Object.assign({}, data)
  monitorDialogVisible.value = true
}

function handleDelete(data: ServerMonitorData) {
  ElMessageBox.confirm(t('deleteTips', { name: data.item }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      new ServerMonitorDelete({ id: data.id }).request().then(() => {
        ElMessage.success('Success')
        getList()
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function submit() {
  form.value?.validate((valid) => {
    if (valid) {
      if (formData.value.id === 0) {
        add()
      } else {
        edit()
      }
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
}

function add() {
  formProps.value.disabled = true
  new ServerMonitorAdd(formData.value)
    .request()
    .then(() => {
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = monitorDialogVisible.value = false
    })
}

function edit() {
  formProps.value.disabled = true
  new ServerMonitorEdit(formData.value)
    .request()
    .then(() => {
      ElMessage.success('Success')
      getList()
    })
    .finally(() => {
      formProps.value.disabled = monitorDialogVisible.value = false
    })
}

function handleRefresh() {
  const end = new Date()
  const start = new Date()
  start.setTime(start.getTime() - 3600 * 1000)
  datetimeRange.value = [dayjs(start), dayjs(end)]
  datetimeChange(datetimeRange.value.map((datetime) => datetime.toDate()))
}

function getList() {
  tableLoading.value = true
  new ServerMonitorList({ serverId })
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function report(chartName: string, values: Date[]) {
  new ServerReport({
    serverId,
    type: chartNameMap[chartName].type,
    datetimeRange: values.map((value) => parseTime(value)).join(','),
  })
    .request()
    .then((response) => {
      echarts.dispose(chartRefs.value[chartName] as HTMLDivElement)
      let chart
      if (isDark.value) {
        chart = echarts.init(
          chartRefs.value[chartName] as HTMLDivElement,
          'dark-mushroom'
        )
      } else {
        chart = echarts.init(chartRefs.value[chartName] as HTMLDivElement)
      }

      let chartOption = deepClone(chartBaseOption)
      chartOption.title.text = chartNameMap[chartName].title
      chartOption.title.subtext = chartNameMap[chartName].subtitle
      for (const key in response.data.map) {
        chartOption.legend.data.push(key)
        itemOptions.value.push(key)
        itemOptions.value = itemOptions.value.filter(
          (e, i, a) => i === a.indexOf(e)
        )
        itemOptions.value.sort()
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
}

function restoreFormData() {
  formData.value = { ...tempFormData }
}
</script>

<style lang="scss" scoped>
@import '@/styles/mixin.scss';

.chart-container {
  width: 100%;
  max-height: calc(100vh - 180px);
}
</style>
