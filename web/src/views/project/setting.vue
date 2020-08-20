<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input v-model="projectName" style="width:200px" placeholder="请输入项目名称" />
        <el-button type="primary" icon="el-icon-search" @click="searchProjectList">搜索</el-button>
      </el-row>
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd">添加</el-button>
    </el-row>
    <el-table
      :key="tableHeight"
      border
      stripe
      highlight-current-row
      :max-height="tableHeight"
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="name" label="项目名称" width="200" />
      <el-table-column prop="url" label="项目地址" width="350">
        <template slot-scope="scope">
          <el-link
            type="primary"
            style="font-size: 12px"
            :underline="false"
            :href="parseGitURL(scope.row.url)"
            target="_blank"
          >
            {{ scope.row.url }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column prop="path" label="部署路径" min-width="200" />
      <el-table-column prop="environment" width="120" label="环境" align="center" />
      <el-table-column prop="branch" width="160" label="分支" align="center" />
      <el-table-column width="90" label="自动部署">
        <template slot-scope="scope">
          <span v-if="scope.row.autoDeploy === 0">关闭</span>
          <span v-else>webhook</span>
          <el-button type="text" icon="el-icon-edit" @click="handleAutoDeploy(scope.row)" />
        </template>
      </el-table-column>
      <el-table-column prop="server" width="80" label="服务器" align="center">
        <template slot-scope="scope">
          <el-button type="text" @click="handleServer(scope.row)">查看</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="user" width="80" label="成员" align="center">
        <template slot-scope="scope">
          <el-button type="text" @click="handleUser(scope.row)">查看</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="operation" label="操作" width="150" align="center" fixed="right">
        <template slot-scope="scope">
          <el-button type="primary" @click="handleEdit(scope.row)">编辑</el-button>
          <el-button type="danger" @click="handleRemove(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px;">
      <el-pagination
        hide-on-single-page
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
    <el-dialog title="项目设置" :visible.sync="dialogVisible" width="60%" class="project-setting-dialog" :close-on-click-modal="false">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="90px">
        <el-tabs v-model="formProps.tab" @tab-click="handleTabClick">
          <el-tab-pane label="基本配置" name="base">
            <el-form-item label="项目名称" prop="name">
              <el-input v-model.trim="formData.name" autocomplete="off" />
            </el-form-item>
            <el-form-item label="项目地址" prop="url">
              <el-row type="flex">
                <el-input v-model.trim="formData.url" autocomplete="off" @change="formProps.branch = []" />
                <el-button :icon="formProps.lsBranchLoading ? 'el-icon-loading' : 'el-icon-view'" type="success" :disabled="formProps.lsBranchLoading" @click="getRemoteBranchList">测试连接</el-button>
              </el-row>
            </el-form-item>
            <el-form-item label="部署路径" prop="path">
              <el-input v-model.trim="formData.path" autocomplete="off" />
            </el-form-item>
            <el-form-item label="环境" prop="environment">
              <el-select v-model="formData.environment" placeholder="选择环境" style="width:100%">
                <el-option label="生产环境" value="生产环境" />
                <el-option label="预发布环境" value="预发布环境" />
                <el-option label="测试环境" value="测试环境" />
                <el-option label="开发环境" value="开发环境" />
              </el-select>
            </el-form-item>
            <el-form-item label="分支" prop="branch">
              <el-row type="flex">
                <el-select
                  v-model="formData.branch"
                  filterable
                  allow-create
                  default-first-option
                  placeholder="请选择或填入"
                  style="width:100%"
                >
                  <el-option
                    v-for="item in formProps.branch"
                    :key="item"
                    :label="item"
                    :value="item"
                  />
                </el-select>
                <el-button :icon="formProps.lsBranchLoading ? 'el-icon-loading' : 'el-icon-search'" type="success" :disabled="formProps.lsBranchLoading" @click="getRemoteBranchList">列出分支</el-button>
              </el-row>
            </el-form-item>
            <el-form-item label="rsync选项" prop="rsyncOption">
              <el-input v-model.trim="formData.rsyncOption" type="textarea" :rows="2" autocomplete="off" placeholder="-rtv --exclude .git --delete-after" />
            </el-form-item>
            <el-form-item v-show="formProps.showServers" label="绑定服务器" prop="serverIds">
              <el-select v-model="formData.serverIds" multiple placeholder="选择服务器，可多选" style="width:100%">
                <el-option
                  v-for="(item, index) in serverOption"
                  :key="index"
                  :label="item.label"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item v-show="formProps.showUsers" label="绑定用户" prop="userIds">
              <el-select v-model="formData.userIds" multiple placeholder="选择用户，可多选" style="width:100%">
                <el-option
                  v-for="(item, index) in userOption.filter(item => [$global.Admin, $global.Manager].indexOf(item.role) === -1)"
                  :key="index"
                  :label="item.userName"
                  :value="item.userId"
                />
              </el-select>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="软链部署(推荐)" name="symlink">
            <el-row style="margin: 0 10px">
              <p>项目先同步到指定目录(rsync 软链目录)，然后ln -s 部署路径 软链目录</p>
              <p>可以避免项目在同步传输文件的过程中，外部访问到部分正在同步的文件</p>
              <p>备份最近10次的部署文件，以便快速回滚</p>
            </el-row>
            <el-form-item label="" label-width="10px">
              <el-radio-group v-model="formProps.symlink">
                <el-radio :label="false">关闭</el-radio>
                <el-radio :label="true">开启</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item v-show="formProps.symlink" label="目录" prop="symlink_path" label-width="50px">
              <el-input v-model.trim="formData.symlinkPath" autocomplete="off" />
            </el-form-item>
            <el-row v-show="formProps.symlink" style="margin: 0 10px">
              <p>如果部署路径已存在在目标服务器，请手动删除该目录<span style="color: red">rm -rf 部署路径</span>，否则软链将会不成功</p>
              <p>如须更换目录，务必手动迁移原先的目录到目标服务器</p>
            </el-row>
          </el-tab-pane>
          <el-tab-pane name="afterPullScript">
            <span slot="label">
              拉取后运行脚本
              <el-tooltip class="item" effect="dark" placement="bottom">
                <div slot="content">
                  拉取代码后在宿主服务器运行的脚本<br>
                  运行方式：打包成一份脚本文件<br>
                  检查服务器是否安装该脚本类型(默认以bash运行)<br>
                </div>
                <i class="el-icon-question" style="padding-left: 3px" />
              </el-tooltip>
            </span>
            <el-form-item prop="afterPullScript" label-width="0px">
              <el-select v-model="formData.afterPullScriptMode" placeholder="脚本类型(默认bash)" style="width:100%" @change="handleScriptModeChange">
                <el-option
                  v-for="(item, index) in scriptModeOption"
                  :key="index"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item prop="afterPullScript" label-width="0px">
              <codemirror ref="afterPullScript" v-model="formData.afterPullScript" :options="cmOption" placeholder="已切换至项目目录..." />
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane name="afterDeployScript">
            <span slot="label">
              部署后运行脚本
              <el-tooltip class="item" effect="dark" placement="bottom">
                <div slot="content">
                  部署后在目标服务器运行的脚本<br>
                  运行方式：打包成一份脚本文件<br>
                  如需重启服务，请注意是否需要nohup<br>
                  检查服务器是否安装该脚本类型(默认以bash运行)
                </div>
                <i class="el-icon-question" style="padding-left: 3px" />
              </el-tooltip>
            </span>
            <el-form-item prop="afterDeployScript" label-width="0px">
              <el-select v-model="formData.afterDeployScriptMode" placeholder="脚本类型(默认bash)" style="width:100%" @change="handleScriptModeChange">
                <el-option
                  v-for="(item, index) in scriptModeOption"
                  :key="index"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item prop="afterDeployScript" label-width="0px">
              <codemirror ref="afterDeployScript" v-model="formData.afterDeployScript" :options="cmOption" />
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane label="高级配置" name="advance">
            <el-form-item label="构建通知" prop="notifyTarget">
              <el-row type="flex">
                <el-select v-model="formData.notifyType" clearable>
                  <el-option label="无" :value="0" />
                  <el-option label="企业微信" :value="1" />
                  <el-option label="钉钉" :value="2" />
                  <el-option label="飞书" :value="3" />
                  <el-option label="自定义" :value="255" />
                </el-select>
                <el-input v-model.trim="formData.notifyTarget" autocomplete="off" placeholder="webhook链接" />
              </el-row>
            </el-form-item>
          </el-tab-pane>
        </el-tabs>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="自动部署设置" :visible.sync="dialogAutoDeployVisible">
      <el-form ref="autoDeployForm" :model="autoDeployFormData">
        <el-row style="margin: 10px">构建触发器：达成某种条件后自动构建发布项目</el-row>
        <el-radio-group v-model="autoDeployFormData.autoDeploy" style="margin: 10px">
          <el-radio :label="0">关闭</el-radio>
          <el-radio :label="1">webhook</el-radio>
        </el-radio-group>
        <el-row v-show="autoDeployFormData.autoDeploy===1" style="margin: 10px">
          前往GitLab、GitHub或Gitee的webhook（可前往谷歌查找各自webhook所在的位置）<br>
          填入连接<span style="color: red">http(s)://域名(IP)/deploy/webhook?project_name={{ autoDeployFormProps.name }}</span><br>
          勾选push event即可
        </el-row>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogAutoDeployVisible = false">取 消</el-button>
        <el-button :disabled="autoDeployFormProps.disabled" type="primary" @click="setAutoDeploy">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="服务器管理" :visible.sync="dialogServerVisible">
      <el-row class="app-bar" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="handleAddServer">添加</el-button>
      </el-row>
      <el-table
        border
        stripe
        highlight-current-row
        :data="tableServerData"
        style="width: 100%"
      >
        <el-table-column prop="serverId" label="服务器ID" width="100" />
        <el-table-column prop="serverName" label="服务器名称" width="100" />
        <el-table-column prop="serverDescription" label="服务器描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="insertTime" width="160" label="绑定时间" />
        <el-table-column prop="updateTime" width="160" label="更新时间" />
        <el-table-column prop="operation" label="操作" width="80">
          <template slot-scope="scope">
            <el-button type="danger" @click="removeServer(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogServerVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="成员管理" :visible.sync="dialogUserVisible">
      <el-row class="app-bar" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="handleAddUser">添加</el-button>
      </el-row>
      <el-table
        border
        stripe
        highlight-current-row
        :data="tableUserData"
        style="width: 100%"
      >
        <el-table-column prop="userId" label="用户ID" />
        <el-table-column prop="userName" label="用户名称" />
        <el-table-column prop="insertTime" width="160" label="绑定时间" />
        <el-table-column prop="updateTime" width="160" label="更新时间" />
        <el-table-column prop="operation" label="操作" width="80">
          <template slot-scope="scope">
            <el-button v-show="scope.row.role !== $global.GroupManager" type="danger" @click="removeUser(scope.row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogUserVisible = false">取 消</el-button>
      </div>
    </el-dialog>
    <el-dialog title="添加服务器" :visible.sync="dialogAddServerVisible">
      <el-form ref="addServerForm" :rules="addServerFormRules" :model="addServerFormData">
        <el-form-item label="绑定服务器" label-width="120px" prop="serverIds">
          <el-select v-model="addServerFormData.serverIds" multiple placeholder="选择服务器，可多选">
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.label"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogAddServerVisible = false">取 消</el-button>
        <el-button :disabled="addServerFormProps.disabled" type="primary" @click="addServer">确 定</el-button>
      </div>
    </el-dialog>
    <el-dialog title="添加用户" :visible.sync="dialogAddUserVisible">
      <el-form ref="addUserForm" :rules="addUserFormRules" :model="addUserFormData">
        <el-form-item label="绑定用户" label-width="120px" prop="userIds">
          <el-select v-model="addUserFormData.userIds" multiple placeholder="选择用户，可多选">
            <el-option
              v-for="(item, index) in userOption.filter(item => [$global.Admin, $global.Manager].indexOf(item.role) === -1)"
              :key="index"
              :label="item.userName"
              :value="item.userId"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogAddUserVisible = false">取 消</el-button>
        <el-button :disabled="addUserFormProps.disabled" type="primary" @click="addUser">确 定</el-button>
      </div>
    </el-dialog>
  </el-row>
</template>
<script>
import tableHeight from '@/mixin/tableHeight'
import { parseGitURL } from '@/utils'
import { getUserOption } from '@/api/namespace'
import { getOption as getServerOption } from '@/api/server'
import {
  getList,
  getTotal,
  getBindServerList,
  getBindUserList,
  getRemoteBranchList,
  add,
  edit,
  remove,
  addServer,
  addUser,
  setAutoDeploy,
  removeServer,
  removeUser
} from '@/api/project'
// require component
import { codemirror } from 'vue-codemirror'
import 'codemirror/mode/shell/shell.js'
import 'codemirror/mode/php/php.js'
import 'codemirror/mode/python/python.js'
import 'codemirror/theme/darcula.css'
// require styles
import 'codemirror/lib/codemirror.css'
import 'codemirror/addon/scroll/simplescrollbars.js'
import 'codemirror/addon/scroll/simplescrollbars.css'
import 'codemirror/addon/display/placeholder.js'
export default {
  components: {
    codemirror
  },
  mixins: [tableHeight],
  data() {
    const validateNotifyTarget = (rule, value, callback) => {
      if (value !== 0 && this.formData.notifyType !== '') {
        callback()
      } else if (value === 0 && this.formData.notifyType === '') {
        callback()
      } else {
        callback(new Error('请选择推送类型'))
      }
    }
    return {
      scriptModeOption: [
        { label: 'sh', value: 'sh', mode: 'text/x-sh' },
        { label: 'zsh', value: 'zsh', mode: 'text/x-sh' },
        { label: 'bash', value: 'bash', mode: 'text/x-sh' },
        { label: 'python', value: 'python', mode: 'text/x-python' },
        { label: 'php', value: 'php', mode: 'text/x-php' }
      ],
      cmOption: {
        tabSize: 4,
        mode: 'text/x-sh',
        lineNumbers: true,
        line: true,
        scrollbarStyle: 'overlay',
        theme: 'darcula'
      },
      projectName: '',
      dialogVisible: false,
      dialogAutoDeployVisible: false,
      dialogServerVisible: false,
      dialogUserVisible: false,
      dialogAddServerVisible: false,
      dialogAddUserVisible: false,
      serverOption: [],
      userOption: [],
      tableData: [],
      pagination: {
        page: 1,
        rows: 17,
        total: 0
      },
      tableServerData: [],
      tableUserData: [],
      formProps: {
        symlink: false,
        disabled: false,
        branch: [],
        lsBranchLoading: false,
        showServers: true,
        showUsers: true,
        tab: 'base'
      },
      tempFormData: {},
      formData: {
        id: 0,
        name: '',
        url: '',
        path: '',
        symlinkPath: '',
        afterPullScriptMode: '',
        afterPullScript: '',
        afterDeployScriptMode: '',
        afterDeployScript: '',
        environment: '生产环境',
        branch: 'master',
        rsyncOption: '-rtv --exclude .git --delete-after',
        serverIds: [],
        userIds: [],
        notifyType: 0,
        notifyTarget: ''
      },
      formRules: {
        name: [
          { required: true, message: '请输入项目名称', trigger: ['blur'] }
        ],
        url: [
          { required: true, message: '请输入项目地址', trigger: ['blur'] }
        ],
        path: [
          { required: true, message: '请输入部署路径', trigger: ['blur'] }
        ],
        environment: [
          { required: true, message: '请选择环境', trigger: ['blur'] }
        ],
        branch: [
          { required: true, message: '请输入分支名称', trigger: ['blur'] }
        ],
        serverIds: [
          { type: 'array', message: '请选择服务器', trigger: 'change' }
        ],
        userIds: [
          { type: 'array', message: '请选择组员', trigger: 'change' }
        ],
        notifyTarget: [
          { validator: validateNotifyTarget, trigger: 'blur' }
        ]
      },
      autoDeployFormProps: {
        disabled: false,
        name: ''
      },
      autoDeployFormData: {
        id: 0,
        autoDeploy: 0
      },
      addServerFormProps: {
        disabled: false
      },
      addServerFormData: {
        projectId: 0,
        serverIds: []
      },
      addServerFormRules: {
        serverIds: [
          { type: 'array', required: true, message: '请选择服务器', trigger: 'change' }
        ]
      },
      addUserFormProps: {
        disabled: false
      },
      addUserFormData: {
        projectId: 0,
        userIds: []
      },
      addUserFormRules: {
        userIds: [
          { type: 'array', required: true, message: '请选择用户', trigger: 'change' }
        ]
      }
    }
  },
  created() {
    this.storeFormData()
    this.getOptions()
    this.getList()
    this.getTotal()
  },
  methods: {
    parseGitURL,
    handleAdd() {
      this.restoreFormData()
      this.formProps.showServers = this.formProps.showUsers = true
      this.formProps.symlink = false
      this.dialogVisible = true
    },

    handleEdit(data) {
      this.formData = Object.assign({}, data)
      this.formData.serverIds = []
      this.formData.userIds = []
      this.formProps.symlink = this.formData.symlinkPath !== ''
      this.formProps.showServers = this.formProps.showUsers = false
      this.formProps.branch = []
      this.dialogVisible = true
    },

    handleRemove(data) {
      this.$confirm('此操作将删除该项目, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        remove(data.id).then((response) => {
          this.$message.success('删除成功')
          this.getList()
          this.getTotal()
        })
      }).catch(() => {
        this.$message.info('已取消删除')
      })
    },

    handleTabClick(vueEvent) {
      const name = vueEvent.name
      // 需要刷新 不然无法出现光标
      if (name === 'afterPullScript') {
        this.$refs.afterPullScript.refresh()
        this.handleScriptModeChange(this.formData.afterPullScriptMode)
      } else if (name === 'afterDeployScript') {
        this.$refs.afterDeployScript.refresh()
        this.handleScriptModeChange(this.formData.afterDeployScriptMode)
      }
    },

    handleScriptModeChange(scriptMode) {
      if (scriptMode !== '') {
        this.cmOption.mode = this.scriptModeOption.find(elem => elem.value === scriptMode)['mode']
      } else {
        this.cmOption.mode = 'text/x-sh'
      }
    },

    handleAutoDeploy(data) {
      this.dialogAutoDeployVisible = true
      this.autoDeployFormProps.name = data.name
      this.autoDeployFormData.id = data.id
      this.autoDeployFormData.autoDeploy = data.autoDeploy
    },

    handleServer(data) {
      this.getBindServerList(data.id)
      // 先把projectID写入添加服务器的表单
      this.addServerFormData.projectId = data.id
      this.dialogServerVisible = true
    },

    handleUser(data) {
      this.getBindUserList(data.id)
      // 先把projectID写入添加用户的表单
      this.addUserFormData.projectId = data.id
      this.dialogUserVisible = true
    },

    handleAddServer() {
      this.dialogAddServerVisible = true
    },

    handleAddUser() {
      this.dialogAddUserVisible = true
    },

    submit() {
      this.$refs.form.validate((valid) => {
        if (valid) {
          if (this.formProps.symlink === false) {
            this.formData.symlinkPath = ''
          }
          if (this.formData.id === 0) {
            this.add()
          } else {
            this.edit()
          }
        } else {
          return false
        }
      })
    },

    add() {
      this.formProps.disabled = true
      add(this.formData).then((response) => {
        this.dialogVisible = false
        this.$message.success('添加成功')
        this.getList()
        this.getTotal()
      }).finally(() => {
        this.formProps.disabled = false
      })
    },

    edit() {
      this.formProps.disabled = true
      edit(this.formData).then((response) => {
        this.dialogVisible = false
        this.$message.success('修改成功')
        this.getList()
      }).finally(() => {
        this.formProps.disabled = false
      })
    },

    setAutoDeploy() {
      this.$refs.autoDeployForm.validate((valid) => {
        if (valid) {
          this.autoDeployFormProps.disabled = true
          setAutoDeploy(this.autoDeployFormData).then((response) => {
            this.dialogAutoDeployVisible = false
            this.$message.success('添加成功')
            this.getList()
          }).finally(() => {
            this.autoDeployFormProps.disabled = false
          })
        } else {
          return false
        }
      })
    },

    addServer() {
      this.$refs.addServerForm.validate((valid) => {
        if (valid) {
          this.addServerFormProps.disabled = true
          addServer(this.addServerFormData).then((response) => {
            this.dialogAddServerVisible = false
            this.$message.success('添加成功')
            this.getBindServerList(this.addServerFormData.projectId)
          }).finally(() => {
            this.addServerFormProps.disabled = false
          })
        } else {
          return false
        }
      })
    },

    addUser() {
      this.$refs.addUserForm.validate((valid) => {
        if (valid) {
          this.addUserFormProps.disabled = true
          addUser(this.addUserFormData).then((response) => {
            this.dialogAddUserVisible = false
            this.$message.success('添加成功')
            this.getBindUserList(this.addUserFormData.projectId)
          }).finally(() => {
            this.addUserFormProps.disabled = false
          })
        } else {
          return false
        }
      })
    },

    removeServer(data) {
      this.$confirm('此操作将永久删除该服务器的绑定关系, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        removeServer(data.id).then((response) => {
          this.$message.success('删除成功')
          this.getBindServerList(data.projectId)
        })
      }).catch(() => {
        this.$message.info('已取消删除')
      })
    },

    removeUser(data) {
      this.$confirm('此操作将永久删除该用户的绑定关系, 是否继续?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      }).then(() => {
        removeUser(data.id).then((response) => {
          this.$message.success('删除成功')
          this.getBindUserList(data.projectId)
        })
      }).catch(() => {
        this.$message.info('已取消删除')
      })
    },

    getOptions() {
      getServerOption().then((response) => {
        this.serverOption = response.data.list
        this.serverOption.map(element => {
          element.label = element.name + (element.description.length > 0 ? '(' + element.description + ')' : '')
          return element
        })
      })
      getUserOption().then((response) => {
        this.userOption = response.data.list
      })
    },

    getList() {
      getList(this.pagination, this.projectName).then((response) => {
        this.tableData = response.data.list
      }).catch(() => {
      })
    },

    getTotal() {
      getTotal(this.projectName).then((response) => {
        this.pagination.total = response.data.total
      })
    },

    getBindServerList(projectID) {
      getBindServerList(projectID).then((response) => {
        this.tableServerData = response.data.list
      })
    },

    getBindUserList(projectID) {
      getBindUserList(projectID).then((response) => {
        this.tableUserData = response.data.list
      })
    },

    getRemoteBranchList() {
      // 已获取过分支
      if (this.formProps.branch.length > 0) {
        this.$message.info('无需重复获取')
        return
      }
      this.formProps.lsBranchLoading = true
      getRemoteBranchList(this.formData.url).then((response) => {
        this.formProps.branch = response.data.branch
        this.$message.success('获取成功')
      }).finally(() => {
        this.formProps.lsBranchLoading = false
      })
    },

    searchProjectList() {
      this.pagination.page = 1
      this.getList()
      this.getTotal()
    },

    handlePageChange(val) {
      this.pagination.page = val
      this.getList()
    },

    storeFormData() {
      this.tempFormData = JSON.parse(JSON.stringify(this.formData))
    },

    restoreFormData() {
      this.formData = JSON.parse(JSON.stringify(this.tempFormData))
    }
  }
}
</script>
<style lang="scss" scoped>

.project-setting-dialog {
  >>> .el-dialog__body {
    padding-top: 10px;
  }
}

</style>

<style>
.CodeMirror {
  border-radius: 4px;
  border: 1px solid #DCDFE6;
  height: 400px;
}
.CodeMirror-placeholder {
  color: #888 !important;
}
</style>
