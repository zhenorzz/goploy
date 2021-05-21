<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('manage')"
    :close-on-click-modal="false"
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
            :label="$t('server')"
            label-width="120px"
            prop="serverIds"
          >
            <el-select v-model="formData.serverIds" multiple>
              <el-option
                v-for="(item, index) in serverOption"
                :key="index"
                :label="item.label"
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
      style="width: 100%"
    >
      <el-table-column prop="serverId" :label="$t('serverId')" width="100" />
      <el-table-column
        prop="serverName"
        :label="$t('serverName')"
        width="120"
      />
      <el-table-column
        prop="serverDescription"
        :label="$t('serverDescription')"
        min-width="200"
        show-overflow-tooltip
      />
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

<script lang="ts">
import {
  CrontabServerData,
  CrontabServerList,
  CrontabServerAdd,
  CrontabServerRemove,
} from '@/api/crontab'
import { ServerOption } from '@/api/server'
import Validator from 'async-validator'
import { ElMessageBox, ElMessage } from 'element-plus'
import { computed, watch, defineComponent, ref, Ref } from 'vue'

export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    crontabId: {
      type: Number,
      default: 0,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const tableData: Ref<CrontabServerList['datagram']['list']> = ref([])
    const dialogVisible = computed({
      get: () => props.modelValue,
      set: (val) => {
        emit('update:modelValue', val)
      },
    })
    const tableLoading = ref(false)
    const getBindServerList = (crontabId: number) => {
      tableLoading.value = true
      new CrontabServerList({ id: crontabId })
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
          getBindServerList(props.crontabId)
        }
      }
    )

    const showAddView = ref(false)
    const handleAdd = () => {
      showAddView.value = true
    }
    const serverOption: Ref<ServerOption['datagram']['list']> = ref([])
    watch(showAddView, (val: boolean) => {
      if (val === true) {
        new ServerOption().request().then((response) => {
          serverOption.value = response.data.list
        })
      }
    })

    return {
      dialogVisible,
      getBindServerList,
      tableLoading,
      tableData,
      showAddView,
      handleAdd,
      serverOption,
    }
  },
  data() {
    return {
      formProps: {
        disabled: false,
      },
      formData: {
        crontabId: 0,
        serverIds: [],
      },
      formRules: {
        serverIds: [
          {
            type: 'array',
            required: true,
            message: 'Server required',
            trigger: 'change',
          },
        ],
      },
    }
  },
  watch: {
    crontabId: function (newVal) {
      this.formData.crontabId = newVal
    },
  },
  methods: {
    add() {
      ;(this.$refs.form as Validator).validate((valid: boolean) => {
        if (valid) {
          this.formProps.disabled = true
          new CrontabServerAdd(this.formData)
            .request()
            .then(() => {
              ElMessage.success('Success')
              this.getBindServerList(this.formData.crontabId)
            })
            .finally(() => {
              this.formProps.disabled = false
            })
        } else {
          return false
        }
      })
    },

    remove(data: CrontabServerData['datagram']['detail']) {
      ElMessageBox.confirm(
        this.$t('crontabPage.removeCrontabServerTips'),
        this.$t('tips'),
        {
          confirmButtonText: this.$t('confirm'),
          cancelButtonText: this.$t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          new CrontabServerRemove({
            crontabServerId: data.id,
            crontabId: data.crontabId,
            serverId: data.serverId,
          })
            .request()
            .then(() => {
              ElMessage.success('Success')
              this.getBindServerList(data.crontabId)
            })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },
  },
})
</script>
