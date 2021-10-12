<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('detail')"
    :fullscreen="$store.state.app.device === 'mobile'"
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
                rollback
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
        style="width: ; flex: 1; align-content: flex-start"
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
              <i v-if="item.state === 1" class="el-icon-check icon-success" />
              <i v-else class="el-icon-close icon-fail" />
              -------------GIT-------------
            </el-row>
            <el-row>Time: {{ item.updateTime }}</el-row>
            <template v-if="item.state !== 0">
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
            <el-row
              v-loading="traceDetail[item.id] === ''"
              style="margin: 5px 0"
            >
              <span style="padding: 5px 0">[goploy ~]#</span>
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
            </el-row>
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
    </el-row>
  </el-dialog>
</template>

<script lang="ts">
import {
  DeployPreviewList,
  DeployRebuild,
  DeployTrace,
  DeployTraceDetail,
  PublishTraceData,
} from '@/api/deploy'
import { NamespaceUserOption } from '@/api/namespace'
import { role } from '@/utils/namespace'
import { empty, parseGitURL, parseTime } from '@/utils'
import { ElMessageBox, ElMessage, ElDatePicker } from 'element-plus'
import { useI18n } from 'vue-i18n'
import dayjs from 'dayjs'
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
    onRebuilt: {
      type: Function,
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
          refreshFilterParams()
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
    const traceDetail = ref<Record<number, string>>({})
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
      traceDetail.value = {}
      new DeployPreviewList(
        {
          projectId: projectId,
          commitDate: filterParams.commitDate
            ? filterParams.commitDate
                .map((date: string) =>
                  dayjs(date).format('YYYY-MM-DD HH:mm:ss')
                )
                .join(',')
            : '',
          deployDate: filterParams.deployDate
            ? filterParams.deployDate
                .map((date: string) =>
                  dayjs(date).format('YYYY-MM-DD HH:mm:ss')
                )
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
    const activeRomoteTracePane = ref('')
    const getPublishTrace = (publishToken: string) => {
      traceLoading.value = true
      new DeployTrace({ lastPublishToken: publishToken })
        .request()
        .then((response) => {
          const publishTraceList = response.data.list.map((element) => {
            if (element.ext !== '') {
              Object.assign(element, JSON.parse(element.ext))
            }
            return element
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
      traceDetail.value[data.id] = ''
      new DeployTraceDetail({ id: data.id }).request().then((response) => {
        traceDetail.value[data.id] =
          response.data.detail === ''
            ? t('deployPage.noDetail')
            : response.data.detail
      })
    }

    const rebuild = (data: PublishTraceData['datagram']['detail']) => {
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
      rebuild,
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
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.$emit('pick', [start, end])
          },
        },
        {
          text: '最近一个月',
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
            picker.$emit('pick', [start, end])
          },
        },
        {
          text: '最近三个月',
          onClick(picker: typeof ElDatePicker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
            picker.$emit('pick', [start, end])
          },
        },
      ],
    }
  },
})
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
@import '@/styles/mixin.scss';
.publish {
  &-filter-label {
    font-size: 12px;
    width: 80px;
  }
  &-preview {
    width: 330px;
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
</style>
