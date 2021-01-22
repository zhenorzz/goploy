<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex">
      <el-input v-model="projectName" style="width:300px" placeholder="Filter the project name" @change="getList" />
    </el-row>
    <el-table
      :key="tableHeight"
      v-loading="tableloading"
      border
      stripe
      highlight-current-row
      :max-height="tableHeight"
      :data="tablePageData"
      style="width: 100%;margin-top: 5px;"
    >
      <el-table-column prop="id" label="ID" width="80" align="center" />
      <el-table-column prop="name" :label="$t('name')" min-width="150" align="center">
        <template slot-scope="scope">
          <b v-if="scope.row.environment === 1" style="color: #F56C6C">{{ scope.row.name }} - {{ $t(`envOption[${scope.row.environment}]`) }}</b>
          <b v-else-if="scope.row.environment === 3" style="color: #E6A23C">{{ scope.row.name }} - {{ $t(`envOption[${scope.row.environment}]`) }}</b>
          <b v-else style="color: #909399">{{ scope.row.name }} - {{ $t(`envOption[${scope.row.environment}]`) }}</b>
        </template>
      </el-table-column>
      <el-table-column prop="branch" :label="$t('branch')" width="150" align="center">
        <template slot-scope="scope">
          <el-link
            style="font-size: 12px"
            :underline="false"
            :href="parseGitURL(scope.row['url']) + '/tree/' + scope.row['branch'].split('/').pop()"
            target="_blank"
          >
            {{ scope.row.branch }}
          </el-link>
        </template>
      </el-table-column>
      <el-table-column prop="commit" label="CommitID" width="100" align="center">
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
            <el-button v-if="scope.row.deployState === 0" type="primary" @click="publish(scope.row)">{{ $t('initial') }}</el-button>
            <el-button v-else-if="hasManagerPermission() && scope.row.deployState === 1" type="primary" @click="resetState(scope.row)">{{ $t('deployPage.resetState') }}</el-button>
            <el-dropdown
              v-else-if="hasGroupManagerPermission() || scope.row.review === 0"
              split-button
              trigger="click"
              type="primary"
              @click="publish(scope.row)"
              @command="(commandFunc) => commandFunc(scope.row)"
            >
              {{ isMember() && scope.row.review === 1 ? $t('submit') : $t('deploy') }}
              <el-dropdown-menu slot="dropdown">
                <el-dropdown-item :command="getBranchList">Commit list</el-dropdown-item>
                <el-dropdown-item :command="getTagList">Tag list</el-dropdown-item>
              </el-dropdown-menu>
            </el-dropdown>
            <el-button v-else type="primary" @click="getBranchList(scope.row)">{{ $t('deploy') }}</el-button>
            <el-dropdown
              v-if="hasGroupManagerPermission() || scope.row.review === 1"
              trigger="click"
              @command="(commandFunc) => commandFunc(scope.row)"
            >
              <el-button type="warning">
                {{ $t('func') }}<i class="el-icon-arrow-down el-icon--right" />
              </el-button>
              <el-dropdown-menu slot="dropdown" style="min-width:84px;text-align:center;">
                <el-dropdown-item v-if="hasGroupManagerPermission()" :command="handleTaskCommand">{{ $t('deployPage.taskDeploy') }}</el-dropdown-item>
                <el-dropdown-item v-if="scope.row.review === 1" :command="handleReviewCommand">{{ $t('deployPage.reviewDeploy') }}</el-dropdown-item>
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
            <el-popover
              placement="bottom-start"
              width="318"
              trigger="click"
            >
              <el-row type="flex" align="middle">
                <label class="publish-filter-label">{{ $t('user') }}</label>
                <el-select v-model="searchPreview.userId" style="flex:1" clearable>
                  <el-option
                    v-for="(item, index) in userOption"
                    :key="index"
                    :label="item.userName"
                    :value="item.userId"
                  />
                </el-select>
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px;">
                <label class="publish-filter-label">Commit</label>
                <el-input v-model.trim="searchPreview.commit" autocomplete="off" style="flex:1" placeholder="Commit" />
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px;">
                <label class="publish-filter-label">{{ $t('branch') }}</label>
                <el-input v-model.trim="searchPreview.branch" autocomplete="off" style="flex:1" :placeholder="$t('branch')" />
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px;">
                <label class="publish-filter-label">{{ $t('filename') }}</label>
                <el-input v-model.trim="searchPreview.filename" autocomplete="off" style="flex:1" :placeholder="$t('filename')" />
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px;">
                <label class="publish-filter-label">{{ $t('state') }}</label>
                <el-select v-model="searchPreview.state" style="flex:1" clearable>
                  <el-option :label="$t('success')" :value="1" />
                  <el-option :label="$t('fail')" :value="0" />
                </el-select>
              </el-row>
              <el-row type="flex" align="middle" style="margin-top: 5px;">
                <label class="publish-filter-label">{{ $t('commitDate') }}</label>
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
              <el-row type="flex" align="middle" style="margin-top: 5px;">
                <label class="publish-filter-label">{{ $t('deployDate') }}</label>
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
              <el-button slot="reference" icon="el-icon-notebook-2" style="width: 220px;">条件筛选({{ previewFilterlength }})</el-button>
            </el-popover>
            <el-button type="warning" icon="el-icon-refresh" @click="refreshSearchPreviewCondition" />
            <el-button type="primary" icon="el-icon-search" style="margin-left: 2px;" @click="searchPreviewList" />
          </el-row>
          <el-radio-group v-model="publishToken" @change="handleTraceChange">
            <el-row v-for="(item, index) in gitTraceList" :key="index">
              <el-row style="margin:5px 0">
                <el-radio class="publish-commit" :label="item.token" border>
                  <span class="publish-name">{{ item.publisherName }}</span> <span class="publish-commitID">commitID: {{ item.commit }}</span>
                  <i v-if="item.publishState === 1" class="el-icon-check icon-success" style="float:right;" />
                  <i v-else class="el-icon-close icon-fail" style="float:right;" />
                </el-radio>
                <el-button type="danger" plain @click="publishByCommit(item)">rebuild</el-button>
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
        <el-row v-loading="traceLoading" class="project-detail" style="flex:1;width:100%">
          <el-row v-for="(item, index) in publishLocalTraceList" :key="index">
            <el-row v-if="item.type === 2">
              <el-row style="margin:5px 0">
                <i v-if="item.state === 1" class="el-icon-check icon-success" />
                <i v-else class="el-icon-close icon-fail" />
                -------------GIT-------------
              </el-row>
              <el-row style="margin:5px 0">Time: {{ item.insertTime }}</el-row>
              <!-- 用数组的形式 兼容以前版本 -->
              <el-row v-if="item.state !== 0">
                <el-row>Branch: {{ item['branch'] }}</el-row>
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
                <el-row><span v-html="enterToBR(item['diff'])" /></el-row>
              </el-row>
              <el-row v-else style="margin:5px 0">
                <span v-html="enterToBR(item.detail)" />
              </el-row>
            </el-row>
            <el-row v-if="item.type === 3">
              <hr>
              <el-row style="margin:5px 0">
                <i v-if="item.state === 1" class="el-icon-check icon-success" />
                <i v-else class="el-icon-close icon-fail" />
                --------After pull--------
              </el-row>
              <el-row style="margin:5px 0">Time: {{ item.insertTime }}</el-row>
              <el-row>Script: <pre v-html="enterToBR(item.script)" /></el-row>
              <el-row v-loading="traceDetail[item.id] === true" style="margin:5px 0">
                [goploy ~]#
                <el-button v-if="item.state === 1 && !(item.id in traceDetail)" type="text" @click="getPublishTraceDetail(item)">{{ $t('deployPage.showDetail') }}</el-button>
                <span v-else v-html="enterToBR(item.detail)" />
              </el-row>
            </el-row>
          </el-row>
          <el-tabs v-model="activeRomoteTracePane">
            <el-tab-pane v-for="(item, serverName) in publishRemoteTraceList" :key="serverName" :label="serverName" :name="serverName">
              <el-row v-for="(trace, key) in item" :key="key">
                <el-row v-if="trace.type === 4">
                  <el-row style="margin:5px 0">
                    <i v-if="trace.state === 1" class="el-icon-check icon-success" />
                    <i v-else class="el-icon-close icon-fail" />
                    ---------Before deploy---------
                  </el-row>
                  <el-row style="margin:5px 0">Time: {{ trace.insertTime }}</el-row>
                  <el-row>Script: <pre v-html="enterToBR(trace.script)" /></el-row>
                  <el-row v-loading="traceDetail[trace.id] === true" style="margin:5px 0">
                    [goploy ~]#
                    <el-button v-if="trace.state === 1 && !(trace.id in traceDetail)" type="text" @click="getPublishTraceDetail(trace)">{{ $t('deployPage.showDetail') }}</el-button>
                    <span v-else v-html="enterToBR(trace.detail)" />
                  </el-row>
                </el-row>
                <el-row v-else-if="trace.type === 5">
                  <el-row style="margin:5px 0">
                    <i v-if="trace.state === 1" class="el-icon-check icon-success" />
                    <i v-else class="el-icon-close icon-fail" />
                    -----------Rsync------------
                  </el-row>
                  <el-row style="margin:5px 0">Time: {{ trace.insertTime }}</el-row>
                  <el-row>Command: {{ trace.command }}</el-row>
                  <el-row v-loading="traceDetail[trace.id] === true" style="margin:5px 0">
                    [goploy ~]#
                    <el-button v-if="trace.state === 1 && !(trace.id in traceDetail)" type="text" @click="getPublishTraceDetail(trace)">{{ $t('deployPage.showDetail') }}</el-button>
                    <span v-else v-html="enterToBR(trace.detail)" />
                  </el-row>
                </el-row>
                <el-row v-else-if="trace.type === 6">
                  <el-row style="margin:5px 0">
                    <i v-if="trace.state === 1" class="el-icon-check icon-success" />
                    <i v-else class="el-icon-close icon-fail" />
                    --------After deploy--------
                  </el-row>
                  <el-row style="margin:5px 0">Time: {{ trace.insertTime }}</el-row>
                  <el-row>Script: {{ trace.script }}</el-row>
                  <el-row v-loading="traceDetail[trace.id] === true" style="margin:5px 0">
                    [goploy ~]#
                    <el-button v-if="trace.state === 1 && !(trace.id in traceDetail)" type="text" @click="getPublishTraceDetail(trace)">{{ $t('deployPage.showDetail') }}</el-button>
                    <span v-else v-html="enterToBR(trace.detail)" />
                  </el-row>
                </el-row>
              </el-row>
            </el-tab-pane>
          </el-tabs>
        </el-row>
      </el-row>
    </el-dialog>
    <el-dialog title="commit" :visible.sync="commitDialogVisible">
      <el-select
        v-model="branch"
        v-loading="branchLoading"
        filterable
        default-first-option
        placeholder="请选择分支"
        style="width:100%;margin-bottom: 10px;"
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
      >
        <el-table-column type="expand">
          <template slot-scope="props">
            <span v-html="enterToBR(props.row.diff)" />
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
        <el-table-column prop="operation" :label="$t('op')" width="260" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button type="danger" @click="publishByCommit(scope.row)">{{ $t('deploy') }}</el-button>
            <el-button v-if="!isMember()" type="primary" @click="handleAddProjectTask(scope.row)">{{ $t('crontab') }}</el-button>
            <el-button v-if="!isMember()" type="warning" @click="handleGreyPublish(scope.row)">{{ $t('grey') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="commitDialogVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog title="tag" :visible.sync="tagDialogVisible">
      <el-table
        v-loading="tagTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="tagTableData"
      >
        <el-table-column type="expand">
          <template slot-scope="props">
            <span v-html="enterToBR(props.row.diff)" />
          </template>
        </el-table-column>
        <el-table-column prop="tag" label="tag">
          <template slot-scope="scope">
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
          <template slot-scope="scope">
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
        <el-table-column prop="author" label="author" width="100" show-overflow-tooltip />
        <el-table-column prop="message" label="message" width="200" show-overflow-tooltip />
        <el-table-column label="time" width="135" align="center">
          <template slot-scope="scope">
            {{ parseTime(scope.row.timestamp) }}
          </template>
        </el-table-column>
        <el-table-column prop="operation" :label="$t('op')" width="160" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button type="danger" @click="publishByCommit(scope.row)">{{ $t('deploy') }}</el-button>
            <el-button v-if="!isMember()" type="warning" @click="handleGreyPublish(scope.row)">{{ $t('grey') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div slot="footer" class="dialog-footer">
        <el-button @click="tagDialogVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('deploy')" :visible.sync="greyServerDialogVisible">
      <el-form ref="greyServerForm" :rules="greyServerFormRules" :model="greyServerFormData">
        <el-form-item :label="$t('server')" label-width="80px" prop="serverIds">
          <el-checkbox-group v-model="greyServerFormData.serverIds">
            <el-checkbox
              v-for="(item, index) in greyServerFormProps.serverOption"
              :key="index"
              :label="item.serverId"
            >
              {{ item.serverName+ '(' + item.serverDescription + ')' }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="greyServerDialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button :disabled="greyServerFormProps.disabled" type="primary" @click="greyPublish">{{ $t('confirm') }}</el-button>
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
        <el-table-column prop="branch" :label="$t('branch')" width="150" align="center">
          <template slot-scope="scope">
            <el-link
              style="font-size: 12px"
              :underline="false"
              :href="parseGitURL(selectedItem.url) + '/tree/' + scope.row.branch.split('/').pop()"
              target="_blank"
            >
              {{ scope.row.branch }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="commitId" label="commit" width="290">
          <template slot-scope="scope">
            <el-link
              type="primary"
              style="font-size: 12px"
              :underline="false"
              :href="parseGitURL(selectedItem.url) + '/commit/' + scope.row.commitId"
              target="_blank"
            >
              {{ scope.row.commitId }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="date" :label="$t('date')" width="150" />
        <el-table-column prop="isRun" :label="$t('task')" width="80">
          <template slot-scope="scope">
            {{ $t(`runOption[${scope.row.isRun}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="state" :label="$t('state')" width="70">
          <template slot-scope="scope">
            {{ $t(`stateOption[${scope.row.state}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="creator" :label="$t('creator')" align="center" />
        <el-table-column prop="editor" :label="$t('editor')" align="center" />
        <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
        <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
        <el-table-column prop="operation" :label="$t('op')" width="100" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button type="danger" :disabled="scope.row.isRun === 1 || scope.row.state === 0" @click="removeProjectTask(scope.row)">{{ $t('delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        hide-on-single-page
        :total="taskPagination.total"
        :page-size="taskPagination.rows"
        :current-page.sync="taskPagination.page"
        style="margin-top:10px; text-align:right;"
        background
        layout="sizes, total, prev, pager, next"
        @current-change="handleTaskPageChange"
      />
      <div slot="footer" class="dialog-footer">
        <el-button @click="taskListDialogVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('review')" :visible.sync="reviewListDialogVisible">
      <el-table
        v-loading="reviewTableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="reviewTableData"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="projectName" :label="$t('projectName')" width="150" />
        <el-table-column prop="branch" :label="$t('branch')" width="150" align="center">
          <template slot-scope="scope">
            <el-link
              style="font-size: 12px"
              :underline="false"
              :href="parseGitURL(selectedItem.url) + '/tree/' + scope.row.branch.split('/').pop()"
              target="_blank"
            >
              {{ scope.row.branch }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="commitId" label="commit" width="290">
          <template slot-scope="scope">
            <el-link
              type="primary"
              style="font-size: 12px"
              :underline="false"
              :href="parseGitURL(selectedItem.url) + '/commit/' + scope.row.commitId"
              target="_blank"
            >
              {{ scope.row.commitId }}
            </el-link>
          </template>
        </el-table-column>
        <el-table-column prop="state" :label="$t('state')" width="50">
          <template slot-scope="scope">
            {{ $t(`deployPage.reviewStateOption[${scope.row.state}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="creator" :label="$t('creator')" align="center" />
        <el-table-column prop="editor" :label="$t('editor')" align="center" />
        <el-table-column prop="insertTime" :label="$t('insertTime')" width="135" align="center" />
        <el-table-column prop="updateTime" :label="$t('updateTime')" width="135" align="center" />
        <el-table-column prop="operation" :label="$t('op')" width="180" align="center" fixed="right">
          <template slot-scope="scope">
            <el-button type="success" :disabled="scope.row.state !== 0" @click="handleProjectReview(scope.row, 1)">{{ $t('approve') }}</el-button>
            <el-button type="danger" :disabled="scope.row.state !== 0" @click="handleProjectReview(scope.row, 2)">{{ $t('deny') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        hide-on-single-page
        :total="reviewPagination.total"
        :page-size="reviewPagination.rows"
        :current-page.sync="reviewPagination.page"
        style="margin-top:10px; text-align:right;"
        background
        layout="sizes, total, prev, pager, next"
        @current-change="handleReviewPageChange"
      />
      <div slot="footer" class="dialog-footer">
        <el-button @click="reviewListDialogVisible = false">{{ $t('cancel') }}</el-button>
      </div>
    </el-dialog>
    <el-dialog :title="$t('setting')" :visible.sync="taskDialogVisible" width="600px">
      <el-form ref="taskForm" :rules="taskFormRules" :model="taskFormData" label-width="120px">
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
import { getList, getPublishTrace, getPublishTraceDetail, getPreview, getCommitList, getBranchList, getTagList, publish, resetState, review, greyPublish } from '@/api/deploy'
import { addTask, removeTask, getTaskList, getBindServerList, getReviewList } from '@/api/project'
import { getUserOption } from '@/api/namespace'
import { empty, parseTime, parseGitURL } from '@/utils'

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
      tagDialogVisible: false,
      greyServerDialogVisible: false,
      taskDialogVisible: false,
      taskListDialogVisible: false,
      reviewDialogVisible: false,
      reviewListDialogVisible: false,
      dialogVisible: false,
      traceLoading: false,
      tableloading: false,
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
        disabled: false,
        pickerOptions: {
          disabledDate(time) {
            return time.getTime() < Date.now() - 3600 * 1000 * 24
          }
        }
      },
      taskFormData: {
        id: 0,
        projectId: '',
        branch: '',
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
      reviewTableLoading: false,
      reviewTableData: [],
      reviewPagination: {
        total: 0,
        page: 1,
        rows: 20
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
        deployDate: []
      },
      pickerOptions: {
        shortcuts: [{
          text: '最近一周',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 7)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近一个月',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 30)
            picker.$emit('pick', [start, end])
          }
        }, {
          text: '最近三个月',
          onClick(picker) {
            const end = new Date()
            const start = new Date()
            start.setTime(start.getTime() - 3600 * 1000 * 24 * 90)
            picker.$emit('pick', [start, end])
          }
        }]
      },
      gitTraceList: [],
      previewPagination: {
        page: 1,
        rows: 11,
        total: 0
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
        serverOption: []
      },
      greyServerFormData: {
        projectId: 0,
        commit: 0,
        serverIds: []
      },
      greyServerFormRules: {
        serverIds: [
          { type: 'array', required: true, message: 'Server required', trigger: 'change' }
        ]
      },
      publishTraceList: [],
      publishLocalTraceList: [],
      publishRemoteTraceList: {},
      traceDetail: {},
      activeRomoteTracePane: ''
    }
  },
  computed: {
    tablePageData: function() {
      return this.tableData.slice((this.pagination.page - 1) * this.pagination.rows, this.pagination.page * this.pagination.rows)
    },
    previewFilterlength: function() {
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
    }
  },
  watch: {
    '$store.getters.ws_message': function(response) {
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
      this.tableloading = true
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
      }).finally(() => {
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
      getPublishTrace(this.publishToken).then((response) => {
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
      }).finally(() => {
        this.traceLoading = false
      })
    },

    getPublishTraceDetail(data) {
      this.$set(this.traceDetail, data.id, true)
      getPublishTraceDetail(data.id).then((response) => {
        data.detail = response.data.detail === '' ? this.$t('deployPage.noDetail') : response.data.detail
      }).finally(() => { this.$set(this.traceDetail, data.id, false) })
    },

    getPreviewList() {
      this.searchPreview.loading = true
      this.traceDetail = {}
      getPreview(this.previewPagination, {
        projectId: this.searchPreview.projectId,
        commitDate: this.searchPreview.commitDate ? this.searchPreview.commitDate.join(',') : '',
        deployDate: this.searchPreview.deployDate ? this.searchPreview.deployDate.join(',') : '',
        branch: this.searchPreview.branch,
        commit: this.searchPreview.commit,
        filename: this.searchPreview.filename,
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
          this.getPublishTrace()
        } else {
          this.publishLocalTraceList = []
          this.publishRemoteTraceList = {}
        }
        this.previewPagination.total = response.data.pagination.total
      }).finally(() => {
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
      // 先把projectID写入添加服务器的表单
      this.greyServerFormData.projectId = data.projectId
      this.greyServerFormData.commit = data.commit
      this.greyServerDialogVisible = true
    },

    getBranchList(data) {
      const id = data.id
      this.selectedItem = data
      this.commitDialogVisible = true
      this.branchLoading = true
      this.branch = ''
      this.branchOption = []
      this.commitTableData = []
      getBranchList(id).then(response => {
        this.branchOption = response.data.branchList.filter(element => {
          return element.indexOf('HEAD') === -1
        })
      }).finally(() => {
        this.branchLoading = false
      })
    },

    getCommitList() {
      this.commitTableLoading = true
      getCommitList(this.selectedItem.id, this.branch).then(response => {
        this.commitTableData = response.data.commitList ? response.data.commitList.map(element => {
          return Object.assign(element, {
            projectId: this.selectedItem.id,
            url: this.selectedItem.url,
            branch: this.branch
          })
        }) : []
      }).finally(() => {
        this.commitTableLoading = false
      })
    },

    getTagList(data) {
      const id = data.id
      this.tagDialogVisible = true
      this.tagTableLoading = true
      getTagList(id).then(response => {
        this.tagTableData = response.data.tagList ? response.data.tagList.map(element => {
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
            shortTag: shortTag
          })
        }) : []
      }).finally(() => {
        this.tagTableLoading = false
      })
    },

    handleAddProjectTask(data) {
      this.taskDialogVisible = true
      this.taskFormData.id = 0
      this.taskFormData.projectId = this.selectedItem.id
      this.taskFormData.branch = data.branch
      this.taskFormData.commitId = data.commit
      this.taskFormData.date = ''
    },

    getTaskList() {
      this.taskTableLoading = true
      getTaskList(this.taskPagination, this.selectedItem.id).then(response => {
        const projectTaskList = response.data.list
        this.taskTableData = projectTaskList.map(element => {
          return Object.assign(element, { projectId: this.selectedItem.id, projectName: this.selectedItem.name })
        })
        this.taskPagination.total = response.data.pagination.total
      }).finally(() => { this.taskTableLoading = false })
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

    submitTask() {
      this.$refs.taskForm.validate((valid) => {
        if (valid) {
          this.taskFormProps.disabled = true
          addTask(this.taskFormData).then(response => {
            this.$message.success('Success')
          }).finally(() => {
            this.taskFormProps.disabled = false
            this.taskDialogVisible = false
          })
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

    getReviewList() {
      this.reviewTableLoading = true
      getReviewList(this.reviewPagination, this.selectedItem.id).then(response => {
        const projectTaskList = response.data.list
        this.reviewTableData = projectTaskList.map(element => {
          return Object.assign(element, { projectId: this.selectedItem.id, projectName: this.selectedItem.name })
        })
        this.reviewPagination.total = response.data.pagination.total
      }).finally(() => { this.reviewTableLoading = false })
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
        this.$confirm(this.$i18n.t('deployPage.reviewTips'), this.$i18n.t('tips'), {
          confirmButtonText: this.$i18n.t('confirm'),
          cancelButtonText: this.$i18n.t('cancel'),
          type: 'warning'
        }).then(() => {
          review(data.id, state).then((response) => {
            const projectIndex = this.tableData.findIndex(element => element.id === data.projectId)
            this.tableData[projectIndex].state = state
            this.reviewListDialogVisible = false
          })
        }).catch(() => {
          this.$message.info('Cancel')
        })
      } else {
        review(data.id, state).then((response) => {
          const projectIndex = this.tableData.findIndex(element => element.id === data.projectId)
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

    publishByCommit(data) {
      console.log(data)
      this.$confirm(this.$i18n.t('deployPage.publishCommitTips', { commit: data.commit }), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        publish(data.projectId, data.branch, data.commit).then((response) => {
          const projectIndex = this.tableData.findIndex(element => element.id === data.projectId)
          this.tableData[projectIndex].deployState = 1
          this.commitDialogVisible = false
          this.tagDialogVisible = false
          this.dialogVisible = false
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    resetState(data) {
      this.$confirm(this.$i18n.t('deployPage.resetStateTips'), this.$i18n.t('tips'), {
        confirmButtonText: this.$i18n.t('confirm'),
        cancelButtonText: this.$i18n.t('cancel'),
        type: 'warning'
      }).then(() => {
        resetState(data.id).then((response) => {
          const projectIndex = this.tableData.findIndex(element => element.id === data.id)
          this.tableData[projectIndex].deployState = 0
          this.tableData[projectIndex].progressPercentage = 0
          this.tableData[projectIndex].tagType = 'info'
          this.tableData[projectIndex].tagText = 'Not deploy'
        })
      }).catch(() => {
        this.$message.info('Cancel')
      })
    },

    greyPublish() {
      this.$refs.greyServerForm.validate((valid) => {
        if (valid) {
          const data = this.greyServerFormData
          this.$confirm(this.$i18n.t('deployPage.publishCommitTips', { commit: data.commit }), this.$i18n.t('tips'), {
            confirmButtonText: this.$i18n.t('confirm'),
            cancelButtonText: this.$i18n.t('cancel'),
            type: 'warning'
          }).then(() => {
            greyPublish(data.projectId, data.commit, data.serverIds).then((response) => {
              const projectIndex = this.tableData.findIndex(element => element.id === data.projectId)
              this.tableData[projectIndex].deployState = 1
              this.commitDialogVisible = false
              this.tagDialogVisible = false
              this.greyServerDialogVisible = false
            })
          }).catch(() => {
            this.$message.info('Cancel')
          })
        } else {
          return false
        }
      })
    },

    enterToBR(detail) {
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

.publish-filter-label {
  font-size: 12px;
  width: 70px;
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

.icon-success {
  color:#67C23A;
  font-size:14px;
  font-weight:900;
}

.icon-fail {
  color:#F56C6C;
  font-size:14px;
  font-weight:900;
}

@media screen and (max-width: 1440px){
  .publish-record {
    >>>.el-dialog {
      width: 75%;
    }
  }
}
</style>
