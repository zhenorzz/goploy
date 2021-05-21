<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="crontabCommand"
          style="width: 200px"
          placeholder="Fitler the command"
        />
        <el-button type="primary" icon="el-icon-search" @click="searchList" />
      </el-row>
      <el-row>
        <el-tooltip effect="dark" :content="$t('import')" placement="top">
          <el-button
            type="primary"
            icon="el-icon-download"
            @click="handleImport"
          />
        </el-tooltip>
        <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
      </el-row>
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
        prop="command"
        :label="$t('command')"
        min-width="140"
        show-overflow-tooltip
      />
      <el-table-column
        prop="description"
        :label="$t('desc')"
        min-width="240"
        show-overflow-tooltip
      />
      <el-table-column prop="creator" :label="$t('creator')" min-width="50" />
      <el-table-column prop="editor" :label="$t('editor')" min-width="50" />
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
        prop="server"
        width="80"
        :label="$t('server')"
        align="center"
      >
        <template #default="scope">
          <el-button type="text" @click="handleServer(scope.row)">{{
            $t('view')
          }}</el-button>
        </template>
      </el-table-column>
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
        label-width="80px"
      >
        <el-form-item :label="$t('time')" prop="date">
          <el-input
            v-model="formData.date"
            autocomplete="off"
            placeholder="* * * * ?"
            @change="onDateChange"
          />
          <span>{{ formProps.dateLocale }}</span>
        </el-form-item>
        <el-form-item :label="$t('script')" prop="script">
          <el-input v-model.trim="formData.script" autocomplete="off" />
        </el-form-item>
        <el-form-item
          v-show="formData.id === 0"
          :label="$t('server')"
          prop="serverIds"
        >
          <el-select
            v-model="formData.serverIds"
            multiple
            style="width: 100%"
            filterable
          >
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.label"
              :value="item.id"
            />
          </el-select>
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
    <el-dialog
      v-model="crontabRemoveVisible"
      :title="$t('remove')"
      width="500px"
    >
      <el-form
        ref="crontabRemoveForm"
        :model="crontabRemoveFormData"
        label-position="top"
      >
        <el-form-item :label="$t('command')">
          <span>{{ crontabRemoveFormProps.command }}</span>
        </el-form-item>
        <el-form-item :label="$t('crontabPage.removeServerCrontabLabel')">
          <el-radio-group v-model="crontabRemoveFormData.radio">
            <el-radio :label="0">{{ $t('boolOption[0]') }}</el-radio>
            <el-radio :label="1">{{ $t('boolOption[1]') }}</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="crontabRemoveVisible = false">{{
          $t('cancel')
        }}</el-button>
        <el-button
          :disabled="crontabRemoveFormProps.disabled"
          type="primary"
          @click="remove"
          >{{ $t('confirm') }}</el-button
        >
      </template>
    </el-dialog>
    <TheImportDialog
      v-model="importVisible"
      :on-success="handleImportSuccess"
    />
    <TheServerDialog v-model="serverVisible" :crontab-id="selectedItem.id" />
  </el-row>
</template>
<script lang="ts">
import tableHeight from '@/mixin/tableHeight'
import cronstrue from 'cronstrue/i18n'
import {
  CrontabList,
  CrontabTotal,
  CrontabAdd,
  CrontabEdit,
  CrontabRemove,
  CrontabData,
} from '@/api/crontab'
import Validator, { RuleItem } from 'async-validator'
import { ElMessage } from 'element-plus'
import { ServerOption } from '@/api/server'
import TheImportDialog from './TheImportDialog.vue'
import TheServerDialog from './TheServerDialog.vue'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Crontab',
  components: { TheImportDialog, TheServerDialog },
  mixins: [tableHeight],
  data() {
    return {
      crontabCommand: '',
      dialogVisible: false,
      crontabRemoveVisible: false,
      serverVisible: false,
      addServerVisible: false,
      importVisible: false,
      selectedItem: {},
      tableLoading: false,
      tableData: [] as CrontabList['datagram']['list'],
      serverOption: [] as ServerOption['datagram']['list'],
      pagination: {
        page: 1,
        rows: 16,
        total: 0,
      },
      formProps: {
        loading: false,
        disabled: false,
        dateLocale: '',
      },
      formData: {
        id: 0,
        command: '',
        date: '',
        script: '',
        serverIds: [],
      },
      formRules: {
        date: [
          {
            required: true,
            validator: (_, value) => {
              try {
                cronstrue.toString(value)
                return true
              } catch (error) {
                return error
              }
            },
            trigger: 'blur',
          } as RuleItem,
        ],
        script: [
          { required: true, message: 'Script required', trigger: 'blur' },
        ],
      },
      crontabRemoveFormData: {
        id: 0,
        radio: 0,
      },
      crontabRemoveFormProps: {
        command: '',
        disabled: false,
      },
    }
  },
  created() {
    this.getList()
    this.getTotal()
  },

  methods: {
    getList() {
      this.tableLoading = true
      new CrontabList({ command: this.crontabCommand }, this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list.map((element) => {
            const commandSplit = element.command.split(' ')
            element.date = commandSplit.slice(0, 5).join(' ')
            element.dateLocale = cronstrue.toString(element.date, {
              locale: this.getLocale(),
            })
            element.script = commandSplit.slice(5).join(' ')
            element.description = `${element.dateLocale}, ${this.$t('run')}: ${
              element.script
            }`
            return element
          })
        })
        .finally(() => {
          this.tableLoading = false
        })
    },

    getTotal() {
      new CrontabTotal({ command: this.crontabCommand })
        .request()
        .then((response) => {
          this.pagination.total = response.data.total
        })
    },

    getServerOption() {
      new ServerOption().request().then((response) => {
        this.serverOption = response.data.list
      })
    },

    searchList() {
      this.pagination.page = 1
      this.getList()
      this.getTotal()
    },

    handleAdd() {
      this.getServerOption()
      this.formData.id = 0
      this.dialogVisible = true
    },

    handleImport() {
      this.importVisible = true
    },

    handleImportSuccess() {
      this.getList()
      this.getTotal()
    },

    handleEdit(data: CrontabData['datagram']['detail']) {
      this.formData.id = data.id
      this.formData.date = data.date || ''
      this.formData.script = data.script || ''
      this.formData.serverIds = []
      this.formProps.dateLocale = data.dateLocale || ''
      this.dialogVisible = true
    },

    handleServer(data: CrontabData['datagram']['detail']) {
      this.selectedItem = data
      this.serverVisible = true
    },

    handleAddServer() {
      this.addServerVisible = true
    },

    handleRemove(data: CrontabData['datagram']['detail']) {
      this.crontabRemoveFormData.id = data.id
      this.crontabRemoveFormProps.command = data.command
      this.crontabRemoveVisible = true
    },

    onDateChange() {
      this.formProps.dateLocale = cronstrue.toString(this.formData.date, {
        locale: this.getLocale(),
      })
    },

    handlePageChange(val = 1) {
      this.pagination.page = val
      this.getList()
    },

    submit() {
      ;(this.$refs.form as Validator).validate((valid: boolean) => {
        if (valid) {
          this.formData.command =
            this.formData.date + ' ' + this.formData.script
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
      new CrontabAdd(this.formData)
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
      new CrontabEdit(this.formData)
        .request()
        .then(() => {
          this.getList()
          ElMessage.success('Success')
        })
        .finally(() => {
          this.formProps.disabled = this.dialogVisible = false
        })
    },

    remove() {
      this.crontabRemoveFormProps.disabled = true
      new CrontabRemove(this.crontabRemoveFormData)
        .request()
        .then(() => {
          this.getList()
          this.getTotal()
          ElMessage.success('Success')
        })
        .finally(() => {
          this.crontabRemoveFormProps.disabled = this.crontabRemoveVisible =
            false
        })
    },

    getLocale() {
      if (this.$i18n.locale === 'zh-cn') {
        return 'zh_CN'
      }
      return this.$i18n.locale
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
