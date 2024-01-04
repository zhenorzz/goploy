<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-input
        v-model="name"
        style="width: 180px"
        placeholder="Filter the name"
      />
      <el-row>
        <el-button
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
      </el-row>
    </el-row>
    <el-row class="app-table">
      <el-table
        v-loading="tableLoading"
        height="100%"
        highlight-current-row
        :data="tablePage.list"
      >
        <el-table-column
          prop="id"
          label="ID"
          min-width="40"
          show-overflow-tooltip
        />
        <el-table-column
          prop="name"
          :label="$t('name')"
          min-width="120"
          show-overflow-tooltip
        />
        <el-table-column
          prop="items"
          min-width="240"
          :label="$t('command')"
          show-overflow-tooltip
        >
          <template #default="scope">
            {{ JSON.stringify(scope.row.items) }}
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
          width="260"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <el-button :icon="CaretRight" @click="handleProcess(scope.row)" />
            <Button
              type="primary"
              :icon="Edit"
              :permissions="[pms.EditServerProcess]"
              @click="handleEdit(scope.row)"
            />
            <Button
              type="info"
              :icon="DocumentCopy"
              :permissions="[pms.AddServerProcess]"
              @click="handleCopy(scope.row)"
            />
            <Button
              type="danger"
              :icon="Delete"
              :permissions="[pms.DeleteServerProcess]"
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
        label-width="80px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item
          :label="$t('name')"
          prop="name"
          :rules="[
            { required: true, message: 'Name required', trigger: 'blur' },
          ]"
        >
          <el-input v-model="formData.name" />
        </el-form-item>
        <el-form-item :label="$t('command')">
          <el-button
            type="primary"
            :icon="Plus"
            plain
            @click="handleCommandAdd"
          />
        </el-form-item>
        <el-form-item v-for="(item, index) in formData.items" :key="index">
          <el-row style="width: 100%">
            <el-row style="flex: 1">
              <el-input v-model="item.name">
                <template #prepend>Alias</template>
              </el-input>
              <el-input v-model="item.command" type="textarea" />
            </el-row>
            <el-button
              type="warning"
              :icon="Minus"
              plain
              @click="handleCommandDel(index)"
            />
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
      v-model="processDialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('process')"
    >
      <el-select
        v-model="serverIds"
        placeholder="Select server"
        placement="top"
        style="width: 100%"
        filterable
        multiple
        @change="serverChange"
      >
        <el-option :value="0" label="Select all" />
        <el-option
          v-for="server in serverOption"
          :key="server.id"
          :label="server.name"
          :value="server.id"
        />
      </el-select>
      <el-button
        v-for="(item, index) in selectedItem.items"
        :key="index"
        :loading="processExecLoading"
        type="primary"
        plain
        style="margin: 10px 10px 10px 0"
        @click="handleProcessCmd(item)"
      >
        {{ item.name }}<el-icon><CaretRight /></el-icon>
      </el-button>
      <el-tabs type="border-card" tab-position="left" style="height: 350px">
        <el-tab-pane v-for="serverId in serverIds" :key="serverId">
          <template #label>
            <el-row align="middle">
              <el-row v-if="processExecRes[serverId]">
                <el-icon
                  v-if="processExecRes[serverId]['execRes']"
                  style="color: var(--el-color-success)"
                >
                  <SuccessFilled />
                </el-icon>
                <el-icon v-else style="color: var(--el-color-danger)">
                  <CircleCloseFilled />
                </el-icon>
              </el-row>
              <span style="margin-left: 5px">
                {{ id2server(serverId)['name'] }}
              </span>
            </el-row>
          </template>
          <el-scrollbar style="height: 320px">
            stdout:
            <pre>{{
              processExecRes[serverId] && processExecRes[serverId]['stdout']
            }}</pre>
            stderr:
            <pre>{{
              processExecRes[serverId] && processExecRes[serverId]['stderr']
            }}</pre>
          </el-scrollbar>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </el-row>
</template>
<script lang="ts">
export default { name: 'ServerProcess' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import Button from '@/components/Permission/Button.vue'
import {
  Refresh,
  Plus,
  Minus,
  Edit,
  DocumentCopy,
  CaretRight,
  SuccessFilled,
  CircleCloseFilled,
  Delete,
} from '@element-plus/icons-vue'
import {
  ServerOption,
  ServerProcessList,
  ServerProcessAdd,
  ServerProcessEdit,
  ServerProcessDelete,
  ServerProcessData,
  ServerExecProcess,
  ServerData,
} from '@/api/server'
import type { ElForm } from 'element-plus'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n({ useScope: 'global' })
const name = ref('')
const serverIds = ref<number[]>([])
const dialogVisible = ref(false)
const processDialogVisible = ref(false)
const processExecLoading = ref(false)
const serverOption = ref<ServerOption['datagram']['list']>([])
const tableLoading = ref(false)
const tableData = ref<ServerProcessList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const selectedItem = ref<ServerProcessData>({} as ServerProcessData)
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = {
  id: 0,
  name: '',
  items: [] as { name: string; command: string }[],
}
const formData = ref(tempFormData)
const formProps = ref({
  loading: false,
  disabled: false,
})

getServerOption()
getList()

function getServerOption() {
  new ServerOption().request().then((response) => {
    serverOption.value = response.data.list
  })
}

function getList() {
  tableLoading.value = true
  tableData.value = []
  new ServerProcessList()
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
  if (name.value !== '') {
    _tableData = _tableData.filter(
      (item) => item.name.indexOf(name.value) !== -1
    )
  }
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
  dialogVisible.value = true
}

function handleEdit(data: ServerProcessData) {
  formData.value = { ...data, items: [...data.items] }
  dialogVisible.value = true
}

function handleCommandAdd() {
  formData.value.items.push({ name: '', command: '' })
}

function handleCommandDel(index: number) {
  formData.value.items.splice(index, 1)
}

function handleCopy(data: ServerProcessData) {
  formData.value = Object.assign({}, data)
  formData.value.id = 0
  dialogVisible.value = true
}

function handleProcess(data: ServerProcessData) {
  selectedItem.value = data
  processDialogVisible.value = true
}

function handleRemove(data: ServerProcessData) {
  ElMessageBox.confirm(t('deleteTips', { name: data.name }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
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
const processExecRes = ref<Record<number, ServerExecProcess['datagram']>>({})

const handleProcessCmd = async (item: { name: string; command: string }) => {
  if (serverIds.value.length === 0) {
    ElMessage.error('Select server')
    return
  }
  let result: string
  try {
    result = await ElMessageBox.confirm(
      t('deployPage.execTips', { command: item.command }),
      t('tips'),
      {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning',
      }
    )
  } catch (error) {
    result = error
  }
  if (result !== 'confirm') {
    return
  }
  processExecLoading.value = true
  Promise.all(
    serverIds.value.map(async (serverId) => {
      return await new ServerExecProcess({
        id: selectedItem.value.id,
        serverId: serverId,
        name: item.name,
      })
        .request()
        .then((response) => {
          processExecRes.value[serverId] = response.data
          return response
        })
    })
  ).then(() => {
    processExecLoading.value = false
  })
}

function handlePageChange(val = 1) {
  pagination.value.page = val
}

function serverChange(_serverIds: number[]) {
  if (_serverIds.includes(0)) {
    serverIds.value = _serverIds = serverOption.value.map((_) => _.id)
  }
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
  new ServerProcessAdd({
    ...formData.value,
    items: JSON.stringify(formData.value.items),
  })
    .request()
    .then(() => {
      getList()
      dialogVisible.value = false
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = false
    })
}

function edit() {
  formProps.value.disabled = true
  new ServerProcessEdit({
    ...formData.value,
    items: JSON.stringify(formData.value.items),
  })
    .request()
    .then(() => {
      getList()
      dialogVisible.value = false
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = false
    })
}

function id2server(serverId: number): ServerData {
  return serverOption.value.find((_) => _.id === serverId) || ({} as ServerData)
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
