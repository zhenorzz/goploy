<template>
  <el-row v-if="recordViewer" class="terminal">
    <div class="close" @click="closeRecordViewer">×</div>
    <div ref="record" style="width: 100%; height: 100%"></div>
  </el-row>
  <el-row v-else class="app-container">
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
          show-overflow-tooltip
        />
        <el-table-column prop="startTime" label="Start time" width="160" />
        <el-table-column prop="endTime" label="End time" width="160" />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="100"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <Button
              :permissions="[permission.ShowTerminalRecord]"
              type="primary"
              text
              @click="handleRecord(scope.row)"
            >
              Record
            </Button>
          </template>
        </el-table-column>
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
export default { name: 'TerminalLog' }
</script>
<script lang="ts" setup>
import permission from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Search } from '@element-plus/icons-vue'
import 'asciinema-player/dist/bundle/asciinema-player.css'
import * as AsciinemaPlayer from 'asciinema-player'
import { TerminalLogData, TerminalLogList, TerminalLogTotal } from '@/api/log'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import { ref, nextTick } from 'vue'
const recordViewer = ref(false)
const record = ref()
const searchParam = ref({ username: '', serverName: '' })
const tableLoading = ref(false)
const tableData = ref<TerminalLogList['datagram']['list']>([])
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
  new TerminalLogList(searchParam.value, pagination.value)
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}
function getTotal() {
  new TerminalLogTotal(searchParam.value).request().then((response) => {
    pagination.value.total = response.data.total
  })
}
function handlePageChange(val = 1) {
  pagination.value.page = val
  getList()
}
function handleRecord(data: TerminalLogData) {
  recordViewer.value = true
  const castUrl = `${location.origin}${
    import.meta.env.VITE_APP_BASE_API
  }/log/getTerminalRecord?${NamespaceKey}=${getNamespaceId()}&recordId=${
    data.id
  }`
  nextTick(() => {
    AsciinemaPlayer.create(castUrl, record.value, {
      fit: false,
      fontSize: '14px',
    })
  })
}
function closeRecordViewer() {
  recordViewer.value = false
}
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
