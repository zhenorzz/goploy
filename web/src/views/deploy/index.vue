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
      <el-select v-model="searchProject.pin" placeholder="Show pin" clearable>
        <el-option label="Pin" :value="true" />
        <el-option label="Unpin" :value="false" />
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
        v-model="searchProject.label"
        :max-collapse-tags="1"
        style="width: 300px"
        multiple
        collapse-tags
        collapse-tags-tooltip
        placeholder="Filter the project label"
        clearable
      >
        <el-option
          v-for="item in labelList"
          :key="item"
          :label="item"
          :value="item"
        />
      </el-select>
      <el-input
        v-model="searchProject.name"
        style="width: 300px"
        placeholder="Filter the project name"
      />
    </el-row>
    <el-row class="app-table">
      <el-scrollbar style="width: 100%">
        <el-row
          v-infinite-scroll="load"
          style="width: 100%"
          :gutter="10"
          :infinite-scroll-disabled="disabled"
        >
          <el-col
            v-for="(row, index) in tablePage.list"
            :key="index"
            style="margin-bottom: 10px"
            :sm="12"
            :md="8"
            :lg="8"
            :xl="6"
          >
            <el-card
              shadow="hover"
              :body-style="{
                padding: '0px',
                backgroundColor: 'var(--el-bg-color)',
                position: 'relative',
              }"
              class="card"
              :class="
                row.deployState === DeployState.Uninitialized
                  ? ''
                  : row.deployState === DeployState.Deploying
                  ? 'deploying'
                  : row.deployState === DeployState.Success
                  ? 'success'
                  : 'fail'
              "
            >
              <div style="padding: 12px">
                <el-row justify="space-between">
                  <svg-icon
                    v-if="row['pin'] === true"
                    style="margin-right: 5px; color: var(--el-color-warning)"
                    icon-class="pin"
                  />
                  <el-tooltip
                    :disabled="row.label == ''"
                    effect="dark"
                    :content="row.label"
                    placement="bottom"
                  >
                    <el-row
                      align="middle"
                      class="card-title__text"
                      :style="envTitleStyle[row.environment || 0]"
                    >
                      <span style="margin-right: 4px">#{{ row.id }}</span>
                      <el-link
                        v-if="isLink(row.name)"
                        :href="row.name"
                        target="_blank"
                        :underline="false"
                        class="card-title__link"
                        style="color: inherit"
                      >
                        {{ row.name }}
                        <el-icon>
                          <Link />
                        </el-icon>
                      </el-link>
                      <span v-else>{{ row.name }}</span>
                      <span style="margin-left: 4px">
                        - {{ $t(`envOption[${row.environment || 0}]`) }}
                      </span>
                    </el-row>
                  </el-tooltip>

                  <el-dropdown
                    trigger="click"
                    @command="(funcName: string) => cardMoreFunc[funcName](row)"
                  >
                    <el-button link :icon="More" />
                    <template #dropdown>
                      <el-dropdown-menu>
                        <el-dropdown-item :command="'handlePinCard'">
                          Pin
                        </el-dropdown-item>
                        <el-dropdown-item :command="'handleUnpinCard'">
                          Unpin
                        </el-dropdown-item>
                      </el-dropdown-menu>
                    </template>
                  </el-dropdown>
                </el-row>
                <el-row style="margin-top: 8px" align="middle">
                  <svg-icon style="margin-right: 5px" icon-class="branch" />
                  <RepoURL
                    style="font-size: 14px"
                    :url="row['url']"
                    :suffix="'/tree/' + row['branch'].split('/').pop()"
                    :text="row.branch"
                  >
                  </RepoURL>
                  <svg-icon style="margin: 0 5px" icon-class="gitCommit" />
                  <RepoURL
                    style="font-size: 14px"
                    :url="row['url']"
                    :suffix="'/commit/' + row['commit']"
                    :text="row['commit'] ? row['commit'].substring(0, 6) : ''"
                  >
                  </RepoURL>
                </el-row>
                <el-row style="margin-top: 10px" align="middle">
                  <svg-icon icon-class="publishTime" />
                  <span style="margin: 0 5px; font-size: 14px">
                    {{ row.updateTime }}
                  </span>
                </el-row>
                <el-row style="margin-top: 10px" justify="space-between">
                  <div>
                    <Button
                      v-if="row.deployState === DeployState.Deploying"
                      :permissions="[pms.DeployResetState]"
                      type="primary"
                      size="small"
                      @click="resetState(row)"
                    >
                      {{ $t('deployPage.resetState') }}
                    </Button>
                    <Button
                      v-else
                      :permissions="[pms.DeployProject]"
                      type="primary"
                      size="small"
                      @click="publish(row)"
                    >
                      {{ $t('deploy') }}
                    </Button>
                    <el-dropdown
                      trigger="click"
                      style="margin-left: 5px"
                      @command="(funcName) => commandFunc[funcName](row)"
                    >
                      <el-button size="small" type="warning">
                        {{ $t('func') }}
                        <el-icon class="el-icon--right">
                          <arrow-down />
                        </el-icon>
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
                  </div>
                  <el-tag :type="row.tagType" effect="plain">
                    {{ row.tagText }}
                  </el-tag>
                </el-row>
              </div>
            </el-card>
          </el-col>
        </el-row>
      </el-scrollbar>
    </el-row>
    <TheDetailDialog
      v-model="dialogVisible"
      :project-row="selectedItem"
      :on-rebuilt="handleRebuilt"
    />
    <TheCommitListDialog
      v-model="commitDialogVisible"
      :project-row="selectedItem"
      @cancel="handleCancelSelectCommit"
    >
      <template #tableOP="scope">
        <Button
          type="primary"
          :permissions="[pms.DeployProject]"
          @click="selectCommit(scope.row)"
        >
          {{ $t('select') }}
        </Button>
      </template>
    </TheCommitListDialog>
    <TheTaskListDialog
      v-model="taskListDialogVisible"
      :project-row="selectedItem"
    />
    <el-dialog
      v-model="publishDialogVisible"
      :title="publishFormProps.title"
      class="publish-dialog"
      align-center
      width="450px"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <template #header>
        <el-row style="white-space: pre-wrap" v-html="publishFormProps.title" />
      </template>
      <el-row>
        <el-row v-if="selectedItem.deployState === DeployState.Uninitialized">
          {{ $t('initial') }}
        </el-row>
        <el-row
          v-for="(variable, index) in publishFormProps.customVariables"
          :key="index"
          type="flex"
          align="middle"
          style="width: 100%; margin-bottom: 10px"
        >
          <el-row style="margin-right: 10px">
            <span>${</span>{{ variable.name }}<span>}</span>
          </el-row>
          <el-select
            v-if="variable.type == 'list'"
            v-model="variable.value"
            style="flex: 1"
          >
            <el-option
              v-for="item in selectedItem.script.customVariables[index].value.split(',')"
              :key="item"
              :label="item"
              :value="item"
            />
          </el-select>
          <el-input
            v-else
            v-model.trim="variable.value"
            style="flex: 1"
            autocomplete="off"
            placeholder="variable value"
          >
          </el-input>
        </el-row>
        <el-row
          v-if="selectedItem.deployState !== DeployState.Uninitialized"
          style="width: 100%"
        >
          <el-checkbox
            v-model="publishFormProps.selectCommit"
            :label="`Select Commit (default lastest ${selectedItem.branch})`"
            @change="handleSelectCommit"
          />
          <el-row
            v-show="publishFormProps.selectCommit"
            style="width: 100%; padding-left: 20px"
          >
            {{ publishFormProps.branch }}
            {{ publishFormProps.commit }}
          </el-row>
        </el-row>
        <el-row
          v-if="selectedItem.deployState !== DeployState.Uninitialized"
          style="width: 100%; margin-top: 10px"
        >
          <el-checkbox
            v-model="publishFormProps.selectServer"
            :label="`Select Server (default all)`"
            @change="handleSelectServer"
          />
          <el-row
            v-show="publishFormProps.selectServer"
            style="width: 100%; padding-left: 20px"
          >
            <el-checkbox-group v-model="publishFormProps.serverIds">
              <el-checkbox
                v-for="(item, index) in publishFormProps.serverOption"
                :key="index"
                :label="item.serverId"
              >
                {{ item.server.name + '(' + item.server.description + ')' }}
              </el-checkbox>
            </el-checkbox-group>
          </el-row>
        </el-row>
      </el-row>
      <template #footer>
        <el-button @click="publishDialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
        <el-button
          :disabled="publishFormProps.disabled"
          type="primary"
          @click="publishFormProps.func"
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
import { Link, More, ArrowDown } from '@element-plus/icons-vue'
import {
  DeployState,
  DeployList,
  DeployPublish,
  DeployResetState,
} from '@/api/deploy'
import { ProjectServerList, ProjectData, LabelList } from '@/api/project'
import RepoURL from '@/components/RepoURL/index.vue'
import { isLink, parseTime, deepClone } from '@/utils'
import TheDetailDialog from './TheDetailDialog.vue'
import TheCommitListDialog from './TheCommitListDialog.vue'
import TheTaskListDialog from './TheTaskListDialog.vue'
import TheReviewListDialog from './TheReviewListDialog.vue'
import TheProcessManagerDialog from './TheProcessManagerDialog.vue'
import TheFileCompareDialog from './TheFileCompareDialog.vue'
import TheFileSyncDialog from './TheFileSyncDialog.vue'
import { computed, watch, h, ref } from 'vue'
import { CommitData } from '@/api/repository'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const store = useStore()
const commitDialogVisible = ref(false)
const taskListDialogVisible = ref(false)
const fileSyncDialogVisible = ref(false)
const fileCompareDialogVisible = ref(false)
const processManagerDialogVisible = ref(false)
const reviewListDialogVisible = ref(false)
const publishDialogVisible = ref(false)
const dialogVisible = ref(false)
const stickList = ref(getStick())
const searchProject = ref({
  sort: getSort(),
  name: '',
  environment: '',
  label: [] as string[],
  pin: '',
})
const envTitleStyle = [
  {
    color: 'var(--el-color-info)',
  },
  {
    color: 'var(--el-color-danger)',
  },
  {
    color: 'var(--el-color-info)',
  },
  {
    color: 'var(--el-color-warning)',
  },
  {
    color: 'var(--el-color-info)',
  },
]
const selectedItem = ref({} as ProjectData)
const noMore = computed(
  () => tablePage.value.total === tablePage.value.list.length
)
const disabled = computed(() => noMore.value)
const tableData = ref<any[]>([])
const labelList = ref<string[]>([])
const pagination = ref({ page: 0, rows: 24 })

const tempPublishFormProps = {
  disabled: false,
  title: '',
  customVariables: [],
  commit: '',
  branch: '',
  serverIds: [],
  selectCommit: false,
  selectServer: false,
  serverOption: [] as ProjectServerList['datagram']['list'],
  func: () => {
    ElMessage.error('Undefined function')
  },
}
const publishFormProps = ref(tempPublishFormProps)

const tablePage = computed(() => {
  let _tableData = tableData.value
  if (searchProject.value.name !== '') {
    _tableData = _tableData.filter(
      (item) => item.name.indexOf(searchProject.value.name) !== -1
    )
  }
  if (searchProject.value.environment !== '') {
    _tableData = _tableData.filter(
      (item) => item.environment === Number(searchProject.value.environment)
    )
  }
  if (searchProject.value.pin !== '') {
    _tableData = _tableData.filter(
      (item) => item.pin === searchProject.value.pin
    )
  }
  if (searchProject.value.label.length > 0) {
    _tableData = _tableData.filter((item) =>
      item.label
        .split(',')
        .find((p: string) => searchProject.value.label.indexOf(p) > -1)
    )
  }
  return {
    list: _tableData.slice(0, pagination.value.page * pagination.value.rows),
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
      tableData.value[projectIndex].tagText = message
      tableData.value[projectIndex].deployState = data['state']
      if (data['state'] === 2) {
        tableData.value[projectIndex].tagType = 'success'
      } else if (data['state'] === 3) {
        tableData.value[projectIndex].tagType = 'danger'
        tableData.value[projectIndex].tagText = 'Fail'
      } else {
        tableData.value[projectIndex].tagType = 'warning'
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
getLabelList()

function getList() {
  new DeployList().request().then((response) => {
    tableData.value = response.data.list.map((item) => {
      let element: any = item
      element.tagType = 'info'
      element.tagText = 'Not deploy'
      if (element.deployState === 2) {
        element.tagType = 'success'
        element.tagText = 'Success'
      } else if (element.deployState === 1) {
        element.tagType = 'warning'
        element.tagText = 'Deploying'
      } else if (element.deployState === 3) {
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
}

function getLabelList() {
  new LabelList().request().then((response) => {
    labelList.value = response.data.list
  })
}

function stickChange() {
  tableData.value = tableData.value.map((_) => {
    _.pin = false
    return _
  })
  for (const stickId of [...stickList.value].reverse()) {
    const moveIndex = tableData.value.findIndex((_) => _.id == stickId)
    if (moveIndex > -1) {
      const moveItem = tableData.value.splice(moveIndex, 1)
      moveItem[0].pin = true
      tableData.value = moveItem.concat(tableData.value)
    }
  }
}

function sortChange(sort: string) {
  setSort(sort)
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
  // custom stick
  stickChange()
}

function load() {
  pagination.value.page++
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

function handleSelectCommit(state: boolean) {
  if (state == false) {
    publishFormProps.value.branch = ''
    publishFormProps.value.commit = ''
    return
  }
  commitDialogVisible.value = true
}

function handleSelectServer() {
  new ProjectServerList({ id: selectedItem.value.id })
    .request()
    .then((response) => {
      publishFormProps.value.serverOption = response.data.list
    })
}

function handleCancelSelectCommit() {
  publishFormProps.value.branch = ''
  publishFormProps.value.commit = ''
  publishFormProps.value.selectCommit = false
}

const cardMoreFunc: { [K: string]: (data: ProjectData) => void } = {
  handlePinCard,
  handleUnpinCard,
}

function handlePinCard(data: ProjectData) {
  let tmp = stickList.value
  tmp = tmp.filter((id) => id != data.id)
  tmp.unshift(data.id)
  stickList.value = tmp
  setStick(JSON.stringify(stickList.value))
  stickChange()
}

function handleUnpinCard(data: ProjectData) {
  stickList.value = stickList.value.filter((id) => id != data.id)
  setStick(JSON.stringify(stickList.value))
  stickChange()
}

const commandFunc: { [K: string]: (data: ProjectData) => void } = {
  handleTaskCommand,
  handleFileCompareCommand,
  handleFileSyncCommand,
  handleProcessManagerCommand,
  handleReviewCommand,
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

function selectCommit(data: CommitData) {
  publishFormProps.value.branch = data.branch
  publishFormProps.value.commit = data.commit
  commitDialogVisible.value = false
}

function publish(data: ProjectData) {
  restorePublishForm()
  let color = ''
  if (data.environment === 1) {
    color = 'color: var(--el-color-danger)'
  } else if (data.environment === 3) {
    color = 'color: var(--el-color-warning)'
  } else {
    color = 'color: var(--el-color-info)'
  }
  selectedItem.value = deepClone(data)
  const customVariables = deepClone(data.script.customVariables)
  publishFormProps.value.customVariables =
    customVariables &&
    customVariables.map((item: any) => {
      if (item.type == 'list') {
        item.value = ''
      }
      return item
    })
  publishDialogVisible.value = true

  let env = t(`envOption[${data.environment}]`)
  publishFormProps.value.title = `<span style="${color}; font-weight: 600; font-size: 15px">${data.name} - ${env}</span>`
  publishFormProps.value.func = () => {
    publishFormProps.value.disabled = true
    new DeployPublish({
      projectId: data.id,
      commit: publishFormProps.value.commit,
      branch: publishFormProps.value.branch,
      customVariables: publishFormProps.value.customVariables,
      serverIds: publishFormProps.value.serverIds,
    })
      .request()
      .then(() => {
        publishDialogVisible.value = false
        const projectIndex = tableData.value.findIndex(
          (element) => element.id === data.id
        )
        tableData.value[projectIndex].deployState = 1
      })
      .finally(() => {
        publishFormProps.value.disabled = false
      })
  }
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

function enterToBR(detail: string) {
  return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
}

function getSort(): string {
  const sortStr = localStorage.getItem('deploy-sort')
  if (sortStr) {
    return sortStr
  }
  return 'idDesc'
}

function setSort(value: string) {
  localStorage.setItem('deploy-sort', value)
}

function getStick(): number[] {
  const stickStr = localStorage.getItem('deploy-stick')
  if (stickStr) {
    return JSON.parse(stickStr)
  }
  return []
}

function setStick(value: string) {
  localStorage.setItem('deploy-stick', value)
}

function restorePublishForm() {
  publishFormProps.value = { ...tempPublishFormProps }
  console.log(publishFormProps.value)
}
</script>
<style lang="scss">
.publish-dialog {
  .el-dialog__header {
    padding: 20px 15px 10px 15px;
    margin-right: 0px;
  }
  .el-dialog__body {
    padding-top: 10px;
    padding-left: 15px;
    padding-bottom: 10px;
  }
  .el-dialog__footer {
    padding: 10px 15px;
  }
}
</style>
<style lang="scss" scoped>
.card {
  border: none;
  padding: 3px;
  position: relative;
  overflow: hidden;
  &:before {
    content: '';
    position: absolute;
    width: 1000px;
    height: 1000px;
  }
  .card-title__text {
    overflow: hidden;
    text-overflow: ellipsis;
    font-size: 14px;
    font-weight: 600;
    white-space: nowrap;
    flex: 1;
  }
  .card-title__link {
    overflow: hidden;
    text-overflow: ellipsis;
    font-weight: 600;
    white-space: nowrap;
  }
}

.deploying:before {
  animation: roll linear 3s infinite;
  background: conic-gradient(var(--el-color-warning), var(--el-bg-color));
  left: 50%;
  top: 50%;
}

.success:before {
  animation: opacity 3s 1;
  animation-fill-mode: forwards;
  background: var(--el-color-success);
  left: -50%;
  top: -50%;
}

.fail:before {
  animation: opacity linear 5s 1;
  animation-fill-mode: forwards;
  background: var(--el-color-danger);
  left: -50%;
  top: -50%;
}

@keyframes roll {
  from {
    transform: translate(-50%, -50%) rotateZ(0deg);
  }
  to {
    transform: translate(-50%, -50%) rotateZ(-360deg);
  }
}

@keyframes opacity {
  from {
    opacity: 1;
  }
  to {
    opacity: 0;
  }
}
</style>
