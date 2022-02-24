<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
    </el-row>
    <el-table
      v-loading="tableLoading"
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
          <el-input v-model="formData.account" autocomplete="off" />
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
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'MemberIndex',
  data() {
    return {
      dialogVisible: false,
      tableLoading: false,
      tableData: [] as UserList['datagram']['list'],
      tempFormData: {},
      pagination: {
        page: 1,
        rows: 18,
        total: 0,
      },
      formProps: {
        disabled: false,
      },
      formData: {
        id: 0,
        account: '',
        password: '',
        name: '',
        contact: '',
        superManager: 0,
      },
      formRules: {
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
      },
    }
  },
  created() {
    this.storeFormData()
    this.getList()
    this.getTotal()
  },
  methods: {
    getList() {
      this.tableLoading = true
      new UserList(this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
    },

    getTotal() {
      new UserTotal().request().then((response) => {
        this.pagination.total = response.data.total
      })
    },

    handleCurrentChange(val: number) {
      this.pagination.page = val
      this.getList()
    },

    handleAdd() {
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleEdit(data: UserData['datagram']) {
      this.restoreFormData()
      this.formData = Object.assign(this.formData, data)
      this.dialogVisible = true
    },

    handleRemove(data: UserData['datagram']) {
      ElMessageBox.confirm(
        this.$t('memberPage.removeUserTips', { name: data.name }),
        this.$t('tips'),
        {
          confirmButtonText: this.$t('confirm'),
          cancelButtonText: this.$t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          new UserRemove(data).request().then(() => {
            ElMessage.success('Success')
            this.getList()
            this.getTotal()
          })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },

    submit() {
      ;(this.$refs.form as Validator).validate((valid: boolean) => {
        if (valid) {
          if (this.formData.id === 0) {
            this.add()
          } else {
            this.edit()
          }
        } else {
          return false
        }
      })
    },

    add() {
      this.formProps.disabled = true
      new UserAdd(this.formData)
        .request()
        .then(() => {
          ElMessage.success('Success')
          this.getList()
          this.getTotal()
          this.dialogVisible = false
        })
        .finally(() => {
          this.formProps.disabled = false
        })
    },

    edit() {
      this.formProps.disabled = true
      new UserEdit(this.formData)
        .request()
        .then(() => {
          ElMessage.success('Success')
          this.getList()
          this.dialogVisible = false
        })
        .finally(() => {
          this.formProps.disabled = false
        })
    },

    storeFormData() {
      this.tempFormData = JSON.parse(JSON.stringify(this.formData))
    },

    restoreFormData() {
      this.formData = JSON.parse(JSON.stringify(this.tempFormData))
    },
  },
})
</script>
