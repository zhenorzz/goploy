<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('manage')"
    class="file-dialog"
    :close-on-click-modal="false"
    :fullscreen="$store.state.app.device === 'mobile'"
    @close="formProps.show = 'file-list'"
  >
    <el-row
      v-if="formProps.show === 'file-list'"
      type="flex"
      justify="space-between"
      align="middle"
      style="margin: 10px 0"
    >
      <span>{{ $t('projectPage.projectFileTips') }}</span>
      <el-button type="primary" :icon="Plus" @click="handleAppendFile" />
    </el-row>
    <el-form :model="formData" class="file-form" @submit.prevent>
      <template v-if="formProps.show === 'file-list'">
        <el-form-item
          v-for="(file, index) in formData.files"
          :key="index"
          label-width="0"
          prop="directory"
        >
          <el-row type="flex" align="middle" style="width: 100%">
            <el-input
              v-model.trim="file.filename"
              autocomplete="off"
              placeholder="path/to/example/.env"
              :disabled="file.state === 'success'"
              :readonly="file.state === 'success'"
              style="flex: 1"
            >
              <template #prepend>{{ formProps.projectPath }}</template>
              <template v-if="file.state === 'success'" #suffix>
                <el-icon
                  class="el-input__icon"
                  color="#67c23a"
                  style="font-size: 16px"
                >
                  <Check />
                </el-icon>
              </template>
              <template v-else-if="file.state === 'loading'" #suffix>
                <el-icon class="el-input__icon is-loading">
                  <Loading />
                </el-icon>
              </template>
              <template v-else #suffix>
                <el-icon
                  class="el-input__icon"
                  color="#f56c6c"
                  style="font-size: 16px"
                >
                  <Close />
                </el-icon>
              </template>
            </el-input>
            <el-upload
              style="margin: 0 12px"
              :action="`${formProps.action}&projectFileId=${file.id}&projectId=${file.projectId}&filename=${file.filename}`"
              :before-upload="(uploadFile) => beforeUpload(uploadFile, index)"
              :on-success="(response) => handleUploadSuccess(response, index)"
              :show-file-list="false"
              :disabled="!validateFilename(file, index)"
              multiple
            >
              <el-button
                :disabled="!validateFilename(file, index)"
                type="primary"
                link
                :icon="Upload"
              />
            </el-upload>
            <el-button
              :disabled="!validateFilename(file, index)"
              type="primary"
              link
              :icon="Edit"
              @click="getProjectFileContent(file, index)"
            />
            <el-button
              type="primary"
              link
              :icon="Delete"
              @click="removeFile(index)"
            />
          </el-row>
        </el-form-item>
      </template>
      <el-form-item
        v-else
        v-loading="formProps.editContentLoading"
        prop="projectFileEdit"
        label-width="0px"
      >
        <v-ace-editor
          v-model:value="formData.content"
          lang="sh"
          :theme="isDark ? 'one_dark' : 'github'"
          style="height: 500px; width: 100%"
        />
      </el-form-item>
    </el-form>
    <template v-if="formProps.show !== 'file-list'" #footer>
      <el-button @click="formProps.show = 'file-list'">
        {{ $t('cancel') }}
      </el-button>
      <el-button
        :disabled="formProps.disabled"
        type="primary"
        @click="fileSubmit"
      >
        {{ $t('confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import {
  Upload,
  Edit,
  Delete,
  Check,
  Plus,
  Close,
  Loading,
} from '@element-plus/icons-vue'
import {
  ProjectFileList,
  ProjectFileContent,
  ProjectFileAdd,
  ProjectFileEdit,
  ProjectFileRemove,
} from '@/api/project'
import { VAceEditor } from 'vue3-ace-editor'
import * as ace from 'ace-builds/src-noconflict/ace'
import { ref, computed, watch } from 'vue'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import { HttpResponse, ID } from '@/api/types'
import { useDark } from '@vueuse/core'
import { useI18n } from 'vue-i18n'
ace.config.set(
  'basePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
ace.config.set(
  'themePath',
  'https://cdn.jsdelivr.net/npm/ace-builds@' + ace.version + '/src-noconflict/'
)
interface FormFileInfo {
  id: number
  projectId: number
  filename: string
  content: string
  state: string
}
const isDark = useDark()
const { t } = useI18n()
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  projectId: {
    type: Number,
    default: 0,
  },
})
const emit = defineEmits(['update:modelValue'])
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val)
  },
})

const formData = ref({
  files: [] as FormFileInfo[],
  content: '',
})
const formProps = ref({
  projectId: 0,
  projectPath: '${PROJECT_PATH}',
  action: `${
    import.meta.env.VITE_APP_BASE_API
  }/project/uploadFile?${NamespaceKey}=${getNamespaceId()}`,
  show: 'file-list',
  editContentLoading: false,
  disabled: false,
  selectedIndex: -1,
})

watch(
  () => props.modelValue,
  (val) => {
    if (val === true) {
      getProjectFileList(props.projectId)
    }
  }
)

watch(
  () => props.projectId,
  (val) => {
    formProps.value.projectId = val
  }
)

function getProjectFileList(projectId: number) {
  formData.value.files = []
  new ProjectFileList({ id: projectId }).request().then((response) => {
    formData.value.files = response.data.list.map((item) => {
      return { ...item, state: 'success', content: '' }
    })
  })
}

function getProjectFileContent(file: FormFileInfo, index: number) {
  formProps.value.selectedIndex = index
  formProps.value.show = 'edit-content'
  if (file.id === 0) {
    formData.value.content = formData.value.files[index]['content']
    return
  }
  formData.value.content = ''
  formProps.value.editContentLoading = true
  new ProjectFileContent({ id: file.id })
    .request()
    .then((response) => {
      formData.value.content = formData.value.files[index]['content'] =
        response.data.content
    })
    .finally(() => {
      formProps.value.editContentLoading = false
    })
}

function handleAppendFile() {
  formData.value.files.push({
    filename: '',
    projectId: formProps.value.projectId,
    state: 'loading',
    content: '',
    id: 0,
  })
}

function beforeUpload(file: File, index: number) {
  formData.value.files[index].state = 'loading'
  ElMessage.info(t('uploading'))
}

function handleUploadSuccess(response: HttpResponse<ID>, index: number) {
  if (response.code !== 0) {
    formData.value.files[index].state = 'fail'
    ElMessage.error(response.message)
  } else {
    formData.value.files[index].id = response.data.id
    formData.value.files[index].state = 'success'
    ElMessage.success('Success')
  }
}
function validateFilename(file: FormFileInfo, index: number) {
  const filename = file.filename
  if (file.state === 'success') {
    return true
  } else if (filename === '') {
    return false
  } else if (filename.substr(filename.length - 1, 1) === '/') {
    return false
  }
  const filenames = formData.value.files.map((item) => item.filename)
  filenames.splice(index, 1)
  return filenames.indexOf(filename) === -1
}

function fileSubmit() {
  const file = formData.value.files[formProps.value.selectedIndex]
  formProps.value.disabled = true
  if (file.id === 0) {
    new ProjectFileAdd({
      projectId: file.projectId,
      filename: file.filename,
      content: formData.value.content,
    })
      .request()
      .then((response) => {
        formData.value.files[formProps.value.selectedIndex].id =
          response.data.id
        formData.value.files[formProps.value.selectedIndex].state = 'success'
        formProps.value.show = 'file-list'
        ElMessage.success('Success')
      })
      .finally(() => {
        formProps.value.disabled = false
      })
  } else {
    new ProjectFileEdit({
      id: file.id,
      content: formData.value.content,
    })
      .request()
      .then(() => {
        formData.value.files[formProps.value.selectedIndex].state = 'success'
        formProps.value.show = 'file-list'
        ElMessage.success('Success')
      })
      .finally(() => {
        formProps.value.disabled = false
      })
  }
}

function removeFile(index: number) {
  if (formData.value.files[index].id === 0) {
    formData.value.files.splice(index, 1)
  } else {
    ElMessageBox.confirm(
      t('deleteTips', { name: formData.value.files[index].filename }),
      t('tips'),
      {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning',
      }
    )
      .then(() => {
        new ProjectFileRemove({
          projectFileId: formData.value.files[index].id,
        })
          .request()
          .then(() => {
            ElMessage.success('Success')
            formData.value.files.splice(index, 1)
          })
      })
      .catch(() => {
        ElMessage.info('Cancel')
      })
  }
}
</script>
<style lang="scss">
@import '@/styles/mixin.scss';
.file-dialog {
  .el-dialog__body {
    padding-top: 20px;
  }
  .file-form {
    height: 520px;
    overflow-y: auto;
    @include scrollBar();
  }
}
</style>
