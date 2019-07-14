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
      <el-table-column prop="operation" label="操作" width="230">
        <template slot-scope="scope">
          <el-button
            :disabled="scope.row.status === '初始化成功'"
            size="small"
            type="success"
            @click="create(scope.row.id)"
          >初始化</el-button>
          <el-button type="primary" @click="handleEdit(scope.row)">编辑</el-button>
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
        <el-form-item label="绑定服务器" label-width="120px" prop="serverIds">
          <el-select v-model="formData.serverIds" multiple placeholder="选择服务器，可多选">
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="绑定组员" label-width="120px" prop="userIds">
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
  </el-row>
</template>
<script>

import { getOption as getUserOption } from '@/api/user'
import { getOption as getServerOption } from '@/api/server'
import { getList, getDetail, add, create } from '@/api/project'
import { parseTime } from '@/utils'

const STATUS = ['未初始化', '初始化中', '初始化成功', '初始化失败']
export default {
  data() {
    return {
      dialogVisible: false,
      serverOption: [],
      userOption: [],
      tableData: [],
      formProps: {
        disabled: false
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
          { type: 'array', required: true, message: '请选择服务器', trigger: 'change' }
        ],
        userIds: [
          { type: 'array', required: true, message: '请选择组员', trigger: 'change' }
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
      this.dialogVisible = true
    },

    handleEdit(data) {
      this.formData = Object.assign({}, data)
      getDetail(data.id).then(response => {
        const projectServerMap = response.data.projectServerMap
        this.formData.serverIds = projectServerMap.map(element => element.serverId)
        const projectUserMap = response.data.projectUserMap
        this.formData.userIds = projectUserMap.map(element => element.userId)
        this.dialogVisible = true
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
