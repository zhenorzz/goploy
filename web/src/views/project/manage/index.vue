<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between">
      <el-row>
        <el-input
          v-model="searchProject.projectName"
          style="width: 200px"
          placeholder="Filter the project name"
        />
        <el-select
          v-model="searchProject.label"
          :max-collapse-tags="1"
          style="width: 300px"
          multiple
          clearable
          collapse-tags
          collapse-tags-tooltip
          placeholder="Filter the project label"
        >
          <el-option
            v-for="item in labelList"
            :key="item"
            :label="item"
            :value="item"
          />
        </el-select>
      </el-row>

      <el-row>
        <el-button
          :loading="tableLoading"
          type="primary"
          :icon="Refresh"
          @click="refresList"
        />
        <Button
          type="primary"
          :icon="Plus"
          :permissions="[pms.AddProject]"
          @click="handleAdd"
        />
      </el-row>
    </el-row>
    <el-row class="app-table">
      <el-table
        v-loading="tableLoading"
        highlight-current-row
        height="100%"
        :data="tablePage.list"
      >
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="name" :label="$t('name')" width="180">
          <template #default="scope">
            <el-link
              v-if="isLink(scope.row.name)"
              :href="scope.row.name"
              type="primary"
              target="_blank"
            >
              {{ scope.row.name }}
              <el-icon>
                <Link />
              </el-icon>
            </el-link>
            <span v-else>{{ scope.row.name }}</span>
          </template>
        </el-table-column>
        <el-table-column :label="$t('projectURL')" min-width="330">
          <template #default="scope">
            <RepoURL :url="scope.row.url" :text="hideURLPwd(scope.row.url)">
            </RepoURL>
          </template>
        </el-table-column>
        <el-table-column
          prop="path"
          :label="$t('projectPath')"
          min-width="200"
        />
        <el-table-column width="120" :label="$t('environment')" align="center">
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
        <el-table-column width="80" :label="$t('review')" align="center">
          <template #default="scope">
            <span v-if="scope.row.review === 0">{{ $t('close') }}</span>
            <span v-else>{{ $t('open') }}</span>
          </template>
        </el-table-column>
        <el-table-column
          width="110"
          align="center"
          :label="$t('autoDeploy')"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <span v-if="scope.row.autoDeploy === 0">{{ $t('close') }}</span>
            <span v-else>{{ $t('open') }}</span>
            <Button
              type="primary"
              link
              :icon="Edit"
              :permissions="[pms.SwitchProjectWebhook]"
              @click="handleAutoDeploy(scope.row)"
            />
          </template>
        </el-table-column>
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="250"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <Button
              type="primary"
              :icon="Edit"
              :permissions="[pms.EditProject]"
              @click="handleEdit(scope.row)"
            />
            <Button
              color="#626aef"
              :icon="Files"
              :dark="isDark"
              :permissions="[pms.ManageRepository]"
              @click="handleShowFiles(scope.row)"
            />
            <el-tooltip
              class="item"
              effect="dark"
              content="Copy"
              placement="bottom"
            >
              <Button
                type="info"
                :icon="DocumentCopy"
                :permissions="[pms.AddProject]"
                @click="handleCopy(scope.row)"
              />
            </el-tooltip>
            <Button
              type="danger"
              :icon="Delete"
              :permissions="[pms.DeleteProject]"
              @click="handleRemove(scope.row)"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-row>
    <el-row type="flex" justify="end" class="app-page">
      <el-pagination
        :total="tablePage.total"
        :page-size="pagination.rows"
        background
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </el-row>
    <el-dialog
      v-model="dialogVisible"
      :title="$t('setting')"
      width="60%"
      class="project-setting-dialog"
      :fullscreen="$store.state.app.device === 'mobile'"
      :close-on-click-modal="false"
    >
      <el-form
        ref="form"
        v-loading="formProps.disabled"
        :model="formData"
        label-width="120px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-tabs v-model="formProps.tab">
          <el-tab-pane name="base">
            <template #label>
              <span style="vertical-align: middle">
                {{ $t('baseSetting') }}
              </span>
            </template>
            <el-form-item
              :label="$t('name')"
              prop="name"
              :rules="[
                { required: true, message: 'Name required', trigger: ['blur'] },
              ]"
            >
              <el-input
                v-model.trim="formData.name"
                autocomplete="off"
                placeholder="goploy"
              />
            </el-form-item>

            <el-form-item
              prop="url"
              :rules="[
                {
                  required: true,
                  message: 'Repository url required',
                  trigger: ['blur'],
                },
              ]"
            >
              <template #label>
                <span style="vertical-align: middle; padding-right: 4px">
                  {{ $t('projectURL') }}
                </span>
                <el-tooltip placement="top">
                  <template #content>
                    ssh://[username:password@]host.xz[:port]/path/to/repo.git/<br />
                    git@host.xz[:port]/path/to/repo.git/<br />
                    http[s]://[username:password@]host.xz[:port]/path/to/repo.git/<br />
                    svn://host.xz[:port]/path --username= --password=<br />
                    ftp[s]://[username:password@]host.xz[:port]/path/to/repo<br />
                    sftp://host.xz[:port]/path/to/repo --user= --keyFile=
                  </template>
                  <el-icon>
                    <question-filled />
                  </el-icon>
                </el-tooltip>
              </template>
              <el-row type="flex" style="width: 100%">
                <el-select v-model="formData.repoType" style="width: 70px">
                  <el-option label="git" value="git" />
                  <el-option label="svn" value="svn" />
                  <el-option label="ftp" value="ftp" />
                  <el-option label="sftp" value="sftp" />
                </el-select>
                <el-input
                  v-model.trim="formData.url"
                  style="flex: 1"
                  autocomplete="off"
                  placeholder="repository url"
                  @change="formProps.branch = []"
                />
                <el-button
                  :icon="View"
                  type="success"
                  :loading="formProps.pinging"
                  @click="pingRepos"
                >
                  {{ $t('projectPage.testConnection') }}
                </el-button>
              </el-row>
            </el-form-item>
            <el-form-item
              :label="$t('projectPath')"
              prop="path"
              :rules="[
                { required: true, message: 'Path required', trigger: ['blur'] },
              ]"
            >
              <el-input
                v-model.trim="formData.path"
                autocomplete="off"
                placeholder="/var/www/goploy"
                @input="() => handleSymlink(formProps.symlink)"
              />
            </el-form-item>
            <el-form-item
              :label="$t('environment')"
              prop="environment"
              :rules="[
                {
                  required: true,
                  message: 'Environment required',
                  trigger: ['blur'],
                },
              ]"
            >
              <el-select v-model="formData.environment" style="width: 100%">
                <el-option :label="$t('envOption[1]')" :value="1" />
                <el-option :label="$t('envOption[2]')" :value="2" />
                <el-option :label="$t('envOption[3]')" :value="3" />
                <el-option :label="$t('envOption[4]')" :value="4" />
              </el-select>
            </el-form-item>
            <el-form-item
              :label="$t('branch')"
              prop="branch"
              :rules="[
                {
                  required: true,
                  message: 'Branch required',
                  trigger: ['blur'],
                },
              ]"
            >
              <el-row type="flex" style="width: 100%">
                <el-select
                  v-model="formData.branch"
                  filterable
                  allow-create
                  default-first-option
                  style="flex: 1"
                >
                  <el-option
                    v-for="branch in formProps.branch"
                    :key="branch"
                    :label="branch"
                    :value="branch"
                  />
                </el-select>
                <el-button
                  :icon="Search"
                  type="success"
                  :loading="formProps.lsBranchLoading"
                  @click="getRemoteBranchList"
                >
                  {{ $t('projectPage.lishBranch') }}
                </el-button>
              </el-row>
            </el-form-item>
            <el-form-item
              :label="$t('projectPage.transferType')"
              prop="transferType"
            >
              <el-radio-group v-model="formData.transferType">
                <el-radio :label="'rsync'">
                  rsync
                  <el-link
                    :underline="false"
                    :href="$t('projectPage.rsyncDoc')"
                    target="_blank"
                    :icon="QuestionFilled"
                    style="color: #666"
                  />
                </el-radio>
                <el-radio :label="'sftp'">
                  sftp
                  <el-link
                    :underline="false"
                    :href="$t('projectPage.sftpDoc')"
                    target="_blank"
                    :icon="QuestionFilled"
                    style="color: #666"
                  />
                </el-radio>
                <el-radio :label="'custom'">
                  custom
                  <el-link
                    :underline="false"
                    :href="$t('projectPage.customDoc')"
                    target="_blank"
                    :icon="QuestionFilled"
                    style="color: #666"
                  />
                </el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item
              :label="$t('projectPage.transferOption')"
              prop="transferOption"
            >
              <el-input
                v-model="formData.transferOption"
                type="textarea"
                :rows="3"
                autocomplete="off"
                placeholder="-rtv --exclude .git --delete-after"
              />
            </el-form-item>
            <el-form-item prop="notifyTarget" :rules="notifyTargetRules">
              <template #label>
                <el-link
                  type="primary"
                  href="https://docs.goploy.icu/#/dependency/notice"
                  target="_blank"
                >
                  {{ $t('projectPage.deployNotice') }}
                </el-link>
              </template>
              <el-row type="flex" style="width: 100%">
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
            <el-form-item :label="$t('server')" prop="serverIds">
              <el-radio-group
                v-model="formData.deployServerMode"
                style="margin-bottom: 5px"
              >
                <el-radio label="parallel">{{ $t('parallel') }}</el-radio>
                <el-radio label="serial">{{ $t('serial') }}</el-radio>
              </el-radio-group>
              <el-select
                v-model="formData.serverIds"
                multiple
                filterable
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
            <el-form-item :label="$t('user')" prop="userIds">
              <el-select
                v-model="formData.userIds"
                multiple
                filterable
                style="width: 100%"
              >
                <el-option
                  v-for="(item, index) in userOption.filter(
                    (item) => item.roleId > 0
                  )"
                  :key="index"
                  :label="item.userName"
                  :value="item.userId"
                />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('tag')" prop="tag">
              <el-select
                v-model="formProps.label"
                style="width: 100%"
                :max-collapse-tags="5"
                allow-create
                :reserve-keyword="false"
                collapse-tags-tooltip
                multiple
                clearable
                filterable
                default-first-option
              >
                <el-option
                  v-for="item in labelList"
                  :key="item"
                  :label="item"
                  :value="item"
                />
              </el-select>
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane name="review">
            <template #label>
              <span style="vertical-align: middle">
                {{ $t('projectPage.publishReview') }}
              </span>
            </template>
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
              style="margin: 0 10px; white-space: pre-wrap"
            >
              {{ $t('projectPage.reviewFooterTips') }}
            </el-row>
          </el-tab-pane>
          <el-tab-pane name="symlink">
            <template #label>
              <span style="vertical-align: middle">
                {{ $t('projectPage.symlinkLabel') }}
              </span>
            </template>
            <el-row style="margin: 0 10px 18px; white-space: pre-line">
              {{ $t('projectPage.symlinkHeaderTips') }}
            </el-row>
            <el-form-item label="Symlink" label-width="120px">
              <el-radio-group
                v-model="formProps.symlink"
                @change="handleSymlink"
              >
                <el-radio :label="false">{{ $t('close') }}</el-radio>
                <el-radio :label="true">
                  {{ $t('open') }}
                </el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item
              v-show="formProps.symlink"
              label="Backup number"
              label-width="120px"
            >
              <el-input-number
                v-model="formData.symlinkBackupNumber"
                :min="1"
              />
            </el-form-item>
            <el-form-item
              v-show="formProps.symlink"
              label="Directory"
              prop="symlink_path"
              label-width="120px"
            >
              <el-input v-model="formData.symlinkPath" readonly disabled>
                <template #append>/uuid-version</template>
              </el-input>
              <div class="el-form-item__error">
                (path =
                project.path.dirname/goploy-symlink/project.path.basename/uuid-version)
              </div>
            </el-form-item>
            <el-row
              v-show="formProps.symlink"
              style="margin: 5px 10px 0; white-space: pre-line"
            >
              {{ $t('projectPage.symlinkFooterTips') }}
            </el-row>
          </el-tab-pane>
          <el-tab-pane name="afterPullScript">
            <template #label>
              <span style="vertical-align: middle; padding-right: 4px">
                {{ $t('projectPage.afterPullScriptLabel') }}
              </span>
              <el-tooltip class="item" effect="dark" placement="bottom">
                <template #content>
                  <div style="white-space: pre-line">
                    {{ $t('projectPage.afterPullScriptTips') }}
                  </div>
                </template>
                <el-icon style="vertical-align: middle" :size="16">
                  <question-filled />
                </el-icon>
              </el-tooltip>
            </template>
            <el-form-item prop="afterPullScript" label-width="0px">
              <el-row type="flex" style="width: 100%">
                <el-select
                  v-model="formData.script.afterPull.mode"
                  :placeholder="
                    $t('projectPage.scriptMode') + '(Default: bash)'
                  "
                  style="flex: 1"
                  @change="
                    (mode) => {
                      handleScriptModeChange(ScriptKey.AfterPull, mode)
                    }
                  "
                >
                  <el-option
                    v-for="(item, index) in scriptLang.Option"
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
                  <div>
                    <el-row>
                      <span>${PROJECT_ID}: </span>
                      <span>{{
                        formData.id > 0 ? formData.id : 'project.id'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_NAME}: </span>
                      <span>
                        {{
                          formData.name !== '' ? formData.name : 'project.name'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_PATH}: </span>
                      <span>
                        {{
                          formData.path !== '' ? formData.path : 'project.path'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_SYMLINK_PATH}: </span>
                      <span>
                        {{
                          formProps.symlink === true
                            ? formData.symlinkPath
                            : 'project.symlink_path'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_BRANCH}: </span>
                      <span>{{
                        formData.branch !== ''
                          ? formData.branch
                          : 'project.branch'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_ENV}: </span>
                      <span>{{
                        formData.environment > 0
                          ? formData.environment
                          : 'project.environment'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_LABEL}: </span>
                      <span>{{
                        formData.label !== '' ? formData.label : 'project.label'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_SERVERS}: </span>
                      <span>JSON format</span>
                    </el-row>
                    <el-row>
                      <span>${REPOSITORY_TYPE}: </span>
                      <span>{{
                        formData.repoType !== ''
                          ? formData.repoType
                          : 'project.repoType'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${REPOSITORY_URL}: </span>
                      <span>project.url</span>
                    </el-row>
                    <el-row>
                      <span>${REPOSITORY_PATH}: </span>
                      <span>project.repo_path</span>
                    </el-row>
                    <el-row>
                      <span>${PUBLISH_TOKEN}: </span>
                      <span>project.last_publish_token</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_ID}: </span>
                      <span>Commit ID</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_SHORT_ID}: </span>
                      <span>Commit ID (6 char)</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_BRANCH}: </span>
                      <span>Commit branch</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_TAG}: </span>
                      <span>Commit tag</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_AUTHOR}: </span>
                      <span>Commit author</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_TIMESTAMP}: </span>
                      <span>Commit timestamp(second)</span>
                    </el-row>
                  </div>
                  <template #reference>
                    <el-button>
                      {{ $t('projectPage.predefinedVar') }}
                    </el-button>
                  </template>
                </el-popover>
              </el-row>
            </el-form-item>
            <el-form-item prop="afterPullScript" label-width="0px">
              <v-ace-editor
                v-model:value="formData.script.afterPull.content"
                :lang="scriptLang.getScriptLang(formData.script.afterPull.mode)"
                :theme="isDark ? 'one_dark' : 'github'"
                style="height: 400px; width: 100%"
                placeholder="Already switched to project directory..."
                :options="{
                  newLineMode:
                    formData.script.afterPull.mode === 'cmd'
                      ? 'windows'
                      : 'unix',
                }"
              />
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane name="afterDeployScript">
            <template #label>
              <span style="vertical-align: middle; padding-right: 4px">
                {{ $t('projectPage.afterDeployScriptLabel') }}
              </span>
              <el-tooltip class="item" effect="dark" placement="bottom">
                <template #content>
                  <div style="white-space: pre-line">
                    {{ $t('projectPage.afterDeployScriptTips') }}
                  </div>
                </template>
                <el-icon style="vertical-align: middle" :size="16">
                  <question-filled />
                </el-icon>
              </el-tooltip>
            </template>
            <el-form-item prop="afterDeployScript" label-width="0px">
              <el-row type="flex" style="width: 100%">
                <el-select
                  v-model="formData.script.afterDeploy.mode"
                  :placeholder="
                    $t('projectPage.scriptMode') + '(Default: bash)'
                  "
                  style="flex: 1"
                  @change="
                    (mode) => {
                      handleScriptModeChange(ScriptKey.AfterDeploy, mode)
                    }
                  "
                >
                  <el-option
                    v-for="(item, index) in scriptLang.Option"
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
                  <div>
                    <el-row>
                      <span>${PROJECT_ID}: </span>
                      <span>{{
                        formData.id > 0 ? formData.id : 'project.id'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_NAME}: </span>
                      <span>
                        {{
                          formData.name !== '' ? formData.name : 'project.name'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_PATH}: </span>
                      <span>
                        {{
                          formData.path !== '' ? formData.path : 'project.path'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_SYMLINK_PATH}: </span>
                      <span>
                        {{
                          formProps.symlink === true
                            ? formData.symlinkPath
                            : 'project.symlink_path'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_BRANCH}: </span>
                      <span>{{
                        formData.branch !== ''
                          ? formData.branch
                          : 'project.branch'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_ENV}: </span>
                      <span>{{
                        formData.environment > 0
                          ? formData.environment
                          : 'project.environment'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_LABEL}: </span>
                      <span>{{
                        formData.label !== '' ? formData.label : 'project.label'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${REPOSITORY_TYPE}: </span>
                      <span>{{
                        formData.repoType !== ''
                          ? formData.repoType
                          : 'project.repoType'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${REPOSITORY_URL}: </span>
                      <span>project.url</span>
                    </el-row>
                    <el-row>
                      <span>${REPOSITORY_PATH}: </span>
                      <span>project.repo_path</span>
                    </el-row>
                    <el-row>
                      <span>${PUBLISH_TOKEN}: </span>
                      <span>project.last_publish_token</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_TOTAL_NUMBER}: </span>
                      <span>{{
                        formData.serverIds.length > 0
                          ? formData.serverIds.length
                          : 'server.total.number'
                      }}</span>
                    </el-row>
                    <el-row v-if="formData.deployServerMode === 'serial'">
                      <span>${SERVER_SERIAL_NUMBER}: </span>
                      <span>start from 1</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_ID}: </span>
                      <span>server.id</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_NAME}: </span>
                      <span>server.name</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_IP}: </span>
                      <span>server.ip</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_PORT}: </span>
                      <span>server.port</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_OWNER}: </span>
                      <span>server.owner</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_PASSWORD}: </span>
                      <span>server.password</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_PATH}: </span>
                      <span>server.path</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_JUMP_IP}: </span>
                      <span>server.jump_ip</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_JUMP_PORT}: </span>
                      <span>server.jump_port</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_JUMP_OWNER}: </span>
                      <span>server.jump_owner</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_JUMP_PASSWORD}: </span>
                      <span>server.jump_password</span>
                    </el-row>
                    <el-row>
                      <span>${SERVER_JUMP_PATH}: </span>
                      <span>server.jump_path</span>
                    </el-row>
                  </div>
                  <template #reference>
                    <el-button>
                      {{ $t('projectPage.predefinedVar') }}
                    </el-button>
                  </template>
                </el-popover>
              </el-row>
            </el-form-item>
            <el-form-item prop="afterDeployScript" label-width="0px">
              <v-ace-editor
                v-model:value="formData.script.afterDeploy.content"
                :lang="
                  scriptLang.getScriptLang(formData.script.afterDeploy.mode)
                "
                :theme="isDark ? 'one_dark' : 'github'"
                style="height: 400px; width: 100%"
                :options="{
                  newLineMode:
                    formData.script.afterDeploy.mode === 'cmd'
                      ? 'windows'
                      : 'unix',
                }"
              />
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane name="deployFinishScript">
            <template #label>
              <span style="vertical-align: middle; padding-right: 4px">
                {{ $t('projectPage.deployFinishScriptLabel') }}
              </span>
              <el-tooltip class="item" effect="dark" placement="bottom">
                <template #content>
                  <div style="white-space: pre-line">
                    {{ $t('projectPage.deployFinishScriptTips') }}
                  </div>
                </template>
                <el-icon style="vertical-align: middle" :size="16">
                  <question-filled />
                </el-icon>
              </el-tooltip>
            </template>
            <el-form-item prop="deployFinishScript" label-width="0px">
              <el-row type="flex" style="width: 100%">
                <el-select
                  v-model="formData.script.deployFinish.mode"
                  :placeholder="
                    $t('projectPage.scriptMode') + '(Default: bash)'
                  "
                  style="flex: 1"
                  @change="
                    (mode) => {
                      handleScriptModeChange(ScriptKey.DeployFinish, mode)
                    }
                  "
                >
                  <el-option
                    v-for="(item, index) in scriptLang.Option"
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
                  <div>
                    <el-row>
                      <span>${PROJECT_ID}: </span>
                      <span>{{
                        formData.id > 0 ? formData.id : 'project.id'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_NAME}: </span>
                      <span>
                        {{
                          formData.name !== '' ? formData.name : 'project.name'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_PATH}: </span>
                      <span>
                        {{
                          formData.path !== '' ? formData.path : 'project.path'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_SYMLINK_PATH}: </span>
                      <span>
                        {{
                          formProps.symlink === true
                            ? formData.symlinkPath
                            : 'project.symlink_path'
                        }}
                      </span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_BRANCH}: </span>
                      <span>{{
                        formData.branch !== ''
                          ? formData.branch
                          : 'project.branch'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_ENV}: </span>
                      <span>{{
                        formData.environment > 0
                          ? formData.environment
                          : 'project.environment'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_LABEL}: </span>
                      <span>{{
                        formData.label !== '' ? formData.label : 'project.label'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${PROJECT_SERVERS}: </span>
                      <span>JSON format</span>
                    </el-row>
                    <el-row>
                      <span>${REPOSITORY_TYPE}: </span>
                      <span>{{
                        formData.repoType !== ''
                          ? formData.repoType
                          : 'project.repoType'
                      }}</span>
                    </el-row>
                    <el-row>
                      <span>${REPOSITORY_URL}: </span>
                      <span>project.url</span>
                    </el-row>
                    <el-row>
                      <span>${REPOSITORY_PATH}: </span>
                      <span>project.repo_path</span>
                    </el-row>
                    <el-row>
                      <span>${PUBLISH_TOKEN}: </span>
                      <span>project.last_publish_token</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_ID}: </span>
                      <span>Commit ID</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_SHORT_ID}: </span>
                      <span>Commit ID (6 char)</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_BRANCH}: </span>
                      <span>Commit branch</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_TAG}: </span>
                      <span>Commit tag</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_AUTHOR}: </span>
                      <span>Commit author</span>
                    </el-row>
                    <el-row>
                      <span>${COMMIT_TIMESTAMP}: </span>
                      <span>Commit timestamp(second)</span>
                    </el-row>
                  </div>
                  <template #reference>
                    <el-button>
                      {{ $t('projectPage.predefinedVar') }}
                    </el-button>
                  </template>
                </el-popover>
              </el-row>
            </el-form-item>
            <el-form-item prop="deployFinishScript" label-width="0px">
              <v-ace-editor
                v-model:value="formData.script.deployFinish.content"
                :lang="
                  scriptLang.getScriptLang(formData.script.deployFinish.mode)
                "
                :theme="isDark ? 'one_dark' : 'github'"
                style="height: 400px; width: 100%"
                placeholder="Already switched to project directory..."
                :options="{
                  newLineMode:
                    formData.script.deployFinish.mode === 'cmd'
                      ? 'windows'
                      : 'unix',
                }"
              />
            </el-form-item>
          </el-tab-pane>
          <el-tab-pane name="customVariable">
            <template #label>
              <span style="vertical-align: middle; padding-right: 4px">
                {{ $t('projectPage.customVariableLabel') }}
              </span>
              <el-tooltip class="item" effect="dark" placement="bottom">
                <template #content>
                  <div style="white-space: pre-line">
                    {{ $t('projectPage.customVariableTips') }}
                  </div>
                </template>
                <el-icon style="vertical-align: middle" :size="16">
                  <question-filled />
                </el-icon>
              </el-tooltip>
            </template>
            <el-row type="flex" justify="end" style="margin-bottom: 10px">
              <el-button
                type="primary"
                :icon="Plus"
                @click="handleCusotmVariableAdd"
              />
            </el-row>
            <el-form-item prop="customVariableScript" label-width="0px">
              <el-row
                v-for="(variable, index) in formData.script.customVariables"
                :key="index"
                type="flex"
                style="width: 100%"
              >
                <el-input
                  v-model.trim="variable.name"
                  style="flex: 1"
                  autocomplete="off"
                  placeholder="variable name"
                >
                  <template #prepend>${</template>
                  <template #append>}</template>
                </el-input>
                <el-select v-model="variable.type" style="width: 200px">
                  <el-option
                    v-for="item in formProps.customVariablesType"
                    :key="item"
                    :label="item"
                    :value="item"
                  />
                </el-select>
                <el-input
                  v-model.trim="variable.value"
                  style="flex: 1"
                  autocomplete="off"
                  :placeholder="
                    variable.type == 'list'
                      ? 'value split by comma'
                      : 'default value'
                  "
                />
                <el-button
                  type="warning"
                  :icon="Minus"
                  plain
                  @click="handleCusotmVariableDelete(index)"
                />
              </el-row>
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
          :loading="formProps.disabled"
          type="primary"
          @click="submit"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
    <el-dialog
      v-model="dialogAutoDeployVisible"
      :title="$t('setting')"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <el-form ref="autoDeployForm" :model="autoDeployFormData">
        <el-row style="margin: 10px">
          {{ $t('projectPage.autoDeployTitle') }}
        </el-row>
        <el-radio-group
          v-model="autoDeployFormData.autoDeploy"
          style="margin: 10px"
        >
          <el-radio :label="0">{{ $t('close') }}</el-radio>
          <el-radio :label="1">webhook</el-radio>
        </el-radio-group>
        <el-row
          v-show="autoDeployFormData.autoDeploy === 1"
          style="margin: 10px; white-space: pre-line"
        >
          <span v-if="selectedItem.repoType === 'svn'">
            {{
              $t('projectPage.autoDeploySVNTips', {
                projectId: autoDeployFormData.id,
                ref: `{\\"ref\\": \\"${selectedItem.branch}\\"}`,
              })
            }}
          </span>
          <span v-else>
            <el-row>
              <RepoURL
                style="font-size: 14px"
                :url="selectedItem.url"
                suffix="/-/settings/integrations"
                text="Gitlab Webhook"
              />、
              <RepoURL
                style="font-size: 14px"
                :url="selectedItem.url"
                suffix="/settings/hooks"
                text="Github Webhook"
              />、
              <RepoURL
                style="font-size: 14px"
                :url="selectedItem.url"
                suffix="/hooks"
                text="Gitee Webhook"
              />
            </el-row>
            {{
              $t('projectPage.autoDeployGitTips', {
                projectId: autoDeployFormData.id,
              })
            }}
          </span>
        </el-row>
      </el-form>
      <template #footer>
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
    <el-dialog
      v-model="dialogFilesVisible"
      :title="$t('projectPage.repositoryDialogTitle')"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <explorer :project="selectedItem"></explorer>
    </el-dialog>
  </el-row>
</template>
<script lang="ts">
export default { name: 'ProjectIndex' }
</script>
<script lang="ts" setup>
import explorer from './components/Explorer.vue'
import pms from '@/permission'
import { scriptLang } from '@/const/const'
import Button from '@/components/Permission/Button.vue'
import {
  Search,
  Refresh,
  View,
  Link,
  Plus,
  Minus,
  Edit,
  QuestionFilled,
  Files,
  DocumentCopy,
  Delete,
} from '@element-plus/icons-vue'
import { VAceEditor } from 'vue3-ace-editor'
import * as ace from 'ace-builds/src-noconflict/ace'
import path from 'path-browserify'
import RepoURL from '@/components/RepoURL/index.vue'
import { isLink, hideURLPwd } from '@/utils'
import { NamespaceUserOption } from '@/api/namespace'
import { ServerOption } from '@/api/server'
import {
  ProjectList,
  ProjectPingRepos,
  ProjectRemoteBranchList,
  ProjectAdd,
  ProjectEdit,
  ProjectUserList,
  ProjectServerList,
  ProjectRemove,
  ProjectAutoDeploy,
  ProjectData,
  LabelList,
} from '@/api/project'
import type { ElRadioGroup, ElForm, FormItemRule } from 'element-plus'
import { ref, computed } from 'vue'
import { useDark } from '@vueuse/core'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const isDark = useDark()

ace.config.set(
  'basePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
ace.config.set(
  'themePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)

const searchProject = ref({
  projectName: '',
  label: [] as string[],
})
const dialogVisible = ref(false)
const dialogAutoDeployVisible = ref(false)
const dialogFilesVisible = ref(false)
const serverOption = ref<ServerOption['datagram']['list']>([])
const userOption = ref<NamespaceUserOption['datagram']['list']>([])
const selectedItem = ref({} as ProjectData)
const tableLoading = ref(false)
const tableData = ref<ProjectList['datagram']['list']>([])
const labelList = ref<LabelList['datagram']['list']>([])
const pagination = ref({ page: 1, rows: 20 })
const form = ref<InstanceType<typeof ElForm>>()
const formProps = ref({
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
  label: [] as string[],
  branch: [] as string[],
  pinging: false,
  lsBranchLoading: false,
  tab: 'base',
  customVariablesType: ['string', 'list'],
})
const tempFormData = {
  id: 0,
  name: '',
  label: '',
  repoType: 'git',
  url: '',
  path: '',
  symlinkPath: '',
  symlinkBackupNumber: 10,
  script: {
    afterPull: { mode: '', content: '' },
    afterDeploy: { mode: '', content: '' },
    deployFinish: { mode: '', content: '' },
    customVariables: [] as { name: string; value: string; type: string }[],
  },
  environment: 1,
  branch: '',
  transferType: 'rsync',
  transferOption: '-rtv --exclude .git',
  deployServerMode: 'parallel',
  serverIds: [] as number[],
  userIds: [] as number[],
  review: 0,
  reviewURL: '',
  notifyType: 0,
  notifyTarget: '',
}
const formData = ref(tempFormData)
const notifyTargetRules: FormItemRule[] = [
  {
    trigger: 'blur',
    validator: (_, value) => {
      if (value !== '' && formData.value.notifyType > 0) {
        return true
      } else if (formData.value.notifyType === 0) {
        return true
      } else {
        return new Error('Select the notice mode')
      }
    },
  },
]
const autoDeployForm = ref<InstanceType<typeof ElForm>>()
const autoDeployFormProps = ref({ disabled: false })
const autoDeployFormData = ref({ id: 0, autoDeploy: 0 })

getOptions()
getList()
getTagList()

const tablePage = computed(() => {
  let _tableData = tableData.value
  if (searchProject.value.projectName !== '') {
    _tableData = tableData.value.filter(
      (item) => item.name.indexOf(searchProject.value.projectName) !== -1
    )
  }
  if (searchProject.value.label.length > 0) {
    _tableData = _tableData.filter((item) =>
      item.label
        .split(',')
        .find((p) => searchProject.value.label.indexOf(p) !== -1)
    )
  }
  return {
    list: _tableData.slice(
      (pagination.value.page - 1) * pagination.value.rows,
      pagination.value.page * pagination.value.rows
    ),
    total: _tableData.length,
  }
})

function getOptions() {
  new ServerOption().request().then((response) => {
    serverOption.value = response.data.list
  })
  new NamespaceUserOption().request().then((response) => {
    userOption.value = response.data.list
  })
}

function getList() {
  tableLoading.value = true
  new ProjectList()
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function getTagList() {
  new LabelList().request().then((response) => {
    labelList.value = response.data.list
  })
}

function handleAdd() {
  if (formData.value.id > 0) {
    restoreFormData()
  }
  formProps.value.symlink = false
  dialogVisible.value = true
}

function handleShowFiles(data: ProjectData) {
  selectedItem.value = data
  dialogFilesVisible.value = true
}

function handleEdit(data: ProjectData) {
  formData.value = Object.assign({}, data)
  formProps.value.symlink = formData.value.symlinkPath !== ''
  formProps.value.branch = []
  formProps.value.reviewURL = ''
  formProps.value.disabled = true
  formData.value.userIds = []
  formData.value.serverIds = []
  formProps.value.label = data.label != '' ? data.label.split(',') : []
  Promise.all([
    new ProjectUserList({ id: data.id }).request(),
    new ProjectServerList({ id: data.id }).request(),
  ]).then((values) => {
    const projectUserList = values[0].data.list
    const ProjectServerList = values[1].data.list
    formData.value.userIds = projectUserList.map((item) => item.userId)
    formData.value.serverIds = ProjectServerList.map((item) => item.serverId)
    formProps.value.disabled = false
  })

  if (formData.value.review === 1 && formData.value.reviewURL.length > 0) {
    formProps.value.reviewURLParam = []
    const url = new URL(formData.value.reviewURL)
    formProps.value.reviewURLParamOption.forEach((item) => {
      if (url.searchParams.has(item.value.split('=')[0])) {
        url.searchParams.delete(item.value.split('=')[0])
        formProps.value.reviewURLParam.push(item.value)
      }
    })
    formProps.value.reviewURL = url.href
  }
  dialogVisible.value = true
}

function handleCopy(data: ProjectData) {
  handleEdit(data)
  formData.value.id = 0
}

function handleRemove(data: ProjectData) {
  ElMessageBox.confirm(t('removeTips', { name: data.name }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      tableLoading.value = true
      new ProjectRemove({ id: data.id }).request().then(() => {
        ElMessage.success('Success')
        getList()
        getTagList()
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function getSymlinkPath(projectPath: string) {
  return path.normalize(
    path.dirname(projectPath) + '/goploy-symlink/' + path.basename(projectPath)
  )
}

const handleSymlink: InstanceType<typeof ElRadioGroup>['onChange'] = (
  value
) => {
  if (value) {
    formData.value.symlinkPath = getSymlinkPath(formData.value.path)
  } else {
    formData.value.symlinkPath = ''
  }
}

function handleAutoDeploy(data: ProjectData) {
  dialogAutoDeployVisible.value = true
  selectedItem.value = data
  autoDeployFormData.value.id = data.id
  autoDeployFormData.value.autoDeploy = data.autoDeploy
}

function handleCusotmVariableAdd() {
  if (!formData.value.script.customVariables) {
    formData.value.script.customVariables = []
  }
  formData.value.script.customVariables.push({
    name: '',
    value: '',
    type: 'string',
  })
}

function handleCusotmVariableDelete(index: number) {
  formData.value.script.customVariables.splice(index, 1)
}

enum ScriptKey {
  AfterPull = 'afterPull',
  AfterDeploy = 'afterDeploy',
  DeployFinish = 'deployFinish',
}

function handleScriptModeChange(type: ScriptKey, mode: string) {
  if (mode === 'cmd') {
    if (
      !formData.value.script[type].content.includes('\r\n') &&
      formData.value.script[type].content.includes('\n')
    ) {
      formData.value.script[type].content = ''
    }
  } else {
    if (formData.value.script[type].content.includes('\r\n')) {
      formData.value.script[type].content = ''
    }
  }
}

function submit() {
  form.value?.validate(async (valid) => {
    if (!valid) {
      return false
    }
    formProps.value.disabled = true
    try {
      await new ProjectPingRepos({
        repoType: formData.value.repoType,
        url: formData.value.url,
      }).request()
    } catch (error) {
      formProps.value.disabled = false
    }

    if (formData.value.review === 1 && formProps.value.reviewURL.length > 0) {
      const url = new URL(formProps.value.reviewURL)
      formProps.value.reviewURLParam.forEach((param) => {
        const [name, value] = param.split('=')
        url.searchParams.set(name, value)
      })
      formData.value.reviewURL = url.href
    } else {
      formData.value.reviewURL = ''
    }
    if (
      formProps.value.label.filter((p) => String(p).includes(',')).length > 0
    ) {
      ElMessage.error('Tag is not allowed to contain , ')
      return false
    }
    ;(formData.value.id === 0
      ? new ProjectAdd({
          ...formData.value,
          label: formProps.value.label.join(','),
        })
      : new ProjectEdit({
          ...formData.value,
          label: formProps.value.label.join(','),
        })
    )
      .request()
      .then(() => {
        dialogVisible.value = false
        ElMessage.success('Success')
        getList()
        getTagList()
      })
      .finally(() => {
        formProps.value.disabled = false
      })

    return true
  })
}

function setAutoDeploy() {
  autoDeployForm.value?.validate((valid) => {
    if (valid) {
      autoDeployFormProps.value.disabled = true
      new ProjectAutoDeploy(autoDeployFormData.value)
        .request()
        .then(() => {
          dialogAutoDeployVisible.value = false
          ElMessage.success('Success')
          getList()
        })
        .finally(() => {
          autoDeployFormProps.value.disabled = false
        })
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
}

function pingRepos() {
  if (formData.value.url === '') {
    ElMessage.error('url can not be blank')
    formProps.value.branch = []
    return
  }
  formProps.value.pinging = true
  new ProjectPingRepos({
    repoType: formData.value.repoType,
    url: formData.value.url,
  })
    .request()
    .then(() => {
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.pinging = false
    })
}

function getRemoteBranchList() {
  if (formData.value.url === '') {
    ElMessage.error('url can not be blank')
    formProps.value.branch = []
    return
  }

  if (formProps.value.branch.length > 0) {
    return
  }
  formProps.value.lsBranchLoading = true
  new ProjectRemoteBranchList({
    repoType: formData.value.repoType,
    url: formData.value.url,
  })
    .request()
    .then((response) => {
      formProps.value.branch = response.data.branch
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.lsBranchLoading = false
    })
}

function refresList() {
  searchProject.value = {
    projectName: '',
    label: [],
  }

  pagination.value.page = 1
  getList()
  getOptions()
}

function handlePageChange(page = 1) {
  pagination.value.page = page
}

function restoreFormData() {
  formData.value = { ...tempFormData }
}
</script>
