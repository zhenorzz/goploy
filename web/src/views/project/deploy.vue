<template>
  <el-row class="app-container">
    <el-table
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" />
      <el-table-column prop="name" label="项目名称" />
      <el-table-column prop="state" label="状态" />
      <el-table-column prop="createTime" label="创建时间" width="160" />
      <el-table-column prop="updateTime" label="更新时间" width="160" />
      <el-table-column prop="operation" label="操作" width="260">
        <template slot-scope="scope">
          <el-button size="small" type="primary" @click="publish(scope.row.id)">构建</el-button>
          <el-button size="small" type="success">详情</el-button>
          <el-button size="small" type="danger">回滚</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-row>
</template>
<script>
import { get } from '@/api/project'
import { publish } from '@/api/deploy'
import { parseTime } from '@/utils'

const STATE = ['构建中', '构建成功', '构建失败', '撤回']
export default {
  data() {
    return {
      tableData: []
    }
  },
  created() {
    this.get()
  },
  methods: {
    get() {
      get().then((response) => {
        const projectList = response.data.projectList
        projectList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
          element.state = STATE[element.state]
        })
        this.tableData = projectList
      })
    },
    publish(id) {
      publish(id).then((response) => {
        this.$message({
          message: response.message,
          type: 'success',
          duration: 5 * 1000
        })
        this.getDeploy()
      })
    }
  }
}
</script>
