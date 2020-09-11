<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
    </el-row>
    <el-table
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="100" />
      <el-table-column prop="name" :label="$t('name')" min-width="140" />
      <el-table-column prop="ip" label="IP" min-width="140">
        <template slot-scope="scope">
          {{ scope.row.ip }}:{{ scope.row.port }}
        </template>
      </el-table-column>
      <el-table-column prop="owner" :label="$t('serverPage.sshKeyOwner')" width="120" show-overflow-tooltip />
      <el-table-column prop="description" :label="$t('desc')" min-width="140" show-overflow-tooltip />
      <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
      <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
      <el-table-column prop="operation" :label="$t('op')" width="180" align="center" fixed="right">
        <template slot-scope="scope">
          <el-button type="primary" icon="el-icon-edit" @click="handleEdit(scope.row)" />
          <el-tooltip class="item" effect="dark" :content="$t('install')" placement="top">
            <el-button type="info" icon="el-icon-set-up" @click="handleInstall(scope.row)" />
          </el-tooltip>
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
      <el-form ref="form" v-loading="formProps.loading" :rules="formRules" :model="formData" label-width="130px">
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="IP" prop="ip">
          <el-input v-model="formData.ip" autocomplete="off" />
        </el-form-item>
        <el-form-item label="port" prop="port">
          <el-input v-model.number="formData.port" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('serverPage.sshKeyOwner')" prop="owner">
          <el-input v-model="formData.owner" autocomplete="off" placeholder="root" />
        </el-form-item>
        <el-form-item :label="$t('desc')" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :autosize="{ minRows: 2}"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-row type="flex" justify="space-between">
          <el-button type="success" @click="check">{{ $t('serverPage.testConnection') }}</el-button>
          <el-row>
            <el-button @click="dialogVisible = false">{{ $t('cancel') }}</el-button>
            <el-button :disabled="formProps.disabled" type="primary" @click="submit">{{ $t('confirm') }}</el-button>
          </el-row>
        </el-row>
      </div>
    </el-dialog>
    <el-dialog :title="$t('install')" :visible.sync="templateDialogVisible">
      <el-row class="template-dialog">
        <el-form ref="templateForm" :rules="templateFormRules" :model="templateFormData" label-width="90px">
          <el-form-item :label="$t('template')" prop="templateId">
            <el-select v-model="templateFormData.templateId" style="width:100%">
              <el-option
                v-for="(item, index) in templateOption"
                :key="index"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
        </el-form>
        <el-row>
          <el-collapse accordion @change="handleTokenChange">
            <el-collapse-item v-for="(item, index) in installPreviewList.slice().reverse()" :key="index" :name="item.token">
              <template slot="title">
                <span style="margin-right: 10px">token: {{ item.token }}</span>
                <el-tag v-if="item.installState === 1" type="success" effect="plain">{{ $t('success') }}</el-tag>
                <el-tag v-else type="danger" effect="plain">{{ $t('fail') }}</el-tag>
              </template>
              <codemirror :value="installTrace" :options="cmOptions" />
            </el-collapse-item>
          </el-collapse>
        </el-row>
      </el-row>
      <div slot="footer" class="dialog-footer">
        <el-button @click="templateDialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="templateFormProps.disabled" type="primary" @click="install">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('detail')" :visible.sync="installDialogVisible">
      <codemirror :value="installLog" :options="cmOptions" />
    </el-dialog>
  </el-row>
</template>
<script>
import { getList, getTotal, getInstallPreview, getInstallList, add, edit, check, remove, install } from '@/api/server'
import { getOption as getTemplateOption } from '@/api/template'
// require component
import { codemirror } from 'vue-codemirror'
import 'codemirror/mode/shell/shell.js'
import 'codemirror/theme/darcula.css'
// require styles
import 'codemirror/lib/codemirror.css'
import 'codemirror/addon/scroll/simplescrollbars.js'
import 'codemirror/addon/scroll/simplescrollbars.css'
import 'codemirror/addon/display/autorefresh.js'

export default {
  name: 'Server',
  components: {
    codemirror
  },
  data() {
    return {
      dialogVisible: false,
      templateDialogVisible: false,
      installDialogVisible: false,
      tableLoading: false,
      tableData: [],
      pagination: {
        page: 1,
        rows: 16,
        total: 0
      },
      templateOption: [],
      installToken: '',
      installPreviewList: [],
      installTraceList: [],
      installLog: '',
      tempFormData: {},
      cmOptions: {
        tabSize: 4,
        mode: 'text/x-sh',
        lineNumbers: false,
        line: false,
        scrollbarStyle: 'overlay',
        theme: 'darcula',
        readOnly: true,
        autoRefresh: true
      },
      formProps: {
        loading: false,
        disabled: false
      },
      formData: {
        id: 0,
        name: '',
        ip: '',
        port: 22,
        owner: '',
        description: ''
      },
      formRules: {
        name: [
          { required: true, message: 'Name required', trigger: 'blur' }
        ],
        ip: [
          { required: true, message: 'IP required', trigger: 'blur' }
        ],
        port: [
          { type: 'number', required: true, min: 0, max: 65535, message: '0 ~ 65535', trigger: 'blur' }
        ],
        owner: [
          { required: true, message: 'SSH-KEY owner required', trigger: 'blur' }
        ],
        description: [
          { max: 255, message: 'Max 255 characters', trigger: 'blur' }
        ]
      },
      templateFormProps: {
        disabled: false
      },
      templateFormData: {
        templateId: '',
        serverName: '',
        serverId: 0
      },
      templateFormRules: {
        templateId: [
          { required: true, message: 'Template required', trigger: 'change' }
        ]
      }
    }
  },

  computed: {
    installTrace: function() {
      let intallTrace = ''
      this.installTraceList.forEach(element => {
        if (element.type === 1) {
          intallTrace += '[goploy~]$ ' + element.command + '\n'
          intallTrace += element.detail + '\n'
        } else if (element.type === 2) {
          intallTrace += '[goploy~]$ ' + element.ssh + '\n'
          intallTrace += element.detail + '\n'
        } else if (element.type === 3) {
          intallTrace += '[' + this.templateFormData.serverName + '~]$ ' + element.script + '\n'
          intallTrace += element.detail + '\n'
        }
      })
      return intallTrace
    }
  },

  watch: {
    '$store.getters.ws_message': function(response) {
      if (response.type !== 2) {
        return
      }
      const data = response.message
      Object.assign(data, JSON.parse(data.ext))
      let intallTrace = ''
      if (data.type === 1) {
        intallTrace += '[goploy~]$ ' + data.command + '\n'
        intallTrace += data.detail + '\n'
      } else if (data.type === 2) {
        intallTrace += '[goploy~]$ ' + data.ssh + '\n'
        intallTrace += data.detail + '\n'
      } else if (data.type === 3) {
        intallTrace += '[' + this.templateFormData.serverName + '~]$ ' + data.script + '\n'
        intallTrace += data.detail + '\n'
      }
      this.installLog += intallTrace
    }
  },

  created() {
    this.storeFormData()
    this.getList()
    this.getTotal()
    this.getTemplateOption()
  },

  methods: {
    getList() {
      this.tableLoading = true
      getList(this.pagination).then((response) => {
        this.tableData = response.data.list
      }).finally(() => {
        this.tableLoading = false
      })
    },

    getTotal() {
      getTotal().then((response) => {
        this.pagination.total = response.data.total
      })
    },

    // 分页事件
    handlePageChange(val) {
      this.pagination.page = val
      this.getList()
    },

    getTemplateOption() {
      getTemplateOption().then((response) => {
        this.templateOption = response.data.list
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
      this.$confirm(this.$i18n.t('serverPage.removeServerTips', { serverName: data.name }), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message.success('Success')
          this.getList()
          this.getTotal()
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    handleInstall(data) {
      this.templateFormData.serverId = data.id
      this.templateFormData.serverName = data.name
      getInstallPreview(data.id).then(response => {
        this.installPreviewList = response.data.installTraceList || []
      })
      this.templateDialogVisible = true
    },

    handleTokenChange(token) {
      if (token === '') return
      getInstallList(token).then(response => {
        const installTraceList = response.data.installTraceList || []
        this.installTraceList = installTraceList.map(element => {
          Object.assign(element, JSON.parse(element.ext))
          return element
        })
      })
    },

    check() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.formProps.loading = true
          this.formProps.disabled = true
          check(this.formData).then(response => {
            this.$message.success('Connected')
          }).finally(_ => {
            this.formProps.loading = false
            this.formProps.disabled = false
          })
        } else {
          return false
        }
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

    install() {
      this.$refs.templateForm.validate((valid) => {
        if (valid) {
          this.templateFormProps.disabled = true
          this.installDialogVisible = true
          this.templateFormProps.disabled = this.templateDialogVisible = false
          this.installLog = ''
          install(this.templateFormData.serverId, this.templateFormData.templateId).then((response) => {
            this.$message.success(response.message)
          })
        } else {
          return false
        }
      })
    },

    formatDetail(detail) {
      return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
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
<style lang="scss" scoped>
@import "@/styles/mixin.scss";
.template-dialog {
  padding-right: 10px;
  height: 400px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
