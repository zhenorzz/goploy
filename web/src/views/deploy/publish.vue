<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex">
      <el-select v-model="projectGroupId" placeholder="选择分组" @change="handleGroupChange">
        <el-option label="默认" :value="0" />
        <el-option
          v-for="(item, index) in projectGroupOption"
          :key="index"
          :label="item.name"
          :value="item.id"
        />
      </el-select>
      <el-input v-model="projectName" style="width:300px" placeholder="请输入项目名称" @change="getList" />
    </el-row>
    <el-table
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="160" />
      <el-table-column prop="name" label="项目名称" />
      <el-table-column prop="group" label="分组">
        <template slot-scope="scope">
          {{ findGroupName(scope.row.projectGroupId) }}
        </template>
      </el-table-column>
      <el-table-column prop="environment" label="环境" />
      <el-table-column prop="branch" label="分支" />
      <el-table-column prop="publisherName" label="构建者" width="160" />
      <el-table-column prop="publishState" label="状态" width="60">
        <template slot-scope="scope">
          <el-tag v-if="scope.row.publishState === 1" type="success" effect="plain">成功</el-tag>
          <el-tag v-else type="danger" effect="plain">失败</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="updateTime" label="上次构建时间" width="160" />
      <el-table-column prop="operation" label="操作" width="220">
        <template slot-scope="scope">
          <el-button type="primary" @click="publish(scope.row.id)">构建</el-button>
          <el-button type="success" @click="handleDetail(scope.row.id)">详情</el-button>
          <el-button type="danger" @click="handleRollback(scope.row.id)">回滚</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="构建记录" :visible.sync="dialogVisible">
      <el-row>
        <el-col :span="8">
          <el-radio-group v-model="gitTraceId" @change="handleDetailChange">
            <el-row v-for="(item, index) in gitTraceList" :key="index">
              <el-row style="margin:5px 0">
                <el-tag v-if="item['remoteState'] === 1" type="success" effect="plain">成功</el-tag>
                <el-tag v-else type="danger" effect="plain">失败</el-tag>
                <el-radio style="margin-left: 10px;margin-right: 5px;" :label="item['id']" border>commitID: {{ item['commit'].substring(0,6) }}</el-radio>
                <el-button type="danger" icon="el-icon-refresh" plain @click="rollback(item)" />
              </el-row>
            </el-row>
          </el-radio-group>
        </el-col>
        <el-col :span="16" class="project-detail">
          <el-row>
            <el-row style="margin:5px 0">git同步信息</el-row>
            <el-row style="margin:5px 0">时间: {{ formatTime(gitTrace['createTime']) }}</el-row>
            <el-row>commit: {{ gitTrace['commit'] }}</el-row>
            <el-row style="margin:5px 0">
              <el-tag v-if=" gitTrace['state'] === 1" type="success" effect="plain">成功</el-tag>
              <el-tag v-else type="danger" effect="plain">失败</el-tag>
              <span v-html="formatDetail(gitTrace['detail'])" />
            </el-row>
          </el-row>
          <hr>
          <el-row>
            <el-row style="margin:5px 0">remote服务器信息</el-row>
            <el-row v-for="(item, index) in remoteTraceList" :key="index">
              <el-row style="margin:5px 0">服务器: {{ item['serverName'] }}</el-row>
              <el-row style="margin:5px 0">日志类型: {{ item['type'] === 1 ? '同步文件' : '运行脚本' }}</el-row>
              <el-row style="margin:5px 0">时间: {{ formatTime(item['createTime']) }}</el-row>
              <el-row style="margin:5px 0">
                <el-tag v-if="item['state'] === 1" type="success" effect="plain">成功</el-tag>
                <el-tag v-else type="danger" effect="plain">失败</el-tag>
                <span v-html="formatDetail(item['detail'])" />
              </el-row>
              <hr>
            </el-row>
          </el-row>
        </el-col>
      </el-row>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="构建进度" :visible.sync="publishDialogVisible">
      <el-row ref="publishSchedule" class="project-detail">
        <el-row>
          <el-row style="margin:5px 0">git同步信息</el-row>
          <el-row v-for="(item, index) in gitLog" :key="index">
            <el-row style="margin:5px 0">
              <el-tag v-show="item['state'] === 0" type="danger" effect="plain">失败</el-tag>
              <span v-html="item['message']" />
            </el-row>
          </el-row>
        </el-row>
        <hr>
        <el-row>
          <el-row style="margin:5px 0">remote服务器信息</el-row>
          <el-row v-for="(serverLog, index) in remoteLog" :key="index">
            <el-row style="margin:5px 0">服务器: {{ index }}</el-row>
            <el-row v-for="(item, serverIndex) in serverLog" :key="serverIndex">
              <el-row style="margin:5px 0">
                <el-tag v-show="item['state'] === 0" type="danger" effect="plain">失败</el-tag>
                <span v-html="item['message']" />
              </el-row>
            </el-row>
            <hr>
          </el-row>
        </el-row>
      </el-row>
    </el-dialog>
    <el-dialog title="commit管理" :visible.sync="commitDialogVisible">
      <el-table
        border
        stripe
        highlight-current-row
        :data="commitTableData"
      >
        <el-table-column prop="commit" label="commit" width="290" />
        <el-table-column prop="author" label="author" />
        <el-table-column prop="date" label="date" />
        <el-table-column prop="message" label="message" />
        <el-table-column prop="operation" label="操作" width="75">
          <template slot-scope="scope">
            <el-button type="danger" @click="rollback(scope.row)">构建</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="commitDialogVisible = false">取 消</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import { getList, getDetail, getSyncDetail, getCommitList, publish, rollback } from '@/api/deploy'
import { getOption as getProjectGroupOption } from '@/api/projectGroup'
import { parseTime } from '@/utils'

export default {
  data() {
    return {
      projectGroupId: parseInt(localStorage.getItem('projectGroupId')) || 0,
      projectGroupOption: [],
      projectName: '',
      gitTraceId: '1',
      publishDialogVisible: false,
      commitDialogVisible: false,
      dialogVisible: false,
      webSocket: null,
      tableData: [],
      commitTableData: [],
      gitTrace: {},
      gitTraceList: [],
      remoteTraceList: [],
      gitLog: [],
      remoteLog: {}
    }
  },
  created() {
    this.getList()
    this.getProjectGroupOption()
    // // 路由跳转时结束websocket链接
    this.$router.afterEach(() => {
      this.webSocket && this.webSocket.close()
    })
  },
  methods: {
    connectWebSocket() {
      if (this.webSocket && this.webSocket.readyState < 2) {
        console.log('reusing the socket connection [state = ' + this.webSocket.readyState + ']: ' + this.webSocket.url)
        return Promise.resolve(this.webSocket)
      }

      return new Promise((resolve, reject) => {
        this.webSocket = new WebSocket('ws://' + window.location.host + process.env.VUE_APP_BASE_API + '/deploy/sync')

        this.webSocket.onopen = () => {
          console.log('socket connection is opened [state = ' + this.webSocket.readyState + ']: ' + this.webSocket.url)
          resolve(this.webSocket)
        }

        this.webSocket.onerror = (err) => {
          console.error('socket connection error : ', err)
          reject(err)
        }

        this.webSocket.onclose = (e) => {
          this.webSocket = null
          console.log('connection closed (' + e.code + ')')
        }

        this.webSocket.onmessage = (e) => {
          const data = JSON.parse(e.data)
          data.message = this.formatDetail(data.message)
          if (data.dataType === 1) {
            this.gitLog.push(data)
          } else {
            if (!this.remoteLog[data.serverName]) {
              this.$set(this.remoteLog, data.serverName, [])
            }
            this.remoteLog[data.serverName].push(data)
          }

          this.$nextTick(() => {
            const contentBox = this.$refs.publishSchedule
            contentBox.$el.scrollTop = contentBox.$el.scrollHeight
          })
        }
      })
    },

    handleGroupChange(projectGroupId) {
      localStorage.setItem('projectGroupId', projectGroupId)
      this.projectGroupId = projectGroupId
      this.getList()
    },

    getProjectGroupOption() {
      getProjectGroupOption().then((response) => {
        this.projectGroupOption = response.data.projectGroupList || []
      })
    },

    getList() {
      getList(this.projectGroupId, this.projectName).then((response) => {
        const projectList = response.data.projectList || []
        projectList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
        })
        this.tableData = projectList
      })
    },

    publish(id) {
      this.$confirm('此操作将部署该项目, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.gitLog = []
        this.remoteLog = {}
        this.publishDialogVisible = true
        this.connectWebSocket().then(server => {
          publish(id).then((response) => {
            setTimeout(() => { this.getList() }, 1000)
          })
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消构建'
        })
      })
    },

    handleDetail(id) {
      getDetail(id).then((response) => {
        this.dialogVisible = true
        this.gitTrace = response.data.gitTrace
        this.gitTraceList = response.data.gitTraceList
        this.remoteTraceList = response.data.remoteTraceList
        this.gitTraceId = this.gitTrace.id
      })
    },

    handleDetailChange(gitTraceId) {
      getSyncDetail(gitTraceId).then(response => {
        this.gitTrace = response.data.gitTrace
        this.remoteTraceList = response.data.remoteTraceList
        this.gitTraceId = this.gitTrace.id
      })
    },

    handleRollback(id) {
      getCommitList(id).then(response => {
        this.commitTableData = response.data.commitList.map(element => {
          const commitInfo = element.split('`')
          return {
            projectId: id,
            commit: commitInfo[0],
            author: commitInfo[1],
            date: commitInfo[2],
            message: commitInfo[3]
          }
        })

        this.commitDialogVisible = true
      })
    },

    rollback(data) {
      this.$confirm('此操作将重新构建' + data.commit + ', 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.gitLog = []
        this.remoteLog = {}
        this.dialogVisible = false
        this.publishDialogVisible = true
        this.connectWebSocket().then(server => {
          rollback(data.projectId, data.commit).then((response) => {
            this.$message({
              message: response.message,
              type: 'success',
              duration: 5 * 1000
            })
            setTimeout(() => { this.getList() }, 1000)
          })
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消重新构建'
        })
      })
    },

    findGroupName(projectGroupId) {
      const projectGroup = this.projectGroupOption.find(element => element.id === projectGroupId)
      return projectGroup ? projectGroup['name'] : '默认'
    },

    formatDetail(detail) {
      return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
    },

    formatTime(timestamp) {
      return parseTime(timestamp)
    }
  }
}
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
@import "@/styles/mixin.scss";
.project-detail {
  height:500px;
  overflow-y: auto;
  @include scrollBar();
}
</style>
