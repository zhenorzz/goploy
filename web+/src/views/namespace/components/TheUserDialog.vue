<template>
  <el-dialog v-model="dialogVisible" :title="$t('manage')">
    <el-row class="app-bar" type="flex" justify="end">
      <el-button type="primary" icon="el-icon-plus" @click="handleAddUser" />
      <el-row
        v-if="showUserAddView"
        type="flex"
        justify="center"
        style="margin-top: 10px; width: 100%"
      >
        <el-form
          ref="addUserForm"
          :inline="true"
          :rules="formRules"
          :model="formData"
        >
          <el-form-item
            :label="$t('user')"
            label-width="60px"
            prop="userIds"
            style="margin-bottom: 5px"
          >
            <el-select v-model="formData.userIds" multiple clearable filterable>
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
              @click="addUser"
            >
              {{ $t('confirm') }}
            </el-button>
            <el-button @click="showUserAddView = false">
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
      >
        <template #default="scope">
          <el-button
            type="danger"
            icon="el-icon-delete"
            @click="removeUser(scope.row)"
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

<script lang="ts">
import {
  NamespaceUserData,
  NamespaceUserList,
  NamespaceUserAdd,
  NamespaceUserRemove,
} from '@/api/namespace'
import { UserOption } from '@/api/user'
import { role } from '@/utils/namespace'
import Validator from 'async-validator'
import { ElMessageBox, ElMessage } from 'element-plus'
import { computed, watch, defineComponent, ref, Ref } from 'vue'

export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    namespaceId: {
      type: Number,
      default: 0,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    let tableData: Ref<NamespaceUserList['datagram']['list']> = ref([])
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

    let userOption: Ref<UserOption['datagram']['list']> = ref([])
    new UserOption().request().then((response) => {
      userOption.value = response.data.list
    })

    let showUserAddView = ref(false)
    const handleAddUser = () => {
      showUserAddView.value = true
    }

    return {
      role,
      dialogVisible,
      getBindUserList,
      tableLoading,
      tableData,
      showUserAddView,
      handleAddUser,
      userOption,
    }
  },
  data() {
    return {
      formProps: {
        disabled: false,
      },
      formData: {
        namespaceId: 0,
        userIds: [],
        role: '',
      },
      formRules: {
        userIds: [
          {
            type: 'array',
            required: true,
            message: 'User required',
            trigger: 'change',
          },
        ],
        role: [{ required: true, message: 'Role required', trigger: 'change' }],
      },
    }
  },
  watch: {
    namespaceId: function (newVal) {
      this.formData.namespaceId = newVal
    },
  },
  methods: {
    addUser() {
      ;(this.$refs.addUserForm as Validator).validate((valid: boolean) => {
        if (valid) {
          this.formProps.disabled = true
          new NamespaceUserAdd(this.formData)
            .request()
            .then(() => {
              this.showUserAddView = false
              ElMessage.success('Success')
              this.getBindUserList(this.formData.namespaceId)
            })
            .finally(() => {
              this.formProps.disabled = false
            })
        } else {
          return false
        }
      })
    },

    removeUser(data: NamespaceUserData['datagram']['detail']) {
      ElMessageBox.confirm(
        this.$t('namespacePage.removeUserTips'),
        this.$t('tips'),
        {
          confirmButtonText: this.$t('confirm'),
          cancelButtonText: this.$t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          new NamespaceUserRemove({ namespaceUserId: data.id })
            .request()
            .then(() => {
              ElMessage.success('Success')
              this.getBindUserList(data.namespaceId)
            })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },
  },
})
</script>
