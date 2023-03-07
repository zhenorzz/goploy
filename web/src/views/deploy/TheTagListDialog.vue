<template>
  <el-dialog
    v-model="dialogVisible"
    title="Tag"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
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
          <div style="white-space: pre-line">
            {{ scope.row.diff }}
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="tag" label="tag">
        <template #default="scope">
          <RepoURL
            style="font-size: 12px"
            :url="projectRow.url"
            :suffix="'/tree/' + scope.row.tag"
            :text="scope.row.tag"
          >
          </RepoURL>
        </template>
      </el-table-column>
      <el-table-column prop="commit" label="commit" width="80">
        <template #default="scope">
          <RepoURL
            style="font-size: 12px"
            :url="projectRow.url"
            :suffix="'/tree/' + scope.row.commit"
            :text="scope.row.commit.substring(0, 6)"
          >
          </RepoURL>
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
      <el-table-column label="time" width="160" align="center">
        <template #default="scope">
          {{ parseTime(scope.row.timestamp) }}
        </template>
      </el-table-column>
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="160"
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
import { RepositoryTagList } from '@/api/repository'
import { parseTime } from '@/utils'
import { computed, watch, ref } from 'vue'
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  projectRow: {
    type: Object,
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
const tableLoading = ref(false)
const tableData = ref<RepositoryTagList['datagram']['list']>([])
watch(
  () => props.modelValue,
  (val: typeof props['modelValue']) => {
    if (val === true) {
      tableLoading.value = true
      new RepositoryTagList({ id: props.projectRow.id })
        .request()
        .then((response) => {
          tableData.value = response.data.list
        })
        .finally(() => {
          tableLoading.value = false
        })
    }
  }
)
</script>
