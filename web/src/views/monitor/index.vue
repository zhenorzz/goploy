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
        border
        stripe
        highlight-current-row
        :data="tablePage.list"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="100" />
        <el-table-column prop="name" :label="$t('name')" min-width="140" />
        <el-table-column prop="domain" label="Domain" min-width="100">
          <template #default="scope">
            {{ scope.row.domain }}:{{ scope.row.port }}
          </template>
        </el-table-column>
        <el-table-column
          prop="second"
          :label="$t('interval') + '(s)'"
          width="80"
        />
        <el-table-column
          prop="times"
          :label="$t('monitorPage.failTimes')"
          width="110"
        />
        <el-table-column prop="notifyType" :label="$t('notice')" width="70">
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
          prop="notifyTimes"
          :label="$t('monitorPage.notifyTimes')"
          width="85"
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
    <el-row type="flex" justify="end" style="margin-top: 10px; width: 100%">
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
        :rules="formRules"
        :model="formData"
        label-width="120px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item prop="url">
          <template #label>
            <el-tooltip placement="top">
              <template #content>scheme:opaque[?query][#fragment]</template>
              <el-button type="text">URL</el-button>
            </el-tooltip>
          </template>
          <el-input v-model="formData.url" autocomplete="off" placeholder="" />
        </el-form-item>
        <el-form-item :label="$t('interval') + '(s)'" prop="second">
          <el-input v-model.number="formData.second" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('monitorPage.failTimes')" prop="times">
          <el-input v-model.number="formData.times" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('notice')" prop="notifyTarget">
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
        <el-form-item :label="$t('monitorPage.notifyTimes')" prop="notifyTimes">
          <el-input v-model.number="formData.notifyTimes" />
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
        <el-row type="flex" justify="space-between">
          <el-button :loading="formProps.loading" type="success" @click="check">
            {{ $t('monitorPage.testAppState') }}
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
const dialogVisible = ref(false)
const monitorName = ref('')
const tableLoading = ref(false)
const tableData = ref<MonitorList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = {
  id: 0,
  name: '',
  url: '',
  second: 3,
  times: 1,
  notifyType: 1,
  notifyTarget: '',
  notifyTimes: 1,
  description: '',
}
const formData = ref(tempFormData)
const formProps = ref({
  loading: false,
  disabled: false,
})
const formRules = <InstanceType<typeof ElForm>['rules']>{
  name: [{ required: true, message: 'Name required', trigger: 'blur' }],
  url: [{ required: true, message: 'URL required', trigger: 'blur' }],
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
  second: [
    {
      type: 'number',
      required: true,
      min: 1,
      message: 'Interval required',
      trigger: 'blur',
    },
  ],
  times: [
    {
      type: 'number',
      required: true,
      min: 1,
      max: 65535,
      message: 'Times required',
      trigger: 'blur',
    },
  ],
  notifyTarget: [{ required: true, message: 'Webhook required' }],
  notifyTimes: [
    {
      type: 'number',
      required: true,
      min: 1,
      max: 65535,
      message: 'Notify times required',
      trigger: 'blur',
    },
  ],
  description: [{ max: 255, message: 'Max 255 characters', trigger: 'blur' }],
}

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

function handleAdd() {
  restoreFormData()
  dialogVisible.value = true
}

function handleEdit(data: MonitorData) {
  formData.value = Object.assign({}, data)
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
        new MonitorToggle({ id: data.id }).request().then(() => {
          ElMessage.success('Stop success')
          getList()
        })
      })
      .catch(() => {
        ElMessage.info('Cancel')
      })
  } else {
    new MonitorToggle({ id: data.id }).request().then(() => {
      ElMessage.success('Open success')
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
      new MonitorCheck(formData.value)
        .request()
        .then(() => {
          ElMessage.success('Connected success')
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
.template-dialog {
  padding-right: 10px;
  height: 400px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
