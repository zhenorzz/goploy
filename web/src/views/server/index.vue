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
      <el-table-column prop="ip" label="Host" min-width="140">
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
          <el-tooltip class="item" effect="dark" content="Connect Terminal" placement="bottom">
            <el-button type="success" icon="el-icon-connection" @click="handleConnect(scope.row)" />
          </el-tooltip>
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
      <el-form ref="form" v-loading="formProps.loading" :rules="formRules" :model="formData" label-width="130px">
        <el-form-item :label="$t('namespace')" prop="namespaceId">
          <el-radio-group v-model="formData.namespaceId">
            <el-radio :label="getNamespace()['id']">当前</el-radio>
            <el-radio :label="0">不限</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Host" prop="ip">
          <el-input v-model="formData.ip" autocomplete="off" />
        </el-form-item>
        <el-form-item label="port" prop="port">
          <el-input v-model.number="formData.port" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('serverPage.sshKeyOwner')" prop="owner">
          <el-input v-model="formData.owner" autocomplete="off" placeholder="root" />
        </el-form-item>
        <el-form-item :label="$t('serverPage.sshKeyPath')" prop="path">
          <el-row type="flex">
            <el-input v-model="formData.path" autocomplete="off" placeholder="/root/.ssh/id_rsa" />
            <el-button
              :icon="'el-icon-copy-document'"
              type="success"
              :disabled="formData.path===''"
              @click="getPublicKey"
            >{{ $t('serverPage.copyPub') }}</el-button>
          </el-row>
        </el-form-item>
        <el-form-item :label="$t('serverPage.sshKeyPassword')" prop="password">
          <el-input v-model="formData.password" autocomplete="off" placeholder="" />
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
    <el-drawer
      ref="drawer"
      :title="term.title"
      :visible.sync="dialogTermVisible"
      @opened="connectTerminal"
      @closed="closeTerminal"
    >
      <div v-if="dialogTermVisible" ref="xterm" class="xterm" />
    </el-drawer>
  </el-row>
</template>
<script>

import { getNamespace } from '@/utils/namespace'
import { getList, getTotal, getPublicKey, add, edit, check, remove } from '@/api/server'
// require component
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { AttachAddon } from 'xterm-addon-attach'
// require styles
import 'xterm/css/xterm.css'

export default {
  name: 'Server',
  data() {
    return {
      dialogTermVisible: false,
      dialogVisible: false,
      tableLoading: false,
      tableData: [],
      pagination: {
        page: 1,
        rows: 16,
        total: 0
      },
      term: {
        window: null,
        ws: null,
        server: {},
        title: ''
      },
      tempFormData: {},
      formProps: {
        loading: false,
        disabled: false
      },
      formData: {
        id: 0,
        namespaceId: '',
        name: '',
        ip: '',
        port: 22,
        owner: '',
        path: '/root/.ssh/id_rsa',
        password: '',
        description: ''
      },
      formRules: {
        namespaceId: [
          { required: true, message: 'Namespace required', trigger: 'blur' }
        ],
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
        path: [
          { required: true, message: 'SSH-KEY path required', trigger: 'blur' }
        ],
        description: [
          { max: 255, message: 'Max 255 characters', trigger: 'blur' }
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
    getNamespace,
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

    getPublicKey() {
      getPublicKey(this.formData.path).then((response) => {
        this.copy(response.data, this.$t('serverPage.copyPubTips'))
      })
    },

    // 分页事件
    handlePageChange(val) {
      this.pagination.page = val
      this.getList()
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

    handleConnect(data) {
      this.term.server = data
      this.term.title = `${data.name}(${data.description})`
      this.dialogTermVisible = true
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

    connectTerminal() {
      const isWindows = ['Windows', 'Win16', 'Win32', 'WinCE'].indexOf(navigator.platform) >= 0
      const term = new Terminal({
        fontSize: 14,
        cursorBlink: true,
        windowsMode: isWindows
      })
      const fitAddon = new FitAddon()
      term.loadAddon(fitAddon)
      term.open(this.$refs['xterm'])
      fitAddon.fit()
      term.focus()
      this.term.ws = new WebSocket(`${location.protocol.replace('http', 'ws')}//${window.location.host + process.env.VUE_APP_BASE_API}/ws/xterm?serverId=${this.term.server.id}&rows=${term.rows}&cols=${term.cols}`)
      const attachAddon = new AttachAddon(this.term.ws)
      term.loadAddon(attachAddon)
      this.term.window = term
      this.term.ws.onopen = () => {
      }
      this.term.ws.onerror = () => {
        this.term.ws = null
      }
    },

    closeTerminal() {
      this.term.window = null
      this.term.ws.close()
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
.xterm {
  width: 100%;
  height: 100%;
}
</style>
