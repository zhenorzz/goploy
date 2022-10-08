<template>
  <el-row class="sftp-container">
    <el-row class="nav" align="middle">
      <el-row style="margin-right: 10px">
        <el-button
          :disabled="backwardHistory.length === 0"
          link
          :icon="Back"
          style="font-size: 14px"
          @click="backward"
        >
        </el-button>
        <el-button
          :disabled="forwardHistory.length === 0"
          link
          :icon="Right"
          style="font-size: 14px"
          @click="forward"
        >
        </el-button>
        <el-button
          :disabled="!wsConnected"
          link
          :icon="Top"
          style="font-size: 14px"
          @click="dotdot(dir)"
        >
        </el-button>
      </el-row>
      <el-row class="nav-path" style="flex: 1">
        <el-input
          v-model="dir"
          :disabled="!wsConnected"
          :readonly="!wsConnected"
          placeholder="Please input absolute path"
          class="input-with-select"
          @keyup.enter="dirOpen(dir)"
        >
          <template #append>
            <el-button :icon="RefreshRight" @click="refresh" />
          </template>
        </el-input>
      </el-row>
    </el-row>
    <el-row class="operator" justify="space-between">
      <el-row class="operator-btn" align="middle">
        <el-upload
          :disabled="dir === ''"
          style="display: inline-block; width: 40px"
          :action="uploadHref"
          :before-upload="beforeUpload"
          multiple
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :show-file-list="false"
        >
          <Button
            :disabled="dir === ''"
            link
            style="color: var(--el-text-color-regular)"
            :permissions="[permission.SFTPUploadFile]"
          >
            {{ $t('upload') }}
          </Button>
        </el-upload>
      </el-row>
      <el-row class="nav-search" align="middle">
        <el-dropdown
          :disabled="fileFilteredList.length === 0"
          @command="handleSort"
        >
          <span style="font-size: 12px; margin-right: 10px">
            {{ $t('sort') }}
          </span>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="default">
                {{ $t('default') }}
              </el-dropdown-item>
              <el-dropdown-item command="nameAsc">
                {{ $t('name') }} {{ $t('asc') }}
              </el-dropdown-item>
              <el-dropdown-item command="nameDesc">
                {{ $t('name') }} {{ $t('desc') }}
              </el-dropdown-item>
              <el-dropdown-item command="sizeAsc">
                {{ $t('size') }} {{ $t('asc') }}
              </el-dropdown-item>
              <el-dropdown-item command="sizeDesc">
                {{ $t('size') }} {{ $t('desc') }}
              </el-dropdown-item>
              <el-dropdown-item command="modTimeAsc">
                {{ $t('modifiedTime') }} {{ $t('asc') }}
              </el-dropdown-item>
              <el-dropdown-item command="modTimeDesc">
                {{ $t('modifiedTime') }} {{ $t('desc') }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-input
          v-model="input"
          :disabled="!wsConnected"
          :readonly="!wsConnected"
          placeholder="Filter file"
          style="flex: 1"
          @input="filterFile"
        />
      </el-row>
    </el-row>
    <el-row v-loading="fileListLoading" class="files">
      <div style="width: 100%">
        <el-empty
          v-show="fileFilteredList.length === 0"
          description="No result"
        ></el-empty>
        <el-dropdown
          v-for="(item, index) in fileFilteredList"
          :key="index"
          trigger="click"
          placement="right"
        >
          <el-row
            v-loading="item['uploading']"
            tabindex="1"
            class="file"
            :class="item.uuid === selectedFile['uuid'] ? 'file-selected' : ''"
            @click="selectFile(item)"
          >
            <svg-icon class="file-type" :icon-class="item.icon" />
            <div class="filename" :title="item.name">{{ item.name }}</div>
          </el-row>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item
                v-if="selectedFile['isDir'] === true"
                @click="dirOpen(`${dir}/${selectedFile['name']}`)"
              >
                {{ $t('open') }}
              </el-dropdown-item>
              <el-dropdown-item
                v-if="selectedFile['isDir'] === false"
                style="padding: 0"
              >
                <Link
                  style="padding: 5px 16px"
                  :href="previewHref"
                  target="_blank"
                  :underline="false"
                  :permissions="[permission.SFTPPreviewFile]"
                >
                  {{ $t('preview') }}
                </Link>
              </el-dropdown-item>
              <el-dropdown-item
                v-if="selectedFile['isDir'] === false"
                style="padding: 0"
              >
                <Link
                  style="padding: 5px 16px"
                  :href="downloadHref"
                  target="_blank"
                  :underline="false"
                  :permissions="[permission.SFTPDownloadFile]"
                >
                  {{ $t('download') }}
                </Link>
              </el-dropdown-item>
              <el-dropdown-item divided @click="fileDetailDialogVisible = true">
                {{ $t('detail') }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-row>
    <el-dialog
      v-model="fileDetailDialogVisible"
      :title="selectedFile['name']"
      custom-class="file-detail"
      width="250px"
    >
      <el-descriptions title="" direction="horizontal" :column="1">
        <el-descriptions-item :label="$t('type')">
          {{ selectedFile['isDir'] === true ? 'dir' : 'file' }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('permission')">
          {{ selectedFile['mode'] }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('size')">
          {{ humanSize(selectedFile['size']) }}
        </el-descriptions-item>
        <el-descriptions-item :label="$t('modifiedTime')">
          {{ selectedFile['modTime'] }}
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </el-row>
</template>

<script lang="ts">
export default { name: 'SFTPExplorer' }
</script>
<script lang="ts" setup>
import permission from '@/permission'
import { Button, Link } from '@/components/Permission'
import { Back, Right, Top, RefreshRight } from '@element-plus/icons-vue'
import svgIds from 'virtual:svg-icons-names'
import path from 'path-browserify'
import { humanSize, parseTime } from '@/utils'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import type { ElUpload } from 'element-plus'
import { ServerData } from '@/api/server'
import { ElMessage, ElMessageBox } from 'element-plus'
import { HttpResponse } from '@/api/types'
import { ref, PropType, computed, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
interface file {
  uuid: number
  isDir: boolean
  modTime: string
  mode: string
  name: string
  size: number
  icon: string
  uploading: boolean
}
const emit = defineEmits(['dir-change'])
const props = defineProps({
  uuid: {
    type: Number,
    default: 0,
  },
  server: {
    type: Object as PropType<ServerData>,
    required: true,
  },
})

let ws: WebSocket
const fileDetailDialogVisible = ref(false)
const serverId = ref(props.server.id)
const wsConnected = ref(false)
const dir = ref('')
const lastDir = ref('')
const backwardHistory = ref<string[]>([])
const forwardHistory = ref<string[]>([])
const fileListLoading = ref(false)
const fileList = ref<file[]>([])
const fileFilteredList = ref<file[]>([])
const selectedFile = ref<file>({} as file)
const input = ref('')
connectServer()
onBeforeUnmount(() => {
  ws?.close()
})

let fileUUID = 0
const uploadHref = computed(() => {
  return `${
    import.meta.env.VITE_APP_BASE_API
  }/server/uploadFile?${NamespaceKey}=${getNamespaceId()}&id=${
    serverId.value
  }&filePath=${dir.value}`
})

const previewHref = computed(() => {
  if (selectedFile.value == undefined) {
    return ''
  }
  const file = path.normalize(`${dir.value}/${selectedFile.value['name']}`)
  return `${
    import.meta.env.VITE_APP_BASE_API
  }/server/previewFile?${NamespaceKey}=${getNamespaceId()}&id=${
    serverId.value
  }&file=${file}`
})

const downloadHref = computed(() => {
  if (selectedFile.value == undefined) {
    return ''
  }
  const file = path.normalize(`${dir.value}/${selectedFile.value['name']}`)
  return `${
    import.meta.env.VITE_APP_BASE_API
  }/server/downloadFile?${NamespaceKey}=${getNamespaceId()}&id=${
    serverId.value
  }&file=${file}`
})

function connectServer() {
  if (ws) {
    ws.close()
    fileList.value = backwardHistory.value = forwardHistory.value = []
  }
  try {
    ws = new WebSocket(
      `${location.protocol.replace('http', 'ws')}//${
        window.location.host + import.meta.env.VITE_APP_BASE_API
      }/ws/sftp?serverId=${serverId.value}&${NamespaceKey}=${getNamespaceId()}`
    )
    ws.onopen = () => {
      wsConnected.value = true
    }
    ws.onerror = (error) => {
      console.log(error)
    }
    ws.onclose = (event) => {
      if (event.reason !== '') {
        ElMessage.error('sftp close, reason: ' + event.reason)
      }
      console.log(event)
    }
    ws.onmessage = (event) => {
      const responseData = JSON.parse(event.data)
      fileListLoading.value = false
      if (responseData['code'] > 0) {
        ElMessage.error(responseData['message'])
        fileList.value = []
      } else {
        fileList.value = []
        const data = responseData['data'] ? responseData['data'] : []
        fileFilteredList.value = fileList.value = data.map((file: file) => {
          file.uuid = fileUUID++
          if (file.isDir) {
            file.icon = 'file-dir'
          } else {
            file.icon = getIcon(file.name)
          }
          return file
        })
        handleSort('default')
      }
    }
  } catch (error) {
    console.error(error)
  }
}

function handleSort(command: string) {
  let compareFunc = (fileA: file, fileB: file): number => {
    if (fileA.isDir > fileB.isDir) {
      // 按某种排序标准进行比较, a 小于 b
      return -1
    }
    if (fileA.isDir < fileB.isDir) {
      return 1
    }
    // a must be equal to b
    return 0
  }
  switch (command) {
    case 'nameAsc':
      compareFunc = (fileA: file, fileB: file): number => {
        return fileA.name.localeCompare(fileB.name)
      }
      break
    case 'nameDesc':
      compareFunc = (fileA: file, fileB: file): number => {
        return fileB.name.localeCompare(fileA.name)
      }
      break
    case 'sizeAsc':
      compareFunc = (fileA: file, fileB: file): number => {
        return fileA.size - fileB.size
      }
      break
    case 'sizeDesc':
      compareFunc = (fileA: file, fileB: file): number => {
        return fileB.size - fileA.size
      }
      break
    case 'modTimeAsc':
      compareFunc = (fileA: file, fileB: file): number => {
        return (
          new Date(fileA.modTime).getTime() - new Date(fileB.modTime).getTime()
        )
      }
      break
    case 'modTimeDesc':
      compareFunc = (fileA: file, fileB: file): number => {
        return (
          new Date(fileB.modTime).getTime() - new Date(fileA.modTime).getTime()
        )
      }
      break
  }
  fileFilteredList.value.sort(compareFunc)
}

function goto(target: string) {
  fileListLoading.value = true
  selectedFile.value = {} as file
  dir.value = path.normalize(target)
  ws?.send(dir.value)
  emit('dir-change', target)
}

function dirOpen(dir: string) {
  if (lastDir.value !== '') {
    if (
      backwardHistory.value.length === 0 ||
      backwardHistory.value[backwardHistory.value.length - 1] !== lastDir.value
    ) {
      backwardHistory.value.push(lastDir.value)
    }
  }
  forwardHistory.value = []
  lastDir.value = dir
  goto(dir)
}

function dotdot(target: string) {
  goto(path.resolve(target, '..'))
}

function backward() {
  const target = backwardHistory.value.pop()
  if (target) {
    lastDir.value = target
    forwardHistory.value.push(dir.value)
    goto(target)
  }
}

function forward() {
  const target = forwardHistory.value.pop()
  if (target) {
    lastDir.value = target
    backwardHistory.value.push(dir.value)
    goto(target)
  }
}

function refresh() {
  fileListLoading.value = true
  ws?.send(dir.value)
}

function filterFile(value: string) {
  fileFilteredList.value = fileList.value.filter((file) =>
    file.name.includes(value)
  )
}

function selectFile(file: file) {
  selectedFile.value = file
}

async function beforeUpload(file: File) {
  const fileIndex = fileList.value.findIndex((_: file) => file.name === _.name)
  if (fileIndex >= 0) {
    const result = await ElMessageBox.confirm(
      `${file.name} is already exist, overwrite?`,
      t('tips'),
      {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning',
      }
    )
    if (result !== 'confirm') {
      return Promise.reject()
    }
    fileList.value.splice(fileIndex, 1)
  }
  fileList.value.push({
    uuid: fileUUID++,
    name: file.name,
    size: file.size,
    mode: '-rw-rw-rw-',
    modTime: parseTime(file.lastModified),
    isDir: false,
    uploading: true,
    icon: getIcon(file.name),
  })
  return Promise.resolve()
}

const handleUploadSuccess: InstanceType<typeof ElUpload>['onSuccess'] =
  function (response: HttpResponse<unknown>, file) {
    const fileIndex = fileList.value.findIndex(
      (_: file) => file.name === _.name
    )
    if (fileIndex >= 0) {
      if (response.code > 0) {
        ElMessage.error(`upload failed, detail: ${response.message}`)
        fileList.value.splice(fileIndex, 1)
      } else {
        fileList.value[fileIndex].uploading = false
      }
    }
    return true
  }

const handleUploadError: InstanceType<typeof ElUpload>['onError'] = function (
  err: Error,
  file
) {
  const fileIndex = fileList.value.findIndex((_: file) => file.name === _.name)
  if (fileIndex >= 0) {
    ElMessage.error(`upload failed, detail: ${err.message}`)
    fileList.value.splice(fileIndex, 1)
  }
  return true
}

function getIcon(filename: string) {
  let file_ext = path.extname(filename)
  file_ext = file_ext.length > 0 ? file_ext.substring(1) : ''
  if (svgIds.includes(`icon-file-${file_ext}`)) {
    return `file-${file_ext}`
  } else {
    return 'file-unknown'
  }
}
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.sftp-container {
  flex-direction: column;
  width: 100%;
  flex: 1;
  .nav {
    padding: 10px;
    border-bottom: 1px solid var(--el-border-color);
  }
  .operator {
    padding: 5px 10px;
    border-bottom: 1px solid var(--el-border-color);
  }
  .files {
    flex: 1;
    overflow-y: auto;
    @include scrollBar();
    .file {
      text-align: center;
      display: inline-block;
      margin: 10px;
      height: 90px;
      width: 100px;
      cursor: pointer;
      &-selected {
        outline: 1px solid var(--el-border-color);
        border-radius: 4px;
      }
    }
    .file-type {
      width: 50px;
      height: 60px;
      margin-bottom: 5px;
    }
    .filename {
      font-size: 12px;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
    }
  }
}
@media only screen and (max-device-width: 400px) {
  .sftp-container {
    .nav {
      flex-direction: column;
      .el-row {
        margin-top: 10px;
      }
    }
  }
}
</style>
<style lang="scss">
.file-detail .el-dialog__body {
  padding: 5px 20px;
}
</style>
