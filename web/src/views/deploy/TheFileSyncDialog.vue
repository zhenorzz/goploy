<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('manage')"
    custom-class="file-dialog"
    :close-on-click-modal="false"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-row
      v-if="formProps.show === 'file-list'"
      type="flex"
      justify="space-between"
      align="middle"
      style="margin: 10px 0"
    >
      <span>{{ $t('projectPage.projectFileTips') }}</span>
      <el-button type="primary" icon="el-icon-plus" @click="handleAppendFile" />
    </el-row>
    <el-form :model="formData" class="file-form" @submit.prevent>
      <template v-if="formProps.show === 'file-list'">
        <el-form-item
          v-for="(file, index) in formData.files"
          :key="index"
          label-width="0"
          prop="directory"
        >
          <el-row type="flex" align="middle">
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
                <i
                  class="el-input__icon el-icon-check"
                  style="color: #67c23a; font-size: 16px; font-weight: 900"
                />
              </template>
              <template v-else-if="file.state === 'loading'" #suffix>
                <i
                  class="el-input__icon el-icon-loading"
                  style="font-size: 14px; font-weight: 900"
                />
              </template>
              <template v-else #suffix>
                <i
                  class="el-input__icon el-icon-close"
                  style="color: #f56c6c; font-size: 16px; font-weight: 900"
                />
              </template>
            </el-input>
            <el-upload
              style="margin: 0 12px"
              :action="`${formProps.action}&projectFileId=${file.id}&projectId=${file.projectId}&filename=${file.filename}`"
              :before-upload="(uploadFile) => beforeUpload(uploadFile, index)"
              :on-success="
                (response, uploadFile, uploadFileList) =>
                  handleUploadSuccess(
                    response,
                    uploadFile,
                    uploadFileList,
                    index
                  )
              "
              :show-file-list="false"
              :disabled="!validateFilename(file, index)"
              multiple
            >
              <el-button
                :disabled="!validateFilename(file, index)"
                type="text"
                icon="el-icon-upload"
              />
            </el-upload>
            <el-button
              :disabled="!validateFilename(file, index)"
              type="text"
              icon="el-icon-edit"
              @click="getProjectFileContent(file, index)"
            />
            <el-button
              type="text"
              icon="el-icon-delete"
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
          theme="github"
          style="height: 500px"
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
  ProjectFileList,
  ProjectFileContent,
  ProjectFileAdd,
  ProjectFileEdit,
  ProjectFileRemove,
} from '@/api/project'
import { VAceEditor } from 'vue3-ace-editor'
import { ElMessageBox, ElMessage } from 'element-plus'
import { ref, computed, watch } from 'vue'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import { HttpResponse, ID } from '@/api/types'
import { useI18n } from 'vue-i18n'
interface FormFileInfo {
  id: number
  projectId: number
  filename: string
  content: string
  state: string
}
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

function handleUploadSuccess(
  response: HttpResponse<ID>,
  file: File,
  fileList: File[],
  index: number
) {
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
      t('projectPage.removeFileTips', {
        filename: formData.value.files[index].filename,
      }),
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
  .el-icon-upload {
    font-size: 14px;
  }
  .file-form {
    height: 520px;
    overflow-y: auto;
    @include scrollBar();
  }
}
</style>
