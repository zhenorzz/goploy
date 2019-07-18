<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd">添加</el-button>
    </el-row>
    <el-table
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="account" label="账号" />
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="mobile" label="手机号码" show-overflow-tooltip />
      <el-table-column prop="role" label="角色">
        <template slot-scope="scope">
          {{ findRoleName(scope.row.roleId) }}
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="160" />
      <el-table-column prop="updateTime" label="更新时间" width="160" />
      <el-table-column prop="operation" label="操作" width="150">
        <template slot-scope="scope">
          <el-button size="small" type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button :disabled="scope.row.id === 1" size="small" type="danger" @click="handleRemove(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px;">
      <el-pagination
        v-show="pagination.total>pagination.rows"
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handleCurrentChange"
      />
    </el-row>
    <el-dialog title="新增成员" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData">
        <el-form-item label="账号" label-width="120px" prop="account">
          <el-input v-model="formData.account" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" label-width="120px" prop="password">
          <el-input v-model="formData.password" autocomplete="off" placeholder="请输入初始密码" />
        </el-form-item>
        <el-form-item label="名称" label-width="120px" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="手机号码" label-width="120px" prop="mobile">
          <el-input v-model="formData.mobile" autocomplete="off" />
        </el-form-item>
        <el-form-item label="角色" label-width="120px" prop="role">
          <el-select v-model="formData.roleId" placeholder="选择角色">
            <el-option
              v-for="(item, index) in roleOption"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { validUsername, validPassword } from '@/utils/validate'
import { getList, add, edit, remove } from '@/api/user'
import { getOption as getRoleOption } from '@/api/role'
import { parseTime } from '@/utils'

export default {
  data() {
    const validateUsername = (rule, value, callback) => {
      if (!validUsername(value)) {
        callback(new Error('请输入正确的用户名'))
      } else {
        callback()
      }
    }
    const validatePassword = (rule, value, callback) => {
      if (!value) {
        callback()
      } else if (!validPassword(value)) {
        callback(new Error('8到16个字符，至少包含字母、数字、特殊符号中的两种'))
      } else {
        callback()
      }
    }
    return {
      dialogVisible: false,
      roleOption: [],
      tableData: [],
      tempFormData: {},
      pagination: {
        page: 1,
        rows: 11,
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
        mobile: '',
        roleId: ''
      },
      formRules: {
        account: [
          { required: true, message: '请输入账号', trigger: 'blur', validator: validateUsername }
        ],
        password: [
          { trigger: 'blur', validator: validatePassword }
        ],
        name: [
          { required: true, message: '请输入名称', trigger: 'blur' }
        ],
        roleId: [
          { required: true, message: '请选择角色', trigger: 'change' }
        ]
      }
    }
  },
  created() {
    this.storeFormData()
    this.getRoleOption()
    this.getUserList()
  },
  methods: {
    getUserList() {
      getList(this.pagination).then((response) => {
        const userList = response.data.userList || []
        userList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
        this.tableData = userList
        this.pagination = response.data.pagination
      })
    },
    getRoleOption() {
      getRoleOption().then((response) => {
        this.roleOption = response.data.roleList
      })
    },

    findRoleName(roleId) {
      return this.roleOption.find(element => element.id === roleId)['name']
    },

    // 分页事件
    handleCurrentChange(val) {
      this.pagination.page = val
      this.getUserList()
    },

    handleAdd() {
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleEdit(data) {
      this.formData = Object.assign({}, data)
      this.dialogVisible = true
    },

    handleRemove(data) {
      this.$confirm('此操作将删除该组员, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message({
            message: response.message,
            type: 'success',
            duration: 5 * 1000
          })
          this.getUserList()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
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
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
        this.getUserList()
        this.dialogVisible = false
      }).finally(() => {
        this.formProps.disabled = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
        this.getUserList()
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
