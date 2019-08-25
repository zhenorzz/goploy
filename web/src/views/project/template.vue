<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd">添加</el-button>
    </el-row>
    <el-table
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="name" label="名称" />
      <el-table-column prop="remark" label="描述" />
      <el-table-column prop="createTime" label="创建时间" width="160" />
      <el-table-column prop="updateTime" label="更新时间" width="160" />
      <el-table-column prop="operation" label="操作" width="150">
        <template slot-scope="scope">
          <el-button size="small" type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button size="small" type="danger" @click="handleDelete(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="模板设置" :visible.sync="dialogVisible">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="80px">
        <el-form-item label="名称" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="描述" prop="remark">
          <el-input v-model="formData.remark" autocomplete="off" />
        </el-form-item>
        <el-form-item label="安装包" prop="package">
          <el-upload
            ref="upload"
            :action="action"
            :on-remove="handleRemove"
            :before-remove="beforeRemove"
            multiple
            :file-list="formProps.fileList"
          >
            <el-button size="small" type="primary">点击上传</el-button>
          </el-upload>
        </el-form-item>
        <el-form-item label="脚本" prop="script">
          <codemirror v-model="formData.script" :options="cmOptions" />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { getList, add, edit, remove, removePackage } from '@/api/template'
import { parseTime } from '@/utils'
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
      tableData: [],
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
        disabled: false,
        fileList: []
      },
      formData: {
        id: 0,
        name: '',
        remark: '',
        package: '',
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
  computed: {
    action: function() {
      let action = process.env.VUE_APP_BASE_API + '/template/upload'
      if (this.formData.id !== 0) {
        action += '?templateId=' + this.formData.id
      }
      return action
    }
  },
  created() {
    this.storeFormData()
    this.getList()
  },
  methods: {
    getList() {
      getList().then((response) => {
        const templateList = response.data.templateList || []
        templateList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
        this.tableData = templateList
      })
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

    handleAdd() {
      this.formProps.fileList = []
      this.restoreFormData()
      this.dialogVisible = true
    },

    handleEdit(data) {
      this.formData = Object.assign({}, data)
      this.formProps.fileList = data.package.split(',').map(fileName => {
        return {
          name: fileName
        }
      })
      this.dialogVisible = true
    },

    handleDelete(data) {
      this.$confirm('此操作将删除该模板, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message({
            message: response.message,
            type: 'success',
            duration: 5 * 1000
          })
          this.getList()
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },

    submit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.formData.package = this.$refs.upload.uploadFiles.map(element => element.name).join(',')
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
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.getList()
        this.$message({
          message: response.message,
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
