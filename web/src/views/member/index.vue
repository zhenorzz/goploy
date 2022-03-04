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
      <el-table-column prop="account" width="130" :label="$t('account')" />
      <el-table-column prop="name" :label="$t('name')" />
      <el-table-column
        prop="contact"
        :label="$t('contact')"
        show-overflow-tooltip
      />
      <el-table-column
        prop="superManager"
        :label="$t('admin')"
        width="60"
        align="center"
      >
        <template #default="scope">
          {{ $t(`boolOption[${scope.row['superManager'] || 0}]`) }}
        </template>
      </el-table-column>
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
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button
            v-if="scope.row.id !== 1 && scope.row.id !== $store.getters.uid"
            type="primary"
            icon="el-icon-edit"
            @click="handleEdit(scope.row)"
          />
          <el-button
            v-if="scope.row.id !== 1 && scope.row.id !== $store.getters.uid"
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
        @current-change="handleCurrentChange"
      />
    </el-row>
    <el-dialog
      v-model="dialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('setting')"
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
        <el-form-item :label="$t('account')" prop="account">
          <el-input
            v-model="formData.account"
            autocomplete="off"
            :readonly="formData.id > 0"
          />
        </el-form-item>
        <el-form-item :label="$t('password')" prop="password">
          <el-input v-model="formData.password" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('contact')" prop="contact">
          <el-input v-model="formData.contact" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('admin')" prop="superManager">
          <el-radio-group v-model="formData.superManager">
            <el-radio :label="1">{{ $t('boolOption[1]') }}</el-radio>
            <el-radio :label="0">{{ $t('boolOption[0]') }}</el-radio>
          </el-radio-group>
          <el-popover
            placement="top-start"
            :title="$t('description')"
            width="300"
            trigger="hover"
          >
            {{ $t('memberPage.permissionDesc') }}
            <template #reference>
              <el-button
                type="text"
                icon="el-icon-question"
                style="color: #666"
              />
            </template>
          </el-popover>
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
  </el-row>
</template>

<script lang="ts">
export default { name: 'MemberIndex' }
</script>
<script lang="ts" setup>
import { validUsername, validPassword } from '@/utils/validate'
import {
  UserData,
  UserList,
  UserTotal,
  UserAdd,
  UserEdit,
  UserRemove,
} from '@/api/user'
import Validator, { RuleItem } from 'async-validator'
import { ElMessageBox, ElMessage } from 'element-plus'
import getTableHeight from '@/composables/tableHeight'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const { tableHeight } = getTableHeight()
const dialogVisible = ref(false)
const tableLoading = ref(false)
const tableData = ref<UserList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 18, total: 0 })
const form = ref<Validator>()
const tempFormData = {
  id: 0,
  account: '',
  password: '',
  name: '',
  contact: '',
  superManager: 0,
}
const formData = ref(tempFormData)
const formProps = ref({
  disabled: false,
})
const formRules = {
  account: [
    {
      trigger: 'blur',
      validator: (_, value) => {
        if (!validUsername(value)) {
          return new Error('Greater than 5 characters')
        } else {
          return true
        }
      },
    } as RuleItem,
  ],
  password: [
    {
      trigger: 'blur',
      validator: (_, value) => {
        if (!value) {
          return true
        } else if (!validPassword(value)) {
          return new Error(
            '8 to 16 characters and a minimum of 2 character sets from these classes: [letters], [numbers], [special characters]'
          )
        } else {
          return true
        }
      },
    } as RuleItem,
  ],
  name: [
    {
      required: true,
      message: 'Name required',
      trigger: 'blur',
    } as RuleItem,
  ],
}

getList()
getTotal()

function getList() {
  tableLoading.value = true
  new UserList(pagination.value)
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function getTotal() {
  new UserTotal().request().then((response) => {
    pagination.value.total = response.data.total
  })
}

function handleCurrentChange(val: number) {
  pagination.value.page = val
  getList()
}

function handleAdd() {
  restoreFormData()
  dialogVisible.value = true
}

function handleEdit(data: UserData['datagram']) {
  restoreFormData()
  formData.value = Object.assign(formData.value, data)
  dialogVisible.value = true
}

function handleRemove(data: UserData['datagram']) {
  ElMessageBox.confirm(
    t('memberPage.removeUserTips', { name: data.name }),
    t('tips'),
    {
      confirmButtonText: t('confirm'),
      cancelButtonText: t('cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      new UserRemove(data).request().then(() => {
        ElMessage.success('Success')
        getList()
        getTotal()
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
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
  new UserAdd(formData.value)
    .request()
    .then(() => {
      ElMessage.success('Success')
      getList()
      getTotal()
      dialogVisible.value = false
    })
    .finally(() => {
      formProps.value.disabled = false
    })
}

function edit() {
  formProps.value.disabled = true
  new UserEdit(formData.value)
    .request()
    .then(() => {
      ElMessage.success('Success')
      getList()
      dialogVisible.value = false
    })
    .finally(() => {
      formProps.value.disabled = false
    })
}

function restoreFormData() {
  formData.value = { ...tempFormData }
}
</script>
