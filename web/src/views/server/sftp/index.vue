<template>
  <el-row class="app-container">
    <el-row class="tabs-container">
      <el-scrollbar style="width: 100%">
        <el-row class="tabs" justify="start" align="middle">
          <div
            v-for="(item, index) in serverList"
            :key="item.uuid"
            class="tabs-item"
            :class="item.uuid === currentUUID ? 'tabs-item-selected' : ''"
          >
            <el-row>
              <div class="tabs-item-serial" @click="selectTab(item)">
                {{ index + 1 }}
              </div>
              <div
                class="tabs-item-name"
                :title="`${item.server.name} : ${item.dir}`"
                @click="selectTab(item)"
              >
                {{ item.server.name }} : {{ item.dir }}
              </div>
              <el-button
                link
                style="font-size: 14px; padding-left: 8px"
                @click="deleteTab(item, index)"
              >
                x
              </el-button>
            </el-row>
          </div>

          <div class="tabs-plus">
            <el-select
              v-model="serverId"
              style="width: 200px"
              filterable
              clearable
              @change="selectServer"
            >
              <el-option
                v-for="item in serverOption"
                :key="item.id"
                :label="item.label"
                :value="item.id"
              />
            </el-select>
          </div>
          <div class="tabs-placeholder"></div>
        </el-row>
      </el-scrollbar>
    </el-row>
    <explorer
      v-for="item in serverList"
      v-show="item.uuid === currentUUID"
      :key="item.uuid"
      :server="item.server"
      @dir-change="handleDirChange"
      @edit-file="handleEditFile"
      @transfer-file="handleTransferFile"
    ></explorer>
    <el-dialog
      v-model="transferFileDialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('serverPage.transferFile')"
      :close-on-click-modal="false"
    >
      <el-form
        ref="form"
        :model="formData"
        label-width="105px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item label="Source server">
          {{ selectedSFTP.server['name'] }}
        </el-form-item>
        <el-form-item
          :label="selectedFile['isDir'] ? 'Source dir' : 'Source file'"
        >
          {{ formData.sourceFile }}
        </el-form-item>
        <el-form-item
          label="Dest server"
          prop="destServerIds"
          :rules="[
            { required: true, message: 'Server required', trigger: 'blur' },
          ]"
        >
          <el-select
            v-model="formData.destServerIds"
            style="width: 100%"
            filterable
            multiple
          >
            <el-option
              v-for="item in serverOption"
              :key="item.id"
              :label="item.label"
              :value="item.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item
          label="Dest Dir"
          prop="destDir"
          :rules="[
            { required: true, message: 'Dest dir required', trigger: 'blur' },
          ]"
        >
          <el-input v-model="formData.destDir" autocomplete="off" />
        </el-form-item>
      </el-form>
      <span style="color: var(--el-color-danger); padding: 0 10px">
        {{ formProps.errorLog }}
      </span>
      <template #footer>
        <el-button
          :disabled="formProps.disabled"
          @click="transferFileDialogVisible = false"
        >
          {{ $t('cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="formProps.disabled"
          @click="transferFile"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
    <el-dialog
      v-model="fileEditorDialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="selectedFile['name']"
      :close-on-click-modal="false"
    >
      <VAceEditor
        ref="fileEditor"
        v-model:value="fileContent"
        v-loading="formProps.disabled"
        :lang="getModeForPath(selectedFile['name']).name"
        :theme="isDark ? 'one_dark' : 'github'"
        style="height: 360px; width: 100%"
      />
      <template #footer>
        <el-button
          :disabled="formProps.disabled"
          @click="fileEditorDialogVisible = false"
        >
          {{ $t('cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="formProps.disabled"
          @click="editFile"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </el-row>
</template>

<script lang="ts">
export default { name: 'ServerSFTP' }
</script>
<script lang="ts" setup>
import explorer from './explorer.vue'
import path from 'path-browserify'
import type { ElForm } from 'element-plus'
import { VAceEditor } from 'vue3-ace-editor'
import * as ace from 'ace-builds/src-noconflict/ace'
import { getModeForPath } from 'ace-builds/src-noconflict/ext-modelist'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import {
  ServerOption,
  ServerData,
  ServerSFTPFile,
  ServerEditFile,
  ServerTransferFile,
} from '@/api/server'
import { useDark } from '@vueuse/core'
import { ref } from 'vue'
interface sftp {
  uuid: number
  server: ServerData
  dir: string
}
ace.config.set(
  'basePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
ace.config.set(
  'themePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
const isDark = useDark()
const currentUUID = ref(0)
const transferFileDialogVisible = ref(false)
const fileEditorDialogVisible = ref(false)
const serverOption = ref<ServerOption['datagram']['list']>([])
const serverId = ref('')
const serverList = ref<sftp[]>([])
const selectedFile = ref<ServerSFTPFile>({} as ServerSFTPFile)
const selectedSFTP = ref<sftp>({} as sftp)
const form = ref<InstanceType<typeof ElForm>>()
const formData = ref({
  sourceFile: '',
  destServerIds: [],
  destDir: '',
})
const formProps = ref({
  disabled: false,
  errorLog: '',
})

const fileContent = ref('')

getServerOption()

function getServerOption() {
  new ServerOption().request().then((response) => {
    serverOption.value = response.data.list
  })
}

function selectTab(sftp: sftp) {
  currentUUID.value = sftp.uuid
  selectedSFTP.value = sftp
}

function deleteTab(sftp: sftp, index: number) {
  serverList.value.splice(index, 1)
  if (currentUUID.value === sftp.uuid) {
    currentUUID.value =
      serverList.value.length === 0
        ? 0
        : serverList.value[serverList.value.length - 1].uuid
  }
}

function selectServer(value: number) {
  const server =
    serverOption.value.find((_) => _.id === value) || ({} as ServerData)
  if (serverList.value.length === 0) {
    currentUUID.value = 0
  } else {
    currentUUID.value = serverList.value[serverList.value.length - 1].uuid + 1
  }
  const serverTab = { uuid: currentUUID.value, server, dir: '' }
  serverList.value.push(serverTab)
  selectTab(serverTab)
  serverId.value = ''
}

function handleDirChange(dir: string) {
  selectedSFTP.value.dir = dir
}

function handleEditFile(file: ServerSFTPFile) {
  selectedFile.value = file
  fileContent.value = ''
  fileEditorDialogVisible.value = true
  formProps.value.disabled = true
  const f = path.normalize(`${selectedSFTP.value.dir}/${file.name}`)
  fetch(
    `${
      import.meta.env.VITE_APP_BASE_API
    }/server/previewFile?${NamespaceKey}=${getNamespaceId()}&id=${
      selectedSFTP.value.server.id
    }&file=${f}`
  )
    .then((response) => response.text())
    .then((content) => {
      fileContent.value = content
    })
    .finally(() => {
      formProps.value.disabled = false
    })
}

function editFile() {
  formProps.value.disabled = true
  new ServerEditFile({
    serverId: selectedSFTP.value.server.id,
    file: path.normalize(
      `${selectedSFTP.value.dir}/${selectedFile.value.name}`
    ),
    content: fileContent.value,
  })
    .request()
    .then(() => {
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = false
    })
}

function handleTransferFile(file: ServerSFTPFile) {
  formData.value.destDir = selectedSFTP.value.dir
  formData.value.sourceFile = path.normalize(
    `${selectedSFTP.value.dir}/${file['name']}`
  )
  formProps.value.errorLog = ''
  selectedFile.value = file
  transferFileDialogVisible.value = true
}

function transferFile() {
  form.value?.validate((valid) => {
    if (valid) {
      formProps.value.disabled = true
      new ServerTransferFile({
        sourceServerId: selectedSFTP.value.server.id,
        sourceIsDir: selectedFile.value.isDir,
        ...formData.value,
      })
        .request()
        .then(() => {
          ElMessage.success('Success')
        })
        .catch((err) => {
          console.log(err.data)
          const data = err.data
          formProps.value.errorLog = data.data.serverName + ': ' + data.message
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
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.app-container {
  background-color: var(--el-bg-color);
  border: 1px solid var(--el-border-color);
}
.tabs-container {
  width: 100%;
  height: 46px;
  .tabs {
    height: auto;
    flex-wrap: nowrap;
    &-item {
      font-size: 14px;
      flex-shrink: 0;
      width: 195px;
      height: 46px;
      line-height: 46px;
      padding: 0 10px;
      border-bottom: 1px solid var(--el-border-color);
      background-color: var(--el-disabled-bg-color);
      &:not(:last-child) {
        border-right: 1px solid var(--el-border-color);
      }
      &-serial {
        cursor: pointer;
        color: var(--el-text-color-secondary);
        text-align: center;
        font-weight: 600;
        width: 20px;
      }
      &-name {
        flex: 1;
        cursor: pointer;
        text-overflow: ellipsis;
        overflow: hidden;
        padding-left: 5px;
        white-space: nowrap;
        display: inline-block;
      }
      &-selected {
        border-bottom: none;
        background-color: var(--el-bg-color);
      }
    }
    &-placeholder {
      height: 46px;
      flex: 1;
      border-bottom: 1px solid var(--el-border-color);
    }
    &-plus {
      text-align: center;
      border-bottom: 1px solid var(--el-border-color);
      border-right: 1px solid var(--el-border-color);
    }
  }
}
</style>
<style lang="scss">
.tabs-plus {
  .el-input__wrapper {
    border-radius: 0px !important;
    box-shadow: none;
  }
  .el-input__inner {
    height: 43px !important;
  }
}
</style>
