<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="projectName"
          style="width: 200px"
          placeholder="Filter the project name"
        />
        <el-button
          :loading="tableloading"
          type="primary"
          icon="el-icon-search"
          @click="searchProjectList"
        />
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
      <el-table-column prop="id" label="ID" width="60" />
      <el-table-column prop="name" :label="$t('name')" width="180" />
      <el-table-column prop="url" :label="$t('projectURL')" width="330">
        <template #default="scope">
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
      <el-table-column
        prop="environment"
        width="120"
        :label="$t('environment')"
        align="center"
      >
        <template #default="scope">
          {{ $t(`envOption[${scope.row.environment || 0}]`) }}
        </template>
      </el-table-column>
      <el-table-column
        prop="branch"
        width="160"
        :label="$t('branch')"
        align="center"
      />
      <el-table-column width="95" :label="$t('autoDeploy')">
        <template #default="scope">
          <span v-if="scope.row.autoDeploy === 0">{{ $t('close') }}</span>
          <span v-else>webhook</span>
          <el-button
            type="text"
            icon="el-icon-edit"
            @click="handleAutoDeploy(scope.row)"
          />
        </template>
      </el-table-column>
      <el-table-column width="80" :label="$t('review')" align="center">
        <template #default="scope">
          <span v-if="scope.row.review === 0">{{ $t('close') }}</span>
          <span v-else>{{ $t('open') }}</span>
        </template>
      </el-table-column>
      <el-table-column
        prop="file"
        width="80"
        :label="$t('file')"
        align="center"
      >
        <template #default="scope">
          <el-button type="text" @click="handleFile(scope.row)">{{
            $t('view')
          }}</el-button>
        </template>
      </el-table-column>
      <el-table-column
        prop="server"
        width="80"
        :label="$t('server')"
        align="center"
      >
        <template #default="scope">
          <el-button type="text" @click="handleServer(scope.row)">
            {{ $t('view') }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column
        prop="user"
        width="80"
        :label="$t('member')"
        align="center"
      >
        <template #default="scope">
          <el-button type="text" @click="handleUser(scope.row)">
            {{ $t('view') }}
          </el-button>
        </template>
      </el-table-column>
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="180"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button
            type="primary"
            icon="el-icon-edit"
            @click="handleEdit(scope.row)"
          />
          <el-tooltip
            class="item"
            effect="dark"
            content="Copy"
            placement="bottom"
          >
            <el-button
              type="info"
              icon="el-icon-document-copy"
              @click="handleCopy(scope.row)"
            />
          </el-tooltip>
          <el-button
            type="danger"
            icon="el-icon-delete"
            @click="handleRemove(scope.row)"
          />
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px">
      <el-pagination
        hide-on-single-page
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        layout="prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
    <el-dialog
      v-model="dialogVisible"
      :title="$t('setting')"
      width="60%"
      custom-class="project-setting-dialog"
      :fullscreen="$store.state.app.device === 'mobile'"
      :close-on-click-modal="false"
    >
      <el-form
        ref="form"
        v-loading="formProps.disabled"
        :rules="formRules"
        :model="formData"
        label-width="120px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-tabs v-model="formProps.tab">
          <el-tab-pane :label="$t('baseSetting')" name="base">
            <el-form-item :label="$t('name')" prop="name">
              <el-input
                v-model.trim="formData.name"
                autocomplete="off"
                placeholder="goploy"
              />
            </el-form-item>
            <el-form-item prop="url">
              <template #label>
                {{ $t('projectURL') }}
                <el-tooltip placement="top">
                  <template #content>
                    ssh://[username:password@]host.xz[:port]/path/to/repo.git/<br />
                    git@host.xz[:port]/path/to/repo.git/<br />
                    http[s]://[username:password@]host.xz[:port]/path/to/repo.git/
                  </template>
                  <i class="el-icon-question" />
                </el-tooltip>
              </template>
              <el-row type="flex">
                <el-input
                  v-model.trim="formData.url"
                  style="flex: 1"
                  autocomplete="off"
                  placeholder="HTTPS、HTTP、SSH"
                  @change="formProps.branch = []"
                />
                <el-button
                  :icon="'el-icon-view'"
                  type="success"
                  :loading="formProps.lsBranchLoading"
                  @click="getRemoteBranchList"
                >
                  {{ $t('projectPage.testConnection') }}
                </el-button>
              </el-row>
            </el-form-item>
            <el-form-item :label="$t('projectPath')" prop="path">
              <el-input
                v-model.trim="formData.path"
                autocomplete="off"
                placeholder="/var/www/goploy"
              />
            </el-form-item>
            <el-form-item :label="$t('environment')" prop="environment">
              <el-select v-model="formData.environment" style="width: 100%">
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
                  style="flex: 1"
                >
                  <el-option
                    v-for="item in formProps.branch"
                    :key="item"
                    :label="item"
                    :value="item"
                  />
                </el-select>
                <el-button
                  :icon="'el-icon-search'"
                  type="success"
                  :loading="formProps.lsBranchLoading"
                  @click="getRemoteBranchList"
                >
                  {{ $t('projectPage.lishBranch') }}
                </el-button>
              </el-row>
            </el-form-item>
            <el-form-item prop="rsyncOption">
              <template #label>
                <el-link
                  type="primary"
                  :href="$t('projectPage.rsyncDoc')"
                  target="_blank"
                >
                  Rsync <i class="el-icon-question" /><br />
                  [OPTION...]
                </el-link>
              </template>
              <el-input
                v-model="formData.rsyncOption"
                type="textarea"
                :rows="3"
                autocomplete="off"
                placeholder="-rtv --exclude .git --delete-after"
              />
            </el-form-item>
            <el-form-item prop="notifyTarget">
              <template #label>
                <el-link
                  type="primary"
                  href="https://zhenorzz.gitee.io/goploy/#/dependency/notice"
                  target="_blank"
                >
                  {{ $t('projectPage.deployNotice') }}
                  <i class="el-icon-question" />
                </el-link>
              </template>
              <el-row type="flex">
                <el-select v-model="formData.notifyType" clearable>
                  <el-option :label="$t('webhookOption[0]')" :value="0" />
                  <el-option :label="$t('webhookOption[1]')" :value="1" />
                  <el-option :label="$t('webhookOption[2]')" :value="2" />
                  <el-option :label="$t('webhookOption[3]')" :value="3" />
                  <el-option :label="$t('webhookOption[255]')" :value="255" />
                </el-select>
                <el-input
                  v-if="formData.notifyType > 0"
                  v-model.trim="formData.notifyTarget"
                  style="flex: 1"
                  autocomplete="off"
                  placeholder="webhook"
                />
              </el-row>
            </el-form-item>
            <el-form-item
              v-show="formProps.showServers"
              :label="$t('server')"
              prop="serverIds"
            >
              <el-select
                v-model="formData.serverIds"
                multiple
                style="width: 100%"
              >
                <el-option
                  v-for="(item, index) in serverOption"
                  :key="index"
                  :label="item.label"
                  :value="item.id"
                />
              </el-select>
            </el-form-item>
            <el-form-item
              v-show="formProps.showUsers"
              :label="$t('user')"
              prop="userIds"
            >
              <el-select
                v-model="formData.userIds"
                multiple
                style="width: 100%"
              >
                <el-option
                  v-for="(item, index) in userOption.filter(
                    (item) =>
                      [role.Admin, role.Manager].indexOf(item.role) === -1
                  )"
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
            <el-form-item
              v-show="formData.review"
              label="URL"
              label-width="60px"
            >
              <el-input
                v-model.trim="formProps.reviewURL"
                autocomplete="off"
                placeholder="http(s)://domain?custom-param=1"
              />
            </el-form-item>
            <el-form-item
              v-show="formData.review"
              :label="$t('param')"
              label-width="60px"
            >
              <el-checkbox-group v-model="formProps.reviewURLParam">
                <el-checkbox
                  v-for="(item, key) in formProps.reviewURLParamOption"
                  :key="key"
                  :label="item.value"
                  :disabled="item['disabled']"
                >
                  {{ item.label }}
                </el-checkbox>
              </el-checkbox-group>
            </el-form-item>
            <el-row
              v-show="formData.review"
              style="margin: 0 10px"
              v-html="$t('projectPage.reviewFooterTips')"
            />
          </el-tab-pane>
          <el-tab-pane :label="$t('projectPage.symlinkLabel')" name="symlink">
            <el-row
              style="margin: 0 10px"
              v-html="$t('projectPage.symlinkHeaderTips')"
            />
            <el-form-item label="" label-width="10px">
              <el-radio-group v-model="formProps.symlink">
                <el-radio :label="false">{{ $t('close') }}</el-radio>
                <el-radio :label="true">{{ $t('open') }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item
              v-show="formProps.symlink"
              :label="$t('directory')"
              prop="symlink_path"
              label-width="80px"
            >
              <el-input
                v-model.trim="formData.symlinkPath"
                autocomplete="off"
              />
            </el-form-item>
            <el-row
              v-show="formProps.symlink"
              style="margin: 0 10px"
              v-html="$t('projectPage.symlinkFooterTips')"
            />
          </el-tab-pane>
          <el-tab-pane name="afterPullScript">
            <template #label>
              {{ $t('projectPage.afterPullScriptLabel') }}
              <el-tooltip class="item" effect="dark" placement="bottom">
                <template #content>
                  <div v-html="$t('projectPage.afterPullScriptTips')" />
                </template>
                <i class="el-icon-question" style="padding-left: 3px" />
              </el-tooltip>
            </template>
            <el-form-item prop="afterPullScript" label-width="0px">
              <el-select
                v-model="formData.afterPullScriptMode"
                :placeholder="$t('projectPage.scriptMode') + '(Default: bash)'"
                style="width: 100%"
              >
                <el-option
                  v-for="(item, index) in scriptLangOption"
                  :key="index"
                  :label="item.label"
                  :value="item.value"
                />
              </el-select>
            </el-form-item>
            <el-form-item prop="afterPullScript" label-width="0px">
              <!-- <span>No support for demo</span> -->
              <v-ace-editor
                v-model:value="formData.afterPullScript"
                :lang="getScriptLang(formData.afterPullScriptMode)"
                theme="github"
                style="height: 400px"
                placeholder="Already switched to project directory..."
              />
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane name="afterDeployScript">
            <template #label>
              {{ $t('projectPage.afterDeployScriptLabel') }}
              <el-tooltip class="item" effect="dark" placement="bottom">
                <template #content>
                  <div v-html="$t('projectPage.afterDeployScriptTips')" />
                </template>
                <i class="el-icon-question" style="padding-left: 3px" />
              </el-tooltip>
            </template>
            <el-form-item prop="afterDeployScript" label-width="0px">
              <el-row type="flex">
                <el-select
                  v-model="formData.afterDeployScriptMode"
                  :placeholder="
                    $t('projectPage.scriptMode') + '(Default: bash)'
                  "
                  style="flex: 1"
                >
                  <el-option
                    v-for="(item, index) in scriptLangOption"
                    :key="index"
                    :label="item.label"
                    :value="item.value"
                  />
                </el-select>
                <el-popover
                  placement="bottom-end"
                  :title="$t('projectPage.predefinedVar')"
                  width="400"
                  trigger="hover"
                >
                  <el-row>
                    <el-row>
                      <span>${PROJECT_NAME}：</span>
                      <span>
                        {{
                          formData.name !== '' ? formData.name : 'project.name'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_PATH}：</span>
                      <span
                        >{{
                          formData.path !== '' ? formData.path : 'project.path'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_SYMLINK_PATH}：</span>
                      <span>
                        {{
                          formData.symlinkPath !== ''
                            ? formData.symlinkPath
                            : 'project.symlink_path'
                        }}
                      </span>
                    </el-row>
                  </el-row>
                  <template #reference>
                    <el-button>
                      {{ $t('projectPage.predefinedVar') }}
                    </el-button>
                  </template>
                </el-popover>
              </el-row>
            </el-form-item>
            <el-form-item prop="afterDeployScript" label-width="0px">
              <!-- <span>No support for demo</span> -->
              <v-ace-editor
                v-model:value="formData.afterDeployScript"
                :lang="getScriptLang(formData.afterDeployScriptMode)"
                theme="github"
                style="height: 400px"
              />
            </el-form-item>
          </el-tab-pane>
        </el-tabs>
      </el-form>
      <template #footer>
        <el-button
          :disabled="formProps.disabled"
          @click="dialogVisible = false"
        >
          {{ $t('cancel') }}
        </el-button>
        <el-button
          :disabled="formProps.disabled"
          type="primary"
          @click="submit"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
    <el-dialog v-model="dialogAutoDeployVisible" :title="$t('setting')">
      <el-form ref="autoDeployForm" :model="autoDeployFormData">
        <el-row style="margin: 10px">{{
          $t('projectPage.autoDeployTitle')
        }}</el-row>
        <el-radio-group
          v-model="autoDeployFormData.autoDeploy"
          style="margin: 10px"
        >
          <el-radio :label="0">{{ $t('close') }}</el-radio>
          <el-radio :label="1">webhook</el-radio>
        </el-radio-group>
        <el-row
          v-show="autoDeployFormData.autoDeploy === 1"
          style="margin: 10px"
          v-html="
            $t('projectPage.autoDeployTips', {
              projectId: autoDeployFormData.id,
            })
          "
        />
      </el-form>
      <template #footer class="dialog-footer">
        <el-button @click="dialogAutoDeployVisible = false">
          {{ $t('cancel') }}
        </el-button>
        <el-button
          :disabled="autoDeployFormProps.disabled"
          type="primary"
          @click="setAutoDeploy"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
    <TheServerDialog
      v-model="dialogServerVisible"
      :project-id="selectedItem.id"
    />
    <TheUserDialog v-model="dialogUserVisible" :project-id="selectedItem.id" />
    <TheFileDialog
      v-model="dialogFileManagerVisible"
      :project-id="selectedItem.id"
    />
  </el-row>
</template>
<script lang="ts">
import tableHeight from '@/mixin/tableHeight'
import { parseGitURL } from '@/utils'
import { NamespaceUserOption } from '@/api/namespace'
import { ServerOption } from '@/api/server'
import {
  ProjectList,
  ProjectTotal,
  ProjectRemoteBranchList,
  ProjectAdd,
  ProjectEdit,
  ProjectRemove,
  ProjectAutoDeploy,
  ProjectData,
} from '@/api/project'
import { role } from '@/utils/namespace'
import { VAceEditor } from 'vue3-ace-editor'
import 'ace-builds/src-noconflict/mode-sh'
import 'ace-builds/src-noconflict/mode-python'
import 'ace-builds/src-noconflict/mode-php'
import 'ace-builds/src-noconflict/theme-github'
import TheServerDialog from './TheServerDialog.vue'
import TheUserDialog from './TheUserDialog.vue'
import TheFileDialog from './TheFileDialog.vue'
import Validator, { RuleItem } from 'async-validator'
import { ElMessageBox, ElMessage } from 'element-plus'
import { defineComponent } from 'vue'
export default defineComponent({
  name: 'ProjectIndex',
  components: {
    VAceEditor,
    TheServerDialog,
    TheUserDialog,
    TheFileDialog,
  },
  mixins: [tableHeight],
  data() {
    return {
      role,
      scriptLangOption: [
        { label: 'sh', value: 'sh', lang: 'sh' },
        { label: 'zsh', value: 'zsh', lang: 'sh' },
        { label: 'bash', value: 'bash', lang: 'sh' },
        { label: 'python', value: 'python', lang: 'python' },
        { label: 'php', value: 'php', lang: 'php' },
      ],
      projectName: '',
      dialogVisible: false,
      dialogAutoDeployVisible: false,
      dialogServerVisible: false,
      dialogUserVisible: false,
      dialogFileManagerVisible: false,
      serverOption: [] as ServerOption['datagram']['list'],
      userOption: [] as NamespaceUserOption['datagram']['list'],
      selectedItem: {},
      tableloading: false,
      tableData: [] as ProjectList['datagram']['list'],
      pagination: {
        page: 1,
        rows: 16,
        total: 0,
      },
      formProps: {
        reviewURLParamOption: [
          {
            label: 'project_id',
            value: 'project_id=__PROJECT_ID__',
          },
          {
            label: 'project_name',
            value: 'project_name=__PROJECT_NAME__',
          },
          {
            label: 'branch',
            value: 'branch=__BRANCH__',
          },
          {
            label: 'environment',
            value: 'environment=__ENVIRONMENT__',
          },
          {
            label: 'commit_id',
            value: 'commit_id=__COMMIT_ID__',
          },
          {
            label: 'publish_time',
            value: 'publish_time=__PUBLISH_TIME__',
          },
          {
            label: 'publisher_id',
            value: 'publisher_id=__PUBLISHER_ID__',
          },
          {
            label: 'publisher_name',
            value: 'publisher_name=__PUBLISHER_NAME__',
          },
          {
            label: 'callback',
            value: 'callback=__CALLBACK__',
            disabled: true,
          },
        ],
        reviewURL: '',
        reviewURLParam: ['callback=__CALLBACK__'],
        symlink: false,
        disabled: false,
        branch: [] as string[],
        lsBranchLoading: false,
        showServers: true,
        showUsers: true,
        tab: 'base',
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
        rsyncOption: '-rtv --exclude .git',
        serverIds: [] as number[],
        userIds: [] as number[],
        review: 0,
        reviewURL: '',
        notifyType: 0,
        notifyTarget: '',
      },
      formRules: {
        name: [{ required: true, message: 'Name required', trigger: ['blur'] }],
        url: [
          {
            required: true,
            message: 'Repository url required',
            trigger: ['blur'],
          },
        ],
        path: [{ required: true, message: 'Path required', trigger: ['blur'] }],
        environment: [
          {
            required: true,
            message: 'Environment required',
            trigger: ['blur'],
          },
        ],
        branch: [
          { required: true, message: 'Branch required', trigger: ['blur'] },
        ],
        notifyTarget: [
          {
            trigger: 'blur',
            validator: (_, value) => {
              if (value !== '' && this.formData.notifyType > 0) {
                return true
              } else if (this.formData.notifyType === 0) {
                return true
              } else {
                return new Error('Select the notice mode')
              }
            },
          } as RuleItem,
        ],
      },
      autoDeployFormProps: {
        disabled: false,
      },
      autoDeployFormData: {
        id: 0,
        autoDeploy: 0,
      },
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

    handleEdit(data: ProjectData['datagram']['detail']) {
      this.formData = Object.assign({}, data)
      this.formProps.symlink = this.formData.symlinkPath !== ''
      this.formProps.showServers = this.formProps.showUsers = false
      this.formProps.branch = []
      this.formProps.reviewURL = ''
      this.formProps.reviewURLParam = []
      if (this.formData.review === 1 && this.formData.reviewURL.length > 0) {
        const url = new URL(this.formData.reviewURL)
        this.formProps.reviewURLParamOption.forEach((item) => {
          if (url.searchParams.has(item.value.split('=')[0])) {
            url.searchParams.delete(item.value.split('=')[0])
            this.formProps.reviewURLParam.push(item.value)
          }
        })
        this.formProps.reviewURL = url.href
      }
      this.dialogVisible = true
    },

    handleCopy(data: ProjectData['datagram']['detail']) {
      this.handleEdit(data)
      this.formData.id = 0
      this.formData.serverIds = []
      this.formData.userIds = []
      this.formProps.showServers = this.formProps.showUsers = true
    },

    handleRemove(data: ProjectData['datagram']['detail']) {
      ElMessageBox.confirm(
        this.$t('projectPage.removeProjectTips', {
          projectName: data.name,
        }),
        this.$t('tips'),
        {
          confirmButtonText: this.$t('confirm'),
          cancelButtonText: this.$t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          this.tableloading = true
          new ProjectRemove({ id: data.id }).request().then(() => {
            ElMessage.success('Success')
            this.getList()
            this.getTotal()
          })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },

    getScriptLang(scriptMode = '') {
      if (scriptMode !== '') {
        const scriptInfo = this.scriptLangOption.find(
          (elem) => elem.value === scriptMode
        )
        return scriptInfo ? scriptInfo['lang'] : ''
      } else {
        return 'sh'
      }
    },

    handleAutoDeploy(data: ProjectData['datagram']['detail']) {
      this.dialogAutoDeployVisible = true
      this.autoDeployFormData.id = data.id
      this.autoDeployFormData.autoDeploy = data.autoDeploy
    },

    handleServer(data: ProjectData['datagram']['detail']) {
      this.selectedItem = data
      this.dialogServerVisible = true
    },

    handleUser(data: ProjectData['datagram']['detail']) {
      this.selectedItem = data
      this.dialogUserVisible = true
    },

    handleFile(data: ProjectData['datagram']['detail']) {
      this.selectedItem = data
      this.dialogFileManagerVisible = true
    },

    submit() {
      ;(this.$refs.form as Validator).validate((valid: boolean) => {
        if (valid) {
          if (this.formProps.symlink === false) {
            this.formData.symlinkPath = ''
          }
          if (
            this.formData.review === 1 &&
            this.formProps.reviewURL.length > 0
          ) {
            const url = new URL(this.formProps.reviewURL)
            this.formProps.reviewURLParam.forEach((param) => {
              const [name, value] = param.split('=')
              url.searchParams.set(name, value)
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
      new ProjectAdd(this.formData)
        .request()
        .then(() => {
          this.dialogVisible = false
          ElMessage.success('Success')
          this.getList()
          this.getTotal()
        })
        .finally(() => {
          this.formProps.disabled = false
        })
    },

    edit() {
      this.formProps.disabled = true
      new ProjectEdit(this.formData)
        .request()
        .then(() => {
          this.dialogVisible = false
          ElMessage.success('Success')
          this.getList()
        })
        .finally(() => {
          this.formProps.disabled = false
        })
    },

    setAutoDeploy() {
      ;(this.$refs.autoDeployForm as Validator).validate((valid: boolean) => {
        if (valid) {
          this.autoDeployFormProps.disabled = true
          new ProjectAutoDeploy(this.autoDeployFormData)
            .request()
            .then(() => {
              this.dialogAutoDeployVisible = false
              ElMessage.success('Success')
              this.getList()
            })
            .finally(() => {
              this.autoDeployFormProps.disabled = false
            })
        } else {
          return false
        }
      })
    },

    getOptions() {
      new ServerOption().request().then((response) => {
        this.serverOption = response.data.list
      })
      new NamespaceUserOption().request().then((response) => {
        this.userOption = response.data.list
      })
    },

    getList() {
      this.tableloading = true
      new ProjectList({ projectName: this.projectName }, this.pagination)
        .request()
        .then((response) => {
          this.tableData = response.data.list
        })
        .finally(() => {
          this.tableloading = false
        })
    },

    getTotal() {
      new ProjectTotal({ projectName: this.projectName })
        .request()
        .then((response) => {
          this.pagination.total = response.data.total
        })
    },

    getRemoteBranchList() {
      if (this.formProps.branch.length > 0) {
        return
      }
      this.formProps.lsBranchLoading = true
      new ProjectRemoteBranchList({ url: this.formData.url })
        .request()
        .then((response) => {
          this.formProps.branch = response.data.branch
          ElMessage.success('Success')
        })
        .finally(() => {
          this.formProps.lsBranchLoading = false
        })
    },

    searchProjectList() {
      this.pagination.page = 1
      this.getList()
      this.getTotal()
    },

    handlePageChange(val = 1) {
      this.pagination.page = val
      this.getList()
    },

    storeFormData() {
      this.tempFormData = JSON.parse(JSON.stringify(this.formData))
    },

    restoreFormData() {
      this.formData = JSON.parse(JSON.stringify(this.tempFormData))
    },
  },
})
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.project-setting-dialog {
  >>> .el-dialog__body {
    padding-top: 10px;
  }
}
.file-dialog {
  >>> .el-dialog__body {
    padding-top: 0px;
    padding-bottom: 0px;
  }
  >>> .el-icon-upload {
    font-size: 14px;
  }
  .file-form {
    height: 520px;
    overflow-y: auto;
    @include scrollBar();
  }
  >>> .CodeMirror {
    height: 500px;
  }
}
</style>
