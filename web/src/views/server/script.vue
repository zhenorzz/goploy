<template>
  <el-row class="app-container">
    <el-row style="height: 100%; width: 100%; flex-direction: column">
      <el-row justify="space-between" align="middle">
        <el-row align="middle">
          <el-select
            v-model="serverIds"
            collapse-tags
            collapse-tags-tooltip
            multiple
            filterable
            clearable
            @change="serverChange"
          >
            <el-option :value="0" label="Select all" />
            <el-option
              v-for="(item, index) in serverOption"
              :key="index"
              :label="item.label"
              :value="item.id"
            />
          </el-select>
          <el-button
            :loading="running"
            :disabled="script == '' || serverIds.length == 0"
            @click="handleRun"
          >
            {{ $t('run') }}
            <el-icon class="el-icon--right"><CaretRight /></el-icon>
          </el-button>
        </el-row>
        <el-row align="middle">
          <el-select
            v-model="templateName"
            filterable
            clearable
            @change="selectTemplate"
          >
            <el-option
              v-for="(item, index) in templateOption"
              :key="index"
              :label="item.name"
              :value="item.id"
            >
              <el-row align="middle" justify="space-between">
                <div
                  style="
                    max-width: 200px;
                    white-space: nowrap;
                    overflow: hidden;
                    text-overflow: ellipsis;
                  "
                  :title="item.name"
                >
                  {{ item.name }}
                </div>
                <el-row align="middle">
                  <el-tooltip
                    v-if="item.description != ''"
                    effect="dark"
                    trigger="click"
                    placement="bottom"
                  >
                    <template #content>
                      <div style="max-width: 200px">{{ item.description }}</div>
                    </template>
                    <el-icon @click.stop=""><InfoFilled /></el-icon>
                  </el-tooltip>
                  <el-icon
                    style="margin-left: 10px"
                    @click.stop="templateRemove(item)"
                  >
                    <CircleCloseFilled />
                  </el-icon>
                </el-row>
              </el-row>
            </el-option>
          </el-select>
          <el-button
            :disabled="script == ''"
            @click="templateDialogVisible = true"
          >
            {{ $t('serverPage.saveTemplate') }}
          </el-button>
          <el-popover
            placement="bottom-end"
            :title="$t('projectPage.predefinedVar')"
            width="400"
            trigger="hover"
          >
            <div>
              <el-row>
                <span>${SERVER_ID}: </span>
                <span>server.id</span>
              </el-row>
              <el-row>
                <span>${SERVER_NAME}: </span>
                <span>server.name</span>
              </el-row>
              <el-row>
                <span>${SERVER_IP}: </span>
                <span>server.ip</span>
              </el-row>
              <el-row>
                <span>${SERVER_PORT}: </span>
                <span>server.port</span>
              </el-row>
              <el-row>
                <span>${SERVER_OWNER}: </span>
                <span>server.owner</span>
              </el-row>
              <el-row>
                <span>${SERVER_PASSWORD}: </span>
                <span>server.password</span>
              </el-row>
              <el-row>
                <span>${SERVER_PATH}: </span>
                <span>server.path</span>
              </el-row>
              <el-row>
                <span>${SERVER_JUMP_IP}: </span>
                <span>server.jump_ip</span>
              </el-row>
              <el-row>
                <span>${SERVER_JUMP_PORT}: </span>
                <span>server.jump_port</span>
              </el-row>
              <el-row>
                <span>${SERVER_JUMP_OWNER}: </span>
                <span>server.jump_owner</span>
              </el-row>
              <el-row>
                <span>${SERVER_JUMP_PASSWORD}: </span>
                <span>server.jump_password</span>
              </el-row>
              <el-row>
                <span>${SERVER_JUMP_PATH}: </span>
                <span>server.jump_path</span>
              </el-row>
            </div>
            <template #reference>
              <el-button>
                {{ $t('projectPage.predefinedVar') }}
              </el-button>
            </template>
          </el-popover>
        </el-row>
      </el-row>
      <el-row v-show="serverIds.length > 0" style="margin: 10px 0">
        <el-tag
          v-for="(_tag, serverId) in serverTags"
          :key="serverId"
          :type="
            _tag.isRun == false
              ? ''
              : _tag.execRes == true
              ? 'success'
              : 'danger'
          "
          closable
          style="padding: 14px 8px; margin-right: 5px"
          @click="showExecRes(_tag)"
          @close="closeTag(serverId)"
        >
          {{ _tag.server.label }}
        </el-tag>
      </el-row>
      <v-ace-editor
        v-model:value="script"
        lang="sh"
        :theme="isDark ? 'one_dark' : 'github'"
        :options="{
          showPrintMargin: false,
        }"
        style="height: 100%; flex: 1; width: 100%"
      />
    </el-row>
    <el-dialog
      v-model="execDialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :close-on-click-modal="false"
    >
      <el-tabs v-model="serverTab" type="border-card">
        <el-tab-pane
          v-for="(_tag, serverId) in serverTags"
          :key="serverId"
          :name="_tag.server.label"
          lazy
        >
          <template #label>
            <el-row align="middle">
              {{ _tag.server.label }}
              <el-icon
                v-if="_tag.execRes == true"
                style="margin-left: 5px; color: var(--el-color-success)"
              >
                <CircleCheck />
              </el-icon>
              <el-icon
                v-else
                style="margin-left: 5px; color: var(--el-color-danger)"
              >
                <CircleClose />
              </el-icon>
            </el-row>
          </template>
          stdout:
          <pre>{{ _tag.stdout }}</pre>
          stderr:
          <pre>{{ _tag.stderr }}</pre>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
    <el-dialog
      v-model="templateDialogVisible"
      :fullscreen="$store.state.app.device === 'mobile'"
      :title="$t('serverPage.saveTemplate')"
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
        <el-form-item
          :label="$t('name')"
          prop="name"
          :rules="[
            { required: true, message: 'Name required', trigger: 'blur' },
          ]"
        >
          <el-input v-model="formData.name" autocomplete="off" />
        </el-form-item>
        <el-form-item
          :label="$t('description')"
          prop="description"
          :rules="[
            { max: 2047, message: 'Max 2047 characters', trigger: 'blur' },
          ]"
        >
          <el-input
            v-model="formData.description"
            type="textarea"
            autocomplete="off"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button
          :disabled="formProps.disabled"
          @click="templateDialogVisible = false"
        >
          {{ $t('cancel') }}
        </el-button>
        <el-button
          type="primary"
          :loading="formProps.disabled"
          @click="handleSaveTemplate"
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
import {
  InfoFilled,
  CircleCloseFilled,
  CaretRight,
  CircleCheck,
  CircleClose,
} from '@element-plus/icons-vue'
import { VAceEditor } from 'vue3-ace-editor'
import type { ElForm } from 'element-plus'
import * as ace from 'ace-builds/src-noconflict/ace'
import { ServerData, ServerOption, ServerExecScript } from '@/api/server'
import {
  TemplateType,
  TemplateData,
  TemplateOption,
  TemplateAdd,
  TemplateRemove,
} from '@/api/template'
import { useDark } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
import { ref } from 'vue'
ace.config.set(
  'basePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
ace.config.set(
  'themePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
interface ServerTag {
  server: ServerData
  isRun: boolean
  execRes: boolean
  stdout: string
  stderr: string
}
const isDark = useDark()
const templateDialogVisible = ref(false)
const execDialogVisible = ref(false)
const serverOption = ref<ServerOption['datagram']['list']>([])
const templateOption = ref<TemplateOption['datagram']['list']>([])
const serverTab = ref('')
const serverTags = ref<Record<number, ServerTag>>({})
const templateName = ref('')
const script = ref('')
const serverIds = ref<number[]>([])
const running = ref(false)
const saving = ref(false)
const form = ref<InstanceType<typeof ElForm>>()
const formData = ref({
  name: '',
  description: '',
})
const formProps = ref({
  disabled: false,
  errorLog: '',
})

const { t } = useI18n()
getServerOption()
function getServerOption() {
  new ServerOption().request().then((response) => {
    serverOption.value = response.data.list
  })
}

getTemplateOption()
function getTemplateOption() {
  new TemplateOption({
    type: TemplateType.Script,
  })
    .request()
    .then((response) => {
      templateOption.value = response.data.list
    })
}

function serverChange(_serverIds: number[]) {
  if (_serverIds.includes(0)) {
    serverIds.value = _serverIds = serverOption.value.map((_) => _.id)
  }
  serverTags.value = {}
  for (const serverId of _serverIds) {
    serverTags.value[serverId] = {
      server: serverOption.value.find((_) => _.id === serverId) as ServerData,
      isRun: false,
      execRes: false,
      stdout: '',
      stderr: '',
    }
  }
}

function selectTemplate(templateId: number) {
  const template = templateOption.value.find(
    (_) => _.id === templateId
  ) as TemplateData
  script.value = template.content
}

function showExecRes(tag: ServerTag) {
  if (tag.isRun == false) {
    return
  }
  serverTab.value = tag.server.label
  execDialogVisible.value = true
}

function closeTag(serverId: number) {
  serverIds.value = serverIds.value.filter((_) => _ != serverId)
  serverTags.value = {}
  for (const serverId of serverIds.value) {
    serverTags.value[serverId] = {
      server: serverOption.value.find((_) => _.id === serverId) as ServerData,
      isRun: false,
      execRes: false,
      stdout: '',
      stderr: '',
    }
  }
}

function handleRun() {
  running.value = true
  new ServerExecScript({
    serverIds: serverIds.value,
    script: script.value,
  })
    .request()
    .then((response) => {
      let execRes = true
      let firstTag!: ServerTag
      for (const serverRes of response.data) {
        Object.assign(serverTags.value[serverRes.serverId], {
          isRun: true,
          ...serverRes,
        })
        if (execRes && serverRes.execRes == false) {
          execRes = false
        }
        if (!firstTag) {
          firstTag = serverTags.value[serverRes.serverId]
        }
      }
      showExecRes(firstTag)
      if (execRes) {
        ElMessage.success('Success')
      } else {
        ElMessage.error('Error')
      }
    })
    .finally(() => {
      running.value = false
    })
}

function handleSaveTemplate() {
  form.value?.validate((valid) => {
    if (valid) {
      saving.value = true
      new TemplateAdd({
        type: TemplateType.Script,
        name: formData.value.name,
        content: script.value,
        description: formData.value.description,
      })
        .request()
        .then(() => {
          ElMessage.success('Success')
          templateDialogVisible.value = false
          getTemplateOption()
        })
        .finally(() => {
          saving.value = false
        })
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
}

function templateRemove(data: TemplateData) {
  ElMessageBox.confirm(t('deleteTips', { name: data.name }), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      new TemplateRemove(data).request().then(() => {
        ElMessage.success('Success')
        getTemplateOption()
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}
</script>
