<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button @click="dialogSftpVisible = true">
        <svg-icon icon-class="ftp" />
      </el-button>
      <el-button @click="dialogTermVisible = true">
        <svg-icon icon-class="terminal" />
      </el-button>
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
      <el-table-column prop="ip" label="Host" min-width="140" sortable>
        <template #default="scope">
          {{ scope.row.ip }}:{{ scope.row.port }}
        </template>
      </el-table-column>
      <el-table-column
        prop="owner"
        :label="$t('serverPage.sshKeyOwner')"
        width="120"
        show-overflow-tooltip
      />
      <el-table-column label="OS" min-width="100" show-overflow-tooltip>
        <template #default="scope">
          <svg-icon
            v-if="scope.row.osInfo !== ''"
            :icon-class="getOSIcon(scope.row.osInfo)"
          />
          {{ getOS(scope.row.osInfo) }}
        </template>
      </el-table-column>
      <el-table-column
        prop="osInfo"
        label="CPU/MEMORY"
        min-width="100"
        show-overflow-tooltip
      >
        <template #default="scope">
          {{ getOSDetail(scope.row.osInfo) }}
        </template>
      </el-table-column>
      <el-table-column
        prop="description"
        :label="$t('desc')"
        min-width="140"
        show-overflow-tooltip
      />
      <el-table-column
        prop="state"
        :label="$t('state')"
        width="70"
        align="center"
      >
        <template #default="scope">
          <el-switch
            :value="scope.row.state === 1"
            active-color="#13ce66"
            inactive-color="#ff4949"
            @change="(value) => onSwitchState(value, scope.$index)"
          >
          </el-switch>
        </template>
      </el-table-column>
      <el-table-column
        prop="insertTime"
        :label="$t('insertTime')"
        width="135"
        align="center"
      />
      <el-table-column
        prop="updateTime"
        :label="$t('updateTime')"
        width="135"
        align="center"
      />
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="180"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button icon="el-icon-time" @click="handleCron(scope.row)" />
          <el-button
            icon="el-icon-data-line"
            @click="handleMonitor(scope.row)"
          />
          <el-button
            type="primary"
            icon="el-icon-edit"
            @click="handleEdit(scope.row)"
          />
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px">
      <el-pagination
        hide-on-single-page
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
    <el-dialog
      v-model="dialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('setting')"
      :close-on-click-modal="false"
    >
      <el-form
        ref="form"
        v-loading="formProps.loading"
        :rules="formRules"
        :model="formData"
        label-width="130px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
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
        <el-form-item label="Port" prop="port">
          <el-input v-model.number="formData.port" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('serverPage.sshKeyOwner')" prop="owner">
          <el-input
            v-model="formData.owner"
            autocomplete="off"
            placeholder="root"
          />
        </el-form-item>
        <el-form-item :label="$t('serverPage.sshKeyPath')" prop="path">
          <el-row type="flex">
            <el-input
              v-model="formData.path"
              style="flex: 1"
              autocomplete="off"
              placeholder="/root/.ssh/id_rsa"
            />
            <el-button
              :icon="'el-icon-copy-document'"
              type="success"
              :loading="formProps.copyPubLoading"
              :disabled="formData.path === ''"
              @click="getPublicKey"
            >
              {{ $t('serverPage.copyPub') }}
            </el-button>
          </el-row>
        </el-form-item>
        <!-- <el-form-item :label="$t('serverPage.sshKeyPassword')" prop="password">
          <el-input
            v-model="formData.password"
            autocomplete="off"
            placeholder=""
          />
        </el-form-item> -->
        <el-form-item :label="$t('desc')" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :autosize="{ minRows: 2 }"
          />
        </el-form-item>
        <el-form-item label="">
          <el-button
            type="text"
            @click="formProps.showAdvance = !formProps.showAdvance"
          >
            高级选项
          </el-button>
        </el-form-item>
        <el-form-item v-show="formProps.showAdvance" label="Jump host">
          <el-input v-model="formData.jumpIP" autocomplete="off" />
        </el-form-item>
        <el-form-item v-show="formProps.showAdvance" label="Jump port">
          <el-input v-model.number="formData.jumpPort" autocomplete="off" />
        </el-form-item>
        <el-form-item
          v-show="formProps.showAdvance"
          :label="$t('serverPage.sshKeyOwner')"
        >
          <el-input
            v-model="formData.jumpOwner"
            autocomplete="off"
            placeholder="root"
          />
        </el-form-item>
        <el-form-item
          v-show="formProps.showAdvance"
          :label="$t('serverPage.sshKeyPath')"
        >
          <el-input
            v-model="formData.jumpPath"
            autocomplete="off"
            placeholder="/root/.ssh/id_rsa"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-row type="flex" justify="space-between">
          <el-button :loading="formProps.loading" type="success" @click="check">
            {{ $t('serverPage.testConnection') }}
          </el-button>
          <el-row>
            <el-button @click="dialogVisible = false">
              {{ $t('cancel') }}
            </el-button>
            <el-button
              :disabled="formProps.disabled"
              type="primary"
              @click="submit"
            >
              {{ $t('confirm') }}
            </el-button>
          </el-row>
        </el-row>
      </template>
    </el-dialog>
    <TheXtermDialog v-model="dialogTermVisible" />
    <TheSftpDialog v-model="dialogSftpVisible" />
  </el-row>
</template>
<script lang="ts">
import { getNamespace } from '@/utils/namespace'
import {
  ServerList,
  ServerTotal,
  ServerPublicKey,
  ServerAdd,
  ServerEdit,
  ServerCheck,
  ServerToggle,
  ServerData,
} from '@/api/server'
import TheXtermDialog from './TheXtermDialog.vue'
import TheSftpDialog from './TheSftpDialog.vue'
import Validator from 'async-validator'
import { defineComponent } from 'vue'
import { copy, humanSize } from '@/utils'
import { ElMessageBox, ElMessage } from 'element-plus'

export default defineComponent({
  name: 'ServerIndex',
  components: { TheXtermDialog, TheSftpDialog },
  data() {
    return {
      dialogTermVisible: false,
      dialogSftpVisible: false,
      dialogVisible: false,
      tableLoading: false,
      tableData: [] as ServerList['datagram']['list'],
      selectedItem: {},
      pagination: {
        page: 1,
        rows: 16,
        total: 0,
      },
      tempFormData: {},
      formProps: {
        loading: false,
        showAdvance: false,
        copyPubLoading: false,
        disabled: false,
      },
      formData: {
        id: 0,
        namespaceId: getNamespace()['id'],
        name: '',
        ip: '',
        port: 22,
        owner: '',
        path: '/root/.ssh/id_rsa',
        password: '',
        jumpIP: '',
        jumpPort: 0,
        jumpOwner: '',
        jumpPath: '',
        jumpPassword: '',
        description: '',
      },
      formRules: {
        namespaceId: [
          { required: true, message: 'Namespace required', trigger: 'blur' },
        ],
        name: [{ required: true, message: 'Name required', trigger: 'blur' }],
        ip: [{ required: true, message: 'IP required', trigger: 'blur' }],
        port: [
          {
            type: 'number',
            required: true,
            min: 0,
            max: 65535,
            message: '0 ~ 65535',
            trigger: 'blur',
          },
        ],
        owner: [
          {
            required: true,
            message: 'SSH-KEY owner required',
            trigger: 'blur',
          },
        ],
        path: [
          { required: true, message: 'SSH-KEY path required', trigger: 'blur' },
        ],
        description: [
          { max: 255, message: 'Max 255 characters', trigger: 'blur' },
        ],
      },
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
      new ServerList(this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
    },

    getTotal() {
      new ServerTotal().request().then((response) => {
        this.pagination.total = response.data.total
      })
    },

    getPublicKey() {
      this.formProps.copyPubLoading = true
      new ServerPublicKey({ path: this.formData.path })
        .request()
        .then((response) => {
          copy(response.data.key)
          ElMessage.success(this.$t('serverPage.copyPubTips'))
        })
        .finally(() => {
          this.formProps.copyPubLoading = false
        })
    },

    // 分页事件
    handlePageChange(val = 1) {
      this.pagination.page = val
      this.getList()
    },

    handleAdd() {
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleEdit(data: ServerData['datagram']) {
      this.formData = Object.assign({}, data)
      this.dialogVisible = true
    },

    handleCron(data: ServerData['datagram']) {
      this.$router.push({ path: '/server/cron', query: { serverId: data.id } })
    },

    handleMonitor(data: ServerData['datagram']) {
      this.$router.push({ path: '/server/agent', query: { serverId: data.id } })
    },

    onSwitchState(value: boolean, index: number) {
      const data = this.tableData[index]
      if (value) {
        new ServerToggle({ id: data.id, state: value ? 1 : 0 })
          .request()
          .then(() => {
            ElMessage.success('Need to bind project again')
            this.tableData[index].state = value ? 1 : 0
          })
      } else {
        ElMessageBox.confirm(
          this.$t('serverPage.removeServerTips', { serverName: data.name }),
          this.$t('tips'),
          {
            confirmButtonText: this.$t('confirm'),
            cancelButtonText: this.$t('cancel'),
            type: 'warning',
          }
        )
          .then(() => {
            new ServerToggle({ id: data.id, state: value ? 1 : 0 })
              .request()
              .then(() => {
                this.tableData[index].state = value ? 1 : 0
              })
          })
          .catch(() => {
            ElMessage.info('Cancel')
          })
      }
    },

    check() {
      ;(this.$refs.form as Validator).validate((valid: boolean) => {
        if (valid) {
          this.formProps.loading = true
          this.formProps.disabled = true
          new ServerCheck(this.formData)
            .request()
            .then(() => {
              ElMessage.success('Connected')
            })
            .finally(() => {
              this.formProps.loading = false
              this.formProps.disabled = false
            })
        } else {
          return false
        }
      })
    },

    submit() {
      ;(this.$refs.form as Validator).validate((valid: boolean) => {
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
      new ServerAdd(this.formData)
        .request()
        .then(() => {
          this.getList()
          this.getTotal()
          ElMessage.success('Success')
        })
        .finally(() => {
          this.formProps.disabled = this.dialogVisible = false
        })
    },

    edit() {
      this.formProps.disabled = true
      new ServerEdit(this.formData)
        .request()
        .then(() => {
          this.getList()
          ElMessage.success('Success')
        })
        .finally(() => {
          this.formProps.disabled = this.dialogVisible = false
        })
    },

    storeFormData() {
      this.tempFormData = JSON.parse(JSON.stringify(this.formData))
    },

    restoreFormData() {
      this.formData = JSON.parse(JSON.stringify(this.tempFormData))
    },

    getOS(osInfo: string): string {
      if (osInfo === '') return ''
      return osInfo.split('|')[0]
    },

    getOSIcon(osInfo: string): string {
      if (osInfo === '') return ''
      else if (osInfo.toLowerCase().includes('centos')) return 'centos'
      else if (osInfo.toLowerCase().includes('ubuntu')) return 'ubuntu'
      else return 'question-mark-blue'
    },

    getOSDetail(osInfo: string): string {
      if (osInfo === '') return ''
      const osArr = osInfo.split('|')
      return osArr[1] + ' cores ' + humanSize(Number(osArr[2]) * 1024)
    },
  },
})
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.template-dialog {
  padding-right: 10px;
  height: 400px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
