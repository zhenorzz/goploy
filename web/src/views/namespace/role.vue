<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="roleName"
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
          :permissions="[pms.AddRole]"
          @click="handleAdd"
        />
      </el-row>
    </el-row>
    <el-row class="app-table">
      <el-table
        v-loading="tableLoading"
        highlight-current-row
        height="100%"
        :data="tablePage.list"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('name')" />
        <el-table-column
          prop="description"
          :label="$t('description')"
          show-overflow-tooltip
        />
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
            <Button
              type="primary"
              :icon="Edit"
              :permissions="[pms.EditRole]"
              @click="handleEdit(scope.row)"
            />
            <Button
              :icon="Setting"
              :permissions="[pms.EditPermission]"
              @click="handlePermission(scope.row)"
            />
            <Button
              type="danger"
              :icon="Delete"
              :permissions="[pms.DeleteRole]"
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
        <el-form-item :label="$t('description')" prop="description">
          <el-input v-model="formData.description" autocomplete="off" />
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
    <el-dialog
      v-model="permissionDialogVisible"
      :title="$t('permission')"
      :close-on-click-modal="false"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <div v-loading="rolePermissionLoading" class="permission-dialog">
        <div v-for="item in permissionList" :key="item.id">
          <el-checkbox
            v-model="item.checked"
            :indeterminate="item.indeterminate"
            :label="$t(`rolePermission.${item.name}`)"
            @change="(val) => handleCheckAllChange(val, item.id)"
          />
          <el-checkbox-group
            v-model="item.checkedGroup"
            style="margin-left: 22px"
            @change="() => handleCheckedChange(item.id)"
          >
            <el-checkbox
              v-for="child in item.children"
              :key="child.id"
              :label="child.id"
            >
              {{ $t(`rolePermission.${child.name}`) }}
            </el-checkbox>
          </el-checkbox-group>
        </div>
      </div>
      <template #footer>
        <el-button @click="permissionDialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
        <el-button
          :disabled="changePermissionDisabled"
          type="primary"
          @click="changePermission"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </el-row>
</template>
<script lang="ts">
export default { name: 'NamespaceRole' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Refresh, Setting, Plus, Edit, Delete } from '@element-plus/icons-vue'
import {
  RoleList,
  RoleAdd,
  RoleEdit,
  RoleRemove,
  RolePermissionList,
  RoleData,
  RolePermissionBindings,
  RolePermissionChange,
} from '@/api/role'
import type { CheckboxValueType, ElForm } from 'element-plus'
import { deepClone } from '@/utils'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'

interface permission {
  id: number
  name: string
  checked: boolean
  indeterminate: boolean
  checkedGroup: Array<number>
  children: Array<{ id: number; name: string }>
}

const { t } = useI18n()
const permissionDialogVisible = ref(false)
const dialogVisible = ref(false)
const roleName = ref('')
const tableLoading = ref(false)
const tableData = ref<RoleList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
let selectedItem = {} as RoleData
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = { id: 0, name: '', description: '' }
const formData = ref(tempFormData)
const formProps = ref({ disabled: false })

const tempPermissionList = {} as Record<number, permission>
const permissionList = ref<Record<number, permission>>({})
const rolePermissionLoading = ref(false)
const changePermissionDisabled = ref(false)

getList()
getPermissionList()

const tablePage = computed(() => {
  let _tableData = tableData.value
  if (roleName.value !== '') {
    _tableData = tableData.value.filter(
      (item) => item.name.indexOf(roleName.value) !== -1
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
  new RoleList()
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function refresList() {
  roleName.value = ''
  pagination.value.page = 1
  getList()
}

function getPermissionList() {
  new RolePermissionList().request().then((response) => {
    response.data.list.forEach((item) => {
      if (item.pid === 0) {
        tempPermissionList[item.id] = {
          id: item.id,
          name: item.name,
          checked: false,
          checkedGroup: [],
          indeterminate: false,
          children: [],
        }
      } else {
        tempPermissionList[item.pid].children.push({
          id: item.id,
          name: item.name,
        })
      }
    })
  })
}

function handlePageChange(val = 1) {
  pagination.value.page = val
}

function handleAdd() {
  restoreFormData()
  dialogVisible.value = true
}

function handleEdit(data: RoleData) {
  formData.value = Object.assign({}, data)
  dialogVisible.value = true
}

function handlePermission(data: RoleData) {
  restorePermissionList()
  selectedItem = data
  changePermissionDisabled.value = rolePermissionLoading.value = true
  new RolePermissionBindings({ roleId: data.id }).request().then((response) => {
    response.data.list.forEach((item) => {
      for (const [_, p] of Object.entries(permissionList.value)) {
        if (p.children.findIndex((_) => _.id === item.permissionId) !== -1) {
          p.checkedGroup.push(item.permissionId)
        }
        handleCheckedChange(p.id)
      }
    })
    changePermissionDisabled.value = rolePermissionLoading.value = false
  })

  permissionDialogVisible.value = true
}

function handleCheckAllChange(val: CheckboxValueType, id: number) {
  permissionList.value[id].checkedGroup = val
    ? permissionList.value[id].children.map((_) => _.id)
    : []
  permissionList.value[id].indeterminate = false
}

function handleCheckedChange(id: number) {
  const checkedCount = permissionList.value[id].checkedGroup.length
  const childrenCount = permissionList.value[id].children.length
  permissionList.value[id].checked = checkedCount === childrenCount
  permissionList.value[id].indeterminate =
    checkedCount > 0 && checkedCount < childrenCount
}

function handleRemove(data: RoleData) {
  ElMessageBox.confirm(t('deleteTips', { name: data.name }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      new RoleRemove(data).request().then(() => {
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
  new RoleAdd(formData.value)
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
  new RoleEdit(formData.value)
    .request()
    .then(() => {
      getList()
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = dialogVisible.value = false
    })
}

function changePermission() {
  let permissionIds: Array<number> = []
  for (const [_, p] of Object.entries(permissionList.value)) {
    permissionIds = permissionIds.concat(p.checkedGroup)
  }

  if (permissionIds.length === 0) {
    ElMessage.warning('Please select permission item')
    return
  }

  rolePermissionLoading.value = true
  new RolePermissionChange({
    roleId: selectedItem.id,
    permissionIds: permissionIds,
  })
    .request()
    .then(() => {
      ElMessage.success('Success')
    })
    .finally(() => {
      rolePermissionLoading.value = permissionDialogVisible.value = false
    })
}

function restoreFormData() {
  formData.value = { ...tempFormData }
}

function restorePermissionList() {
  permissionList.value = deepClone(tempPermissionList)
}
</script>
<style lang="scss">
@import '@/styles/mixin.scss';
.permission-dialog {
  height: 520px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
