<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex">
      <el-input v-model="projectName" style="width:300px" placeholder="请输入项目名称" @change="getList" />
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
      <el-table-column prop="name" label="项目名称" min-width="160" align="center">
        <template slot-scope="scope">
          <b v-if="scope.row.environment === '生产环境'" style="color: #F56C6C">{{ scope.row.name }} - {{ scope.row.environment }}</b>
          <b v-else-if="scope.row.environment === '测试环境'" style="color: #E6A23C">{{ scope.row.name }} - {{ scope.row.environment }}</b>
          <b v-else style="color: #909399">{{ scope.row.name }} - {{ scope.row.environment }}</b>
        </template>
      </el-table-column>
      <el-table-column prop="branch" label="分支" align="center" />
      <el-table-column prop="commit" label="commitID" width="150" align="center">
        <template slot-scope="scope">
          <el-tooltip effect="dark" :content="scope.row['commit']" placement="top">
            <span>{{ scope.row['commit'] ? scope.row['commit'].substring(0, 6) : '' }}</span>
          </el-tooltip>
        </template>
      </el-table-column>
      <el-table-column prop="deployState" label="构建状态" width="230" align="center">
        <template slot-scope="scope">
          <el-tag :type="scope.row.tagType" effect="plain">{{ scope.row.tagText }}</el-tag>
          <el-progress :percentage="scope.row.progressPercentage" :status="scope.row.progressStatus" />
        </template>
      </el-table-column>
      <el-table-column prop="updateTime" label="上次构建时间" width="160" align="center" />
      <el-table-column prop="operation" label="操作" width="255" fixed="right" align="center">
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
              构建
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item :command="scope.row">选择具体commit</el-dropdown-item>
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
              定时
              <el-dropdown-menu slot="dropdown" style="min-width:84px;text-align:center;">
                <el-dropdown-item :command="scope.row">定时构建管理</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
            <el-button type="success" @click="handleDetail(scope.row)">详情</el-button>
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
    <el-dialog title="构建记录" :visible.sync="dialogVisible" class="publish-record">
      <el-row type="flex">
        <el-row v-loading="searchPreview.loading" class="publish-preview">
          <el-select v-model="searchPreview.userId" style="width:142px" placeholder="选择用户" clearable>
            <el-option
              v-for="(item, index) in userOption"
              :key="index"
              :label="item.name"
              :value="item.id"
            />
          </el-select>
          <el-select v-model="searchPreview.state" placeholder="状态" style="width:95px" clearable>
            <el-option label="成功" :value="1" />
            <el-option label="失败" :value="0" />
          </el-select>
          <el-button type="primary" icon="el-icon-search" @click="searchPreviewList">搜索</el-button>
          <el-radio-group v-model="publishToken" @change="handleDetailChange">
            <el-row v-for="(item, index) in gitTraceList" :key="index">
              <el-row style="margin:5px 0">
                <el-radio class="publish-commit" :label="item.token" border>
                  <span class="publish-name">{{ item.publisherName }}</span> <span class="publish-commitID">commitID: {{ item.commit }}</span>
                  <span v-if="item.publishState === 1" style="color:#67C23A;float:right;">成功</span>
                  <span v-else style="color:#F56C6C;float:right;">失败</span>
                </el-radio>
                <el-button type="danger" plain @click="rollback(item)">rebuild</el-button>
              </el-row>
            </el-row>
          </el-radio-group>
          <el-pagination
            :total="previewPagination.total"
            :page-size="previewPagination.rows"
            :current-page.sync="previewPagination.page"
            prev-text="上一页"
            next-text="下一页"
            style="text-align:right;margin-right:20px"
            layout="total, prev, next"
            @current-change="handlePreviewPageChange"
          />
        </el-row>
        <el-row class="project-detail" style="flex:1;width:100%">
          <el-row v-for="(item, index) in publishLocalTraceList" :key="index">
            <el-row v-if="item.type === 2">
              <el-row style="margin:5px 0">git同步信息</el-row>
              <el-row style="margin:5px 0">时间: {{ item.insertTime }}</el-row>
              <!-- 用数组的形式 兼容以前版本 -->
              <el-row v-if="item.state !== 0">
                <el-row>commit: {{ item['commit'] }}</el-row>
                <el-row>message: {{ item['message'] }}</el-row>
                <el-row>author: {{ item['author'] }}</el-row>
                <el-row>datetime: {{ item['timestamp'] ? parseTime(item['timestamp']) : '' }}</el-row>
                <el-row><span v-html="formatDetail(item['diff'])" /></el-row>
              </el-row>
              <el-row v-else style="margin:5px 0">
                <el-tag type="danger" effect="plain">失败</el-tag>
                <span v-html="formatDetail(item.detail)" />
              </el-row>
            </el-row>
            <el-row v-if="item.type === 3">
              <hr>
              <el-row style="margin:5px 0">获取代码后脚本信息</el-row>
              <el-row style="margin:5px 0">时间: {{ item.insertTime }}</el-row>
              <el-row>脚本: <pre v-html="formatDetail(item.script)" /></el-row>
              <el-row style="margin:5px 0">
                <el-tag v-if="item.state === 0" type="danger" effect="plain">失败</el-tag>
                <el-row>[goploy ~]# <span v-html="formatDetail(item.detail)" /></el-row>
              </el-row>
            </el-row>
            <el-row v-if="item.type === 6">
              <hr>
              <el-row>
                <el-row style="margin:5px 0">remote服务器信息</el-row>
                <el-row style="margin:5px 0">服务器: {{ item.serverName }}</el-row>
                <el-row style="margin:5px 0">时间: {{ item.insertTime }}</el-row>
                <el-row style="margin:5px 0">脚本: {{ item.script }}</el-row>
                <el-row style="margin:5px 0">
                  <el-tag v-if="item.state === 0" type="danger" effect="plain">失败</el-tag>
                  <span v-html="formatDetail(item.detail)" />
                </el-row>
              </el-row>
            </el-row>
          </el-row>
          <hr>
          <el-row style="margin:5px 0">目标服务器</el-row>
          <el-tabs v-model="activeRomoteTracePane">
            <el-tab-pane v-for="(item, serverName) in publishRemoteTraceList" :key="serverName" :label="serverName" :name="serverName">
              <el-row v-for="(trace, key) in item" :key="key">
                <el-row v-if="trace.type === 4">
                  <el-row style="margin:5px 0">部署前脚本</el-row>
                  <el-row style="margin:5px 0">时间: {{ trace.insertTime }}</el-row>
                  <el-row>脚本: <pre v-html="formatDetail(trace.script)" /></el-row>
                  <el-row style="margin:5px 0">
                    <el-tag v-if="trace.state === 0" type="danger" effect="plain">失败</el-tag>
                    <el-row>[goploy ~]# <span v-html="formatDetail(trace.detail)" /></el-row>
                  </el-row>
                </el-row>
                <el-row v-else-if="trace.type === 5">
                  <el-row style="margin:5px 0">rsync同步文件</el-row>
                  <el-row style="margin:5px 0">时间: {{ trace.insertTime }}</el-row>
                  <el-row>命令: {{ trace.command }}</el-row>
                  <el-row style="margin:5px 0">
                    <el-tag v-if="trace.state === 0" type="danger" effect="plain">失败</el-tag>
                    <span v-html="formatDetail(trace.detail)" />
                  </el-row>
                </el-row>
                <el-row v-else>
                  <el-row style="margin:5px 0">部署后脚本</el-row>
                  <el-row style="margin:5px 0">时间: {{ trace.insertTime }}</el-row>
                  <el-row>脚本: {{ trace.script }}</el-row>
                  <el-row style="margin:5px 0">
                    <el-tag v-if="trace.state === 0" type="danger" effect="plain">失败</el-tag>
                    <el-row>[goploy ~]# <span v-html="formatDetail(trace.detail)" /></el-row>
                  </el-row>
                </el-row>
              </el-row>
            </el-tab-pane>
          </el-tabs>
        </el-row>
      </el-row>
    </el-dialog>
    <el-dialog title="commit管理" :visible.sync="commitDialogVisible">
      <el-table
        v-loading="commitTableLoading"
        element-loading-text="获取最近的commit历史，请稍候"
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
        <el-table-column prop="commit" label="commit" width="290" />
        <el-table-column prop="author" label="author" />
        <el-table-column label="提交时间" width="135">
          <template slot-scope="scope">
            {{ parseTime(scope.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="message" label="message" width="200" show-overflow-tooltip />
        <el-table-column prop="operation" label="操作" width="80" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button type="danger" @click="rollback(scope.row)">构建</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="commitDialogVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="定时构建设置" :visible.sync="taskDialogVisible" width="600px">
      <el-form ref="taskForm" :rules="taskFormRules" :model="taskFormData" label-width="120px">
        <el-form-item label="项目名称">
          <span>{{ taskFormProps.projectName }}</span>
        </el-form-item>
        <el-form-item label="commitId" prop="commitId">
          <el-select v-model="taskFormData.commitId" placeholder="选择CommitID" style="width: 400px">
            <el-option
              v-for="(item, index) in taskFormProps.commitOptions"
              :key="index"
              :label="item.commit+'('+item.author+')'"
              :value="item.commit"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="时间" prop="date">
          <el-date-picker
            v-model="taskFormData.date"
            :picker-options="taskFormProps.pickerOptions"
            type="datetime"
            placeholder="选择日期时间"
            value-format="yyyy-MM-dd HH:mm:ss"
            style="width: 400px"
          />
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="taskDialogVisible = false">取 消</el-button>
        <el-button :disabled="taskFormProps.disabled" type="primary" @click="submitTask">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="定时构建管理" :visible.sync="taskListDialogVisible">
      <el-table
        v-loading="taskTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="taskTableData"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="projectName" label="项目名称" width="150" />
        <el-table-column prop="commitId" label="commit" width="290" />
        <el-table-column prop="date" label="日期" width="150" />
        <el-table-column prop="isRun" label="任务" width="60">
          <template slot-scope="scope">
            {{ scope.row.isRun === 1 ? '已运行': '未运行' }}
          </template>
        </el-table-column>
        <el-table-column prop="state" label="状态" width="50">
          <template slot-scope="scope">
            {{ scope.row.state === 1 ? '有效': '无效' }}
          </template>
        </el-table-column>
        <el-table-column prop="creator" label="创建人" />
        <el-table-column prop="editor" label="修改人" />
        <el-table-column prop="insertTime" label="插入时间" width="135" />
        <el-table-column prop="updateTime" label="更新时间" width="135" />
        <el-table-column prop="operation" label="操作" width="150" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button type="primary" :disabled="scope.row.isRun === 1 || scope.row.state === 0" @click="handleEditProjectTask(scope.row)">修改</el-button>
            <el-button type="danger" :disabled="scope.row.isRun === 1 || scope.row.state === 0" @click="removeProjectTask(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="taskListDialogVisible = false">取 消</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import tableHeight from '@/mixin/tableHeight'
import { getList, getDetail, getPreview, getCommitList, publish } from '@/api/deploy'
import { addTask, editTask, removeTask, getTaskList } from '@/api/project'
import { getOption as getUserOption } from '@/api/user'
import { parseTime } from '@/utils'

export default {
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
          { required: true, message: '请选择CommitID', trigger: 'change' }
        ],
        date: [
          { required: true, message: '请选择日期', trigger: 'change' }
        ]
      },
      searchPreview: {
        loading: false,
        projectId: '',
        userId: '',
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
          this.tableData[projectIndex].tagText = '失败'
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
          element.tagText = '未构建'
          if (element.deployState === 2) {
            element.progressPercentage = 100
            element.progressStatus = 'success'
            element.tagType = 'success'
            element.tagText = '成功'
          } else if (element.deployState === 1) {
            element.progressPercentage = 60
            element.progressStatus = 'warning'
            element.tagType = 'warning'
            element.tagText = '构建中'
          }
          try {
            Object.assign(element, JSON.parse(element.publishExt))
          } catch (error) {
            console.log('项目尚未构建')
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
      if (data.environment === '生产环境') {
        color = 'color: #F56C6C'
      } else if (data.environment === '测试环境') {
        color = 'color: #E6A23C'
      } else {
        color = 'color: #909399'
      }
      this.$confirm('', '提示', {
        message: h('p', null, [
          h('span', null, '此操作将部署 ' + data.name),
          h('b', { style: color }, '(' + data.environment + ')'),
          h('span', null, ', 是否继续? ')
        ]),
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.gitLog = []
        this.remoteLog = {}
        publish(id, '').then((response) => {
          const projectIndex = this.tableData.findIndex(element => element.id === id)
          this.tableData[projectIndex].deployState = 1
        })
      }).catch(() => {
        this.$message.info('已取消构建')
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
          return Object.assign(element, { projectId: id })
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
              this.$message.success('添加成功')
            }).finally(() => {
              this.taskFormProps.disabled = false
              this.taskDialogVisible = false
            })
          } else {
            editTask(this.taskFormData).then(response => {
              this.$message.success('修改成功')
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
      this.$confirm('此操作删除' + data.projectName + '的定时任务, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
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
        this.$message.info('已取消操作')
      })
    },

    rollback(data) {
      this.$confirm('此操作将重新构建' + data.commit + ', 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        publish(data.projectId, data.commit).then((response) => {
          const projectIndex = this.tableData.findIndex(element => element.id === data.projectId)
          this.tableData[projectIndex].deployState = 1
          this.commitDialogVisible = false
          this.dialogVisible = false
        })
      }).catch(() => {
        this.$message.info('已取消重新构建')
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
    overflow:hidden; //超出的文本隐藏
    vertical-align: top;
    text-overflow:ellipsis; //溢出用省略号显示
    white-space:nowrap; //溢出不换行
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
