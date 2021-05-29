<template>
  <el-dialog
    v-model="dialogVisible"
    title="commit"
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
    <el-table
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      max-height="447px"
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column type="expand">
        <template #default="props">
          <div style="white-space: pre-line">
            {{ props.row.diff }}
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="commit" label="commit" width="290">
        <template #default="scope">
          <el-link
            type="primary"
            style="font-size: 12px"
            :underline="false"
            :href="`${gitURL}/commit/${scope.row.commit}`"
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
    <template #footer class="dialog-footer">
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts">
import { RepositoryBranchList, RepositoryCommitList } from '@/api/repository'
import { role } from '@/utils/namespace'
import { parseGitURL, parseTime } from '@/utils'
import { computed, watch, defineComponent, ref } from 'vue'

export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    projectRow: {
      type: Object,
      required: true,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const dialogVisible = computed({
      get: () => props.modelValue,
      set: (val) => {
        emit('update:modelValue', val)
      },
    })
    const gitURL = ref<string>('')
    const branchLoading = ref(false)
    const branchOption = ref<RepositoryBranchList['datagram']['list']>([])
    const branch = ref('')
    const tableData = ref<RepositoryCommitList['datagram']['list']>([])
    watch(
      () => props.modelValue,
      (val: typeof props['modelValue']) => {
        if (val === true) {
          gitURL.value = parseGitURL(props.projectRow.url)
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

    return {
      dialogVisible,
      role,
      gitURL,
      parseTime,
      branchLoading,
      branchOption,
      branch,
      tableLoading,
      tableData,
      getCommitList,
    }
  },
})
</script>
