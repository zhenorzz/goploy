<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('manage')"
    :close-on-click-modal="false"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
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
                  (item) => item.superManager !== 1
                )"
                :key="index"
                :label="item.name"
                :value="item.id"
              />
            </el-select>
          </el-form-item>
          <el-form-item
            :label="$t('role')"
            label-width="60px"
            prop="role"
            style="margin-bottom: 5px"
          >
            <el-select v-model="formData.role">
              <el-option
                v-for="(_value, index) in [
                  role.Manager,
                  role.GroupManager,
                  role.Member,
                ]"
                :key="index"
                :label="_value"
                :value="_value"
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
      :data="tableData.filter((row) => row.role !== role.Admin)"
      style="width: 100%"
    >
      <el-table-column prop="userId" :label="$t('userId')" />
      <el-table-column prop="userName" :label="$t('userName')" />
      <el-table-column prop="role" :label="$t('role')" />
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
            type="danger"
            icon="el-icon-delete"
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
import {
  NamespaceUserData,
  NamespaceUserList,
  NamespaceUserAdd,
  NamespaceUserRemove,
} from '@/api/namespace'
import { UserOption } from '@/api/user'
import { getRole } from '@/utils/namespace'
import Validator from 'async-validator'
import { ElMessageBox, ElMessage } from 'element-plus'
import { computed, watch, ref, Ref } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
const role = getRole()
const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false,
  },
  namespaceId: {
    type: Number,
    default: 0,
  },
})
const emit = defineEmits(['update:modelValue'])
const form = ref<Validator>()
const formProps = ref({ disabled: false })
const formData = ref({ namespaceId: 0, userIds: [], role: '' })
const formRules = {
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
let tableData = ref<NamespaceUserList['datagram']['list']>([])
const dialogVisible = computed({
  get: () => props.modelValue,
  set: (val) => {
    emit('update:modelValue', val)
  },
})
let tableLoading = ref(false)
const getBindUserList = (namespaceId: number) => {
  tableLoading.value = true
  new NamespaceUserList({ id: namespaceId })
    .request()
    .then((response) => {
      tableData.value = response.data.list
    })
    .finally(() => {
      tableLoading.value = false
    })
}
watch(
  () => props.modelValue,
  (val: typeof props['modelValue']) => {
    if (val === true) {
      getBindUserList(props.namespaceId)
    }
  }
)
watch(
  () => props.namespaceId,
  (val: typeof props['namespaceId']) => {
    formData.value.namespaceId = val
  }
)

const userLoading = ref(false)
let userOption: Ref<UserOption['datagram']['list']> = ref([])
let showAddView = ref(false)
watch(showAddView, (val: boolean) => {
  if (val === true) {
    userLoading.value = true
    new UserOption()
      .request()
      .then((response) => {
        userOption.value = response.data.list
      })
      .finally(() => {
        userLoading.value = false
      })
  }
})
const handleAdd = () => {
  showAddView.value = true
}

function add() {
  form.value?.validate((valid: boolean) => {
    if (valid) {
      formProps.value.disabled = true
      new NamespaceUserAdd(formData.value)
        .request()
        .then(() => {
          showAddView.value = false
          ElMessage.success('Success')
          getBindUserList(formData.value.namespaceId)
        })
        .finally(() => {
          formProps.value.disabled = false
        })
    } else {
      return false
    }
  })
}

function remove(data: NamespaceUserData['datagram']) {
  ElMessageBox.confirm(t('namespacePage.removeUserTips'), t('tips'), {
    confirmButtonText: t('confirm'),
    cancelButtonText: t('cancel'),
    type: 'warning',
  })
    .then(() => {
      new NamespaceUserRemove({ namespaceUserId: data.id })
        .request()
        .then(() => {
          ElMessage.success('Success')
          getBindUserList(data.namespaceId)
        })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}
</script>
