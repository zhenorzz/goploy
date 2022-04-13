<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="account"
          style="width: 200px"
          placeholder="Filter the account"
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
        border
        stripe
        highlight-current-row
        height="100%"
        :data="tableData"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="account" label="Account" width="110" />
        <el-table-column prop="remoteAddr" label="Remote addr" width="160" />
        <el-table-column prop="referer" label="Referer" width="200" />
        <el-table-column
          prop="userAgent"
          label="User agent"
          width="160"
          show-overflow-tooltip
        />
        <el-table-column prop="reason" label="Reason" show-overflow-tooltip />
        <el-table-column prop="loginTime" label="Login time" width="155" />
      </el-table>
    </el-row>
    <el-row type="flex" justify="end" style="margin-top: 10px; width: 100%">
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
export default { name: 'LoginLog' }
</script>
<script lang="ts" setup>
import { Search } from '@element-plus/icons-vue'
import { LoginLogList, LoginLogTotal } from '@/api/log'
import { ref } from 'vue'
const account = ref('')
const tableLoading = ref(false)
const tableData = ref<LoginLogList['datagram']['list']>([])
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
  new LoginLogList({ account: account.value }, pagination.value)
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function getTotal() {
  new LoginLogTotal({ account: account.value }).request().then((response) => {
    pagination.value.total = response.data.total
  })
}

function handlePageChange(val = 1) {
  pagination.value.page = val
  getList()
}
</script>
