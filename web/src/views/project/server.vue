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
      <el-table-column prop="name" label="服务器" />
      <el-table-column prop="ip" label="IP" />
      <el-table-column prop="port" label="端口" />
      <el-table-column prop="owner" label="sshKey所有者" show-overflow-tooltip />
      <el-table-column prop="group" label="分组" width="100">
        <template slot-scope="scope">
          {{ findGroupName(scope.row.groupId) }}
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="创建时间" width="160" />
      <el-table-column prop="updateTime" label="更新时间" width="160" />
      <el-table-column prop="operation" label="操作" width="220">
        <template slot-scope="scope">
          <el-button size="small" type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button size="small" type="warning" @click="handleInstall(scope.row)">安装</el-button>
          <el-button size="small" type="danger" @click="handleRemove(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="服务器设置" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="120px">
        <el-form-item label="服务器名称" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="IP" prop="ip">
          <el-input v-model="formData.ip" autocomplete="off" />
        </el-form-item>
        <el-form-item label="port" prop="port">
          <el-input v-model.number="formData.port" autocomplete="off" />
        </el-form-item>
        <el-form-item label="sshKey所有者" prop="owner">
          <el-input v-model="formData.owner" autocomplete="off" placeholder="root" />
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
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="安装模板" :visible.sync="templateDialogVisible">
      <el-form ref="templateForm" :rules="templateFormRules" :model="templateFormData" label-width="120px">
        <el-form-item label="选择模板" prop="templateId">
          <el-select v-model="templateFormData.templateId" placeholder="选择模板" style="width:100%">
            <el-option
              v-for="(item, index) in templateOption"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="templateDialogVisible = false">取 消</el-button>
        <el-button :disabled="templateFormProps.disabled" type="primary" @click="install">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { getList, add, edit, remove, install } from '@/api/server'
import { getOption as getGroupOption } from '@/api/group'
import { getOption as getTemplateOption } from '@/api/template'
import { parseTime } from '@/utils'

export default {
  data() {
    return {
      dialogVisible: false,
      templateDialogVisible: false,
      tableData: [],
      groupOption: [],
      templateOption: [],
      tempFormData: {},
      formProps: {
        disabled: false
      },
      formData: {
        id: 0,
        name: '',
        ip: '',
        port: 22,
        owner: '',
        groupId: 0
      },
      formRules: {
        name: [
          { required: true, message: '请输入服务器名称', trigger: 'blur' }
        ],
        ip: [
          { required: true, message: '请输入服务器ip', trigger: 'blur' }
        ],
        port: [
          { required: true, message: '请输入服务器port', trigger: 'blur' }
        ],
        owner: [
          { required: true, message: '请输入SSH-KEY的所有者', trigger: 'blur' }
        ]
      },
      templateFormProps: {
        disabled: false
      },
      templateFormData: {
        templateId: '',
        serverId: 0
      },
      templateFormRules: {
        templateId: [
          { required: true, message: '请选择模板', trigger: 'change' }
        ]
      }
    }
  },
  created() {
    this.storeFormData()
    this.getList()
    this.getGroupOption()
    this.getTemplateOption()
  },
  methods: {
    getList() {
      getList().then((response) => {
        const serverList = response.data.serverList || []
        serverList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
        this.tableData = serverList
      })
    },

    getGroupOption() {
      getGroupOption().then((response) => {
        this.groupOption = response.data.groupList || []
      })
    },

    getTemplateOption() {
      getTemplateOption().then((response) => {
        this.templateOption = response.data.templateList || []
      })
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
      this.$confirm('此操作将删除该服务器, 是否继续?', '提示', {
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
          this.getList()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    handleInstall(data) {
      this.templateFormData.serverId = data.id
      this.templateDialogVisible = true
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
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.getList()
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    install() {
      this.$refs.templateForm.validate((valid) => {
        if (valid) {
          this.templateFormProps.disabled = true
          install(this.templateFormData.serverId, this.templateFormData.templateId).then((response) => {
            this.$message({
              message: response.message,
              duration: 5 * 1000
            })
          }).finally(() => {
            this.templateFormProps.disabled = this.templateDialogVisible = false
          })
        } else {
          return false
        }
      })
    },

    findGroupName(groupId) {
      const group = this.groupOption.find(element => element.id === groupId)
      return group ? group['name'] : '默认'
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
