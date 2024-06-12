<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('manage')"
    :close-on-click-modal="false"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-row class="app-bar" type="flex" justify="end">
      <Button
        type="primary"
        :icon="Plus"
        :permissions="[permission.AddNamespaceUser]"
        @click="handleAdd"
      />
      <el-row
        v-if="showAddView"
        type="flex"
        justify="center"
        style="margin-top: 10px; width: 100%"
      >
        <el-form ref="form" :inline="true" :model="formData">
          <el-form-item
            :label="$t('user')"
            label-width="60px"
            prop="userIds"
            style="margin-bottom: 5px"
            :rules="[
              {
                type: 'array',
                required: true,
                message: 'User required',
                trigger: 'change',
              },
            ]"
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
            prop="roleId"
            style="margin-bottom: 5px"
            :rules="[
              { required: true, message: 'Role required', trigger: 'change' },
            ]"
          >
            <el-select
              v-model="formData.roleId"
              :loading="roleLoading"
              filterable
            >
              <el-option
                v-for="(item, index) in roleOption"
                :key="index"
                :label="item.name"
                :value="item.id"
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
    >
      <el-table-column prop="userName" :label="$t('userName')" />
      <el-table-column prop="roleName" :label="$t('role')" />
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
        width="80"
        align="center"
        fixed="right"
      >
        <template #default="scope">
          <Button
            type="danger"
            :icon="Delete"
            :permissions="[permission.DeleteNamespaceUser]"
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
import permission from '@/permission'
import Button from '@/components/Permission/Button.vue'
import { Plus, Delete } from '@element-plus/icons-vue'
import {
  NamespaceUserData,
  NamespaceUserList,
  NamespaceUserAdd,
  NamespaceUserRemove,
} from '@/api/namespace'
import { RoleOption } from '@/api/role'
import { UserOption } from '@/api/user'
import type { ElForm } from 'element-plus'
import { computed, watch, ref } from 'vue'
import { useI18n } from 'vue-i18n'
const { t } = useI18n()
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
const form = ref<InstanceType<typeof ElForm>>()
const formProps = ref({ disabled: false })
const formData = ref({ namespaceId: 0, userIds: [], roleId: '' })

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
      tableData.value = response.data.list.filter((_) => _.roleId > 0)
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
const userOption = ref<UserOption['datagram']['list']>([])
const roleLoading = ref(false)
const roleOption = ref<RoleOption['datagram']['list']>([])
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
    roleLoading.value = true
    new RoleOption()
      .request()
      .then((response) => {
        roleOption.value = response.data.list
      })
      .finally(() => {
        roleLoading.value = false
      })
  }
})
const handleAdd = () => {
  showAddView.value = true
}

function add() {
  form.value?.validate((valid) => {
    if (valid) {
      formProps.value.disabled = true
      new NamespaceUserAdd({
        ...formData.value,
        roleId: Number(formData.value.roleId),
      })
        .request()
        .then(() => {
          showAddView.value = false
          ElMessage.success('Success')
          getBindUserList(formData.value.namespaceId)
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

function remove(data: NamespaceUserData) {
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
          getBindUserList(props.namespaceId)
        })
    })
    .catch(() => {
      ElMessage.info('Cancel')
    })
}
</script>
