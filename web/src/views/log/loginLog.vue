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
      <el-table-column prop="account" label="Account" width="80" />
      <el-table-column prop="remoteAddr" label="Remote addr" width="160" />
      <el-table-column prop="referer" label="Referer" width="200" />
      <el-table-column
        prop="userAgent"
        label="User agent"
        width="160"
        show-overflow-tooltip
      />
      <el-table-column prop="reason" label="Reason" show-overflow-tooltip />
      <el-table-column prop="loginTime" label="Login time" width="135" />
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
import { LoginLogList, LoginLogTotal } from '@/api/log'
import tableHeight from '@/mixin/tableHeight'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'LoginLog',
  mixins: [tableHeight],
  data() {
    return {
      tableLoading: false,
      account: '',
      tableData: [] as LoginLogList['datagram']['list'],
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
      new LoginLogList({ account: this.account }, this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
    },
    getTotal() {
      new LoginLogTotal({ account: this.account })
        .request()
        .then((response) => {
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
