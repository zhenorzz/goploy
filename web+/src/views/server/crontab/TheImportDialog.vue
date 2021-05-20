<template>
  <el-dialog v-model="dialogVisible" :title="$t('import')">
    <el-row>
      <el-row type="flex" style="width: 100%">
        <el-select
          v-model="importProps.serverId"
          style="flex: 1; margin-right: 5px"
        >
          <el-option
            v-for="(item, index) in serverOption"
            :key="index"
            :label="
              item.name + (item.description ? `(${item.description})` : '')
            "
            :value="item.id"
          />
        </el-select>
        <el-button
          :disabled="importProps.disabled"
          :icon="importProps.disabled ? 'el-icon-loading' : 'el-icon-refresh'"
          type="primary"
          @click="getRemoteServerList"
        >
          {{ $t('read') }}
        </el-button>
      </el-row>
      <el-table
        border
        stripe
        highlight-current-row
        :empty-text="$t('crontabPage.importTips')"
        :data="tableData"
        style="width: 100%; margin-top: 10px"
        @selection-change="onSelectionChange"
      >
        <el-table-column type="selection" width="40" />
        <el-table-column
          prop="command"
          :label="$t('command')"
          min-width="140"
          show-overflow-tooltip
        />
        <el-table-column
          prop="description"
          :label="$t('desc')"
          min-width="240"
          show-overflow-tooltip
        />
      </el-table>
    </el-row>
    <template #footer>
      <el-button @click="dialogVisible = false">
        {{ $t('cancel') }}
      </el-button>
      <el-button
        :disabled="importProps.disabled"
        type="primary"
        @click="importCrontab"
      >
        {{ $t('confirm') }}
      </el-button>
    </template>
  </el-dialog>
</template>

<script lang="ts">
import { CrontabImport, CrontabsInRemoteServer } from '@/api/crontab'
import { ServerOption } from '@/api/server'
import cronstrue from 'cronstrue/i18n'
import { ElMessage } from 'element-plus'
import { useI18n } from 'vue-i18n'
import { computed, watch, defineComponent, ref, reactive } from 'vue'

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
    onSuccess: {
      type: Function,
      required: true,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const { t, locale } = useI18n()

    const dialogVisible = computed({
      get: () => props.modelValue,
      set: (val) => {
        emit('update:modelValue', val)
      },
    })
    let serverOption = ref<ServerOption['datagram']['list']>([])
    watch(
      () => props.modelValue,
      (val: typeof props['modelValue']) => {
        if (val === true) {
          new ServerOption().request().then((response) => {
            serverOption.value = response.data.list
          })
        }
      }
    )

    const importProps = reactive({
      serverId: '',
      disabled: false,
      loading: false,
    })
    let tableData = ref<Record<string, string>[]>([])
    const getLocale = () => {
      if (locale.value === 'zh-cn') {
        return 'zh_CN'
      }
      return locale.value
    }
    const getRemoteServerList = () => {
      if (importProps.serverId.length <= 0) {
        ElMessage.warning(t('crontabPage.selectServerTips'))
        return
      }
      importProps.disabled = true
      new CrontabsInRemoteServer({ serverId: Number(importProps.serverId) })
        .request()
        .then((response) => {
          tableData.value = response.data.list.map((command) => {
            const element: Record<string, string> = {}
            const commandSplit = command.split(' ')
            element.command = command
            element.date = commandSplit.slice(0, 5).join(' ')
            element.dateLocale = cronstrue.toString(element.date, {
              locale: getLocale(),
            })
            element.script = commandSplit.slice(5).join(' ')
            element.description = `${element.dateLocale}, ${t('run')}: ${
              element.script
            }`
            return element
          })
        })
        .finally(() => {
          importProps.disabled = false
        })
    }

    let selectedItems = ref<Record<string, string>[]>([])
    const onSelectionChange = (items: Record<string, string>[]) => {
      selectedItems.value = items
    }

    const importCrontab = () => {
      if (selectedItems.value.length === 0) {
        ElMessage.warning(t('crontabPage.selectItemTips'))
        return
      }
      new CrontabImport({
        serverId: Number(importProps.serverId),
        commands: selectedItems.value.map((element) => element.command),
      })
        .request()
        .then(() => {
          ElMessage.success('Success')
          props.onSuccess()
        })
        .finally(() => {
          dialogVisible.value = false
        })
    }

    return {
      dialogVisible,
      serverOption,
      getRemoteServerList,
      importProps,
      tableData,
      onSelectionChange,
      importCrontab,
    }
  },
})
</script>
