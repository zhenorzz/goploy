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
    <el-form ref="fileForm" :model="formData" class="file-form" @submit.prevent>
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
              ref="upload"
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

<script lang="ts">
import {
  ProjectFileList,
  ProjectFileContent,
  ProjectFileAdd,
  ProjectFileEdit,
  ProjectFileRemove,
} from '@/api/project'
import { VAceEditor } from 'vue3-ace-editor'
import { getRole } from '@/utils/namespace'
import { ElMessageBox, ElMessage } from 'element-plus'
import { computed, watch, defineComponent, reactive } from 'vue'
import { NamespaceKey, getNamespaceId } from '@/utils/namespace'
import { HttpResponse, ID } from '@/api/types'

interface FormFileInfo {
  id: number
  projectId: number
  filename: string
  content: string
  state: string
}

export default defineComponent({
  components: {
    VAceEditor,
  },
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    projectId: {
      type: Number,
      default: 0,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const dialogVisible = computed({
      get: () => props.modelValue,
      set: (val) => {
        emit('update:modelValue', val)
      },
    })
    const formData = reactive({
      files: [] as FormFileInfo[],
      content: '',
      projectId: 0,
    })
    const getProjectFileList = (projectId: number) => {
      new ProjectFileList({ id: projectId }).request().then((response) => {
        formData.files = response.data.list.map((item) => {
          return { ...item, state: 'success', content: '' }
        })
      })
    }
    watch(
      () => props.modelValue,
      (val: typeof props['modelValue']) => {
        if (val === true) {
          getProjectFileList(props.projectId)
        }
      }
    )

    return {
      role: getRole(),
      dialogVisible,
      getProjectFileList,
      formData,
    }
  },
  data() {
    return {
      formProps: {
        projectPath: '${PROJECT_PATH}',
        action: `${
          import.meta.env.VITE_APP_BASE_API
        }/project/uploadFile?${NamespaceKey}=${getNamespaceId()}`,
        show: 'file-list',
        editContentLoading: false,
        disabled: false,
        selectedIndex: -1,
      },
    }
  },
  watch: {
    projectId: function (newVal) {
      this.formData.projectId = newVal
    },
  },
  methods: {
    getProjectFileContent(file: FormFileInfo, index: number) {
      this.formProps.selectedIndex = index
      this.formProps.show = 'edit-content'
      if (file.id === 0) {
        this.formData.content = this.formData.files[index]['content']
        return
      }
      this.formData.content = ''
      this.formProps.editContentLoading = true
      new ProjectFileContent({ id: file.id })
        .request()
        .then((response) => {
          this.formData.content = this.formData.files[index]['content'] =
            response.data.content
        })
        .finally(() => {
          this.formProps.editContentLoading = false
        })
    },
    handleAppendFile() {
      this.formData.files.push({
        filename: '',
        projectId: this.projectId,
        state: 'loading',
        content: '',
        id: 0,
      })
    },
    beforeUpload(file: File, index: number) {
      this.formData.files[index].state = 'loading'
      ElMessage.info(this.$t('uploading'))
    },

    handleUploadSuccess(
      response: HttpResponse<ID>,
      file: File,
      fileList: File[],
      index: number
    ) {
      if (response.code !== 0) {
        this.formData.files[index].state = 'fail'
        ElMessage.error(response.message)
      } else {
        this.formData.files[index].id = response.data.id
        this.formData.files[index].state = 'success'
        ElMessage.success('Success')
      }
    },
    validateFilename(file: FormFileInfo, index: number) {
      const filename = file.filename
      if (file.state === 'success') {
        return true
      } else if (filename === '') {
        return false
      } else if (filename.substr(filename.length - 1, 1) === '/') {
        return false
      }
      const filenames = this.formData.files.map((item) => item.filename)
      filenames.splice(index, 1)
      return filenames.indexOf(filename) === -1
    },

    fileSubmit() {
      const file = this.formData.files[this.formProps.selectedIndex]
      this.formProps.disabled = true
      if (file.id === 0) {
        new ProjectFileAdd({
          projectId: file.projectId,
          filename: file.filename,
          content: this.formData.content,
        })
          .request()
          .then((response) => {
            this.formData.files[this.formProps.selectedIndex].id =
              response.data.id
            this.formData.files[this.formProps.selectedIndex].state = 'success'
            this.formProps.show = 'file-list'
            ElMessage.success('Success')
          })
          .finally(() => {
            this.formProps.disabled = false
          })
      } else {
        new ProjectFileEdit({
          id: file.id,
          content: this.formData.content,
        })
          .request()
          .then(() => {
            this.formData.files[this.formProps.selectedIndex].state = 'success'
            this.formProps.show = 'file-list'
            ElMessage.success('Success')
          })
          .finally(() => {
            this.formProps.disabled = false
          })
      }
    },

    removeFile(index: number) {
      if (this.formData.files[index].id === 0) {
        this.formData.files.splice(index, 1)
      } else {
        ElMessageBox.confirm(
          this.$t('projectPage.removeFileTips', {
            filename: this.formData.files[index].filename,
          }),
          this.$t('tips'),
          {
            confirmButtonText: this.$t('confirm'),
            cancelButtonText: this.$t('cancel'),
            type: 'warning',
          }
        )
          .then(() => {
            new ProjectFileRemove({
              projectFileId: this.formData.files[index].id,
            })
              .request()
              .then(() => {
                ElMessage.success('Success')
                this.formData.files.splice(index, 1)
              })
          })
          .catch(() => {
            ElMessage.info('Cancel')
          })
      }
    },
  },
})
</script>
