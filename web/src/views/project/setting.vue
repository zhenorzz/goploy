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
      <el-table-column prop="name" label="项目名称" width="200" />
      <el-table-column prop="url" label="项目地址" width="350" />
      <el-table-column prop="path" label="部署路径" />
      <el-table-column prop="group" label="分组" width="100">
        <template slot-scope="scope">
          {{ findGroupName(scope.row.groupId) }}
        </template>
      </el-table-column>
      <el-table-column prop="environment" width="120" label="环境" />
      <el-table-column prop="branch" width="160" label="分支" />
      <el-table-column prop="operation" label="操作" width="350">
        <template slot-scope="scope">
          <el-button type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button type="success" @click="handleServer(scope.row)">服务器管理</el-button>
          <el-button type="warning" @click="handleUser(scope.row)">成员管理</el-button>
          <el-button type="danger" @click="handleRemove(scope.row)">删除</el-button>
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
    <el-dialog title="项目设置" :visible.sync="dialogVisible" width="60%">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="120px">
        <el-row>
          <el-col :span="8">
            <el-form-item label="项目名称" prop="name">
              <el-input v-model="formData.name" autocomplete="off" />
            </el-form-item>
            <el-form-item label="项目地址" prop="url">
              <el-input v-model="formData.url" autocomplete="off" />
            </el-form-item>
            <el-form-item label="部署路径" prop="path">
              <el-input v-model="formData.path" autocomplete="off" />
            </el-form-item>
            <el-form-item label="环境" prop="environment">
              <el-select v-model="formData.environment" placeholder="选择环境" style="width:100%">
                <el-option label="生产环境" value="生产环境" />
                <el-option label="测试环境" value="测试环境" />
                <el-option label="开发环境" value="开发环境" />
              </el-select>
            </el-form-item>
            <el-form-item label="分支" prop="branch">
              <el-input v-model="formData.branch" autocomplete="off" />
            </el-form-item>
            <el-form-item label="rsync选项" prop="rsyncOption">
              <el-input v-model="formData.rsyncOption" type="textarea" :rows="4" autocomplete="off" placeholder="-rtv --exclude .git --delete-after" />
            </el-form-item>
            <el-form-item label="绑定分组" prop="groupId">
              <el-select v-model="formData.groupId" placeholder="选择分组" style="width:100%">
                <el-option label="默认" :value="0" />
                <el-option
                  v-for="(item, index) in groupOption"
                  :key="index"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item v-show="formProps.showServers" label="绑定服务器" prop="serverIds">
              <el-select v-model="formData.serverIds" multiple placeholder="选择服务器，可多选" style="width:100%">
                <el-option
                  v-for="(item, index) in serverOption"
                  :key="index"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item v-show="formProps.showUsers" label="绑定组员" prop="userIds">
              <el-select v-model="formData.userIds" multiple placeholder="选择组员，可多选" style="width:100%">
                <el-option
                  v-for="(item, index) in userOption"
                  :key="index"
                  :label="item.name"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="16">
            <el-form-item label="拉取后运行脚本" prop="afterPullScrpit">
              <codemirror v-model="formData.afterPullScript" :options="cmOptions" />
            </el-form-item>
            <el-form-item label="部署后运行脚本" prop="afterDeployScrpit">
              <codemirror v-model="formData.afterDeployScript" :options="cmOptions" />
            </el-form-item>
          </el-col>
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="服务器管理" :visible.sync="dialogServerVisible">
      <el-row class="app-bar" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="handleAddServer">添加</el-button>
      </el-row>
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
        <el-table-column prop="operation" label="操作" width="75">
          <template slot-scope="scope">
            <el-button type="danger" @click="removeProjectServer(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogServerVisible = false">取 消</el-button>
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
        :data="tableUserData"
        style="width: 100%"
      >
        <el-table-column prop="userId" label="用户ID" />
        <el-table-column prop="userName" label="用户名称" />
        <el-table-column prop="createTime" width="160" label="绑定时间" />
        <el-table-column prop="updateTime" width="160" label="更新时间" />
        <el-table-column prop="operation" label="操作" width="75">
          <template slot-scope="scope">
            <el-button type="danger" @click="removeProjectUser(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogUserVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="添加服务器" :visible.sync="dialogAddServerVisible">
      <el-form ref="addServerForm" :rules="addServerFormRules" :model="addServerFormData">
        <el-form-item label="绑定服务器" label-width="120px" prop="serverIds">
          <el-select v-model="addServerFormData.serverIds" multiple placeholder="选择服务器，可多选">
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogAddServerVisible = false">取 消</el-button>
        <el-button :disabled="addServerFormProps.disabled" type="primary" @click="addServer">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="添加用户" :visible.sync="dialogAddUserVisible">
      <el-form ref="addUserForm" :rules="addUserFormRules" :model="addUserFormData">
        <el-form-item label="绑定组员" label-width="120px" prop="userIds">
          <el-select v-model="addUserFormData.userIds" multiple placeholder="选择组员，可多选">
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
        <el-button @click="dialogAddUserVisible = false">取 消</el-button>
        <el-button :disabled="addUserFormProps.disabled" type="primary" @click="addUser">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>

import { getOption as getUserOption } from '@/api/user'
import { getOption as getServerOption } from '@/api/server'
import { getOption as getGroupOption } from '@/api/group'
import { getList, getBindServerList, getBindUserList, add, edit, create, remove, addServer, addUser, removeProjectServer, removeProjectUser } from '@/api/project'
import { parseTime } from '@/utils'
// require component
import { codemirror } from 'vue-codemirror'
import 'codemirror/mode/shell/shell.js'
import 'codemirror/theme/darcula.css'
// require styles
import 'codemirror/lib/codemirror.css'
import 'codemirror/addon/scroll/simplescrollbars.js'
import 'codemirror/addon/scroll/simplescrollbars.css'
export default {
  components: {
    codemirror
  },
  data() {
    return {
      cmOptions: {
        tabSize: 4,
        mode: 'text/x-sh',
        lineNumbers: true,
        line: true,
        scrollbarStyle: 'overlay',
        theme: 'darcula'
      },
      dialogVisible: false,
      dialogServerVisible: false,
      dialogUserVisible: false,
      dialogAddServerVisible: false,
      dialogAddUserVisible: false,
      serverOption: [],
      userOption: [],
      groupOption: [],
      tableData: [],
      pagination: {
        page: 1,
        rows: 11,
        total: 0
      },
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
        groupId: 0,
        name: '',
        url: '',
        path: '',
        afterPullScript: '',
        afterDeployScript: '',
        environment: '生产环境',
        branch: 'master',
        rsyncOption: '-rtv --exclude .git --delete-after',
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
        environment: [
          { required: true, message: '请选择环境', trigger: ['blur'] }
        ],
        branch: [
          { required: true, message: '请输入分支名称', trigger: ['blur'] }
        ],
        serverIds: [
          { type: 'array', message: '请选择服务器', trigger: 'change' }
        ],
        userIds: [
          { type: 'array', message: '请选择组员', trigger: 'change' }
        ]
      },
      addServerFormProps: {
        disabled: false
      },
      addServerFormData: {
        projectId: 0,
        serverIds: []
      },
      addServerFormRules: {
        serverIds: [
          { type: 'array', required: true, message: '请选择服务器', trigger: 'change' }
        ]
      },
      addUserFormProps: {
        disabled: false
      },
      addUserFormData: {
        projectId: 0,
        userIds: []
      },
      addUserFormRules: {
        userIds: [
          { type: 'array', required: true, message: '请选择用户', trigger: 'change' }
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
      this.formData.serverIds = []
      this.formData.userIds = []
      this.formProps.showServers = this.formProps.showUsers = false
      this.dialogVisible = true
    },

    handleRemove(data) {
      this.$confirm('此操作将删除该项目, 是否继续?', '提示', {
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
          this.getProjectList()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    handleServer(data) {
      this.getBindServerList(data.id)
      // 先把projectID写入添加服务器的表单
      this.addServerFormData.projectId = data.id
      this.dialogServerVisible = true
    },

    handleGroup(data) {
      this.getBindServerList(data.id)
      // 先把projectID写入添加服务器的表单
      this.addServerFormData.projectId = data.id
      this.dialogServerVisible = true
    },

    handleUser(data) {
      this.getBindUserList(data.id)
      // 先把projectID写入添加用户的表单
      this.addUserFormData.projectId = data.id
      this.dialogUserVisible = true
    },

    handleAddServer() {
      this.dialogAddServerVisible = true
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

    create(projectId) {
      create(projectId).then((response) => {
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
      })
    },

    addServer() {
      this.$refs.addServerForm.validate((valid) => {
        if (valid) {
          this.addServerFormProps.disabled = true
          addServer(this.addServerFormData).then((response) => {
            this.dialogAddServerVisible = false
            this.$message({
              message: response.message,
              type: 'success',
              duration: 5 * 1000
            })
            this.getBindServerList(this.addServerFormData.projectId)
          }).finally(() => {
            this.addServerFormProps.disabled = false
          })
        } else {
          return false
        }
      })
    },

    addUser() {
      this.$refs.addUserForm.validate((valid) => {
        if (valid) {
          this.addUserFormProps.disabled = true
          addUser(this.addUserFormData).then((response) => {
            this.dialogAddUserVisible = false
            this.$message({
              message: response.message,
              type: 'success',
              duration: 5 * 1000
            })
            this.getBindUserList(this.addUserFormData.projectId)
          }).finally(() => {
            this.addUserFormProps.disabled = false
          })
        } else {
          return false
        }
      })
    },

    removeProjectServer(data) {
      this.$confirm('此操作将永久删除该服务器的绑定关系, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        removeProjectServer(data.id).then((response) => {
          this.$message({
            message: response.message,
            type: 'success',
            duration: 5 * 1000
          })
          this.getBindServerList(data.projectId)
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    removeProjectUser(data) {
      this.$confirm('此操作将永久删除该用户的绑定关系, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        removeProjectUser(data.id).then((response) => {
          this.$message({
            message: response.message,
            type: 'success',
            duration: 5 * 1000
          })
          this.getBindUserList(data.projectId)
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    get() {
      this.getProjectList()
      getServerOption().then((response) => {
        this.serverOption = response.data.serverList || []
      })
      getUserOption().then((response) => {
        this.userOption = response.data.userList || []
      })
      getGroupOption().then((response) => {
        this.groupOption = response.data.groupList || []
      })
    },

    getProjectList() {
      getList(this.pagination).then((response) => {
        const projectList = response.data.projectList
        projectList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
        this.tableData = projectList
        this.pagination = response.data.pagination
      }).catch(() => {
      })
    },

    getBindServerList(projectID) {
      getBindServerList(projectID).then((response) => {
        this.tableServerData = response.data.projectServerMap || []
        this.tableServerData.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
      })
    },

    getBindUserList(projectID) {
      getBindUserList(projectID).then((response) => {
        this.tableUserData = response.data.projectUserMap || []
        this.tableUserData.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
      })
    },

    // 分页事件
    handlePageChange(val) {
      this.pagination.page = val
      this.getProjectList()
    },

    findGroupName(groupId) {
      const projectGroup = this.groupOption.find(element => element.id === groupId)
      return projectGroup ? projectGroup['name'] : '默认'
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
<style>
.CodeMirror {
  border-radius: 4px;
  border: 1px solid #DCDFE6;
  height: 200px;
}
</style>
