<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="serverName"
          style="width: 200px"
          placeholder="Filter the name"
        />
        <el-input
          v-model="serverHost"
          style="width: 200px"
          placeholder="Filter the host"
        />
      </el-row>
      <el-row>
        <Button
          type="primary"
          style="margin-right: 10px"
          :permissions="[pms.InstallAgent]"
          @click="handleInstallAgent"
        >
          {{ $t('serverPage.installAgent') }}
        </Button>
        <el-upload
          :action="uploadHref"
          accept=".csv"
          :show-file-list="false"
          :before-upload="beforeUpload"
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
        >
          <Button
            type="primary"
            :loading="uploading"
            :permissions="[pms.ImportCSV]"
          >
            {{ $t('serverPage.importCSV') }}
          </Button>
        </el-upload>
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
          :permissions="[pms.AddServer]"
          @click="handleAdd"
        />
      </el-row>
    </el-row>
    <el-table
      :key="tableHeight"
      v-loading="tableLoading"
      :max-height="tableHeight"
      border
      stripe
      highlight-current-row
      :data="tablePage.list"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="ID" width="100" />
      <el-table-column prop="name" :label="$t('name')" min-width="140" />
      <el-table-column prop="ip" label="Host" min-width="140" sortable>
        <template #default="scope">
          {{ scope.row.ip }}:{{ scope.row.port }}
        </template>
      </el-table-column>
      <el-table-column
        prop="owner"
        :label="$t('serverPage.sshKeyOwner')"
        width="140"
        show-overflow-tooltip
      />
      <el-table-column label="OS" min-width="100" show-overflow-tooltip>
        <template #default="scope">
          <svg-icon
            v-if="scope.row.osInfo !== ''"
            :icon-class="getOSIcon(scope.row.osInfo)"
          />
          {{ getOS(scope.row.osInfo) }}
        </template>
      </el-table-column>
      <el-table-column
        prop="osInfo"
        label="CPU/MEMORY"
        min-width="100"
        show-overflow-tooltip
      >
        <template #default="scope">
          {{ getOSDetail(scope.row.osInfo) }}
        </template>
      </el-table-column>
      <el-table-column
        prop="description"
        :label="$t('description')"
        min-width="140"
        show-overflow-tooltip
      />
      <el-table-column
        prop="state"
        :label="$t('state')"
        width="110"
        align="center"
      >
        <template #default="scope">
          {{ $t(`stateOption[${scope.row.state || 0}]`) }}
          <Switch
            :value="scope.row.state === 1"
            active-color="#13ce66"
            inactive-color="#ff4949"
            :permissions="[pms.SwitchServerState]"
            @change="(value) => onSwitchState(value as boolean, scope.$index)"
          />
        </template>
      </el-table-column>
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
        width="190"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <Button
            :icon="DataAnalysis"
            :permissions="[pms.ShowServerMonitorPage]"
            @click="handleMonitor(scope.row)"
          />
          <Button
            type="info"
            :icon="DocumentCopy"
            :permissions="[pms.AddServer]"
            @click="handleCopy(scope.row)"
          />
          <Button
            type="primary"
            :icon="Edit"
            :permissions="[pms.EditServer]"
            @click="handleEdit(scope.row)"
          />
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px; width: 100%">
      <el-pagination
        hide-on-single-page
        :total="tablePage.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
    <el-dialog
      v-model="dialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('setting')"
      :close-on-click-modal="false"
    >
      <el-form
        ref="form"
        v-loading="formProps.loading"
        :rules="formRules"
        :model="formData"
        label-width="130px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item :label="$t('namespace')" prop="namespaceId">
          <el-radio-group v-model="formData.namespaceId">
            <el-radio :label="getNamespace()['id']">
              {{ $t('current') }}
            </el-radio>
            <el-radio :label="0">{{ $t('unlimited') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Host" prop="ip">
          <el-input v-model="formData.ip" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Port" prop="port">
          <el-input v-model.number="formData.port" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('serverPage.sshKeyOwner')" prop="owner">
          <el-input
            v-model="formData.owner"
            autocomplete="off"
            placeholder="root"
          />
        </el-form-item>
        <el-form-item :label="$t('serverPage.sshKeyPath')" prop="path">
          <el-row type="flex" style="width: 100%">
            <el-input
              v-model="formData.path"
              style="flex: 1"
              autocomplete="off"
              placeholder="/root/.ssh/id_rsa"
            />
            <el-button
              :icon="DocumentCopy"
              type="success"
              :loading="formProps.copyPubLoading"
              :disabled="formData.path === ''"
              @click="getPublicKey"
            >
              {{ $t('serverPage.copyPub') }}
            </el-button>
          </el-row>
        </el-form-item>
        <!-- <el-form-item :label="$t('serverPage.sshKeyPassword')" prop="password">
          <el-input
            v-model="formData.password"
            autocomplete="off"
            placeholder=""
          />
        </el-form-item> -->
        <el-form-item :label="$t('description')" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :autosize="{ minRows: 2 }"
          />
        </el-form-item>
        <el-form-item label="">
          <el-button
            type="text"
            @click="formProps.showAdvance = !formProps.showAdvance"
          >
            {{ $t('serverPage.advance') }}
          </el-button>
        </el-form-item>
        <el-form-item v-show="formProps.showAdvance" label="Jump host">
          <el-input v-model="formData.jumpIP" autocomplete="off" />
        </el-form-item>
        <el-form-item v-show="formProps.showAdvance" label="Jump port">
          <el-input v-model.number="formData.jumpPort" autocomplete="off" />
        </el-form-item>
        <el-form-item
          v-show="formProps.showAdvance"
          :label="$t('serverPage.sshKeyOwner')"
        >
          <el-input
            v-model="formData.jumpOwner"
            autocomplete="off"
            placeholder="root"
          />
        </el-form-item>
        <el-form-item
          v-show="formProps.showAdvance"
          :label="$t('serverPage.sshKeyPath')"
        >
          <el-input
            v-model="formData.jumpPath"
            autocomplete="off"
            placeholder="/root/.ssh/id_rsa"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-row type="flex" justify="space-between">
          <el-button :loading="formProps.loading" type="success" @click="check">
            {{ $t('serverPage.testConnection') }}
          </el-button>
          <el-row>
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
          </el-row>
        </el-row>
      </template>
    </el-dialog>
    <el-dialog
      v-model="agentDialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('setting')"
      :close-on-click-modal="false"
    >
      <el-form
        ref="agentForm"
        :rules="agentFormRules"
        :model="agentFormData"
        label-width="130px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item label="Install path" prop="installPath">
          <el-input v-model="agentFormData.installPath" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Use" prop="tool">
          <el-radio v-model="agentFormData.tool" label="wget">wget</el-radio>
          <el-radio v-model="agentFormData.tool" label="curl">curl</el-radio>
        </el-form-item>
        <el-form-item label="Report URL" prop="reportURL">
          <el-input v-model="agentFormData.reportURL" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Turn on web" prop="webState">
          <el-radio v-model="agentFormData.webState" :label="1"> Yes </el-radio>
          <el-radio v-model="agentFormData.webState" :label="0"> No </el-radio>
        </el-form-item>
        <el-form-item
          v-show="agentFormData.webState === 1"
          label="Web port"
          prop="webPort"
        >
          <el-input v-model="agentFormData.webPort" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="agentDialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
        <el-button
          :disabled="agentFormProps.disabled"
          type="primary"
          @click="installAgent"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </el-row>
</template>

<script lang="ts">
export default { name: 'ServerIndex' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import { Button, Switch } from '@/components/Permission'
import {
  Refresh,
  Plus,
  Edit,
  DocumentCopy,
  DataAnalysis,
} from '@element-plus/icons-vue'
import type { ElForm } from 'element-plus'
import { getNamespace } from '@/utils/namespace'
import {
  ServerList,
  ServerPublicKey,
  ServerAdd,
  ServerEdit,
  ServerCheck,
  ServerToggle,
  ServerData,
  ServerInstallAgent,
} from '@/api/server'
import getTableHeight from '@/composables/tableHeight'
import { HttpResponse } from '@/api/types'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import { ref, computed } from 'vue'
import { copy, humanSize } from '@/utils'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const router = useRouter()
const dialogVisible = ref(false)
const agentDialogVisible = ref(false)
const serverName = ref('')
const serverHost = ref('')
const { tableHeight } = getTableHeight()
const tableLoading = ref(false)
const tableData = ref<ServerList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 16 })
const selectedItems = ref<ServerList['datagram']['list']>([])
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = {
  id: 0,
  namespaceId: getNamespace()['id'],
  name: '',
  ip: '',
  port: 22,
  owner: '',
  path: '/root/.ssh/id_rsa',
  password: '',
  jumpIP: '',
  jumpPort: 0,
  jumpOwner: '',
  jumpPath: '',
  jumpPassword: '',
  description: '',
}
const formData = ref(tempFormData)
const formProps = ref({
  loading: false,
  showAdvance: false,
  copyPubLoading: false,
  disabled: false,
})
const formRules = <InstanceType<typeof ElForm>['rules']>{
  namespaceId: [
    { required: true, message: 'Namespace required', trigger: 'blur' },
  ],
  name: [{ required: true, message: 'Name required', trigger: 'blur' }],
  ip: [{ required: true, message: 'IP required', trigger: 'blur' }],
  port: [
    {
      type: 'number',
      required: true,
      min: 0,
      max: 65535,
      message: '0 ~ 65535',
      trigger: 'blur',
    },
  ],
  owner: [
    {
      required: true,
      message: 'SSH-KEY owner required',
      trigger: 'blur',
    },
  ],
  path: [{ required: true, message: 'SSH-KEY path required', trigger: 'blur' }],
  description: [{ max: 255, message: 'Max 255 characters', trigger: 'blur' }],
}

const uploading = ref(false)
const uploadHref = computed(() => {
  return `${
    import.meta.env.VITE_APP_BASE_API
  }/server/import?${NamespaceKey}=${getNamespaceId()}`
})

const agentForm = ref<InstanceType<typeof ElForm>>()
const agentFormData = ref({
  installPath: '',
  reportURL: '',
  tool: 'wget',
  webState: 1,
  webPort: '',
})
const agentFormProps = ref({
  disabled: false,
})
const agentFormRules = <InstanceType<typeof ElForm>['rules']>{
  installPath: [
    { required: true, message: 'Install path required', trigger: 'blur' },
  ],
  reportURL: [
    { required: true, message: 'Report url required', trigger: 'blur' },
  ],
}

getList()

const tablePage = computed(() => {
  let _tableData = tableData.value
  if (serverName.value !== '') {
    _tableData = tableData.value.filter(
      (item) => item.name.indexOf(serverName.value) !== -1
    )
  }
  if (serverHost.value !== '') {
    _tableData = tableData.value.filter(
      (item) => item.ip.indexOf(serverHost.value) !== -1
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

function getList() {
  tableLoading.value = true
  new ServerList()
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function refresList() {
  serverName.value = ''
  serverHost.value = ''
  pagination.value.page = 1
  getList()
}

function getPublicKey() {
  formProps.value.copyPubLoading = true
  new ServerPublicKey({ path: formData.value.path })
    .request()
    .then((response) => {
      copy(response.data.key)
      ElMessage.success(t('serverPage.copyPubTips'))
    })
    .finally(() => {
      formProps.value.copyPubLoading = false
    })
}

function handleInstallAgent() {
  if (selectedItems.value.length === 0) {
    ElMessage.warning('Please, select the item')
  } else {
    agentDialogVisible.value = true
  }
}

function handlePageChange(val = 1) {
  pagination.value.page = val
  getList()
}

function handleSelectionChange(value: ServerData[]) {
  selectedItems.value = value
}

function beforeUpload() {
  uploading.value = true
  return true
}

function handleUploadSuccess(response: HttpResponse<string>) {
  if (response.code > 0) {
    ElMessage.error(`upload failed, detail: ${response.message}`)
  } else {
    ElMessage.success('Success')
  }
  getList()
  uploading.value = false
  return true
}

function handleUploadError(err: Error) {
  ElMessage.error(err.message)
  uploading.value = false
  return true
}

function handleAdd() {
  restoreFormData()
  dialogVisible.value = true
}

function handleEdit(data: ServerData) {
  formData.value = Object.assign({}, data)
  dialogVisible.value = true
}

function handleCopy(data: ServerData) {
  formData.value = Object.assign({}, data)
  formData.value.id = 0
  dialogVisible.value = true
}

function handleMonitor(data: ServerData) {
  router.push({ path: '/server/agent', query: { serverId: data.id } })
}

function onSwitchState(value: boolean, index: number) {
  const data = tableData.value[index]
  if (value) {
    new ServerToggle({ id: data.id, state: value ? 1 : 0 })
      .request()
      .then(() => {
        ElMessage.success('Need to bind project again')
        tableData.value[index].state = value ? 1 : 0
      })
  } else {
    ElMessageBox.confirm(
      t('serverPage.removeServerTips', { serverName: data.name }),
      t('tips'),
      {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning',
      }
    )
      .then(() => {
        new ServerToggle({ id: data.id, state: value ? 1 : 0 })
          .request()
          .then(() => {
            tableData.value[index].state = value ? 1 : 0
          })
      })
      .catch(() => {
        ElMessage.info('Cancel')
      })
  }
}

function check() {
  form.value?.validate((valid) => {
    if (valid) {
      formProps.value.loading = true
      formProps.value.disabled = true
      new ServerCheck(formData.value)
        .request()
        .then(() => {
          ElMessage.success('Connected')
        })
        .finally(() => {
          formProps.value.loading = false
          formProps.value.disabled = false
        })
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
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
  new ServerAdd(formData.value)
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
  new ServerEdit(formData.value)
    .request()
    .then(() => {
      getList()
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = dialogVisible.value = false
    })
}

function installAgent() {
  agentForm.value?.validate((valid) => {
    if (valid) {
      agentFormProps.value.disabled = true
      if (agentFormData.value.webState === 0) {
        agentFormData.value.webPort = ''
      }

      new ServerInstallAgent({
        ids: selectedItems.value.map((_) => _.id),
        ...agentFormData.value,
      })
        .request()
        .then(() => {
          ElMessage.warning(t('serverPage.installAgentTips'))
          agentDialogVisible.value = false
        })
        .finally(() => {
          agentFormProps.value.disabled = false
        })
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
}

function restoreFormData() {
  formData.value = { ...tempFormData }
}

function getOS(osInfo: string): string {
  if (osInfo === '') return ''
  return osInfo.split('|')[0]
}

function getOSIcon(osInfo: string): string {
  if (osInfo === '') return ''
  else if (osInfo.toLowerCase().includes('centos')) return 'centos'
  else if (osInfo.toLowerCase().includes('ubuntu')) return 'ubuntu'
  else return 'question-mark-blue'
}

function getOSDetail(osInfo: string): string {
  if (osInfo === '') return ''
  const osArr = osInfo.split('|')
  return osArr[1] + ' cores ' + humanSize(Number(osArr[2]) * 1024)
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
