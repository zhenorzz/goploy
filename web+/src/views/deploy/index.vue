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
      style="width: 100%; margin-top: 5px"
    >
      <el-table-column prop="id" label="ID" width="80" align="center" />
      <el-table-column
        prop="name"
        :label="$t('name')"
        min-width="150"
        align="center"
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
        width="230"
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
      />
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="310"
        fixed="right"
        align="center"
      >
        <template #default="scope">
          <el-row class="operation-btn">
            <el-button
              v-if="scope.row.deployState === 0"
              type="primary"
              @click="publish(scope.row)"
              >{{ $t('initial') }}</el-button
            >
            <el-button
              v-else-if="hasManagerPermission() && scope.row.deployState === 1"
              type="primary"
              @click="resetState(scope.row)"
              >{{ $t('deployPage.resetState') }}</el-button
            >
            <el-dropdown
              v-else-if="hasGroupManagerPermission() || scope.row.review === 0"
              split-button
              trigger="click"
              type="primary"
              @click="publish(scope.row)"
              @command="(commandFunc) => commandFunc(scope.row)"
            >
              {{
                isMember() && scope.row.review === 1
                  ? $t('submit')
                  : $t('deploy')
              }}
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item :command="handleCommitCommand">
                    Commit list
                  </el-dropdown-item>
                  <el-dropdown-item :command="getTagList">
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
              v-if="hasGroupManagerPermission() || scope.row.review === 1"
              trigger="click"
              @command="(commandFunc) => commandFunc(scope.row)"
            >
              <el-button type="warning">
                {{ $t('func') }}<i class="el-icon-arrow-down el-icon--right" />
              </el-button>
              <template #dropdown>
                <el-dropdown-menu style="min-width: 84px; text-align: center">
                  <el-dropdown-item
                    v-if="hasGroupManagerPermission()"
                    :command="handleTaskCommand"
                  >
                    {{ $t('deployPage.taskDeploy') }}
                  </el-dropdown-item>
                  <el-dropdown-item
                    v-if="scope.row.review === 1"
                    :command="handleReviewCommand"
                  >
                    {{ $t('deployPage.reviewDeploy') }}
                  </el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
            <el-button type="success" @click="handleDetail(scope.row)">
              {{ $t('detail') }}
            </el-button>
          </el-row>
        </template>
      </el-table-column>
    </el-table>
    <el-row type="flex" justify="end" style="margin-top: 10px">
      <el-pagination
        v-model:current-page="pagination.page"
        :total="pagination.total"
        :page-size="pagination.rows"
        background
        :page-sizes="[20, 50, 100]"
        layout="sizes, total, prev, pager, next, jumper"
        @size-change="handleSizeChange"
        @current-change="handlePageChange"
      />
    </el-row>
    <el-dialog
      v-model="dialogVisible"
      :title="$t('detail')"
      custom-class="publish-record"
    >
      <el-row type="flex">
        <el-row v-loading="searchPreview.loading" class="publish-preview">
          <el-row>
            <el-popover placement="bottom-start" width="318" trigger="click">
              <el-row type="flex" align="middle">
                <label class="publish-filter-label">{{ $t('user') }}</label>
                <el-select
                  v-model="searchPreview.userId"
                  style="flex: 1"
                  clearable
                >
                  <el-option
                    v-for="(item, index) in userOption"
                    :key="index"
                    :label="item.userName"
                    :value="item.userId"
                  />
                </el-select>
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px">
                <label class="publish-filter-label">Commit</label>
                <el-input
                  v-model.trim="searchPreview.commit"
                  autocomplete="off"
                  style="flex: 1"
                  placeholder="Commit"
                />
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px">
                <label class="publish-filter-label">{{ $t('branch') }}</label>
                <el-input
                  v-model.trim="searchPreview.branch"
                  autocomplete="off"
                  style="flex: 1"
                  :placeholder="$t('branch')"
                />
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px">
                <label class="publish-filter-label">{{ $t('filename') }}</label>
                <el-input
                  v-model.trim="searchPreview.filename"
                  autocomplete="off"
                  style="flex: 1"
                  :placeholder="$t('filename')"
                />
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px">
                <label class="publish-filter-label">{{ $t('state') }}</label>
                <el-select
                  v-model="searchPreview.state"
                  style="flex: 1"
                  clearable
                >
                  <el-option :label="$t('success')" :value="1" />
                  <el-option :label="$t('fail')" :value="0" />
                </el-select>
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px">
                <label class="publish-filter-label">{{
                  $t('commitDate')
                }}</label>
                <el-date-picker
                  v-model="searchPreview.commitDate"
                  class="dmp-date-picker"
                  popper-class="dmp-date-picker-popper"
                  :picker-options="pickerOptions"
                  type="daterange"
                  value-format="yyyy-MM-dd HH:mm:ss"
                  range-separator="-"
                  :start-placeholder="$t('startDate')"
                  :end-placeholder="$t('endDate')"
                  :default-time="['00:00:00', '23:59:59']"
                  style="flex: 1"
                  align="center"
                />
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px">
                <label class="publish-filter-label">{{
                  $t('deployDate')
                }}</label>
                <el-date-picker
                  v-model="searchPreview.deployDate"
                  class="dmp-date-picker"
                  popper-class="dmp-date-picker-popper"
                  :picker-options="pickerOptions"
                  type="daterange"
                  value-format="yyyy-MM-dd HH:mm:ss"
                  range-separator="-"
                  :start-placeholder="$t('startDate')"
                  :end-placeholder="$t('endDate')"
                  :default-time="['00:00:00', '23:59:59']"
                  style="flex: 1"
                  align="center"
                />
              </el-row>
              <template #reference>
                <el-button icon="el-icon-notebook-2" style="width: 220px">
                  Filter({{ previewFilterlength }})
                </el-button>
              </template>
            </el-popover>
            <el-button
              type="warning"
              icon="el-icon-refresh"
              @click="refreshSearchPreviewCondition"
            />
            <el-button
              :loading="searchPreview.loading"
              type="primary"
              icon="el-icon-search"
              style="margin-left: 2px"
              @click="searchPreviewList"
            />
          </el-row>
          <el-radio-group v-model="publishToken" @change="handleTraceChange">
            <el-row v-for="(item, index) in gitTraceList" :key="index">
              <el-row style="margin: 5px 0">
                <el-radio class="publish-commit" :label="item.token" border>
                  <span class="publish-name">{{ item.publisherName }}</span>
                  <span class="publish-commitID">
                    commitID: {{ item.commit }}
                  </span>
                  <i
                    v-if="item.publishState === 1"
                    class="el-icon-check icon-success"
                    style="float: right"
                  />
                  <i
                    v-else
                    class="el-icon-close icon-fail"
                    style="float: right"
                  />
                </el-radio>
                <el-button type="danger" plain @click="rebuild(item)">
                  rebuild
                </el-button>
              </el-row>
            </el-row>
          </el-radio-group>
          <el-pagination
            v-model:current-page="previewPagination.page"
            :total="previewPagination.total"
            :page-size="previewPagination.rows"
            style="text-align: right; margin-right: 20px"
            layout="total, prev, next"
            @current-change="handlePreviewPageChange"
          />
        </el-row>
        <el-row
          v-loading="traceLoading"
          class="project-detail"
          style="flex: 1; width: 100%"
        >
          <el-row v-for="(item, index) in publishLocalTraceList" :key="index">
            <el-row v-if="item.type === 2">
              <el-row style="margin: 5px 0">
                <i v-if="item.state === 1" class="el-icon-check icon-success" />
                <i v-else class="el-icon-close icon-fail" />
                -------------GIT-------------
              </el-row>
              <el-row style="margin: 5px 0">Time: {{ item.updateTime }}</el-row>
              <el-row v-if="item.state !== 0">
                <el-row>Branch: {{ item['branch'] }}</el-row>
                <el-row>
                  Commit:
                  <el-link
                    type="primary"
                    :underline="false"
                    :href="
                      parseGitURL(searchPreview.url) +
                      '/commit/' +
                      item['commit']
                    "
                    target="_blank"
                  >
                    {{ item['commit'] }}
                  </el-link>
                </el-row>
                <el-row>Message: {{ item['message'] }}</el-row>
                <el-row>Author: {{ item['author'] }}</el-row>
                <el-row>
                  Datetime:
                  {{ item['timestamp'] ? parseTime(item['timestamp']) : '' }}
                </el-row>
                <el-row><span v-html="enterToBR(item['diff'])" /></el-row>
              </el-row>
              <el-row v-else style="margin: 5px 0">
                <span v-html="enterToBR(item.detail)" />
              </el-row>
            </el-row>
            <el-row v-if="item.type === 3">
              <hr />
              <el-row style="margin: 5px 0">
                <i v-if="item.state === 1" class="el-icon-check icon-success" />
                <i v-else class="el-icon-close icon-fail" />
                --------After pull--------
              </el-row>
              <el-row style="margin: 5px 0">Time: {{ item.updateTime }}</el-row>
              <el-row>
                Script:
                <pre v-html="enterToBR(item.script)"></pre>
              </el-row>
              <el-row
                v-loading="traceDetail[item.id] === true"
                style="margin: 5px 0"
              >
                [goploy ~]#
                <el-button
                  v-if="item.state === 1 && !(item.id in traceDetail)"
                  type="text"
                  @click="getPublishTraceDetail(item)"
                >
                  {{ $t('deployPage.showDetail') }}
                </el-button>
                <span v-else v-html="enterToBR(item.detail)" />
              </el-row>
            </el-row>
          </el-row>
          <el-tabs v-model="activeRomoteTracePane">
            <el-tab-pane
              v-for="(item, serverName) in publishRemoteTraceList"
              :key="serverName"
              :label="serverName"
              :name="serverName"
            >
              <el-row v-for="(trace, key) in item" :key="key">
                <el-row v-if="trace.type === 4">
                  <el-row style="margin: 5px 0">
                    <i
                      v-if="trace.state === 1"
                      class="el-icon-check icon-success"
                    />
                    <i v-else class="el-icon-close icon-fail" />
                    ---------Before deploy---------
                  </el-row>
                  <el-row style="margin: 5px 0">
                    Time: {{ trace.updateTime }}
                  </el-row>
                  <el-row>
                    Script:
                    <pre v-html="enterToBR(trace.script)"></pre>
                  </el-row>
                  <el-row
                    v-loading="traceDetail[trace.id] === true"
                    style="margin: 5px 0"
                  >
                    [goploy ~]#
                    <el-button
                      v-if="trace.state === 1 && !(trace.id in traceDetail)"
                      type="text"
                      @click="getPublishTraceDetail(trace)"
                    >
                      {{ $t('deployPage.showDetail') }}
                    </el-button>
                    <span v-else v-html="enterToBR(trace.detail)" />
                  </el-row>
                </el-row>
                <el-row v-else-if="trace.type === 5">
                  <el-row style="margin: 5px 0">
                    <i
                      v-if="trace.state === 1"
                      class="el-icon-check icon-success"
                    />
                    <i v-else class="el-icon-close icon-fail" />
                    -----------Rsync------------
                  </el-row>
                  <el-row style="margin: 5px 0">
                    Time: {{ trace.updateTime }}
                  </el-row>
                  <el-row>Command: {{ trace.command }}</el-row>
                  <el-row
                    v-loading="traceDetail[trace.id] === true"
                    style="margin: 5px 0"
                  >
                    [goploy ~]#
                    <el-button
                      v-if="trace.state === 1 && !(trace.id in traceDetail)"
                      type="text"
                      @click="getPublishTraceDetail(trace)"
                    >
                      {{ $t('deployPage.showDetail') }}
                    </el-button>
                    <span v-else v-html="enterToBR(trace.detail)" />
                  </el-row>
                </el-row>
                <el-row v-else-if="trace.type === 6">
                  <el-row style="margin: 5px 0">
                    <i
                      v-if="trace.state === 1"
                      class="el-icon-check icon-success"
                    />
                    <i v-else class="el-icon-close icon-fail" />
                    --------After deploy--------
                  </el-row>
                  <el-row style="margin: 5px 0"
                    >Time: {{ trace.updateTime }}</el-row
                  >
                  <el-row>Script: {{ trace.script }}</el-row>
                  <el-row
                    v-loading="traceDetail[trace.id] === true"
                    style="margin: 5px 0"
                  >
                    [goploy ~]#
                    <el-button
                      v-if="trace.state === 1 && !(trace.id in traceDetail)"
                      type="text"
                      @click="getPublishTraceDetail(trace)"
                      >{{ $t('deployPage.showDetail') }}</el-button
                    >
                    <span v-else v-html="enterToBR(trace.detail)" />
                  </el-row>
                </el-row>
              </el-row>
            </el-tab-pane>
          </el-tabs>
        </el-row>
      </el-row>
    </el-dialog>
    <el-dialog v-model="commitDialogVisible" title="commit">
      <el-select
        v-model="branch"
        v-loading="branchLoading"
        filterable
        default-first-option
        placeholder="请选择分支"
        style="width: 100%; margin-bottom: 10px"
        @change="getCommitList"
      >
        <el-option
          v-for="item in branchOption"
          :key="item"
          :label="item"
          :value="item"
        />
      </el-select>
      <el-table
        v-loading="commitTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="commitTableData"
        style="width: 100%"
      >
        <el-table-column type="expand">
          <template #default="props">
            <span v-html="enterToBR(props.row.diff)" />
          </template>
        </el-table-column>
        <el-table-column prop="commit" label="commit" width="290">
          <template #default="scope">
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
        <el-table-column
          prop="message"
          label="message"
          width="200"
          show-overflow-tooltip
        />
        <el-table-column label="time" width="135" align="center">
          <template #default="scope">
            {{ parseTime(scope.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="180"
          align="center"
          fixed="right"
        >
          <template #default="scope">
            <el-button
              v-if="commitDialogReferer !== 'task'"
              type="danger"
              @click="publishByCommit(scope.row)"
            >
              {{ $t('deploy') }}
            </el-button>
            <el-popover
              :ref="`task${scope.row.commit}`"
              placement="bottom"
              trigger="click"
              width="270"
            >
              <el-date-picker
                v-model="scope.row.date"
                :picker-options="taskFormProps.pickerOptions"
                type="datetime"
                value-format="yyyy-MM-dd HH:mm:ss"
                style="width: 180px"
              />
              <el-button
                type="primary"
                :disabled="!scope.row['date']"
                @click="submitTask(scope.row)"
              >
                {{ $t('confirm') }}
              </el-button>
              <template #reference>
                <el-button
                  v-if="!isMember() && commitDialogReferer === 'task'"
                  type="primary"
                >
                  {{ $t('crontab') }}
                </el-button>
              </template>
            </el-popover>
            <el-button
              v-if="!isMember() && commitDialogReferer !== 'task'"
              type="warning"
              @click="handleGreyPublish(scope.row)"
            >
              {{ $t('grey') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer class="dialog-footer">
        <el-button @click="commitDialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
      </template>
    </el-dialog>
    <el-dialog v-model="tagDialogVisible" title="tag">
      <el-table
        v-loading="tagTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="tagTableData"
      >
        <el-table-column type="expand">
          <template #default="props">
            <span v-html="enterToBR(props.row.diff)" />
          </template>
        </el-table-column>
        <el-table-column prop="tag" label="tag">
          <template #default="scope">
            <el-link
              type="primary"
              style="font-size: 12px"
              :underline="false"
              :href="parseGitURL(scope.row.url) + '/tree/' + scope.row.shortTag"
              target="_blank"
            >
              {{ scope.row.shortTag }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="commit" label="commit" width="80">
          <template #default="scope">
            <el-link
              type="primary"
              style="font-size: 12px"
              :underline="false"
              :href="parseGitURL(scope.row.url) + '/commit/' + scope.row.commit"
              target="_blank"
            >
              {{ scope.row.commit.substring(0, 6) }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column
          prop="author"
          label="author"
          width="100"
          show-overflow-tooltip
        />
        <el-table-column
          prop="message"
          label="message"
          width="200"
          show-overflow-tooltip
        />
        <el-table-column label="time" width="135" align="center">
          <template #default="scope">
            {{ parseTime(scope.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="160"
          align="center"
          fixed="right"
        >
          <template #default="scope">
            <el-button type="danger" @click="publishByCommit(scope.row)">
              {{ $t('deploy') }}
            </el-button>
            <el-button
              v-if="!isMember()"
              type="warning"
              @click="handleGreyPublish(scope.row)"
            >
              {{ $t('grey') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <template #footer class="dialog-footer">
        <el-button @click="tagDialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
      </template>
    </el-dialog>
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
    <el-dialog v-model="taskListDialogVisible" :title="$t('manage')">
      <el-row class="app-bar" type="flex" justify="end">
        <el-button
          type="primary"
          icon="el-icon-plus"
          @click="handleAddProjectTask"
        />
      </el-row>
      <el-table
        v-loading="taskTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="taskTableData"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column
          prop="projectName"
          :label="$t('projectName')"
          width="150"
        />
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
                parseGitURL(selectedItem.url) +
                '/tree/' +
                scope.row.branch.split('/').pop()
              "
              target="_blank"
            >
              {{ scope.row.branch }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="commitId" label="commit" width="290">
          <template #default="scope">
            <el-link
              type="primary"
              style="font-size: 12px"
              :underline="false"
              :href="
                parseGitURL(selectedItem.url) + '/commit/' + scope.row.commitId
              "
              target="_blank"
            >
              {{ scope.row.commitId }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="date" :label="$t('date')" width="150" />
        <el-table-column prop="isRun" :label="$t('task')" width="80">
          <template #default="scope">
            {{ $t(`runOption[${scope.row.isRun}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="state" :label="$t('state')" width="70">
          <template #default="scope">
            {{ $t(`stateOption[${scope.row.state}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="creator" :label="$t('creator')" align="center" />
        <el-table-column prop="editor" :label="$t('editor')" align="center" />
        <el-table-column
          prop="insertTime"
          :label="$t('insertTime')"
          width="135"
          align="center"
        />
        <el-table-column
          prop="updateTime"
          :label="$t('updateTime')"
          width="135"
          align="center"
        />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="100"
          align="center"
          fixed="right"
        >
          <template #default="scope">
            <el-button
              type="danger"
              :disabled="scope.row.isRun === 1 || scope.row.state === 0"
              @click="removeProjectTask(scope.row)"
            >
              {{ $t('delete') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="taskPagination.page"
        hide-on-single-page
        :total="taskPagination.total"
        :page-size="taskPagination.rows"
        style="margin-top: 10px; text-align: right"
        background
        layout="sizes, total, prev, pager, next"
        @current-change="handleTaskPageChange"
      />
      <template #footer class="dialog-footer">
        <el-button @click="taskListDialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
      </template>
    </el-dialog>
    <el-dialog v-model="reviewListDialogVisible" :title="$t('review')">
      <el-table
        v-loading="reviewTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="reviewTableData"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column
          prop="projectName"
          :label="$t('projectName')"
          width="150"
        />
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
                parseGitURL(selectedItem.url) +
                '/tree/' +
                scope.row.branch.split('/').pop()
              "
              target="_blank"
            >
              {{ scope.row.branch }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="commitId" label="commit" width="290">
          <template #default="scope">
            <el-link
              type="primary"
              style="font-size: 12px"
              :underline="false"
              :href="
                parseGitURL(selectedItem.url) + '/commit/' + scope.row.commitId
              "
              target="_blank"
            >
              {{ scope.row.commitId }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="state" :label="$t('state')" width="50">
          <template #default="scope">
            {{ $t(`deployPage.reviewStateOption[${scope.row.state}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="creator" :label="$t('creator')" align="center" />
        <el-table-column prop="editor" :label="$t('editor')" align="center" />
        <el-table-column
          prop="insertTime"
          :label="$t('insertTime')"
          width="135"
          align="center"
        />
        <el-table-column
          prop="updateTime"
          :label="$t('updateTime')"
          width="135"
          align="center"
        />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="180"
          align="center"
          fixed="right"
        >
          <template #default="scope">
            <el-button
              type="success"
              :disabled="scope.row.state !== 0"
              @click="handleProjectReview(scope.row, 1)"
            >
              {{ $t('approve') }}
            </el-button>
            <el-button
              type="danger"
              :disabled="scope.row.state !== 0"
              @click="handleProjectReview(scope.row, 2)"
            >
              {{ $t('deny') }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="reviewPagination.page"
        hide-on-single-page
        :total="reviewPagination.total"
        :page-size="reviewPagination.rows"
        style="margin-top: 10px; text-align: right"
        background
        layout="sizes, total, prev, pager, next"
        @current-change="handleReviewPageChange"
      />
      <template #footer class="dialog-footer">
        <el-button @click="reviewListDialogVisible = false">
          {{ $t('cancel') }}
        </el-button>
      </template>
    </el-dialog>
  </el-row>
</template>
<script>
import tableHeight from '@/mixin/tableHeight'
import {
  getList,
  getPublishTrace,
  getPublishTraceDetail,
  getPreview,
  getCommitList,
  getBranchList,
  getTagList,
  publish,
  rebuild,
  resetState,
  review,
  greyPublish,
} from '@/api/deploy'
import {
  addTask,
  removeTask,
  getTaskList,
  getBindServerList,
  getReviewList,
} from '@/api/project'
import { getUserOption } from '@/api/namespace'
import { empty, parseTime, parseGitURL } from '@/utils'
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'Deploy',
  mixins: [tableHeight],
  data() {
    return {
      userId: '',
      userOption: [],
      publishToken: '',
      commitDialogVisible: false,
      tagDialogVisible: false,
      greyServerDialogVisible: false,
      taskListDialogVisible: false,
      reviewDialogVisible: false,
      reviewListDialogVisible: false,
      dialogVisible: false,
      traceLoading: false,
      tableloading: false,
      dateVisible: false,
      commitDialogReferer: '',
      tableData: [],
      searchProject: {
        name: '',
        environment: '',
        autoDeploy: '',
      },
      pagination: {
        total: 0,
        page: 1,
        rows: 20,
      },
      taskTableLoading: false,
      taskTableData: [],
      taskPagination: {
        total: 0,
        page: 1,
        rows: 20,
      },
      taskFormProps: {
        dateVisible: true,
        pickerOptions: {
          disabledDate(time) {
            return time.getTime() < Date.now() - 3600 * 1000 * 24
          },
        },
      },
      reviewTableLoading: false,
      reviewTableData: [],
      reviewPagination: {
        total: 0,
        page: 1,
        rows: 20,
      },
      searchPreview: {
        loading: false,
        projectId: '',
        userId: '',
        url: '',
        state: '',
        filename: '',
        branch: '',
        commit: '',
        commitDate: [],
        deployDate: [],
      },
      pickerOptions: {
        shortcuts: [
          {
            text: '最近一周',
            onClick(picker) {
              const end = new Date()
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
              picker.$emit('pick', [start, end])
            },
          },
          {
            text: '最近一个月',
            onClick(picker) {
              const end = new Date()
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
              picker.$emit('pick', [start, end])
            },
          },
          {
            text: '最近三个月',
            onClick(picker) {
              const end = new Date()
              const start = new Date()
              start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
              picker.$emit('pick', [start, end])
            },
          },
        ],
      },
      gitTraceList: [],
      previewPagination: {
        page: 1,
        rows: 11,
        total: 0,
      },
      commitTableLoading: false,
      branchOption: [],
      branchLoading: false,
      branch: '',
      commitTableData: [],
      tagTableLoading: false,
      tagTableData: [],
      greyServerFormProps: {
        disabled: false,
        serverOption: [],
      },
      greyServerFormData: {
        projectId: 0,
        commit: 0,
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
      publishTraceList: [],
      publishLocalTraceList: [],
      publishRemoteTraceList: {},
      traceDetail: {},
      activeRomoteTracePane: '',
    }
  },
  computed: {
    tablePageData: function () {
      let tableData = this.tableData
      if (this.searchProject.name !== '') {
        tableData = this.tableData.filter(
          (item) => item.name.indexOf(this.searchProject.name) !== -1
        )
      }
      if (this.searchProject.environment !== '') {
        tableData = this.tableData.filter(
          (item) => item.environment === this.searchProject.environment
        )
      }
      if (this.searchProject.autoDeploy !== '') {
        tableData = this.tableData.filter(
          (item) => item.autoDeploy === this.searchProject.autoDeploy
        )
      }
      return tableData.slice(
        (this.pagination.page - 1) * this.pagination.rows,
        this.pagination.page * this.pagination.rows
      )
    },
    previewFilterlength: function () {
      let number = 0
      for (const key in this.searchPreview) {
        if (['projectId', 'loading', 'url'].indexOf(key) !== -1) {
          continue
        }
        if (!empty(this.searchPreview[key])) {
          number++
        }
      }
      return number
    },
  },
  watch: {
    '$store.getters.ws_message': function (response) {
      if (response.type !== 1) {
        return
      }
      const data = response.message
      data.message = this.enterToBR(data.message)
      if (data.state === 0) {
        this.$notify.error({
          title: data.projectName,
          dangerouslyUseHTMLString: true,
          message: data.message,
          duration: 0,
        })
      }
      const projectIndex = this.tableData.findIndex(
        (element) => element.id === data.projectId
      )
      if (projectIndex !== -1) {
        const percent = 12.5 * data.state
        this.tableData[projectIndex].progressPercentage = percent
        this.tableData[projectIndex].progressStatus = 'warning'
        this.tableData[projectIndex].tagType = 'warning'
        this.tableData[projectIndex].tagText = data.message
        this.tableData[projectIndex].deployState = 1
        if (percent === 0) {
          this.tableData[projectIndex].progressStatus = 'exception'
          this.tableData[projectIndex].tagType = 'danger'
          this.tableData[projectIndex].tagText = 'Fail'
          this.tableData[projectIndex].deployState = 3
        } else if (percent > 98) {
          this.tableData[projectIndex].progressStatus = 'success'
          this.tableData[projectIndex].tagType = 'success'
          this.tableData[projectIndex].deployState = 2
        }

        if (data['ext']) {
          Object.assign(this.tableData[projectIndex], data['ext'])
        }
        this.tableData[projectIndex].publisherName = data.username
        this.tableData[projectIndex].updateTime = parseTime(new Date())
      }
    },
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
      this.tableloading = true
      getList()
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
          this.pagination.total = this.tableData.length
        })
        .finally(() => {
          this.tableloading = false
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

    getPublishTrace() {
      this.traceLoading = true
      getPublishTrace(this.publishToken)
        .then((response) => {
          const publishTraceList = response.data.publishTraceList || []
          this.publishTraceList = publishTraceList.map((element) => {
            if (element.ext !== '')
              Object.assign(element, JSON.parse(element.ext))
            return element
          })

          this.publishLocalTraceList = this.publishTraceList.filter(
            (element) => element.type < 4
          )
          this.publishRemoteTraceList = {}
          for (const trace of this.publishTraceList) {
            if (trace.type < 4) continue
            if (!this.publishRemoteTraceList[trace.serverName]) {
              this.publishRemoteTraceList[trace.serverName] = []
            }
            this.publishRemoteTraceList[trace.serverName].push(trace)
          }
          this.activeRomoteTracePane = Object.keys(
            this.publishRemoteTraceList
          )[0]
        })
        .finally(() => {
          this.traceLoading = false
        })
    },

    getPublishTraceDetail(data) {
      this.$set(this.traceDetail, data.id, true)
      getPublishTraceDetail(data.id)
        .then((response) => {
          data.detail =
            response.data.detail === ''
              ? this.$t('deployPage.noDetail')
              : response.data.detail
        })
        .finally(() => {
          this.$set(this.traceDetail, data.id, false)
        })
    },

    getPreviewList() {
      this.searchPreview.loading = true
      this.traceDetail = {}
      getPreview(this.previewPagination, {
        projectId: this.searchPreview.projectId,
        commitDate: this.searchPreview.commitDate
          ? this.searchPreview.commitDate.join(',')
          : '',
        deployDate: this.searchPreview.deployDate
          ? this.searchPreview.deployDate.join(',')
          : '',
        branch: this.searchPreview.branch,
        commit: this.searchPreview.commit,
        filename: this.searchPreview.filename,
        userId: this.searchPreview.userId || 0,
        state: this.searchPreview.state === '' ? -1 : this.searchPreview.state,
      })
        .then((response) => {
          const gitTraceList = response.data.gitTraceList || []
          this.gitTraceList = gitTraceList.map((element) => {
            if (element.ext !== '')
              Object.assign(element, JSON.parse(element.ext))
            element.commit = element['commit']
              ? element['commit'].substring(0, 6)
              : ''
            return element
          })
          if (this.gitTraceList.length > 0) {
            this.publishToken = this.gitTraceList[0].token
            this.getPublishTrace()
          } else {
            this.publishLocalTraceList = []
            this.publishRemoteTraceList = {}
          }
          this.previewPagination.total = response.data.pagination.total
        })
        .finally(() => {
          this.searchPreview.loading = false
        })
    },

    refreshSearchPreviewCondition() {
      this.searchPreview.userId = ''
      this.searchPreview.state = ''
      this.searchPreview.filename = ''
      this.searchPreview.branch = ''
      this.searchPreview.commit = ''
      this.searchPreview.commitDate = []
      this.searchPreview.deployDate = []
    },

    searchPreviewList() {
      this.handlePreviewPageChange(1)
    },

    handleDetail(data) {
      this.dialogVisible = true
      this.searchPreview.projectId = data.id
      this.searchPreview.url = data.url
      this.searchPreview.userId = ''
      this.getPreviewList()
    },

    handlePreviewPageChange(page) {
      this.previewPagination.page = page
      this.getPreviewList()
    },

    handleTraceChange(lastPublishToken) {
      this.publishToken = lastPublishToken
      this.getPublishTrace()
    },

    handleGreyPublish(data) {
      getBindServerList(data.projectId).then((response) => {
        this.greyServerFormProps.serverOption = response.data.list
      })
      // add projectID to server form
      this.greyServerFormData.projectId = data.projectId
      this.greyServerFormData.commit = data.commit
      this.greyServerDialogVisible = true
    },

    handleCommitCommand(data) {
      this.commitDialogReferer = 'commit'
      this.selectedItem = data
      this.commitDialogVisible = true
      this.getBranchList()
    },

    getBranchList() {
      const id = this.selectedItem.id
      this.branchLoading = true
      this.branch = ''
      this.branchOption = []
      this.commitTableData = []
      getBranchList(id)
        .then((response) => {
          this.branchOption = response.data.branchList.filter((element) => {
            return element.indexOf('HEAD') === -1
          })
        })
        .finally(() => {
          this.branchLoading = false
        })
    },

    getCommitList() {
      this.commitTableLoading = true
      getCommitList(this.selectedItem.id, this.branch)
        .then((response) => {
          this.commitTableData = response.data.commitList
            ? response.data.commitList.map((element) => {
                return Object.assign(element, {
                  projectId: this.selectedItem.id,
                  url: this.selectedItem.url,
                  branch: this.branch,
                })
              })
            : []
        })
        .finally(() => {
          this.commitTableLoading = false
        })
    },

    getTagList(data) {
      const id = data.id
      this.tagDialogVisible = true
      this.tagTableLoading = true
      getTagList(id)
        .then((response) => {
          this.tagTableData = response.data.tagList
            ? response.data.tagList.map((element) => {
                let shortTag = element.tag.replace(/[()]/g, '')
                for (const tag of shortTag.split(',')) {
                  if (tag.indexOf('tag: ') !== -1) {
                    shortTag = tag.replace('tag: ', '').trim()
                    break
                  }
                }

                return Object.assign(element, {
                  projectId: id,
                  url: data.url,
                  shortTag: shortTag,
                })
              })
            : []
        })
        .finally(() => {
          this.tagTableLoading = false
        })
    },

    getTaskList() {
      this.taskTableLoading = true
      getTaskList(this.taskPagination, this.selectedItem.id)
        .then((response) => {
          const projectTaskList = response.data.list
          this.taskTableData = projectTaskList.map((element) => {
            return Object.assign(element, {
              projectId: this.selectedItem.id,
              projectName: this.selectedItem.name,
            })
          })
          this.taskPagination.total = response.data.pagination.total
        })
        .finally(() => {
          this.taskTableLoading = false
        })
    },

    handleTaskCommand(data) {
      this.selectedItem = data
      this.taskPagination.page = 1
      this.taskListDialogVisible = true
      this.getTaskList()
    },

    handleTaskPageChange(page) {
      this.taskPagination.page = page
      this.getTaskList()
    },

    handleAddProjectTask() {
      this.commitDialogReferer = 'task'
      this.commitDialogVisible = true
      this.getBranchList()
    },

    submitTask(data) {
      addTask(data)
        .then(() => {
          this.$message.success('Success')
        })
        .finally(() => {
          this.$refs[`task${data.commit}`].doClose()
          this.taskPagination.page = 1
          this.commitDialogVisible = false
          this.getTaskList()
        })
    },

    removeProjectTask(data) {
      this.$confirm(
        this.$i18n.t('deployPage.removeProjectTaskTips', {
          projectName: data.projectName,
        }),
        this.$i18n.t('tips'),
        {
          confirmButtonText: this.$i18n.t('confirm'),
          cancelButtonText: this.$i18n.t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          removeTask(data.id).then((_) => {
            const projectTaskIndex = this.taskTableData.findIndex(
              (element) => element.id === data.id
            )
            this.taskTableData[projectTaskIndex]['state'] = 0
            this.taskTableData[projectTaskIndex]['editor'] =
              this.$store.getters.name
            this.taskTableData[projectTaskIndex]['editorId'] =
              this.$store.getters.uid
            this.taskTableData[projectTaskIndex]['updateTime'] = parseTime(
              new Date()
            )
          })
        })
        .catch(() => {
          this.$message.info('Cancel')
        })
    },

    getReviewList() {
      this.reviewTableLoading = true
      getReviewList(this.reviewPagination, this.selectedItem.id)
        .then((response) => {
          const projectTaskList = response.data.list
          this.reviewTableData = projectTaskList.map((element) => {
            return Object.assign(element, {
              projectId: this.selectedItem.id,
              projectName: this.selectedItem.name,
            })
          })
          this.reviewPagination.total = response.data.pagination.total
        })
        .finally(() => {
          this.reviewTableLoading = false
        })
    },

    handleReviewCommand(data) {
      this.selectedItem = data
      this.reviewPagination.page = 1
      this.reviewListDialogVisible = true
      this.getReviewList()
    },

    handleReviewPageChange(page) {
      this.reviewPagination.page = page
      this.getReviewList()
    },

    handleProjectReview(data, state) {
      if (state === 1) {
        this.$confirm(
          this.$i18n.t('deployPage.reviewTips'),
          this.$i18n.t('tips'),
          {
            confirmButtonText: this.$i18n.t('confirm'),
            cancelButtonText: this.$i18n.t('cancel'),
            type: 'warning',
          }
        )
          .then(() => {
            review(data.id, state).then((response) => {
              const projectIndex = this.tableData.findIndex(
                (element) => element.id === data.projectId
              )
              this.tableData[projectIndex].state = state
              this.reviewListDialogVisible = false
            })
          })
          .catch(() => {
            this.$message.info('Cancel')
          })
      } else {
        review(data.id, state).then((response) => {
          const projectIndex = this.tableData.findIndex(
            (element) => element.id === data.projectId
          )
          this.tableData[projectIndex].state = state
          this.reviewListDialogVisible = false
        })
      }
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
          h(
            'b',
            { style: color },
            data.name + ' - ' + this.$i18n.t(`envOption[${data.environment}]`)
          ),
        ]),
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning',
      })
        .then(() => {
          this.gitLog = []
          this.remoteLog = {}
          publish(id, '').then((_) => {
            const projectIndex = this.tableData.findIndex(
              (element) => element.id === id
            )
            this.tableData[projectIndex].deployState = 1
          })
        })
        .catch(() => {
          this.$message.info('Cancel')
        })
    },

    publishByCommit(data) {
      this.$confirm(
        this.$i18n.t('deployPage.publishCommitTips', { commit: data.commit }),
        this.$i18n.t('tips'),
        {
          confirmButtonText: this.$i18n.t('confirm'),
          cancelButtonText: this.$i18n.t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          publish(data.projectId, data.branch, data.commit).then((response) => {
            const projectIndex = this.tableData.findIndex(
              (element) => element.id === data.projectId
            )
            this.tableData[projectIndex].deployState = 1
            this.commitDialogVisible = false
            this.tagDialogVisible = false
          })
        })
        .catch(() => {
          this.$message.info('Cancel')
        })
    },

    rebuild(data) {
      this.$confirm(
        this.$i18n.t('deployPage.publishCommitTips', { commit: data.commit }),
        this.$i18n.t('tips'),
        {
          confirmButtonText: this.$i18n.t('confirm'),
          cancelButtonText: this.$i18n.t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          this.searchPreview.loading = true
          rebuild(data.projectId, data.token).then((response) => {
            if (response.data === 'symlink') {
              this.$message.success('Success')
            } else {
              const projectIndex = this.tableData.findIndex(
                (element) => element.id === data.projectId
              )
              this.tableData[projectIndex].deployState = 1
            }
            this.dialogVisible = false
          })
        })
        .catch(() => {
          this.$message.info('Cancel')
        })
    },

    resetState(data) {
      this.$confirm(
        this.$i18n.t('deployPage.resetStateTips'),
        this.$i18n.t('tips'),
        {
          confirmButtonText: this.$i18n.t('confirm'),
          cancelButtonText: this.$i18n.t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          resetState(data.id).then((response) => {
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
          this.$message.info('Cancel')
        })
    },

    greyPublish() {
      this.$refs.greyServerForm.validate((valid) => {
        if (valid) {
          const data = this.greyServerFormData
          this.$confirm(
            this.$i18n.t('deployPage.publishCommitTips', {
              commit: data.commit,
            }),
            this.$i18n.t('tips'),
            {
              confirmButtonText: this.$i18n.t('confirm'),
              cancelButtonText: this.$i18n.t('cancel'),
              type: 'warning',
            }
          )
            .then(() => {
              greyPublish(data.projectId, data.commit, data.serverIds).then(
                (response) => {
                  const projectIndex = this.tableData.findIndex(
                    (element) => element.id === data.projectId
                  )
                  this.tableData[projectIndex].deployState = 1
                  this.commitDialogVisible = false
                  this.tagDialogVisible = false
                  this.greyServerDialogVisible = false
                }
              )
            })
            .catch(() => {
              this.$message.info('Cancel')
            })
        } else {
          return false
        }
      })
    },

    enterToBR(detail) {
      return detail ? detail.replace(/\n|(\r\n)/g, '<br>') : ''
    },
  },
})
</script>
<style rel="stylesheet/scss" lang="scss" scoped>
@import '@/styles/mixin.scss';

.publish {
  &-preview {
    width: 330px;
    margin-left: 10px;
  }
  &-commit {
    margin-right: 5px;
    padding-right: 8px;
    width: 240px;
    line-height: 12px;
  }
  &-commitID {
    display: inline-block;
    vertical-align: top;
  }
  &-name {
    width: 60px;
    display: inline-block;
    text-align: center;
    overflow: hidden;
    vertical-align: top;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
}

.publish-filter-label {
  font-size: 12px;
  width: 70px;
}

.project-detail {
  padding-left: 5px;
  height: 470px;
  overflow-y: auto;
  @include scrollBar();
}

.operation-btn {
  >>> .el-button {
    line-height: 1.15;
  }
}

.icon-success {
  color: #67c23a;
  font-size: 14px;
  font-weight: 900;
}

.icon-fail {
  color: #f56c6c;
  font-size: 14px;
  font-weight: 900;
}

@media screen and (max-width: 1440px) {
  .publish-record {
    >>> .el-dialog {
      width: 75%;
    }
  }
}
</style>
