<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('deployPage.fileCompare')"
    :fullscreen="$store.state.app.device === 'mobile'"
    :close-on-click-modal="false"
  >
    <el-row type="flex">
      <el-input
        v-model.trim="filePath"
        style="flex: 1"
        @change="handlePathChange"
      >
        <template #prepend>{$project_path}/</template>
        <template #append>
          <el-button
            type="primary"
            :icon="FolderOpened"
            @click="handleSelectPath(filePath)"
          />
        </template>
      </el-input>
      <el-button type="primary" @click="handleCompare">Compare</el-button>
    </el-row>
    <el-table
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      max-height="447px"
      style="margin-top: 10px; width: 100%"
      :data="tableData"
    >
      <el-table-column label="server">
        <template #default="scope"> {{ scope.row.serverName }} </template>
      </el-table-column>
      <el-table-column
        prop="status"
        label="Status"
        width="200"
        show-overflow-tooltip
      />
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="80"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button
            v-if="scope.row.isModified"
            type="primary"
            @click="handleDiff(scope.row)"
          >
            diff
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <template #footer>
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
  <el-dialog
    v-model="fileVisible"
    title="File"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-table
      v-loading="fileLoading"
      border
      stripe
      highlight-current-row
      :data="fileList"
      style="width: 100%"
      :max-height="460"
    >
      <el-table-column prop="name" :label="$t('name')" min-width="100">
        <template #default="scope">
          <el-icon v-if="scope.row.isDir" style="vertical-align: middle">
            <folder-opened />
          </el-icon>
          {{ scope.row.name }}
        </template>
      </el-table-column>
      <el-table-column prop="size" :label="$t('size')" width="100">
        <template #default="scope">
          {{ humanSize(scope.row.size) }}
        </template>
      </el-table-column>
      <el-table-column prop="mode" label="mode" width="100" />
      <el-table-column
        prop="modTime"
        :label="$t('modifiedTime')"
        width="160"
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
          <template v-if="scope.row.uploading">
            <el-icon><loading /></el-icon>
          </template>
          <template v-else>
            <el-button
              v-if="scope.row.isDir"
              style="margin-right: 10px"
              type="text"
              :icon="Right"
              @click="handleSelectPath(filePath + scope.row.name + '/')"
            />
            <el-button
              v-else
              style="margin-right: 10px"
              type="text"
              :icon="Check"
              @click="handleSelectFile(filePath + scope.row.name)"
            />
          </template>
        </template>
      </el-table-column>
    </el-table>
    <template #footer>
      <el-button @click="fileVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
  <el-dialog
    v-model="fileDiffVisible"
    :title="$t('deployPage.fileCompare')"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <div class="file-content">
      <div v-loading="diffLoading">
        <div
          v-for="(item, index) in changeLines"
          :key="index"
          class="file-line"
        >
          <div
            class="file-line-number"
            :style="{
              'background-color': item.type
                ? item.type === '+'
                  ? '#ccFFD8'
                  : '#ffd7d5'
                : '',
            }"
          >
            {{ item.lineNumber }}
          </div>
          <div
            class="file-line-type"
            :style="{
              'background-color': item.type
                ? item.type === '+'
                  ? '#e6ffec'
                  : '#ffebe9'
                : '',
            }"
          >
            {{ item.type }}
          </div>
          <div
            class="file-line-text"
            :style="{
              'background-color': item.type
                ? item.type === '+'
                  ? '#e6ffec'
                  : '#ffebe9'
                : '',
            }"
          >
            {{ item.text }}
          </div>
        </div>
      </div>
    </div>
    <template #footer>
      <el-button @click="fileDiffVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { FolderOpened, Check, Right, Loading } from '@element-plus/icons-vue'
import { diffLines } from 'diff'
import path from 'path-browserify'
import { FileCompare, FileDiff } from '@/api/deploy'
import { ReposFileList, ProjectData } from '@/api/project'
import { humanSize } from '@/utils'
import { computed, PropType, ref } from 'vue'
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
const filePath = ref('')
const handlePathChange = () => {
  filePath.value = path
    .normalize(filePath.value)
    .split('/')
    .filter((dir) => dir !== '..')
    .join('/')
}
const fileVisible = ref(false)
const fileLoading = ref(false)
const fileList = ref([])

const handleSelectPath = (_path: string) => {
  fileVisible.value = true
  fileLoading.value = true
  _path = _path ? path.normalize(_path) : '/'
  if (!_path.endsWith('/')) {
    _path = path.normalize(path.dirname(_path) + '/')
  }
  filePath.value = _path
  new ReposFileList({
    id: props.projectRow.id,
    path: _path,
  })
    .request()
    .then((response) => {
      fileList.value = response.data
    })
    .finally(() => {
      fileLoading.value = false
    })
}

const handleSelectFile = (_path: string) => {
  filePath.value = _path
  fileVisible.value = false
  handleCompare()
}

const tableLoading = ref(false)
const tableData = ref([])
const handleCompare = () => {
  if (filePath.value === '') {
    ElMessage.warning('file path can not be empty')
    return
  }
  tableLoading.value = true
  new FileCompare({
    projectId: props.projectRow.id,
    filePath: filePath.value,
  })
    .request()
    .then((response) => {
      tableData.value = response.data
    })
    .finally(() => {
      tableLoading.value = false
    })
}

const fileDiffVisible = ref(false)
const diffLoading = ref(false)
const changeLines = ref(
  [] as { text: string; lineNumber: string; type: string }[]
)
const handleDiff = (data: { serverId: number }) => {
  fileDiffVisible.value = true
  changeLines.value = []
  diffLoading.value = true
  new FileDiff({
    projectId: props.projectRow.id,
    serverId: data.serverId,
    filePath: filePath.value,
  })
    .request()
    .then((response) => {
      let lineNumber = 0
      diffLines(response.data.srcText, response.data.distText).forEach(
        (item) => {
          const strArr = item.value?.split('\n').filter((item) => item) || []
          const type = (item.added && '+') || (item.removed && '-') || ''
          strArr.forEach((text) => {
            const thisLineNumber = !item.removed ? ++lineNumber : ''
            changeLines.value.push({
              text,
              type,
              lineNumber: thisLineNumber.toString(),
            })
          })
        }
      )
    })
    .finally(() => {
      diffLoading.value = false
    })
}
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
@import '@/styles/mixin.scss';
.file {
  &-content {
    height: 470px;
    overflow-y: auto;
    @include scrollBar();
  }
  &-line {
    display: flex;
    flex-direction: row;
    font-size: 12px;
    line-height: 20px;
  }
  &-line-number {
    width: 50px;
    text-align: right;
    padding-right: 10px;
    padding-left: 10px;
    color: #6e7781;
  }
  &-line-type {
    width: 22px;
    text-align: center;
  }
  &-line-text {
    flex: 1;
    white-space: pre-wrap;
  }
}
</style>
