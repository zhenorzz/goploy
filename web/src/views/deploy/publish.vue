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
          <el-button size="small" type="primary" @click="publish(scope.row.id)">构建</el-button>
          <el-button size="small" type="success" @click="handleDetail(scope.row.id)">详情</el-button>
          <el-button size="small" type="danger" @click="wsSend">回滚</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-dialog title="上一次构建记录" :visible.sync="dialogVisible">
      <el-row class="project-detail">
        <el-row>
          <el-row style="margin:5px 0">git同步信息</el-row>
          <el-row style="margin:5px 0">时间: {{ formatTime(gitTrace['createTime']) }}</el-row>
          <el-row style="margin:5px 0">状态: {{ gitTrace['state'] === 1 ? '成功' : '失败' }}</el-row>
          <el-row style="margin:5px 0" v-html="formatDetail(gitTrace['detail'])" />
        </el-row>
        <hr>
        <el-row>
          <el-row style="margin:5px 0">remote服务器信息</el-row>
          <el-row v-for="(item, index) in remoteTraceList" :key="index">
            <el-row style="margin:5px 0">服务器: {{ item['serverName'] }}</el-row>
            <el-row style="margin:5px 0">日志类型: {{ item['type'] === 1 ? '同步文件' : '运行脚本' }}</el-row>
            <el-row style="margin:5px 0">时间: {{ formatTime(item['createTime']) }}</el-row>
            <el-row style="margin:5px 0">状态: {{ item['state'] === 1 ? '成功' : '失败' }}</el-row>
            <el-row style="margin:5px 0" v-html="formatDetail(item['detail'])" />
            <hr>
          </el-row>
        </el-row>
      </el-row>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="构建进度" :visible.sync="publishDialogVisible">
      <el-row class="project-detail">
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
import { getList, getDetail, publish } from '@/api/deploy'
import { parseTime } from '@/utils'

const STATE = ['构建中', '构建成功', '构建失败', '撤回']
export default {
  data() {
    return {
      publishDialogVisible: false,
      dialogVisible: false,
      webSocket: null,
      tableData: [],
      gitTrace: {},
      remoteTraceList: [],
      gitLog: [],
      remoteLog: {}
    }
  },
  created() {
    this.getList()
    this.initWebSocket()
  },
  methods: {
    initWebSocket() {
      try {
        this.webSocket = new WebSocket('ws://localhost:3000/deploy/sync')
        this.initEventHandle()
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
        data.messgae = this.formatDetail(data.messgae)
        if (data.dataType === 1) {
          this.gitLog.push(data)
        } else {
          if (!this.remoteLog[data.serverName]) {
            this.$set(this.remoteLog, data.serverName, [])
          }
          this.remoteLog[data.serverName].push(data)
        }
      }
      this.webSocket.onclose = (e) => {
        this.webSocket = null
        console.log('connection closed (' + e.code + ')')
      }
      // 路由跳转时结束websocket链接
      this.$router.afterEach(() => {
        this.webSocket.close()
      })
    },
    wsSend() {
      this.webSocket.send('hello world')
    },
    getList() {
      getList().then((response) => {
        const projectList = response.data.projectList || []
        projectList.forEach((element) => {
          element.createTime = parseTime(element.createTime)
          element.updateTime = parseTime(element.updateTime)
          element.state = STATE[element.state]
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
        publish(id).then((response) => {
          this.$message({
            message: response.message,
            type: 'success',
            duration: 5 * 1000
          })
          this.getList()
          this.publishDialogVisible = true
        })
      }).catch(() => {
        this.$message({
          type: 'info',
          message: '已取消删除'
        })
      })
    },
    handleDetail(id) {
      getDetail(id).then((response) => {
        this.dialogVisible = true
        this.gitTrace = response.data.gitTrace
        this.remoteTraceList = response.data.remoteTraceList
      })
    },
    formatDetail(detail) {
      return detail ? detail.replace(/\n/g, '<br>') : ''
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
