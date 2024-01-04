<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="namespaceName"
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
          :permissions="[permission.AddNamespace]"
          @click="handleAdd"
        />
      </el-row>
    </el-row>
    <el-row class="app-table">
      <el-table
        :key="tableHeight"
        v-loading="tableLoading"
        highlight-current-row
        height="100%"
        :data="tablePage.list"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('name')" />
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
          prop="user"
          :label="$t('member')"
          width="90"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <el-button type="primary" link @click="handleUser(scope.row)">
              {{ $t('view') }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="100"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <Button
              type="primary"
              :icon="Edit"
              :permissions="[permission.EditNamespace]"
              @click="handleEdit(scope.row)"
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
    >
      <el-form
        ref="form"
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
import permission from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Refresh, Plus, Edit } from '@element-plus/icons-vue'
import {
  NamespaceList,
  NamespaceAdd,
  NamespaceEdit,
  NamespaceData,
} from '@/api/namespace'
import getTableHeight from '@/composables/tableHeight'
import type { ElForm } from 'element-plus'
import TheUserDialog from './components/TheUserDialog.vue'
import { ref, computed } from 'vue'

const { tableHeight } = getTableHeight()
const dialogVisible = ref(false)
const dialogUserVisible = ref(false)
const namespaceName = ref('')
const tableLoading = ref(false)
const tableData = ref<NamespaceList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const selectedItem = ref<NamespaceData>()
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = { id: 0, name: '' }
const formData = ref(tempFormData)
const formProps = ref({ disabled: false })

getList()

const tablePage = computed(() => {
  let _tableData = tableData.value
  if (namespaceName.value !== '') {
    _tableData = tableData.value.filter(
      (item) => item.name.indexOf(namespaceName.value) !== -1
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
  new NamespaceList()
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function refresList() {
  namespaceName.value = ''
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

function handleEdit(data: NamespaceData) {
  formData.value = Object.assign({}, data)
  dialogVisible.value = true
}

function handleUser(data: NamespaceData) {
  selectedItem.value = data
  dialogUserVisible.value = true
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
  new NamespaceAdd(formData.value)
    .request()
    .then(() => {
      getList()
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
