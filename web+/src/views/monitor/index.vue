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
      <el-table-column prop="domain" label="Domain" min-width="100">
        <template #default="scope">
          {{ scope.row.domain }}:{{ scope.row.port }}
        </template>
      </el-table-column>
      <el-table-column
        prop="second"
        :label="$t('interval') + '(s)'"
        width="80"
      />
      <el-table-column
        prop="times"
        :label="$t('monitorPage.failTimes')"
        width="100"
      />
      <el-table-column prop="notifyType" :label="$t('notice')" width="70">
        <template #default="scope">
          <span v-if="scope.row.notifyType === 1">
            {{ $t('webhookOption[1]') }}
          </span>
          <span v-else-if="scope.row.notifyType === 2">
            {{ $t('webhookOption[2]') }}
          </span>
          <span v-else-if="scope.row.notifyType === 3">
            {{ $t('webhookOption[3]') }}
          </span>
          <span v-else-if="scope.row.notifyType === 255">
            {{ $t('webhookOption[255]') }}
          </span>
        </template>
      </el-table-column>
      <el-table-column
        prop="notifyTimes"
        :label="$t('monitorPage.notifyTimes')"
        width="80"
      />
      <el-table-column
        prop="state"
        :label="$t('state')"
        width="65"
        align="center"
      >
        <template #default="scope">
          <el-switch
            :value="scope.row.state === 1"
            active-color="#13ce66"
            inactive-color="#ff4949"
            @change="handleToggle(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column
        prop="errorContent"
        :label="$t('monitorPage.errorContent')"
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
        width="130"
        align="center"
        fixed="right"
      >
        <template #default="scope">
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
        label-width="120px"
      >
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Domain/IP" prop="domain">
          <el-input
            v-model="formData.domain"
            autocomplete="off"
            placeholder="Skip http(s)"
          />
        </el-form-item>
        <el-form-item label="port" prop="port">
          <el-input v-model.number="formData.port" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('interval') + '(s)'" prop="second">
          <el-input v-model.number="formData.second" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('monitorPage.failTimes')" prop="times">
          <el-input v-model.number="formData.times" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('notice')" prop="notifyTarget">
          <el-row type="flex">
            <el-select v-model="formData.notifyType">
              <el-option :label="$t('webhookOption[1]')" :value="1" />
              <el-option :label="$t('webhookOption[2]')" :value="2" />
              <el-option :label="$t('webhookOption[3]')" :value="3" />
              <el-option :label="$t('webhookOption[255]')" :value="255" />
            </el-select>
            <el-input
              v-model.trim="formData.notifyTarget"
              style="flex: 1"
              autocomplete="off"
              placeholder="webhook"
            />
          </el-row>
        </el-form-item>
        <el-form-item :label="$t('monitorPage.notifyTimes')" prop="notifyTimes">
          <el-input v-model.number="formData.notifyTimes" />
        </el-form-item>
        <el-form-item :label="$t('desc')" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :autosize="{ minRows: 2 }"
          />
        </el-form-item>
      </el-form>
      <template #footer class="dialog-footer">
        <el-row type="flex" justify="space-between">
          <el-button :loading="formProps.loading" type="success" @click="check">
            {{ $t('monitorPage.testAppState') }}
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
  </el-row>
</template>
<script>
import {
  getList,
  getTotal,
  add,
  edit,
  check,
  toggle,
  remove,
} from '@/api/monitor'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Monitor',
  data() {
    return {
      dialogVisible: false,
      tableLoading: false,
      tableData: [],
      pagination: {
        page: 1,
        rows: 16,
        total: 0,
      },
      tempFormData: {},
      formProps: {
        loading: false,
        disabled: false,
      },
      formData: {
        id: 0,
        name: '',
        domain: '',
        port: 80,
        second: 3,
        times: 1,
        notifyType: 1,
        notifyTarget: '',
        notifyTimes: 1,
      },
      formRules: {
        name: [{ required: true, message: 'Name required', trigger: 'blur' }],
        domain: [
          { required: true, message: 'Host or IP required', trigger: 'blur' },
        ],
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
        second: [
          {
            type: 'number',
            required: true,
            min: 1,
            message: 'Interval required',
            trigger: 'blur',
          },
        ],
        times: [
          {
            type: 'number',
            required: true,
            min: 1,
            max: 65535,
            message: 'Times required',
            trigger: 'blur',
          },
        ],
        notifyTarget: [{ required: true, message: 'Webhook required' }],
        notifyTimes: [
          {
            type: 'number',
            required: true,
            min: 1,
            max: 65535,
            message: 'Notify times required',
            trigger: 'blur',
          },
        ],
        description: [
          { max: 255, message: 'Max 255 characters', trigger: 'blur' },
        ],
      },
    }
  },

  watch: {
    '$store.getters.ws_message': function (response) {
      if (response.type !== 3) {
        return
      }
      const data = response.message
      const monitorIndex = this.tableData.findIndex(
        (element) => element.id === data.monitorId
      )
      if (monitorIndex !== -1) {
        this.tableData[monitorIndex].errorContent = data.errorContent
        this.tableData[monitorIndex].state = data.state
      }
    },
  },

  created() {
    this.storeFormData()
    this.getList()
    this.getTotal()
  },

  methods: {
    getList() {
      this.tableLoading = true
      getList(this.pagination)
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
    },

    getTotal() {
      getTotal().then((response) => {
        this.pagination.total = response.data.total
      })
    },

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

    handleToggle(data) {
      if (data.state === 1) {
        this.$confirm(
          this.$i18n.t('monitorPage.toggleStateTips', {
            monitorName: data.name,
          }),
          this.$i18n.t('tips'),
          {
            confirmButtonText: this.$i18n.t('confirm'),
            cancelButtonText: this.$i18n.t('cancel'),
            type: 'warning',
          }
        )
          .then(() => {
            toggle(data.id).then((response) => {
              this.$message.success('Stop success')
              this.getList()
            })
          })
          .catch(() => {
            this.$message.info('Cancel')
          })
      } else {
        toggle(data.id).then((response) => {
          this.$message.success('Open success')
          this.getList()
        })
      }
    },

    handleRemove(data) {
      this.$confirm(
        this.$i18n.t('monitorPage.removeMontiorTips', {
          monitorName: data.name,
        }),
        this.$i18n.t('tips'),
        {
          confirmButtonText: this.$i18n.t('confirm'),
          cancelButtonText: this.$i18n.t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          remove(data.id).then((response) => {
            this.$message.success('Success')
            this.getList()
            this.getTotal()
          })
        })
        .catch(() => {
          this.$message.info('Cancel')
        })
    },

    check() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.formProps.loading = true
          this.formProps.disabled = true
          check(this.formData)
            .then((response) => {
              this.$message.success('Connected success')
            })
            .finally((_) => {
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
      add(this.formData)
        .then((response) => {
          this.getList()
          this.getTotal()
          this.$message.success('Success')
        })
        .finally(() => {
          this.formProps.disabled = this.dialogVisible = false
        })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData)
        .then((response) => {
          this.getList()
          this.$message.success('Success')
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
