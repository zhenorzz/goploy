<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-col :span="8">
        <el-select
          v-model="serverId"
          placeholder="Select server"
          style="width: 160px"
          filterable
          @change="selectServer"
        >
          <el-option
            v-for="server in serverOption"
            :key="server.id"
            :label="server.label"
            :value="server.id"
          />
        </el-select>
      </el-col>
      <el-col v-if="serverId !== ''" :span="16" style="text-align: right">
        <el-button
          :loading="tableLoading"
          type="primary"
          :icon="Plus"
          @click="handleAdd"
        >
          New
        </el-button>
        <el-button
          :loading="tableLoading"
          type="primary"
          :icon="Film"
          @click="refresList"
        >
          Backup
        </el-button>
        <el-button
          :loading="tableLoading"
          type="warning"
          :icon="Download"
          @click="refresList"
        >
          Import
        </el-button>
        <el-button
          :loading="tableLoading"
          type="warning"
          :icon="Upload"
          @click="refresList"
        >
          Export
        </el-button>
        <el-button
          :loading="tableLoading"
          type="success"
          :icon="Refresh"
          @click="getRemoteCrontabList"
        >
          Get from crontab
        </el-button>
        <el-button
          :loading="tableLoading"
          type="success"
          :icon="Document"
          @click="refresList"
        >
          Save to crontab
        </el-button>
      </el-col>
    </el-row>
    <el-row class="app-table">
      <el-table
        v-loading="tableLoading"
        height="100%"
        highlight-current-row
        :data="tablePage.list"
      >
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column
          prop="expression"
          :label="$t('expression')"
          min-width="120"
          show-overflow-tooltip
        />
        <el-table-column
          prop="description"
          :label="$t('description')"
          min-width="120"
          show-overflow-tooltip
        />
        <el-table-column
          prop="insertTime"
          :label="$t('insertTime')"
          width="155"
          align="center"
        />
        <el-table-column
          prop="updateTime"
          :label="$t('updateTime')"
          width="155"
          align="center"
        />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="130"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <Button
              type="primary"
              :icon="Edit"
              :permissions="[pms.EditCron]"
              @click="handleEdit(scope.row)"
            />
            <Button
              type="danger"
              :icon="Delete"
              :permissions="[pms.DeleteCron]"
              @click="handleRemove(scope.row)"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-row>
    <el-row type="flex" justify="end" class="app-page">
      <el-pagination
        :total="tablePage.total"
        :page-size="pagination.rows"
        background
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
    <el-dialog
      v-model="dialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('setting')"
    >
      <el-form
        ref="form"
        v-loading="formProps.loading"
        :rules="formRules"
        :model="formData"
        label-width="120px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item :label="$t('command')" prop="command">
          <el-input v-model="formData.command" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Quick Schedule">
          <el-row style="width: 100%">
            <el-button type="primary" @click="handleQuickSet('startup')">
              Startup
            </el-button>
            <el-button type="primary" @click="handleQuickSet('hourly')">
              Hourly
            </el-button>
            <el-button type="primary" @click="handleQuickSet('daily')">
              Daily
            </el-button>
            <el-button type="primary" @click="handleQuickSet('weekly')">
              Weekly
            </el-button>
            <el-button type="primary" @click="handleQuickSet('monthly')">
              Monthly
            </el-button>
            <el-button type="primary" @click="handleQuickSet('yearly')">
              Yearly
            </el-button>
          </el-row>
        </el-form-item>
        <el-form-item>
          <el-row style="width: 100%" align="bottom">
            <el-row style="width: 60px; margin-right: 10px">
              <span>Minute</span>
              <el-input v-model="formData.minute" />
            </el-row>
            <el-row style="width: 60px; margin-right: 10px">
              <span>Hour</span>
              <el-input v-model="formData.hour" />
            </el-row>
            <el-row style="width: 60px; margin-right: 10px">
              <span>Day</span>
              <el-input v-model="formData.day" />
            </el-row>
            <el-row style="width: 60px; margin-right: 10px">
              <span>Month</span>
              <el-input v-model="formData.month" />
            </el-row>
            <el-row style="width: 60px; margin-right: 10px">
              <span>Week</span>
              <el-input v-model="formData.week" />
            </el-row>
            <el-button
              :disabled="
                formData.minute == '' ||
                formData.hour == '' ||
                formData.day == '' ||
                formData.month == '' ||
                formData.week == ''
              "
              type="primary"
              @click="handleSetTime"
            >
              Set
            </el-button>
          </el-row>
          <span>{{ formProps.dateLocale }}</span>
        </el-form-item>
        <el-form-item :label="$t('expression')" prop="expression">
          <el-input v-model="formData.expression" disabled />
        </el-form-item>
        <el-form-item :label="$t('description')">
          <el-input v-model="formData.description" autocomplete="off" />
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
  </el-row>
</template>
<script lang="ts">
export default { name: 'ServerCron' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import Button from '@/components/Permission/Button.vue'
import {
  Film,
  Upload,
  Download,
  Document,
  Refresh,
  Plus,
  Edit,
  Delete,
} from '@element-plus/icons-vue'
import cronstrue from 'cronstrue/i18n'
import { ServerRemoteCrontabList, ServerOption } from '@/api/server'
import { CronList, CronAdd, CronEdit, CronRemove, CronData } from '@/api/cron'
import type { ElForm } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
const { locale, t } = useI18n({ useScope: 'global' })
const serverId = ref('')
const dialogVisible = ref(false)
const serverOption = ref<ServerOption['datagram']['list']>([])
const tableLoading = ref(false)
const tableData = ref<CronList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = {
  id: 0,
  serverId: 0,
  name: '',
  command: '/bin/sh /data/sh/hello.sh',
  minute: '*/1',
  hour: '*',
  day: '*',
  month: '*',
  week: '*',
  expression: '',
  description: '',
}
const formData = ref(tempFormData)
const formProps = ref({
  loading: false,
  disabled: false,
  dateLocale: '',
})
const formRules: InstanceType<typeof ElForm>['rules'] = {
  expression: [
    {
      required: true,
      validator: (_, value) => {
        if (value.trim().split(/\s+/).length != 6) {
          return new Error('6 parts are required.')
        }
        try {
          cronstrue.toString(value)
          return true
        } catch (error) {
          if (typeof error === 'string') {
            return new Error(error)
          } else if (error instanceof Error) {
            return error
          }
        }
      },
      trigger: 'blur',
    },
  ],
  command: [{ required: true, message: 'Command required', trigger: 'blur' }],
}

getServerOption()

function selectServer() {
  getList()
}

function getServerOption() {
  new ServerOption().request().then((response) => {
    serverOption.value = response.data.list
  })
}

function getRemoteCrontabList() {
  tableLoading.value = true
  tableData.value = []
  new ServerRemoteCrontabList({ serverId: Number(serverId.value) })
    .request()
    .then((response) => {
      console.log(response)
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function getList() {
  tableLoading.value = true
  tableData.value = []
  new CronList({ serverId: Number(serverId.value) })
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

const tablePage = computed(() => {
  let _tableData = tableData.value
  return {
    list: _tableData.slice(
      (pagination.value.page - 1) * pagination.value.rows,
      pagination.value.page * pagination.value.rows
    ),
    total: _tableData.length,
  }
})

function refresList() {
  pagination.value.page = 1
  getList()
}

function handleAdd() {
  restoreFormData()
  formData.value.serverId = Number(serverId.value)
  dialogVisible.value = true
}

function handleEdit(data: CronData) {
  formData.value = data
  dialogVisible.value = true
}

function handleRemove(data: CronData) {
  ElMessageBox.confirm(
    t('serverPage.deleteTips', { name: data.command }),
    t('tips'),
    {
      confirmButtonText: t('confirm'),
      cancelButtonText: t('cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      new CronRemove({ id: data.id }).request().then(() => {
        getList()
        ElMessage.success('Success')
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function handleQuickSet(schedule: string) {
  switch (schedule) {
    case 'startup':
      formData.value.expression = `@reboot ${formData.value.command}`
      formData.value.minute = ''
      formData.value.hour = ''
      formData.value.day = ''
      formData.value.month = ''
      formData.value.week = ''
      onExpressionChange()
      break
    case 'hourly':
      formData.value.minute = '0'
      formData.value.hour = '*'
      formData.value.day = '*'
      formData.value.month = '*'
      formData.value.week = '*'
      handleSetTime()
      break
    case 'daily':
      formData.value.minute = '0'
      formData.value.hour = '0'
      formData.value.day = '*'
      formData.value.month = '*'
      formData.value.week = '*'
      handleSetTime()
      break
    case 'monthly':
      formData.value.minute = '0'
      formData.value.hour = '0'
      formData.value.day = '1'
      formData.value.month = '*'
      formData.value.week = '*'
      handleSetTime()
      break
    case 'weekly':
      formData.value.minute = '0'
      formData.value.hour = '0'
      formData.value.day = '*'
      formData.value.month = '*'
      formData.value.week = '0'
      handleSetTime()
      break
    case 'yearly':
      formData.value.minute = '0'
      formData.value.hour = '0'
      formData.value.day = '1'
      formData.value.month = '1'
      formData.value.week = '*'
      handleSetTime()
      break
    default:
      break
  }
}

function handleSetTime() {
  formData.value.expression = `${formData.value.minute} ${formData.value.hour} ${formData.value.day} ${formData.value.month} ${formData.value.week} ${formData.value.command}`
  onExpressionChange()
}

function onExpressionChange() {
  if (formData.value.expression.startsWith('@reboot')) {
    formProps.value.dateLocale = 'Startup'
  } else {
    formProps.value.dateLocale = cronstrue.toString(
      `${formData.value.minute} ${formData.value.hour} ${formData.value.day} ${formData.value.month} ${formData.value.week}`,
      {
        use24HourTimeFormat: true,
        locale: getLocale(),
      }
    )
  }
}

function handlePageChange(val = 1) {
  pagination.value.page = val
}

function submit() {
  form.value?.validate((valid) => {
    formData.value.expression = formData.value.expression.trim()
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
  new CronAdd(formData.value)
    .request()
    .then(() => {
      getList()
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = dialogVisible.value = false
    })
}

function edit() {
  formProps.value.disabled = true
  new CronEdit(formData.value)
    .request()
    .then(() => {
      getList()
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = dialogVisible.value = false
    })
}

function getLocale() {
  if (locale.value === 'zh-cn') {
    return 'zh_CN'
  }
  return locale.value
}

function restoreFormData() {
  formData.value = { ...tempFormData }
}
</script>

<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.template-dialog {
  padding-right: 10px;
  height: 400px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
