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
          icon="el-icon-search"
          @click="searchList"
        />
      </el-row>
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
      <el-table-column prop="username" label="Username" width="80" />
      <el-table-column prop="serverName" label="Server Name" width="160" />
      <el-table-column prop="remoteAddr" label="Remote addr" width="160" />
      <el-table-column
        prop="userAgent"
        label="User agent"
        width="160"
        show-overflow-tooltip
      />
      <el-table-column prop="type" label="Type" width="100" />
      <el-table-column prop="path" label="Path" show-overflow-tooltip />
      <el-table-column prop="reason" label="Reason" show-overflow-tooltip />
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px; width: 100%">
      <el-pagination
        style=""
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
  </el-row>
</template>
<script lang="ts">
import { SftpLogList, SftpLogTotal } from '@/api/log'
import tableHeight from '@/mixin/tableHeight'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'SftpLog',
  mixins: [tableHeight],
  data() {
    return {
      tableLoading: false,
      searchParam: {
        username: '',
        serverName: '',
      },
      tableData: [] as SftpLogList['datagram']['list'],
      pagination: {
        page: 1,
        rows: 19,
        total: 0,
      },
    }
  },

  created() {
    this.getList()
    this.getTotal()
  },

  methods: {
    searchList() {
      this.pagination.page = 1
      this.getList()
      this.getTotal()
    },
    getList() {
      this.tableLoading = true
      this.tableData = []
      new SftpLogList(this.searchParam, this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
    },
    getTotal() {
      new SftpLogTotal(this.searchParam).request().then((response) => {
        this.pagination.total = response.data.total
      })
    },
    handlePageChange(val = 1) {
      this.pagination.page = val
      this.getList()
    },
  },
})
</script>
