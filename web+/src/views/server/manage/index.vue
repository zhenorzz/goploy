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
      <el-table-column
        prop="description"
        :label="$t('desc')"
        min-width="140"
        show-overflow-tooltip
      />
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
        fixed="right"
      >
        <template #default="scope">
          <el-tooltip
            class="item"
            effect="dark"
            content="Connect Terminal"
            placement="bottom"
          >
            <el-button
              type="success"
              icon="el-icon-connection"
              @click="handleConnect(scope.row)"
            />
          </el-tooltip>
          <el-button
            type="primary"
            icon="el-icon-edit"
            @click="handleEdit(scope.row)"
          />
          <el-button
            type="danger"
            icon="el-icon-delete"
            @click="handleRemove(scope.row)"
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
    <el-dialog v-model="dialogVisible" :title="$t('setting')">
      <el-form
        ref="form"
        v-loading="formProps.loading"
        :rules="formRules"
        :model="formData"
        label-width="130px"
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
        <el-form-item label="port" prop="port">
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
        <el-form-item :label="$t('serverPage.sshKeyPassword')" prop="password">
          <el-input
            v-model="formData.password"
            autocomplete="off"
            placeholder=""
          />
        </el-form-item>
        <el-form-item :label="$t('desc')" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :autosize="{ minRows: 2 }"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-row type="flex" justify="space-between">
          <el-button
            :loading="formProps.loading"
            type="success"
            @click="check"
            >{{ $t('serverPage.testConnection') }}</el-button
          >
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
    <TheXtermDrawer v-model="dialogTermVisible" :server-row="selectedItem" />
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
  ServerRemove,
  ServerData,
} from '@/api/server'
import TheXtermDrawer from './TheXtermDrawer.vue'
import Validator from 'async-validator'
import { defineComponent } from 'vue'
import { copy } from '@/utils'
import { ElMessageBox, ElMessage } from 'element-plus'

export default defineComponent({
  name: 'ServerIndex',
  components: { TheXtermDrawer },
  data() {
    return {
      dialogTermVisible: false,
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

    handleEdit(data: ServerData['datagram']['detail']) {
      this.formData = Object.assign({}, data)
      this.dialogVisible = true
    },

    handleRemove(data: ServerData['datagram']['detail']) {
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
          new ServerRemove({ id: data.id }).request().then(() => {
            ElMessage.success('Success')
            this.getList()
            this.getTotal()
          })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },

    handleConnect(data: ServerData['datagram']) {
      this.selectedItem = data
      this.dialogTermVisible = true
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
