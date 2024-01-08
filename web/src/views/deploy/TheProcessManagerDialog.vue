<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('deployPage.processManager')"
    :fullscreen="$store.state.app.device === 'mobile'"
    :close-on-click-modal="false"
  >
    <el-row type="flex">
      <el-select
        v-model="projectProcessId"
        v-loading="processLoading"
        filterable
        style="flex: 1"
        @change="handleProcessChange"
      >
        <el-option
          v-for="item in processOption"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        >
          <el-row type="flex" justify="space-between">
            <span style="">{{ item.name }}</span>
            <el-row>
              <el-button
                style=""
                type="primary"
                text
                :icon="Edit"
                @click.stop="handleEdit(item)"
              />
              <el-button
                type="primary"
                text
                :icon="Delete"
                @click.stop="handleDelete(item.id)"
              />
            </el-row>
          </el-row>
        </el-option>
      </el-select>
      <el-button type="primary" :icon="Plus" @click="handleAdd" />
    </el-row>
    <el-table
      ref="table"
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      max-height="447px"
      style="margin-top: 10px; width: 100%"
      :data="tableData"
    >
      <el-table-column type="expand">
        <template #default="{}">
          <el-row
            v-if="commandRes.hasOwnProperty('execRes')"
            style="padding: 8px 16px; flex-direction: column; line-height: 20px"
          >
            <el-row>
              {{ $t('deployPage.execRes') }}:
              <span
                :class="commandRes.execRes ? 'exec-success' : 'exec-fail'"
                style="padding-left: 5px"
              >
                {{ commandRes.execRes ? $t('success') : $t('fail') }}
              </span>
            </el-row>
            <el-row style="white-space: pre-wrap">
              {{ $t('deployPage.execTime') }}: {{ commandRes.startTime }} -
              {{ commandRes.endTime }}
            </el-row>
            <el-row style="white-space: pre-wrap">
              {{ commandRes.stdout }}
            </el-row>
            <el-row style="white-space: pre-wrap">
              {{ commandRes.stderr }}
            </el-row>
          </el-row>
          <el-row v-else style="padding: 0 8px"> {{ 'not run' }} </el-row>
        </template>
      </el-table-column>
      <el-table-column label="Server">
        <template #default="scope"> {{ scope.row.server.name }} </template>
      </el-table-column>
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="400"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button
            :loading="commandLoading"
            type="primary"
            @click="handleProcessCmd(scope.row, 'status')"
          >
            status
          </el-button>
          <el-button
            :loading="commandLoading"
            type="success"
            @click="handleProcessCmd(scope.row, 'start')"
          >
            start
          </el-button>
          <el-button
            :loading="commandLoading"
            type="warning"
            @click="handleProcessCmd(scope.row, 'restart')"
          >
            restart
          </el-button>
          <el-button
            :loading="commandLoading"
            type="danger"
            @click="handleProcessCmd(scope.row, 'stop')"
          >
            stop
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <template #footer>
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
  <el-dialog
    v-model="processVisible"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <template #header>
      {{ $t('deployPage.processManager') }}
      <el-popover
        placement="bottom-start"
        :title="$t('projectPage.predefinedVar')"
        width="400"
        trigger="hover"
      >
        <div>
          <el-row>
            <span>${PROJECT_NAME}：</span>
            <span>project.name</span>
          </el-row>
          <el-row>
            <span>${PROJECT_PATH}：</span>
            <span>project.path</span>
          </el-row>
          <el-row>
            <span>${PROJECT_SYMLINK_PATH}：</span>
            <span>project.symlink_path</span>
          </el-row>
        </div>
        <template #reference>
          <el-button type="primary" :text="true">
            {{ $t('projectPage.predefinedVar') }}
          </el-button>
        </template>
      </el-popover>
    </template>
    <el-form
      ref="form"
      :model="formData"
      label-width="130px"
      :label-position="$store.state.app.device === 'desktop' ? 'right' : 'top'"
    >
      <el-form-item :label="$t('name')" prop="name" required>
        <el-input v-model="formData.name" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Status">
        <el-input v-model="formData.status" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Start">
        <el-input v-model="formData.start" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Stop">
        <el-input v-model="formData.stop" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Restart">
        <el-input v-model="formData.restart" autocomplete="off" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="processVisible = false">
        {{ $t('cancel') }}
      </el-button>
      <el-button :disabled="formProps.disabled" type="primary" @click="submit">
        {{ $t('confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>
<script lang="ts" setup>
import { Edit, Delete, Plus } from '@element-plus/icons-vue'
import { ManageProcess } from '@/api/deploy'
import {
  ProjectProcessData,
  ProjectProcessList,
  ProjectProcessAdd,
  ProjectProcessEdit,
  ProjectServerData,
  ProjectServerList,
  ProjectProcessDelete,
  ProjectData,
} from '@/api/project'
import type { ElForm } from 'element-plus'
import { computed, watch, ref, PropType } from 'vue'
import { useI18n } from 'vue-i18n'
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  projectRow: {
    type: Object as PropType<ProjectData>,
    required: true,
  },
})
const emit = defineEmits(['update:modelValue'])
const { t } = useI18n()
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val)
  },
})

watch(
  () => props.modelValue,
  (val: typeof props['modelValue']) => {
    if (val === true) {
      getList()
    }
  }
)

const processLoading = ref(false)
const projectProcessId = ref<number>()
const processOption = ref<ProjectProcessList['datagram']['list']>([])
const getList = () => {
  tableData.value = []
  const _processId = localStorage.getItem(
    `${props.projectRow.id}-latest-use-process`
  )
  projectProcessId.value = undefined
  if (_processId) {
    projectProcessId.value = Number(_processId)
    handleProcessChange(projectProcessId.value)
  }
  processLoading.value = true
  processOption.value = []
  new ProjectProcessList(
    { projectId: props.projectRow.id },
    { page: 1, rows: 999 }
  )
    .request()
    .then((response) => {
      processOption.value = response.data.list
    })
    .finally(() => {
      processLoading.value = false
    })
}

const table = ref()
const tableLoading = ref(false)
const tableData = ref<ProjectServerList['datagram']['list']>([])
const handleProcessChange = (processId: number) => {
  localStorage.setItem(
    `${props.projectRow.id}-latest-use-process`,
    processId.toString()
  )
  if (tableData.value.length > 0) {
    return
  }
  tableLoading.value = true
  new ProjectServerList({ id: props.projectRow.id })
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}
const commandRes = ref<ManageProcess['datagram']>(
  {} as ManageProcess['datagram']
)
const commandLoading = ref(false)
const handleProcessCmd = (data: ProjectServerData, command: string) => {
  ElMessageBox.confirm(t('deployPage.execTips', { command }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      commandLoading.value = true
      new ManageProcess({
        serverId: data.serverId,
        projectProcessId: Number(projectProcessId.value),
        command,
      })
        .request()
        .then((response) => {
          commandRes.value = response.data
          table.value.toggleRowExpansion(data, true)
        })
        .finally(() => {
          commandLoading.value = false
        })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

const processVisible = ref(false)
const form = ref<InstanceType<typeof ElForm>>()
const formData = ref({
  id: 0,
  projectId: 0,
  name: '',
  status: '',
  start: '',
  stop: '',
  restart: '',
})
const formProps = ref({
  disabled: false,
})

watch(
  () => props.projectRow,
  (val) => {
    formData.value.projectId = val.id
  }
)

function handleAdd() {
  processVisible.value = true
  formData.value.id = 0
}
function handleEdit(data: ProjectProcessData) {
  processVisible.value = true
  formData.value.id = data.id
  formData.value.name = data.name
  formData.value.status = data.status
  formData.value.start = data.start
  formData.value.stop = data.stop
  formData.value.restart = data.restart
}
function handleDelete(id: number) {
  ElMessageBox.confirm(t('deployPage.deleteProcessTips'), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      new ProjectProcessDelete({ id }).request().then(() => {
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
  new ProjectProcessAdd(formData.value)
    .request()
    .then(() => {
      processVisible.value = false
      ElMessage.success('Success')
      getList()
    })
    .finally(() => {
      formProps.value.disabled = false
    })
}

function edit() {
  formProps.value.disabled = true
  new ProjectProcessEdit(formData.value)
    .request()
    .then(() => {
      processVisible.value = false
      ElMessage.success('Success')
      getList()
    })
    .finally(() => {
      formProps.value.disabled = false
    })
}
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
.exec-success {
  color: #67c23a;
}

.exec-fail {
  color: #f56c6c;
}
</style>
