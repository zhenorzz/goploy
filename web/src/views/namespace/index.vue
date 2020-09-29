<template>
  <el-row class="app-container">
    <el-row v-show="$store.getters.superManager" class="app-bar" type="flex" justify="end">
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
      <el-table-column prop="id" label="ID" width="160" />
      <el-table-column prop="name" :label="$t('name')" />
      <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
      <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
      <el-table-column prop="user" :label="$t('member')" width="80" align="center">
        <template slot-scope="scope">
          <el-button type="text" @click="handleUser(scope.row)">{{ $t('view') }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="operation" :label="$t('op')" width="80" align="center">
        <template slot-scope="scope">
          <el-button type="primary" icon="el-icon-edit" @click="handleEdit(scope.row)" />
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px;">
      <el-pagination
        hide-on-single-page
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
    <el-dialog :title="$t('setting')" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="80px">
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('manage')" :visible.sync="dialogUserVisible">
      <el-row class="app-bar" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="handleAddUser" />
      </el-row>
      <el-table
        border
        stripe
        highlight-current-row
        :data="tableUserData.filter(row => row.role!== $global.Admin)"
        style="width: 100%"
      >
        <el-table-column prop="userId" :label="$t('userId')" />
        <el-table-column prop="userName" :label="$t('userName')" />
        <el-table-column prop="role" :label="$t('role')" />
        <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
        <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
        <el-table-column prop="operation" :label="$t('op')" width="80" align="center">
          <template slot-scope="scope">
            <el-button type="danger" icon="el-icon-delete" @click="removeUser(scope.row)" />
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogUserVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('add')" :visible.sync="dialogAddUserVisible">
      <el-form ref="addUserForm" :rules="addUserFormRules" :model="addUserFormData">
        <el-form-item :label="$t('user')" label-width="120px" prop="userIds">
          <el-select v-model="addUserFormData.userIds" multiple clearable filterable>
            <el-option
              v-for="(item, index) in userOption.filter(item => item.superManager !== 1)"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('role')" label-width="120px" prop="role">
          <el-select v-model="addUserFormData.role">
            <el-option
              v-for="(role, index) in [$global.Manager, $global.GroupManager, $global.Member]"
              :key="index"
              :label="role"
              :value="role"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogAddUserVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="addUserFormProps.disabled" type="primary" @click="addUser">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { getList, getTotal, getBindUserList, add, edit, addUser, removeUser } from '@/api/namespace'
import { getOption as getUserOption } from '@/api/user'
import tableHeight from '@/mixin/tableHeight'

export default {
  name: 'Namespace',
  mixins: [tableHeight],
  data() {
    return {
      dialogVisible: false,
      dialogUserVisible: false,
      dialogAddUserVisible: false,
      userOption: [],
      tableLoading: false,
      tableData: [],
      tableUserData: [],
      pagination: {
        page: 1,
        rows: 16,
        total: 0
      },
      tempFormData: {},
      formProps: {
        disabled: false
      },
      formData: {
        id: 0
      },
      formRules: {
        name: [
          { required: true, message: 'Name required', trigger: 'blur' }
        ]
      },
      addUserFormProps: {
        disabled: false
      },
      addUserFormData: {
        namespaceId: 0,
        userIds: [],
        role: ''
      },
      addUserFormRules: {
        userIds: [
          { type: 'array', required: true, message: 'User required', trigger: 'change' }
        ],
        role: [
          { required: true, message: 'Role required', trigger: 'change' }
        ]
      }
    }
  },
  created() {
    this.storeFormData()
    this.getList()
    this.getTotal()
    this.getUserOption()
  },
  methods: {
    getList() {
      this.tableLoading = true
      getList(this.pagination).then((response) => {
        this.tableData = response.data.list
      }).finally(() => {
        this.tableLoading = false
      })
    },

    getTotal() {
      getTotal().then((response) => {
        this.pagination.total = response.data.total
      })
    },

    getBindUserList(namespaceID) {
      getBindUserList(namespaceID).then((response) => {
        this.tableUserData = response.data.list
      })
    },

    getUserOption() {
      getUserOption().then((response) => {
        this.userOption = response.data.list
      })
    },

    handlePageChange(val) {
      this.pagination.page = val
      this.getList()
    },

    handleAdd() {
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleEdit(data) {
      this.formData = Object.assign({}, data)
      this.dialogVisible = true
    },

    handleUser(data) {
      this.getBindUserList(data.id)
      this.addUserFormData.namespaceId = data.id
      this.dialogUserVisible = true
    },

    handleAddUser() {
      this.dialogAddUserVisible = true
    },

    submit() {
      this.$refs.form.validate((valid) => {
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
      add(this.formData).then((response) => {
        this.getList()
        this.getTotal()
        this.$message.success('Need to login again')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.getList()
        this.$message.success('Success')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    addUser() {
      this.$refs.addUserForm.validate((valid) => {
        if (valid) {
          this.addUserFormProps.disabled = true
          addUser(this.addUserFormData).then((response) => {
            this.dialogAddUserVisible = false
            this.$message.success('Success')
            this.getBindUserList(this.addUserFormData.namespaceId)
          }).finally(() => {
            this.addUserFormProps.disabled = false
          })
        } else {
          return false
        }
      })
    },

    removeUser(data) {
      this.$confirm(this.$i18n.t('namespacePage.removeUserTips'), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        removeUser(data.id).then((response) => {
          this.$message.success('Success')
          this.getBindUserList(data.namespaceId)
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    storeFormData() {
      this.tempFormData = JSON.parse(JSON.stringify(this.formData))
    },

    restoreFormData() {
      this.formData = JSON.parse(JSON.stringify(this.tempFormData))
    }
  }
}
</script>
