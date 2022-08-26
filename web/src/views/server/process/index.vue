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
          :permissions="[pms.AddServerProcess]"
          @click="handleAdd"
        />
      </el-col>
    </el-row>
    <el-row class="app-table">
      <el-table
        ref="table"
        v-loading="tableLoading"
        height="100%"
        highlight-current-row
        :data="tablePage.list"
        style="width: 100%"
      >
        <el-table-column type="expand">
          <template #default="{}">
            <div style="padding: 0 20px">
              <el-row style="margin-left: 4px">
                {{ $t('deployPage.execRes') }}:
                <span
                  :class="commandRes.execRes ? 'exec-success' : 'exec-fail'"
                  style="padding-left: 5px"
                >
                  {{ commandRes.execRes }}
                </span>
              </el-row>
              <el-row style="white-space: pre-wrap">
                stdout: {{ commandRes.stdout }}
              </el-row>
              <el-row style="white-space: pre-wrap">
                stderr:{{ commandRes.stderr }}
              </el-row>
            </div>
          </template>
        </el-table-column>
        <el-table-column
          prop="name"
          :label="$t('name')"
          min-width="120"
          show-overflow-tooltip
        />
        <el-table-column prop="status" label="Status" align="center">
          <template #default="scope">
            <el-button
              :loading="selectedItem['id'] === scope.row.id"
              type="primary"
              text
              @click="handleProcessCmd(scope.row, 'status')"
            >
              status
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="start" label="Start" align="center">
          <template #default="scope">
            <el-button
              :loading="selectedItem['id'] === scope.row.id"
              type="success"
              text
              @click="handleProcessCmd(scope.row, 'start')"
            >
              start
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="restart" label="restart" align="center">
          <template #default="scope">
            <el-button
              :loading="selectedItem['id'] === scope.row.id"
              type="warning"
              text
              @click="handleProcessCmd(scope.row, 'restart')"
            >
              restart
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="stop" label="stop" align="center">
          <template #default="scope">
            <el-button
              :loading="selectedItem['id'] === scope.row.id"
              type="danger"
              text
              @click="handleProcessCmd(scope.row, 'stop')"
            >
              stop
            </el-button>
          </template>
        </el-table-column>
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
        label-width="80px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item label="Status">
          <el-input v-model="formData.status" />
        </el-form-item>
        <el-form-item label="Start">
          <el-input v-model="formData.start" />
        </el-form-item>
        <el-form-item label="Stop">
          <el-input v-model="formData.stop" />
        </el-form-item>
        <el-form-item label="Restart">
          <el-input v-model="formData.restart" />
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
export default { name: 'ServerProcess' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Refresh, Plus, Edit, Delete } from '@element-plus/icons-vue'
import {
  ServerOption,
  ServerProcessList,
  ServerProcessAdd,
  ServerProcessEdit,
  ServerProcessDelete,
  ServerProcessData,
  ServerExecProcess,
} from '@/api/server'
import type { ElForm } from 'element-plus'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n({ useScope: 'global' })
const serverId = ref('')
const dialogVisible = ref(false)
const serverOption = ref<ServerOption['datagram']['list']>([])
const table = ref()
const tableLoading = ref(false)
const tableData = ref<ServerProcessList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const selectedItem = ref<ServerProcessData>({})
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = {
  id: 0,
  serverId: 0,
  name: '',
  status: '',
  start: '',
  stop: '',
  restart: '',
}
const formData = ref(tempFormData)
const formProps = ref({
  loading: false,
  disabled: false,
})
const formRules = <InstanceType<typeof ElForm>['rules']>{
  name: [{ required: true, message: 'Name required', trigger: 'blur' }],
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
  new ServerProcessList({ serverId: Number(serverId.value) })
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

function handleEdit(data: ServerProcessData) {
  formData.value = data
  dialogVisible.value = true
}

function handleRemove(data: ServerProcessData) {
  ElMessageBox.confirm(
    t('serverPage.deleteTips', { name: data.name }),
    t('tips'),
    {
      confirmButtonText: t('confirm'),
      cancelButtonText: t('cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      new ServerProcessDelete({ id: data.id }).request().then(() => {
        getList()
        ElMessage.success('Success')
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}
const commandRes = ref<ServerExecProcess['datagram']>({
  execRes: true,
  stdout: '',
  stderr: '',
})

const handleProcessCmd = (data: ServerProcessData, command: string) => {
  ElMessageBox.confirm(t('deployPage.execTips', { command }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      selectedItem.value = data
      new ServerExecProcess({
        id: data.id,
        serverId: data.serverId,
        command,
      })
        .request()
        .then((response) => {
          commandRes.value = response.data
          table.value.toggleRowExpansion(data, true)
        })
        .finally(() => {
          selectedItem.value = {}
        })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function handlePageChange(val = 1) {
  pagination.value.page = val
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
  new ServerProcessAdd(formData.value)
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
  new ServerProcessEdit(formData.value)
    .request()
    .then(() => {
      getList()
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = dialogVisible.value = false
    })
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

.exec-success {
  color: #67c23a;
}

.exec-fail {
  color: #f56c6c;
}
</style>
