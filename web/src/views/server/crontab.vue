<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input v-model="crontabCommand" style="width:200px" placeholder="请输入命令关键词" />
        <el-button type="primary" icon="el-icon-search" @click="searchList">搜索</el-button>
      </el-row>
      <el-row>
        <el-button type="primary" icon="el-icon-download" @click="handleImport">导入</el-button>
        <el-button type="primary" icon="el-icon-plus" @click="handleAdd">添加</el-button>
      </el-row>
    </el-row>
    <el-table
      :max-height="tableHeight"
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="command" label="命令" min-width="140" show-overflow-tooltip />
      <el-table-column prop="description" label="描述" min-width="240" show-overflow-tooltip />
      <el-table-column prop="creator" label="创建人" min-width="50" />
      <el-table-column prop="editor" label="修改人" min-width="50" />
      <el-table-column prop="insertTime" label="创建时间" width="135" />
      <el-table-column prop="updateTime" label="更新时间" width="135" />
      <el-table-column prop="operation" label="操作" width="255" fixed="right">
        <template slot-scope="scope">
          <el-button type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button type="warning" @click="handleServer(scope.row)">查看服务器</el-button>
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
    <el-dialog title="crontab设置" :visible.sync="dialogVisible">
      <el-form ref="form" v-loading="formProps.loading" :rules="formRules" :model="formData" label-width="80px">
        <el-form-item label="时间" prop="date">
          <el-input v-model="formData.date" autocomplete="off" placeholder="* * * * ?" @change="onDateChange" />
          <span>{{ formProps.dateCN }}</span>
        </el-form-item>
        <el-form-item label="脚本" prop="script">
          <el-input v-model.trim="formData.script" autocomplete="off" />
        </el-form-item>
        <el-form-item v-show="formData.id === 0" label="服务器" prop="serverIds">
          <el-select v-model="formData.serverIds" multiple placeholder="选择服务器，可多选" style="width:100%" filterable>
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.label"
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
    <el-dialog title="crontab移除" :visible.sync="crontabRemoveVisible" width="400px">
      <el-form ref="crontabRemoveForm" :model="crontabRemoveFormData" label-width="80px">
        <el-form-item label="命令" label-width="45px">
          <span>{{ crontabRemoveFormProps.command }}</span>
        </el-form-item>
        <el-form-item label="删除服务器Crontab任务" label-width="170px">
          <el-radio-group v-model="crontabRemoveFormData.radio">
            <el-radio :label="0">否</el-radio>
            <el-radio :label="1">是</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="crontabRemoveVisible = false">取 消</el-button>
        <el-button :disabled="crontabRemoveFormProps.disabled" type="primary" @click="remove">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="crontab导入" :visible.sync="importVisible">
      <el-row>
        <el-row type="flex">
          <el-select v-model="importProps.serverId" placeholder="选择服务器" style="width:100%;margin-right:5px;">
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.label"
              :value="item.id"
            />
          </el-select>
          <el-button :disabled="importProps.disabled" :icon="importProps.disabled ? 'el-icon-loading' : 'el-icon-search'" type="primary" @click="getRemoteServerList">读取</el-button>
        </el-row>
        <el-table
          border
          stripe
          highlight-current-row
          empty-text="请先选择服务器读取"
          :data="serverTableData"
          style="width: 100%; margin-top: 10px;"
          @selection-change="handleCrontabSelectionChange"
        >
          <el-table-column
            type="selection"
            width="40"
          />
          <el-table-column prop="command" label="命令" min-width="140" show-overflow-tooltip />
          <el-table-column prop="description" label="描述" min-width="240" show-overflow-tooltip />
        </el-table>
      </el-row>
      <div slot="footer" class="dialog-footer">
        <el-button @click="importVisible = false">取 消</el-button>
        <el-button :disabled="importProps.disabled" type="primary" @click="importCrontab">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="服务器管理" :visible.sync="serverVisible">
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
        <el-table-column prop="serverId" label="服务器ID" width="100" />
        <el-table-column prop="serverName" label="服务器名称" width="100" />
        <el-table-column prop="serverDescription" label="服务器描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="insertTime" width="160" label="绑定时间" />
        <el-table-column prop="updateTime" width="160" label="更新时间" />
        <el-table-column prop="operation" label="操作" width="80">
          <template slot-scope="scope">
            <el-button type="danger" @click="removeCrontabServer(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="serverVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="添加服务器" :visible.sync="addServerVisible">
      <el-form ref="addServerForm" :rules="addServerFormRules" :model="addServerFormData">
        <el-form-item label="绑定服务器" label-width="120px" prop="serverIds">
          <el-select v-model="addServerFormData.serverIds" multiple placeholder="选择服务器，可多选">
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.label"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="addServerVisible = false">取 消</el-button>
        <el-button :disabled="addServerFormProps.disabled" type="primary" @click="addServer">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import tableHeight from '@/mixin/tableHeight'
import cronstrue from 'cronstrue/i18n'
import { getList, getTotal, getRemoteServerList, getBindServerList, add, edit, remove, importCrontab, addServer, removeCrontabServer } from '@/api/crontab'
import { getOption as getServerOption } from '@/api/server'

export default {
  name: 'Crontab',
  mixins: [tableHeight],
  data() {
    const validateDate = (rule, value, callback) => {
      try {
        cronstrue.toString(value, { locale: 'zh_CN' })
        callback()
      } catch (error) {
        callback(error)
      }
    }
    return {
      crontabCommand: '',
      dialogVisible: false,
      crontabRemoveVisible: false,
      serverVisible: false,
      addServerVisible: false,
      importVisible: false,
      selectedItems: [],
      tableData: [],
      tableServerData: [],
      serverOption: [],
      pagination: {
        page: 1,
        rows: 16,
        total: 0
      },
      formProps: {
        loading: false,
        disabled: false,
        dateCN: ''
      },
      formData: {
        id: 0,
        command: '',
        date: '',
        script: '',
        serverIds: []
      },
      formRules: {
        date: [
          { required: true, validator: validateDate, trigger: 'blur' }
        ],
        script: [
          { required: true, message: '请输入脚本', trigger: 'blur' }
        ]
      },
      importProps: {
        serverId: '',
        disabled: false,
        loading: false
      },
      crontabRemoveFormData: {
        id: 0,
        radio: 0
      },
      crontabRemoveFormProps: {
        command: '',
        disabled: false
      },
      addServerFormProps: {
        disabled: false
      },
      addServerFormData: {
        crontabId: 0,
        serverIds: []
      },
      addServerFormRules: {
        serverIds: [
          { type: 'array', required: true, message: '请选择服务器', trigger: 'change' }
        ]
      },
      serverTableData: []
    }
  },
  created() {
    this.getList()
    this.getTotal()
    this.getServerOption()
  },

  methods: {
    getList() {
      getList(this.pagination, this.crontabCommand).then((response) => {
        this.tableData = response.data.list.map(element => {
          const commandSplit = element.command.split(' ')
          element.date = commandSplit.slice(0, 5).join(' ')
          element.dateCN = cronstrue.toString(element.date, { locale: 'zh_CN' })
          element.script = commandSplit.slice(5).join(' ')
          element.description = element.dateCN + ', 运行: ' + element.script
          return element
        })
      })
    },

    getTotal() {
      getTotal(this.crontabCommand).then((response) => {
        this.pagination.total = response.data.total
      })
    },

    getServerOption() {
      getServerOption().then((response) => {
        this.serverOption = response.data.list
        this.serverOption.map(element => {
          element.label = element.name + (element.description.length > 0 ? '(' + element.description + ')' : '')
          return element
        })
      })
    },

    getRemoteServerList() {
      if (this.importProps.serverId <= 0) {
        this.$message.warning('请先选择服务器')
        return
      }
      this.importProps.disabled = true
      getRemoteServerList(this.importProps.serverId).then(response => {
        this.serverTableData = response.data.list.map(command => {
          const element = {}
          const commandSplit = command.split(' ')
          element.command = command
          element.date = commandSplit.slice(0, 5).join(' ')
          element.dateCN = cronstrue.toString(element.date, { locale: 'zh_CN' })
          element.script = commandSplit.slice(5).join(' ')
          element.description = element.dateCN + ', 运行: ' + element.script
          return element
        })
      }).finally(() => { this.importProps.disabled = false })
    },

    getBindServerList(crontabID) {
      getBindServerList(crontabID).then((response) => {
        this.tableServerData = response.data.list
      })
    },

    searchList() {
      this.pagination.page = 1
      this.getList()
      this.getTotal()
    },

    handleAdd() {
      this.formData.id = 0
      this.dialogVisible = true
    },

    handleImport() {
      this.importVisible = true
    },

    handleEdit(data) {
      this.formData.id = data.id
      this.formData.date = data.date
      this.formData.script = data.script
      this.formData.serverIds = []
      this.formProps.dateCN = data.dateCN
      this.dialogVisible = true
    },

    handleServer(data) {
      this.getBindServerList(data.id)
      this.addServerFormData.crontabId = data.id
      this.serverVisible = true
    },

    handleAddServer() {
      this.addServerVisible = true
    },

    handleRemove(data) {
      this.crontabRemoveFormData.id = data.id
      this.crontabRemoveFormProps.command = data.command
      this.crontabRemoveVisible = true
    },

    handleCrontabSelectionChange(items) {
      this.selectedItems = items
    },

    importCrontab() {
      if (this.selectedItems.length === 0) {
        this.$message.warning('请先选择需要导入的条目')
        return
      }
      importCrontab({ commands: this.selectedItems.map(element => element.command) }).then(response => {
        this.getList()
        this.getTotal()
        this.$message.success('导入成功')
      }).finally(() => { this.importVisible = false })
    },

    onDateChange() {
      this.formProps.dateCN = cronstrue.toString(this.formData.date, { locale: 'zh_CN' })
    },

    handlePageChange(val) {
      this.pagination.page = val
      this.getList()
    },

    submit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.formData.command = this.formData.date + ' ' + this.formData.script
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
        this.$message.success('编辑成功')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    remove() {
      this.crontabRemoveFormProps.disabled = true
      remove(this.crontabRemoveFormData).then((response) => {
        this.getList()
        this.getTotal()
        this.$message.success('删除成功')
      }).finally(() => {
        this.crontabRemoveFormProps.disabled = this.crontabRemoveVisible = false
      })
    },

    addServer() {
      this.$refs.addServerForm.validate((valid) => {
        if (valid) {
          this.addServerFormProps.disabled = true
          addServer(this.addServerFormData).then((response) => {
            this.addServerVisible = false
            this.$message.success('添加成功')
            this.getBindServerList(this.addServerFormData.crontabId)
          }).finally(() => {
            this.addServerFormProps.disabled = false
          })
        } else {
          return false
        }
      })
    },

    removeCrontabServer(data) {
      this.$confirm('此操作将永久删除该服务器的绑定关系, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        removeCrontabServer({ crontabServerId: data.id, crontabId: data.crontabId, serverId: data.serverId }).then((response) => {
          this.$message.success('删除成功')
          this.getBindServerList(data.crontabId)
        })
      }).catch(() => {
        this.$message.info('已取消删除')
      })
    }
  }
}
</script>
<style lang="scss" scoped>
@import "@/styles/mixin.scss";
.template-dialog {
  padding-right: 10px;
  height: 400px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
