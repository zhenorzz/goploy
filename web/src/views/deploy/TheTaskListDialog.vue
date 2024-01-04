<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('manage')"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <template v-if="type === 'list'">
      <el-row class="app-bar" type="flex" justify="end">
        <Button
          :permissions="[pms.DeployTask]"
          type="primary"
          :icon="Plus"
          @click="handleAddProjectTask"
        />
      </el-row>
      <el-table
        v-loading="tableLoading"
        border
        stripe
        highlight-current-row
        max-height="447px"
        :data="tableData"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column
          prop="projectName"
          :label="$t('projectName')"
          width="150"
        >
          <template #default>
            {{ projectRow.name }}
          </template>
        </el-table-column>
        <el-table-column
          prop="branch"
          :label="$t('branch')"
          width="150"
          align="center"
        >
          <template #default="scope">
            <RepoURL
              style="font-size: 12px"
              :url="projectRow.url"
              :suffix="'/tree/' + scope.row.branch.split('/').pop()"
              :text="scope.row.branch"
            >
            </RepoURL>
          </template>
        </el-table-column>
        <el-table-column prop="commit" label="commit" width="290">
          <template #default="scope">
            <RepoURL
              style="font-size: 12px"
              :url="projectRow.url"
              :suffix="`/commit/${scope.row.commit}`"
              :text="scope.row.commit"
            >
            </RepoURL>
          </template>
        </el-table-column>
        <el-table-column prop="date" :label="$t('date')" width="160" />
        <el-table-column prop="isRun" :label="$t('task')" width="80">
          <template #default="scope">
            {{ $t(`runOption[${scope.row.isRun || 0}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="state" :label="$t('state')" width="70">
          <template #default="scope">
            {{ $t(`stateOption[${scope.row.state || 0}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="creator" :label="$t('creator')" align="center" />
        <el-table-column prop="editor" :label="$t('editor')" align="center" />
        <el-table-column
          prop="insertTime"
          :label="$t('insertTime')"
          width="160"
          align="center"
        />
        <el-table-column
          prop="updateTime"
          :label="$t('updateTime')"
          width="160"
          align="center"
        />
        <el-table-column
          prop="operation"
          :label="$t('op')"
          width="80"
          align="center"
          :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
        >
          <template #default="scope">
            <Button
              :permissions="[pms.DeployTask]"
              type="danger"
              :disabled="scope.row.isRun === 1 || scope.row.state === 0"
              :icon="Delete"
              @click="removeProjectTask(scope.row)"
            />
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="pagination.page"
        hide-on-single-page
        :total="pagination.total"
        :page-size="pagination.rows"
        style="margin-top: 10px; text-align: right"
        background
        layout="total, prev, pager, next"
        @current-change="handlePageChange"
      />
    </template>
    <template v-else>
      <el-row class="app-bar" type="flex" justify="start">
        <el-button type="primary" :icon="Back" @click="handleBackList" />
      </el-row>
      <el-form
        ref="form"
        :model="formData"
        label-width="65px"
        :label-position="$store.state.app.device === 'desktop' ? 'left' : 'top'"
      >
        <el-form-item :label="$t('branch')" prop="branch">
          <el-select
            v-model="formData.branch"
            style="width: 100%"
            :loading="formProps.branchLoading"
            @change="handleBranchChange"
          >
            <el-option
              v-for="(item, index) in formProps.branchOption"
              :key="index"
              :label="index"
              :value="item"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="Commit" prop="commit">
          <el-select
            v-model="formData.commit"
            :loading="formProps.commitLoading"
            style="width: 100%"
          >
            <el-option
              v-for="(item, index) in formProps.commitOption"
              :key="index"
              :label="item.commit"
              :value="item.commit"
            >
              <el-row style="width: 100%">
                <span :title="item.commit">
                  [commit: {{ item.commit ? item.commit.substring(0, 6) : '' }}]
                </span>
                <span style="margin-left: 5px">
                  [author: {{ item.author }}]
                </span>
                <span
                  style="
                    flex: 1;
                    margin-left: 5px;
                    white-space: nowrap;
                    overflow: hidden;
                    text-overflow: ellipsis;
                  "
                  :title="item.message"
                >
                  [desc:{{ item.message }}]
                </span>
              </el-row>
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item
          :label="$t('date')"
          prop="date"
          :rules="[
            { required: true, message: 'Date required', trigger: 'blur' },
          ]"
        >
          <el-date-picker
            v-model="formData.date"
            type="datetime"
            :disabled-date="
              (time: Date) => time.getTime() < Date.now() - 3600 * 1000 * 24
            "
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
    </template>
    <template #footer>
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
      <el-button
        v-show="type === 'add'"
        :disabled="formProps.disabled"
        type="primary"
        @click="submitProjectTask"
      >
        {{ $t('confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>
<script lang="ts" setup>
import pms from '@/permission'
import { Button } from '@/components/Permission'
import { Plus, Back, Delete } from '@element-plus/icons-vue'
import RepoURL from '@/components/RepoURL/index.vue'
import {
  ProjectData,
  ProjectTaskList,
  ProjectTaskAdd,
  ProjectTaskRemove,
  ProjectTaskData,
} from '@/api/project'
import type { ElForm } from 'element-plus'
import { RepositoryBranchList, RepositoryCommitList } from '@/api/repository'
import { parseTime } from '@/utils'
import { PropType, computed, watch, ref, reactive } from 'vue'
import dayjs from 'dayjs'
import { useStore } from 'vuex'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const store = useStore()
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  projectRow: {
    type: Object as PropType<ProjectData>,
    required: true,
  },
})
const emit = defineEmits(['update:modelValue'])
watch(
  () => props.modelValue,
  (val: typeof props['modelValue']) => {
    if (val === true) {
      type.value = 'list'
      handlePageChange()
    }
  }
)
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val)
  },
})
const type = ref('list')
const tableLoading = ref(false)
const tableData = ref<ProjectTaskList['datagram']['list']>([])
const pagination = reactive({ page: 1, rows: 11, total: 0 })
const form = ref<InstanceType<typeof ElForm>>()
const formData = ref({
  branch: '',
  commit: '',
  date: '',
})
const formProps = ref({
  branchLoading: false,
  branchOption: [] as RepositoryBranchList['datagram']['list'],
  commitLoading: false,
  commitOption: [] as RepositoryCommitList['datagram']['list'],
  disabled: false,
})

const getTaskList = () => {
  tableLoading.value = true
  new ProjectTaskList({ id: props.projectRow.id }, pagination)
    .request()
    .then((response) => {
      tableData.value = response.data.list
      pagination.total = response.data.pagination.total
    })
    .finally(() => {
      tableLoading.value = false
    })
}

const handlePageChange = (page = 1) => {
  pagination.page = page
  getTaskList()
}

const handleAddProjectTask = () => {
  type.value = 'add'
  formData.value.branch = ''
  formData.value.commit = ''
  formProps.value.branchLoading = true
  formProps.value.branchOption = []
  formProps.value.commitOption = []
  new RepositoryBranchList({ id: props.projectRow.id })
    .request()
    .then((response) => {
      formProps.value.branchOption = response.data.list.filter((element) => {
        return element.indexOf('HEAD') === -1
      })
    })
    .finally(() => {
      formProps.value.branchLoading = false
    })
}
const handleBackList = () => {
  type.value = 'list'
}

const handleBranchChange = () => {
  formProps.value.commitLoading = true
  new RepositoryCommitList({
    id: props.projectRow.id,
    branch: formData.value.branch,
  })
    .request()
    .then((response) => {
      formProps.value.commitOption = response.data.list.map((element) => {
        return Object.assign(element, {
          projectId: props.projectRow.id,
          branch: formData.value.branch,
        })
      })
    })
    .finally(() => {
      formProps.value.commitLoading = false
    })
}

function submitProjectTask() {
  form.value?.validate((valid) => {
    if (valid) {
      formProps.value.disabled = true
      const date = dayjs(formData.value.date).format('YYYY-MM-DD HH:mm:ss')
      new ProjectTaskAdd({
        projectId: props.projectRow.id,
        branch: formData.value.branch,
        commit: formData.value.commit,
        date,
      })
        .request()
        .then(() => {
          ElMessage.success('Success')
          type.value = 'list'
          handlePageChange()
        })
        .finally(() => {
          formProps.value.disabled = false
        })
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
}

function removeProjectTask(data: ProjectTaskData) {
  ElMessageBox.confirm(
    t('deployPage.removeProjectTaskTips', {
      projectName: props.projectRow.name,
    }),
    t('tips'),
    {
      confirmButtonText: t('confirm'),
      cancelButtonText: t('cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      new ProjectTaskRemove({ id: data.id }).request().then(() => {
        const projectTaskIndex = tableData.value.findIndex(
          (element) => element.id === data.id
        )
        tableData.value[projectTaskIndex]['state'] = 0
        tableData.value[projectTaskIndex]['editor'] = store.getters.name
        tableData.value[projectTaskIndex]['editorId'] = store.getters.uid
        tableData.value[projectTaskIndex]['updateTime'] = parseTime(
          new Date().getTime()
        )
        ElMessage.success('The task is disabled')
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}
</script>
