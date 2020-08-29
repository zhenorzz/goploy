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
      <el-table-column prop="superManager" label="超管">
        <template slot-scope="scope">
          {{ scope.row.superManager === 1 ? "是":"否" }}
        </template>
      </el-table-column>
      <el-table-column prop="insertTime" label="创建时间" width="160" />
      <el-table-column prop="updateTime" label="更新时间" width="160" />
      <el-table-column prop="operation" label="操作" width="150">
        <template slot-scope="scope">
          <el-button v-if="scope.row.id !== 1 && scope.row.id !== $store.getters.uid" size="small" type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button v-if="scope.row.id !== 1 && scope.row.id !== $store.getters.uid" size="small" type="danger" @click="handleRemove(scope.row)">删除</el-button>
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
    <el-dialog title="成员设置" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="80px">
        <el-form-item label="账号" prop="account">
          <el-input v-model="formData.account" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" prop="password">
          <el-input v-model="formData.password" autocomplete="off" placeholder="请输入初始密码" />
        </el-form-item>
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="手机号码" prop="mobile">
          <el-input v-model="formData.mobile" autocomplete="off" />
        </el-form-item>
        <el-form-item label="超管" prop="super_manager">
          <el-radio-group v-model="formData.super_manager">
            <el-radio :label="1">是</el-radio>
            <el-radio :label="0">否</el-radio>
          </el-radio-group>
          <el-popover
            placement="top-start"
            title="权限说明"
            width="300"
            trigger="hover"
          >
            超管具有所有空间和项目权限
            <el-button slot="reference" type="text" icon="el-icon-question" style="color: #666;" />
          </el-popover>
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
import { getList, getTotal, add, edit, remove } from '@/api/user'

export default {
  name: 'Member',
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
        mobile: '',
        role: 'member',
        super_manager: 0
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

    // 分页事件
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
      this.$confirm('此操作将删除该组员, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message.success('删除成功')
          this.getList()
          this.getTotal()
        })
      }).catch(() => {
        this.$message.info('已取消删除')
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
        this.$message.success('添加成功')
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
        this.$message.success('修改成功')
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
