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
      <el-table-column prop="name" label="项目名称" />
      <el-table-column prop="url" label="项目地址" />
      <el-table-column prop="path" label="部署路径" />
      <el-table-column prop="status" label="状态" />
      <el-table-column prop="createTime" width="160" label="创建时间" />
      <el-table-column prop="updateTime" width="160" label="更新时间" />
      <el-table-column prop="operation" label="操作" width="430">
        <template slot-scope="scope">
          <el-button
            :disabled="scope.row.status === '初始化成功'"
            size="small"
            type="success"
            @click="create(scope.row.id)"
          >初始化</el-button>
          <el-button type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button type="primary" @click="handleServer(scope.row)">服务器管理</el-button>
          <el-button type="primary" @click="handleUser(scope.row)">成员管理</el-button>
          <el-button type="danger">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="项目设置" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData">
        <el-form-item label="项目名称" label-width="120px" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="项目地址" label-width="120px" prop="url">
          <el-input v-model="formData.url" autocomplete="off" />
        </el-form-item>
        <el-form-item label="部署路径" label-width="120px" prop="path">
          <el-input v-model="formData.path" autocomplete="off" />
        </el-form-item>
        <el-form-item v-show="formProps.showServers" label="绑定服务器" label-width="120px" prop="serverIds">
          <el-select v-model="formData.serverIds" multiple placeholder="选择服务器，可多选">
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item v-show="formProps.showUsers" label="绑定组员" label-width="120px" prop="userIds">
          <el-select v-model="formData.userIds" multiple placeholder="选择组员，可多选">
            <el-option
              v-for="(item, index) in userOption"
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
    <el-dialog title="服务器管理" :visible.sync="dialogServerVisible">
      <el-table
        border
        stripe
        highlight-current-row
        :data="tableServerData"
        style="width: 100%"
      >
        <el-table-column prop="serverId" label="服务器ID" />
        <el-table-column prop="serverName" label="服务器名称" />
        <el-table-column prop="createTime" width="160" label="绑定时间" />
        <el-table-column prop="updateTime" width="160" label="更新时间" />
        <el-table-column prop="operation" label="操作" width="130">
          <template slot-scope="scope">
            <el-button type="danger">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogServerVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="成员管理" :visible.sync="dialogUserVisible">
      <el-table
        border
        stripe
        highlight-current-row
        :data="tableUserData"
        style="width: 100%"
      >
        <el-table-column prop="userId" label="用户ID" />
        <el-table-column prop="userName" label="用户名称" />
        <el-table-column prop="createTime" width="160" label="绑定时间" />
        <el-table-column prop="updateTime" width="160" label="更新时间" />
        <el-table-column prop="operation" label="操作" width="130">
          <template slot-scope="scope">
            <el-button type="danger">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogUserVisible = false">取 消</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>

import { getOption as getUserOption } from '@/api/user'
import { getOption as getServerOption } from '@/api/server'
import { getList, getBindServerList, getBindUserList, add, edit, create } from '@/api/project'
import { parseTime } from '@/utils'

const STATUS = ['未初始化', '初始化中', '初始化成功', '初始化失败']
export default {
  data() {
    return {
      dialogVisible: false,
      dialogServerVisible: false,
      dialogUserVisible: false,
      serverOption: [],
      userOption: [],
      tableData: [],
      tableServerData: [],
      tableUserData: [],
      formProps: {
        disabled: false,
        showServers: true,
        showUsers: true
      },
      tempFormData: {},
      formData: {
        id: 0,
        name: '',
        url: '',
        path: '',
        serverIds: [],
        userIds: []
      },
      formRules: {
        name: [
          { required: true, message: '请输入项目名称', trigger: ['blur'] }
        ],
        url: [
          { required: true, message: '请输入项目地址', trigger: ['blur'] }
        ],
        path: [
          { required: true, message: '请输入部署路径', trigger: ['blur'] }
        ],
        serverIds: [
          { type: 'array', message: '请选择服务器', trigger: 'change' }
        ],
        userIds: [
          { type: 'array', message: '请选择组员', trigger: 'change' }
        ]
      }
    }
  },
  created() {
    this.storeFormData()
    this.get()
  },
  methods: {
    handleAdd() {
      this.restoreFormData()
      this.formProps.showServers = this.formProps.showUsers = true
      this.dialogVisible = true
    },

    handleEdit(data) {
      this.formData = Object.assign({}, data)
      this.formProps.showServers = this.formProps.showUsers = false
      this.dialogVisible = true
    },

    handleServer(data) {
      getBindServerList(data.id).then((response) => {
        this.tableServerData = response.data.projectServerMap || []
        this.tableServerData.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
        this.dialogServerVisible = true
      })
    },

    handleUser(data) {
      getBindUserList(data.id).then((response) => {
        this.tableUserData = response.data.projectUserMap || []
        this.tableUserData.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
        this.dialogUserVisible = true
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
        this.dialogVisible = false
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
        this.getProjectList()
      }).finally(() => {
        this.formProps.disabled = false
      })
    },
    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.dialogVisible = false
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
        this.getProjectList()
      }).finally(() => {
        this.formProps.disabled = false
      })
    },
    get() {
      this.getProjectList()
      getServerOption().then((response) => {
        this.serverOption = response.data.serverList
      })
      getUserOption().then((response) => {
        this.userOption = response.data.userList
      })
    },
    getProjectList() {
      getList().then((response) => {
        const projectList = response.data.projectList
        projectList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
          element.status = STATUS[element.status]
        })
        this.tableData = projectList
      }).catch(() => {
      })
    },
    create(projectId) {
      create(projectId).then((response) => {
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
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
