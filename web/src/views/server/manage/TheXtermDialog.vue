<template>
  <el-dialog
    v-model="dialogVisible"
    title="WebSSH"
    :close-on-click-modal="false"
    :fullscreen="$store.state.app.device === 'mobile'"
  >
    <el-select
      v-model="serverId"
      v-loading="serverLoading"
      filterable
      default-first-option
      placeholder="please select server"
      style="width: 100%"
      @change="handleSelectServer"
    >
      <el-option
        v-for="item in serverOption"
        :key="item.id"
        :label="`${item.name}(${item.description})`"
        :value="item.id"
      />
    </el-select>
    <el-tabs v-model="activeName" closable @tab-remove="removeTab">
      <el-tab-pane
        v-for="item in tabList"
        :key="item.name"
        :label="item.label"
        :name="item.name"
      >
        <div
          :ref="
            (el) => {
              item.el = el
            }
          "
          class="xterm"
          :style="
            $store.state.app.device === 'mobile'
              ? 'height: calc(100vh - 202px)'
              : 'height:500px'
          "
        />
      </el-tab-pane>
    </el-tabs>
    <el-input
      v-show="tabList.length > 0"
      v-model="command"
      type="textarea"
      placeholder="Send to all windows"
      @keyup.enter="enterCommand"
    />
  </el-dialog>
</template>

<script lang="ts">
import 'xterm/css/xterm.css'
import { xterm } from './xterm'
import { computed, watch, defineComponent, ref, nextTick } from 'vue'
import { ServerOption } from '@/api/server'

export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
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

    const serverLoading = ref(false)
    const serverOption = ref<ServerOption['datagram']['list']>([])
    watch(dialogVisible, (val: boolean) => {
      if (val === true) {
        serverLoading.value = true
        new ServerOption()
          .request()
          .then((response) => {
            serverOption.value = response.data.list
          })
          .finally(() => {
            serverLoading.value = false
          })
      }
    })
    const activeName = ref('')
    const tabList = ref<{ label: string; name: string; el: HTMLDivElement }[]>(
      []
    )
    const serverId = ref('')
    let xterms: xterm[] = []
    const handleSelectServer = () => {
      const selectedServer = serverOption.value.find(
        (_) => _.id === Number(serverId.value)
      )
      if (!selectedServer) {
        return
      }
      const tabLen = tabList.value.length
      const item = {
        label: selectedServer.ip,
        name: '' + tabLen,
        el: {} as HTMLDivElement,
      }
      tabList.value.push(item)
      activeName.value = item.name
      serverId.value = ''
      nextTick(() => {
        xterms[tabLen] = new xterm(item.el, selectedServer.id)
        xterms[tabLen].connect()
      })
    }

    const removeTab = (targetName: string) => {
      if (activeName.value === targetName) {
        tabList.value.forEach((tab, index) => {
          if (tab.name === targetName) {
            let nextTab = tabList.value[index + 1] || tabList.value[index - 1]
            if (nextTab) {
              activeName.value = nextTab.name
            }
          }
        })
      }
      xterms[Number(targetName)].close()
      tabList.value = tabList.value.filter((tab) => tab.name !== targetName)
    }

    const command = ref('')
    const enterCommand = () => {
      xterms.forEach((item: xterm) => {
        item.send(command.value)
      })
      command.value = ''
    }

    return {
      dialogVisible,
      serverLoading,
      serverOption,
      serverId,
      tabList,
      activeName,
      handleSelectServer,
      removeTab,
      command,
      enterCommand,
    }
  },
})
</script>
<style lang="scss" scoped>
.xterm {
  width: 100%;
  height: 100%;
}
</style>
