<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('manage')"
    :close-on-click-modal="false"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" :icon="Plus" @click="handleAdd" />
      <el-row
        v-if="showAddView"
        type="flex"
        justify="center"
        style="margin-top: 10px; width: 100%"
      >
        <el-form ref="form" :inline="true" :rules="formRules" :model="formData">
          <el-form-item
            :label="$t('user')"
            label-width="60px"
            prop="userIds"
            style="margin-bottom: 5px"
          >
            <el-select
              v-model="formData.userIds"
              :loading="userLoading"
              multiple
              clearable
              filterable
            >
              <el-option
                v-for="(item, index) in userOption.filter(
                  (item) => [role.Admin, role.Manager].indexOf(item.role) === -1
                )"
                :key="index"
                :label="item.userName"
                :value="item.userId"
              />
            </el-select>
          </el-form-item>
          <el-form-item style="margin-right: 0px; margin-bottom: 5px">
            <el-button
              type="primary"
              :disabled="formProps.disabled"
              @click="add"
            >
              {{ $t('confirm') }}
            </el-button>
            <el-button @click="showAddView = false">
              {{ $t('cancel') }}
            </el-button>
          </el-form-item>
        </el-form>
      </el-row>
    </el-row>
    <el-table
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      :data="tableData"
      style="width: 100%"
    >
      <el-table-column prop="userId" :label="$t('userId')" width="100" />
      <el-table-column prop="userName" :label="$t('userName')" />
      <el-table-column
        prop="insertTime"
        :label="$t('insertTime')"
        width="135"
        align="center"
      />
      <el-table-column
        prop="updateTime"
        :label="$t('updateTime')"
        width="135"
        align="center"
      />
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="80"
        align="center"
        fixed="right"
      >
        <template #default="scope">
          <el-button
            v-show="role.hasManagerPermission()"
            type="danger"
            :icon="Delete"
            @click="remove(scope.row)"
          />
        </template>
      </el-table-column>
    </el-table>
    <template #footer>
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts" setup>
import { Plus, Delete } from '@element-plus/icons-vue'
import {
  ProjectUserData,
  ProjectUserList,
  ProjectUserAdd,
  ProjectUserRemove,
} from '@/api/project'
import { NamespaceUserOption } from '@/api/namespace'
import { getRole } from '@/utils/namespace'
import type { ElForm } from 'element-plus'
import { ElMessageBox, ElMessage } from 'element-plus'
import { computed, watch, ref } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const role = getRole()
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

watch(
  () => props.modelValue,
  (val: typeof props['modelValue']) => {
    if (val === true) {
      getBindUserList(props.projectId)
    }
  }
)

watch(
  () => props.projectId,
  (val) => {
    formData.value.projectId = val
  }
)

const showAddView = ref(false)
const userLoading = ref(false)
const userOption = ref<NamespaceUserOption['datagram']['list']>([])
watch(showAddView, (val: boolean) => {
  if (val === true) {
    userLoading.value = true
    new NamespaceUserOption()
      .request()
      .then((response) => {
        userOption.value = response.data.list
      })
      .finally(() => {
        userLoading.value = false
      })
  }
})

const tableLoading = ref(false)
const tableData = ref<ProjectUserList['datagram']['list']>([])

const form = ref<InstanceType<typeof ElForm>>()
const formProps = ref({ disabled: false })
const formData = ref({ projectId: 0, userIds: [] })
const formRules = <InstanceType<typeof ElForm>['rules']>{
  userIds: [
    {
      type: 'array',
      required: true,
      message: 'User required',
      trigger: 'change',
    },
  ],
  role: [{ required: true, message: 'Role required', trigger: 'change' }],
}

function getBindUserList(projectId: number) {
  tableLoading.value = true
  new ProjectUserList({ id: projectId })
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}

function handleAdd() {
  showAddView.value = true
}

function add() {
  form.value?.validate((valid) => {
    if (valid) {
      formProps.value.disabled = true
      new ProjectUserAdd(formData.value)
        .request()
        .then(() => {
          showAddView.value = false
          ElMessage.success('Success')
          getBindUserList(formData.value.projectId)
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

function remove(data: ProjectUserData) {
  ElMessageBox.confirm(t('namespacePage.removeUserTips'), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      new ProjectUserRemove({ projectUserId: data.id }).request().then(() => {
        ElMessage.success('Success')
        getBindUserList(data.projectId)
      })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}
</script>
