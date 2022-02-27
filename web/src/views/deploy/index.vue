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
              @command="(funcName) => callCommandFunc(funcName, scope.row)"
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
              @command="(funcName) => callCommandFunc(funcName, scope.row)"
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
      ref="test"
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
      <template #footer class="dialog-footer">
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
import tableHeight from '@/mixin/tableHeight'
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
import { ElMessageBox, ElMessage } from 'element-plus'
import Validator from 'async-validator'
import { h, defineComponent } from 'vue'
import { CommitData } from '@/api/repository'

export default defineComponent({
  name: 'DeployIndex',
  components: {
    TheDetailDialog,
    TheCommitListDialog,
    TheTagListDialog,
    TheTaskListDialog,
    TheReviewListDialog,
    TheFileCompareDialog,
    TheProcessManagerDialog,
  },
  mixins: [tableHeight],
  data() {
    return {
      role: getRole(),
      commitDialogVisible: false,
      tagDialogVisible: false,
      greyServerDialogVisible: false,
      taskListDialogVisible: false,
      fileCompareDialogVisible: false,
      processManagerDialogVisible: false,
      reviewDialogVisible: false,
      reviewListDialogVisible: false,
      dialogVisible: false,
      tableloading: false,
      tableDefaultSort: {} as { prop: string; order: string },
      selectedItem: {} as ProjectData['datagram'],
      tableData: [] as DeployList['datagram']['list'],
      searchProject: {
        name: '',
        environment: '',
        autoDeploy: '',
      },
      pagination: {
        page: 1,
        rows: 20,
      },
      greyServerFormProps: {
        disabled: false,
        serverOption: [] as ProjectServerList['datagram']['list'],
      },
      greyServerFormData: {
        projectId: 0,
        commit: '',
        serverIds: [],
      },
      greyServerFormRules: {
        serverIds: [
          {
            type: 'array',
            required: true,
            message: 'Server required',
            trigger: 'change',
          },
        ],
      },
    }
  },
  computed: {
    tablePageData(): DeployList['datagram']['list'] {
      let tableData = this.tableData
      if (this.searchProject.name !== '') {
        tableData = this.tableData.filter(
          (item) => item.name.indexOf(this.searchProject.name) !== -1
        )
      }
      if (this.searchProject.environment !== '') {
        tableData = this.tableData.filter(
          (item) => item.environment === Number(this.searchProject.environment)
        )
      }
      if (this.searchProject.autoDeploy !== '') {
        tableData = this.tableData.filter(
          (item) => item.autoDeploy === Number(this.searchProject.autoDeploy)
        )
      }
      return tableData.slice(
        (this.pagination.page - 1) * this.pagination.rows,
        this.pagination.page * this.pagination.rows
      )
    },
  },
  watch: {
    '$store.state.websocket.message': function (response) {
      if (response.type !== 1) {
        return
      }
      const data = response.message
      const message = this.enterToBR(data.message)
      const projectIndex = this.tableData.findIndex(
        (element) => element.id === data.projectId
      )
      if (projectIndex !== -1) {
        const percent = 20 * data.state
        this.tableData[projectIndex].progressPercentage = percent
        this.tableData[projectIndex].progressStatus = 'warning'
        this.tableData[projectIndex].tagType = 'warning'
        this.tableData[projectIndex].tagText = message
        this.tableData[projectIndex].deployState = 1
        if (percent === 0) {
          this.tableData[projectIndex].progressStatus = 'exception'
          this.tableData[projectIndex].tagType = 'danger'
          this.tableData[projectIndex].tagText = 'Fail'
          this.tableData[projectIndex].deployState = 3
        } else if (percent === 100) {
          this.tableData[projectIndex].progressStatus = 'success'
          this.tableData[projectIndex].tagType = 'success'
          this.tableData[projectIndex].deployState = 2
        }

        if (data['ext']) {
          Object.assign(this.tableData[projectIndex], data['ext'])
        }
        this.tableData[projectIndex].publisherName = data.username
        this.tableData[projectIndex].updateTime = parseTime(
          new Date().getTime()
        )
      }
    },
  },
  created() {
    this.tableDefaultSort = this.getTableSort()
    this.getList()
  },
  mounted() {
    console.log(this.$refs['test'])
  },
  methods: {
    parseTime,
    parseGitURL,
    getList() {
      this.tableloading = true
      new DeployList()
        .request()
        .then((response) => {
          this.tableData = response.data.list.map((element) => {
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
          this.sortChange(this.tableDefaultSort)
        })
        .finally(() => {
          this.tableloading = false
        })
    },

    sortChange({ prop, order }: { prop: string; order: string }) {
      this.setTableSort(prop, order)
      if (!prop && !order) {
        prop = 'id'
        order = 'descending'
      }
      if (prop === 'name') {
        prop = 'environment'
      }
      this.tableData = this.tableData.sort(
        (
          row1: ProjectData['datagram'],
          row2: ProjectData['datagram']
        ): number => {
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
    },

    handleSizeChange(val = 1) {
      this.pagination.rows = val
      this.handlePageChange(1)
    },

    handlePageChange(page = 1) {
      this.pagination.page = page
    },

    handleDetail(data: ProjectData['datagram']) {
      this.selectedItem = data
      this.dialogVisible = true
    },

    handleRebuilt() {
      const projectIndex = this.tableData.findIndex(
        (element) => element.id === this.selectedItem.id
      )
      this.tableData[projectIndex].deployState = 1
    },

    handleGreyPublish(data: CommitData['datagram']) {
      new ProjectServerList({ id: this.selectedItem.id })
        .request()
        .then((response) => {
          this.greyServerFormProps.serverOption = response.data.list
        })
      // add projectID to server form
      this.greyServerFormData.projectId = this.selectedItem.id
      this.greyServerFormData.commit = data.commit
      this.greyServerDialogVisible = true
    },

    callCommandFunc(funcName: string, data: ProjectData['datagram']) {
      ;(this[funcName] as (data: ProjectData['datagram']) => void)(data)
    },

    handleCommitCommand(data: ProjectData['datagram']) {
      this.selectedItem = data
      this.commitDialogVisible = true
    },

    handleTagCommand(data: ProjectData['datagram']) {
      this.selectedItem = data
      this.tagDialogVisible = true
    },

    handleTaskCommand(data: ProjectData['datagram']) {
      this.selectedItem = data
      this.taskListDialogVisible = true
    },

    handleFileCompareCommand(data: ProjectData['datagram']) {
      this.selectedItem = data
      this.fileCompareDialogVisible = true
    },

    handleProcessManagerCommand(data: ProjectData['datagram']) {
      this.selectedItem = data
      this.processManagerDialogVisible = true
    },

    handleReviewCommand(data: ProjectData['datagram']) {
      this.selectedItem = data
      this.reviewListDialogVisible = true
    },

    publish(data: ProjectData['datagram']) {
      const id = data.id
      let color = ''
      if (data.environment === 1) {
        color = 'color: #F56C6C'
      } else if (data.environment === 3) {
        color = 'color: #E6A23C'
      } else {
        color = 'color: #909399'
      }
      ElMessageBox.confirm('', this.$t('tips'), {
        message: h('p', null, [
          h('span', null, 'Deploy Project: '),
          h(
            'b',
            { style: color },
            data.name + ' - ' + this.$t(`envOption[${data.environment}]`)
          ),
        ]),
        confirmButtonText: this.$t('confirm'),
        cancelButtonText: this.$t('cancel'),
        type: 'warning',
      })
        .then(() => {
          new DeployPublish({ projectId: id, commit: '', branch: '' })
            .request()
            .then(() => {
              const projectIndex = this.tableData.findIndex(
                (element) => element.id === id
              )
              this.tableData[projectIndex].deployState = 1
            })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },

    publishByCommit(data: CommitData['datagram']) {
      ElMessageBox.confirm(
        this.$t('deployPage.publishCommitTips', { commit: data.commit }),
        this.$t('tips'),
        {
          confirmButtonText: this.$t('confirm'),
          cancelButtonText: this.$t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          new DeployPublish({
            projectId: this.selectedItem.id,
            branch: data.branch,
            commit: data.commit,
          })
            .request()
            .then(() => {
              const projectIndex = this.tableData.findIndex(
                (element) => element.id === this.selectedItem.id
              )
              this.tableData[projectIndex].deployState = 1
              this.commitDialogVisible = false
              this.tagDialogVisible = false
            })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },

    resetState(data: ProjectData['datagram']) {
      ElMessageBox.confirm(
        this.$t('deployPage.resetStateTips'),
        this.$t('tips'),
        {
          confirmButtonText: this.$t('confirm'),
          cancelButtonText: this.$t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          new DeployResetState({ projectId: data.id }).request().then(() => {
            const projectIndex = this.tableData.findIndex(
              (element) => element.id === data.id
            )
            this.tableData[projectIndex].deployState = 0
            this.tableData[projectIndex].progressPercentage = 0
            this.tableData[projectIndex].tagType = 'info'
            this.tableData[projectIndex].tagText = 'Not deploy'
          })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },

    greyPublish() {
      ;(this.$refs.greyServerForm as Validator).validate((valid: boolean) => {
        if (valid) {
          const data = this.greyServerFormData
          ElMessageBox.confirm(
            this.$t('deployPage.publishCommitTips', {
              commit: data.commit,
            }),
            this.$t('tips'),
            {
              confirmButtonText: this.$t('confirm'),
              cancelButtonText: this.$t('cancel'),
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
                  const projectIndex = this.tableData.findIndex(
                    (element) => element.id === data.projectId
                  )
                  this.tableData[projectIndex].deployState = 1
                  this.commitDialogVisible = false
                  this.tagDialogVisible = false
                  this.greyServerDialogVisible = false
                })
            })
            .catch(() => {
              ElMessage.info('Cancel')
            })
        } else {
          return false
        }
      })
    },

    enterToBR(detail: string) {
      return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
    },

    getTableSort(): { prop: string; order: string } {
      const sortJsonStr = localStorage.getItem('deploy-table-sort')
      if (sortJsonStr) {
        return JSON.parse(sortJsonStr) as { prop: string; order: string }
      }
      return { prop: 'id', order: 'descending' }
    },

    setTableSort(prop: string, order: string) {
      localStorage.setItem('deploy-table-sort', JSON.stringify({ prop, order }))
    },
  },
})
</script>
