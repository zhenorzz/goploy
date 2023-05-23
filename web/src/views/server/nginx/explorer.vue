<template>
  <el-row class="nginx-container">
    <el-row class="nav" align="middle">
      <el-row class="nav-path" style="flex: 1">
        <el-input
          v-model="dir"
          v-loading="dirLoading"
          placeholder="Please input nginx execution path"
          class="input-with-select"
          @keyup.enter="dirOpen(dir)"
        >
          <template #append>
            <el-button :icon="RefreshRight" @click="dirOpen(dir)" />
          </template>
        </el-input>
      </el-row>
      <el-row style="margin-left: 10px">
        <el-button
          :loading="commandLoading"
          type="warning"
          @click="handleNginxCmd(dir, 'check')"
        >
          check
        </el-button>
        <el-button
          :loading="commandLoading"
          type="success"
          @click="handleNginxCmd(dir, 'reload')"
        >
          reload
        </el-button>
      </el-row>
    </el-row>
    <el-row class="operator">
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
        <el-table
          v-show="fileFilteredList.length !== 0"
          v-loading="fileListLoading"
          height="100%"
          highlight-current-row
          :data="fileFilteredList"
        >
          <el-table-column prop="name" :label="$t('name')" min-width="120" />
          <el-table-column
            prop="modTime"
            :label="$t('updateTime')"
            min-width="80"
            align="center"
          />
          <el-table-column
            prop="operation"
            :label="$t('op')"
            min-width="50"
            align="center"
          >
            <template #default="scope">
              <Button
                type="primary"
                :loading="EditProps.editButtonLoading"
                :permissions="[permission.EditNginxConfig]"
                @click="handleEdit(scope.row)"
              >
                {{ $t('edit') }}
              </Button>
              <Button
                type="primary"
                :permissions="[permission.AddNginxConfig]"
                @click="handleCopy(scope.row)"
              >
                {{ $t('copy') }}
              </Button>
              <Button
                type="primary"
                :permissions="[permission.EditNginxConfig]"
                @click="handleRename(scope.row)"
              >
                {{ $t('rename') }}
              </Button>
              <Button
                type="danger"
                :permissions="[permission.DeleteNginxConfig]"
                @click="handleDelete(scope.row)"
              >
                {{ $t('delete') }}
              </Button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-row>
    <el-row class="footer" justify="space-between">
      <div>
        {{ fileFilteredList.length }} {{ $t('serverPage.sftpFileCount') }}
      </div>
    </el-row>
    <el-dialog
      v-model="editDialogVisible"
      :title="$t('edit')"
      :close-on-click-modal="false"
      :fullscreen="$store.state.app.device === 'mobile'"
      @close="editDialogVisible = false"
    >
      <v-ace-editor
        v-model:value="fileContent"
        lang="nginx"
        :theme="isDark ? 'one_dark' : 'github'"
        style="height: 500px; width: 100%"
      />
      <template #footer>
        <el-button
          :disabled="EditProps.disabled"
          @click="editDialogVisible = false"
        >
          {{ $t('cancel') }}
        </el-button>
        <el-button
          :loading="EditProps.disabled"
          type="primary"
          @click="editConfig"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </el-row>
</template>

<script lang="ts">
export default { name: 'NginxExplorer' }
</script>
<script lang="ts" setup>
import { RefreshRight, Sort } from '@element-plus/icons-vue'
import {
  ServerData,
  ServerNginxConfigList,
  ServerNginxData,
  NginxConfigContent,
  NginxConfigEdit,
  NginxConfigCopy,
  NginxConfigRename,
  ServerNginxPath,
  NginxConfigDelete,
} from '@/api/server'
import permission from '@/permission'
import { Button } from '@/components/Permission'
import { ElMessage, ElMessageBox } from 'element-plus'
import { ref, PropType } from 'vue'
import { useI18n } from 'vue-i18n'
import { ManageNginx } from '@/api/server'
import { useDark } from '@vueuse/core'
import { VAceEditor } from 'vue3-ace-editor'
import * as ace from 'ace-builds/src-noconflict/ace'
import path from 'path-browserify'
ace.config.set(
  'basePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
ace.config.set(
  'themePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
const { t } = useI18n()

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

const serverId = ref(props.server.id)
const dir = ref('')
const lastDir = ref('')
const fileListLoading = ref(false)
const fileList = ref<ServerNginxData[]>([])
const fileFilteredList = ref<ServerNginxData[]>([])
const input = ref('')
const fileContent = ref('')
const editDialogVisible = ref(false)
const isDark = useDark()
const EditProps = ref({
  editButtonLoading: false,
  disabled: false,
  selectedConfig: {} as ServerNginxData,
})

const commandLoading = ref(false)
const dirLoading = ref(false)

function handleSort(command: string) {
  let compareFunc = (
    fileA: ServerNginxData,
    fileB: ServerNginxData
  ): number => {
    return 0
  }
  switch (command) {
    case 'nameAsc':
      compareFunc = (
        fileA: ServerNginxData,
        fileB: ServerNginxData
      ): number => {
        return fileA.name.localeCompare(fileB.name)
      }
      break
    case 'nameDesc':
      compareFunc = (
        fileA: ServerNginxData,
        fileB: ServerNginxData
      ): number => {
        return fileB.name.localeCompare(fileA.name)
      }
      break
    case 'sizeAsc':
      compareFunc = (
        fileA: ServerNginxData,
        fileB: ServerNginxData
      ): number => {
        return fileA.size - fileB.size
      }
      break
    case 'sizeDesc':
      compareFunc = (
        fileA: ServerNginxData,
        fileB: ServerNginxData
      ): number => {
        return fileB.size - fileA.size
      }
      break
    case 'modTimeAsc':
      compareFunc = (
        fileA: ServerNginxData,
        fileB: ServerNginxData
      ): number => {
        return (
          new Date(fileA.modTime).getTime() - new Date(fileB.modTime).getTime()
        )
      }
      break
    case 'modTimeDesc':
      compareFunc = (
        fileA: ServerNginxData,
        fileB: ServerNginxData
      ): number => {
        return (
          new Date(fileB.modTime).getTime() - new Date(fileA.modTime).getTime()
        )
      }
      break
  }
  fileFilteredList.value.sort(compareFunc)
}

init()

function init() {
  dirLoading.value = true
  new ServerNginxPath({ serverId: serverId.value })
    .request()
    .then((response) => {
      if (response.data.path !== '') {
        dir.value = response.data.path
        dirOpen(response.data.path)
      }
    })
    .finally(() => {
      dirLoading.value = false
    })
}

function dirOpen(dir: string) {
  lastDir.value = dir
  emit('dir-change', dir)
  getNginxConfigList(dir)
}

function filterFile(value: string) {
  fileFilteredList.value = fileList.value.filter((file) =>
    file.name.includes(value)
  )
}

function handleCopy(data: ServerNginxData) {
  ElMessageBox.prompt('', t('copy') + ' ' + data.name, {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    inputPattern: /.+/,
    inputErrorMessage: 'Name required',
  })
    .then(({ value }) => {
      fileListLoading.value = true
      new NginxConfigCopy({
        serverId: serverId.value,
        dir: data.dir,
        srcName: data.name,
        dstName: value,
      })
        .request()
        .then(() => {
          getNginxConfigList(lastDir.value)
        })
        .finally(() => {
          fileListLoading.value = false
        })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function handleEdit(data: ServerNginxData) {
  EditProps.value.editButtonLoading = true
  EditProps.value.selectedConfig = data
  new NginxConfigContent({
    serverId: serverId.value,
    dir: data.dir,
    filename: data.name,
  })
    .request()
    .then((response) => {
      editDialogVisible.value = true
      fileContent.value = response.data.content
    })
    .finally(() => {
      EditProps.value.editButtonLoading = false
    })
}

function getNginxConfigList(dir: string) {
  fileListLoading.value = true
  fileList.value = []
  new ServerNginxConfigList({ serverId: serverId.value, dir: dir })
    .request()
    .then((response) => {
      fileFilteredList.value = fileList.value = response.data.list
    })
    .finally(() => {
      fileListLoading.value = false
    })
}

function handleNginxCmd(dir: string, command: string) {
  ElMessageBox.confirm(t('serverPage.execTips', { command }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      commandLoading.value = true
      new ManageNginx({
        serverId: serverId.value,
        path: dir,
        command,
      })
        .request()
        .then((response) => {
          ElMessage.success(
            response.data.output === '' ? 'Success' : response.data.output
          )
        })
        .finally(() => {
          commandLoading.value = false
        })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}

function editConfig() {
  const file = EditProps.value.selectedConfig
  EditProps.value.disabled = true
  new NginxConfigEdit({
    serverId: serverId.value,
    dir: file.dir,
    content: fileContent.value,
    filename: file.name,
  })
    .request()
    .then(() => {
      editDialogVisible.value = false
      ElMessage.success('Success')
    })
    .finally(() => {
      EditProps.value.disabled = false
    })
}

function handleRename(data: ServerNginxData) {
  ElMessageBox.prompt('', t('rename') + data.name, {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    inputPattern: /.+/,
    inputErrorMessage: 'Name required',
  })
    .then(({ value }) => {
      fileListLoading.value = true
      new NginxConfigRename({
        serverId: serverId.value,
        dir: data.dir,
        currentName: data.name,
        newName: value,
      })
        .request()
        .then(() => {
          ElMessage.success('Success')
          const pos = fileList.value.findIndex(
            (item) => item.name === data.name
          )
          fileList.value[pos].name = value
        })
        .finally(() => {
          fileListLoading.value = false
        })
    })
    .catch()
}

function handleDelete(data: ServerNginxData) {
  const file = path.normalize(data.dir + '/' + data.name)
  ElMessageBox.confirm(t('deleteTips', { name: file }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      fileListLoading.value = true
      new NginxConfigDelete({
        serverId: serverId.value,
        file: file,
      })
        .request()
        .then(() => {
          ElMessage.success('Success')
          const pos = fileList.value.findIndex(
            (item) => item.name === data.name
          )
          fileList.value.splice(pos, 1)
        })
        .finally(() => {
          fileListLoading.value = false
        })
    })
    .catch()
}
</script>
<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.nginx-container {
  flex-direction: column;
  width: 100%;
  flex: 1;
  min-height: 1px;
  .nav {
    padding: 10px;
    border-bottom: 1px solid var(--el-border-color);
  }
  .operator {
    padding: 5px 10px;
    border-bottom: 1px solid var(--el-border-color);
    .nav-search {
      margin-left: auto;
    }
  }
  .files {
    flex: 1;
    overflow-y: auto;
    @include scrollBar();
  }
  .footer {
    padding: 8px 15px;
    font-size: 13px;
    color: var(--el-text-color-regular);
  }
}
@media only screen and (max-device-width: 400px) {
  .nginx-container {
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
