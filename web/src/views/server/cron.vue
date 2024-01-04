<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-col :span="16">
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
            :label="server.name"
            :value="server.id"
          />
        </el-select>
      </el-col>
      <el-col v-if="serverId !== ''" :span="8" style="text-align: right">
        <el-button
          style="margin-left: 10px"
          :loading="tableLoading"
          type="primary"
          :icon="Refresh"
          @click="refresList"
        />
        <Button
          type="primary"
          :icon="Plus"
          :permissions="[pms.AddCron]"
          @click="handleAdd"
        />
      </el-col>
    </el-row>
    <el-row class="app-table">
      <el-table
        v-loading="tableLoading"
        height="100%"
        highlight-current-row
        :data="tablePage.list"
      >
        <el-table-column
          prop="expression"
          :label="$t('expression')"
          min-width="120"
          show-overflow-tooltip
        />
        <el-table-column
          prop="command"
          :label="$t('command')"
          min-width="140"
          show-overflow-tooltip
        />
        <el-table-column prop="singleMode" label="Single mode" width="110">
          <template #default="scope">
            <span v-if="scope.row.singleMode === 0">no</span>
            <span v-else>yes</span>
          </template>
        </el-table-column>
        <el-table-column prop="logLevel" label="Log level" width="90">
          <template #default="scope">
            <span v-if="scope.row.logLevel === 0">none</span>
            <span v-else-if="scope.row.logLevel === 1">stdout</span>
            <span v-else-if="scope.row.logLevel === 2">stdout+stderr</span>
          </template>
        </el-table-column>
        <el-table-column
          prop="description"
          :label="$t('description')"
          min-width="240"
          show-overflow-tooltip
        />
        <el-table-column prop="creator" :label="$t('creator')" min-width="80" />
        <el-table-column prop="editor" :label="$t('editor')" min-width="80" />
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
          width="190"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <el-button :icon="Tickets" @click="handleShowLogs(scope.row)" />
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
        :model="formData"
        label-width="120px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item
          :label="$t('expression')"
          prop="expression"
          :rules="formRules.expression"
        >
          <el-input
            v-model="formData.expression"
            autocomplete="off"
            placeholder="* * * * * ? with second"
            @change="onExpressionChange"
          />
          <span>{{ formProps.dateLocale }}</span>
        </el-form-item>
        <el-form-item
          :label="$t('command')"
          prop="command"
          :rules="[
            { required: true, message: 'Command required', trigger: 'blur' },
          ]"
        >
          <el-input v-model="formData.command" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Single mode">
          <el-radio-group v-model="formData.singleMode">
            <el-radio :label="0">no</el-radio>
            <el-radio :label="1">yes</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="Log level">
          <el-radio-group v-model="formData.logLevel">
            <el-radio :label="0">none</el-radio>
            <el-radio :label="1">stdout</el-radio>
            <el-radio :label="2">stdout+stderr</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('description')" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :autosize="{ minRows: 2 }"
          />
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
      v-model="cronDialogVisible"
      :title="$t('detail')"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <el-table
        v-loading="cronTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="cronTableData"
      >
        <el-table-column prop="reportTime" label="Report time" width="160" />
        <el-table-column prop="message" label="Message" />
        <el-table-column prop="result" label="Result" width="120">
          <template #default="scope">
            <span v-if="scope.row.execCode > 0" style="color: #f56c6c">
              {{ $t('fail') }} ({{ scope.row.execCode }})
            </span>
            <span v-else style="color: #67c23a">
              {{ $t('success') }}
            </span>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </el-row>
</template>
<script lang="ts">
export default { name: 'ServerCron' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Refresh, Plus, Edit, Delete, Tickets } from '@element-plus/icons-vue'
import cronstrue from 'cronstrue/i18n'
import { ServerOption } from '@/api/server'
import {
  CronList,
  CronLogs,
  CronAdd,
  CronEdit,
  CronRemove,
  CronData,
} from '@/api/cron'
import type { ElForm } from 'element-plus'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
const { locale, t } = useI18n({ useScope: 'global' })
const serverId = ref('')
const dialogVisible = ref(false)
const cronDialogVisible = ref(false)
const serverOption = ref<ServerOption['datagram']['list']>([])
const tableLoading = ref(false)
const tableData = ref<CronList['datagram']['list']>([])
const cronTableLoading = ref(false)
const cronTableData = ref<CronLogs['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = {
  id: 0,
  serverId: 0,
  expression: '',
  command: '',
  singleMode: 0,
  logLevel: 0,
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

function handleShowLogs(data: CronData) {
  new CronLogs({
    serverId: Number(serverId.value),
    cronId: data.id,
    page: 1,
    rows: 50,
  })
    .request()
    .then((response) => {
      cronTableData.value = response.data.list
    })
    .finally(() => {
      cronTableLoading.value = false
    })
  cronDialogVisible.value = true
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
  ElMessageBox.confirm(t('deleteTips', { name: data.command }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
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

function onExpressionChange() {
  if (formData.value.expression.trim().split(/\s+/).length != 6) {
    return
  }
  formProps.value.dateLocale = cronstrue.toString(formData.value.expression, {
    locale: getLocale(),
  })
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
