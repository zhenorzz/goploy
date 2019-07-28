<template>
  <el-row class="app-container">
    <el-table
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="160" />
      <el-table-column prop="name" label="项目名称" />
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
      <el-table-column prop="operation" label="操作" width="150">
        <template slot-scope="scope">
          <el-button type="primary" @click="publish(scope.row.id)">构建</el-button>
          <el-button type="success" @click="handleDetail(scope.row.id)">详情</el-button>
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
  </el-row>
</template>
<script>
import { getList, getDetail, getSyncDetail, publish, rollback } from '@/api/deploy'
import { parseTime } from '@/utils'

export default {
  data() {
    return {
      gitTraceId: '1',
      publishDialogVisible: false,
      dialogVisible: false,
      webSocket: null,
      tableData: [],
      gitTrace: {},
      gitTraceList: [],
      remoteTraceList: [],
      gitLog: [],
      remoteLog: {}
    }
  },
  created() {
    this.getList()
  },
  methods: {
    initWebSocket() {
      try {
        if (this.webSocket === null) {
          this.webSocket = new WebSocket('ws://' + window.location.host + process.env.VUE_APP_BASE_API + '/deploy/sync')
          this.initEventHandle()
        }
      } catch (e) {
        console.log(e)
      }
    },
    initEventHandle() {
      this.webSocket.onopen = () => {
        console.log('WebSocket连接成功')
      }

      this.webSocket.onerror = () => {
        console.log('WebSocket连接发生错误')
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
      this.webSocket.onclose = (e) => {
        this.webSocket = null
        console.log('connection closed (' + e.code + ')')
      }
      // 路由跳转时结束websocket链接
      this.$router.afterEach(() => {
        this.webSocket && this.webSocket.close()
      })
    },
    getList() {
      getList().then((response) => {
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
        this.initWebSocket()
        publish(id).then((response) => {
          this.$message({
            message: response.message,
            type: 'success',
            duration: 5 * 1000
          })
          setTimeout(() => { this.getList() }, 1000)
          this.publishDialogVisible = true
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

    rollback(data) {
      this.$confirm('此操作将重新构建' + data.commit + ', 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        this.gitLog = []
        this.remoteLog = {}
        this.initWebSocket()
        rollback(data.projectId, data.commit).then((response) => {
          this.$message({
            message: response.message,
            type: 'success',
            duration: 5 * 1000
          })
          setTimeout(() => { this.getList() }, 1000)
          this.dialogVisible = false
          this.publishDialogVisible = true
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
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
