<template>
  <el-row class="app-container">
    <el-row
      v-show="$store.state.user.superManager"
      class="app-bar"
      type="flex"
      justify="end"
    >
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
    </el-row>
    <el-table
      :key="tableHeight"
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      :max-height="tableHeight"
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" :label="$t('name')" />
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
        prop="user"
        :label="$t('member')"
        width="80"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button type="text" @click="handleUser(scope.row)">
            {{ $t('view') }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="80"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button
            type="primary"
            icon="el-icon-edit"
            @click="handleEdit(scope.row)"
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
    >
      <el-form
        ref="form"
        :rules="formRules"
        :model="formData"
        label-width="80px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button
          :disabled="formProps.disabled"
          type="primary"
          @click="submit"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
    <TheUserDialog
      v-model="dialogUserVisible"
      :namespace-id="selectedItem?.id"
    />
  </el-row>
</template>
<script lang="ts">
export default { name: 'NamespaceIndex' }
</script>
<script lang="ts" setup>
import {
  NamespaceList,
  NamespaceTotal,
  NamespaceAdd,
  NamespaceEdit,
  NamespaceData,
} from '@/api/namespace'
import getTableHeight from '@/composables/tableHeight'
import { ElMessage } from 'element-plus'
import Validator from 'async-validator'
import TheUserDialog from './components/TheUserDialog.vue'
import { ref } from 'vue'

const { tableHeight } = getTableHeight()
const dialogVisible = ref(false)
const dialogUserVisible = ref(false)
const tableLoading = ref(false)
const tableData = ref<NamespaceList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 17, total: 0 })
const selectedItem = ref<NamespaceData['datagram']>()
const form = ref<Validator>()
const tempFormData = { id: 0, name: '' }
const formData = ref(tempFormData)
const formProps = ref({ disabled: false })
const formRules = {
  name: [{ required: true, message: 'Name required', trigger: 'blur' }],
}
getList()
getTotal()

function getList() {
  tableLoading.value = true
  new NamespaceList(pagination.value)
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function getTotal() {
  new NamespaceTotal().request().then((response) => {
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

function handleEdit(data: NamespaceData['datagram']) {
  formData.value = Object.assign({}, data)
  dialogVisible.value = true
}

function handleUser(data: NamespaceData['datagram']) {
  selectedItem.value = data
  dialogUserVisible.value = true
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
  new NamespaceAdd(formData.value)
    .request()
    .then(() => {
      getList()
      getTotal()
    })
    .finally(() => {
      formProps.value.disabled = dialogVisible.value = false
    })
}

function edit() {
  formProps.value.disabled = true
  new NamespaceEdit(formData.value)
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
