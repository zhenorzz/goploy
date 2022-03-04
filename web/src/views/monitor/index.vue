<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
    </el-row>
    <el-table
      :key="tableHeight"
      v-loading="tableLoading"
      :max-height="tableHeight"
      border
      stripe
      highlight-current-row
      :data="tableData"
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
        width="100"
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
        width="80"
      />
      <el-table-column
        prop="state"
        :label="$t('state')"
        width="65"
        align="center"
      >
        <template #default="scope">
          <el-switch
            :value="scope.row.state === 1"
            active-color="#13ce66"
            inactive-color="#ff4949"
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
        width="130"
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
            @click="handleRemove(scope.row)"
          />
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px; width: 100%">
      <el-pagination
        hide-on-single-page
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
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
            URL
            <el-tooltip placement="top">
              <template #content>scheme:opaque[?query][#fragment]</template>
              <i class="el-icon-question" />
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
          <el-row type="flex">
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
import {
  MonitorList,
  MonitorTotal,
  MonitorAdd,
  MonitorEdit,
  MonitorCheck,
  MonitorToggle,
  MonitorRemove,
  MonitorData,
} from '@/api/monitor'
import getTableHeight from '@/composables/tableHeight'
import Validator from 'async-validator'
import { ElMessageBox, ElMessage } from 'element-plus'
import { ref, watch } from 'vue'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const store = useStore()
const { tableHeight } = getTableHeight()
const dialogVisible = ref(false)
const tableLoading = ref(false)
const tableData = ref<MonitorList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 16, total: 0 })
const form = ref<Validator>()
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
const formRules = {
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
  () => store.getters.ws_message,
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
getTotal()
function getList() {
  tableLoading.value = true
  new MonitorList(pagination.value)
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function getTotal() {
  new MonitorTotal().request().then((response) => {
    pagination.value.total = response.data.total
  })
}

function handlePageChange(val = 1) {
  pagination.value.page = val
  getList()
}

function handleAdd() {
  restoreFormData()
  dialogVisible.value = true
}

function handleEdit(data: MonitorData['datagram']) {
  formData.value = Object.assign({}, data)
  dialogVisible.value = true
}

function handleToggle(data: MonitorData['datagram']) {
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

function handleRemove(data: MonitorData['datagram']) {
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
        getTotal()
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function check() {
  form.value?.validate((valid: boolean) => {
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
    } else {
      return false
    }
  })
}

function submit() {
  form.value?.validate((valid: boolean) => {
    if (valid) {
      if (formData.value.id === 0) {
        add()
      } else {
        edit()
      }
    } else {
      return false
    }
  })
}

function add() {
  formProps.value.disabled = true
  new MonitorAdd(formData.value)
    .request()
    .then(() => {
      getList()
      getTotal()
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
