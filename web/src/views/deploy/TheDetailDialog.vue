<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('detail')"
    :fullscreen="$store.state.app.device === 'mobile'"
    width="75%"
    @close="onClose"
  >
    <el-row type="flex">
      <div v-loading="filterloading" class="publish-preview">
        <div>
          <el-popover
            v-model:visible="filterInpurtVisible"
            placement="bottom-start"
            width="318"
            trigger="manual"
          >
            <el-row type="flex" justify="space-between" align="middle">
              <span style="font-size: 16px">Filter</span>
              <el-button
                type="text"
                :icon="Close"
                style="color: #999; font-size: 16px; margin-bottom: 10px"
                @click="filterInpurtVisible = false"
              />
            </el-row>
            <el-row type="flex" align="middle">
              <label class="publish-filter-label">{{ $t('user') }}</label>
              <el-select
                v-model="filterParams.userId"
                style="flex: 1"
                clearable
              >
                <el-option
                  v-for="(item, index) in userOption"
                  :key="index"
                  :label="item.userName"
                  :value="item.userId"
                />
              </el-select>
            </el-row>
            <el-row type="flex" align="middle" style="margin-top: 5px">
              <label class="publish-filter-label">Commit</label>
              <el-input
                v-model.trim="filterParams.commit"
                autocomplete="off"
                style="flex: 1"
                placeholder="Commit"
              />
            </el-row>
            <el-row type="flex" align="middle" style="margin-top: 5px">
              <label class="publish-filter-label">{{ $t('branch') }}</label>
              <el-input
                v-model.trim="filterParams.branch"
                autocomplete="off"
                style="flex: 1"
                :placeholder="$t('branch')"
              />
            </el-row>
            <el-row type="flex" align="middle" style="margin-top: 5px">
              <label class="publish-filter-label">{{ $t('filename') }}</label>
              <el-input
                v-model.trim="filterParams.filename"
                autocomplete="off"
                style="flex: 1"
                :placeholder="$t('filename')"
              />
            </el-row>
            <el-row type="flex" align="middle" style="margin-top: 5px">
              <label class="publish-filter-label">{{ $t('state') }}</label>
              <el-select v-model="filterParams.state" style="flex: 1" clearable>
                <el-option :label="$t('success')" :value="1" />
                <el-option :label="$t('fail')" :value="0" />
              </el-select>
            </el-row>
            <el-row type="flex" align="middle" style="margin-top: 5px">
              <label class="publish-filter-label">{{ $t('commitDate') }}</label>
              <el-date-picker
                v-model="filterParams.commitDate"
                :shortcuts="shortcuts"
                type="daterange"
                format="YYYY-MM-DD"
                range-separator="-"
                :start-placeholder="$t('startDate')"
                :end-placeholder="$t('endDate')"
                :default-time="[
                  new Date(2000, 1, 1, 0, 0, 0),
                  new Date(2000, 2, 1, 23, 59, 59),
                ]"
                style="flex: 1"
              />
            </el-row>
            <el-row type="flex" align="middle" style="margin-top: 5px">
              <label class="publish-filter-label">
                {{ $t('deployDate') }}
              </label>
              <el-date-picker
                v-model="filterParams.deployDate"
                :shortcuts="shortcuts"
                type="daterange"
                format="YYYY-MM-DD"
                range-separator="-"
                :start-placeholder="$t('startDate')"
                :end-placeholder="$t('endDate')"
                :default-time="[
                  new Date(2000, 1, 1, 0, 0, 0),
                  new Date(2000, 2, 1, 23, 59, 59),
                ]"
                style="flex: 1"
              />
            </el-row>
            <el-row type="flex" justify="end">
              <el-button
                type="danger"
                style="margin-top: 10px"
                @click="clearFilterParams"
              >
                {{ $t('clear') }}
              </el-button>
              <el-button
                :loading="filterloading"
                type="primary"
                style="margin-top: 10px"
                @click="searchPreviewList"
              >
                {{ $t('search') }}
              </el-button>
            </el-row>
            <template #reference>
              <el-button
                :icon="Search"
                :loading="filterloading"
                style="width: 270px"
                @click="filterInpurtVisible = !filterInpurtVisible"
              >
                Filter({{ filterlength }})
              </el-button>
            </template>
          </el-popover>
          <el-button
            type="primary"
            :icon="Refresh"
            @click="searchPreviewList"
          />
        </div>
        <el-radio-group v-model="publishToken" @change="handleTraceChange">
          <el-row
            v-for="(item, index) in gitTraceList"
            :key="index"
            style="margin: 5px 0; width: 100%"
          >
            <el-radio class="publish-commit" :label="item.token" border>
              <span class="publish-name">{{ item.publisherName }}</span>
              <span
                v-if="projectRow.repoType === 'svn'"
                class="publish-commitID"
                :title="item['commit']"
              >
                revision: {{ item['commit'].substring(0, 6) }}
              </span>
              <span
                v-else-if="projectRow.repoType === 'sftp'"
                class="publish-commitID"
                :title="item['token']"
              >
                uuid: {{ item['token'].substring(0, 6) }}
              </span>
              <span
                v-else-if="projectRow.repoType === 'ftp'"
                class="publish-commitID"
                :title="item['token']"
              >
                uuid: {{ item['token'].substring(0, 6) }}
              </span>
              <span v-else class="publish-commitID" :title="item['commit']">
                commitID: {{ item['commit'].substring(0, 6) }}
              </span>
              <el-icon
                v-if="
                  item.token === projectRow.lastPublishToken &&
                  projectRow.deployState === 1
                "
                class="is-loading"
                style="font-size: 15px; float: right"
              >
                <Loading />
              </el-icon>
              <SvgIcon
                v-else-if="
                  item.token === projectRow.lastPublishToken &&
                  projectRow.deployState === 3
                "
                style="color: #f56c6c; font-size: 15px; float: right"
                icon-class="close"
              ></SvgIcon>
              <SvgIcon
                v-else-if="item.state === 1"
                style="color: #67c23a; font-size: 15px; float: right"
                icon-class="check"
              ></SvgIcon>
              <SvgIcon
                v-else
                style="color: #f56c6c; font-size: 15px; float: right"
                icon-class="close"
              ></SvgIcon>
            </el-radio>
            <el-button
              v-if="
                ['git', 'svn'].includes(projectRow.repoType) ||
                projectRow.symlinkPath != ''
              "
              type="danger"
              plain
              @click="rebuild(item)"
            >
              rollback
            </el-button>
          </el-row>
        </el-radio-group>
        <el-pagination
          v-model:current-page="pagination.page"
          :total="pagination.total"
          :page-size="pagination.rows"
          style="text-align: right; margin-right: 20px"
          layout="total, prev, next"
          @current-change="handlePageChange"
        />
      </div>
      <el-row
        v-loading="traceLoading"
        class="project-detail"
        style="flex: 1; align-content: flex-start"
        :style="{
          width: '100%',
          alignContent: 'flex-start',
          flex: $store.state.app.device === 'mobile' ? '' : 1,
        }"
      >
        <div
          v-for="(item, index) in publishLocalTraceList"
          :key="index"
          style="width: 100%"
        >
          <template v-if="item.type === 2">
            <el-row style="margin: 5px 0">
              <div class="project-title">
                <span style="margin-right: 5px">
                  {{ projectRow.repoType.toUpperCase() }}
                </span>
                <span v-if="item.state === 1" class="icon-success"></span>
                <span v-else class="icon-fail"></span>
              </div>
            </el-row>
            <el-row>Time: {{ item.updateTime }}</el-row>
            <template v-if="item.state !== 0">
              <el-row>Branch: {{ item['branch'] }}</el-row>
              <el-row>
                Commit:
                <RepoURL
                  :url="projectRow.url"
                  :suffix="`/commit/${item['commit']}`"
                  :text="item['commit']"
                >
                </RepoURL>
              </el-row>
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
            <el-row style="margin: 5px 0" class="project-title">
              <span style="margin-right: 5px">After pull</span>
              <span v-if="item.state === 1" class="icon-success"></span>
              <span v-else class="icon-fail"></span>
            </el-row>
            <el-row>Time: {{ item.updateTime }}</el-row>
            <el-row style="width: 100%">
              <div>Script:</div>
              <pre style="white-space: pre-line">{{ item.script }}</pre>
            </el-row>
            <div v-loading="traceDetail[item.id] === ''" style="margin: 5px 0">
              <span>[goploy ~]#</span>
              <el-button
                v-if="item.state === 1 && !(item.id in traceDetail)"
                type="text"
                @click="getPublishTraceDetail(item)"
              >
                {{ $t('deployPage.showDetail') }}
              </el-button>
              <span v-else style="white-space: pre-line; padding: 5px 0">
                {{ traceDetail[item.id] }}
              </span>
            </div>
          </div>
        </div>
        <el-tabs v-model="activeRomoteTracePane" style="width: 100%">
          <el-tab-pane
            v-for="(item, serverName) in publishRemoteTraceList"
            :key="serverName"
            :label="serverName"
            :name="serverName"
          >
            <div v-for="(trace, key) in item" :key="key">
              <template v-if="trace.type === 4">
                <el-row style="margin: 5px 0" class="project-title">
                  <span style="margin-right: 5px">Before deploy</span>
                  <span v-if="trace.state === 1" class="icon-success"></span>
                  <span v-else class="icon-fail"></span>
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
                <el-row style="margin: 5px 0" class="project-title">
                  <span style="margin-right: 5px">
                    {{ projectRow.transferType.toUpperCase() }}
                  </span>
                  <span v-if="trace.state === 1" class="icon-success"></span>
                  <span v-else class="icon-fail"></span>
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
                <el-row style="margin: 5px 0" class="project-title">
                  <span style="margin-right: 5px">After deploy</span>
                  <span v-if="trace.state === 1" class="icon-success"></span>
                  <span v-else class="icon-fail"></span>
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
    </el-row>
  </el-dialog>
</template>

<script lang="ts" setup>
import { Search, Refresh, Close, Loading } from '@element-plus/icons-vue'
import RepoURL from '@/components/RepoURL/index.vue'
import {
  DeployPreviewList,
  DeployRebuild,
  DeployTrace,
  DeployTraceDetail,
  PublishTraceData,
  PublishTraceExt,
} from '@/api/deploy'
import { ProjectData } from '@/api/project'
import { NamespaceUserOption } from '@/api/namespace'
import { empty, parseTime } from '@/utils'
import type { ElRadioGroup } from 'element-plus'
import { ElMessageBox, ElMessage, ElDatePicker } from 'element-plus'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
import { computed, watch, ref, reactive, PropType } from 'vue'
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  projectRow: {
    type: Object as PropType<ProjectData>,
    required: true,
  },
  onRebuilt: {
    type: Function,
    required: true,
  },
})
const emit = defineEmits(['update:modelValue'])
const { t } = useI18n()
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val)
  },
})
const shortcuts = [
  {
    text: t('lastWeek'),
    onClick(picker: typeof ElDatePicker) {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
      picker.emit('pick', [dayjs(start), dayjs(end)])
    },
  },
  {
    text: t('lastMonth'),
    onClick(picker: typeof ElDatePicker) {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
      picker.emit('pick', [dayjs(start), dayjs(end)])
    },
  },
  {
    text: t('last3Months'),
    onClick(picker: typeof ElDatePicker) {
      const end = new Date()
      const start = new Date()
      start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
      picker.emit('pick', [dayjs(start), dayjs(end)])
    },
  },
]
const userOption = ref<NamespaceUserOption['datagram']['list']>([])
watch(
  () => props.modelValue,
  (val: typeof props['modelValue']) => {
    if (val === true) {
      clearFilterParams()
      getPreviewList(props.projectRow.id)
      new NamespaceUserOption().request().then((response) => {
        userOption.value = response.data.list
      })
    }
  }
)

const filterloading = ref(false)
const filterInpurtVisible = ref(false)
const filterParams = reactive<Record<string, any>>({
  userId: '',
  state: '',
  filename: '',
  branch: '',
  commit: '',
  commitDate: [],
  deployDate: [],
})
const gitTraceList = ref<(PublishTraceData & PublishTraceExt)[]>([])
const pagination = reactive({ page: 1, rows: 11, total: 0 })
const traceDetail = ref<Record<number, string>>({})
const publishToken = ref('')
const publishLocalTraceList = ref<(PublishTraceData & PublishTraceExt)[]>([])
const publishRemoteTraceList = ref<
  Record<string, (PublishTraceData & PublishTraceExt)[]>
>({})
const filterlength = computed(() => {
  let number = 0
  for (const key in filterParams) {
    if (['projectId', 'loading', 'url'].indexOf(key) !== -1) {
      continue
    }
    if (!empty(filterParams[key])) {
      number++
    }
  }
  return number
})

const clearFilterParams = () => {
  filterParams.userId = ''
  filterParams.state = ''
  filterParams.filename = ''
  filterParams.branch = ''
  filterParams.commit = ''
  filterParams.commitDate = []
  filterParams.deployDate = []
}

const getPreviewList = (projectId: number) => {
  filterloading.value = true
  traceDetail.value = {}
  new DeployPreviewList(
    {
      projectId: projectId,
      commitDate: filterParams.commitDate
        ? filterParams.commitDate
            .map((date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'))
            .join(',')
        : '',
      deployDate: filterParams.deployDate
        ? filterParams.deployDate
            .map((date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss'))
            .join(',')
        : '',
      branch: filterParams.branch,
      commit: filterParams.commit,
      filename: filterParams.filename,
      userId: filterParams.userId !== '' ? Number(filterParams.userId) : 0,
      state: filterParams.state === '' ? -1 : Number(filterParams.state),
    },
    pagination
  )
    .request()
    .then((response) => {
      gitTraceList.value = response.data.list.map((item) => {
        let element = <PublishTraceData & PublishTraceExt>item
        if (element.ext) {
          element.commit = element.ext.replaceAll('"', '')
        }
        return element
      })
      if (gitTraceList.value.length > 0) {
        publishToken.value = gitTraceList.value[0].token
        getPublishTrace(publishToken.value)
      } else {
        publishLocalTraceList.value = []
        publishRemoteTraceList.value = {}
      }
      pagination.total = response.data.pagination.total
    })
    .finally(() => {
      filterloading.value = false
    })
}

const searchPreviewList = () => {
  filterInpurtVisible.value = false
  handlePageChange(1)
}

const handlePageChange = (page: number) => {
  pagination.page = page
  getPreviewList(props.projectRow.id)
}

const traceLoading = ref(false)
const activeRomoteTracePane = ref('')
let timeout: ReturnType<typeof setInterval>
function getPublishTrace(publishToken: string) {
  traceLoading.value = true && !timeout
  new DeployTrace({ lastPublishToken: publishToken })
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
      if (props.projectRow.lastPublishToken === publishToken) {
        if (props.projectRow.deployState === 1 && !timeout) {
          timeout = setInterval(() => {
            getPublishTrace(publishToken)
          }, 1000)
        } else if (props.projectRow.deployState !== 1) {
          clearInterval(timeout)
        }
      } else {
        clearInterval(timeout)
      }
    })
    .finally(() => {
      traceLoading.value = false
    })
}

const onClose = () => {
  filterInpurtVisible.value = false
  clearInterval(timeout)
}

const handleTraceChange: InstanceType<typeof ElRadioGroup>['onChange'] = (
  lastPublishToken
) => {
  publishToken.value = <string>lastPublishToken
  getPublishTrace(publishToken.value)
}

const getPublishTraceDetail = (data: PublishTraceData) => {
  traceDetail.value[data.id] = ''
  new DeployTraceDetail({ id: data.id }).request().then((response) => {
    traceDetail.value[data.id] =
      response.data.detail === ''
        ? t('deployPage.noDetail')
        : response.data.detail
  })
}

const rebuild = (data: PublishTraceData & PublishTraceExt) => {
  ElMessageBox.confirm(
    t('deployPage.publishCommitTips', { commit: data.commit }),
    t('tips'),
    {
      confirmButtonText: t('confirm'),
      cancelButtonText: t('cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      filterloading.value = true
      new DeployRebuild({ projectId: data.projectId, token: data.token })
        .request()
        .then((response) => {
          if (response.data === 'symlink') {
            ElMessage.success('Success')
          } else {
            props.onRebuilt()
          }
          dialogVisible.value = false
        })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
@import '@/styles/mixin.scss';
.publish {
  &-filter-label {
    width: 80px;
  }
  &-preview {
    width: 330px;
  }
  &-commit {
    margin-right: 5px;
    padding-right: 8px;
    flex: 1;
  }
  &-commitID {
    display: inline-block;
    vertical-align: top;
  }
  &-name {
    width: 60px;
    display: inline-block;
    text-align: center;
    overflow: hidden;
    vertical-align: top;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.icon-success {
  color: #67c23a;
  font-size: 15px;
  font-weight: 500;
  &::before {
    content: '\2713';
  }
}

.icon-fail {
  color: #f56c6c;
  font-size: 15px;
  font-weight: 500;
  &::before {
    content: '\2717';
  }
}
.project {
  &-detail {
    padding-left: 5px;
    height: 490px;
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
<style lang="scss">
.publish-commit {
  .el-radio__label {
    width: 100%;
  }
}
</style>
