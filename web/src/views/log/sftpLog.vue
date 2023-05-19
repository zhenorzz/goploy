<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="searchParam.username"
          style="width: 200px"
          placeholder="Filter the username"
        />
        <el-input
          v-model="searchParam.serverName"
          style="width: 200px"
          placeholder="Filter the server name"
        />
        <el-button
          :loading="tableLoading"
          type="primary"
          :icon="Search"
          @click="searchList"
        />
      </el-row>
    </el-row>
    <el-row class="app-table">
      <el-table
        v-loading="tableLoading"
        highlight-current-row
        height="100%"
        :data="tableData"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="Username" width="100" />
        <el-table-column prop="serverName" label="Server Name" width="160" />
        <el-table-column prop="remoteAddr" label="Remote addr" width="160" />
        <el-table-column
          prop="userAgent"
          label="User agent"
          width="160"
          show-overflow-tooltip
        />
        <el-table-column prop="type" label="Type" width="110" />
        <el-table-column prop="path" label="Path" show-overflow-tooltip />
        <el-table-column prop="reason" label="Reason" show-overflow-tooltip />
        <el-table-column prop="insertTime" label="Time" width="160" />
      </el-table>
    </el-row>
    <el-row type="flex" justify="end" class="app-page">
      <el-pagination
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
  </el-row>
</template>
<script lang="ts">
export default { name: 'SftpLog' }
</script>
<script lang="ts" setup>
import { Search } from '@element-plus/icons-vue'
import { SftpLogList, SftpLogTotal } from '@/api/log'
import { ref } from 'vue'
const searchParam = ref({ username: '', serverName: '' })
const tableLoading = ref(false)
const tableData = ref<SftpLogList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20, total: 0 })

getList()
getTotal()

function searchList() {
  pagination.value.page = 1
  getList()
  getTotal()
}
function getList() {
  tableLoading.value = true
  tableData.value = []
  new SftpLogList(searchParam.value, pagination.value)
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}
function getTotal() {
  new SftpLogTotal(searchParam.value).request().then((response) => {
    pagination.value.total = response.data.total
  })
}
function handlePageChange(val = 1) {
  pagination.value.page = val
  getList()
}
</script>
