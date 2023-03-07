<template>
  <el-dialog
    v-model="dialogVisible"
    title="Commit"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-select
      v-model="branch"
      v-loading="branchLoading"
      filterable
      default-first-option
      placeholder="please select branch"
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
    <slot></slot>
    <el-table
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      max-height="447px"
      :data="tableData"
    >
      <el-table-column type="expand">
        <template #default="scope">
          <div style="padding-left: 50px; white-space: pre-line">
            {{ scope.row.diff }}
          </div>
        </template>
      </el-table-column>
      <el-table-column
        prop="message"
        width="155"
        label="message"
        show-overflow-tooltip
      />
      <el-table-column prop="commit" label="commit" width="290">
        <template #default="scope">
          <RepoURL
            style="font-size: 12px"
            :url="projectRow.url"
            :suffix="`/commit/${scope.row.commit}`"
            :text="scope.row.commit.substring(0, 6)"
          >
          </RepoURL>
        </template>
      </el-table-column>
      <el-table-column prop="author" width="155" label="author" />
      <el-table-column label="time" width="160" align="center">
        <template #default="scope">
          {{ parseTime(scope.row.timestamp) }}
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
          <slot
            name="tableOP"
            :row="scope.row"
            :$index="scope.$index"
            :column="scope.column"
          ></slot>
        </template>
      </el-table-column>
    </el-table>
    <template #footer>
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
</template>
<script lang="ts" setup>
import RepoURL from '@/components/RepoURL/index.vue'
import { RepositoryBranchList, RepositoryCommitList } from '@/api/repository'
import { ProjectData } from '@/api/project'
import { parseTime } from '@/utils'
import { PropType, computed, watch, ref } from 'vue'
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
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val)
  },
})
const branchLoading = ref(false)
const branchOption = ref<RepositoryBranchList['datagram']['list']>([])
const branch = ref('')
const tableData = ref<RepositoryCommitList['datagram']['list']>([])
watch(
  () => props.modelValue,
  (val: typeof props['modelValue']) => {
    if (val === true) {
      branchLoading.value = true
      branchOption.value = []
      branch.value = ''
      tableData.value = []
      new RepositoryBranchList({ id: props.projectRow.id })
        .request()
        .then((response) => {
          branchOption.value = response.data.list.filter((element) => {
            return element.indexOf('HEAD') === -1
          })
        })
        .finally(() => {
          branchLoading.value = false
        })
    }
  }
)

const tableLoading = ref(false)
const getCommitList = () => {
  tableLoading.value = true
  new RepositoryCommitList({
    id: props.projectRow.id,
    branch: branch.value,
  })
    .request()
    .then((response) => {
      tableData.value = response.data.list.map((element) => {
        return Object.assign(element, {
          projectId: props.projectRow.id,
          branch: branch.value,
        })
      })
    })
    .finally(() => {
      tableLoading.value = false
    })
}
</script>
