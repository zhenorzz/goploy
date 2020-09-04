<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex">
      <el-input v-model="projectName" style="width:300px" placeholder="Filter the project name" @change="getList" />
    </el-row>
    <el-table
      :key="tableHeight"
      border
      stripe
      highlight-current-row
      :max-height="tableHeight"
      :data="tablePageData"
      style="width: 100%;margin-top: 5px;"
    >
      <el-table-column prop="id" label="ID" width="80" align="center" />
      <el-table-column prop="name" :label="$t('name')" min-width="160" align="center">
        <template slot-scope="scope">
          <b v-if="scope.row.environment === 1" style="color: #F56C6C">{{ scope.row.name }} - {{ $t(`envOption[${scope.row.environment}]`) }}</b>
          <b v-else-if="scope.row.environment === 3" style="color: #E6A23C">{{ scope.row.name }} - {{ $t(`envOption[${scope.row.environment}]`) }}</b>
          <b v-else style="color: #909399">{{ scope.row.name }} - {{ $t(`envOption[${scope.row.environment}]`) }}</b>
        </template>
      </el-table-column>
      <el-table-column prop="branch" :label="$t('branch')" align="center">
        <template slot-scope="scope">
          <el-link
            style="font-size: 12px"
            :underline="false"
            :href="parseGitURL(scope.row.url) + '/tree/' + scope.row.branch"
            target="_blank"
          >
            {{ scope.row.branch }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column prop="commit" label="CommitID" width="150" align="center">
        <template slot-scope="scope">
          <el-tooltip effect="dark" :content="scope.row['commit']" placement="top">
            <el-link
              type="primary"
              style="font-size: 12px"
              :underline="false"
              :href="parseGitURL(scope.row['url']) + '/commit/' + scope.row['commit']"
              target="_blank"
            >
              {{ scope.row['commit'] ? scope.row['commit'].substring(0, 6) : '' }}
            </el-link>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column prop="deployState" :label="$t('state')" width="230" align="center">
        <template slot-scope="scope">
          <el-tag :type="scope.row.tagType" effect="plain">{{ scope.row.tagText }}</el-tag>
          <el-progress :percentage="scope.row.progressPercentage" :status="scope.row.progressStatus" />
        </template>
      </el-table-column>
      <el-table-column prop="updateTime" :label="$t('time')" width="160" align="center" />
      <el-table-column prop="operation" :label="$t('op')" width="310" fixed="right" align="center">
        <template slot-scope="scope">
          <el-row class="operation-btn">
            <el-dropdown
              split-button
              trigger="click"
              :disabled="scope.row.deployState === 1"
              type="primary"
              @click="publish(scope.row)"
              @command="handlePublishCommand"
            >
              {{ $t('deploy') }}
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item :command="scope.row">List commit</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
            <el-dropdown
              v-if="hasGroupManagerPermission()"
              split-button
              trigger="click"
              :disabled="scope.row.deployState === 1"
              type="warning"
              @click="handleAddProjectTask(scope.row)"
              @command="handleProjectTaskCommand"
            >
              {{ $t('crontab') }}
              <el-dropdown-menu slot="dropdown" style="min-width:84px;text-align:center;">
                <el-dropdown-item :command="scope.row">{{ $t('manage') }}</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
            <el-button type="success" @click="handleDetail(scope.row)">{{ $t('detail') }}</el-button>
          </el-row>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      hide-on-single-page
      :total="pagination.total"
      :page-size="pagination.rows"
      :current-page.sync="pagination.page"
      style="margin-top:10px; text-align:right;"
      background
      :page-sizes="[20, 50, 100]"
      layout="sizes, total, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handlePageChange"
    />
    <el-dialog :title="$t('detail')" :visible.sync="dialogVisible" class="publish-record">
      <el-row type="flex">
        <el-row v-loading="searchPreview.loading" class="publish-preview">
          <el-row>
            <el-select v-model="searchPreview.userId" :placeholder="$t('user')" style="width: 170px;" clearable>
              <el-option
                v-for="(item, index) in userOption"
                :key="index"
                :label="item.userName"
                :value="item.userId"
              />
            </el-select>
            <el-select v-model="searchPreview.state" :placeholder="$t('state')" style="width: 95px;" clearable>
              <el-option :label="$t('success')" :value="1" />
              <el-option :label="$t('fail')" :value="0" />
            </el-select>
            <el-button type="primary" icon="el-icon-search" @click="searchPreviewList" />
          </el-row>
          <el-radio-group v-model="publishToken" @change="handleDetailChange">
            <el-row v-for="(item, index) in gitTraceList" :key="index">
              <el-row style="margin:5px 0">
                <el-radio class="publish-commit" :label="item.token" border>
                  <span class="publish-name">{{ item.publisherName }}</span> <span class="publish-commitID">commitID: {{ item.commit }}</span>
                  <i v-if="item.publishState === 1" class="el-icon-check" style="color:#67C23A;float:right;font-size:14px;font-weight:900;" />
                  <i v-else class="el-icon-close" style="color:#F56C6C;float:right;font-size:14px;font-weight:900;" />
                  <!-- <span v-if="item.publishState === 1" style="color:#67C23A;float:right;">{{ $t('success') }}</span>
                  <span v-else style="color:#F56C6C;float:right;">{{ $t('fail') }}</span> -->
                </el-radio>
                <el-button type="danger" plain @click="rollback(item)">rebuild</el-button>
              </el-row>
            </el-row>
          </el-radio-group>
          <el-pagination
            :total="previewPagination.total"
            :page-size="previewPagination.rows"
            :current-page.sync="previewPagination.page"
            style="text-align:right;margin-right:20px"
            layout="total, prev, next"
            @current-change="handlePreviewPageChange"
          />
        </el-row>
        <el-row class="project-detail" style="flex:1;width:100%">
          <el-row v-for="(item, index) in publishLocalTraceList" :key="index">
            <el-row v-if="item.type === 2">
              <el-row style="margin:5px 0">---------------GIT----------------</el-row>
              <el-row style="margin:5px 0">Time: {{ item.insertTime }}</el-row>
              <!-- 用数组的形式 兼容以前版本 -->
              <el-row v-if="item.state !== 0">
                <el-row>Commit:
                  <el-link
                    type="primary"
                    :underline="false"
                    :href="parseGitURL(searchPreview.url) + '/commit/' + item['commit']"
                    target="_blank"
                  >
                    {{ item['commit'] }}
                  </el-link>
                </el-row>
                <el-row>Message: {{ item['message'] }}</el-row>
                <el-row>Author: {{ item['author'] }}</el-row>
                <el-row>Datetime: {{ item['timestamp'] ? parseTime(item['timestamp']) : '' }}</el-row>
                <el-row><span v-html="formatDetail(item['diff'])" /></el-row>
              </el-row>
              <el-row v-else style="margin:5px 0">
                <el-tag type="danger" effect="plain">{{ $t('fail') }}</el-tag>
                <span v-html="formatDetail(item.detail)" />
              </el-row>
            </el-row>
            <el-row v-if="item.type === 3">
              <hr>
              <el-row style="margin:5px 0">------------After pull------------</el-row>
              <el-row style="margin:5px 0">Time: {{ item.insertTime }}</el-row>
              <el-row>Script: <pre v-html="formatDetail(item.script)" /></el-row>
              <el-row style="margin:5px 0">
                <el-tag v-if="item.state === 0" type="danger" effect="plain">{{ $t('fail') }}</el-tag>
                <el-row v-if="item.detail.length > 0">[goploy ~]# <span v-html="formatDetail(item.detail)" /></el-row>
              </el-row>
            </el-row>
          </el-row>
          <el-tabs v-model="activeRomoteTracePane">
            <el-tab-pane v-for="(item, serverName) in publishRemoteTraceList" :key="serverName" :label="serverName" :name="serverName">
              <el-row v-for="(trace, key) in item" :key="key">
                <el-row v-if="trace.type === 4">
                  <el-row style="margin:5px 0">----------Before deploy--------</el-row>
                  <el-row style="margin:5px 0">Time: {{ trace.insertTime }}</el-row>
                  <el-row>Script: <pre v-html="formatDetail(trace.script)" /></el-row>
                  <el-row style="margin:5px 0">
                    <el-tag v-if="trace.state === 0" type="danger" effect="plain">{{ $t('fail') }}</el-tag>
                    <el-row>[goploy ~]# <span v-html="formatDetail(trace.detail)" /></el-row>
                  </el-row>
                </el-row>
                <el-row v-else-if="trace.type === 5">
                  <el-row style="margin:5px 0">--------------Rsync-------------</el-row>
                  <el-row style="margin:5px 0">Time: {{ trace.insertTime }}</el-row>
                  <el-row>Command: {{ trace.command }}</el-row>
                  <el-row style="margin:5px 0">
                    <el-tag v-if="trace.state === 0" type="danger" effect="plain">{{ $t('fail') }}</el-tag>
                    <span v-html="formatDetail(trace.detail)" />
                  </el-row>
                </el-row>
                <el-row v-else>
                  <el-row style="margin:5px 0">----------After deploy---------</el-row>
                  <el-row style="margin:5px 0">Time: {{ trace.insertTime }}</el-row>
                  <el-row>Script: {{ trace.script }}</el-row>
                  <el-row style="margin:5px 0">
                    <el-tag v-if="trace.state === 0" type="danger" effect="plain">{{ $t('fail') }}</el-tag>
                    <el-row v-if="trace.detail.length > 0">[goploy ~]# <span v-html="formatDetail(trace.detail)" /></el-row>
                  </el-row>
                </el-row>
              </el-row>
            </el-tab-pane>
          </el-tabs>
        </el-row>
      </el-row>
    </el-dialog>
    <el-dialog title="commit" :visible.sync="commitDialogVisible">
      <el-table
        v-loading="commitTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="commitTableData"
      >
        <el-table-column type="expand">
          <template slot-scope="props">
            <span v-html="formatDetail(props.row.diff)" />
          </template>
        </el-table-column>
        <el-table-column prop="commit" label="commit" width="290">
          <template slot-scope="scope">
            <el-link
              type="primary"
              style="font-size: 12px"
              :underline="false"
              :href="parseGitURL(scope.row.url) + '/commit/' + scope.row.commit"
              target="_blank"
            >
              {{ scope.row.commit }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="author" label="author" />
        <el-table-column prop="message" label="message" width="200" show-overflow-tooltip />
        <el-table-column label="time" width="135" align="center">
          <template slot-scope="scope">
            {{ parseTime(scope.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="operation" :label="$t('op')" width="80" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button type="danger" @click="rollback(scope.row)">{{ $t('deploy') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="commitDialogVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('manage')" :visible.sync="taskListDialogVisible">
      <el-table
        v-loading="taskTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="taskTableData"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="projectName" :label="$t('projectName')" width="150" />
        <el-table-column prop="commitId" label="commit" width="290" />
        <el-table-column prop="date" :label="$t('date')" width="150" />
        <el-table-column prop="isRun" :label="$t('task')" width="60">
          <template slot-scope="scope">
            {{ $t(`runOption[${scope.row.isRun}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="state" :label="$t('state')" width="50">
          <template slot-scope="scope">
            {{ $t(`stateOption[${scope.row.state}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="creator" :label="$t('creator')" />
        <el-table-column prop="editor" :label="$t('editor')" />
        <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
        <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
        <el-table-column prop="operation" :label="$t('op')" width="150" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button type="primary" :disabled="scope.row.isRun === 1 || scope.row.state === 0" @click="handleEditProjectTask(scope.row)">{{ $t('edit') }}</el-button>
            <el-button type="danger" :disabled="scope.row.isRun === 1 || scope.row.state === 0" @click="removeProjectTask(scope.row)">{{ $t('delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="taskListDialogVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('setting')" :visible.sync="taskDialogVisible" width="600px">
      <el-form ref="taskForm" :rules="taskFormRules" :model="taskFormData" label-width="120px">
        <el-form-item :label="$t('projectName')">
          <span>{{ taskFormProps.projectName }}</span>
        </el-form-item>
        <el-form-item label="commitId" prop="commitId">
          <el-select v-model="taskFormData.commitId" placeholder="CommitID" style="width: 400px">
            <el-option
              v-for="(item, index) in taskFormProps.commitOptions"
              :key="index"
              :label="item.commit+'('+item.author+')'"
              :value="item.commit"
            />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('time')" prop="date">
          <el-date-picker
            v-model="taskFormData.date"
            :picker-options="taskFormProps.pickerOptions"
            type="datetime"
            value-format="yyyy-MM-dd HH:mm:ss"
            style="width: 400px"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="taskDialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="taskFormProps.disabled" type="primary" @click="submitTask">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import tableHeight from '@/mixin/tableHeight'
import { getList, getDetail, getPreview, getCommitList, publish } from '@/api/deploy'
import { addTask, editTask, removeTask, getTaskList } from '@/api/project'
import { getUserOption } from '@/api/namespace'
import { parseTime, parseGitURL } from '@/utils'

export default {
  name: 'Deploy',
  mixins: [tableHeight],
  data() {
    return {
      userId: '',
      userOption: [],
      projectName: '',
      publishToken: '',
      commitDialogVisible: false,
      taskDialogVisible: false,
      taskListDialogVisible: false,
      dialogVisible: false,
      tableData: [],
      pagination: {
        total: 0,
        page: 1,
        rows: 20
      },
      taskTableLoading: false,
      taskTableData: [],
      taskPagination: {
        total: 0,
        page: 1,
        rows: 20
      },
      taskFormProps: {
        projectName: '',
        disabled: false,
        commitOptions: [],
        pickerOptions: {
          disabledDate(time) {
            return time.getTime() < Date.now() - 3600 * 1000 * 24
          }
        }
      },
      taskFormData: {
        id: 0,
        projectId: '',
        commitId: '',
        date: ''
      },
      taskFormRules: {
        commitId: [
          { required: true, message: 'CommitID required', trigger: 'change' }
        ],
        date: [
          { required: true, message: 'Date required', trigger: 'change' }
        ]
      },
      searchPreview: {
        loading: false,
        projectId: '',
        userId: '',
        url: '',
        state: ''
      },
      gitTraceList: [],
      previewPagination: {
        page: 1,
        rows: 11,
        total: 0
      },
      commitTableLoading: false,
      commitTableData: [],
      publishTraceList: [],
      publishLocalTraceList: [],
      publishRemoteTraceList: {},
      activeRomoteTracePane: ''
    }
  },
  computed: {
    tablePageData: function() {
      return this.tableData.slice((this.pagination.page - 1) * this.pagination.rows, this.pagination.page * this.pagination.rows)
    }
  },
  watch: {
    '$store.getters.ws_message': function(response) {
      if (response.type !== 1) {
        return
      }
      const data = response.message
      data.message = this.formatDetail(data.message)
      if (data.state === 0) {
        this.$notify.error({
          title: data.projectName,
          dangerouslyUseHTMLString: true,
          message: data.message,
          duration: 0
        })
      }
      const projectIndex = this.tableData.findIndex(element => element.id === data.projectId)
      if (projectIndex !== -1) {
        const percent = 12.5 * data.state
        this.tableData[projectIndex].progressPercentage = percent
        this.tableData[projectIndex].progressStatus = 'warning'
        this.tableData[projectIndex].tagType = 'warning'
        this.tableData[projectIndex].tagText = data.message

        if (percent === 0) {
          this.tableData[projectIndex].progressStatus = 'exception'
          this.tableData[projectIndex].tagType = 'danger'
          this.tableData[projectIndex].tagText = 'Fail'
        } else if (percent > 98) {
          this.tableData[projectIndex].progressStatus = 'success'
          this.tableData[projectIndex].tagType = 'success'
        }

        if (data['ext']) {
          Object.assign(this.tableData[projectIndex], data['ext'])
        }
        this.tableData[projectIndex].deployState = data.state
        this.tableData[projectIndex].publisherName = data.username
        this.tableData[projectIndex].updateTime = parseTime(new Date())
      }
    }
  },
  created() {
    this.getList()
    this.getUserOption()
  },
  methods: {
    parseTime,
    parseGitURL,
    getUserOption() {
      getUserOption().then((response) => {
        this.userOption = response.data.list
      })
    },

    getList() {
      getList(this.projectName).then((response) => {
        this.tableData = response.data.list.map(element => {
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
          }
          try {
            Object.assign(element, JSON.parse(element.publishExt))
          } catch (error) {
            console.log('Project not deploy')
          }
          return element
        })
        this.pagination.total = this.tableData.length
      })
    },

    handleSizeChange(val) {
      this.pagination.rows = val
      this.handlePageChange(1)
    },

    handlePageChange(page) {
      this.pagination.page = page
      this.getList()
    },

    publish(data) {
      const id = data.id
      const h = this.$createElement
      let color = ''
      if (data.environment === 1) {
        color = 'color: #F56C6C'
      } else if (data.environment === 3) {
        color = 'color: #E6A23C'
      } else {
        color = 'color: #909399'
      }
      this.$confirm('', this.$i18n.t('tips'), {
        message: h('p', null, [
          h('span', null, 'Deploy Project: '),
          h('b', { style: color }, data.name + ' - ' + this.$i18n.t(`envOption[${data.environment}]`))
        ]),
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        this.gitLog = []
        this.remoteLog = {}
        publish(id, '').then((response) => {
          const projectIndex = this.tableData.findIndex(element => element.id === id)
          this.tableData[projectIndex].deployState = 1
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    getDetail() {
      getDetail(this.publishToken).then((response) => {
        const publishTraceList = response.data.publishTraceList || []
        this.publishTraceList = publishTraceList.map(element => {
          if (element.ext !== '') Object.assign(element, JSON.parse(element.ext))
          return element
        })

        this.publishLocalTraceList = this.publishTraceList.filter(element => element.type < 4)
        this.publishRemoteTraceList = {}
        for (const trace of this.publishTraceList) {
          if (trace.type < 4) continue
          if (!this.publishRemoteTraceList[trace.serverName]) {
            this.publishRemoteTraceList[trace.serverName] = []
          }
          this.publishRemoteTraceList[trace.serverName].push(trace)
        }
        this.activeRomoteTracePane = Object.keys(this.publishRemoteTraceList)[0]
      })
    },

    getPreviewList() {
      this.searchPreview.loading = true
      getPreview(this.previewPagination, {
        projectId: this.searchPreview.projectId,
        userId: this.searchPreview.userId || 0,
        state: this.searchPreview.state === '' ? -1 : this.searchPreview.state
      }).then((response) => {
        const gitTraceList = response.data.gitTraceList || []
        this.gitTraceList = gitTraceList.map(element => {
          if (element.ext !== '') Object.assign(element, JSON.parse(element.ext))
          element.commit = element['commit'] ? element['commit'].substring(0, 6) : ''
          return element
        })
        if (this.gitTraceList.length > 0) {
          this.publishToken = this.gitTraceList[0].token
          this.getDetail()
        }
        this.previewPagination.total = response.data.pagination.total
      }).finally(() => {
        this.searchPreview.loading = false
      })
    },

    searchPreviewList() {
      this.handlePreviewPageChange(1)
    },

    handleDetail(data) {
      this.dialogVisible = true
      this.searchPreview.projectId = data.id
      this.searchPreview.url = data.url
      this.searchPreview.userId = ''
      this.searchPreview.state = ''
      this.getPreviewList()
    },

    handlePreviewPageChange(page) {
      this.previewPagination.page = page
      this.getPreviewList()
    },

    handleDetailChange(lastPublishToken) {
      this.publishToken = lastPublishToken
      this.getDetail()
    },

    handlePublishCommand(data) {
      const id = data.id
      this.commitDialogVisible = true
      this.commitTableLoading = true
      getCommitList(id).then(response => {
        this.commitTableData = response.data.commitList.map(element => {
          return Object.assign(element, {
            projectId: id,
            url: data.url
          })
        })
      }).finally(() => {
        this.commitTableLoading = false
      })
    },

    handleAddProjectTask(data) {
      this.taskDialogVisible = true
      this.taskFormData.id = 0
      if (this.taskFormData.projectId !== data.id) {
        this.taskFormData.projectId = data.id
        this.taskFormProps.projectName = data.name
        this.taskFormData.commitId = ''
        this.taskFormData.date = ''
        const id = data.id
        getCommitList(id).then(response => {
          this.taskFormProps.commitOptions = response.data.commitList || []
        })
      }
    },

    handleEditProjectTask(data) {
      this.taskDialogVisible = true
      this.taskFormData.id = data.id
      this.taskFormData.commitId = data.commitId
      this.taskFormData.date = data.date
      if (this.taskFormData.projectId !== data.projectId) {
        this.taskFormProps.projectName = data.projectName
        getCommitList(data.projectId).then(response => {
          this.taskFormProps.commitOptions = response.data.commitList || []
        })
      }
    },

    handleProjectTaskCommand(data) {
      this.taskListDialogVisible = true
      this.taskTableLoading = true
      getTaskList(this.taskPagination, data.id).then(response => {
        const projectTaskList = response.data.projectTaskList || []
        this.taskTableData = projectTaskList.map(element => {
          return Object.assign(element, { projectId: data.id, projectName: data.name })
        })
        this.taskPagination.total = response.data.pagination.total
      }).finally(() => { this.taskTableLoading = false })
    },

    submitTask() {
      this.$refs.taskForm.validate((valid) => {
        if (valid) {
          this.taskFormProps.disabled = true
          if (this.taskFormData.id === 0) {
            addTask(this.taskFormData).then(response => {
              this.$message.success('Success')
            }).finally(() => {
              this.taskFormProps.disabled = false
              this.taskDialogVisible = false
            })
          } else {
            editTask(this.taskFormData).then(response => {
              this.$message.success('Success')
              const projectTaskIndex = this.taskTableData.findIndex(element => element.id === this.taskFormData.id)
              this.taskTableData[projectTaskIndex]['commitId'] = this.taskFormData.commitId
              this.taskTableData[projectTaskIndex]['date'] = this.taskFormData.date
              this.taskTableData[projectTaskIndex]['editor'] = this.$store.getters.name
              this.taskTableData[projectTaskIndex]['editorId'] = this.$store.getters.uid
              this.taskTableData[projectTaskIndex]['updateTime'] = parseTime(new Date())
            }).finally(() => {
              this.taskFormProps.disabled = false
              this.taskDialogVisible = false
            })
          }
        } else {
          return false
        }
      })
    },

    removeProjectTask(data) {
      this.$confirm(this.$i18n.t('deployPage.removeProjectTaskTips', { projectName: data.projectName }), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        removeTask(data.id).then((response) => {
          const projectTaskIndex = this.taskTableData.findIndex(element => element.id === data.id)
          this.taskTableData[projectTaskIndex]['state'] = 0
          this.taskTableData[projectTaskIndex]['editor'] = this.$store.getters.name
          this.taskTableData[projectTaskIndex]['editorId'] = this.$store.getters.uid
          this.taskTableData[projectTaskIndex]['updateTime'] = parseTime(new Date())
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    rollback(data) {
      this.$confirm(this.$i18n.t('deployPage.rollbackTips', { commit: data.commit }), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        publish(data.projectId, data.commit).then((response) => {
          const projectIndex = this.tableData.findIndex(element => element.id === data.projectId)
          this.tableData[projectIndex].deployState = 1
          this.commitDialogVisible = false
          this.dialogVisible = false
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    formatDetail(detail) {
      return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
    }
  }
}
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
@import "@/styles/mixin.scss";

.publish {
  &-preview {
    width: 330px;
    margin-left: 10px;
  }
  &-commit{
    margin-right: 5px;
    padding-right:8px;
    width: 240px;
    line-height: 12px;
  }
  &-commitID{
    display: inline-block;
    vertical-align: top;
  }
  &-name {
    width: 60px;
    display: inline-block;
    text-align: center;
    overflow:hidden;
    vertical-align: top;
    text-overflow:ellipsis;
    white-space:nowrap;
  }
}

.project-detail {
  padding-left:5px;
  height:470px;
  overflow-y: auto;
  @include scrollBar();
}

.operation-btn {
  >>>.el-button {
    line-height: 1.15;
  }
}

@media screen and (max-width: 1440px){
  .publish-record {
    >>>.el-dialog {
      width: 75%;
    }
  }
}
</style>
