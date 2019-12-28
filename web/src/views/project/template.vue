<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end" align="middle">
      <el-button v-if="activeTableName==='template'" type="primary" icon="el-icon-plus" @click="handleTemplateAdd">添加</el-button>
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
        <el-button size="small" type="primary">点击上传</el-button>
      </el-upload>
    </el-row>
    <el-tabs v-model="activeTableName" type="border-card" style="box-shadow:none;">
      <el-tab-pane label="模板" name="template">
        <el-table
          border
          stripe
          highlight-current-row
          :data="templateTableData"
        >
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="remark" label="描述" />
          <el-table-column prop="createTime" label="创建时间" width="160" />
          <el-table-column prop="updateTime" label="更新时间" width="160" />
          <el-table-column prop="operation" label="操作" width="150">
            <template slot-scope="scope">
              <el-button size="small" type="primary" @click="handleTemplateEdit(scope.row)">编辑</el-button>
              <el-button size="small" type="danger" @click="handleTemplateDelete(scope.row)">删除</el-button>
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
      <el-tab-pane label="安装包" name="package">
        <el-table
          border
          stripe
          highlight-current-row
          :data="packageTableData"
        >
          <el-table-column prop="name" label="名称" />
          <el-table-column prop="humanSize" label="大小" />
          <el-table-column prop="createTime" label="创建时间" width="160" />
          <el-table-column prop="updateTime" label="更新时间" width="160" />
          <el-table-column prop="operation" label="操作" width="90">
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
                <el-button size="small" type="primary">重传</el-button>
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
    <el-dialog title="模板设置" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="描述" prop="remark">
          <el-input v-model="formData.remark" autocomplete="off" />
        </el-form-item>
        <el-form-item label="安装包" prop="package">
          <el-select
            v-model="formData.packageIds"
            placeholder="选择安装包"
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
        <el-form-item label="脚本" prop="script">
          注意：安装包上传至目标服务器的/tmp目录
          <codemirror v-model="formData.script" :options="cmOptions" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submitTemplate">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import {
  getList as getTemplateList,
  add as addTemplate,
  edit as editTemplate,
  remove,
  removePackage
} from '@/api/template'
import {
  getList as getPackageList
} from '@/api/package'

import { parseTime, humanSize } from '@/utils'
// require component
import { codemirror } from 'vue-codemirror'
import 'codemirror/mode/shell/shell.js'
import 'codemirror/theme/darcula.css'
// require styles
import 'codemirror/lib/codemirror.css'
import 'codemirror/addon/scroll/simplescrollbars.js'
import 'codemirror/addon/scroll/simplescrollbars.css'
export default {
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
          { required: true, message: '名称', trigger: 'blur' }
        ],
        script: [
          { required: true, message: '请输入脚本', trigger: 'blur' }
        ]
      }
    }
  },
  created() {
    this.storeFormData()
    this.getTemplateList()
    this.getPackageList()
  },
  methods: {
    getTemplateList() {
      getTemplateList(this.tplPagination).then((response) => {
        const templateList = response.data.templateList || []
        templateList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
        this.templateTableData = templateList
        this.tplPagination = response.data.pagination
      })
    },

    handleTplPageChange(val) {
      this.tplPagination.page = val
      this.getTemplateList()
    },

    getPackageList() {
      getPackageList(this.pkgPagination).then((response) => {
        const packageList = response.data.packageList || []
        packageList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
          element.humanSize = humanSize(element.size)
        })
        this.packageOption = this.packageTableData = packageList
        this.pkgPagination = response.data.pagination
      })
    },

    handlePkgPageChange(val) {
      this.pkgPagination.page = val
      this.getPackageList()
    },

    beforeUpload(file) {
      this.$message.info('正在上传')
    },

    handleUploadSuccess(response, file, fileList) {
      if (response.code !== 0) {
        this.$message.error(response.message)
      } else {
        this.$message.success('上传成功')
        this.getPackageList()
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
      this.$confirm('此操作将删除该模板, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message({
            message: '删除成功',
            type: 'success',
            duration: 5 * 1000
          })
          this.getTemplateList()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
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
        this.$message({
          message: '添加成功',
          type: 'success',
          duration: 5 * 1000
        })
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    editTemplate() {
      this.formProps.disabled = true
      editTemplate(this.formData).then((response) => {
        this.getTemplateList()
        this.$message({
          message: '修改成功',
          type: 'success',
          duration: 5 * 1000
        })
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
