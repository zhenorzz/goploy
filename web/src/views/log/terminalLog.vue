<template>
  <el-row class="app-container">
    <el-row v-if="recordViewer" class="terminal">
      <div class="close" @click="closeRecordViewer">×</div>
      <div ref="record" style="width: 100%; height: 100%"></div>
    </el-row>
    <el-row v-show="!recordViewer">
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
          show-overflow-tooltip
        />
        <el-table-column prop="startTime" label="Start time" width="135" />
        <el-table-column prop="endTime" label="End time" width="135" />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="100"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <el-button type="text" @click="handleRecord(scope.row)">
              Record
            </el-button>
          </template>
        </el-table-column>
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
  </el-row>
</template>
<script lang="ts">
import 'asciinema-player/dist/bundle/asciinema-player.css'
import * as AsciinemaPlayer from 'asciinema-player'
import { TerminalLogData, TerminalLogList, TerminalLogTotal } from '@/api/log'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import tableHeight from '@/mixin/tableHeight'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'TerminalLog',
  mixins: [tableHeight],
  data() {
    return {
      recordViewer: false,
      tableLoading: false,
      searchParam: {
        username: '',
        serverName: '',
      },
      tableData: [] as TerminalLogList['datagram']['list'],
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
      new TerminalLogList(this.searchParam, this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
    },
    getTotal() {
      new TerminalLogTotal(this.searchParam).request().then((response) => {
        this.pagination.total = response.data.total
      })
    },
    handlePageChange(val = 1) {
      this.pagination.page = val
      this.getList()
    },
    handleRecord(data: TerminalLogData['datagram']) {
      this.recordViewer = true
      const castUrl = `${location.origin}${
        import.meta.env.VITE_APP_BASE_API
      }/log/getTerminalRecord?${NamespaceKey}=${getNamespaceId()}&recordId=${
        data.id
      }`
      this.$nextTick(() => {
        AsciinemaPlayer.create(castUrl, this.$refs['record'], {
          fit: false,
          fontSize: '14px',
        })
      })
    },
    closeRecordViewer() {
      this.recordViewer = false
    },
  },
})
</script>
<style lang="scss" scoped>
.terminal {
  height: calc(100vh - 124px);
  width: 100%;
  padding: 10px;
  .close {
    position: absolute;
    top: 20px;
    right: 30px;
    color: #fff;
    z-index: 1000;
    font-size: 18px;
    cursor: pointer;
  }
}
</style>

<style>
.asciinema-player {
  width: 100%;
}
</style>
