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
      <el-table-column prop="role" label="角色" />
      <el-table-column prop="createTime" label="创建时间" width="160" />
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
        <el-form-item label="角色" prop="roleId">
          <el-select v-model="formData.role" placeholder="选择角色">
            <el-option
              v-for="(role, index) in roleOption"
              :key="index"
              :label="role"
              :value="role"
            />
          </el-select>
          <el-popover
            placement="top-start"
            title="权限说明"
            width="300"
            trigger="hover"
          >
            <p>admin具有一切权限</p>
            <p>manager不具有成员管理</p>
            <p>group-manager管理和构建分组的项目，也能绑定不属于分组的项目(只能构建)</p>
            <p>member只允许构建绑定的项目</p>
            <el-button slot="reference" type="text" icon="el-icon-question" style="color: #666;" />
          </el-popover>
        </el-form-item>
        <el-form-item v-show="formData.role==='group-manager'" label="管理分组" prop="groupId">
          <el-select
            v-model="formProps.groupIds"
            multiple
            placeholder="选择分组"
            style="width:100%"
          >
            <el-option
              v-for="(item, index) in groupOption"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-show="formData.role==='member' || formData.role==='group-manager'" label="绑定项目" prop="projectIds">
          <el-cascader
            v-model="formProps.projectIds"
            style="width: 100%"
            :options="filterProjectOption"
            :props="{ multiple: true }"
            collapse-tags
            clearable
          />
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
import { getOption as getGroupOption } from '@/api/group'
import { getOption as getProjectOption, getBindProjectList } from '@/api/project'
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
      groupOption: [],
      projectOption: [],
      tableData: [],
      tempFormData: {},
      pagination: {
        page: 1,
        rows: 11,
        total: 0
      },
      formProps: {
        disabled: false,
        groupIds: [],
        projectIds: [],
        projectOption: []
      },
      formData: {
        id: 0,
        account: '',
        password: '',
        name: '',
        mobile: '',
        role: 'member',
        manageGroupStr: '',
        projectIds: []
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
  computed: {
    filterProjectOption: function() {
      return this.projectOption.filter(element => this.formProps.groupIds.indexOf(element.value) === -1)
    }
  },
  created() {
    this.storeFormData()
    this.getRoleOption()
    this.getGroupOption()
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
        this.roleOption = response.data.roleList || []
      })
    },

    getGroupOption() {
      getGroupOption().then((response) => {
        this.groupOption = response.data.groupList || []
        if (this.groupOption.length !== 0) {
          getProjectOption().then(response => {
            this.projectOption = this.groupOption.map(element => {
              return {
                value: element.id,
                label: element.name,
                children: []
              }
            })
            this.projectOption.unshift({
              value: 0,
              label: '默认',
              children: []
            })
            response.data.projectList.forEach(element => {
              const groupIndex = this.projectOption.findIndex(group => group.value === element.groupId)
              this.projectOption[groupIndex].children.push({
                value: element.id,
                label: element.name
              })
            })
          })
        }
      })
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
      this.restoreFormData()
      this.formData = Object.assign(this.formData, data)
      this.formProps.groupIds = data.manageGroupStr.split(',').filter(element => element !== '' && element !== 'all').map(element => {
        return parseInt(element)
      })
      getBindProjectList(data.id).then((response) => {
        this.formData.projectIds = response.data.projectUserMap ? response.data.projectUserMap.map(element => {
          return element.projectId
        }) : []
        this.formProps.projectIds = response.data.projectUserMap ? response.data.projectUserMap.map(element => {
          return [element.groupId, element.projectId]
        }) : []
      })
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

    handleGroupRadioChange(value) {
      if (value === 2) {
        this.formProps.showGroupSelect = true
      } else {
        this.formProps.showGroupSelect = false
      }
    },

    submit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          if (this.formData.role === 'admin' || this.formData.role === 'manager') {
            this.formData.manageGroupStr = 'all'
          } else if (this.formData.role === 'group-manager') {
            this.formData.manageGroupStr = this.formProps.groupIds.sort((x, y) => x - y).join(',')
          } else {
            this.formData.manageGroupStr = ''
          }
          this.formData.projectIds = []
          if (this.formData.role === 'group-manager' || this.formData.role === 'member') {
            this.formProps.projectIds.forEach(element => {
              this.formData.projectIds.push(element[1])
            })
          }
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
