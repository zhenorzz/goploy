<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex">
      <el-select
        v-model="searchProject.environment"
        placeholder="environment"
        clearable
      >
        <el-option :label="$t('envOption[1]')" :value="1" />
        <el-option :label="$t('envOption[2]')" :value="2" />
        <el-option :label="$t('envOption[3]')" :value="3" />
        <el-option :label="$t('envOption[4]')" :value="4" />
      </el-select>
      <el-select
        v-model="searchProject.autoDeploy"
        placeholder="auto deploy"
        clearable
      >
        <el-option :label="$t('close')" :value="0" />
        <el-option label="webhook" :value="1" />
      </el-select>
      <el-input
        v-model="searchProject.name"
        style="width: 300px"
        placeholder="Filter the project name"
      />
    </el-row>
    <el-table
      :key="tableHeight"
      v-loading="tableloading"
      border
      stripe
      highlight-current-row
      :max-height="tableHeight - 34"
      :data="tablePageData"
      :default-sort="tableDefaultSort"
      style="width: 100%; margin-top: 5px"
      @sort-change="sortChange"
    >
      <el-table-column
        sortable="custom"
        prop="id"
        label="ID"
        width="80"
        align="center"
      />
      <el-table-column
        prop="name"
        :label="$t('name')"
        min-width="150"
        align="center"
        sortable="custom"
      >
        <template #default="scope">
          <b v-if="scope.row.environment === 1" style="color: #f56c6c">
            {{ scope.row.name }} -
            {{ $t(`envOption[${scope.row.environment || 0}]`) }}
          </b>
          <b v-else-if="scope.row.environment === 3" style="color: #e6a23c">
            {{ scope.row.name }} -
            {{ $t(`envOption[${scope.row.environment || 0}]`) }}
          </b>
          <b v-else style="color: #909399">
            {{ scope.row.name }} -
            {{ $t(`envOption[${scope.row.environment || 0}]`) }}
          </b>
        </template>
      </el-table-column>
      <el-table-column
        prop="branch"
        :label="$t('branch')"
        width="150"
        align="center"
      >
        <template #default="scope">
          <el-link
            style="font-size: 12px"
            :underline="false"
            :href="
              parseGitURL(scope.row['url']) +
              '/tree/' +
              scope.row['branch'].split('/').pop()
            "
            target="_blank"
          >
            {{ scope.row.branch }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column
        prop="commit"
        label="CommitID"
        width="100"
        align="center"
      >
        <template #default="scope">
          <el-tooltip
            effect="dark"
            :content="scope.row['commit']"
            placement="top"
          >
            <el-link
              type="primary"
              style="font-size: 12px"
              :underline="false"
              :href="
                parseGitURL(scope.row['url']) + '/commit/' + scope.row['commit']
              "
              target="_blank"
            >
              {{
                scope.row['commit'] ? scope.row['commit'].substring(0, 6) : ''
              }}
            </el-link>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column
        prop="deployState"
        :label="$t('state')"
        width="200"
        align="center"
      >
        <template #default="scope">
          <el-tag :type="scope.row.tagType" effect="plain">
            {{ scope.row.tagText }}
          </el-tag>
          <el-progress
            :percentage="scope.row.progressPercentage"
            :status="scope.row.progressStatus"
          />
        </template>
      </el-table-column>
      <el-table-column
        prop="updateTime"
        :label="$t('time')"
        width="160"
        align="center"
        sortable="custom"
      />
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="310"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        align="center"
      >
        <template #default="scope">
          <div class="operation-btn">
            <el-button
              v-if="scope.row.deployState === 0"
              type="primary"
              @click="publish(scope.row)"
            >
              {{ $t('initial') }}
            </el-button>
            <el-button
              v-else-if="
                role.hasManagerPermission() && scope.row.deployState === 1
              "
              type="primary"
              @click="resetState(scope.row)"
            >
              {{ $t('deployPage.resetState') }}
            </el-button>
            <el-dropdown
              v-else-if="
                role.hasGroupManagerPermission() || scope.row.review === 0
              "
              split-button
              trigger="click"
              type="primary"
              @click="publish(scope.row)"
              @command="(funcName) => commandFunc[funcName](scope.row)"
            >
              {{
                role.isMember() && scope.row.review === 1
                  ? $t('submit')
                  : $t('deploy')
              }}
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="'handleCommitCommand'">
                    Commit list
                  </el-dropdown-item>
                  <el-dropdown-item :command="'handleTagCommand'">
                    Tag list
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button
              v-else
              type="primary"
              @click="handleCommitCommand(scope.row)"
            >
              {{ $t('deploy') }}
            </el-button>
            <el-dropdown
              v-if="role.hasGroupManagerPermission() || scope.row.review === 1"
              trigger="click"
              style="margin-left: 5px"
              @command="(funcName) => commandFunc[funcName](scope.row)"
            >
              <el-button type="warning">
                {{ $t('func') }}<i class="el-icon-arrow-down el-icon--right" />
              </el-button>
              <template #dropdown>
                <el-dropdown-menu style="min-width: 84px; text-align: center">
                  <el-dropdown-item
                    v-if="role.hasGroupManagerPermission()"
                    :command="'handleTaskCommand'"
                  >
                    {{ $t('deployPage.taskDeploy') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="role.hasGroupManagerPermission()"
                    :command="'handleFileCompareCommand'"
                  >
                    {{ $t('deployPage.fileCompare') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="role.hasManagerPermission()"
                    :command="'handleProcessManagerCommand'"
                  >
                    {{ $t('deployPage.processManager') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="scope.row.review === 1"
                    :command="'handleReviewCommand'"
                  >
                    {{ $t('deployPage.reviewDeploy') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button
              type="success"
              style="margin-left: 5px"
              @click="handleDetail(scope.row)"
            >
              {{ $t('detail') }}
            </el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="width: 100%; margin-top: 5px">
      <el-pagination
        v-model:current-page="pagination.page"
        :total="tablePageData.length"
        :page-size="pagination.rows"
        background
        :page-sizes="[20, 50, 100]"
        layout="total, prev, pager, next"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </el-row>
    <TheDetailDialog
      v-model="dialogVisible"
      :project-row="selectedItem"
      :on-rebuilt="handleRebuilt"
    />
    <TheCommitListDialog
      v-model="commitDialogVisible"
      :project-row="selectedItem"
    >
      <template #tableOP="scope">
        <el-button type="danger" @click="publishByCommit(scope.row)">
          {{ $t('deploy') }}
        </el-button>
        <el-button
          v-if="!role.isMember()"
          type="warning"
          @click="handleGreyPublish(scope.row)"
        >
          {{ $t('grey') }}
        </el-button>
      </template>
    </TheCommitListDialog>
    <TheTagListDialog v-model="tagDialogVisible" :project-row="selectedItem">
      <template #tableOP="scope">
        <el-button type="danger" @click="publishByCommit(scope.row)">
          {{ $t('deploy') }}
        </el-button>
        <el-button
          v-if="!role.isMember()"
          type="warning"
          @click="handleGreyPublish(scope.row)"
        >
          {{ $t('grey') }}
        </el-button>
      </template>
    </TheTagListDialog>
    <TheTaskListDialog
      v-model="taskListDialogVisible"
      :project-row="selectedItem"
    />
    <el-dialog v-model="greyServerDialogVisible" :title="$t('deploy')">
      <el-form
        ref="greyServerForm"
        :rules="greyServerFormRules"
        :model="greyServerFormData"
      >
        <el-form-item :label="$t('server')" label-width="80px" prop="serverIds">
          <el-checkbox-group v-model="greyServerFormData.serverIds">
            <el-checkbox
              v-for="(item, index) in greyServerFormProps.serverOption"
              :key="index"
              :label="item.serverId"
            >
              {{ item.serverName + '(' + item.serverDescription + ')' }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="greyServerDialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
        <el-button
          :disabled="greyServerFormProps.disabled"
          type="primary"
          @click="greyPublish"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
    <TheReviewListDialog
      v-model="reviewListDialogVisible"
      :project-row="selectedItem"
    />
    <TheFileCompareDialog
      v-model="fileCompareDialogVisible"
      :project-row="selectedItem"
    />
    <TheProcessManagerDialog
      v-model="processManagerDialogVisible"
      :project-row="selectedItem"
    />
  </el-row>
</template>
<script lang="ts">
export default { name: 'DeployIndex' }
</script>
<script lang="ts" setup>
import {
  DeployList,
  DeployPublish,
  DeployResetState,
  DeployGreyPublish,
} from '@/api/deploy'
import { ProjectServerList, ProjectData } from '@/api/project'
import { getRole } from '@/utils/namespace'
import { parseTime, parseGitURL } from '@/utils'
import TheDetailDialog from './TheDetailDialog.vue'
import TheCommitListDialog from './TheCommitListDialog.vue'
import TheTagListDialog from './TheTagListDialog.vue'
import TheTaskListDialog from './TheTaskListDialog.vue'
import TheReviewListDialog from './TheReviewListDialog.vue'
import TheProcessManagerDialog from './TheProcessManagerDialog.vue'
import TheFileCompareDialog from './TheFileCompareDialog.vue'
import getTableHeight from '@/composables/tableHeight'
import type { ElForm } from 'element-plus'
import { ElMessageBox, ElMessage } from 'element-plus'
import { computed, watch, h, ref } from 'vue'
import { CommitData } from '@/api/repository'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const store = useStore()
const role = getRole()
const commitDialogVisible = ref(false)
const tagDialogVisible = ref(false)
const greyServerDialogVisible = ref(false)
const taskListDialogVisible = ref(false)
const fileCompareDialogVisible = ref(false)
const processManagerDialogVisible = ref(false)
const reviewListDialogVisible = ref(false)
const dialogVisible = ref(false)
const searchProject = ref({ name: '', environment: '', autoDeploy: '' })
const selectedItem = ref({} as ProjectData['datagram'])
const { tableHeight } = getTableHeight()
const tableDefaultSort = getTableSort()
const tableloading = ref(false)
const tableData = ref<DeployList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const greyServerForm = ref<InstanceType<typeof ElForm>>()
const greyServerFormProps = ref({
  disabled: false,
  serverOption: [] as ProjectServerList['datagram']['list'],
})
const greyServerFormData = ref({
  projectId: 0,
  commit: '',
  serverIds: [],
})
const greyServerFormRules = {
  serverIds: [
    {
      type: 'array',
      required: true,
      message: 'Server required',
      trigger: 'change',
    },
  ],
}
const tablePageData = computed(() => {
  let _tableData = tableData.value
  if (searchProject.value.name !== '') {
    _tableData = tableData.value.filter(
      (item) => item.name.indexOf(searchProject.value.name) !== -1
    )
  }
  if (searchProject.value.environment !== '') {
    _tableData = tableData.value.filter(
      (item) => item.environment === Number(searchProject.value.environment)
    )
  }
  if (searchProject.value.autoDeploy !== '') {
    _tableData = tableData.value.filter(
      (item) => item.autoDeploy === Number(searchProject.value.autoDeploy)
    )
  }
  return _tableData.slice(
    (pagination.value.page - 1) * pagination.value.rows,
    pagination.value.page * pagination.value.rows
  )
})
watch(
  () => store.state.websocket.message,
  function (response) {
    if (response.type !== 1) {
      return
    }
    const data = response.message
    const message = enterToBR(data.message)
    const projectIndex = tableData.value.findIndex(
      (element) => element.id === data.projectId
    )
    if (projectIndex !== -1) {
      const percent = 20 * data.state
      tableData.value[projectIndex].progressPercentage = percent
      tableData.value[projectIndex].progressStatus = 'warning'
      tableData.value[projectIndex].tagType = 'warning'
      tableData.value[projectIndex].tagText = message
      tableData.value[projectIndex].deployState = 1
      if (percent === 0) {
        tableData.value[projectIndex].progressStatus = 'exception'
        tableData.value[projectIndex].tagType = 'danger'
        tableData.value[projectIndex].tagText = 'Fail'
        tableData.value[projectIndex].deployState = 3
      } else if (percent === 100) {
        tableData.value[projectIndex].progressStatus = 'success'
        tableData.value[projectIndex].tagType = 'success'
        tableData.value[projectIndex].deployState = 2
      }

      if (data['ext']) {
        Object.assign(tableData.value[projectIndex], data['ext'])
      }
      tableData.value[projectIndex].publisherName = data.username
      tableData.value[projectIndex].updateTime = parseTime(new Date().getTime())
    }
  }
)

getList()

function getList() {
  tableloading.value = true
  new DeployList()
    .request()
    .then((response) => {
      tableData.value = response.data.list.map((element) => {
        element.progressPercentage = 0
        element.tagType = 'info'
        element.tagText = 'Not deploy'
        if (element.deployState === 2) {
          element.progressPercentage = 100
          element.progressStatus = 'success'
          element.tagType = 'success'
          element.tagText = 'Success'
        } else if (element.deployState === 1) {
          element.progressPercentage = 60
          element.progressStatus = 'warning'
          element.tagType = 'warning'
          element.tagText = 'Deploying'
        } else if (element.deployState === 3) {
          element.progressPercentage = 0
          element.progressStatus = 'exception'
          element.tagType = 'danger'
          element.tagText = 'Fail'
        }
        try {
          Object.assign(element, JSON.parse(element.publishExt))
        } catch (error) {
          console.log('Project not deploy')
        }
        return element
      })
      sortChange(tableDefaultSort)
    })
    .finally(() => {
      tableloading.value = false
    })
}

function sortChange({ prop, order }: { prop: string; order: string }) {
  setTableSort(prop, order)
  if (!prop && !order) {
    prop = 'id'
    order = 'descending'
  }
  if (prop === 'name') {
    prop = 'environment'
  }
  tableData.value = tableData.value.sort(
    (row1: ProjectData['datagram'], row2: ProjectData['datagram']): number => {
      let val1
      let val2
      if (order === 'ascending') {
        val1 = row1[prop]
        val2 = row2[prop]
      } else if (order === 'descending') {
        val1 = row2[prop]
        val2 = row1[prop]
      }
      if (val1 < val2) {
        return -1
      } else if (val1 > val2) {
        return 1
      } else {
        return 0
      }
    }
  )
}

function handleSizeChange(val = 1) {
  pagination.value.rows = val
  handlePageChange(1)
}

function handlePageChange(page = 1) {
  pagination.value.page = page
}

function handleDetail(data: ProjectData['datagram']) {
  selectedItem.value = data
  dialogVisible.value = true
}

function handleRebuilt() {
  const projectIndex = tableData.value.findIndex(
    (element) => element.id === selectedItem.value.id
  )
  tableData.value[projectIndex].deployState = 1
}

function handleGreyPublish(data: CommitData['datagram']) {
  new ProjectServerList({ id: selectedItem.value.id })
    .request()
    .then((response) => {
      greyServerFormProps.value.serverOption = response.data.list
    })
  // add projectID to server form
  greyServerFormData.value.projectId = selectedItem.value.id
  greyServerFormData.value.commit = data.commit
  greyServerDialogVisible.value = true
}

const commandFunc: { [K: string]: (data: ProjectData['datagram']) => void } = {
  handleCommitCommand,
  handleTagCommand,
  handleTaskCommand,
  handleFileCompareCommand,
  handleProcessManagerCommand,
  handleReviewCommand,
}

function handleCommitCommand(data: ProjectData['datagram']) {
  selectedItem.value = data
  commitDialogVisible.value = true
}

function handleTagCommand(data: ProjectData['datagram']) {
  selectedItem.value = data
  tagDialogVisible.value = true
}

function handleTaskCommand(data: ProjectData['datagram']) {
  selectedItem.value = data
  taskListDialogVisible.value = true
}

function handleFileCompareCommand(data: ProjectData['datagram']) {
  selectedItem.value = data
  fileCompareDialogVisible.value = true
}

function handleProcessManagerCommand(data: ProjectData['datagram']) {
  selectedItem.value = data
  processManagerDialogVisible.value = true
}

function handleReviewCommand(data: ProjectData['datagram']) {
  selectedItem.value = data
  reviewListDialogVisible.value = true
}

function publish(data: ProjectData['datagram']) {
  const id = data.id
  let color = ''
  if (data.environment === 1) {
    color = 'color: #F56C6C'
  } else if (data.environment === 3) {
    color = 'color: #E6A23C'
  } else {
    color = 'color: #909399'
  }
  ElMessageBox.confirm('', t('tips'), {
    message: h('p', null, [
      h('span', null, 'Deploy Project: '),
      h(
        'b',
        { style: color },
        data.name + ' - ' + t(`envOption[${data.environment}]`)
      ),
    ]),
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      new DeployPublish({ projectId: id, commit: '', branch: '' })
        .request()
        .then(() => {
          const projectIndex = tableData.value.findIndex(
            (element) => element.id === id
          )
          tableData.value[projectIndex].deployState = 1
        })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function publishByCommit(data: CommitData['datagram']) {
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
      new DeployPublish({
        projectId: selectedItem.value.id,
        branch: data.branch,
        commit: data.commit,
      })
        .request()
        .then(() => {
          const projectIndex = tableData.value.findIndex(
            (element) => element.id === selectedItem.value.id
          )
          tableData.value[projectIndex].deployState = 1
          commitDialogVisible.value = false
          tagDialogVisible.value = false
        })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function resetState(data: ProjectData['datagram']) {
  ElMessageBox.confirm(t('deployPage.resetStateTips'), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      new DeployResetState({ projectId: data.id }).request().then(() => {
        const projectIndex = tableData.value.findIndex(
          (element) => element.id === data.id
        )
        tableData.value[projectIndex].deployState = 0
        tableData.value[projectIndex].progressPercentage = 0
        tableData.value[projectIndex].tagType = 'info'
        tableData.value[projectIndex].tagText = 'Not deploy'
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function greyPublish() {
  greyServerForm.value?.validate((valid) => {
    if (valid) {
      const data = greyServerFormData.value
      ElMessageBox.confirm(
        t('deployPage.publishCommitTips', {
          commit: data.commit,
        }),
        t('tips'),
        {
          confirmButtonText: t('confirm'),
          cancelButtonText: t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          new DeployGreyPublish({
            projectId: data.projectId,
            commit: data.commit,
            serverIds: data.serverIds,
          })
            .request()
            .then(() => {
              const projectIndex = tableData.value.findIndex(
                (element) => element.id === data.projectId
              )
              tableData.value[projectIndex].deployState = 1
              commitDialogVisible.value = false
              tagDialogVisible.value = false
              greyServerDialogVisible.value = false
            })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
}

function enterToBR(detail: string) {
  return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
}

function getTableSort(): { prop: string; order: string } {
  const sortJsonStr = localStorage.getItem('deploy-table-sort')
  if (sortJsonStr) {
    return JSON.parse(sortJsonStr) as { prop: string; order: string }
  }
  return { prop: 'id', order: 'descending' }
}

function setTableSort(prop: string, order: string) {
  localStorage.setItem('deploy-table-sort', JSON.stringify({ prop, order }))
}
</script>
