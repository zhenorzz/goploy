<template>
  <el-dialog
    v-model="dialogVisible"
    :title="$t('deployPage.processManager')"
    :fullscreen="$store.state.app.device === 'mobile'"
    :close-on-click-modal="false"
  >
    <el-row type="flex">
      <el-select
        v-model="projectProcessId"
        v-loading="processLoading"
        filterable
        style="flex: 1"
        @change="handleProcessChange"
      >
        <el-option
          v-for="item in processOption"
          :key="item.id"
          :label="item.name"
          :value="item.id"
        >
          <el-row type="flex" justify="space-between">
            <span style="">{{ item.name }}</span>
            <el-button
              style=""
              type="text"
              icon="el-icon-delete"
              @click.stop="handleDelete(item.id)"
            />
          </el-row>
        </el-option>
      </el-select>
      <el-button type="primary" icon="el-icon-plus" @click="handleAdd" />
    </el-row>
    <el-table
      ref="table"
      v-loading="tableLoading"
      border
      stripe
      highlight-current-row
      max-height="447px"
      style="margin-top: 10px; width: 100%"
      :data="tableData"
    >
      <el-table-column type="expand">
        <template #default="{}">
          <el-row>
            {{ $t('deployPage.execRes') }}:{{ commandRes.execRes }}
          </el-row>
          <el-row style="white-space: pre-wrap">{{ commandRes.stdout }}</el-row>
          <el-row style="white-space: pre-wrap">{{ commandRes.stderr }}</el-row>
        </template>
      </el-table-column>
      <el-table-column label="Server">
        <template #default="scope"> {{ scope.row.serverName }} </template>
      </el-table-column>
      <el-table-column
        prop="operation"
        :label="$t('op')"
        width="400"
        align="center"
        :fixed="$store.state.app.device === 'mobile' ? false : 'right'"
      >
        <template #default="scope">
          <el-button
            :loading="commandLoading"
            type="primary"
            @click="handleProcessCmd(scope.row, 'status')"
          >
            status
          </el-button>
          <el-button
            :loading="commandLoading"
            type="success"
            @click="handleProcessCmd(scope.row, 'start')"
          >
            start
          </el-button>
          <el-button
            :loading="commandLoading"
            type="warning"
            @click="handleProcessCmd(scope.row, 'restart')"
          >
            restart
          </el-button>
          <el-button
            :loading="commandLoading"
            type="danger"
            @click="handleProcessCmd(scope.row, 'stop')"
          >
            stop
          </el-button>
        </template>
      </el-table-column>
    </el-table>
    <template #footer class="dialog-footer">
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
    </template>
  </el-dialog>
  <el-dialog
    v-model="processVisible"
    title="File diff"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-form
      ref="form"
      :model="formData"
      label-width="130px"
      :label-position="$store.state.app.device === 'desktop' ? 'right' : 'top'"
    >
      <el-form-item :label="$t('name')" prop="name" required>
        <el-input v-model="formData.name" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Status">
        <el-input v-model="formData.status" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Start">
        <el-input v-model="formData.start" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Stop">
        <el-input v-model="formData.stop" autocomplete="off" />
      </el-form-item>
      <el-form-item label="Restart">
        <el-input v-model="formData.restart" autocomplete="off" />
      </el-form-item>
    </el-form>
    <template #footer class="dialog-footer">
      <el-button @click="processVisible = false">
        {{ $t('cancel') }}
      </el-button>
      <el-button :disabled="formProps.disabled" type="primary" @click="add">
        {{ $t('confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts">
import { ManageProcess } from '@/api/deploy'
import {
  ProjectProcessList,
  ProjectProcessAdd,
  ProjectServerData,
  ProjectServerList,
  ProjectProcessDelete,
} from '@/api/project'
import Validator from 'async-validator'
import { ElMessageBox, ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { computed, defineComponent, ref, watch } from 'vue'

export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    projectRow: {
      type: Object,
      required: true,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const { t } = useI18n()
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
          getList()
        }
      }
    )
    const processLoading = ref(false)
    const projectProcessId = ref<number>()
    const processOption = ref<ProjectProcessList['datagram']['list']>([])
    const getList = () => {
      processLoading.value = true
      projectProcessId.value = undefined
      processOption.value = []
      tableData.value = []
      new ProjectProcessList(
        { projectId: props.projectRow.id },
        { page: 1, rows: 999 }
      )
        .request()
        .then((response) => {
          processOption.value = response.data.list
        })
        .finally(() => {
          processLoading.value = false
        })
    }

    const table = ref(null)
    const tableLoading = ref(false)
    const tableData = ref<ProjectServerList['datagram']['list']>([])
    const handleProcessChange = () => {
      tableLoading.value = true
      new ProjectServerList({ id: props.projectRow.id })
        .request()
        .then((response) => {
          tableData.value = response.data.list
        })
        .finally(() => {
          tableLoading.value = false
        })
    }
    const commandRes = ref<ManageProcess['datagram']>({
      execRes: true,
      stdout: '',
      stderr: '',
    })
    const commandLoading = ref(false)
    const handleProcessCmd = (
      data: ProjectServerData['datagram'],
      command: string
    ) => {
      ElMessageBox.confirm(t('deployPage.execTips', { command }), t('tips'), {
        confirmButtonText: t('confirm'),
        cancelButtonText: t('cancel'),
        type: 'warning',
      })
        .then(() => {
          commandLoading.value = true
          new ManageProcess({
            serverId: data.serverId,
            projectProcessId: Number(projectProcessId.value),
            command,
          })
            .request()
            .then((response) => {
              commandRes.value = response.data
              table.value.toggleRowExpansion(data, true)
            })
            .finally(() => {
              commandLoading.value = false
            })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    }
    return {
      dialogVisible,
      getList,
      projectProcessId,
      processOption,
      processLoading,
      handleProcessChange,
      table,
      tableLoading,
      tableData,
      handleProcessCmd,
      commandLoading,
      commandRes,
    }
  },
  data() {
    return {
      processVisible: false,
      formProps: {
        disabled: false,
      },
      formData: {
        projectId: 0,
        name: '',
        status: '',
        start: '',
        stop: '',
        restart: '',
      },
    }
  },
  watch: {
    projectRow: function (newVal) {
      this.formData.projectId = newVal.id
    },
  },
  methods: {
    handleAdd() {
      this.processVisible = true
    },
    handleDelete(id: number) {
      ElMessageBox.confirm(
        this.$t('deployPage.deleteProcessTips'),
        this.$t('tips'),
        {
          confirmButtonText: this.$t('confirm'),
          cancelButtonText: this.$t('cancel'),
          type: 'warning',
        }
      )
        .then(() => {
          new ProjectProcessDelete({ id }).request().then(() => {
            ElMessage.success('Success')
            this.getList()
          })
        })
        .catch(() => {
          ElMessage.info('Cancel')
        })
    },
    add() {
      ;(this.$refs.form as Validator).validate((valid: boolean) => {
        if (valid) {
          this.formProps.disabled = true
          new ProjectProcessAdd(this.formData)
            .request()
            .then(() => {
              this.processVisible = false
              ElMessage.success('Success')
              this.getList()
            })
            .finally(() => {
              this.formProps.disabled = false
            })
        } else {
          return false
        }
      })
    },
  },
})
</script>

<style rel="stylesheet/scss" lang="scss" scoped>
@import '@/styles/mixin.scss';
.file {
  flex: 1;
  white-space: pre-wrap;
}
</style>
