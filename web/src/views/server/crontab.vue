<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input v-model="crontabCommand" style="width:200px" placeholder="Fitler the command" />
        <el-button type="primary" icon="el-icon-search" @click="searchList" />
      </el-row>
      <el-row>
        <el-tooltip effect="dark" :content="$t('import')" placement="top">
          <el-button type="primary" icon="el-icon-download" @click="handleImport" />
        </el-tooltip>
        <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
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
      <el-table-column prop="command" :label="$t('command')" min-width="140" show-overflow-tooltip />
      <el-table-column prop="description" :label="$t('desc')" min-width="240" show-overflow-tooltip />
      <el-table-column prop="creator" :label="$t('creator')" min-width="50" />
      <el-table-column prop="editor" :label="$t('editor')" min-width="50" />
      <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
      <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
      <el-table-column prop="server" width="80" :label="$t('server')" align="center">
        <template slot-scope="scope">
          <el-button type="text" @click="handleServer(scope.row)">{{ $t('view') }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="operation" :label="$t('op')" width="130" align="center" fixed="right">
        <template slot-scope="scope">
          <el-button type="primary" icon="el-icon-edit" @click="handleEdit(scope.row)" />
          <el-button type="danger" icon="el-icon-delete" @click="handleRemove(scope.row)" />
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
    <el-dialog :title="$t('setting')" :visible.sync="dialogVisible">
      <el-form ref="form" v-loading="formProps.loading" :rules="formRules" :model="formData" label-width="80px">
        <el-form-item :label="$t('time')" prop="date">
          <el-input v-model="formData.date" autocomplete="off" placeholder="* * * * ?" @change="onDateChange" />
          <span>{{ formProps.dateLocale }}</span>
        </el-form-item>
        <el-form-item :label="$t('script')" prop="script">
          <el-input v-model.trim="formData.script" autocomplete="off" />
        </el-form-item>
        <el-form-item v-show="formData.id === 0" :label="$t('server')" prop="serverIds">
          <el-select v-model="formData.serverIds" multiple style="width:100%" filterable>
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
        <el-button @click="dialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('remove')" :visible.sync="crontabRemoveVisible" width="400px">
      <el-form ref="crontabRemoveForm" :model="crontabRemoveFormData" label-width="80px">
        <el-form-item :label="$t('command')" label-width="45px">
          <span>{{ crontabRemoveFormProps.command }}</span>
        </el-form-item>
        <el-form-item :label="$t('crontabPage.removeServerCrontabLabel')" label-width="170px">
          <el-radio-group v-model="crontabRemoveFormData.radio">
            <el-radio :label="0">{{ $t('boolOption[0]') }}</el-radio>
            <el-radio :label="1">{{ $t('boolOption[1]') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="crontabRemoveVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="crontabRemoveFormProps.disabled" type="primary" @click="remove">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('import')" :visible.sync="importVisible">
      <el-row>
        <el-row type="flex">
          <el-select v-model="importProps.serverId" style="width:100%;margin-right:5px;">
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.label"
              :value="item.id"
            />
          </el-select>
          <el-button
            :disabled="importProps.disabled"
            :icon="importProps.disabled ? 'el-icon-loading' : 'el-icon-refresh'"
            type="primary"
            @click="getRemoteServerList"
          >{{ $t('read') }}</el-button>
        </el-row>
        <el-table
          border
          stripe
          highlight-current-row
          :empty-text="$t('crontabPage.importTips')"
          :data="serverTableData"
          style="width: 100%; margin-top: 10px;"
          @selection-change="handleCrontabSelectionChange"
        >
          <el-table-column
            type="selection"
            width="40"
          />
          <el-table-column prop="command" :label="$t('command')" min-width="140" show-overflow-tooltip />
          <el-table-column prop="description" :label="$t('desc')" min-width="240" show-overflow-tooltip />
        </el-table>
      </el-row>
      <div slot="footer" class="dialog-footer">
        <el-button @click="importVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="importProps.disabled" type="primary" @click="importCrontab">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('manage')" :visible.sync="serverVisible">
      <el-row class="app-bar" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="handleAddServer" />
      </el-row>
      <el-table
        border
        stripe
        highlight-current-row
        :data="tableServerData"
        style="width: 100%"
      >
        <el-table-column prop="serverId" :label="$t('serverId')" width="100" />
        <el-table-column prop="serverName" :label="$t('serverName')" width="100" />
        <el-table-column prop="serverDescription" :label="$t('serverDescription')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="insertTime" :label="$t('insertTime')" width="160" align="center" />
        <el-table-column prop="updateTime" :label="$t('updateTime')" width="160" align="center" />
        <el-table-column prop="operation" :label="$t('op')" width="80" align="center">
          <template slot-scope="scope">
            <el-button type="danger" icon="el-icon-delete" @click="removeCrontabServer(scope.row)" />
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="serverVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('add')" :visible.sync="addServerVisible">
      <el-form ref="addServerForm" :rules="addServerFormRules" :model="addServerFormData">
        <el-form-item :label="$t('server')" label-width="120px" prop="serverIds">
          <el-select v-model="addServerFormData.serverIds" multiple>
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
        <el-button @click="addServerVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="addServerFormProps.disabled" type="primary" @click="addServer">{{ $t('confirm') }}</el-button>
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
        cronstrue.toString(value)
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
        dateLocale: ''
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
          { required: true, message: 'Script required', trigger: 'blur' }
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
          { type: 'array', required: true, message: 'Server required', trigger: 'change' }
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
          element.dateLocale = cronstrue.toString(element.date, { locale: this.getLocale() })
          element.script = commandSplit.slice(5).join(' ')
          element.description = `${element.dateLocale}, ${this.$i18n.t('run')}: ${element.script}`
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
        this.$message.warning(this.$i18n.t('crontabPage.selectServerTips'))
        return
      }
      this.importProps.disabled = true
      getRemoteServerList(this.importProps.serverId).then(response => {
        this.serverTableData = response.data.list.map(command => {
          const element = {}
          const commandSplit = command.split(' ')
          element.command = command
          element.date = commandSplit.slice(0, 5).join(' ')
          element.dateLocale = cronstrue.toString(element.date, { locale: this.getLocale() })
          element.script = commandSplit.slice(5).join(' ')
          element.description = `${element.dateLocale}, ${this.$i18n.t('run')}: ${element.script}`
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
      this.formProps.dateLocale = data.dateLocale
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
        this.$message.warning(this.$i18n.t('crontabPage.selectItemTips'))
        return
      }
      importCrontab({ commands: this.selectedItems.map(element => element.command) }).then(response => {
        this.getList()
        this.getTotal()
        this.$message.success('Success')
      }).finally(() => { this.importVisible = false })
    },

    onDateChange() {
      this.formProps.dateLocale = cronstrue.toString(this.formData.date, { locale: this.getLocale() })
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
        this.$message.success('Success')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.getList()
        this.$message.success('Success')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    remove() {
      this.crontabRemoveFormProps.disabled = true
      remove(this.crontabRemoveFormData).then((response) => {
        this.getList()
        this.getTotal()
        this.$message.success('Success')
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
            this.$message.success('Success')
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
      this.$confirm(this.$i18n.t('crontabPage.removeCrontabServerTips'), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        removeCrontabServer({ crontabServerId: data.id, crontabId: data.crontabId, serverId: data.serverId }).then((response) => {
          this.$message.success('Success')
          this.getBindServerList(data.crontabId)
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    getLocale() {
      if (this.$i18n.locale === 'zh') {
        return 'zh_CN'
      }
      return this.$i18n.locale
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
