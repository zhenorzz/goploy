<template>
  <el-row class="app-container">
    <el-row class="header">
      <el-scrollbar>
        <el-row class="nav" justify="start" align="middle">
          <div
            v-for="(item, index) in terminalList"
            :key="item.uuid"
            class="nav-item"
            :class="
              item.uuid === currentTerminalUUID ? 'nav-item-selected' : ''
            "
          >
            <el-row>
              <div class="nav-item-serial" @click="selectTerminal(item)">
                {{ index + 1 }}
              </div>
              <div
                class="nav-item-name"
                :title="`${item.server.label}`"
                @click="selectTerminal(item)"
              >
                {{ item.server.label }}
              </div>
              <el-button
                link
                style="color: #bfcbd9; padding-left: 8px; height: 42px"
                @click="deleteTerminal(item, index)"
              >
                x
              </el-button>
            </el-row>
          </div>
          <div class="nav-plus">
            <el-select
              v-model="serverId"
              style="width: 200px"
              filterable
              clearable
              @change="selectServer"
            >
              <el-option
                v-for="item in serverOption"
                :key="item.id"
                :label="item.label"
                :value="item.id"
              />
            </el-select>
          </div>
        </el-row>
      </el-scrollbar>
    </el-row>
    <el-row class="terminal">
      <div
        v-for="item in terminalList"
        v-show="item.uuid === currentTerminalUUID"
        :key="item.uuid"
        :ref="
          (el) => {
            if (el) {
              terminalRefs[item.uuid] = el
            }
          }
        "
        style="width: 100%; height: 100%"
      ></div>
    </el-row>
    <el-row class="footer">
      <el-input
        v-model="command"
        :disabled="terminalList.length === 0"
        placeholder="Click here send to all windows"
        class="terminal-cmd"
        @keyup.enter="enterCommand"
      />
    </el-row>
  </el-row>
</template>
<script lang="ts">
export default { name: 'ServerTerminal' }
</script>
<script lang="ts" setup>
import 'xterm/css/xterm.css'
import { ServerOption, ServerData } from '@/api/server'
import { xterm } from './xterm'
import { ref, nextTick, ComponentPublicInstance, onBeforeUnmount } from 'vue'
interface terminal {
  uuid: number
  xterm?: xterm
  server: ServerData
}
const terminalList = ref<terminal[]>([])
const currentTerminalUUID = ref(0)
const serverOption = ref<ServerOption['datagram']['list']>([])
const serverId = ref('')
const command = ref('')
const terminalRefs = ref<Record<string, Element | ComponentPublicInstance>>({})

getServerOption()

onBeforeUnmount(() => {
  for (const terminal of terminalList.value) {
    terminal.xterm?.close()
  }
})

function getServerOption() {
  new ServerOption().request().then((response) => {
    serverOption.value = response.data.list
  })
}

function selectServer(value: number) {
  const server =
    serverOption.value.find((_) => _.id === value) || ({} as ServerData)
  if (terminalList.value.length === 0) {
    currentTerminalUUID.value = 0
  } else {
    currentTerminalUUID.value =
      terminalList.value[terminalList.value.length - 1].uuid + 1
  }
  terminalList.value.push({ uuid: currentTerminalUUID.value, server })
  serverId.value = ''
  nextTick(() => {
    const x = new xterm(
      terminalRefs.value[currentTerminalUUID.value] as HTMLDivElement,
      server.id
    )
    x.connect()
    terminalList.value[terminalList.value.length - 1].xterm = x
  })
}
function selectTerminal(terminal: terminal) {
  currentTerminalUUID.value = terminal.uuid
}
function deleteTerminal(terminal: terminal, index: number) {
  terminal.xterm?.close()
  terminalList.value.splice(index, 1)
  if (currentTerminalUUID.value === terminal.uuid) {
    currentTerminalUUID.value =
      terminalList.value.length === 0
        ? 0
        : terminalList.value[terminalList.value.length - 1].uuid
  }
}
function enterCommand() {
  terminalList.value.forEach((terminal) => {
    terminal.xterm?.send(command.value + '\n')
  })
  command.value = ''
}
</script>

<style lang="scss" scoped>
.header {
  width: 100%;
  background: #0e0f12;
}
.nav {
  height: auto;
  flex-wrap: nowrap;
  &-item {
    color: #bfcbd9;
    font-size: 14px;
    flex-shrink: 0;
    width: 185px;
    height: 45px;
    line-height: 45px;
    padding: 0 10px;
    &-serial {
      cursor: pointer;
      color: var(--el-text-color-secondary);
      text-align: center;
      font-weight: 600;
      width: 20px;
    }
    &-name {
      flex: 1;
      cursor: pointer;
      text-overflow: ellipsis;
      overflow: hidden;
      white-space: nowrap;
      display: inline-block;
    }
    &-selected {
      background: #1d2935;
    }
  }
  &-plus {
    text-align: center;
    // border-bottom: 1px solid var(--el-border-color);
    // border-right: 1px solid var(--el-border-color);
  }
}

.terminal {
  flex: 1;
  height: 100%;
  padding: 10px;
  background: #1d2935;
}

.footer {
  height: 28px;
  width: 100%;
  background: #0e0f12;
}

@media only screen and (max-device-width: 992px) {
  .terminal {
    width: calc(100vw - 40px);
  }
}
</style>
<style lang="scss">
@import '@/styles/mixin.scss';
.xterm .xterm-viewport {
  overflow: auto;
  @include scrollBar();
}
.nav-popover.el-popper {
  background: #0e0f12;
  border-color: #0e0f12;
}
.terminal-cmd {
  .el-input {
    &__wrapper {
      background-color: #0e0f12 !important;
      border-color: #304156 !important;
      box-shadow: none !important;
    }
    &__inner {
      color: #f0f2f5 !important;
    }
  }
}

.nav-plus {
  .el-input__wrapper {
    border-radius: 0px !important;
    box-shadow: none;
    background-color: #0e0f12 !important;
    border-color: #304156 !important;
  }
  .el-input__inner {
    height: 43px !important;
  }
}

.server-list {
  height: 216px;
  margin-top: 10px;
}

.server-filter {
  .el-input {
    &__wrapper {
      background-color: #0e0f12 !important;
    }
    &__inner {
      color: #f0f2f5 !important;
    }
  }
}

.server-item {
  width: 100%;
  text-align: left;
  display: inline-block;
  color: #cfd3dc !important;
  &:hover {
    color: var(--el-button-hover-text-color) !important;
  }
}
</style>
