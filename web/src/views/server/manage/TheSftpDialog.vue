<template>
  <el-dialog
    v-model="dialogVisible"
    title="WebSFTP"
    :close-on-click-modal="false"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-select
      v-model="serverId"
      v-loading="serverLoading"
      filterable
      default-first-option
      placeholder="please select server"
      style="width: 100%; margin-bottom: 10px"
      @change="handleSelectServer"
    >
      <el-option
        v-for="item in serverOption"
        :key="item.id"
        :label="`${item.name}(${item.description})`"
        :value="item.id"
      />
    </el-select>
    <el-input
      v-model="dir"
      style="flex: 1; margin-right: 4px"
      :disabled="!wsConnected"
      :readonly="!wsConnected"
      @keyup.enter="goto(dir)"
    >
      <template #prepend>Path</template>
      <template #append>
        <i v-if="!wsConnected" class="el-icon-loading"></i>
        <el-button
          v-if="wsConnected"
          type="primary"
          icon="el-icon-right"
          @click="goto(dir)"
        />
        <el-button
          v-if="wsConnected"
          type="primary"
          icon="el-icon-refresh-left"
          @click="back"
        />
        <el-button
          v-if="wsConnected"
          type="primary"
          icon="el-icon-refresh"
          @click="refresh"
        />
        <el-upload
          v-if="wsConnected"
          style="display: inline-block; width: 30px; text-align: right"
          :action="uploadHref(serverId)"
          :before-upload="beforeUpload"
          multiple
          :on-success="handleUploadSuccess"
          :on-error="handleUploadError"
          :show-file-list="false"
        >
          <el-button type="primary" icon="el-icon-upload" />
        </el-upload>
      </template>
    </el-input>
    <el-table
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
      :max-height="460"
    >
      <el-table-column prop="name" :label="$t('name')" min-width="100">
        <template #default="scope">
          <i v-if="scope.row.isDir" class="el-icon-folder-opened"></i>
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
          <template v-if="scope.row.uploading">
            <i class="el-icon-loading"></i>
          </template>
          <template v-else>
            <el-button
              style="margin-right: 10px"
              :disabled="!scope.row.isDir"
              type="text"
              icon="el-icon-right"
              @click="goto(`${dir}/${scope.row.name}`)"
            />
            <el-link
              :disabled="scope.row.isDir"
              :href="downloadHref(`${dir}/${scope.row.name}`)"
              target="_blank"
              :underline="false"
            >
              <i class="el-icon-download"></i>
            </el-link>
          </template>
        </template>
      </el-table-column>
    </el-table>
  </el-dialog>
</template>

<script lang="ts">
import path from 'path-browserify'
import { ElMessageBox, ElMessage } from 'element-plus'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import { computed, watch, defineComponent, ref } from 'vue'
import { ServerOption } from '@/api/server'
import { useI18n } from 'vue-i18n'
import { humanSize, parseTime } from '@/utils'
import { HttpResponse } from '@/api/types'
export interface sftpFile {
  name: string
  size: number
  mode: string
  modTime: string
  isDir: boolean
  uploading: boolean
}
export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const { t } = useI18n()
    const dialogVisible = computed({
      get: () => props.modelValue,
      set: (val) => {
        emit('update:modelValue', val)
      },
    })

    const serverLoading = ref(false)
    const serverOption = ref<ServerOption['datagram']['list']>([])
    watch(dialogVisible, (val: boolean) => {
      if (val === true) {
        serverLoading.value = true
        new ServerOption()
          .request()
          .then((response) => {
            serverOption.value = response.data.list
          })
          .finally(() => {
            serverLoading.value = false
          })
      }
    })
    const serverId = ref('')
    const tableData = ref([] as sftpFile[])
    const wsConnected = ref(false)
    let ws: WebSocket
    const handleSelectServer = () => {
      const selectedServer = serverOption.value.find(
        (_) => _.id === Number(serverId.value)
      )
      if (!selectedServer) {
        return
      }
      if (ws) {
        ws.close()
        tableData.value = []
        dirHistory = []
        dir.value = '/'
      }
      try {
        ws = new WebSocket(
          `${location.protocol.replace('http', 'ws')}//${
            window.location.host + import.meta.env.VITE_APP_BASE_API
          }/ws/sftp?serverId=${
            selectedServer.id
          }&${NamespaceKey}=${getNamespaceId()}`
        )
      } catch (error) {
        console.error(error)
      }

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
        tableLoading.value = false
        if (responseData['code'] > 0) {
          ElMessage.error(responseData['message'])
          tableData.value = []
        } else {
          tableData.value = responseData['data']
          if (dir.value !== '/') {
            tableData.value.unshift({
              name: '..',
              size: 0,
              mode: '',
              modTime: '',
              isDir: true,
              uploading: false,
            })
          }
        }
      }
    }
    const tableLoading = ref(false)
    const dir = ref('/')
    let dirHistory: string[] = []
    const goto = (target: string) => {
      tableLoading.value = true
      dir.value = path.normalize(target)
      dirHistory.push(dir.value)
      ws.send(target)
    }
    const back = () => {
      if (dirHistory.length > 1) {
        dirHistory.pop()
        dir.value = dirHistory[dirHistory.length - 1]
        ws.send(dir.value)
      }
    }
    const refresh = () => {
      tableLoading.value = true
      ws.send(dir.value)
    }
    const downloadHref = (file: string) => {
      file = path.normalize(file)
      return `${
        import.meta.env.VITE_APP_BASE_API
      }/server/downloadFile?${NamespaceKey}=${getNamespaceId()}&id=${
        serverId.value
      }&file=${file}`
    }

    const uploadHref = (serverId: string) => {
      return `${
        import.meta.env.VITE_APP_BASE_API
      }/server/uploadFile?${NamespaceKey}=${getNamespaceId()}&id=${serverId}&filePath=${
        dir.value
      }`
    }

    const handleUploadSuccess = (
      response: HttpResponse<string>,
      file: File
    ) => {
      const tableIndex = tableData.value.findIndex((_) => file.name === _.name)
      if (tableIndex >= 0) {
        if (response.code > 0) {
          ElMessage.error(`upload failed, detail: ${response.message}`)
          tableData.value.splice(tableIndex, 1)
        } else {
          tableData.value[tableIndex].uploading = false
        }
      }
      return true
    }

    const handleUploadError = (err: Error, file: File) => {
      const tableIndex = tableData.value.findIndex((_) => file.name === _.name)
      if (tableIndex >= 0) {
        ElMessage.error(`upload failed, detail: ${err.message}`)
        tableData.value.splice(tableIndex, 1)
      }
      return true
    }

    const beforeUpload = async (file: File) => {
      const tableIndex = tableData.value.findIndex((_) => file.name === _.name)
      if (tableIndex >= 0) {
        const overwriteFile = await ElMessageBox.confirm('', t('tips'), {
          message: `${file.name} is already exist, would you like to overwrite it?`,
          confirmButtonText: t('confirm'),
          cancelButtonText: t('cancel'),
          type: 'warning',
        })
          .then(() => {
            tableData.value.splice(tableIndex, 1)
            return true
          })
          .catch(() => {
            ElMessage.info(t('cancel'))
            return false
          })
        if (overwriteFile === false) {
          return Promise.reject()
        }
      }
      tableData.value.unshift({
        name: file.name,
        size: file.size,
        mode: '-rw-rw-rw-',
        modTime: parseTime(file.lastModified),
        isDir: false,
        uploading: true,
      })
      return Promise.resolve()
    }

    return {
      dialogVisible,
      serverLoading,
      serverOption,
      serverId,
      handleSelectServer,
      wsConnected,
      dir,
      tableLoading,
      tableData,
      humanSize,
      goto,
      back,
      refresh,
      downloadHref,
      uploadHref,
      handleUploadSuccess,
      handleUploadError,
      beforeUpload,
    }
  },
})
</script>
<style lang="scss" scoped>
.xterm {
  width: 100%;
  height: 100%;
}
</style>
