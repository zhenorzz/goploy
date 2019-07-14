<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="dialogFormVisible = true">添加</el-button>
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
          <el-button size="small" type="primary">编辑</el-button>
          <el-button size="small" type="danger">删除</el-button>
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
    <el-dialog title="新增成员" :visible.sync="dialogFormVisible">
      <el-form ref="form" :rules="form.rules" :model="form">
        <el-form-item label="账号" label-width="120px" prop="account">
          <el-input v-model="form.account" autocomplete="off" />
        </el-form-item>
        <el-form-item label="密码" label-width="120px" prop="password">
          <el-input v-model="form.password" autocomplete="off" placeholder="请输入初始密码" />
        </el-form-item>
        <el-form-item label="名称" label-width="120px" prop="name">
          <el-input v-model="form.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="手机号码" label-width="120px" prop="mobile">
          <el-input v-model="form.email" autocomplete="off" />
        </el-form-item>
        <el-form-item label="角色" label-width="120px" prop="role">
          <el-select v-model="form.roleId" placeholder="选择角色">
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
        <el-button @click="dialogFormVisible = false">取 消</el-button>
        <el-button :disabled="form.disabled" type="primary" @click="add">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { getList, add } from '@/api/user'
import { getOption as getRoleOption } from '@/api/role'
import { parseTime } from '@/utils'

export default {
  data() {
    return {
      dialogFormVisible: false,
      roleOption: [],
      tableData: [],
      pagination: {
        page: 1,
        rows: 11,
        total: 0
      },
      form: {
        disabled: false,
        account: '',
        password: '',
        name: '',
        mobile: '',
        roleId: '',
        rules: {
          account: [
            { required: true, message: '请输入账号', trigger: 'blur' }
          ],
          password: [
            { required: true, min: 5, message: '不少于5位数', trigger: 'blur' }
          ],
          name: [
            { required: true, message: '请输入名称', trigger: 'blur' }
          ],
          roleId: [
            { required: true, message: '请选择角色', trigger: 'change' }
          ]
        }
      }
    }
  },
  created() {
    this.getRoleOption()
    this.getUserList()
  },
  methods: {
    getUserList() {
      getList(this.pagination).then((response) => {
        const userList = response.data.userList
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
    add() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.form.disabled = true
          add(this.form.account, this.form.password, this.form.name, this.form.email, this.form.role).then((response) => {
            this.form.disabled = this.dialogFormVisible = false
            this.$message({
              message: response.message,
              type: 'success',
              duration: 5 * 1000
            })
            this.getUserList()
          }).catch(() => {
            this.form.disabled = this.dialogFormVisible = false
          })
        } else {
          this.form.disabled = this.dialogFormVisible = false
          return false
        }
      })
    }
  }
}
</script>
