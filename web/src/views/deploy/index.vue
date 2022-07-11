<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" align="middle">
      <el-select
        v-model="searchProject.sort"
        style="width: 130px"
        placeholder="Sort"
        clearable
        @change="sortChange"
      >
        <el-option :label="'ID Asc'" value="idAsc" />
        <el-option :label="'ID Desc'" value="idDesc" />
        <el-option :label="'Name Asc'" value="nameAsc" />
        <el-option :label="'Name Desc'" value="nameDesc" />
        <el-option :label="'Env Asc'" value="envAsc" />
        <el-option :label="'Env Desc'" value="envDesc" />
      </el-select>
      <el-select
        v-model="searchProject.environment"
        placeholder="Environment"
        clearable
      >
        <el-option :label="$t('envOption[1]')" :value="1" />
        <el-option :label="$t('envOption[2]')" :value="2" />
        <el-option :label="$t('envOption[3]')" :value="3" />
        <el-option :label="$t('envOption[4]')" :value="4" />
      </el-select>
      <el-select
        v-model="searchProject.autoDeploy"
        placeholder="Auto deploy"
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
    <el-row class="app-table">
      <el-scrollbar style="width: 100%">
        <el-row style="width: 100%" :gutter="10">
          <el-col
            v-for="(row, index) in tablePage.list"
            :key="index"
            style="margin-bottom: 10px"
            :sm="12"
            :md="8"
            :lg="8"
            :xl="6"
          >
            <el-card shadow="hover" :body-style="{ padding: '0px' }">
              <el-row
                align="middle"
                justify="space-between"
                style="position: relative; padding: 15px"
              >
                <el-row
                  style="flex: 1; flex-direction: column; flex-wrap: nowrap"
                >
                  <el-row>
                    <span
                      style="
                        font-size: 15px;
                        font-weight: 600;
                        white-space: nowrap;
                      "
                    >
                      <span
                        v-if="row.environment === 1"
                        style="color: var(--el-color-danger)"
                      >
                        {{ row.name }} -
                        {{ $t(`envOption[${row.environment || 0}]`) }}
                      </span>
                      <span
                        v-else-if="row.environment === 3"
                        style="color: var(--el-color-warning)"
                      >
                        {{ row.name }} -
                        {{ $t(`envOption[${row.environment || 0}]`) }}
                      </span>
                      <span v-else style="color: var(--el-color-info)">
                        {{ row.name }} -
                        {{ $t(`envOption[${row.environment || 0}]`) }}
                      </span>
                    </span>
                  </el-row>
                  <el-row style="margin-top: 8px" align="middle">
                    <svg-icon style="margin-right: 5px" icon-class="branch" />
                    <RepoURL
                      :url="row['url']"
                      :suffix="'/tree/' + row['branch'].split('/').pop()"
                      :text="row.branch"
                    >
                    </RepoURL>
                    <svg-icon style="margin: 0 5px" icon-class="gitCommit" />
                    <RepoURL
                      :url="row['url']"
                      :suffix="'/commit/' + row['commit']"
                      :text="row['commit'] ? row['commit'].substring(0, 6) : ''"
                    >
                    </RepoURL>
                  </el-row>
                  <el-row style="margin-top: 8px" align="middle">
                    <svg-icon icon-class="publishTime" />
                    <span style="margin: 0 5px; font-size: 14px">
                      {{ row.updateTime }}
                    </span>
                    <el-tag :type="row.tagType" size="small" effect="plain">
                      {{ row.tagText }}
                    </el-tag>
                  </el-row>
                </el-row>
                <el-progress
                  style="margin: 5px 0; width: 100%"
                  :percentage="row.progressPercentage"
                  :status="row.progressStatus"
                />
                <el-row justify="end">
                  <Button
                    v-if="row.deployState === 0"
                    :permissions="[pms.DeployProject]"
                    type="primary"
                    size="small"
                    @click="publish(row)"
                  >
                    {{ $t('initial') }}
                  </Button>
                  <Button
                    v-else-if="row.deployState === 1"
                    :permissions="[pms.DeployResetState]"
                    type="primary"
                    size="small"
                    @click="resetState(row)"
                  >
                    {{ $t('deployPage.resetState') }}
                  </Button>
                  <Dropdown
                    v-else
                    :permissions="[pms.DeployProject]"
                    :split-button="row.review === 1 ? false : true"
                    trigger="click"
                    type="primary"
                    size="small"
                    @click="publish(row)"
                    @command="(funcName: string) => commandFunc[funcName](row)"
                  >
                    <el-button
                      v-if="row.review === 1"
                      size="small"
                      type="primary"
                    >
                      {{ $t('submit') }}
                      <el-icon class="el-icon--right"><arrow-down /></el-icon>
                    </el-button>
                    <span v-else>{{ $t('deploy') }}</span>
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
                  </Dropdown>
                  <el-dropdown
                    trigger="click"
                    style="margin-left: 5px"
                    @command="(funcName) => commandFunc[funcName](row)"
                  >
                    <el-button size="small" type="warning">
                      {{ $t('func') }}
                      <el-icon class="el-icon--right"><arrow-down /></el-icon>
                    </el-button>
                    <template #dropdown>
                      <el-dropdown-menu
                        style="min-width: 84px; text-align: center"
                      >
                        <DropdownItem
                          :permissions="[pms.DeployTask]"
                          :command="'handleTaskCommand'"
                        >
                          {{ $t('deployPage.taskDeploy') }}
                        </DropdownItem>
                        <DropdownItem
                          :permissions="[pms.FileCompare]"
                          :command="'handleFileCompareCommand'"
                        >
                          {{ $t('deployPage.fileCompare') }}
                        </DropdownItem>
                        <DropdownItem
                          :permissions="[pms.FileSync]"
                          :command="'handleFileSyncCommand'"
                        >
                          {{ $t('deployPage.fileSync') }}
                        </DropdownItem>
                        <DropdownItem
                          :permissions="[pms.ProcessManager]"
                          :command="'handleProcessManagerCommand'"
                        >
                          {{ $t('deployPage.processManager') }}
                        </DropdownItem>
                        <DropdownItem
                          v-if="row.review === 1"
                          :permissions="[pms.DeployReview]"
                          :command="'handleReviewCommand'"
                        >
                          {{ $t('deployPage.reviewDeploy') }}
                        </DropdownItem>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                  <Button
                    type="success"
                    size="small"
                    style="margin-left: 5px"
                    :permissions="[pms.DeployDetail]"
                    @click="handleDetail(row)"
                  >
                    {{ $t('detail') }}
                  </Button>
                </el-row>
                <el-row
                  style="
                    top: 15px;
                    right: 15px;
                    position: absolute;
                    color: var(--el-text-color-secondary);
                  "
                >
                  # {{ row.id }}
                </el-row>
              </el-row>
            </el-card>
          </el-col>
        </el-row>
      </el-scrollbar>
      <!-- <el-table
        v-loading="tableloading"
        border
        stripe
        height="0"
        highlight-current-row
        :data="tablePage.list"
        :default-sort="tableDefaultSort"
        @sort-change="sortChange"
      >
        <el-table-column
          sortable="custom"
          prop="id"
          label="ID"
          width="70"
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
            <span v-if="scope.row.environment === 1" style="color: #f56c6c">
              {{ scope.row.name }} -
              {{ $t(`envOption[${scope.row.environment || 0}]`) }}
            </span>
            <span
              v-else-if="scope.row.environment === 3"
              style="color: #e6a23c"
            >
              {{ scope.row.name }} -
              {{ $t(`envOption[${scope.row.environment || 0}]`) }}
            </span>
            <span v-else style="color: #909399">
              {{ scope.row.name }} -
              {{ $t(`envOption[${scope.row.environment || 0}]`) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column
          prop="branch"
          :label="$t('branch')"
          width="180"
          align="center"
        >
          <template #default="scope">
            <RepoURL
              :url="scope.row['url']"
              :suffix="'/tree/' + scope.row['branch'].split('/').pop()"
              :text="scope.row.branch"
            >
            </RepoURL>
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
              <RepoURL
                :url="scope.row['url']"
                :suffix="'/commit/' + scope.row['commit']"
                :text="
                  scope.row['commit'] ? scope.row['commit'].substring(0, 6) : ''
                "
              >
              </RepoURL>
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
          width="155"
          align="center"
          sortable="custom"
        />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="300"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
          align="center"
        >
          <template #default="scope">
            <div class="operation-btn">
              <Button
                v-if="scope.row.deployState === 0"
                :permissions="[pms.DeployProject]"
                type="primary"
                @click="publish(scope.row)"
              >
                {{ $t('initial') }}
              </Button>
              <Button
                v-else-if="scope.row.deployState === 1"
                :permissions="[pms.DeployResetState]"
                type="primary"
                @click="resetState(scope.row)"
              >
                {{ $t('deployPage.resetState') }}
              </Button>
              <Dropdown
                v-else
                :permissions="[pms.DeployProject]"
                :split-button="scope.row.review === 1 ? false : true"
                trigger="click"
                type="primary"
                @click="publish(scope.row)"
                @command="(funcName: string) => commandFunc[funcName](scope.row)"
              >
                <el-button v-if="scope.row.review === 1" type="primary">
                  {{ $t('submit') }}
                  <el-icon class="el-icon--right"><arrow-down /></el-icon>
                </el-button>
                <span v-else>{{ $t('deploy') }}</span>
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
              </Dropdown>
              <el-dropdown
                trigger="click"
                style="margin-left: 5px"
                @command="(funcName) => commandFunc[funcName](scope.row)"
              >
                <el-button type="warning">
                  {{ $t('func') }}
                  <el-icon class="el-icon--right"><arrow-down /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu style="min-width: 84px; text-align: center">
                    <DropdownItem
                      :permissions="[pms.DeployTask]"
                      :command="'handleTaskCommand'"
                    >
                      {{ $t('deployPage.taskDeploy') }}
                    </DropdownItem>
                    <DropdownItem
                      :permissions="[pms.FileCompare]"
                      :command="'handleFileCompareCommand'"
                    >
                      {{ $t('deployPage.fileCompare') }}
                    </DropdownItem>
                    <DropdownItem
                      :permissions="[pms.FileSync]"
                      :command="'handleFileSyncCommand'"
                    >
                      {{ $t('deployPage.fileSync') }}
                    </DropdownItem>
                    <DropdownItem
                      :permissions="[pms.ProcessManager]"
                      :command="'handleProcessManagerCommand'"
                    >
                      {{ $t('deployPage.processManager') }}
                    </DropdownItem>
                    <DropdownItem
                      v-if="scope.row.review === 1"
                      :permissions="[pms.DeployReview]"
                      :command="'handleReviewCommand'"
                    >
                      {{ $t('deployPage.reviewDeploy') }}
                    </DropdownItem>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
              <Button
                type="success"
                style="margin-left: 5px"
                :permissions="[pms.DeployDetail]"
                @click="handleDetail(scope.row)"
              >
                {{ $t('detail') }}
              </Button>
            </div>
          </template>
        </el-table-column>
      </el-table> -->
    </el-row>
    <el-row type="flex" justify="end" style="width: 100%; margin-top: 5px">
      <el-pagination
        v-model:current-page="pagination.page"
        :total="tablePage.total"
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
        <Button
          type="primary"
          :permissions="[pms.DeployProject]"
          @click="publishByCommit(scope.row)"
        >
          {{ $t('deploy') }}
        </Button>
        <Button
          type="warning"
          :permissions="[pms.GreyDeploy]"
          @click="handleGreyPublish(scope.row)"
        >
          {{ $t('grey') }}
        </Button>
      </template>
    </TheCommitListDialog>
    <TheTagListDialog v-model="tagDialogVisible" :project-row="selectedItem">
      <template #tableOP="scope">
        <Button
          type="primary"
          :permissions="[pms.DeployProject]"
          @click="publishByCommit(scope.row)"
        >
          {{ $t('deploy') }}
        </Button>
        <Button
          type="warning"
          :permissions="[pms.GreyDeploy]"
          @click="handleGreyPublish(scope.row)"
        >
          {{ $t('grey') }}
        </Button>
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
    <TheFileSyncDialog
      v-model="fileSyncDialogVisible"
      :project-id="selectedItem.id"
    />
  </el-row>
</template>
<script lang="ts">
export default { name: 'DeployIndex' }
</script>
<script lang="ts" setup>
import pms from '@/permission'
import { Button, Dropdown, DropdownItem } from '@/components/Permission'
import { ArrowDown } from '@element-plus/icons-vue'
import {
  DeployList,
  DeployPublish,
  DeployResetState,
  DeployGreyPublish,
} from '@/api/deploy'
import { ProjectServerList, ProjectData } from '@/api/project'
import RepoURL from '@/components/RepoURL/index.vue'
import { parseTime } from '@/utils'
import TheDetailDialog from './TheDetailDialog.vue'
import TheCommitListDialog from './TheCommitListDialog.vue'
import TheTagListDialog from './TheTagListDialog.vue'
import TheTaskListDialog from './TheTaskListDialog.vue'
import TheReviewListDialog from './TheReviewListDialog.vue'
import TheProcessManagerDialog from './TheProcessManagerDialog.vue'
import TheFileCompareDialog from './TheFileCompareDialog.vue'
import TheFileSyncDialog from './TheFileSyncDialog.vue'
import type { ElForm } from 'element-plus'
import { ElMessageBox, ElMessage } from 'element-plus'
import { computed, watch, h, ref } from 'vue'
import { CommitData } from '@/api/repository'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const store = useStore()
const commitDialogVisible = ref(false)
const tagDialogVisible = ref(false)
const greyServerDialogVisible = ref(false)
const taskListDialogVisible = ref(false)
const fileSyncDialogVisible = ref(false)
const fileCompareDialogVisible = ref(false)
const processManagerDialogVisible = ref(false)
const reviewListDialogVisible = ref(false)
const dialogVisible = ref(false)
const searchProject = ref({
  sort: getTableSort(),
  name: '',
  environment: '',
  autoDeploy: '',
})
const selectedItem = ref({} as ProjectData)
const tableloading = ref(false)
const tableData = ref<any[]>([])
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
const greyServerFormRules = <InstanceType<typeof ElForm>['rules']>{
  serverIds: [
    {
      type: 'array',
      required: true,
      message: 'Server required',
      trigger: 'change',
    },
  ],
}
const tablePage = computed(() => {
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
  return {
    list: _tableData.slice(
      (pagination.value.page - 1) * pagination.value.rows,
      pagination.value.page * pagination.value.rows
    ),
    total: _tableData.length,
  }
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
      tableData.value = response.data.list.map((item) => {
        let element: any = item
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
      sortChange(searchProject.value.sort)
    })
    .finally(() => {
      tableloading.value = false
    })
}

function sortChange(sort: string) {
  setTableSort(sort)
  let prop: string
  let order: string

  switch (sort) {
    case 'idAsc':
      prop = 'id'
      order = 'asc'
      break
    case 'idDesc':
      prop = 'id'
      order = 'desc'
      break
    case 'nameAsc':
      prop = 'name'
      order = 'asc'
      break
    case 'nameDesc':
      prop = 'name'
      order = 'desc'
      break
    case 'envAsc':
      prop = 'environment'
      order = 'asc'
      break
    case 'envDesc':
      prop = 'environment'
      order = 'desc'
      break
    default:
      prop = 'id'
      order = 'desc'
      break
  }
  tableData.value = tableData.value.sort(
    (row1: ProjectData, row2: ProjectData): number => {
      let val1 = row1[prop]
      let val2 = row2[prop]
      if (order === 'desc') {
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

function handleDetail(data: ProjectData) {
  selectedItem.value = data
  dialogVisible.value = true
}

function handleRebuilt() {
  const projectIndex = tableData.value.findIndex(
    (element) => element.id === selectedItem.value.id
  )
  tableData.value[projectIndex].deployState = 1
}

function handleGreyPublish(data: CommitData) {
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

const commandFunc: { [K: string]: (data: ProjectData) => void } = {
  handleCommitCommand,
  handleTagCommand,
  handleTaskCommand,
  handleFileCompareCommand,
  handleFileSyncCommand,
  handleProcessManagerCommand,
  handleReviewCommand,
}

function handleCommitCommand(data: ProjectData) {
  selectedItem.value = data
  commitDialogVisible.value = true
}

function handleTagCommand(data: ProjectData) {
  selectedItem.value = data
  tagDialogVisible.value = true
}

function handleTaskCommand(data: ProjectData) {
  selectedItem.value = data
  taskListDialogVisible.value = true
}

function handleFileCompareCommand(data: ProjectData) {
  selectedItem.value = data
  fileCompareDialogVisible.value = true
}

function handleFileSyncCommand(data: ProjectData) {
  selectedItem.value = data
  fileSyncDialogVisible.value = true
}

function handleProcessManagerCommand(data: ProjectData) {
  selectedItem.value = data
  processManagerDialogVisible.value = true
}

function handleReviewCommand(data: ProjectData) {
  selectedItem.value = data
  reviewListDialogVisible.value = true
}

function publish(data: ProjectData) {
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

function publishByCommit(data: CommitData) {
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

function resetState(data: ProjectData) {
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

function getTableSort(): string {
  const sortStr = localStorage.getItem('deploy-sort')
  if (sortStr) {
    return sortStr
  }
  return 'idDesc'
}

function setTableSort(value: string) {
  localStorage.setItem('deploy-sort', value)
}
</script>
