<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('detail')"
    custom-class="publish-record"
  >
    <el-row type="flex">
      <div v-loading="filterloading" class="publish-preview">
        <div>
          <el-popover placement="bottom-start" width="318" trigger="click">
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
                class="dmp-date-picker"
                popper-class="dmp-date-picker-popper"
                :shortcuts="shortcuts"
                type="daterange"
                format="yyyy-MM-dd HH:mm:ss"
                range-separator="-"
                :start-placeholder="$t('startDate')"
                :end-placeholder="$t('endDate')"
                :default-time="['00:00:00', '23:59:59']"
                style="flex: 1"
              />
            </el-row>
            <el-row type="flex" align="middle" style="margin-top: 5px">
              <label class="publish-filter-label">
                {{ $t('deployDate') }}
              </label>
              <el-date-picker
                v-model="filterParams.deployDate"
                class="dmp-date-picker"
                popper-class="dmp-date-picker-popper"
                :shortcuts="shortcuts"
                type="daterange"
                format="yyyy-MM-dd HH:mm:ss"
                range-separator="-"
                :start-placeholder="$t('startDate')"
                :end-placeholder="$t('endDate')"
                :default-time="['00:00:00', '23:59:59']"
                style="flex: 1"
              />
            </el-row>
            <template #reference>
              <el-button icon="el-icon-notebook-2" style="width: 220px">
                Filter({{ filterlength }})
              </el-button>
            </template>
          </el-popover>
          <el-button
            type="warning"
            icon="el-icon-refresh"
            @click="refreshFilterParams"
          />
          <el-button
            :loading="filterloading"
            type="primary"
            icon="el-icon-search"
            style="margin-left: 2px"
            @click="searchPreviewList"
          />
        </div>
        <el-radio-group v-model="publishToken" @change="handleTraceChange">
          <el-row v-for="(item, index) in gitTraceList" :key="index">
            <el-row style="margin: 5px 0">
              <el-radio class="publish-commit" :label="item.token" border>
                <span class="publish-name">{{ item.publisherName }}</span>
                <span class="publish-commitID">
                  commitID: {{ item.commit }}
                </span>
                <i
                  v-if="item.publishState === 1"
                  class="el-icon-check icon-success"
                  style="float: right"
                />
                <i
                  v-else
                  class="el-icon-close icon-fail"
                  style="float: right"
                />
              </el-radio>
              <el-button type="danger" plain @click="rebuild(item)">
                rebuild
              </el-button>
            </el-row>
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
        style="width: 100%; flex: 1; align-content: flex-start"
      >
        <el-row v-for="(item, index) in publishLocalTraceList" :key="index">
          <el-row v-if="item.type === 2">
            <el-row style="margin: 5px 0">
              <i v-if="item.state === 1" class="el-icon-check icon-success" />
              <i v-else class="el-icon-close icon-fail" />
              -------------GIT-------------
            </el-row>
            <el-row style="margin: 5px 0">Time: {{ item.updateTime }}</el-row>
            <el-row v-if="item.state !== 0">
              <el-row>Branch: {{ item['branch'] }}</el-row>
              <el-row>
                Commit:
                <el-link
                  type="primary"
                  :underline="false"
                  :href="`${gitURL}/commit/${item['commit']}`"
                  target="_blank"
                >
                  {{ item['commit'] }}
                </el-link>
              </el-row>
              <el-row>Message: {{ item['message'] }}</el-row>
              <el-row>Author: {{ item['author'] }}</el-row>
              <el-row>
                Datetime:
                {{ item['timestamp'] ? parseTime(item['timestamp']) : '' }}
              </el-row>
              <el-row><span v-html="enterToBR(item['diff'])" /></el-row>
            </el-row>
            <el-row v-else style="margin: 5px 0">
              <span v-html="enterToBR(item.detail)" />
            </el-row>
          </el-row>
          <el-row v-if="item.type === 3">
            <hr />
            <el-row style="margin: 5px 0">
              <i v-if="item.state === 1" class="el-icon-check icon-success" />
              <i v-else class="el-icon-close icon-fail" />
              --------After pull--------
            </el-row>
            <el-row style="margin: 5px 0">Time: {{ item.updateTime }}</el-row>
            <el-row>
              Script:
              <pre v-html="enterToBR(item.script)"></pre>
            </el-row>
            <el-row
              v-loading="traceDetail[item.id] === true"
              style="margin: 5px 0"
            >
              [goploy ~]#
              <el-button
                v-if="item.state === 1 && !(item.id in traceDetail)"
                type="text"
                @click="getPublishTraceDetail(item)"
              >
                {{ $t('deployPage.showDetail') }}
              </el-button>
              <span v-else v-html="enterToBR(item.detail)" />
            </el-row>
          </el-row>
        </el-row>
        <el-tabs v-model="activeRomoteTracePane">
          <el-tab-pane
            v-for="(item, serverName) in publishRemoteTraceList"
            :key="serverName"
            :label="serverName"
            :name="serverName"
          >
            <el-row v-for="(trace, key) in item" :key="key">
              <el-row v-if="trace.type === 4">
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
                  <pre v-html="enterToBR(trace.script)"></pre>
                </el-row>
                <el-row
                  v-loading="traceDetail[trace.id] === true"
                  style="margin: 5px 0"
                >
                  [goploy ~]#
                  <el-button
                    v-if="trace.state === 1 && !(trace.id in traceDetail)"
                    type="text"
                    @click="getPublishTraceDetail(trace)"
                  >
                    {{ $t('deployPage.showDetail') }}
                  </el-button>
                  <span v-else v-html="enterToBR(trace.detail)" />
                </el-row>
              </el-row>
              <el-row v-else-if="trace.type === 5">
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
                <el-row
                  v-loading="traceDetail[trace.id] === true"
                  style="margin: 5px 0"
                >
                  [goploy ~]#
                  <el-button
                    v-if="trace.state === 1 && !(trace.id in traceDetail)"
                    type="text"
                    @click="getPublishTraceDetail(trace)"
                  >
                    {{ $t('deployPage.showDetail') }}
                  </el-button>
                  <span v-else v-html="enterToBR(trace.detail)" />
                </el-row>
              </el-row>
              <el-row v-else-if="trace.type === 6">
                <el-row style="margin: 5px 0">
                  <i
                    v-if="trace.state === 1"
                    class="el-icon-check icon-success"
                  />
                  <i v-else class="el-icon-close icon-fail" />
                  --------After deploy--------
                </el-row>
                <el-row style="margin: 5px 0"
                  >Time: {{ trace.updateTime }}</el-row
                >
                <el-row>Script: {{ trace.script }}</el-row>
                <el-row
                  v-loading="traceDetail[trace.id] === true"
                  style="margin: 5px 0"
                >
                  [goploy ~]#
                  <el-button
                    v-if="trace.state === 1 && !(trace.id in traceDetail)"
                    type="text"
                    @click="getPublishTraceDetail(trace)"
                    >{{ $t('deployPage.showDetail') }}</el-button
                  >
                  <span v-else v-html="enterToBR(trace.detail)" />
                </el-row>
              </el-row>
            </el-row>
          </el-tab-pane>
        </el-tabs>
      </el-row>
    </el-row>
  </el-dialog>
</template>

<script lang="ts">
import {
  DeployPreviewList,
  DeployTrace,
  DeployTraceDetail,
  PublishTraceData,
} from '@/api/deploy'
import { NamespaceUserOption } from '@/api/namespace'
import { role } from '@/utils/namespace'
import { empty, parseGitURL, parseTime } from '@/utils'
import { useI18n } from 'vue-i18n'
import { computed, watch, defineComponent, ref, reactive } from 'vue'

export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    projectRow: {
      type: Object,
      required: true,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const { t } = useI18n()
    const dialogVisible = computed({
      get: () => props.modelValue,
      set: (val) => {
        emit('update:modelValue', val)
      },
    })
    const userOption = ref<NamespaceUserOption['datagram']['list']>([])
    const gitURL = ref<string>('')
    watch(
      () => props.modelValue,
      (val: typeof props['modelValue']) => {
        if (val === true) {
          getPreviewList(props.projectRow.id)
          new NamespaceUserOption().request().then((response) => {
            userOption.value = response.data.list
          })
          gitURL.value = parseGitURL(props.projectRow.url)
        }
      }
    )
    const filterloading = ref(false)
    const filterParams = reactive<Record<string, any>>({
      loading: false,
      userId: '',
      state: '',
      filename: '',
      branch: '',
      commit: '',
      commitDate: [],
      deployDate: [],
    })
    const gitTraceList = ref<DeployPreviewList['datagram']['list']>([])
    const pagination = reactive({ page: 1, rows: 11, total: 0 })
    let traceDetail = reactive<Record<string, boolean>>({})
    const publishToken = ref('')
    const publishLocalTraceList = ref<DeployTrace['datagram']['list']>([])
    const publishRemoteTraceList = ref<
      Record<string, DeployTrace['datagram']['list']>
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
    const refreshFilterParams = () => {
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
      traceDetail = {}
      new DeployPreviewList(
        {
          projectId: projectId,
          commitDate: filterParams.commitDate
            ? filterParams.commitDate.join(',')
            : '',
          deployDate: filterParams.deployDate
            ? filterParams.deployDate.join(',')
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
          gitTraceList.value = response.data.list.map((element) => {
            if (element.ext !== '') {
              Object.assign(element, JSON.parse(element.ext))
              element.commit = element['commit']
                ? element['commit'].substring(0, 6)
                : ''
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
      handlePageChange(1)
    }

    const handlePageChange = (page: number) => {
      pagination.page = page
      getPreviewList(props.projectRow.id)
    }

    const traceLoading = ref(false)
    const publishTraceList = ref<DeployTrace['datagram']['list']>([])
    const activeRomoteTracePane = ref('')
    const getPublishTrace = (publishToken: string) => {
      traceLoading.value = true
      new DeployTrace({ lastPublishToken: publishToken })
        .request()
        .then((response) => {
          publishTraceList.value = response.data.list.map((element) => {
            if (element.ext !== '')
              Object.assign(element, JSON.parse(element.ext))
            return element
          })

          publishLocalTraceList.value = publishTraceList.value.filter(
            (element) => element.type < 4
          )
          publishRemoteTraceList.value = {}
          for (const trace of publishTraceList.value) {
            if (trace.type < 4) continue
            if (!publishRemoteTraceList.value[trace.serverName]) {
              publishRemoteTraceList.value[trace.serverName] = []
            }
            publishRemoteTraceList.value[trace.serverName].push(trace)
          }
          activeRomoteTracePane.value = Object.keys(
            publishRemoteTraceList.value
          )[0]
        })
        .finally(() => {
          traceLoading.value = false
        })
    }

    const handleTraceChange = (lastPublishToken: string) => {
      publishToken.value = lastPublishToken
      getPublishTrace(publishToken.value)
    }

    const getPublishTraceDetail = (
      data: PublishTraceData['datagram']['detail']
    ) => {
      traceDetail[data.id] = true
      new DeployTraceDetail({ id: data.id })
        .request()
        .then((response) => {
          data.detail =
            response.data.detail === ''
              ? t('deployPage.noDetail')
              : response.data.detail
        })
        .finally(() => {
          traceDetail[data.id] = false
        })
    }

    return {
      dialogVisible,
      role,
      parseTime,
      userOption,
      filterloading,
      filterParams,
      filterlength,
      refreshFilterParams,
      pagination,
      getPreviewList,
      searchPreviewList,
      handlePageChange,
      traceDetail,
      gitTraceList,
      publishToken,
      publishLocalTraceList,
      publishRemoteTraceList,
      traceLoading,
      activeRomoteTracePane,
      getPublishTrace,
      handleTraceChange,
      getPublishTraceDetail,
      gitURL,
    }
  },
  data() {
    return {
      shortcuts: [
        {
          text: '最近一周',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.$emit('pick', [start, end])
          },
        },
        {
          text: '最近一个月',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
            picker.$emit('pick', [start, end])
          },
        },
        {
          text: '最近三个月',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
            picker.$emit('pick', [start, end])
          },
        },
      ],
    }
  },
  methods: {
    enterToBR(detail: string) {
      return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
    },
  },
})
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
@import '@/styles/mixin.scss';
.publish {
  &-filter-label {
    font-size: 12px;
    width: 70px;
  }
  &-preview {
    width: 330px;
    margin-left: 10px;
  }
  &-commit {
    margin-right: 5px;
    padding-right: 8px;
    width: 246px;
    line-height: 12px;
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

@media screen and (max-width: 1440px) {
  .publish-record {
    :deep(.el-dialog) {
      width: 75%;
    }
  }
}
</style>
