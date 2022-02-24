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
          v-model="searchParam.projectName"
          style="width: 200px"
          placeholder="Filter the project name"
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
      <el-table-column prop="token" label="Token" width="300" />
      <el-table-column prop="publisherName" label="Username" width="80" />
      <el-table-column prop="projectName" label="Project Name" width="160" />
      <el-table-column prop="state" label="State" align="center" width="80">
        <template #default="scope">
          <span v-if="scope.row.state === 1" style="color: #67c23a">
            {{ $t('success') }}
          </span>
          <span v-else style="color: #f56c6c">{{ $t('fail') }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="detail" label="Reason" show-overflow-tooltip />
      <el-table-column prop="insertTime" label="insertTime" width="135" />
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="100"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button type="text" @click="handleDetail(scope.row)">
            {{ $t('detail') }}
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
    <el-dialog
      v-model="dialogVisible"
      :title="$t('detail')"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <el-row v-loading="traceLoading" class="project-detail">
        <div
          v-for="(item, index) in publishLocalTraceList"
          :key="index"
          style="width: 100%"
        >
          <template v-if="item.type === 2">
            <el-row style="margin: 5px 0">
              <i v-if="item.state === 1" class="el-icon-check icon-success" />
              <i v-else class="el-icon-close icon-fail" />
              -------------GIT-------------
            </el-row>
            <el-row>Time: {{ item.updateTime }}</el-row>
            <template v-if="item.state !== 0">
              <el-row>Branch: {{ item['branch'] }}</el-row>
              <el-row>Commit:{{ item['commit'] }}</el-row>
              <el-row>Message: {{ item['message'] }}</el-row>
              <el-row>Author: {{ item['author'] }}</el-row>
              <el-row>
                Datetime:
                {{ item['timestamp'] ? parseTime(item['timestamp']) : '' }}
              </el-row>
              <el-row>
                <span style="white-space: pre-line">{{ item['diff'] }}</span>
              </el-row>
            </template>
            <el-row v-else style="margin: 5px 0">
              <span style="white-space: pre-line; padding: 5px 0">
                {{ item.detail }}
              </span>
            </el-row>
          </template>
          <div v-if="item.type === 3">
            <hr />
            <el-row align="middle">
              <i v-if="item.state === 1" class="el-icon-check icon-success" />
              <i v-else class="el-icon-close icon-fail" />
              --------After pull--------
            </el-row>
            <el-row>Time: {{ item.updateTime }}</el-row>
            <el-row style="width: 100%">
              <div>Script:</div>
              <pre style="white-space: pre-line">{{ item.script }}</pre>
            </el-row>
            <div v-loading="traceDetail[item.id] === ''" style="margin: 5px 0">
              <span style="padding: 5px 0">[goploy ~]#</span>
              <el-button
                v-if="item.state === 1 && !(item.id in traceDetail)"
                type="text"
                @click="getPublishTraceDetail(item)"
              >
                {{ $t('deployPage.showDetail') }}
              </el-button>
              <div v-else style="white-space: pre-line; padding: 5px 0">
                {{ traceDetail[item.id] }}
              </div>
            </div>
          </div>
        </div>
        <el-tabs v-model="activeRomoteTracePane">
          <el-tab-pane
            v-for="(item, serverName) in publishRemoteTraceList"
            :key="serverName"
            :label="serverName"
            :name="serverName"
          >
            <div v-for="(trace, key) in item" :key="key">
              <template v-if="trace.type === 4">
                <el-row style="margin: 5px 0">
                  <i
                    v-if="trace.state === 1"
                    class="el-icon-check icon-success"
                  />
                  <i v-else class="el-icon-close icon-fail" />
                  ---------Before deploy---------
                </el-row>
                <el-row style="margin: 5px 0">
                  Time: {{ trace.updateTime }}
                </el-row>
                <el-row>
                  Script:
                  <pre style="white-space: pre-line">{{ trace.script }}</pre>
                </el-row>
                <div v-loading="traceDetail[trace.id] === ''">
                  <span style="padding: 5px 0">[goploy ~]#</span>
                  <el-button
                    v-if="trace.state === 1 && !(trace.id in traceDetail)"
                    type="text"
                    @click="getPublishTraceDetail(trace)"
                  >
                    {{ $t('deployPage.showDetail') }}
                  </el-button>
                  <div v-else style="white-space: pre-line; padding: 5px 0">
                    {{ traceDetail[trace.id] }}
                  </div>
                </div>
              </template>
              <template v-else-if="trace.type === 5">
                <el-row style="margin: 5px 0">
                  <i
                    v-if="trace.state === 1"
                    class="el-icon-check icon-success"
                  />
                  <i v-else class="el-icon-close icon-fail" />
                  -----------Rsync------------
                </el-row>
                <el-row style="margin: 5px 0">
                  Time: {{ trace.updateTime }}
                </el-row>
                <el-row>Command: {{ trace.command }}</el-row>
                <div v-loading="traceDetail[trace.id] === ''">
                  <span style="padding: 5px 0">[goploy ~]#</span>
                  <el-button
                    v-if="trace.state === 1 && !(trace.id in traceDetail)"
                    type="text"
                    @click="getPublishTraceDetail(trace)"
                  >
                    {{ $t('deployPage.showDetail') }}
                  </el-button>
                  <div v-else style="white-space: pre-line; padding: 5px 0">
                    {{ traceDetail[trace.id] }}
                  </div>
                </div>
              </template>
              <template v-else-if="trace.type === 6">
                <el-row style="margin: 5px 0">
                  <i
                    v-if="trace.state === 1"
                    class="el-icon-check icon-success"
                  />
                  <i v-else class="el-icon-close icon-fail" />
                  --------After deploy--------
                </el-row>
                <el-row style="margin: 5px 0">
                  Time: {{ trace.updateTime }}
                </el-row>
                <el-row>Script: {{ trace.script }}</el-row>
                <div
                  v-loading="traceDetail[trace.id] === ''"
                  style="margin: 5px 0"
                >
                  <span>[goploy ~]#</span>
                  <el-button
                    v-if="trace.state === 1 && !(trace.id in traceDetail)"
                    type="text"
                    @click="getPublishTraceDetail(trace)"
                  >
                    {{ $t('deployPage.showDetail') }}
                  </el-button>
                  <div v-else style="white-space: pre-line; padding: 5px 0">
                    {{ traceDetail[trace.id] }}
                  </div>
                </div>
              </template>
            </div>
          </el-tab-pane>
        </el-tabs>
      </el-row>
    </el-dialog>
  </el-row>
</template>
<script lang="ts">
import { PublishLogData, PublishLogList, PublishLogTotal } from '@/api/log'
import { DeployTrace, DeployTraceDetail, PublishTraceData } from '@/api/deploy'
import { parseTime } from '@/utils'
import tableHeight from '@/mixin/tableHeight'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'PublishLog',
  mixins: [tableHeight],
  data() {
    return {
      dialogVisible: false,
      tableLoading: false,
      searchParam: {
        username: '',
        projectName: '',
      },
      tableData: [] as PublishLogList['datagram']['list'],
      pagination: {
        page: 1,
        rows: 17,
        total: 0,
      },
      traceLoading: false,
      traceDetail: {} as Record<number, string>,
      activeRomoteTracePane: '',
      publishLocalTraceList: [] as DeployTrace['datagram']['list'],
      publishRemoteTraceList: {} as Record<
        string,
        DeployTrace['datagram']['list']
      >,
    }
  },

  created() {
    this.getList()
    this.getTotal()
  },

  methods: {
    parseTime,
    searchList() {
      this.pagination.page = 1
      this.getList()
      this.getTotal()
    },
    getList() {
      this.tableLoading = true
      this.tableData = []
      new PublishLogList(this.searchParam, this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableLoading = false
        })
    },
    getTotal() {
      new PublishLogTotal(this.searchParam).request().then((response) => {
        this.pagination.total = response.data.total
      })
    },
    handlePageChange(val = 1) {
      this.pagination.page = val
      this.getList()
    },
    handleDetail(data: PublishLogData['datagram']) {
      this.dialogVisible = true
      this.traceLoading = true
      new DeployTrace({ lastPublishToken: data.token })
        .request()
        .then((response) => {
          const publishTraceList = response.data.list.map((element) => {
            if (element.ext !== '') {
              Object.assign(element, JSON.parse(element.ext))
            }
            return element
          })

          this.publishLocalTraceList = publishTraceList.filter(
            (element) => element.type < 4
          )
          this.publishRemoteTraceList = {}
          for (const trace of publishTraceList) {
            if (trace.detail !== '') {
              this.traceDetail[trace.id] = trace.detail
            }
            if (trace.type < 4) continue
            if (!this.publishRemoteTraceList[trace.serverName]) {
              this.publishRemoteTraceList[trace.serverName] = []
            }
            this.publishRemoteTraceList[trace.serverName].push(trace)
          }
          this.activeRomoteTracePane = Object.keys(
            this.publishRemoteTraceList
          )[0]
        })
        .finally(() => {
          this.traceLoading = false
        })
    },

    getPublishTraceDetail(data: PublishTraceData['datagram']) {
      this.traceDetail[data.id] = ''
      new DeployTraceDetail({ id: data.id }).request().then((response) => {
        this.traceDetail[data.id] =
          response.data.detail === ''
            ? this.$t('deployPage.noDetail')
            : response.data.detail
      })
    },
  },
})
</script>
<style scoped lang="scss">
@import '@/styles/mixin.scss';
.icon-success {
  color: #67c23a;
  font-size: 14px;
  font-weight: 900;
}

.icon-fail {
  color: #f56c6c;
  font-size: 14px;
  font-weight: 900;
}
.project-detail {
  padding-left: 5px;
  height: 470px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
