<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
    </el-row>
    <el-table
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="account" :label="$t('account')" />
      <el-table-column prop="name" :label="$t('name')" />
      <el-table-column prop="contact" :label="$t('contact')" show-overflow-tooltip />
      <el-table-column prop="superManager" :label="$t('admin')" width="50" align="center">
        <template slot-scope="scope">
          {{ $t(`boolOption[${scope.row.superManager}]`) }}
        </template>
      </el-table-column>
      <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
      <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
      <el-table-column prop="operation" :label="$t('op')" width="130" align="center">
        <template slot-scope="scope">
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
    <el-row type="flex" justify="end" style="margin-top: 10px;">
      <el-pagination
        hide-on-single-page
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handleCurrentChange"
      />
    </el-row>
    <el-dialog :title="$t('setting')" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="80px">
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
        <el-form-item :label="$t('admin')" prop="super_manager">
          <el-radio-group v-model="formData.super_manager">
            <el-radio :label="1">{{ $t('boolOption[1]') }}</el-radio>
            <el-radio :label="0">{{ $t('boolOption[0]') }}</el-radio>
          </el-radio-group>
          <el-popover
            placement="top-start"
            :title="$t('desc')"
            width="300"
            trigger="hover"
          >
            {{ $t('memberPage.permissionDesc') }}
            <el-button slot="reference" type="text" icon="el-icon-question" style="color: #666;" />
          </el-popover>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { validUsername, validPassword } from '@/utils/validate'
import { getList, getTotal, add, edit, remove } from '@/api/user'

export default {
  name: 'Member',
  data() {
    const validateUsername = (rule, value, callback) => {
      if (!validUsername(value)) {
        callback(new Error('Greater than 5 characters'))
      } else {
        callback()
      }
    }
    const validatePassword = (rule, value, callback) => {
      if (!value) {
        callback()
      } else if (!validPassword(value)) {
        callback(new Error('8 to 16 characters and a minimum of 2 character sets from these classes: [letters], [numbers], [special characters]'))
      } else {
        callback()
      }
    }
    return {
      dialogVisible: false,
      tableData: [],
      tempFormData: {},
      pagination: {
        page: 1,
        rows: 16,
        total: 0
      },
      formProps: {
        disabled: false
      },
      formData: {
        id: 0,
        account: '',
        password: '',
        name: '',
        contact: '',
        super_manager: 0
      },
      formRules: {
        account: [
          { required: true, message: 'Account required', trigger: 'blur', validator: validateUsername }
        ],
        password: [
          { trigger: 'blur', validator: validatePassword }
        ],
        name: [
          { required: true, message: 'Name required', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.storeFormData()
    this.getList()
    this.getTotal()
  },
  methods: {
    getList() {
      getList(this.pagination).then((response) => {
        this.tableData = response.data.list
      })
    },

    getTotal() {
      getTotal(this.crontabCommand).then((response) => {
        this.pagination.total = response.data.total
      })
    },

    handleCurrentChange(val) {
      this.pagination.page = val
      this.getList()
    },

    handleAdd() {
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleEdit(data) {
      this.restoreFormData()
      this.formData = Object.assign(this.formData, data)
      this.dialogVisible = true
    },

    handleRemove(data) {
      this.$confirm(this.$i18n.t('memberPage.removeUserTips', { name: data.name }), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message.success('Success')
          this.getList()
          this.getTotal()
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
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
        this.$message.success('Success')
        this.getList()
        this.getTotal()
        this.dialogVisible = false
      }).finally(() => {
        this.formProps.disabled = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.$message.success('Success')
        this.getList()
        this.dialogVisible = false
      }).finally(() => {
        this.formProps.disabled = false
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
