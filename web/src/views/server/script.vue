<template>
  <el-row class="app-container">
    <el-col :span="14">
      <el-row style="height: 100%; width: 100%; flex-direction: column">
        <el-row justify="space-between" align="middle">
          <el-row align="middle">
            <el-select v-model="scriptMode" style="width: 90px">
              <el-option
                v-for="(item, index) in scriptLangOption"
                :key="index"
                :label="item.label"
                :value="item.value"
              />
            </el-select>
            <el-select
              v-model="templateName"
              style="flex: 1"
              filterable
              clearable
            >
              <el-option value="Option1" label="Option1"></el-option>
            </el-select>
            <el-button>保存模板</el-button>
          </el-row>
          <el-button @click="handleRun">
            {{ $t('run') }}
            <el-icon class="el-icon--right"><CaretRight /></el-icon>
          </el-button>
        </el-row>
        <v-ace-editor
          v-model:value="script"
          :lang="getScriptLang('')"
          :theme="isDark ? 'one_dark' : 'github'"
          :options="{
            showPrintMargin: false,
          }"
          style="height: 100%; flex: 1; width: 100%"
        />
      </el-row>
    </el-col>
    <el-col :span="10">
      <el-tabs type="border-card" style="height: 100%; width: 100%">
        <el-tab-pane
          v-for="(window, serverId) in serverWindows"
          :key="serverId"
          :label="window.server.name"
        >
          {{ window.server.name }}
        </el-tab-pane>
      </el-tabs>
    </el-col>
    <el-dialog
      v-model="dialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('serverPage.transferFile')"
      :close-on-click-modal="false"
    >
      <el-form
        ref="form"
        :rules="formRules"
        :model="formData"
        label-width="105px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item :label="$t('server')" prop="serverIds">
          <el-select
            v-model="formData.serverIds"
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
        <el-form-item :label="$t('serverPage.execDir')" prop="dir">
          <el-input v-model="formData.dir" autocomplete="off" />
        </el-form-item>
        <el-form-item :label="$t('command')" prop="command">
          <el-input v-model="formData.command" autocomplete="off" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button
          :disabled="formProps.disabled"
          @click="dialogVisible = false"
        >
          {{ $t('cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="formProps.disabled"
          @click="handleExec"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </el-row>
</template>
<script lang="ts">
export default { name: 'ServerScript' }
</script>
<script lang="ts" setup>
import { CaretRight, Plus } from '@element-plus/icons-vue'
import { VAceEditor } from 'vue3-ace-editor'
import type { ElForm } from 'element-plus'
import * as ace from 'ace-builds/src-noconflict/ace'
import { ServerData, ServerOption, ServerExecScript } from '@/api/server'
import { useDark } from '@vueuse/core'
import { ref } from 'vue'
ace.config.set(
  'basePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
ace.config.set(
  'themePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
interface serverWindow {
  server: ServerData
  output: string
}
const isDark = useDark()
const dialogVisible = ref(true)
const serverOption = ref<ServerOption['datagram']['list']>([])
const serverWindows = ref<Record<number, serverWindow>>({})
const templateName = ref('')
const scriptMode = ref('bash')
const script = ref('')
const scriptLangOption = [
  { label: 'sh', value: 'sh', lang: 'sh' },
  { label: 'zsh', value: 'zsh', lang: 'sh' },
  { label: 'bash', value: 'bash', lang: 'sh' },
  { label: 'python', value: 'python', lang: 'python' },
  { label: 'php', value: 'php', lang: 'php' },
  { label: 'bat', value: 'cmd', lang: 'batchfile' },
]

const form = ref<InstanceType<typeof ElForm>>()
const formData = ref({
  serverIds: [],
  dir: '/tmp',
  command: 'bash ${filename}',
})
const formProps = ref({
  disabled: false,
  errorLog: '',
})
const formRules: InstanceType<typeof ElForm>['rules'] = {
  serverIds: [{ required: true, message: 'Server required', trigger: 'blur' }],
  dir: [{ required: true, message: 'Dir required', trigger: 'blur' }],
  command: [{ required: true, message: 'Command required', trigger: 'blur' }],
}

getServerOption()

function getServerOption() {
  new ServerOption().request().then((response) => {
    serverOption.value = response.data.list
  })
}

function handleRun() {
  dialogVisible.value = true
}

function handleExec() {
  form.value?.validate((valid) => {
    if (valid) {
      new ServerExecScript(formData.value)
        .request()
        .then(() => {
          serverWindows.value = {}
          for (const serverId of formData.value.serverIds) {
            serverWindows.value[serverId] = {
              server:
                serverOption.value.find((_) => _.id === serverId) ||
                ({} as ServerData),
              output: '',
            }
          }
        })
        .finally(() => {
          formProps.value.disabled = false
        })

      dialogVisible.value = false
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
}

function getScriptLang(scriptMode = '') {
  if (scriptMode !== '') {
    const scriptInfo = scriptLangOption.find(
      (elem) => elem.value === scriptMode
    )
    return scriptInfo ? scriptInfo['lang'] : ''
  } else {
    return 'sh'
  }
}
</script>

<style lang="scss" scoped>
@import '@/styles/mixin.scss';
.app-container {
  background-color: var(--el-bg-color);
  flex-direction: row;
}
</style>
