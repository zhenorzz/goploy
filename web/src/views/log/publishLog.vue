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
        <el-table-column prop="token" label="Token" width="300" />
        <el-table-column prop="publisherName" label="Username" width="100" />
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
        <el-table-column prop="insertTime" label="insertTime" width="155" />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="100"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <Button
              type="primary"
              link
              :permissions="[pms.DeployDetail]"
              @click="handleDetail(scope.row)"
            >
              {{ $t('detail') }}
            </Button>
          </template>
        </el-table-column>
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
              <div class="project-title">
                <span style="margin-right: 5px">Repo</span>
                <span v-if="item.state === 1" class="icon-success"></span>
                <span v-else class="icon-fail"></span>
              </div>
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
            <el-row style="margin: 5px 0">
              <div class="project-title">
                <span style="margin-right: 5px">After pull</span>
                <span v-if="item.state === 1" class="icon-success"></span>
                <span v-else class="icon-fail"></span>
              </div>
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
                type="primary"
                link
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
                  <div class="project-title">
                    <span style="margin-right: 5px">Before deploy</span>
                    <span v-if="trace.state === 1" class="icon-success"></span>
                    <span v-else class="icon-fail"></span>
                  </div>
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
                    type="primary"
                    link
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
                  <div class="project-title">
                    <span style="margin-right: 5px">Sync</span>
                    <span v-if="trace.state === 1" class="icon-success"></span>
                    <span v-else class="icon-fail"></span>
                  </div>
                </el-row>
                <el-row style="margin: 5px 0">
                  Time: {{ trace.updateTime }}
                </el-row>
                <el-row>Command: {{ trace.command }}</el-row>
                <div v-loading="traceDetail[trace.id] === ''">
                  <span style="padding: 5px 0">[goploy ~]#</span>
                  <el-button
                    v-if="trace.state === 1 && !(trace.id in traceDetail)"
                    type="primary"
                    link
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
                  <div class="project-title">
                    <span style="margin-right: 5px">After deploy</span>
                    <span v-if="trace.state === 1" class="icon-success"></span>
                    <span v-else class="icon-fail"></span>
                  </div>
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
                    type="primary"
                    link
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
export default { name: 'PublishLog' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Search } from '@element-plus/icons-vue'
import { PublishLogData, PublishLogList, PublishLogTotal } from '@/api/log'
import {
  DeployTrace,
  DeployTraceDetail,
  PublishTraceData,
  PublishTraceExt,
} from '@/api/deploy'
import { parseTime } from '@/utils'
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()

const dialogVisible = ref(false)
const searchParam = ref({ username: '', projectName: '' })
const tableLoading = ref(false)
const tableData = ref<PublishLogList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20, total: 0 })
const traceLoading = ref(false)
const traceDetail = ref({} as Record<number, string>)
const activeRomoteTracePane = ref('')
const publishLocalTraceList = ref<(PublishTraceData & PublishTraceExt)[]>([])
const publishRemoteTraceList = ref(
  {} as Record<string, (PublishTraceData & PublishTraceExt)[]>
)
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
  new PublishLogList(searchParam.value, pagination.value)
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}
function getTotal() {
  new PublishLogTotal(searchParam.value).request().then((response) => {
    pagination.value.total = response.data.total
  })
}
function handlePageChange(val = 1) {
  pagination.value.page = val
  getList()
}
function handleDetail(data: PublishLogData) {
  dialogVisible.value = true
  traceLoading.value = true
  new DeployTrace({ lastPublishToken: data.token })
    .request()
    .then((response) => {
      const publishTraceList = response.data.list.map((element) => {
        if (element.ext !== '') {
          Object.assign(element, JSON.parse(element.ext))
        }
        return <PublishTraceData & PublishTraceExt>element
      })

      publishLocalTraceList.value = publishTraceList.filter(
        (element) => element.type < 4
      )
      publishRemoteTraceList.value = {}
      for (const trace of publishTraceList) {
        if (trace.detail !== '') {
          traceDetail.value[trace.id] = trace.detail
        }
        if (trace.type < 4) continue
        if (!publishRemoteTraceList.value[trace.serverName]) {
          publishRemoteTraceList.value[trace.serverName] = []
        }
        publishRemoteTraceList.value[trace.serverName].push(trace)
      }
      activeRomoteTracePane.value = Object.keys(publishRemoteTraceList.value)[0]
    })
    .finally(() => {
      traceLoading.value = false
    })
}

function getPublishTraceDetail(data: PublishTraceData) {
  traceDetail.value[data.id] = ''
  new DeployTraceDetail({ id: data.id }).request().then((response) => {
    traceDetail.value[data.id] =
      response.data.detail === ''
        ? t('deployPage.noDetail')
        : response.data.detail
  })
}
</script>

<style scoped lang="scss">
@import '@/styles/mixin.scss';
.icon-success {
  color: #67c23a;
  font-size: 15px;
  &::before {
    content: '\2713';
  }
}

.icon-fail {
  color: #f56c6c;
  font-size: 15px;
  &::before {
    content: '\2717';
  }
}
.project {
  &-detail {
    padding-left: 5px;
    height: 470px;
    overflow-y: auto;
    @include scrollBar();
  }
  &-title {
    display: flex;
    flex-direction: row;
    width: 100%;
    &:before,
    &:after {
      content: '';
      flex: 1 1;
      border-bottom: 1px solid;
      margin: auto;
    }
    &:before {
      margin-right: 10px;
    }
    &:after {
      margin-left: 10px;
    }
  }
}
</style>
