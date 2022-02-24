<template>
  <el-row class="app-container main">
    <el-row class="sftp-container">
      <el-row class="nav" align="middle">
        <el-row style="margin-right: 10px">
          <el-button
            :disabled="backwardHistory.length === 0"
            type="text"
            icon="el-icon-back"
            style="color: #303133; font-size: 14px"
            @click="backward"
          >
          </el-button>
          <el-button
            :disabled="forwardHistory.length === 0"
            type="text"
            icon="el-icon-right"
            style="color: #303133; font-size: 14px"
            @click="forward"
          >
          </el-button>
          <el-button
            :disabled="!wsConnected"
            type="text"
            icon="el-icon-top"
            style="color: #303133; font-size: 14px"
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
            <template #prepend>
              <el-select
                v-model="serverId"
                placeholder="Select server"
                style="width: 140px"
                filterable
                @change="selectServer"
              >
                <el-option
                  v-for="server in serverOption"
                  :key="server.id"
                  :label="server.name"
                  :value="server.id"
                />
              </el-select>
            </template>
            <template #append>
              <el-button icon="el-icon-refresh-right" @click="refresh" />
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
            <el-button
              :disabled="dir === ''"
              type="text"
              style="color: #606266"
            >
              {{ $t('upload') }}
            </el-button>
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
            v-for="(file, index) in fileFilteredList"
            :key="index"
            :offset="-30"
            trigger="click"
            placement="right"
          >
            <el-row
              v-loading="file['uploading']"
              tabindex="1"
              class="file"
              :class="file.uuid === selectedFile['uuid'] ? 'file-selected' : ''"
              @click="selectFile(file)"
            >
              <svg-icon class="file-type" :icon-class="file.icon" />
              <div class="filename" :title="file.name">{{ file.name }}</div>
            </el-row>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item v-if="selectedFile['isDir'] === true">
                  <el-button
                    type="text"
                    @click="dirOpen(`${dir}/${selectedFile['name']}`)"
                  >
                    {{ $t('open') }}
                  </el-button>
                </el-dropdown-item>
                <el-dropdown-item v-if="selectedFile['isDir'] === false">
                  <el-link
                    :href="previewHref"
                    target="_blank"
                    :underline="false"
                    style="font-size: 12px"
                  >
                    {{ $t('preview') }}
                  </el-link>
                </el-dropdown-item>
                <el-dropdown-item v-if="selectedFile['isDir'] === false">
                  <el-link
                    :href="downloadHref"
                    target="_blank"
                    :underline="false"
                    style="font-size: 12px"
                  >
                    {{ $t('download') }}
                  </el-link>
                </el-dropdown-item>
                <el-dropdown-item divided>
                  <el-button
                    type="text"
                    @click="fileDetailDialogVisible = true"
                  >
                    {{ $t('detail') }}
                  </el-button>
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-row>
    </el-row>
    <el-dialog v-model="fileDetailDialogVisible" title="" width="250px">
      <el-descriptions title="" direction="horizontal" :column="1">
        <el-descriptions-item :label="$t('name')">
          {{ selectedFile['name'] }}
        </el-descriptions-item>
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
import svgIds from 'virtual:svg-icons-names'
import path from 'path-browserify'
import { humanSize, parseTime } from '@/utils'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import { ElMessage } from 'element-plus'
import { ServerOption } from '@/api/server'
import { HttpResponse } from '@/api/types'
import { defineComponent } from 'vue'

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

export default defineComponent({
  name: 'ServerSFTP',
  data() {
    return {
      fileDetailDialogVisible: false,
      serverOption: [] as ServerOption['datagram']['list'],
      serverId: '',
      ws: null as null | WebSocket,
      wsConnected: false,
      dir: '',
      lastDir: '',
      backwardHistory: [] as string[],
      forwardHistory: [] as string[],
      fileListLoading: false,
      fileList: [] as file[],
      fileFilteredList: [] as file[],
      selectedFile: {} as file,
      input: '',
      fileUUID: 0,
    }
  },

  computed: {
    uploadHref: function () {
      return `${
        import.meta.env.VITE_APP_BASE_API
      }/server/uploadFile?${NamespaceKey}=${getNamespaceId()}&id=${
        this.serverId
      }&filePath=${this.dir}`
    },
    previewHref: function () {
      const file = path.normalize(`${this.dir}/${this.selectedFile['name']}`)
      return `${
        import.meta.env.VITE_APP_BASE_API
      }/server/previewFile?${NamespaceKey}=${getNamespaceId()}&id=${
        this.serverId
      }&file=${file}`
    },
    downloadHref: function () {
      const file = path.normalize(`${this.dir}/${this.selectedFile['name']}`)
      return `${
        import.meta.env.VITE_APP_BASE_API
      }/server/downloadFile?${NamespaceKey}=${getNamespaceId()}&id=${
        this.serverId
      }&file=${file}`
    },
  },

  created() {
    this.getServerOption()
  },

  methods: {
    humanSize,
    getServerOption() {
      new ServerOption().request().then((response) => {
        this.serverOption = response.data.list
      })
    },
    selectServer() {
      if (this.ws) {
        this.ws.close()
        this.fileList = this.backwardHistory = this.forwardHistory = []
      }
      try {
        this.ws = new WebSocket(
          `${location.protocol.replace('http', 'ws')}//${
            window.location.host + import.meta.env.VITE_APP_BASE_API
          }/ws/sftp?serverId=${
            this.serverId
          }&${NamespaceKey}=${getNamespaceId()}`
        )
        this.ws.onopen = () => {
          this.wsConnected = true
        }
        this.ws.onerror = (error) => {
          console.log(error)
        }
        this.ws.onclose = (event) => {
          if (event.reason !== '') {
            ElMessage.error('sftp close, reason: ' + event.reason)
          }
          console.log(event)
        }
        this.ws.onmessage = (event) => {
          const responseData = JSON.parse(event.data)
          this.fileListLoading = false
          if (responseData['code'] > 0) {
            ElMessage.error(responseData['message'])
            this.fileList = []
          } else {
            this.fileList = []
            const data = responseData['data'] ? responseData['data'] : []
            this.fileFilteredList = this.fileList = data.map((file: file) => {
              file.uuid = this.fileUUID++
              if (file.isDir) {
                file.icon = 'file-dir'
              } else {
                file.icon = this.getIcon(file.name)
              }
              return file
            })
            this.handleSort('default')
          }
        }
      } catch (error) {
        console.error(error)
      }
    },

    handleSort(command: string) {
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
              new Date(fileA.modTime).getTime() -
              new Date(fileB.modTime).getTime()
            )
          }
          break
        case 'modTimeDesc':
          compareFunc = (fileA: file, fileB: file): number => {
            return (
              new Date(fileB.modTime).getTime() -
              new Date(fileA.modTime).getTime()
            )
          }
          break
      }
      this.fileFilteredList.sort(compareFunc)
    },

    goto(target: string) {
      this.fileListLoading = true
      this.selectedFile = {} as file
      this.dir = path.normalize(target)
      this.ws?.send(this.dir)
    },

    dirOpen(dir: string) {
      if (this.lastDir !== '') {
        if (
          this.backwardHistory.length === 0 ||
          this.backwardHistory[this.backwardHistory.length - 1] !== this.lastDir
        ) {
          this.backwardHistory.push(this.lastDir)
        }
      }
      this.forwardHistory = []
      this.lastDir = dir
      this.goto(dir)
    },

    dotdot(target: string) {
      this.goto(path.resolve(target, '..'))
    },

    backward() {
      const target = this.backwardHistory.pop()
      if (target) {
        this.lastDir = target
        this.forwardHistory.push(this.dir)
        this.goto(target)
      }
    },

    forward() {
      const target = this.forwardHistory.pop()
      if (target) {
        this.lastDir = target
        this.backwardHistory.push(this.dir)
        this.goto(target)
      }
    },

    refresh() {
      this.fileListLoading = true
      this.ws?.send(this.dir)
    },

    filterFile(value: string) {
      this.fileFilteredList = this.fileList.filter((file) =>
        file.name.includes(value)
      )
    },

    selectFile(file: file) {
      this.selectedFile = file
    },

    async beforeUpload(file: File) {
      const fileIndex = this.fileList.findIndex(
        (_: file) => file.name === _.name
      )
      if (fileIndex >= 0) {
        ElMessage.warning(`${file.name} is already exist`)
        return Promise.reject()
      }
      this.fileList.push({
        uuid: this.fileUUID++,
        name: file.name,
        size: file.size,
        mode: '-rw-rw-rw-',
        modTime: parseTime(file.lastModified),
        isDir: false,
        uploading: true,
        icon: this.getIcon(file.name),
      })
      return Promise.resolve()
    },

    handleUploadSuccess(response: HttpResponse<string>, file: File) {
      const fileIndex = this.fileList.findIndex(
        (_: file) => file.name === _.name
      )
      if (fileIndex >= 0) {
        if (response.code > 0) {
          ElMessage.error(`upload failed, detail: ${response.message}`)
          this.fileList.splice(fileIndex, 1)
        } else {
          this.fileList[fileIndex].uploading = false
        }
      }
      return true
    },

    handleUploadError(err: Error, file: File) {
      const fileIndex = this.fileList.findIndex(
        (_: file) => file.name === _.name
      )
      if (fileIndex >= 0) {
        ElMessage.error(`upload failed, detail: ${err.message}`)
        this.fileList.splice(fileIndex, 1)
      }
      return true
    },

    getIcon(filename: string) {
      let file_ext = path.extname(filename)
      file_ext = file_ext.length > 0 ? file_ext.substring(1) : ''
      if (svgIds.includes(`icon-file-${file_ext}`)) {
        return `file-${file_ext}`
      } else {
        return 'file-unknown'
      }
    },
  },
})
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.main {
  flex-direction: column;
  height: calc(100vh - 84px);
}
.sftp-container {
  flex-direction: column;
  border: 1px solid #d8dce5;
  height: 100%;
  .nav {
    padding: 10px;
    border-bottom: 1px solid #d8dce5;
  }
  .operator {
    padding: 5px 10px;
    border-bottom: 1px solid #d8dce5;
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
        outline: 1px solid #d8dce5;
        border-radius: 4px;
      }
    }
    .file-type {
      width: 50px;
      height: 60px;
      margin-bottom: 5px;
    }
    .filename {
      color: #303133;
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
<style>
.input-with-select .el-input-group__prepend {
  background-color: #fff;
}
</style>
