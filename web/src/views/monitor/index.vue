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
      <el-table-column prop="name" label="应用名称" min-width="140" />
      <el-table-column prop="domain" label="Domain" min-width="140">
        <template slot-scope="scope">
          {{ scope.row.domain }}:{{ scope.row.port }}
        </template>
      </el-table-column>
      <el-table-column prop="second" label="间隔(s)" width="80" />
      <el-table-column prop="times" label="连续失败次数" width="100" />
      <el-table-column prop="notifyType" label="通知方式" width="70">
        <template slot-scope="scope">
          <span v-if="scope.row.notifyType === 1">企业微信</span>
          <span v-else-if="scope.row.notifyType === 2">钉钉</span>
          <span v-else-if="scope.row.notifyType === 3">飞书</span>
          <span v-else-if="scope.row.notifyType === 255">自定义</span>
        </template>
      </el-table-column>
      <el-table-column prop="description" label="描述" min-width="140" show-overflow-tooltip />
      <el-table-column prop="state" label="状态" width="45">
        <template slot-scope="scope">
          {{ scope.row.state === 1 ? '开启' : '暂停' }}
        </template>
      </el-table-column>
      <el-table-column prop="insertTime" label="创建时间" width="160" />
      <el-table-column prop="updateTime" label="更新时间" width="160" />
      <el-table-column prop="operation" label="操作" width="220" fixed="right">
        <template slot-scope="scope">
          <el-button type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button :type="scope.row.state === 1 ? 'warning' : 'success'" @click="handleToggle(scope.row)">
            {{ scope.row.state === 1 ? '暂停' : '开启' }}
          </el-button>
          <el-button type="danger" @click="handleRemove(scope.row)">删除</el-button>
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
    <el-dialog title="应用设置" :visible.sync="dialogVisible">
      <el-form ref="form" v-loading="formProps.loading" :element-loading-text="formProps.loadingMessage" :rules="formRules" :model="formData" label-width="120px">
        <el-form-item label="应用名称" prop="name">
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item label="Domain/IP" prop="domain">
          <el-input v-model="formData.domain" autocomplete="off" placeholder="(不需要带http)" />
        </el-form-item>
        <el-form-item label="port" prop="port">
          <el-input v-model.number="formData.port" autocomplete="off" />
        </el-form-item>
        <el-form-item label="间隔(s)" prop="second">
          <el-input v-model.number="formData.second" autocomplete="off" />
        </el-form-item>
        <el-form-item label="连续失败次数" prop="times">
          <el-input v-model="formData.times" autocomplete="off" />
        </el-form-item>
        <el-form-item label="通知" prop="notifyTarget">
          <el-row type="flex">
            <el-select v-model="formData.notifyType">
              <el-option label="企业微信" :value="1" />
              <el-option label="钉钉" :value="2" />
              <el-option label="飞书" :value="3" />
              <el-option label="自定义" :value="255" />
            </el-select>
            <el-input v-model.trim="formData.notifyTarget" autocomplete="off" placeholder="webhook链接" />
          </el-row>
        </el-form-item>
        <el-form-item label="描述" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :autosize="{ minRows: 2}"
            placeholder="请输入描述"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-row type="flex" justify="space-between">
          <el-button type="success" @click="check">测试应用状态</el-button>
          <el-row>
            <el-button @click="dialogVisible = false">取 消</el-button>
            <el-button :disabled="formProps.disabled" type="primary" @click="submit">确 定</el-button>
          </el-row>
        </el-row>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { getList, getTotal, add, edit, check, toggle, remove } from '@/api/monitor'

export default {
  name: 'Monitor',
  data() {
    return {
      dialogVisible: false,
      tableData: [],
      pagination: {
        page: 1,
        rows: 16,
        total: 0
      },
      tempFormData: {},
      formProps: {
        loading: false,
        loadingMessage: '测试连接中，请稍后',
        disabled: false
      },
      formData: {
        id: 0,
        name: '',
        domain: '',
        port: 80,
        second: 3,
        times: 1,
        notifyType: 1,
        notifyTarget: ''
      },
      formRules: {
        name: [
          { required: true, message: '请输入应用名称', trigger: 'blur' }
        ],
        domain: [
          { required: true, message: '请输入host或者ip(不需要带http)', trigger: 'blur' }
        ],
        port: [
          { type: 'number', required: true, min: 0, max: 65535, message: '请输入正确服务器端口', trigger: 'blur' }
        ],
        second: [
          { type: 'number', required: true, min: 1, message: '请输入时间间隔', trigger: 'blur' }
        ],
        times: [
          { type: 'number', required: true, min: 1, message: '请输入连续失败次数', trigger: 'blur' }
        ],
        notifyTarget: [
          { required: true, message: '请填写webhook链接' }
        ],
        description: [
          { max: 255, message: '描述最多255个字符', trigger: 'blur' }
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
    getList() {
      getList(this.pagination).then((response) => {
        this.tableData = response.data.list
      })
    },

    getTotal() {
      getTotal().then((response) => {
        this.pagination.total = response.data.total
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

    handleToggle(data) {
      if (data.state === 1) {
        this.$confirm('此操作将暂停监控该应用, 是否继续?', '提示', {
          confirmButtonText: '确定',
          cancelButtonText: '取消',
          type: 'warning'
        }).then(() => {
          toggle(data.id).then((response) => {
            this.$message.success('暂停成功')
            this.getList()
          })
        }).catch(() => {
          this.$message.info('已取消暂停')
        })
      } else {
        toggle(data.id).then((response) => {
          this.$message.success('开启成功')
          this.getList()
        })
      }
    },

    handleRemove(data) {
      this.$confirm('此操作将不再监控该应用, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message.success('删除成功')
          this.getList()
          this.getTotal()
        })
      }).catch(() => {
        this.$message.info('已取消删除')
      })
    },

    check() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          this.formProps.loading = true
          this.formProps.disabled = true
          check(this.formData).then(response => {
            this.$message.success('连接成功')
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
        this.$message.success('添加成功')
      }).finally(() => {
        this.formProps.disabled = this.dialogVisible = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.getList()
        this.$message.success('编辑成功')
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
<style lang="scss" scoped>
@import "@/styles/mixin.scss";
.template-dialog {
  padding-right: 10px;
  height: 400px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
