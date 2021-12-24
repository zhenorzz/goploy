<template>
  <el-row class="app-container">
    <el-row
      v-show="$store.state.user.superManager"
      class="app-bar"
      type="flex"
      justify="end"
    >
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
    </el-row>
    <el-table
      :key="tableHeight"
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      :max-height="tableHeight"
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="80" />
      <el-table-column prop="name" :label="$t('name')" />
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
        prop="user"
        :label="$t('member')"
        width="80"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button type="text" @click="handleUser(scope.row)">
            {{ $t('view') }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="80"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
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
      :title="$t('setting')"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <el-form
        ref="form"
        :rules="formRules"
        :model="formData"
        label-width="80px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button
          :disabled="formProps.disabled"
          type="primary"
          @click="submit"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
    <TheUserDialog
      v-model="dialogUserVisible"
      :namespace-id="selectedItem.id"
    />
  </el-row>
</template>
<script lang="ts">
import {
  NamespaceList,
  NamespaceTotal,
  NamespaceAdd,
  NamespaceEdit,
  NamespaceData,
} from '@/api/namespace'
import { getRole } from '@/utils/namespace'
import tableHeight from '@/mixin/tableHeight'
import { ElMessage } from 'element-plus'
import Validator from 'async-validator'
import TheUserDialog from './components/TheUserDialog.vue'
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'NamespaceIndex',
  components: { TheUserDialog },
  mixins: [tableHeight],
  data() {
    return {
      role: getRole(),
      dialogVisible: false,
      dialogUserVisible: false,
      selectedItem: {},
      tableLoading: false,
      tableData: [] as NamespaceList['datagram']['list'],
      pagination: {
        page: 1,
        rows: 16,
        total: 0,
      },
      tempFormData: {},
      formProps: {
        disabled: false,
      },
      formData: {
        id: 0,
        name: '',
      },
      formRules: {
        name: [{ required: true, message: 'Name required', trigger: 'blur' }],
      },
    }
  },
  created() {
    this.storeFormData()
    this.getList()
    this.getTotal()
  },
  methods: {
    getList() {
      this.tableLoading = true
      new NamespaceList(this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
    },

    getTotal() {
      new NamespaceTotal().request().then((response) => {
        this.pagination.total = response.data.total
      })
    },

    handlePageChange(val = 1) {
      this.pagination.page = val
      this.getList()
    },

    handleAdd() {
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleEdit(data: NamespaceData['datagram']) {
      this.formData = Object.assign({}, data)
      this.dialogVisible = true
    },

    handleUser(data: NamespaceData['datagram']) {
      this.selectedItem = data
      this.dialogUserVisible = true
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
      new NamespaceAdd(this.formData)
        .request()
        .then(() => {
          this.getList()
          this.getTotal()
          ElMessage.success('Need to login again')
        })
        .finally(() => {
          this.formProps.disabled = this.dialogVisible = false
        })
    },

    edit() {
      this.formProps.disabled = true
      new NamespaceEdit(this.formData)
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
