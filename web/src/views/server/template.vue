<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end" align="middle">
      <el-button v-if="activeTableName==='template'" type="primary" icon="el-icon-plus" @click="handleTemplateAdd" />
      <el-upload
        v-else
        ref="upload"
        :action="action"
        :before-upload="beforeUpload"
        :on-success="handleUploadSuccess"
        :on-remove="handleRemove"
        :before-remove="beforeRemove"
        :show-file-list="false"
        multiple
      >
        <el-button type="primary" icon="el-icon-upload" />
      </el-upload>
    </el-row>
    <el-tabs v-model="activeTableName" type="border-card" style="box-shadow:none;">
      <el-tab-pane :label="$t('template')" name="template">
        <el-table
          border
          stripe
          highlight-current-row
          :data="templateTableData"
        >
          <el-table-column prop="name" :label="$t('name')" />
          <el-table-column prop="remark" :label="$t('desc')" />
          <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
          <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
          <el-table-column prop="operation" :label="$t('op')" width="130" align="center">
            <template slot-scope="scope">
              <el-button type="primary" icon="el-icon-edit" @click="handleTemplateEdit(scope.row)" />
              <el-button type="danger" icon="el-icon-delete" @click="handleTemplateDelete(scope.row)" />
            </template>
          </el-table-column>
        </el-table>
        <el-row type="flex" justify="end" style="margin-top: 10px;">
          <el-pagination
            hide-on-single-page
            :total="tplPagination.total"
            :page-size="tplPagination.rows"
            background
            layout="prev, pager, next"
            @current-change="handleTplPageChange"
          />
        </el-row>
      </el-tab-pane>
      <el-tab-pane :label="$t('package')" name="package">
        <el-table
          border
          stripe
          highlight-current-row
          :data="packageTableData"
        >
          <el-table-column prop="name" :label="$t('name')" />
          <el-table-column prop="humanSize" :label="$t('size')" />
          <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
          <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
          <el-table-column prop="operation" :label="$t('op')" width="110" align="center">
            <template slot-scope="scope">
              <el-upload
                ref="upload"
                :action="action+'?packageId='+scope.row.id"
                :before-upload="beforeUpload"
                :on-success="handleUploadSuccess"
                :on-remove="handleRemove"
                :before-remove="beforeRemove"
                :show-file-list="false"
                multiple
              >
                <el-button type="primary">{{ $t('reUpload') }}</el-button>
              </el-upload>
            </template>
          </el-table-column>
        </el-table>
        <el-row type="flex" justify="end" style="margin-top: 10px;">
          <el-pagination
            hide-on-single-page
            :total="pkgPagination.total"
            :page-size="pkgPagination.rows"
            background
            layout="prev, pager, next"
            @current-change="handlePkgPageChange"
          />
        </el-row>
      </el-tab-pane>
    </el-tabs>
    <el-dialog :title="$t('setting')" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="100px">
        <el-form-item :label="$t('name')" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('desc')" prop="remark">
          <el-input v-model="formData.remark" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('package')" prop="package">
          <el-select
            v-model="formData.packageIds"
            multiple
            clearable
            filterable
            style="width:100%"
          >
            <el-option
              v-for="(item, index) in packageOption"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('script')" prop="script">
          {{ $t('templatePage.scriptNotice') }}
          <codemirror v-model="formData.script" :options="cmOptions" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submitTemplate">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import {
  getList as getTemplateList,
  getTotal as getTemplateTotal,
  add as addTemplate,
  edit as editTemplate,
  remove,
  removePackage
} from '@/api/template'
import {
  getList as getPackageList,
  getTotal as getPackageTotal
} from '@/api/package'

import { humanSize } from '@/utils'
// require component
import { codemirror } from 'vue-codemirror'
import 'codemirror/mode/shell/shell.js'
import 'codemirror/theme/darcula.css'
// require styles
import 'codemirror/lib/codemirror.css'
import 'codemirror/addon/scroll/simplescrollbars.js'
import 'codemirror/addon/scroll/simplescrollbars.css'
export default {
  name: 'Template',
  components: {
    codemirror
  },
  data() {
    return {
      dialogVisible: false,
      activeTableName: 'template',
      templateTableData: [],
      tplPagination: {
        page: 1,
        rows: 11,
        total: 0
      },
      packageTableData: [],
      pkgPagination: {
        page: 1,
        rows: 11,
        total: 0
      },
      packageOption: [],
      action: process.env.VUE_APP_BASE_API + '/package/upload',
      tempFormData: {},
      cmOptions: {
        tabSize: 4,
        mode: 'text/x-sh',
        lineNumbers: true,
        line: true,
        scrollbarStyle: 'overlay',
        theme: 'darcula'
      },
      formProps: {
        disabled: false
      },
      formData: {
        id: 0,
        name: '',
        remark: '',
        packageIds: [],
        packageIdStr: '',
        script: ''
      },
      formRules: {
        name: [
          { required: true, message: 'Name required', trigger: 'blur' }
        ],
        script: [
          { required: true, message: 'Script required', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.storeFormData()
    this.getTemplateList()
    this.getTemplateTotal()
    this.getPackageList()
    this.getPackageTotal()
  },
  methods: {
    getTemplateList() {
      getTemplateList(this.tplPagination).then((response) => {
        this.templateTableData = response.data.list
      })
    },

    getTemplateTotal() {
      getTemplateTotal().then((response) => {
        this.tplPagination.total = response.data.total
      })
    },

    handleTplPageChange(val) {
      this.tplPagination.page = val
      this.getTemplateList()
    },

    getPackageList() {
      getPackageList(this.pkgPagination).then((response) => {
        const packageList = response.data.list
        packageList.forEach((element) => {
          element.humanSize = humanSize(element.size)
        })
        this.packageOption = this.packageTableData = packageList
      })
    },

    getPackageTotal() {
      getPackageTotal().then((response) => {
        this.pkgPagination.total = response.data.total
      })
    },

    handlePkgPageChange(val) {
      this.pkgPagination.page = val
      this.getPackageList()
    },

    beforeUpload(file) {
      this.$message.info(this.$i18n.t('uploading'))
    },

    handleUploadSuccess(response, file, fileList) {
      if (response.code !== 0) {
        this.$message.error(response.message)
      } else {
        this.$message.success('Success')
        this.getPackageList()
        this.getPackageTotal()
      }
    },

    handleRemove(file, fileList) {
      console.log(file, fileList)
    },

    beforeRemove(file, fileList) {
      return removePackage(this.formData.id, file.name).then((response) => {
        return Promise.resolve(response)
      }).catch(err => {
        return Promise.reject(err)
      })
    },

    handleTemplateAdd() {
      this.formProps.fileList = []
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleTemplateEdit(data) {
      this.formData = Object.assign(this.formData, data)
      this.formData.packageIds = data.packageIdStr.split(',').map(element => parseInt(element))
      this.dialogVisible = true
    },

    handleTemplateDelete(data) {
      this.$confirm(this.$i18n.t('templatePage.templateDeleteTips', { templateName: data.name }), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message.success('Success')
          this.getTemplateList()
          this.getTemplateTotal()
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    submitTemplate() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.formData.packageIdStr = this.formData.packageIds.join(',')
          if (this.formData.id === 0) {
            this.addTemplate()
          } else {
            this.editTemplate()
          }
        } else {
          return false
        }
      })
    },

    addTemplate() {
      this.formProps.disabled = true
      addTemplate(this.formData).then((response) => {
        this.getTemplateList()
        this.getTemplateTotal()
        this.$message.success('Success')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    editTemplate() {
      this.formProps.disabled = true
      editTemplate(this.formData).then((response) => {
        this.getTemplateList()
        this.$message('Success')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
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
<style>
.CodeMirror {
  border-radius: 4px;
  border: 1px solid #DCDFE6;
  height: 200px;
}
</style>
