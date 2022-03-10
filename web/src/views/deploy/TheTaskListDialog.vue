<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('manage')"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-row class="app-bar" type="flex" justify="end">
      <el-button
        type="primary"
        icon="el-icon-plus"
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
      <el-table-column prop="date" :label="$t('date')" width="150" />
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
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
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
      v-model:current-page="pagination.page"
      hide-on-single-page
      :total="pagination.total"
      :page-size="pagination.rows"
      style="margin-top: 10px; text-align: right"
      background
      layout="total, prev, pager, next"
      @current-change="handlePageChange"
    />
    <template #footer>
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
  <TheCommitListDialog v-model="commitDialogVisible" :project-row="projectRow">
    <template #tableOP="scope">
      <el-popover
        :ref="
          (el) => {
            if (el) {
              taskPopoverRefs[scope.row.commit] = el
            }
          }
        "
        placement="bottom"
        trigger="click"
        width="270"
      >
        <el-date-picker
          v-model="scope.row.date"
          type="datetime"
          :disabled-date="
            (time) => time.getTime() < Date.now() - 3600 * 1000 * 24
          "
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
          <el-button v-if="!getRole().isMember()" type="primary">
            {{ $t('crontab') }}
          </el-button>
        </template>
      </el-popover>
    </template>
  </TheCommitListDialog>
</template>
<script lang="ts" setup>
import RepoURL from '@/components/RepoURL/index.vue'
import {
  ProjectData,
  ProjectTaskList,
  ProjectTaskAdd,
  ProjectTaskRemove,
  ProjectTaskData,
} from '@/api/project'
import TheCommitListDialog from './TheCommitListDialog.vue'
import { ElMessageBox, ElMessage } from 'element-plus'
import { getRole } from '@/utils/namespace'
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
    type: Object as PropType<ProjectData['datagram']>,
    required: true,
  },
})
const emit = defineEmits(['update:modelValue'])
const taskPopoverRefs = ref<Record<string, any>>({})
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val)
  },
})
watch(
  () => props.modelValue,
  (val: typeof props['modelValue']) => {
    if (val === true) {
      handlePageChange()
    }
  }
)
const tableLoading = ref(false)
const tableData = ref<ProjectTaskList['datagram']['list']>([])
const pagination = reactive({ page: 1, rows: 11, total: 0 })
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

const commitDialogVisible = ref(false)
const handleAddProjectTask = () => {
  commitDialogVisible.value = true
}

function submitTask(data: ProjectTaskData['datagram']) {
  const date = dayjs(data.date).format('YYYY-MM-DD HH:mm:ss')
  new ProjectTaskAdd({ ...data, date })
    .request()
    .then(() => {
      ElMessage.success('Success')
    })
    .finally(() => {
      taskPopoverRefs.value[data.commit].doDestroy(true)
      commitDialogVisible.value = false
      handlePageChange()
    })
}

function removeProjectTask(data: ProjectTaskData['datagram']) {
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
