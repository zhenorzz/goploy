<template>
  <el-drawer
    ref="drawer"
    v-model="drawerVisible"
    :title="`${serverRow.name}(${serverRow.description})`"
    @opened="connectTerminal"
    @closed="closeTerminal"
  >
    <div v-if="drawerVisible" ref="xterm" class="xterm" />
  </el-drawer>
</template>

<script lang="ts">
import 'xterm/css/xterm.css'
import { Terminal } from 'xterm'
import { FitAddon } from 'xterm-addon-fit'
import { AttachAddon } from 'xterm-addon-attach'
import { computed, defineComponent, ref } from 'vue'

export default defineComponent({
  props: {
    modelValue: {
      type: Boolean,
      default: false,
    },
    serverRow: {
      type: Object,
      required: true,
    },
  },
  emits: ['update:modelValue'],
  setup(props, { emit }) {
    const drawerVisible = computed({
      get: () => props.modelValue,
      set: (val) => {
        emit('update:modelValue', val)
      },
    })
    const terminal = ref<Terminal | null>()
    const websocket = ref<WebSocket | null>()
    // html dom
    const xterm = ref()
    const connectTerminal = () => {
      const isWindows =
        ['Windows', 'Win16', 'Win32', 'WinCE'].indexOf(navigator.platform) >= 0
      const term = new Terminal({
        fontSize: 14,
        cursorBlink: true,
        windowsMode: isWindows,
      })
      const fitAddon = new FitAddon()
      term.loadAddon(fitAddon)
      term.open(xterm.value)
      fitAddon.fit()
      term.focus()
      websocket.value = new WebSocket(
        `${location.protocol.replace('http', 'ws')}//${
          window.location.host + import.meta.env.VITE_APP_BASE_API
        }/ws/xterm?serverId=${props.serverRow.id}&rows=${term.rows}&cols=${
          term.cols
        }`
      )
      const attachAddon = new AttachAddon(websocket.value)
      term.loadAddon(attachAddon)
      terminal.value = term
      websocket.value.onerror = () => {
        websocket.value = null
      }
    }

    const closeTerminal = () => {
      terminal.value = null
      websocket.value?.close()
    }
    return {
      xterm,
      drawerVisible,
      connectTerminal,
      closeTerminal,
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
