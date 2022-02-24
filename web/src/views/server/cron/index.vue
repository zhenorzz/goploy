<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-col :span="16">
        <el-select
          v-model="serverId"
          placeholder="Select server"
          style="width: 160px"
          filterable
          @change="selectServer"
        >
          <el-option
            v-for="server in serverOption"
            :key="server.id"
            :label="server.name"
            :value="server.id"
          />
        </el-select>
      </el-col>
      <el-col v-if="serverId !== ''" :span="8" style="text-align: right">
        <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
      </el-col>
    </el-row>
    <el-table
      v-loading="tableLoading"
      :max-height="tableHeight"
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column
        prop="expression"
        :label="$t('expression')"
        min-width="120"
        show-overflow-tooltip
      />
      <el-table-column
        prop="command"
        :label="$t('command')"
        min-width="140"
        show-overflow-tooltip
      />
      <el-table-column prop="singleMode" label="Single mode" width="100">
        <template #default="scope">
          <span v-if="scope.row.singleMode === 0">no</span>
          <span v-else>yes</span>
        </template>
      </el-table-column>
      <el-table-column prop="logLevel" label="Log level" width="80">
        <template #default="scope">
          <span v-if="scope.row.logLevel === 0">none</span>
          <span v-else-if="scope.row.logLevel === 1">stdout</span>
          <span v-else-if="scope.row.logLevel === 2">stdout+stderr</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="description"
        :label="$t('description')"
        min-width="240"
        show-overflow-tooltip
      />
      <el-table-column prop="creator" :label="$t('creator')" min-width="80" />
      <el-table-column prop="editor" :label="$t('editor')" min-width="80" />
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
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
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
    <el-row type="flex" justify="end" style="margin-top: 10px; width: 100%">
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
    >
      <el-form
        ref="form"
        v-loading="formProps.loading"
        :rules="formRules"
        :model="formData"
        label-width="120px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item :label="$t('expression')" prop="expression">
          <el-input
            v-model="formData.expression"
            autocomplete="off"
            placeholder="* * * * * ? with second"
            @change="onExpressionChange"
          />
          <span>{{ formProps.dateLocale }}</span>
        </el-form-item>
        <el-form-item :label="$t('command')" prop="command">
          <el-input v-model="formData.command" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Single mode">
          <el-radio-group v-model="formData.singleMode">
            <el-radio :label="0">no</el-radio>
            <el-radio :label="1">yes</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="Log level">
          <el-radio-group v-model="formData.logLevel">
            <el-radio :label="0">none</el-radio>
            <el-radio :label="1">stdout</el-radio>
            <el-radio :label="2">stdout+stderr</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item :label="$t('description')" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :autosize="{ minRows: 2 }"
          />
        </el-form-item>
      </el-form>
      <template #footer>
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
      </template>
    </el-dialog>
  </el-row>
</template>
<script lang="ts">
import tableHeight from '@/mixin/tableHeight'
import cronstrue from 'cronstrue/i18n'
import { ServerOption } from '@/api/server'
import { CronList, CronAdd, CronEdit, CronRemove, CronData } from '@/api/cron'
import Validator, { RuleItem } from 'async-validator'
import { ElMessage } from 'element-plus'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'ServerCron',
  mixins: [tableHeight],
  data() {
    return {
      serverId: '',
      crontabCommand: '',
      dialogVisible: false,
      serverOption: [] as ServerOption['datagram']['list'],
      selectedItem: {},
      tableLoading: false,
      tableData: [] as CronList['datagram']['list'],
      pagination: {
        page: 1,
        rows: 16,
        total: 0,
      },
      tempFormData: {},
      formData: {
        id: 0,
        serverId: 0,
        expression: '',
        command: '',
        singleMode: 0,
        logLevel: 0,
        description: '',
      },
      formProps: {
        loading: false,
        disabled: false,
        dateLocale: '',
      },
      formRules: {
        expression: [
          {
            required: true,
            validator: (_, value) => {
              if (value.trim().split(/\s+/).length != 6) {
                return new Error('6 parts are required.')
              }
              try {
                cronstrue.toString(value)
                return true
              } catch (error) {
                if (typeof error === 'string') {
                  return new Error(error)
                } else if (error instanceof Error) {
                  return error
                }
              }
            },
            trigger: 'blur',
          } as RuleItem,
        ],
        command: [
          { required: true, message: 'Command required', trigger: 'blur' },
        ],
      },
    }
  },

  created() {
    this.getServerOption()
    this.storeFormData()
  },

  methods: {
    selectServer() {
      this.getList()
    },

    getServerOption() {
      new ServerOption().request().then((response) => {
        this.serverOption = response.data.list
      })
    },

    getList() {
      this.tableLoading = true
      this.tableData = []
      new CronList({ serverId: Number(this.serverId) }, this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
    },

    handleAdd() {
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleEdit(data: CronData['datagram']) {
      this.formData = data
      this.dialogVisible = true
    },

    handleRemove(data: CronData['datagram']) {
      new CronRemove({ id: data.id }).request().then(() => {
        this.getList()
        ElMessage.success('Success')
      })
    },

    onExpressionChange() {
      if (this.formData.expression.trim().split(/\s+/).length != 6) {
        return
      }
      this.formProps.dateLocale = cronstrue.toString(this.formData.expression, {
        locale: this.getLocale(),
      })
    },

    handlePageChange(val = 1) {
      this.pagination.page = val
      this.getList()
    },

    submit() {
      ;(this.$refs.form as Validator).validate((valid: boolean) => {
        this.formData.expression = this.formData.expression.trim()
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
      new CronAdd(this.formData)
        .request()
        .then(() => {
          this.getList()
          ElMessage.success('Success')
        })
        .finally(() => {
          this.formProps.disabled = this.dialogVisible = false
        })
    },

    edit() {
      this.formProps.disabled = true
      new CronEdit(this.formData)
        .request()
        .then(() => {
          this.getList()
          ElMessage.success('Success')
        })
        .finally(() => {
          this.formProps.disabled = this.dialogVisible = false
        })
    },

    getLocale() {
      if (this.$i18n.locale === 'zh-cn') {
        return 'zh_CN'
      }
      return this.$i18n.locale
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
