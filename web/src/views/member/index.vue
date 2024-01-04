<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="userName"
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
          :permissions="[permission.AddMember]"
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
        <el-table-column prop="account" width="150" :label="$t('account')" />
        <el-table-column prop="name" :label="$t('name')" />
        <el-table-column
          prop="contact"
          :label="$t('contact')"
          show-overflow-tooltip
        />
        <el-table-column
          prop="superManager"
          :label="$t('admin')"
          width="70"
          align="center"
        >
          <template #default="scope">
            {{ $t(`boolOption[${scope.row['superManager'] || 0}]`) }}
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
          width="130"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <Button
              v-if="scope.row.id !== 1 && scope.row.id !== $store.state.user.id"
              type="primary"
              :icon="Edit"
              :permissions="[permission.EditMember]"
              @click="handleEdit(scope.row)"
            />
            <Button
              v-if="scope.row.id !== 1 && scope.row.id !== $store.state.user.id"
              type="danger"
              :icon="Delete"
              :permissions="[permission.DeleteMember]"
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
        <el-form-item
          v-show="$store.state.user.superManager === 1"
          :label="$t('admin')"
          prop="superManager"
        >
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
                :icon="QuestionFilled"
                style="margin-left: 10px; color: #666"
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
import permission from '@/permission'
import Button from '@/components/Permission/Button.vue'
import {
  Refresh,
  Plus,
  Edit,
  Delete,
  QuestionFilled,
} from '@element-plus/icons-vue'
import { validUsername, validPassword } from '@/utils/validate'
import { UserData, UserList, UserAdd, UserEdit, UserRemove } from '@/api/user'
import type { ElForm } from 'element-plus'
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const dialogVisible = ref(false)
const userName = ref('')
const tableLoading = ref(false)
const tableData = ref<UserList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const form = ref<InstanceType<typeof ElForm>>()
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
const formRules: InstanceType<typeof ElForm>['rules'] = {
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
    },
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
    },
  ],
  name: [
    {
      required: true,
      message: 'Name required',
      trigger: 'blur',
    },
  ],
}

getList()

const tablePage = computed(() => {
  let _tableData = tableData.value
  if (userName.value !== '') {
    _tableData = tableData.value.filter(
      (item) => item.name.indexOf(userName.value) !== -1
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
  new UserList()
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function refresList() {
  userName.value = ''
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

function handleEdit(data: UserData) {
  restoreFormData()
  formData.value = Object.assign(formData.value, data)
  dialogVisible.value = true
}

function handleRemove(data: UserData) {
  ElMessageBox.confirm(t('removeTips', { name: data.name }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      new UserRemove(data).request().then(() => {
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
  new UserAdd(formData.value)
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
