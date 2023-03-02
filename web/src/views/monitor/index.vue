<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="monitorName"
          style="width: 200px"
          placeholder="Filter the name"
        />
      </el-row>
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
          :permissions="[pms.AddMonitor]"
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
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column prop="name" :label="$t('name')" min-width="120" />
        <el-table-column
          prop="target"
          label="Domain"
          min-width="140"
          show-overflow-tooltip
        >
        </el-table-column>
        <el-table-column
          prop="second"
          :label="$t('interval') + '(s)'"
          width="95"
        />
        <el-table-column
          prop="times"
          :label="$t('monitorPage.failTimes')"
          width="115"
        />
        <el-table-column prop="notifyType" :label="$t('notice')" width="90">
          <template #default="scope">
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
          </template>
        </el-table-column>
        <el-table-column
          prop="state"
          :label="$t('state')"
          width="120"
          align="center"
        >
          <template #default="scope">
            {{ $t(`switchOption[${scope.row.state || 0}]`) }}
            <Switch
              :value="scope.row.state === 1"
              active-color="#13ce66"
              inactive-color="#ff4949"
              :permissions="[pms.EditMonitor]"
              @change="handleToggle(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column
          prop="errorContent"
          :label="$t('monitorPage.errorContent')"
          min-width="140"
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
          fixed="right"
        >
          <template #default="scope">
            <Button
              type="primary"
              :icon="Edit"
              :permissions="[pms.EditMonitor]"
              @click="handleEdit(scope.row)"
            />
            <Button
              type="danger"
              :icon="Delete"
              :permissions="[pms.DeleteMonitor]"
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
      :title="$t('setting')"
      :fullscreen="$store.state.app.device === 'mobile'"
      :close-on-click-modal="false"
    >
      <el-form
        ref="form"
        v-loading="formProps.loading"
        :class="$store.state.app.device === 'desktop' ? 'monitor-dialog' : ''"
        :model="formData"
        label-width="120px"
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
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item
          :label="$t('type')"
          prop="type"
          :rules="[
            { required: true, message: 'Type required', trigger: 'blur' },
          ]"
        >
          <el-select
            v-model="formData.type"
            style="width: 100%"
            @change="handleTypeChange"
          >
            <el-option :label="$t('monitorPage.typeOption[1]')" :value="1" />
            <el-option :label="$t('monitorPage.typeOption[2]')" :value="2" />
            <el-option :label="$t('monitorPage.typeOption[3]')" :value="3" />
            <el-option :label="$t('monitorPage.typeOption[4]')" :value="4" />
            <el-option :label="$t('monitorPage.typeOption[5]')" :value="5" />
          </el-select>
        </el-form-item>
        <template v-if="4 > formData.type && formData.type > 0">
          <el-form-item :label="$t('target')">
            <el-select
              v-model="formProps.items"
              style="width: 100%"
              allow-create
              filterable
              multiple
              default-first-option
              clearable
            >
            </el-select>
          </el-form-item>
        </template>
        <template v-else-if="formData.type === 4">
          <el-form-item :label="$t('target')">
            <el-select
              v-model="formProps.items"
              multiple
              filterable
              style="width: 100%"
            >
              <el-option
                v-for="(item, index) in formProps.serverOption"
                :key="index"
                :label="item.label"
                :value="item.id.toString()"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('process')">
            <el-input
              v-model="formProps.process"
              autocomplete="off"
              placeholder="The name within ps -ef"
            />
          </el-form-item>
        </template>
        <template v-else-if="formData.type === 5">
          <el-form-item :label="$t('target')">
            <el-select
              v-model="formProps.items"
              multiple
              filterable
              style="width: 100%"
            >
              <el-option
                v-for="(item, index) in formProps.serverOption"
                :key="index"
                :label="item.label"
                :value="item.id.toString()"
              />
            </el-select>
          </el-form-item>
          <el-form-item :label="$t('script')">
            <VAceEditor
              v-model:value="formProps.script"
              lang="sh"
              theme="github"
              style="height: 360px; width: 100%"
              :options="{ newLineMode: 'unix' }"
            />
          </el-form-item>
        </template>
        <el-form-item :label="$t('timeout') + '(s)'">
          <el-input
            v-model="formProps.timeout"
            autocomplete="off"
            placeholder=""
          />
        </el-form-item>
        <el-form-item
          :label="$t('interval') + '(s)'"
          prop="second"
          :rules="[
            {
              type: 'number',
              required: true,
              min: 1,
              message: 'Interval required',
              trigger: 'blur',
            },
          ]"
        >
          <el-radio-group v-model="formData.second">
            <el-radio :label="60">1 min</el-radio>
            <el-radio :label="300">5 min</el-radio>
            <el-radio :label="900">15 min</el-radio>
            <el-radio :label="1800">30 min</el-radio>
            <el-radio :label="3600">60 min</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item
          :label="$t('monitorPage.failTimes')"
          prop="times"
          :rules="[
            {
              type: 'number',
              required: true,
              min: 1,
              max: 65535,
              message: 'Times required',
              trigger: 'blur',
            },
          ]"
        >
          <el-radio-group v-model="formData.times">
            <el-radio :label="1">1</el-radio>
            <el-radio :label="2">2</el-radio>
            <el-radio :label="3">3</el-radio>
            <el-radio :label="4">4</el-radio>
            <el-radio :label="5">5</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('monitorPage.silentCycle')">
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
          <el-row type="flex" style="width: 100%">
            <el-select v-model="formData.notifyType">
              <el-option :label="$t('webhookOption[1]')" :value="1" />
              <el-option :label="$t('webhookOption[2]')" :value="2" />
              <el-option :label="$t('webhookOption[3]')" :value="3" />
              <el-option :label="$t('webhookOption[255]')" :value="255" />
            </el-select>
            <el-input
              v-model.trim="formData.notifyTarget"
              style="flex: 1"
              autocomplete="off"
              placeholder="webhook"
            />
          </el-row>
        </el-form-item>
        <el-form-item
          :label="$t('description')"
          prop="description"
          :rules="[
            { max: 255, message: 'Max 255 characters', trigger: 'blur' },
          ]"
        >
          <el-input
            v-model="formData.description"
            type="textarea"
            :autosize="{ minRows: 2 }"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-row type="flex" justify="space-between">
          <el-button :loading="formProps.loading" type="success" @click="check">
            {{ $t('monitorPage.testState') }}
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
  </el-row>
</template>
<script lang="ts">
export default { name: 'MonitorIndex' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import { Button, Switch } from '@/components/Permission'
import { Refresh, Plus, Edit, Delete } from '@element-plus/icons-vue'
import { VAceEditor } from 'vue3-ace-editor'
import * as ace from 'ace-builds/src-noconflict/ace'
import { ServerOption } from '@/api/server'
import {
  MonitorList,
  MonitorAdd,
  MonitorEdit,
  MonitorCheck,
  MonitorToggle,
  MonitorRemove,
  MonitorData,
} from '@/api/monitor'
import type { ElForm } from 'element-plus'
import { ElMessageBox, ElMessage } from 'element-plus'
import { ref, watch, computed } from 'vue'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const store = useStore()
ace.config.set(
  'basePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
ace.config.set(
  'themePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
const dialogVisible = ref(false)
const monitorName = ref('')
const tableLoading = ref(false)
const tableData = ref<MonitorList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = {
  id: 0,
  name: '',
  type: 1,
  target: '',
  second: 60,
  times: 1,
  silentCycle: 1440,
  notifyType: 1,
  notifyTarget: '',
  description: '',
}
const formData = ref(tempFormData)
const formProps = ref({
  loading: false,
  disabled: false,
  serverLoading: false,
  serverOption: [] as ServerOption['datagram']['list'],
  items: [],
  timeout: 10,
  process: '',
  script: '',
})

watch(
  () => store.state.websocket.message,
  function (response) {
    if (response.type !== 3) {
      return
    }
    const data = response.message
    const monitorIndex = tableData.value.findIndex(
      (element) => element.id === data.monitorId
    )
    if (monitorIndex !== -1) {
      tableData.value[monitorIndex].errorContent = data.errorContent
      tableData.value[monitorIndex].state = data.state
    }
  }
)

getList()

const tablePage = computed(() => {
  let _tableData = tableData.value
  if (monitorName.value !== '') {
    _tableData = tableData.value.filter(
      (item) => item.name.indexOf(monitorName.value) !== -1
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
  new MonitorList()
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function refresList() {
  monitorName.value = ''
  pagination.value.page = 1
  getList()
}

function handlePageChange(val = 1) {
  pagination.value.page = val
}

function handleTypeChange(type: number) {
  if (type > 3) {
    formProps.value.serverLoading = true
    new ServerOption()
      .request()
      .then((response) => {
        formProps.value.serverOption = response.data.list
      })
      .finally(() => {
        formProps.value.serverLoading = false
      })
  }
}

function handleAdd() {
  restoreFormData()
  formProps.value.items = []
  dialogVisible.value = true
}

function handleEdit(data: MonitorData) {
  formData.value = Object.assign({}, data)
  formProps.value = Object.assign(formProps.value, JSON.parse(data.target))
  handleTypeChange(formData.value.type)
  dialogVisible.value = true
}

function handleToggle(data: MonitorData) {
  if (data.state === 1) {
    ElMessageBox.confirm(
      t('monitorPage.toggleStateTips', {
        monitorName: data.name,
      }),
      t('tips'),
      {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning',
      }
    )
      .then(() => {
        new MonitorToggle({ id: data.id, state: 0 }).request().then(() => {
          ElMessage.success(t('close'))
          getList()
        })
      })
      .catch(() => {
        ElMessage.info('Cancel')
      })
  } else {
    new MonitorToggle({ id: data.id, state: 1 }).request().then(() => {
      ElMessage.success(t('open'))
      getList()
    })
  }
}

function handleRemove(data: MonitorData) {
  ElMessageBox.confirm(
    t('monitorPage.removeMontiorTips', {
      monitorName: data.name,
    }),
    t('tips'),
    {
      confirmButtonText: t('confirm'),
      cancelButtonText: t('cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      new MonitorRemove({ id: data.id }).request().then(() => {
        ElMessage.success('Success')
        getList()
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function check() {
  form.value?.validate((valid) => {
    if (valid) {
      formProps.value.loading = true
      formProps.value.disabled = true
      new MonitorCheck({
        type: formData.value.type,
        items: formProps.value.items,
        timeout: formProps.value.timeout,
        process: formProps.value.process,
        script: formProps.value.script,
      })
        .request()
        .then(() => {
          ElMessage.success(t('pass'))
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
    if (formProps.value.items.length === 0) {
      ElMessage.error('Target at least one item')
      return
    }
    if (4 > formData.value.type && formData.value.type > 0) {
      formData.value.target = JSON.stringify({
        items: formProps.value.items,
        timeout: formProps.value.timeout || 0,
      })
    } else if (formData.value.type === 4) {
      if (formProps.value.process.length === 0) {
        ElMessage.error('Process empty')
        return
      }
      formData.value.target = JSON.stringify({
        items: formProps.value.items,
        timeout: formProps.value.timeout || 0,
        process: formProps.value.process,
      })
    } else if (formData.value.type === 5) {
      if (formProps.value.script.length === 0) {
        ElMessage.error('Script empty')
        return
      }
      formData.value.target = JSON.stringify({
        items: formProps.value.items,
        timeout: formProps.value.timeout || 0,
        script: formProps.value.script,
      })
    }

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
  new MonitorAdd(formData.value)
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
  new MonitorEdit(formData.value)
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
.monitor-dialog {
  padding-right: 10px;
  max-height: 50vh;
  overflow-y: auto;
  @include scrollBar();
}
</style>
