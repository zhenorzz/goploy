<template>
  <el-row class="app-container">
    <el-row v-show="$store.getters.superManager" class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd">添加</el-button>
    </el-row>
    <el-table
      :key="tableHeight"
      border
      stripe
      highlight-current-row
      :max-height="tableHeight"
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="160" />
      <el-table-column prop="name" label="空间名" />
      <el-table-column prop="insertTime" label="创建时间" width="160" />
      <el-table-column prop="updateTime" label="更新时间" width="160" />
      <el-table-column prop="operation" label="操作" width="172">
        <template slot-scope="scope">
          <el-button type="warning" @click="handleUser(scope.row)">成员管理</el-button>
          <el-button size="small" type="primary" @click="handleEdit(scope.row)">编辑</el-button>
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
    <el-dialog title="设置" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="成员管理" :visible.sync="dialogUserVisible">
      <el-row class="app-bar" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="handleAddUser">添加</el-button>
      </el-row>
      <el-table
        border
        stripe
        highlight-current-row
        :data="tableUserData.filter(row => row.role!== $global.Admin)"
        style="width: 100%"
      >
        <el-table-column prop="userId" label="用户ID" />
        <el-table-column prop="userName" label="用户名称" />
        <el-table-column prop="role" label="角色" />
        <el-table-column prop="insertTime" width="160" label="绑定时间" />
        <el-table-column prop="updateTime" width="160" label="更新时间" />
        <el-table-column prop="operation" label="操作" width="80">
          <template slot-scope="scope">
            <el-button type="danger" @click="removeUser(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogUserVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="添加用户" :visible.sync="dialogAddUserVisible">
      <el-form ref="addUserForm" :rules="addUserFormRules" :model="addUserFormData">
        <el-form-item label="绑定用户" label-width="120px" prop="userIds">
          <el-select v-model="addUserFormData.userIds" multiple clearable filterable placeholder="选择用户，可多选">
            <el-option
              v-for="(item, index) in userOption.filter(item => item.superManager !== 1)"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="角色" label-width="120px" prop="role">
          <el-select v-model="addUserFormData.role" placeholder="选择角色">
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
        <el-button @click="dialogAddUserVisible = false">取 消</el-button>
        <el-button :disabled="addUserFormProps.disabled" type="primary" @click="addUser">确 定</el-button>
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
          { required: true, message: '请输入名称', trigger: 'blur' }
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
          { type: 'array', required: true, message: '请选择用户', trigger: 'change' }
        ],
        role: [
          { required: true, message: '请选择角色', trigger: 'change' }
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
      getList(this.pagination).then((response) => {
        this.tableData = response.data.list
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

    // 分页事件
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
      // 先把namespaceID写入添加用户的表单
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
        this.$message.success('添加成功')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.getList()
        this.$message.success('修改成功')
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
            this.$message.success('添加成功')
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
      this.$confirm('此操作将永久删除该用户的绑定关系, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        removeUser(data.id).then((response) => {
          this.$message.success('删除成功')
          this.getBindUserList(data.namespaceId)
        })
      }).catch(() => {
        this.$message.info('已取消删除')
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
