<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input v-model="projectName" style="width:200px" placeholder="Filter the project name" />
        <el-button type="primary" icon="el-icon-search" @click="searchProjectList" />
      </el-row>
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
    </el-row>
    <el-table
      :key="tableHeight"
      v-loading="tableloading"
      border
      stripe
      highlight-current-row
      :max-height="tableHeight"
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="id" label="ID" width="100" />
      <el-table-column prop="name" :label="$t('name')" width="200" />
      <el-table-column prop="url" :label="$t('projectURL')" width="350">
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
      <el-table-column prop="path" :label="$t('projectPath')" min-width="200" />
      <el-table-column prop="environment" width="120" :label="$t('environment')" align="center">
        <template slot-scope="scope">
          {{ $t(`envOption[${scope.row.environment}]`) }}
        </template>
      </el-table-column>
      <el-table-column prop="branch" width="160" :label="$t('branch')" align="center" />
      <el-table-column width="95" :label="$t('autoDeploy')">
        <template slot-scope="scope">
          <span v-if="scope.row.autoDeploy === 0">{{ $t('close') }}</span>
          <span v-else>webhook</span>
          <el-button type="text" icon="el-icon-edit" @click="handleAutoDeploy(scope.row)" />
        </template>
      </el-table-column>
      <el-table-column width="60" :label="$t('review')" align="center">
        <template slot-scope="scope">
          <span v-if="scope.row.review === 0">{{ $t('close') }}</span>
          <span v-else>{{ $t('open') }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="server" width="80" :label="$t('server')" align="center">
        <template slot-scope="scope">
          <el-button type="text" @click="handleServer(scope.row)">{{ $t('view') }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="user" width="80" :label="$t('member')" align="center">
        <template slot-scope="scope">
          <el-button type="text" @click="handleUser(scope.row)">{{ $t('view') }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="operation" :label="$t('op')" width="180" align="center" fixed="right">
        <template slot-scope="scope">
          <el-button type="primary" icon="el-icon-edit" @click="handleEdit(scope.row)" />
          <el-button type="info" icon="el-icon-document-copy" @click="handleCopy(scope.row)" />
          <el-button type="danger" icon="el-icon-delete" @click="handleRemove(scope.row)" />
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
    <el-dialog :title="$t('setting')" :visible.sync="dialogVisible" width="60%" class="project-setting-dialog" :close-on-click-modal="false">
      <el-form ref="form" :rules="formRules" :model="formData" label-width="110px">
        <el-tabs v-model="formProps.tab" @tab-click="handleTabClick">
          <el-tab-pane :label="$t('baseSetting')" name="base">
            <el-form-item :label="$t('name')" prop="name">
              <el-input v-model.trim="formData.name" autocomplete="off" placeholder="goploy" />
            </el-form-item>
            <el-form-item :label="$t('projectURL')" prop="url">
              <el-row type="flex">
                <el-input v-model.trim="formData.url" autocomplete="off" placeholder="HTTPS、HTTP、SSH" @change="formProps.branch = []" />
                <el-button
                  :icon="formProps.lsBranchLoading ? 'el-icon-loading' : 'el-icon-view'"
                  type="success"
                  :disabled="formProps.lsBranchLoading"
                  @click="getRemoteBranchList"
                >{{ $t('projectPage.testConnection') }}</el-button>
              </el-row>
            </el-form-item>
            <el-form-item :label="$t('projectPath')" prop="path">
              <el-input v-model.trim="formData.path" autocomplete="off" placeholder="/var/www/goploy" />
            </el-form-item>
            <el-form-item :label="$t('environment')" prop="environment">
              <el-select v-model="formData.environment" style="width:100%">
                <el-option :label="$t('envOption[1]')" :value="1" />
                <el-option :label="$t('envOption[2]')" :value="2" />
                <el-option :label="$t('envOption[3]')" :value="3" />
                <el-option :label="$t('envOption[4]')" :value="4" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('branch')" prop="branch">
              <el-row type="flex">
                <el-select
                  v-model="formData.branch"
                  filterable
                  allow-create
                  default-first-option
                  style="width:100%"
                >
                  <el-option
                    v-for="item in formProps.branch"
                    :key="item"
                    :label="item"
                    :value="item"
                  />
                </el-select>
                <el-button
                  :icon="formProps.lsBranchLoading ? 'el-icon-loading' : 'el-icon-search'"
                  type="success"
                  :disabled="formProps.lsBranchLoading"
                  @click="getRemoteBranchList"
                >{{ $t('projectPage.lishBranch') }}</el-button>
              </el-row>
            </el-form-item>
            <el-form-item prop="rsyncOption">
              <span slot="label">Rsync<br> [OPTION...]</span>
              <el-input v-model.trim="formData.rsyncOption" type="textarea" :rows="3" autocomplete="off" placeholder="-rtv --exclude .git --delete-after" />
            </el-form-item>
            <el-form-item :label="$t('projectPage.deployNotice')" prop="notifyTarget">
              <el-row type="flex">
                <el-select v-model="formData.notifyType" clearable>
                  <el-option :label="$t('webhookOption[0]')" :value="0" />
                  <el-option :label="$t('webhookOption[1]')" :value="1" />
                  <el-option :label="$t('webhookOption[2]')" :value="2" />
                  <el-option :label="$t('webhookOption[3]')" :value="3" />
                  <el-option :label="$t('webhookOption[255]')" :value="255" />
                </el-select>
                <el-input v-model.trim="formData.notifyTarget" autocomplete="off" placeholder="webhook" />
              </el-row>
            </el-form-item>
            <el-form-item v-show="formProps.showServers" :label="$t('server')" prop="serverIds">
              <el-select v-model="formData.serverIds" multiple style="width:100%">
                <el-option
                  v-for="(item, index) in serverOption"
                  :key="index"
                  :label="item.label"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item v-show="formProps.showUsers" :label="$t('user')" prop="userIds">
              <el-select v-model="formData.userIds" multiple style="width:100%">
                <el-option
                  v-for="(item, index) in userOption.filter(item => [$global.Admin, $global.Manager].indexOf(item.role) === -1)"
                  :key="index"
                  :label="item.userName"
                  :value="item.userId"
                />
              </el-select>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane :label="$t('projectPage.publishReview')" name="review">
            <el-form-item label="" label-width="10px">
              <el-radio-group v-model="formData.review">
                <el-radio :label="0">{{ $t('close') }}</el-radio>
                <el-radio :label="1">{{ $t('open') }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item v-show="formData.review" label="URL" label-width="50px">
              <el-input v-model.trim="formProps.reviewURL" autocomplete="off" placeholder="http(s)://domain?custom-param=1" />
            </el-form-item>
            <el-form-item v-show="formData.review" :label="$t('param')" label-width="50px">
              <el-checkbox-group v-model="formProps.reviewURLParam">
                <el-checkbox
                  v-for="(item, key) in formProps.reviewURLParamOption"
                  :key="key"
                  :label="item.value"
                  :disabled="item['disabled']"
                >{{ item.label }}</el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-row v-show="formData.review" style="margin: 0 10px" v-html="$t('projectPage.reviewFooterTips')" />
          </el-tab-pane>
          <el-tab-pane :label="$t('projectPage.symlinkLabel')" name="symlink">
            <el-row style="margin: 0 10px" v-html="$t('projectPage.symlinkHeaderTips')" />
            <el-form-item label="" label-width="10px">
              <el-radio-group v-model="formProps.symlink">
                <el-radio :label="false">{{ $t('close') }}</el-radio>
                <el-radio :label="true">{{ $t('open') }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item v-show="formProps.symlink" :label="$t('directory')" prop="symlink_path" label-width="80px">
              <el-input v-model.trim="formData.symlinkPath" autocomplete="off" />
            </el-form-item>
            <el-row v-show="formProps.symlink" style="margin: 0 10px" v-html="$t('projectPage.symlinkFooterTips')" />
          </el-tab-pane>
          <el-tab-pane name="afterPullScript">
            <span slot="label">
              {{ $t('projectPage.afterPullScriptLabel') }}
              <el-tooltip class="item" effect="dark" placement="bottom">
                <div slot="content" v-html="$t('projectPage.afterPullScriptTips')" />
                <i class="el-icon-question" style="padding-left: 3px" />
              </el-tooltip>
            </span>
            <el-form-item prop="afterPullScript" label-width="0px">
              <el-select v-model="formData.afterPullScriptMode" :placeholder="$t('projectPage.scriptMode')+'(Default: bash)'" style="width:100%" @change="handleScriptModeChange">
                <el-option
                  v-for="(item, index) in scriptModeOption"
                  :key="index"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item prop="afterPullScript" label-width="0px">
              <codemirror ref="afterPullScript" v-model="formData.afterPullScript" :options="cmOption" placeholder="Already switched to project directory..." />
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane name="afterDeployScript">
            <span slot="label">
              {{ $t('projectPage.afterDeployScriptLabel') }}
              <el-tooltip class="item" effect="dark" placement="bottom">
                <div slot="content" v-html="$t('projectPage.afterDeployScriptTips')" />
                <i class="el-icon-question" style="padding-left: 3px" />
              </el-tooltip>
            </span>
            <el-form-item prop="afterDeployScript" label-width="0px">
              <el-select v-model="formData.afterDeployScriptMode" :placeholder="$t('projectPage.scriptMode')+'(Default: bash)'" style="width:100%" @change="handleScriptModeChange">
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
        </el-tabs>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="formProps.disabled" type="primary" @click="submit">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('setting')" :visible.sync="dialogAutoDeployVisible">
      <el-form ref="autoDeployForm" :model="autoDeployFormData">
        <el-row style="margin: 10px">{{ $t('projectPage.autoDeployTitle') }}</el-row>
        <el-radio-group v-model="autoDeployFormData.autoDeploy" style="margin: 10px">
          <el-radio :label="0">{{ $t('close') }}</el-radio>
          <el-radio :label="1">webhook</el-radio>
        </el-radio-group>
        <el-row v-show="autoDeployFormData.autoDeploy===1" style="margin: 10px" v-html="$t('projectPage.autoDeployTips', {projectId: autoDeployFormData.id})" />
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogAutoDeployVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="autoDeployFormProps.disabled" type="primary" @click="setAutoDeploy">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('manage')" :visible.sync="dialogServerVisible">
      <el-row class="app-bar" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="handleAddServer" />
      </el-row>
      <el-table
        border
        stripe
        highlight-current-row
        :data="tableServerData"
        style="width: 100%"
      >
        <el-table-column prop="serverId" :label="$t('serverId')" width="100" />
        <el-table-column prop="serverName" :label="$t('serverName')" min-width="120" />
        <el-table-column prop="serverDescription" :label="$t('serverDescription')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
        <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
        <el-table-column prop="operation" :label="$t('op')" width="80" align="center">
          <template slot-scope="scope">
            <el-button type="danger" icon="el-icon-delete" @click="removeServer(scope.row)" />
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogServerVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('manage')" :visible.sync="dialogUserVisible">
      <el-row class="app-bar" type="flex" justify="end">
        <el-button type="primary" icon="el-icon-plus" @click="handleAddUser" />
      </el-row>
      <el-table
        border
        stripe
        highlight-current-row
        :data="tableUserData"
        style="width: 100%"
      >
        <el-table-column prop="userId" :label="$t('userId')" width="100" />
        <el-table-column prop="userName" :label="$t('userName')" />
        <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
        <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
        <el-table-column prop="operation" :label="$t('op')" width="80" align="center">
          <template slot-scope="scope">
            <el-button v-show="scope.row.role !== $global.GroupManager" type="danger" icon="el-icon-delete" @click="removeUser(scope.row)" />
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="dialogUserVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('add')" :visible.sync="dialogAddServerVisible">
      <el-form ref="addServerForm" :rules="addServerFormRules" :model="addServerFormData">
        <el-form-item :label="$t('server')" label-width="120px" prop="serverIds">
          <el-select v-model="addServerFormData.serverIds" multiple>
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
        <el-button @click="dialogAddServerVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="addServerFormProps.disabled" type="primary" @click="addServer">{{ $t('confirm') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('add')" :visible.sync="dialogAddUserVisible">
      <el-form ref="addUserForm" :rules="addUserFormRules" :model="addUserFormData">
        <el-form-item :label="$t('member')" label-width="120px" prop="userIds">
          <el-select v-model="addUserFormData.userIds" multiple>
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
        <el-button @click="dialogAddUserVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="addUserFormProps.disabled" type="primary" @click="addUser">{{ $t('confirm') }}</el-button>
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
  name: 'Project',
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
        callback(new Error('Select the notice mode'))
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
      tableloading: false,
      tableData: [],
      pagination: {
        page: 1,
        rows: 17,
        total: 0
      },
      tableServerData: [],
      tableUserData: [],
      formProps: {
        reviewURLParamOption: [
          {
            label: 'project_id',
            value: 'project_id=__PROJECT_ID__'
          },
          {
            label: 'project_name',
            value: 'project_name=__PROJECT_NAME__'
          },
          {
            label: 'branch',
            value: 'branch=__BRANCH__'
          },
          {
            label: 'environment',
            value: 'environment=__ENVIRONMENT__'
          },
          {
            label: 'commit_id',
            value: 'commit_id=__COMMIT_ID__'
          },
          {
            label: 'publish_time',
            value: 'publish_time=__PUBLISH_TIME__'
          },
          {
            label: 'publisher_id',
            value: 'publisher_id=__PUBLISHER_ID__'
          },
          {
            label: 'publisher_name',
            value: 'publisher_name=__PUBLISHER_NAME__'
          },
          {
            label: 'callback',
            value: 'callback=__CALLBACK__',
            disabled: true
          }
        ],
        reviewURL: '',
        reviewURLParam: ['callback=__CALLBACK__'],
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
        environment: 1,
        branch: 'master',
        rsyncOption: '-rtv --exclude .git --delete-after',
        serverIds: [],
        userIds: [],
        review: 0,
        reviewURL: '',
        notifyType: 0,
        notifyTarget: ''
      },
      formRules: {
        name: [
          { required: true, message: 'Name required', trigger: ['blur'] }
        ],
        url: [
          { required: true, message: 'Repository url required', trigger: ['blur'] }
        ],
        path: [
          { required: true, message: 'Path required', trigger: ['blur'] }
        ],
        environment: [
          { required: true, message: 'Environment required', trigger: ['blur'] }
        ],
        branch: [
          { required: true, message: 'Branch required', trigger: ['blur'] }
        ],
        notifyTarget: [
          { validator: validateNotifyTarget, trigger: 'blur' }
        ]
      },
      autoDeployFormProps: {
        disabled: false
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
          { type: 'array', required: true, message: 'Server required', trigger: 'change' }
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
          { type: 'array', required: true, message: 'User required', trigger: 'change' }
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
      this.formProps.reviewURL = ''
      this.formProps.reviewURLParam = []
      if (this.formData.review === 1 && this.formData.reviewURL.length > 0) {
        const url = new URL(this.formData.reviewURL)
        this.formProps.reviewURLParamOption.forEach(item => {
          if (url.searchParams.has(item.value.split('=')[0])) {
            url.searchParams.delete(item.value.split('=')[0])
            this.formProps.reviewURLParam.push(item.value)
          }
        })
        this.formProps.reviewURL = url.href
      }
      this.dialogVisible = true
    },

    handleCopy(data) {
      this.handleEdit(data)
      this.formData.id = 0
    },

    handleRemove(data) {
      this.$confirm(this.$i18n.t('projectPage.removeProjectTips', { projectName: data.name }), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        this.tableloading = true
        remove(data.id).then((response) => {
          this.$message.success('Success')
          this.getList()
          this.getTotal()
        })
      }).catch(() => {
        this.$message.info('Cancel')
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
          if (this.formData.review === 1 && this.formProps.reviewURL.length > 0) {
            const url = new URL(this.formProps.reviewURL)
            this.formProps.reviewURLParam.forEach(param => {
              url.searchParams.set(...param.split('='))
            })
            this.formData.reviewURL = url.href
          } else {
            this.formData.reviewURL = ''
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
        this.$message.success('Success')
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
        this.$message.success('Success')
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
            this.$message.success('Success')
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
            this.$message.success('Success')
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
            this.$message.success('Success')
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
      this.$confirm(this.$i18n.t('projectPage.removeServerTips', { serverName: data.serverName }), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        removeServer(data.id).then((response) => {
          this.$message.success('Success')
          this.getBindServerList(data.projectId)
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    removeUser(data) {
      this.$confirm(this.$i18n.t('projectPage.removeUserTips', { userName: data.userName }), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        removeUser(data.id).then((response) => {
          this.$message.success('Success')
          this.getBindUserList(data.projectId)
        })
      }).catch(() => {
        this.$message.info('Cancel')
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
      this.tableloading = true
      getList(this.pagination, this.projectName).then((response) => {
        this.tableData = response.data.list
      }).finally(() => {
        this.tableloading = false
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
      if (this.formProps.branch.length > 0) {
        return
      }
      this.formProps.lsBranchLoading = true
      getRemoteBranchList(this.formData.url).then((response) => {
        this.formProps.branch = response.data.branch
        this.$message.success('Success')
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
