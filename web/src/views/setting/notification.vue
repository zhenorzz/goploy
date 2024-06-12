<template>
  <el-row class="app-container">
    <el-row class="app-bar" type="flex" justify="space-between"></el-row>
    <el-row class="app-table">
      <el-table
        :key="tableHeight"
        v-loading="tableLoading"
        highlight-current-row
        height="100%"
        :data="tableData"
        style="width: 100%"
      >
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="useBy" :label="'Use By'" width="100" />
        <el-table-column prop="title" :label="$t('title')" width="200" />
        <el-table-column prop="type" :label="$t('type')" width="100">
          <template #default="scope">
            {{ $t(`webhookOption[${scope.row.type || 0}]`) }}
          </template>
        </el-table-column>
        <el-table-column prop="template" :label="$t('template')" />
        <el-table-column
          prop="insertTime"
          :label="$t('insertTime')"
          width="160"
          align="center"
        />
        <el-table-column
          prop="updateTime"
          :label="$t('updateTime')"
          width="160"
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
            <Button
              type="primary"
              :icon="Edit"
              :permissions="[permission.EditNotification]"
              @click="handleEdit(scope.row)"
            />
          </template>
        </el-table-column>
      </el-table>
    </el-row>
    <el-dialog
      v-model="dialogVisible"
      :title="$t('setting')"
      :fullscreen="$store.state.app.device === 'mobile'"
    >
      <el-form
        ref="form"
        :model="formData"
        label-width="80px"
        :label-position="
          $store.state.app.device === 'desktop' ? 'right' : 'top'
        "
      >
        <el-form-item
          :label="$t('title')"
          prop="title"
          :rules="[
            { required: true, message: 'Title required', trigger: 'blur' },
          ]"
        >
          <el-input v-model="formData.title" autocomplete="off" />
        </el-form-item>
        <el-form-item
          :label="$t('template')"
          prop="template"
          :rules="[
            { required: true, message: 'Template required', trigger: 'blur' },
          ]"
        >
          <el-input
            v-model="formData.template"
            type="textarea"
            :autosize="{ minRows: 8 }"
            autocomplete="off"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('cancel') }}</el-button>
        <el-button
          :disabled="formProps.disabled"
          type="primary"
          @click="submit"
        >
          {{ $t('confirm') }}
        </el-button>
      </template>
    </el-dialog>
  </el-row>
</template>
<script lang="ts">
export default { name: 'NotificationSetting' }
</script>
<script lang="ts" setup>
import permission from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Edit } from '@element-plus/icons-vue'
import {
  NotificationList,
  NotificationEdit,
  NotificationData,
} from '@/api/notification'
import getTableHeight from '@/composables/tableHeight'
import type { ElForm } from 'element-plus'
import { ref } from 'vue'

const { tableHeight } = getTableHeight()
const dialogVisible = ref(false)
const tableLoading = ref(false)
const tableData = ref<NotificationList['datagram']['list']>([])
const selectedItem = ref<NotificationData>()
const form = ref<InstanceType<typeof ElForm>>()
const tempFormData = { id: 0, title: '', template: '' }
const formData = ref(tempFormData)
const formProps = ref({ disabled: false })

getList()

function getList() {
  tableLoading.value = true
  new NotificationList()
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function handleEdit(data: NotificationData) {
  formData.value = Object.assign({}, data)
  dialogVisible.value = true
}

function submit() {
  form.value?.validate((valid) => {
    if (valid) {
      edit()
      return Promise.resolve(true)
    } else {
      return Promise.reject(false)
    }
  })
}

function edit() {
  formProps.value.disabled = true
  new NotificationEdit(formData.value)
    .request()
    .then(() => {
      getList()
      ElMessage.success('Success')
    })
    .finally(() => {
      formProps.value.disabled = dialogVisible.value = false
    })
}
</script>
