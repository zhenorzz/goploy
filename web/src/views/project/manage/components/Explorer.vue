<template>
  <el-row class="explorer-container">
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
          :disabled="dir === '/'"
          link
          :icon="Top"
          style="font-size: 14px"
          @click="dotdot"
        >
        </el-button>
      </el-row>
      <el-row class="nav-path" style="flex: 1">
        <el-input
          v-model="dir"
          placeholder="Please input absolute path"
          class="input-with-select"
          @keyup.enter="dirOpen(dir)"
        >
          <template #append>
            <el-button :icon="RefreshRight" @click="dirOpen(dir)" />
          </template>
        </el-input>
      </el-row>
    </el-row>
    <el-row class="operator" justify="space-between">
      <el-row class="operator-btn" align="middle">
        <el-button
          :disabled="dir === ''"
          :icon="Delete"
          text
          style="color: var(--el-text-color-regular); padding: 8px 5px"
          @click="deleteDir"
        >
          {{ $t('projectPage.deleteDir') }}
        </el-button>
      </el-row>
      <el-row class="nav-search" align="middle">
        <el-dropdown
          :disabled="fileFilteredList.length === 0"
          @command="handleSort"
        >
          <el-button text :icon="Sort">
            {{ $t('sort') }}
          </el-button>
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
            </el-dropdown-menu>
          </template>
        </el-dropdown>
        <el-input
          v-model="input"
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
            @dblclick="
              selectedFile['isDir'] === true &&
                dirOpen(`${dir}/${selectedFile['name']}`)
            "
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
                <el-link
                  style="padding: 5px 16px"
                  :href="previewHref"
                  target="_blank"
                  :underline="false"
                >
                  {{ $t('preview') }}
                </el-link>
              </el-dropdown-item>
              <el-dropdown-item
                v-if="selectedFile['isDir'] === false"
                style="padding: 0"
              >
                <el-link
                  style="padding: 5px 16px"
                  :href="downloadHref"
                  target="_blank"
                  :underline="false"
                >
                  {{ $t('download') }}
                </el-link>
              </el-dropdown-item>
              <el-dropdown-item @click="deleteFile">
                {{ $t('delete') }}
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>
      </div>
    </el-row>
    <el-row class="footer" justify="space-between">
      <span style="padding: 0 2px">{{ selectedFile['name'] }}</span>
      <div>
        {{ fileFilteredList.length }} {{ $t('serverPage.sftpFileCount') }}
      </div>
    </el-row>
  </el-row>
</template>
<script lang="ts" setup>
import {
  Back,
  Right,
  Top,
  RefreshRight,
  Delete,
  Sort,
} from '@element-plus/icons-vue'
import svgIds from 'virtual:svg-icons-names'
import path from 'path-browserify'
import {
  RepositoryFile,
  RepositoryFileList,
  RepositoryDeleteFile,
} from '@/api/repository'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import { ref, PropType, computed, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { ProjectData } from '@/api/project'
const { t } = useI18n()
const props = defineProps({
  project: {
    type: Object as PropType<ProjectData>,
    required: true,
  },
})

const dir = ref('')
const lastDir = ref('')
const backwardHistory = ref<string[]>([])
const forwardHistory = ref<string[]>([])
const fileListLoading = ref(false)
const fileList = ref<RepositoryFile[]>([])
const fileFilteredList = ref<RepositoryFile[]>([])
const selectedFile = ref<RepositoryFile>({} as RepositoryFile)
const input = ref('')
let fileUUID = 0
const previewHref = computed(() => {
  if (selectedFile.value == undefined) {
    return ''
  }
  const file = path.normalize(`${dir.value}/${selectedFile.value['name']}`)
  return `${
    import.meta.env.VITE_APP_BASE_API
  }/repository/previewFile?${NamespaceKey}=${getNamespaceId()}&id=${
    props.project.id
  }&file=${file}`
})

const downloadHref = computed(() => {
  if (selectedFile.value == undefined) {
    return ''
  }
  const file = path.normalize(`${dir.value}/${selectedFile.value['name']}`)
  return `${
    import.meta.env.VITE_APP_BASE_API
  }/repository/downloadFile?${NamespaceKey}=${getNamespaceId()}&id=${
    props.project.id
  }&file=${file}`
})

watch(
  [() => props.project],
  () => {
    forwardHistory.value = []
    backwardHistory.value = []
    lastDir.value = ''
    dirOpen('/')
  },
  { immediate: true }
)

function getList() {
  fileListLoading.value = true
  fileList.value = []
  fileFilteredList.value = []
  new RepositoryFileList({ id: props.project.id, dir: dir.value })
    .request()
    .then((response) => {
      fileFilteredList.value = fileList.value = response.data.list.map(
        (file) => {
          file.uuid = fileUUID++
          if (file.isDir) {
            file.icon = 'file-dir'
          } else {
            file.icon = getIcon(file.name)
          }
          return file
        }
      )
      handleSort('default')
    })
    .finally(() => {
      fileListLoading.value = false
    })
}

function goto(target: string) {
  fileListLoading.value = true
  selectedFile.value = {} as RepositoryFile
  dir.value = path.normalize(target)
  getList()
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

function dotdot() {
  if (dir.value == '/') {
    goto(dir.value)
  } else {
    goto(path.resolve(dir.value, '..'))
  }
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

function selectFile(file: RepositoryFile) {
  selectedFile.value = file
}

function filterFile(value: string) {
  fileFilteredList.value = fileList.value.filter((file) =>
    file.name.includes(value)
  )
}

function handleSort(command: string) {
  let compareFunc = (fileA: RepositoryFile, fileB: RepositoryFile): number => {
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
      compareFunc = (fileA: RepositoryFile, fileB: RepositoryFile): number => {
        return fileA.name.localeCompare(fileB.name)
      }
      break
    case 'nameDesc':
      compareFunc = (fileA: RepositoryFile, fileB: RepositoryFile): number => {
        return fileB.name.localeCompare(fileA.name)
      }
      break
  }
  fileFilteredList.value.sort(compareFunc)
}

function deleteFile() {
  const file = path.normalize(dir.value + '/' + selectedFile.value.name)
  ElMessageBox.confirm(
    t('deleteTips', { name: '${REPOSITORY_PATH}' + file }),
    t('tips'),
    {
      confirmButtonText: t('confirm'),
      cancelButtonText: t('cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      fileListLoading.value = true
      new RepositoryDeleteFile({
        id: props.project.id,
        file: file,
      })
        .request()
        .then(() => {
          const pos = fileList.value.findIndex(
            (item) => item.name === selectedFile.value.name
          )
          fileList.value.splice(pos, 1)
        })
        .finally(() => {
          fileListLoading.value = false
        })
    })
    .catch()
}

function deleteDir() {
  const file = path.normalize(dir.value)
  ElMessageBox.confirm(
    t('deleteTips', { name: '${REPOSITORY_PATH}' + file }),
    t('tips'),
    {
      confirmButtonText: t('confirm'),
      cancelButtonText: t('cancel'),
      type: 'warning',
    }
  )
    .then(() => {
      fileListLoading.value = true
      new RepositoryDeleteFile({
        id: props.project.id,
        file: file,
      })
        .request()
        .then(() => {
          dotdot()
        })
        .finally(() => {
          fileListLoading.value = false
        })
    })
    .catch()
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
.explorer-container {
  flex-direction: column;
  width: 100%;
  flex: 1;
  min-height: 1px;
  border: 1px solid var(--el-border-color);
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
      padding: 0 5px;
      font-size: 12px;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
    }
  }
  .footer {
    padding: 8px 15px;
    font-size: 13px;
    color: var(--el-text-color-regular);
  }
}
@media only screen and (max-device-width: 400px) {
  .explorer-container {
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
